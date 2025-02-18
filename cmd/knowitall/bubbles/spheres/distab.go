package spheres

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newDisTab(parent *breakdown) disassemblyTab {
	return disassemblyTab{height: parent.viewHeight, width: parent.viewWidth}
}

type disassemblyTab struct {
	disassembled  Disassembly
	height, width int
}

func (this *disassemblyTab) update(msg tea.Msg, parent *breakdown) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.height = msg.Height
		this.width = msg.Width
	case Disassembly:
		this.disassembled = msg
	}
	return nil
}

func (_ *disassemblyTab) init(*breakdown) tea.Cmd { return nil }

func (this disassemblyTab) view(w io.Writer) {
	fmt.Fprintln(w, "Disassembly")
	if this.disassembled.Id == 0 {
		return
	}

	fmt.Fprintln(w, this.disassembled.Name, this.disassembled.Id)
	if this.disassembled.Err != nil {
		fmt.Fprintln(w, renderF(errStyle, "Error: %s", this.disassembled.Err))
	}

	consts := make([]string, len(this.disassembled.Dis.Constants))
	for i := range consts {
		consts[i] = this.disassembled.Dis.Constants[i].String()
	}

	nodeChunkStyle := edgeChunkStyle.Height(this.height - 10).Width((this.width - 8) / 2).Padding(1)
	constCol := nodeChunkStyle.Render(lipgloss.JoinVertical(lipgloss.Left, consts...))
	codeCol := nodeChunkStyle.Render(this.disassembled.Dis.Dis)
	display := lipgloss.JoinHorizontal(lipgloss.Top, codeCol, constCol)
	fmt.Fprint(w, display)
}

func renderF(style lipgloss.Style, tpl string, v ...any) string {
	return style.Render(fmt.Sprintf(tpl, v...))
}
