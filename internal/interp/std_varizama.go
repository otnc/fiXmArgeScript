package interp

import (
	"cmp"
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// varizamaPackage は :varizama (配列)。docs/std/varizama.md。
// naddy・vaniell・lea は渡した配列を書き換え、ほかは新しい配列を返す (JavaScript の push と toSorted の分け方)。
func varizamaPackage() *Module {
	return newStd(map[string]stdDef{
		"naddy":   {fn: varizamaNaddy},
		"muyon":   {fn: varizamaMuyon},
		"orei":    {fn: varizamaOrei},
		"tarip":   {fn: varizamaTarip},
		"onkhe":   {fn: varizamaOnkhe},
		"eda":     {fn: varizamaEda},
		"vaniell": {fn: varizamaVaniell},
		"lea":     {fn: varizamaLea},
		"matill":  {fn: varizamaMatill},
	}, nil)
}

func varizamaNaddy(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	xs.Elems = append(xs.Elems, a.vals[1])
	return null, nil
}

func varizamaMuyon(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	n := len(xs.Elems)
	out := make([]Value, n)
	for i, x := range xs.Elems {
		out[n-1-i] = x
	}
	return NewArray(out), nil
}

// varizamaOrei は安定な並べ替え。比べる関数を省くと、数は小さい順、文字列は文字の番号の順。
func varizamaOrei(a *argv) (Value, error) {
	if err := a.count(1, 2); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	compare := func(x, y Value) (int, error) { return defaultCompare(a, x, y) }
	if a.has(1) {
		f, err := a.fn(1)
		if err != nil {
			return nil, err
		}
		compare = func(x, y Value) (int, error) {
			r, err := a.call(f, x, y)
			if err != nil {
				return 0, err
			}
			c, ok := toFloat(r)
			if !ok {
				return 0, a.fail(diag.ImJune)
			}
			return cmp.Compare(c, 0), nil
		}
	} else if err := checkSortable(a, xs.Elems); err != nil {
		return nil, err
	}
	// 比べる関数は fxma の関数なので、エラーや巻き戻しを返すことがある。
	// 最初のエラーを覚えて、それより後は関数を呼ばずに並べ替えを終わらせる
	var sortErr error
	out := slices.Clone(xs.Elems)
	slices.SortStableFunc(out, func(x, y Value) int {
		if sortErr != nil {
			return 0
		}
		c, err := compare(x, y)
		sortErr = err
		return c
	})
	if sortErr != nil {
		return nil, sortErr
	}
	return NewArray(out), nil
}

// checkSortable は、比べる関数なしで並べられる配列 (数だけ、または文字列だけ) かを調べる。
func checkSortable(a *argv, xs []Value) error {
	numbers, strs := 0, 0
	for _, x := range xs {
		switch x.(type) {
		case Int, Float:
			numbers++
		case Str:
			strs++
		default:
			return a.fail(diag.ImJune)
		}
	}
	if numbers > 0 && strs > 0 {
		return a.fail(diag.ImJune)
	}
	return nil
}

func defaultCompare(a *argv, x, y Value) (int, error) {
	if xs, ok := x.(Str); ok {
		return strings.Compare(string(xs), string(y.(Str))), nil
	}
	return cmpNumbers(a, x, y)
}

func varizamaTarip(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	f, err := a.fn(1)
	if err != nil {
		return nil, err
	}
	out := make([]Value, 0, len(xs.Elems))
	for _, x := range xs.Elems {
		v, err := a.call(f, x)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return NewArray(out), nil
}

func varizamaOnkhe(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	f, err := a.fn(1)
	if err != nil {
		return nil, err
	}
	var out []Value
	for _, x := range xs.Elems {
		keep, err := a.callBool(f, x)
		if err != nil {
			return nil, err
		}
		if keep {
			out = append(out, x)
		}
	}
	return NewArray(out), nil
}

// varizamaEda は位置を探す。関数を渡すと、tadact を返す最初の要素の位置 (JavaScript の findIndex)。
func varizamaEda(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	target := a.vals[1]
	for i, x := range xs.Elems {
		found := equal(x, target)
		if isFunc(target) {
			if found, err = a.callBool(target, x); err != nil {
				return nil, err
			}
		}
		if found {
			return Int(i), nil
		}
	}
	return Int(-1), nil
}

// position は添字を配列の中の位置にする。負の数は後ろから数える (xs no -1 と同じ)。
func position(a *argv, i, n int64) (int, error) {
	if i < 0 {
		i += n
	}
	if i < 0 || i >= n {
		return 0, a.in.errorf(diag.NiiExi, a.pos, Display(a.vals[1]))
	}
	return int(i), nil
}

func varizamaVaniell(a *argv) (Value, error) {
	if err := a.count(1, 2); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	n := int64(len(xs.Elems))
	if !a.has(1) {
		if n == 0 {
			return nil, a.fail(diag.NiiExi)
		}
		last := xs.Elems[n-1]
		xs.Elems = xs.Elems[:n-1]
		return last, nil
	}
	raw, err := a.index(1)
	if err != nil {
		return nil, err
	}
	i, err := position(a, raw, n)
	if err != nil {
		return nil, err
	}
	v := xs.Elems[i]
	xs.Elems = append(xs.Elems[:i], xs.Elems[i+1:]...)
	return v, nil
}

func varizamaLea(a *argv) (Value, error) {
	if err := a.count(3, 3); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	raw, err := a.index(1)
	if err != nil {
		return nil, err
	}
	n := int64(len(xs.Elems))
	if raw < 0 {
		// 後ろから数えるときは、-1 が最後の要素の前 (JavaScript の splice と同じ)
		raw += n
	}
	if raw < 0 || raw > n {
		return nil, a.in.errorf(diag.NiiExi, a.pos, Display(a.vals[1]))
	}
	i := int(raw)
	xs.Elems = append(xs.Elems[:i], append([]Value{a.vals[2]}, xs.Elems[i:]...)...)
	return null, nil
}

func varizamaMatill(a *argv) (Value, error) {
	if err := a.count(3, 3); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	f, err := a.fn(1)
	if err != nil {
		return nil, err
	}
	acc := a.vals[2]
	for _, x := range xs.Elems {
		if acc, err = a.call(f, acc, x); err != nil {
			return nil, err
		}
	}
	return acc, nil
}
