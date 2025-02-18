package spheres

import (
	"fmt"
	"io"
	"strings"
	"sudonters/libzootr/components"
	"sudonters/libzootr/zecs"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/etc-sudonters/substrate/skelly/bitset32"
)

type Handles struct {
	Visited, Pending, Reached bitset32.Bitset

	Edges []VisitedEdge
}

type VisitedEdge struct {
	Index   int
	Crossed bool
}

func (this Handles) visited(yield func(zecs.Entity) bool) {
	iter := bitset32.IterT[zecs.Entity](&this.Visited)
	for entity := range iter.All {
		if !yield(entity) {
			return
		}
	}
}

func (this Handles) pending(id zecs.Entity) bool {
	return bitset32.IsSet(&this.Pending, id)
}

func (this Handles) reached(id zecs.Entity) bool {
	return bitset32.IsSet(&this.Reached, id)
}

type visitFlag rune

const (
	VISITED = ' '
	PENDING = 'P'
	REACHED = 'R'
)

const namedEdgeRowTemplate = "%06d %s"

type NamedEdge struct {
	Id   zecs.Entity
	Name components.Name
}

func (this NamedEdge) String() string {
	return fmt.Sprintf(namedEdgeRowTemplate, this.Id, this.Name)
}

type NamedNode struct {
	Name components.Name
	Id   zecs.Entity
}

const namedNodeRowTemplate = "%s %04d %s"

func (this NamedNode) asRow(flag rune) string {
	return fmt.Sprintf(
		namedNodeRowTemplate,
		string(flag), this.Id, this.Name,
	)
}

type NamedToken struct {
	Name components.Name
	Id   zecs.Entity
	Qty  int
}

type Details struct {
	I   int
	Err error

	Nodes  []NamedNode
	Tokens []NamedToken
	Edges  []NamedEdge

	Adult, Child Handles
}

func summarize(this Details) Summary {
	var s Summary
	s.I = this.I

	if this.Err != nil {
		errStr := this.Err.Error()
		if len(errStr) > 15 {
			errStr = fmt.Sprintf("%s...", errStr[:15])
		}
		s.Err = errStr
	}

	s.Visited = len(this.Nodes)
	s.Pending = this.Adult.Pending.Len() + this.Child.Pending.Len()
	s.Reached = this.Adult.Reached.Len() + this.Child.Reached.Len()
	s.Collected = len(this.Tokens)

	return s
}

type Summary struct {
	I   int
	Err string

	Collected, Pending, Reached, Visited int
}

func (this Summary) FilterValue() string {
	return fmt.Sprintf("%d", this.I)
}

type summaryDelegate struct{}

func (_ summaryDelegate) Render(w io.Writer, parent list.Model, idx int, item list.Item) {
	var view strings.Builder
	summary := item.(Summary)
	fmt.Fprintf(&view, "Sphere %3d\n", summary.I)
	fmt.Fprintf(&view, "V: %03d P: %03d R: %03d C: %03d\n", summary.Visited, summary.Pending, summary.Reached, summary.Collected)
	if summary.Err != "" {
		fmt.Fprintf(&view, "Err: %q", summary.Err)
	}
	style := inactiveStyle
	if summary.I == parent.Index() {
		style = currentStyle
	} else if summary.Err != "" {
		style = errStyle
	} else if idx == len(parent.Items())-1 {
		style = latestStyle
	}

	fmt.Fprint(w, style.Render(view.String()))
}

func (_ summaryDelegate) Height() int                         { return 3 }
func (_ summaryDelegate) Spacing() int                        { return 1 }
func (_ summaryDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

var (
	inactiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	currentStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	latestStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)
