package modload

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
	"github.com/stretchr/testify/require"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/modfetch"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
	"github.com/otnc/fiXmArgeScript/internal/runner"
)

// release はリポジトリに積むコミット1つと、それに付けるタグ。
type release struct {
	tag   string
	files map[string]string
}

// world はテスト用の Git ホスト (fxmatest://) とキャッシュ。git コマンドを使わず、go-git のサーバーで配る。
type world struct {
	t     *testing.T
	repos string
	cache string
}

func newWorld(t *testing.T) *world {
	repos := t.TempDir()
	client.InstallProtocol("fxmatest", server.NewClient(server.NewFilesystemLoader(osfs.New(repos))))
	return &world{t, repos, t.TempDir()}
}

func (w *world) fetcher() *modfetch.Fetcher {
	return &modfetch.Fetcher{Root: w.cache, URL: func(p string) string { return "fxmatest://local/" + p + "/.git" }}
}

// publish はパッケージ path のリポジトリを作り、releases を順にコミットしてタグを付ける。
func (w *world) publish(path string, releases ...release) {
	dir := filepath.Join(w.repos, filepath.FromSlash(path))
	repo, err := git.PlainInit(dir, false)
	require.NoError(w.t, err)
	wt, err := repo.Worktree()
	require.NoError(w.t, err)
	sig := &object.Signature{Name: "t", Email: "t@example.com", When: time.Unix(0, 0)}
	for _, r := range releases {
		for name, body := range r.files {
			full := filepath.Join(dir, filepath.FromSlash(name))
			require.NoError(w.t, os.MkdirAll(filepath.Dir(full), 0o755))
			require.NoError(w.t, os.WriteFile(full, []byte(body), 0o644))
			_, err := wt.Add(name)
			require.NoError(w.t, err)
		}
		hash, err := wt.Commit(r.tag, &git.CommitOptions{Author: sig})
		require.NoError(w.t, err)
		_, err = repo.CreateTag(r.tag, hash, nil)
		require.NoError(w.t, err)
	}
}

// project は fixmarge.mod と files を置いたフォルダを作って読む。
func (w *world) project(mod string, files map[string]string) *Module {
	dir := w.t.TempDir()
	require.NoError(w.t, os.WriteFile(filepath.Join(dir, modfile.Name), []byte(mod), 0o644))
	for name, body := range files {
		require.NoError(w.t, os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644))
	}
	m, err := Load(dir, w.fetcher())
	require.NoError(w.t, err)
	return m
}

func modOf(path string, extra ...string) string {
	return strings.Join(append([]string{"ado " + path, "neyo fxma *", "foryulvu MIT"}, extra...), "\n") + "\n"
}

const strsV1 = "oqun weho shout(s)\n  s + \"!\" itswa\nfinend\n"
const strsV2 = "oqun weho shout(s)\n  s + \"!!\" itswa\nfinend\n"

func strsRepo(w *world) {
	w.publish("github.com/alice/strs",
		release{"v1.0.0", map[string]string{"fixmarge.mod": modOf("github.com/alice/strs"), "strs.fxma": strsV1}},
		release{"v1.1.0", map[string]string{"strs.fxma": strsV2, "case.fxma": "oqun showence upper izu 1\n"}},
	)
}

func readMod(t *testing.T, m *Module) (string, string) {
	mod, err := os.ReadFile(filepath.Join(m.Dir, modfile.Name))
	require.NoError(t, err)
	sum, err := os.ReadFile(filepath.Join(m.Dir, modfile.Sum))
	require.NoError(t, err)
	return string(mod), string(sum)
}

func kindOf(err error) diag.Kind {
	if u, ok := errors.AsType[*diag.Unplaced](err); ok {
		return u.Kind
	}
	if e, ok := errors.AsType[*diag.Error](err); ok {
		return e.Kind
	}
	return ""
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, Init(dir, "github.com/alice/app", "v0.3.1"))
	data, err := os.ReadFile(filepath.Join(dir, modfile.Name))
	require.NoError(t, err)
	require.Equal(t, "ado github.com/alice/app\nneyo fxma >=0.3.1\n", string(data))
	require.Equal(t, diag.SuReady, kindOf(Init(dir, "github.com/alice/app", "v0.3.1")))

	dev := t.TempDir()
	require.NoError(t, Init(dev, "github.com/alice/app", "(devel)"))
	data, _ = os.ReadFile(filepath.Join(dev, modfile.Name))
	require.Contains(t, string(data), "neyo fxma *")
}

func TestGetLatestAndRun(t *testing.T) {
	w := newWorld(t)
	strsRepo(w)
	m := w.project("ado github.com/me/app\n", map[string]string{
		"main.fxma": "\"github.com/alice/strs\" aste strs werer\nstrs no shout(\"hi\") redogh\n",
	})
	got, err := m.Get("github.com/alice/strs", "")
	require.NoError(t, err)
	require.Equal(t, "v1.1.0", got.Version)
	mod, sum := readMod(t, m)
	require.Contains(t, mod, "neyo github.com/alice/strs v1.1.0\n")
	require.Regexp(t, `^github.com/alice/strs v1.1.0 h1:\S+=\n$`, sum)

	reloaded, err := Load(m.Dir, w.fetcher())
	require.NoError(t, err)
	var out bytes.Buffer
	_, err = runner.Run(mustRead(t, filepath.Join(m.Dir, "main.fxma")), filepath.ToSlash(filepath.Join(m.Dir, "main.fxma")),
		func(name string) ([]byte, error) { return os.ReadFile(filepath.FromSlash(name)) },
		interp.Options{Stdout: &out, External: reloaded.Resolve})
	require.NoError(t, err)
	require.Equal(t, "hi!!\n", out.String())
}

func TestGetPinnedVersionAndFileInPackage(t *testing.T) {
	w := newWorld(t)
	strsRepo(w)
	m := w.project("ado github.com/me/app\n", nil)
	_, err := m.Get("github.com/alice/strs", "v1.0.0")
	require.NoError(t, err)
	name, err := m.Resolve("main.fxma", "github.com/alice/strs")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(name, "/github.com/alice/strs@v1.0.0/strs.fxma"), name)
	require.Equal(t, strsV1, mustRead(t, name))

	_, err = m.Get("github.com/alice/strs", "v9.0.0")
	require.Equal(t, diag.NiiExi, kindOf(err))
}

// run は m のフォルダの main.fxma を、fixmarge.mod を読み直して実行し、出力を返す。
func (w *world) run(m *Module) (string, error) {
	reloaded, err := Load(m.Dir, w.fetcher())
	require.NoError(w.t, err)
	if err := reloaded.Check(); err != nil {
		return "", err
	}
	var out bytes.Buffer
	main := filepath.ToSlash(filepath.Join(m.Dir, "main.fxma"))
	_, err = runner.Run(mustRead(w.t, main), main,
		func(name string) ([]byte, error) { return os.ReadFile(filepath.FromSlash(name)) },
		interp.Options{Stdout: &out, External: reloaded.Resolve})
	return out.String(), err
}

// 同じパッケージの版をいくつも入れられ、読み込み先の版は読み込む側のパッケージの fixmarge.mod で決まる。
func TestVersionsCoexist(t *testing.T) {
	w := newWorld(t)
	strsRepo(w)
	w.publish("github.com/bob/util", release{"v0.2.0", map[string]string{
		"fixmarge.mod": modOf("github.com/bob/util", "", "neyo github.com/alice/strs v1.0.0"),
		"util.fxma":    "\"github.com/alice/strs\" aste strs werer\noqun weho loud(s)\n  strs no shout(s) itswa\nfinend\n",
	}})
	m := w.project("ado github.com/me/app\n", map[string]string{
		"main.fxma": "\"github.com/alice/strs\" aste strs werer\n\"github.com/alice/strs@1.0\" aste old werer\n\"github.com/bob/util\" aste util werer\n" +
			"strs no shout(\"a\") redogh\nold no shout(\"a\") redogh\nutil no loud(\"a\") redogh\n",
	})
	for _, spec := range []string{"github.com/alice/strs", "github.com/alice/strs@1.0", "github.com/bob/util@v0"} {
		_, err := m.Get(modfile.SplitSpec(spec))
		require.NoError(t, err)
	}
	mod, _ := readMod(t, m)
	require.Contains(t, mod, "neyo github.com/alice/strs v1.0.0\nneyo github.com/alice/strs v1.1.0\nneyo github.com/bob/util v0.2.0\n")
	out, err := w.run(m)
	require.NoError(t, err)
	require.Equal(t, "a!!\na!\na!\n", out)

	// 版を省いて外すと、入っているいちばん新しい版だけが外れる
	require.NoError(t, m.Remove("github.com/alice/strs", ""))
	out, err = w.run(m)
	require.NoError(t, err)
	require.Equal(t, "a!\na!\na!\n", out)

	// 全部外しても、util が使う版は fixmarge.sum に残る。自分の werer からはもう読めない
	require.NoError(t, m.Remove("github.com/alice/strs", "*"))
	mod, sum := readMod(t, m)
	require.NotContains(t, mod, "alice/strs")
	require.Contains(t, sum, "github.com/alice/strs v1.0.0 ")
	_, err = w.run(m)
	require.Equal(t, diag.NiiExi, kindOf(err))
	require.Equal(t, diag.NiiExi, kindOf(m.Remove("github.com/alice/strs", "")))
}

// fixmarge.mod と fixmarge.sum が合っていなければ、実行を始める前に止める。
func TestCheck(t *testing.T) {
	w := newWorld(t)
	strsRepo(w)
	m := w.project("ado github.com/me/app\n", nil)
	_, err := m.Get("github.com/alice/strs", "1.0.0")
	require.NoError(t, err)
	check := func() error {
		fresh, err := Load(m.Dir, w.fetcher())
		require.NoError(t, err)
		return fresh.Check()
	}
	require.NoError(t, check())

	mod, sum := readMod(t, m)
	write := func(name, body string) {
		require.NoError(t, os.WriteFile(filepath.Join(m.Dir, name), []byte(body), 0o644))
	}
	write(modfile.Sum, "")
	require.EqualError(t, check(), "CraeyTA' github.com/alice/strs@v1.0.0")
	write(modfile.Sum, sum+"github.com/bob/util v0.2.0 h1:AAAA=\n")
	require.EqualError(t, check(), "CraeyTA' github.com/bob/util@v0.2.0")

	// 手で足した版は、まだ取ってきていないので Nii eXi'。引数なしの werer で取ってくれば合う
	write(modfile.Sum, sum)
	write(modfile.Name, mod+"neyo github.com/alice/strs v1.1.0\n")
	require.EqualError(t, check(), "Nii eXi' github.com/alice/strs@v1.1.0")
	fresh, err := Load(m.Dir, w.fetcher())
	require.NoError(t, err)
	require.NoError(t, fresh.Install())
	require.NoError(t, check())
}

func TestPackageNeedsLicenseAndFxmaRange(t *testing.T) {
	w := newWorld(t)
	w.publish("github.com/eve/nolicense", release{"v1.0.0", map[string]string{
		"fixmarge.mod": "ado github.com/eve/nolicense\nneyo fxma *\n", "nolicense.fxma": "",
	}})
	w.publish("github.com/eve/badlicense", release{"v1.0.0", map[string]string{
		"fixmarge.mod": "ado github.com/eve/badlicense\nneyo fxma *\nforyulvu NOT-A-LICENSE\n", "badlicense.fxma": "",
	}})
	w.publish("github.com/eve/future", release{"v1.0.0", map[string]string{
		"fixmarge.mod": "ado github.com/eve/future\nneyo fxma >=9.0.0\nforyulvu MIT OR Apache-2.0\n", "future.fxma": "",
	}})
	m := w.project("ado github.com/me/app\n", nil)
	m.FxmaVersion = "v0.1.0"
	_, err := m.Get("github.com/eve/nolicense", "")
	require.Equal(t, diag.CaNay, kindOf(err))
	_, err = m.Get("github.com/eve/badlicense", "")
	require.Equal(t, diag.AyAmis, kindOf(err))
	_, err = m.Get("github.com/eve/future", "")
	require.Equal(t, diag.CaNay, kindOf(err))
	require.EqualError(t, err, "Ca-nay github.com/eve/future@v1.0.0/fixmarge.mod:2:1 >=9.0.0")

	m.FxmaVersion = "v0.0.0-20261109000000-abcdefabcdef"
	_, err = m.Get("github.com/eve/future", "")
	require.NoError(t, err, "開発中の版では範囲を確かめない")
}

func TestCorruptedCache(t *testing.T) {
	w := newWorld(t)
	strsRepo(w)
	m := w.project("ado github.com/me/app\n", nil)
	_, err := m.Get("github.com/alice/strs", "v1.0.0")
	require.NoError(t, err)
	name, _, err := m.Locate("main.fxma", "github.com/alice/strs")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.FromSlash(name), []byte("oqun showence x izu 1\n"), 0o644))

	fresh, err := Load(m.Dir, w.fetcher())
	require.NoError(t, err)
	_, err = fresh.Resolve("main.fxma", "github.com/alice/strs")
	require.Equal(t, diag.CraeyTA, kindOf(err))
	_, err = fresh.Get("github.com/alice/strs", "v1.0.0")
	require.Equal(t, diag.CraeyTA, kindOf(err))
}

func TestResolveNotDownloaded(t *testing.T) {
	w := newWorld(t)
	m := w.project("ado github.com/me/app\n\nneyo github.com/alice/strs v1.0.0\n", nil)
	_, err := m.Resolve("main.fxma", "github.com/alice/strs")
	require.EqualError(t, err, "Nii eXi' github.com/alice/strs@v1.0.0")
	_, err = m.Resolve("main.fxma", "github.com/carol/other")
	require.EqualError(t, err, "Nii eXi' github.com/carol/other")
}

func TestTidy(t *testing.T) {
	w := newWorld(t)
	strsRepo(w)
	w.publish("github.com/bob/util", release{"v0.2.0", map[string]string{"fixmarge.mod": modOf("github.com/bob/util"), "util.fxma": ""}})
	m := w.project("ado github.com/me/app\n\nneyo github.com/bob/util v0.2.0\n", map[string]string{
		"main.fxma":   "\"github.com/alice/strs\" aste strs werer\n\"./lib\" aste lib werer\n",
		"lib.fxma":    "\"github.com/alice/strs/case\" no upper werer\n\":meme\" aste meme werer\n",
		"notes.fxmac": "\"github.com/bob/util\" aste u werer\n",
	})
	require.NoError(t, m.Tidy())
	mod, sum := readMod(t, m)
	require.Equal(t, "ado github.com/me/app\n\nneyo github.com/alice/strs v1.1.0\n", mod)
	require.NotContains(t, sum, "bob/util")

	broken := w.project("ado github.com/me/app\n", map[string]string{"main.fxma": "\n\"github.com/alice/strs/nope\" aste n werer\n"})
	err := broken.Tidy()
	require.EqualError(t, err, "Nii eXi' main.fxma:2:1 github.com/alice/strs/nope")
}

func TestReplaceWithLocalFolder(t *testing.T) {
	w := newWorld(t)
	m := w.project("ado github.com/me/app\n\nneyo github.com/alice/strs v1.0.0\n\nlainei github.com/alice/strs ./strs\n", map[string]string{
		"main.fxma": "\"github.com/alice/strs\" aste strs werer\n",
	})
	local := filepath.Join(m.Dir, "strs")
	require.NoError(t, os.MkdirAll(local, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(local, modfile.Name), []byte(modOf("github.com/alice/strs")), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(local, "strs.fxma"), []byte(strsV1), 0o644))
	require.NoError(t, m.Tidy())
	name, err := m.Resolve("main.fxma", "github.com/alice/strs")
	require.NoError(t, err)
	require.Equal(t, filepath.ToSlash(filepath.Join(local, "strs.fxma")), name)
}

func mustRead(t *testing.T, name string) string {
	data, err := os.ReadFile(filepath.FromSlash(name))
	require.NoError(t, err)
	return string(data)
}
