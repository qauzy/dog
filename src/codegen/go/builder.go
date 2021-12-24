package codegen_go

import (
	"fmt"
	gast "go/ast"
	"go/parser"
	"go/token"
)

func (this *Translation) getExpr(src string) (expr gast.Expr) {
	src = fmt.Sprintf(`package p; func f() { %s }`, src)

	f, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		this.TranslationBug(err)
	}

	// RHS refers to undefined globals; LHS does not.
	exprStmt := f.Decls[0].(*gast.FuncDecl).Body.List[0].(*gast.ExprStmt)
	expr = exprStmt.X
	return
}

func (this *Translation) getFunc(src string) (fn *gast.FuncDecl) {
	src = fmt.Sprintf(`package p; %v`, src)

	f, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		this.TranslationBug(err)
	}
	fn = f.Decls[0].(*gast.FuncDecl)
	return
}
