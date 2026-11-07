package interp

import (
	"errors"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// testFailure は sentky の中の shoubek が偽になったこと。巻き戻さず、その sentky を有罪にして止める (Jest の expect と同じ)。
// youlugive では捕まえない。
type testFailure struct {
	file string
	pos  token.Pos
}

func (t *testFailure) Error() string {
	return lexicon.Spell("shoubek") + " " + t.file + ":" + t.pos.String()
}

// runTest は sentky を1つ実行して判決を出す。落ちた場所は判決の下に字下げして出す。
// shoubek をテストとして扱うのは本体に直接書いたものだけで、呼んだ関数の中の shoubek はいつもどおり巻き戻す。
func (in *Interp) runTest(s *ast.Test, env *Env) error {
	// 本体の途中で frames が伸びて作り直されることがあるので、指す先を覚えずに毎回引く
	in.frames[len(in.frames)-1].test = true
	_, _, err := in.execBlock(s.Body, newEnv(env, false))
	in.frames[len(in.frames)-1].test = false
	// 巻き戻しと終了は、テストの結果ではなく実行全体の流れ
	if _, ok := errors.AsType[*backtrack](err); ok {
		return err
	}
	if _, ok := errors.AsType[*Exit](err); ok {
		return err
	}
	verdict := lexicon.Phrase("nii", "guilme")
	if err != nil {
		in.failed++
		verdict = lexicon.Spell("guilme")
	}
	in.emit(false, lexicon.Spell("hunggement")+" "+s.Name+": "+verdict+"\n")
	if err != nil {
		in.emit(false, "  "+uncaught(err).Error()+"\n")
	}
	return nil
}
