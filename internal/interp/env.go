package interp

// Env はブロック1つ分の変数。JavaScript の let と同じくブロックごとに分かれる。
//
// ブロックは入るたびに作るので、作るのを軽くする。ほとんどのブロックは変数が数個なので作った順に並べて持ち、
// 多くなったとき (ファイルの一番外側など) だけ名前の索引を作る。
type Env struct {
	vars   []variable
	index  map[string]int // 名前 → vars の位置。vars が indexedFrom 個以上のときだけ作る
	parent *Env
	// funcTop は関数の一番外側のブロックか。ここより外の変数には izu できない (spec/syntax.md の変数のスコープ)。
	funcTop bool
	// inherited は fixmarge hak の前の入力で作った名前。weho、showence、werer で作り直せる (Chrome のコンソールと同じ)。
	inherited map[string]bool
}

type variable struct {
	name    string
	value   Value
	isConst bool // showence
}

// indexedFrom は、変数がこの数になったら索引を作る境目。これより少なければ、順に比べるほうが速い。
const indexedFrom = 8

func newEnv(parent *Env, funcTop bool) *Env {
	return &Env{parent: parent, funcTop: funcTop}
}

// own はこのブロックで作った変数を返す。なければ nil。
// 返したポインタは、次に define するまでしか使えない (vars が伸びると置き場所が変わる)。
func (e *Env) own(name string) *variable {
	if e.index != nil {
		if i, ok := e.index[name]; ok {
			return &e.vars[i]
		}
		return nil
	}
	for i := range e.vars {
		if e.vars[i].name == name {
			return &e.vars[i]
		}
	}
	return nil
}

func (e *Env) lookup(name string) (*variable, bool) {
	for env := e; env != nil; env = env.parent {
		if v := env.own(name); v != nil {
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
		if v := env.own(name); v != nil {
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
	delete(e.inherited, name)
	if old := e.own(name); old != nil {
		old.value, old.isConst = v, isConst
		return
	}
	e.vars = append(e.vars, variable{name: name, value: v, isConst: isConst})
	switch {
	case e.index != nil:
		e.index[name] = len(e.vars) - 1
	case len(e.vars) == indexedFrom:
		e.index = make(map[string]int, 2*indexedFrom)
		for i := range e.vars {
			e.index[e.vars[i].name] = i
		}
	}
}

// declared は、このブロックでもう作った名前かを返す。前の入力から引き継いだ名前は作り直せるので数えない。
func (e *Env) declared(name string) bool {
	return e.own(name) != nil && !e.inherited[name]
}

// inheritAll は今ある名前を全部、前の入力から引き継いだものにする。fixmarge hak で入力を1つ実行する前に呼ぶ。
func (e *Env) inheritAll() {
	e.inherited = map[string]bool{}
	for i := range e.vars {
		e.inherited[e.vars[i].name] = true
	}
}
