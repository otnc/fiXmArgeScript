package lsp

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

// editor はテストでエディタの役をする。サーバーから届いた診断をファイルごとに覚える。
type editor struct {
	protocol.UnimplementedClient
	diagnostics chan *protocol.PublishDiagnosticsParams
}

func (e *editor) PublishDiagnostics(_ context.Context, p *protocol.PublishDiagnosticsParams) error {
	e.diagnostics <- p
	return nil
}

func connect(t *testing.T) (protocol.Server, *editor) {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	serverEnd, clientEnd := net.Pipe()
	go func() { _ = Serve(ctx, serverEnd, "test") }()
	e := &editor{diagnostics: make(chan *protocol.PublishDiagnosticsParams, 16)}
	_, conn, server := protocol.NewClient(ctx, e, jsonrpc2.NewStream(clientEnd))
	t.Cleanup(func() {
		cancel()
		_ = conn.Close()
		_ = serverEnd.Close()
	})
	_, err := server.Initialize(ctx, &protocol.InitializeParams{})
	require.NoError(t, err)
	return server, e
}

const util = `":meme" aste meme werer

oqun weho average(xs)
  meme no zelln(xs) / nang(xs) itswa
finend
`

const main = `":coword" aste coword werer
"./util" no average werer
"./util" aste util werer

scores izu [72, 95]
average(scores) redogh
util no average(scores) redogh
unused izu 1
"#{scores}" redogh
`

// open はファイルをディスクに書き、エディタで開いたことにして、届いた診断を返す。
func open(t *testing.T, server protocol.Server, e *editor, path, text string) []protocol.Diagnostic {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(text), 0o644))
	err := server.DidOpen(t.Context(), &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{URI: uri.File(path), LanguageID: "fixmarge", Text: text},
	})
	require.NoError(t, err)
	p := <-e.diagnostics
	require.Equal(t, uri.File(path), p.URI)
	return p.Diagnostics
}

func at(path string, line, char uint32) protocol.TextDocumentPositionParams {
	return protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri.File(path)},
		Position:     protocol.Position{Line: line, Character: char},
	}
}

func setup(t *testing.T) (protocol.Server, *editor, string, string) {
	server, e := connect(t)
	dir := t.TempDir()
	mainPath, utilPath := filepath.Join(dir, "main.fxma"), filepath.Join(dir, "util.fxma")
	require.NoError(t, os.WriteFile(utilPath, []byte(util), 0o644))
	diags := open(t, server, e, mainPath, main)
	// 使っていない coword と unused。werer で読み込んだ名前は、文の頭ではなく名前に波線を引く
	require.Len(t, diags, 2)
	for _, d := range diags {
		require.Equal(t, protocol.String(lexicon.Spell("wasda")), d.Code)
	}
	require.Equal(t, protocol.Range{Start: protocol.Position{Line: 0, Character: 15}, End: protocol.Position{Line: 0, Character: 21}}, diags[0].Range)
	require.Equal(t, protocol.Range{Start: protocol.Position{Line: 7}, End: protocol.Position{Line: 7, Character: 6}}, diags[1].Range)
	return server, e, mainPath, utilPath
}

func TestDiagnosticsForSyntaxError(t *testing.T) {
	server, e := connect(t)
	diags := open(t, server, e, filepath.Join(t.TempDir(), "bad.fxma"), "x izu (1\n")
	require.Len(t, diags, 1)
	require.Equal(t, protocol.DiagnosticSeverityError, diags[0].Severity)
}

func TestDefinition(t *testing.T) {
	server, _, mainPath, utilPath := setup(t)
	averageInUtil := &protocol.Location{
		URI:   uri.File(utilPath),
		Range: protocol.Range{Start: protocol.Position{Line: 2, Character: 10}, End: protocol.Position{Line: 2, Character: 17}},
	}
	for _, pos := range []protocol.TextDocumentPositionParams{at(mainPath, 5, 2), at(mainPath, 6, 10)} {
		got, err := server.Definition(t.Context(), &protocol.DefinitionParams{TextDocumentPositionParams: pos})
		require.NoError(t, err)
		require.Equal(t, averageInUtil, got)
	}
	// #{} の中の名前も、ファイルの中の定義に飛ぶ
	got, err := server.Definition(t.Context(), &protocol.DefinitionParams{TextDocumentPositionParams: at(mainPath, 8, 4)})
	require.NoError(t, err)
	require.Equal(t, &protocol.Location{
		URI:   uri.File(mainPath),
		Range: protocol.Range{Start: protocol.Position{Line: 4}, End: protocol.Position{Line: 4, Character: 6}},
	}, got)
}

// 外部パッケージは、読み込み先に書いた版 (@1 と @2) ごとに、.fixmarge の別のフォルダの定義に飛ぶ。
func TestDefinitionInVersionedPackage(t *testing.T) {
	server, e := connect(t)
	dir := t.TempDir()
	write := func(name, text string) string {
		path := filepath.Join(dir, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(text), 0o644))
		return path
	}
	write("fixmarge.mod", "ado github.com/me/app\n\nneyo github.com/alice/strs v1.2.0\nneyo github.com/alice/strs v2.0.1\n")
	const strs = "oqun weho shout(s)\n  s itswa\nfinend\n"
	v1 := write(".fixmarge/github.com/alice/strs@v1.2.0/strs.fxma", strs)
	v2 := write(".fixmarge/github.com/alice/strs@v2.0.1/strs.fxma", "\n"+strs)
	mainPath := filepath.Join(dir, "main.fxma")
	diags := open(t, server, e, mainPath, "\"github.com/alice/strs@1\" aste old werer\n\"github.com/alice/strs@2\" no shout werer\nold no shout(\"a\") redogh\nshout(\"b\") redogh\n")
	require.Empty(t, diags)

	for _, c := range []struct {
		pos  protocol.TextDocumentPositionParams
		file string
		line uint32
	}{
		{at(mainPath, 2, 8), v1, 0},
		{at(mainPath, 3, 2), v2, 1},
	} {
		got, err := server.Definition(t.Context(), &protocol.DefinitionParams{TextDocumentPositionParams: c.pos})
		require.NoError(t, err)
		require.Equal(t, &protocol.Location{
			URI:   uri.File(c.file),
			Range: protocol.Range{Start: protocol.Position{Line: c.line, Character: 10}, End: protocol.Position{Line: c.line, Character: 15}},
		}, got)
	}

	// 補完にも、その版が公開している名前が出る
	list, err := server.Completion(t.Context(), &protocol.CompletionParams{TextDocumentPositionParams: at(mainPath, 2, 7)})
	require.NoError(t, err)
	require.Contains(t, labels(t, list), "shout")
}

func TestHover(t *testing.T) {
	server, _, mainPath, _ := setup(t)
	got, err := server.Hover(t.Context(), &protocol.HoverParams{TextDocumentPositionParams: at(mainPath, 5, 17)})
	require.NoError(t, err)
	require.Contains(t, got.Contents.(*protocol.MarkupContent).Value, lexicon.Spell("redogh"))
	require.Contains(t, got.Contents.(*protocol.MarkupContent).Value, lexicon.Meaning("redogh"))

	// average( の間にカーソルがあれば、( ではなく average の説明
	got, err = server.Hover(t.Context(), &protocol.HoverParams{TextDocumentPositionParams: at(mainPath, 5, 7)})
	require.NoError(t, err)
	require.Contains(t, got.Contents.(*protocol.MarkupContent).Value, `werer "./util" の average`)
}

func labels(t *testing.T, r protocol.CompletionResult) []string {
	t.Helper()
	var out []string
	for _, item := range r.(protocol.CompletionItemSlice) {
		out = append(out, item.Label)
	}
	return out
}

func TestCompletion(t *testing.T) {
	server, _, mainPath, _ := setup(t)
	got, err := server.Completion(t.Context(), &protocol.CompletionParams{TextDocumentPositionParams: at(mainPath, 4, 0)})
	require.NoError(t, err)
	require.Subset(t, labels(t, got), []string{"redogh", "nang", "scores", "average"})

	// util no ▮ には読み込み先が公開した名前だけを出す
	got, err = server.Completion(t.Context(), &protocol.CompletionParams{TextDocumentPositionParams: at(mainPath, 6, 8)})
	require.NoError(t, err)
	require.Equal(t, []string{"average"}, labels(t, got))
}

func TestCompletionStdMembers(t *testing.T) {
	server, e := connect(t)
	path := filepath.Join(t.TempDir(), "main.fxma")
	open(t, server, e, path, "\":coword\" aste coword werer\ncoword no c\n")
	got, err := server.Completion(t.Context(), &protocol.CompletionParams{TextDocumentPositionParams: at(path, 1, 11)})
	require.NoError(t, err)
	require.Contains(t, labels(t, got), "confsoy")
	require.NotContains(t, labels(t, got), "redogh")
}

func TestFormatting(t *testing.T) {
	server, e := connect(t)
	path := filepath.Join(t.TempDir(), "main.fxma")
	open(t, server, e, path, "x   izu 1\n")
	got, err := server.Formatting(t.Context(), &protocol.DocumentFormattingParams{TextDocument: protocol.TextDocumentIdentifier{URI: uri.File(path)}})
	require.NoError(t, err)
	require.Equal(t, []protocol.TextEdit{{
		Range:   protocol.Range{End: protocol.Position{Line: 1}},
		NewText: "x izu 1\n",
	}}, got)
}

func TestSemanticTokens(t *testing.T) {
	server, _, mainPath, _ := setup(t)
	got, err := server.SemanticTokensFull(t.Context(), &protocol.SemanticTokensParams{TextDocument: protocol.TextDocumentIdentifier{URI: uri.File(mainPath)}})
	require.NoError(t, err)
	require.Zero(t, len(got.Data)%5)
	type tok struct{ line, char, length, typ, mods uint32 }
	var toks []tok
	var line, char uint32
	for i := 0; i < len(got.Data); i += 5 {
		if got.Data[i] > 0 {
			char = 0
		}
		line += got.Data[i]
		char += got.Data[i+1]
		toks = append(toks, tok{line, char, got.Data[i+2], got.Data[i+3], got.Data[i+4]})
	}
	// scores izu [72, 95] の scores は変数の定義、izu はキーワード
	require.Contains(t, toks, tok{4, 0, 6, typeVariable, modDeclaration})
	require.Contains(t, toks, tok{4, 7, 3, typeKeyword, 0})
	// util no average の util はモジュール
	require.Contains(t, toks, tok{6, 0, 4, typeNamespace, 0})
	// "#{scores}" の中の scores は変数で、まわりは文字列
	require.Contains(t, toks, tok{8, 0, 3, typeString, 0})
	require.Contains(t, toks, tok{8, 3, 6, typeVariable, 0})
	require.True(t, slices.IsSortedFunc(toks, func(a, b tok) int {
		if a.line != b.line {
			return int(a.line) - int(b.line)
		}
		return int(a.char) - int(b.char)
	}))
}
