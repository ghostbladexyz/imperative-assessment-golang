package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
	"github.com/pleft/imperative-assessment-golang/internal/receipts"
	"github.com/pleft/imperative-assessment-golang/internal/runner"
	frontend "github.com/pleft/imperative-assessment-golang/internal/web"
)

const maxRequestBytes = runner.MaxSourceBytes + 32*1024

type api struct {
	runner   *runner.Runner
	receipts *receipts.Manager
}

type runRequest struct {
	LevelID int      `json:"levelId"`
	Code    string   `json:"code"`
	TestIDs []string `json:"testIds"`
}

type formatRequest struct {
	Code string `json:"code"`
}

type validateRequest struct {
	Receipts map[string]string `json:"receipts"`
}

func main() {
	address := flag.String("addr", "127.0.0.1:8080", "local address to listen on")
	openBrowser := flag.Bool("open", false, "open the assessment in the default browser")
	flag.Parse()

	if err := assessment.Validate(); err != nil {
		log.Fatal(err)
	}
	dataDir, err := localDataDir()
	if err != nil {
		log.Fatal(err)
	}
	receiptManager, err := receipts.New(dataDir)
	if err != nil {
		log.Fatalf("initialize pass receipts: %v", err)
	}
	handler, err := routes(&api{
		runner:   runner.New("go", 2, receiptManager),
		receipts: receiptManager,
	})
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              *address,
		Handler:           securityHeaders(handler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      35 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("listen on %s: %v", *address, err)
	}
	url := "http://" + listener.Addr().String()
	log.Printf("Imperative Go Practice Assessment is ready at %s", url)
	log.Printf("Trusted local practice only: the Go runner is not hardened for public deployment.")
	if *openBrowser {
		go func() {
			time.Sleep(250 * time.Millisecond)
			_ = launchBrowser(url)
		}()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()
	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func routes(api *api) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(writer http.ResponseWriter, _ *http.Request) {
		writeJSON(writer, http.StatusOK, map[string]any{"ok": true, "goVersion": strings.TrimSpace(goVersion())})
	})
	mux.HandleFunc("GET /api/levels", func(writer http.ResponseWriter, _ *http.Request) {
		writeJSON(writer, http.StatusOK, map[string]any{"levels": assessment.PublicLevels(), "schemaVersion": 1})
	})
	mux.HandleFunc("POST /api/run", api.run)
	mux.HandleFunc("POST /api/format", api.format)
	mux.HandleFunc("POST /api/receipts/validate", api.validateReceipts)

	dist, err := fs.Sub(frontend.Files, "dist")
	if err != nil {
		return nil, err
	}
	files := http.FileServer(http.FS(dist))
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if strings.HasPrefix(request.URL.Path, "/api/") {
			writeJSON(writer, http.StatusNotFound, map[string]string{"error": "API endpoint not found"})
			return
		}
		path := strings.TrimPrefix(request.URL.Path, "/")
		if path != "" {
			if _, err := fs.Stat(dist, path); err == nil {
				files.ServeHTTP(writer, request)
				return
			}
		}
		request.URL.Path = "/"
		files.ServeHTTP(writer, request)
	})
	return mux, nil
}

func (api *api) run(writer http.ResponseWriter, request *http.Request) {
	var input runRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		return
	}
	level, found := assessment.FindLevel(input.LevelID)
	if !found {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "unknown level"})
		return
	}
	writeJSON(writer, http.StatusOK, api.runner.Run(request.Context(), level, input.Code, input.TestIDs))
}

func (api *api) format(writer http.ResponseWriter, request *http.Request) {
	var input formatRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
	defer cancel()
	formatted, err := api.runner.Format(ctx, input.Code)
	if err != nil {
		writeJSON(writer, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"code": formatted})
}

func (api *api) validateReceipts(writer http.ResponseWriter, request *http.Request) {
	var input validateRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		return
	}
	valid := make([]int, 0, len(input.Receipts))
	for levelText, encoded := range input.Receipts {
		levelID, err := strconv.Atoi(levelText)
		receipt, ok := api.receipts.Validate(encoded)
		if err == nil && ok && receipt.LevelID == levelID {
			valid = append(valid, levelID)
		}
	}
	writeJSON(writer, http.StatusOK, map[string]any{"validLevelIds": valid})
}

func decodeJSON(writer http.ResponseWriter, request *http.Request, destination any) error {
	request.Body = http.MaxBytesReader(writer, request.Body, maxRequestBytes)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "invalid request: " + err.Error()})
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "request must contain exactly one JSON value"})
		return errors.New("trailing JSON")
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "DENY")
		writer.Header().Set("Referrer-Policy", "no-referrer")
		writer.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'; img-src 'self' data:; connect-src 'self'; font-src 'self'")
		next.ServeHTTP(writer, request)
	})
}

func localDataDir() (string, error) {
	if configured := os.Getenv("IMPERATIVE_ASSESSMENT_DATA"); configured != "" {
		return filepath.Abs(configured)
	}
	config, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(config, "imperative-go-assessment"), nil
}

func goVersion() string {
	output, err := exec.Command("go", "version").Output()
	if err != nil {
		return "Go executable unavailable"
	}
	return string(output)
}

func launchBrowser(url string) error {
	var command *exec.Cmd
	switch {
	case runtime.GOOS == "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case runtime.GOOS == "darwin":
		command = exec.Command("open", url)
	default:
		command = exec.Command("xdg-open", url)
	}
	return command.Start()
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("assessment: ")
}
