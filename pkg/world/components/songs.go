package components

type (
	OcarinaButton struct{}
	OcarinaNote   rune

	Song struct {
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

func (c Song) String() string { return "Song" }
