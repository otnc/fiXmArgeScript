package karady

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestParseTarget(t *testing.T) {
	if got, err := ParseTarget("darwin-arm64"); err != nil || got != (Target{"darwin", "arm64"}) {
		t.Errorf("ParseTarget = %v, %v", got, err)
	}
	for _, s := range []string{"plan9-amd64", "linux", "linux-386"} {
		if _, err := ParseTarget(s); err == nil {
			t.Errorf("ParseTarget(%q) should fail", s)
		}
	}
}

func TestRuntimesGet(t *testing.T) {
	linux := tarGz(t, "fixmarge_0.1.0_linux_amd64/fixmarge", "linux binary")
	windows := zipped(t, "fixmarge.exe", "windows binary")
	release := map[string][]byte{
		"fixmarge_0.1.0_linux_amd64.tar.gz":  linux,
		"fixmarge_0.1.0_windows_arm64.zip":   windows,
		"fixmarge_0.1.0_darwin_arm64.tar.gz": []byte("tampered"),
	}
	sums := sum(linux) + "  fixmarge_0.1.0_linux_amd64.tar.gz\n" + sum(windows) + "  fixmarge_0.1.0_windows_arm64.zip\n" + sum([]byte("original")) + "  fixmarge_0.1.0_darwin_arm64.tar.gz\n"
	release["checksums.txt"] = []byte(sums)
	downloads := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downloads++
		data, ok := release[filepath.Base(r.URL.Path)]
		if !ok || filepath.Base(filepath.Dir(r.URL.Path)) != "v0.1.0" {
			http.NotFound(w, r)
			return
		}
		w.Write(data)
	}))
	defer srv.Close()
	r := &Runtimes{Root: t.TempDir(), BaseURL: srv.URL}

	for target, want := range map[Target]string{{"linux", "amd64"}: "linux binary", {"windows", "arm64"}: "windows binary"} {
		p, err := r.Get("v0.1.0", target)
		if err != nil {
			t.Fatalf("%v: %v", target, err)
		}
		if data, _ := os.ReadFile(p); string(data) != want {
			t.Errorf("%v: got %q", target, data)
		}
	}
	before := downloads
	if _, err := r.Get("v0.1.0", Target{"linux", "amd64"}); err != nil || downloads != before {
		t.Errorf("cached Get downloaded again (%v)", err)
	}
	if _, err := r.Get("v0.1.0", Target{"darwin", "arm64"}); !isUnplaced(err, diag.CraeyTA) {
		t.Errorf("tampered: err = %v, want CraeyTA'", err)
	}
	if _, err := r.Get("v9.9.9", Target{"linux", "amd64"}); !isUnplaced(err, diag.NiiExi) {
		t.Errorf("missing version: err = %v, want Nii eXi'", err)
	}
}

func isUnplaced(err error, kind diag.Kind) bool {
	u, ok := errors.AsType[*diag.Unplaced](err)
	return ok && u.Kind == kind
}

func sum(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func tarGz(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(content))}); err != nil {
		t.Fatal(err)
	}
	tw.Write([]byte(content))
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

func zipped(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte(content))
	zw.Close()
	return buf.Bytes()
}
