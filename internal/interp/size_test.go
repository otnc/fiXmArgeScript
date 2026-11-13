package interp

import (
	"errors"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// 上限ちょうどの配列はつなげず、1つ足りなければつなげる。ゼロの配列は中身に触るまでメモリを使わないので、本物の大きさで確かめる。
func TestConcatArrayLimit(t *testing.T) {
	in := newInterp("main.fxma", Options{})
	in.restart()
	full := NewArray(make([]Value, maxArrayLen))
	one := NewArray([]Value{Int(1)})

	_, _, err := in.concat(full, one, token.Pos{})
	if de, ok := errors.AsType[*diag.Error](err); !ok || de.Kind != diag.Bavellabion {
		t.Fatalf("上限を超える配列: %v, want Bavellabion", err)
	}
	if _, _, err := in.concat(NewArray(full.Elems[1:]), one, token.Pos{}); err != nil {
		t.Fatalf("上限ちょうどの配列: %v", err)
	}
}
