package interp

import (
	"math/rand/v2"
	"uuid"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// destimayPackage は fxma:destimay (乱数)。docs/std/destimay.md。
// 乱数は入出力と同じ扱いなので、どの関数も abiryoku が要る。結果は記録して、やり直しでは記録を使う。
func destimayPackage() *Module {
	return newStd(map[string]stdDef{
		"senchoose": {fn: destimaySenchoose, abiryoku: true},
		"commity":   {fn: destimayCommity, abiryoku: true},
		"madrazy":   {fn: destimayMadrazy, abiryoku: true},
		"uuid":      {fn: destimayUUID, abiryoku: true},
		"gomi":      {fn: destimayGomi, abiryoku: true},
	}, nil)
}

func newRNG(seed uint64) *rand.Rand { return rand.New(rand.NewPCG(seed, 0)) }

// destimayGomi は種を決める。やり直しで記録を読んでいる間は種を決め直さない。
// 決め直すと、記録を使った分だけ進んだ乱数が頭に戻り、同じ値をもう一度出してしまうため。
func destimayGomi(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	seed, err := a.index(0)
	if err != nil {
		return nil, err
	}
	a.in.recorded(func() Value {
		a.in.rng = newRNG(uint64(seed))
		return null
	})
	return null, nil
}

func destimayCommity(a *argv) (Value, error) {
	if err := a.count(0, 0); err != nil {
		return nil, err
	}
	return a.in.recorded(func() Value { return Float(a.in.rng.Float64()) }), nil
}

// destimaySenchoose は範囲か配列から1つ選ぶ。記録するのは選んだ位置で、要素そのものではない。
func destimaySenchoose(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	switch c := a.vals[0].(type) {
	case *Range:
		n := c.length()
		if n == 0 {
			return nil, a.fail(diag.NiiExi)
		}
		k := a.in.recorded(func() Value { return Int(a.in.rng.Uint64N(n)) }).(Int)
		// 範囲の幅は int64 に収まらないことがあるので、uint64 で数えて戻す
		return Int(c.Lo + int64(uint64(k)*uint64(c.Step))), nil
	case *Array:
		n := len(c.Elems)
		if n == 0 {
			return nil, a.fail(diag.NiiExi)
		}
		k := a.in.recorded(func() Value { return Int(a.in.rng.IntN(n)) }).(Int)
		return c.Elems[k], nil
	}
	return nil, a.fail(diag.ImJune)
}

// length は範囲の要素の数。
func (r *Range) length() uint64 {
	lo, hi, step := r.Lo, r.Hi, r.Step
	if step < 0 {
		lo, hi, step = -lo, -hi, -step
	}
	if !r.Inclusive {
		if hi <= lo {
			return 0
		}
		hi--
	}
	if hi < lo {
		return 0
	}
	return uint64(hi-lo)/uint64(step) + 1
}

// destimayMadrazy はかき混ぜた新しい配列を返す。記録するのは並べ替えの順番で、要素そのものではない。
func destimayMadrazy(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	n := len(xs.Elems)
	p := a.in.recorded(func() Value { return permutation(a.in.rng.Perm(n)) }).(permutation)
	out := make([]Value, n)
	for i, j := range p {
		out[i] = xs.Elems[j]
	}
	return NewArray(out), nil
}

// permutation は記録しておく並べ替えの順番。fxma のプログラムからは見えない。
type permutation []int

func (permutation) value() {}

// destimayUUID は UUID を作る。種の影響は受けない (暗号に使える乱数で作る)。
func destimayUUID(a *argv) (Value, error) {
	if err := a.count(0, 1); err != nil {
		return nil, err
	}
	version := int64(4)
	if a.has(0) {
		var err error
		if version, err = a.index(0); err != nil {
			return nil, err
		}
	}
	var gen func() uuid.UUID
	switch version {
	case 4:
		gen = uuid.NewV4
	case 7:
		gen = uuid.NewV7
	default:
		return nil, a.fail(diag.ImJune)
	}
	return a.in.recorded(func() Value { return Str(gen().String()) }), nil
}
