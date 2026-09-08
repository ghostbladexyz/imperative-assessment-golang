package assessment

import (
	"fmt"
	"strings"
)

type checkSpec struct {
	name, purpose, input, need string
	labels                     []string
}

func officialChecks(slug string) []VisibleTest {
	checks := officialCheckSpecs[slug]
	result := make([]VisibleTest, 0, len(checks))
	for index, check := range checks {
		labels := check.labels
		if len(labels) == 0 {
			labels = []string{check.name}
		}
		result = append(result, VisibleTest{
			ID: fmt.Sprintf("%s-%02d", slug, index+1), Name: check.name,
			Purpose: check.purpose, Input: check.input, Expected: check.need, labels: labels,
		})
	}
	return result
}

func cases(names ...string) []checkSpec {
	result := make([]checkSpec, 0, len(names))
	for _, name := range names {
		detail, found := officialDetails[name]
		if !found {
			panic("missing official check detail: " + name)
		}
		purpose := "Fixed fixture from the official grader."
		if strings.Contains(name, "secret") || strings.Contains(name, "seeded") || strings.Contains(detail.input, "generated") {
			purpose = "Seeded official check; a failure shows this run's concrete values and replay seed."
		}
		result = append(result, checkSpec{name: name, purpose: purpose, input: detail.input, need: detail.need})
	}
	return result
}

type checkDetail struct{ input, need string }

var officialDetails = map[string]checkDetail{
	"tp4/distinct_ok":      {`argv: ["3", "1", "2"]`, `exit 0; stdout and stderr empty`},
	"tp4/duplicate":        {`argv: ["1", "2", "2", "3"]`, `exit 1; stderr exactly "Error\n"; stdout empty`},
	"tp4/non_integer":      {`argv: ["0", "one", "2"]`, `exit 1; stderr exactly "Error\n"; stdout empty`},
	"tp4/empty_ok":         {`argv: []`, `exit 0; stdout and stderr empty`},
	"tp4/single_ok":        {`argv: ["42"]`, `exit 0; stdout and stderr empty`},
	"tp4/float_rejected":   {`argv: ["1.5", "2", "3"]`, `exit 1; stderr exactly "Error\n"; stdout empty`},
	"tp4/dup_leading_zero": {`argv: ["1", "01"]`, `numeric duplicate rejected: exit 1 and stderr "Error\n"`},
	"tp4/dup_signed_zero":  {`argv: ["-0", "0"]`, `numeric duplicate rejected: exit 1 and stderr "Error\n"`},
	"tp4/secret-valid":     {`argv: [n, n+1, -(n+2)], n generated in [1000,8999]`, `accept three distinct decimal integers`},
	"tp4/secret-duplicate": {`argv: [n, "0"+n], n generated in [1000,8999]`, `reject equal numeric values despite different spelling`},
	"tp4/secret-invalid":   {`argv: [n, "invalid"+k], generated n and k`, `reject the invalid integer; exit 1 and stderr "Error\n"`},
	"tp4/secret-overflow":  {`argv: ["999999999999999999999999999999"]`, `reject integer overflow; exit 1 and stderr "Error\n"`},
	"tp4/secret-decimal":   {`argv: ["0x10"]`, `reject non-decimal grammar; exit 1 and stderr "Error\n"`},
	"tp4/secret-wide":      {`argv: ["+10", "2147483648", "-2147483649"]`, `accept distinct 64-bit decimal integers`},

	"rb1/00-empty":                  {`file bytes: ""`, `stdout exactly ""; no stderr; exit 0`},
	"rb1/01-blank":                  {`file bytes: "   \t  "`, `one output line: 0`},
	"rb1/02-trailing-input-newline": {`file bytes: "1 2\n"`, `one output line: 3, with exact final newline`},
	"rb1/03-int64-boundaries":       {`lines: 9223372036854775807; -9223372036854775808; both together`, `echo each boundary sum, then -1; byte-exact stdout`},
	"rb1/04-whitespace":             {`line: +7<NBSP>-3  9; next line: "  -0 +001\t2"`, `sums are 13 and 3; all Unicode/ASCII whitespace accepted`},
	"rb1/05-first-invalid":          {`lines: "5 first_invalid 7" and "9"`, `first line reports first_invalid exactly; second line still outputs 9`},
	"rb1/06-overflow-token":         {`lines: 9223372036854775808 and -9223372036854775809`, `both overflowing tokens reported as invalid, without panic`},
	"rb1/07-multiple-invalid":       {`line: "1 first\"bad 2 second\\bad"`, `report only the first invalid token with exact quoting`},
	"rb1/08-seeded-valid":           {`6 generated lines; each has 1–8 integers in [-1,000,000,1,000,000]`, `byte-exact per-line int64 sums; no stderr`},
	"rb1/09-seeded-invalid":         {`lines include generated token bad_<64-bit hex>`, `report the generated invalid token exactly and continue line processing`},

	"GET / 200 home":                         {`GET /`, `status 200; body exactly "home"`},
	"POST /submit named 200 submitted":       {`POST /submit, form name=u<generated integer>`, `status 200; body exactly "submitted"`},
	"POST /submit missing name 400":          {`POST /submit, empty body`, `status 400`},
	"POST /submit empty name 400":            {`POST /submit, form name=`, `status 400`},
	"GET /submit 405":                        {`GET /submit`, `status 405`},
	"PUT /submit 405":                        {`PUT /submit`, `status 405`},
	"PATCH /submit 405":                      {`PATCH /submit`, `status 405`},
	"DELETE /submit 405":                     {`DELETE /submit`, `status 405`},
	"POST / 405":                             {`POST /`, `status 405`},
	"PUT / 405":                              {`PUT /`, `status 405`},
	"PATCH / 405":                            {`PATCH /`, `status 405`},
	"DELETE / 405":                           {`DELETE /`, `status 405`},
	"GET unknown 404":                        {`GET /x<generated integer>`, `status 404`},
	"POST unknown 404":                       {`POST /x<generated integer>`, `status 404`},
	"GET / serves form 200":                  {`GET /`, `status 200; body contains "<form"`},
	"POST /generate uppercases 200":          {`POST /generate, text=<5–9 generated lowercase letters>`, `status 200; body is the exact uppercase input`},
	"POST /generate empty 400":               {`POST /generate, form text=`, `status 400`},
	"GET unknown path 404":                   {`GET /nope`, `status 404`},
	"POST /generate panic text 500 no crash": {`POST /generate, form text=panic`, `status 500 and handler must not panic`},
	"GET / again still 200":                  {`GET / after the panic-text case`, `status 200; handler remains healthy`},

	"rb2/empty":       {`file text: ""`, `no words; stdout empty`},
	"rb2/one":         {`one generated token w<12 hex digits>`, `that token followed by count 1`},
	"rb2/ties":        {`beta alpha beta alpha gamma`, `alpha 2, beta 2, gamma 1; ties alphabetical`},
	"rb2/punctuation": {`go go, go go, rust! rust`, `punctuation stays in tokens; frequency then alphabetical ordering`},
	"rb2/whitespace":  {`"a\t\tb\n a<NBSP>c c"`, `Unicode/ASCII whitespace split; deterministic counts`},
	"rb2/seeded":      {`generated shuffled bag of 19 unique w<hex> tokens, repeated 1–5 times`, `frequency-descending then alphabetical output, identical across 20 runs`},

	"sy1/seeded-0": {`1 generated row: type f/d/l, mode, owner/group, size, epoch, name`, `byte-exact ls-style permissions and Europe/Athens timestamp`},
	"sy1/seeded-1": {`2 generated rows; second filename contains a space`, `two byte-exact formatted rows in input order`},
	"sy1/seeded-2": {`3 generated rows; alternating filenames may contain spaces`, `three byte-exact formatted rows in input order`},
	"sy1/seeded-3": {`4 generated rows; modes from 000,001,111,204,640,755,777`, `four byte-exact formatted rows in input order`},
	"sy1/seeded-4": {`5 generated rows; epochs include 1970, 2000, 2024, 2024 year-end`, `five byte-exact formatted rows in input order`},

	"tp3/a":                     {`input file text: "HI"`, `render H then I from banner.txt, row by row, byte-exact`},
	"tp3/b":                     {`input file text: "GO HI"`, `render all five glyphs, including space, row by row`},
	"tp3/c":                     {`input file text: literal "HI\\nGO"`, `two rendered blocks in order, one block per literal \\n segment`},
	"tp3/secret-mutated-banner": {`generated banner mutates every glyph row; input uses ASCII 32–126 plus "\\n\\n~A0!"`, `look up glyphs from supplied banner; preserve empty-segment/block ordering exactly`},

	"tp1/a":               {`Hello WORLD (low) and (cap) friends`, `Hello world And friends`},
	"tp1/b":               {`this is so cool (up, 2)`, `this is SO COOL`},
	"tp1/c":               {`the quick brown fox (rev, 3)`, `the fox brown quick`},
	"tp1/d":               {`alpha beta gamma (rev, 2) (up)`, `alpha gamma BETA`},
	"tp1/e":               {`hi (up, 5)`, `HI`},
	"tp1/f":               {`foo bar (rev, 5)`, `bar foo`},
	"tp1/g":               {`red green , blue (rev, 3)`, `blue green, red`},
	"tp1/h":               {`one two three (up, 2) (rev, 2)`, `one THREE TWO`},
	"tp1/secret-twin":     {`generated words with a randomized (rev, n) directive`, `reverse exactly the selected word positions; punctuation positions stay fixed`},
	"tp1/secret-case":     {`generated punctuation span plus case directive`, `apply case only to eligible preceding words; byte-exact punctuation`},
	"tp1/secret-newline":  {`generated transformation input ending in newline`, `correct transformation with exactly one final newline`},
	"tp1/secret-zero":     {`generated case directive with count 0`, `non-positive count changes no words and does not panic`},
	"tp1/secret-negative": {`generated reversal directive with negative count`, `non-positive count changes no words and does not panic`},

	"tp2/a":              {`1E (hex) files and 1010 (bin) bytes`, `30 files and 10 bytes`},
	"tp2/b":              {`A apple costs a euro`, `An apple costs an euro`},
	"tp2/c":              {`Hello , world . How are you ?`, `Hello, world. How are you?`},
	"tp2/cons":           {`a dog and A cat`, `a dog and A cat`},
	"tp2/d":              {`Wait ... what !?`, `Wait... what!?`},
	"tp2/e":              {`I have FF (hex) apples , and a orange .`, `I have 255 apples, and an orange.`},
	"tp2/f":              {`zzz (hex) and ff (hex) done`, `zzz and 255 done`},
	"tp2/g":              {`ff , (hex) done`, `255, done`},
	"tp2/h":              {`1111 (bin) items cost A euro`, `15 items cost An euro`},
	"tp2/secret-twin":    {`generated number 1–4000 as hex/bin, plus a/A and generated vowel word`, `convert the number, repair the article, and assemble exact text`},
	"tp2/secret-convert": {`(bin) 102 (bin) and <generated valid number> (hex|bin)`, `leave invalid/missing operands intact; convert the valid operand`},
	"tp2/secret-article": {`a , orange and a A (hex)`, `respect punctuation boundaries and article/conversion interaction`},
	"tp2/secret-newline": {`A (hex) followed by newline`, `10 followed by exactly one newline`},

	"al1":               {`4 ants; paths start-a-end and start-b-end`, `print any valid shortest path from start to end`},
	"al1b":              {`3 ants; routes start-x-y-end and start-z-end`, `choose direct-shorter route start-z-end`},
	"al1c":              {`3 ants; start-a and b-end are disconnected`, `report ERROR for unreachable end`},
	"al1d":              {`4 ants; routes of 3 and 4 tunnels`, `choose shortest start-a-b-end path`},
	"al1/direct":        {`generated direct start-end colony`, `print the direct shortest path`},
	"al1/secret-seeded": {`generated colony with alternate routes`, `print a valid shortest path; any shortest tie is accepted`},
	"al2":               {`6 ants; path lengths 3 and 4`, `print movements with the oracle's minimum makespan`},
	"al2b":              {`5 ants; path lengths 2 and 3`, `minimum-makespan schedule with exact valid movement format`},
	"al2c":              {`9 ants; path lengths 2, 3, and 4`, `minimum-makespan schedule across useful paths`},
	"al2d":              {`3 ants; single path start-a-end`, `queue ants correctly on the single lane`},
	"al2/secret-seeded": {`generated colony and ant count`, `minimum-makespan schedule matching the compiled oracle`},
	"al2/zero-ants":     {`generated valid colony with 0 ants`, `zero movement lines; no error`},
	"al2/no-path":       {`generated disconnected colony`, `report ERROR`},
	"al2/order-a":       {`generated paths in ordering A`, `optimal output independent of path declaration order`},
	"al2/order-b":       {`same graph as order-a with permuted path order`, `same optimal scheduling behavior as order-a`},

	"dedup and switch":                       {`user reacts twice, then switches like ↔ dislike on one post`, `exactly one row for (user,post); count reflects latest value`},
	"randomized counts":                      {`generated users/posts with mixed like/dislike values`, `exact like/dislike counts for only the requested post`},
	"untouched post is zero":                 {`count on a post with no rows`, `likes=0, dislikes=0, error=nil`},
	"one row per user post":                  {`repeated replacements for one generated (user,post)`, `database contains exactly one row for that pair`},
	"atomic rollback on commit error":        {`transaction driver forces Commit to fail`, `return the commit error and expose none of the replacement`},
	"atomic rollback on begin error":         {`transaction driver forces Begin to fail`, `return the begin error; perform no write`},
	"atomic rollback on write error":         {`transaction driver fails UPDATE/INSERT replacement`, `return the write error and preserve the old committed reaction`},
	"concurrent writes keep one row":         {`two concurrent setReaction calls for one (user,post)`, `serialize atomically; finish with exactly one valid row`},
	"uncommitted replacement is not visible": {`count runs while replacement transaction is paused`, `reader never observes an intermediate delete/partial replacement`},
	"setReaction sql error":                  {`closed/failing database passed to setReaction`, `return the SQL error instead of swallowing or panicking`},
	"count sql error":                        {`driver returns an iteration/query error`, `return the SQL error; do not report misleading counts`},
}

var officialCheckSpecs = map[string][]checkSpec{
	"validate-stack": cases(
		"tp4/distinct_ok", "tp4/duplicate", "tp4/non_integer", "tp4/empty_ok", "tp4/single_ok",
		"tp4/float_rejected", "tp4/dup_leading_zero", "tp4/dup_signed_zero", "tp4/secret-valid",
		"tp4/secret-duplicate", "tp4/secret-invalid", "tp4/secret-overflow", "tp4/secret-decimal", "tp4/secret-wide",
	),
	"safe-sum": cases(
		"rb1/00-empty", "rb1/01-blank", "rb1/02-trailing-input-newline", "rb1/03-int64-boundaries",
		"rb1/04-whitespace", "rb1/05-first-invalid", "rb1/06-overflow-token", "rb1/07-multiple-invalid",
		"rb1/08-seeded-valid", "rb1/09-seeded-invalid",
	),
	"tetris": {{name: "al4/exhaustive", purpose: "Checks all 19 pieces, undo behavior, and 12 seeded multi-piece boards.", input: "19 single-piece shapes/rotations; undo O→J; 12 generated boards of 2–4 pieces", need: "valid non-overlapping placement of every piece in the minimum square; exact checker-approved board", labels: []string{"al4/exhaustive", "al4/*"}}},
	"method-routing": cases(
		"GET / 200 home", "POST /submit named 200 submitted", "POST /submit missing name 400",
		"POST /submit empty name 400", "GET /submit 405", "PUT /submit 405",
		"PATCH /submit 405", "DELETE /submit 405", "POST / 405", "PUT / 405",
		"PATCH / 405", "DELETE / 405", "GET unknown 404", "POST unknown 404",
	),
	"status-matrix": cases(
		"GET / serves form 200", "POST /generate uppercases 200", "POST /generate empty 400",
		"GET unknown path 404", "POST /generate panic text 500 no crash", "GET / again still 200",
	),
	"wordcount":    cases("rb2/empty", "rb2/one", "rb2/ties", "rb2/punctuation", "rb2/whitespace", "rb2/seeded"),
	"ls-format":    cases("sy1/seeded-0", "sy1/seeded-1", "sy1/seeded-2", "sy1/seeded-3", "sy1/seeded-4"),
	"ascii-render": cases("tp3/a", "tp3/b", "tp3/c", "tp3/secret-mutated-banner"),
	"reloaded-rev": cases(
		"tp1/a", "tp1/b", "tp1/c", "tp1/d", "tp1/e", "tp1/f", "tp1/g", "tp1/h",
		"tp1/secret-twin", "tp1/secret-case", "tp1/secret-newline", "tp1/secret-zero", "tp1/secret-negative",
	),
	"reloaded-format": cases(
		"tp2/a", "tp2/b", "tp2/c", "tp2/cons", "tp2/d", "tp2/e", "tp2/f", "tp2/g", "tp2/h",
		"tp2/secret-twin", "tp2/secret-convert", "tp2/secret-article", "tp2/secret-newline",
	),
	"lemin-path": cases("al1", "al1b", "al1c", "al1d", "al1/direct", "al1/secret-seeded"),
	"lemin-why": cases(
		"al2", "al2b", "al2c", "al2d",
		"al2/secret-seeded", "al2/zero-ants", "al2/no-path", "al2/order-a", "al2/order-b",
	),
	"bounded-fanout": {
		{name: "cc1/empty", purpose: "Empty input completes cleanly.", input: "[]", need: "empty output; deterministic across 2 runs"},
		{name: "cc1/single", purpose: "One value is squared.", input: "[generated integer in -999..999]", need: "that integer squared at index 0; deterministic across 2 runs"},
		{name: "cc1/signed-duplicates", purpose: "Signed duplicate values retain order.", input: "[0, -7, 7, -7, 11]", need: "[0, 49, 49, 49, 121] in the same order"},
		{name: "cc1/seeded-small", purpose: "Results remain ordered across concurrent workers.", input: "17 generated integers in -5000..5000", need: "all squares in input order; deterministic across 2 runs"},
		{name: "cc1/seeded-large", purpose: "The bounded worker pool handles a larger workload.", input: "193 generated integers in -30000..30000", need: "all squares in input order; deterministic across 2 runs"},
		{name: "cc1/structural", purpose: "Source structure is inspected in addition to runtime output.", input: "submitted squareAll AST", need: "launch a goroutine that consumes the jobs channel and writes indexed output"},
	},
	"consume-join": {
		{name: "join Artist1 (seeded)", purpose: "Joins the first randomized artist and sorts locations.", input: "Artist1_<suffix>; generated 19xx album; relation [edinburgh, amsterdam, delhi]", need: "firstAlbum: <19xx>\nlocations: amsterdam, delhi, edinburgh\n", labels: []string{"join Artist1 *"}},
		{name: "join Artist2 (seeded)", purpose: "Joins the second randomized artist and sorts locations.", input: "Artist2_<suffix>; generated 20xx album; relation [cairo, berlin]", need: "firstAlbum: <20xx>\nlocations: berlin, cairo\n", labels: []string{"join Artist2 *"}},
		{name: "unknown artist returns error", purpose: "Returns an error for an absent artist.", input: "Nobody_<same generated suffix>", need: "non-nil error; no stale joined output"},
	},
	"push-swap": {{name: "al3/exhaustive", purpose: "Checks every permutation of 0–5 seeded values, allowed operations, and the 11-operation ceiling.", input: "all 154 permutations of prefixes of 5 generated distinct signed integers", need: "only allowed operations; stack A sorted; B empty; fewer than 12 operations; no stderr", labels: []string{"al3/exhaustive", "al3/case-*"}}},
	"broadcast": {{name: "race-safe broadcast suite", purpose: "Checks exact delivery, sender exclusion, registry cleanup, continued service, and Go race detection.", input: "net.Pipe clients alice, bob, carol; messages hello-<suffix>, world-<suffix>, still-<suffix>", need: "send exact [name]: message\n to every other client; remove disconnects under mutex; pass go test -race", labels: []string{"race-safe broadcast suite"}}},
	"reactions": cases(
		"dedup and switch", "randomized counts", "untouched post is zero", "one row per user post",
		"atomic rollback on commit error", "atomic rollback on begin error", "atomic rollback on write error",
		"concurrent writes keep one row", "uncommitted replacement is not visible", "setReaction sql error", "count sql error",
	),
}
