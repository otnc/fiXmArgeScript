package grammar

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

const vimHeader = `" 語の表 (internal/lexicon) から internal/grammar が作る。手で直さず FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
scriptencoding utf-8
`

// VimFiletype は Vim のファイルタイプの名前。下書きと清書で別にする (語の書き方が違うため)。
func VimFiletype(mode lexer.Mode) string {
	if mode == lexer.Fair {
		return "fxmac"
	}
	return "fxma"
}

// VimSyntax は Vim の syntax ファイルを返す。グループの並びは TextMate 文法と同じで、
// Vim は同じ位置で始まるものが重なると後に定義したほうを使うので、優先したいものほど後ろに書く。
func VimSyntax(mode lexer.Mode) []byte {
	var b strings.Builder
	b.WriteString(vimHeader)
	b.WriteString("if exists('b:current_syntax')\n  finish\nendif\n\n")
	b.WriteString("syntax iskeyword @,48-57,_\n\n")
	match := func(group, pattern string) { fmt.Fprintf(&b, "syntax match %s /%s/\n", group, pattern) }

	match("fxmaOperator", `\*\*\|\/\/\|<<\|>>\|<=\|>=\|\.\.=\=\|[-+*\/%<>&|^~]`)
	match("fxmaNumber", `\<\d\%(_\=\d\)*\%(\.\d\%(_\=\d\)*\)\=\%([eE][+-]\=\d\+\)\=\>`)
	match("fxmaCall", `\<\h\w*\ze\s*(`)
	match("fxmaNamespace", vimWords(uses(lexicon.Package), mode))
	match("fxmaBuiltin", vimWords(uses(lexicon.Builtin|lexicon.Function|lexicon.Compound), mode))
	match("fxmaKeyword", vimWords(otherKeywords(), mode))
	match("fxmaModifier", vimWords(modifiers, mode))
	match("fxmaOperatorWord", vimWords(operators, mode))
	match("fxmaConstant", vimWords(constants, mode))
	match("fxmaDeclare", vimWords(declares, mode))
	match("fxmaFunctionName", `\%(`+vimWords(declares, mode)+`\s\+\)\@<=\h\w*`)
	b.WriteString("syntax match fxmaComment /#.*$/ contains=@Spell\n")
	b.WriteString(`syntax match fxmaEscape /\\\%(u{\x\+}\|.\)/ contained
syntax region fxmaInterpolation matchgroup=fxmaInterpolationDelimiter start=/#{/ end=/}/ contained contains=TOP
syntax region fxmaString matchgroup=fxmaStringDelimiter start=/"/ skip=/\\./ end=/"/ contains=fxmaEscape,fxmaInterpolation,@Spell
syntax region fxmaString matchgroup=fxmaStringDelimiter start=/"""/ end=/"""/ contains=fxmaEscape,fxmaInterpolation,@Spell
`)

	b.WriteString("\n")
	for _, link := range [][2]string{
		{"fxmaOperator", "Operator"},
		{"fxmaNumber", "Number"},
		{"fxmaCall", "Function"},
		{"fxmaNamespace", "Type"},
		{"fxmaBuiltin", "Function"},
		{"fxmaKeyword", "Keyword"},
		{"fxmaModifier", "StorageClass"},
		{"fxmaOperatorWord", "Operator"},
		{"fxmaConstant", "Constant"},
		{"fxmaDeclare", "Keyword"},
		{"fxmaFunctionName", "Function"},
		{"fxmaComment", "Comment"},
		{"fxmaEscape", "SpecialChar"},
		{"fxmaInterpolationDelimiter", "Delimiter"},
		{"fxmaStringDelimiter", "String"},
		{"fxmaString", "String"},
	} {
		fmt.Fprintf(&b, "highlight default link %s %s\n", link[0], link[1])
	}
	fmt.Fprintf(&b, "\nlet b:current_syntax = '%s'\n", VimFiletype(mode))
	return []byte(b.String())
}

// VimAdoFiletype は fxma.ado の Vim のファイルタイプの名前。
const VimAdoFiletype = "fxmaado"

// VimAdoSyntax は fxma.ado の syntax ファイルを返す。TextMate の AdoTextMate と同じく、行の頭の項目名と文字列だけに色を付ける。
func VimAdoSyntax() []byte {
	var b strings.Builder
	b.WriteString(vimHeader)
	b.WriteString("if exists('b:current_syntax')\n  finish\nendif\n\n")
	fmt.Fprintf(&b, "syntax match fxmaadoField /^\\s*\\zs%s/\n", vimWords(modFields(), lexer.Draft))
	b.WriteString("syntax region fxmaadoString start=/\"/ skip=/\\\\./ end=/\"/ oneline\n\n")
	b.WriteString("highlight default link fxmaadoField Keyword\n")
	b.WriteString("highlight default link fxmaadoString String\n")
	fmt.Fprintf(&b, "\nlet b:current_syntax = '%s'\n", VimAdoFiletype)
	return []byte(b.String())
}

// VimWords は字下げ (indent/fxma.vim) が使う、ブロックの語の正規表現を下書きと清書の両方について返す。
func VimWords() []byte {
	var b strings.Builder
	b.WriteString(vimHeader)
	b.WriteString("\n\" ブロックを開く語・閉じる語・区切る語の正規表現。ファイルタイプ (fxma か fxmac) ごとに返す。\n")
	b.WriteString("function! fxma#words#Block(filetype) abort\n")
	for _, mode := range []lexer.Mode{lexer.Fair, lexer.Draft} {
		indent := "  "
		if mode == lexer.Fair {
			b.WriteString("  if a:filetype ==# 'fxmac'\n")
			indent = "    "
		}
		pairs := [][2]string{
			{"opener", vimWords([]string{"weho", "teitance", "sentky", "curoi", "azagara", "mof"}, mode)},
			{"continued", vimWords([]string{"siva"}, mode) + `\s\+` + vimWords([]string{"mof"}, mode)},
			{"closer", vimWords([]string{"finend"}, mode)},
			{"dedent", vimWords(dedentWords, mode)},
		}
		b.WriteString(indent + "return {\n")
		for _, p := range pairs {
			fmt.Fprintf(&b, "%s\\ '%s': '%s',\n", indent, p[0], strings.ReplaceAll(p[1], "'", "''"))
		}
		b.WriteString(indent + "\\ }\n")
		if mode == lexer.Fair {
			b.WriteString("  endif\n")
		}
	}
	b.WriteString("endfunction\n")

	b.WriteString("\n\" 行の頭に来たら1段浅く書く語。打ち終えたら字下げし直す (indentkeys)。\n")
	b.WriteString("function! fxma#words#Dedent(filetype) abort\n")
	fmt.Fprintf(&b, "  if a:filetype ==# 'fxmac'\n    return %s\n  endif\n  return %s\nendfunction\n",
		vimList(dedentWords, lexer.Fair), vimList(dedentWords, lexer.Draft))
	return []byte(b.String())
}

var dedentWords = []string{"finend", "siva", "youlugive", "ennote"}

// vimList は語の並びを、下書きか清書の綴りの Vim のリストにする。
func vimList(skels []string, mode lexer.Mode) string {
	var items []string
	for _, s := range skels {
		e, _, _ := lexicon.Lookup(lexicon.Skeleton(s))
		name := e.Name
		if mode == lexer.Fair {
			name = e.Spelling
		}
		items = append(items, "'"+strings.ReplaceAll(name, "'", "''")+"'")
	}
	return "[" + strings.Join(items, ", ") + "]"
}

// vimWords は words と同じく、語の並びを語として一致する Vim の正規表現にする。
func vimWords(skels []string, mode lexer.Mode) string {
	var alts []string
	keywordsOnly := true
	for _, s := range skels {
		e, _, ok := lexicon.Lookup(lexicon.Skeleton(s))
		if !ok {
			panic("grammar: 表にない語: " + s)
		}
		keywordsOnly = keywordsOnly && e.IsKeyword()
		if mode == lexer.Fair {
			alts = append(alts, vimQuote(e.Spelling))
		} else {
			alts = append(alts, e.Name)
		}
	}
	slices.SortFunc(alts, func(a, b string) int { return cmp.Or(cmp.Compare(len(b), len(a)), strings.Compare(a, b)) })
	alts = slices.Compact(alts)
	group := `\%(` + strings.Join(alts, `\|`) + `\)`
	switch {
	case mode == lexer.Fair:
		return `\C\%([[:alnum:]_'’~·-]\)\@<!` + group + `\w\@!`
	case keywordsOnly:
		return `\C\<` + group + `\>`
	}
	return `\c\<` + group + `\>`
}

// vimQuote は綴りを Vim の正規表現 (magic) の中でそのまま一致する形にする。/ は syntax match の区切り。
func vimQuote(s string) string {
	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune(`\.*[]~^$/`, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}
