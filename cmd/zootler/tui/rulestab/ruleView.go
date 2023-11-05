package rulestab

import (
	"sudonters/zootler/internal/astrender"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/rules/ast"

	"github.com/charmbracelet/lipgloss"
)

type changeRuleView entity.Model
type addRuleView struct {
	id    entity.Model
	v     ruleView
	focus bool
}

type ruleView struct {
	name    string
	rawRule string
	ast     ast.Expression
	err     error
}

var ruleRowNameStyle = lipgloss.NewStyle().Width(15).Align(lipgloss.Left)

func (r ruleView) View() string {
	var rows []string
	makeRow := func(name, content string) {
		rows = append(rows, lipgloss.JoinHorizontal(
			lipgloss.Top,
			ruleRowNameStyle.Render(name),
			"| ",
			content,
		))
	}

	makeRow("Name", r.name)
	makeRow("Raw Rule", r.rawRule)

	if r.err == nil {
		s := astrender.NewSexpr(astrender.LipglossColorScheme())
		ast.Visit(s, r.ast)
		makeRow("s-expr", s.String())

		p := astrender.NewPretty()
		ast.Visit(p, r.ast)
		makeRow("Expanded", p.String())
	} else {
		makeRow("Error", r.err.Error())
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
