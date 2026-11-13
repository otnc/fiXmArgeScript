package interp

import (
	"errors"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// Session は fixmarge hak の対話の状態。入力を1つずつ実行し、前の入力で作った変数を引き継ぐ。
//
// 普段は前の入力が残した状態の続きで実行する。shoubek で巻き戻すときと、エラーになった入力をなかったことにするときだけ、
// 前の入力を頭から実行し直して状態を作り直す (やり直し方式)。前の入力の選択肢は確定していて記録もあるので、
// 実行し直しても同じ道を通り、外の世界には二度触らない。
type Session struct {
	in      *Interp
	entries [][]ast.Stmt // 最後まで実行できた入力
	env     *Env         // 一番外側のブロック。入力をまたいで使う
}

// NewSession は対話を始める。file はエラーの場所に出す名前。
func NewSession(file string, opts Options) *Session {
	in := newInterp(file, opts)
	return &Session{in: in, env: in.restart()}
}

// checkpoint は入力1つを実行する前の記録の長さ。その入力をなかったことにするときに、ここまで切り詰める。
type checkpoint struct {
	choices, inputs, records int
}

// Eval は入力1つを実行する。入力が式の文で終わっていたら、その値を出す (vanita は出さない)。
// 後の入力の shoubek が偽になっても前の入力には戻らず、Co-Jundic になる (RFC の fixmarge hak)。
// エラーになった入力はなかったことにする。ただし、もう出した出力と外の世界への操作は残る。
func (s *Session) Eval(stmts []ast.Stmt) error {
	in := s.in
	before := checkpoint{choices: in.frozen, inputs: len(in.inputs), records: len(in.records)}
	err := s.exec(stmts)
	for {
		bt, ok := errors.AsType[*backtrack](err)
		if !ok {
			break
		}
		if !in.advancePlan() {
			err = diag.New(diag.CoJundic, bt.file, bt.pos, "")
			break
		}
		if err = s.replay(); err == nil {
			err = s.exec(stmts)
		}
	}
	if err != nil {
		in.flush()
		if _, exited := errors.AsType[*Exit](err); !exited {
			s.rollback(before)
		}
		return uncaught(err)
	}
	in.commit()
	s.entries = append(s.entries, stmts)
	return nil
}

// exec は入力1つを、前の入力の続きとして実行する。
func (s *Session) exec(stmts []ast.Stmt) error {
	s.env.inheritAll()
	body, last := stmts, (*ast.ExprStmt)(nil)
	if n := len(stmts); n > 0 {
		if es, ok := stmts[n-1].(*ast.ExprStmt); ok {
			body, last = stmts[:n-1], es
		}
	}
	if _, _, err := s.in.execBlock(body, s.env); err != nil {
		return err
	}
	if last != nil {
		v, err := s.in.eval(last.X, s.env)
		if err != nil {
			return err
		}
		if v != null {
			// 文字列は配列の中と同じく " で囲み、redogh の出力と見分けられるようにする
			text, ok := nested(v)
			if !ok {
				return s.in.errorf(diag.Bavellabion, last.X.Position(), "")
			}
			s.in.emit(false, text+"\n")
		}
	}
	// 入力を実行し終えたら、ファイルの一番外側と同じく、動いている task と外の処理を待つ
	return s.in.drain()
}

// replay は最後まで実行できた入力を頭から実行し直して、一番外側のブロックを作り直す。出力は書き出し済みなので出ない。
func (s *Session) replay() error {
	s.env = s.in.restart()
	for _, stmts := range s.entries {
		if err := s.exec(stmts); err != nil {
			return err
		}
	}
	return nil
}

// rollback はエラーになった入力の記録を捨てて、前の入力が終わったときの状態に戻す。
func (s *Session) rollback(c checkpoint) {
	in := s.in
	in.plan = in.plan[:c.choices]
	in.frozen = c.choices
	in.inputs = in.inputs[:c.inputs]
	in.records = in.records[:c.records]
	// 前の入力は一度最後まで通り、同じ記録で実行し直すので、またエラーなく通る
	_ = s.replay()
	// 捨てた入力が出した出力はもう画面にあるので、次の入力の出力から数え直す
	in.committed = in.emitted
}
