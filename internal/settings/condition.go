package settings

// this is actually two uint8
// the upper bits describe what the condition is
// the lower bits describe how much of that condition
type quantitycondition struct {
	q uint16
}

func (q quantitycondition) Decode() (which Condition, qty uint8) {
	which = Condition((q.q & 0xFF00) >> 8)
	qty = uint8(0x00FF & q.q)
	return
}

func encodeqty(which Condition, qty uint8) quantitycondition {
	return quantitycondition{(uint16(which) << 8) | uint16(qty)}
}

type LacsCondition quantitycondition
type BridgeCondition quantitycondition
type GanonBKCondition quantitycondition

func CreateLacs(cond Condition, qty uint8) LacsCondition {
	return LacsCondition(encodeqty(cond, qty))
}

func CreateBridge(cond Condition, qty uint8) BridgeCondition {
	return BridgeCondition(encodeqty(cond, qty))
}

func CreateGanonBK(cond Condition, qty uint8) GanonBKCondition {
	return GanonBKCondition(encodeqty(cond, qty))
}

type Condition uint8

const (
	CondMedallions Condition = iota + 1
	CondStones
	CondRewards
	CondTokens
	CondHearts
	CondDefault
)
