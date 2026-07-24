package runner

import (
	"context"
	"fmt"
	"testing"
)

func TestZone01ReferenceSolutionsPass(t *testing.T) {
	solutions := zone01ReferenceSolutions()
	if len(solutions) != 44 {
		t.Fatalf("got %d reference solutions, want 44", len(solutions))
	}
	localRunner := New("go", 4, nil)
	for originalID := 22; originalID <= 65; originalID++ {
		originalID := originalID
		level := levelByOriginalTestID(t, fmt.Sprintf("l%d-01", originalID))
		t.Run(level.Title, func(t *testing.T) {
			result := localRunner.Run(context.Background(), level, solutions[originalID], nil)
			if !result.Passed {
				t.Fatalf(
					"exercise %d reference failed (%d/%d): compile=%q runtime=%q results=%#v",
					originalID,
					result.PassedCount,
					result.TotalCount,
					result.CompileError,
					result.RuntimeError,
					result.Results,
				)
			}
		})
	}
}

func zone01ReferenceSolutions() map[int]string {
	return map[int]string{
		22: `func OnlyA() string { return "a" }`,
		23: `func PrintIfNot(value string) string {
	if len(value) < 3 { return "G\n" }
	return "Invalid Input\n"
}`,
		24: `func PrintIf(value string) string {
	if len(value) == 0 || len(value) >= 3 { return "G\n" }
	return "Invalid Input\n"
}`,
		25: `func RectPerimeter(width, height int) int {
	if width < 0 || height < 0 { return -1 }
	return 2*width + 2*height
}`,
		26: `func CountChar(value string, target rune) int {
	count := 0
	for _, current := range value {
		if current == target { count++ }
	}
	return count
}`,
		27: `func CheckNumber(value string) bool {
	for _, current := range value {
		if current >= '0' && current <= '9' { return true }
	}
	return false
}`,
		28: `func RetainFirstHalf(value string) string {
	if len(value) < 2 { return value }
	return value[:len(value)/2]
}`,
		29: `func CountAlpha(value string) int {
	count := 0
	for _, current := range value {
		if (current >= 'a' && current <= 'z') || (current >= 'A' && current <= 'Z') { count++ }
	}
	return count
}`,
		30: `func FirstWord(value string) string {
	result := ""
	started := false
	for i := 0; i < len(value); i++ {
		if value[i] == ' ' {
			if started { break }
			continue
		}
		started = true
		result += string(value[i])
	}
	return result + "\n"
}`,
		31: `func LastWord(value string) string {
	result := ""
	current := ""
	for i := 0; i < len(value); i++ {
		if value[i] == ' ' {
			if current != "" { result = current; current = "" }
			continue
		}
		current += string(value[i])
	}
	if current != "" { result = current }
	return result + "\n"
}`,
		32: `func FishAndChips(value int) string {
	if value < 0 { return "error: number is negative" }
	if value%6 == 0 { return "fish and chips" }
	if value%2 == 0 { return "fish" }
	if value%3 == 0 { return "chips" }
	return "error: non divisible"
}`,
		33: `func DigitLen(value, base int) int {
	if base < 2 || base > 36 { return -1 }
	if value < 0 { value = -value }
	count := 0
	for value > 0 { count++; value /= base }
	return count
}`,
		34: `func SearchReplace(value, old, replacement string) string {
	if len(old) != 1 || len(replacement) != 1 { return value }
	result := ""
	for i := 0; i < len(value); i++ {
		if value[i] == old[0] { result += replacement } else { result += string(value[i]) }
	}
	return result
}`,
		35: `func RepeatAlpha(value string) string {
	result := ""
	for _, current := range value {
		count := 1
		if current >= 'a' && current <= 'z' { count = int(current-'a') + 1 }
		if current >= 'A' && current <= 'Z' { count = int(current-'A') + 1 }
		for repeat := 0; repeat < count; repeat++ { result += string(current) }
	}
	return result
}`,
		36: `func Gcd(left, right uint) uint {
	if left == 0 || right == 0 { return 0 }
	for right != 0 { left, right = right, left%right }
	return left
}`,
		37: `func CamelToSnakeCase(value string) string {
	if value == "" { return "" }
	for i := 0; i < len(value); i++ {
		current := value[i]
		if !((current >= 'a' && current <= 'z') || (current >= 'A' && current <= 'Z')) { return value }
		if current >= 'A' && current <= 'Z' {
			if i == len(value)-1 { return value }
			if i > 0 && value[i-1] >= 'A' && value[i-1] <= 'Z' { return value }
		}
	}
	result := ""
	for i := 0; i < len(value); i++ {
		if i > 0 && value[i] >= 'A' && value[i] <= 'Z' { result += "_" }
		result += string(value[i])
	}
	return result
}`,
		38: `func HashCode(value string) string {
	result := ""
	size := len(value)
	for _, current := range value {
		hashed := (int(current) + size) % 127
		if hashed < 33 { hashed += 33 }
		result += string(rune(hashed))
	}
	return result
}`,
		39: `func ThirdTimeIsACharm(value string) string {
	result := ""
	for i := 2; i < len(value); i += 3 { result += string(value[i]) }
	return result + "\n"
}`,
		40: `func FromTo(from, to int) string {
	if from < 0 || from > 99 || to < 0 || to > 99 { return "Invalid\n" }
	result := ""
	step := 1
	if from > to { step = -1 }
	for value := from; ; value += step {
		if value < 10 {
			result += "0" + string(rune('0'+value))
		} else {
			result += string(rune('0'+value/10)) + string(rune('0'+value%10))
		}
		if value == to { break }
		result += ", "
	}
	return result + "\n"
}`,
		41: `func IsCapitalized(value string) bool {
	if len(value) == 0 { return false }
	atStart := true
	for i := 0; i < len(value); i++ {
		if value[i] == ' ' { atStart = true; continue }
		if atStart {
			if value[i] >= 'a' && value[i] <= 'z' { return false }
			atStart = false
		}
	}
	return true
}`,
		42: `func FindPrevPrime(value int) int {
	for candidate := value; candidate >= 2; candidate-- {
		prime := true
		for divisor := 2; divisor*divisor <= candidate; divisor++ {
			if candidate%divisor == 0 { prime = false; break }
		}
		if prime { return candidate }
	}
	return 0
}`,
		43: `func Itoa(value int) string {
	if value == 0 { return "0" }
	negative := value < 0
	result := ""
	for value != 0 {
		digit := value % 10
		if digit < 0 { digit = -digit }
		result = string(rune('0'+digit)) + result
		value /= 10
	}
	if negative { result = "-" + result }
	return result
}`,
		44: `func CleanStr(value string) string {
	result := ""
	inWord := false
	for i := 0; i < len(value); i++ {
		if value[i] == ' ' || value[i] == '\t' { inWord = false; continue }
		if !inWord && result != "" { result += " " }
		result += string(value[i])
		inWord = true
	}
	return result + "\n"
}`,
		45: `func ExpandStr(value string) string {
	result := ""
	inWord := false
	for i := 0; i < len(value); i++ {
		if value[i] == ' ' || value[i] == '\t' { inWord = false; continue }
		if !inWord && result != "" { result += "   " }
		result += string(value[i])
		inWord = true
	}
	if result == "" { return "" }
	return result + "\n"
}`,
		46: `func WeAreUnique(left, right string) int {
	if left == "" && right == "" { return -1 }
	leftSet := map[rune]bool{}
	rightSet := map[rune]bool{}
	for _, value := range left { leftSet[value] = true }
	for _, value := range right { rightSet[value] = true }
	count := 0
	for value := range leftSet { if !rightSet[value] { count++ } }
	for value := range rightSet { if !leftSet[value] { count++ } }
	return count
}`,
		47: `import "fmt"
func ZipString(value string) string {
	result := ""
	count := 0
	var previous rune
	for _, current := range value {
		if count > 0 && current != previous {
			result += fmt.Sprintf("%d%c", count, previous)
			count = 0
		}
		previous = current
		count++
	}
	if count > 0 { result += fmt.Sprintf("%d%c", count, previous) }
	return result
}`,
		48: `import "fmt"
func PrintRevCombo() string {
	result := ""
	for first := 9; first >= 2; first-- {
		for second := first-1; second >= 1; second-- {
			for third := second-1; third >= 0; third-- {
				if result != "" { result += ", " }
				result += fmt.Sprintf("%d%d%d", first, second, third)
			}
		}
	}
	return result + "\n"
}`,
		49: `import "fmt"
func PrintMemory(value [10]byte) string {
	result := ""
	for i, current := range value {
		result += fmt.Sprintf("%02x", current)
		if i == 3 || i == 7 || i == 9 { result += "\n" } else { result += " " }
	}
	for _, current := range value {
		if current >= 32 && current <= 126 { result += string(current) } else { result += "." }
	}
	return result + "\n"
}`,
		50: `func ConcatSlice(left, right []int) []int {
	result := []int{}
	result = append(result, left...)
	result = append(result, right...)
	return result
}`,
		51: `func SaveAndMiss(value string, size int) string {
	if size <= 0 { return value }
	result := ""
	for start := 0; start < len(value); start += size*2 {
		end := start + size
		if end > len(value) { end = len(value) }
		result += value[start:end]
	}
	return result
}`,
		52: `func HiddenP(needle, haystack string) bool {
	index := 0
	for cursor := 0; cursor < len(haystack) && index < len(needle); cursor++ {
		if needle[index] == haystack[cursor] { index++ }
	}
	return index == len(needle)
}`,
		53: `func WdMatch(needle, haystack string) string {
	if needle == "" { return "" }
	index := 0
	for cursor := 0; cursor < len(haystack) && index < len(needle); cursor++ {
		if needle[index] == haystack[cursor] { index++ }
	}
	if index == len(needle) { return needle }
	return ""
}`,
		54: `func Inter(left, right string) string {
	result := ""
	seen := map[rune]bool{}
	available := map[rune]bool{}
	for _, value := range right { available[value] = true }
	for _, value := range left {
		if available[value] && !seen[value] { seen[value] = true; result += string(value) }
	}
	return result
}`,
		55: `func Union(left, right string) string {
	result := ""
	seen := map[rune]bool{}
	for _, text := range []string{left, right} {
		for _, value := range text {
			if !seen[value] { seen[value] = true; result += string(value) }
		}
	}
	return result
}`,
		56: `func ConcatAlternate(left, right []int) []int {
	if len(left) < len(right) { left, right = right, left }
	result := []int{}
	for index, value := range left {
		result = append(result, value)
		if index < len(right) { result = append(result, right[index]) }
	}
	return result
}`,
		57: `func Chunk(values []int, size int) [][]int {
	result := [][]int{}
	if size <= 0 { return result }
	for start := 0; start < len(values); start += size {
		end := start + size
		if end > len(values) { end = len(values) }
		part := []int{}
		part = append(part, values[start:end]...)
		result = append(result, part)
	}
	return result
}`,
		58: `func ReverseStrCap(value string) string {
	result := []byte(value)
	for i := 0; i < len(result); i++ {
		if result[i] >= 'A' && result[i] <= 'Z' { result[i] += 'a'-'A' }
		if result[i] != ' ' && (i == len(result)-1 || result[i+1] == ' ') {
			if result[i] >= 'a' && result[i] <= 'z' { result[i] -= 'a'-'A' }
		}
	}
	return string(result)
}`,
		59: `func CanJump(values []uint) bool {
	if len(values) == 0 { return false }
	position := 0
	for position < len(values)-1 {
		if values[position] == 0 { return false }
		position += int(values[position])
		if position >= len(values) { return false }
	}
	return true
}`,
		60: `func AddPrimeSum(value int) int {
	sum := 0
	for candidate := 2; candidate <= value; candidate++ {
		prime := true
		for divisor := 2; divisor*divisor <= candidate; divisor++ {
			if candidate%divisor == 0 { prime = false; break }
		}
		if prime { sum += candidate }
	}
	return sum
}`,
		61: `func FPrime(value int) string {
	if value < 2 { return "" }
	result := ""
	for divisor := 2; value > 1; divisor++ {
		for value%divisor == 0 {
			if result != "" { result += "*" }
			result += ItoaFactor(divisor)
			value /= divisor
		}
	}
	return result
}
func ItoaFactor(value int) string {
	result := ""
	for value > 0 { result = string(rune('0'+value%10))+result; value /= 10 }
	return result
}`,
		62: `func FifthAndSkip(value string) string {
	if value == "" { return "\n" }
	clean := []rune{}
	for _, current := range value { if current != ' ' { clean = append(clean, current) } }
	if len(clean) < 5 { return "Invalid Input\n" }
	result := ""
	for index := 0; index < len(clean); {
		for kept := 0; kept < 5 && index < len(clean); kept++ {
			result += string(clean[index]); index++
		}
		if index < len(clean) { index++ }
		if index < len(clean) { result += " " }
	}
	return result + "\n"
}`,
		63: `func RevConcatAlternate(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	leftIndex, rightIndex := len(left)-1, len(right)-1
	for leftIndex >= 0 || rightIndex >= 0 {
		if leftIndex > rightIndex { result = append(result, left[leftIndex]); leftIndex--; continue }
		if rightIndex > leftIndex { result = append(result, right[rightIndex]); rightIndex--; continue }
		if leftIndex >= 0 { result = append(result, left[leftIndex]); leftIndex-- }
		if rightIndex >= 0 { result = append(result, right[rightIndex]); rightIndex-- }
	}
	return result
}`,
		64: `func NotDecimal(value string) string {
	if value == "" { return "\n" }
	dot := -1
	for i := 0; i < len(value); i++ {
		if value[i] == '.' {
			if dot != -1 { return value + "\n" }
			dot = i
			continue
		}
		if value[i] == '-' && i == 0 { continue }
		if value[i] < '0' || value[i] > '9' { return value + "\n" }
	}
	if dot == -1 { return value + "\n" }
	before, after := value[:dot], value[dot+1:]
	nonZero := 0
	for nonZero < len(after) && after[nonZero] == '0' { nonZero++ }
	if nonZero == len(after) { return before + "\n" }
	after = after[nonZero:]
	if before == "0" { before = "" }
	if before == "-0" { before = "-" }
	return before + after + "\n"
}`,
		65: `func Slice(values []string, bounds ...int) []string {
	if len(bounds) == 0 { return values }
	first := bounds[0]
	if first < 0 { first = len(values)+first }
	if len(bounds) == 1 {
		if first < 0 { return values }
		if first > len(values) { first = len(values) }
		return values[first:]
	}
	second := bounds[1]
	if first < 0 { first = 0 }
	if first > len(values) { first = len(values) }
	if second < 0 { second = len(values)+second }
	if second < 0 { second = 0 }
	if second > len(values) { second = len(values) }
	if first > second { return nil }
	return values[first:second]
}`,
	}
}
