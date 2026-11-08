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
	"math/rand/v2"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Options は外から渡す入出力。コアは画面やファイルに直接触らない。
type Options struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	// Load は werer で読み込むファイルを読み、構文解析と検査をして返す。ファイルがなければ fs.ErrNotExist を包んだエラーを返す。
	// nil ならファイルの読み込みは Ca-nay (標準パッケージは読める)。
	Load func(name string) (*parser.File, error)
	// WorkDir は今いる場所 (izanavited no defutai、memoid の相対パス)。空ならブラウザ版と同じく /。
	WorkDir string
	// FS は fxma:memoid が使うファイル。nil なら memoid の関数は Ca-nay (ブラウザ版)。
	FS afero.Fs
	// Process は fxma:ieviljio が使う今のプロセスの情報。nil なら空の情報で、魔女は生めない (ブラウザ版)。
	Process *Process
	// HTTPClient は fxma:gemon が使う。nil なら通信は Ca-nay。
	HTTPClient *http.Client
	// Now と Location は時刻とタイムゾーン (fxma:toim)。nil なら time.Now と time.Local。
	Now      func() time.Time
	Location *time.Location
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

	// 時刻と乱数の結果。入力と同じく、やり直しのときは記録を使う (std/README.md)
	records      []outcome
	recordCursor int
	rng          *rand.Rand

	frames []frame

	// testing は fxma coatain で動かしているか。failed は有罪になった sentky の数で、やり直しのたびに数え直す
	testing bool
	failed  int

	// 読み込んだモジュール。やり直しのたびに作り直す。構文木は使い回す
	modules map[string]*Module
	loading map[string]bool

	// greatWitch は ieviljio no diviljio() の大魔女。やり直しのたびに作り直す
	greatWitch *Dict
	parsed     map[string]*parser.File
}

type output struct {
	stderr bool
	text   string
}

// frame は呼び出し中の関数の性質。
type frame struct {
	abiryoku bool   // 入出力をしてよいか
	async    bool   // praylee を書けるか (itsver の関数とファイルの一番外側)
	file     string // 実行中の文が書かれたファイル。エラーの位置に使う
	test     bool   // sentky の本体を実行しているか。ここの shoubek は巻き戻さずテストとして扱う
}

// thrown は accuhearts で投げた例外と、捕まえられる実行時エラーを包んだもの。
type thrown struct {
	value *Dict
	file  string
	pos   token.Pos
}

func (t *thrown) Error() string { return "thrown" }

// backtrack は shoubek が偽になったことを知らせる。youlugive では捕まえない。
type backtrack struct {
	file string
	pos  token.Pos
}

func (b *backtrack) Error() string { return "backtrack" }

// Run はファイルを実行する。捕まらなかったエラーを返す。
func Run(f *parser.File, opts Options) error {
	_, err := run(f, opts, false)
	return err
}

// RunTests は _coatain.fxma を実行し、sentky ごとに判決を出す。有罪が1つでもあれば failed が true。
func RunTests(f *parser.File, opts Options) (failed bool, err error) {
	in, err := run(f, opts, true)
	return in.failed > 0, err
}

func run(f *parser.File, opts Options, testing bool) (*Interp, error) {
	in := &Interp{testing: testing, file: f.Name, opts: opts, reader: bufio.NewReader(opts.Stdin), parsed: map[string]*parser.File{}, rng: newRNG(rand.Uint64())}
	for {
		err := in.runOnce(f)
		if bt, ok := errors.AsType[*backtrack](err); ok {
			if in.advancePlan() {
				continue
			}
			in.flush()
			return in, diag.New(diag.CoJundic, bt.file, bt.pos, "")
		}
		in.flush()
		return in, uncaught(err)
	}
}

// uncaught は捕まらなかったエラーを表示する形にする。accuhearts で投げた例外だけは KeI·tion で、対象に理由を書く。
func uncaught(err error) error {
	if th, ok := errors.AsType[*thrown](err); ok {
		reason, _ := th.value.Get("lyuzun")
		return diag.New(diag.KeItion, th.file, th.pos, Display(reason))
	}
	return err
}

func (in *Interp) runOnce(f *parser.File) error {
	in.counts = in.counts[:0]
	in.emitted = 0
	in.pending = nil
	in.inputCursor = 0
	in.recordCursor = 0
	in.frames = []frame{{abiryoku: true, async: true, file: in.file}}
	in.modules = map[string]*Module{}
	in.loading = map[string]bool{}
	in.greatWitch = nil
	in.failed = 0
	_, _, err := in.execBlock(f.Stmts, newEnv(nil, false))
	return err
}

func (in *Interp) workDir() string {
	if in.opts.WorkDir == "" {
		return "/"
	}
	return filepath.ToSlash(in.opts.WorkDir)
}

func (in *Interp) frame() frame { return in.frames[len(in.frames)-1] }

// choose は mekamo で選ぶ候補の番号を返す。やり直しのときは記録どおりに選ぶ。
func (in *Interp) choose(n int, pos token.Pos) (int, error) {
	i := len(in.counts)
	in.counts = append(in.counts, n)
	if i < len(in.plan) {
		if in.plan[i] >= n {
			// やり直しで候補の数が変わるのは、実行が決定的でないということ
			return 0, in.errorf(diag.SaiVI, pos, "meKaMo")
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

// outcome は外の世界とやり取りした結果1つ。
type outcome struct {
	v   any
	err error
}

// once は外の世界とやり取りする処理 (時刻、乱数、ファイル) を1回だけ実行する。
// やり直しのときは処理をせず、前の結果を返す。ファイルを書くなどの取り消せない処理を二度しないためと、
// 前と同じ結果にするため。
// 記録する値は書き換えられないものだけにする。fxma の配列や辞書を記録すると、前の実行で書き換えた中身が残ってしまう。
func (in *Interp) once(produce func() (any, error)) (any, error) {
	if in.recordCursor < len(in.records) {
		o := in.records[in.recordCursor]
		in.recordCursor++
		return o.v, o.err
	}
	v, err := produce()
	in.records = append(in.records, outcome{v, err})
	in.recordCursor++
	return v, err
}

// recorded は失敗しない once。
func (in *Interp) recorded(produce func() Value) Value {
	v, _ := in.once(func() (any, error) { return produce(), nil })
	return v.(Value)
}

func (in *Interp) now() time.Time {
	if in.opts.Now != nil {
		return in.opts.Now()
	}
	return time.Now()
}

func (in *Interp) location() *time.Location {
	if in.opts.Location != nil {
		return in.opts.Location
	}
	return time.Local
}

func (in *Interp) errorf(kind diag.Kind, pos token.Pos, target string) error {
	return diag.New(kind, in.frame().file, pos, target)
}
