package modfile

import (
	"regexp"
	"slices"
	"strings"

	"golang.org/x/mod/semver"
)

// AllVersions は、パッケージの版をすべて指す書き方 (fixmarge unneyo <パッケージ>@*)。
const AllVersions = "*"

// specPattern は版の指定の形。メジャーだけ、メジャーとマイナーだけ、全部、プレリリース付きのどれか。頭の v は書いても書かなくてもよい。
var specPattern = regexp.MustCompile(`^v?\d+(\.\d+(\.\d+(-[0-9A-Za-z.-]+)?)?)?$`)

// SplitSpec は "github.com/alice/strs/case@1.2" を、読み込み先と版の指定に分ける。版の指定がなければ spec は空。
func SplitSpec(s string) (path, spec string) {
	path, spec, _ = strings.CutLast(s, "@")
	if path == "" {
		return s, ""
	}
	return path, spec
}

// ValidSpec は版の指定として読めるかを返す。空 (省略) も読める。
func ValidSpec(spec string) bool {
	return spec == "" || specPattern.MatchString(spec)
}

// Pick は versions (v1.2.0 の形) から、spec に合ういちばん新しい版を選ぶ。合うものがなければ空。
// 正式版を先に選び、合う正式版がないときだけプレリリースを選ぶ。プレリリースまで書いた指定は、その版だけに合う。
func Pick(versions []string, spec string) string {
	matching := slices.DeleteFunc(slices.Clone(versions), func(v string) bool { return !matches(v, spec) })
	slices.SortFunc(matching, semver.Compare)
	for _, v := range slices.Backward(matching) {
		if semver.Prerelease(v) == "" {
			return v
		}
	}
	if len(matching) == 0 {
		return ""
	}
	return matching[len(matching)-1]
}

// matches は版 v (v1.2.0 の形) が spec に合うかを返す。
func matches(v, spec string) bool {
	if spec == "" {
		return true
	}
	want := "v" + strings.TrimPrefix(spec, "v")
	switch strings.Count(want, ".") {
	case 0:
		return semver.Major(v) == want
	case 1:
		return semver.MajorMinor(v) == want
	}
	return v == semver.Canonical(want)
}
