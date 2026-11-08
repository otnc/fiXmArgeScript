package lint

import (
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/scope"
)

// unused は Was·Da。使っていない変数と読み込み。_ で始まる名前は、わざと使わない名前として見ない (ESLint と同じ)。
// 引数と関数の名前、公開した名前は、使わなくても指摘しない。
func (l *linter) unused(stmts []ast.Stmt) {
	for _, b := range scope.Resolve(stmts).Bindings {
		if len(b.Reads) > 0 || b.Export || strings.HasPrefix(b.Name, "_") {
			continue
		}
		switch b.Kind {
		case scope.Func, scope.Param:
			continue
		}
		l.report(WasDa, b.Pos, b.Name)
	}
}
