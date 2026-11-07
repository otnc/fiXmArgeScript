// Package build は下書き (.fxma) を清書 (.fxmac) にする。
// 字句の並びと、字句の間の空白とコメントはそのまま写し、辞書の語だけを正式な綴りに置き換える。
package build

import (
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Build は下書きのソースから清書のソースを作る。
func Build(src, file string) (string, []diag.Warning, error) {
	tokens, warnings, err := lexer.Lex(src, file, lexer.Draft)
	if err != nil {
		return "", nil, err
	}
	var b strings.Builder
	if err := render(&b, src, tokens, 0, len(src), file); err != nil {
		return "", nil, err
	}
	return b.String(), warnings, nil
}

// compoundOps は izu の直前に付けて複合代入になる記号。清書では空白を挟まずに書く。
var compoundOps = map[string]bool{
	"+": true, "-": true, "*": true, "/": true, "//": true, "%": true, "**": true,
	"&": true, "|": true, "^": true, "<<": true, ">>": true,
}

func render(b *strings.Builder, src string, tokens []token.Token, from, to int, file string) error {
	prev := from
	for i, t := range tokens {
		if t.Kind == token.EOF {
			break
		}
		gap := src[prev:t.Start]
		if t.IsKeyword("izu") && i > 0 && tokens[i-1].Kind == token.Op && compoundOps[tokens[i-1].Name] && !strings.Contains(gap, "\n") {
			gap = ""
		}
		b.WriteString(gap)
		switch {
		case t.Lexicon:
			b.WriteString(lexicon.Spell(lexicon.Skeleton(t.Name)))
		case t.Kind == token.String:
			if err := renderString(b, src, t, file); err != nil {
				return err
			}
		default:
			b.WriteString(src[t.Start:t.End])
		}
		prev = t.End
	}
	b.WriteString(src[prev:to])
	return nil
}

// renderString は文字列をそのまま写し、#{} の中のコードだけを清書する。
func renderString(b *strings.Builder, src string, t token.Token, file string) error {
	cursor := t.Start
	for _, p := range t.Parts {
		if !p.IsCode {
			continue
		}
		end := p.Start + len(p.Text)
		b.WriteString(src[cursor:p.Start])
		tokens, _, err := lexer.LexRange(src, p.Start, end, p.Pos, file, lexer.Draft)
		if err != nil {
			return err
		}
		if err := render(b, src, tokens, p.Start, end, file); err != nil {
			return err
		}
		cursor = end
	}
	b.WriteString(src[cursor:t.End])
	return nil
}
