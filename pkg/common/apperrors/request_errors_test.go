package apperrors

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"unicode"
)

/*
This test file ensures that the application’s error categories stay in lockstep with the ErrorCode enum defined in Structs.sol. 
It parses the Solidity enum directly from the contract source and verifies both the order and the count against the Go-side 
category definitions. If anyone adds, removes, or reorders an error code in either place without updating the other, the test 
fails—giving us a fast safety net against drifting contract/go mappings, since the pipeline check will also fail.
*/

// ===== helpers =====

func camelToUpperSnake(s string) string {
	var out []rune
	var prevLower bool
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && (prevLower || (i+1 < len(s) && unicode.IsLower(rune(s[i+1])))) {
				out = append(out, '_')
			}
			out = append(out, unicode.ToUpper(r))
			prevLower = false
		} else if unicode.IsDigit(r) {
			if i > 0 && out[len(out)-1] != '_' {
				out = append(out, '_')
			}
			out = append(out, r)
			prevLower = false
		} else if r == '_' {
			if len(out) == 0 || out[len(out)-1] != '_' {
				out = append(out, '_')
			}
			prevLower = false
		} else {
			out = append(out, unicode.ToUpper(r))
			prevLower = true
		}
	}
	return string(out)
}

func stripPrefix(s, pref string) string {
	if strings.HasPrefix(s, pref) {
		return s[len(pref):]
	}
	return s
}

// ===== parse Go: const ( category = iota ... ) =====

func readGoCategoryOrder(t *testing.T, goFilePath string) []string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, goFilePath, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse Go file %s: %v", goFilePath, err)
	}

	var names []string

	ast.Inspect(file, func(n ast.Node) bool {
		gen, ok := n.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			return true
		}

		hasCategoryIota := false

		for i, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			if i == 0 {
				typeIsCategory := false
				if ident, ok := vs.Type.(*ast.Ident); ok && ident.Name == "category" {
					typeIsCategory = true
				}
				valueHasIota := false
				for _, v := range vs.Values {
					if id, ok := v.(*ast.Ident); ok && id.Name == "iota" {
						valueHasIota = true
						break
					}
				}
				hasCategoryIota = typeIsCategory && valueHasIota
			}

			if hasCategoryIota {
				for _, n := range vs.Names {
					names = append(names, n.Name)
				}
			}
		}

		return !hasCategoryIota
	})

	if len(names) == 0 {
		t.Fatalf("didn't find const block for 'category = iota' in %s", goFilePath)
	}
	return names
}

// ===== parse Solidity: enum ErrorCode { ... } =====

var enumRe = regexp.MustCompile(`enum\s+ErrorCode\s*\{([^}]*)\}`)

func readSolidityEnumOrder(t *testing.T, solPath string) []string {
	t.Helper()
	b, err := os.ReadFile(solPath)
	if err != nil {
		t.Fatalf("read Solidity %s: %v", solPath, err)
	}
	m := enumRe.FindSubmatch(b)
	if m == nil {
		t.Fatalf("didn't find 'enum ErrorCode' in %s", solPath)
	}
	body := string(m[1])

	lines := strings.Split(body, ",")
	var out []string
	for _, piece := range lines {
		line := piece
		if i := strings.Index(line, "//"); i >= 0 {
			line = line[:i]
		}
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		out = append(out, name)
	}
	if len(out) == 0 {
		t.Fatalf("enum ErrorCode empty at %s", solPath)
	}
	return out
}

// ===== test =====

func Test_Category_Aligned_With_Solidity_ErrorCode(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	root := filepath.Dir(thisFile)

	goFile := filepath.Join(root, "request_errors.go")
	solFile := filepath.Clean(filepath.Join(root, "../../../contracts/contracts/Structs.sol"))

	goNames := readGoCategoryOrder(t, goFile)       // e.g. ["categoryNoError", "categoryUnknown", ...]
	solNames := readSolidityEnumOrder(t, solFile)   // e.g. ["NO_ERROR", "UNKNOWN", ...]

	normGo := make([]string, 0, len(goNames))
	for _, gn := range goNames {
		base := stripPrefix(gn, "category")
		if base == gn {
			t.Fatalf("const '%s' doesn't start with 'category'", gn)
		}
		norm := camelToUpperSnake(base)
		norm = strings.Trim(norm, "_")
		normGo = append(normGo, norm)
	}

	if len(normGo) != len(solNames) {
		t.Fatalf("different sizes: Go(%d)=%v vs Solidity(%d)=%v", len(normGo), normGo, len(solNames), solNames)
	}

	for i := range solNames {
		if normGo[i] != solNames[i] {
			t.Fatalf("misaligned at index %d: Go=%s, Solidity=%s\nGo Order: %v\nSolidity Order: %v",
				i, normGo[i], solNames[i], normGo, solNames)
		}
	}
}