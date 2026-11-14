package interp

import (
	"cmp"
	"math"
	"slices"
	"strings"
	"unicode/utf8"

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
		if v, ok, err := in.concat(a, b, pos); ok {
			return v, err
		}
	case "*":
		s, n, ok := strAndInt(a, b)
		if !ok {
			s, n, ok = strAndInt(b, a)
		}
		if ok {
			repeated, fits := repeat(s, n)
			if !fits {
				return nil, in.errorf(diag.Bavellabion, pos, op)
			}
			return repeated, nil
		}
	}
	return nil, in.errorf(diag.ImJune, pos, op)
}

func strAndInt(a, b Value) (Str, Int, bool) {
	s, isStr := a.(Str)
	n, isInt := b.(Int)
	return s, n, isStr && isInt
}

// maxStringBytes は、プログラムが作る文字列とバイト列の上限。超えると Bavellabion (メモリを取り切って落ちる前に止める)。
// JavaScript が長すぎる文字列で RangeError を投げるのと同じ役目。
const maxStringBytes = 1 << 30

// maxArrayLen は、つないで作る配列の要素数の上限。要素1つが16バイトなので、maxStringBytes と同じ大きさになる。
const maxArrayLen = maxStringBytes / 16

// concat は + でつなぐ。s izu s + s のように繰り返すと1回ごとに倍になるので、大きさを先に確かめる。
func (in *Interp) concat(a, b Value, pos token.Pos) (Value, bool, error) {
	switch a := a.(type) {
	case Str:
		if b, ok := b.(Str); ok {
			if len(a)+len(b) > maxStringBytes {
				return nil, true, in.errorf(diag.Bavellabion, pos, "+")
			}
			return a + b, true, nil
		}
	case *Array:
		if b, ok := b.(*Array); ok {
			if len(a.Elems)+len(b.Elems) > maxArrayLen {
				return nil, true, in.errorf(diag.Bavellabion, pos, "+")
			}
			return NewArray(slices.Concat(a.Elems, b.Elems)), true, nil
		}
	}
	return nil, false, nil
}

// text は Display と同じ形を返す。maxStringBytes を超えるなら、target を対象にした Bavellabion。
func (in *Interp) text(v Value, pos token.Pos, target string) (string, error) {
	s, ok := display(v)
	if !ok {
		return "", in.errorf(diag.Bavellabion, pos, target)
	}
	return s, nil
}

func repeat(s Str, n Int) (Str, bool) {
	if n <= 0 || s == "" {
		return "", true
	}
	if n > Int(maxStringBytes/len(s)) {
		return "", false
	}
	return Str(strings.Repeat(string(s), int(n))), true
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
			return ordered(op, cmp.Compare(x, y)), nil
		}
	}
	if x, ok := toFloat(a); ok {
		if y, ok := toFloat(b); ok {
			if math.IsNaN(x) || math.IsNaN(y) {
				return false, nil
			}
			return ordered(op, cmp.Compare(x, y)), nil
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

// equal は senaji。Python の == と同じだが、真偽値と数は別の型として扱う (spec/keywords.md の等しさ)。
// 配列と辞書はいくらでも深く入れ子にできるので、再帰ではなく、これから比べる組を積んでおいて順に比べる。
func equal(a, b Value) bool {
	var todo [][2]Value
	// seen はもう比べ始めた配列と辞書の組。自分自身を中に持つものをもう一度比べ始めたら、
	// ほかの要素で決まるものとして飛ばす (比べ続けると終わらないため)
	var seen map[[2]Value]bool
	begin := func(pair [2]Value) bool {
		if pair[0] == pair[1] || seen[pair] {
			return false
		}
		if seen == nil {
			seen = map[[2]Value]bool{}
		}
		seen[pair] = true
		return true
	}
	for {
		switch x := a.(type) {
		case *Array:
			y, ok := b.(*Array)
			if !ok || len(x.Elems) != len(y.Elems) {
				return false
			}
			if begin([2]Value{x, y}) {
				for i := range x.Elems {
					todo = append(todo, [2]Value{x.Elems[i], y.Elems[i]})
				}
			}
		case *Dict:
			y, ok := b.(*Dict)
			if !ok || x.Len() != y.Len() {
				return false
			}
			if begin([2]Value{x, y}) {
				for _, k := range x.keys {
					yv, ok := y.values[k]
					if !ok {
						return false
					}
					todo = append(todo, [2]Value{x.values[k], yv})
				}
			}
		default:
			if !equalAtoms(a, b) {
				return false
			}
		}
		if len(todo) == 0 {
			return true
		}
		a, b = todo[len(todo)-1][0], todo[len(todo)-1][1]
		todo = todo[:len(todo)-1]
	}
}

// equalAtoms は配列と辞書以外の値を比べる。
func equalAtoms(a, b Value) bool {
	if x, ok := toFloat(a); ok {
		y, ok := toFloat(b)
		if !ok {
			return false
		}
		if xi, ok := a.(Int); ok {
			if yi, ok := b.(Int); ok {
				return xi == yi
			}
		}
		return x == y
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
		return c.contains(i) && c.onStep(i), nil
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

// onStep は i が始まりから刻みの整数倍だけ離れているかを返す。差は int64 に収まらないことがあるので、uint64 で数える。
func (r *Range) onStep(i int64) bool {
	if r.Step < 0 {
		return (uint64(r.Lo)-uint64(i))%-uint64(r.Step) == 0
	}
	return (uint64(i)-uint64(r.Lo))%uint64(r.Step) == 0
}

// next は i の次の要素を返す。int64 からあふれるなら、もう次はない。
func (r *Range) next(i int64) (int64, bool) {
	n := i + r.Step
	return n, (n > i) == (r.Step > 0) && r.contains(n)
}

// position は添字 key を、長さ n の並びの中の位置にする。負の数は後ろから数える (xs no -1)。
// 添字に使えない値なら imJune、並びの外なら Nii eXi'。
func (in *Interp) position(key Value, n int, pos token.Pos) (int, error) {
	i, ok := toIndex(key)
	if !ok {
		return 0, in.errorf(diag.ImJune, pos, Display(key))
	}
	if i < 0 {
		i += int64(n)
	}
	if i < 0 || i >= int64(n) {
		return 0, in.errorf(diag.NiiExi, pos, Display(key))
	}
	return int(i), nil
}

// index は A no B。
func (in *Interp) index(container, key Value, pos token.Pos) (Value, error) {
	switch c := container.(type) {
	case *Array:
		i, err := in.position(key, len(c.Elems), pos)
		if err != nil {
			return nil, err
		}
		return c.Elems[i], nil
	case Str:
		i, err := in.position(key, utf8.RuneCountInString(string(c)), pos)
		if err != nil {
			return nil, err
		}
		for _, r := range string(c) {
			if i == 0 {
				return Str(string(r)), nil
			}
			i--
		}
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
