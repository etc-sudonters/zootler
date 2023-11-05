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
	"github.com/charmbracelet/lipgloss"
)

func createRuleList(p entity.Pool, s tea.WindowSizeMsg) listpanel.Model {
	haveRules, err := p.Query([]entity.Selector{
		entity.With[logic.RawRule]{},
		entity.With[world.Name]{},
	})

	if err != nil {
		panic(err)
	}

	items := make([]list.Item, 0, len(haveRules))

	for _, ent := range haveRules {
		var rule logic.RawRule
		var name world.Name

		if err := ent.Get(&rule); err != nil {
			panic(err)
		}
		if err := ent.Get(&name); err != nil {
			panic(err)
		}

		item := hasRuleItem{
			raw:    string(rule),
			name:   string(name),
			entity: ent,
		}
		items = append(items, item)
	}

	l := list.New(items, hasRuleDelegate{}, s.Width, s.Height)
	l.Title = fmt.Sprintf("Available Rules")
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.SetShowPagination(false)

	panel := listpanel.New(listpanel.WithList(l))
	(&panel).Focus()
	return panel
}

type hasRuleItem struct {
	name      string
	raw       string
	hasParsed bool
	entity    entity.View
}

func (d hasRuleItem) FilterValue() string {
	tags := d.name
	if d.hasParsed {
		tags += " parsed"
	}
	return tags
}

type hasRuleDelegate struct {
}

func (d hasRuleDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var buf strings.Builder
	hasRule := item.(hasRuleItem)

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

	name := hasRule.name
	if strings.Contains(name, "->") {
		var fromName world.FromName
		var toName world.ToName

		if err := hasRule.entity.Get(&fromName); err != nil {
			panic(err)
		}

		if err := hasRule.entity.Get(&toName); err != nil {
			panic(err)
		}
		name = fmt.Sprintf("%s\n%s", fromName, toName)
	} else {
		buf.WriteString(hasRule.name)
		buf.WriteRune('\n')
	}

	fmt.Fprint(w, lipgloss.JoinHorizontal(lipgloss.Top, buf.String(), name))
}

func (d hasRuleDelegate) Height() int {
	return 2
}

func (d hasRuleDelegate) Spacing() int {
	return 1
}

func (d hasRuleDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	parseRule := key.NewBinding(key.WithKeys("P"))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, parseRule):
			if m.SelectedItem() == nil {
				break
			}
			item := m.SelectedItem().(hasRuleItem)
			if item.hasParsed {
				cmds = append(cmds, func() tea.Msg { return changeRuleView(item.entity.Model()) })
			} else {
				item.hasParsed = true
				cmds = append(cmds, m.SetItem(m.Index(), item))
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
