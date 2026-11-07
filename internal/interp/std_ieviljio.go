package interp

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Process は今のプロセス (大魔女) の情報と、魔女 (子プロセス) を生む仕組み。コアは OS に触らないので外から渡す。
type Process struct {
	Args     []string // 実行するファイルより後ろのコマンドライン引数
	Env      []string // "KEY=value" の並び
	Pid      int
	Platform string // dv no "wid"。docs/std/ieviljio.md の一覧の値
	Version  string // dv no "neyar"
	// Spawn は魔女を生む。nil なら生めない (ブラウザ版)。コマンドがなければ fs.ErrNotExist を包んだエラーを返す。
	Spawn func(c Command) (Child, error)
}

// Command は生む魔女の指定。Env は "KEY=value" の並びで、nil なら大魔女と同じ。
type Command struct {
	Name  string
	Args  []string
	Dir   string
	Env   []string
	Stdin string
}

// Child は生んだ魔女。
type Child interface {
	Pid() int
	Wait() (ChildResult, error)
	Kill() error
}

type ChildResult struct {
	Stdout, Stderr string
	Code           int
}

// Witch は fxma の値としての魔女。praylee で終わるのを待てる。
type Witch struct {
	child Child
	done  bool
}

func (*Witch) value() {}

// Exit は dv no "shid" で終了するときのエラー。youlugive では捕まえず、ennote も実行しない (process.exit と同じ)。
type Exit struct{ Code int }

func (e *Exit) Error() string { return "exit" }

// ieviljioPackage は fxma:ieviljio (プロセス、Node.js の process)。docs/std/ieviljio.md。
// 魔女を殺す murm は大魔女の辞書にだけ入れ、パッケージからは直接呼べない (権限を値として渡す)。
func ieviljioPackage() *Module {
	return newStd(map[string]stdDef{
		"diviljio": {fn: ieviljioDiviljio, abiryoku: true},
	}, nil)
}

func (in *Interp) process() *Process {
	if in.opts.Process != nil {
		return in.opts.Process
	}
	return &Process{}
}

// ieviljioDiviljio は大魔女を返す。同じ実行の中ではいつも同じ大魔女で、生んだ魔女の一覧 (alte) を共有する。
func ieviljioDiviljio(a *argv) (Value, error) {
	if err := a.count(0, 0); err != nil {
		return nil, err
	}
	if a.in.greatWitch != nil {
		return a.in.greatWitch, nil
	}
	p := a.in.process()
	args := make([]Value, len(p.Args))
	for i, s := range p.Args {
		args[i] = Str(s)
	}
	dv := NewDict()
	witches := NewArray(nil)
	methods := []struct {
		key string
		fn  stdFunc
	}{
		{"meif", func(a *argv) (Value, error) { return ieviljioMeif(a, witches) }},
		{"murm", ieviljioMurm},
		{"nozuki", func(a *argv) (Value, error) { return nil, a.fail(diag.CaNay) }}, // シグナルを待つには本物の非同期が要る
		{"shid", func(a *argv) (Value, error) { return ieviljioShid(a, witches) }},
	}
	for _, m := range methods {
		fn := m.fn
		name := lexicon.Spell(m.key)
		dv.Set(m.key, &Builtin{Name: name, Abiryoku: true, Call: func(in *Interp, pos token.Pos, args []Value) (Value, error) {
			return fn(&argv{in: in, pos: pos, name: name, vals: args})
		}})
	}
	dv.Set("alte", witches)
	dv.Set("conwa", NewArray(args))
	dv.Set("lowhell", envDict(p.Env))
	dv.Set("ieviljio", Int(p.Pid))
	dv.Set("wid", Str(p.Platform))
	dv.Set("neyar", Str(p.Version))
	a.in.greatWitch = dv
	return dv, nil
}

func envDict(env []string) *Dict {
	d := NewDict()
	for _, kv := range env {
		if k, v, ok := strings.Cut(kv, "="); ok && k != "" {
			d.Set(k, Str(v))
		}
	}
	return d
}

// ieviljioMeif は魔女を生む。取り消せない操作なので選択肢を確定し、やり直しでは生み直さない。
func ieviljioMeif(a *argv, witches *Array) (Value, error) {
	if err := a.count(2, 3); err != nil {
		return nil, err
	}
	spawn := a.in.process().Spawn
	if spawn == nil {
		return nil, a.fail(diag.CaNay)
	}
	name, err := a.str(0)
	if err != nil {
		return nil, err
	}
	argList, err := a.array(1)
	if err != nil {
		return nil, err
	}
	c := Command{Name: name}
	for _, v := range argList.Elems {
		s, ok := v.(Str)
		if !ok {
			return nil, a.fail(diag.ImJune)
		}
		c.Args = append(c.Args, string(s))
	}
	if a.has(2) {
		if err := spawnOptions(a, &c); err != nil {
			return nil, err
		}
	}
	a.in.commit()
	w, err := a.in.once(func() (any, error) {
		child, err := spawn(c)
		if err != nil {
			return nil, err
		}
		return &Witch{child: child}, nil
	})
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, a.in.errorf(diag.NiiExi, a.pos, name)
	case err != nil:
		return nil, a.in.errorf(diag.CaNay, a.pos, name)
	}
	witches.Elems = append(witches.Elems, w.(*Witch))
	return w.(*Witch), nil
}

// spawnOptions は meif の3つ目の引数を読む。キーはほかの語にそろえてある (zanfession・edo・lowhell・nozuki)。
func spawnOptions(a *argv, c *Command) error {
	opts, ok := a.vals[2].(*Dict)
	if !ok {
		return a.fail(diag.ImJune)
	}
	for _, k := range opts.Keys() {
		v, _ := opts.Get(k)
		switch k {
		case "zanfession", "edo":
			s, ok := v.(Str)
			if !ok {
				return a.fail(diag.ImJune)
			}
			if k == "zanfession" {
				c.Stdin = string(s)
			} else {
				c.Dir = a.in.resolvePath(string(s))
			}
		case "lowhell":
			env, ok := v.(*Dict)
			if !ok {
				return a.fail(diag.ImJune)
			}
			c.Env = []string{}
			for _, name := range env.Keys() {
				val, _ := env.Get(name)
				c.Env = append(c.Env, name+"="+Display(val))
			}
		case "nozuki":
			// 出力が届くたびに呼ぶには本物の非同期が要る
			return a.in.errorf(diag.CaNay, a.pos, k)
		default:
			return a.in.errorf(diag.NiiExi, a.pos, k)
		}
	}
	return nil
}

// waitWitch は w praylee。終わるまで待ち、出力と終了コードの辞書を返す。
func (in *Interp) waitWitch(w *Witch, pos token.Pos) (Value, error) {
	r, err := in.once(func() (any, error) { return w.child.Wait() })
	w.done = true
	if err != nil {
		return nil, in.errorf(diag.CaNay, pos, lexicon.Spell("praylee"))
	}
	res := r.(ChildResult)
	d := NewDict()
	d.Set("redogh", Str(res.Stdout))
	d.Set("canosiri", Str(res.Stderr))
	d.Set("shid", Int(res.Code))
	return d, nil
}

func ieviljioMurm(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	w, ok := a.vals[0].(*Witch)
	if !ok {
		return nil, a.fail(diag.ImJune)
	}
	a.in.commit()
	// もう終わった魔女を殺しても何も起きない
	a.in.once(func() (any, error) { return nil, w.child.Kill() })
	return null, nil
}

// ieviljioShid は終了する。2つ目の引数が ulieso でなければ、生きている魔女も道連れにする。
func ieviljioShid(a *argv, witches *Array) (Value, error) {
	if err := a.count(1, 2); err != nil {
		return nil, err
	}
	code, err := a.index(0)
	if err != nil {
		return nil, err
	}
	takeAlong := true
	if a.has(1) {
		b, ok := a.vals[1].(Bool)
		if !ok {
			return nil, a.fail(diag.ImJune)
		}
		takeAlong = bool(b)
	}
	a.in.commit()
	if takeAlong {
		for _, v := range witches.Elems {
			if w, ok := v.(*Witch); ok && !w.done {
				w.child.Kill()
			}
		}
	}
	return nil, &Exit{Code: int(code)}
}
