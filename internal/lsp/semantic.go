package lsp

import (
	"cmp"
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/scope"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// 色の種類と修飾。並びが LSP の legend の番号になる。
var (
	tokenTypes     = []string{"keyword", "variable", "parameter", "function", "namespace", "property", "string", "number", "operator"}
	tokenModifiers = []string{"declaration", "readonly", "defaultLibrary"}
)

const (
	typeKeyword = iota
	typeVariable
	typeParameter
	typeFunction
	typeNamespace
	typeProperty
	typeString
	typeNumber
	typeOperator
)

const (
	modDeclaration = 1 << iota
	modReadonly
	modDefaultLibrary
)

// span は色を付ける範囲1つ。start と end はバイト位置。
type span struct {
	start, end int
	typ, mods  int
}

// semanticTokens は LSP の semanticTokens/full の data を作る。
// 字句に分けられないとき (書きかけの文字列など) は空を返し、エディタの TextMate 文法の色に任せる。
func (d *document) semanticTokens() []uint32 {
	var spans []span
	for _, t := range d.tokens {
		spans = append(spans, d.spansOf(t)...)
	}
	slices.SortFunc(spans, func(a, b span) int { return cmp.Compare(a.start, b.start) })
	return d.encode(spans)
}

func (d *document) spansOf(t token.Token) []span {
	switch t.Kind {
	case token.Keyword:
		return []span{{t.Start, t.End, typeKeyword, 0}}
	case token.Int, token.Float:
		return []span{{t.Start, t.End, typeNumber, 0}}
	case token.Op:
		return []span{{t.Start, t.End, typeOperator, 0}}
	case token.String:
		return d.stringSpans(t)
	case token.Ident:
		typ, mods := d.classifyIdent(t)
		return []span{{t.Start, t.End, typ, mods}}
	}
	return nil
}

// stringSpans は文字列から #{} の中のコードを除いた部分を返す。中のコードは flatten で別の字句になっている。
func (d *document) stringSpans(t token.Token) []span {
	var out []span
	start := t.Start
	for _, p := range t.Parts {
		if !p.IsCode {
			continue
		}
		out = append(out, span{start, p.Start, typeString, 0})
		start = p.Start + len(p.Text)
	}
	return append(out, span{start, t.End, typeString, 0})
}

func (d *document) classifyIdent(t token.Token) (typ, mods int) {
	if b, ok := d.refs[t.Pos]; ok {
		if name, ok := d.nameToken(b); ok && name.Pos == t.Pos {
			mods |= modDeclaration
		}
		switch b.Kind {
		case scope.Func:
			return typeFunction, mods
		case scope.Param:
			return typeParameter, mods
		case scope.Const:
			return typeVariable, mods | modReadonly
		case scope.Module:
			return typeNamespace, mods
		}
		return typeVariable, mods
	}
	if m, ok := d.member(t.Pos); ok {
		if pkg, isStd := strings.CutPrefix(m.Module.Path, "fxma:"); isStd {
			if _, isFunc := interp.StdFunc(lexicon.Skeleton(pkg), m.Name); isFunc {
				return typeFunction, modDefaultLibrary
			}
		}
		return typeProperty, 0
	}
	if e, _, ok := lexicon.Lookup(lexicon.Skeleton(t.Name)); ok && t.Lexicon {
		switch {
		case e.Uses&lexicon.Package != 0:
			return typeNamespace, modDefaultLibrary
		case e.Uses&(lexicon.Builtin|lexicon.Function|lexicon.Compound) != 0:
			return typeFunction, modDefaultLibrary
		}
	}
	return typeVariable, 0
}

// encode は範囲を LSP の形 (前の字句からの行の差、列の差、長さ、種類、修飾) にする。
// 複数行にまたがる範囲 (""" の文字列) は、行ごとに分ける。
func (d *document) encode(spans []span) []uint32 {
	out := []uint32{}
	var prevLine, prevChar uint32
	for _, s := range spans {
		for start := s.start; start < s.end; {
			end := s.end
			if nl := strings.IndexByte(d.text[start:s.end], '\n'); nl >= 0 {
				end = start + nl
			}
			from, to := d.position(start), d.position(end)
			if to.Character > from.Character {
				char := from.Character
				if from.Line == prevLine {
					char -= prevChar
				}
				out = append(out, from.Line-prevLine, char, to.Character-from.Character, uint32(s.typ), uint32(s.mods))
				prevLine, prevChar = from.Line, from.Character
			}
			start = end + 1
		}
	}
	return out
}
