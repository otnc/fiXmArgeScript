package modfile

import "testing"

func TestPick(t *testing.T) {
	versions := []string{"v1.0.0", "v1.2.0", "v1.2.1", "v1.3.0-beta.0", "v2.0.0-rc.0", "v2.0.0-rc.1"}
	cases := []struct{ spec, want string }{
		{"", "v1.2.1"}, // 正式版があれば、もっと新しいプレリリースより先
		{"1", "v1.2.1"},
		{"v1", "v1.2.1"},
		{"1.2", "v1.2.1"},
		{"1.2.0", "v1.2.0"},
		{"1.3", "v1.3.0-beta.0"}, // 合う正式版がなければプレリリース
		{"2", "v2.0.0-rc.1"},
		{"2.0.0-rc.0", "v2.0.0-rc.0"},
		{"2.0.0", ""}, // プレリリースは、そこまで書いたときだけ合う
		{"3", ""},
	}
	for _, c := range cases {
		if got := Pick(versions, c.spec); got != c.want {
			t.Errorf("Pick(%q) = %q, want %q", c.spec, got, c.want)
		}
	}
}

func TestSpecSyntax(t *testing.T) {
	for _, spec := range []string{"", "1", "v1", "1.2", "1.2.3", "1.2.3-rc.0", "v1.2.3-beta.0"} {
		if !ValidSpec(spec) {
			t.Errorf("ValidSpec(%q) = false", spec)
		}
	}
	for _, spec := range []string{"*", "latest", "1.", "1.2.3.4", "^1.2", "1-rc.0"} {
		if ValidSpec(spec) {
			t.Errorf("ValidSpec(%q) = true", spec)
		}
	}
	path, spec := SplitSpec("github.com/alice/strs/case@1.2")
	if path != "github.com/alice/strs/case" || spec != "1.2" {
		t.Errorf("SplitSpec = %q, %q", path, spec)
	}
	if path, spec := SplitSpec("github.com/alice/strs"); path != "github.com/alice/strs" || spec != "" {
		t.Errorf("SplitSpec = %q, %q", path, spec)
	}
}
