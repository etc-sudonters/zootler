package ioutil

import "fmt"

type ConsolePainter interface {
	Paint(string) string
}

type BackgroundColor color
type ForegroundColor color

type color int
type canvas int
type Reset struct{}

var canvasFore canvas = 38
var canvasBack canvas = 48

func (r Reset) String() string {
	return "\033[0m"
}

func (f ForegroundColor) Paint(s string) string {
	return fmt.Sprintf("%s%s%s", paint(color(f), canvasFore), s, Reset{})
}

func (b BackgroundColor) Paint(s string) string {
	return fmt.Sprintf("%s%s%s", paint(color(b), canvasBack), s, Reset{})
}

func paint(clr color, cnv canvas) string {
	return fmt.Sprintf("\033[%d:5:%dm", cnv, clr)
}
