package codegen_go

import (
	"dog/ast"
	gast "go/ast"
)

// 翻译函数抽象语法树
//
// param: fi
// return:
func (this *Translation) transFunc(fi ast.Method) (fn *gast.FuncDecl) {
	if method, ok := fi.(*ast.MethodSingle); ok {
		//处理函数参数
		params := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}
		for _, p := range method.Formals {
			params.List = append(params.List, this.transField(p))
		}
		//处理返回值
		results := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}
		rel := &ast.FieldSingle{
			Access:  0,
			Tp:      method.RetType,
			Name:    "result",
			IsField: false,
			Stms:    nil,
		}
		//如果是void则没有返回值
		ret := this.transField(rel)
		if ret.Type != nil {
			results.List = append(results.List, ret)
		} else {
			results = nil
		}
		var body = &gast.BlockStmt{
			Lbrace: 0,
			List:   nil,
			Rbrace: 0,
		}
		//处理函数体
		for _, stm := range method.Stms {

			if stm.IsTriple() {
				sss := this.transTriple(stm)
				if sss != nil && len(sss) > 0 {
					body.List = append(body.List, sss...)
				}
			} else {
				ss := this.transStm(stm)

				if ss != nil {
					body.List = append(body.List, ss)
				}
			}

		}
		//函数体为空
		if body.List == nil {
			body = nil
		}

		fn = &gast.FuncDecl{
			Doc:  nil,
			Recv: nil,
			Name: gast.NewIdent(Capitalize(method.Name)),
			Type: &gast.FuncType{
				Func:    0,
				Params:  params,
				Results: results,
			},
			Body: body,
		}

	}
	return
}

// 翻译抽象语法树
//
// param: fi
// return:
func (this *Translation) transFuncLit(v *ast.Question) (fn *gast.FuncLit) {

	fn = &gast.FuncLit{
		Type: &gast.FuncType{
			Func:    0,
			Params:  nil,
			Results: nil,
		},
		Body: nil,
	}

	one := this.transExp(v.One)
	two := this.transExp(v.Two)

	resultOne := &gast.ReturnStmt{
		Return:  0,
		Results: []gast.Expr{one},
	}

	resultTwo := &gast.ReturnStmt{
		Return:  0,
		Results: []gast.Expr{two},
	}

	stmt := &gast.IfStmt{
		If:   0,
		Init: nil,
		Cond: this.transExp(v.E),
		Body: &gast.BlockStmt{
			Lbrace: 0,
			List:   []gast.Stmt{resultOne},
			Rbrace: 0,
		},
		Else: &gast.BlockStmt{
			Lbrace: 0,
			List:   []gast.Stmt{resultTwo},
			Rbrace: 0,
		},
	}

	fn.Body = &gast.BlockStmt{
		Lbrace: 0,
		List:   []gast.Stmt{stmt},
		Rbrace: 0,
	}
	return
}

func (this *Translation) transLambda(fi ast.Exp) (fn *gast.FuncLit) {
	if lam, ok := fi.(*ast.Lambda); ok {
		//处理函数参数
		params := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}
		for _, p := range lam.Formals {
			params.List = append(params.List, this.transField(p))
		}
		//处理返回值
		results := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}
		//rel := &ast.FieldSingle{
		//	Access:  0,
		//	Tp:      method.RetType,
		//	Name:    "result",
		//	IsField: false,
		//	Stms:    nil,
		//}
		////如果是void则没有返回值
		//ret := transField(rel)
		//if ret.Type != nil {
		//	results.List = append(results.List, ret)
		//} else {
		//	results = nil
		//}
		var body = &gast.BlockStmt{
			Lbrace: 0,
			List:   nil,
			Rbrace: 0,
		}
		//处理函数体
		for _, stm := range lam.Stms {

			if stm.IsTriple() {
				sss := this.transTriple(stm)
				if sss != nil && len(sss) > 0 {
					body.List = append(body.List, sss...)
				}
			} else {
				ss := this.transStm(stm)

				if ss != nil {
					body.List = append(body.List, ss)
				}
			}

		}
		//函数体为空
		if body.List == nil {
			body = nil
		}

		fn = &gast.FuncLit{
			Type: &gast.FuncType{
				Func:    0,
				Params:  params,
				Results: results,
			},
			Body: body,
		}

	}
	return
}
