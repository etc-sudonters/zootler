package explore

import (
	"fmt"
	"strings"
	"sudonters/libzootr/cmd/knowitall/leaves"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func newEditor() editor {
	var e editor
	e.text = textarea.New()
	return e
}

type editor struct {
	text textarea.Model

	writing bool
}

func (_ editor) Init() tea.Cmd { return nil }

func (this editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.text.SetHeight((msg.Height - 5) / 2)
		this.text.SetWidth(msg.Width - 5)
	case EditRule:
		if msg.Err != nil {
			return this, leaves.WriteStatusMsg("edit rule failed: %s", msg.Err.Error())
		}
		this.text.SetValue(string(msg.Source))
		cmd := this.insertMode()
		return this, cmd
	case tea.KeyMsg:
		if this.writing && msg.Type == tea.KeyEsc {
			cmd := this.normalMode()
			return this, cmd
		} else if this.writing {
			var cmd tea.Cmd
			this.text, cmd = this.text.Update(msg)
			return this, cmd
		} else if msg.String() == "i" {
			this.writing = true
			cmd := this.insertMode()
			return this, cmd
		}
	}

	return this, nil
}

func (this editor) View() string {
	var view strings.Builder
	fmt.Fprintln(&view, "EDITOR")
	view.WriteString("\n\n")
	fmt.Fprint(&view, this.text.View())
	return view.String()
}

func (this *editor) insertMode() tea.Cmd {
	this.writing = true
	cmd := this.text.Focus()
	return tea.Batch(cmd, leaves.WriteStatusMsg("insert mode"))
}

func (this *editor) normalMode() tea.Cmd {
	this.writing = false
	return leaves.WriteStatusMsg("normal mode")
}
