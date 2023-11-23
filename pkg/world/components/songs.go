package components

type (
	OcarinaButton rune

	Song struct {
		Notes []OcarinaButton
	}
)

const (
	OcarinaA OcarinaButton = 'A'
	OcarinaL               = '<'
	OcarinaR               = '>'
	OcarinaU               = '^'
	OcarinaD               = 'v'
)

func (c Song) String() string { return "Song" }
