package lexicon

import (
	"encoding/csv"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestSkeleton(t *testing.T) {
	cases := map[string]string{
		"Iz'u":       "izu",
		"S~ond'":     "sond",
		"Syu Tain":   "syutain",
		"Kei-P":      "keip",
		"Vani·tA":    "vanita",
		"kee’on":     "keeon",
		"izuDefutai": "izudefutai",
	}
	for in, want := range cases {
		if got := Skeleton(in); got != want {
			t.Errorf("Skeleton(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCompoundSpelling(t *testing.T) {
	cases := map[string]string{
		"izuDefutai":    "Iz'udeFuTai",
		"madrazySanand": "madRazysAnAnd",
		"izanavitedExi": "i·zana-vitedeXi'",
		"laineiAlte":    "Lai NeiAlte",
	}
	for name, want := range cases {
		e, _, ok := Lookup(Skeleton(name))
		if !ok || e.Name != name || e.Spelling != want {
			t.Errorf("Lookup(%q) = %+v, %v; want spelling %q", name, e, ok, want)
		}
	}
}

func TestDialect(t *testing.T) {
	e, dialect, ok := Lookup("lowhel")
	if !ok || !dialect || e.Spelling != "LowhEll" {
		t.Errorf("Lookup(lowhel) = %+v, %v, %v", e, dialect, ok)
	}
	for from, to := range dialects {
		if _, ok := bySkeleton[from]; ok {
			t.Errorf("方言 %q が標準語の骨格とぶつかっている", from)
		}
		if _, ok := bySkeleton[to]; !ok {
			t.Errorf("方言 %q の行き先 %q が表にない", from, to)
		}
	}
}

func TestMatchSpellingPrefersLongest(t *testing.T) {
	e, ok := MatchSpelling("Iz'udeFuTai(p)")
	if !ok || e.Name != "izuDefutai" {
		t.Errorf("MatchSpelling = %+v, %v", e, ok)
	}
}

// TestDictionary は、表の綴りが辞書 CSV にそのままあることと、辞書の骨格がぶつからないことを確かめる。
// CSV はライセンスがなく同梱できないので、FXMA_DICT_CHECK=1 のときだけ取得する (CI で使う)。
func TestDictionary(t *testing.T) {
	if os.Getenv("FXMA_DICT_CHECK") != "1" {
		t.Skip("FXMA_DICT_CHECK=1 のときだけ辞書 CSV を取得して確かめる")
	}
	const url = "https://raw.githubusercontent.com/Asankilp/fiXmArge-dictionary/main/fixmarge-EN.csv"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	dict := map[string]string{}
	for _, row := range rows[1:] {
		spelling := strings.TrimSpace(row[0])
		if spelling == "" {
			continue
		}
		skel := Skeleton(spelling)
		if prev, dup := dict[skel]; dup && prev != spelling {
			t.Errorf("辞書の骨格がぶつかっている: %q と %q", prev, spelling)
		}
		dict[skel] = spelling
	}

	for _, w := range words {
		if got, ok := dict[Skeleton(w.Spelling)]; !ok || got != w.Spelling {
			t.Errorf("%q が辞書にない (辞書の綴り: %q)", w.Spelling, got)
		}
	}
	for from := range dialects {
		if _, ok := dict[from]; ok {
			t.Errorf("方言 %q が辞書の別の語とぶつかっている", from)
		}
	}
}
