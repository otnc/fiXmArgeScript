// Package build は下書き (.fxma) を清書 (.fxmac) にする。
// 字句の並びと、字句の間の空白とコメントはそのまま写し、辞書の語だけを正式な綴りに置き換える。
// 清書は文章のように読めるよう、行頭の字下げだけは消す (ブロックは finend で閉じるので、字下げに意味はない)。
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
	if err := render(&b, src, tokens, 0, len(src), file, true); err != nil {
		return "", nil, err
	}
	return b.String(), warnings, nil
}

// render は from から to までを清書する。flat なら行頭の字下げを消す。
// #{} の中は文字列の一部なので flat にしない。
func render(b *strings.Builder, src string, tokens []token.Token, from, to int, file string, flat bool) error {
	prev := from
	// 改行は字句なので、改行の字句の直後のすき間は行頭から始まる
	lineStart := from == 0
	for i, t := range tokens {
		if t.Kind == token.EOF {
			break
		}
		gap := src[prev:t.Start]
		if t.IsKeyword("izu") && i > 0 && tokens[i-1].Kind == token.Op && token.CompoundOps[tokens[i-1].Name] && !strings.Contains(gap, "\n") {
			gap = ""
		}
		if flat {
			gap = unindent(gap, lineStart)
		}
		lineStart = t.Kind == token.Newline
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
	rest := src[prev:to]
	if flat {
		rest = unindent(rest, lineStart)
	}
	b.WriteString(rest)
	return nil
}

// unindent は、字句のすき間 (空白・改行・コメント) から、行頭の空白とタブを消す。
// atStart なら、すき間の頭も行頭として扱う。
func unindent(gap string, atStart bool) string {
	var b strings.Builder
	lineStart := atStart
	for _, c := range gap {
		if lineStart && (c == ' ' || c == '\t') {
			continue
		}
		lineStart = c == '\n'
		b.WriteRune(c)
	}
	return b.String()
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
		if err := render(b, src, tokens, p.Start, end, file, false); err != nil {
			return err
		}
		cursor = end
	}
	b.WriteString(src[cursor:t.End])
	return nil
}
