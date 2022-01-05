package codegen_go

import (
	"dog/ast"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"go/token"
	"regexp"
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
			if !fi.IsStatic() && cc.GetType() == ast.CLASS_TYPE {
				if fi.GetName() == "SerialVersionUID" {
					continue
				}
				gfi := this.transField(fi)
				if _, ok := fi.GetDecType().(*ast.ClassType); ok {
					gfi.Type = &gast.StarExpr{X: gfi.Type}
				}

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
	//构建New函数
	if ConstructNewFunc {
		this.GolangFile.Decls = append(this.GolangFile.Decls, this.getNewService(c))

	}

	return
}

// 枚举转换
//
// param: c
func (this *Translation) transEnum(c ast.Class) {
	this.CurrentClass = c
	defer func() {
		this.CurrentClass = nil
	}()

	if OneFold {
		this.PkgName = c.GetName()
		this.GolangFile.Name.Name = c.GetName()
	}
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
				Names:   []*gast.Ident{gast.NewIdent(fi.GetName())},
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

				cause.List = append(cause.List, gast.NewIdent(sf.Name))
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

//
func (this *Translation) transInterface(c ast.Class) {
	this.CurrentClass = c
	defer func() {
		this.CurrentClass = nil
	}()

	it := &gast.GenDecl{
		Doc:    nil,
		TokPos: 0,
		Tok:    token.TYPE,
		Lparen: 0,
		Specs:  nil,
		Rparen: 0,
	}

	sp := &gast.TypeSpec{
		Doc:     nil,
		Name:    gast.NewIdent(c.GetName()),
		Assign:  0,
		Type:    nil,
		Comment: nil,
	}
	Type := &gast.InterfaceType{
		Interface: 0,
		Methods: &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		},
		Incomplete: false,
	}
	sp.Type = Type

	for _, m := range c.ListMethods() {
		gmeth := this.transFunc(m)
		if gmeth.Name.Name == "FindById" {
			continue
		}
		if gmeth.Type.Results != nil {
			gmeth.Type.Results.List = append(gmeth.Type.Results.List, this.getErrRet())
		} else {
			gmeth.Type.Results = &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			}
			gmeth.Type.Results.List = append(gmeth.Type.Results.List, this.getErrRet())
		}

		field := &gast.Field{
			Doc:     nil,
			Names:   []*gast.Ident{gmeth.Name},
			Type:    gmeth.Type,
			Tag:     nil,
			Comment: nil,
		}

		Type.Methods.List = append(Type.Methods.List, field)
	}

	//Save
	gmeth := this.getSaveDao(c)
	field := &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{gmeth.Name},
		Type:    gmeth.Type,
		Tag:     nil,
		Comment: nil,
	}
	Type.Methods.List = append(Type.Methods.List, field)

	//FindById
	gmeth = this.getFindByIdDao(c)
	field = &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{gmeth.Name},
		Type:    gmeth.Type,
		Tag:     nil,
		Comment: nil,
	}
	Type.Methods.List = append(Type.Methods.List, field)

	//DeleteById
	gmeth = this.getDeleteByIdDao(c)
	field = &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{gmeth.Name},
		Type:    gmeth.Type,
		Tag:     nil,
		Comment: nil,
	}
	Type.Methods.List = append(Type.Methods.List, field)

	//FindAll
	gmeth = this.getFindAllDao(c)
	field = &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{gmeth.Name},
		Type:    gmeth.Type,
		Tag:     nil,
		Comment: nil,
	}
	Type.Methods.List = append(Type.Methods.List, field)

	it.Specs = append(it.Specs, sp)

	this.GolangFile.Decls = append(this.GolangFile.Decls, it)

	this.buildDao(c)

}

// dao接口实现
//
// param: c
func (this *Translation) buildDao(c ast.Class) {
	cl := &gast.GenDecl{
		Doc:    nil,
		TokPos: 0,
		Tok:    token.TYPE,
		Lparen: 0,
		Specs:  nil,
		Rparen: 0,
	}

	sp := &gast.TypeSpec{
		Doc:     nil,
		Name:    gast.NewIdent(DeCapitalize(c.GetName())),
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

	Type.Fields.List = append(Type.Fields.List, this.getField(gast.NewIdent(""), gast.NewIdent("*db.DB")))

	cl.Specs = append(cl.Specs, sp)

	this.GolangFile.Decls = append(this.GolangFile.Decls, cl)

	this.GolangFile.Decls = append(this.GolangFile.Decls, this.getNewDaoFunc(c))

	//接口实现
	for _, m := range c.ListMethods() {
		gmeth := this.transFunc(m)

		if gmeth.Name.Name == "FindById" {
			continue
		}

		//添加jpa查询实现
		reg := regexp.MustCompile(`Find(All)?(\w+)?By(\w+)`)
		n := reg.FindStringSubmatch(gmeth.Name.Name)
		//Find
		var q = "err = this.DBRead()"
		if len(n) == 4 && n[1] == "" && n[2] == "" && gmeth.Type.Params.List != nil {
			ss := strings.Split(n[3], "And")
			for k, v := range ss {
				if strings.HasSuffix(v, "Not") {
					v = strings.Replace(v, "Not", "", 1)
					q += fmt.Sprintf(".Where(\"%s = ?\", %s)", SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				} else {
					q += fmt.Sprintf(".Where(\"%s = ?\", %s)", SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				}
			}
			q += ".First(&result).Error"
			gmeth.Body.List = nil
			gmeth.Body.List = append(gmeth.Body.List, &gast.ExprStmt{gast.NewIdent(q)})
			ret := &gast.ReturnStmt{
				Return:  0,
				Results: nil,
			}
			gmeth.Body.List = append(gmeth.Body.List, ret)
			//FindAll
		} else if len(n) == 4 && n[1] == "All" && gmeth.Type.Params.List != nil {
			ss := strings.Split(n[3], "And")
			for k, v := range ss {
				if strings.HasSuffix(v, "Not") {
					v = strings.Replace(v, "Not", "", 1)
					q += fmt.Sprintf(".Where(\"%s <> ?\",%s)", SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				} else {
					q += fmt.Sprintf(".Where(\"%s = ?\",%s)", SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				}
			}
			q += ".Find(&result).Error"
			gmeth.Body.List = nil
			gmeth.Body.List = append(gmeth.Body.List, &gast.ExprStmt{gast.NewIdent(q)})
			ret := &gast.ReturnStmt{
				Return:  0,
				Results: nil,
			}
			gmeth.Body.List = append(gmeth.Body.List, ret)

		} else if len(n) == 4 && n[1] != "All" && gmeth.Type.Params.List != nil {
			ss := strings.Split(n[3], "And")
			for k, v := range ss {
				if strings.HasSuffix(v, "Not") {
					v = strings.Replace(v, "Not", "", 1)
					q += fmt.Sprintf(".Where(\"%s <> ?\",%s)", SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				} else {
					q += fmt.Sprintf(".Where(\"%s = ?\",%s)", SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				}
			}
			q += ".First(&result).Error"
			gmeth.Body.List = nil
			gmeth.Body.List = append(gmeth.Body.List, &gast.ExprStmt{gast.NewIdent(q)})
			ret := &gast.ReturnStmt{
				Return:  0,
				Results: nil,
			}
			gmeth.Body.List = append(gmeth.Body.List, ret)
		}

		log.Warn(n, gmeth.Name.Name, len(n))
		for _, v := range gmeth.Recv.List {
			//实现接口的struct名字小写开口
			v.Type = &gast.StarExpr{X: gast.NewIdent(DeCapitalize(c.GetName()))}
		}

		//每个函数末尾加err 返回
		if gmeth.Type.Results != nil {
			gmeth.Type.Results.List = append(gmeth.Type.Results.List, this.getErrRet())
		} else {
			gmeth.Type.Results = &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			}
			gmeth.Type.Results.List = append(gmeth.Type.Results.List, this.getErrRet())
		}

		this.GolangFile.Decls = append(this.GolangFile.Decls, gmeth)
	}
	//TODO 增加Save,FindAll,FindById等接口
	this.GolangFile.Decls = append(this.GolangFile.Decls, this.getSaveDao(c))
	this.GolangFile.Decls = append(this.GolangFile.Decls, this.getFindByIdDao(c))
	this.GolangFile.Decls = append(this.GolangFile.Decls, this.getDeleteByIdDao(c))
	this.GolangFile.Decls = append(this.GolangFile.Decls, this.getFindAllDao(c))

}

// 完善查询函数 FindAllxxx,Findxxx
//
// param: fn
func (this *Translation) OptimizeDaoFun(fn *gast.FuncDecl) {
	//
	//regexp.Regexp{}
	//
	//var act = string
	//if strings.HasPrefix(fn.Name.Name, "FindAll") {
	//	act = "Find"
	//} else if strings.HasPrefix(fn.Name.Name, "Find") {
	//	act = "First"
	//}

}

func (this *Translation) getField(name *gast.Ident, tp *gast.Ident) (gfi *gast.Field) {
	gfi = &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{name},
		Type:    tp,
		Tag:     nil,
		Comment: nil,
	}
	return
}
