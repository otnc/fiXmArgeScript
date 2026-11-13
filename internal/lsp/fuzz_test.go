package lsp

import (
	"testing"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// FuzzDocument は、どんなテキストのどの位置で補完や説明を求められても、言語サーバーが落ちないことを確かめる。
// エディタは書きかけのテキストを送ってくるので、壊れたソースが普通に届く。
func FuzzDocument(f *testing.F) {
	f.Add("\":meme\" aste meme werer\nweho f(x)\n  meme no zelln([x]) itswa\nfinend\nf(1) redogh\n", uint32(2), uint32(10))
	f.Add("x izu \"a #{1 +\n", uint32(0), uint32(9))
	f.Add("sentky \"t\"\n  tadact shoubek\nfinend\n\"./util\" no a aste b werer\nb.\n", uint32(4), uint32(2))

	s := &Server{docs: map[uri.URI]*document{}}
	f.Fuzz(func(t *testing.T, text string, line, char uint32) {
		d := newDocument(uri.File("/work/main.fxma"), text)
		p := protocol.Position{Line: line % 64, Character: char % 256}
		d.diagnostics()
		d.semanticTokens()
		d.hover(p)
		s.completion(d, p)
		s.definition(d, p)
	})
}
