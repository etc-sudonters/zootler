package astrender

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/etc-sudonters/substrate/dontio"
)

func DontTheme() ColorScheme {
	return ColorScheme{
		Node:       dontPaint{},
		Property:   dontPaint{},
		String:     dontPaint{},
		Number:     dontPaint{},
		Boolean:    dontPaint{},
		Identifier: dontPaint{},
		Keyword:    dontPaint{},
		Function:   dontPaint{},
		Brackets:   nil,
	}

}

type dontPaint struct{}

func (d dontPaint) Paint(s string) string { return s }

type ColorScheme struct {
	Node       dontio.Painter
	Property   dontio.Painter
	String     dontio.Painter
	Number     dontio.Painter
	Boolean    dontio.Painter
	Identifier dontio.Painter
	Keyword    dontio.Painter
	Function   dontio.Painter
	Brackets   []dontio.Painter
}

func (c ColorScheme) BracketFor(i int) dontio.Painter {
	numColors := len(c.Brackets)
	if numColors == 0 {
		return dontPaint{}
	}

	return c.Brackets[i%numColors]
}

const (
	nodeColor    dontio.ForegroundColor = 244
	propColor    dontio.ForegroundColor = 252
	strColor     dontio.ForegroundColor = 112
	numColor     dontio.ForegroundColor = 33
	boolColor    dontio.ForegroundColor = 160
	identColor   dontio.ForegroundColor = 99
	keywordColor dontio.ForegroundColor = 208
	fnColor      dontio.ForegroundColor = 228
)

func DefaultColorScheme() ColorScheme {
	return ColorScheme{
		Node:       nodeColor,
		Property:   propColor,
		String:     strColor,
		Number:     numColor,
		Boolean:    boolColor,
		Identifier: identColor,
		Keyword:    keywordColor,
		Function:   fnColor,
		Brackets: []dontio.Painter{
			dontio.ForegroundColor(1),
			dontio.ForegroundColor(2),
			dontio.ForegroundColor(3),
			dontio.ForegroundColor(4),
			dontio.ForegroundColor(5),
			dontio.ForegroundColor(6),
			dontio.ForegroundColor(7),
		},
	}

}

const (
	nodeLipgloss    lipgloss.ANSIColor = 244
	propLipgloss    lipgloss.ANSIColor = 252
	strLipgloss     lipgloss.ANSIColor = 112
	numLipgloss     lipgloss.ANSIColor = 33
	boolLipgloss    lipgloss.ANSIColor = 160
	identLipgloss   lipgloss.ANSIColor = 99
	keywordLipgloss lipgloss.ANSIColor = 208
	fnLipgloss      lipgloss.ANSIColor = 228
)

func LipglossColorScheme() ColorScheme {
	return ColorScheme{
		Node:       lipglossPainter(lipgloss.NewStyle().Foreground(nodeLipgloss)),
		Property:   lipglossPainter(lipgloss.NewStyle().Foreground(propLipgloss)),
		String:     lipglossPainter(lipgloss.NewStyle().Foreground(strLipgloss)),
		Number:     lipglossPainter(lipgloss.NewStyle().Foreground(numLipgloss)),
		Boolean:    lipglossPainter(lipgloss.NewStyle().Foreground(boolLipgloss)),
		Identifier: lipglossPainter(lipgloss.NewStyle().Foreground(identLipgloss)),
		Keyword:    lipglossPainter(lipgloss.NewStyle().Foreground(keywordLipgloss)),
		Function:   lipglossPainter(lipgloss.NewStyle().Foreground(fnLipgloss)),
		Brackets:   lipglossBrackets,
	}
}

type lipglossPainter lipgloss.Style

func (l lipglossPainter) Paint(s string) string {
	return lipgloss.Style(l).Render(s)
}

var lipglossBrackets = []dontio.Painter{
	lipglossPainter(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1))),
	lipglossPainter(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2))),
	lipglossPainter(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(3))),
	lipglossPainter(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(4))),
	lipglossPainter(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(5))),
	lipglossPainter(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(6))),
	lipglossPainter(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(7))),
}
