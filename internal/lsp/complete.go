package lsp

import (
	"strings"

	"go.lsp.dev/protocol"

	"github.com/otnc/fiXmArgeScript/internal/explain"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/scope"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// completionUses は補完に出す語の使い道。エラー名やコマンド名などは、コードの名前としては書かないので出さない。
const completionUses = lexicon.Keyword | lexicon.Builtin | lexicon.Package | lexicon.Function | lexicon.Compound

// completion は補完の候補を返す。モジュールの後ろ (meme no ▮) なら中の名前、それ以外は語の表とファイルの中の名前。
func (s *Server) completion(d *document, p protocol.Position) []protocol.CompletionItem {
	if m, ok := d.moduleBefore(d.offsetOf(p)); ok {
		return s.memberItems(d, m)
	}
	var items []protocol.CompletionItem
	for _, e := range lexicon.Entries() {
		if e.Uses&completionUses == 0 {
			continue
		}
		items = append(items, protocol.CompletionItem{
			Label:  d.label(e),
			Kind:   entryKind(e),
			Detail: protocol.NewOptional(e.Spelling + " " + lexicon.Meaning(e.Name)),
		})
	}
	seen := map[string]bool{}
	if d.scope != nil {
		for _, b := range d.scope.Bindings {
			if seen[b.Name] || strings.Contains(b.Name, "·") {
				continue
			}
			seen[b.Name] = true
			items = append(items, protocol.CompletionItem{
				Label:  b.Name,
				Kind:   bindingKind(b),
				Detail: protocol.NewOptional(kindLabel(b)),
			})
		}
	}
	return items
}

// moduleBefore は、off の直前が「モジュール no」ならそのモジュールを返す。打ちかけの名前があれば、その前を見る。
func (d *document) moduleBefore(off int) (*scope.Binding, bool) {
	start := off
	for start > 0 && isWordByte(d.text[start-1]) {
		start--
	}
	var before []token.Token
	for _, t := range d.tokens {
		if t.End <= start && t.Kind != token.Newline && t.Kind != token.EOF {
			before = append(before, t)
		}
	}
	if len(before) < 2 || !before[len(before)-1].IsKeyword("no") {
		return nil, false
	}
	b, ok := d.refs[before[len(before)-2].Pos]
	return b, ok && b.Kind == scope.Module
}

func isWordByte(c byte) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func (s *Server) memberItems(d *document, m *scope.Binding) []protocol.CompletionItem {
	var items []protocol.CompletionItem
	if pkg, isStd := strings.CutPrefix(m.Path, "fxma:"); isStd {
		pkg = lexicon.Skeleton(pkg)
		for _, name := range interp.StdMembers(pkg) {
			kind := protocol.CompletionItemKindConstant
			if _, isFunc := interp.StdFunc(pkg, name); isFunc {
				kind = protocol.CompletionItemKindFunction
			}
			item := protocol.CompletionItem{Label: name, Kind: kind}
			if e, _, ok := lexicon.Lookup(lexicon.Skeleton(name)); ok {
				item.Label = d.label(e)
				item.Detail = protocol.NewOptional(e.Spelling + " " + lexicon.Meaning(e.Name))
			}
			items = append(items, item)
		}
		return items
	}
	for _, ex := range s.exports(d, m.Path) {
		kind := protocol.CompletionItemKindVariable
		if ex.binding.Kind == scope.Func {
			kind = protocol.CompletionItemKindFunction
		}
		items = append(items, protocol.CompletionItem{Label: ex.binding.Name, Kind: kind, Detail: protocol.NewOptional(m.Path)})
	}
	return items
}

// label は補完で挿入する形。下書きは骨格 (2語の名前は camelCase)、清書は正式な綴り。
func (d *document) label(e lexicon.Entry) string {
	if lexer.ModeOf(d.path) == lexer.Fair {
		return e.Spelling
	}
	return e.Name
}

func entryKind(e lexicon.Entry) protocol.CompletionItemKind {
	switch {
	case e.IsKeyword():
		return protocol.CompletionItemKindKeyword
	case e.Uses&lexicon.Package != 0:
		return protocol.CompletionItemKindModule
	}
	return protocol.CompletionItemKindFunction
}

func bindingKind(b *scope.Binding) protocol.CompletionItemKind {
	switch b.Kind {
	case scope.Func:
		return protocol.CompletionItemKindFunction
	case scope.Const:
		return protocol.CompletionItemKindConstant
	case scope.Module:
		return protocol.CompletionItemKindModule
	}
	return protocol.CompletionItemKindVariable
}

// kindLabel は定義の種類の説明。ホバーと補完の詳細に出す。
func kindLabel(b *scope.Binding) string {
	switch b.Kind {
	case scope.Const:
		return "showence (証拠品)"
	case scope.Func:
		return "weho (関数)"
	case scope.Param:
		return "引数"
	case scope.LoopVar:
		return "aste の変数"
	case scope.CatchVar:
		return "youlugive の変数"
	case scope.Import:
		return "werer \"" + b.Path + "\" の " + b.Remote
	case scope.Module:
		return "werer \"" + b.Path + "\""
	}
	return "変数"
}

// hover は位置の語の説明を返す。ファイルの中の名前なら定義の種類、辞書の語なら訳、エラー名なら fxma nai と同じ説明。
func (d *document) hover(p protocol.Position) (string, protocol.Range, bool) {
	t, ok := d.tokenAt(p)
	if !ok || (t.Kind != token.Ident && t.Kind != token.Keyword) {
		return "", protocol.Range{}, false
	}
	if b, ok := d.refs[t.Pos]; ok && t.Kind == token.Ident {
		return "```fxma\n" + b.Name + "\n```\n" + kindLabel(b), d.tokenRange(t), true
	}
	name := t.Name
	if m, ok := d.member(t.Pos); ok {
		name = m.Name
	}
	e, _, ok := lexicon.Lookup(lexicon.Skeleton(name))
	if !ok {
		return "", protocol.Range{}, false
	}
	lines := []string{"**" + e.Spelling + "** (" + e.Name + ")", "", lexicon.Meaning(e.Name)}
	if e.Uses&lexicon.Error != 0 {
		if _, text, ok := explain.Lookup(e.Name); ok {
			lines = append(lines, "", text)
		}
	}
	return strings.Join(lines, "\n"), d.tokenRange(t), true
}
