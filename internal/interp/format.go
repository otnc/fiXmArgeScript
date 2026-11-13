package interp

import (
	"encoding/hex"
	"math"
	"strconv"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

// Display は redogh や #{} で出す形を返す。文字列はそのまま出す。
func Display(v Value) string {
	s, _ := display(v)
	return s
}

// display は Display の形と、それが maxStringBytes に収まったかを返す。収まらなければ途中までの形。
func display(v Value) (string, bool) {
	if s, ok := v.(Str); ok {
		return string(s), true
	}
	return nested(v)
}

// inner は配列や辞書の中に入った値の形。文字列は " で囲む (JSON に近い形)。
func inner(v Value) string {
	s, _ := nested(v)
	return s
}

// nested は inner の形と、それが maxStringBytes に収まったかを返す。
func nested(v Value) (string, bool) {
	switch v.(type) {
	case *Array, *Dict:
		return container(v)
	}
	return atom(v), true
}

// pendingText は container がこれから書くもの1つ。v が nil なら text をそのまま書き、closes があればその配列か辞書を書き終える。
type pendingText struct {
	v      Value
	text   string
	closes Value
}

// container は配列と辞書の形。自分自身を中に持つものは、Python と同じく [...] と {...} で出す。
// 入れ子はいくらでも深く作れるので、再帰ではなく、これから書くものを積んでおいて順に書く。
// 同じ配列を何度も入れた値 (x izu [x, x] の繰り返し) は、小さくても形は倍々に大きくなるので、maxStringBytes で打ち切って false を返す。
func container(root Value) (string, bool) {
	var b strings.Builder
	open := map[Value]bool{} // 今書いている途中の配列と辞書
	todo := []pendingText{{v: root}}
	for len(todo) > 0 {
		if b.Len() > maxStringBytes {
			return b.String(), false
		}
		next := todo[len(todo)-1]
		todo = todo[:len(todo)-1]
		switch v := next.v.(type) {
		case nil:
			b.WriteString(next.text)
			delete(open, next.closes)
		case *Array:
			if open[v] {
				b.WriteString("[...]")
				continue
			}
			open[v] = true
			b.WriteByte('[')
			todo = append(todo, pendingText{text: "]", closes: v})
			for i := len(v.Elems) - 1; i >= 0; i-- {
				todo = append(todo, pendingText{v: v.Elems[i]})
				if i > 0 {
					todo = append(todo, pendingText{text: ", "})
				}
			}
		case *Dict:
			if open[v] {
				b.WriteString("{...}")
				continue
			}
			open[v] = true
			b.WriteByte('{')
			todo = append(todo, pendingText{text: "}", closes: v})
			for i := len(v.keys) - 1; i >= 0; i-- {
				k := v.keys[i]
				todo = append(todo, pendingText{v: v.values[k]}, pendingText{text: quote(k) + ": "})
				if i > 0 {
					todo = append(todo, pendingText{text: ", "})
				}
			}
		default:
			b.WriteString(atom(v))
		}
	}
	return b.String(), b.Len() <= maxStringBytes
}

// atom は配列と辞書以外の値の形。
func atom(v Value) string {
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
	return sign + s + "e" + expSign + strconv.Itoa(max(n-1, 1-n))
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
