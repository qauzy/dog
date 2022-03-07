package codegen_go

import (
	gast "go/ast"
)

type FakeBlock struct {
	gast.BlockStmt
	List []gast.Stmt
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
