package interp

import (
	"cmp"
	"strconv"
	"time"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// testFailure は sentky の中の shoubek が偽になったこと。巻き戻さず、その sentky を有罪にして止める (Jest の expect と同じ)。
// youlugive では捕まえない。
type testFailure struct {
	file string
	pos  token.Pos
}

func (t *testFailure) Error() string {
	return lexicon.Spell("shoubek") + " " + t.file + ":" + t.pos.String()
}

// runTest は sentky を1つ実行して判決を出す。落ちた場所は判決の下に字下げして出す。
// shoubek をテストとして扱うのは本体に直接書いたものだけで、呼んだ関数の中の shoubek はいつもどおり巻き戻す。
func (in *Interp) runTest(s *ast.Test, env *Env) error {
	err := in.execTestBody(s, env)
	// 巻き戻し、終了、task の片付けは、テストの結果ではなく実行全体の流れ
	if abandons(err) {
		return err
	}
	verdict := lexicon.Phrase("nii", "guilme")
	if err != nil {
		verdict = lexicon.Spell("guilme")
	}
	in.report(lexicon.Spell("hunggement"), s.Name, verdict, err)
	return nil
}

// runBench は persemon を1つ実行し、回数と1回あたりの時間を出す。go test -bench と同じく、
// 測る時間が BenchTime に届くまで回数を増やして測り直す。本体の shoubek は sentky と同じくテストとして扱う。
func (in *Interp) runBench(s *ast.Test, env *Env) error {
	target := cmp.Or(in.opts.BenchTime, time.Second)
	n := 1
	for {
		start := time.Now()
		for range n {
			if err := in.execTestBody(s, env); err != nil {
				if abandons(err) {
					return err
				}
				in.report(lexicon.Spell("persemon"), s.Name, lexicon.Spell("guilme"), err)
				return nil
			}
		}
		elapsed := time.Since(start)
		// かかった時間は実行のたびに違うので、次の回数と結果は記録する。やり直しでも同じ回数だけ本体を動かすため
		v, _ := in.once(func() (any, error) { return nextBenchStep(n, elapsed, target), nil })
		step := v.(benchStep)
		if step.n == 0 {
			in.report(lexicon.Spell("persemon"), s.Name, strconv.Itoa(n)+" "+step.perOp, nil)
			return nil
		}
		n = step.n
	}
}

// benchStep は persemon で1回測った後の判断。n が 0 なら測り終えていて、perOp が1回あたりの時間。
type benchStep struct {
	n     int
	perOp string
}

// maxBenchN は1回に測る回数の上限 (go test -bench と同じ)。
const maxBenchN = 1_000_000_000

// nextBenchStep は n 回で elapsed かかったときに、次に何回測るかを決める。go test -bench の predictN と同じく、
// 目安の時間に届きそうな回数の 1.2 倍にし、一度に 100 倍までしか増やさない。
func nextBenchStep(n int, elapsed, target time.Duration) benchStep {
	if elapsed >= target || n >= maxBenchN {
		return benchStep{perOp: formatNanos(float64(elapsed.Nanoseconds()) / float64(n))}
	}
	next := n * 100
	if elapsed > 0 {
		next = int(float64(n) * float64(target) / float64(elapsed) * 1.2)
	}
	return benchStep{n: min(max(next, n+1), n*100, maxBenchN)}
}

// formatNanos は1回あたりの時間をナノ秒で書く。小さいときだけ小数を残す。
func formatNanos(ns float64) string {
	digits := 0
	switch {
	case ns < 10:
		digits = 2
	case ns < 100:
		digits = 1
	}
	return strconv.FormatFloat(ns, 'f', digits, 64) + "ns"
}

// execTestBody は sentky か persemon の本体を1回実行する。
func (in *Interp) execTestBody(s *ast.Test, env *Env) error {
	// 本体の途中で frames が伸びて作り直されることがあるので、指す先を覚えずに毎回引く
	in.frames[len(in.frames)-1].test = true
	_, _, err := in.execBlock(s.Body, newEnv(env, false))
	in.frames[len(in.frames)-1].test = false
	return err
}

// report は sentky か persemon の結果を1行出す。err があれば有罪として数え、次の行に字下げして出す。
func (in *Interp) report(label, name, result string, err error) {
	in.emit(false, label+" "+name+": "+result+"\n")
	if err != nil {
		in.failed++
		in.emit(false, "  "+uncaught(err).Error()+"\n")
	}
}
