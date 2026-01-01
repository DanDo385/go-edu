// reset-exercises.go - Resets exercise.go files to TODO state
//
// Usage:
//   go run scripts/reset-exercises.go [options]
//
// Options:
//   -target=all    Reset all exercises (default)
//   -target=minis  Reset only minis exercises
//   -target=geth   Reset only geth exercises
//   -dry-run       Show what would be changed without writing files

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	target := flag.String("target", "all", "Target directory: all, minis, or geth")
	dryRun := flag.Bool("dry-run", false, "Show what would be changed without writing files")
	flag.Parse()

	// Find workspace root
	workspaceRoot, err := findWorkspaceRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding workspace root: %v\n", err)
		os.Exit(1)
	}

	// Determine which directories to process
	var dirs []string
	switch *target {
	case "all":
		dirs = []string{
			filepath.Join(workspaceRoot, "minis"),
			filepath.Join(workspaceRoot, "geth"),
		}
	case "minis":
		dirs = []string{filepath.Join(workspaceRoot, "minis")}
	case "geth":
		dirs = []string{filepath.Join(workspaceRoot, "geth")}
	default:
		fmt.Fprintf(os.Stderr, "Invalid target: %s. Use 'all', 'minis', or 'geth'\n", *target)
		os.Exit(1)
	}

	// Find and process all solution.reference.go files
	var processed, failed int
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == "solution.reference.go" {
				exercisePath := filepath.Join(filepath.Dir(path), "exercise.go")
				if err := resetExercise(path, exercisePath, *dryRun); err != nil {
					fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", path, err)
					failed++
				} else {
					processed++
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking %s: %v\n", dir, err)
		}
	}

	fmt.Printf("Processed: %d, Failed: %d\n", processed, failed)
	if *dryRun {
		fmt.Println("(dry-run mode - no files were modified)")
	}
}

func findWorkspaceRoot() (string, error) {
	// Start from current directory and look for workspace markers
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// If we're in the scripts directory, go up one level
	if filepath.Base(dir) == "scripts" {
		dir = filepath.Dir(dir)
	}

	// Check if minis and geth directories exist
	minisPath := filepath.Join(dir, "minis")
	gethPath := filepath.Join(dir, "geth")
	if _, err := os.Stat(minisPath); err == nil {
		if _, err := os.Stat(gethPath); err == nil {
			return dir, nil
		}
	}

	// Try parent directory
	parentDir := filepath.Dir(dir)
	minisPath = filepath.Join(parentDir, "minis")
	gethPath = filepath.Join(parentDir, "geth")
	if _, err := os.Stat(minisPath); err == nil {
		if _, err := os.Stat(gethPath); err == nil {
			return parentDir, nil
		}
	}

	return "", fmt.Errorf("cannot find workspace root (directory containing minis/ and geth/)")
}

// resetExercise generates an exercise.go file from solution.reference.go
func resetExercise(solutionPath, exercisePath string, dryRun bool) error {
	// Parse the solution file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, solutionPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// Extract package name
	pkgName := node.Name.Name

	// Collect imports
	imports := collectImports(node)

	// Collect type declarations and function signatures
	types := collectTypes(node, fset)
	funcs := collectFunctions(node, fset)

	// Generate exercise file content
	content := generateExercise(pkgName, imports, types, funcs)

	// Format the content
	formatted, err := format.Source([]byte(content))
	if err != nil {
		// If formatting fails, write unformatted content
		formatted = []byte(content)
	}

	if dryRun {
		fmt.Printf("Would write: %s\n", exercisePath)
		return nil
	}

	// Write the exercise file
	if err := os.WriteFile(exercisePath, formatted, 0644); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	fmt.Printf("Reset: %s\n", exercisePath)
	return nil
}

func collectImports(node *ast.File) []string {
	var imports []string
	for _, imp := range node.Imports {
		importStr := ""
		if imp.Name != nil {
			importStr = imp.Name.Name + " "
		}
		importStr += imp.Path.Value
		imports = append(imports, importStr)
	}
	sort.Strings(imports)
	return imports
}

type typeInfo struct {
	name       string
	code       string
	comment    string
	fullDecl   string // Full type declaration including type params
	isExported bool
}

func collectTypes(node *ast.File, fset *token.FileSet) []typeInfo {
	var types []typeInfo

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			isExported := ast.IsExported(typeSpec.Name.Name)

			// Build full type declaration
			var buf bytes.Buffer
			buf.WriteString("type ")
			buf.WriteString(typeSpec.Name.Name)

			// Handle type parameters (generics)
			if typeSpec.TypeParams != nil {
				buf.WriteString("[")
				for i, field := range typeSpec.TypeParams.List {
					if i > 0 {
						buf.WriteString(", ")
					}
					for j, name := range field.Names {
						if j > 0 {
							buf.WriteString(", ")
						}
						buf.WriteString(name.Name)
					}
					buf.WriteString(" ")
					var typeBuf bytes.Buffer
					printer.Fprint(&typeBuf, fset, field.Type)
					buf.WriteString(typeBuf.String())
				}
				buf.WriteString("]")
			}

			buf.WriteString(" ")

			// Print the type itself
			var typeBuf bytes.Buffer
			printer.Fprint(&typeBuf, fset, typeSpec.Type)
			buf.WriteString(typeBuf.String())

			types = append(types, typeInfo{
				name:       typeSpec.Name.Name,
				fullDecl:   buf.String(),
				isExported: isExported,
			})
		}
	}
	return types
}

type funcInfo struct {
	name       string
	receiver   string
	typeParams string
	signature  string
	params     string
	returns    string
	docComment string
	todoItems  []string
}

func collectFunctions(node *ast.File, fset *token.FileSet) []funcInfo {
	var funcs []funcInfo

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Skip unexported functions (but keep methods)
		if funcDecl.Recv == nil && !ast.IsExported(funcDecl.Name.Name) {
			continue
		}

		fi := funcInfo{
			name: funcDecl.Name.Name,
		}

		// Get receiver if present
		if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
			var buf bytes.Buffer
			printer.Fprint(&buf, fset, funcDecl.Recv.List[0].Type)
			recvType := buf.String()
			recvName := "r"
			if len(funcDecl.Recv.List[0].Names) > 0 {
				recvName = funcDecl.Recv.List[0].Names[0].Name
			}
			fi.receiver = fmt.Sprintf("(%s %s)", recvName, recvType)
		}

		// Get type parameters (generics) if present
		if funcDecl.Type.TypeParams != nil && len(funcDecl.Type.TypeParams.List) > 0 {
			fi.typeParams = formatTypeParams(funcDecl.Type.TypeParams, fset)
		}

		// Get parameters - format as string
		fi.params = formatFieldList(funcDecl.Type.Params, fset)

		// Get return types
		fi.returns = formatFieldList(funcDecl.Type.Results, fset)

		// Build signature
		sig := "func "
		if fi.receiver != "" {
			sig += fi.receiver + " "
		}
		sig += fi.name
		if fi.typeParams != "" {
			sig += "[" + fi.typeParams + "]"
		}
		sig += "(" + fi.params + ")"
		if fi.returns != "" {
			// Multiple returns need parentheses
			if strings.Contains(fi.returns, ",") {
				sig += " (" + fi.returns + ")"
			} else {
				sig += " " + fi.returns
			}
		}
		fi.signature = sig

		// Extract TODO items from the function body (look for comments with STEP or step pattern)
		fi.todoItems = extractTodoItems(funcDecl, fset)

		funcs = append(funcs, fi)
	}

	return funcs
}

// formatTypeParams formats type parameters (generics) as a string
func formatTypeParams(fl *ast.FieldList, fset *token.FileSet) string {
	if fl == nil || len(fl.List) == 0 {
		return ""
	}

	var parts []string
	for _, field := range fl.List {
		var typeBuf bytes.Buffer
		printer.Fprint(&typeBuf, fset, field.Type)
		typeStr := typeBuf.String()

		if len(field.Names) == 0 {
			parts = append(parts, typeStr)
		} else {
			for _, name := range field.Names {
				parts = append(parts, name.Name+" "+typeStr)
			}
		}
	}
	return strings.Join(parts, ", ")
}

// formatFieldList converts an ast.FieldList to a comma-separated string
func formatFieldList(fl *ast.FieldList, fset *token.FileSet) string {
	if fl == nil || len(fl.List) == 0 {
		return ""
	}

	var parts []string
	for _, field := range fl.List {
		var typeBuf bytes.Buffer
		printer.Fprint(&typeBuf, fset, field.Type)
		typeStr := typeBuf.String()

		if len(field.Names) == 0 {
			// Unnamed parameter/return (just type)
			parts = append(parts, typeStr)
		} else {
			// Named parameters
			for _, name := range field.Names {
				parts = append(parts, name.Name+" "+typeStr)
			}
		}
	}
	return strings.Join(parts, ", ")
}

func extractTodoItems(funcDecl *ast.FuncDecl, fset *token.FileSet) []string {
	var todos []string

	if funcDecl.Body == nil {
		return todos
	}

	// Look for comments that mention STEP or numbered steps
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		return true
	})

	// Extract high-level steps from function body structure
	for _, stmt := range funcDecl.Body.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			// Look for validation patterns
			if isValidationCheck(s) {
				todos = append(todos, "Validate inputs")
			}
		case *ast.AssignStmt:
			// Look for key assignments
			for _, lhs := range s.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if isImportantVar(ident.Name) {
						todos = append(todos, fmt.Sprintf("Initialize %s", ident.Name))
					}
				}
			}
		case *ast.ForStmt, *ast.RangeStmt:
			todos = append(todos, "Process items in loop")
		case *ast.ReturnStmt:
			todos = append(todos, "Return result")
		}
	}

	// Remove duplicates and limit
	return uniqueTodos(todos)
}

func isValidationCheck(stmt *ast.IfStmt) bool {
	// Check if this looks like a validation (nil check, empty check, etc.)
	if binExpr, ok := stmt.Cond.(*ast.BinaryExpr); ok {
		if binExpr.Op == token.EQL || binExpr.Op == token.NEQ {
			if ident, ok := binExpr.Y.(*ast.Ident); ok {
				return ident.Name == "nil"
			}
		}
	}
	return false
}

func isImportantVar(name string) bool {
	importantVars := []string{"result", "client", "ctx", "config", "cfg", "data", "resp", "response"}
	for _, v := range importantVars {
		if strings.EqualFold(name, v) {
			return true
		}
	}
	return false
}

func uniqueTodos(todos []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, t := range todos {
		if !seen[t] {
			seen[t] = true
			result = append(result, t)
		}
	}
	// Limit to 5 items
	if len(result) > 5 {
		result = result[:5]
	}
	return result
}

func generateExercise(pkgName string, imports []string, types []typeInfo, funcs []funcInfo) string {
	var buf bytes.Buffer

	// Build tag
	buf.WriteString("//go:build !solution && !reference\n\n")

	// Package declaration
	buf.WriteString(fmt.Sprintf("package %s\n\n", pkgName))

	// Imports
	if len(imports) > 0 {
		if len(imports) == 1 {
			buf.WriteString(fmt.Sprintf("import %s\n\n", imports[0]))
		} else {
			buf.WriteString("import (\n")
			for _, imp := range imports {
				buf.WriteString(fmt.Sprintf("\t%s\n", imp))
			}
			buf.WriteString(")\n\n")
		}
	}

	// Types (preserve from solution)
	for _, t := range types {
		buf.WriteString(t.fullDecl + "\n\n")
	}

	// Functions with TODO placeholders
	for _, f := range funcs {
		// Write function doc comment
		buf.WriteString(fmt.Sprintf("// %s implements the exercise.\n", f.name))
		buf.WriteString("//\n")
		buf.WriteString("// TODO: Implement this function\n")

		// Write function signature with panic placeholder
		buf.WriteString(f.signature + " {\n")
		buf.WriteString("\t// TODO: Implement\n")

		// Generate appropriate return statement
		returnStmt := generateReturnStmt(f.returns)
		if returnStmt != "" {
			buf.WriteString(fmt.Sprintf("\t%s\n", returnStmt))
		}

		buf.WriteString("}\n\n")
	}

	return buf.String()
}

func generateReturnStmt(returns string) string {
	if returns == "" {
		return ""
	}

	// Parse the return types
	returns = strings.TrimSpace(returns)
	returns = strings.TrimPrefix(returns, "(")
	returns = strings.TrimSuffix(returns, ")")

	if returns == "" {
		return ""
	}

	// Split by comma, handling nested types
	parts := splitReturnTypes(returns)

	var zeroVals []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Remove parameter name if present (e.g., "err error" -> "error")
		// But be careful not to split types with brackets like *Cache[K, V]
		// Only split if space comes before any [ or ( characters
		if idx := strings.Index(part, " "); idx != -1 {
			// Check if there are any brackets before the space
			bracketIdx := strings.IndexAny(part, "[(")
			if bracketIdx == -1 || idx < bracketIdx {
				// The space is before any brackets, so it's a parameter name
				part = part[idx+1:]
			}
		}
		zeroVals = append(zeroVals, zeroValue(part))
	}

	return "return " + strings.Join(zeroVals, ", ")
}

func splitReturnTypes(s string) []string {
	var parts []string
	depth := 0
	start := 0

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case '[':
			depth++
		case ']':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}

func zeroValue(typeName string) string {
	typeName = strings.TrimSpace(typeName)

	// Pointer types (including generic pointers like *Cache[K, V])
	if strings.HasPrefix(typeName, "*") {
		return "nil"
	}

	// Slice types
	if strings.HasPrefix(typeName, "[]") {
		return "nil"
	}

	// Map types
	if strings.HasPrefix(typeName, "map[") {
		return "nil"
	}

	// Channel types
	if strings.HasPrefix(typeName, "chan ") || strings.HasPrefix(typeName, "<-chan") {
		return "nil"
	}

	// Interface types
	if typeName == "error" || typeName == "interface{}" || typeName == "any" {
		return "nil"
	}

	// Function types
	if strings.HasPrefix(typeName, "func") {
		return "nil"
	}

	// Type parameters (single letter or short uppercase identifiers like K, V, T, etc.)
	// These are common in generic functions
	if len(typeName) <= 2 && typeName[0] >= 'A' && typeName[0] <= 'Z' {
		// For type parameters, use a zero value pattern
		return "*new(" + typeName + ")"
	}

	// Common types
	switch typeName {
	case "int", "int8", "int16", "int32", "int64":
		return "0"
	case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return "0"
	case "float32", "float64":
		return "0"
	case "complex64", "complex128":
		return "0"
	case "bool":
		return "false"
	case "string":
		return `""`
	case "byte":
		return "0"
	case "rune":
		return "0"
	}

	// Common Ethereum/big types that are pointers
	if typeName == "common.Address" {
		return "common.Address{}"
	}
	if typeName == "common.Hash" {
		return "common.Hash{}"
	}
	if strings.HasSuffix(typeName, "Address") || strings.HasSuffix(typeName, "Hash") {
		return typeName + "{}"
	}

	// If it's a struct type (starts with uppercase or contains dot)
	if len(typeName) > 0 && (typeName[0] >= 'A' && typeName[0] <= 'Z') {
		return typeName + "{}"
	}
	if strings.Contains(typeName, ".") && !strings.HasPrefix(typeName, "*") {
		return typeName + "{}"
	}

	// Default to nil for unknown types (likely interfaces or pointers)
	return "nil"
}
