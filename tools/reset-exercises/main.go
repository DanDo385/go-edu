package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	targetDir := flag.String("target", ".", "Directory to scan")
	flag.Parse()

	err := filepath.Walk(*targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, "solution.reference.go") {
			if err := processFile(path); err != nil {
				fmt.Printf("Failed to process %s: %v\n", path, err)
			} else {
				fmt.Printf("Reset %s\n", strings.Replace(path, "solution.reference.go", "exercise.go", 1))
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking path: %v\n", err)
		os.Exit(1)
	}
}

type Edit struct {
	Offset int
	Length int
	Text   string
}

func processFile(srcPath string) error {
	destPath := strings.Replace(srcPath, "solution.reference.go", "exercise.go", 1)
	
	contentBytes, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, srcPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var edits []Edit

	// 1. Comment Groups (Build Tags and Cleanup)
	for _, cg := range file.Comments {
		// Check for build tags first
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "//go:build") {
				start := fset.Position(c.Pos()).Offset
				end := fset.Position(c.End()).Offset
				edits = append(edits, Edit{
					Offset: start,
					Length: end - start,
					Text:   "//go:build !solution && !reference",
				})
			}
		}

		// Check if group is verbose
		isGroupVerbose := false
		for _, c := range cg.List {
			if isVerbose(c.Text) {
				isGroupVerbose = true
				break
			}
		}

		if isGroupVerbose {
			// Ensure we don't delete build tags if they happened to be in this group (unlikely but safe to check)
			// Actually, if we replaced build tag above, we have an edit for it.
			// If we now delete the whole group, we overwrite the build tag edit?
			// Or if we generate an edit that covers the whole group, it conflicts with build tag edit?
			// We should check if the group contains build tags.
			hasBuildTag := false
			for _, c := range cg.List {
				if strings.HasPrefix(c.Text, "//go:build") {
					hasBuildTag = true
					break
				}
			}

			if !hasBuildTag {
				start := fset.Position(cg.Pos()).Offset
				end := fset.Position(cg.End()).Offset
				edits = append(edits, Edit{
					Offset: start,
					Length: end - start,
					Text:   "",
				})
			}
		}
	}

	// 2. Package Doc (Simplify)
	if file.Doc != nil {
		start := fset.Position(file.Doc.Pos()).Offset
		end := fset.Position(file.Doc.End()).Offset
		
		origDoc := content[start:end]
		newDoc := simplifyDoc(origDoc)
		
		edits = append(edits, Edit{
			Offset: start,
			Length: end - start,
			Text:   newDoc,
		})
	}

	// 3. Function Bodies
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			start := fset.Position(fn.Body.Lbrace).Offset
			end := fset.Position(fn.Body.Rbrace).Offset + 1 // Include '}'
			
			bodySrc := content[start:end]
			newBody := generateBodyFromSource(bodySrc)
			
			edits = append(edits, Edit{
				Offset: start,
				Length: end - start,
				Text:   newBody,
			})
		}
	}

	// Filter overlapping edits (remove contained edits)
	var validEdits []Edit
	for i, e1 := range edits {
		isContained := false
		for j, e2 := range edits {
			if i == j {
				continue
			}
			if e1.Offset >= e2.Offset && e1.Offset+e1.Length <= e2.Offset+e2.Length {
                if e1.Offset == e2.Offset && e1.Length == e2.Length {
                    if j < i { isContained = true; break }
                } else {
				    isContained = true
				    break
                }
			}
		}
		if !isContained {
			validEdits = append(validEdits, e1)
		}
	}
	edits = validEdits

	// Apply edits from bottom to top
	sort.Slice(edits, func(i, j int) bool {
		return edits[i].Offset > edits[j].Offset
	})

	newContent := content
	for _, edit := range edits {
		if edit.Offset < 0 || edit.Offset+edit.Length > len(newContent) {
			continue 
		}
		newContent = newContent[:edit.Offset] + edit.Text + newContent[edit.Offset+edit.Length:]
	}

	return os.WriteFile(destPath, []byte(newContent), 0644)
}

func isVerbose(text string) bool {
	keywords := []string{"DEBUGGING", "Exercise", "Complexity", "Three-Input Iteration Table", "Alternatives & Trade-offs"}
	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func simplifyDoc(text string) string {
	isBlock := strings.HasPrefix(text, "/*")
	lines := strings.Split(text, "\n")
	var newLines []string
	count := 0
	
	for i, line := range lines {
		if isVerbose(line) {
			break
		}
		newLines = append(newLines, line)
		count++
		if count > 10 {
			if isBlock && i < len(lines)-1 {
				if !strings.HasSuffix(line, "*/") {
					newLines = append(newLines, "*/")
				}
			}
			break
		}
	}
	if isBlock {
		res := strings.Join(newLines, "\n")
		if !strings.HasSuffix(strings.TrimSpace(res), "*/") {
			res += "\n*/"
		}
		return res
	}
	return strings.Join(newLines, "\n")
}

func generateBodyFromSource(src string) string {
	if len(src) < 2 {
		return src
	}
	inner := src[1 : len(src)-1]
	lines := strings.Split(inner, "\n")
	var steps []string
	
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		
		if strings.Contains(trimmed, "STEP") {
			if strings.HasPrefix(trimmed, "//") {
				steps = append(steps, trimmed)
			}
		} else if strings.Contains(trimmed, "====") {
			if i+1 < len(lines) && !strings.Contains(lines[i+1], "====") {
				steps = append(steps, trimmed)
				i++
				if i < len(lines) {
					title := strings.TrimSpace(lines[i])
					if strings.HasPrefix(title, "//") {
						steps = append(steps, title)
					}
				}
			} else {
				if strings.HasPrefix(trimmed, "//") {
					steps = append(steps, trimmed)
				}
			}
		}
	}
	
	var blocks []string
	currentBlock := ""
	
	for _, s := range steps {
		if strings.Contains(s, "====") {
			if currentBlock != "" {
				blocks = append(blocks, currentBlock)
				currentBlock = ""
			}
			currentBlock += "\t" + s + "\n"
		} else {
			currentBlock += "\t" + s + "\n"
		}
	}
	if currentBlock != "" {
		blocks = append(blocks, currentBlock)
	}
	
	var sb strings.Builder
	sb.WriteString("{\n")
	
	if len(blocks) > 0 {
		for _, block := range blocks {
			sb.WriteString(block)
			sb.WriteString("\t// TODO: Implement\n\n")
		}
	} else {
		sb.WriteString("\t// TODO: Implement this function\n")
	}
	
	sb.WriteString("\tpanic(\"unimplemented\")\n")
	sb.WriteString("}")
	
	return sb.String()
}
