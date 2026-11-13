package interp

import (
	"errors"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

// flow は文を実行した後にどこへ進むか。
type flow int

const (
	flowNormal flow = iota
	flowBreak
	flowContinue
	flowReturn
)

// execBlock はブロックの文を順に実行する。weho は JavaScript と同じく、先にまとめて定義しておく (巻き上げ)。
func (in *Interp) execBlock(stmts []ast.Stmt, env *Env) (flow, Value, error) {
	for _, s := range stmts {
		d, ok := s.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if env.declared(d.Fn.Name) {
			return flowNormal, nil, in.errorf(diag.SuReady, d.Pos, d.Fn.Name)
		}
		env.define(d.Fn.Name, &Func{Lit: d.Fn, Env: env, File: in.frame().file}, false)
	}
	for _, s := range stmts {
		f, v, err := in.exec(s, env)
		if err != nil || f != flowNormal {
			return f, v, err
		}
	}
	return flowNormal, nil, nil
}

func (in *Interp) exec(s ast.Stmt, env *Env) (flow, Value, error) {
	switch s := s.(type) {
	case *ast.ExprStmt:
		_, err := in.eval(s.X, env)
		return flowNormal, nil, err
	case *ast.Assign:
		return flowNormal, nil, in.assign(s, env)
	case *ast.Print:
		if !in.frame().abiryoku {
			return flowNormal, nil, in.errorf(diag.KiNbOO, s.Pos, "")
		}
		v, err := in.eval(s.X, env)
		if err != nil {
			return flowNormal, nil, err
		}
		in.emit(s.Stderr, Display(v)+"\n")
		return flowNormal, nil, nil
	case *ast.Return:
		if s.X == nil {
			return flowReturn, null, nil
		}
		v, err := in.eval(s.X, env)
		return flowReturn, v, err
	case *ast.Throw:
		return flowNormal, nil, in.throw(s, env)
	case *ast.Assert:
		ok, err := in.evalBool(s.X, env, s.Pos, "shoubek")
		switch {
		case err != nil || ok:
		case in.frame().test:
			err = &testFailure{file: in.frame().file, pos: s.Pos}
		default:
			err = &backtrack{file: in.frame().file, pos: s.Pos}
		}
		return flowNormal, nil, err
	case *ast.If:
		for i, cond := range s.Conds {
			ok, err := in.evalBool(cond, env, cond.Position(), "mof")
			if err != nil {
				return flowNormal, nil, err
			}
			if ok {
				return in.execBlock(s.Blocks[i], newEnv(env, false))
			}
		}
		return in.execBlock(s.Else, newEnv(env, false))
	case *ast.While:
		return in.execWhile(s, env)
	case *ast.For:
		return in.execFor(s, env)
	case *ast.FuncDecl:
		return flowNormal, nil, nil // execBlock で定義済み
	case *ast.Try:
		return in.execTry(s, env)
	case *ast.Continue:
		return flowContinue, nil, nil
	case *ast.Break:
		return flowBreak, nil, nil
	case *ast.Test:
		switch {
		case in.mode == modeRun:
			return flowNormal, nil, in.errorf(diag.CaNay, s.Pos, s.Keyword())
		case s.Bench != (in.mode == modeBench):
			return flowNormal, nil, nil
		case s.Bench:
			return flowNormal, nil, in.runBench(s, env)
		}
		return flowNormal, nil, in.runTest(s, env)
	case *ast.Import:
		return flowNormal, nil, in.execImport(s, env)
	}
	return flowNormal, nil, in.errorf(diag.SaiVI, s.Position(), "")
}

// loopBody はループの本体の結果を、ループを続けるかどうかに読み替える。
func loopBody(f flow, v Value, err error) (stop bool, result flow, value Value, e error) {
	switch {
	case err != nil:
		return true, flowNormal, nil, err
	case f == flowBreak:
		return true, flowNormal, nil, nil
	case f == flowReturn:
		return true, flowReturn, v, nil
	}
	return false, flowNormal, nil, nil
}

func (in *Interp) execWhile(s *ast.While, env *Env) (flow, Value, error) {
	for {
		ok, err := in.evalBool(s.Cond, env, s.Cond.Position(), "azagara")
		if err != nil || !ok {
			return flowNormal, nil, err
		}
		if stop, f, v, err := loopBody(in.execBlock(s.Body, newEnv(env, false))); stop {
			return f, v, err
		}
	}
}

func (in *Interp) execFor(s *ast.For, env *Env) (flow, Value, error) {
	iter, err := in.eval(s.Iter, env)
	if err != nil {
		return flowNormal, nil, err
	}
	var result flow
	var value Value
	err = in.each(iter, s.Iter, func(item Value) (bool, error) {
		// ループの変数は1周ごとに新しく作る (JavaScript の for (const x of xs))
		body := newEnv(env, false)
		body.define(s.Var, item, false)
		stop, f, v, err := loopBody(in.execBlock(s.Body, body))
		result, value = f, v
		return stop, err
	})
	return result, value, err
}

// each は aste で回せる値の要素を順に f に渡す。f が true を返したら止める。
func (in *Interp) each(v Value, at ast.Node, f func(Value) (bool, error)) error {
	switch v := v.(type) {
	case *Array:
		// 回している途中で naddy されても JavaScript と同じく追いかける
		for i := 0; i < len(v.Elems); i++ {
			if stop, err := f(v.Elems[i]); stop || err != nil {
				return err
			}
		}
	case Str:
		for _, r := range string(v) {
			if stop, err := f(Str(string(r))); stop || err != nil {
				return err
			}
		}
	case *Dict:
		for _, k := range v.Keys() {
			if stop, err := f(Str(k)); stop || err != nil {
				return err
			}
		}
	case *Range:
		for i := v.Lo; v.contains(i); i += v.Step {
			if stop, err := f(Int(i)); stop || err != nil {
				return err
			}
		}
	default:
		return in.errorf(diag.ImJune, at.Position(), lexicon.Spell("aste"))
	}
	return nil
}

func (in *Interp) execTry(s *ast.Try, env *Env) (flow, Value, error) {
	f, v, err := in.execBlock(s.Body, newEnv(env, false))
	if exc, ok := catchable(err); ok && s.HasCatch {
		catchEnv := newEnv(env, false)
		if s.CatchVar != "" {
			catchEnv.define(s.CatchVar, exc, false)
		}
		f, v, err = in.execBlock(s.Catch, catchEnv)
	}
	if !s.HasFinally || abandons(err) {
		return f, v, err
	}
	ff, fv, ferr := in.execBlock(s.Finally, newEnv(env, false))
	if ferr != nil || ff != flowNormal {
		return ff, fv, ferr
	}
	return f, v, err
}

// abandons は、今の実行を捨てるエラーかを返す。ennote は実行しない。
// 巻き戻すときと task を片付けるときは変数も出力も捨てるので、実行しても意味がない。
// dv no "shid" で終わるときも、JavaScript の process.exit と同じく実行しない。
func abandons(err error) bool {
	_, isBacktrack := errors.AsType[*backtrack](err)
	_, isExit := errors.AsType[*Exit](err)
	return isBacktrack || isExit || errors.Is(err, errAbandoned)
}

// catchable は youlugive で捕まえられるエラーなら、例外の辞書にして返す。
func catchable(err error) (*Dict, bool) {
	if th, ok := errors.AsType[*thrown](err); ok {
		return th.value, true
	}
	if de, ok := errors.AsType[*diag.Error](err); ok && de.Kind.Catchable() {
		return exception(de.Error(), de.Kind.Spelling()), true
	}
	return nil, false
}

// exception は e no "lyuzun" と e no "chidd" で中身を読める例外の辞書を作る。
func exception(reason, kind string) *Dict {
	d := NewDict()
	d.Set("lyuzun", Str(reason))
	d.Set("chidd", Str(kind))
	return d
}

func (in *Interp) throw(s *ast.Throw, env *Env) error {
	v, err := in.eval(s.X, env)
	if err != nil {
		return err
	}
	switch v := v.(type) {
	case Str:
		return &thrown{value: exception(string(v), diag.AccUhEarts.Spelling()), file: in.frame().file, pos: s.Pos}
	case *Dict:
		// youlugive で受け取った例外を投げ直す
		if _, ok := v.Get("lyuzun"); ok {
			return &thrown{value: v, file: in.frame().file, pos: s.Pos}
		}
	}
	return in.errorf(diag.ImJune, s.Pos, lexicon.Spell("accuhearts"))
}

func (in *Interp) assign(s *ast.Assign, env *Env) error {
	value, err := in.eval(s.Value, env)
	if err != nil {
		return err
	}
	switch t := s.Target.(type) {
	case *ast.Ident:
		return in.assignName(s, t, value, env)
	case *ast.Index:
		return in.assignIndex(s, t, value, env)
	case *ast.Slice:
		return in.assignSlice(s, t, value, env)
	}
	return in.errorf(diag.SaiVI, s.Pos, "")
}

func (in *Interp) assignName(s *ast.Assign, t *ast.Ident, value Value, env *Env) error {
	v, res := env.resolveAssign(t.Name)
	if s.Const {
		if res != assignNew && !env.inherited[t.Name] {
			return in.errorf(diag.SuReady, t.Pos, t.Name)
		}
		env.define(t.Name, value, true)
		return nil
	}
	switch res {
	case assignAcrossFunc:
		return in.errorf(diag.SuReady, t.Pos, t.Name)
	case assignNew:
		if s.Op != "" {
			return in.errorf(diag.NiiExi, t.Pos, t.Name)
		}
		env.define(t.Name, value, false)
		return nil
	}
	if v.isConst {
		return in.errorf(diag.BlasTOac, t.Pos, t.Name)
	}
	if s.Op != "" {
		var err error
		if value, err = in.binary(s.Op, v.value, value, s.Pos); err != nil {
			return err
		}
	}
	v.value = value
	return nil
}

func (in *Interp) assignIndex(s *ast.Assign, t *ast.Index, value Value, env *Env) error {
	container, err := in.eval(t.X, env)
	if err != nil {
		return err
	}
	if _, ok := container.(*Module); ok {
		// モジュールの中身は読み込んだ側からは書き換えられない (JavaScript の import と同じ)
		return in.errorf(diag.BlasTOac, t.Pos, callName(t.Key))
	}
	key, err := in.eval(t.Key, env)
	if err != nil {
		return err
	}
	if s.Op != "" {
		old, err := in.index(container, key, t.Pos)
		if err != nil {
			return err
		}
		if value, err = in.binary(s.Op, old, value, s.Pos); err != nil {
			return err
		}
	}
	switch c := container.(type) {
	case *Array:
		i, ok := toIndex(key)
		if !ok {
			return in.errorf(diag.ImJune, t.Pos, Display(key))
		}
		n := int64(len(c.Elems))
		if i < 0 {
			i += n
		}
		if i < 0 || i >= n {
			return in.errorf(diag.NiiExi, t.Pos, Display(key))
		}
		c.Elems[i] = value
		return nil
	case *Dict:
		k, ok := key.(Str)
		if !ok {
			return in.errorf(diag.ImJune, t.Pos, Display(key))
		}
		c.Set(string(k), value)
		return nil
	}
	// 文字列は書き換えられない
	return in.errorf(diag.ImJune, t.Pos, lexicon.Spell("no"))
}

func (in *Interp) assignSlice(s *ast.Assign, t *ast.Slice, value Value, env *Env) error {
	if s.Op != "" {
		return in.errorf(diag.ImJune, s.Pos, s.Op+lexicon.Spell("izu"))
	}
	container, err := in.eval(t.X, env)
	if err != nil {
		return err
	}
	arr, ok := container.(*Array)
	repl, ok2 := value.(*Array)
	if !ok || !ok2 {
		return in.errorf(diag.ImJune, t.Pos, lexicon.Spell("no"))
	}
	lo, hi, err := in.sliceBounds(t, len(arr.Elems), env)
	if err != nil {
		return err
	}
	elems := append([]Value(nil), arr.Elems[:lo]...)
	elems = append(elems, repl.Elems...)
	arr.Elems = append(elems, arr.Elems[hi:]...)
	return nil
}
