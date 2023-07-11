package errs

import (
	"strings"
)

func WriteLn(msg string) Fragment {
	if !strings.HasSuffix(msg, "\n") {
		msg = msg + "\n"
	}
	return Write(msg)
}

func Write(msg string) Fragment {
	return func(s *strings.Builder) {
		s.WriteString(msg)
	}
}

type Fragment func(*strings.Builder)

type Builder []Fragment

func (b Builder) Append(frags ...Fragment) Builder {
	return append(b, frags...)
}

func (b Builder) Error() string {
	msg := &strings.Builder{}
	for i := range b {
		b[i](msg)
	}

	return msg.String()
}

func WriteBefore(msg string, frag Fragment) Fragment {
	b := []byte(msg)
	return func(w *strings.Builder) {
		w.Write(b)
		frag(w)
	}
}

func WriteAfter(msg string, frag Fragment) Fragment {
	b := []byte(msg)
	return func(w *strings.Builder) {
		frag(w)
		w.Write(b)
	}
}
