package codegen_go

import (
	"dog/ast"
	gast "go/ast"
	"go/token"
)

//dao构造函数
func (this *Translation) getNewDaoFunc(c ast.Class) (fn *gast.FuncDecl) {
	var init gast.Stmt // 构造函数的初始化语句
	//处理函数参数
	params := &gast.FieldList{
		Opening: 0,
		List:    nil,
		Closing: 0,
	}
	params.List = append(params.List, this.getField(gast.NewIdent("db"), gast.NewIdent("*db.DB")))
	//处理返回值
	results := &gast.FieldList{
		Opening: 0,
		List:    nil,
		Closing: 0,
	}

	results.List = append(results.List, this.getField(gast.NewIdent("dao"), gast.NewIdent(c.GetName())))

	var body = &gast.BlockStmt{
		Lbrace: 0,
		List:   nil,
		Rbrace: 0,
	}

	//初始化语句
	val := &gast.UnaryExpr{
		OpPos: 0,
		Op:    token.AND,
		X: &gast.CompositeLit{
			Type:       gast.NewIdent(DeCapitalize(c.GetName())),
			Lbrace:     0,
			Elts:       []gast.Expr{gast.NewIdent("db")},
			Rbrace:     0,
			Incomplete: false,
		},
	}

	init = &gast.AssignStmt{
		Lhs:    []gast.Expr{gast.NewIdent("dao")},
		TokPos: 0,
		Tok:    token.ASSIGN,
		Rhs:    []gast.Expr{val},
	}

	body.List = append(body.List, init)

	retStm := &gast.ReturnStmt{
		Return:  0,
		Results: nil,
	}

	body.List = append(body.List, retStm)

	fn = &gast.FuncDecl{
		Doc:  nil,
		Recv: nil,
		Name: gast.NewIdent("New" + c.GetName()),
		Type: &gast.FuncType{
			Func:    0,
			Params:  params,
			Results: results,
		},
		Body: body,
	}
	return
}

func (this *Translation) getSaveDao(c ast.Class) (fn *gast.FuncDecl) {
	fn = &gast.FuncDecl{
		Recv: &gast.FieldList{
			List: []*gast.Field{
				&gast.Field{
					Names: []*gast.Ident{gast.NewIdent("this")},
					Type: &gast.StarExpr{
						X: gast.NewIdent(DeCapitalize(c.GetName()))},
				},
			},
		},
		Name: gast.NewIdent("Save"),
		Type: &gast.FuncType{

			Params: &gast.FieldList{
				List: []*gast.Field{
					&gast.Field{
						Names: []*gast.Ident{gast.NewIdent("m")},
						Type:  gast.NewIdent("interface{}"),
					},
				},
			},
			Results: &gast.FieldList{

				List: []*gast.Field{&gast.Field{
					Names: []*gast.Ident{gast.NewIdent("err")},
					Type:  gast.NewIdent("error"),
				},
				},
			},
		},
		Body: &gast.BlockStmt{

			List: []gast.Stmt{
				&gast.AssignStmt{
					Lhs: []gast.Expr{
						gast.NewIdent("err"),
					},
					Tok: token.ASSIGN,
					Rhs: []gast.Expr{&gast.SelectorExpr{
						X: &gast.CallExpr{
							Fun: &gast.SelectorExpr{
								X: &gast.CallExpr{
									Fun: &gast.SelectorExpr{
										X:   gast.NewIdent("this"),
										Sel: gast.NewIdent("DBWrite"),
									},
								},
								Sel: gast.NewIdent("Save"),
							},
							Args: []gast.Expr{gast.NewIdent("m")},
						},
						Sel: gast.NewIdent("Error"),
					},
					},
				},
				1: &gast.ReturnStmt{
					Return:  0,
					Results: nil,
				},
			},
		},
	}

	return
}

//func getExtraFun() {
//	f, err := parser.ParseFile(token.NewFileSet(), "", `package p; func f() { x, y, z := x, y, z }`, 0)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// RHS refers to undefined globals; LHS does not.
//	as := f.Decls[0].(*gast.FuncDecl).Body.List[0].(*gast.AssignStmt)
//	for _, v := range as.Rhs {
//		id := v.(*gast.Ident)
//		if id.Obj != nil {
//			t.Errorf("rhs %s has Obj, should not", id.Name)
//		}
//	}
//	for _, v := range as.Lhs {
//		id := v.(*gast.Ident)
//		if id.Obj == nil {
//			t.Errorf("lhs %s does not have Obj, should", id.Name)
//		}
//	}
//}
