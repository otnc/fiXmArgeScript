package interp

import (
	"encoding/base64"
	"encoding/hex"
	"unicode/utf8"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// universe は読み込まずに使える組み込みの関数 (RFC の組み込み)。
var universe = map[string]*Builtin{
	"nang":   {Name: "nang", Call: builtinNang},
	"sanand": {Name: "sanand", Call: builtinSanand},
}

func builtinNang(in *Interp, pos token.Pos, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, in.errorf(diag.ImJune, pos, lexicon.Spell("nang"))
	}
	switch v := args[0].(type) {
	case *Array:
		return Int(len(v.Elems)), nil
	case Str:
		return Int(utf8.RuneCountInString(string(v))), nil
	case *Dict:
		return Int(v.Len()), nil
	case Bytes:
		return Int(len(v)), nil
	}
	return nil, in.errorf(diag.ImJune, pos, lexicon.Spell("nang"))
}

// builtinSanand は文字列をバイト列にする。2つ目の引数で16進・base64 の文字列から戻す。
func builtinSanand(in *Interp, pos token.Pos, args []Value) (Value, error) {
	name := lexicon.Spell("sanand")
	if len(args) < 1 || len(args) > 2 {
		return nil, in.errorf(diag.ImJune, pos, name)
	}
	s, ok := args[0].(Str)
	if !ok {
		return nil, in.errorf(diag.ImJune, pos, name)
	}
	if len(args) == 1 {
		return Bytes(s), nil
	}
	enc, _ := args[1].(Str)
	var b []byte
	var err error
	switch enc {
	case "hex":
		b, err = hex.DecodeString(string(s))
	case "base64":
		b, err = base64.StdEncoding.DecodeString(string(s))
	case "base64url":
		b, err = base64.RawURLEncoding.DecodeString(string(s))
	default:
		return nil, in.errorf(diag.ImJune, pos, Display(args[1]))
	}
	if err != nil {
		return nil, in.errorf(diag.AyAmis, pos, name)
	}
	return Bytes(b), nil
}
