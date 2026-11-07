package interp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// jsonPackage は fxma:json。docs/std/json.md。
func jsonPackage() *Module {
	return newStd(map[string]stdDef{
		"fineclowt": {fn: jsonFineclowt},
		"salook":    {fn: jsonSalook},
	}, nil)
}

// errNotJSON は JSON にできない値を見つけたこと。
var errNotJSON = errors.New("not json")

// jsonFineclowt は JSON を読む。辞書のキーを書かれた順に並べるため、json.Unmarshal ではなく字句を順に読む。
func jsonFineclowt(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	s, err := a.str(0)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	v, err := decodeJSON(dec)
	if err == nil {
		// 値の後ろに何か残っていたら、JSON.parse と同じく正しくない
		if _, rest := dec.Token(); rest != io.EOF {
			err = errNotJSON
		}
	}
	if err != nil {
		return nil, a.fail(diag.AyAmis)
	}
	return v, nil
}

func decodeJSON(dec *json.Decoder) (Value, error) {
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := t.(type) {
	case json.Number:
		// JavaScript の JSON.parse と同じく、数はいつも小数
		f, err := strconv.ParseFloat(string(t), 64)
		if err != nil {
			return nil, err
		}
		return Float(f), nil
	case string:
		return Str(t), nil
	case bool:
		return Bool(t), nil
	case nil:
		return null, nil
	case json.Delim:
		if t == '[' {
			return decodeJSONArray(dec)
		}
		return decodeJSONObject(dec)
	}
	return nil, errNotJSON
}

func decodeJSONArray(dec *json.Decoder) (Value, error) {
	var elems []Value
	for dec.More() {
		v, err := decodeJSON(dec)
		if err != nil {
			return nil, err
		}
		elems = append(elems, v)
	}
	_, err := dec.Token() // ]
	return NewArray(elems), err
}

func decodeJSONObject(dec *json.Decoder) (Value, error) {
	d := NewDict()
	for dec.More() {
		k, err := dec.Token()
		if err != nil {
			return nil, err
		}
		v, err := decodeJSON(dec)
		if err != nil {
			return nil, err
		}
		d.Set(k.(string), v)
	}
	_, err := dec.Token() // }
	return d, err
}

// jsonSalook は JSON.stringify と同じ形で書き出す。字下げは JavaScript と同じく 10 までで、1 未満なら1行。
func jsonSalook(a *argv) (Value, error) {
	if err := a.count(1, 2); err != nil {
		return nil, err
	}
	indent := ""
	if a.has(1) {
		n, err := a.index(1)
		if err != nil {
			return nil, err
		}
		indent = strings.Repeat(" ", int(max(0, min(n, 10))))
	}
	var b strings.Builder
	if err := encodeJSON(&b, a.vals[0], indent, ""); err != nil {
		return nil, a.fail(diag.ImJune)
	}
	return Str(b.String()), nil
}

func encodeJSON(b *strings.Builder, v Value, indent, prefix string) error {
	switch v := v.(type) {
	case Int:
		b.WriteString(strconv.FormatInt(int64(v), 10))
	case Float:
		f := float64(v)
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return errNotJSON
		}
		b.WriteString(formatFloat(f))
	case Str:
		b.WriteString(jsonQuote(string(v)))
	case Bool:
		b.WriteString(strconv.FormatBool(bool(v)))
	case Null:
		b.WriteString("null")
	case *Array:
		return encodeJSONList(b, '[', ']', len(v.Elems), indent, prefix, func(i int, inner string) error {
			return encodeJSON(b, v.Elems[i], indent, inner)
		})
	case *Dict:
		colon := ":"
		if indent != "" {
			colon = ": "
		}
		return encodeJSONList(b, '{', '}', v.Len(), indent, prefix, func(i int, inner string) error {
			k := v.keys[i]
			b.WriteString(jsonQuote(k) + colon)
			return encodeJSON(b, v.values[k], indent, inner)
		})
	default:
		return errNotJSON
	}
	return nil
}

// encodeJSONList は配列と辞書の、括弧・区切り・字下げを書く。空なら JavaScript と同じく [] と {}。
func encodeJSONList(b *strings.Builder, open, close byte, n int, indent, prefix string, item func(i int, inner string) error) error {
	b.WriteByte(open)
	inner := prefix + indent
	for i := range n {
		if i > 0 {
			b.WriteByte(',')
		}
		if indent != "" {
			b.WriteString("\n" + inner)
		}
		if err := item(i, inner); err != nil {
			return err
		}
	}
	if indent != "" && n > 0 {
		b.WriteString("\n" + prefix)
	}
	b.WriteByte(close)
	return nil
}

// jsonQuote は JSON.stringify と同じく、" と \ と制御文字だけをエスケープする。
func jsonQuote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(&b, `\u%04x`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
	return b.String()
}
