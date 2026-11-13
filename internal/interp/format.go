package interp

import (
	"encoding/hex"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

// Display は redogh や #{} で出す形を返す。文字列はそのまま出す。
func Display(v Value) string {
	if s, ok := v.(Str); ok {
		return string(s)
	}
	return inner(v)
}

// inner は配列や辞書の中に入った値の形。文字列は " で囲む (JSON に近い形)。
func inner(v Value) string { return nested(v, nil) }

// nested は inner の中身。open は今たどっている途中の配列と辞書で、自分自身を中に持つものをもう一度たどらないために使う
// (Python と同じく [...] と {...} で出す)。
func nested(v Value, open []Value) string {
	switch v := v.(type) {
	case Int:
		return strconv.FormatInt(int64(v), 10)
	case Float:
		return formatFloat(float64(v))
	case Str:
		return quote(string(v))
	case Bool:
		if v {
			return lexicon.Spell("tadact")
		}
		return lexicon.Spell("ulieso")
	case Null:
		return lexicon.Spell("vanita")
	case *Array:
		if slices.Contains(open, Value(v)) {
			return "[...]"
		}
		open = append(open, v)
		parts := make([]string, len(v.Elems))
		for i, e := range v.Elems {
			parts[i] = nested(e, open)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case *Dict:
		if slices.Contains(open, Value(v)) {
			return "{...}"
		}
		open = append(open, v)
		parts := make([]string, 0, v.Len())
		for _, k := range v.keys {
			parts = append(parts, quote(k)+": "+nested(v.values[k], open))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case Bytes:
		return lexicon.Spell("sanand") + "(" + quote(hex.EncodeToString(v)) + ", " + quote("hex") + ")"
	case *Range:
		op := "..."
		if v.Inclusive {
			op = ".."
		}
		s := strconv.FormatInt(v.Lo, 10) + op + strconv.FormatInt(v.Hi, 10)
		if v.Step != 1 {
			s += " " + lexicon.Spell("batte") + " " + strconv.FormatInt(v.Step, 10)
		}
		return s
	case *Func:
		name := lexicon.Spell("weho")
		if v.Lit.Name != "" {
			name += " " + v.Lit.Name
		}
		return name
	case *Builtin:
		return lexicon.Spell("weho") + " " + v.Name
	case *Promise:
		return lexicon.Spell("itsver")
	case *Witch:
		return lexicon.Spell("ieviljio") + " " + strconv.Itoa(v.child.Pid())
	case *Module:
		return quote(v.Path) + " " + lexicon.Spell("werer")
	}
	return "?"
}

// formatFloat は JavaScript の Number の文字列化と同じ形にする (RFC：数は coword no salook と同じ形)。
func formatFloat(f float64) string {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "-Infinity"
	case f == 0:
		return "0"
	}
	sign := ""
	if f < 0 {
		sign, f = "-", -f
	}
	// 一番短い桁の並びと指数を取り出し、ECMAScript の Number::toString の規則で並べ直す
	e := strconv.FormatFloat(f, 'e', -1, 64)
	mantissa, expPart, _ := strings.Cut(e, "e")
	digits := strings.Replace(mantissa, ".", "", 1)
	exp, _ := strconv.Atoi(expPart)
	k, n := len(digits), exp+1

	switch {
	case k <= n && n <= 21:
		return sign + digits + strings.Repeat("0", n-k)
	case 0 < n && n <= 21:
		return sign + digits[:n] + "." + digits[n:]
	case -6 < n && n <= 0:
		return sign + "0." + strings.Repeat("0", -n) + digits
	}
	s := digits[:1]
	if k > 1 {
		s += "." + digits[1:]
	}
	expSign := "+"
	if n-1 < 0 {
		expSign = "-"
	}
	return sign + s + "e" + expSign + strconv.Itoa(abs(n-1))
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func quote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\t':
			b.WriteString(`\t`)
		case '\r':
			b.WriteString(`\r`)
		default:
			if r < 0x20 {
				b.WriteString(`\u{` + strconv.FormatInt(int64(r), 16) + `}`)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
	return b.String()
}
