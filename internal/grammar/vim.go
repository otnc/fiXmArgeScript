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
		return "fixmargec"
	}
	return "fixmarge"
}

// VimSyntax は Vim の syntax ファイルを返す。グループの並びは TextMate 文法と同じで、
// Vim は同じ位置で始まるものが重なると後に定義したほうを使うので、優先したいものほど後ろに書く。
func VimSyntax(mode lexer.Mode) []byte {
	var b strings.Builder
	b.WriteString(vimHeader)
	b.WriteString("if exists('b:current_syntax')\n  finish\nendif\n\n")
	b.WriteString("syntax iskeyword @,48-57,_\n\n")
	match := func(group, pattern string) { fmt.Fprintf(&b, "syntax match %s /%s/\n", group, pattern) }

	match("fixmargeOperator", `\*\*\|\/\/\|<<\|>>\|<=\|>=\|\.\.=\=\|[-+*\/%<>&|^~]`)
	match("fixmargeNumber", `\<\d\%(_\=\d\)*\%(\.\d\%(_\=\d\)*\)\=\%([eE][+-]\=\d\+\)\=\>`)
	match("fixmargeCall", `\<\h\w*\ze\s*(`)
	match("fixmargeNamespace", vimWords(uses(lexicon.Package), mode))
	match("fixmargeBuiltin", vimWords(uses(lexicon.Builtin|lexicon.Function|lexicon.Compound), mode))
	match("fixmargeKeyword", vimWords(otherKeywords(), mode))
	match("fixmargeModifier", vimWords(modifiers, mode))
	match("fixmargeOperatorWord", vimWords(operators, mode))
	match("fixmargeConstant", vimWords(constants, mode))
	match("fixmargeDeclare", vimWords(declares, mode))
	match("fixmargeFunctionName", `\%(`+vimWords(declares, mode)+`\s\+\)\@<=\h\w*`)
	b.WriteString("syntax match fixmargeComment /#.*$/ contains=@Spell\n")
	b.WriteString(`syntax match fixmargeEscape /\\\%(u{\x\+}\|.\)/ contained
syntax region fixmargeInterpolation matchgroup=fixmargeInterpolationDelimiter start=/#{/ end=/}/ contained contains=TOP
syntax region fixmargeString matchgroup=fixmargeStringDelimiter start=/"/ skip=/\\./ end=/"/ contains=fixmargeEscape,fixmargeInterpolation,@Spell
syntax region fixmargeString matchgroup=fixmargeStringDelimiter start=/"""/ end=/"""/ contains=fixmargeEscape,fixmargeInterpolation,@Spell
`)

	b.WriteString("\n")
	for _, link := range [][2]string{
		{"fixmargeOperator", "Operator"},
		{"fixmargeNumber", "Number"},
		{"fixmargeCall", "Function"},
		{"fixmargeNamespace", "Type"},
		{"fixmargeBuiltin", "Function"},
		{"fixmargeKeyword", "Keyword"},
		{"fixmargeModifier", "StorageClass"},
		{"fixmargeOperatorWord", "Operator"},
		{"fixmargeConstant", "Constant"},
		{"fixmargeDeclare", "Keyword"},
		{"fixmargeFunctionName", "Function"},
		{"fixmargeComment", "Comment"},
		{"fixmargeEscape", "SpecialChar"},
		{"fixmargeInterpolationDelimiter", "Delimiter"},
		{"fixmargeStringDelimiter", "String"},
		{"fixmargeString", "String"},
	} {
		fmt.Fprintf(&b, "highlight default link %s %s\n", link[0], link[1])
	}
	fmt.Fprintf(&b, "\nlet b:current_syntax = '%s'\n", VimFiletype(mode))
	return []byte(b.String())
}

// VimModFiletype は fixmarge.mod の Vim のファイルタイプの名前。
const VimModFiletype = "fixmargemod"

// VimModSyntax は fixmarge.mod の syntax ファイルを返す。TextMate の ModTextMate と同じく、行の頭の項目名と文字列だけに色を付ける。
func VimModSyntax() []byte {
	var b strings.Builder
	b.WriteString(vimHeader)
	b.WriteString("if exists('b:current_syntax')\n  finish\nendif\n\n")
	fmt.Fprintf(&b, "syntax match fixmargemodField /^\\s*\\zs%s/\n", vimWords(modFields(), lexer.Draft))
	b.WriteString("syntax region fixmargemodString start=/\"/ skip=/\\\\./ end=/\"/ oneline\n\n")
	b.WriteString("highlight default link fixmargemodField Keyword\n")
	b.WriteString("highlight default link fixmargemodString String\n")
	fmt.Fprintf(&b, "\nlet b:current_syntax = '%s'\n", VimModFiletype)
	return []byte(b.String())
}

// VimWords は字下げ (indent/fixmarge.vim) が使う、ブロックの語の正規表現を下書きと清書の両方について返す。
func VimWords() []byte {
	var b strings.Builder
	b.WriteString(vimHeader)
	b.WriteString("\n\" ブロックを開く語・閉じる語・区切る語の正規表現。ファイルタイプ (fixmarge か fixmargec) ごとに返す。\n")
	b.WriteString("function! fixmarge#words#Block(filetype) abort\n")
	for _, mode := range []lexer.Mode{lexer.Fair, lexer.Draft} {
		indent := "  "
		if mode == lexer.Fair {
			b.WriteString("  if a:filetype ==# 'fixmargec'\n")
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
	b.WriteString("function! fixmarge#words#Dedent(filetype) abort\n")
	fmt.Fprintf(&b, "  if a:filetype ==# 'fixmargec'\n    return %s\n  endif\n  return %s\nendfunction\n",
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
