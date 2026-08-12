package assessment

// errorInspectorExercise defines classification across wrapped, joined, sentinel, and typed errors.
func errorInspectorExercise() advancedSpec {
	tests := advancedTests(10,
		[4]string{"Nil", "Classifies nil without fabricating a failure.", `nil`, `{"kind":"none","field":"","temporary":false}`},
		[4]string{"Wrapped unavailable", "Finds a sentinel through wrapping.", `fmt.Errorf("fetch: %w", ErrUnavailable)`, `{"kind":"unavailable","field":"","temporary":true}`},
		[4]string{"Typed field error", "Extracts typed error data through wrapping.", `wrapped FieldError`, `{"kind":"field","field":"email","temporary":false}`},
		[4]string{"Joined errors", "Inspects every branch of an errors.Join tree.", `errors.Join(FieldError, ErrUnavailable)`, `{"kind":"field","field":"name","temporary":true}`},
		[4]string{"Authorization precedence", "Uses authorization as the highest-priority kind.", `errors.Join(ErrUnavailable, ErrUnauthorized)`, `{"kind":"unauthorized","field":"","temporary":true}`},
		[4]string{"Unknown", "Keeps unclassified errors explicit.", `errors.New("other")`, `{"kind":"unknown","field":"","temporary":false}`},
	)
	return advancedSpec{
		id: 10, title: "Error Inspector", topic: "Error handling · wrapping · errors.Is · errors.As",
		signature: "InspectError(err error) Failure",
		objective: "Classify complex Go error trees without comparing error strings.",
		input:     "A nil, wrapped, joined, sentinel, or typed error.",
		output:    "A deterministic Failure classification with authorization precedence and independent temporary status.",
		constraints: []string{
			"Kind precedence is unauthorized, field, unavailable, unknown; nil is none.",
			"Temporary is true whenever ErrUnavailable occurs anywhere in the tree, even when another Kind wins.",
			"Use errors.Is and errors.As so wrapping and errors.Join remain transparent.",
		},
		hints:    []string{"Compute Temporary independently from Kind.", "Check the highest-priority sentinel first.", "Use a *FieldError target with errors.As to retrieve Field."},
		pitfalls: []string{"Comparing err.Error()", "Inspecting only errors.Unwrap's single chain", "Losing temporary status when another classification wins"},
		packages: []string{"errors"},
		starter: `import "errors"

var (
	ErrUnavailable  = errors.New("unavailable")
	ErrUnauthorized = errors.New("unauthorized")
)

type FieldError struct { Field string }
func (err *FieldError) Error() string { return "invalid field: " + err.Field }

type Failure struct {
	Kind      string ` + "`json:\"kind\"`" + `
	Field     string ` + "`json:\"field\"`" + `
	Temporary bool   ` + "`json:\"temporary\"`" + `
}

func InspectError(err error) Failure {
	// TODO: inspect the complete error tree without matching text.
	return Failure{}
}
`,
		tests: tests,
		build: func(selected []VisibleTest) string {
			loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	cases := []error{nil, fmt.Errorf("fetch: %w", ErrUnavailable), fmt.Errorf("input: %w", &FieldError{Field:"email"}), errors.Join(&FieldError{Field:"name"}, ErrUnavailable), errors.Join(ErrUnavailable, ErrUnauthorized), errors.New("other")}
	for _, current := range tests { current := current; assessmentRun(current.ID, func() string { encoded,_:=json.Marshal(InspectError(cases[current.Payload.Case])); return string(encoded) }) }`
			return harness(commonImports("errors"), advancedCaseDeclarations, loop, selected)
		},
	}
}

// compensatingStepsExercise defines reverse-order cleanup and joined error preservation.
func compensatingStepsExercise() advancedSpec {
	tests := advancedTests(11,
		[4]string{"Success", "Runs every step and performs no compensation.", `ctx, two successful steps`, `run-1,run-2|ok`},
		[4]string{"Run failure", "Compensates completed steps in reverse order.", `ctx, second step fails`, `run-1,run-2,undo-1|run`},
		[4]string{"Cleanup failure", "Preserves both run and compensation errors.", `ctx, run and undo fail`, `run-1,run-2,undo-1|run+undo`},
		[4]string{"Reverse cleanup", "Undoes multiple completed steps in strict reverse order.", `ctx, fourth step fails`, `run-1,run-2,run-3,run-4,undo-3,undo-2,undo-1|run`},
		[4]string{"Cancelled before start", "Checks the parent context before invoking a step.", `cancelled ctx`, `|context canceled`},
	)
	return advancedSpec{
		id: 11, title: "Compensating Steps", topic: "Error handling · cleanup · errors.Join · context",
		signature: "ExecuteSteps(ctx context.Context, steps []Step) error",
		objective: "Run fallible steps and compensate completed work when execution stops.",
		input:     "A context and ordered Run/Undo function pairs.",
		output:    "nil after full success, or an error tree containing the triggering error and every reverse-order Undo failure.",
		constraints: []string{
			"Check ctx before each Run; a cancelled context is the triggering error.",
			"After failure, call Undo only for successfully completed steps and in reverse order.",
			"Attempt every required Undo even after one fails, and combine errors with errors.Join.",
		},
		hints:    []string{"Track how many Run calls completed successfully.", "Walk completed steps backward after the first failure.", "Collect the trigger first, append Undo errors, then errors.Join the collection."},
		pitfalls: []string{"Undoing the step whose Run failed", "Stopping cleanup after its first error", "Returning an Undo error while losing the original failure"},
		builtins: []string{"append", "len"},
		packages: []string{"context", "errors"},
		starter: `import (
	"context"
	"errors"
)

type Step struct {
	Run  func(context.Context) error
	Undo func(context.Context) error
}

func ExecuteSteps(ctx context.Context, steps []Step) error {
	// TODO: run forward, compensate backward, and preserve every failure.
	_ = errors.Join
	return nil
}
`,
		tests: tests,
		build: buildCompensationHarness,
	}
}

// buildCompensationHarness records call order because correct cleanup sequencing is part of the interface.
func buildCompensationHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `
var assessmentRunFailure = errors.New("run failed")
var assessmentUndoFailure = errors.New("undo failed")`
	loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests { current:=current; assessmentRun(current.ID, func() string {
		trace:=[]string{}; makeStep:=func(id int, runErr, undoErr error) Step { return Step{Run:func(context.Context)error{trace=append(trace,fmt.Sprintf("run-%d",id));return runErr},Undo:func(context.Context)error{trace=append(trace,fmt.Sprintf("undo-%d",id));return undoErr}} }
		ctx:=context.Background(); steps:=[]Step{}
		switch current.Payload.Case { case 0: steps=[]Step{makeStep(1,nil,nil),makeStep(2,nil,nil)}; case 1: steps=[]Step{makeStep(1,nil,nil),makeStep(2,assessmentRunFailure,nil)}; case 2: steps=[]Step{makeStep(1,nil,assessmentUndoFailure),makeStep(2,assessmentRunFailure,nil)}; case 3: steps=[]Step{makeStep(1,nil,nil),makeStep(2,nil,nil),makeStep(3,nil,nil),makeStep(4,assessmentRunFailure,nil)}; default: var cancel context.CancelFunc; ctx,cancel=context.WithCancel(ctx);cancel();steps=[]Step{makeStep(1,nil,nil)} }
		err:=ExecuteSteps(ctx,steps); label:="ok"; if errors.Is(err,assessmentRunFailure)&&errors.Is(err,assessmentUndoFailure){label="run+undo"}else if errors.Is(err,assessmentRunFailure){label="run"}else if errors.Is(err,context.Canceled){label="context canceled"}else if err!=nil{label="wrong"}; return strings.Join(trace,",")+"|"+label
	}) }`
	return harness(commonImports("context", "errors", "strings"), declarations, loop, selected)
}

// profileClientExercise defines a bounded and strict HTTP client decode contract without real network access.
func profileClientExercise() advancedSpec {
	tests := advancedTests(12,
		[4]string{"Success", "Decodes a valid bounded response.", `ctx, client, url`, `{"id":"42","name":"Ada"}|<nil>`},
		[4]string{"Remote status", "Classifies non-2xx responses without decoding them.", `503 response`, `{"id":"","name":""}|remote response`},
		[4]string{"Unknown JSON field", "Rejects response schema drift.", `JSON with extra field`, `{"id":"","name":""}|invalid response`},
		[4]string{"Trailing JSON", "Rejects multiple JSON values.", `two JSON objects`, `{"id":"","name":""}|invalid response`},
		[4]string{"Oversized body", "Enforces the 64 KiB response limit.", `large response`, `{"id":"","name":""}|invalid response`},
		[4]string{"Transport error", "Wraps and preserves the client's transport error.", `failing transport`, `{"id":"","name":""}|transport failure`},
	)
	return advancedSpec{
		id: 12, title: "Profile Client", topic: "HTTP · clients · bounded JSON · transport errors",
		signature: "FetchProfile(ctx context.Context, client *http.Client, url string) (Profile, error)",
		objective: "Fetch and strictly decode a small JSON profile through a supplied HTTP client.",
		input:     "A context, non-nil http.Client, and URL.",
		output:    "A valid Profile for a 2xx response, or a classified/wrapped error and zero Profile.",
		constraints: []string{
			"Create a GET request with NewRequestWithContext and execute it through client.Do.",
			"Map non-2xx status to ErrRemoteResponse and transport failures must remain discoverable with errors.Is.",
			"Read no more than 64 KiB plus one detection byte; strictly decode one JSON object with non-empty id and name.",
		},
		hints:    []string{"Use io.LimitReader(response.Body, 64*1024+1).", "DisallowUnknownFields and require the second Decode to return io.EOF.", "Always close a non-nil response body."},
		pitfalls: []string{"Using http.Get instead of the supplied client", "Reading an unbounded body", "Accepting a valid object followed by junk"},
		builtins: []string{"len", "make"},
		packages: []string{"bytes", "context", "encoding/json", "errors", "fmt", "io", "net/http"},
		starter: `import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrRemoteResponse  = errors.New("remote response")
	ErrInvalidResponse = errors.New("invalid response")
)

type Profile struct { ID string ` + "`json:\"id\"`" + `; Name string ` + "`json:\"name\"`" + ` }

func FetchProfile(ctx context.Context, client *http.Client, url string) (Profile, error) {
	// TODO: issue a contextual request and strictly decode a bounded response.
	_ = bytes.NewReader
	_ = json.NewDecoder
	_ = fmt.Errorf
	_ = io.LimitReader
	return Profile{}, ErrInvalidResponse
}
`,
		tests: tests,
		build: buildProfileHarness,
	}
}

// buildProfileHarness supplies RoundTrippers so client behavior is deterministic and network-free.
func buildProfileHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `
var assessmentTransportFailure = errors.New("transport failure")
type assessmentRoundTrip func(*http.Request)(*http.Response,error)
func (roundTrip assessmentRoundTrip) RoundTrip(request *http.Request)(*http.Response,error){return roundTrip(request)}`
	loop := `var tests []wireCase; _=json.Unmarshal([]byte(raw),&tests)
	for _,current:=range tests{current:=current;assessmentRun(current.ID,func()string{
		status:=200;body:="{\"id\":\"42\",\"name\":\"Ada\"}";transportErr:=error(nil)
		switch current.Payload.Case{case 1:status=503;case 2:body="{\"id\":\"42\",\"name\":\"Ada\",\"extra\":true}";case 3:body+=" {}";case 4:body="{\"id\":\"42\",\"name\":\""+strings.Repeat("x",70*1024)+"\"}";case 5:transportErr=assessmentTransportFailure}
		client:=&http.Client{Transport:assessmentRoundTrip(func(request *http.Request)(*http.Response,error){if request.Method!="GET"||request.Context()==nil{return nil,errors.New("bad request")};if transportErr!=nil{return nil,transportErr};return &http.Response{StatusCode:status,Body:io.NopCloser(strings.NewReader(body)),Header:make(http.Header)},nil})}
		value,err:=FetchProfile(context.Background(),client,"https://profiles.test/42");encoded,_:=json.Marshal(value);label:=fmt.Sprint(err);if errors.Is(err,assessmentTransportFailure){label="transport failure"};return string(encoded)+"|"+label
	})}`
	return harness(commonImports("context", "errors", "net/http", "strings"), declarations, loop, selected)
}

// conditionalDocumentExercise defines conditional GET and HEAD behavior around an injected loader.
func conditionalDocumentExercise() advancedSpec {
	tests := advancedTests(13,
		[4]string{"GET", "Returns a body and ETag for a document.", `GET /docs/a`, `200|"v3"|alpha`},
		[4]string{"HEAD", "Returns GET metadata without a body.", `HEAD /docs/a`, `200|"v3"|`},
		[4]string{"Not modified", "Honors an exact If-None-Match value.", `GET with If-None-Match`, `304|"v3"|`},
		[4]string{"Missing", "Maps ErrDocumentNotFound to 404.", `GET /docs/missing`, `404||`},
		[4]string{"Loader failure", "Maps unknown loader failures to 500.", `GET /docs/fail`, `500||`},
		[4]string{"Routing", "Separates unsupported methods from unknown paths.", `POST /docs/a; GET /other`, `405,404`},
	)
	return advancedSpec{
		id: 13, title: "Conditional Document", topic: "HTTP · ETag · conditional GET · HEAD",
		signature: "NewDocumentHandler(load func(context.Context, string) (Document, error)) http.Handler",
		objective: "Serve cache-aware documents through an injected loader.",
		input:     "GET or HEAD /docs/{id}, optionally with If-None-Match.",
		output:    "Correct status, quoted version ETag, and body semantics for GET, HEAD, and 304.",
		constraints: []string{
			"Reject unrelated/empty/nested paths with 404 and methods other than GET/HEAD with 405.",
			"Map ErrDocumentNotFound to 404 and other loader errors to 500.",
			"Set ETag to the quoted form `\"vN\"`; exact If-None-Match returns 304, and HEAD never writes Body.",
		},
		hints:    []string{"Validate the route before loading.", "Build the ETag with strconv.Quote and strconv.Itoa.", "Set ETag before checking If-None-Match."},
		pitfalls: []string{"Writing a body for HEAD", "Returning 200 for a matching ETag", "Treating /docs/a/more as document a"},
		packages: []string{"context", "errors", "net/http", "strconv", "strings"},
		starter: `import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var ErrDocumentNotFound = errors.New("document not found")
type Document struct { ID string; Body string; Version int }

func NewDocumentHandler(load func(context.Context, string) (Document, error)) http.Handler {
	// TODO: implement route validation, status mapping, ETag, GET, and HEAD.
	_ = strconv.Quote
	_ = strings.TrimPrefix
	return http.NotFoundHandler()
}
`,
		tests: tests,
		build: buildDocumentHarness,
	}
}

// buildDocumentHarness keeps conditional-request assertions at the HTTP seam.
func buildDocumentHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `
var assessmentDocumentFailure=errors.New("load failed")
func assessmentDocumentRequest(handler http.Handler,method,path,etag string)(int,string,string){request:=httptest.NewRequest(method,path,nil);if etag!=""{request.Header.Set("If-None-Match",etag)};recorder:=httptest.NewRecorder();handler.ServeHTTP(recorder,request);return recorder.Code,recorder.Header().Get("ETag"),recorder.Body.String()}`
	loop := `var tests []wireCase;_=json.Unmarshal([]byte(raw),&tests);for _,current:=range tests{current:=current;assessmentRun(current.ID,func()string{handler:=NewDocumentHandler(func(ctx context.Context,id string)(Document,error){if id=="missing"{return Document{},ErrDocumentNotFound};if id=="fail"{return Document{},assessmentDocumentFailure};return Document{ID:id,Body:"alpha",Version:3},nil});switch current.Payload.Case{case 0:s,e,b:=assessmentDocumentRequest(handler,"GET","/docs/a","");return fmt.Sprintf("%d|%s|%s",s,e,b);case 1:s,e,b:=assessmentDocumentRequest(handler,"HEAD","/docs/a","");return fmt.Sprintf("%d|%s|%s",s,e,b);case 2:s,e,b:=assessmentDocumentRequest(handler,"GET","/docs/a","\"v3\"");return fmt.Sprintf("%d|%s|%s",s,e,b);case 3:s,_,_:=assessmentDocumentRequest(handler,"GET","/docs/missing","");return fmt.Sprintf("%d||",s);case 4:s,_,_:=assessmentDocumentRequest(handler,"GET","/docs/fail","");return fmt.Sprintf("%d||",s);default:first,_,_:=assessmentDocumentRequest(handler,"POST","/docs/a","");second,_,_:=assessmentDocumentRequest(handler,"GET","/other","");return fmt.Sprintf("%d,%d",first,second)}})}`
	return harness(commonImports("context", "errors", "net/http", "net/http/httptest"), declarations, loop, selected)
}

// weightedMazeExercise defines Dijkstra-style path finding with strict grid validation.
func weightedMazeExercise() advancedSpec {
	tests := advancedTests(14,
		[4]string{"Single cell", "Includes the start cell cost.", `[][]int{{7}}`, `7|<nil>`},
		[4]string{"Cheapest detour", "Chooses cost rather than the fewest steps.", `weighted grid`, `7|<nil>`},
		[4]string{"Obstacles", "Routes around -1 blocked cells.", `blocked grid`, `9|<nil>`},
		[4]string{"No path", "Classifies an unreachable destination.", `disconnected grid`, `0|no path`},
		[4]string{"Ragged grid", "Rejects inconsistent row lengths.", `ragged grid`, `0|invalid grid`},
		[4]string{"Invalid weight", "Rejects negative weights other than -1.", `grid containing -2`, `0|invalid grid`},
	)
	return advancedSpec{
		id: 14, title: "Weighted Maze", topic: "Algorithms · Dijkstra · grids · validation",
		signature:   "CheapestPath(grid [][]int) (int, error)",
		objective:   "Find the minimum traversal cost from the top-left to bottom-right grid cell.",
		input:       "A non-empty rectangular grid; -1 is blocked and non-negative values are cell-entry costs.",
		output:      "Minimum cost including start/end, ErrNoPath, or ErrInvalidGrid.",
		constraints: []string{"Move only up, down, left, or right.", "Validate the entire rectangular grid before searching; values below -1 are invalid.", "A blocked start/end or disconnected destination returns ErrNoPath."},
		hints:       []string{"Store the best known cost per coordinate.", "A min-heap should prioritize the smallest accumulated cost.", "Skip a popped state when a cheaper cost is already recorded."},
		pitfalls:    []string{"Using breadth-first search on weighted cells", "Omitting the start cost", "Treating every negative value as a valid obstacle"},
		builtins:    []string{"append", "len", "make"},
		packages:    []string{"container/heap", "errors"},
		starter: `import (
	"container/heap"
	"errors"
)

var ( ErrInvalidGrid = errors.New("invalid grid"); ErrNoPath = errors.New("no path") )

func CheapestPath(grid [][]int) (int, error) {
	// TODO: validate the grid and find the minimum weighted path.
	_ = heap.Init
	return 0, ErrInvalidGrid
}
`,
		tests: tests,
		build: func(selected []VisibleTest) string {
			loop := `var tests []wireCase;_=json.Unmarshal([]byte(raw),&tests);cases:=[][][]int{{{7}},{{1,50,1},{2,2,1},{50,1,1}},{{1,-1,5},{2,-1,2},{2,2,2}},{{1,-1},{-1,1}},{{1,2},{3}},{{1,-2}}};for _,current:=range tests{current:=current;assessmentRun(current.ID,func()string{value,err:=CheapestPath(cases[current.Payload.Case]);return fmt.Sprintf("%d|%v",value,err)})}`
			return harness(commonImports(), advancedCaseDeclarations, loop, selected)
		},
	}
}

// channelFanInExercise defines lifecycle-safe merging for dynamic channel sets.
func channelFanInExercise() advancedSpec {
	tests := advancedTests(15,
		[4]string{"Merge", "Forwards every value from every input.", `ctx, first, second`, `[1 2 3 4]|closed`},
		[4]string{"Uneven producers", "Does not serialize fast and slow producers.", `ctx, fast, slow`, `[1 2 9]|closed`},
		[4]string{"Nil input", "Ignores nil channels rather than waiting forever.", `ctx, nil, values`, `[5 6]|closed`},
		[4]string{"No inputs", "Returns an already-closed channel.", `ctx`, `[]|closed`},
		[4]string{"Cancellation", "Stops forwarders and closes output after cancellation.", `cancelled ctx, blocked input`, `[]|closed`},
	)
	return advancedSpec{
		id: 15, title: "Channel Fan-In", topic: "Concurrency · channels · cancellation · lifecycle",
		signature:   "FanIn(ctx context.Context, inputs ...<-chan int) <-chan int",
		objective:   "Merge any number of input channels into one cancellation-aware output channel.",
		input:       "A context and zero or more receive-only channels, including possible nil channels.",
		output:      "Every value produced before cancellation in arbitrary order, followed by exactly one output close.",
		constraints: []string{"Use independent forwarding so one idle producer cannot block the others.", "Ignore nil inputs; zero usable inputs returns a closed channel.", "On cancellation, every forwarder must stop and output closes only after all forwarders exit."},
		hints:       []string{"Start one goroutine per non-nil input.", "Each forwarder selects on both receiving/sending and ctx.Done.", "A final goroutine can wait for the WaitGroup and close output once."},
		pitfalls:    []string{"Reading inputs sequentially", "Closing output from a worker", "Starting a goroutine that blocks forever on a nil input"},
		builtins:    []string{"close", "make"},
		packages:    []string{"context", "sync"},
		starter: `import (
	"context"
	"sync"
)

func FanIn(ctx context.Context, inputs ...<-chan int) <-chan int {
	// TODO: forward concurrently and close output after every worker exits.
	_ = sync.WaitGroup{}
	output := make(chan int)
	close(output)
	return output
}
`,
		tests: tests,
		build: buildFanInHarness,
	}
}

// buildFanInHarness sorts collected values because the FanIn interface intentionally leaves ordering unspecified.
func buildFanInHarness(selected []VisibleTest) string {
	loop := `var tests []wireCase;_=json.Unmarshal([]byte(raw),&tests);for _,current:=range tests{current:=current;assessmentRun(current.ID,func()string{ctx,cancel:=context.WithCancel(context.Background());defer cancel();makeInput:=func(values ...int)<-chan int{channel:=make(chan int,len(values));for _,value:=range values{channel<-value};close(channel);return channel};var output <-chan int;switch current.Payload.Case{case 0:output=FanIn(ctx,makeInput(1,3),makeInput(2,4));case 1:slow:=make(chan int);go func(){defer close(slow);time.Sleep(10*time.Millisecond);slow<-9}();output=FanIn(ctx,makeInput(1,2),slow);case 2:output=FanIn(ctx,nil,makeInput(5,6));case 3:output=FanIn(ctx);default:blocked:=make(chan int);output=FanIn(ctx,blocked);cancel()};values:=[]int{};timeout:=time.After(300*time.Millisecond);for{select{case value,open:=<-output:if !open{sort.Ints(values);return fmt.Sprintf("%v|closed",values)};values=append(values,value);case <-timeout:return "timeout"}}})}`
	return harness(commonImports("context", "sort"), advancedCaseDeclarations, loop, selected)
}

// firstSuccessExercise defines racing work with loser cancellation and stable aggregate failures.
func firstSuccessExercise() advancedSpec {
	tests := advancedTests(16,
		[4]string{"Fastest success", "Returns the first successful completion rather than task order.", `ctx, slow, fast`, `fast|ok`},
		[4]string{"Ignore early error", "Continues after failures while another task can succeed.", `ctx, failing, successful`, `value|ok`},
		[4]string{"All fail", "Returns an error tree containing every task failure.", `ctx, three failures`, `|all-errors`},
		[4]string{"Concurrent start", "Starts tasks concurrently.", `ctx, two barrier tasks`, `winner|ok`},
		[4]string{"Cancel losers", "Cancels and waits for losing tasks before returning.", `ctx, winner, observer`, `winner|loser-cancelled`},
		[4]string{"No tasks", "Rejects an empty task set.", `ctx`, `|no tasks`},
	)
	return advancedSpec{
		id: 16, title: "First Success", topic: "Concurrency · racing · cancellation · aggregate errors",
		signature:   "FirstSuccess(ctx context.Context, tasks ...func(context.Context) (string, error)) (string, error)",
		objective:   "Race fallible tasks and return the first success while cleaning up every loser.",
		input:       "A context and one or more context-aware tasks.",
		output:      "The first successful completion, or errors.Join of every failure in task order.",
		constraints: []string{"Start all tasks concurrently.", "On success, cancel the derived context and wait for every task before returning.", "When all tasks fail, join their errors in original task order; no tasks returns ErrNoTasks."},
		hints:       []string{"Give each result its task index.", "A buffered result channel lets every goroutine report even after a winner exists.", "Collect all reports so cleanup is complete, while remembering the first success observed."},
		pitfalls:    []string{"Returning before losers observe cancellation", "Treating the first error as final", "Joining errors in completion order"},
		builtins:    []string{"append", "len", "make"},
		packages:    []string{"context", "errors", "sync"},
		starter: `import (
	"context"
	"errors"
	"sync"
)

var ErrNoTasks = errors.New("no tasks")

func FirstSuccess(ctx context.Context, tasks ...func(context.Context) (string, error)) (string, error) {
	// TODO: race every task, cancel losers, wait, and preserve ordered failures.
	_ = sync.WaitGroup{}
	return "", ErrNoTasks
}
`,
		tests: tests,
		build: buildFirstSuccessHarness,
	}
}

// buildFirstSuccessHarness uses barriers and cancellation observers to make goroutine lifecycle visible.
func buildFirstSuccessHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `
var assessmentFirstErrors=[]error{errors.New("first"),errors.New("second"),errors.New("third")}`
	loop := `var tests []wireCase;_=json.Unmarshal([]byte(raw),&tests);for _,current:=range tests{current:=current;assessmentRun(current.ID,func()string{ctx:=context.Background();switch current.Payload.Case{case 0:value,err:=FirstSuccess(ctx,func(context.Context)(string,error){time.Sleep(20*time.Millisecond);return "slow",nil},func(context.Context)(string,error){return "fast",nil});return value+"|"+map[bool]string{true:"ok",false:fmt.Sprint(err)}[err==nil];case 1:value,err:=FirstSuccess(ctx,func(context.Context)(string,error){return "",assessmentFirstErrors[0]},func(context.Context)(string,error){time.Sleep(time.Millisecond);return "value",nil});if err!=nil{return "|wrong"};return value+"|ok";case 2:_,err:=FirstSuccess(ctx,func(context.Context)(string,error){return "",assessmentFirstErrors[0]},func(context.Context)(string,error){return "",assessmentFirstErrors[1]},func(context.Context)(string,error){return "",assessmentFirstErrors[2]});if errors.Is(err,assessmentFirstErrors[0])&&errors.Is(err,assessmentFirstErrors[1])&&errors.Is(err,assessmentFirstErrors[2]){return "|all-errors"};return "|missing";case 3:ready:=make(chan struct{},2);release:=make(chan struct{});go func(){<-ready;<-ready;close(release)}();task:=func(ctx context.Context)(string,error){ready<-struct{}{};select{case<-release:return "winner",nil;case<-ctx.Done():return "",ctx.Err()}};value,err:=FirstSuccess(ctx,task,task);if err!=nil{return "|wrong"};return value+"|ok";case 4:observed:=make(chan struct{});value,err:=FirstSuccess(ctx,func(context.Context)(string,error){time.Sleep(time.Millisecond);return "winner",nil},func(ctx context.Context)(string,error){<-ctx.Done();close(observed);return "",ctx.Err()});if err!=nil{return "|wrong"};select{case<-observed:return value+"|loser-cancelled";default:return value+"|loser-running"};default:value,err:=FirstSuccess(ctx);if errors.Is(err,ErrNoTasks){return value+"|no tasks"};return value+"|wrong"}})}`
	return harness(commonImports("context", "errors"), declarations, loop, selected)
}

// overdueInvoicesExercise defines disciplined row iteration and error handling over database/sql.
func overdueInvoicesExercise() advancedSpec {
	tests := advancedTests(17,
		[4]string{"Rows", "Scans every row and closes the result set.", `ctx, db, cutoff`, `[{"id":1,"customer":"Ada","due":"2026-01-01","cents":900},{"id":2,"customer":"Lin","due":"2026-01-02","cents":1200}]|query,close|ok`},
		[4]string{"Empty", "Returns an empty non-nil slice.", `empty rows`, `[]|query,close|ok`},
		[4]string{"Query failure", "Wraps a query error and has no rows to close.", `query fails`, `[]|query|database failure`},
		[4]string{"Scan failure", "Discards partial rows and closes the result set.", `malformed second row`, `[]|query,close|scan failure`},
		[4]string{"Rows failure", "Checks rows.Err after iteration and discards partial data.", `rows fail`, `[]|query,close|database failure`},
	)
	return advancedSpec{
		id: 17, title: "Overdue Invoices", topic: "SQL · queries · row lifecycle · scan errors",
		signature:   "ListOverdue(ctx context.Context, db *sql.DB, cutoff time.Time) ([]Invoice, error)",
		objective:   "Query and scan overdue invoices with complete database/sql row error handling.",
		input:       "A context, database handle, and cutoff timestamp.",
		output:      "Ordered invoices or an empty non-nil slice with the wrapped query, scan, or rows error.",
		constraints: []string{"Use one parameterized query ordered by due_at then id.", "Scan id, customer, due_at, and cents; format due_at in UTC as YYYY-MM-DD.", "Close rows, check rows.Err, and discard all partial data after any failure."},
		hints:       []string{"Defer rows.Close immediately after QueryContext succeeds.", "Scan due_at into time.Time.", "Check rows.Err only after the Next loop finishes."},
		pitfalls:    []string{"Ignoring rows.Err", "Returning partial invoices after scan failure", "Interpolating cutoff into SQL"},
		builtins:    []string{"append"},
		packages:    []string{"context", "database/sql", "fmt", "time"},
		starter: `import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Invoice struct { ID int64 ` + "`json:\"id\"`" + `; Customer string ` + "`json:\"customer\"`" + `; Due string ` + "`json:\"due\"`" + `; Cents int64 ` + "`json:\"cents\"`" + ` }

func ListOverdue(ctx context.Context, db *sql.DB, cutoff time.Time) ([]Invoice, error) {
	// TODO: query, scan, close, check rows.Err, and discard partial failures.
	_ = fmt.Errorf
	return []Invoice{}, nil
}
`,
		tests: tests,
		build: buildInvoicesHarness,
	}
}

// buildInvoicesHarness supplies scripted rows so query, scan, close, and terminal row failures are independently observable.
func buildInvoicesHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `
var assessmentInvoiceFailure=errors.New("database failure");var assessmentInvoiceScenario *assessmentInvoiceScript
type assessmentInvoiceScript struct{mode string;events []string}
type assessmentInvoiceDriver struct{};func(assessmentInvoiceDriver)Open(string)(driver.Conn,error){return &assessmentInvoiceConn{assessmentInvoiceScenario},nil}
type assessmentInvoiceConn struct{script *assessmentInvoiceScript};func(*assessmentInvoiceConn)Prepare(string)(driver.Stmt,error){return nil,assessmentInvoiceFailure};func(*assessmentInvoiceConn)Close()error{return nil};func(*assessmentInvoiceConn)Begin()(driver.Tx,error){return nil,assessmentInvoiceFailure}
func(connection *assessmentInvoiceConn)QueryContext(ctx context.Context,query string,args []driver.NamedValue)(driver.Rows,error){connection.script.events=append(connection.script.events,"query");if len(args)!=1||!strings.Contains(strings.ToLower(query),"order by"){return nil,assessmentInvoiceFailure};if connection.script.mode=="query"{return nil,assessmentInvoiceFailure};return &assessmentInvoiceRows{script:connection.script},nil}
type assessmentInvoiceRows struct{script *assessmentInvoiceScript;index int};func(*assessmentInvoiceRows)Columns()[]string{return []string{"id","customer","due_at","cents"}};func(rows *assessmentInvoiceRows)Close()error{rows.script.events=append(rows.script.events,"close");return nil}
func(rows *assessmentInvoiceRows)Next(values []driver.Value)error{if rows.script.mode=="empty"{return io.EOF};if rows.index>=2{if rows.script.mode=="rows"{return assessmentInvoiceFailure};return io.EOF};values[0]=int64(rows.index+1);values[1]=[]string{"Ada","Lin"}[rows.index];values[2]=time.Date(2026,1,rows.index+1,0,0,0,0,time.UTC);values[3]=int64([]int{900,1200}[rows.index]);if rows.script.mode=="scan"&&rows.index==1{values[0]="bad"};rows.index++;return nil}
func init(){sql.Register("imperative-invoices",assessmentInvoiceDriver{})}`
	loop := `var tests []wireCase;_=json.Unmarshal([]byte(raw),&tests);modes:=[]string{"","empty","query","scan","rows"};for _,current:=range tests{current:=current;assessmentRun(current.ID,func()string{script:=&assessmentInvoiceScript{mode:modes[current.Payload.Case]};assessmentInvoiceScenario=script;db,_:=sql.Open("imperative-invoices","");defer db.Close();value,err:=ListOverdue(context.Background(),db,time.Now());encoded,_:=json.Marshal(value);label:="ok";if err!=nil{label="wrong";if current.Payload.Case==3{label="scan failure"}else if errors.Is(err,assessmentInvoiceFailure){label="database failure"}};return string(encoded)+"|"+strings.Join(script.events,",")+"|"+label})}`
	return harness(commonImports("context", "database/sql", "database/sql/driver", "errors", "strings"), declarations, loop, selected)
}

// bulkOrderExercise defines validation and multi-row transactional insertion.
func bulkOrderExercise() advancedSpec {
	tests := advancedTests(18,
		[4]string{"Commit", "Inserts the order and every line before committing.", `ctx, db, order`, `ok|begin,order,line,line,commit`},
		[4]string{"Invalid order", "Rejects invalid input before opening a transaction.", `order with no lines`, `invalid order|`},
		[4]string{"Header failure", "Rolls back when the order insert fails.", `header insert fails`, `database failure|begin,order,rollback`},
		[4]string{"Line failure", "Rolls back after any line insert failure.", `second line fails`, `database failure|begin,order,line,line,rollback`},
		[4]string{"Commit failure", "Returns a failed commit.", `commit fails`, `database failure|begin,order,line,line,commit`},
	)
	return advancedSpec{
		id: 18, title: "Bulk Order", topic: "SQL · transactions · validation · repeated statements",
		signature:   "CreateOrder(ctx context.Context, db *sql.DB, order Order) error",
		objective:   "Insert one order and all of its lines atomically.",
		input:       "A context, database handle, and validated order with at least one line.",
		output:      "nil after commit, ErrInvalidOrder before database access, or a wrapped database error after rollback.",
		constraints: []string{"Order ID and Customer must be positive/non-empty; every line needs non-empty SKU, positive quantity, and non-negative cents.", "Insert the order header first, then lines in input order using parameterized ExecContext calls inside one transaction.", "Rollback every unsuccessful transaction and return commit failures."},
		hints:       []string{"Validate the complete object before BeginTx.", "Defer Rollback immediately after BeginTx succeeds.", "Use one loop over Lines and stop at the first failed ExecContext."},
		pitfalls:    []string{"Beginning a transaction before validation", "Inserting lines through db instead of tx", "Continuing after one line insert fails"},
		builtins:    []string{"len"},
		packages:    []string{"context", "database/sql", "errors", "fmt"},
		starter: `import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrInvalidOrder = errors.New("invalid order")
type OrderLine struct { SKU string; Quantity int; Cents int64 }
type Order struct { ID int64; Customer string; Lines []OrderLine }

func CreateOrder(ctx context.Context, db *sql.DB, order Order) error {
	// TODO: validate every field, transact, insert in order, and roll back failures.
	_ = fmt.Errorf
	return ErrInvalidOrder
}
`,
		tests: tests,
		build: buildBulkOrderHarness,
	}
}

// buildBulkOrderHarness records transaction events through a scripted database/sql driver.
func buildBulkOrderHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `
var assessmentOrderFailure=errors.New("database failure");var assessmentOrderScenario *assessmentOrderScript
type assessmentOrderScript struct{mode string;events []string;line int}
type assessmentOrderDriver struct{};func(assessmentOrderDriver)Open(string)(driver.Conn,error){return &assessmentOrderConn{assessmentOrderScenario},nil}
type assessmentOrderConn struct{script *assessmentOrderScript};func(*assessmentOrderConn)Prepare(string)(driver.Stmt,error){return nil,assessmentOrderFailure};func(*assessmentOrderConn)Close()error{return nil};func(connection *assessmentOrderConn)Begin()(driver.Tx,error){connection.script.events=append(connection.script.events,"begin");return &assessmentOrderTx{connection.script},nil};func(connection *assessmentOrderConn)BeginTx(context.Context,driver.TxOptions)(driver.Tx,error){return connection.Begin()}
func(connection *assessmentOrderConn)ExecContext(ctx context.Context,query string,args []driver.NamedValue)(driver.Result,error){normalized:=strings.ToLower(query);if strings.Contains(normalized,"order_lines"){connection.script.events=append(connection.script.events,"line");connection.script.line++;if connection.script.mode=="line"&&connection.script.line==2{return nil,assessmentOrderFailure};if len(args)<4{return nil,assessmentOrderFailure}}else{connection.script.events=append(connection.script.events,"order");if connection.script.mode=="order"{return nil,assessmentOrderFailure};if len(args)<2{return nil,assessmentOrderFailure}};return driver.RowsAffected(1),nil}
type assessmentOrderTx struct{script *assessmentOrderScript};func(tx *assessmentOrderTx)Commit()error{tx.script.events=append(tx.script.events,"commit");if tx.script.mode=="commit"{return assessmentOrderFailure};return nil};func(tx *assessmentOrderTx)Rollback()error{tx.script.events=append(tx.script.events,"rollback");return nil}
func init(){sql.Register("imperative-orders",assessmentOrderDriver{})}`
	loop := `var tests []wireCase;_=json.Unmarshal([]byte(raw),&tests);modes:=[]string{"","invalid","order","line","commit"};for _,current:=range tests{current:=current;assessmentRun(current.ID,func()string{script:=&assessmentOrderScript{mode:modes[current.Payload.Case]};assessmentOrderScenario=script;db,_:=sql.Open("imperative-orders","");defer db.Close();order:=Order{ID:7,Customer:"Ada",Lines:[]OrderLine{{SKU:"A",Quantity:1,Cents:100},{SKU:"B",Quantity:2,Cents:250}}};if current.Payload.Case==1{order.Lines=[]OrderLine{}};err:=CreateOrder(context.Background(),db,order);label:="ok";if errors.Is(err,ErrInvalidOrder){label="invalid order"}else if errors.Is(err,assessmentOrderFailure){label="database failure"}else if err!=nil{label="wrong"};return label+"|"+strings.Join(script.events,",")})}`
	return harness(commonImports("context", "database/sql", "database/sql/driver", "errors", "strings"), declarations, loop, selected)
}
