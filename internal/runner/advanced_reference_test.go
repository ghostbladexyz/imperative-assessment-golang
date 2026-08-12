package runner

import (
	"context"
	"testing"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

// TestAdvancedReferenceSolutionsPass proves each authored contract has a source-policy-compliant implementation.
func TestAdvancedReferenceSolutionsPass(t *testing.T) {
	localRunner := New("go", 1, nil)
	for sourceID, solution := range advancedReferenceSolutions() {
		level, found := assessment.FindExercise(assessment.ExerciseKey("advanced/" + sourceID))
		if !found {
			t.Fatalf("advanced exercise %s is missing", sourceID)
		}
		result := localRunner.Run(context.Background(), level, solution, nil)
		if !result.Passed {
			t.Fatalf(
				"advanced/%s reference failed (%d/%d): compile=%q runtime=%q results=%#v",
				sourceID, result.PassedCount, result.TotalCount,
				result.CompileError, result.RuntimeError, result.Results,
			)
		}
	}
}

// advancedReferenceSolutions keys executable examples by immutable advanced source identity.
func advancedReferenceSolutions() map[string]string {
	return map[string]string{
		"1":  advancedPortsSolution,
		"2":  advancedLedgerSolution,
		"3":  advancedGlyphsSolution,
		"4":  advancedDependencySolution,
		"5":  advancedWidgetSolution,
		"6":  advancedWorkersSolution,
		"7":  advancedTransferSolution,
		"8":  advancedMessagesSolution,
		"9":  advancedRoomsSolution,
		"10": advancedErrorInspectorSolution,
		"11": advancedCompensationSolution,
		"12": advancedProfileClientSolution,
		"13": advancedDocumentSolution,
		"14": advancedMazeSolution,
		"15": advancedFanInSolution,
		"16": advancedFirstSuccessSolution,
		"17": advancedInvoicesSolution,
		"18": advancedBulkOrderSolution,
	}
}

const advancedPortsSolution = `import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

var ErrInvalidPortSpec = errors.New("invalid port specification")

func ParsePorts(spec string) ([]int, error) {
	result := []int{}
	seen := make(map[int]bool)
	tokens := strings.Split(spec, ",")
	for _, raw := range tokens {
		token := strings.TrimSpace(raw)
		if token == "" { return []int{}, ErrInvalidPortSpec }
		parts := strings.Split(token, "-")
		if len(parts) > 2 { return []int{}, ErrInvalidPortSpec }
		parse := func(text string) (int, bool) {
			text = strings.TrimSpace(text)
			if text == "" { return 0, false }
			for _, char := range text { if char < '0' || char > '9' { return 0, false } }
			value, err := strconv.Atoi(text)
			return value, err == nil && value >= 1 && value <= 65535
		}
		first, valid := parse(parts[0])
		if !valid { return []int{}, ErrInvalidPortSpec }
		last := first
		if len(parts) == 2 {
			last, valid = parse(parts[1])
			if !valid || last < first || last-first+1 > 1024 { return []int{}, ErrInvalidPortSpec }
		}
		for port := first; port <= last; port++ { seen[port] = true }
	}
	for port := range seen { result = append(result, port) }
	sort.Ints(result)
	return result, nil
}`

const advancedLedgerSolution = `import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidEvent = errors.New("invalid event")
type UserSummary struct { User string; Total int64; Last string }

func SummarizeEvents(input string) ([]UserSummary, error) {
	type aggregate struct { total int64; last time.Time }
	values := make(map[string]aggregate)
	for index, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") { continue }
		reader := csv.NewReader(strings.NewReader(line))
		record, err := reader.Read()
		if err != nil { return []UserSummary{}, fmt.Errorf("line %d: %w", index+1, ErrInvalidEvent) }
		if _, err = reader.Read(); err != io.EOF || len(record) != 3 { return []UserSummary{}, fmt.Errorf("line %d: %w", index+1, ErrInvalidEvent) }
		user := strings.TrimSpace(record[1])
		stamp, stampErr := time.Parse(time.RFC3339, strings.TrimSpace(record[0]))
		delta, deltaErr := strconv.ParseInt(strings.TrimSpace(record[2]), 10, 64)
		if user == "" || stampErr != nil || deltaErr != nil { return []UserSummary{}, fmt.Errorf("line %d: %w", index+1, ErrInvalidEvent) }
		current := values[user]
		current.total += delta
		if current.last.IsZero() || stamp.After(current.last) { current.last = stamp }
		values[user] = current
	}
	result := []UserSummary{}
	for user, value := range values { result = append(result, UserSummary{User:user, Total:value.total, Last:value.last.UTC().Format(time.RFC3339)}) }
	sort.Slice(result, func(i,j int) bool { return result[i].User < result[j].User })
	return result, nil
}`

const advancedGlyphsSolution = `import (
	"errors"
	"fmt"
	"strings"
)

var ( ErrMissingGlyph = errors.New("missing glyph"); ErrGlyphHeight = errors.New("inconsistent glyph height") )

func RenderGlyphs(input string, glyphs map[rune][]string) (string, error) {
	height := 0
	for _, glyph := range glyphs {
		if len(glyph) == 0 { return "", ErrGlyphHeight }
		if height == 0 { height = len(glyph) } else if len(glyph) != height { return "", ErrGlyphHeight }
	}
	if height == 0 { return "", ErrGlyphHeight }
	for _, line := range strings.Split(input, "\n") { for _, char := range line { if _, found := glyphs[char]; !found { return "", fmt.Errorf("%w: U+%04X", ErrMissingGlyph, char) } } }
	var output strings.Builder
	for _, line := range strings.Split(input, "\n") {
		if line == "" { output.WriteByte('\n'); continue }
		for row := 0; row < height; row++ { for _, char := range line { output.WriteString(glyphs[char][row]) }; output.WriteByte('\n') }
	}
	return output.String(), nil
}`

const advancedDependencySolution = `import ("errors"; "fmt"; "sort")
var ( ErrInvalidGraph = errors.New("invalid dependency graph"); ErrDependencyCycle = errors.New("dependency cycle") )
type Task struct { Name string; DependsOn []string }
func BuildOrder(tasks []Task) ([]string, error) {
	indegree := make(map[string]int); next := make(map[string][]string)
	for _, task := range tasks {
		if task.Name == "" { return []string{}, fmt.Errorf("%w: empty task", ErrInvalidGraph) }
		if _, found := indegree[task.Name]; found { return []string{}, fmt.Errorf("%w: duplicate task %s", ErrInvalidGraph, task.Name) }
		indegree[task.Name] = 0
	}
	for _, task := range tasks {
		seen := make(map[string]bool)
		for _, dependency := range task.DependsOn {
			if _, found := indegree[dependency]; !found { return []string{}, fmt.Errorf("%w: missing task %s", ErrInvalidGraph, dependency) }
			if dependency == task.Name || seen[dependency] { return []string{}, fmt.Errorf("%w: dependency %s", ErrInvalidGraph, dependency) }
			seen[dependency] = true; indegree[task.Name]++; next[dependency] = append(next[dependency], task.Name)
		}
	}
	ready := []string{}; for name, degree := range indegree { if degree == 0 { ready = append(ready, name) } }; sort.Strings(ready)
	result := []string{}
	for len(ready) > 0 { name := ready[0]; ready = ready[1:]; result = append(result,name); for _, dependent := range next[name] { indegree[dependent]--; if indegree[dependent] == 0 { ready=append(ready,dependent); sort.Strings(ready) } } }
	if len(result) != len(tasks) { return []string{}, ErrDependencyCycle }
	return result,nil
}`

const advancedWidgetSolution = `import ("encoding/json"; "errors"; "io"; "net/http"; "strings")
var ( ErrWidgetNotFound = errors.New("widget not found"); ErrWidgetConflict = errors.New("widget conflict") )
type Widget struct { ID string ` + "`json:\"id\"`" + `; Name string ` + "`json:\"name\"`" + ` }
type WidgetStore interface { Get(string)(Widget,error); Put(Widget) error }
func NewWidgetHandler(store WidgetStore) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func(){ if recover() != nil { writer.WriteHeader(http.StatusInternalServerError) } }()
		if !strings.HasPrefix(request.URL.Path,"/widgets/") || strings.TrimPrefix(request.URL.Path,"/widgets/") == "" || strings.Contains(strings.TrimPrefix(request.URL.Path,"/widgets/"),"/") { writer.WriteHeader(http.StatusNotFound); return }
		id := strings.TrimPrefix(request.URL.Path,"/widgets/")
		switch request.Method {
		case http.MethodGet:
			widget, err := store.Get(id)
			if err != nil { if errors.Is(err,ErrWidgetNotFound){writer.WriteHeader(404)}else{writer.WriteHeader(500)}; return }
			writer.Header().Set("Content-Type","application/json"); writer.WriteHeader(200); _=json.NewEncoder(writer).Encode(widget)
		case http.MethodPut:
			var input struct{Name string ` + "`json:\"name\"`" + `}; decoder:=json.NewDecoder(request.Body); decoder.DisallowUnknownFields()
			if decoder.Decode(&input)!=nil || strings.TrimSpace(input.Name)=="" { writer.WriteHeader(400); return }
			if err:=decoder.Decode(&struct{}{}); err!=io.EOF { writer.WriteHeader(400); return }
			err:=store.Put(Widget{ID:id,Name:input.Name}); if errors.Is(err,ErrWidgetConflict){writer.WriteHeader(409);return}; if err!=nil{writer.WriteHeader(500);return}; writer.WriteHeader(204)
		default: writer.WriteHeader(405)
		}
	})
}`

const advancedWorkersSolution = `import ("context"; "errors"; "sync")
var ErrInvalidWorkerCount = errors.New("invalid worker count")
func MapConcurrent(ctx context.Context, values []int, workers int, transform func(context.Context,int)(int,error)) ([]int,error) {
	if workers<=0{return []int{},ErrInvalidWorkerCount}; if len(values)==0{return []int{},nil}; if err:=ctx.Err();err!=nil{return []int{},err}
	ctx,cancel:=context.WithCancel(ctx); defer cancel(); results:=make([]int,len(values)); jobs:=make(chan int); var group sync.WaitGroup; var once sync.Once; var first error
	worker:=func(){defer group.Done();for index:=range jobs{value,err:=transform(ctx,values[index]);if err!=nil{once.Do(func(){first=err;cancel()});continue};results[index]=value}}
	if workers>len(values){workers=len(values)}; group.Add(workers); for i:=0;i<workers;i++{go worker()}
	for index:=range values{select{case jobs<-index:case <-ctx.Done():break}};close(jobs);group.Wait();if first!=nil{return []int{},first};if err:=ctx.Err();err!=nil{return []int{},err};return results,nil
}`

const advancedTransferSolution = `import("context";"database/sql";"errors";"fmt")
var(ErrInvalidTransfer=errors.New("invalid transfer");ErrInsufficientFunds=errors.New("insufficient funds"))
func TransferCredits(ctx context.Context,db *sql.DB,fromID,toID,amount int64)error{
	if db==nil||fromID<=0||toID<=0||fromID==toID||amount<=0{return ErrInvalidTransfer};tx,err:=db.BeginTx(ctx,nil);if err!=nil{return fmt.Errorf("begin: %w",err)};defer tx.Rollback();var balance int64
	if err=tx.QueryRowContext(ctx,"SELECT balance FROM accounts WHERE id = ?",fromID).Scan(&balance);err!=nil{return fmt.Errorf("balance: %w",err)};if balance<amount{return ErrInsufficientFunds}
	if _,err=tx.ExecContext(ctx,"UPDATE accounts SET balance = balance - ? WHERE id = ?",amount,fromID);err!=nil{return fmt.Errorf("debit: %w",err)}
	if _,err=tx.ExecContext(ctx,"UPDATE accounts SET balance = balance + ? WHERE id = ?",amount,toID);err!=nil{return fmt.Errorf("credit: %w",err)}
	if err=tx.Commit();err!=nil{return fmt.Errorf("commit: %w",err)};return nil
}`

const advancedMessagesSolution = `import("bufio";"bytes";"encoding/json";"errors";"fmt";"io";"strings")
var(ErrInvalidMessage=errors.New("invalid message");ErrDuplicateMessageID=errors.New("duplicate message id");ErrMessageTooLarge=errors.New("message too large"))
type Message struct{ID string ` + "`json:\"id\"`" + `;Payload string ` + "`json:\"payload\"`" + `}
func DecodeMessages(source io.Reader,max int)([]Message,error){result:=[]Message{};if max<=0{return result,ErrInvalidMessage};reader:=bufio.NewReader(source);seen:=make(map[string]bool);lineNumber:=0
	for{line,readErr:=reader.ReadString('\n');if len(line)>0{lineNumber++;content:=strings.TrimSuffix(strings.TrimSuffix(line,"\n"),"\r");if len(content)>max{return result,fmt.Errorf("line %d: %w",lineNumber,ErrMessageTooLarge)};trimmed:=strings.TrimSpace(content);if trimmed!=""{var message Message;decoder:=json.NewDecoder(bytes.NewBufferString(trimmed));decoder.DisallowUnknownFields();if decoder.Decode(&message)!=nil||message.ID==""||message.Payload==""{return result,fmt.Errorf("line %d: %w",lineNumber,ErrInvalidMessage)};if err:=decoder.Decode(&struct{}{});err!=io.EOF{return result,fmt.Errorf("line %d: %w",lineNumber,ErrInvalidMessage)};if seen[message.ID]{return result,fmt.Errorf("line %d: %w",lineNumber,ErrDuplicateMessageID)};seen[message.ID]=true;result=append(result,message)}}
		if readErr==io.EOF{return result,nil};if readErr!=nil{return result,fmt.Errorf("read messages: %w",readErr)}
	}}
`

const advancedRoomsSolution = `import("container/heap";"errors";"fmt";"sort")
var ErrInvalidMeeting=errors.New("invalid meeting")
type Meeting struct{ID string;Start int;End int}
func AssignRooms(meetings []Meeting)(map[string]int,error){result:=map[string]int{};seen:=make(map[string]bool);items:=append([]Meeting{},meetings...);for _,item:=range items{if item.ID==""||item.Start>=item.End{return map[string]int{},fmt.Errorf("%w: %s",ErrInvalidMeeting,item.ID)};if seen[item.ID]{return map[string]int{},fmt.Errorf("%w: duplicate %s",ErrInvalidMeeting,item.ID)};seen[item.ID]=true};sort.Slice(items,func(i,j int)bool{if items[i].Start!=items[j].Start{return items[i].Start<items[j].Start};if items[i].End!=items[j].End{return items[i].End<items[j].End};return items[i].ID<items[j].ID});type occupied struct{room,end int};active:=[]occupied{};available:=[]int{};next:=1;_ = heap.Init
	for _,item:=range items{remaining:=[]occupied{};for _,current:=range active{if current.end<=item.Start{available=append(available,current.room)}else{remaining=append(remaining,current)}};active=remaining;sort.Ints(available);room:=next;if len(available)>0{room=available[0];available=available[1:]}else{next++};result[item.ID]=room;active=append(active,occupied{room:room,end:item.End})};return result,nil}
`

const advancedErrorInspectorSolution = `import "errors"
var(ErrUnavailable=errors.New("unavailable");ErrUnauthorized=errors.New("unauthorized"))
type FieldError struct{Field string};func(err *FieldError)Error()string{return "invalid field: "+err.Field}
type Failure struct{Kind string ` + "`json:\"kind\"`" + `;Field string ` + "`json:\"field\"`" + `;Temporary bool ` + "`json:\"temporary\"`" + `}
func InspectError(err error)Failure{if err==nil{return Failure{Kind:"none"}};result:=Failure{Kind:"unknown",Temporary:errors.Is(err,ErrUnavailable)};if errors.Is(err,ErrUnauthorized){result.Kind="unauthorized";return result};var field *FieldError;if errors.As(err,&field){result.Kind="field";result.Field=field.Field;return result};if result.Temporary{result.Kind="unavailable"};return result}`

const advancedCompensationSolution = `import("context";"errors")
type Step struct{Run func(context.Context)error;Undo func(context.Context)error}
func ExecuteSteps(ctx context.Context,steps []Step)error{completed:=0;var trigger error;for index,step:=range steps{if err:=ctx.Err();err!=nil{trigger=err;break};if err:=step.Run(ctx);err!=nil{trigger=err;break};completed=index+1};if trigger==nil{return nil};failures:=[]error{trigger};for index:=completed-1;index>=0;index--{if steps[index].Undo==nil{continue};if err:=steps[index].Undo(ctx);err!=nil{failures=append(failures,err)}};return errors.Join(failures...)}`

const advancedProfileClientSolution = `import("bytes";"context";"encoding/json";"errors";"fmt";"io";"net/http")
var(ErrRemoteResponse=errors.New("remote response");ErrInvalidResponse=errors.New("invalid response"));type Profile struct{ID string ` + "`json:\"id\"`" + `;Name string ` + "`json:\"name\"`" + `}
func FetchProfile(ctx context.Context,client *http.Client,url string)(Profile,error){if client==nil{return Profile{},ErrInvalidResponse};request,err:=http.NewRequestWithContext(ctx,http.MethodGet,url,nil);if err!=nil{return Profile{},fmt.Errorf("request: %w",err)};response,err:=client.Do(request);if err!=nil{return Profile{},fmt.Errorf("fetch: %w",err)};defer response.Body.Close();if response.StatusCode<200||response.StatusCode>=300{return Profile{},ErrRemoteResponse};data,err:=io.ReadAll(io.LimitReader(response.Body,64*1024+1));if err!=nil||len(data)>64*1024{return Profile{},ErrInvalidResponse};decoder:=json.NewDecoder(bytes.NewReader(data));decoder.DisallowUnknownFields();var profile Profile;if decoder.Decode(&profile)!=nil||profile.ID==""||profile.Name==""{return Profile{},ErrInvalidResponse};if err:=decoder.Decode(&struct{}{});err!=io.EOF{return Profile{},ErrInvalidResponse};return profile,nil}`

const advancedDocumentSolution = `import("context";"errors";"net/http";"strconv";"strings")
var ErrDocumentNotFound=errors.New("document not found");type Document struct{ID string;Body string;Version int}
func NewDocumentHandler(load func(context.Context,string)(Document,error))http.Handler{return http.HandlerFunc(func(writer http.ResponseWriter,request *http.Request){if !strings.HasPrefix(request.URL.Path,"/docs/"){writer.WriteHeader(404);return};id:=strings.TrimPrefix(request.URL.Path,"/docs/");if id==""||strings.Contains(id,"/"){writer.WriteHeader(404);return};if request.Method!="GET"&&request.Method!="HEAD"{writer.WriteHeader(405);return};document,err:=load(request.Context(),id);if errors.Is(err,ErrDocumentNotFound){writer.WriteHeader(404);return};if err!=nil{writer.WriteHeader(500);return};etag:=strconv.Quote("v"+strconv.Itoa(document.Version));writer.Header().Set("ETag",etag);if request.Header.Get("If-None-Match")==etag{writer.WriteHeader(304);return};writer.WriteHeader(200);if request.Method=="GET"{_,_=writer.Write([]byte(document.Body))}})}`

const advancedMazeSolution = `import("container/heap";"errors")
var(ErrInvalidGrid=errors.New("invalid grid");ErrNoPath=errors.New("no path"));type mazeNode struct{cost,row,col int};type mazeQueue []mazeNode;func(queue mazeQueue)Len()int{return len(queue)};func(queue mazeQueue)Less(i,j int)bool{return queue[i].cost<queue[j].cost};func(queue mazeQueue)Swap(i,j int){queue[i],queue[j]=queue[j],queue[i]};func(queue *mazeQueue)Push(value any){*queue=append(*queue,value.(mazeNode))};func(queue *mazeQueue)Pop()any{old:=*queue;value:=old[len(old)-1];*queue=old[:len(old)-1];return value}
func CheapestPath(grid [][]int)(int,error){if len(grid)==0||len(grid[0])==0{return 0,ErrInvalidGrid};width:=len(grid[0]);for _,row:=range grid{if len(row)!=width{return 0,ErrInvalidGrid};for _,value:=range row{if value < -1{return 0,ErrInvalidGrid}}};if grid[0][0]==-1||grid[len(grid)-1][width-1]==-1{return 0,ErrNoPath};best:=make([][]int,len(grid));for row:=range best{best[row]=make([]int,width);for column:=range best[row]{best[row][column]=int(^uint(0)>>1)}};best[0][0]=grid[0][0];queue:=&mazeQueue{{cost:grid[0][0]}};heap.Init(queue);moves:=[][2]int{{1,0},{-1,0},{0,1},{0,-1}};for queue.Len()>0{current:=heap.Pop(queue).(mazeNode);if current.cost!=best[current.row][current.col]{continue};if current.row==len(grid)-1&&current.col==width-1{return current.cost,nil};for _,move:=range moves{row,column:=current.row+move[0],current.col+move[1];if row<0||row>=len(grid)||column<0||column>=width||grid[row][column]==-1{continue};cost:=current.cost+grid[row][column];if cost<best[row][column]{best[row][column]=cost;heap.Push(queue,mazeNode{cost:cost,row:row,col:column})}}};return 0,ErrNoPath}`

const advancedFanInSolution = `import("context";"sync")
func FanIn(ctx context.Context,inputs ...<-chan int)<-chan int{output:=make(chan int);var group sync.WaitGroup;for _,input:=range inputs{if input==nil{continue};group.Add(1);go func(source <-chan int){defer group.Done();for{select{case<-ctx.Done():return;case value,open:=<-source:if !open{return};select{case output<-value:case<-ctx.Done():return}}}}(input)};go func(){group.Wait();close(output)}();return output}`

const advancedFirstSuccessSolution = `import("context";"errors";"sync")
var ErrNoTasks=errors.New("no tasks");type successResult struct{index int;value string;err error}
func FirstSuccess(ctx context.Context,tasks ...func(context.Context)(string,error))(string,error){if len(tasks)==0{return "",ErrNoTasks};child,cancel:=context.WithCancel(ctx);defer cancel();results:=make(chan successResult,len(tasks));var group sync.WaitGroup;group.Add(len(tasks));for index,task:=range tasks{go func(index int,task func(context.Context)(string,error)){defer group.Done();value,err:=task(child);results<-successResult{index,value,err}}(index,task)};failures:=make([]error,len(tasks));winner:="";won:=false;for range tasks{result:=<-results;if result.err==nil&&!won{winner=result.value;won=true;cancel()}else if result.err!=nil{failures[result.index]=result.err}};group.Wait();if won{return winner,nil};ordered:=[]error{};for _,err:=range failures{if err!=nil{ordered=append(ordered,err)}};return "",errors.Join(ordered...)}`

const advancedInvoicesSolution = `import("context";"database/sql";"fmt";"time")
type Invoice struct{ID int64 ` + "`json:\"id\"`" + `;Customer string ` + "`json:\"customer\"`" + `;Due string ` + "`json:\"due\"`" + `;Cents int64 ` + "`json:\"cents\"`" + `}
func ListOverdue(ctx context.Context,db *sql.DB,cutoff time.Time)([]Invoice,error){rows,err:=db.QueryContext(ctx,"SELECT id, customer, due_at, cents FROM invoices WHERE due_at < ? ORDER BY due_at, id",cutoff);if err!=nil{return []Invoice{},fmt.Errorf("query: %w",err)};defer rows.Close();result:=[]Invoice{};for rows.Next(){var invoice Invoice;var due time.Time;if err:=rows.Scan(&invoice.ID,&invoice.Customer,&due,&invoice.Cents);err!=nil{return []Invoice{},fmt.Errorf("scan: %w",err)};invoice.Due=due.UTC().Format("2006-01-02");result=append(result,invoice)};if err:=rows.Err();err!=nil{return []Invoice{},fmt.Errorf("rows: %w",err)};return result,nil}`

const advancedBulkOrderSolution = `import("context";"database/sql";"errors";"fmt")
var ErrInvalidOrder=errors.New("invalid order");type OrderLine struct{SKU string;Quantity int;Cents int64};type Order struct{ID int64;Customer string;Lines []OrderLine}
func CreateOrder(ctx context.Context,db *sql.DB,order Order)error{if db==nil||order.ID<=0||order.Customer==""||len(order.Lines)==0{return ErrInvalidOrder};for _,line:=range order.Lines{if line.SKU==""||line.Quantity<=0||line.Cents<0{return ErrInvalidOrder}};tx,err:=db.BeginTx(ctx,nil);if err!=nil{return fmt.Errorf("begin: %w",err)};defer tx.Rollback();if _,err=tx.ExecContext(ctx,"INSERT INTO orders (id, customer) VALUES (?, ?)",order.ID,order.Customer);err!=nil{return fmt.Errorf("order: %w",err)};for _,line:=range order.Lines{if _,err=tx.ExecContext(ctx,"INSERT INTO order_lines (order_id, sku, quantity, cents) VALUES (?, ?, ?, ?)",order.ID,line.SKU,line.Quantity,line.Cents);err!=nil{return fmt.Errorf("line: %w",err)}};if err=tx.Commit();err!=nil{return fmt.Errorf("commit: %w",err)};return nil}`
