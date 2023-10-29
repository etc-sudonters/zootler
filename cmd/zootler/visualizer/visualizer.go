package visualizer

import (
	"context"
	"sudonters/zootler/internal/entity/bitpool"
	"sudonters/zootler/internal/reitertools"
	"sudonters/zootler/pkg/world"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/etc-sudonters/substrate/stageleft"
)

func Visualize(w world.World) visualizer {
	return visualizer{w}
}

type visualizer struct {
	w world.World
}

func (v visualizer) Run(ctx context.Context) error {
	var frame frame

	tbl, err := bitpool.ExtractComponentTable(v.w.Entities.Pool)
	if err != nil {
		return err
	}

	frame.keys = keymap{
		quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		change: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "change focus")),
	}

	frame.menu = menu{
		keys: menuKeys{
			up:     key.NewBinding(key.WithKeys("up", "j"), key.WithHelp("up/j", "move up (wraps)")),
			down:   key.NewBinding(key.WithKeys("down", "k"), key.WithHelp("down/k", "moves down (wraps)")),
			choose: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "press enter to do nothing")),
		},
		components: reitertools.ToSlice(tbl.Rows()),
	}

	_, err = tea.NewProgram(
		&frame,
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return stageleft.AttachExitCode(err, stageleft.ExitCode(97))
	}
	return nil
}

type frameFocus uint64

const (
	focusedOnMenu frameFocus = iota
	focusedOnPrimary
)

type frame struct {
	menu    tea.Model
	primary tea.Model
	focused frameFocus

	keys keymap
	help help.Model
}

func (f *frame) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (f *frame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keys.quit):
			cmds = append(cmds, tea.Quit)
		}
	case interface{ menuOnly() }:
		if c := f.updateMenu(msg); c != nil {
			cmds = append(cmds, c)
		}
		break
	case interface{ primaryOnly() }:
		if c := f.updatePrimary(msg); c != nil {
			cmds = append(cmds, c)
		}
		break
	}

	// msg was already handled, don't double dispatch
	if cmds == nil {
		switch f.focused {
		case focusedOnMenu:
			if c := f.updateMenu(msg); c != nil {
				cmds = append(cmds, c)
			}
			break
		case focusedOnPrimary:
			if c := f.updatePrimary(msg); c != nil {
				cmds = append(cmds, c)
			}
		}
	}

	var cmd tea.Cmd
	if len(cmds) != 0 {
		cmd = tea.Batch(cmds...)
	}

	return f, cmd
}

func (f *frame) updateMenu(msg tea.Msg) tea.Cmd {
	if f.menu == nil {
		return nil
	}
	m, c := f.menu.Update(msg)
	f.menu = m
	return c
}

func (f *frame) updatePrimary(msg tea.Msg) tea.Cmd {
	if f.primary == nil {
		return nil
	}
	m, c := f.primary.Update(msg)
	f.primary = m
	return c
}

func (f *frame) View() string {
	return f.menu.View()
}

type keymap struct {
	change, quit key.Binding
}
