package rulestab

import (
	"fmt"
	"io"
	"strings"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/tui/listpanel"
	"sudonters/zootler/pkg/logic"
	"sudonters/zootler/pkg/rules/parser"
	"sudonters/zootler/pkg/world"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func createRuleList(p entity.Pool, s tea.WindowSizeMsg) listpanel.Model {
	haveRules, err := p.Query([]entity.Selector{
		entity.With[logic.RawRule]{},
	})
	if err != nil {
		panic(err)
	}

	items := make([]list.Item, len(haveRules))

	for i, hasRule := range haveRules {
		var raw logic.RawRule
		var name world.Name
		hasRule := hasRule
		hasRule.Get(&name)
		hasRule.Get(&raw)
		items[i] = hasRuleItem{
			name:   string(name),
			raw:    string(raw),
			entity: hasRule,
		}
	}

	l := list.New(items, hasRuleDelegate{}, s.Width, s.Height)
	l.Title = fmt.Sprintf("Available Rules")
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	panel := listpanel.New(listpanel.WithList(l))
	panel.Focus()
	return panel
}

type hasRuleItem struct {
	name      string
	raw       string
	hasParsed bool
	entity    entity.View
}

func (d hasRuleItem) FilterValue() string {
	return d.name
}

type hasRuleDelegate struct {
}

func (d hasRuleDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var buf strings.Builder
	hasRule := m.SelectedItem().(hasRuleItem)

	if hasRule.hasParsed {
		buf.WriteRune('P')
	} else {
		buf.WriteRune(' ')
	}

	buf.WriteRune(' ')

	if index == m.Index() {
		buf.WriteRune('>')
	} else {
		buf.WriteRune(' ')
	}

	buf.WriteRune(' ')
	buf.WriteString(hasRule.name)
	fmt.Fprint(w, buf.String())
}

func (d hasRuleDelegate) Height() int {
	return 1
}

func (d hasRuleDelegate) Spacing() int {
	return 0
}

func (d hasRuleDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	item := m.SelectedItem().(hasRuleItem)
	parseRule := key.NewBinding(key.WithKeys("P"))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, parseRule):
			if item.hasParsed {
				cmds = append(cmds, func() tea.Msg { return changeRuleView(item.entity.Model()) })
			} else {
				cmds = append(cmds, func() tea.Msg {
					ast, err := parser.Parse(item.raw)
					return addRuleView{
						id: item.entity.Model(),
						v: ruleView{
							name:    item.name,
							rawRule: item.raw,
							ast:     ast,
							err:     err,
						},
						focus: true,
					}
				})
			}
		}
	}

	return tea.Batch(cmds...)
}
