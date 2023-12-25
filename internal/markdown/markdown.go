package markdown

import (
	"fmt"
	"strings"
)

type Document struct {
	sb strings.Builder
}

func (m *Document) String() string {
	return m.sb.String()
}

func (m *Document) Code(language, code string, args ...any) *Document {
	m.sb.WriteString("```")
	m.sb.WriteString(language)
	m.sb.WriteString("\n")
	m.sb.WriteString(strings.ReplaceAll(fmt.Sprintf(code, args...), "`", "\\`"))
	m.sb.WriteString("```")
	return m
}

func (m *Document) InlineCode(code string, args ...any) *Document {
	m.sb.WriteString("`")
	m.sb.WriteString(strings.ReplaceAll(fmt.Sprintf(code, args...), "`", "\\`"))
	m.sb.WriteString("`")
	return m
}

func (m *Document) Normal(s string) *Document {
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "{", "\\{")
	s = strings.ReplaceAll(s, "}", "\\}")
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	s = strings.ReplaceAll(s, "#", "\\#")
	s = strings.ReplaceAll(s, "+", "\\+")
	s = strings.ReplaceAll(s, "-", "\\-")
	s = strings.ReplaceAll(s, ".", "\\.")
	s = strings.ReplaceAll(s, "!", "\\!")
	m.sb.WriteString(s)
	return m
}
