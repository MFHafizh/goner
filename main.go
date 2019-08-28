package main

import (
	"encoding/json"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var source string
var issue sonarIssues
var fset *token.FileSet

func main() {
	flag.Parse()
	source = flag.Args()[0]
	/*
		vendor := regexp.MustCompile(`([\\/])?vendor([\\/])?`)
		var packages []string
		path := "/Users/nakama/go/src/github.com/tokopedia/feeds/..."
		_, _ = getPackagePath(path, vendor)
		fmt.Println(packages)
	*/
	src, err := ioutil.ReadFile(source)
	issue = sonarIssues{[]sonarIssue{}}
	if err != nil {
		panic(err)
	}
	fset = token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	//ast.Print(fset, f)
	ast.Inspect(f, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.BinaryExpr:
			checkStrConcate(n)
		case *ast.CallExpr:
			checkStrFormat(n)
		}
		return true
	})
	res, err := json.MarshalIndent(issue, "", "\t")
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat("/result"); os.IsNotExist(err) {
		os.Mkdir("."+string(filepath.Separator)+"/result", os.ModePerm)
	}
	err = ioutil.WriteFile("result/report.json", res, 0644)
	if err != nil {
		panic(err)
	}
}

func checkStrConcate(n ast.Node) {
	patterns := regexp.MustCompile(`(?)(SELECT|DELETE|INSERT|UPDATE|INTO|FROM|WHERE) `)
	if node, ok := n.(*ast.BinaryExpr); ok {
		if start, ok := node.X.(*ast.BasicLit); ok {
			if str, e := strconv.Unquote(start.Value); e == nil {
				if !patterns.MatchString(str) {
					return
				} else if node.Op.String() == "+" {
					i := sonarIssue{
						EngineID: "custom scanner",
						RuleID:   "100",
						PrimaryLocation: location{
							Message:   "SQL Queries Using String Concate",
							FilePath:  source,
							TextRange: textRange{StartLine: fset.Position(node.Pos()).Line, EndLine: fset.Position(node.End()).Line},
						},
						Type:          "VULNERABILITY",
						Severity:      "MAJOR",
						EffortMinutes: 5,
					}
					issue.SonarIssues = append(issue.SonarIssues, i)
				}
			}
		}
	}
}
func checkStrFormat(n ast.Node) {
	patterns := regexp.MustCompile(`(?)(SELECT|DELETE|INSERT|UPDATE|INTO|FROM|WHERE) `)
	if node, ok := n.(*ast.CallExpr); ok {
		if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
			if inStrFormat(sel.Sel.Name) {
				if arg, ok := node.Args[0].(*ast.BasicLit); ok {
					if patterns.MatchString(arg.Value) {
						i := sonarIssue{
							EngineID: "custom scanner",
							RuleID:   "101",
							PrimaryLocation: location{
								Message:   "SQL Queries Using String Format",
								FilePath:  source,
								TextRange: textRange{StartLine: fset.Position(node.Pos()).Line, EndLine: fset.Position(node.End()).Line},
							},
							Type:          "VULNERABILITY",
							Severity:      "MAJOR",
							EffortMinutes: 5,
						}
						issue.SonarIssues = append(issue.SonarIssues, i)
					}
				}
			}
		}
	}
}

func inStrFormat(s string) bool {
	strFormating := []string{"fmt", "Sprint", "Sprintf", "Sprintln", "Fprintf"}
	for _, str := range strFormating {
		if s == str {
			return true
		}
	}
	return false
}
func process(packagesPath ...string) {

}

func getPackagePath(root string, excludedPath *regexp.Regexp) ([]string, error) {
	if strings.HasSuffix(root, "...") {
		root = root[0 : len(root)-3]
	} else {
		return []string{root}, nil
	}
	paths := map[string]bool{}
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			path = filepath.Dir(path)
			if excludedPath != nil && excludedPath.MatchString(path) {
				return nil
			}
			paths[path] = true
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	result := []string{}
	for path := range paths {
		result = append(result, path)
	}
	return result, nil
}
