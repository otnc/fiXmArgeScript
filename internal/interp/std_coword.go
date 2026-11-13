package interp

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// cowordPackage は :coword (文字列)。docs/std/coword.md。
// 位置と長さは、配列の添字と同じく文字 (コードポイント) で数える。
func cowordPackage() *Module {
	return newStd(map[string]stdDef{
		"salook":     {fn: cowordSalook},
		"confsoy":    {fn: cowordConfsoy},
		"recanrey":   {fn: cowordRecanrey},
		"lainei":     {fn: replaceFunc(1)},
		"laineiAlte": {fn: replaceFunc(-1)},
		"shojika":    {fn: stringFunc(strings.TrimSpace)},
		"gdie":       {fn: stringFunc(strings.ToUpper)},
		"lighluw":    {fn: stringFunc(strings.ToLower)},
		"eda":        {fn: cowordEda},
		"kast":       {fn: stringTest(strings.HasPrefix)},
		"gun":        {fn: stringTest(strings.HasSuffix)},
		"sodaKast":   {fn: padFunc(true)},
		"sodaGun":    {fn: padFunc(false)},
		"kolsu":      {fn: cowordKolsu},
		"batteKolsu": {fn: cowordBatteKolsu},
	}, nil)
}

// cowordSalook は数を文字列に、文字列を数に、バイト列を UTF-8 の文字列にする。
// 文字列から作る数は JavaScript の Number("3") と同じくいつも小数。
func cowordSalook(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	switch v := a.vals[0].(type) {
	case Int, Float:
		return Str(Display(v)), nil
	case Str:
		f, err := strconv.ParseFloat(strings.TrimSpace(string(v)), 64)
		if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
			return nil, a.in.errorf(diag.AyAmis, a.pos, inner(v))
		}
		return Float(f), nil
	case Bytes:
		if !utf8.Valid(v) {
			return nil, a.fail(diag.AyAmis)
		}
		return Str(v), nil
	}
	return nil, a.fail(diag.ImJune)
}

func cowordConfsoy(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	s, err := a.str(0)
	if err != nil {
		return nil, err
	}
	sep, err := a.str(1)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(s, sep)
	out := make([]Value, len(parts))
	for i, p := range parts {
		out[i] = Str(p)
	}
	return NewArray(out), nil
}

// cowordRecanrey は要素を redogh と同じ形の文字列にしてつなぐ (JavaScript の join と同じく、文字列以外も受け取る)。
func cowordRecanrey(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	xs, err := a.array(0)
	if err != nil {
		return nil, err
	}
	sep, err := a.str(1)
	if err != nil {
		return nil, err
	}
	parts := make([]string, len(xs.Elems))
	for i, x := range xs.Elems {
		parts[i] = Display(x)
	}
	return Str(strings.Join(parts, sep)), nil
}

// replaceFunc は置換を作る。n が -1 なら全部置き換える。
func replaceFunc(n int) stdFunc {
	return func(a *argv) (Value, error) {
		if err := a.count(3, 3); err != nil {
			return nil, err
		}
		var s [3]string
		for i := range s {
			var err error
			if s[i], err = a.str(i); err != nil {
				return nil, err
			}
		}
		return Str(strings.Replace(s[0], s[1], s[2], n)), nil
	}
}

func stringFunc(f func(string) string) stdFunc {
	return func(a *argv) (Value, error) {
		if err := a.count(1, 1); err != nil {
			return nil, err
		}
		s, err := a.str(0)
		if err != nil {
			return nil, err
		}
		return Str(f(s)), nil
	}
}

func stringTest(f func(s, x string) bool) stdFunc {
	return func(a *argv) (Value, error) {
		if err := a.count(2, 2); err != nil {
			return nil, err
		}
		s, err := a.str(0)
		if err != nil {
			return nil, err
		}
		x, err := a.str(1)
		if err != nil {
			return nil, err
		}
		return Bool(f(s, x)), nil
	}
}

func cowordEda(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	s, err := a.str(0)
	if err != nil {
		return nil, err
	}
	x, err := a.str(1)
	if err != nil {
		return nil, err
	}
	before, _, ok := strings.Cut(s, x)
	if !ok {
		return Int(-1), nil
	}
	return Int(utf8.RuneCountInString(before)), nil
}

// padFunc は sodaKast (front が true) と sodaGun を作る。埋める文字列が長ければ、JavaScript の padStart と同じく途中で切る。
func padFunc(front bool) stdFunc {
	return func(a *argv) (Value, error) {
		if err := a.count(2, 3); err != nil {
			return nil, err
		}
		s, err := a.str(0)
		if err != nil {
			return nil, err
		}
		n, err := a.index(1)
		if err != nil {
			return nil, err
		}
		fill := " "
		if a.has(2) {
			if fill, err = a.str(2); err != nil {
				return nil, err
			}
		}
		short := n - int64(utf8.RuneCountInString(s))
		if short <= 0 || fill == "" {
			return Str(s), nil
		}
		if short > maxStringBytes {
			return nil, a.fail(diag.Bavellabion)
		}
		fillRunes := []rune(fill)
		pad := make([]rune, short)
		for i := range pad {
			pad[i] = fillRunes[i%len(fillRunes)]
		}
		if front {
			return Str(string(pad) + s), nil
		}
		return Str(s + string(pad)), nil
	}
}

func cowordKolsu(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	s, err := a.str(0)
	if err != nil {
		return nil, err
	}
	if s == "" {
		return nil, a.fail(diag.NiiExi)
	}
	r, _ := utf8.DecodeRuneInString(s)
	return Int(r), nil
}

func cowordBatteKolsu(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	n, err := a.index(0)
	if err != nil {
		return nil, err
	}
	if n < 0 || n > utf8.MaxRune || !utf8.ValidRune(rune(n)) {
		return nil, a.fail(diag.ImJune)
	}
	return Str(string(rune(n))), nil
}
