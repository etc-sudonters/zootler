package symfile

type SymbolFile struct {
	LittleEndian bool

	Body    []Section
	Hash    [4]uint64
	Headers []SectionHeader
	Version uint64
}

type SectionHeader struct {
	NamePtr uint32
	Name    string
	Index   int
	Address uint64
	Size    uint64
	Type    SectionType
}

type Section interface {
	Type() SectionType
	Bytes() []uint8
}

type SectionString struct {
	raw     []uint8
	strings map[uint32]string
}

type SectionRuleDefinitions struct {
	raw   []uint8
	rules []RuleDefinition
}

type SectionRuleProgram struct {
	raw []uint32
}

type SectionSymbol struct {
	raw     []uint8
	symbols []Symbol
}

type SectionType uint8
type SymbolType uint8

type RuleDefinition struct {
	Origin, Dest Symbol
	Address      uint32
	Size         uint16
}

type Symbol struct {
	Type  SymbolType
	Size  uint16
	Name  uint32
	Value uint64
}

const (
	_ SectionType = iota
	SECT_STR
	SECT_RULE_DEF
	SECT_RULES
	SECT_SYMBOLS
)
