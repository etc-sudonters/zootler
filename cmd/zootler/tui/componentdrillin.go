package tui

import (
	"fmt"
	"io"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/bitpool"
	"sudonters/zootler/internal/entity/componenttable"
	"sudonters/zootler/internal/tui/listpanel"
	"sudonters/zootler/internal/tui/panels"
	"sudonters/zootler/pkg/world"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func createComponentDrillIn(r componenttable.RowData, pool entity.Pool, s panels.Size) listpanel.Model {
	listItems := make([]list.Item, 0, r.Len())
	comps := r.Components()

	for comps.MoveNext() {
		cur := comps.Current()
		listItems = append(listItems, drillinItem{cur})
	}

	l := list.New(listItems, drillInDelegate{pool}, s.Width, s.Height)
	l.Title = fmt.Sprintf("Component: %s (%d/%d)", r.Type().Name(), r.Len(), r.Capacity())
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	return listpanel.New(listpanel.WithList(l))
}

type drillinItem struct {
	componenttable.RowEntry
}

func (c drillinItem) FilterValue() string { return fmt.Sprintf("%+v", c.Component) }

type drillInDelegate struct {
	p entity.Pool
}

func (c drillInDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var name world.Name

	component := item.(drillinItem)
	c.p.Get(component.Entity, []interface{}{&name})
	repr := fmt.Sprintf("%q (Entity: %d)\n%+v", name, component.Entity, component.Component)

	render := itemStyle.Render
	if index == m.Index() {
		render = func(s ...string) string {
			return selectedStyle.Render(s...)
		}
	}

	fmt.Fprint(w, render(repr))
}

func (c drillInDelegate) Height() int {
	return 2
}

func (c drillInDelegate) Spacing() int {
	return 0
}

func (c drillInDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			return panels.CreateListPanel(func(s panels.Size) listpanel.Model {
				tbl, err := bitpool.ExtractComponentTable(c.p)
				if err != nil {
					panic(err)
				}
				item := m.SelectedItem().(drillinItem)
				return createComponentsForEntity(item.Entity, tbl, s)
			})
		}
	}

	return nil
}
