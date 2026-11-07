package parser

import (
	"slices"
	"strconv"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// maxMacroDepth は呪いが呪いを呼べる深さの上限。超えると Bavellabion。
const maxMacroDepth = 64

// macro は curoi で定義した呪い1つ。本体は字句の並びのまま持ち、呼ばれた所に引数を当てはめて差し込む (Rust の macro_rules! と同じ型はめ式)。
type macro struct {
	params []string
	body   []token.Token
}

// expander は字句の並びの中の呪いを展開する。構文解析の前に1度だけ行う。
type expander struct {
	p      *parser
	macros map[string]*macro
	// renamed は呪いの中で作った名前に付ける番号。展開のたびに別の名前にして、呼び出し側の変数とぶつからないようにする
	renamed int
}

// expandMacros は curoi の定義を字句の並びから取り除き、呼び出しを展開する。
// 定義はファイルのどこに書いても、そのファイル全体で使える (weho の巻き上げと同じ)。
func (p *parser) expandMacros() error {
	e := &expander{p: p, macros: map[string]*macro{}}
	rest, err := e.collect(p.tokens)
	if err != nil {
		return err
	}
	if len(e.macros) == 0 {
		return nil
	}
	p.macros = e
	p.tokens, err = e.expand(rest, 0)
	return err
}

// collect は curoi の定義を取り出し、残りの字句を返す。定義の終わりは構文解析で探す。
func (e *expander) collect(tokens []token.Token) ([]token.Token, error) {
	var rest []token.Token
	for i := 0; i < len(tokens); {
		t := tokens[i]
		if !t.IsKeyword("curoi") {
			rest = append(rest, t)
			i++
			continue
		}
		sub := &parser{src: e.p.src, file: e.p.file, mode: e.p.mode, tokens: tokens, i: i + 1}
		name, err := sub.expectIdent()
		if err != nil {
			return nil, err
		}
		params, err := sub.macroParams()
		if err != nil {
			return nil, err
		}
		bodyStart := sub.i
		if _, err := sub.parseBlock(); err != nil {
			return nil, err
		}
		if _, dup := e.macros[name.Name]; dup {
			return nil, diag.New(diag.SuReady, e.p.file, name.Pos, e.p.text(name))
		}
		// parseBlock は finend まで読むので、本体はその手前まで
		e.macros[name.Name] = &macro{params: params, body: trimNewlines(tokens[bodyStart : sub.i-1])}
		i = sub.i
	}
	return rest, nil
}

func (p *parser) macroParams() ([]string, error) {
	if err := p.expectOp("("); err != nil {
		return nil, err
	}
	var params []string
	for !p.acceptOp(")") {
		if len(params) > 0 {
			if err := p.expectOp(","); err != nil {
				return nil, err
			}
		}
		t, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		params = append(params, t.Name)
	}
	return params, nil
}

func trimNewlines(tokens []token.Token) []token.Token {
	for len(tokens) > 0 && tokens[0].Kind == token.Newline {
		tokens = tokens[1:]
	}
	for len(tokens) > 0 && tokens[len(tokens)-1].Kind == token.Newline {
		tokens = tokens[:len(tokens)-1]
	}
	return tokens
}

// expand は呼び出しを本体に置き換える。置き換えた結果の中の呼び出しも、深さの上限まで展開する。
func (e *expander) expand(tokens []token.Token, depth int) ([]token.Token, error) {
	var out []token.Token
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		m, ok := e.macros[t.Name]
		// util no twice(...) のようにモジュールの中の名前として書いたものは呪いではない
		isCall := ok && t.Kind == token.Ident && i+1 < len(tokens) && tokens[i+1].IsOp("(") && (i == 0 || !tokens[i-1].IsKeyword("no"))
		if !isCall {
			out = append(out, t)
			continue
		}
		if depth >= maxMacroDepth {
			return nil, diag.New(diag.Bavellabion, e.p.file, t.Pos, e.p.text(t))
		}
		args, end, err := e.splitArgs(tokens, i+1)
		if err != nil {
			return nil, err
		}
		if len(args) != len(m.params) {
			return nil, diag.New(diag.AyAmis, e.p.file, t.Pos, e.p.text(t))
		}
		expanded, err := e.expand(e.substitute(m, args), depth+1)
		if err != nil {
			return nil, err
		}
		out = append(out, expanded...)
		i = end
	}
	return out, nil
}

// splitArgs は open の ( から対応する ) までを、一番外側の , で引数に分ける。end は ) の位置。
func (e *expander) splitArgs(tokens []token.Token, open int) (args [][]token.Token, end int, err error) {
	depth := 0
	var current []token.Token
	for i := open; i < len(tokens); i++ {
		t := tokens[i]
		switch {
		case t.IsOp("(") || t.IsOp("[") || t.IsOp("{"):
			depth++
			if depth == 1 {
				continue
			}
		case t.IsOp(")") || t.IsOp("]") || t.IsOp("}"):
			depth--
			if depth == 0 {
				if len(current) > 0 || len(args) > 0 {
					args = append(args, current)
				}
				return args, i, nil
			}
		case t.IsOp(",") && depth == 1:
			args = append(args, current)
			current = nil
			continue
		}
		current = append(current, t)
	}
	return nil, 0, diag.New(diag.AyAmis, e.p.file, tokens[open].Pos, e.p.text(tokens[open]))
}

// substitute は本体の引数の名前を、渡されたコードに置き換える。
func (e *expander) substitute(m *macro, args [][]token.Token) []token.Token {
	e.renamed++
	locals := createdNames(m.body)
	var out []token.Token
	for _, t := range m.body {
		if t.Kind != token.Ident {
			out = append(out, t)
			continue
		}
		if i := slices.Index(m.params, t.Name); i >= 0 {
			out = append(out, group(args[i])...)
			continue
		}
		if locals[t.Name] {
			// 下書きに書けない文字を混ぜて、呼び出し側のどの名前ともぶつからないようにする。エラーには元の綴りが出る
			t.Name += "·" + strconv.Itoa(e.renamed)
		}
		out = append(out, t)
	}
	return out
}

// createdNames は本体の中で作る変数の名前を集める。izu の左、aste の後ろ、youlugive の後ろに書いた名前。
func createdNames(body []token.Token) map[string]bool {
	names := map[string]bool{}
	for i, t := range body {
		if t.Kind != token.Ident {
			continue
		}
		prev := token.Token{}
		if i > 0 {
			prev = body[i-1]
		}
		next := token.Token{}
		if i+1 < len(body) {
			next = body[i+1]
		}
		startsLine := i == 0 || prev.Kind == token.Newline || prev.IsKeyword("showence")
		if (startsLine && (next.IsKeyword("izu") || next.Kind == token.Op && i+2 < len(body) && body[i+2].IsKeyword("izu"))) ||
			prev.IsKeyword("aste") || prev.IsKeyword("youlugive") {
			names[t.Name] = true
		}
	}
	return names
}

// expressionKeywords は式の途中に書けるキーワード。引数がこれらと記号だけでできていれば、式として括弧で包む。
var expressionKeywords = []string{"no", "senaji", "sond", "owk", "nii", "exi", "tadact", "ulieso", "vanita", "batte"}

// group は2語以上の式の引数を括弧で包む。square(1 + 2) を 1 + 2 * 1 + 2 にしないため。
// 文を渡したとき ("Mof" redogh) は包まずにそのまま差し込む。
func group(arg []token.Token) []token.Token {
	if len(arg) <= 1 {
		return arg
	}
	for _, t := range arg {
		if t.Kind == token.Keyword && !slices.Contains(expressionKeywords, t.Name) {
			return arg
		}
	}
	open := token.Token{Kind: token.Op, Name: "(", Pos: arg[0].Pos, Start: arg[0].Start, End: arg[0].Start}
	last := arg[len(arg)-1]
	closing := token.Token{Kind: token.Op, Name: ")", Pos: last.Pos, Start: last.End, End: last.End}
	return append(append([]token.Token{open}, arg...), closing)
}
