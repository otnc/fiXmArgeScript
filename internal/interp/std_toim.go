package interp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// toimPackage は fxma:toim (時間)。docs/std/toim.md。
// 時刻は 1970年1月1日 (UTC) からのミリ秒の整数で表す (JavaScript の Date.now)。
func toimPackage() *Module {
	return newStd(map[string]stdDef{
		"ieywoo":    {fn: toimIeywoo, abiryoku: true},
		"saimoku":   {fn: toimSaimoku, abiryoku: true},
		"confsoy":   {fn: toimConfsoy},
		"salook":    {fn: toimSalook},
		"fineclowt": {fn: toimFineclowt},
	}, nil)
}

func toimIeywoo(a *argv) (Value, error) {
	if err := a.count(0, 0); err != nil {
		return nil, err
	}
	return a.in.recorded(func() Value { return Int(a.in.now().UnixMilli()) }), nil
}

// toimSaimoku は待つ。本物の非同期を作るまでは、呼んだ時点で待ち終えてから Promise を返す。
// やり直しのときは、もう待ったので待たない。
func toimSaimoku(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	ms, err := a.num(0)
	if err != nil {
		return nil, err
	}
	a.in.recorded(func() Value {
		// setTimeout と同じく、負の数は 0 として扱う
		time.Sleep(time.Duration(max(0, ms) * float64(time.Millisecond)))
		return null
	})
	return &Promise{Result: null}, nil
}

// timeArg は i 番目の引数の時刻を、タイムゾーンを当てた time.Time にする。
// 最後の引数が "utc" なら UTC で考え、その引数を取り除く。
func timeArg(a *argv, maxArgs int) (time.Time, error) {
	loc := a.in.location()
	if n := len(a.vals); n > 1 {
		if s, ok := a.vals[n-1].(Str); ok && s == "utc" {
			loc = time.UTC
			a.vals = a.vals[:n-1]
		}
	}
	if err := a.count(1, maxArgs); err != nil {
		return time.Time{}, err
	}
	ms, err := a.index(0)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(ms).In(loc), nil
}

func toimConfsoy(a *argv) (Value, error) {
	t, err := timeArg(a, 1)
	if err != nil {
		return nil, err
	}
	parts := []int{t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond() / int(time.Millisecond), int(t.Weekday())}
	out := make([]Value, len(parts))
	for i, p := range parts {
		out[i] = Int(p)
	}
	return NewArray(out), nil
}

// isoLayout は salook(t) の形。JavaScript の toISOString と違い、時差を付ける。
const isoLayout = "YYYY-MM-DD[T]HH:mm:ss.SSSZ"

func toimSalook(a *argv) (Value, error) {
	t, err := timeArg(a, 2)
	if err != nil {
		return nil, err
	}
	layout := isoLayout
	if a.has(1) {
		if layout, err = a.str(1); err != nil {
			return nil, err
		}
	}
	return Str(formatTime(t, layout)), nil
}

// timeTokens は書式の部品。長いものから試す (Day.js の format と同じ書き方)。
var timeTokens = []struct {
	token  string
	format func(t time.Time) string
}{
	{"YYYY", func(t time.Time) string { return fmt.Sprintf("%04d", t.Year()) }},
	{"SSS", func(t time.Time) string { return fmt.Sprintf("%03d", t.Nanosecond()/int(time.Millisecond)) }},
	{"MM", func(t time.Time) string { return fmt.Sprintf("%02d", int(t.Month())) }},
	{"DD", func(t time.Time) string { return fmt.Sprintf("%02d", t.Day()) }},
	{"HH", func(t time.Time) string { return fmt.Sprintf("%02d", t.Hour()) }},
	{"mm", func(t time.Time) string { return fmt.Sprintf("%02d", t.Minute()) }},
	{"ss", func(t time.Time) string { return fmt.Sprintf("%02d", t.Second()) }},
	{"M", func(t time.Time) string { return strconv.Itoa(int(t.Month())) }},
	{"D", func(t time.Time) string { return strconv.Itoa(t.Day()) }},
	{"H", func(t time.Time) string { return strconv.Itoa(t.Hour()) }},
	{"m", func(t time.Time) string { return strconv.Itoa(t.Minute()) }},
	{"s", func(t time.Time) string { return strconv.Itoa(t.Second()) }},
	{"d", func(t time.Time) string { return strconv.Itoa(int(t.Weekday())) }},
	{"Z", func(t time.Time) string { return t.Format("-07:00") }},
}

func formatTime(t time.Time, layout string) string {
	var b strings.Builder
	for rest := layout; rest != ""; {
		if strings.HasPrefix(rest, "[") {
			if end := strings.IndexByte(rest, ']'); end > 0 {
				b.WriteString(rest[1:end])
				rest = rest[end+1:]
				continue
			}
		}
		matched := false
		for _, tok := range timeTokens {
			if strings.HasPrefix(rest, tok.token) {
				b.WriteString(tok.format(t))
				rest = rest[len(tok.token):]
				matched = true
				break
			}
		}
		if !matched {
			b.WriteByte(rest[0])
			rest = rest[1:]
		}
	}
	return b.String()
}

// parseLayouts は fineclowt が読める形。JavaScript の Date.parse と同じく、時差のない日時はそのパソコンの時刻、日付だけなら UTC として読む。
var parseLayouts = []struct {
	layout string
	local  bool
}{
	{time.RFC3339Nano, false},
	{"2006-01-02T15:04:05.999999999", true},
	{"2006-01-02T15:04", true},
	{"2006-01-02", false},
}

func toimFineclowt(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	s, err := a.str(0)
	if err != nil {
		return nil, err
	}
	for _, p := range parseLayouts {
		loc := time.UTC
		if p.local {
			loc = a.in.location()
		}
		if t, err := time.ParseInLocation(p.layout, s, loc); err == nil {
			return Int(t.UnixMilli()), nil
		}
	}
	return nil, a.in.errorf(diag.AyAmis, a.pos, inner(Str(s)))
}
