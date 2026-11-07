package interp

import (
	"math"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

func (in *Interp) unary(op string, v Value, pos token.Pos) (Value, error) {
	switch v := v.(type) {
	case Int:
		switch op {
		case "-":
			if v == math.MinInt64 {
				return nil, in.errorf(diag.Bavellabion, pos, op)
			}
			return -v, nil
		case "+":
			return v, nil
		case "~":
			return ^v, nil
		}
	case Float:
		switch op {
		case "-":
			return -v, nil
		case "+":
			return v, nil
		}
	}
	return nil, in.errorf(diag.ImJune, pos, op)
}

// binary は算術とビットの演算。数の扱いは Python と同じで、整数があふれたら Bavellabion。
func (in *Interp) binary(op string, a, b Value, pos token.Pos) (Value, error) {
	if x, ok := a.(Int); ok {
		if y, ok := b.(Int); ok {
			return in.intOp(op, int64(x), int64(y), pos)
		}
	}
	x, aNum := toFloat(a)
	y, bNum := toFloat(b)
	if aNum && bNum {
		return in.floatOp(op, x, y, pos)
	}
	switch op {
	case "+":
		switch a := a.(type) {
		case Str:
			if b, ok := b.(Str); ok {
				return a + b, nil
			}
		case *Array:
			if b, ok := b.(*Array); ok {
				elems := append(append([]Value(nil), a.Elems...), b.Elems...)
				return NewArray(elems), nil
			}
		}
	case "*":
		if s, ok := a.(Str); ok {
			if n, ok := b.(Int); ok {
				return repeat(s, n), nil
			}
		}
		if s, ok := b.(Str); ok {
			if n, ok := a.(Int); ok {
				return repeat(s, n), nil
			}
		}
	}
	return nil, in.errorf(diag.ImJune, pos, op)
}

func repeat(s Str, n Int) Str {
	if n <= 0 {
		return ""
	}
	return Str(strings.Repeat(string(s), int(n)))
}

func (in *Interp) intOp(op string, x, y int64, pos token.Pos) (Value, error) {
	overflow := func() (Value, error) { return nil, in.errorf(diag.Bavellabion, pos, op) }
	zero := func() (Value, error) { return nil, in.errorf(diag.KYoNAihim, pos, op) }
	switch op {
	case "+":
		s := x + y
		if (s > x) != (y > 0) {
			return overflow()
		}
		return Int(s), nil
	case "-":
		d := x - y
		if (d < x) != (y > 0) {
			return overflow()
		}
		return Int(d), nil
	case "*":
		p, ok := mulInt(x, y)
		if !ok {
			return overflow()
		}
		return Int(p), nil
	case "/":
		if y == 0 {
			return zero()
		}
		return Float(float64(x) / float64(y)), nil
	case "//", "%":
		if y == 0 {
			return zero()
		}
		if x == math.MinInt64 && y == -1 {
			if op == "%" {
				return Int(0), nil
			}
			return overflow()
		}
		q, r := x/y, x%y
		// Python と同じく、商は小さい方に切り捨て、余りは割る数と同じ符号にする
		if r != 0 && (r < 0) != (y < 0) {
			q--
			r += y
		}
		if op == "//" {
			return Int(q), nil
		}
		return Int(r), nil
	case "**":
		return in.intPow(x, y, pos)
	case "&":
		return Int(x & y), nil
	case "|":
		return Int(x | y), nil
	case "^":
		return Int(x ^ y), nil
	case "<<":
		if y < 0 {
			return nil, in.errorf(diag.ImJune, pos, op)
		}
		if y >= 64 || (x<<y)>>y != x {
			if x == 0 {
				return Int(0), nil
			}
			return overflow()
		}
		return Int(x << y), nil
	case ">>":
		if y < 0 {
			return nil, in.errorf(diag.ImJune, pos, op)
		}
		if y >= 64 {
			y = 63
		}
		return Int(x >> y), nil
	}
	return nil, in.errorf(diag.ImJune, pos, op)
}

// mulInt は掛け算をして、64ビットに収まったかを返す。
func mulInt(x, y int64) (int64, bool) {
	if x == 0 || y == 0 {
		return 0, true
	}
	p := x * y
	if p/y != x || (x == -1 && y == math.MinInt64) || (y == -1 && x == math.MinInt64) {
		return 0, false
	}
	return p, true
}

func (in *Interp) intPow(x, y int64, pos token.Pos) (Value, error) {
	if y < 0 {
		if x == 0 {
			return nil, in.errorf(diag.KYoNAihim, pos, "**")
		}
		return Float(math.Pow(float64(x), float64(y))), nil
	}
	result, base := int64(1), x
	for y > 0 {
		var ok bool
		if y&1 == 1 {
			if result, ok = mulInt(result, base); !ok {
				return nil, in.errorf(diag.Bavellabion, pos, "**")
			}
		}
		y >>= 1
		if y > 0 {
			if base, ok = mulInt(base, base); !ok {
				return nil, in.errorf(diag.Bavellabion, pos, "**")
			}
		}
	}
	return Int(result), nil
}

func (in *Interp) floatOp(op string, x, y float64, pos token.Pos) (Value, error) {
	switch op {
	case "+":
		return Float(x + y), nil
	case "-":
		return Float(x - y), nil
	case "*":
		return Float(x * y), nil
	case "/":
		if y == 0 {
			return nil, in.errorf(diag.KYoNAihim, pos, op)
		}
		return Float(x / y), nil
	case "//":
		if y == 0 {
			return nil, in.errorf(diag.KYoNAihim, pos, op)
		}
		return Float(math.Floor(x / y)), nil
	case "%":
		if y == 0 {
			return nil, in.errorf(diag.KYoNAihim, pos, op)
		}
		r := math.Mod(x, y)
		if r != 0 && (r < 0) != (y < 0) {
			r += y
		}
		return Float(r), nil
	case "**":
		if x == 0 && y < 0 {
			return nil, in.errorf(diag.KYoNAihim, pos, op)
		}
		return Float(math.Pow(x, y)), nil
	}
	return nil, in.errorf(diag.ImJune, pos, op)
}

func toFloat(v Value) (float64, bool) {
	switch v := v.(type) {
	case Int:
		return float64(v), true
	case Float:
		return float64(v), true
	}
	return 0, false
}

// toIndex は添字に使える数を整数にする。2.0 のように端数のない小数も使える (JavaScript の xs[2.0])。
func toIndex(v Value) (int64, bool) {
	switch v := v.(type) {
	case Int:
		return int64(v), true
	case Float:
		f := float64(v)
		if f == math.Trunc(f) && math.Abs(f) < 1<<62 {
			return int64(f), true
		}
	}
	return 0, false
}

// compare は比較の演算子1つ分を調べる。
func (in *Interp) compare(op string, a, b Value, pos token.Pos) (bool, error) {
	switch op {
	case "senaji":
		return equal(a, b), nil
	case "exi":
		return in.contains(b, a, pos)
	}
	if x, ok := a.(Int); ok {
		if y, ok := b.(Int); ok {
			return ordered(op, cmpInt(int64(x), int64(y))), nil
		}
	}
	if x, ok := toFloat(a); ok {
		if y, ok := toFloat(b); ok {
			if math.IsNaN(x) || math.IsNaN(y) {
				return false, nil
			}
			return ordered(op, cmpFloat(x, y)), nil
		}
	}
	if x, ok := a.(Str); ok {
		if y, ok := b.(Str); ok {
			// Go の文字列の比較は UTF-8 のバイト順で、これは文字の番号の順と同じ
			return ordered(op, strings.Compare(string(x), string(y))), nil
		}
	}
	return false, in.errorf(diag.ImJune, pos, op)
}

func ordered(op string, c int) bool {
	switch op {
	case "<":
		return c < 0
	case ">":
		return c > 0
	case "<=":
		return c <= 0
	}
	return c >= 0
}

func cmpInt(x, y int64) int {
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	}
	return 0
}

func cmpFloat(x, y float64) int {
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	}
	return 0
}

// equal は senaji。Python の == と同じだが、真偽値と数は別の型として扱う (RFC の等しさ)。
func equal(a, b Value) bool {
	if x, ok := toFloat(a); ok {
		if y, ok := toFloat(b); ok {
			if xi, ok := a.(Int); ok {
				if yi, ok := b.(Int); ok {
					return xi == yi
				}
			}
			return x == y
		}
		return false
	}
	switch a := a.(type) {
	case Str:
		b, ok := b.(Str)
		return ok && a == b
	case Bool:
		b, ok := b.(Bool)
		return ok && a == b
	case Null:
		_, ok := b.(Null)
		return ok
	case Bytes:
		b, ok := b.(Bytes)
		return ok && string(a) == string(b)
	case *Array:
		b, ok := b.(*Array)
		if !ok || len(a.Elems) != len(b.Elems) {
			return false
		}
		for i := range a.Elems {
			if !equal(a.Elems[i], b.Elems[i]) {
				return false
			}
		}
		return true
	case *Dict:
		b, ok := b.(*Dict)
		if !ok || a.Len() != b.Len() {
			return false
		}
		for _, k := range a.keys {
			bv, ok := b.values[k]
			if !ok || !equal(a.values[k], bv) {
				return false
			}
		}
		return true
	case *Range:
		b, ok := b.(*Range)
		return ok && *a == *b
	}
	// 関数などは同じものかどうかで比べる
	return a == b
}

// contains は x exi xs。
func (in *Interp) contains(container, x Value, pos token.Pos) (bool, error) {
	switch c := container.(type) {
	case *Array:
		for _, e := range c.Elems {
			if equal(e, x) {
				return true, nil
			}
		}
		return false, nil
	case Str:
		s, ok := x.(Str)
		if !ok {
			return false, in.errorf(diag.ImJune, pos, lexicon.Spell("exi"))
		}
		return strings.Contains(string(c), string(s)), nil
	case *Dict:
		s, ok := x.(Str)
		if !ok {
			return false, nil
		}
		_, has := c.Get(string(s))
		return has, nil
	case *Range:
		i, ok := toIndex(x)
		if !ok {
			return false, nil
		}
		return c.contains(i) && (i-c.Lo)%c.Step == 0, nil
	}
	return false, in.errorf(diag.ImJune, pos, lexicon.Spell("exi"))
}

// contains は i が範囲の端の内側にあるかを返す。刻みに乗っているかは見ない。
func (r *Range) contains(i int64) bool {
	if r.Step > 0 {
		if r.Inclusive {
			return r.Lo <= i && i <= r.Hi
		}
		return r.Lo <= i && i < r.Hi
	}
	if r.Inclusive {
		return r.Hi <= i && i <= r.Lo
	}
	return r.Hi < i && i <= r.Lo
}

// index は A no B。
func (in *Interp) index(container, key Value, pos token.Pos) (Value, error) {
	switch c := container.(type) {
	case *Array:
		i, ok := toIndex(key)
		if !ok {
			return nil, in.errorf(diag.ImJune, pos, Display(key))
		}
		n := int64(len(c.Elems))
		if i < 0 {
			i += n
		}
		if i < 0 || i >= n {
			return nil, in.errorf(diag.NiiExi, pos, Display(key))
		}
		return c.Elems[i], nil
	case Str:
		i, ok := toIndex(key)
		if !ok {
			return nil, in.errorf(diag.ImJune, pos, Display(key))
		}
		runes := []rune(string(c))
		n := int64(len(runes))
		if i < 0 {
			i += n
		}
		if i < 0 || i >= n {
			return nil, in.errorf(diag.NiiExi, pos, Display(key))
		}
		return Str(string(runes[i])), nil
	case *Dict:
		k, ok := key.(Str)
		if !ok {
			return nil, in.errorf(diag.NiiExi, pos, Display(key))
		}
		v, ok := c.Get(string(k))
		if !ok {
			return nil, in.errorf(diag.NiiExi, pos, string(k))
		}
		return v, nil
	}
	return nil, in.errorf(diag.ImJune, pos, lexicon.Spell("no"))
}
