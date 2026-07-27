package assessment

import (
	"fmt"
	"strings"
)

// piscineLevels contains the unique exercises found in kinoz01/zone01-Piscine.
// Duplicate solution variants and exercises already represented elsewhere in
// the curriculum are intentionally omitted.
func piscineLevels() []Level {
	specs := make([]practiceSpec, 0, 106)
	specs = append(specs, piscineQuestTwoSpecs()...)
	specs = append(specs, piscineQuestThreeAndFourSpecs()...)
	specs = append(specs, piscineQuestFiveSpecs()...)
	specs = append(specs, piscineQuestSixAndSevenSpecs()...)
	specs = append(specs, piscineQuestEightAndNineSpecs()...)
	specs = append(specs, piscineAdvancedSpecs()...)

	levels := make([]Level, 0, len(specs))
	for _, spec := range specs {
		levels = append(levels, compilePracticeExercise(sourcePiscine, spec))
	}
	return levels
}

func piscineQuestTwoSpecs() []practiceSpec {
	return []practiceSpec{
		piscinePrintSpec(1001, "Only Z", "Piscine Quest 2", "OnlyZ() string",
			"Return the single lowercase letter z.", "No input.", "Exactly \"z\".",
			nil, noFields, "OnlyZ()", `""`,
			[]practiceCase{pc("Exact output", ``, jsonString("z"), map[string]any{})}),
		piscinePrintSpec(1002, "Only One", "Piscine Quest 2", "OnlyOne() string",
			"Return the single character 1.", "No input.", "Exactly \"1\".",
			nil, noFields, "OnlyOne()", `""`,
			[]practiceCase{pc("Exact output", ``, jsonString("1"), map[string]any{})}),
		piscinePrintSpec(1003, "Only B", "Piscine Quest 2", "OnlyB() string",
			"Return the single uppercase letter B.", "No input.", "Exactly \"B\".",
			nil, noFields, "OnlyB()", `""`,
			[]practiceCase{pc("Exact output", ``, jsonString("B"), map[string]any{})}),
		piscinePrintSpec(1004, "Only F", "Piscine Quest 2", "OnlyF() string",
			"Return the single lowercase letter f.", "No input.", "Exactly \"f\".",
			nil, noFields, "OnlyF()", `""`,
			[]practiceCase{pc("Exact output", ``, jsonString("f"), map[string]any{})}),
		piscinePrintSpec(1005, "Hello World", "Piscine Quest 2", "HelloWorld() string",
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
		piscinePrintSpec(1008, "Print Alphabet", "Piscine Quest 2", "PrintAlphabet() string",
			"Generate the lowercase Latin alphabet.", "No input.", "abcdefghijklmnopqrstuvwxyz followed by newline.",
			nil, noFields, "PrintAlphabet()", `""`,
			[]practiceCase{pc("Alphabet", ``, jsonString("abcdefghijklmnopqrstuvwxyz\n"), map[string]any{})}),
		piscinePrintSpec(1009, "Print Reverse Alphabet", "Piscine Quest 2", "PrintReverseAlphabet() string",
			"Generate the lowercase Latin alphabet in reverse.", "No input.", "zyxwvutsrqponmlkjihgfedcba followed by newline.",
			nil, noFields, "PrintReverseAlphabet()", `""`,
			[]practiceCase{pc("Reverse alphabet", ``, jsonString("zyxwvutsrqponmlkjihgfedcba\n"), map[string]any{})}),
		piscinePrintSpec(1010, "Print Digits", "Piscine Quest 2", "PrintDigits() string",
			"Generate decimal digits in ascending order.", "No input.", "0123456789 followed by newline.",
			nil, noFields, "PrintDigits()", `""`,
			[]practiceCase{pc("Digits", ``, jsonString("0123456789\n"), map[string]any{})}),
		piscinePrintSpec(1011, "Countdown", "Piscine Quest 2", "Countdown() string",
			"Generate decimal digits in descending order.", "No input.", "9876543210 followed by newline.",
			nil, noFields, "Countdown()", `""`,
			[]practiceCase{pc("Countdown", ``, jsonString("9876543210\n"), map[string]any{})}),
		piscinePrintSpec(1012, "Is Negative", "Piscine Quest 2", "IsNegative(value int) string",
			"Classify an integer by its sign.", "One integer.", "\"T\\n\" for negative values and \"F\\n\" otherwise.",
			nil, intField, "IsNegative(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Negative", -1, "T\n"),
				intOutputCase("Zero", 0, "F\n"),
				intOutputCase("Positive", 1, "F\n"),
				intOutputCase("Large negative", -999, "T\n"),
			}),
		piscinePrintSpec(1013, "Print Number", "Piscine Quest 2", "PrintNbr(value int) string",
			"Convert any int to decimal output without strconv.", "One integer.", "Its exact base-10 representation.",
			nil, intField, "PrintNbr(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Zero", 0, "0"),
				intOutputCase("Positive", 42, "42"),
				intOutputCase("Negative", -42, "-42"),
				intOutputCase("Large", 987654321, "987654321"),
			}),
		piscinePrintSpec(1014, "Print Combination", "Piscine Quest 2", "PrintComb() string",
			"Generate every strictly ascending combination of three distinct digits.", "No input.", "Combinations from 012 through 789 separated by comma and space, then newline.",
			nil, noFields, "PrintComb()", `""`,
			[]practiceCase{pc("All combinations", ``, jsonString(ascendingTriplesExpected()), map[string]any{})}),
		piscinePrintSpec(1015, "Print Combination Two", "Piscine Quest 2", "PrintComb2() string",
			"Generate every ordered pair of distinct two-digit numbers.", "No input.", "Pairs from 00 01 through 98 99 separated by comma and space, then newline.",
			nil, noFields, "PrintComb2()", `""`,
			[]practiceCase{pc("All pairs", ``, jsonString(ascendingPairsExpected()), map[string]any{})}),
		piscinePrintSpec(1016, "Descending Combination", "Piscine Quest 2", "DescendComb() string",
			"Generate every pair of distinct two-digit numbers in descending order.", "No input.", "Pairs from 99 98 down through 01 00 separated by comma and space, then newline.",
			nil, noFields, "DescendComb()", `""`,
			[]practiceCase{pc("All descending pairs", ``, jsonString(descendingPairsExpected()), map[string]any{})}),
		piscinePrintSpec(1017, "Alternating Alphabet", "Piscine Quest 2", "DisplayAlphaM() string",
			"Generate the alphabet with odd positions lowercase and even positions uppercase.", "No input.", "aBcDeF...yZ followed by newline.",
			nil, noFields, "DisplayAlphaM()", `""`,
			[]practiceCase{pc("Alternating alphabet", ``, jsonString("aBcDeFgHiJkLmNoPqRsTuVwXyZ\n"), map[string]any{})}),
		piscinePrintSpec(1018, "Reverse Alternating Alphabet", "Piscine Quest 2", "DisplayAlphaReverseM() string",
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
			intStringCases([]string{"", "go", "γ", "🙂", "café"}, []int{0, 2, 1, 1, 4})),
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

func piscineQuestFiveSpecs() []practiceSpec {
	return []practiceSpec{
		piscineSpec(1033, "Concatenate", "Piscine Quest 5", "Concat(left string, right string) string",
			"Join two strings without a separator.", "Two strings.", "left immediately followed by right.",
			nil, twoStringsFields, "Concat(current.Payload.Left, current.Payload.Right)", `""`,
			[]practiceCase{
				twoStringOutputCase("Both empty", "", "", ""),
				twoStringOutputCase("Left empty", "", "go", "go"),
				twoStringOutputCase("Right empty", "go", "", "go"),
				twoStringOutputCase("Both", "go", "pher", "gopher"),
			}),
		piscineSpec(1034, "Basic Join", "Piscine Quest 5", "BasicJoin(values []string) string",
			"Concatenate every string in a slice.", "One string slice.", "All elements joined without a separator.",
			nil, stringsField, "BasicJoin(current.Payload.Values)", `""`,
			[]practiceCase{
				stringsOutputCase("Empty", []string{}, ""),
				stringsOutputCase("One", []string{"go"}, "go"),
				stringsOutputCase("Several", []string{"go", "pher"}, "gopher"),
				stringsOutputCase("Includes empty", []string{"a", "", "b"}, "ab"),
			}),
		piscineSpec(1035, "Compare Strings", "Piscine Quest 5", "Compare(left string, right string) int",
			"Compare two strings lexicographically.", "Two strings.", "-1 when left is smaller, 0 when equal, and 1 when greater.",
			nil, twoStringsFields, "Compare(current.Payload.Left, current.Payload.Right)", `0`,
			[]practiceCase{
				twoStringIntCase("Equal", "abc", "abc", 0),
				twoStringIntCase("Less", "abc", "abd", -1),
				twoStringIntCase("Greater", "z", "a", 1),
				twoStringIntCase("Prefix", "go", "gopher", -1),
			}),
		piscineSpec(1036, "First Rune", "Piscine Quest 5", "FirstRune(value string) rune",
			"Return the first Unicode code point.", "A non-empty UTF-8 string.", "Its first rune.",
			nil, stringField, "FirstRune(current.Payload.Value)", `0`,
			[]practiceCase{
				runeOutputCase("ASCII", "go", 'g'),
				runeOutputCase("Greek", "γεια", 'γ'),
				runeOutputCase("Emoji", "🙂ok", '🙂'),
			}),
		piscineSpec(1037, "Last Rune", "Piscine Quest 5", "LastRune(value string) rune",
			"Return the final Unicode code point.", "A non-empty UTF-8 string.", "Its last rune.",
			nil, stringField, "LastRune(current.Payload.Value)", `0`,
			[]practiceCase{
				runeOutputCase("ASCII", "go", 'o'),
				runeOutputCase("Greek", "γεια", 'α'),
				runeOutputCase("Emoji", "ok🙂", '🙂'),
			}),
		piscineSpec(1038, "Nth Rune", "Piscine Quest 5", "NRune(value string, position int) rune",
			"Return a rune by one-based position.", "A UTF-8 string and position.", "The selected rune, or 0 when the position is invalid.",
			nil, stringPositionFields, "NRune(current.Payload.Value, current.Payload.Position)", `0`,
			[]practiceCase{
				runePositionCase("First", "hello", 1, 'h'),
				runePositionCase("Middle", "γεια", 2, 'ε'),
				runePositionCase("Too high", "go", 3, 0),
				runePositionCase("Zero", "go", 0, 0),
			}),
		piscineSpec(1039, "String Index", "Piscine Quest 5", "Index(value string, target string) int",
			"Find the first byte index of a substring.", "A string and substring.", "The first index, 0 for an empty target, or -1 when absent.",
			[]string{"len"}, targetStringFields, "Index(current.Payload.Value, current.Payload.Target)", `-1`,
			[]practiceCase{
				indexCase("Empty target", "hello", "", 0),
				indexCase("Beginning", "hello", "he", 0),
				indexCase("Middle", "hello", "ll", 2),
				indexCase("Absent", "hello", "x", -1),
				indexCase("Too long", "go", "gopher", -1),
			}),
		piscineSpec(1040, "Is Alphanumeric", "Piscine Quest 5", "IsAlpha(value string) bool",
			"Validate that a string contains only ASCII letters and digits.", "One string.", "true for alphanumeric or empty input.",
			nil, stringField, "IsAlpha(current.Payload.Value)", `false`,
			boolStringCases([]string{"", "abcXYZ", "abc123", "hello world", "café"}, []bool{true, true, true, false, false})),
		piscineSpec(1041, "Is Lowercase", "Piscine Quest 5", "IsLower(value string) bool",
			"Validate that every character is an ASCII lowercase letter.", "One string.", "true for lowercase-only or empty input.",
			nil, stringField, "IsLower(current.Payload.Value)", `false`,
			boolStringCases([]string{"", "abc", "a1", "Abc", "é"}, []bool{true, true, false, false, false})),
		piscineSpec(1042, "Is Numeric", "Piscine Quest 5", "IsNumeric(value string) bool",
			"Validate that every character is an ASCII decimal digit.", "One string.", "true for digit-only or empty input.",
			nil, stringField, "IsNumeric(current.Payload.Value)", `false`,
			boolStringCases([]string{"", "0123", "-1", "1.0", "１２"}, []bool{true, true, false, false, false})),
		piscineSpec(1043, "Is Printable", "Piscine Quest 5", "IsPrintable(value string) bool",
			"Validate that every rune is printable.", "One string.", "true when all runes are at least the space character.",
			nil, stringField, "IsPrintable(current.Payload.Value)", `false`,
			boolStringCases([]string{"", "hello", "a b", "line\nbreak", "\t"}, []bool{true, true, true, false, false})),
		piscineSpec(1044, "Is Uppercase", "Piscine Quest 5", "IsUpper(value string) bool",
			"Validate that every character is an ASCII uppercase letter.", "One string.", "true for uppercase-only or empty input.",
			nil, stringField, "IsUpper(current.Payload.Value)", `false`,
			boolStringCases([]string{"", "ABC", "A1", "AbC", "É"}, []bool{true, true, false, false, false})),
		piscineSpec(1045, "To Lower", "Piscine Quest 5", "ToLower(value string) string",
			"Convert ASCII uppercase letters to lowercase.", "One string.", "The converted string with non-uppercase data unchanged.",
			nil, stringField, "ToLower(current.Payload.Value)", `""`,
			stringCases([]string{"", "HELLO", "Go 123!", "already"}, []string{"", "hello", "go 123!", "already"})),
		piscineSpec(1046, "To Upper", "Piscine Quest 5", "ToUpper(value string) string",
			"Convert ASCII lowercase letters to uppercase.", "One string.", "The converted string with non-lowercase data unchanged.",
			nil, stringField, "ToUpper(current.Payload.Value)", `""`,
			stringCases([]string{"", "hello", "Go 123!", "ALREADY"}, []string{"", "HELLO", "GO 123!", "ALREADY"})),
		piscineSpec(1047, "Capitalize", "Piscine Quest 5", "Capitalize(value string) string",
			"Uppercase each alphanumeric word's first character and lowercase the rest.", "One string.", "The transformed string; non-alphanumeric characters separate words.",
			[]string{"len"}, stringField, "Capitalize(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "Hello! How are you?", "56street", "a-b_C", "GOlang"},
				[]string{"", "Hello! How Are You?", "56street", "A-B_C", "Golang"},
			)),
		piscineSpec(1048, "Join", "Piscine Quest 5", "Join(values []string, separator string) string",
			"Join strings with a separator between neighboring elements.", "A string slice and separator.", "The joined string with no leading or trailing separator.",
			[]string{"len"}, stringsSeparatorFields, "Join(current.Payload.Values, current.Payload.Separator)", `""`,
			[]practiceCase{
				joinCase("Empty", []string{}, ",", ""),
				joinCase("One", []string{"go"}, ",", "go"),
				joinCase("Several", []string{"go", "lang"}, "-", "go-lang"),
				joinCase("Empty separator", []string{"a", "b"}, "", "ab"),
			}),
		piscineSpec(1049, "Trim Atoi", "Piscine Quest 5", "TrimAtoi(value string) int",
			"Extract all decimal digits from a string and apply an earlier minus sign.", "One string.", "The collected integer, negative only when - occurs before the first digit.",
			nil, stringField, "TrimAtoi(current.Payload.Value)", `0`,
			intStringCases(
				[]string{"", "123", "abc-12def3", "a1-23", "--42", "no digits"},
				[]int{0, 123, -123, 123, -42, 0},
			)),
		piscinePrintSpec(1050, "Print Number In Order", "Piscine Quest 5", "PrintNbrInOrder(value int) string",
			"Sort the decimal digits of a non-negative integer.", "One non-negative integer.", "Its digits in ascending order.",
			nil, intField, "PrintNbrInOrder(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Zero", 0, "0"),
				intOutputCase("Sorted", 1234, "1234"),
				intOutputCase("Mixed", 3214, "1234"),
				intOutputCase("Repeated", 90901, "00199"),
			}),
		piscinePrintSpec(1051, "Print Number In Base", "Piscine Quest 5", "PrintNbrBase(value int, base string) string",
			"Format an integer using a caller-provided digit alphabet.", "An integer and base alphabet.", "The representation, or \"NV\" when the base is invalid.",
			[]string{"len"}, intBaseStringFields, "PrintNbrBase(current.Payload.Value, current.Payload.Base)", `""`,
			[]practiceCase{
				numberBaseCase("Binary", 10, "01", "1010"),
				numberBaseCase("Hex", 255, "0123456789abcdef", "ff"),
				numberBaseCase("Negative", -10, "01", "-1010"),
				numberBaseCase("Short base", 10, "0", "NV"),
				numberBaseCase("Duplicate digit", 10, "001", "NV"),
			}),
		piscineSpec(1052, "Alphabet Mirror", "Piscine Quest 5", "AlphaMirror(value string) string",
			"Replace each ASCII letter with its opposite in the alphabet.", "One string.", "a↔z, b↔y, preserving case and other characters.",
			nil, stringField, "AlphaMirror(current.Payload.Value)", `""`,
			stringCases([]string{"", "abc", "XYZ", "Hello!", "azAZ"}, []string{"", "zyx", "CBA", "Svool!", "zaZA"})),
		piscineSpec(1053, "ROT13", "Piscine Quest 5", "Rot13(value string) string",
			"Rotate ASCII letters by 13 places.", "One string.", "The ROT13 transformation with case preserved.",
			nil, stringField, "Rot13(current.Payload.Value)", `""`,
			stringCases([]string{"", "abc", "Hello, World!", "uryyb"}, []string{"", "nop", "Uryyb, Jbeyq!", "hello"})),
		piscineSpec(1054, "ROT14", "Piscine Quest 5", "Rot14(value string) string",
			"Rotate ASCII letters by 14 places.", "One string.", "The ROT14 transformation with case preserved.",
			nil, stringField, "Rot14(current.Payload.Value)", `""`,
			stringCases([]string{"", "abc", "XYZ", "Hello!"}, []string{"", "opq", "LMN", "Vszzc!"})),
		piscineSpec(1055, "Switch Case", "Piscine Quest 5", "SwitchCase(value string) string",
			"Reverse the case of every ASCII letter.", "One string.", "Uppercase becomes lowercase and lowercase becomes uppercase.",
			nil, stringField, "SwitchCase(current.Payload.Value)", `""`,
			stringCases([]string{"", "Hello!", "ABC xyz", "123"}, []string{"", "hELLO!", "abc XYZ", "123"})),
		piscineSpec(1056, "Pig Latin", "Piscine Quest 5", "PigLatin(value string) string",
			"Transform one word using the piscine Pig Latin rules.", "One lowercase ASCII word.", "Append \"ay\" after an initial vowel; otherwise move the leading consonant cluster and append \"ay\".",
			[]string{"len"}, stringField, "PigLatin(current.Payload.Value)", `""`,
			stringCases([]string{"", "apple", "pig", "smile", "rhythm"}, []string{"", "appleay", "igpay", "ilesmay", "No vowels"})),
	}
}

func piscineQuestSixAndSevenSpecs() []practiceSpec {
	return []practiceSpec{
		piscineSpec(1057, "Parameter Count", "Piscine Quest 6", "ParamCount(arguments []string) int",
			"Count the command-line arguments supplied to a program.", "A slice containing arguments, without the program name.", "The number of arguments.",
			[]string{"len"}, stringsField, "ParamCount(current.Payload.Values)", `0`,
			[]practiceCase{
				stringsIntCase("None", []string{}, 0),
				stringsIntCase("One", []string{"go"}, 1),
				stringsIntCase("Several", []string{"go", "is", "fun"}, 3),
			}),
		piscineSpec(1058, "First Parameter", "Piscine Quest 6", "FirstParam(arguments []string) string",
			"Return the first command-line argument.", "A slice containing arguments, without the program name.", "The first argument followed by newline, or newline when none exists.",
			[]string{"len"}, stringsField, "FirstParam(current.Payload.Values)", `"\n"`,
			[]practiceCase{
				stringsOutputCase("None", []string{}, "\n"),
				stringsOutputCase("One", []string{"go"}, "go\n"),
				stringsOutputCase("Several", []string{"first", "last"}, "first\n"),
			}),
		piscineSpec(1059, "Last Parameter", "Piscine Quest 6", "LastParam(arguments []string) string",
			"Return the final command-line argument.", "A slice containing arguments, without the program name.", "The last argument followed by newline, or newline when none exists.",
			[]string{"len"}, stringsField, "LastParam(current.Payload.Values)", `"\n"`,
			[]practiceCase{
				stringsOutputCase("None", []string{}, "\n"),
				stringsOutputCase("One", []string{"go"}, "go\n"),
				stringsOutputCase("Several", []string{"first", "last"}, "last\n"),
			}),
		piscinePrintSpec(1060, "Print Parameters", "Piscine Quest 6", "PrintParams(arguments []string) string",
			"Format command-line arguments in their original order.", "A slice containing arguments, without the program name.", "Every argument on its own line.",
			nil, stringsField, "PrintParams(current.Payload.Values)", `""`,
			[]practiceCase{
				stringsOutputCase("None", []string{}, ""),
				stringsOutputCase("One", []string{"go"}, "go\n"),
				stringsOutputCase("Several", []string{"go", "lang"}, "go\nlang\n"),
			}),
		piscinePrintSpec(1061, "Reverse Parameters", "Piscine Quest 6", "ReverseParams(arguments []string) string",
			"Format command-line arguments from last to first.", "A slice containing arguments, without the program name.", "Every argument on its own line in reverse order.",
			[]string{"len"}, stringsField, "ReverseParams(current.Payload.Values)", `""`,
			[]practiceCase{
				stringsOutputCase("None", []string{}, ""),
				stringsOutputCase("One", []string{"go"}, "go\n"),
				stringsOutputCase("Several", []string{"one", "two", "three"}, "three\ntwo\none\n"),
			}),
		piscinePrintSpec(1062, "Sort Parameters", "Piscine Quest 6", "SortParams(arguments []string) []string",
			"Sort command-line arguments by ascending ASCII order.", "A slice containing arguments, without the program name.", "The sorted arguments, one per line.",
			[]string{"append", "len"}, stringsField, "SortParams(current.Payload.Values)", `nil`,
			[]practiceCase{
				stringsOutputCase("None", []string{}, ""),
				stringsOutputCase("Mixed", []string{"z", "a", "m"}, "a\nm\nz\n"),
				stringsOutputCase("Case-sensitive", []string{"a", "B", "A"}, "A\nB\na\n"),
			}),
		piscinePrintSpec(1063, "Number Convert Alpha", "Piscine Quest 6", "NbrConvertAlpha(values []int, upper bool) string",
			"Convert numbers 1 through 26 into alphabet letters.", "Integer values and a case flag.", "Letters for valid values; invalid values become spaces. Append newline.",
			nil, intsBoolFields, "NbrConvertAlpha(current.Payload.Values, current.Payload.Upper)", `"\n"`,
			[]practiceCase{
				intsBoolStringCase("Lowercase", []int{8, 5, 12, 12, 15}, false, "hello\n"),
				intsBoolStringCase("Uppercase", []int{7, 15}, true, "GO\n"),
				intsBoolStringCase("Invalid", []int{1, 0, 26, 27}, false, "a z \n"),
			}),
		piscinePrintSpec(1064, "Program Name", "Piscine Quest 6", "ProgramName(name string) string",
			"Format the executable name as the piscine program would print it.", "One program name.", "The name followed by newline.",
			nil, stringField, "ProgramName(current.Payload.Value)", `"\n"`,
			stringCases([]string{"main", "./app", "go-task"}, []string{"main\n", "./app\n", "go-task\n"})),
		piscineSpec(1065, "Flags", "Piscine Quest 6", "ApplyFlags(value string, insert string, order bool) string",
			"Apply insert and order flags to a string.", "A value, optional text to append, and an order flag.", "The transformed string followed by newline.",
			[]string{"len"}, flagsFields, "ApplyFlags(current.Payload.Value, current.Payload.Insert, current.Payload.Order)", `"\n"`,
			[]practiceCase{
				flagsCase("Unchanged", "go", "", false, "go\n"),
				flagsCase("Insert", "go", "lang", false, "golang\n"),
				flagsCase("Order", "dcba", "", true, "abcd\n"),
				flagsCase("Both", "dc", "ba", true, "abcd\n"),
			}),
		piscineSpec(1066, "Append Range", "Piscine Quest 7", "AppendRange(minimum int, maximum int) []int",
			"Build an ascending half-open integer range using append.", "Two integer bounds.", "All integers from minimum through maximum-1; empty when minimum is not smaller.",
			[]string{"append"}, minMaxFields, "AppendRange(current.Payload.Minimum, current.Payload.Maximum)", `[]int{}`,
			[]practiceCase{
				rangeCase("Positive", 1, 4, []int{1, 2, 3}),
				rangeCase("Across zero", -2, 2, []int{-2, -1, 0, 1}),
				rangeCase("Equal", 3, 3, []int{}),
				rangeCase("Descending", 4, 1, []int{}),
			}),
		piscineSpec(1067, "Make Range", "Piscine Quest 7", "MakeRange(minimum int, maximum int) []int",
			"Build an ascending half-open integer range using make.", "Two integer bounds.", "All integers from minimum through maximum-1; nil when minimum is not smaller.",
			[]string{"make"}, minMaxFields, "MakeRange(current.Payload.Minimum, current.Payload.Maximum)", `nil`,
			[]practiceCase{
				rangeCase("Positive", 1, 4, []int{1, 2, 3}),
				rangeCase("Across zero", -2, 2, []int{-2, -1, 0, 1}),
				rangeNilCase("Equal", 3, 3),
				rangeNilCase("Descending", 4, 1),
			}),
		piscineSpec(1068, "Add Front", "Piscine Quest 7", "AddFront(value string, values []string) []string",
			"Insert a non-empty string at the front of a slice.", "A value and string slice.", "The prefixed slice; an empty value leaves the slice unchanged.",
			[]string{"append", "len", "make"}, valueStringsFields, "AddFront(current.Payload.Value, current.Payload.Values)", `nil`,
			[]practiceCase{
				valueStringsCase("Empty slice", "go", []string{}, []string{"go"}),
				valueStringsCase("Existing", "go", []string{"lang"}, []string{"go", "lang"}),
				valueStringsCase("Empty value", "", []string{"go"}, []string{"go"}),
			}),
		piscineSpec(1069, "Concatenate Parameters", "Piscine Quest 7", "ConcatParams(arguments []string) string",
			"Join command-line arguments with newline separators.", "A string slice.", "Arguments separated by newline, without a trailing newline.",
			[]string{"len"}, stringsField, "ConcatParams(current.Payload.Values)", `""`,
			[]practiceCase{
				stringsOutputCase("None", []string{}, ""),
				stringsOutputCase("One", []string{"go"}, "go"),
				stringsOutputCase("Several", []string{"go", "lang"}, "go\nlang"),
			}),
		piscineSpec(1070, "Split Whitespace", "Piscine Quest 7", "SplitWhiteSpaces(value string) []string",
			"Split a string on spaces, tabs, and newlines.", "One string.", "Non-empty words in encounter order.",
			[]string{"append", "len"}, stringField, "SplitWhiteSpaces(current.Payload.Value)", `nil`,
			[]practiceCase{
				stringSliceCase("Empty", "", []string{}),
				stringSliceCase("Spaces", "  go   lang ", []string{"go", "lang"}),
				stringSliceCase("Mixed", "go\tis\nfun", []string{"go", "is", "fun"}),
			}),
		piscineSpec(1071, "Split", "Piscine Quest 7", "Split(value string, separator string) []string",
			"Split a string at each exact separator occurrence.", "A string and non-empty separator.", "All fields, including empty fields at the edges.",
			[]string{"append", "len"}, valueSeparatorFields, "Split(current.Payload.Value, current.Payload.Separator)", `nil`,
			[]practiceCase{
				splitCase("No match", "golang", ",", []string{"golang"}),
				splitCase("Several", "go,is,fun", ",", []string{"go", "is", "fun"}),
				splitCase("Edges", ",go,", ",", []string{"", "go", ""}),
				splitCase("Wide separator", "a--b", "--", []string{"a", "b"}),
			}),
		piscinePrintSpec(1072, "Print Words Tables", "Piscine Quest 7", "PrintWordsTables(values []string) string",
			"Format every word from a table on its own line.", "A string slice.", "Every element followed by newline.",
			nil, stringsField, "PrintWordsTables(current.Payload.Values)", `""`,
			[]practiceCase{
				stringsOutputCase("Empty", []string{}, ""),
				stringsOutputCase("One", []string{"go"}, "go\n"),
				stringsOutputCase("Several", []string{"go", "lang"}, "go\nlang\n"),
			}),
		piscineSpec(1073, "Convert Base", "Piscine Quest 7", "ConvertBase(value string, fromBase string, toBase string) string",
			"Convert a valid integer representation between custom bases.", "A number and two valid digit alphabets.", "The equivalent representation in the target base.",
			[]string{"len"}, convertBaseFields, "ConvertBase(current.Payload.Value, current.Payload.FromBase, current.Payload.ToBase)", `""`,
			[]practiceCase{
				convertBaseCase("Binary to decimal", "1010", "01", "0123456789", "10"),
				convertBaseCase("Decimal to hex", "255", "0123456789", "0123456789abcdef", "ff"),
				convertBaseCase("Negative", "-10", "0123456789", "01", "-1010"),
				convertBaseCase("Zero", "0", "0123456789", "01", "0"),
			}),
		piscineSpec(1074, "Atoi Base", "Piscine Quest 7", "AtoiBase(value string, base string) int",
			"Parse a signed integer using a custom base alphabet.", "A number and valid digit alphabet.", "Its base-10 integer value.",
			[]string{"len"}, valueBaseFields, "AtoiBase(current.Payload.Value, current.Payload.Base)", `0`,
			[]practiceCase{
				valueBaseIntCase("Binary", "1010", "01", 10),
				valueBaseIntCase("Hex", "ff", "0123456789abcdef", 255),
				valueBaseIntCase("Negative", "-101", "01", -5),
				valueBaseIntCase("Zero", "0", "0123456789", 0),
			}),
		piscineSpec(1075, "String Chunks", "Piscine Quest 7", "StringChunks(value string, size int) []string",
			"Divide a string into fixed-size chunks.", "A string and positive chunk size.", "Chunks in order; the final chunk may be shorter. Invalid sizes return nil.",
			[]string{"append", "len"}, stringSizeFields, "StringChunks(current.Payload.Value, current.Payload.Size)", `nil`,
			[]practiceCase{
				stringChunksCase("Exact", "abcdef", 2, []string{"ab", "cd", "ef"}),
				stringChunksCase("Remainder", "abcde", 2, []string{"ab", "cd", "e"}),
				stringChunksCase("Large", "go", 5, []string{"go"}),
				stringChunksNilCase("Invalid", "go", 0),
			}),
	}
}

func piscineQuestEightAndNineSpecs() []practiceSpec {
	return []practiceSpec{
		piscineSpec(1076, "Is Even", "Piscine Quest 8", "IsEven(value int) bool",
			"Determine whether an integer is divisible by two.", "One integer.", "true for even values, including zero and negative values.",
			nil, intField, "IsEven(current.Payload.Value)", `false`,
			[]practiceCase{intBoolCase("Zero", 0, true), intBoolCase("Odd", 3, false), intBoolCase("Negative even", -4, true)}),
		piscineSpec(1077, "Point", "Piscine Quest 8", "SetPoint() []int",
			"Create the piscine point whose x coordinate is 42 and y coordinate is 21.", "No input.", "A two-item slice containing x then y.",
			nil, noFields, "SetPoint()", `nil`,
			[]practiceCase{pc("Coordinates", ``, jsonString([]int{42, 21}), map[string]any{})}),
		piscineSpec(1078, "Display File", "Piscine Quest 8", "DisplayFile(arguments []string, files map[string]string) string",
			"Display one virtual file while handling missing or extra arguments.", "File-name arguments and a map of available file contents.", "The file content or the required diagnostic message.",
			[]string{"len"}, virtualFilesFields, "DisplayFile(current.Payload.Arguments, current.Payload.Files)", `""`,
			[]practiceCase{
				virtualFilesCase("Missing name", []string{}, map[string]string{"a": "A"}, "File name missing\n"),
				virtualFilesCase("Too many names", []string{"a", "b"}, map[string]string{"a": "A"}, "Too many arguments\n"),
				virtualFilesCase("Available", []string{"a"}, map[string]string{"a": "hello\n"}, "hello\n"),
				virtualFilesCase("Unavailable", []string{"x"}, map[string]string{"a": "A"}, "Cannot read file\n"),
			}),
		piscinePrintSpec(1079, "Cat", "Piscine Quest 8", "Cat(names []string, files map[string]string, stdin string) string",
			"Concatenate virtual files, or return standard input when no names are supplied.", "File names, available contents, and simulated standard input.", "Contents in order; missing files contribute a diagnostic line.",
			[]string{"len"}, catFields, "Cat(current.Payload.Names, current.Payload.Files, current.Payload.Stdin)", `""`,
			[]practiceCase{
				catCase("Standard input", []string{}, map[string]string{}, "hello\n", "hello\n"),
				catCase("One file", []string{"a"}, map[string]string{"a": "A\n"}, "", "A\n"),
				catCase("Several", []string{"a", "b"}, map[string]string{"a": "A", "b": "B"}, "", "AB"),
				catCase("Missing", []string{"x"}, map[string]string{}, "", "Cannot read x\n"),
			}),
		piscineSpec(1080, "ZTail", "Piscine Quest 8", "ZTail(count int, names []string, files map[string]string) string",
			"Return the last count bytes of each virtual file.", "A positive byte count, file names, and available contents.", "Tail content in order, with file headings when several names are supplied.",
			[]string{"len"}, tailFields, "ZTail(current.Payload.Count, current.Payload.Names, current.Payload.Files)", `""`,
			[]practiceCase{
				tailCase("Short file", 5, []string{"a"}, map[string]string{"a": "go"}, "go"),
				tailCase("Suffix", 3, []string{"a"}, map[string]string{"a": "golang"}, "ang"),
				tailCase("Zero", 0, []string{"a"}, map[string]string{"a": "go"}, ""),
				tailCase("Missing", 3, []string{"x"}, map[string]string{}, "Cannot read x\n"),
			}),
		piscineSpec(1081, "Any", "Piscine Quest 8", "Any(values []string, mode string) bool",
			"Report whether any value satisfies a selected predicate.", "String values and mode: empty, numeric, or lowercase.", "true when at least one item matches.",
			nil, stringsModeFields, "Any(current.Payload.Values, current.Payload.Mode)", `false`,
			[]practiceCase{
				stringsModeBoolCase("Any empty", []string{"go", ""}, "empty", true),
				stringsModeBoolCase("None empty", []string{"go"}, "empty", false),
				stringsModeBoolCase("Any numeric", []string{"go", "42"}, "numeric", true),
				stringsModeBoolCase("Any lowercase", []string{"GO", "go"}, "lowercase", true),
			}),
		piscineSpec(1082, "Count If", "Piscine Quest 8", "CountIf(values []string, mode string) int",
			"Count values that satisfy a selected predicate.", "String values and mode: empty, numeric, or lowercase.", "The number of matching items.",
			nil, stringsModeFields, "CountIf(current.Payload.Values, current.Payload.Mode)", `0`,
			[]practiceCase{
				stringsModeIntCase("Empty", []string{"", "go", ""}, "empty", 2),
				stringsModeIntCase("Numeric", []string{"1", "go", "42"}, "numeric", 2),
				stringsModeIntCase("Lowercase", []string{"go", "GO", "lang"}, "lowercase", 2),
			}),
		piscineSpec(1083, "For Each", "Piscine Quest 8", "ForEach(values []int, operation string) []int",
			"Apply the same operation to every integer.", "Integer values and operation: double, square, or negate.", "A result slice with one transformed value per input.",
			[]string{"append"}, intsOperationFields, "ForEach(current.Payload.Values, current.Payload.Operation)", `nil`,
			[]practiceCase{
				intsOperationCase("Double", []int{1, 2, 3}, "double", []int{2, 4, 6}),
				intsOperationCase("Square", []int{-2, 3}, "square", []int{4, 9}),
				intsOperationCase("Negate", []int{-1, 0, 2}, "negate", []int{1, 0, -2}),
			}),
		piscineSpec(1084, "Map", "Piscine Quest 8", "Map(values []int, predicate string) []bool",
			"Evaluate the same predicate for every integer.", "Integer values and predicate: positive, even, or prime.", "A boolean result for each value.",
			[]string{"append"}, intsPredicateFields, "Map(current.Payload.Values, current.Payload.Predicate)", `nil`,
			[]practiceCase{
				intsPredicateCase("Positive", []int{-1, 0, 2}, "positive", []bool{false, false, true}),
				intsPredicateCase("Even", []int{1, 2, 3}, "even", []bool{false, true, false}),
				intsPredicateCase("Prime", []int{1, 2, 4, 5}, "prime", []bool{false, true, false, true}),
			}),
		piscineSpec(1085, "Fold Int", "Piscine Quest 8", "FoldInt(values []int, initial int, operation string) int",
			"Fold integers into an accumulator from left to right.", "Values, an initial accumulator, and operation: add, multiply, or maximum.", "The final accumulator.",
			nil, foldFields, "FoldInt(current.Payload.Values, current.Payload.Initial, current.Payload.Operation)", `0`,
			[]practiceCase{
				foldCase("Add", []int{1, 2, 3}, 0, "add", 6),
				foldCase("Multiply", []int{2, 3, 4}, 1, "multiply", 24),
				foldCase("Maximum", []int{-2, 7, 3}, -10, "maximum", 7),
			}),
		piscineSpec(1086, "Reduce Int", "Piscine Quest 8", "ReduceInt(values []int, operation string) int",
			"Reduce a non-empty integer slice from its first value.", "Values and operation: add, multiply, or maximum.", "The reduced value; zero for empty input.",
			[]string{"len"}, intsOperationFields, "ReduceInt(current.Payload.Values, current.Payload.Operation)", `0`,
			[]practiceCase{
				intsOperationIntCase("Empty", []int{}, "add", 0),
				intsOperationIntCase("Add", []int{1, 2, 3}, "add", 6),
				intsOperationIntCase("Multiply", []int{2, 3, 4}, "multiply", 24),
				intsOperationIntCase("Maximum", []int{-2, 7, 3}, "maximum", 7),
			}),
		piscineSpec(1087, "Is Sorted", "Piscine Quest 8", "IsSorted(values []int, order string) bool",
			"Check whether integers follow an ascending or descending comparator.", "Integer values and order: ascending or descending.", "true when every neighboring pair follows the requested order.",
			[]string{"len"}, intsOrderFields, "IsSorted(current.Payload.Values, current.Payload.Order)", `false`,
			[]practiceCase{
				intsOrderCase("Empty", []int{}, "ascending", true),
				intsOrderCase("Ascending", []int{1, 2, 2, 3}, "ascending", true),
				intsOrderCase("Not ascending", []int{2, 1}, "ascending", false),
				intsOrderCase("Descending", []int{3, 2, 1}, "descending", true),
			}),
		piscineSpec(1088, "Sort Word Array", "Piscine Quest 8", "SortWordArray(values []string) []string",
			"Sort words lexicographically without the sort package.", "A string slice.", "A new slice in ascending order.",
			[]string{"append", "len"}, stringsField, "SortWordArray(current.Payload.Values)", `nil`,
			[]practiceCase{
				stringsSliceCase("Empty", []string{}, []string{}),
				stringsSliceCase("Words", []string{"zone", "go", "alpha"}, []string{"alpha", "go", "zone"}),
				stringsSliceCase("Case", []string{"a", "B", "A"}, []string{"A", "B", "a"}),
			}),
		piscineSpec(1089, "Do Operation", "Piscine Quest 9", "DoOperation(left int, operator string, right int) string",
			"Evaluate one integer arithmetic operation.", "Two integers and one of +, -, *, /, or %.", "The result, or the piscine error text for invalid operations and division by zero.",
			nil, operationFields, "DoOperation(current.Payload.Left, current.Payload.Operator, current.Payload.Right)", `""`,
			[]practiceCase{
				operationCase("Add", 4, "+", 2, "6\n"),
				operationCase("Divide", 7, "/", 2, "3\n"),
				operationCase("Divide by zero", 7, "/", 0, "No division by 0\n"),
				operationCase("Invalid", 7, "x", 2, "0\n"),
			}),
		piscineSpec(1090, "Maximum", "Piscine Quest 9", "Maximum(values []int) int",
			"Return the greatest integer in a slice.", "An integer slice.", "Its maximum, or zero for empty input.",
			[]string{"len"}, intsField, "Maximum(current.Payload.Values)", `0`,
			[]practiceCase{
				intSliceIntCase("Empty", []int{}, 0),
				intSliceIntCase("Positive", []int{1, 9, 3}, 9),
				intSliceIntCase("Negative", []int{-9, -2, -7}, -2),
				intSliceIntCase("Repeated", []int{4, 4}, 4),
			}),
	}
}

func piscineAdvancedSpecs() []practiceSpec {
	return []practiceSpec{
		piscineSpec(1091, "Abort", "Piscine Advanced", "Abort(a int, b int, c int, d int, e int) int",
			"Return the median of five integers.", "Five integers.", "The value that would occupy the middle position when sorted.",
			nil, fiveIntsFields, "Abort(current.Payload.A, current.Payload.B, current.Payload.C, current.Payload.D, current.Payload.E)", `0`,
			[]practiceCase{
				fiveIntsCase("Ordered", []int{1, 2, 3, 4, 5}, 3),
				fiveIntsCase("Mixed", []int{9, 1, 7, 3, 5}, 5),
				fiveIntsCase("Duplicates", []int{2, 1, 2, 3, 2}, 2),
			}),
		piscineSpec(1092, "Food Delivery Time", "Piscine Advanced", "FoodDeliveryTime(order string) int",
			"Calculate a delivery estimate from the foods in an order.", "A comma-separated order using burger, chips, nuggets, and drink.", "The longest preparation time among valid items; -1 for an invalid order.",
			[]string{"len"}, stringField, "FoodDeliveryTime(current.Payload.Value)", `-1`,
			intStringCases([]string{"burger", "chips,drink", "nuggets,burger", "", "pizza"}, []int{15, 10, 15, -1, -1})),
		piscineSpec(1093, "Unmatch", "Piscine Advanced", "Unmatch(values []int) int",
			"Find the value without a matching duplicate.", "Integers where every value except at most one appears an even number of times.", "The unmatched value, or -1 when all values pair.",
			nil, intsField, "Unmatch(current.Payload.Values)", `-1`,
			[]practiceCase{
				intSliceIntCase("All paired", []int{1, 1, 2, 2}, -1),
				intSliceIntCase("One unmatched", []int{1, 2, 1}, 2),
				intSliceIntCase("Zero", []int{0, 3, 3}, 0),
				intSliceIntCase("Empty", []int{}, -1),
			}),
		piscineSpec(1094, "Reverse Bits", "Piscine Advanced", "ReverseBits(value byte) byte",
			"Reverse the order of all eight bits in a byte.", "One byte value.", "The bit-reversed byte.",
			nil, byteField, "ReverseBits(current.Payload.Value)", `0`,
			[]practiceCase{
				byteCase("Zero", 0, 0),
				byteCase("One", 1, 128),
				byteCase("Pattern", 0b00110100, 0b00101100),
				byteCase("All ones", 255, 255),
			}),
		piscineSpec(1095, "Swap Bits", "Piscine Advanced", "SwapBits(value byte) byte",
			"Swap the high and low four-bit halves of a byte.", "One byte value.", "The byte with its nibbles exchanged.",
			nil, byteField, "SwapBits(current.Payload.Value)", `0`,
			[]practiceCase{
				byteCase("Zero", 0, 0),
				byteCase("Low nibble", 0x0f, 0xf0),
				byteCase("Pattern", 0x3c, 0xc3),
				byteCase("Equal halves", 0xaa, 0xaa),
			}),
		piscineSpec(1096, "Print Bits", "Piscine Advanced", "PrintBits(value byte) string",
			"Format all eight bits of a byte.", "One byte value.", "Exactly eight binary digits followed by newline.",
			nil, byteField, "PrintBits(current.Payload.Value)", `""`,
			[]practiceCase{
				byteStringCase("Zero", 0, "00000000\n"),
				byteStringCase("One", 1, "00000001\n"),
				byteStringCase("Pattern", 0x3c, "00111100\n"),
				byteStringCase("Maximum", 255, "11111111\n"),
			}),
		piscineSpec(1097, "Two's Complement", "Piscine Advanced", "TwosComplement(value int, width int) string",
			"Format an integer as a fixed-width two's-complement bit pattern.", "An integer and bit width from 1 through 64.", "Exactly width binary digits.",
			nil, valueWidthFields, "TwosComplement(current.Payload.Value, current.Payload.Width)", `""`,
			[]practiceCase{
				valueWidthCase("Positive", 5, 8, "00000101"),
				valueWidthCase("Minus one", -1, 8, "11111111"),
				valueWidthCase("Negative", -5, 8, "11111011"),
				valueWidthCase("Four bits", -2, 4, "1110"),
			}),
		piscineSpec(1098, "Print Hex", "Piscine Advanced", "PrintHex(value int) string",
			"Format a non-negative integer in lowercase hexadecimal.", "One non-negative integer.", "The hexadecimal digits followed by newline.",
			nil, intField, "PrintHex(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Zero", 0, "0\n"),
				intOutputCase("Ten", 10, "a\n"),
				intOutputCase("Byte", 255, "ff\n"),
				intOutputCase("Large", 4096, "1000\n"),
			}),
		piscineSpec(1099, "Multiplication Table", "Piscine Advanced", "MultiplicationTable(value int) string",
			"Generate the multiplication table from one through nine.", "One integer.", "Nine lines in the form n x value = result.",
			nil, intField, "MultiplicationTable(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Two", 2, multiplicationTableExpected(2)),
				intOutputCase("Zero", 0, multiplicationTableExpected(0)),
				intOutputCase("Negative", -3, multiplicationTableExpected(-3)),
			}),
		piscineSpec(1100, "Reverse Words", "Piscine Advanced", "ReverseWords(value string) string",
			"Reverse the order of whitespace-separated words.", "One string.", "Words joined by single spaces followed by newline.",
			[]string{"append", "len"}, stringField, "ReverseWords(current.Payload.Value)", `"\n"`,
			stringCases(
				[]string{"", "hello", "hello world", "  go\tis\nfun "},
				[]string{"\n", "hello\n", "world hello\n", "fun is go\n"},
			)),
		piscineSpec(1101, "Rotate String", "Piscine Advanced", "RotateString(value string) string",
			"Move the first whitespace-separated word to the end.", "One string.", "Normalized words after one left rotation, followed by newline.",
			[]string{"append", "len"}, stringField, "RotateString(current.Payload.Value)", `"\n"`,
			stringCases(
				[]string{"", "hello", "hello world", "  one two three "},
				[]string{"\n", "hello\n", "world hello\n", "two three one\n"},
			)),
		piscineSpec(1102, "Roman Numbers", "Piscine Advanced", "RomanNumbers(value int) string",
			"Convert an integer from 1 through 3999 to Roman numerals.", "One integer.", "The uppercase Roman representation, or ERROR for an out-of-range value.",
			nil, intField, "RomanNumbers(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("One", 1, "I"),
				intOutputCase("Subtractive", 944, "CMXLIV"),
				intOutputCase("Maximum", 3999, "MMMCMXCIX"),
				intOutputCase("Invalid", 0, "ERROR"),
			}),
		piscineSpec(1103, "RPN Calculator", "Piscine Advanced", "RPNCalculator(expression string) string",
			"Evaluate a space-separated reverse Polish notation expression.", "Integers and +, -, *, /, or % operators.", "The integer result followed by newline, or Error for malformed input.",
			[]string{"append", "len"}, expressionField, "RPNCalculator(current.Payload.Expression)", `""`,
			[]practiceCase{
				expressionCase("Add", "2 3 +", "5\n"),
				expressionCase("Nested", "5 1 2 + 4 * + 3 -", "14\n"),
				expressionCase("Negative", "-3 2 *", "-6\n"),
				expressionCase("Malformed", "2 +", "Error\n"),
			}),
		piscineSpec(1104, "Grouping", "Piscine Advanced", "Grouping(pattern string, value string) []string",
			"Find substrings matching a simplified parenthesized alternation pattern.", "A pattern such as (a|b) and a value to scan.", "Matches in encounter order.",
			[]string{"append", "len"}, patternValueFields, "Grouping(current.Payload.Pattern, current.Payload.Value)", `nil`,
			[]practiceCase{
				groupingCase("Alternatives", "(cat|dog)", "a cat and a dog", []string{"cat", "dog"}),
				groupingCase("Repeated", "(go|lang)", "gogolang", []string{"go", "go", "lang"}),
				groupingCase("No match", "(a|b)", "xyz", []string{}),
				groupingCase("Literal", "(01|edu)", "01-edu", []string{"01", "edu"}),
			}),
		piscineSpec(1105, "Find Pairs", "Piscine Advanced", "FindPairs(values []int, target int) [][]int",
			"Find disjoint pairs whose sum equals a target.", "Integer values and target sum.", "Pairs in encounter order; each input position may be used once.",
			[]string{"append", "len"}, intsTargetFields, "FindPairs(current.Payload.Values, current.Payload.Target)", `nil`,
			[]practiceCase{
				findPairsCase("None", []int{1, 2}, 10, [][]int{}),
				findPairsCase("One", []int{1, 2, 3}, 4, [][]int{{1, 3}}),
				findPairsCase("Several", []int{1, 5, 3, 3, 2, 4}, 6, [][]int{{1, 5}, {3, 3}, {2, 4}}),
				findPairsCase("Duplicates", []int{2, 2, 2, 2}, 4, [][]int{{2, 2}, {2, 2}}),
			}),
		piscineSpec(1106, "Word Flip", "Piscine Advanced", "WordFlip(value string) string",
			"Reverse the characters of each word while preserving word order.", "One whitespace-separated string.", "Flipped words joined by single spaces followed by newline.",
			[]string{"append", "len"}, stringField, "WordFlip(current.Payload.Value)", `"\n"`,
			stringCases(
				[]string{"", "hello", "hello world", "  go\tis\nfun "},
				[]string{"\n", "olleh\n", "olleh dlrow\n", "og si nuf\n"},
			)),
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
		"Piscine · "+strings.TrimPrefix(difficulty, "Piscine "),
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

func piscinePrintSpec(
	id int,
	title, difficulty, signature, objective, input, output string,
	builtins []string,
	inputFields, call, zero string,
	cases []practiceCase,
) practiceSpec {
	return printedPracticeSpec(piscineSpec(
		id, title, difficulty, signature, objective, input, output,
		builtins, inputFields, call, zero, cases,
	))
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
