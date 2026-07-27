package assessment

import "fmt"

// zone01Levels adapts the checkpoint catalogue to the browser runner. Exercises
// that were originally command-line programs return their output so every case
// can be tested without hiding stdout from the learner.
func zone01Levels() []Level {
	specs := []practiceSpec{
		zoneSpec(22, "Only A", "Output · constants", "Checkpoint 5%",
			"OnlyA() string",
			"Return the single lowercase letter a.",
			"No input.", "Exactly \"a\".",
			nil, noFields, "OnlyA()", `""`,
			[]practiceCase{
				pc("Exact output", ``, jsonString("a"), map[string]any{}),
			}),
		zoneSpec(23, "Print If Not", "Conditionals · strings", "Checkpoint 10%",
			"PrintIfNot(value string) string",
			"Choose an exact message from the byte length of a string.",
			"One string.", "\"G\\n\" when its length is below 3; otherwise \"Invalid Input\\n\".",
			[]string{"len"}, stringField, "PrintIfNot(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "a", "ab", "abc", "abcdef"},
				[]string{"G\n", "G\n", "G\n", "Invalid Input\n", "Invalid Input\n"},
			)),
		zoneSpec(24, "Print If", "Conditionals · strings", "Checkpoint 10%",
			"PrintIf(value string) string",
			"Return G for an empty or sufficiently long string.",
			"One string.", "\"G\\n\" for empty input or length at least 3; otherwise \"Invalid Input\\n\".",
			[]string{"len"}, stringField, "PrintIf(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "a", "14", "abc", "abcdef"},
				[]string{"G\n", "Invalid Input\n", "Invalid Input\n", "G\n", "G\n"},
			)),
		zoneSpec(25, "Rectangle Perimeter", "Arithmetic · validation", "Checkpoint 10%",
			"RectPerimeter(width int, height int) int",
			"Calculate a rectangle's perimeter.",
			"Width and height as integers.", "2 × width + 2 × height, or -1 when either value is negative.",
			nil, widthHeightFields, "RectPerimeter(current.Payload.Width, current.Payload.Height)", `0`,
			[]practiceCase{
				pc("Square", `2, 2`, `8`, map[string]any{"width": 2, "height": 2}),
				pc("Rectangle", `10, 2`, `24`, map[string]any{"width": 10, "height": 2}),
				pc("Zero width", `0, 8`, `16`, map[string]any{"width": 0, "height": 8}),
				pc("Negative width", `-1, 8`, `-1`, map[string]any{"width": -1, "height": 8}),
				pc("Negative height", `8, -1`, `-1`, map[string]any{"width": 8, "height": -1}),
			}),
		zoneSpec(26, "Count Character", "Strings · runes", "Checkpoint 10%",
			"CountChar(value string, target rune) int",
			"Count exact occurrences of one rune.",
			"A UTF-8 string and a rune.", "The number of matching runes.",
			nil, stringRuneFields, "CountChar(current.Payload.Value, current.Payload.Target)", `0`,
			[]practiceCase{
				runeCase("Repeated", "Hello World", 'l', 3),
				runeCase("Absent", "gopher", 'x', 0),
				runeCase("Spaces", "   ", ' ', 3),
				runeCase("Digit", "The 7 deadly sins", '7', 1),
				runeCase("Unicode", "γεια γ", 'γ', 2),
			}),
		zoneSpec(27, "Check Number", "Strings · scanning", "Checkpoint 10%",
			"CheckNumber(value string) bool",
			"Detect whether a string contains an ASCII digit.",
			"One string.", "true when at least one character is between '0' and '9'.",
			nil, stringField, "CheckNumber(current.Payload.Value)", `false`,
			boolStringCases(
				[]string{"", "Hello", "Hello1", "007", "１２"},
				[]bool{false, false, true, true, false},
			)),
		zoneSpec(28, "Retain First Half", "Strings · slicing", "Checkpoint 10%",
			"RetainFirstHalf(value string) string",
			"Return the first half of a string.",
			"One string.", "The first floor(length/2) bytes; a one-byte string remains unchanged.",
			[]string{"len"}, stringField, "RetainFirstHalf(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "A", "Go", "Hello", "Hello World"},
				[]string{"", "A", "G", "He", "Hello"},
			)),
		zoneSpec(29, "Count Alpha", "Strings · ASCII", "Checkpoint 10%",
			"CountAlpha(value string) int",
			"Count ASCII alphabetic characters.",
			"One string.", "The number of bytes in A-Z or a-z.",
			nil, stringField, "CountAlpha(current.Payload.Value)", `0`,
			intStringCases(
				[]string{"", "Hello world", "H e l l o", "H1e2l3l4o", "café"},
				[]int{0, 10, 5, 5, 3},
			)),
		zoneSpec(30, "First Word", "Strings · scanning", "Checkpoint 20%",
			"FirstWord(value string) string",
			"Extract the first space-delimited word.",
			"One string.", "The first word followed by a newline; just a newline when no word exists.",
			[]string{"len"}, stringField, "FirstWord(current.Payload.Value)", `"\n"`,
			stringCases(
				[]string{"", "hello there", "   hello   there", "one", "   "},
				[]string{"\n", "hello\n", "hello\n", "one\n", "\n"},
			)),
		zoneSpec(31, "Last Word", "Strings · scanning", "Checkpoint 20%",
			"LastWord(value string) string",
			"Extract the last space-delimited word.",
			"One string.", "The last word followed by a newline; just a newline when no word exists.",
			[]string{"len"}, stringField, "LastWord(current.Payload.Value)", `"\n"`,
			stringCases(
				[]string{"", "hello there", " lorem,ipsum ", "one", "   "},
				[]string{"\n", "there\n", "lorem,ipsum\n", "one\n", "\n"},
			)),
		zoneSpec(32, "Fish And Chips", "Conditionals · modulo", "Checkpoint 20%",
			"FishAndChips(value int) string",
			"Classify a number by divisibility by 2 and 3.",
			"One integer.", "The exact fish/chips message, with negative and non-divisible errors.",
			nil, intField, "FishAndChips(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Fish", 4, "fish"),
				intOutputCase("Chips", 9, "chips"),
				intOutputCase("Both", 6, "fish and chips"),
				intOutputCase("Negative", -6, "error: number is negative"),
				intOutputCase("Neither", 5, "error: non divisible"),
			}),
		zoneSpec(33, "Digit Length", "Arithmetic · bases", "Checkpoint 20%",
			"DigitLen(value int, base int) int",
			"Count the digits needed to represent an integer in a base.",
			"An integer and a base.", "The repeated-division count; -1 when base is outside 2..36.",
			nil, intBaseFields, "DigitLen(current.Payload.Value, current.Payload.Base)", `0`,
			[]practiceCase{
				intBaseCase("Decimal", 100, 10, 3),
				intBaseCase("Binary", 100, 2, 7),
				intBaseCase("Negative", -100, 16, 2),
				intBaseCase("Zero", 0, 10, 0),
				intBaseCase("Invalid base", 100, 1, -1),
			}),
		zoneSpec(34, "Search Replace", "Strings · replacement", "Checkpoint 20%",
			"SearchReplace(value string, old string, replacement string) string",
			"Replace every exact occurrence of one character.",
			"A string plus one-character old and replacement strings.", "The replaced string; unchanged when old is absent.",
			[]string{"len"}, replaceFields, "SearchReplace(current.Payload.Value, current.Payload.Old, current.Payload.Replacement)", `""`,
			[]practiceCase{
				replaceCase("One match", "hella there", "a", "o", "hello there"),
				replaceCase("Several matches", "hallo thara", "a", "e", "hello there"),
				replaceCase("Absent", "abcd", "z", "l", "abcd"),
				replaceCase("Empty input", "", "a", "b", ""),
				replaceCase("Repeated", "aaaa", "a", "x", "xxxx"),
			}),
		zoneSpec(35, "Repeat Alpha", "Strings · nested loops", "Checkpoint 20%",
			"RepeatAlpha(value string) string",
			"Repeat each ASCII letter by its one-based alphabet position.",
			"One string.", "Letters repeated case-insensitively by position; non-letters copied once.",
			nil, stringField, "RepeatAlpha(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "abc", "Choumi.", "aZ", "abacadaba 01!"},
				[]string{"", "abbccc", "CCChhhhhhhhooooooooooooooouuuuuuuuuuuuuuuuuuuuummmmmmmmmmmmmiiiiiiiii.", "aZZZZZZZZZZZZZZZZZZZZZZZZZZ", "abbacccaddddabba 01!"},
			)),
		zoneSpec(36, "Greatest Common Divisor", "Arithmetic · Euclid", "Checkpoint 20%",
			"Gcd(left uint, right uint) uint",
			"Find the greatest common divisor of two unsigned integers.",
			"Two unsigned integers.", "Their greatest common divisor, or 0 when either input is 0.",
			nil, twoUintsFields, "Gcd(current.Payload.Left, current.Payload.Right)", `0`,
			[]practiceCase{
				uintCase("Coprime", 17, 3, 1),
				uintCase("Common factor", 42, 12, 6),
				uintCase("Same", 9, 9, 9),
				uintCase("Left zero", 0, 9, 0),
				uintCase("Right zero", 9, 0, 0),
			}),
		zoneSpec(37, "Camel To Snake Case", "Strings · validation", "Checkpoint 20%",
			"CamelToSnakeCase(value string) string",
			"Convert valid lowerCamelCase or UpperCamelCase text to snake_case.",
			"An ASCII string.", "The string with underscores before capitals, or the original when camelCase rules are broken.",
			[]string{"len"}, stringField, "CamelToSnakeCase(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "HelloWorld", "helloWorld", "camelToSnakeCase", "CAMELtoSnackCASE", "hey2"},
				[]string{"", "Hello_World", "hello_World", "camel_To_Snake_Case", "CAMELtoSnackCASE", "hey2"},
			)),
		zoneSpec(38, "Hash Code", "ASCII · modular arithmetic", "Checkpoint 20%",
			"HashCode(value string) string",
			"Transform each character with the checkpoint hash equation.",
			"An ASCII string.", "(character + byte length) modulo 127, adding 33 when the result is not printable.",
			[]string{"len"}, stringField, "HashCode(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "A", "AB", "BAC", "Hello World"},
				[]string{"", "B", "CD", "EDF", "Spwwz+bz}wo"},
			)),
		zoneSpec(39, "Third Time Is A Charm", "Strings · indexing", "Checkpoint 35%",
			"ThirdTimeIsACharm(value string) string",
			"Collect every third byte from a string.",
			"One string.", "Bytes 3, 6, 9, ... followed by a newline.",
			[]string{"len"}, stringField, "ThirdTimeIsACharm(current.Payload.Value)", `"\n"`,
			stringCases(
				[]string{"", "1", "12", "123456789", "a b c d e f"},
				[]string{"\n", "\n", "\n", "369\n", "b e\n"},
			)),
		zoneSpec(40, "From To", "Loops · formatting", "Checkpoint 35%",
			"FromTo(from int, to int) string",
			"Format an inclusive ascending or descending two-digit range.",
			"Two integers.", "Comma-separated zero-padded values plus newline, or \"Invalid\\n\" outside 0..99.",
			nil, fromToFields, "FromTo(current.Payload.From, current.Payload.To)", `""`,
			[]practiceCase{
				fromToCase("Ascending", 1, 4, "01, 02, 03, 04\n"),
				fromToCase("Descending", 4, 1, "04, 03, 02, 01\n"),
				fromToCase("Same", 10, 10, "10\n"),
				fromToCase("Invalid high", 100, 10, "Invalid\n"),
				fromToCase("Invalid low", -1, 10, "Invalid\n"),
			}),
		zoneSpec(41, "Is Capitalized", "Strings · word boundaries", "Checkpoint 35%",
			"IsCapitalized(value string) bool",
			"Check whether every space-delimited word starts with uppercase or non-alphabetic data.",
			"One string.", "false for empty input or any word beginning with a lowercase ASCII letter.",
			[]string{"len"}, stringField, "IsCapitalized(current.Payload.Value)", `false`,
			boolStringCases(
				[]string{"", "Hello How Are You", "Hello! How are you?", "Whats 4this 100K?", "!!!!Whatsthis4"},
				[]bool{false, true, false, true, true},
			)),
		zoneSpec(42, "Find Previous Prime", "Arithmetic · primality", "Checkpoint 35%",
			"FindPrevPrime(value int) int",
			"Find the largest prime less than or equal to an integer.",
			"One integer.", "The closest previous prime, or 0 when none exists.",
			nil, intField, "FindPrevPrime(current.Payload.Value)", `0`,
			[]practiceCase{
				intCase("Prime", 5, 5),
				intCase("Composite", 4, 3),
				intCase("Two", 2, 2),
				intCase("One", 1, 0),
				intCase("Negative", -10, 0),
			}),
		zoneSpec(43, "Integer To ASCII", "Arithmetic · conversion", "Checkpoint 35%",
			"Itoa(value int) string",
			"Convert an integer to decimal text without strconv.",
			"One integer.", "Its base-10 representation including a minus sign when negative.",
			nil, intField, "Itoa(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Zero", 0, "0"),
				intOutputCase("Positive", 12345, "12345"),
				intOutputCase("Negative", -1234, "-1234"),
				intOutputCase("One digit", 7, "7"),
				intOutputCase("Large", 987654321, "987654321"),
			}),
		zoneSpec(44, "Clean String", "Strings · whitespace", "Checkpoint 35%",
			"CleanStr(value string) string",
			"Normalize spaces and tabs between words.",
			"One string.", "Words joined by one space and followed by newline; blank input becomes newline.",
			[]string{"len"}, stringField, "CleanStr(current.Payload.Value)", `"\n"`,
			stringCases(
				[]string{"", "   ", "hello", " only    it's  harder   ", "\thow\t funny "},
				[]string{"\n", "\n", "hello\n", "only it's harder\n", "how funny\n"},
			)),
		zoneSpec(45, "Expand String", "Strings · whitespace", "Checkpoint 35%",
			"ExpandStr(value string) string",
			"Normalize spaces and tabs to exactly three spaces between words.",
			"One string.", "Words joined by three spaces and followed by newline; no words produce an empty string.",
			[]string{"len"}, stringField, "ExpandStr(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "   ", "hello", " only  it's harder   ", "\thow\t funny "},
				[]string{"", "", "hello\n", "only   it's   harder\n", "how   funny\n"},
			)),
		zoneSpec(46, "We Are Unique", "Strings · sets", "Checkpoint 35%",
			"WeAreUnique(left string, right string) int",
			"Count distinct characters present in only one of two strings.",
			"Two strings.", "The symmetric-difference character count; -1 when both strings are empty.",
			nil, twoStringsFields, "WeAreUnique(current.Payload.Left, current.Payload.Right)", `0`,
			[]practiceCase{
				twoStringIntCase("Both empty", "", "", -1),
				twoStringIntCase("Overlap", "foo", "boo", 2),
				twoStringIntCase("Disjoint", "abc", "def", 6),
				twoStringIntCase("Same", "abc", "abc", 0),
				twoStringIntCase("Duplicates", "aaaa", "bbb", 2),
			}),
		zoneSpec(47, "Zip String", "Strings · run-length encoding", "Checkpoint 35%",
			"ZipString(value string) string",
			"Run-length encode consecutive characters.",
			"One string.", "Each run written as its decimal count followed by the character.",
			nil, stringField, "ZipString(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "a", "aaa", "abbccc", "YouuungFellllas"},
				[]string{"", "1a", "3a", "1a2b3c", "1Y1o3u1n1g1F1e4l1a1s"},
			)),
		zoneSpec(48, "Print Reverse Combo", "Loops · enumeration", "Checkpoint 35%",
			"PrintRevCombo() string",
			"Generate every strictly descending three-digit combination from 987 to 210.",
			"No input.", "Combinations in descending order, separated by comma and space, followed by newline.",
			nil, noFields, "PrintRevCombo()", `""`,
			[]practiceCase{
				pc("Complete sequence", ``, jsonString(reverseComboExpected()), map[string]any{}),
			}),
		zoneSpec(49, "Print Memory", "Bytes · hexadecimal", "Checkpoint 35%",
			"PrintMemory(value [10]byte) string",
			"Render ten bytes as hexadecimal and printable ASCII.",
			"An array of exactly 10 bytes.", "Three lowercase hex rows (4, 4, 2 bytes), then an ASCII row; every row ends with newline.",
			[]string{"len"}, bytesField, "PrintMemory(current.Payload.Value)", `""`,
			[]practiceCase{
				memoryCase("Example", [10]byte{'h', 'e', 'l', 'l', 'o', 16, 21, '*'}, "68 65 6c 6c\n6f 10 15 2a\n00 00\nhello..*..\n"),
				memoryCase("Zeros", [10]byte{}, "00 00 00 00\n00 00 00 00\n00 00\n..........\n"),
				memoryCase("Printable", [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}, "30 31 32 33\n34 35 36 37\n38 39\n0123456789\n"),
				memoryCase("Boundaries", [10]byte{31, 32, 33, 126, 127}, "1f 20 21 7e\n7f 00 00 00\n00 00\n. !~......\n"),
				memoryCase("High bytes", [10]byte{255, 128, 'A'}, "ff 80 41 00\n00 00 00 00\n00 00\n..A.......\n"),
			}),
		zoneSpec(50, "Concat Slice", "Slices · append", "Checkpoint 50%",
			"ConcatSlice(left []int, right []int) []int",
			"Concatenate two integer slices into a new slice.",
			"Two integer slices.", "All left values followed by all right values.",
			[]string{"append"}, twoIntSlicesFields, "ConcatSlice(current.Payload.Left, current.Payload.Right)", `[]int{}`,
			[]practiceCase{
				twoSliceCase("Both empty", []int{}, []int{}, []int{}),
				twoSliceCase("Left empty", []int{}, []int{1, 2}, []int{1, 2}),
				twoSliceCase("Right empty", []int{1, 2}, []int{}, []int{1, 2}),
				twoSliceCase("Both", []int{1, 2, 3}, []int{4, 5}, []int{1, 2, 3, 4, 5}),
				twoSliceCase("Duplicates", []int{1, 1}, []int{1}, []int{1, 1, 1}),
			}),
		zoneSpec(51, "Save And Miss", "Strings · grouped scanning", "Checkpoint 50%",
			"SaveAndMiss(value string, size int) string",
			"Keep one group of bytes, skip the next, and repeat.",
			"A string and group size.", "The kept groups; the original string when size is zero or negative.",
			[]string{"len"}, stringSizeFields, "SaveAndMiss(current.Payload.Value, current.Payload.Size)", `""`,
			[]practiceCase{
				stringSizeCase("Example", "123456789", 3, "123789"),
				stringSizeCase("Empty", "", 3, ""),
				stringSizeCase("Zero", "hello", 0, "hello"),
				stringSizeCase("Negative", "hello", -2, "hello"),
				stringSizeCase("Partial group", "abcdefgh", 3, "abcgh"),
			}),
		zoneSpec(52, "Hidden P", "Strings · subsequences", "Checkpoint 50%",
			"HiddenP(needle string, haystack string) bool",
			"Check whether one string is a subsequence of another.",
			"A needle and a haystack.", "true when every needle byte appears in order in the haystack.",
			[]string{"len"}, needleFields, "HiddenP(current.Payload.Needle, current.Payload.Haystack)", `false`,
			[]practiceCase{
				hiddenCase("Empty needle", "", "anything", true),
				hiddenCase("Exact", "abc", "abc", true),
				hiddenCase("Hidden", "abc", "2altrb53c.sse", true),
				hiddenCase("Wrong order", "abc", "btarc", false),
				hiddenCase("Missing repeat", "DD", "DABC", false),
			}),
		zoneSpec(53, "Word Match", "Strings · subsequences", "Checkpoint 50%",
			"WdMatch(needle string, haystack string) string",
			"Return a word only when it can be formed in order from another string.",
			"A needle and a haystack.", "needle when it is a subsequence; otherwise an empty string.",
			[]string{"len"}, needleFields, "WdMatch(current.Payload.Needle, current.Payload.Haystack)", `""`,
			[]practiceCase{
				hiddenStringCase("Empty", "", "abc", ""),
				hiddenStringCase("Exact", "123", "123", "123"),
				hiddenStringCase("Hidden", "faya", "fgvvfdxcacpolhyghbreda", "faya"),
				hiddenStringCase("Missing", "faya", "fgvvfdxcacpolhyghbred", ""),
				hiddenStringCase("Repeated", "error", "rrerrrfiiljdfxjyuifrrvcoojh", ""),
			}),
		zoneSpec(54, "Intersection", "Strings · deduplication", "Checkpoint 50%",
			"Inter(left string, right string) string",
			"Collect distinct characters present in both strings.",
			"Two strings.", "Shared characters without duplicates, ordered by first appearance in left.",
			nil, twoStringsFields, "Inter(current.Payload.Left, current.Payload.Right)", `""`,
			[]practiceCase{
				twoStringOutputCase("Empty", "", "abc", ""),
				twoStringOutputCase("None", "abc", "XYZ", ""),
				twoStringOutputCase("Example", "padinton", "paqefwtdjetyiytjneytjoeyjnejeyj", "padinto"),
				twoStringOutputCase("Duplicates", "aabbcc", "banana", "ab"),
				twoStringOutputCase("Order", "cba", "abc", "cba"),
			}),
		zoneSpec(55, "Union", "Strings · deduplication", "Checkpoint 50%",
			"Union(left string, right string) string",
			"Collect distinct characters appearing in either string.",
			"Two strings.", "Unique characters in first-appearance order across left then right.",
			nil, twoStringsFields, "Union(current.Payload.Left, current.Payload.Right)", `""`,
			[]practiceCase{
				twoStringOutputCase("Both empty", "", "", ""),
				twoStringOutputCase("Left only", "aab", "", "ab"),
				twoStringOutputCase("Right only", "", "bba", "ba"),
				twoStringOutputCase("Overlap", "abc", "bcd", "abcd"),
				twoStringOutputCase("Example", "rien", "cette phrase ne cache rien", "rienct phas"),
			}),
		zoneSpec(56, "Concat Alternate", "Slices · interleaving", "Checkpoint 50%",
			"ConcatAlternate(left []int, right []int) []int",
			"Interleave two slices, beginning with the longer slice.",
			"Two integer slices.", "Alternating values while both remain, then the longer slice's tail.",
			[]string{"append", "len"}, twoIntSlicesFields, "ConcatAlternate(current.Payload.Left, current.Payload.Right)", `[]int{}`,
			[]practiceCase{
				twoSliceCase("Empty", []int{}, []int{}, []int{}),
				twoSliceCase("Equal", []int{1, 2, 3}, []int{4, 5, 6}, []int{1, 4, 2, 5, 3, 6}),
				twoSliceCase("Right longer", []int{1, 2, 3}, []int{4, 5, 6, 7}, []int{4, 1, 5, 2, 6, 3, 7}),
				twoSliceCase("Left longer", []int{2, 4, 6}, []int{1}, []int{2, 1, 4, 6}),
				twoSliceCase("One empty", []int{1, 2}, []int{}, []int{1, 2}),
			}),
		zoneSpec(57, "Chunk", "Slices · partitioning", "Checkpoint 50%",
			"Chunk(values []int, size int) [][]int",
			"Partition a slice into consecutive chunks.",
			"An integer slice and a non-negative size.", "A non-nil slice of chunks; empty when size is zero.",
			[]string{"append", "len"}, intsSizeFields, "Chunk(current.Payload.Values, current.Payload.Size)", `[][]int{}`,
			[]practiceCase{
				chunkCase("Empty", []int{}, 3, [][]int{}),
				chunkCase("Zero size", []int{1, 2}, 0, [][]int{}),
				chunkCase("Exact", []int{0, 1, 2, 3}, 2, [][]int{{0, 1}, {2, 3}}),
				chunkCase("Remainder", []int{0, 1, 2, 3, 4}, 2, [][]int{{0, 1}, {2, 3}, {4}}),
				chunkCase("Large size", []int{1, 2}, 10, [][]int{{1, 2}}),
			}),
		zoneSpec(58, "Reverse String Capitalization", "Strings · casing", "Checkpoint 50%",
			"ReverseStrCap(value string) string",
			"Uppercase the last letter of every word and lowercase its other letters.",
			"One ASCII string.", "The transformed string with spacing preserved.",
			[]string{"len"}, stringField, "ReverseStrCap(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "First SMALL TesT", "SEconD Test IS a LItTLE EasIEr", " Go ", "abc  DEF"},
				[]string{"", "firsT smalL tesT", "seconD tesT iS A littlE easieR", " gO ", "abC  deF"},
			)),
		zoneSpec(59, "Can Jump", "Slices · simulation", "Checkpoint 50%",
			"CanJump(values []uint) bool",
			"Follow exact jumps and determine whether they land on the final index.",
			"A slice of non-negative jump sizes.", "true only when repeated exact jumps reach and stay at the last index.",
			[]string{"len"}, uintsField, "CanJump(current.Payload.Values)", `false`,
			[]practiceCase{
				uintSliceCase("Empty", []uint{}, false),
				uintSliceCase("Single", []uint{0}, true),
				uintSliceCase("Reachable", []uint{2, 3, 1, 1, 4}, true),
				uintSliceCase("Blocked", []uint{3, 2, 1, 0, 4}, false),
				uintSliceCase("Overshoots", []uint{3, 0, 0}, false),
			}),
		zoneSpec(60, "Add Prime Sum", "Arithmetic · primality", "Checkpoint 50%",
			"AddPrimeSum(value int) int",
			"Sum every prime number less than or equal to a value.",
			"One integer.", "The prime sum, or 0 for non-positive input.",
			nil, intField, "AddPrimeSum(current.Payload.Value)", `0`,
			[]practiceCase{
				intCase("Negative", -2, 0),
				intCase("Zero", 0, 0),
				intCase("One", 1, 0),
				intCase("Five", 5, 10),
				intCase("Seven", 7, 17),
			}),
		zoneSpec(61, "Prime Factors", "Arithmetic · factorization", "Checkpoint 50%",
			"FPrime(value int) string",
			"Format a positive integer's prime factors in ascending order.",
			"One integer.", "Factors separated by *, or an empty string when value is below 2.",
			nil, intField, "FPrime(current.Payload.Value)", `""`,
			[]practiceCase{
				intOutputCase("Zero", 0, ""),
				intOutputCase("One", 1, ""),
				intOutputCase("Prime", 9539, "9539"),
				intOutputCase("Composite", 42, "2*3*7"),
				intOutputCase("Repeated factors", 225, "3*3*5*5"),
			}),
		zoneSpec(62, "Fifth And Skip", "Strings · grouped scanning", "Checkpoint 65%",
			"FifthAndSkip(value string) string",
			"Ignore spaces, keep groups of five characters, and discard each sixth character.",
			"One string.", "Five-character groups separated by spaces plus newline; special empty/short results apply.",
			[]string{"append", "len"}, stringField, "FifthAndSkip(current.Payload.Value)", `""`,
			stringCases(
				[]string{"", "1234", "12345", "abcdef", "abcdefghijklmnopqrstuwxyz", "This is a short sentence"},
				[]string{"\n", "Invalid Input\n", "12345\n", "abcde\n", "abcde ghijk mnopq stuwx z\n", "Thisi ashor sente ce\n"},
			)),
		zoneSpec(63, "Reverse Concat Alternate", "Slices · reverse interleaving", "Checkpoint 65%",
			"RevConcatAlternate(left []int, right []int) []int",
			"Interleave two slices from the end, giving the longer slice's excess values first.",
			"Two integer slices.", "A reverse-order alternation that begins with left when remaining lengths are equal.",
			[]string{"append", "len", "make"}, twoIntSlicesFields, "RevConcatAlternate(current.Payload.Left, current.Payload.Right)", `[]int{}`,
			[]practiceCase{
				twoSliceCase("Empty", []int{}, []int{}, []int{}),
				twoSliceCase("Equal", []int{1, 2, 3}, []int{4, 5, 6}, []int{3, 6, 2, 5, 1, 4}),
				twoSliceCase("Right longer", []int{1, 2, 3}, []int{4, 5, 6, 7, 8, 9}, []int{9, 8, 7, 3, 6, 2, 5, 1, 4}),
				twoSliceCase("Left longer", []int{1, 2, 3, 9, 8}, []int{4, 5}, []int{8, 9, 3, 2, 5, 1, 4}),
				twoSliceCase("One empty", []int{1, 2, 3}, []int{}, []int{3, 2, 1}),
			}),
		zoneSpec(64, "Not Decimal", "Strings · decimal parsing", "Checkpoint 65%",
			"NotDecimal(value string) string",
			"Remove a valid decimal point and insignificant fractional zeros.",
			"A decimal-looking string.", "The checkpoint conversion followed by newline; invalid input is returned unchanged plus newline.",
			[]string{"len"}, stringField, "NotDecimal(current.Payload.Value)", `"\n"`,
			stringCases(
				[]string{"", "0.1", "174.2", "1.20525856", "-0.0f00d00", "-19.525856", "1952"},
				[]string{"\n", "1\n", "1742\n", "120525856\n", "-0.0f00d00\n", "-19525856\n", "1952\n"},
			)),
		zoneSpec(65, "Slice", "Slices · variadic bounds", "Checkpoint 65%",
			"Slice(values []string, bounds ...int) []string",
			"Return a subsection of a string slice using positive or negative bounds.",
			"A string slice and zero, one, or two integer bounds.", "The selected slice, following the checkpoint's negative-index and invalid-range rules.",
			[]string{"len"}, sliceBoundsFields, "Slice(current.Payload.Values, current.Payload.Bounds...)", `nil`,
			[]practiceCase{
				sliceCase("No bounds", []string{"a", "b"}, []int{}, []string{"a", "b"}),
				sliceCase("Start", []string{"coding", "algorithm", "ascii", "package", "golang"}, []int{1}, []string{"algorithm", "ascii", "package", "golang"}),
				sliceCase("Range", []string{"coding", "algorithm", "ascii", "package", "golang"}, []int{2, 4}, []string{"ascii", "package"}),
				sliceCase("Negative", []string{"coding", "algorithm", "ascii", "package", "golang"}, []int{-2, -1}, []string{"package"}),
				sliceNilCase("Invalid order", []string{"coding", "algorithm", "ascii"}, []int{2, 0}),
			}),
	}

	levels := make([]Level, 0, len(specs))
	for _, spec := range specs {
		levels = append(levels, compilePracticeExercise(sourceZone01, spec))
	}
	return levels
}

func zoneSpec(
	id int,
	title, topic, difficulty, signature, objective, input, output string,
	builtins []string,
	inputFields, call, zero string,
	cases []practiceCase,
) practiceSpec {
	return practiceSpec{
		id: id, title: title, topic: topic, difficulty: difficulty,
		signature: signature, objective: objective, input: input, output: output,
		constraints: []string{
			"Keep the required function name, parameters, and return type.",
			"Match the stated edge cases and output exactly.",
			"Use loops, conditionals, and basic Go values rather than hidden dependencies.",
		},
		hints: []string{
			"Start with the smallest example, then handle the boundary cases.",
			"Track the current position or partial result explicitly.",
			"Compare your return value with the exact examples, including whitespace.",
		},
		pitfalls: []string{
			"Printing when the adapted browser exercise requires a return value",
			"Changing the required function signature",
			"Missing empty, zero, or boundary input",
		},
		builtins: builtins, starter: functionStarter(signature, zero),
		inputFields: inputFields, call: call, cases: cases,
	}
}

func zonePrintSpec(
	id int,
	title, topic, difficulty, signature, objective, input, output string,
	builtins []string,
	inputFields, call, zero string,
	cases []practiceCase,
) practiceSpec {
	return printedPracticeSpec(zoneSpec(
		id, title, topic, difficulty, signature, objective, input, output,
		builtins, inputFields, call, zero, cases,
	))
}

func functionStarter(signature, zero string) string {
	return fmt.Sprintf("func %s {\n\t// TODO: implement the checkpoint behavior.\n\treturn %s\n}\n", signature, zero)
}

func stringCases(inputs, outputs []string) []practiceCase {
	result := make([]practiceCase, 0, len(inputs))
	for index := range inputs {
		result = append(result, pc(
			fmt.Sprintf("Case %d", index+1),
			jsonString(inputs[index]),
			jsonString(outputs[index]),
			map[string]any{"value": inputs[index]},
		))
	}
	return result
}

func boolStringCases(inputs []string, outputs []bool) []practiceCase {
	result := make([]practiceCase, 0, len(inputs))
	for index := range inputs {
		result = append(result, pc(
			fmt.Sprintf("Case %d", index+1),
			jsonString(inputs[index]),
			fmt.Sprint(outputs[index]),
			map[string]any{"value": inputs[index]},
		))
	}
	return result
}

func intStringCases(inputs []string, outputs []int) []practiceCase {
	result := make([]practiceCase, 0, len(inputs))
	for index := range inputs {
		result = append(result, pc(
			fmt.Sprintf("Case %d", index+1),
			jsonString(inputs[index]),
			fmt.Sprint(outputs[index]),
			map[string]any{"value": inputs[index]},
		))
	}
	return result
}

func runeCase(name, value string, target rune, expected int) practiceCase {
	return pc(name, fmt.Sprintf("%s, %q", jsonString(value), target), fmt.Sprint(expected), map[string]any{"value": value, "target": target})
}

func intOutputCase(name string, value int, expected string) practiceCase {
	return pc(name, fmt.Sprint(value), jsonString(expected), map[string]any{"value": value})
}

func intCase(name string, value, expected int) practiceCase {
	return pc(name, fmt.Sprint(value), fmt.Sprint(expected), map[string]any{"value": value})
}

func intBaseCase(name string, value, base, expected int) practiceCase {
	return pc(name, fmt.Sprintf("%d, %d", value, base), fmt.Sprint(expected), map[string]any{"value": value, "base": base})
}

func replaceCase(name, value, old, replacement, expected string) practiceCase {
	input := fmt.Sprintf("%s, %s, %s", jsonString(value), jsonString(old), jsonString(replacement))
	return pc(name, input, jsonString(expected), map[string]any{"value": value, "old": old, "replacement": replacement})
}

func uintCase(name string, left, right, expected uint) practiceCase {
	return pc(name, fmt.Sprintf("%d, %d", left, right), fmt.Sprint(expected), map[string]any{"left": left, "right": right})
}

func fromToCase(name string, from, to int, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%d, %d", from, to), jsonString(expected), map[string]any{"from": from, "to": to})
}

func twoStringIntCase(name, left, right string, expected int) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(left), jsonString(right)), fmt.Sprint(expected), map[string]any{"left": left, "right": right})
}

func memoryCase(name string, value [10]byte, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%v", value), jsonString(expected), map[string]any{"value": value})
}

func twoSliceCase(name string, left, right, expected []int) practiceCase {
	return pc(name, fmt.Sprintf("%v, %v", left, right), jsonString(expected), map[string]any{"left": left, "right": right})
}

func stringSizeCase(name, value string, size int, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %d", jsonString(value), size), jsonString(expected), map[string]any{"value": value, "size": size})
}

func hiddenCase(name, needle, haystack string, expected bool) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(needle), jsonString(haystack)), fmt.Sprint(expected), map[string]any{"needle": needle, "haystack": haystack})
}

func hiddenStringCase(name, needle, haystack, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(needle), jsonString(haystack)), jsonString(expected), map[string]any{"needle": needle, "haystack": haystack})
}

func twoStringOutputCase(name, left, right, expected string) practiceCase {
	return pc(name, fmt.Sprintf("%s, %s", jsonString(left), jsonString(right)), jsonString(expected), map[string]any{"left": left, "right": right})
}

func chunkCase(name string, values []int, size int, expected [][]int) practiceCase {
	return pc(name, fmt.Sprintf("%v, %d", values, size), jsonString(expected), map[string]any{"values": values, "size": size})
}

func uintSliceCase(name string, values []uint, expected bool) practiceCase {
	return pc(name, fmt.Sprintf("%v", values), fmt.Sprint(expected), map[string]any{"values": values})
}

func sliceCase(name string, values []string, bounds []int, expected []string) practiceCase {
	return pc(name, fmt.Sprintf("%v, %v", values, bounds), jsonString(expected), map[string]any{"values": values, "bounds": bounds})
}

func sliceNilCase(name string, values []string, bounds []int) practiceCase {
	return pc(name, fmt.Sprintf("%v, %v", values, bounds), `null`, map[string]any{"values": values, "bounds": bounds})
}

func reverseComboExpected() string {
	result := ""
	for first := 9; first >= 2; first-- {
		for second := first - 1; second >= 1; second-- {
			for third := second - 1; third >= 0; third-- {
				if result != "" {
					result += ", "
				}
				result += fmt.Sprintf("%d%d%d", first, second, third)
			}
		}
	}
	return result + "\n"
}

const (
	noFields           = ""
	widthHeightFields  = "Width int `json:\"width\"`\n\tHeight int `json:\"height\"`"
	stringRuneFields   = "Value string `json:\"value\"`\n\tTarget rune `json:\"target\"`"
	intBaseFields      = "Value int `json:\"value\"`\n\tBase int `json:\"base\"`"
	replaceFields      = "Value string `json:\"value\"`\n\tOld string `json:\"old\"`\n\tReplacement string `json:\"replacement\"`"
	twoUintsFields     = "Left uint `json:\"left\"`\n\tRight uint `json:\"right\"`"
	fromToFields       = "From int `json:\"from\"`\n\tTo int `json:\"to\"`"
	twoStringsFields   = "Left string `json:\"left\"`\n\tRight string `json:\"right\"`"
	bytesField         = "Value [10]byte `json:\"value\"`"
	twoIntSlicesFields = "Left []int `json:\"left\"`\n\tRight []int `json:\"right\"`"
	stringSizeFields   = "Value string `json:\"value\"`\n\tSize int `json:\"size\"`"
	needleFields       = "Needle string `json:\"needle\"`\n\tHaystack string `json:\"haystack\"`"
	intsSizeFields     = "Values []int `json:\"values\"`\n\tSize int `json:\"size\"`"
	uintsField         = "Values []uint `json:\"values\"`"
	sliceBoundsFields  = "Values []string `json:\"values\"`\n\tBounds []int `json:\"bounds\"`"
)
