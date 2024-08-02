package components

type (
	OcarinaButton struct{}
	OcarinaNote   rune

	OcarinaSong struct {
		Notes []OcarinaNote
	}
)

const (
	OcarinaA OcarinaNote = 'A'
	OcarinaL             = '<'
	OcarinaR             = '>'
	OcarinaU             = '^'
	OcarinaD             = 'v'
)

func (c OcarinaSong) String() string { return "Song" }
