package spheres

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/zecs"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type edgesTab struct {
	adultView list.Model
	childView list.Model
	selected  magicbean.Age
	rendered  bool

	chunkHeight, chunkWidth int
}

func newEdgesTab(parent *breakdown) edgesTab {
	var new edgesTab

	new.selected = magicbean.AgeAdult
	var renderer list.ItemDelegate = renderedStr{}

	var adultItems, childItems []list.Item
	if parent.sphere != nil {
		adultItems = makeEdgeItems(parent.sphere.Adult.Edges, parent.sphere.Edges)
		childItems = makeEdgeItems(parent.sphere.Child.Edges, parent.sphere.Edges)

	}
	new.adultView = list.New(adultItems, renderer, 0, 0)
	new.childView = list.New(childItems, renderer, 0, 0)
	listDefaults(&new.adultView)
	listDefaults(&new.childView)
	new.adultView.Styles.NoItems = lipgloss.NewStyle().SetString("No locations reached")
	new.childView.Styles.NoItems = lipgloss.NewStyle().SetString("No locations reached")
	(&new).setSize(tea.WindowSizeMsg{Height: parent.viewHeight, Width: parent.viewWidth})
	return new
}

func (this *edgesTab) setSize(size tea.WindowSizeMsg) {
	width := size.Width - 10
	height := ((size.Height - 10) / 2) - 4
	this.chunkHeight, this.chunkWidth = height, width
	(&this.adultView).SetHeight(height)
	(&this.adultView).SetWidth(width)
	(&this.childView).SetHeight(height)
	(&this.childView).SetWidth(width)
}

func (this *edgesTab) init(current *breakdown) tea.Cmd {
	return nil
}

func (this *edgesTab) update(msg tea.Msg, current *breakdown) tea.Cmd {
	if current.sphere == nil {
		return nil
	}

	var cmds []tea.Cmd
	addCmd := func(c tea.Cmd) { cmds = append(cmds, c) }

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.setSize(msg)
		goto exit
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeySpace:
			switch this.selected {
			case magicbean.AgeAdult:
				this.selected = magicbean.AgeChild
			case magicbean.AgeChild:
				this.selected = magicbean.AgeAdult
			}
			goto exit
		case tea.KeyEnter:
			var row displayRow
			switch this.selected {
			case magicbean.AgeAdult:
				row = this.adultView.SelectedItem().(displayRow)
			case magicbean.AgeChild:
				row = this.childView.SelectedItem().(displayRow)
			}
			addCmd(RequestDisassembly(row.id))
			goto exit
		}

	}

	switch this.selected {
	case magicbean.AgeAdult:
		var cmd tea.Cmd
		this.adultView, cmd = this.adultView.Update(msg)
		addCmd(cmd)
	case magicbean.AgeChild:
		var cmd tea.Cmd
		this.childView, cmd = this.childView.Update(msg)
		addCmd(cmd)
	}

exit:
	return tea.Batch(cmds...)
}

func (this edgesTab) view(w io.Writer, current *breakdown) {
	if current.sphere == nil {
		fmt.Fprintln(w, "no sphere loaded")
		return
	}

	style := edgeChunkStyle.Height(this.chunkHeight).Width(this.chunkWidth)

	fmt.Fprint(w, style.Render(this.renderChunk("Adult", this.selected == magicbean.AgeAdult, this.adultView)))
	fmt.Fprint(w, style.Render(this.renderChunk("Child", this.selected == magicbean.AgeChild, this.childView)))
}

func (this *edgesTab) renderChunk(name string, selected bool, list list.Model) string {
	var chunk strings.Builder
	fmt.Fprintf(&chunk, "%s:\n", name)
	if !selected {
		fmt.Fprintln(&chunk, inactiveStyle.Render("-----collasped-----"))
	} else {
		fmt.Fprintln(&chunk, list.View())
	}
	return chunk.String()
}

type displayRow struct {
	display string
	id      zecs.Entity
}

type renderedStr struct{}

func (_ displayRow) FilterValue() string {
	return ""
}

func (_ renderedStr) Height() int  { return 1 }
func (_ renderedStr) Spacing() int { return 0 }

func (_ renderedStr) Render(w io.Writer, parent list.Model, idx int, item list.Item) {
	row := item.(displayRow)

	cursor := " "
	if idx == parent.Index() {
		cursor = ">"
	}

	fmt.Fprintf(w, "%s %s", cursor, row.display)
}

func (_ renderedStr) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

func makeEdgeItems(visited []VisitedEdge, defs []NamedEdge) []list.Item {
	rows := make([]list.Item, len(visited))
	for i, edge := range visited {
		def := defs[edge.Index]
		mark := "C"
		if !edge.Crossed {
			mark = "F"
		}
		row := displayRow{
			id:      def.Id,
			display: fmt.Sprintf("%s %s", mark, def.String()),
		}
		rows[i] = row
	}
	slices.SortFunc(rows, func(a, b list.Item) int {
		rowA, rowB := a.(displayRow), b.(displayRow)
		return -strings.Compare(rowA.display, rowB.display)
	})

	return rows
}

var edgeChunkStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).MarginBottom(2)
