package codegen_go

import (
	"dog/ast"
	gast "go/ast"
	"go/token"
	"strings"
)

//
//
// param: c
// return:
func (this *Translation) transClass(c ast.Class) (cl *gast.GenDecl) {
	this.CurrentClass = c
	if cc, ok := c.(*ast.ClassSingle); ok {
		cl = &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}

		sp := &gast.TypeSpec{
			Doc:     nil,
			Name:    gast.NewIdent(cc.Name),
			Assign:  0,
			Type:    nil,
			Comment: nil,
		}
		Type := &gast.StructType{
			Struct: 0,
			Fields: &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			},
			Incomplete: false,
		}
		sp.Type = Type

		for _, fi := range cc.Fields {
			//FIXME 是否排除static
			if !fi.IsStatic() && !cc.IsEnum() {
				gfi := this.transField(fi)
				Type.Fields.List = append(Type.Fields.List, gfi)
				this.constructFieldFunc(fi)
			}
		}
		for _, m := range cc.Methods {
			gmeth := this.transFunc(m)

			this.GolangFile.Decls = append(this.GolangFile.Decls, gmeth)
		}

		cl.Specs = append(cl.Specs, sp)

	}
	return
}

// 枚举转换
//
// param: c
func (this *Translation) transEnum(c ast.Class) {
	this.CurrentClass = c

	if cc, ok := c.(*ast.ClassSingle); ok {
		//1 定义枚举类型为int
		t := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		sp := &gast.TypeSpec{
			Doc:     nil,
			Name:    gast.NewIdent(cc.GetName()),
			Assign:  0,
			Type:    gast.NewIdent("int"),
			Comment: nil,
		}
		t.Specs = append(t.Specs, sp)
		this.GolangFile.Decls = append(this.GolangFile.Decls, t)

		//2 解析枚举元素
		v := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.CONST,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		this.GolangFile.Decls = append(this.GolangFile.Decls, v)
		for idx, fi := range cc.Fields {
			value := &gast.ValueSpec{
				Doc:     nil,
				Names:   []*gast.Ident{gast.NewIdent(strings.ToUpper(cc.GetName()) + "_" + fi.GetName())},
				Type:    nil,
				Values:  nil,
				Comment: nil,
			}
			if idx == 0 {
				value.Type = gast.NewIdent(cc.GetName())
				value.Values = append(value.Values, gast.NewIdent("iota"))
			}

			v.Specs = append(v.Specs, value)
		}

		//3 枚举元素值的String转换

		//处理类接收
		recv := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}

		gfi := &gast.Field{
			Doc:   nil,
			Names: []*gast.Ident{gast.NewIdent("this")},
			Type: &gast.Ident{
				NamePos: 0,
				Name:    this.CurrentClass.GetName(),
				Obj:     gast.NewObj(gast.Typ, this.CurrentClass.GetName()),
			},
			Tag:     nil,
			Comment: nil,
		}

		recv.List = append(recv.List, gfi)

		//处理返回值
		resType := &gast.Field{
			Doc:     nil,
			Names:   nil,
			Type:    gast.NewIdent("string"),
			Tag:     nil,
			Comment: nil,
		}

		results := &gast.FieldList{
			Opening: 0,
			List:    []*gast.Field{resType},
			Closing: 0,
		}

		//处理函数体

		swBlock := &gast.BlockStmt{
			Lbrace: 0,
			List:   nil,
			Rbrace: 0,
		}

		for _, fi := range cc.Fields {
			if sf, ok := fi.(*ast.FieldSingle); ok && sf.Value != nil {
				cause := &gast.CaseClause{
					Case:  0,
					List:  nil,
					Colon: 0,
					Body:  nil,
				}
				getStm := &gast.ReturnStmt{
					Return:  0,
					Results: []gast.Expr{this.transExp(sf.Value)},
				}

				cause.List = append(cause.List, gast.NewIdent(strings.ToUpper(cc.GetName())+"_"+sf.Name))
				cause.Body = append(cause.Body, getStm)

				swBlock.List = append(swBlock.List, cause)
			}

		}

		var getBody = &gast.BlockStmt{
			Lbrace: 0,
			List:   nil,
			Rbrace: 0,
		}

		swStm := &gast.SwitchStmt{
			Switch: 0,
			Init:   nil,
			Tag:    gast.NewIdent("this"),
			Body:   swBlock,
		}
		//switch代码块
		getBody.List = append(getBody.List, swStm)

		retStm := &gast.ReturnStmt{
			Return:  0,
			Results: []gast.Expr{gast.NewIdent("\"\"")},
		}

		//空return代码块
		getBody.List = append(getBody.List, retStm)

		stringFun := &gast.FuncDecl{
			Doc:  nil,
			Recv: recv,
			Name: gast.NewIdent("String"),
			Type: &gast.FuncType{
				Func:    0,
				Params:  nil,
				Results: results,
			},
			Body: getBody,
		}

		this.GolangFile.Decls = append(this.GolangFile.Decls, stringFun)

	}
}
