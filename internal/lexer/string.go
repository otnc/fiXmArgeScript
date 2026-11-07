package lexer

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// lexString は "..." と """...""" を読む。#{} の中はコードとして範囲だけ覚え、読むのは構文解析とビルドに任せる。
func (l *lexer) lexString() error {
	pos, start := l.pos(), l.i
	quote := `"`
	if strings.HasPrefix(l.src[l.i:], `"""`) {
		quote = `"""`
	}
	l.advanceTo(l.i + len(quote))

	var parts []token.StringPart // ここでは文字の部分はエスケープを解く前の形
	litStart := l.i
	for {
		if l.i >= len(l.src) {
			return l.errorf(diag.AyAmis, pos, quote)
		}
		if strings.HasPrefix(l.src[l.i:], quote) {
			parts = append(parts, token.StringPart{Text: l.src[litStart:l.i]})
			l.advanceTo(l.i + len(quote))
			break
		}
		c := l.src[l.i]
		switch {
		case c == '\n' && quote == `"`:
			return l.errorf(diag.AyAmis, pos, quote)
		case c == '\\':
			l.advance()
			if l.i < len(l.src) {
				l.advance()
			}
		case c == '#' && l.peek(1) == '{':
			parts = append(parts, token.StringPart{Text: l.src[litStart:l.i]})
			l.advanceTo(l.i + 2)
			codeStart, codePos := l.i, l.pos()
			if err := l.skipCode(pos); err != nil {
				return err
			}
			parts = append(parts, token.StringPart{Text: l.src[codeStart:l.i], IsCode: true, Start: codeStart, Pos: codePos})
			l.advance() // }
			litStart = l.i
		default:
			l.advance()
		}
	}

	if quote == `"""` {
		dedent(parts)
	}
	var out []token.StringPart
	for _, p := range parts {
		if p.IsCode {
			out = append(out, p)
			continue
		}
		text, bad := unescape(p.Text)
		if bad != "" {
			return l.errorf(diag.AyAmis, pos, bad)
		}
		if text != "" {
			out = append(out, token.StringPart{Text: text})
		}
	}
	l.tokens = append(l.tokens, token.Token{Kind: token.String, Pos: pos, Start: start, End: l.i, Parts: out})
	return nil
}

// skipCode は #{ の後ろから、対応する } の手前まで進む。
func (l *lexer) skipCode(stringPos token.Pos) error {
	depth := 0
	for l.i < len(l.src) {
		switch l.src[l.i] {
		case '"':
			if err := l.skipQuoted(); err != nil {
				return err
			}
			continue
		case '{':
			depth++
		case '}':
			if depth == 0 {
				return nil
			}
			depth--
		}
		l.advance()
	}
	return l.errorf(diag.AyAmis, stringPos, "#{")
}

// skipQuoted は #{} の中に入れ子になった文字列を読み飛ばす。
func (l *lexer) skipQuoted() error {
	pos := l.pos()
	quote := `"`
	if strings.HasPrefix(l.src[l.i:], `"""`) {
		quote = `"""`
	}
	l.advanceTo(l.i + len(quote))
	for l.i < len(l.src) {
		if strings.HasPrefix(l.src[l.i:], quote) {
			l.advanceTo(l.i + len(quote))
			return nil
		}
		switch {
		case l.src[l.i] == '\\':
			l.advance()
		case l.src[l.i] == '#' && l.peek(1) == '{':
			l.advanceTo(l.i + 2)
			if err := l.skipCode(pos); err != nil {
				return err
			}
		}
		if l.i < len(l.src) {
			l.advance()
		}
	}
	return l.errorf(diag.AyAmis, pos, quote)
}

// dedent は """ の文字列から、最初と最後の改行と、全部の行に共通する字下げを取り除く (CoffeeScript と同じ)。
func dedent(parts []token.StringPart) {
	first, last := &parts[0], &parts[len(parts)-1]
	if i := strings.IndexByte(first.Text, '\n'); i >= 0 && strings.TrimLeft(first.Text[:i], " \t") == "" {
		first.Text = first.Text[i+1:]
	}
	if i := strings.LastIndexByte(last.Text, '\n'); i >= 0 && strings.TrimLeft(last.Text[i+1:], " \t") == "" {
		last.Text = last.Text[:i]
	}

	// コードの部分は字下げではないので、空白でない印に置き換えて行を調べる
	var flat strings.Builder
	for _, p := range parts {
		if p.IsCode {
			flat.WriteByte(0)
		} else {
			flat.WriteString(p.Text)
		}
	}
	indent := -1
	for _, line := range strings.Split(flat.String(), "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		if n := len(line) - len(trimmed); indent < 0 || n < indent {
			indent = n
		}
	}
	if indent <= 0 {
		return
	}

	atLineStart := true
	for i := range parts {
		p := &parts[i]
		if p.IsCode {
			atLineStart = false
			continue
		}
		var b strings.Builder
		for j := 0; j < len(p.Text); {
			if atLineStart {
				n := 0
				for n < indent && j < len(p.Text) && (p.Text[j] == ' ' || p.Text[j] == '\t') {
					j++
					n++
				}
				atLineStart = false
				continue
			}
			c := p.Text[j]
			b.WriteByte(c)
			j++
			if c == '\n' {
				atLineStart = true
			}
		}
		p.Text = b.String()
	}
}

// unescape はエスケープを解く。知らないエスケープがあれば、それを bad に返す。
func unescape(raw string) (text, bad string) {
	var b strings.Builder
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if c != '\\' {
			b.WriteByte(c)
			continue
		}
		if i+1 >= len(raw) {
			return "", `\`
		}
		i++
		switch raw[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case '\\', '"', '#':
			b.WriteByte(raw[i])
		case 'u':
			end := strings.IndexByte(raw[i:], '}')
			if !strings.HasPrefix(raw[i:], "u{") || end < 0 {
				return "", `\u`
			}
			n, err := strconv.ParseUint(raw[i+2:i+end], 16, 32)
			if err != nil || !utf8.ValidRune(rune(n)) {
				return "", raw[i-1 : i+end+1]
			}
			b.WriteRune(rune(n))
			i += end
		default:
			return "", raw[i-1 : i+1]
		}
	}
	return b.String(), ""
}
