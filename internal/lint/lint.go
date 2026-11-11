// Package lint は、エラーにはならないけれど間違いのもとになる書き方を見つける (fixmarge lavi、法の悪魔)。
// 規則は docs/RFC.md の「法の悪魔」の10個。実行せずに構文木を調べ、確かにそうだと言えるときだけ指摘する。
package lint

import (
	"cmp"
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/modpath"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// 規則の名前。値は lexicon の骨格で、行末の # lavi: と fixmarge.mod の youlugive にもこの形で書く。
const (
	WasDa    = "wasda"
	DowGuil  = "dowguil"
	ShiDes   = "shides"
	EmBasy   = "embasy"
	AnxAn    = "anxan"
	SuZai    = "suzai"
	MehBeRac = "mehberac"
	JokeDan  = "jokedan"
	SAiMoku  = "saimoku"
	DesLee   = "deslee"
)

// Finding は指摘1つ。
type Finding struct {
	Rule   string
	File   string
	Pos    token.Pos
	Target string
}

// String は LaVI <規則> <ファイル>:<行>:<列> <対象> の形。
func (f Finding) String() string {
	s := lexicon.Spell("lavi") + " " + lexicon.Spell(f.Rule) + " " + f.File + ":" + f.Pos.String()
	if f.Target != "" {
		s += " " + f.Target
	}
	return s
}

// File はファイル1つを調べる。src は行末の # lavi: を読むために使う。allowed は fixmarge.mod の youlugive で許した規則。
func File(f *parser.File, src string, allowed []string) []Finding {
	l := newLinter(f)
	for _, w := range f.Warnings {
		l.report(DowGuil, w.Pos, w.From)
	}
	l.unused(f.Stmts)
	l.unreachable(f.Stmts)
	l.emptyBlocks(f.Stmts)
	l.needlessAbiryoku(f.Stmts)
	l.floatingPromises(f.Stmts)
	l.uncheckedEda(f.Stmts)
	l.comparisons(f.Stmts)
	l.hopelessShoubek(f.Stmts)

	silenced := silencedLines(src)
	var out []Finding
	for _, fd := range l.findings {
		if !slices.Contains(allowed, fd.Rule) && !slices.Contains(silenced[fd.Pos.Line], fd.Rule) {
			out = append(out, fd)
		}
	}
	slices.SortStableFunc(out, func(a, b Finding) int {
		return cmp.Or(cmp.Compare(a.Pos.Line, b.Pos.Line), cmp.Compare(a.Pos.Col, b.Pos.Col))
	})
	return out
}

// silencedLines は行末の # lavi: 規則 を読み、行の番号ごとに止める規則を返す。
func silencedLines(src string) map[int][]string {
	out := map[int][]string{}
	for i, line := range strings.Split(src, "\n") {
		_, rules, ok := strings.Cut(line, "# lavi:")
		if !ok {
			continue
		}
		out[i+1] = strings.FieldsFunc(rules, func(r rune) bool { return r == ' ' || r == ',' || r == '\t' || r == '\r' })
	}
	return out
}

type importedName struct{ path, name string }

type linter struct {
	file     string
	findings []Finding
	aliases  map[string]string       // "./util" aste util の util → 読み込み先
	named    map[string]importedName // "./util" no double の double → 読み込み先と中の名前
	funcs    map[string]*ast.FuncLit // weho で定義した関数
}

func newLinter(f *parser.File) *linter {
	l := &linter{file: f.Name, aliases: map[string]string{}, named: map[string]importedName{}, funcs: map[string]*ast.FuncLit{}}
	ast.InspectStmts(f.Stmts, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Import:
			if n.Alias != "" {
				l.aliases[n.Alias] = n.Path
			}
			for _, in := range n.Names {
				local := cmp.Or(in.Alias, in.Name)
				l.named[local] = importedName{n.Path, in.Name}
			}
		case *ast.FuncDecl:
			l.funcs[n.Fn.Name] = n.Fn
		}
		return true
	})
	return l
}

func (l *linter) report(rule string, pos token.Pos, target string) {
	l.findings = append(l.findings, Finding{Rule: rule, File: l.file, Pos: pos, Target: target})
}

// callee は呼び出す関数について、実行前に分かること。known が false なら分からない。
type callee struct {
	known    bool
	abiryoku bool
	itsver   bool
	name     string
}

// universe は読み込まずに使える組み込み。どれも入出力をしない。
var universe = []string{"nang", "sanand"}

func (l *linter) resolve(fn ast.Expr) callee {
	switch fn := fn.(type) {
	case *ast.FuncLit:
		return callee{known: true, abiryoku: fn.Abiryoku, itsver: fn.Itsver}
	case *ast.Ident:
		if f, ok := l.funcs[fn.Name]; ok {
			return callee{known: true, abiryoku: f.Abiryoku, itsver: f.Itsver, name: fn.Name}
		}
		if slices.Contains(universe, fn.Name) {
			return callee{known: true, name: fn.Name}
		}
		if in, ok := l.named[fn.Name]; ok {
			return stdCallee(in.path, in.name)
		}
	case *ast.Index:
		x, ok1 := fn.X.(*ast.Ident)
		key, ok2 := fn.Key.(*ast.Ident)
		if path, ok3 := l.aliases[name(x, ok1)]; ok1 && ok2 && ok3 {
			return stdCallee(path, key.Name)
		}
	}
	return callee{}
}

func name(id *ast.Ident, ok bool) string {
	if !ok {
		return ""
	}
	return id.Name
}

func stdCallee(path, name string) callee {
	pkg, ok := strings.CutPrefix(path, modpath.StdPrefix)
	if !ok {
		return callee{}
	}
	info, ok := interp.StdFunc(pkg, name)
	if !ok {
		return callee{}
	}
	return callee{known: true, abiryoku: info.Abiryoku, itsver: info.Itsver, name: name}
}

// isEda は呼び出しが coword no eda か varizama no eda (見つからないと -1 を返す) かを返す。
func (l *linter) isEda(x ast.Expr) bool {
	if p, ok := x.(*ast.Paren); ok {
		return l.isEda(p.X)
	}
	call, ok := x.(*ast.Call)
	if !ok {
		return false
	}
	var path, fn string
	switch f := call.Fn.(type) {
	case *ast.Ident:
		in := l.named[f.Name]
		path, fn = in.path, in.name
	case *ast.Index:
		x, ok1 := f.X.(*ast.Ident)
		key, ok2 := f.Key.(*ast.Ident)
		if !ok1 || !ok2 {
			return false
		}
		path, fn = l.aliases[x.Name], key.Name
	}
	return fn == "eda" && (path == ":coword" || path == ":varizama")
}
