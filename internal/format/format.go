// Package format は下書きを整える (fixmarge beaukshi)。
// 字句の並びを1行ずつ並べ直し、字句のすき間からはコメントだけを拾う。字句そのものは書き換えない。
package format

import (
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// indentUnit は字下げ1段。仕様書とサンプルに合わせて空白2つ (Prettier の既定と同じ)。
const indentUnit = "  "

// Format は下書きを整えた結果を返す。整える前と後で字句の並びが変わったら、壊さないようエラーにする。
func Format(src, file string) (string, error) {
	tokens, _, err := lexer.Lex(src, file, lexer.Draft)
	if err != nil {
		return "", err
	}
	out, err := (&printer{src: src, keepGaps: true}).format(tokens, file)
	if err != nil {
		return "", err
	}
	if err := sameTokens(tokens, out, file); err != nil {
		return "", err
	}
	return out, nil
}

// Expand は呪いを展開した結果を返す (fixmarge beaukshi --curoi)。見るためのもので、コメントは残らない。
// 呪いの中で作った名前は、付け替えた名前 (tmp·1) のまま出す。
func Expand(src, file string) (string, error) {
	tokens, err := parser.Expand(src, file)
	if err != nil {
		return "", err
	}
	return (&printer{src: src}).format(tokens, file)
}

// sameTokens は、整えた結果を読み直した字句の並びが元と同じかを調べる。改行は続いた空行をまとめるので、続く改行は1つとして比べる。
func sameTokens(before []token.Token, out, file string) error {
	after, _, err := lexer.Lex(out, file, lexer.Draft)
	if err != nil {
		return err
	}
	a, b := significant(before), significant(after)
	if !slices.EqualFunc(a, b, func(x, y token.Token) bool { return x.Kind == y.Kind && x.Name == y.Name }) {
		return diag.New(diag.SaiVI, file, token.Pos{Line: 1, Col: 1}, "")
	}
	return nil
}

func significant(tokens []token.Token) []token.Token {
	var out []token.Token
	for _, t := range tokens {
		if t.Kind == token.Newline && (len(out) == 0 || out[len(out)-1].Kind == token.Newline) {
			continue
		}
		out = append(out, t)
	}
	for len(out) > 0 && (out[len(out)-1].Kind == token.Newline || out[len(out)-1].Kind == token.EOF) {
		out = out[:len(out)-1]
	}
	return out
}

type printer struct {
	src      string
	keepGaps bool // コメントと、括弧の中の改行を残すか

	b      strings.Builder
	depth  int  // ブロックの深さ
	blank  bool // 空行を1つためている
	output bool // 何か書いたか (ファイルの頭の空行を消すため)
	// tooDeep は字下げが ast.MaxDepth を超えた所の字句。深すぎる入れ子は、行ごとの字下げだけで出力が膨れ上がるので、そこで書くのをやめる
	tooDeep *token.Token
}

// open は1つの文の中で開いているもの (括弧かブロック) と、それを開いた物理行の番号。
type open struct {
	block bool
	line  int
}

// format は tokens を整えて返す。入れ子が深すぎれば Bavellabion。
func (p *printer) format(tokens []token.Token, file string) (string, error) {
	out := p.print(tokens)
	if t := p.tooDeep; t != nil {
		return "", diag.New(diag.Bavellabion, file, t.Pos, p.text(*t))
	}
	return out, nil
}

func (p *printer) print(tokens []token.Token) string {
	var line []token.Token
	prevEnd := 0
	for _, t := range tokens {
		if p.tooDeep != nil {
			break
		}
		if t.Kind == token.Newline || t.Kind == token.EOF {
			if len(line) > 0 {
				p.logicalLine(line, p.gap(line[len(line)-1].End, t.Start))
			} else {
				p.betweenLines(p.gap(prevEnd, t.Start))
			}
			line = nil
			prevEnd = t.End
			if t.Kind == token.EOF {
				break
			}
			continue
		}
		if len(line) == 0 {
			p.betweenLines(p.gap(prevEnd, t.Start))
		}
		line = append(line, t)
	}
	if p.output {
		p.b.WriteByte('\n')
	}
	return p.b.String()
}

// betweenLines は文と文の間のすき間を書く。字句解析は続く改行を1つにまとめるので、空行とコメントだけの行はここに入っている。
// 空行は続けても1つにまとめる。
func (p *printer) betweenLines(g string) {
	rows := strings.Split(g, "\n")
	for k, row := range rows {
		switch c := commentIn(row); {
		case c != "":
			p.startLine(p.depth)
			p.b.WriteString(c)
		case k < len(rows)-1:
			p.blank = p.output
		}
	}
}

func (p *printer) gap(from, to int) string {
	if !p.keepGaps || from > to {
		return ""
	}
	return p.src[from:to]
}

// logicalLine は改行の字句で区切った1つの文を書く。trailing は文の後ろのすき間 (行末のコメントがあればここ)。
func (p *printer) logicalLine(tokens []token.Token, trailing string) {
	comment := commentIn(trailing)
	indent := p.depth
	if dedentsLine(tokens[0]) {
		// 開いたブロックがないのに siva や finend が来ても、字下げは 0 より浅くしない
		indent = max(0, indent-1)
	}
	if indent > ast.MaxDepth {
		p.tooDeep = &tokens[0]
		return
	}
	p.startLine(indent)

	var stack []open
	physLine := 0
	for i, t := range tokens {
		if i > 0 {
			g := p.gap(tokens[i-1].End, t.Start)
			if strings.Contains(g, "\n") {
				p.breakInside(g, tokens, i, stack, indent)
				if p.tooDeep != nil {
					return
				}
				physLine++
			} else if needsSpace(tokens, i) {
				p.b.WriteByte(' ')
			}
		}
		p.b.WriteString(p.text(t))
		stack = track(stack, tokens, i, physLine)
	}
	if comment != "" {
		p.b.WriteString("  " + comment)
	}
	// 文の終わりで開いたままのブロック (mof x の行など) は、次の行からの深さになる
	for _, o := range stack {
		if o.block {
			p.depth++
		}
	}
	p.depth = max(0, p.depth-closedOutside(tokens))
}

func (p *printer) startLine(indent int) {
	if p.output {
		p.b.WriteByte('\n')
		if p.blank {
			p.b.WriteByte('\n')
		}
	}
	p.blank = false
	p.output = true
	p.b.WriteString(strings.Repeat(indentUnit, indent))
}

// breakInside は括弧の中の改行を書く。すき間のコメントは残し、空行は消す。
// 次の行の字下げは、まだ開いているものを開いた物理行の数だけ深くする (Prettier と同じく、1行で何個開いても1段)。
func (p *printer) breakInside(g string, tokens []token.Token, next int, stack []open, indent int) {
	rows := strings.Split(g, "\n")
	if c := commentIn(rows[0]); c != "" {
		p.b.WriteString("  " + c)
	}
	// 行の頭で続けて閉じるもの (finend や ) が並ぶ所) は、閉じた後の深さで書く
	view := stack
	for j := next; j < len(tokens) && isCloser(tokens[j]) && len(view) > 0; j++ {
		if j > next && strings.Contains(p.gap(tokens[j-1].End, tokens[j].Start), "\n") {
			break
		}
		view = view[:len(view)-1]
	}
	levels := indent + distinctLines(view)
	if levels > ast.MaxDepth {
		p.tooDeep = &tokens[next]
		return
	}
	for _, row := range rows[1 : len(rows)-1] {
		if c := commentIn(row); c != "" {
			p.b.WriteString("\n" + strings.Repeat(indentUnit, levels) + c)
		}
	}
	p.b.WriteString("\n" + strings.Repeat(indentUnit, levels))
}

func distinctLines(stack []open) int {
	n, last := 0, -1
	for _, o := range stack {
		if o.line != last {
			n++
			last = o.line
		}
	}
	return n
}

// track は i 番目の字句で開いたものを積み、閉じたものを下ろす。
func track(stack []open, tokens []token.Token, i, physLine int) []open {
	t := tokens[i]
	switch {
	case t.IsOp("(") || t.IsOp("[") || t.IsOp("{"):
		return append(stack, open{line: physLine})
	case isBlockOpener(tokens, i):
		return append(stack, open{block: true, line: physLine})
	case (t.IsOp(")") || t.IsOp("]") || t.IsOp("}")) && len(stack) > 0:
		return stack[:len(stack)-1]
	case t.IsKeyword("finend") && len(stack) > 0 && stack[len(stack)-1].block:
		return stack[:len(stack)-1]
	}
	return stack
}

// closedOutside は、この文の中で開いていないブロックを閉じた finend の数 (前の行で開いたブロックを閉じる)。
func closedOutside(tokens []token.Token) int {
	var stack []open
	closed := 0
	for i, t := range tokens {
		if t.IsKeyword("finend") && !slices.ContainsFunc(stack, func(o open) bool { return o.block }) {
			closed++
			continue
		}
		stack = track(stack, tokens, i, 0)
	}
	return closed
}

// isBlockOpener は、その語でブロックが始まるかを返す。siva mof は前の mof のブロックの続き。
func isBlockOpener(tokens []token.Token, i int) bool {
	t := tokens[i]
	if t.Kind != token.Keyword {
		return false
	}
	switch t.Name {
	case "weho", "teitance", "sentky", "persemon", "curoi", "azagara":
		return true
	case "mof":
		return i == 0 || !tokens[i-1].IsKeyword("siva")
	}
	return false
}

// dedentsLine は、行の頭に来たら1段浅く書く語 (ブロックの終わりと、ブロックの区切り)。
func dedentsLine(t token.Token) bool {
	return t.Kind == token.Keyword && slices.Contains([]string{"finend", "siva", "youlugive", "ennote"}, t.Name)
}

func isCloser(t token.Token) bool {
	return t.IsOp(")") || t.IsOp("]") || t.IsOp("}") || t.IsKeyword("finend")
}

// commentIn はすき間の中のコメントを、右の空白を除いて返す。
func commentIn(gap string) string {
	i := strings.IndexByte(gap, '#')
	if i < 0 {
		return ""
	}
	return strings.TrimRight(gap[i:], " \t\r")
}

// text は字句の書かれた形。展開で作った括弧と付け替えた名前は、ソースにないので名前から書く。
func (p *printer) text(t token.Token) string {
	if t.Start == t.End || strings.Contains(t.Name, "·") {
		return t.Name
	}
	return p.src[t.Start:t.End]
}

// valueKeywords は値になるキーワード。この後ろの - は引き算。
var valueKeywords = []string{"tadact", "ulieso", "vanita"}

func endsValue(t token.Token) bool {
	switch t.Kind {
	case token.Ident, token.Int, token.Float, token.String:
		return true
	case token.Keyword:
		return slices.Contains(valueKeywords, t.Name)
	}
	return t.IsOp(")") || t.IsOp("]") || t.IsOp("}")
}

func startsValue(t token.Token) bool {
	switch t.Kind {
	case token.Ident, token.Int, token.Float, token.String:
		return true
	case token.Keyword:
		return slices.Contains(valueKeywords, t.Name)
	}
	return t.IsOp("(") || t.IsOp("[") || t.IsOp("{") || t.IsOp("-")
}

// isUnary は i 番目の記号が、-1 のような値の前の記号かを返す。
func isUnary(tokens []token.Token, i int) bool {
	t := tokens[i]
	if !t.IsOp("-") && !t.IsOp("+") && !t.IsOp("~") {
		return false
	}
	return i == 0 || !endsValue(tokens[i-1])
}

// needsSpace は i-1 番目と i 番目の字句の間に空白を書くかを返す。
func needsSpace(tokens []token.Token, i int) bool {
	prev, cur := tokens[i-1], tokens[i]
	switch {
	case prev.IsOp("(") || prev.IsOp("[") || prev.IsOp("{"):
		return false
	case cur.IsOp(")") || cur.IsOp("]") || cur.IsOp("}") || cur.IsOp(",") || cur.IsOp(":"):
		return false
	case prev.IsOp(",") || prev.IsOp(":"):
		return true
	case cur.IsOp("..") || cur.IsOp("..."):
		return !endsValue(prev)
	case prev.IsOp("..") || prev.IsOp("..."):
		return !startsValue(cur)
	case cur.IsOp("("):
		// f(x) と dv no "meif"(x) は呼び出し。weho (n) や mof (a) は空白を空ける
		return !(prev.Kind == token.Ident || prev.Kind == token.String || prev.IsOp(")") || prev.IsOp("]"))
	case cur.IsKeyword("izu") && prev.Kind == token.Op && token.CompoundOps[prev.Name]:
		return false
	case isUnary(tokens, i-1):
		return false
	}
	return true
}
