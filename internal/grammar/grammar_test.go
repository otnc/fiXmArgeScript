package grammar

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otnc/fiXmArgeScript/internal/lexer"
)

// TestUpToDate は、VS Code 拡張の TextMate 文法・言語の設定と Vim のファイルが語の表と合っているかを確かめる。
// 語を足したら FXMA_UPDATE=1 go test ./internal/grammar/ で書き直す。
func TestUpToDate(t *testing.T) {
	root := filepath.Join("..", "..")
	mod, err := ModTextMate()
	require.NoError(t, err)
	files := map[string][]byte{
		filepath.Join("autoload", "fixmarge", "words.vim"):      VimWords(),
		filepath.Join("editors", "vscode", "syntaxes", ModFile): mod,
		filepath.Join("syntax", VimModFiletype+".vim"):          VimModSyntax(),
	}
	for _, mode := range []lexer.Mode{lexer.Draft, lexer.Fair} {
		textMate, err := TextMate(mode)
		require.NoError(t, err)
		files[filepath.Join("editors", "vscode", "syntaxes", File(mode))] = textMate
		files[filepath.Join("syntax", VimFiletype(mode)+".vim")] = VimSyntax(mode)
		config, err := LanguageConfiguration(mode)
		require.NoError(t, err)
		files[filepath.Join("editors", "vscode", LanguageConfigurationFile(mode))] = config
	}
	for name, want := range files {
		path := filepath.Join(root, name)
		if os.Getenv("FXMA_UPDATE") != "" {
			require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
			require.NoError(t, os.WriteFile(path, want, 0o644))
			continue
		}
		got, err := os.ReadFile(path)
		require.NoError(t, err, "FXMA_UPDATE=1 go test ./internal/grammar/ で作ってください")
		require.Equal(t, string(want), string(got), "FXMA_UPDATE=1 go test ./internal/grammar/ で書き直してください: %s", name)
	}
}
