package codegen_go

import (
	"dog/ast"
	"dog/cfg"
	"dog/util"
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
	log.Infof("解析类:%v", c.GetName())
	if cfg.OneFold {
		this.PkgName = strings.Replace(c.GetName(), "Controller", "", -1)
		this.GolangFile.Name.Name = strings.Replace(c.GetName(), "Controller", "", -1)
	}

	this.currentClass = c
	this.Push(c)
	this.classStack.Push(c)
	defer func() {
		this.classStack.Pop()
		this.Pop()
		if this.classStack.Peek() != nil {
			this.currentClass = this.classStack.Peek().(ast.Class)
		} else {
			this.currentClass = nil
		}
	}()

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
		//处理类泛型
		var tpList = &gast.FieldList{}
		for _, generic := range cc.Generics {
			log.Debugf("解析泛型：%v", generic)
			var tp = gast.NewIdent("any")
			if generic.Extends != "" {
				tp = gast.NewIdent(generic.Extends)
			}

			fi := &gast.Field{
				Doc:     nil,
				Names:   []*gast.Ident{gast.NewIdent(generic.Name)},
				Type:    tp,
				Tag:     nil,
				Comment: nil,
			}

			tpList.List = append(tpList.List, fi)
			sp.TypeParams = tpList
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

			_, ok = fi.(*ast.FieldEnum)
			if !fi.IsStatic() && (cc.GetType() == ast.CLASS_TYPE || (cc.GetType() == ast.ENUM_TYPE && !ok)) {
				if fi.GetName() == "SerialVersionUID" {
					continue
				}

				gfi := this.transField(fi)
				gfi.Type = this.checkStar(gfi.Type)

				Type.Fields.List = append(Type.Fields.List, gfi)

				//构建成员的Get,Set函数
				this.constructFieldFunc(gfi)
			}
		}

		//处理Extends
		if cc.Extends != nil && c.GetType() != ast.ENUM_TYPE {

			Extends := &gast.Field{
				Doc:     nil,
				Type:    this.transType(cc.Extends),
				Tag:     nil,
				Comment: nil,
			}
			Type.Fields.List = append(Type.Fields.List, Extends)
		}

		//如果是枚举类型加一个枚举序号变量 和 枚举序号获取函数
		if cc.GetType() == ast.ENUM_TYPE {
			gfi := &gast.Field{
				Doc:     nil,
				Names:   []*gast.Ident{gast.NewIdent("ordinal")},
				Type:    gast.NewIdent("int"),
				Tag:     nil,
				Comment: nil,
			}
			Type.Fields.List = append(Type.Fields.List, gfi)

			//枚举序号获取函数
			or := this.getOrdinalFN(c)
			this.GolangFile.Decls = append(this.GolangFile.Decls, or)

			val := this.getValuelFN(c)
			this.GolangFile.Decls = append(this.GolangFile.Decls, val)

			sc := this.getScanFN(c)
			this.GolangFile.Decls = append(this.GolangFile.Decls, sc)

			mj := this.getMarshalJSONFN(c)
			this.GolangFile.Decls = append(this.GolangFile.Decls, mj)

			uj := this.getUnmarshalJSONFN(c)
			this.GolangFile.Decls = append(this.GolangFile.Decls, uj)

			vals := this.getValuesFN(c)
			this.GolangFile.Decls = append(this.GolangFile.Decls, vals)

		}

		for _, m := range cc.Methods {
			gmeth := this.transFunc(m)

			if m.IsConstruct() && len(cc.Generics) > 0 {
				gmeth.Type.TypeParams = tpList
			}
			this.GolangFile.Decls = append(this.GolangFile.Decls, gmeth)
		}

		cl.Specs = append(cl.Specs, sp)

	}
	//构建New函数
	if cfg.ConstructNewFunc && c.GetType() == ast.CLASS_TYPE {
		this.GolangFile.Decls = append(this.GolangFile.Decls, this.getNewService(c))

	}

	return
}

func (this *Translation) getValuesFN(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func Values() (result []######) {
	return []######{%%%%%%}
}`
	src = strings.Replace(src, "######", c.GetName(), 2)
	var elements string
	for _, fi := range c.ListFields() {
		if sf, ok := fi.(*ast.FieldEnum); ok {
			if elements != "" {
				elements += ","
			}
			elements += sf.Name
		}

	}
	src = strings.Replace(src, "%%%%%%", elements, 1)
	fn = this.getFunc(src)
	return
}

func (this *Translation) getOrdinalFN(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *######) Ordinal() (result int) {
		return this.ordinal
}`
	src = strings.Replace(src, "######", c.GetName(), 1)
	fn = this.getFunc(src)

	return
}
func (this *Translation) getValuelFN(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *######) Value() (driver.Value, error)  {
		return this.Ordinal(),nil
}`
	src = strings.Replace(src, "######", c.GetName(), 1)
	fn = this.getFunc(src)

	return
}
func (this *Translation) getScanFN(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *######) Scan(v interface{}) error  {
	switch vt := v.(type) {
	case int:
		this.ordinal = vt
		switch vt {

		}
	default:
		this = nil
	}
	return nil
}`
	src = strings.Replace(src, "######", c.GetName(), 1)
	fn = this.getFunc(src)

	if len(fn.Body.List) > 0 {

		if sw1, ok := fn.Body.List[0].(*gast.TypeSwitchStmt); ok {
			if len(sw1.Body.List) == 2 {

				if cl1, ok := sw1.Body.List[0].(*gast.CaseClause); ok {
					if sw2, ok := cl1.Body[1].(*gast.SwitchStmt); ok {
						for idx, fi := range c.ListFields() {
							if sf, ok := fi.(*ast.FieldEnum); ok && len(sf.Values) >= 1 {
								cause := &gast.CaseClause{
									Case:  0,
									List:  nil,
									Colon: 0,
									Body:  nil,
								}
								getStm := &gast.AssignStmt{
									Lhs: []gast.Expr{gast.NewIdent("this.CnName")},
									Tok: token.ASSIGN,
									Rhs: []gast.Expr{this.transExp(sf.Values[0])},
								}
								cause.List = append(cause.List, gast.NewIdent(fmt.Sprintf("%v", idx+1)))
								cause.Body = append(cause.Body, getStm)
								sw2.Body.List = append(sw2.Body.List, cause)

							}

						}
					}

				}

			}

		}
	}

	return
}

func (this *Translation) getMarshalJSONFN(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *######) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", this.ordinal)), nil
}`
	src = strings.Replace(src, "######", c.GetName(), 1)
	fn = this.getFunc(src)

	return
}
func (this *Translation) getUnmarshalJSONFN(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *######) UnmarshalJSON(data []byte) (err error) {
	if data == nil || len(data) == 2 {
		return
	}
	this.ordinal, err = strconv.Atoi(string(data))
	if err != nil {
		return
	}
	switch this.ordinal {

	}

}`
	src = strings.Replace(src, "######", c.GetName(), 1)
	fn = this.getFunc(src)
	sw2 := &gast.SwitchStmt{
		Switch: 0,
		Init:   nil,
		Tag:    gast.NewIdent("this.ordinal"),
		Body:   &gast.BlockStmt{},
	}

	for idx, fi := range c.ListFields() {
		if sf, ok := fi.(*ast.FieldEnum); ok && len(sf.Values) >= 1 {
			cause := &gast.CaseClause{
				Case:  0,
				List:  nil,
				Colon: 0,
				Body:  nil,
			}
			getStm := &gast.AssignStmt{
				Lhs: []gast.Expr{gast.NewIdent("this.CnName")},
				Tok: token.ASSIGN,
				Rhs: []gast.Expr{this.transExp(sf.Values[0])},
			}
			cause.List = append(cause.List, gast.NewIdent(fmt.Sprintf("%v", idx+1)))
			cause.Body = append(cause.Body, getStm)
			sw2.Body.List = append(sw2.Body.List, cause)

		}

	}
	fn.Body.List = append(fn.Body.List, sw2)
	fn.Body.List = append(fn.Body.List, &gast.ReturnStmt{})
	return
}

// 枚举转换
//
// param: c
func (this *Translation) transEnum(c ast.Class) {
	this.currentClass = c
	defer func() {
		this.currentClass = nil
	}()

	if cfg.OneFold {
		this.PkgName = c.GetName()
		this.GolangFile.Name.Name = c.GetName()
	}
	if cc, ok := c.(*ast.ClassSingle); ok {

		//2 解析枚举元素
		v := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.VAR,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		this.GolangFile.Decls = append(this.GolangFile.Decls, v)
		var enumIdx = 1
		for idx, fi := range cc.Fields {
			if fiEn, ok := fi.(*ast.FieldEnum); ok {
				value := &gast.ValueSpec{
					Doc:     nil,
					Names:   []*gast.Ident{gast.NewIdent(fi.GetName())},
					Type:    nil,
					Values:  nil,
					Comment: nil,
				}
				if len(fiEn.Values) > 0 {
					//value.Type = gast.NewIdent(cc.GetName())
					val := &gast.CompositeLit{
						Type:       gast.NewIdent(cc.GetName()),
						Lbrace:     0,
						Elts:       nil,
						Rbrace:     0,
						Incomplete: false,
					}

					for _, vv := range fiEn.Values {
						val.Elts = append(val.Elts, this.transExp(vv))
					}
					//末尾添加一个枚举序号
					val.Elts = append(val.Elts, gast.NewIdent(fmt.Sprintf("%v", enumIdx)))
					enumIdx++

					value.Values = append(value.Values, val)
				} else {
					if idx == 0 {
						value.Type = gast.NewIdent(cc.GetName())
						value.Values = append(value.Values, gast.NewIdent("iota"))
					}
				}

				v.Specs = append(v.Specs, value)
			}

		}
		cl := this.transClass(c)

		this.GolangFile.Decls = append(this.GolangFile.Decls, cl)

		//this.buildEnumString(cc)
		//this.buildEnumName(cc)
		//this.buildEnumGetCode(cc)
	}
}
func (this *Translation) buildEnumName(cc *ast.ClassSingle) {

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
			Name:    this.currentClass.GetName(),
			Obj:     gast.NewObj(gast.Typ, this.currentClass.GetName()),
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
	//String
	swBlock := &gast.BlockStmt{
		Lbrace: 0,
		List:   nil,
		Rbrace: 0,
	}
	for _, fi := range cc.Fields {
		if sf, ok := fi.(*ast.FieldEnum); ok && len(sf.Values) >= 2 {
			//String()
			cause := &gast.CaseClause{
				Case:  0,
				List:  nil,
				Colon: 0,
				Body:  nil,
			}
			getStm := &gast.ReturnStmt{
				Return:  0,
				Results: []gast.Expr{this.transExp(sf.Values[1])},
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

func (this *Translation) buildEnumString(cc *ast.ClassSingle) {

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
			Name:    this.currentClass.GetName(),
			Obj:     gast.NewObj(gast.Typ, this.currentClass.GetName()),
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
	//String
	swBlock := &gast.BlockStmt{
		Lbrace: 0,
		List:   nil,
		Rbrace: 0,
	}
	for _, fi := range cc.Fields {
		if sf, ok := fi.(*ast.FieldEnum); ok {
			//String()
			cause := &gast.CaseClause{
				Case:  0,
				List:  nil,
				Colon: 0,
				Body:  nil,
			}
			getStm := &gast.ReturnStmt{
				Return:  0,
				Results: []gast.Expr{gast.NewIdent("\"" + fi.GetName() + "\"")},
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
		Name: gast.NewIdent("Name"),
		Type: &gast.FuncType{
			Func:    0,
			Params:  nil,
			Results: results,
		},
		Body: getBody,
	}
	this.GolangFile.Decls = append(this.GolangFile.Decls, stringFun)

}

func (this *Translation) buildEnumGetCode(cc *ast.ClassSingle) {

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
			Name:    this.currentClass.GetName(),
			Obj:     gast.NewObj(gast.Typ, this.currentClass.GetName()),
		},
		Tag:     nil,
		Comment: nil,
	}

	recv.List = append(recv.List, gfi)

	//处理返回值
	resType := &gast.Field{
		Doc:     nil,
		Names:   nil,
		Type:    gast.NewIdent("int64"),
		Tag:     nil,
		Comment: nil,
	}

	results := &gast.FieldList{
		Opening: 0,
		List:    []*gast.Field{resType},
		Closing: 0,
	}

	//处理函数体

	var getBody = &gast.BlockStmt{
		Lbrace: 0,
		List:   nil,
		Rbrace: 0,
	}

	retStm := &gast.ReturnStmt{
		Return:  0,
		Results: []gast.Expr{gast.NewIdent("int64(this)")},
	}

	//空return代码块
	getBody.List = append(getBody.List, retStm)

	stringFun := &gast.FuncDecl{
		Doc:  nil,
		Recv: recv,
		Name: gast.NewIdent("GetCode"),
		Type: &gast.FuncType{
			Func:    0,
			Params:  nil,
			Results: results,
		},
		Body: getBody,
	}
	this.GolangFile.Decls = append(this.GolangFile.Decls, stringFun)

}

//
func (this *Translation) transInterface(c ast.Class) {
	this.currentClass = c
	defer func() {
		this.currentClass = nil
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
		Name:    gast.NewIdent(util.DeCapitalize(c.GetName())),
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
					q += fmt.Sprintf(".Where(\"%s = ?\", %s)", util.SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				} else {
					q += fmt.Sprintf(".Where(\"%s = ?\", %s)", util.SnakeString(v), gmeth.Type.Params.List[k].Names[0])
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
					q += fmt.Sprintf(".Where(\"%s <> ?\",%s)", util.SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				} else {
					q += fmt.Sprintf(".Where(\"%s = ?\",%s)", util.SnakeString(v), gmeth.Type.Params.List[k].Names[0])
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
					q += fmt.Sprintf(".Where(\"%s <> ?\",%s)", util.SnakeString(v), gmeth.Type.Params.List[k].Names[0])
				} else {
					q += fmt.Sprintf(".Where(\"%s = ?\",%s)", util.SnakeString(v), gmeth.Type.Params.List[k].Names[0])
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
			v.Type = &gast.StarExpr{X: gast.NewIdent(util.DeCapitalize(c.GetName()))}
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
