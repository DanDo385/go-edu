package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	var scope string
	var dryRun bool
	var verbose bool

	flag.StringVar(&scope, "scope", "all", `Reset scope: "all", "minis", or "geth"`)
	flag.BoolVar(&dryRun, "dry-run", false, "Print planned changes, but do not write files")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	roots, err := scopeRoots(scope)
	must(err)

	var changed int
	for _, root := range roots {
		if verbose {
			fmt.Printf("Scanning %s/\n", root)
		}
		n, err := resetRoot(root, dryRun, verbose)
		must(err)
		changed += n
	}

	if dryRun {
		fmt.Printf("Dry run: would update %d exercise.go file(s)\n", changed)
		return
	}
	fmt.Printf("Updated %d exercise.go file(s)\n", changed)
}

func scopeRoots(scope string) ([]string, error) {
	switch scope {
	case "all":
		return []string{"geth", "minis"}, nil
	case "geth":
		return []string{"geth"}, nil
	case "minis":
		return []string{"minis"}, nil
	default:
		return nil, fmt.Errorf("invalid -scope %q (expected all|minis|geth)", scope)
	}
}

func resetRoot(root string, dryRun bool, verbose bool) (int, error) {
	var changed int
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(path) != "exercise.go" {
			return nil
		}

		out, ok, err := generateExerciseFromExercise(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if !ok {
			return nil
		}

		prev, readErr := os.ReadFile(path)
		if readErr == nil && bytes.Equal(bytes.TrimSpace(prev), bytes.TrimSpace(out)) {
			return nil
		}

		changed++
		if verbose || dryRun {
			fmt.Printf("%s\n", path)
		}
		if dryRun {
			return nil
		}
		if err := os.WriteFile(path, out, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		return nil
	})
	return changed, err
}

func generateExerciseFromExercise(exercisePath string) ([]byte, bool, error) {
	fset := token.NewFileSet()
	src, err := parser.ParseFile(fset, exercisePath, nil, parser.ParseComments)
	if err != nil {
		return nil, false, err
	}

	// If this file is already "instructions only" (no decls besides imports),
	// leave it as-is. Some exercises intentionally keep their instructions in
	// comments inside exercise.go (e.g. build-tags exercises).
	if hasNoCodeDecls(src.Decls) {
		return nil, false, nil
	}

	// Map import name (as used in selectors) -> import spec.
	importByName := map[string]*ast.ImportSpec{}
	for _, imp := range src.Imports {
		name := importName(imp)
		importByName[name] = imp
	}

	decls := exerciseDecls(src.Decls)
	usedImportNames := collectUsedImportNames(decls)
	importSpecs := selectImportSpecs(importByName, usedImportNames)

	outFile := &ast.File{
		Name: ast.NewIdent(src.Name.Name),
	}

	if len(importSpecs) > 0 {
		outFile.Decls = append(outFile.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: importSpecs,
		})
	}

	for _, d := range decls {
		outFile.Decls = append(outFile.Decls, stubDecl(d))
	}

	var buf bytes.Buffer
	buf.WriteString("//go:build !solution && !reference\n\n")
	buf.WriteString(todoHeader())
	buf.WriteByte('\n')

	formatted, err := formatNode(fset, outFile)
	if err != nil {
		return nil, false, err
	}
	buf.Write(formatted)
	if !bytes.HasSuffix(buf.Bytes(), []byte("\n")) {
		buf.WriteByte('\n')
	}
	return buf.Bytes(), true, nil
}

func hasNoCodeDecls(decls []ast.Decl) bool {
	for _, d := range decls {
		switch dd := d.(type) {
		case *ast.GenDecl:
			if dd.Tok != token.IMPORT {
				return false
			}
		case *ast.FuncDecl:
			_ = dd
			return false
		}
	}
	return true
}

func exerciseDecls(decls []ast.Decl) []ast.Decl {
	out := make([]ast.Decl, 0, len(decls))
	for _, d := range decls {
		if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			continue
		}
		out = append(out, d)
	}
	return out
}

func stubDecl(d ast.Decl) ast.Decl {
	switch dd := d.(type) {
	case *ast.FuncDecl:
		out := shallowCopyFuncDecl(dd)
		out.Doc = todoFuncDoc(dd)
		out.Body = &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun:  ast.NewIdent("panic"),
						Args: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: `"TODO: implement"`}},
					},
				},
			},
		}
		return out
	case *ast.GenDecl:
		out := shallowCopyGenDecl(dd)
		out.Doc = nil
		// Strip per-spec comments too.
		for _, s := range out.Specs {
			switch ss := s.(type) {
			case *ast.TypeSpec:
				ss.Doc = nil
				ss.Comment = nil
			case *ast.ValueSpec:
				ss.Doc = nil
				ss.Comment = nil
			}
		}
		return out
	default:
		return d
	}
}

func shallowCopyFuncDecl(in *ast.FuncDecl) *ast.FuncDecl {
	out := *in
	return &out
}

func shallowCopyGenDecl(in *ast.GenDecl) *ast.GenDecl {
	out := *in
	return &out
}

func todoHeader() string {
	return strings.Join([]string{
		"// TODO:",
		"// - Read the tests in exercise_test.go to understand expected behavior.",
		"// - Implement the exported API in this file.",
		"// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).",
	}, "\n")
}

func todoFuncDoc(fn *ast.FuncDecl) *ast.CommentGroup {
	name := "<unnamed>"
	if fn.Name != nil {
		name = fn.Name.Name
	}
	text := fmt.Sprintf("// TODO: implement %s.", name)
	return &ast.CommentGroup{List: []*ast.Comment{{Text: text}}}
}

func importName(spec *ast.ImportSpec) string {
	if spec.Name != nil {
		return spec.Name.Name
	}
	path := strings.Trim(spec.Path.Value, `"`)
	base := filepath.Base(path)
	return base
}

func collectUsedImportNames(decls []ast.Decl) map[string]struct{} {
	used := map[string]struct{}{}

	addSelectorPkgs := func(n ast.Node) {
		ast.Inspect(n, func(nn ast.Node) bool {
			sel, ok := nn.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok || id == nil {
				return true
			}
			used[id.Name] = struct{}{}
			return true
		})
	}

	for _, d := range decls {
		switch dd := d.(type) {
		case *ast.GenDecl:
			// For exported const/var/type declarations, we need imports used in types and values.
			addSelectorPkgs(dd)
		case *ast.FuncDecl:
			// Only consider receiver + signature, not the original function body.
			if dd.Recv != nil {
				addSelectorPkgs(dd.Recv)
			}
			if dd.Type != nil {
				addSelectorPkgs(dd.Type)
			}
		}
	}
	return used
}

func selectImportSpecs(importByName map[string]*ast.ImportSpec, used map[string]struct{}) []ast.Spec {
	var names []string
	for n := range used {
		if _, ok := importByName[n]; ok {
			names = append(names, n)
		}
	}
	sort.Strings(names)

	out := make([]ast.Spec, 0, len(names))
	for _, n := range names {
		imp := importByName[n]
		// Copy import spec so we don't mutate the parsed AST.
		cp := *imp
		cp.Doc = nil
		cp.Comment = nil
		out = append(out, &cp)
	}
	return out
}

func formatNode(fset *token.FileSet, file *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
