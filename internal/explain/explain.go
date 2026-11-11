// Package explain は fixmarge nai が出す日本語の説明を持つ。コンソールで日本語を出すのは fixmarge nai と fixmarge neyar だけ。
package explain

import (
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

type entry struct {
	Kind diag.Kind
	Text string
}

// entries は RFC のエラー一覧と同じ内容。
var entries = []entry{
	{diag.AyAmis, "過ち。構文が正しくありません。実行中に JSON・日付・数として読めない文字列を読んだときも出ます (こちらは youlugive で捕まえられます)。\n下書きのキーワードは小文字だけで書きます。等しさは == ではなく senaji、等しくないは nii a senaji b と書きます。"},
	{diag.NiiExi, "存在しない。定義されていない変数や関数を使ったか、配列の範囲の外や辞書にないキーを読みました。\nブロックの中で作った変数は finend の後では使えません。"},
	{diag.ImJune, "不純。型が合いません。数と文字列の足し算や、条件に tadact と ulieso 以外を書いたときに出ます。\n値があるかは x exi、空かどうかは nang(xs) senaji 0 で調べます。"},
	{diag.KYoNAihim, "虚無。0 で割りました。"},
	{diag.Bavellabion, "呪いの言葉があふれる。整数が64ビットの範囲を超えたか、関数の呼び出しや呪い (マクロ) の展開が深すぎます。"},
	{diag.BlasTOac, "冒涜的。証拠品 (showence) に再代入しました。証拠を書き換えるのは裁判への冒涜です。"},
	{diag.KiNbOO, "禁忌。abiryoku の付いていない関数の中で入出力をしたか、abiryoku の関数を呼びました。\n関数の定義に abiryoku を付けてください。"},
	{diag.CoJundic, "矛盾。どの選択肢 (mekamo) でも shoubek を満たせませんでした。praylee で待っているのに、ほかに動ける処理も外から届くものもなく、二度と再開できないとき (自分自身の結果を待つなど) にも出ます。youlugive では捕まえられません。"},
	{diag.KeItion, "処刑。accuhearts で投げた例外が、誰にも許されませんでした (youlugive で捕まりませんでした)。"},
	{diag.SuReady, "既に。同じ名前を二重に定義したか、関数の中から外の変数に izu したか、読み込みが循環しました。\n外の値を変えたいときは、戻り値か、配列や辞書を使います。"},
	{diag.CraeyTA, "壊れた。ダウンロードしたパッケージが fixmarge.sum と合いません。"},
	{diag.SaiVI, "災い。インタプリタ自体のバグです。起きたら報告してください。"},
	{diag.CaNay, "できない。その場所では使えない書き方か、まだ作られていない機能です。\npraylee は itsver の関数とファイルの一番外側で、sentky は _coatain.fxma のファイルで使えます。"},
	{diag.AccUhEarts, "告発。accuhearts で投げた例外の種類です。e no \"chidd\" で読むとこの名前になります。"},
	{diag.DowGuil, "疑惑 (警告)。方言の綴りを使いました。実行はできますが、清書では標準語に直ります。"},
}

// Lookup はエラー名から説明を探す。名前は清書の綴りでも骨格でもよい。
func Lookup(name string) (spelling, text string, ok bool) {
	skel := lexicon.Skeleton(name)
	for _, e := range entries {
		if lexicon.Skeleton(e.Kind.Spelling()) == skel {
			return e.Kind.Spelling(), e.Text, true
		}
	}
	return "", "", false
}
