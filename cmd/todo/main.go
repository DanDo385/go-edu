package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `todo: reset exercise.go files back to TODO form

This regenerates each exercise.go from the corresponding reference solution file.

Usage:
  todo <path> [<path> ...]

Examples:
  todo geth
  todo minis
  todo geth/01-stack
  todo minis/30-build-tags-conditional-compilation/internal/buildtagsconditionalcompilation/exercise.go
`)
	}
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	var hadErr bool
	for _, p := range paths {
		if err := resetPath(p); err != nil {
			hadErr = true
			_, _ = fmt.Fprintf(os.Stderr, "todo: %v\n", err)
		}
	}
	if hadErr {
		os.Exit(1)
	}
}

func resetPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return resetDir(path)
	}
	return resetFile(path)
}

func resetDir(dir string) error {
	// Map exercise.go path -> chosen solution file path.
	// Prefer solution.reference.go when both exist.
	chosen := map[string]string{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		base := filepath.Base(p)
		if base != "solution.reference.go" && base != "solution_no_err.reference.go" {
			return nil
		}
		ex := exercisePathFromSolution(p)
		if ex == "" {
			return nil
		}
		if base == "solution.reference.go" {
			chosen[ex] = p
			return nil
		}
		if _, ok := chosen[ex]; !ok {
			chosen[ex] = p
		}
		return nil
	})
	if err != nil {
		return err
	}

	var exerciseFiles []string
	for ex := range chosen {
		exerciseFiles = append(exerciseFiles, ex)
	}
	sort.Strings(exerciseFiles)
	if len(exerciseFiles) == 0 {
		return fmt.Errorf("no reference solution files found under %q", dir)
	}

	var errs []error
	for _, ex := range exerciseFiles {
		sol := chosen[ex]
		if err := generateExerciseFromSolution(sol, ex); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", ex, err))
			continue
		}
		fmt.Printf("reset %s\n", ex)
	}
	return errors.Join(errs...)
}

func resetFile(path string) error {
	base := filepath.Base(path)
	switch base {
	case "solution.reference.go", "solution_no_err.reference.go":
		ex := exercisePathFromSolution(path)
		if ex == "" {
			return fmt.Errorf("cannot map %q to exercise.go", path)
		}
		if err := generateExerciseFromSolution(path, ex); err != nil {
			return err
		}
		fmt.Printf("reset %s\n", ex)
		return nil
	case "exercise.go":
		baseDir := strings.TrimSuffix(path, "exercise.go")
		sol := baseDir + "solution.reference.go"
		if _, err := os.Stat(sol); err != nil {
			alt := baseDir + "solution_no_err.reference.go"
			if _, err2 := os.Stat(alt); err2 != nil {
				return fmt.Errorf("exercise.go given but no reference solution file found next to it (tried %q and %q)", sol, alt)
			}
			sol = alt
		}
		if err := generateExerciseFromSolution(sol, path); err != nil {
			return err
		}
		fmt.Printf("reset %s\n", path)
		return nil
	default:
		return fmt.Errorf("unsupported file %q (expected exercise.go or a reference solution file)", path)
	}
}

func exercisePathFromSolution(solutionPath string) string {
	if strings.HasSuffix(solutionPath, "solution.reference.go") {
		return strings.TrimSuffix(solutionPath, "solution.reference.go") + "exercise.go"
	}
	if strings.HasSuffix(solutionPath, "solution_no_err.reference.go") {
		return strings.TrimSuffix(solutionPath, "solution_no_err.reference.go") + "exercise.go"
	}
	return ""
}

func generateExerciseFromSolution(solutionFile, exerciseFile string) error {
	src, err := os.ReadFile(solutionFile)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, solutionFile, src, parser.ParseComments)
	if err != nil {
		return err
	}

	usedPkgs := collectUsedPkgIdents(f)
	keptImports, err := filterImports(solutionFile, fset, f.Imports, usedPkgs)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	out.WriteString("//go:build !solution && !reference\n\n")
	out.WriteString("package ")
	out.WriteString(f.Name.Name)
	out.WriteString("\n\n")

	if len(keptImports) > 0 {
		out.WriteString("import (\n")
		for _, imp := range keptImports {
			var b bytes.Buffer
			if err := printer.Fprint(&b, fset, imp); err != nil {
				return err
			}
			out.WriteString("\t")
			out.WriteString(strings.TrimSpace(b.String()))
			out.WriteString("\n")
		}
		out.WriteString(")\n\n")
	}

	if top := minimalTopBlockComment(f); strings.TrimSpace(top) != "" {
		out.WriteString(top)
		out.WriteString("\n\n")
	}

	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gen.Tok != token.CONST && gen.Tok != token.TYPE && gen.Tok != token.VAR {
			continue
		}
		var b bytes.Buffer
		if err := printer.Fprint(&b, fset, gen); err != nil {
			return err
		}
		out.WriteString(strings.TrimSpace(b.String()))
		out.WriteString("\n\n")
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Type == nil || fn.Body == nil {
			continue
		}

		out.WriteString("// ")
		out.WriteString(fn.Name.Name)
		out.WriteString(" - TODO: implement this function\n")

		sig := extractFuncSignature(fset, src, fn)
		if strings.TrimSpace(sig) == "" {
			return fmt.Errorf("failed to extract function signature for %s", fn.Name.Name)
		}
		out.WriteString(strings.TrimSpace(sig))
		out.WriteString(" {\n")

		out.WriteString("\t// TODO: Implement this function\n")
		out.WriteString("\t// Refer to solution.reference.go for the complete implementation with detailed explanations\n")

		if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
			zeroDecls, ret := zeroReturnForResults(fset, fn.Type.Results)
			for _, line := range zeroDecls {
				out.WriteString("\t")
				out.WriteString(line)
				out.WriteString("\n")
			}
			out.WriteString("\t")
			out.WriteString(ret)
			out.WriteString("\n")
		}

		out.WriteString("}\n\n")
	}

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		return fmt.Errorf("format generated exercise: %w", err)
	}
	return os.WriteFile(exerciseFile, formatted, 0o644)
}

func extractFuncSignature(fset *token.FileSet, src []byte, fn *ast.FuncDecl) string {
	if fn == nil || fn.Body == nil {
		return ""
	}
	start := fset.Position(fn.Pos()).Offset
	lbrace := fset.Position(fn.Body.Lbrace).Offset
	if start < 0 || lbrace < 0 || start >= len(src) || lbrace > len(src) || start >= lbrace {
		return ""
	}
	return string(src[start:lbrace])
}

func minimalTopBlockComment(f *ast.File) string {
	var first *ast.CommentGroup
	for _, cg := range f.Comments {
		if cg == nil || len(cg.List) == 0 {
			continue
		}
		if strings.HasPrefix(cg.List[0].Text, "/*") {
			first = cg
			break
		}
	}
	if first == nil {
		return ""
	}

	raw := first.Text()
	lines := strings.Split(raw, "\n")
	var out []string
	out = append(out, "/*")
	keepNextBullets := false
	bulletRe := regexp.MustCompile(`^\d+\.`)

	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" || trim == "/*" || trim == "*/" {
			continue
		}
		if strings.HasPrefix(trim, "Problem:") ||
			strings.HasPrefix(trim, "Requirements:") ||
			strings.HasPrefix(trim, "Constraints:") ||
			strings.HasPrefix(trim, "Algorithm:") ||
			strings.HasPrefix(trim, "Time/Space") {
			out = append(out, line)
			keepNextBullets = true
			continue
		}
		if keepNextBullets {
			if strings.HasPrefix(strings.TrimLeft(trim, " \t"), "-") || bulletRe.MatchString(trim) {
				out = append(out, line)
				continue
			}
			keepNextBullets = false
		}
		if strings.HasPrefix(trim, "Problem:") {
			out = append(out, line)
		}
	}
	out = append(out, "*/")
	return strings.Join(out, "\n")
}

func collectUsedPkgIdents(f *ast.File) map[string]bool {
	used := make(map[string]bool)

	addFromExpr := func(expr ast.Expr) {
		if expr == nil {
			return
		}
		ast.Inspect(expr, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if ok {
				used[id.Name] = true
			}
			return true
		})
	}

	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gen.Tok != token.CONST && gen.Tok != token.TYPE && gen.Tok != token.VAR {
			continue
		}
		ast.Inspect(gen, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if ok {
				used[id.Name] = true
			}
			return true
		})
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Type == nil {
			continue
		}
		if fn.Recv != nil {
			for _, field := range fn.Recv.List {
				addFromExpr(field.Type)
			}
		}
		if fn.Type.Params != nil {
			for _, field := range fn.Type.Params.List {
				addFromExpr(field.Type)
			}
		}
		if fn.Type.Results != nil {
			for _, field := range fn.Type.Results.List {
				addFromExpr(field.Type)
			}
		}
	}

	return used
}

var goListNameCache = map[string]string{}

func filterImports(solutionFile string, fset *token.FileSet, imports []*ast.ImportSpec, usedPkgs map[string]bool) ([]*ast.ImportSpec, error) {
	var kept []*ast.ImportSpec
	for _, imp := range imports {
		if imp == nil || imp.Path == nil {
			continue
		}
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid import path in %q: %w", solutionFile, err)
		}

		if imp.Name != nil && (imp.Name.Name == "_" || imp.Name.Name == ".") {
			kept = append(kept, imp)
			continue
		}

		local := ""
		if imp.Name != nil && imp.Name.Name != "" {
			local = imp.Name.Name
		} else {
			name, err := goListPackageName(path)
			if err != nil {
				return nil, fmt.Errorf("go list import %q: %w", path, err)
			}
			local = name
		}

		if usedPkgs[local] {
			kept = append(kept, imp)
		}
	}
	return kept, nil
}

func goListPackageName(importPath string) (string, error) {
	if v, ok := goListNameCache[importPath]; ok {
		return v, nil
	}
	cmd := exec.Command("go", "list", "-f", "{{.Name}}", importPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
	}
	name := strings.TrimSpace(string(out))
	if name == "" {
		return "", fmt.Errorf("empty package name for %q", importPath)
	}
	goListNameCache[importPath] = name
	return name, nil
}

func zeroReturnForResults(fset *token.FileSet, results *ast.FieldList) (zeroDeclLines []string, returnLine string) {
	var names []string
	zeroIdx := 0
	for _, field := range results.List {
		t := typeString(fset, field.Type)
		// A single result field can declare multiple named results:
		//   func f() (vals []int, min, max int)
		// In that case we need one zero var per name.
		n := 1
		if len(field.Names) > 0 {
			n = len(field.Names)
		}
		for j := 0; j < n; j++ {
			name := fmt.Sprintf("zero%d", zeroIdx)
			zeroIdx++
			zeroDeclLines = append(zeroDeclLines, fmt.Sprintf("var %s %s", name, t))
			names = append(names, name)
		}
	}
	return zeroDeclLines, "return " + strings.Join(names, ", ")
}

func typeString(fset *token.FileSet, expr ast.Expr) string {
	var b bytes.Buffer
	_ = printer.Fprint(&b, fset, expr)
	return strings.TrimSpace(b.String())
}

