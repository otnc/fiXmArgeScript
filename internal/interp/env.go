package interp

// Env はブロック1つ分の変数。JavaScript の let と同じくブロックごとに分かれる。
type Env struct {
	vars   map[string]*variable
	parent *Env
	// funcTop は関数の一番外側のブロックか。ここより外の変数には izu できない (RFC の変数のスコープ)。
	funcTop bool
	// inherited は fxma hak の前の入力で作った名前。weho、showence、werer で作り直せる (Chrome のコンソールと同じ)。
	inherited map[string]bool
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
	delete(e.inherited, name)
}

// declared は、このブロックでもう作った名前かを返す。前の入力から引き継いだ名前は作り直せるので数えない。
func (e *Env) declared(name string) bool {
	_, ok := e.vars[name]
	return ok && !e.inherited[name]
}

// inheritAll は今ある名前を全部、前の入力から引き継いだものにする。fxma hak で入力を1つ実行する前に呼ぶ。
func (e *Env) inheritAll() {
	e.inherited = map[string]bool{}
	for name := range e.vars {
		e.inherited[name] = true
	}
}
