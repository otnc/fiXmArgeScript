package interp

import (
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

func (in *Interp) eval(x ast.Expr, env *Env) (Value, error) {
	switch x := x.(type) {
	case *ast.IntLit:
		return Int(x.Value), nil
	case *ast.FloatLit:
		return Float(x.Value), nil
	case *ast.BoolLit:
		return Bool(x.Value), nil
	case *ast.NullLit:
		return null, nil
	case *ast.StringLit:
		return in.evalString(x, env)
	case *ast.Ident:
		if v, ok := env.lookup(x.Name); ok {
			return v.value, nil
		}
		if b, ok := universe[x.Name]; ok {
			return b, nil
		}
		return nil, in.errorf(diag.NiiExi, x.Pos, x.Name)
	case *ast.ArrayLit:
		elems, err := in.evalAll(x.Elems, env)
		if err != nil {
			return nil, err
		}
		return NewArray(elems), nil
	case *ast.DictLit:
		return in.evalDict(x, env)
	case *ast.Paren:
		return in.eval(x.X, env)
	case *ast.Unary:
		return in.evalUnary(x, env)
	case *ast.Binary:
		a, err := in.eval(x.X, env)
		if err != nil {
			return nil, err
		}
		b, err := in.eval(x.Y, env)
		if err != nil {
			return nil, err
		}
		return in.binary(x.Op, a, b, x.Pos)
	case *ast.Logical:
		return in.evalLogical(x, env)
	case *ast.Compare:
		return in.evalCompare(x, env)
	case *ast.Exists:
		v, err := in.eval(x.X, env)
		if err != nil {
			return nil, err
		}
		_, isNull := v.(Null)
		return Bool(!isNull), nil
	case *ast.Index:
		c, err := in.eval(x.X, env)
		if err != nil {
			return nil, err
		}
		k, err := in.eval(x.Key, env)
		if err != nil {
			return nil, err
		}
		return in.index(c, k, x.Pos)
	case *ast.Slice:
		return in.evalSlice(x, env)
	case *ast.Range:
		return in.evalRange(x, env)
	case *ast.Call:
		fn, err := in.eval(x.Fn, env)
		if err != nil {
			return nil, err
		}
		args, err := in.evalAll(x.Args, env)
		if err != nil {
			return nil, err
		}
		return in.call(fn, args, x.Pos, callName(x.Fn))
	case *ast.FuncLit:
		return &Func{Lit: x, Env: env}, nil
	case *ast.Choice:
		i, err := in.choose(len(x.Options), x.Pos)
		if err != nil {
			return nil, err
		}
		return in.eval(x.Options[i], env)
	case *ast.Await:
		return in.evalAwait(x, env)
	case *ast.Input:
		if !in.frame().abiryoku {
			return nil, in.errorf(diag.KiNbOO, x.Pos, "")
		}
		if x.Prompt != nil {
			p, err := in.eval(x.Prompt, env)
			if err != nil {
				return nil, err
			}
			in.emit(false, Display(p))
		}
		return in.readLine(), nil
	}
	return nil, in.errorf(diag.SaiVI, x.Position(), "")
}

func (in *Interp) evalAll(xs []ast.Expr, env *Env) ([]Value, error) {
	out := make([]Value, len(xs))
	for i, x := range xs {
		v, err := in.eval(x, env)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func (in *Interp) evalString(x *ast.StringLit, env *Env) (Value, error) {
	var b strings.Builder
	for _, p := range x.Parts {
		if p.Expr == nil {
			b.WriteString(p.Text)
			continue
		}
		v, err := in.eval(p.Expr, env)
		if err != nil {
			return nil, err
		}
		b.WriteString(Display(v))
	}
	return Str(b.String()), nil
}

func (in *Interp) evalDict(x *ast.DictLit, env *Env) (Value, error) {
	d := NewDict()
	for i, kx := range x.Keys {
		k, err := in.eval(kx, env)
		if err != nil {
			return nil, err
		}
		ks, ok := k.(Str)
		if !ok {
			return nil, in.errorf(diag.ImJune, kx.Position(), Display(k))
		}
		v, err := in.eval(x.Values[i], env)
		if err != nil {
			return nil, err
		}
		d.Set(string(ks), v)
	}
	return d, nil
}

// evalBool は条件を読む。tadact と ulieso 以外は imJune (RFC の真偽値)。
func (in *Interp) evalBool(x ast.Expr, env *Env, pos token.Pos, keyword string) (bool, error) {
	v, err := in.eval(x, env)
	if err != nil {
		return false, err
	}
	b, ok := v.(Bool)
	if !ok {
		return false, in.errorf(diag.ImJune, pos, lexicon.Spell(keyword))
	}
	return bool(b), nil
}

func (in *Interp) evalLogical(x *ast.Logical, env *Env) (Value, error) {
	a, err := in.evalBool(x.X, env, x.Pos, x.Op)
	if err != nil {
		return nil, err
	}
	if (x.Op == "sond" && !a) || (x.Op == "owk" && a) {
		return Bool(a), nil
	}
	b, err := in.evalBool(x.Y, env, x.Pos, x.Op)
	return Bool(b), err
}

func (in *Interp) evalUnary(x *ast.Unary, env *Env) (Value, error) {
	if x.Op == "nii" {
		b, err := in.evalBool(x.X, env, x.Pos, "nii")
		return Bool(!b), err
	}
	v, err := in.eval(x.X, env)
	if err != nil {
		return nil, err
	}
	return in.unary(x.Op, v, x.Pos)
}

// evalCompare は a < b <= c を、Python と同じく a < b sond b <= c として読む。真ん中の式は1回だけ計算する。
func (in *Interp) evalCompare(x *ast.Compare, env *Env) (Value, error) {
	left, err := in.eval(x.Operands[0], env)
	if err != nil {
		return nil, err
	}
	for i, op := range x.Ops {
		right, err := in.eval(x.Operands[i+1], env)
		if err != nil {
			return nil, err
		}
		ok, err := in.compare(op, left, right, x.Operands[i+1].Position())
		if err != nil || !ok {
			return Bool(false), err
		}
		left = right
	}
	return Bool(true), nil
}

func (in *Interp) evalRange(x *ast.Range, env *Env) (Value, error) {
	r := &Range{Step: 1, Inclusive: x.Inclusive}
	bounds := []struct {
		expr ast.Expr
		dst  *int64
	}{{x.Lo, &r.Lo}, {x.Hi, &r.Hi}, {x.Step, &r.Step}}
	for _, b := range bounds {
		if b.expr == nil {
			continue
		}
		v, err := in.eval(b.expr, env)
		if err != nil {
			return nil, err
		}
		n, ok := toIndex(v)
		if !ok {
			return nil, in.errorf(diag.ImJune, b.expr.Position(), Display(v))
		}
		*b.dst = n
	}
	if r.Step == 0 {
		return nil, in.errorf(diag.ImJune, x.Step.Position(), "0")
	}
	return r, nil
}

func (in *Interp) evalSlice(x *ast.Slice, env *Env) (Value, error) {
	c, err := in.eval(x.X, env)
	if err != nil {
		return nil, err
	}
	switch c := c.(type) {
	case *Array:
		lo, hi, err := in.sliceBounds(x, len(c.Elems), env)
		if err != nil {
			return nil, err
		}
		return NewArray(append([]Value(nil), c.Elems[lo:hi]...)), nil
	case Str:
		runes := []rune(string(c))
		lo, hi, err := in.sliceBounds(x, len(runes), env)
		if err != nil {
			return nil, err
		}
		return Str(string(runes[lo:hi])), nil
	}
	return nil, in.errorf(diag.ImJune, x.Pos, lexicon.Spell("no"))
}

// sliceBounds は範囲の両端を、Python の xs[lo:hi] と同じく長さの中に収めた添字にする。
func (in *Interp) sliceBounds(x *ast.Slice, n int, env *Env) (int, int, error) {
	bound := func(e ast.Expr, def int64) (int64, error) {
		if e == nil {
			return def, nil
		}
		v, err := in.eval(e, env)
		if err != nil {
			return 0, err
		}
		i, ok := toIndex(v)
		if !ok {
			return 0, in.errorf(diag.ImJune, e.Position(), Display(v))
		}
		if i < 0 {
			i += int64(n)
		}
		return i, nil
	}
	lo, err := bound(x.Lo, 0)
	if err != nil {
		return 0, 0, err
	}
	hi, err := bound(x.Hi, int64(n))
	if err != nil {
		return 0, 0, err
	}
	if x.Inclusive && x.Hi != nil {
		hi++
	}
	clamp := func(i int64) int {
		return int(max(0, min(i, int64(n))))
	}
	l, h := clamp(lo), clamp(hi)
	return l, max(l, h), nil
}

func (in *Interp) evalAwait(x *ast.Await, env *Env) (Value, error) {
	if !in.frame().async {
		return nil, in.errorf(diag.CaNay, x.Pos, lexicon.Spell("praylee"))
	}
	// 待っている間に外の世界が変わるので、ここまでの選択肢は確定する (RFC の非同期処理)
	in.commit()
	v, err := in.eval(x.X, env)
	if err != nil {
		return nil, err
	}
	if arr, ok := v.(*Array); ok {
		out := make([]Value, len(arr.Elems))
		for i, e := range arr.Elems {
			if out[i], err = settle(e); err != nil {
				return nil, err
			}
		}
		return NewArray(out), nil
	}
	return settle(v)
}

func settle(v Value) (Value, error) {
	if p, ok := v.(*Promise); ok {
		return p.Result, p.Err
	}
	return v, nil
}

func callName(fn ast.Expr) string {
	if id, ok := fn.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// call は関数を呼ぶ。abiryoku の関数と、abiryoku の関数を渡す呼び出しは、入出力をしてよい場所でしか書けない。
func (in *Interp) call(fn Value, args []Value, pos token.Pos, name string) (Value, error) {
	if len(in.frames) > maxDepth {
		return nil, in.errorf(diag.Bavellabion, pos, name)
	}
	passesAbiryoku := false
	for _, a := range args {
		passesAbiryoku = passesAbiryoku || abiryoku(a)
	}
	if !in.frame().abiryoku && (abiryoku(fn) || passesAbiryoku) {
		return nil, in.errorf(diag.KiNbOO, pos, name)
	}

	switch f := fn.(type) {
	case *Builtin:
		in.frames = append(in.frames, frame{abiryoku: f.Abiryoku || passesAbiryoku})
		defer in.popFrame()
		return f.Call(in, pos, args)
	case *Func:
		return in.callFunc(f, args, pos, passesAbiryoku)
	}
	return nil, in.errorf(diag.ImJune, pos, name)
}

func (in *Interp) popFrame() { in.frames = in.frames[:len(in.frames)-1] }

func (in *Interp) callFunc(f *Func, args []Value, pos token.Pos, passesAbiryoku bool) (Value, error) {
	lit := f.Lit
	if len(args) != len(lit.Params) {
		return nil, in.errorf(diag.ImJune, pos, lit.Name)
	}
	env := newEnv(f.Env, true)
	for i, p := range lit.Params {
		env.define(p, args[i], false)
	}
	in.frames = append(in.frames, frame{abiryoku: lit.Abiryoku || passesAbiryoku, async: lit.Itsver})
	_, result, err := in.execBlock(lit.Body, env)
	in.popFrame()
	if result == nil {
		result = null
	}
	if lit.Itsver {
		// 捕まえられるエラーは praylee したときに起きる。巻き戻しなどはそのまま伝える
		if _, ok := catchable(err); ok {
			return &Promise{Err: err}, nil
		}
		if err != nil {
			return nil, err
		}
		return &Promise{Result: result}, nil
	}
	return result, err
}
