package grammar

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otnc/fiXmArgeScript/internal/lexer"
)

// TestUpToDate は、editors/vscode/syntaxes の文法が語の表と合っているかを確かめる。
// 語を足したら FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
func TestUpToDate(t *testing.T) {
	for _, mode := range []lexer.Mode{lexer.Draft, lexer.Fair} {
		want, err := TextMate(mode)
		require.NoError(t, err)
		path := filepath.Join("..", "..", "editors", "vscode", "syntaxes", File(mode))
		if os.Getenv("FXMA_UPDATE") != "" {
			require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
			require.NoError(t, os.WriteFile(path, want, 0o644))
			continue
		}
		got, err := os.ReadFile(path)
		require.NoError(t, err, "FXMA_UPDATE=1 go test ./internal/grammar/ で作ってください")
		require.Equal(t, string(want), string(got), "FXMA_UPDATE=1 go test ./internal/grammar/ で書き直してください")
	}
}
