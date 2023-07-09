package ioutils

import "io"

type CountingWriter struct {
	W io.Writer
	N int64
}

func (c *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = c.W.Write(p)
	c.N += int64(n)
	return
}

type ErrorCarryingWriter struct {
	W   io.Writer
	Err error
}

func (e *ErrorCarryingWriter) Write(p []byte) (n int, err error) {
	if e.Err != nil {
		return 0, e.Err
	}

	n, err = e.W.Write(p)
	if err != nil {
		e.Err = err
	}
	return
}
