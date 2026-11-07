// Package ast は構文木の型を持つ。
package ast

import "github.com/otnc/fiXmArgeScript/internal/token"

type Node interface {
	Position() token.Pos
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

// At は構文木の各部品に位置を持たせるために埋め込む。
type At struct{ Pos token.Pos }

func (a At) Position() token.Pos { return a.Pos }

// 式

type (
	Ident struct {
		At
		Name string
	}
	IntLit struct {
		At
		Value int64
	}
	FloatLit struct {
		At
		Value float64
	}
	// StringLit は文字列。Parts のうち Expr が nil でないものは #{} の中の式。
	StringLit struct {
		At
		Parts []StringPiece
	}
	BoolLit struct {
		At
		Value bool
	}
	NullLit  struct{ At }
	ArrayLit struct {
		At
		Elems []Expr
	}
	DictLit struct {
		At
		Keys   []Expr
		Values []Expr
	}
	// Paren は括弧で囲んだ式。mekamo の候補の区切りを括弧の中まで広げないために残す。
	Paren struct {
		At
		X Expr
	}
	Unary struct {
		At
		Op string // "-" "+" "~" "nii"
		X  Expr
	}
	Binary struct {
		At
		Op   string // 算術とビットの記号
		X, Y Expr
	}
	// Logical は sond と owk。
	Logical struct {
		At
		Op   string // "sond" "owk"
		X, Y Expr
	}
	// Compare は a < b <= c のようにつながる比較。Ops は "<" "senaji" "exi" など。
	Compare struct {
		At
		Ops      []string
		Operands []Expr
	}
	// Exists は右に何もない x exi。
	Exists struct {
		At
		X Expr
	}
	// Index は A no B。B が識別子で A がモジュールなら、B はモジュールの中の名前として読む。
	Index struct {
		At
		X   Expr
		Key Expr
	}
	// Slice は A no lo..hi。lo と hi は省ける。
	Slice struct {
		At
		X         Expr
		Lo, Hi    Expr
		Inclusive bool
	}
	Range struct {
		At
		Lo, Hi    Expr
		Step      Expr // batte。省くと nil
		Inclusive bool
	}
	Call struct {
		At
		Fn   Expr
		Args []Expr
	}
	FuncLit struct {
		At
		Name     string // 名前のない関数なら空
		Params   []string
		Body     []Stmt
		Abiryoku bool
		Itsver   bool
	}
	// Choice は a owk b owk c mekamo。
	Choice struct {
		At
		Options []Expr
	}
	Await struct {
		At
		X Expr
	}
	// Input は zanfession。Prompt は省くと nil。
	Input struct {
		At
		Prompt Expr
	}
)

type StringPiece struct {
	Text string
	Expr Expr
}

// 文

type (
	ExprStmt struct {
		At
		X Expr
	}
	// Assign は izu。Op は複合代入の記号で、ただの代入なら空。
	Assign struct {
		At
		Target Expr
		Op     string
		Value  Expr
		Const  bool // showence
		Export bool // oqun
	}
	Print struct {
		At
		X      Expr
		Stderr bool // canosiri
	}
	Return struct {
		At
		X Expr // 省くと nil
	}
	Throw struct {
		At
		X Expr
	}
	Assert struct {
		At
		X Expr
	}
	If struct {
		At
		Conds  []Expr
		Blocks [][]Stmt
		Else   []Stmt // siva がなければ nil
	}
	While struct {
		At
		Cond Expr
		Body []Stmt
	}
	For struct {
		At
		Iter Expr
		Var  string
		Body []Stmt
	}
	FuncDecl struct {
		At
		Fn     *FuncLit
		Export bool
	}
	Try struct {
		At
		Body       []Stmt
		HasCatch   bool
		CatchVar   string // 省くと空
		Catch      []Stmt
		HasFinally bool
		Finally    []Stmt
	}
	Test struct {
		At
		Name string
		Body []Stmt
	}
	Continue struct{ At }
	Break    struct{ At }
	// Import は werer。Alias があれば名前を付けて読み込み、なければ Names を個別に読み込む。
	Import struct {
		At
		Path  string
		Alias string
		Names []ImportName
	}
)

type ImportName struct {
	Name  string
	Alias string
}

func (*Ident) exprNode()     {}
func (*IntLit) exprNode()    {}
func (*FloatLit) exprNode()  {}
func (*StringLit) exprNode() {}
func (*BoolLit) exprNode()   {}
func (*NullLit) exprNode()   {}
func (*ArrayLit) exprNode()  {}
func (*DictLit) exprNode()   {}
func (*Paren) exprNode()     {}
func (*Unary) exprNode()     {}
func (*Binary) exprNode()    {}
func (*Logical) exprNode()   {}
func (*Compare) exprNode()   {}
func (*Exists) exprNode()    {}
func (*Index) exprNode()     {}
func (*Slice) exprNode()     {}
func (*Range) exprNode()     {}
func (*Call) exprNode()      {}
func (*FuncLit) exprNode()   {}
func (*Choice) exprNode()    {}
func (*Await) exprNode()     {}
func (*Input) exprNode()     {}

func (*ExprStmt) stmtNode() {}
func (*Assign) stmtNode()   {}
func (*Print) stmtNode()    {}
func (*Return) stmtNode()   {}
func (*Throw) stmtNode()    {}
func (*Assert) stmtNode()   {}
func (*If) stmtNode()       {}
func (*While) stmtNode()    {}
func (*For) stmtNode()      {}
func (*FuncDecl) stmtNode() {}
func (*Try) stmtNode()      {}
func (*Test) stmtNode()     {}
func (*Continue) stmtNode() {}
func (*Break) stmtNode()    {}
func (*Import) stmtNode()   {}
