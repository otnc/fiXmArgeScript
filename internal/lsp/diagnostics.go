package lsp

import (
	"errors"

	"go.lsp.dev/protocol"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/lint"
)

const source = "fxma"

// diagnostics はエラーと lavi の指摘を LSP の診断にする。
// 文言はコンソールと同じくフィクスマージ語だけにし、位置はエディタが範囲で示すので書かない。
func (d *document) diagnostics() []protocol.Diagnostic {
	out := []protocol.Diagnostic{}
	if d.err != nil {
		out = append(out, d.errorDiagnostic(d.err))
	}
	for _, f := range d.findings {
		out = append(out, d.findingDiagnostic(f))
	}
	return out
}

func (d *document) errorDiagnostic(err error) protocol.Diagnostic {
	e, ok := errors.AsType[*diag.Error](err)
	if !ok {
		return protocol.Diagnostic{
			Severity: protocol.DiagnosticSeverityError,
			Source:   protocol.NewOptional(source),
			Message:  protocol.String(err.Error()),
		}
	}
	return protocol.Diagnostic{
		Range:    d.rangeAt(e.Pos),
		Severity: protocol.DiagnosticSeverityError,
		Code:     protocol.String(e.Kind.Spelling()),
		Source:   protocol.NewOptional(source),
		Message:  protocol.String(withTarget(e.Kind.Spelling(), e.Target)),
	}
}

func (d *document) findingDiagnostic(f lint.Finding) protocol.Diagnostic {
	rule := lexicon.Spell(f.Rule)
	diagnostic := protocol.Diagnostic{
		Range:    d.rangeAt(f.Pos),
		Severity: protocol.DiagnosticSeverityWarning,
		Code:     protocol.String(rule),
		Source:   protocol.NewOptional(source),
		Message:  protocol.String(withTarget(lexicon.Spell("lavi")+" "+rule, f.Target)),
	}
	// 使っていない変数は、エディタが薄く表示できるように印を付ける。
	// 位置が文の頭 (werer など) のこともあるので、範囲は名前そのものにする
	if f.Rule == lint.WasDa {
		diagnostic.Tags = protocol.NewDiagnosticTags(protocol.DiagnosticTagUnnecessary)
		if t, ok := d.identFrom(f.Pos, f.Target); ok {
			diagnostic.Range = d.tokenRange(t)
		}
	}
	return diagnostic
}

func withTarget(s, target string) string {
	if target == "" {
		return s
	}
	return s + " " + target
}
