package objects

import "fmt"

const (
	PtrToken   PtrTag = 0xCA
	PtrRegion  PtrTag = 0xFE
	PtrPlace   PtrTag = 0xDE
	PtrTrans   PtrTag = 0xAD
	PtrFunc    PtrTag = 0xBE
	PtrSetting PtrTag = 0xEF
)

func (this PtrTag) String() string {
	switch this {
	case PtrToken:
		return "PtrToken"
	case PtrRegion:
		return "PtrRegion"
	case PtrPlace:
		return "PtrPlace"
	case PtrTrans:
		return "PtrTrans"
	case PtrFunc:
		return "PtrFunc"
	case PtrSetting:
		return "PtrSetting"
	default:
		panic(fmt.Errorf("unknown ptr tag %X", uint8(this)))
	}
}
