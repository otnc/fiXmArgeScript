package interp

import (
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// stdPackages は "fxma:名前" で読み込める標準パッケージ。キーは lexicon の骨格。
// 標準パッケージの関数が werer の処理を呼び返すので、変数の初期化の循環を避けて init で入れる。
var stdPackages map[string]func() *Module

func init() {
	stdPackages = map[string]func() *Module{
		"meme":       memePackage,
		"coword":     cowordPackage,
		"varizama":   varizamaPackage,
		"json":       jsonPackage,
		"izanavited": izanavitedPackage,
	}
}

// stdFunc は標準パッケージの関数の中身。
type stdFunc func(a *argv) (Value, error)

// stdDef は標準パッケージの関数1つ。
type stdDef struct {
	fn       stdFunc
	abiryoku bool
}

// newStd は標準パッケージを作る。中の名前はすべて公開する。
func newStd(funcs map[string]stdDef, consts map[string]Value) *Module {
	m := &Module{members: map[string]Value{}, exported: map[string]bool{}}
	for name, def := range funcs {
		spelled := spellName(name)
		fn := def.fn
		m.members[name] = &Builtin{
			Name:     spelled,
			Abiryoku: def.abiryoku,
			Call: func(in *Interp, pos token.Pos, args []Value) (Value, error) {
				return fn(&argv{in: in, pos: pos, name: spelled, vals: args})
			},
		}
		m.exported[name] = true
	}
	for name, v := range consts {
		m.members[name] = v
		m.exported[name] = true
	}
	return m
}

// spellName は表にある名前なら正式な綴りを、sin や pi のような固有名詞ならそのまま返す。
func spellName(name string) string {
	if e, _, ok := lexicon.Lookup(lexicon.Skeleton(name)); ok {
		return e.Spelling
	}
	return name
}

// argv は標準パッケージの関数が受け取った引数。型が合わなければ、関数の名前を対象にした imJune を返す。
type argv struct {
	in   *Interp
	pos  token.Pos
	name string
	vals []Value
}

func (a *argv) fail(kind diag.Kind) error { return a.in.errorf(kind, a.pos, a.name) }

func (a *argv) count(lo, hi int) error {
	if len(a.vals) < lo || len(a.vals) > hi {
		return a.fail(diag.ImJune)
	}
	return nil
}

func (a *argv) has(i int) bool { return i < len(a.vals) }

func (a *argv) num(i int) (float64, error) {
	f, ok := toFloat(a.vals[i])
	if !ok {
		return 0, a.fail(diag.ImJune)
	}
	return f, nil
}

func (a *argv) index(i int) (int64, error) {
	n, ok := toIndex(a.vals[i])
	if !ok {
		return 0, a.fail(diag.ImJune)
	}
	return n, nil
}

func (a *argv) str(i int) (string, error) {
	s, ok := a.vals[i].(Str)
	if !ok {
		return "", a.fail(diag.ImJune)
	}
	return string(s), nil
}

func (a *argv) array(i int) (*Array, error) {
	xs, ok := a.vals[i].(*Array)
	if !ok {
		return nil, a.fail(diag.ImJune)
	}
	return xs, nil
}

func (a *argv) fn(i int) (Value, error) {
	if !isFunc(a.vals[i]) {
		return nil, a.fail(diag.ImJune)
	}
	return a.vals[i], nil
}

// call は引数で受け取った関数を呼ぶ。abiryoku の関数を渡されていれば、呼び出しの frame で入出力が許されている。
func (a *argv) call(f Value, args ...Value) (Value, error) {
	return a.in.call(f, args, a.pos, a.name)
}

// callBool は条件として使う関数を呼ぶ。tadact と ulieso 以外を返したら imJune (RFC の真偽値)。
func (a *argv) callBool(f Value, args ...Value) (bool, error) {
	v, err := a.call(f, args...)
	if err != nil {
		return false, err
	}
	b, ok := v.(Bool)
	if !ok {
		return false, a.fail(diag.ImJune)
	}
	return bool(b), nil
}

func isFunc(v Value) bool {
	switch v.(type) {
	case *Func, *Builtin:
		return true
	}
	return false
}
