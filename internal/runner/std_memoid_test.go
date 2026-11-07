package runner

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// fileIs は、実行した後に name の中身が want になっていることを確かめる。
func fileIs(name, want string) func(t *testing.T, fsys afero.Fs) {
	return func(t *testing.T, fsys afero.Fs) {
		t.Helper()
		got, err := afero.ReadFile(fsys, name)
		if err != nil || string(got) != want {
			t.Errorf("%s: %q, %v, want %q", name, got, err, want)
		}
	}
}

func notExists(name string) func(t *testing.T, fsys afero.Fs) {
	return func(t *testing.T, fsys afero.Fs) {
		t.Helper()
		if ok, _ := afero.Exists(fsys, name); ok {
			t.Errorf("%s が残っている", name)
		}
	}
}

func TestStdMemoid(t *testing.T) {
	head := "\"fxma:memoid\" aste memoid werer\n"
	cases := []testCase{
		{name: "read", src: head + "memoid no fineclowt(\"a.txt\") redogh\nmemoid no fineclowt(\"/work/a.txt\", \"sanand\") redogh\n", disk: map[string]string{"a.txt": "魔女"}, stdout: "魔女\nsAnAnd(\"e9ad94e5a5b3\", \"hex\")\n"},
		{name: "read json", src: head + "d izu memoid no fineclowt(\"d.json\", \"json\")\nd no \"k\" redogh\n", disk: map[string]string{"d.json": "{\"k\": [1, 2]}"}, stdout: "[1, 2]\n"},
		{name: "missing", src: head + "memoid no fineclowt(\"nothing.txt\")\n", err: diag.NiiExi},
		{name: "write without dir", src: head + "memoid no surc(\"out/x.txt\", \"a\")\n", err: diag.NiiExi},
		{name: "write", src: head + "memoid no surc(\"x.txt\", \"a\")\nmemoid no naddy(\"x.txt\", sanand(\"b\"))\nmemoid no fineclowt(\"x.txt\") redogh\n", stdout: "ab\n", check: fileIs("/work/x.txt", "ab")},
		{name: "write json", src: head + "memoid no surc(\"d.json\", {\"a\": [1]}, \"json\", 2)\n", check: fileIs("/work/d.json", "{\n  \"a\": [\n    1\n  ]\n}\n")},
		{name: "dirs", src: head + "memoid no meifKuwpace(\"a/b\")\nmemoid no meifKuwpace(\"a/b\")\nmemoid no surc(\"a/b/c.txt\", \"c\")\nmemoid no surc(\"a/d.txt\", \"d\")\nmemoid no fineclowtKuwpace(\"a\") redogh\nmemoid no salook(\"a/b/c.txt\") no \"nang\" redogh\nmemoid no salook(\"a\") no \"izuKuwpace\" redogh\n", stdout: "[\"b\", \"d.txt\"]\n1\nTadACt\n"},
		{name: "copy dir", src: head + "memoid no sagain(\"src\", \"dst/copy\")\n", disk: map[string]string{"src/a.txt": "a", "src/sub/b.txt": "b"}, check: func(t *testing.T, fsys afero.Fs) {
			fileIs("/work/dst/copy/a.txt", "a")(t, fsys)
			fileIs("/work/dst/copy/sub/b.txt", "b")(t, fsys)
			fileIs("/work/src/a.txt", "a")(t, fsys)
		}},
		{name: "move", src: head + "memoid no gokuu(\"a.txt\", \"moved/b.txt\")\nmemoid no izanavitedExi(\"a.txt\") redogh\n", disk: map[string]string{"a.txt": "a"}, stdout: "ULieSo\n", check: fileIs("/work/moved/b.txt", "a")},
		{name: "remove dir", src: head + "memoid no vaniell(\"d\")\n", disk: map[string]string{"d/a.txt": "a"}, check: notExists("/work/d")},
		{name: "remove missing", src: head + "memoid no vaniell(\"nothing\")\n", err: diag.NiiExi},
		{name: "edo", src: head + "memoid no edo() redogh\n", stdout: "/work\n"},
		{name: "needs abiryoku", src: head + "weho f()\n  memoid no fineclowt(\"a.txt\")\nfinend\nf()\n", disk: map[string]string{"a.txt": "a"}, err: diag.KiNbOO},
		// 書いた後で巻き戻ると、書く前の選択肢には戻れない (取り消せない操作は選択肢を確定する)
		{name: "write commits choices", src: head + "x izu 1 owk 2 mekamo\nmemoid no surc(\"x.txt\", \"#{x}\")\nx senaji 2 shoubek\n", err: diag.CoJundic, check: fileIs("/work/x.txt", "1")},
		// 書いた後の選択肢で巻き戻っても、やり直しで同じ追記を二度しない
		{name: "no double append", src: head + "memoid no naddy(\"log.txt\", \"a\")\ny izu 1 owk 2 mekamo\ny senaji 2 shoubek\nmemoid no fineclowt(\"log.txt\") redogh\n", stdout: "a\n", check: fileIs("/work/log.txt", "a")},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}
