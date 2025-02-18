package leaves

import (
	"context"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
)

func NewApp(ctx context.Context, std *dontio.Std, mount tea.Model) App {
	if mount == nil {
		panic("cannot mount nil model")
	}

	ctx, cancel := context.WithCancelCause(ctx)

	return App{mounted: mount, std: std, ctx: ctx, cancelCause: cancel}
}

type App struct {
	maybeQuitting bool
	mounted       tea.Model
	std           *dontio.Std
	ctx           context.Context
	cancelCause   context.CancelCauseFunc
	Err           error
}

type StdOutMsg string
type StdErrMsg string
type Wrote struct {
	Err error
	N   int
}

func WriteToStdOutF(msg string, v ...any) tea.Cmd {
	return func() tea.Msg {
		return StdOutMsg(fmt.Sprintf(msg, v...))
	}
}

func WriteToStdErrF(msg string, v ...any) tea.Cmd {
	return func() tea.Msg {
		return StdErrMsg(fmt.Sprintf(msg, v...))
	}
}

func writeTo(w io.Writer, msg string) tea.Cmd {
	return func() tea.Msg {
		var wrote Wrote
		wrote.N, wrote.Err = fmt.Fprintln(w, msg)
		return wrote
	}
}

func (this App) Init() tea.Cmd {
	return this.mounted.Init()
}

func (this App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if ctxErr := this.ctx.Err(); ctxErr != nil {
		this.Err = ctxErr
		return this, tea.Batch(tea.Quit, WriteToStdErrF("context canceled: %s", ctxErr))
	}

	switch msg := msg.(type) {
	case StdErrMsg:
		return this, writeTo(this.std.Err, string(msg))
	case StdOutMsg:
		return this, writeTo(this.std.Out, string(msg))
	case Wrote:
		if msg.Err != nil {
			err := slipup.Describe(msg.Err, "failed to write")
			this.Err = err
			this.cancelCause(err)
			return this, tea.Batch(tea.Quit, cmd)
		}
		return this, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return this, tea.Batch(tea.Quit, WriteToStdErrF("ctrl+c interrupt, exiting immediately"))
		case tea.KeyRunes:
			str := string(msg.Runes)
			wasQuiting := this.maybeQuitting
			if wasQuiting && str == "!" {
				return this, tea.Batch(tea.Quit, WriteToStdErrF("graceful shutdown requested"))
			}

			this.maybeQuitting = "Q" == str
			if this.maybeQuitting {
				cmd = tea.Batch(cmd, WriteToStdErrF("maybe quiting..."))
			} else if wasQuiting {
				cmd = tea.Batch(cmd, WriteToStdErrF("not quitting"))
			}
		}
	}

	var mountedCmd tea.Cmd
	this.mounted, mountedCmd = this.mounted.Update(msg)
	return this, tea.Batch(cmd, mountedCmd)
}

func (this App) View() string {
	return this.mounted.View()
}
