package interp

import (
	"encoding/json/jsontext"
	"errors"
	"io"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// jsonPackage は :json。docs/std/json.md。
func jsonPackage() *Module {
	return newStd(map[string]stdDef{
		"fineclowt": {fn: jsonFineclowt},
		"salook":    {fn: jsonSalook},
	}, nil)
}

func jsonFineclowt(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	s, err := a.str(0)
	if err != nil {
		return nil, err
	}
	v, err := parseJSON(s)
	if err != nil {
		return nil, a.fail(diag.AyAmis)
	}
	return v, nil
}

// parseJSON は JSON を読む。辞書のキーを書かれた順に並べるため、字句を順に読む。
// 同じキーが2回あれば、JSON.parse と同じく後の値を使う。
func parseJSON(s string) (Value, error) {
	dec := jsontext.NewDecoder(strings.NewReader(s), jsontext.AllowDuplicateNames(true))
	v, err := decodeJSON(dec)
	if err != nil {
		return nil, err
	}
	// 値の後ろに何か残っていたら、JSON.parse と同じく正しくない
	if _, rest := dec.ReadToken(); rest != io.EOF {
		return nil, errors.New("trailing data")
	}
	return v, nil
}

func decodeJSON(dec *jsontext.Decoder) (Value, error) {
	t, err := dec.ReadToken()
	if err != nil {
		return nil, err
	}
	switch t.Kind() {
	case '0':
		// JavaScript の JSON.parse と同じく、数はいつも小数
		f, err := t.Float()
		return Float(f), err
	case '"':
		return Str(t.String()), nil
	case 't', 'f':
		return Bool(t.Bool()), nil
	case 'n':
		return null, nil
	case '[':
		var elems []Value
		for dec.PeekKind() != ']' {
			v, err := decodeJSON(dec)
			if err != nil {
				return nil, err
			}
			elems = append(elems, v)
		}
		_, err := dec.ReadToken()
		return NewArray(elems), err
	}
	d := NewDict()
	for dec.PeekKind() != '}' {
		k, err := dec.ReadToken()
		if err != nil {
			return nil, err
		}
		// 字句は次に読むと使えなくなるので、先に文字列にしておく
		key := k.String()
		v, err := decodeJSON(dec)
		if err != nil {
			return nil, err
		}
		d.Set(key, v)
	}
	_, err = dec.ReadToken()
	return d, err
}

func jsonSalook(a *argv) (Value, error) {
	if err := a.count(1, 2); err != nil {
		return nil, err
	}
	s, err := jsonArg(a, 0, 1)
	if err != nil {
		return nil, err
	}
	return Str(s), nil
}

// jsonArg は i 番目の引数を JSON にする。indentAt 番目に字下げの幅があれば使う。
func jsonArg(a *argv, i, indentAt int) (string, error) {
	var indent int64
	if a.has(indentAt) {
		var err error
		if indent, err = a.index(indentAt); err != nil {
			return "", err
		}
	}
	s, err := formatJSON(a.vals[i], indent)
	if errors.Is(err, errTooLong) {
		return "", a.fail(diag.Bavellabion)
	}
	if err != nil {
		return "", a.fail(diag.ImJune)
	}
	return s, nil
}

var errTooLong = errors.New("too long")

// boundedText は maxStringBytes までしか書けない書き先。同じ配列を何度も入れた値は、JSON にすると倍々に大きくなる。
type boundedText struct{ strings.Builder }

func (b *boundedText) Write(p []byte) (int, error) {
	if b.Len()+len(p) > maxStringBytes {
		return 0, errTooLong
	}
	return b.Builder.Write(p)
}

// formatJSON は JSON.stringify と同じ形で書き出す。字下げは JavaScript と同じく 10 までで、1 未満なら1行。
func formatJSON(v Value, indent int64) (string, error) {
	opts := []jsontext.Options{jsontext.Multiline(false)}
	if indent >= 1 {
		opts = []jsontext.Options{jsontext.WithIndent(strings.Repeat(" ", int(min(indent, 10))))}
	}
	var b boundedText
	if err := encodeJSON(jsontext.NewEncoder(&b, opts...), v); err != nil {
		return "", err
	}
	// Encoder は値ごとに改行を付けるので外す
	return strings.TrimSuffix(b.String(), "\n"), nil
}

// encodeJSON は値を字句にして書く。関数やバイト列のように JSON にできない値はエラー。
func encodeJSON(enc *jsontext.Encoder, v Value) error {
	switch v := v.(type) {
	case Int:
		return enc.WriteToken(jsontext.Int(int64(v)))
	case Float:
		// NaN と無限大は JSON にできないので、jsontext がエラーを返す
		return enc.WriteToken(jsontext.Float(float64(v)))
	case Str:
		return enc.WriteToken(jsontext.String(string(v)))
	case Bool:
		return enc.WriteToken(jsontext.Bool(bool(v)))
	case Null:
		return enc.WriteToken(jsontext.Null)
	case *Array:
		if err := enc.WriteToken(jsontext.BeginArray); err != nil {
			return err
		}
		for _, e := range v.Elems {
			if err := encodeJSON(enc, e); err != nil {
				return err
			}
		}
		return enc.WriteToken(jsontext.EndArray)
	case *Dict:
		if err := enc.WriteToken(jsontext.BeginObject); err != nil {
			return err
		}
		for _, k := range v.keys {
			if err := enc.WriteToken(jsontext.String(k)); err != nil {
				return err
			}
			if err := encodeJSON(enc, v.values[k]); err != nil {
				return err
			}
		}
		return enc.WriteToken(jsontext.EndObject)
	}
	return errors.New("not json")
}
