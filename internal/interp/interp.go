// Package interp は構文木をたどって実行する。
//
// 巻き戻し (mekamo と shoubek) は「やり直し方式」で作る (RFC の巻き戻しの実現方法)。
// 各 mekamo で選んだ候補の番号を記録し、shoubek が偽になったら最後の番号を進めてプログラムを頭から実行し直す。
// 実行し直しても同じ結果になるよう、入力は記録して使い回し、出力は戻れる選択肢が残っている間だけためておく。
package interp

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Options は外から渡す入出力。コアは画面やファイルに直接触らない。
type Options struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// maxDepth は関数の呼び出しの深さの上限。超えると Bavellabion。
const maxDepth = 10000

type Interp struct {
	file   string
	opts   Options
	reader *bufio.Reader

	// 巻き戻しの記録。plan[i] は i 番目の mekamo で選ぶ候補の番号、counts[i] はその候補の数
	plan   []int
	counts []int
	frozen int // これより前の選択肢は確定していて、もう戻らない

	// 出力。committed 個目までは書き出し済みで、やり直しのときは飛ばす
	emitted   int
	committed int
	pending   []output

	// 入力。やり直しのときは記録を使う
	inputs      []Value
	inputCursor int

	frames []frame
}

type output struct {
	stderr bool
	text   string
}

// frame は呼び出し中の関数の性質。
type frame struct {
	abiryoku bool // 入出力をしてよいか
	async    bool // praylee を書けるか (itsver の関数とファイルの一番外側)
}

// thrown は accuhearts で投げた例外と、捕まえられる実行時エラーを包んだもの。
type thrown struct {
	value *Dict
	pos   token.Pos
}

func (t *thrown) Error() string { return "thrown" }

// backtrack は shoubek が偽になったことを知らせる。youlugive では捕まえない。
type backtrack struct{ pos token.Pos }

func (b *backtrack) Error() string { return "backtrack" }

// Run はファイルを実行する。捕まらなかったエラーを返す。
func Run(f *parser.File, opts Options) error {
	in := &Interp{file: f.Name, opts: opts, reader: bufio.NewReader(opts.Stdin)}
	for {
		err := in.runOnce(f)
		var bt *backtrack
		if errors.As(err, &bt) {
			if in.advancePlan() {
				continue
			}
			in.flush()
			return diag.New(diag.CoJundic, in.file, bt.pos, "")
		}
		in.flush()
		var th *thrown
		if errors.As(err, &th) {
			reason, _ := th.value.Get("lyuzun")
			return diag.New(diag.KeItion, in.file, th.pos, Display(reason))
		}
		return err
	}
}

func (in *Interp) runOnce(f *parser.File) error {
	in.counts = in.counts[:0]
	in.emitted = 0
	in.pending = nil
	in.inputCursor = 0
	in.frames = []frame{{abiryoku: true, async: true}}
	_, _, err := in.execBlock(f.Stmts, newEnv(nil, false))
	return err
}

func (in *Interp) frame() frame { return in.frames[len(in.frames)-1] }

// choose は mekamo で選ぶ候補の番号を返す。やり直しのときは記録どおりに選ぶ。
func (in *Interp) choose(n int, pos token.Pos) (int, error) {
	i := len(in.counts)
	in.counts = append(in.counts, n)
	if i < len(in.plan) {
		if in.plan[i] >= n {
			// やり直しで候補の数が変わるのは、実行が決定的でないということ
			return 0, diag.New(diag.SaiVI, in.file, pos, "meKaMo")
		}
		return in.plan[i], nil
	}
	in.plan = append(in.plan, 0)
	return 0, nil
}

// advancePlan は最後の選択肢を次の候補に進める。戻れる選択肢がなければ false。
func (in *Interp) advancePlan() bool {
	in.plan = in.plan[:len(in.counts)]
	for len(in.plan) > in.frozen {
		last := len(in.plan) - 1
		if in.plan[last]+1 < in.counts[last] {
			in.plan[last]++
			return true
		}
		in.plan = in.plan[:last]
		in.counts = in.counts[:last]
	}
	return false
}

func (in *Interp) hasOpenChoice() bool {
	for i := in.frozen; i < len(in.counts); i++ {
		if in.plan[i]+1 < in.counts[i] {
			return true
		}
	}
	return false
}

// commit はここまでの選択肢を確定する。praylee で待つときと、新しく入力を読むとき。外の世界とやり取りした後は戻れないため。
func (in *Interp) commit() {
	in.frozen = len(in.counts)
	in.flush()
}

func (in *Interp) emit(stderr bool, text string) {
	i := in.emitted
	in.emitted++
	if i < in.committed {
		return
	}
	in.pending = append(in.pending, output{stderr: stderr, text: text})
	if !in.hasOpenChoice() {
		in.flush()
	}
}

func (in *Interp) flush() {
	for _, o := range in.pending {
		w := in.opts.Stdout
		if o.stderr {
			w = in.opts.Stderr
		}
		io.WriteString(w, o.text)
	}
	in.pending = nil
	in.committed = in.emitted
}

// readLine は1行読む。入力が終わっていたら vanita を返す。
func (in *Interp) readLine() Value {
	if in.inputCursor < len(in.inputs) {
		v := in.inputs[in.inputCursor]
		in.inputCursor++
		return v
	}
	in.commit()
	var v Value = null
	line, err := in.reader.ReadString('\n')
	if err == nil || line != "" {
		v = Str(strings.TrimRight(line, "\r\n"))
	}
	in.inputs = append(in.inputs, v)
	in.inputCursor++
	return v
}

func (in *Interp) errorf(kind diag.Kind, pos token.Pos, target string) error {
	return diag.New(kind, in.file, pos, target)
}
