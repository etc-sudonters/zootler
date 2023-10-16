package errs

import (
	"io"
	"runtime/debug"
	"strings"
)

func WritePanicTrace(w io.Writer) {
	w.Write(debug.Stack())
}

func ShowPanicTrace() string {
	b := &strings.Builder{}
	WritePanicTrace(b)
	return b.String()
}
