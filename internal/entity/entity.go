package entity

import (
	"fmt"
)

// a member of a pool's population
type Model uint64

const INVALID_ENTITY Model = 0

func (m Model) String() string {
	return fmt.Sprintf("Model{%d}", m)
}
