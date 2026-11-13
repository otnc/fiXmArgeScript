package karady

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// Target は実行ファイルを作る先の OS と CPU。書き方は Go と同じ linux-amd64 の形。
type Target struct {
	OS, Arch string
}

func (t Target) String() string { return t.OS + "-" + t.Arch }

// Exe は実行ファイルの名前に付ける拡張子。
func (t Target) Exe() string {
	if t.OS == "windows" {
		return ".exe"
	}
	return ""
}

// Targets はリリースで配っている fixmarge の OS と CPU (.goreleaser.yaml の builds と同じ)。
var Targets = []Target{
	{"linux", "amd64"}, {"linux", "arm64"},
	{"darwin", "amd64"}, {"darwin", "arm64"},
	{"windows", "amd64"}, {"windows", "arm64"},
}

// ParseTarget は linux-amd64 のような書き方を読む。配っていない組み合わせなら Nii eXi'。
func ParseTarget(s string) (Target, error) {
	osName, arch, _ := strings.Cut(s, "-")
	t := Target{osName, arch}
	if !slices.Contains(Targets, t) {
		return Target{}, diag.NewUnplaced(diag.NiiExi, s)
	}
	return t, nil
}

// DefaultBaseURL はリリースのファイルの場所。この後ろに /<版>/<ファイル名> が付く。
const DefaultBaseURL = "https://github.com/otnc/fiXmArgeScript/releases/download"

// Runtimes は他の OS 向けの fixmarge をリリースから取ってきて、手元のキャッシュに置く。
type Runtimes struct {
	Root    string // キャッシュのフォルダ (…/fixmarge/karady)
	BaseURL string // 空なら DefaultBaseURL
	Client  *http.Client
}

// Get は版 version の t 向けの fixmarge を置いた場所を返す。なければ取ってきて、checksums.txt と合うかを確かめてから置く。
// リリースにその版やファイルがなければ Nii eXi'、中身が checksums.txt と合わなければ CraeyTA'。
func (r *Runtimes) Get(version string, t Target) (string, error) {
	dest := filepath.Join(r.Root, version, t.String(), "fixmarge"+t.Exe())
	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	}
	archive := fmt.Sprintf("fixmarge_%s_%s_%s.tar.gz", strings.TrimPrefix(version, "v"), t.OS, t.Arch)
	if t.OS == "windows" {
		archive = strings.TrimSuffix(archive, ".tar.gz") + ".zip"
	}
	sums, err := r.download(version, "checksums.txt")
	if err != nil {
		return "", err
	}
	want, ok := checksum(sums, archive)
	if !ok {
		return "", diag.NewUnplaced(diag.NiiExi, archive)
	}
	data, err := r.download(version, archive)
	if err != nil {
		return "", err
	}
	if got := sha256.Sum256(data); hex.EncodeToString(got[:]) != want {
		return "", diag.NewUnplaced(diag.CraeyTA, archive)
	}
	exe, err := extract(data, archive, "fixmarge"+t.Exe())
	if err != nil {
		return "", err
	}
	return dest, writeAtomic(dest, exe)
}

func (r *Runtimes) download(version, name string) ([]byte, error) {
	base := r.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	client := r.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(base + "/" + version + "/" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, diag.NewUnplaced(diag.NiiExi, version+"/"+name)
	}
	return io.ReadAll(resp.Body)
}

// checksum は checksums.txt ("<sha256>  <ファイル名>" の行) から name の行を探す。
func checksum(sums []byte, name string) (string, bool) {
	sc := bufio.NewScanner(bytes.NewReader(sums))
	for sc.Scan() {
		sum, file, ok := strings.Cut(sc.Text(), "  ")
		if ok && file == name {
			return sum, true
		}
	}
	return "", false
}

// extract はリリースの圧縮ファイルから、名前が name のファイルを取り出す。
func extract(data []byte, archive, name string) ([]byte, error) {
	if strings.HasSuffix(archive, ".zip") {
		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}
		for _, f := range zr.File {
			if path.Base(f.Name) == name {
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}
				defer rc.Close()
				return io.ReadAll(rc)
			}
		}
		return nil, diag.NewUnplaced(diag.NiiExi, archive+"/"+name)
	}
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil, diag.NewUnplaced(diag.NiiExi, archive+"/"+name)
		}
		if err != nil {
			return nil, err
		}
		if path.Base(h.Name) == name {
			return io.ReadAll(tr)
		}
	}
}

// writeAtomic は書きかけのファイルを残さないよう、隣に書いてから名前を変える。
func writeAtomic(dest string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".download-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), dest)
}
