package codegen_go

import (
	"go/ast"
)

type FakeBlock struct {
	ast.BlockStmt
	List []ast.Stmt
}

type MethodReference struct {
	X ast.Expr // left operand
	Y ast.Expr // right operand
	ast.BinaryExpr
}

//func (FakeBlock) stmtNode() {
//
//}
//
//func (*FakeBlock) Pos() (p token.Pos) {
//	return
//}
//
//func (*FakeBlock) End() (p token.Pos) {
//	return
//}
//
//func (this *FakeBlock) accept(v gast.Visitor) {
//}
//func (this *FakeBlock) _stm() {
//}
