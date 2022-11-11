package depslint

import (
	"bufio"
	"go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var flagRoot string

const doc = "depslint is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "depslint",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func init() {
	// ルートURLを明示指定しなくとも動くようにしたい。
	// 例えば os.Getwd などで取得したいところだが、 ./... 指定でvetを実行すると、カレントディレクトリも変化してしまう。
	Analyzer.Flags.StringVar(&flagRoot, "root", "", "root url")
}

func run(pass *analysis.Pass) (any, error) {
	if flagRoot == "" {
		return nil, nil
	}

	// 以下の処理は繰り返し実行されてしまっているが、go vet実行全体で1回だけ実行されるようにしたい。
	// initに処理を逃がすなど試みたが、 ./... 指定で実行するとinitも複数回実行される模様。
	rootPkg, err := GetRootPkgName(flagRoot)
	if err != nil {
		return nil, err
	}
	rules := ParsePuml(rootPkg, filepath.Join(flagRoot, "./.depslint.puml"))

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
	}

	pass.Report = analysisutil.ReportWithoutIgnore(pass)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ImportSpec:
			CheckImport(pass, n, rules)
		}
	})

	return nil, nil
}

func GetRootPkgName(rootPath string) (string, error) {
	modfilename := filepath.Join(rootPath, "./go.mod")
	data, err := ioutil.ReadFile(modfilename)
	if err != nil {
		return "", err
	}

	f, err := modfile.Parse(modfilename, data, nil)
	if err != nil {
		return "", err
	}

	return f.Module.Mod.Path, nil
}

func CheckImport(pass *analysis.Pass, n *ast.ImportSpec, rules []DependencyRule) {
	if IsDisabledLint(n) {
		return
	}

	from := pass.Pkg.Path()
	to := strings.Trim(n.Path.Value, "\"")
	for _, r := range rules {
		if r.From == to && r.To == from {
			pass.Reportf(n.Pos(), "\x1b[31mdependency violation:\x1b[0m %s -> %s", from, to)
		}
	}
}

func IsDisabledLint(n *ast.ImportSpec) bool {
	if n.Doc == nil {
		return false
	}

	r := regexp.MustCompile(`//\s?lint:ignore\s+(.*,|)depslint`)
	for _, com := range n.Doc.List {
		if r.MatchString(com.Text) {
			return true
		}
	}
	return false
}

type DependencyRule struct {
	From string
	To   string
}

func ParsePuml(rootPkg, path string) []DependencyRule {
	f, _ := os.Open(path)
	bu := bufio.NewReaderSize(f, 1024)

	var rules []DependencyRule
	scanner := bufio.NewScanner(bu)
	for scanner.Scan() {
		rule := ParseRule(rootPkg, scanner.Text())
		if rule != nil {
			rules = append(rules, *rule)
		}
	}

	return rules
}

func ParseRule(rootPkg, lineTxt string) *DependencyRule {
	r := regexp.MustCompile(`-+>|<-+`)
	token := r.Find([]byte(lineTxt))
	if token == nil {
		return nil
	}

	txts := strings.Split(lineTxt, string(token))
	for i, v := range txts {
		txts[i] = strings.Trim(v, " []")
	}

	isForward := strings.HasSuffix(string(token), ">")
	if isForward {
		return &DependencyRule{
			From: rootPkg + "/" + txts[0], To: rootPkg + "/" + txts[1],
		}
	} else {
		return &DependencyRule{
			From: rootPkg + "/" + txts[1], To: rootPkg + "/" + txts[0],
		}
	}
}
