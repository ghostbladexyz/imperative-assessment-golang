package assessment

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"strconv"
	"strings"
)

var goBuiltins = map[string]struct{}{
	"append": {}, "cap": {}, "clear": {}, "close": {}, "complex": {},
	"copy": {}, "delete": {}, "imag": {}, "len": {}, "make": {},
	"max": {}, "min": {}, "new": {}, "panic": {}, "print": {},
	"println": {}, "real": {}, "recover": {},
}

// ValidateSourcePolicy rejects imports and Go built-in calls outside the
// exercise's explicit allowlist before untrusted code reaches a compiler.
func ValidateSourcePolicy(level Level, source string) error {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "solution.go", source, 0)
	if err != nil {
		return nil
	}

	allowedPackages := toSet(level.Instructions.AllowedPackages)
	for _, spec := range file.Imports {
		path, unquoteErr := strconv.Unquote(spec.Path.Value)
		if unquoteErr != nil {
			continue
		}
		if _, allowed := allowedPackages[path]; !allowed {
			position := fileSet.Position(spec.Pos())
			return fmt.Errorf(
				"solution.go:%d:%d: package %q is not allowed in this exercise; allowed packages: %s",
				position.Line,
				position.Column,
				path,
				displayAllowlist(level.Instructions.AllowedPackages),
			)
		}
		if spec.Name != nil && (spec.Name.Name == "." || spec.Name.Name == "_") {
			position := fileSet.Position(spec.Pos())
			return fmt.Errorf(
				"solution.go:%d:%d: dot and blank imports are not allowed",
				position.Line,
				position.Column,
			)
		}
	}

	allowedBuiltins := toSet(level.Instructions.AllowedBuiltins)
	selectorNames := make(map[*ast.Ident]struct{})
	ast.Inspect(file, func(node ast.Node) bool {
		if selector, ok := node.(*ast.SelectorExpr); ok {
			selectorNames[selector.Sel] = struct{}{}
		}
		return true
	})
	var policyErr error
	ast.Inspect(file, func(node ast.Node) bool {
		if policyErr != nil {
			return false
		}
		identifier, ok := node.(*ast.Ident)
		if !ok {
			return true
		}
		if identifier.Obj != nil {
			return true
		}
		if _, selectorName := selectorNames[identifier]; selectorName {
			return true
		}
		if _, builtin := goBuiltins[identifier.Name]; !builtin {
			return true
		}
		if _, allowed := allowedBuiltins[identifier.Name]; allowed {
			return true
		}
		position := fileSet.Position(identifier.Pos())
		policyErr = fmt.Errorf(
			"solution.go:%d:%d: built-in %q is not allowed in this exercise; allowed built-ins: %s",
			position.Line,
			position.Column,
			identifier.Name,
			displayAllowlist(level.Instructions.AllowedBuiltins),
		)
		return false
	})
	return policyErr
}

func toSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func displayAllowlist(values []string) string {
	if len(values) == 0 {
		return "none"
	}
	copyOfValues := slices.Clone(values)
	slices.Sort(copyOfValues)
	return strings.Join(copyOfValues, ", ")
}
