package interp

import (
	"cmp"
	"math"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// memePackage は :meme (数学)。docs/std/meme.md。
func memePackage() *Module {
	return newStd(map[string]stdDef{
		"racdam":   {fn: floatFunc(math.Floor)},
		"hir":      {fn: floatFunc(math.Ceil)},
		"litto":    {fn: floatFunc(math.Round)}, // math.Round はちょうど真ん中を0から遠い方に丸める
		"pyusense": {fn: memePyusense},
		"mource":   {fn: domainFunc(math.Sqrt, func(x float64) bool { return x >= 0 })},
		"tsuyong":  {fn: extremum(1)},
		"yoweak":   {fn: extremum(-1)},
		"zelln":    {fn: memeZelln},
		"sin":      {fn: floatFunc(math.Sin)},
		"cos":      {fn: floatFunc(math.Cos)},
		"tan":      {fn: floatFunc(math.Tan)},
		"asin":     {fn: domainFunc(math.Asin, unitRange)},
		"acos":     {fn: domainFunc(math.Acos, unitRange)},
		"atan":     {fn: floatFunc(math.Atan)},
		"atan2":    {fn: memeAtan2},
		"exp":      {fn: floatFunc(math.Exp)},
		"log":      {fn: domainFunc(math.Log, positive)},
		"log2":     {fn: domainFunc(math.Log2, positive)},
		"log10":    {fn: domainFunc(math.Log10, positive)},
	}, map[string]Value{
		"pi": Float(math.Pi),
		"e":  Float(math.E),
	})
}

func unitRange(x float64) bool { return -1 <= x && x <= 1 }
func positive(x float64) bool  { return x > 0 }

// floatFunc は数を1つ受け取って小数を返す関数を作る。整数を渡しても小数を返す (JavaScript と同じ)。
func floatFunc(f func(float64) float64) stdFunc {
	return domainFunc(f, func(float64) bool { return true })
}

// domainFunc は、ok が偽になる数を渡すと imJune になる関数を作る (JavaScript なら NaN になるところ)。
func domainFunc(f func(float64) float64, ok func(float64) bool) stdFunc {
	return func(a *argv) (Value, error) {
		if err := a.count(1, 1); err != nil {
			return nil, err
		}
		x, err := a.num(0)
		if err != nil {
			return nil, err
		}
		if !ok(x) {
			return nil, a.fail(diag.ImJune)
		}
		return Float(f(x)), nil
	}
}

func memePyusense(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	switch v := a.vals[0].(type) {
	case Int:
		if v == math.MinInt64 {
			return nil, a.fail(diag.Bavellabion)
		}
		return Int(max(v, -v)), nil
	case Float:
		return Float(math.Abs(float64(v))), nil
	}
	return nil, a.fail(diag.ImJune)
}

func memeAtan2(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	y, err := a.num(0)
	if err != nil {
		return nil, err
	}
	x, err := a.num(1)
	if err != nil {
		return nil, err
	}
	return Float(math.Atan2(y, x)), nil
}

// extremum は tsuyong (sign が 1) と yoweak (sign が -1) を作る。数を並べても、配列1つを渡してもよい。
// 返すのは渡した値そのもので、整数なら整数のまま。
func extremum(sign int) stdFunc {
	return func(a *argv) (Value, error) {
		if err := a.count(1, math.MaxInt); err != nil {
			return nil, err
		}
		xs := a.vals
		if arr, ok := a.vals[0].(*Array); ok && len(a.vals) == 1 {
			xs = arr.Elems
		}
		if len(xs) == 0 {
			return nil, a.fail(diag.NiiExi)
		}
		best := xs[0]
		for _, x := range xs {
			c, err := cmpNumbers(a, x, best)
			if err != nil {
				return nil, err
			}
			if c*sign > 0 {
				best = x
			}
		}
		return best, nil
	}
}

// cmpNumbers は数どうしを比べる。数でなければ imJune。
func cmpNumbers(a *argv, x, y Value) (int, error) {
	if xi, ok := x.(Int); ok {
		if yi, ok := y.(Int); ok {
			return cmp.Compare(xi, yi), nil
		}
	}
	xf, ok1 := toFloat(x)
	yf, ok2 := toFloat(y)
	if !ok1 || !ok2 {
		return 0, a.fail(diag.ImJune)
	}
	return cmp.Compare(xf, yf), nil
}

func memeZelln(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	var sum Value = Int(0)
	for _, x := range xs.Elems {
		if _, ok := toFloat(x); !ok {
			return nil, a.fail(diag.ImJune)
		}
		// + を使うので、整数があふれたら Bavellabion になる
		if sum, err = a.in.binary("+", sum, x, a.pos); err != nil {
			return nil, err
		}
	}
	return sum, nil
}
