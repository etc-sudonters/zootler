package main

type WorldGraphLoader struct {
	Helpers, Path string
	IncludeMQ     bool
}

func (w WorldGraphLoader) Load() error { return nil }
