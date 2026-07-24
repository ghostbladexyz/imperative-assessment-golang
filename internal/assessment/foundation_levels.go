package assessment

import "fmt"

type practiceCase struct {
	name     string
	input    string
	expected string
	payload  any
}

type practiceSpec struct {
	id          int
	title       string
	topic       string
	difficulty  string
	signature   string
	objective   string
	input       string
	output      string
	constraints []string
	hints       []string
	pitfalls    []string
	builtins    []string
	packages    []string
	starter     string
	inputFields string
	call        string
	cases       []practiceCase
}

func foundationalLevels() []Level {
	specs := []practiceSpec{
		{
			id: 1, title: "Echo", topic: "Functions · strings", difficulty: "Beginner",
			signature: "Echo(value string) string",
			objective: "Return the string you receive.",
			input:     "One string.", output: "The same string, unchanged.",
			constraints: []string{"Do not add or remove characters.", "Return empty input as an empty string."},
			hints:       []string{"The parameter already contains the answer.", "Use return followed by a value.", "No loop or package is needed."},
			pitfalls:    []string{"Printing instead of returning", "Returning a fixed example"},
			starter: `func Echo(value string) string {
	// TODO: return value.
	return ""
}
`, inputFields: stringField, call: "Echo(current.Payload.Value)",
			cases: []practiceCase{
				pc("Empty", `""`, `""`, map[string]any{"value": ""}),
				pc("Word", `"go"`, `"go"`, map[string]any{"value": "go"}),
				pc("Spaces", `"hello world"`, `"hello world"`, map[string]any{"value": "hello world"}),
				pc("Symbols", `"!@#"`, `"!@#"`, map[string]any{"value": "!@#"}),
				pc("Unicode", `"γειά"`, `"γειά"`, map[string]any{"value": "γειά"}),
			},
		},
		{
			id: 2, title: "Plus One", topic: "Functions · integers", difficulty: "Beginner",
			signature: "Increment(value int) int",
			objective: "Add one to an integer.",
			input:     "One integer.", output: "The input plus one.",
			constraints: []string{"Return the result; do not print it.", "Negative values are valid."},
			hints:       []string{"Use the + operator.", "Integer arithmetic needs no package.", "Return the expression directly."},
			pitfalls:    []string{"Printing the result", "Changing the function signature"},
			starter: `func Increment(value int) int {
	// TODO: add one.
	return 0
}
`, inputFields: intField, call: "Increment(current.Payload.Value)",
			cases: []practiceCase{
				pc("Zero", `0`, `1`, map[string]any{"value": 0}),
				pc("Positive", `4`, `5`, map[string]any{"value": 4}),
				pc("Negative", `-3`, `-2`, map[string]any{"value": -3}),
				pc("Minus one", `-1`, `0`, map[string]any{"value": -1}),
				pc("Large", `999`, `1000`, map[string]any{"value": 999}),
			},
		},
		{
			id: 3, title: "Positive Check", topic: "Booleans · comparisons", difficulty: "Beginner",
			signature: "IsPositive(value int) bool",
			objective: "Report whether an integer is greater than zero.",
			input:     "One integer.", output: "true only when the value is greater than zero.",
			constraints: []string{"Zero is not positive.", "Return a bool."},
			hints:       []string{"Compare value with zero.", "A comparison already produces a bool.", "You can return the comparison directly."},
			pitfalls:    []string{"Treating zero as positive", "Returning the strings \"true\" or \"false\""},
			starter: `func IsPositive(value int) bool {
	// TODO: compare value with zero.
	return false
}
`, inputFields: intField, call: "IsPositive(current.Payload.Value)",
			cases: []practiceCase{
				pc("Zero", `0`, `false`, map[string]any{"value": 0}),
				pc("One", `1`, `true`, map[string]any{"value": 1}),
				pc("Negative", `-1`, `false`, map[string]any{"value": -1}),
				pc("Large positive", `500`, `true`, map[string]any{"value": 500}),
				pc("Large negative", `-500`, `false`, map[string]any{"value": -500}),
			},
		},
		{
			id: 4, title: "Larger Number", topic: "Conditionals · integers", difficulty: "Beginner+",
			signature: "MaxInt(left int, right int) int",
			objective: "Return the larger of two integers.",
			input:     "Two integers.", output: "The greater value; either value when they are equal.",
			constraints: []string{"Handle equal values.", "Negative values are valid."},
			hints:       []string{"Compare left and right.", "Return from either branch.", "Equality needs no special result."},
			pitfalls:    []string{"Always returning the left value", "Comparing absolute values"},
			starter: `func MaxInt(left int, right int) int {
	// TODO: return the larger value.
	return 0
}
`, inputFields: twoIntsFields, call: "MaxInt(current.Payload.Left, current.Payload.Right)",
			cases: []practiceCase{
				pc("Left larger", `4, 2`, `4`, map[string]any{"left": 4, "right": 2}),
				pc("Right larger", `2, 4`, `4`, map[string]any{"left": 2, "right": 4}),
				pc("Equal", `7, 7`, `7`, map[string]any{"left": 7, "right": 7}),
				pc("Negatives", `-8, -3`, `-3`, map[string]any{"left": -8, "right": -3}),
				pc("Across zero", `-1, 0`, `0`, map[string]any{"left": -1, "right": 0}),
			},
		},
		{
			id: 5, title: "Absolute Value", topic: "Conditionals · arithmetic", difficulty: "Beginner+",
			signature: "Abs(value int) int",
			objective: "Return the non-negative magnitude of an integer.",
			input:     "One integer.", output: "The value without a negative sign.",
			constraints: []string{"Zero remains zero.", "Do not import math for integer input."},
			hints:       []string{"Only negative values need changing.", "Negating a negative integer makes it positive.", "Use an if statement or a direct branch."},
			pitfalls:    []string{"Negating positive values", "Using floating-point math"},
			starter: `func Abs(value int) int {
	// TODO: remove a negative sign.
	return 0
}
`, inputFields: intField, call: "Abs(current.Payload.Value)",
			cases: []practiceCase{
				pc("Zero", `0`, `0`, map[string]any{"value": 0}),
				pc("Positive", `9`, `9`, map[string]any{"value": 9}),
				pc("Negative", `-9`, `9`, map[string]any{"value": -9}),
				pc("One", `1`, `1`, map[string]any{"value": 1}),
				pc("Minus one", `-1`, `1`, map[string]any{"value": -1}),
			},
		},
		{
			id: 6, title: "Number Clamp", topic: "Conditionals · bounds", difficulty: "Beginner+",
			signature: "Clamp(value int, minimum int, maximum int) int",
			objective: "Keep a number inside an inclusive range.",
			input:     "A value and valid minimum/maximum bounds.", output: "minimum below the range, maximum above it, otherwise value.",
			constraints: []string{"Bounds are inclusive.", "minimum is never greater than maximum."},
			hints:       []string{"Check the lower bound first.", "Then check the upper bound.", "Return value when neither bound is crossed."},
			pitfalls:    []string{"Swapping minimum and maximum", "Excluding values equal to a bound"},
			starter: `func Clamp(value int, minimum int, maximum int) int {
	// TODO: enforce both bounds.
	return value
}
`, inputFields: clampFields, call: "Clamp(current.Payload.Value, current.Payload.Minimum, current.Payload.Maximum)",
			cases: []practiceCase{
				pc("Inside", `5, 0, 10`, `5`, map[string]any{"value": 5, "minimum": 0, "maximum": 10}),
				pc("Below", `-2, 0, 10`, `0`, map[string]any{"value": -2, "minimum": 0, "maximum": 10}),
				pc("Above", `12, 0, 10`, `10`, map[string]any{"value": 12, "minimum": 0, "maximum": 10}),
				pc("At minimum", `0, 0, 10`, `0`, map[string]any{"value": 0, "minimum": 0, "maximum": 10}),
				pc("At maximum", `10, 0, 10`, `10`, map[string]any{"value": 10, "minimum": 0, "maximum": 10}),
			},
		},
		{
			id: 7, title: "Byte Counter", topic: "Strings · len", difficulty: "Beginner+",
			signature: "ByteCount(value string) int",
			objective: "Count the bytes stored in a string.",
			input:     "One UTF-8 string.", output: "Its byte length.",
			constraints: []string{"Count bytes, not visible characters.", "Use Go's len built-in."},
			hints:       []string{"len works on strings.", "UTF-8 characters may occupy multiple bytes.", "No loop is required."},
			pitfalls:    []string{"Counting only non-space bytes", "Assuming every character is one byte"},
			builtins:    []string{"len"},
			starter: `func ByteCount(value string) int {
	// TODO: return the string length in bytes.
	return 0
}
`, inputFields: stringField, call: "ByteCount(current.Payload.Value)",
			cases: []practiceCase{
				pc("Empty", `""`, `0`, map[string]any{"value": ""}),
				pc("ASCII", `"gopher"`, `6`, map[string]any{"value": "gopher"}),
				pc("Spaces", `"a b"`, `3`, map[string]any{"value": "a b"}),
				pc("Greek", `"γ"`, `2`, map[string]any{"value": "γ"}),
				pc("Emoji", `"🙂"`, `4`, map[string]any{"value": "🙂"}),
			},
		},
		{
			id: 8, title: "First Byte", topic: "Strings · indexing", difficulty: "Beginner+",
			signature: "FirstByte(value string) string",
			objective: "Return the first ASCII byte as a one-character string.",
			input:     "An ASCII string, possibly empty.", output: "The first character, or an empty string.",
			constraints: []string{"Handle empty input before indexing.", "Inputs in this exercise are ASCII."},
			hints:       []string{"Check len before value[0].", "String indexing returns a byte.", "Convert the byte with string(value[0])."},
			pitfalls:    []string{"Indexing an empty string", "Returning a byte instead of a string"},
			builtins:    []string{"len"},
			starter: `func FirstByte(value string) string {
	// TODO: handle empty input, then return the first byte.
	return ""
}
`, inputFields: stringField, call: "FirstByte(current.Payload.Value)",
			cases: []practiceCase{
				pc("Empty", `""`, `""`, map[string]any{"value": ""}),
				pc("One", `"a"`, `"a"`, map[string]any{"value": "a"}),
				pc("Word", `"gopher"`, `"g"`, map[string]any{"value": "gopher"}),
				pc("Space", `" hello"`, `" "`, map[string]any{"value": " hello"}),
				pc("Digit", `"9 lives"`, `"9"`, map[string]any{"value": "9 lives"}),
			},
		},
		{
			id: 9, title: "Repeat Text", topic: "Loops · strings", difficulty: "Beginner+",
			signature: "Repeat(value string, count int) string",
			objective: "Build a string by repeating another string.",
			input:     "A string and a non-negative repetition count.", output: "value repeated count times.",
			constraints: []string{"Zero repetitions return an empty string.", "Use a loop."},
			hints:       []string{"Start with an empty result.", "Run a loop count times.", "Add value to result on every iteration."},
			pitfalls:    []string{"Repeating count+1 times", "Adding separators that were not requested"},
			starter: `func Repeat(value string, count int) string {
	// TODO: append value count times.
	return ""
}
`, inputFields: stringCountFields, call: "Repeat(current.Payload.Value, current.Payload.Count)",
			cases: []practiceCase{
				pc("Zero", `"go", 0`, `""`, map[string]any{"value": "go", "count": 0}),
				pc("Once", `"go", 1`, `"go"`, map[string]any{"value": "go", "count": 1}),
				pc("Three", `"go", 3`, `"gogogo"`, map[string]any{"value": "go", "count": 3}),
				pc("Empty text", `"", 5`, `""`, map[string]any{"value": "", "count": 5}),
				pc("Symbol", `"!", 4`, `"!!!!"`, map[string]any{"value": "!", "count": 4}),
			},
		},
		{
			id: 10, title: "Slice Sum", topic: "Slices · loops", difficulty: "Foundation",
			signature: "Sum(values []int) int",
			objective: "Add every integer in a slice.",
			input:     "A slice of integers.", output: "Their sum; zero for an empty slice.",
			constraints: []string{"Do not modify the input.", "Negative values are valid."},
			hints:       []string{"Create a total starting at zero.", "Range over values.", "Add each value to total."},
			pitfalls:    []string{"Starting the total at one", "Returning after the first item"},
			starter: `func Sum(values []int) int {
	// TODO: add every value.
	return 0
}
`, inputFields: intsField, call: "Sum(current.Payload.Values)",
			cases: []practiceCase{
				pc("Empty", `[]int{}`, `0`, map[string]any{"values": []int{}}),
				pc("One", `[]int{7}`, `7`, map[string]any{"values": []int{7}}),
				pc("Several", `[]int{1,2,3}`, `6`, map[string]any{"values": []int{1, 2, 3}}),
				pc("Mixed", `[]int{-5,2,3}`, `0`, map[string]any{"values": []int{-5, 2, 3}}),
				pc("Negatives", `[]int{-2,-4}`, `-6`, map[string]any{"values": []int{-2, -4}}),
			},
		},
		{
			id: 11, title: "Match Counter", topic: "Slices · comparisons", difficulty: "Foundation",
			signature: "CountValue(values []int, target int) int",
			objective: "Count how often a target appears.",
			input:     "A slice of integers and a target.", output: "The number of equal elements.",
			constraints: []string{"Count every occurrence.", "Return zero when there are no matches."},
			hints:       []string{"Start a counter at zero.", "Compare every value with target.", "Increment only on equality."},
			pitfalls:    []string{"Stopping after the first match", "Returning an index"},
			starter: `func CountValue(values []int, target int) int {
	// TODO: count target.
	return 0
}
`, inputFields: intsTargetFields, call: "CountValue(current.Payload.Values, current.Payload.Target)",
			cases: []practiceCase{
				pc("Empty", `[]int{}, 2`, `0`, map[string]any{"values": []int{}, "target": 2}),
				pc("None", `[]int{1,3}, 2`, `0`, map[string]any{"values": []int{1, 3}, "target": 2}),
				pc("One", `[]int{1,2,3}, 2`, `1`, map[string]any{"values": []int{1, 2, 3}, "target": 2}),
				pc("Repeated", `[]int{2,2,1,2}, 2`, `3`, map[string]any{"values": []int{2, 2, 1, 2}, "target": 2}),
				pc("Negative", `[]int{-1,0,-1}, -1`, `2`, map[string]any{"values": []int{-1, 0, -1}, "target": -1}),
			},
		},
		{
			id: 12, title: "Contains", topic: "Slices · early return", difficulty: "Foundation",
			signature: "Contains(values []int, target int) bool",
			objective: "Report whether a slice contains a target.",
			input:     "A slice of integers and a target.", output: "true when at least one element equals target.",
			constraints: []string{"An empty slice returns false.", "You may return as soon as a match is found."},
			hints:       []string{"Inspect every value until a match.", "Return true inside the matching branch.", "Return false after the loop."},
			pitfalls:    []string{"Returning false after checking only one item", "Comparing indexes with target"},
			starter: `func Contains(values []int, target int) bool {
	// TODO: search for target.
	return false
}
`, inputFields: intsTargetFields, call: "Contains(current.Payload.Values, current.Payload.Target)",
			cases: []practiceCase{
				pc("Empty", `[]int{}, 1`, `false`, map[string]any{"values": []int{}, "target": 1}),
				pc("First", `[]int{1,2,3}, 1`, `true`, map[string]any{"values": []int{1, 2, 3}, "target": 1}),
				pc("Last", `[]int{1,2,3}, 3`, `true`, map[string]any{"values": []int{1, 2, 3}, "target": 3}),
				pc("Missing", `[]int{1,2,3}, 4`, `false`, map[string]any{"values": []int{1, 2, 3}, "target": 4}),
				pc("Negative", `[]int{0,-2}, -2`, `true`, map[string]any{"values": []int{0, -2}, "target": -2}),
			},
		},
		{
			id: 13, title: "Reverse Text", topic: "Strings · indexing", difficulty: "Foundation+",
			signature: "Reverse(value string) string",
			objective: "Reverse an ASCII string.",
			input:     "An ASCII string.", output: "Its bytes in reverse order.",
			constraints: []string{"Empty input returns empty output.", "Inputs in this exercise are ASCII."},
			hints:       []string{"Start at len(value)-1.", "Move the index toward zero.", "Build a new result string."},
			pitfalls:    []string{"Starting at len(value)", "Dropping the byte at index zero"},
			builtins:    []string{"len"},
			starter: `func Reverse(value string) string {
	// TODO: walk backward through value.
	return ""
}
`, inputFields: stringField, call: "Reverse(current.Payload.Value)",
			cases: []practiceCase{
				pc("Empty", `""`, `""`, map[string]any{"value": ""}),
				pc("One", `"a"`, `"a"`, map[string]any{"value": "a"}),
				pc("Word", `"gopher"`, `"rehpog"`, map[string]any{"value": "gopher"}),
				pc("Spaces", `"a b"`, `"b a"`, map[string]any{"value": "a b"}),
				pc("Palindrome", `"level"`, `"level"`, map[string]any{"value": "level"}),
			},
		},
		{
			id: 14, title: "Palindrome", topic: "Strings · two pointers", difficulty: "Foundation+",
			signature: "IsPalindrome(value string) bool",
			objective: "Check whether an ASCII string reads the same backward.",
			input:     "An ASCII string with exact case and spacing.", output: "true when mirrored bytes match.",
			constraints: []string{"Comparison is case-sensitive.", "Empty and one-byte strings are palindromes."},
			hints:       []string{"Compare the first and last bytes.", "Move both indexes inward.", "Return false on the first mismatch."},
			pitfalls:    []string{"Ignoring case when not requested", "Skipping the final pair"},
			builtins:    []string{"len"},
			starter: `func IsPalindrome(value string) bool {
	// TODO: compare mirrored bytes.
	return false
}
`, inputFields: stringField, call: "IsPalindrome(current.Payload.Value)",
			cases: []practiceCase{
				pc("Empty", `""`, `true`, map[string]any{"value": ""}),
				pc("One", `"x"`, `true`, map[string]any{"value": "x"}),
				pc("Odd", `"level"`, `true`, map[string]any{"value": "level"}),
				pc("Even", `"noon"`, `true`, map[string]any{"value": "noon"}),
				pc("Mismatch", `"gopher"`, `false`, map[string]any{"value": "gopher"}),
			},
		},
		{
			id: 15, title: "Even Filter", topic: "Slices · append", difficulty: "Foundation+",
			signature: "FilterEven(values []int) []int",
			objective: "Keep only even integers in their original order.",
			input:     "A slice of integers.", output: "A non-nil slice containing each even value.",
			constraints: []string{"Preserve order and duplicates.", "Return []int{} when no values match."},
			hints:       []string{"Start with result := []int{}.", "An even value has value%2 == 0.", "Use result = append(result, value)."},
			pitfalls:    []string{"Using result.append", "Returning nil for no matches"},
			builtins:    []string{"append"},
			starter: `func FilterEven(values []int) []int {
	result := []int{}
	// TODO: append even values.
	return result
}
`, inputFields: intsField, call: "FilterEven(current.Payload.Values)",
			cases: []practiceCase{
				pc("Empty", `[]int{}`, `[]`, map[string]any{"values": []int{}}),
				pc("None", `[]int{1,3}`, `[]`, map[string]any{"values": []int{1, 3}}),
				pc("All", `[]int{2,4}`, `[2,4]`, map[string]any{"values": []int{2, 4}}),
				pc("Mixed", `[]int{1,2,3,4}`, `[2,4]`, map[string]any{"values": []int{1, 2, 3, 4}}),
				pc("Negative", `[]int{-2,-1,0}`, `[-2,0]`, map[string]any{"values": []int{-2, -1, 0}}),
			},
		},
		{
			id: 16, title: "Minimum", topic: "Slices · state", difficulty: "Foundation+",
			signature: "Minimum(values []int) int",
			objective: "Find the smallest integer in a slice.",
			input:     "A slice of integers.", output: "The smallest value, or zero for an empty slice.",
			constraints: []string{"Handle negative values.", "Do not sort the input."},
			hints:       []string{"Return zero when len(values)==0.", "Start the minimum at values[0].", "Replace it whenever a smaller value appears."},
			pitfalls:    []string{"Starting the minimum at zero", "Returning the final value instead of the minimum"},
			builtins:    []string{"len"},
			starter: `func Minimum(values []int) int {
	// TODO: handle empty input and track the minimum.
	return 0
}
`, inputFields: intsField, call: "Minimum(current.Payload.Values)",
			cases: []practiceCase{
				pc("Empty", `[]int{}`, `0`, map[string]any{"values": []int{}}),
				pc("One", `[]int{8}`, `8`, map[string]any{"values": []int{8}}),
				pc("First", `[]int{1,5,9}`, `1`, map[string]any{"values": []int{1, 5, 9}}),
				pc("Last", `[]int{9,5,1}`, `1`, map[string]any{"values": []int{9, 5, 1}}),
				pc("Negative", `[]int{4,-7,2}`, `-7`, map[string]any{"values": []int{4, -7, 2}}),
			},
		},
		{
			id: 17, title: "Unique Strings", topic: "Maps · stable order", difficulty: "Intermediate-",
			signature: "UniqueStrings(values []string) []string",
			objective: "Remove duplicate strings while keeping first-seen order.",
			input:     "A slice of case-sensitive strings.", output: "A non-nil slice containing each distinct value once.",
			constraints: []string{"Preserve first appearance order.", "Empty strings are ordinary values."},
			hints:       []string{"Use a map[string]bool to remember seen values.", "Append only when a value is new.", "Iterate over the input, not the map."},
			pitfalls:    []string{"Returning map iteration order", "Removing empty strings"},
			builtins:    []string{"append", "make"},
			starter: `func UniqueStrings(values []string) []string {
	result := []string{}
	// TODO: remember values already appended.
	return result
}
`, inputFields: stringsField, call: "UniqueStrings(current.Payload.Values)",
			cases: []practiceCase{
				pc("Empty", `[]string{}`, `[]`, map[string]any{"values": []string{}}),
				pc("Distinct", `[]string{"a","b"}`, `["a","b"]`, map[string]any{"values": []string{"a", "b"}}),
				pc("Repeated", `[]string{"a","a"}`, `["a"]`, map[string]any{"values": []string{"a", "a"}}),
				pc("Stable", `[]string{"b","a","b"}`, `["b","a"]`, map[string]any{"values": []string{"b", "a", "b"}}),
				pc("Empty value", `[]string{"","a",""}`, `["","a"]`, map[string]any{"values": []string{"", "a", ""}}),
			},
		},
		{
			id: 18, title: "String Lengths", topic: "Slices · transformation", difficulty: "Intermediate-",
			signature: "Lengths(values []string) []int",
			objective: "Convert strings into their byte lengths.",
			input:     "A slice of UTF-8 strings.", output: "A non-nil []int with one byte length per input.",
			constraints: []string{"Preserve input order.", "The output length must equal the input length."},
			hints:       []string{"Start with []int{}.", "Range over each string.", "Append len(value)."},
			pitfalls:    []string{"Skipping empty strings", "Returning one combined total"},
			builtins:    []string{"append", "len"},
			starter: `func Lengths(values []string) []int {
	result := []int{}
	// TODO: append each byte length.
	return result
}
`, inputFields: stringsField, call: "Lengths(current.Payload.Values)",
			cases: []practiceCase{
				pc("Empty slice", `[]string{}`, `[]`, map[string]any{"values": []string{}}),
				pc("One", `[]string{"go"}`, `[2]`, map[string]any{"values": []string{"go"}}),
				pc("Mixed", `[]string{"","a","abc"}`, `[0,1,3]`, map[string]any{"values": []string{"", "a", "abc"}}),
				pc("Spaces", `[]string{"a b"," "}`, `[3,1]`, map[string]any{"values": []string{"a b", " "}}),
				pc("UTF-8 bytes", `[]string{"γ","🙂"}`, `[2,4]`, map[string]any{"values": []string{"γ", "🙂"}}),
			},
		},
		{
			id: 19, title: "Frequency Map", topic: "Maps · counting", difficulty: "Intermediate",
			signature: "Frequencies(values []string) map[string]int",
			objective: "Count how often every string appears.",
			input:     "A slice of case-sensitive strings.", output: "A non-nil map from each string to its count.",
			constraints: []string{"Count empty strings.", "Do not normalize case."},
			hints:       []string{"Create a map[string]int.", "Range over values.", "Increment result[value]."},
			pitfalls:    []string{"Returning nil for empty input", "Resetting a count to one on every match"},
			builtins:    []string{"make"},
			starter: `func Frequencies(values []string) map[string]int {
	result := map[string]int{}
	// TODO: count each value.
	return result
}
`, inputFields: stringsField, call: "Frequencies(current.Payload.Values)",
			cases: []practiceCase{
				pc("Empty", `[]string{}`, `{}`, map[string]any{"values": []string{}}),
				pc("One", `[]string{"a"}`, `{"a":1}`, map[string]any{"values": []string{"a"}}),
				pc("Repeated", `[]string{"a","a"}`, `{"a":2}`, map[string]any{"values": []string{"a", "a"}}),
				pc("Several", `[]string{"b","a","b"}`, `{"a":1,"b":2}`, map[string]any{"values": []string{"b", "a", "b"}}),
				pc("Case", `[]string{"Go","go"}`, `{"Go":1,"go":1}`, map[string]any{"values": []string{"Go", "go"}}),
			},
		},
		{
			id: 20, title: "Rotate Left", topic: "Slices · modular indexing", difficulty: "Intermediate",
			signature: "RotateLeft(values []int, steps int) []int",
			objective: "Rotate a slice left by a non-negative number of steps.",
			input:     "A slice and a non-negative step count.", output: "A new non-nil slice in rotated order.",
			constraints: []string{"Do not modify the input.", "Steps may be greater than the slice length."},
			hints:       []string{"Return []int{} for empty input.", "Reduce steps with steps % len(values).", "Append the tail before the head."},
			pitfalls:    []string{"Dividing by zero for empty input", "Dropping values when steps exceed the length"},
			builtins:    []string{"append", "len"},
			starter: `func RotateLeft(values []int, steps int) []int {
	result := []int{}
	// TODO: normalize steps and append both sections.
	return result
}
`, inputFields: intsStepsFields, call: "RotateLeft(current.Payload.Values, current.Payload.Steps)",
			cases: []practiceCase{
				pc("Empty", `[]int{}, 3`, `[]`, map[string]any{"values": []int{}, "steps": 3}),
				pc("Zero", `[]int{1,2,3}, 0`, `[1,2,3]`, map[string]any{"values": []int{1, 2, 3}, "steps": 0}),
				pc("One", `[]int{1,2,3}, 1`, `[2,3,1]`, map[string]any{"values": []int{1, 2, 3}, "steps": 1}),
				pc("Length", `[]int{1,2,3}, 3`, `[1,2,3]`, map[string]any{"values": []int{1, 2, 3}, "steps": 3}),
				pc("Beyond", `[]int{1,2,3}, 5`, `[3,1,2]`, map[string]any{"values": []int{1, 2, 3}, "steps": 5}),
			},
		},
		{
			id: 21, title: "Balanced Brackets", topic: "Stacks · parsing", difficulty: "Intermediate",
			signature: "BalancedBrackets(value string) bool",
			objective: "Validate nested (), [], and {} bracket pairs.",
			input:     "An ASCII string containing brackets and other characters.", output: "true when every bracket is correctly matched and nested.",
			constraints: []string{"Ignore non-bracket characters.", "Closing brackets must match the latest opening bracket."},
			hints:       []string{"Use a byte slice as a stack.", "Push opening brackets.", "On a closing bracket, check and remove the stack's final item."},
			pitfalls:    []string{"Checking only bracket counts", "Accepting a non-empty stack at the end"},
			builtins:    []string{"append", "len"},
			starter: `func BalancedBrackets(value string) bool {
	stack := []byte{}
	// TODO: push openings and match closings.
	return len(stack) == 0
}
`, inputFields: stringField, call: "BalancedBrackets(current.Payload.Value)",
			cases: []practiceCase{
				pc("Empty", `""`, `true`, map[string]any{"value": ""}),
				pc("Simple", `"()"`, `true`, map[string]any{"value": "()"}),
				pc("Nested", `"{[()]}"`, `true`, map[string]any{"value": "{[()]}"}),
				pc("Wrong order", `"([)]"`, `false`, map[string]any{"value": "([)]"}),
				pc("Unclosed", `"(text"`, `false`, map[string]any{"value": "(text"}),
			},
		},
	}

	levels := make([]Level, 0, len(specs))
	for _, spec := range specs {
		levels = append(levels, makePracticeLevel(spec))
	}
	return levels
}

func makePracticeLevel(spec practiceSpec) Level {
	instructions := baseInstructions(
		spec.objective,
		"func "+spec.signature,
		spec.input,
		spec.output,
		"Complete the function body. Keep the signature unchanged.",
	)
	instructions.Constraints = append([]string(nil), spec.constraints...)
	instructions.Hints = append([]string(nil), spec.hints...)
	instructions.CommonPitfalls = append([]string(nil), spec.pitfalls...)
	instructions.Examples = make([]Example, 0, 4)
	for index := 0; index < 4; index++ {
		current := spec.cases[index%len(spec.cases)]
		instructions.Examples = append(instructions.Examples, Example{
			Input:  functionName(spec.signature) + "(" + current.input + ")",
			Output: current.expected,
		})
	}
	packages := spec.packages
	if len(packages) == 0 {
		packages = []string{"fmt"}
	}
	setSourcePolicy(&instructions, spec.builtins, packages)

	tests := make([]VisibleTest, 0, len(spec.cases))
	for index, current := range spec.cases {
		tests = append(tests, test(
			fmt.Sprintf("l%d-%02d", spec.id, index+1),
			current.name,
			"Checks "+current.name+".",
			current.input,
			current.expected,
			current.payload,
		))
	}
	inputFields, call := spec.inputFields, spec.call
	return Level{
		ID: spec.id, Title: spec.title, Topic: spec.topic,
		Difficulty: spec.difficulty, Signature: spec.signature,
		StarterCode: spec.starter, Instructions: instructions, Tests: tests,
		build: func(selected []VisibleTest) string {
			return buildPracticeHarness(inputFields, call, selected)
		},
	}
}

func buildPracticeHarness(inputFields, call string, tests []VisibleTest) string {
	declarations := `type inputCase struct {
	` + inputFields + `
}
type wireCase struct {
	ID string ` + "`json:\"id\"`" + `
	Payload inputCase ` + "`json:\"payload\"`" + `
}`
	loop := `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value := ` + call + `
			encoded, _ := json.Marshal(value)
			return string(encoded)
		})
	}`
	return harness(commonImports(), declarations, loop, tests)
}

func pc(name, input, expected string, payload any) practiceCase {
	return practiceCase{name: name, input: input, expected: expected, payload: payload}
}

func functionName(signature string) string {
	for index, char := range signature {
		if char == '(' {
			return signature[:index]
		}
	}
	return signature
}

const (
	stringField       = "Value string `json:\"value\"`"
	intField          = "Value int `json:\"value\"`"
	twoIntsFields     = "Left int `json:\"left\"`\n\tRight int `json:\"right\"`"
	clampFields       = "Value int `json:\"value\"`\n\tMinimum int `json:\"minimum\"`\n\tMaximum int `json:\"maximum\"`"
	stringCountFields = "Value string `json:\"value\"`\n\tCount int `json:\"count\"`"
	intsField         = "Values []int `json:\"values\"`"
	intsTargetFields  = "Values []int `json:\"values\"`\n\tTarget int `json:\"target\"`"
	stringsField      = "Values []string `json:\"values\"`"
	intsStepsFields   = "Values []int `json:\"values\"`\n\tSteps int `json:\"steps\"`"
)
