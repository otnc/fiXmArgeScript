// Package lexicon は fxma が使う辞書の語の表と、骨格と正式な綴りの変換を持つ。
// キーワード、エラー名、標準パッケージの名前は、すべてこの表から引く。
package lexicon

import (
	"sort"
	"strings"
	"unicode"
)

// Use は語の使い道。1つの語が複数の使い道を持てる。
type Use uint32

const (
	Keyword Use = 1 << iota
	Builtin
	Package
	Function
	Key
	Part
	Error
	LintRule
	Console
	Command
	Option
	ModField
	Installer
	Compound
)

// Word は辞書の1語。
type Word struct {
	Spelling string // 正式な綴り (清書の形)
	Meaning  string // 辞書の訳
	Uses     Use
}

// Entry は下書きに書ける名前1つ。辞書の1語か、2語の名前。
type Entry struct {
	Name     string // 下書きに書く形。1語なら骨格、2語の名前なら camelCase
	Spelling string // 清書に書く形
	Uses     Use
}

func (e Entry) IsKeyword() bool { return e.Uses&Keyword != 0 }

var (
	bySkeleton = map[string]Entry{}
	// longestFirst は清書を読むときの最長一致のための並び。
	longestFirst []Entry
)

func init() {
	for _, w := range words {
		add(Entry{Name: Skeleton(w.Spelling), Spelling: w.Spelling, Uses: w.Uses})
	}
	for _, name := range compounds {
		var spelling strings.Builder
		for _, part := range SplitCamel(name) {
			spelling.WriteString(mustLookup(part).Spelling)
		}
		add(Entry{Name: name, Spelling: spelling.String(), Uses: Compound})
	}
	for _, e := range bySkeleton {
		longestFirst = append(longestFirst, e)
	}
	sort.Slice(longestFirst, func(i, j int) bool {
		a, b := longestFirst[i].Spelling, longestFirst[j].Spelling
		if len(a) != len(b) {
			return len(a) > len(b)
		}
		return a < b
	})
}

func add(e Entry) {
	skel := Skeleton(e.Name)
	if _, dup := bySkeleton[skel]; dup {
		panic("lexicon: 骨格がぶつかっている: " + skel)
	}
	bySkeleton[skel] = e
}

// Skeleton は綴りから記号と空白を抜いて小文字にした形を返す。
func Skeleton(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r < unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

// SplitCamel は camelCase の名前を小文字の部品に分ける。
func SplitCamel(name string) []string {
	var parts []string
	start := 0
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			parts = append(parts, strings.ToLower(name[start:i]))
			start = i
		}
	}
	return append(parts, strings.ToLower(name[start:]))
}

// Lookup は骨格から語を探す。方言の骨格は標準語に読み替え、dialect に true を返す。
func Lookup(skel string) (e Entry, dialect bool, ok bool) {
	if std, isDialect := dialects[skel]; isDialect {
		e, ok = bySkeleton[std]
		return e, true, ok
	}
	e, ok = bySkeleton[skel]
	return e, false, ok
}

// Spell は骨格から正式な綴りを返す。表にない骨格はプログラムの誤りなので panic する。
func Spell(skel string) string {
	return mustLookup(skel).Spelling
}

// Phrase は複数の語を空白でつないだ表示 (Nii eXi' など) を返す。
func Phrase(skels ...string) string {
	spelled := make([]string, len(skels))
	for i, s := range skels {
		spelled[i] = Spell(s)
	}
	return strings.Join(spelled, " ")
}

func mustLookup(skel string) Entry {
	e, ok := bySkeleton[skel]
	if !ok {
		panic("lexicon: 表にない語: " + skel)
	}
	return e
}

// MatchSpelling は s の先頭に一致する正式な綴りのうち、いちばん長いものを返す。
func MatchSpelling(s string) (Entry, bool) {
	for _, e := range longestFirst {
		if strings.HasPrefix(s, e.Spelling) {
			return e, true
		}
	}
	return Entry{}, false
}

// Words は表の語をすべて返す。
func Words() []Word {
	return append([]Word(nil), words...)
}

// Entries は2語の名前も含めた、下書きに書ける名前をすべて返す。
func Entries() []Entry {
	return append([]Entry(nil), longestFirst...)
}
