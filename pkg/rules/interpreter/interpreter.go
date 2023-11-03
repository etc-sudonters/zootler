package interpreter

import "sudonters/zootler/pkg/rules/parser"

type TreeWalk struct {
}

type treewalker struct{}

func (t *treewalker) VisitAttrAccess(node *parser.AttrAccess) {
	panic("not implemented")
}
func (t *treewalker) VisitBinOp(node *parser.BinOp) {
	panic("not implemented")
}
func (t *treewalker) VisitBoolOp(node *parser.BoolOp) {
	panic("not implemented")
}
func (t *treewalker) VisitBoolean(node *parser.Boolean) {
	panic("not implemented")
}
func (t *treewalker) VisitCall(node *parser.Call) {
	panic("not implemented")
}
func (t *treewalker) VisitIdentifier(node *parser.Identifier) {
	panic("not implemented")
}
func (t *treewalker) VisitNumber(node *parser.Number) {
	panic("not implemented")
}
func (t *treewalker) VisitString(node *parser.String) {
	panic("not implemented")
}
func (t *treewalker) VisitSubscript(node *parser.Subscript) {
	panic("not implemented")
}
func (t *treewalker) VisitTuple(node *parser.Tuple) {
	panic("not implemented")
}
func (t *treewalker) VisitUnary(node *parser.UnaryOp) {
	panic("not implemented")
}
