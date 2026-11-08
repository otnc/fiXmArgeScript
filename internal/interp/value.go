package interp

import (
	"slices"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Value は fxma の値。
type Value interface {
	value()
}

type (
	Int   int64
	Float float64
	Str   string
	Bool  bool
	Null  struct{}

	// Array と Dict は参照で扱う (JavaScript と同じ)。
	Array struct{ Elems []Value }
	Dict  struct {
		keys   []string // 入れた順
		values map[string]Value
	}
	// Bytes はバイト列 (sAnAnd)。
	Bytes []byte
	Range struct {
		Lo, Hi, Step int64
		Inclusive    bool
	}
	Func struct {
		Lit  *ast.FuncLit
		Env  *Env
		File string // 定義したファイル
	}
	Builtin struct {
		Name     string
		Abiryoku bool
		Itsver   bool
		Call     func(in *Interp, pos token.Pos, args []Value) (Value, error)
	}
	// Promise は itsver の関数を呼んだ結果と、外の処理の結果。settled になるまで waiters の task が待つ。
	Promise struct {
		settled bool
		Result  Value
		Err     error
		waiters []*task
		awaited bool // praylee されたか
	}
)

func (Int) value()      {}
func (Float) value()    {}
func (Str) value()      {}
func (Bool) value()     {}
func (Null) value()     {}
func (*Array) value()   {}
func (*Dict) value()    {}
func (Bytes) value()    {}
func (*Range) value()   {}
func (*Func) value()    {}
func (*Builtin) value() {}
func (*Promise) value() {}
func (*Module) value()  {}

var null = Null{}

func NewArray(elems []Value) *Array { return &Array{Elems: elems} }

func NewDict() *Dict { return &Dict{values: map[string]Value{}} }

func (d *Dict) Get(k string) (Value, bool) {
	v, ok := d.values[k]
	return v, ok
}

func (d *Dict) Set(k string, v Value) {
	if _, ok := d.values[k]; !ok {
		d.keys = append(d.keys, k)
	}
	d.values[k] = v
}

// Keys はキーを入れた順に返す。写しを返すので、回している途中で Set してもよい。
func (d *Dict) Keys() []string { return slices.Clone(d.keys) }

func (d *Dict) Len() int { return len(d.keys) }

// abiryoku は値が abiryoku 付きの関数かを返す。
func abiryoku(v Value) bool {
	switch f := v.(type) {
	case *Func:
		return f.Lit.Abiryoku
	case *Builtin:
		return f.Abiryoku
	}
	return false
}
