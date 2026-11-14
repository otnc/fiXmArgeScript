package interp

import (
	"errors"
	"slices"
	"time"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// 非同期は JavaScript と同じく1本の流れで動かす (spec/async.md)。
//
// itsver の関数を1回呼ぶごとに task を1つ作り、goroutine で動かす。ただし動くのはいつも1つだけで、
// praylee のところでだけ次の task に順番を渡す。順番は決まった規則 (runq の先頭から) で決めるので、やり直しでも同じになる。
// 外の世界とのやり取り (タイマー、通信、子プロセスの終わり) は別の goroutine で待ち、終わった知らせを events で受け取る。
// どの知らせをどの順で受け取ったかは記録し、やり直しでは記録どおりに受け取る。

// task は itsver の関数を1回呼んだもの。ファイルの一番外側も1つの task (main) として扱う。
type task struct {
	frames []frame
	// wake は順番が来たことを知らせる。nil なら続け、エラーならそのエラーで止まる (プログラム全体を止めるとき)。
	// main には自分で知らせることもあるので、1つだけためられるようにする
	wake chan error
	done chan struct{} // goroutine が終わったら閉じる
	// waitFile と waitPos は今待っている praylee の場所。待ちぼうけの Co-Jundic に使う
	waitFile string
	waitPos  token.Pos
}

func newTask(frames []frame) *task {
	return &task{frames: frames, wake: make(chan error, 1), done: make(chan struct{})}
}

// errAbandoned は、やり直しや終了で、止まっている task を片付けるときに渡す。youlugive でも ennote でも扱わない。
var errAbandoned = errors.New("abandoned")

// event は外の処理が終わった知らせ。op は startOp で付けた通し番号。
type event struct {
	op  int
	v   any
	err error
}

// spawn は itsver の関数の本体を新しい task で動かし、結果の Promise を返す。
// JavaScript の async function と同じく、本体は最初の praylee まですぐに動き、そこで呼んだ側に戻る。
func (in *Interp) spawn(fr frame, body func() (Value, error)) (*Promise, error) {
	p := &Promise{}
	// 呼び出しの深さを数え続けるため、呼んだ側の frames を引き継ぐ
	t := newTask(append(slices.Clone(in.frames), fr))
	in.tasks = append(in.tasks, t)
	me := in.cur
	in.runq = append([]*task{me}, in.runq...)
	go in.runTask(t, p, body)
	in.resume(t)
	return p, in.waitTurn(me)
}

func (in *Interp) runTask(t *task, p *Promise, body func() (Value, error)) {
	defer close(t.done)
	if err := <-t.wake; err != nil {
		return
	}
	result, err := body()
	if errors.Is(err, errAbandoned) {
		return
	}
	in.tasks = slices.DeleteFunc(in.tasks, func(x *task) bool { return x == t })
	if _, ok := catchable(err); err != nil && (!ok || p == nil) {
		// 巻き戻し、終了、検査のエラーなどは、プログラム全体の流れなので main に渡す。
		// 結果を受け取る Promise がない task (callback) は、どのエラーでも止める
		in.stopWith(err)
		return
	}
	if p != nil {
		in.settle(p, result, err)
	}
	next, err := in.nextTask()
	if err != nil {
		in.stopWith(err)
		return
	}
	in.resume(next)
}

// resume は task t に順番を渡す。渡した後の goroutine は Interp に触らない。
func (in *Interp) resume(t *task) {
	in.cur.frames = in.frames
	in.cur = t
	in.frames = t.frames
	t.wake <- nil
}

func (in *Interp) waitTurn(me *task) error {
	return <-me.wake
}

// suspend は今の task を止めて次の task に順番を渡し、また順番が来るまで待つ。
// 呼ぶ前に、今の task を runq か Promise の waiters に入れておく。
func (in *Interp) suspend() error {
	me := in.cur
	next, err := in.nextTask()
	if err != nil {
		in.stopWith(err)
		return in.waitTurn(me)
	}
	if next == me {
		return nil
	}
	in.resume(next)
	return in.waitTurn(me)
}

// nextTask は次に動かす task を runq から取る。runq が空なら外の処理が終わるのを待つ。
// 待つものもなければ、もう誰も再開できないので、main が待っている praylee の場所で Co-Jundic にする。
func (in *Interp) nextTask() (*task, error) {
	for len(in.runq) == 0 {
		if in.pendingOps == 0 {
			return nil, diag.New(diag.CoJundic, in.main.waitFile, in.main.waitPos, "")
		}
		in.deliver(in.nextEvent())
	}
	t := in.runq[0]
	in.runq = in.runq[1:]
	return t, nil
}

// stopWith はプログラム全体を止めるエラーを main に渡す。main はそのエラーで今の場所から戻る。
func (in *Interp) stopWith(err error) {
	in.main.wake <- err
}

// await は praylee で Promise を待つ。JavaScript と同じく、もう終わっていても一度ほかの task に順番を回す。
func (in *Interp) await(p *Promise, pos token.Pos) (Value, error) {
	p.awaited = true
	me := in.cur
	if p.settled {
		in.runq = append(in.runq, me)
	} else {
		p.waiters = append(p.waiters, me)
	}
	me.waitFile, me.waitPos = in.frame().file, pos
	if err := in.suspend(); err != nil {
		return nil, err
	}
	return p.Result, p.Err
}

// settle は Promise を終わらせ、待っていた task を runq に入れる。
func (in *Interp) settle(p *Promise, result Value, err error) {
	if result == nil {
		result = null
	}
	p.settled, p.Result, p.Err = true, result, err
	in.runq = append(in.runq, p.waiters...)
	p.waiters = nil
	if err != nil {
		in.rejected = append(in.rejected, p)
	}
}

// rejectedPromise は、もう失敗した Promise を返す。
func (in *Interp) rejectedPromise(err error) *Promise {
	p := &Promise{}
	in.settle(p, nil, err)
	return p
}

// startOp は外の世界とやり取りする処理 work を別の goroutine で始め、結果の Promise を返す。
// convert は work の結果を Promise の結果にする。
// 始めた処理は取り消せないので、選択肢を確定する。やり直しのときは始め直さず、終わった知らせを記録から受け取る。
func (in *Interp) startOp(work func() (any, error), convert func(any, error) (Value, error)) *Promise {
	in.commit()
	v, _ := in.once(func() (any, error) {
		// 番号はやり直しても使い回さない。前のやり直しで始めた処理の知らせと取り違えないため
		id := in.nextOp
		in.nextOp++
		closed := in.closed
		go func() {
			v, err := work()
			select {
			case in.events <- event{id, v, err}:
			case <-closed:
			}
		}()
		return id, nil
	})
	p := &Promise{}
	in.ops[v.(int)] = func(v any, err error) {
		result, err := convert(v, err)
		in.settle(p, result, err)
	}
	in.pendingOps++
	return p
}

// timer は toim no saimoku の待ち。JavaScript の setTimeout と同じく、締め切りの早い順、同じなら作った順に終わる。
// goroutine ごとに寝かせると終わる順が入れ替わるので、締め切りの一覧を持って nextEvent で見る。
type timer struct {
	op       int
	deadline time.Time
}

// startTimer は d だけ待つ Promise を返す。
func (in *Interp) startTimer(d time.Duration) *Promise {
	in.commit()
	v, _ := in.once(func() (any, error) {
		id := in.nextOp
		in.nextOp++
		return timer{id, time.Now().Add(d)}, nil
	})
	t := v.(timer)
	// 締め切りが同じなら作った順 (番号の順) にする
	i, _ := slices.BinarySearchFunc(in.timers, t, func(a, b timer) int {
		if c := a.deadline.Compare(b.deadline); c != 0 {
			return c
		}
		return a.op - b.op
	})
	in.timers = slices.Insert(in.timers, i, t)
	p := &Promise{}
	in.ops[t.op] = func(any, error) { in.settle(p, null, nil) }
	in.pendingOps++
	return p
}

// nextEvent は外の処理が終わった知らせを1つ受け取る。タイマーの締め切りが先に来れば、そのタイマーの知らせにする。
func (in *Interp) nextEvent() event {
	v, _ := in.once(func() (any, error) {
		if len(in.timers) == 0 {
			return <-in.events, nil
		}
		first := in.timers[0]
		select {
		case e := <-in.events:
			return e, nil
		case <-time.After(time.Until(first.deadline)):
			return event{op: first.op}, nil
		}
	})
	return v.(event)
}

// deliver は知らせを Promise か listener に届ける。fixmarge hak でなかったことにした入力が始めた処理の知らせは、届け先がないので捨てる。
func (in *Interp) deliver(e event) {
	if handle, ok := in.listeners[e.op]; ok {
		handle(e.v)
		return
	}
	settle, ok := in.ops[e.op]
	if !ok {
		return
	}
	delete(in.ops, e.op)
	in.timers = slices.DeleteFunc(in.timers, func(t timer) bool { return t.op == e.op })
	in.pendingOps--
	settle(e.v, e.err)
}

// listen は何度も知らせが届く外の処理 (サーバー、見張り、シグナル、魔女の出力) を始める。
// start は send で知らせを送り、done が閉じたら止まる。知らせが届くたびに handle を呼ぶ。
// keepAlive なら、プログラムはこの処理がある限り終わらない (Node.js のサーバーと fs.watch と同じ)。
// 始めた処理は取り消せないので選択肢を確定し、やり直しでは始め直さない。
func (in *Interp) listen(keepAlive bool, start func(send func(any), done <-chan struct{}) error, handle func(any)) error {
	in.commit()
	v, err := in.once(func() (any, error) {
		id := in.nextOp
		in.nextOp++
		closed := in.closed
		send := func(v any) {
			select {
			case in.events <- event{op: id, v: v}:
			case <-closed:
			}
		}
		return id, start(send, closed)
	})
	if err != nil {
		return err
	}
	in.listeners[v.(int)] = handle
	if keepAlive {
		in.pendingOps++
	}
	return nil
}

// callback は外からの知らせで fxma の関数 f を呼ぶ task を作り、runq の最後に入れる。
// f が Promise を返したら待つ。after があれば f の結果を渡す (サーバーの返事)。
// f の中の捕まえていないエラーは、Node.js と同じくプログラム全体を止める。
func (in *Interp) callback(f Value, args []Value, file string, pos token.Pos, name string, after func(Value) error) {
	t := newTask([]frame{{abiryoku: true, file: file}})
	in.tasks = append(in.tasks, t)
	in.runq = append(in.runq, t)
	go in.runTask(t, nil, func() (Value, error) {
		v, err := in.call(f, args, pos, name)
		if err != nil {
			return nil, err
		}
		if p, ok := v.(*Promise); ok {
			if v, err = in.await(p, pos); err != nil {
				return nil, err
			}
		}
		if after != nil {
			err = after(v)
		}
		return v, err
	})
}

// drain は一番外側を実行し終えた後、ほかの task と外の処理が全部終わるまで待つ (Node.js と同じ)。
// praylee されないまま失敗した Promise があれば、そのエラーで止める (Node.js の unhandled rejection と同じ)。
func (in *Interp) drain() error {
	for {
		switch {
		case len(in.runq) > 0:
			in.runq = append(in.runq, in.cur)
			if err := in.suspend(); err != nil {
				return err
			}
		case in.pendingOps > 0:
			in.deliver(in.nextEvent())
		default:
			for _, p := range in.rejected {
				if !p.awaited {
					return p.Err
				}
			}
			in.rejected = nil
			return nil
		}
	}
}

// abandonTasks は止まっている task を1つずつ起こして片付ける。起こした task は errAbandoned で戻るだけで、
// ennote も実行しない。同時に2つ動かないよう、1つ片付け終わるのを待ってから次を起こす。
func (in *Interp) abandonTasks() {
	saved := in.frames
	for _, t := range in.tasks {
		in.frames = t.frames
		t.wake <- errAbandoned
		<-t.done
	}
	in.tasks = nil
	in.frames = saved
}
