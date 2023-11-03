package tui

import (
	"fmt"
	"io"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/componenttable"
	"sudonters/zootler/internal/mirrors"
	"sudonters/zootler/internal/tui/listpanel"
	"sudonters/zootler/internal/tui/panels"
	"sudonters/zootler/pkg/world"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func createComponentsForEntity(e entity.Model, tbl *componenttable.Table, s panels.Size) listpanel.Model {
	name, _ := tbl.Get(e, mirrors.TypeOf[world.Name]())

	listItems := make([]list.Item, 0, tbl.Len())
	rows := tbl.Rows()

	for rows.MoveNext() {
		r := rows.Current()
		v := r.Get(e)
		if v == nil {
			continue
		}
		listItems = append(listItems, entityComponentItem{
			name: r.Type().Name(),
			id:   r.Id(),
			v:    v,
		})
	}

	l := list.New(listItems, entityComponentDelegate{}, s.Width, s.Height)
	l.Title = fmt.Sprintf("Entity %s (%d)", name.(world.Name), len(listItems))
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	return listpanel.New(listpanel.WithList(l))
}

type entityComponentItem struct {
	name string
	id   entity.ComponentId
	v    entity.Component
}

func (e entityComponentItem) FilterValue() string {
	return fmt.Sprintf("%+v", e.v)
}

type entityComponentDelegate struct{}

func (e entityComponentDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(entityComponentItem)
	repr := fmt.Sprintf("Type: %q\tID: %d\n%+v", i.name, i.id, i.v)

	render := itemStyle.Render
	if index == m.Index() {
		render = func(s ...string) string {
			return selectedStyle.Render(s...)
		}
	}

	fmt.Fprint(w, render(repr))
}

func (e entityComponentDelegate) Height() int {
	return 2
}

func (e entityComponentDelegate) Spacing() int {
	return 0
}

func (e entityComponentDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
