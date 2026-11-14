package release

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

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

func TestGet(t *testing.T) {
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
	r := &Fetcher{Root: t.TempDir(), BaseURL: srv.URL}

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
