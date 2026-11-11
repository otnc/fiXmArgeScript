package runner

import (
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

const utilSrc = "oqun weho double(n)\n  n * 2 itswa\nfinend\nweho secret()\n  1 itswa\nfinend\noqun showence limit izu 15\n"

func TestModules(t *testing.T) {
	util := map[string]string{"util.fxma": utilSrc}
	cases := []testCase{
		{name: "alias", src: "\"./util\" aste util werer\nutil no double(21) redogh\nutil no limit redogh\n", files: util, stdout: "42\n15\n"},
		{name: "names", src: "\"./util\" no double, limit werer\ndouble(limit) redogh\n", files: util, stdout: "30\n"},
		{name: "rename", src: "\"./util\" no double aste d werer\nd(2) redogh\n", files: util, stdout: "4\n"},
		{name: "not exported", src: "\"./util\" aste util werer\nutil no secret() redogh\n", files: util, err: diag.CaNay},
		{name: "not exported by name", src: "\"./util\" no secret werer\n", files: util, err: diag.CaNay},
		{name: "missing member", src: "\"./util\" aste util werer\nutil no triple(1) redogh\n", files: util, err: diag.NiiExi},
		{name: "missing file", src: "\"./nothing\" aste n werer\n", err: diag.NiiExi},
		{name: "imported name is const", src: "\"./util\" no double werer\ndouble izu 1\n", files: util, err: diag.BlasTOac},
		{name: "member is const", src: "\"./util\" aste util werer\nutil no limit izu 1\n", files: util, err: diag.BlasTOac},
		{name: "same name twice", src: "\":coword\" no eda werer\n\":varizama\" no eda werer\n", err: diag.SuReady},
		{name: "same name with aste", src: "\":coword\" no eda werer\n\":varizama\" no eda aste xsEda werer\neda(\"abc\", \"c\") redogh\nxsEda([1, 2], 2) redogh\n", stdout: "2\n1\n"},
		{name: "cycle", src: "\"./a\" aste a werer\n", files: map[string]string{"a.fxma": "\"./b\" aste b werer\n", "b.fxma": "\"./a\" aste a werer\n"}, err: diag.SuReady},
		{name: "statement in module", src: "\"./bad\" aste bad werer\n", files: map[string]string{"bad.fxma": "1 redogh\n"}, err: diag.CaNay},
		{name: "module imports module", src: "\"./lib/twice\" aste t werer\nt no quad(3) redogh\n", files: map[string]string{
			"lib/twice.fxma": "\"../util\" no double werer\noqun weho quad(n)\n  double(double(n)) itswa\nfinend\n",
			"util.fxma":      utilSrc,
		}, stdout: "12\n"},
		{name: "abiryoku from module", src: "\"./io\" aste io werer\nweho f()\n  io no hello()\nfinend\nf()\n", files: map[string]string{"io.fxma": "oqun abiryoku weho hello()\n  \"hi\" redogh\nfinend\n"}, err: diag.KiNbOO},
		{name: "backtrack through module", src: "\"./util\" no double werer\nx izu 1 owk 2 owk 3 mekamo\ndouble(x) senaji 6 shoubek\nx redogh\n", files: util, stdout: "3\n"},
		{name: "external package not yet", src: "\"github.com/alice/strs\" aste strs werer\n", err: diag.CaNay},
		{name: "unknown std", src: "\":nothing\" aste n werer\n", err: diag.NiiExi},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestStd(t *testing.T) {
	cases := []testCase{
		{name: "meme rounding", src: "\":meme\" aste meme werer\nmeme no racdam(2.7) redogh\nmeme no hir(2.1) redogh\nmeme no litto(2.5) redogh\nmeme no litto(-2.5) redogh\nmeme no pyusense(-3) redogh\n", stdout: "2\n3\n3\n-3\n3\n"},
		{name: "meme extremum", src: "\":meme\" aste meme werer\nmeme no tsuyong(1, 5, 3) redogh\nmeme no yoweak([4, 2.5, 9]) redogh\nmeme no zelln([1, 2, 3]) redogh\nmeme no zelln([]) redogh\n", stdout: "5\n2.5\n6\n0\n"},
		{name: "meme empty max", src: "\":meme\" aste meme werer\nmeme no tsuyong([]) redogh\n", err: diag.NiiExi},
		{name: "meme constants", src: "\":meme\" aste meme werer\nmeme no pi redogh\nmeme no mource(16) redogh\nmeme no log10(1000) redogh\n", stdout: "3.141592653589793\n4\n3\n"},
		{name: "meme log domain", src: "\":meme\" aste meme werer\nmeme no log(0) redogh\n", err: diag.ImJune},
		{name: "meme pi is const", src: "\":meme\" aste meme werer\nmeme no pi izu 3\n", err: diag.BlasTOac},
		{name: "coword salook", src: "\":coword\" aste coword werer\ncoword no salook(\"3\") redogh\ncoword no salook(\"3\") / 2 redogh\ncoword no salook(1.5) + \"!\" redogh\ncoword no salook(sanand(\"魔女\")) redogh\n", stdout: "3\n1.5\n1.5!\n魔女\n"},
		{name: "coword salook bad", src: "\":coword\" aste coword werer\ncoword no salook(\"abc\") redogh\n", err: diag.AyAmis},
		{name: "coword split join", src: "\":coword\" aste coword werer\nxs izu coword no confsoy(\"a,b,c\", \",\")\nxs redogh\ncoword no recanrey(xs, \"-\") redogh\ncoword no recanrey([1, 2.5], \"+\") redogh\n", stdout: "[\"a\", \"b\", \"c\"]\na-b-c\n1+2.5\n"},
		{name: "coword replace", src: "\":coword\" aste coword werer\ncoword no lainei(\"aaa\", \"a\", \"b\") redogh\ncoword no laineiAlte(\"aaa\", \"a\", \"b\") redogh\n", stdout: "baa\nbbb\n"},
		{name: "coword case and trim", src: "\":coword\" aste coword werer\ncoword no gdie(\"abc\") redogh\ncoword no lighluw(\"ABC\") redogh\n\"[#{coword no shojika(\"  x \")}]\" redogh\n", stdout: "ABC\nabc\n[x]\n"},
		{name: "coword position", src: "\":coword\" aste coword werer\ncoword no eda(\"魔女裁判\", \"裁\") redogh\ncoword no eda(\"abc\", \"z\") redogh\ncoword no kast(\"abc\", \"ab\") redogh\ncoword no gun(\"abc\", \"ab\") redogh\n", stdout: "2\n-1\nTadACt\nULieSo\n"},
		{name: "coword pad", src: "\":coword\" aste coword werer\ncoword no sodaKast(\"5\", 3, \"0\") redogh\n\"[#{coword no sodaGun(\"ab\", 4)}]\" redogh\ncoword no sodaKast(\"x\", 6, \"ab\") redogh\n", stdout: "005\n[ab  ]\nababax\n"},
		{name: "coword code point", src: "\":coword\" aste coword werer\ncoword no kolsu(\"魔\") redogh\ncoword no batteKolsu(39764) redogh\n", stdout: "39764\n魔\n"},
		{name: "varizama mutate", src: "\":varizama\" aste varizama werer\nxs izu [1, 2, 3]\nvarizama no naddy(xs, 4)\nvarizama no vaniell(xs) redogh\nvarizama no vaniell(xs, 0) redogh\nvarizama no lea(xs, 1, 9)\nxs redogh\n", stdout: "4\n1\n[2, 9, 3]\n"},
		{name: "varizama new arrays", src: "\":varizama\" aste varizama werer\nxs izu [3, 1, 2]\nvarizama no muyon(xs) redogh\nvarizama no orei(xs) redogh\nxs redogh\n", stdout: "[2, 1, 3]\n[1, 2, 3]\n[3, 1, 2]\n"},
		{name: "varizama orei with function", src: "\":varizama\" aste varizama werer\nxs izu [[2, \"a\"], [1, \"b\"], [2, \"c\"], [1, \"d\"]]\nvarizama no orei(xs, weho (a, b) (a no 0) - (b no 0) itswa finend) redogh\n", stdout: "[[1, \"b\"], [1, \"d\"], [2, \"a\"], [2, \"c\"]]\n"},
		{name: "varizama orei mixed", src: "\":varizama\" aste varizama werer\nvarizama no orei([1, \"a\"]) redogh\n", err: diag.ImJune},
		{name: "varizama higher order", src: "\":varizama\" aste varizama werer\nxs izu [1, 2, 3, 4]\nvarizama no tarip(xs, weho (x) x * 10 itswa finend) redogh\nvarizama no onkhe(xs, weho (x) x % 2 senaji 0 itswa finend) redogh\nvarizama no matill(xs, weho (acc, x) acc + x itswa finend, 0) redogh\nvarizama no eda(xs, 3) redogh\nvarizama no eda(xs, weho (x) x > 1 itswa finend) redogh\nvarizama no eda(xs, 9) redogh\n", stdout: "[10, 20, 30, 40]\n[2, 4]\n10\n2\n1\n-1\n"},
		{name: "varizama onkhe needs bool", src: "\":varizama\" aste varizama werer\nvarizama no onkhe([1], weho (x) x itswa finend) redogh\n", err: diag.ImJune},
		{name: "varizama passes abiryoku", src: "\":varizama\" aste varizama werer\nabiryoku weho show(x)\n  x redogh\n  x itswa\nfinend\nvarizama no tarip([1, 2], show)\n", stdout: "1\n2\n"},
		{name: "varizama abiryoku in plain function", src: "\":varizama\" aste varizama werer\nabiryoku weho show(x)\n  x redogh\nfinend\nweho f()\n  varizama no tarip([1], show)\nfinend\nf()\n", err: diag.KiNbOO},
		{name: "std function as value", src: "\":varizama\" aste varizama werer\n\":coword\" aste coword werer\nvarizama no tarip([1, 2], coword no salook) redogh\n", stdout: "[\"1\", \"2\"]\n"},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestStdJSON(t *testing.T) {
	head := "\":json\" aste json werer\n"
	cases := []testCase{
		{name: "parse", src: head + "d izu json no fineclowt(\"{\\\"b\\\": 1, \\\"a\\\": [true, null, \\\"x\\\"]}\")\nd redogh\nd no \"b\" + 1 redogh\n", stdout: "{\"b\": 1, \"a\": [TadACt, Vani·tA, \"x\"]}\n2\n"},
		{name: "numbers are floats", src: head + "json no fineclowt(\"1\") / 2 redogh\n", stdout: "0.5\n"},
		{name: "invalid", src: head + "json no fineclowt(\"{\\\"a\\\": }\")\n", err: diag.AyAmis},
		{name: "trailing", src: head + "json no fineclowt(\"1 2\")\n", err: diag.AyAmis},
		{name: "stringify", src: head + `json no salook({"a": [1, 2.5, "x\n"], "b": vanita, "c": {}}) redogh` + "\n", stdout: `{"a":[1,2.5,"x\n"],"b":null,"c":{}}` + "\n"},
		{name: "indent", src: head + "json no salook({\"a\": [1], \"b\": []}, 2) redogh\n", stdout: "{\n  \"a\": [\n    1\n  ],\n  \"b\": []\n}\n"},
		{name: "not json", src: head + "json no salook([sanand(\"a\")])\n", err: diag.ImJune},
		{name: "round trip", src: head + "s izu \"[1,{\\\"k\\\":\\\"魔女\\\"}]\"\njson no salook(json no fineclowt(s)) senaji s redogh\n", stdout: "TadACt\n"},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestStdIzanavited(t *testing.T) {
	head := "\":izanavited\" aste izanavited werer\n"
	cases := []testCase{
		{name: "spec examples", src: head + "izanavited no recanrey(\"a\", \"b\", \"c.txt\") redogh\nizanavited no ado(\"a/b/c.txt\") redogh\nizanavited no chidd(\"a/b/c.txt\") redogh\nizanavited no eatahguy(\"a/b\", \"a/c\") redogh\n", stdout: "a/b/c.txt\nc.txt\n.txt\n../c\n"},
		{name: "backslash", src: head + `izanavited no saymal("a\\b\\..\\c") redogh` + "\n" + `izanavited no kuwpace("a\\b\\c.txt") redogh` + "\n", stdout: "a/c\na/b\n"},
		{name: "node edge cases", src: head + "izanavited no chidd(\".bashrc\") redogh\nizanavited no saymal(\"a/./b/\") redogh\nizanavited no kuwpace(\"a/b/\") redogh\nizanavited no recanrey() redogh\n", stdout: "\na/b/\na\n.\n"},
		{name: "is absolute", src: head + "izanavited no izuDefutai(\"/x\") redogh\n" + `izanavited no izuDefutai("C:\\x") redogh` + "\nizanavited no izuDefutai(\"x\") redogh\n", stdout: "TadACt\nTadACt\nULieSo\n"},
		{name: "defutai", src: head + "izanavited no defutai(\"a/../b\") redogh\n", stdout: "/work/b\n"},
		{name: "defutai needs abiryoku", src: head + "weho f()\n  izanavited no defutai(\"a\")\nfinend\nf()\n", err: diag.KiNbOO},
		{name: "relative escapes", src: head + "izanavited no eatahguy(\"../a\", \"b\")\n", err: diag.ImJune},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}
