package codegen_go

import (
	"dog/ast"
	"dog/cfg"
	"dog/util"
	gast "go/ast"
	"go/token"
)

func (this *Translation) getClassTpeExpr() gast.Expr {
	if this.currentClass == nil {
		return nil
	}

	var clsType gast.Expr = &gast.Ident{
		NamePos: 0,
		Name:    this.currentClass.GetName(),
		Obj:     gast.NewObj(gast.Typ, this.currentClass.GetName()),
	}

	//处理接收对象的泛型
	if this.currentClass.ListGeneric() != nil {
		var indexListExpr = &gast.IndexListExpr{
			X:       clsType,
			Lbrack:  0,
			Indices: nil,
			Rbrack:  0,
		}
		for _, vv := range this.currentClass.ListGeneric() {
			indexListExpr.Indices = append(indexListExpr.Indices, gast.NewIdent(vv.Name))
		}
		clsType = indexListExpr
	}

	return clsType
}

func (this *Translation) constructFieldFunc(gfi *gast.Field) {
	if !cfg.ConstructFieldFunc {
		return
	}
	var clsType = this.getClassTpeExpr()
	if clsType == nil {
		return
	}

	var recv *gast.FieldList
	//处理类接收
	recv = &gast.FieldList{
		Opening: 0,
		List:    nil,
		Closing: 0,
	}

	recvFi := &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{gast.NewIdent("this")},
		Type:    &gast.StarExpr{X: clsType},
		Tag:     nil,
		Comment: nil,
	}

	recv.List = append(recv.List, recvFi)

	//member values setter function
	params := &gast.FieldList{
		Opening: 0,
		List:    nil,
		Closing: 0,
	}
	var paramFi = *gfi
	params.List = append(params.List, &paramFi)
	var setBody = &gast.BlockStmt{
		Lbrace: 0,
		List:   nil,
		Rbrace: 0,
	}

	setStm := &gast.AssignStmt{
		Lhs:    []gast.Expr{gast.NewIdent("this." + util.Capitalize(gfi.Names[0].Name))},
		TokPos: 0,
		Tok:    token.ASSIGN,
		Rhs:    []gast.Expr{paramFi.Names[0]},
	}

	setBody.List = append(setBody.List, setStm)

	setRetStm := &gast.ReturnStmt{
		Return:  0,
		Results: []gast.Expr{gast.NewIdent("this")},
	}
	setBody.List = append(setBody.List, setRetStm)

	//return value
	setResult := &gast.FieldList{
		Opening: 0,
		List:    nil,
		Closing: 0,
	}

	ret := &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{gast.NewIdent("result")},
		Type:    &gast.StarExpr{X: clsType},
		Tag:     nil,
		Comment: nil,
	}

	setResult.List = append(setResult.List, ret)

	//function definition
	setFun := &gast.FuncDecl{
		Doc:  nil,
		Recv: recv,
		Name: gast.NewIdent("Set" + util.Capitalize(gfi.Names[0].Name)),
		Type: &gast.FuncType{
			Func:    0,
			Params:  params,
			Results: setResult,
		},
		Body: setBody,
	}

	this.GolangFile.Decls = append(this.GolangFile.Decls, setFun)

	//成员值获取函数
	var getBody = &gast.BlockStmt{
		Lbrace: 0,
		List:   nil,
		Rbrace: 0,
	}

	getStm := &gast.ReturnStmt{
		Return:  0,
		Results: []gast.Expr{gast.NewIdent("this." + util.Capitalize(gfi.Names[0].Name))},
	}

	getBody.List = append(getBody.List, getStm)

	//处理返回值
	results := &gast.FieldList{
		Opening: 0,
		List:    []*gast.Field{gfi},
		Closing: 0,
	}

	getFun := &gast.FuncDecl{
		Doc:  nil,
		Recv: recv,
		Name: gast.NewIdent("Get" + util.Capitalize(gfi.Names[0].Name)),
		Type: &gast.FuncType{
			Func:    0,
			Params:  nil,
			Results: results,
		},
		Body: getBody,
	}

	this.GolangFile.Decls = append(this.GolangFile.Decls, getFun)

	return
}

// 翻译函数抽象语法树
//
// param: fi
// return:
func (this *Translation) transFunc(fi ast.Method) (fn *gast.FuncDecl) {
	this.currentMethod = fi
	this.methodStack.Push(fi)
	this.Push(fi)
	defer func() {
		this.Pop()
		this.methodStack.Pop()
		if this.methodStack.Peek() != nil {
			this.currentMethod = this.methodStack.Peek().(ast.Method)
		} else {
			this.currentMethod = nil
		}
	}()
	if method, ok := fi.(*ast.MethodSingle); ok {
		var recv *gast.FieldList

		var init gast.Stmt // 构造函数的初始化语句
		//处理函数参数
		params := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}
		//添加gin
		if cfg.AppendContext {
			params.List = append(params.List, this.getField(gast.NewIdent("ctx"), gast.NewIdent("*gin.Context")))
		}

		for _, p := range method.Formals {
			pa := this.transFormals(p)

			params.List = append(params.List, pa)
		}
		//处理返回值
		results := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}

		var body = &gast.BlockStmt{
			Lbrace: 0,
			List:   nil,
			Rbrace: 0,
		}
		//代表类的描述信息
		var tp = this.getClassTpeExpr()

		//处理类接收
		if !fi.IsConstruct() && !fi.IsStatic() {
			//处理类接收
			recv = &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			}

			gfi := &gast.Field{
				Doc:     nil,
				Names:   []*gast.Ident{gast.NewIdent("this")},
				Type:    &gast.StarExpr{X: tp},
				Tag:     nil,
				Comment: nil,
			}
			recv.List = append(recv.List, gfi)
		}

		//处理类接收
		if !fi.IsConstruct() {

			//处理返回值
			ret := &gast.Field{
				Doc:     nil,
				Names:   []*gast.Ident{gast.NewIdent("result")},
				Type:    this.transType(method.RetType),
				Tag:     nil,
				Comment: nil,
			}

			//if type is SelectorExpr,add *
			ret.Type = this.checkStar(ret.Type)

			//if return type is void,then no return
			if ret.Type != nil {
				results.List = append(results.List, ret)
			}

			if this.currentMethod.IsThrows() || cfg.AddErrReruen {
				results.List = append(results.List, this.getErrRet())
			}

			if results.List == nil {
				results = nil
			}

		} else {
			call := &gast.CallExpr{
				Fun:      gast.NewIdent("new"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}

			call.Args = append(call.Args, tp)
			init = &gast.AssignStmt{
				Lhs:    []gast.Expr{gast.NewIdent("this")},
				TokPos: 0,
				Tok:    token.ASSIGN,
				Rhs:    []gast.Expr{call},
			}

			body.List = append(body.List, init)

			//处理返回值
			ret := &gast.Field{
				Doc:     nil,
				Names:   []*gast.Ident{gast.NewIdent("this")},
				Type:    &gast.StarExpr{X: tp},
				Tag:     nil,
				Comment: nil,
			}
			results.List = append(results.List, ret)

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
					//为了兼容处理要翻译为多个stm的语句
					if blk, ok := ss.(*FakeBlock); ok {
						for _, st := range blk.List {
							body.List = append(body.List, st)
						}
					} else {
						body.List = append(body.List, ss)
					}
				}

				if stm.GetExtra() != nil {
					stms := this.transStm(stm.GetExtra())
					//为了兼容处理要翻译为多个stm的语句
					if blk, ok := stms.(*FakeBlock); ok {
						for _, st := range blk.List {
							body.List = append(body.List, st)
						}
					} else {
						body.List = append(body.List, stms)
					}
				}
			}
		}
		//函数体为空
		if body.List == nil {
			body = nil
		}

		//处理构造函数
		if fi.IsConstruct() {
			ret := &gast.ReturnStmt{
				Return:  0,
				Results: nil,
			}
			body.List = append(body.List, ret)
		}
		cm := &gast.CommentGroup{[]*gast.Comment{{
			Slash: 0,
			Text:  method.Comment,
		}}}

		if method.Comment == "" {
			cm = nil
		}
		if cfg.DropResult {
			results = nil
		}

		fn = &gast.FuncDecl{
			Doc:  cm,
			Recv: recv,
			Name: gast.NewIdent(util.Capitalize(method.Name.Name)),
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
func (this *Translation) transFuncLit(v *ast.Question) (call *gast.CallExpr) {

	fn := &gast.FuncLit{
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
	call = &gast.CallExpr{
		Fun:      fn,
		Lparen:   0,
		Args:     nil,
		Ellipsis: 0,
		Rparen:   0,
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
				stms := this.transStm(stm)
				if stms != nil {
					//为了兼容处理要翻译为多个stm的语句
					if blk, ok := stms.(*FakeBlock); ok {
						for _, st := range blk.List {
							body.List = append(body.List, st)
						}
					} else {
						body.List = append(body.List, stms)
					}
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
