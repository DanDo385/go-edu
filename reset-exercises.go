package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  directory: 'geth', 'minis', or 'all'\n")
		os.Exit(1)
	}

	target := os.Args[1]
	var dirs []string

	switch target {
	case "all":
		dirs = []string{"geth", "minis"}
	case "geth", "minis":
		dirs = []string{target}
	default:
		fmt.Fprintf(os.Stderr, "Error: directory must be 'geth', 'minis', or 'all'\n")
		os.Exit(1)
	}

	count := 0
	for _, dir := range dirs {
		n, err := resetExercisesInDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", dir, err)
			os.Exit(1)
		}
		count += n
	}

	fmt.Printf("✓ Reset %d exercise.go files successfully\n", count)
}

func resetExercisesInDir(dir string) (int, error) {
	count := 0
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Look for solution.reference.go files
		if d.Name() == "solution.reference.go" {
			dir := filepath.Dir(path)
			exercisePath := filepath.Join(dir, "exercise.go")

			// Check if exercise.go exists
			if _, err := os.Stat(exercisePath); os.IsNotExist(err) {
				return nil // Skip if exercise.go doesn't exist
			}

			if err := resetExerciseFile(exercisePath, path); err != nil {
				return fmt.Errorf("failed to reset %s: %w", exercisePath, err)
			}

			fmt.Printf("Reset: %s\n", exercisePath)
			count++
		}

		return nil
	})
	return count, err
}

func resetExerciseFile(exercisePath, solutionPath string) error {
	// Read the solution file
	solutionContent, err := os.ReadFile(solutionPath)
	if err != nil {
		return err
	}

	// Parse the solution file
	fset := token.NewFileSet()
	solutionAST, err := parser.ParseFile(fset, solutionPath, solutionContent, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse solution file: %w", err)
	}

	// Extract package name
	packageName := solutionAST.Name.Name

	// Extract imports
	var importGroups []string
	var stdlibImports []string
	var thirdPartyImports []string

	for _, imp := range solutionAST.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		importStr := importPath
		if imp.Name != nil {
			importStr = fmt.Sprintf("%s %q", imp.Name.Name, importPath)
		} else {
			importStr = fmt.Sprintf("%q", importPath)
		}

		// Categorize imports
		if !strings.Contains(importPath, ".") {
			stdlibImports = append(stdlibImports, importStr)
		} else {
			thirdPartyImports = append(thirdPartyImports, importStr)
		}
	}

	// Combine imports with blank line between stdlib and third-party
	if len(stdlibImports) > 0 {
		importGroups = append(importGroups, strings.Join(stdlibImports, "\n\t"))
	}
	if len(thirdPartyImports) > 0 {
		importGroups = append(importGroups, strings.Join(thirdPartyImports, "\n\t"))
	}

	// Extract STEP comments from solution to create TODOs
	stepRegex := regexp.MustCompile(`(?m)^\s*//\s*STEP\s+\d+:\s*(.+)$`)
	matches := stepRegex.FindAllStringSubmatch(string(solutionContent), -1)
	var todos []string
	for _, match := range matches {
		if len(match) > 1 {
			stepDesc := strings.TrimSpace(match[1])
			// Remove trailing dashes and extra formatting
			stepDesc = regexp.MustCompile(`\s*-\s*.*$`).ReplaceAllString(stepDesc, "")
			stepDesc = strings.TrimSpace(stepDesc)
			if stepDesc != "" {
				todos = append(todos, fmt.Sprintf("\t// TODO: %s", stepDesc))
			}
		}
	}

	// If no STEP comments found, create generic TODOs based on function structure
	if len(todos) == 0 {
		todos = []string{
			"\t// TODO: Implement this function",
			"\t// TODO: Handle errors appropriately",
			"\t// TODO: Add necessary validations",
		}
	}

	// Extract function signatures
	var funcDecls []*ast.FuncDecl
	ast.Inspect(solutionAST, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name != "init" && !strings.HasPrefix(x.Name.Name, "Test") && !strings.HasPrefix(x.Name.Name, "Benchmark") {
				funcDecls = append(funcDecls, x)
			}
		}
		return true
	})

	// Generate the new exercise.go content
	var buf strings.Builder
	buf.WriteString("//go:build !solution && !reference\n\n")
	buf.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	if len(importGroups) > 0 {
		buf.WriteString("import (\n\t")
		buf.WriteString(strings.Join(importGroups, "\n\t"))
		buf.WriteString("\n)\n\n")
	}

	// Write function stubs with TODOs
	for i, fn := range funcDecls {
		if i > 0 {
			buf.WriteString("\n")
		}

		// Write function signature
		buf.WriteString("func ")
		if fn.Recv != nil {
			buf.WriteString("(")
			for j, field := range fn.Recv.List {
				if j > 0 {
					buf.WriteString(", ")
				}
				for k, name := range field.Names {
					if k > 0 {
						buf.WriteString(", ")
					}
					buf.WriteString(name.Name)
				}
				if len(field.Names) > 0 {
					buf.WriteString(" ")
				}
				buf.WriteString(exprToString(field.Type))
			}
			buf.WriteString(") ")
		}
		buf.WriteString(fn.Name.Name)
		buf.WriteString("(")
		if fn.Type.Params != nil {
			for j, param := range fn.Type.Params.List {
				if j > 0 {
					buf.WriteString(", ")
				}
				for k, name := range param.Names {
					if k > 0 {
						buf.WriteString(", ")
					}
					if name != nil {
						buf.WriteString(name.Name)
					}
				}
				if len(param.Names) > 0 {
					buf.WriteString(" ")
				}
				buf.WriteString(exprToString(param.Type))
			}
		}
		buf.WriteString(")")
		if fn.Type.Results != nil {
			buf.WriteString(" ")
			if len(fn.Type.Results.List) > 1 {
				buf.WriteString("(")
			}
			for j, result := range fn.Type.Results.List {
				if j > 0 {
					buf.WriteString(", ")
				}
				for k, name := range result.Names {
					if k > 0 {
						buf.WriteString(", ")
					}
					if name != nil {
						buf.WriteString(name.Name)
					}
				}
				if len(result.Names) > 0 {
					buf.WriteString(" ")
				}
				buf.WriteString(exprToString(result.Type))
			}
			if len(fn.Type.Results.List) > 1 {
				buf.WriteString(")")
			}
		}
		buf.WriteString(" {\n")

		// Add TODOs
		if len(todos) > 0 {
			// Use all TODOs for the first function, or distribute them
			todosToUse := todos
			if len(funcDecls) > 1 && i > 0 {
				// For subsequent functions, use a subset or generic TODOs
				todosToUse = []string{"\t// TODO: Implement this function"}
			}
			for _, todo := range todosToUse {
				buf.WriteString(todo)
				buf.WriteString("\n")
			}
		} else {
			buf.WriteString("\t// TODO: Implement this function\n")
		}

		buf.WriteString("\tpanic(\"not implemented\")\n")
		buf.WriteString("}\n")
	}

	// Write the file
	return os.WriteFile(exercisePath, []byte(buf.String()), 0644)
}

func exprToString(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return exprToString(x.X) + "." + x.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(x.X)
	case *ast.ArrayType:
		if x.Len == nil {
			return "[]" + exprToString(x.Elt)
		}
		return "[" + exprToString(x.Len) + "]" + exprToString(x.Elt)
	case *ast.MapType:
		return "map[" + exprToString(x.Key) + "]" + exprToString(x.Value)
	case *ast.ChanType:
		dir := ""
		switch x.Dir {
		case ast.SEND:
			dir = "chan<- "
		case ast.RECV:
			dir = "<-chan "
		default:
			dir = "chan "
		}
		return dir + exprToString(x.Value)
	case *ast.FuncType:
		var buf strings.Builder
		buf.WriteString("func(")
		if x.Params != nil {
			for i, param := range x.Params.List {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(exprToString(param.Type))
			}
		}
		buf.WriteString(")")
		if x.Results != nil {
			if len(x.Results.List) > 1 {
				buf.WriteString(" (")
			} else {
				buf.WriteString(" ")
			}
			for i, result := range x.Results.List {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(exprToString(result.Type))
			}
			if len(x.Results.List) > 1 {
				buf.WriteString(")")
			}
		}
		return buf.String()
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "..." + exprToString(x.Elt)
	case *ast.BasicLit:
		return x.Value
	default:
		return "interface{}"
	}
}
