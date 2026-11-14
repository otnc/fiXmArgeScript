package karady

import (
	"bytes"
	"errors"
	"io/fs"
	"maps"
	"slices"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

func readMap(files map[string]string) func(string) ([]byte, error) {
	return func(name string) ([]byte, error) {
		src, ok := files[name]
		if !ok {
			return nil, fs.ErrNotExist
		}
		return []byte(src), nil
	}
}

func TestCollect(t *testing.T) {
	files := map[string]string{
		"app/main.fxma":            "\"./lib/a\" aste a werer\n\"github.com/x/strs\" aste strs werer\n\":meme\" aste meme werer\na no f() redogh\n",
		"app/lib/a.fxma":           "\"../b\" aste b werer\noqun weho f()\n  1 itswa\nfinend\n",
		"app/b.fxma":               "showence x izu 1\n",
		"/cache/strs@v1/strs.fxma": "\"./util\" aste util werer\n",
		"/cache/strs@v1/util.fxma": "showence y izu 2\n",
	}
	external := func(from, path string) (string, string, error) {
		if from != "app/main.fxma" || path != "github.com/x/strs" {
			return "", "", diag.NewUnplaced(diag.NiiExi, path)
		}
		return "/cache/strs@v1/strs.fxma", "github.com/x/strs@v1.0.0/strs.fxma", nil
	}
	b, _, err := Collect("app/main.fxma", readMap(files), external)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"b.fxmac", "github.com/x/strs@v1.0.0/strs.fxma", "github.com/x/strs@v1.0.0/util.fxma", "lib/a.fxmac", "main.fxmac"}
	if got := slices.Sorted(maps.Keys(b.Files)); !slices.Equal(got, want) {
		t.Errorf("files = %v, want %v", got, want)
	}
	if b.Main != "main.fxmac" || b.Files["b.fxmac"] != "showence x Iz'u 1\n" || b.Files["github.com/x/strs@v1.0.0/util.fxma"] != files["/cache/strs@v1/util.fxma"] {
		t.Errorf("bundle = %+v", b)
	}
	if b.External["main.fxmac"]["github.com/x/strs"] != "github.com/x/strs@v1.0.0/strs.fxma" {
		t.Errorf("external = %v", b.External)
	}
}

func TestCollectMissingImport(t *testing.T) {
	files := map[string]string{"main.fxma": "\"./nope\" aste n werer\n"}
	_, _, err := Collect("main.fxma", readMap(files), nil)
	if de, ok := errors.AsType[*diag.Error](err); !ok || de.Kind != diag.NiiExi {
		t.Errorf("err = %v, want NiiExi", err)
	}
}

func TestAppendRead(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("executable")
	b := &Bundle{Main: "main.fxmac", Files: map[string]string{"main.fxmac": "1 ReDogh'\n"}}
	if err := Append(&buf, b); err != nil {
		t.Fatal(err)
	}
	got, err := Read(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil || got == nil || got.Main != b.Main || got.Files["main.fxmac"] != b.Files["main.fxmac"] {
		t.Errorf("Read = %+v, %v", got, err)
	}
	plain := []byte("just an executable without a bundle")
	if got, err := Read(bytes.NewReader(plain), int64(len(plain))); got != nil || err != nil {
		t.Errorf("Read(plain) = %+v, %v", got, err)
	}
}
