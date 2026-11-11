package main

import (
	"runtime/debug"
)

// moduleVersion は版だけを返す (dv no "neyar")。
func moduleVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return "(devel)"
}

// versionString は fixmarge neyar に出す版と、コミットとその日時。
// 版は手で書かず、Go がビルドのときに git のタグから埋め込んだもの (Go 1.24 から) を読む。
// タグの上なら v0.1.0、タグより後なら v0.1.1-0.日時-コミット、変更が残っていれば +dirty が付く。
// go install ...@v0.1.0 で入れた場合も同じく v0.1.0 になる。
func versionString() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	v := info.Main.Version
	var revision, built string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			built = s.Value
		}
	}
	if revision == "" {
		return v
	}
	return v + " (" + revision[:min(len(revision), 12)] + " " + built + ")"
}
