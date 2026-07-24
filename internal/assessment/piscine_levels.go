package assessment

import (
	"fmt"
	"strings"
)

// piscineLevels contains the unique exercises found in kinoz01/zone01-Piscine.
// Duplicate solution variants and exercises already represented elsewhere in
// the curriculum are intentionally omitted.
func piscineLevels() []Level {
	specs := make([]practiceSpec, 0, 32)
specs = append(specs, piscineQuestTwoSpecs()...)
specs = append(specs, piscineQuestThreeAndFourSpecs()...)

	levels := make([]Level, 0, len(specs))
	for _, spec := range specs {
		levels = append(levels, makePracticeLevel(spec))
	}
	return levels
}

func piscineQuestTwoSpecs() []practiceSpec {
	return []practiceSpec{
		piscineSpec(1001, "Only Z", "Piscine Quest 2", "OnlyZ() string",
			"Return the single lowercase letter z.", "No input.", "Exactly \"z\".",
			nil, noFields, "OnlyZ()", `""`,
			[]practiceCase{pc("Exact output", ``, jsonString("z"), map[string]any{})}),
		piscineSpec(1002, "Only One", "Piscine Quest 2", "OnlyOne() string",
			"Return the single character 1.", "No input.", "Exactly \"1\".",
			nil, noFields, "OnlyOne()", `""`,
			[]practiceCase{pc("Exact output", ``, jsonString("1"), map[string]any{})}),
		piscineSpec(1003, "Only B", "Piscine Quest 2", "OnlyB() string",
			"Return the single uppercase letter B.", "No input.", "Exactly \"B\".",
			nil, noFields, "OnlyB()", `""`,
			[]practiceCase{pc("Exact output", ``, jsonString("B"), map[string]any{})}),
		piscineSpec(1004, "Only F", "Piscine Quest 2", "OnlyF() string",
			"Return the single lowercase letter f.", "No input.", "Exactly \"f\".",
			nil, noFields, "OnlyF()", `""`,
			[]practiceCase{pc("Exact output", ``, jsonString("f"), map[string]any{})}),
		piscineSpec(1005, "Hello World", "Piscine Quest 2", "HelloWorld() string",
			"Produce the traditional first-program greeting.", "No input.", "\"Hello World!\\n\".",
			nil, noFields, "HelloWorld()", `""`,
			[]practiceCase{pc("Greeting", ``, jsonString("Hello World!\n"), map[string]any{})}),
		piscineSpec(1006, "Display A", "Piscine Quest 2", "DisplayA(value string) string",
			"Report whether a string contains the lowercase letter a.", "One string.", "\"a\\n\" when a is present; otherwise just newline.",
			nil, stringField, "DisplayA(current.Payload.Value)", `"\n"`,
			stringCases([]string{"", "xyz", "cat", "A", "banana"}, []string{"\n", "\n", "a\n", "\n", "a\n"})),
		piscineSpec(1007, "Display Z", "Piscine Quest 2", "DisplayZ(value string) string",
			"Report whether a string contains the lowercase letter z.", "One string.", "\"z\\n\" when z is present; otherwise just newline.",
			nil, stringField, "DisplayZ(current.Payload.Value)", `"\n"`,
			stringCases([]string{"", "abc", "fizz", "Z", "zebra"}, []string{"\n", "\n", "z\n", "\n", "z\n"})),
		piscineSpec(1008, "Print Alphabet", "Piscine Quest 2", "PrintAlphabet() string",
			"Generate the lowercase Latin alphabet.", "No input.", "abcdefghijklmnopqrstuvwxyz followed by newline.",
			nil, noFields, "PrintAlphabet()", `""`,
			[]practiceCase{pc("Alphabet", ``, jsonString("abcdefghijklmnopqrstuvwxyz\n"), map[string]any{})}),
		piscineSpec(1009, "Print Reverse Alphabet", "Piscine Quest 2", "PrintReverseAlphabet() string",
			"Generate the lowercase Latin alphabet in reverse.", "No input.", "zyxwvutsrqponmlkjihgfedcba followed by newline.",
			nil, noFields, "PrintReverseAlphabet()", `""`,
			[]practiceCase{pc("Reverse alphabet", ``, jsonString("zyxwvutsrqponmlkjihgfedcba\n"), map[string]any{})}),
		piscineSpec(1010, "Print Digits", "Piscine Quest 2", "PrintDigits() string",
			"Generate decimal digits in ascending order.", "No input.", "0123456789 followed by newline.",
			nil, noFields, "PrintDigits()", `""`,
			[]practiceCase{pc("Digits", ``, jsonString("0123456789\n"), map[string]any{})}),
		piscineSpec(1011, "Countdown", "Piscine Quest 2", "Countdown() string",
			"Generate decimal digits in descending order.", "No input.", "9876543210 followed by newline.",
			nil, noFields, "Countdown()", `""`,
			[]practiceCase{pc("Countdown", ``, jsonString("9876543210\n"), map[string]any{})}),
		piscineSpec(1012, "Is Negative", "Piscine Quest 2", "IsNegative(value int) string",
			"Classify an integer by its sign.", "One integer.", "\"T\\n\" for negative values and \"F\\n\" otherwise.",
			nil, intField, "IsNegative(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Negative", -1, "T\n"),
				intOutputCase("Zero", 0, "F\n"),
				intOutputCase("Positive", 1, "F\n"),
				intOutputCase("Large negative", -999, "T\n"),
			}),
		piscineSpec(1013, "Print Number", "Piscine Quest 2", "PrintNbr(value int) string",
			"Convert any int to decimal output without strconv.", "One integer.", "Its exact base-10 representation.",
			nil, intField, "PrintNbr(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Zero", 0, "0"),
				intOutputCase("Positive", 42, "42"),
				intOutputCase("Negative", -42, "-42"),
				intOutputCase("Large", 987654321, "987654321"),
			}),
		piscineSpec(1014, "Print Combination", "Piscine Quest 2", "PrintComb() string",
			"Generate every strictly ascending combination of three distinct digits.", "No input.", "Combinations from 012 through 789 separated by comma and space, then newline.",
			nil, noFields, "PrintComb()", `""`,
			[]practiceCase{pc("All combinations", ``, jsonString(ascendingTriplesExpected()), map[string]any{})}),
		piscineSpec(1015, "Print Combination Two", "Piscine Quest 2", "PrintComb2() string",
			"Generate every ordered pair of distinct two-digit numbers.", "No input.", "Pairs from 00 01 through 98 99 separated by comma and space, then newline.",
			nil, noFields, "PrintComb2()", `""`,
			[]practiceCase{pc("All pairs", ``, jsonString(ascendingPairsExpected()), map[string]any{})}),
		piscineSpec(1016, "Descending Combination", "Piscine Quest 2", "DescendComb() string",
			"Generate every pair of distinct two-digit numbers in descending order.", "No input.", "Pairs from 99 98 down through 01 00 separated by comma and space, then newline.",
			nil, noFields, "DescendComb()", `""`,
			[]practiceCase{pc("All descending pairs", ``, jsonString(descendingPairsExpected()), map[string]any{})}),
		piscineSpec(1017, "Alternating Alphabet", "Piscine Quest 2", "DisplayAlphaM() string",
			"Generate the alphabet with odd positions lowercase and even positions uppercase.", "No input.", "aBcDeF...yZ followed by newline.",
			nil, noFields, "DisplayAlphaM()", `""`,
			[]practiceCase{pc("Alternating alphabet", ``, jsonString("aBcDeFgHiJkLmNoPqRsTuVwXyZ\n"), map[string]any{})}),
		piscineSpec(1018, "Reverse Alternating Alphabet", "Piscine Quest 2", "DisplayAlphaReverseM() string",
			"Generate the reverse alphabet with alternating case.", "No input.", "zYxW...bA followed by newline.",
			nil, noFields, "DisplayAlphaReverseM()", `""`,
			[]practiceCase{pc("Reverse alternating alphabet", ``, jsonString("zYxWvUtSrQpOnMlKjIhGfEdCbA\n"), map[string]any{})}),
	}
}

func piscineQuestThreeAndFourSpecs() []practiceSpec {
	return []practiceSpec{
		piscineSpec(1019, "Basic Atoi", "Piscine Quest 3", "BasicAtoi(value string) int",
			"Convert a digit-only string to an integer.", "A non-empty string of decimal digits.", "Its integer value.",
			nil, stringField, "BasicAtoi(current.Payload.Value)", `0`,
			intStringCases([]string{"0", "7", "12345", "00042"}, []int{0, 7, 12345, 42})),
		piscineSpec(1020, "Validated Basic Atoi", "Piscine Quest 3", "BasicAtoi2(value string) int",
			"Convert decimal digits and reject any other character.", "One string.", "Its value, or 0 when any character is not a digit.",
			nil, stringField, "BasicAtoi2(current.Payload.Value)", `0`,
			intStringCases([]string{"", "12345", "00042", "12 3", "hello"}, []int{0, 12345, 42, 0, 0})),
		piscineSpec(1021, "ASCII To Integer", "Piscine Quest 3", "Atoi(value string) int",
			"Parse an optional leading sign and decimal digits.", "One string.", "The signed integer, or 0 for invalid syntax.",
			[]string{"len"}, stringField, "Atoi(current.Payload.Value)", `0`,
			intStringCases([]string{"0", "12345", "-1234", "+1234", "++1", "12 3"}, []int{0, 12345, -1234, 1234, 0, 0})),
		piscineSpec(1022, "Point One", "Piscine Quest 3", "PointOne(value *int) int",
			"Set an integer through a pointer.", "A pointer to an integer.", "Set the pointed value to 1 and return it.",
			nil, intField, "func() int { value := current.Payload.Value; return PointOne(&value) }()", `0`,
			[]practiceCase{
				intCase("Zero", 0, 1),
				intCase("Positive", 7, 1),
				intCase("Negative", -7, 1),
			}),
		piscineSpec(1023, "Swap", "Piscine Quest 3", "Swap(left *int, right *int) []int",
			"Swap two integers through pointers.", "Pointers to two integers.", "A two-item slice containing the swapped values.",
			nil, twoIntsFields, "func() []int { left, right := current.Payload.Left, current.Payload.Right; return Swap(&left, &right) }()", `nil`,
			[]practiceCase{
				twoIntsSliceCase("Different", 1, 2, []int{2, 1}),
				twoIntsSliceCase("Equal", 5, 5, []int{5, 5}),
				twoIntsSliceCase("Signs", -1, 2, []int{2, -1}),
			}),
		piscineSpec(1024, "Division And Modulo", "Piscine Quest 3", "DivMod(left int, right int) []int",
			"Calculate integer quotient and remainder.", "Two integers with a non-zero divisor.", "A two-item slice containing quotient and remainder.",
			nil, twoIntsFields, "DivMod(current.Payload.Left, current.Payload.Right)", `nil`,
			[]practiceCase{
				twoIntsSliceCase("Exact", 10, 2, []int{5, 0}),
				twoIntsSliceCase("Remainder", 13, 5, []int{2, 3}),
				twoIntsSliceCase("Negative", -13, 5, []int{-2, -3}),
			}),
		piscineSpec(1025, "Ultimate Division And Modulo", "Piscine Quest 3", "UltimateDivMod(left int, right int) []int",
			"Replace two values with their quotient and remainder.", "Two integers with a non-zero divisor.", "A two-item slice containing the resulting values.",
			nil, twoIntsFields, "UltimateDivMod(current.Payload.Left, current.Payload.Right)", `nil`,
			[]practiceCase{
				twoIntsSliceCase("Exact", 10, 2, []int{5, 0}),
				twoIntsSliceCase("Remainder", 13, 5, []int{2, 3}),
				twoIntsSliceCase("Negative", -13, 5, []int{-2, -3}),
			}),
		piscineSpec(1026, "Ultimate Point One", "Piscine Quest 3", "UltimatePointOne(value *int) int",
			"Reach an integer through three pointer levels and set it to 1.", "A pointer to an integer; create the extra pointer levels inside the function.", "Set and return 1.",
			nil, intField, "func() int { value := current.Payload.Value; return UltimatePointOne(&value) }()", `0`,
			[]practiceCase{
				intCase("Zero", 0, 1),
				intCase("Positive", 42, 1),
				intCase("Negative", -42, 1),
			}),
		piscineSpec(1027, "Sort Integer Table", "Piscine Quest 3", "SortIntegerTable(values []int) []int",
			"Sort integers in ascending order without using sort.", "One integer slice.", "A sorted slice containing the same values.",
			[]string{"len"}, intsField, "SortIntegerTable(current.Payload.Values)", `nil`,
			[]practiceCase{
				intSliceCase("Empty", []int{}, []int{}),
				intSliceCase("Sorted", []int{1, 2, 3}, []int{1, 2, 3}),
				intSliceCase("Reverse", []int{3, 2, 1}, []int{1, 2, 3}),
				intSliceCase("Duplicates", []int{2, 1, 2}, []int{1, 2, 2}),
			}),
		piscineSpec(1028, "Rune Length", "Piscine Quest 3", "StrLen(value string) int",
			"Count runes rather than UTF-8 bytes.", "One UTF-8 string.", "The number of Unicode code points.",
			nil, stringField, "StrLen(current.Payload.Value)", `0`,
			intStringCases([]string{"", "go", "?", "??", "caf‚"}, []int{0, 2, 1, 1, 4})),
		piscineSpec(1029, "Prime Check", "Piscine Quest 4", "IsPrime(value int) bool",
			"Determine whether an integer is prime.", "One integer.", "true only for prime values greater than 1.",
			nil, intField, "IsPrime(current.Payload.Value)", `false`,
			[]practiceCase{
				intBoolCase("Negative", -3, false),
				intBoolCase("One", 1, false),
				intBoolCase("Two", 2, true),
				intBoolCase("Composite", 9, false),
				intBoolCase("Prime", 97, true),
			}),
		piscineSpec(1030, "Find Next Prime", "Piscine Quest 4", "FindNextPrime(value int) int",
			"Find the smallest prime greater than or equal to an integer.", "One integer.", "The nearest prime at or above the input.",
			nil, intField, "FindNextPrime(current.Payload.Value)", `0`,
			[]practiceCase{
				intCase("Below primes", -3, 2),
				intCase("Prime", 5, 5),
				intCase("Composite", 8, 11),
				intCase("Large", 100, 101),
			}),
		piscineSpec(1031, "Power Of Two", "Piscine Quest 4", "IsPowerOfTwo(value int) bool",
			"Check whether a positive integer is an exact power of two.", "One integer.", "true for 1, 2, 4, 8, ... and false otherwise.",
			nil, intField, "IsPowerOfTwo(current.Payload.Value)", `false`,
			[]practiceCase{
				intBoolCase("Zero", 0, false),
				intBoolCase("One", 1, true),
				intBoolCase("Power", 64, true),
				intBoolCase("Between", 63, false),
				intBoolCase("Negative", -2, false),
			}),
		piscineSpec(1032, "Collatz Countdown", "Piscine Quest 4", "CollatzCountdown(value int) int",
			"Count how many Collatz steps are needed to reach 1.", "A positive integer.", "The number of transformations, or -1 for non-positive input.",
			nil, intField, "CollatzCountdown(current.Payload.Value)", `0`,
			[]practiceCase{
				intCase("Invalid", 0, -1),
				intCase("Already one", 1, 0),
				intCase("Even", 2, 1),
				intCase("Sequence", 12, 9),
			}),
	}
}

func piscineSpec(
	id int,
	title, difficulty, signature, objective, input, output string,
	builtins []string,
	inputFields, call, zero string,
	cases []practiceCase,
) practiceSpec {
	return zoneSpec(
		id,
		title,
		"Piscine ú "+strings.TrimPrefix(difficulty, "Piscine "),
		difficulty,
		signature,
		objective,
		input,
		output,
		builtins,
		inputFields,
		call,
		zero,
		cases,
	)
}

func ascendingTriplesExpected() string {
	var result strings.Builder
	for first := 0; first <= 7; first++ {
		for second := first + 1; second <= 8; second++ {
			for third := second + 1; third <= 9; third++ {
				if result.Len() > 0 {
					result.WriteString(", ")
				}
				fmt.Fprintf(&result, "%d%d%d", first, second, third)
			}
		}
	}
	result.WriteByte('\n')
	return result.String()
}

func ascendingPairsExpected() string {
	var result strings.Builder
	for left := 0; left <= 98; left++ {
		for right := left + 1; right <= 99; right++ {
			if result.Len() > 0 {
				result.WriteString(", ")
			}
			fmt.Fprintf(&result, "%02d %02d", left, right)
		}
	}
	result.WriteByte('\n')
	return result.String()
}

func descendingPairsExpected() string {
	var result strings.Builder
	for left := 99; left >= 1; left-- {
		for right := left - 1; right >= 0; right-- {
			if result.Len() > 0 {
				result.WriteString(", ")
			}
			fmt.Fprintf(&result, "%02d %02d", left, right)
		}
	}
	result.WriteByte('\n')
	return result.String()
}

func multiplicationTableExpected(value int) string {
	var result strings.Builder
	for multiplier := 1; multiplier <= 9; multiplier++ {
		fmt.Fprintf(&result, "%d x %d = %d\n", multiplier, value, multiplier*value)
	}
	return result.String()
}

func twoIntsSliceCase(name string, left, right int, expected []int) practiceCase {
	return pc(name, fmt.Sprintf("%d, %d", left, right), jsonString(expected), map[string]any{"left": left, "right": right})
}

func intSliceCase(name string, values, expected []int) practiceCase {
	return pc(name, fmt.Sprint(values), jsonString(expected), map[string]any{"values": values})
}

func intBoolCase(name string, value int, expected bool) practiceCase {
	return pc(name, fmt.Sprint(value), fmt.Sprint(expected), map[string]any{"value": value})
}

func stringsOutputCase(name string, values []string, expected string) practiceCase {
	return pc(name, jsonString(values), jsonString(expected), map[string]any{"values": values})
}

func stringsIntCase(name string, values []string, expected int) practiceCase {
	return pc(name, jsonString(values), fmt.Sprint(expected), map[string]any{"values": values})
}

func stringsSliceCase(name string, values, expected []string) practiceCase {
	return pc(name, jsonString(values), jsonString(expected), map[string]any{"values": values})
}

func runeOutputCase(name, value string, expected rune) practiceCase {
	return pc(name, jsonString(value), fmt.Sprint(expected), map[string]any{"value": value})
}

func runePositionCase(name, value string, position int, expected rune) practiceCase {
	return pc(name, fmt.Sprintf("%s, %d", jsonString(value), position), fmt.Sprint(expected), map[string]any{"value": value, "position": position})
}

func indexCase(name, value, target string, expected int) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(value), jsonString(target)), fmt.Sprint(expected), map[string]any{"value": value, "target": target})
}

func joinCase(name string, values []string, separator, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(values), jsonString(separator)), jsonString(expected), map[string]any{"values": values, "separator": separator})
}

func numberBaseCase(name string, value int, base, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%d, %s", value, jsonString(base)), jsonString(expected), map[string]any{"value": value, "base": base})
}

func intsBoolStringCase(name string, values []int, upper bool, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%v, %t", values, upper), jsonString(expected), map[string]any{"values": values, "upper": upper})
}

func flagsCase(name, value, insert string, order bool, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s, %t", jsonString(value), jsonString(insert), order), jsonString(expected), map[string]any{"value": value, "insert": insert, "order": order})
}

func rangeCase(name string, minimum, maximum int, expected []int) practiceCase {
	return pc(name, fmt.Sprintf("%d, %d", minimum, maximum), jsonString(expected), map[string]any{"minimum": minimum, "maximum": maximum})
}

func rangeNilCase(name string, minimum, maximum int) practiceCase {
	return pc(name, fmt.Sprintf("%d, %d", minimum, maximum), "null", map[string]any{"minimum": minimum, "maximum": maximum})
}

func valueStringsCase(name, value string, values, expected []string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(value), jsonString(values)), jsonString(expected), map[string]any{"value": value, "values": values})
}

func stringSliceCase(name, value string, expected []string) practiceCase {
	return pc(name, jsonString(value), jsonString(expected), map[string]any{"value": value})
}

func splitCase(name, value, separator string, expected []string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(value), jsonString(separator)), jsonString(expected), map[string]any{"value": value, "separator": separator})
}

func convertBaseCase(name, value, fromBase, toBase, expected string) practiceCase {
	input := fmt.Sprintf("%s, %s, %s", jsonString(value), jsonString(fromBase), jsonString(toBase))
	return pc(name, input, jsonString(expected), map[string]any{"value": value, "fromBase": fromBase, "toBase": toBase})
}

func valueBaseIntCase(name, value, base string, expected int) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(value), jsonString(base)), fmt.Sprint(expected), map[string]any{"value": value, "base": base})
}

func stringChunksCase(name, value string, size int, expected []string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %d", jsonString(value), size), jsonString(expected), map[string]any{"value": value, "size": size})
}

func stringChunksNilCase(name, value string, size int) practiceCase {
	return pc(name, fmt.Sprintf("%s, %d", jsonString(value), size), "null", map[string]any{"value": value, "size": size})
}

func virtualFilesCase(name string, arguments []string, files map[string]string, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(arguments), jsonString(files)), jsonString(expected), map[string]any{"arguments": arguments, "files": files})
}

func catCase(name string, names []string, files map[string]string, stdin, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s, %s", jsonString(names), jsonString(files), jsonString(stdin)), jsonString(expected), map[string]any{"names": names, "files": files, "stdin": stdin})
}

func tailCase(name string, count int, names []string, files map[string]string, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%d, %s, %s", count, jsonString(names), jsonString(files)), jsonString(expected), map[string]any{"count": count, "names": names, "files": files})
}

func stringsModeBoolCase(name string, values []string, mode string, expected bool) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(values), jsonString(mode)), fmt.Sprint(expected), map[string]any{"values": values, "mode": mode})
}

func stringsModeIntCase(name string, values []string, mode string, expected int) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(values), jsonString(mode)), fmt.Sprint(expected), map[string]any{"values": values, "mode": mode})
}

func intsOperationCase(name string, values []int, operation string, expected []int) practiceCase {
	return pc(name, fmt.Sprintf("%v, %s", values, jsonString(operation)), jsonString(expected), map[string]any{"values": values, "operation": operation})
}

func intsPredicateCase(name string, values []int, predicate string, expected []bool) practiceCase {
	return pc(name, fmt.Sprintf("%v, %s", values, jsonString(predicate)), jsonString(expected), map[string]any{"values": values, "predicate": predicate})
}

func foldCase(name string, values []int, initial int, operation string, expected int) practiceCase {
	return pc(name, fmt.Sprintf("%v, %d, %s", values, initial, jsonString(operation)), fmt.Sprint(expected), map[string]any{"values": values, "initial": initial, "operation": operation})
}

func intsOperationIntCase(name string, values []int, operation string, expected int) practiceCase {
	return pc(name, fmt.Sprintf("%v, %s", values, jsonString(operation)), fmt.Sprint(expected), map[string]any{"values": values, "operation": operation})
}

func intsOrderCase(name string, values []int, order string, expected bool) practiceCase {
	return pc(name, fmt.Sprintf("%v, %s", values, jsonString(order)), fmt.Sprint(expected), map[string]any{"values": values, "order": order})
}

func operationCase(name string, left int, operator string, right int, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%d, %s, %d", left, jsonString(operator), right), jsonString(expected), map[string]any{"left": left, "operator": operator, "right": right})
}

func intSliceIntCase(name string, values []int, expected int) practiceCase {
	return pc(name, fmt.Sprint(values), fmt.Sprint(expected), map[string]any{"values": values})
}

func fiveIntsCase(name string, values []int, expected int) practiceCase {
	payload := map[string]any{"a": values[0], "b": values[1], "c": values[2], "d": values[3], "e": values[4]}
	return pc(name, fmt.Sprint(values), fmt.Sprint(expected), payload)
}

func byteCase(name string, value, expected byte) practiceCase {
	return pc(name, fmt.Sprint(value), fmt.Sprint(expected), map[string]any{"value": value})
}

func byteStringCase(name string, value byte, expected string) practiceCase {
	return pc(name, fmt.Sprint(value), jsonString(expected), map[string]any{"value": value})
}

func valueWidthCase(name string, value, width int, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%d, %d", value, width), jsonString(expected), map[string]any{"value": value, "width": width})
}

func expressionCase(name, expression, expected string) practiceCase {
	return pc(name, jsonString(expression), jsonString(expected), map[string]any{"expression": expression})
}

func groupingCase(name, pattern, value string, expected []string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(pattern), jsonString(value)), jsonString(expected), map[string]any{"pattern": pattern, "value": value})
}

func findPairsCase(name string, values []int, target int, expected [][]int) practiceCase {
	return pc(name, fmt.Sprintf("%v, %d", values, target), jsonString(expected), map[string]any{"values": values, "target": target})
}

const (
	stringPositionFields   = "Value string `json:\"value\"`\n\tPosition int `json:\"position\"`"
	targetStringFields     = "Value string `json:\"value\"`\n\tTarget string `json:\"target\"`"
	stringsSeparatorFields = "Values []string `json:\"values\"`\n\tSeparator string `json:\"separator\"`"
	intBaseStringFields    = "Value int `json:\"value\"`\n\tBase string `json:\"base\"`"
	intsBoolFields         = "Values []int `json:\"values\"`\n\tUpper bool `json:\"upper\"`"
	flagsFields            = "Value string `json:\"value\"`\n\tInsert string `json:\"insert\"`\n\tOrder bool `json:\"order\"`"
	minMaxFields           = "Minimum int `json:\"minimum\"`\n\tMaximum int `json:\"maximum\"`"
	valueStringsFields     = "Value string `json:\"value\"`\n\tValues []string `json:\"values\"`"
	valueSeparatorFields   = "Value string `json:\"value\"`\n\tSeparator string `json:\"separator\"`"
	convertBaseFields      = "Value string `json:\"value\"`\n\tFromBase string `json:\"fromBase\"`\n\tToBase string `json:\"toBase\"`"
	valueBaseFields        = "Value string `json:\"value\"`\n\tBase string `json:\"base\"`"
	virtualFilesFields     = "Arguments []string `json:\"arguments\"`\n\tFiles map[string]string `json:\"files\"`"
	catFields              = "Names []string `json:\"names\"`\n\tFiles map[string]string `json:\"files\"`\n\tStdin string `json:\"stdin\"`"
	tailFields             = "Count int `json:\"count\"`\n\tNames []string `json:\"names\"`\n\tFiles map[string]string `json:\"files\"`"
	stringsModeFields      = "Values []string `json:\"values\"`\n\tMode string `json:\"mode\"`"
	intsOperationFields    = "Values []int `json:\"values\"`\n\tOperation string `json:\"operation\"`"
	intsPredicateFields    = "Values []int `json:\"values\"`\n\tPredicate string `json:\"predicate\"`"
	foldFields             = "Values []int `json:\"values\"`\n\tInitial int `json:\"initial\"`\n\tOperation string `json:\"operation\"`"
	intsOrderFields        = "Values []int `json:\"values\"`\n\tOrder string `json:\"order\"`"
	operationFields        = "Left int `json:\"left\"`\n\tOperator string `json:\"operator\"`\n\tRight int `json:\"right\"`"
	fiveIntsFields         = "A int `json:\"a\"`\n\tB int `json:\"b\"`\n\tC int `json:\"c\"`\n\tD int `json:\"d\"`\n\tE int `json:\"e\"`"
	byteField              = "Value byte `json:\"value\"`"
	valueWidthFields       = "Value int `json:\"value\"`\n\tWidth int `json:\"width\"`"
	expressionField        = "Expression string `json:\"expression\"`"
	patternValueFields     = "Pattern string `json:\"pattern\"`\n\tValue string `json:\"value\"`"
)
