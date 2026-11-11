// Package lsp はエディタのための言語サーバー (fixmarge trench) を持つ。
// 診断、整形、補完、ホバー、定義へ移動、色分け (セマンティックトークン) を、字句解析から lint までの部品を組み合わせて出す。
package lsp

import (
	"context"
	"errors"
	"io"
	"sync"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/otnc/fiXmArgeScript/internal/format"
)

// Server は言語サーバー。開いているファイルを持ち、変わるたびに解析し直して診断を送る。
type Server struct {
	protocol.UnimplementedServer

	version string
	// client と conn は接続を作った後に入れる。エディタからの知らせはその前に届くことがあるので、ready が閉じるまで待つ
	ready  chan struct{}
	client protocol.Client
	conn   jsonrpc2.Conn

	mu   sync.Mutex
	docs map[uri.URI]*document
}

// Serve は rwc (標準入出力など) で LSP を話し、エディタが接続を閉じるか exit を送るまで続ける。
func Serve(ctx context.Context, rwc io.ReadWriteCloser, version string) error {
	s := &Server{version: version, ready: make(chan struct{}), docs: map[uri.URI]*document{}}
	_, s.conn, s.client = protocol.NewServer(ctx, s, jsonrpc2.NewStream(rwc))
	close(s.ready)
	<-s.conn.Done()
	if err := s.conn.Err(); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func (s *Server) Initialize(context.Context, *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync:           protocol.TextDocumentSyncKindFull,
			CompletionProvider:         &protocol.CompletionOptions{},
			HoverProvider:              protocol.Boolean(true),
			DefinitionProvider:         protocol.Boolean(true),
			DocumentFormattingProvider: protocol.Boolean(true),
			SemanticTokensProvider: &protocol.SemanticTokensOptions{
				Legend: protocol.SemanticTokensLegend{TokenTypes: tokenTypes, TokenModifiers: tokenModifiers},
				Full:   protocol.Boolean(true),
			},
		},
		ServerInfo: protocol.ServerInfo{Name: "fixmarge trench", Version: protocol.NewOptional(s.version)},
	}, nil
}

func (s *Server) Initialized(context.Context, *protocol.InitializedParams) error { return nil }

func (s *Server) Shutdown(context.Context) error { return nil }

func (s *Server) Exit(context.Context) error {
	<-s.ready
	return s.conn.Close()
}

func (s *Server) DidOpen(ctx context.Context, p *protocol.DidOpenTextDocumentParams) error {
	return s.update(ctx, p.TextDocument.URI, p.TextDocument.Text)
}

func (s *Server) DidChange(ctx context.Context, p *protocol.DidChangeTextDocumentParams) error {
	// 同期は全文なので、最後の変更が新しい中身そのもの
	for _, c := range p.ContentChanges {
		if whole, ok := c.(*protocol.TextDocumentContentChangeWholeDocument); ok {
			if err := s.update(ctx, p.TextDocument.URI, whole.Text); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) DidClose(ctx context.Context, p *protocol.DidCloseTextDocumentParams) error {
	s.mu.Lock()
	delete(s.docs, p.TextDocument.URI)
	s.mu.Unlock()
	<-s.ready
	return s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{URI: p.TextDocument.URI, Diagnostics: []protocol.Diagnostic{}})
}

func (s *Server) update(ctx context.Context, u uri.URI, text string) error {
	d := newDocument(u, text)
	s.mu.Lock()
	s.docs[u] = d
	s.mu.Unlock()
	<-s.ready
	return s.client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{URI: u, Diagnostics: d.diagnostics()})
}

func (s *Server) document(u uri.URI) *document {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.docs[u]
}

func (s *Server) Formatting(_ context.Context, p *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	d := s.document(p.TextDocument.URI)
	if d == nil {
		return nil, nil
	}
	// 構文が正しくなければ整えられない。エラーは診断で出ているので、ここでは何もしない
	formatted, err := format.Format(d.text, d.path)
	if err != nil || formatted == d.text {
		return []protocol.TextEdit{}, nil
	}
	return []protocol.TextEdit{{Range: d.fullRange(), NewText: formatted}}, nil
}

func (s *Server) Completion(_ context.Context, p *protocol.CompletionParams) (protocol.CompletionResult, error) {
	d := s.document(p.TextDocument.URI)
	if d == nil {
		return protocol.CompletionItemSlice{}, nil
	}
	return protocol.CompletionItemSlice(s.completion(d, p.Position)), nil
}

func (s *Server) Hover(_ context.Context, p *protocol.HoverParams) (*protocol.Hover, error) {
	d := s.document(p.TextDocument.URI)
	if d == nil {
		return nil, nil
	}
	text, r, ok := d.hover(p.Position)
	if !ok {
		return nil, nil
	}
	return &protocol.Hover{Contents: &protocol.MarkupContent{Kind: protocol.MarkupKindMarkdown, Value: text}, Range: &r}, nil
}

func (s *Server) Definition(_ context.Context, p *protocol.DefinitionParams) (protocol.DefinitionResult, error) {
	d := s.document(p.TextDocument.URI)
	if d == nil {
		return nil, nil
	}
	loc, ok := s.definition(d, p.Position)
	if !ok {
		return nil, nil
	}
	return &loc, nil
}

func (s *Server) SemanticTokensFull(_ context.Context, p *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	d := s.document(p.TextDocument.URI)
	if d == nil {
		return &protocol.SemanticTokens{Data: []uint32{}}, nil
	}
	return &protocol.SemanticTokens{Data: d.semanticTokens()}, nil
}
