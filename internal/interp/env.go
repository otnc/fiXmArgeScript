package interp

// Env はブロック1つ分の変数。JavaScript の let と同じくブロックごとに分かれる。
type Env struct {
	vars   map[string]*variable
	parent *Env
	// funcTop は関数の一番外側のブロックか。ここより外の変数には izu できない (RFC の変数のスコープ)。
	funcTop bool
}

type variable struct {
	value   Value
	isConst bool // showence
}

func newEnv(parent *Env, funcTop bool) *Env {
	return &Env{vars: map[string]*variable{}, parent: parent, funcTop: funcTop}
}

func (e *Env) lookup(name string) (*variable, bool) {
	for env := e; env != nil; env = env.parent {
		if v, ok := env.vars[name]; ok {
			return v, true
		}
	}
	return nil, false
}

// assignResult は izu の左の名前をどう扱うか。
type assignResult int

const (
	assignNew        assignResult = iota // 今のブロックに新しく作る
	assignExisting                       // 関数の中で見つかった変数に代入する
	assignAcrossFunc                     // 関数の外の変数。Su·Ready
)

// resolveAssign は izu の左の名前を探す。関数の境目を越えた先で見つかったら assignAcrossFunc を返す。
func (e *Env) resolveAssign(name string) (*variable, assignResult) {
	crossed := false
	for env := e; env != nil; env = env.parent {
		if v, ok := env.vars[name]; ok {
			if crossed {
				return v, assignAcrossFunc
			}
			return v, assignExisting
		}
		if env.funcTop {
			crossed = true
		}
	}
	return nil, assignNew
}

func (e *Env) define(name string, v Value, isConst bool) {
	e.vars[name] = &variable{value: v, isConst: isConst}
}
