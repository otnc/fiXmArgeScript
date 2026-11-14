package grammar

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"slices"

	"github.com/otnc/fiXmArgeScript/internal/lexer"
)

// LanguageConfigurationFile は VS Code の言語の設定ファイルの名前。字下げの語の綴りが違うので、下書きと清書で別にする。
func LanguageConfigurationFile(mode lexer.Mode) string {
	if mode == lexer.Fair {
		return "fixmargec.language-configuration.json"
	}
	return "fixmarge.language-configuration.json"
}

type languageConfiguration struct {
	Comments         comments          `json:"comments"`
	Brackets         [][2]string       `json:"brackets"`
	AutoClosingPairs []autoClosingPair `json:"autoClosingPairs"`
	SurroundingPairs [][2]string       `json:"surroundingPairs"`
	IndentationRules *indentationRules `json:"indentationRules,omitempty"`
}

type comments struct {
	LineComment string `json:"lineComment"`
}

type autoClosingPair struct {
	Open  string   `json:"open"`
	Close string   `json:"close"`
	NotIn []string `json:"notIn,omitempty"`
}

type indentationRules struct {
	IncreaseIndentPattern string `json:"increaseIndentPattern"`
	DecreaseIndentPattern string `json:"decreaseIndentPattern"`
}

// LanguageConfiguration は VS Code の言語の設定 (括弧と字下げ) を JSON で返す。
// 字下げは Vim の autoload/fixmarge/indent.vim と同じ語を使う。ただし VS Code は1行を正規表現1つで見るので、語を数えずに
// 「開く語があって同じ行に finend がない」行の次を深くする。清書は字下げしない (fixmarge sinruits が行頭の空白を消す)。
func LanguageConfiguration(mode lexer.Mode) ([]byte, error) {
	pairs := [][2]string{{"(", ")"}, {"[", "]"}, {"{", "}"}}
	c := languageConfiguration{
		Comments: comments{LineComment: "#"},
		Brackets: pairs,
		AutoClosingPairs: []autoClosingPair{
			{Open: "(", Close: ")"},
			{Open: "[", Close: "]"},
			{Open: "{", Close: "}"},
			{Open: `"`, Close: `"`, NotIn: []string{"string"}},
		},
		SurroundingPairs: slices.Concat(pairs, [][2]string{{`"`, `"`}}),
	}
	if mode == lexer.Draft {
		closer := words([]string{"finend"}, mode)
		dedent := words(dedentWords, mode)
		opener := words([]string{"weho", "teitance", "sentky", "persemon", "curoi", "azagara", "mof"}, mode)
		c.IndentationRules = &indentationRules{
			// コメントの中の語で深くしないように、# より前だけを見る。
			IncreaseIndentPattern: `^(?![^#]*` + closer + `)(?:\s*` + dedent + `|[^#]*` + opener + `)`,
			DecreaseIndentPattern: `^\s*` + dedent,
		}
	}
	out, err := json.Marshal(c, jsontext.WithIndent("  "))
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}
