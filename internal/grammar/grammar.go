// Package grammar は、エディタの色分けに使うファイルを語の表 (lexicon) から作る。
// VS Code 拡張 (editors/vscode) の TextMate 文法と、Vim・Neovim の syntax と字下げの語 (リポジトリの直下)。
// 字句解析に沿った正確な色は LSP のセマンティックトークンが上書きするので、ここでは fixmarge を入れていなくても読める程度の色を付ける。
package grammar

import (
	"cmp"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"regexp"
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

// キーワードの色の系統。どれにも入らないキーワードは keyword.control。
var (
	constants = []string{"tadact", "ulieso", "vanita"}
	operators = []string{"izu", "no", "exi", "nii", "sond", "owk", "senaji", "batte"}
	declares  = []string{"weho", "curoi"}
	modifiers = []string{"abiryoku", "itsver", "oqun", "showence"}
)

type grammar struct {
	Schema     string          `json:"$schema"`
	Name       string          `json:"name"`
	ScopeName  string          `json:"scopeName"`
	Patterns   []rule          `json:"patterns"`
	Repository map[string]rule `json:"repository"`
}

type rule struct {
	Name          string             `json:"name,omitempty"`
	Match         string             `json:"match,omitempty"`
	Begin         string             `json:"begin,omitempty"`
	End           string             `json:"end,omitempty"`
	Captures      map[string]capture `json:"captures,omitempty"`
	BeginCaptures map[string]capture `json:"beginCaptures,omitempty"`
	EndCaptures   map[string]capture `json:"endCaptures,omitempty"`
	Include       string             `json:"include,omitempty"`
	Patterns      []rule             `json:"patterns,omitempty"`
}

type capture struct {
	Name string `json:"name"`
}

func include(name string) rule { return rule{Include: "#" + name} }

// File は文法のファイルの名前。下書きと清書で別の文法にする (語の書き方が違うため)。
func File(mode lexer.Mode) string {
	if mode == lexer.Fair {
		return "fixmargec.tmLanguage.json"
	}
	return "fixmarge.tmLanguage.json"
}

// TextMate は下書きか清書の文法を JSON で返す。
func TextMate(mode lexer.Mode) ([]byte, error) {
	g := grammar{
		Schema:    "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
		Name:      "fiXmArge Script",
		ScopeName: "source.fixmarge",
		Patterns: []rule{
			include("comment"), include("string"), include("number"), include("declaration"),
			include("constant"), include("operator-word"), include("keyword"), include("support"),
			include("call"), include("operator"),
		},
		Repository: map[string]rule{
			"comment": {Name: "comment.line.number-sign.fixmarge", Match: `#.*$`},
			"string": {Patterns: []rule{
				stringRule(`"""`, "string.quoted.triple.fixmarge"),
				stringRule(`"`, "string.quoted.double.fixmarge"),
			}},
			"escape": {Name: "constant.character.escape.fixmarge", Match: `\\(?:u\{[0-9A-Fa-f]+\}|.)`},
			"interpolation": {
				Begin:         `#\{`,
				End:           `\}`,
				BeginCaptures: map[string]capture{"0": {"punctuation.section.interpolation.begin.fixmarge"}},
				EndCaptures:   map[string]capture{"0": {"punctuation.section.interpolation.end.fixmarge"}},
				Name:          "meta.interpolation.fixmarge",
				Patterns:      []rule{{Include: "$self"}},
			},
			"number": {Name: "constant.numeric.fixmarge", Match: `\b\d(?:_?\d)*(?:\.\d(?:_?\d)*)?(?:[eE][+-]?\d+)?\b`},
			"declaration": {
				Match: words(declares, mode) + `\s+` + identifier,
				Captures: map[string]capture{
					"1": {"storage.type.function.fixmarge"},
					"2": {"entity.name.function.fixmarge"},
				},
			},
			"constant":      {Name: "constant.language.fixmarge", Match: words(constants, mode)},
			"operator-word": {Name: "keyword.operator.word.fixmarge", Match: words(operators, mode)},
			"keyword": {Patterns: []rule{
				{Name: "storage.modifier.fixmarge", Match: words(modifiers, mode)},
				{Name: "keyword.control.fixmarge", Match: words(otherKeywords(), mode)},
			}},
			"support": {Patterns: []rule{
				{Name: "support.function.fixmarge", Match: words(uses(lexicon.Builtin|lexicon.Function|lexicon.Compound), mode)},
				{Name: "entity.name.namespace.fixmarge", Match: words(uses(lexicon.Package), mode)},
			}},
			"call": {
				Match:    identifier + `(?=\s*\()`,
				Captures: map[string]capture{"1": {"entity.name.function.fixmarge"}},
			},
			"operator": {Name: "keyword.operator.fixmarge", Match: `\*\*|//|<<|>>|<=|>=|\.\.=?|[-+*/%<>&|^~]`},
		},
	}
	if mode == lexer.Fair {
		g.Name += " (清書)"
		g.ScopeName = "source.fixmargec"
	}
	return marshal(g)
}

// ModFile は fixmarge.mod の文法のファイルの名前。
const ModFile = "fixmarge-mod.tmLanguage.json"

// ModTextMate は fixmarge.mod の文法を JSON で返す。行の頭の項目名 (lexicon の ModField) と文字列だけに色を付ける。
func ModTextMate() ([]byte, error) {
	fields := modFields()
	return marshal(grammar{
		Schema:    "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
		Name:      "fiXmArge Script (fixmarge.mod)",
		ScopeName: "source.fixmarge-mod",
		Patterns:  []rule{include("field"), include("string")},
		Repository: map[string]rule{
			"field": {
				Match:    `^\s*` + words(fields, lexer.Draft),
				Captures: map[string]capture{"1": {"keyword.other.fixmarge-mod"}},
			},
			"string": {
				Name:  "string.quoted.double.fixmarge-mod",
				Match: `"(?:[^"\\]|\\.)*"`,
			},
		},
	})
}

// modFields は fixmarge.mod の項目名 (キーワードを兼ねる youlugive も含む)。
func modFields() []string {
	var out []string
	for _, e := range lexicon.Entries() {
		if e.Uses&lexicon.ModField != 0 {
			out = append(out, e.Name)
		}
	}
	return out
}

func marshal(g grammar) ([]byte, error) {
	out, err := json.Marshal(g, json.Deterministic(true), jsontext.WithIndent("  "))
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}

func stringRule(quote, scope string) rule {
	return rule{
		Name:          scope,
		Begin:         quote,
		End:           quote,
		BeginCaptures: map[string]capture{"0": {"punctuation.definition.string.begin.fixmarge"}},
		EndCaptures:   map[string]capture{"0": {"punctuation.definition.string.end.fixmarge"}},
		Patterns:      []rule{include("escape"), include("interpolation")},
	}
}

// identifier は自分で付けた名前1つを捕まえる。表の語 (記号の入った清書の綴りなど) には support が先に色を付ける。
const identifier = `([A-Za-z_]\w*)`

// words は語の並びを、語として一致する正規表現にする。長い綴りを先に試す。
// 下書きのキーワードは小文字だけで書く。ほかの名前は大文字小文字を区別しない (字句解析と同じ)。
// 清書は正式な綴りそのもので、' や - で終わる綴りもあるので \b ではなく前後の文字で区切る。
func words(skels []string, mode lexer.Mode) string {
	var alts []string
	keywordsOnly := true
	for _, s := range skels {
		e, _, ok := lexicon.Lookup(lexicon.Skeleton(s))
		if !ok {
			panic("grammar: 表にない語: " + s)
		}
		keywordsOnly = keywordsOnly && e.IsKeyword()
		if mode == lexer.Fair {
			alts = append(alts, regexp.QuoteMeta(e.Spelling))
		} else {
			alts = append(alts, e.Name)
		}
	}
	slices.SortFunc(alts, func(a, b string) int { return cmp.Or(cmp.Compare(len(b), len(a)), strings.Compare(a, b)) })
	alts = slices.Compact(alts)
	group := strings.Join(alts, "|")
	switch {
	case mode == lexer.Fair:
		return `(?<![\w'’~·-])(` + group + `)(?![\w])`
	case keywordsOnly:
		return `\b(` + group + `)\b`
	}
	return `\b((?i:` + group + `))\b`
}

func otherKeywords() []string {
	special := slices.Concat(constants, operators, declares, modifiers)
	var out []string
	for _, e := range lexicon.Entries() {
		if e.IsKeyword() && !slices.Contains(special, e.Name) {
			out = append(out, e.Name)
		}
	}
	return out
}

// uses は使い道 u を持つ名前のうち、キーワードでないものを返す。
func uses(u lexicon.Use) []string {
	var out []string
	for _, e := range lexicon.Entries() {
		if e.Uses&u != 0 && !e.IsKeyword() {
			out = append(out, e.Name)
		}
	}
	return out
}
