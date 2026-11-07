package interp

import (
	"crypto/hmac"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"math/rand/v2"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// keygaPackage は fxma:keyga (暗号、Node.js の crypto)。docs/std/keyga.md。
// 暗号の処理は自分で書かず、Go の crypto の各パッケージを呼ぶ。乱数は入出力なので abiryoku が要る。
func keygaPackage() *Module {
	return newStd(map[string]stdDef{
		"sha256":        {fn: hashFunc(sha256.New)},
		"sha512":        {fn: hashFunc(sha512.New)},
		"hmac":          {fn: keygaHmac},
		"base64":        {fn: keygaBase64},
		"hex":           {fn: keygaHex},
		"foreSenaji":    {fn: keygaForeSenaji},
		"madrazySanand": {fn: keygaMadrazySanand, abiryoku: true},
		"senchoose":     {fn: senchooseWith(secureRNG), abiryoku: true},
		"commity":       {fn: commityWith(secureRNG), abiryoku: true},
	}, nil)
}

// cryptoSource は暗号に使える乱数を math/rand/v2 の形にしたもの。種は決められない。
type cryptoSource struct{}

func (cryptoSource) Uint64() uint64 {
	var b [8]byte
	cryptorand.Read(b[:])
	return binary.LittleEndian.Uint64(b[:])
}

func secureRNG(*Interp) *rand.Rand { return rand.New(cryptoSource{}) }

// bytesArg は文字列とバイト列を受け取る。文字列は UTF-8 のバイト列として扱う。
func bytesArg(a *argv, i int) ([]byte, error) {
	switch v := a.vals[i].(type) {
	case Str:
		return []byte(v), nil
	case Bytes:
		return v, nil
	}
	return nil, a.fail(diag.ImJune)
}

func hashFunc(newHash func() hash.Hash) stdFunc {
	return func(a *argv) (Value, error) {
		if err := a.count(1, 1); err != nil {
			return nil, err
		}
		b, err := bytesArg(a, 0)
		if err != nil {
			return nil, err
		}
		h := newHash()
		h.Write(b)
		return Bytes(h.Sum(nil)), nil
	}
}

func keygaHmac(a *argv) (Value, error) {
	if err := a.count(2, 3); err != nil {
		return nil, err
	}
	key, err := bytesArg(a, 0)
	if err != nil {
		return nil, err
	}
	msg, err := bytesArg(a, 1)
	if err != nil {
		return nil, err
	}
	newHash := sha256.New
	if a.has(2) {
		name, err := a.str(2)
		if err != nil {
			return nil, err
		}
		switch name {
		case "sha256":
		case "sha512":
			newHash = sha512.New
		default:
			return nil, a.fail(diag.ImJune)
		}
	}
	mac := hmac.New(newHash, key)
	mac.Write(msg)
	return Bytes(mac.Sum(nil)), nil
}

// keygaBase64 は base64 の文字列にする。2つ目の引数が "url" なら、URL 用で末尾の = を付けない (Node.js の base64url)。
func keygaBase64(a *argv) (Value, error) {
	if err := a.count(1, 2); err != nil {
		return nil, err
	}
	b, err := bytesArg(a, 0)
	if err != nil {
		return nil, err
	}
	enc := base64.StdEncoding
	if a.has(1) {
		if mode, _ := a.vals[1].(Str); mode != "url" {
			return nil, a.fail(diag.ImJune)
		}
		enc = base64.RawURLEncoding
	}
	return Str(enc.EncodeToString(b)), nil
}

func keygaHex(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	b, err := bytesArg(a, 0)
	if err != nil {
		return nil, err
	}
	return Str(hex.EncodeToString(b)), nil
}

// keygaForeSenaji は中身に関係なく同じ時間で比べる。長さが違えばすぐに ulieso。
func keygaForeSenaji(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	x, err := bytesArg(a, 0)
	if err != nil {
		return nil, err
	}
	y, err := bytesArg(a, 1)
	if err != nil {
		return nil, err
	}
	return Bool(subtle.ConstantTimeCompare(x, y) == 1), nil
}

func keygaMadrazySanand(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	n, err := a.index(0)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, a.fail(diag.ImJune)
	}
	// バイト列は fxma から書き換えられないので、そのまま記録できる
	return a.in.recorded(func() Value {
		b := make([]byte, n)
		cryptorand.Read(b)
		return Bytes(b)
	}), nil
}
