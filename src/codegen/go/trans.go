package codegen_go

import (
	"bytes"
	"dog/ast"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"go/format"
	"go/token"
)

func TransGo(p ast.File) (f *gast.File) {

	return transFile(p)
}

func transFile(p ast.File) (f *gast.File) {
	if pp, ok := p.(*ast.FileSingle); ok {
		f = &gast.File{}
		f.Name = gast.NewIdent(pp.Name)
		for _, c := range pp.Classes {
			cl := transClass(c, f)
			f.Decls = append(f.Decls, cl)
		}
		// 输出Go代码
		header := `// Code generated by log-gen. DO NOT EDIT.`
		buffer := bytes.NewBufferString(header)
		//	fset := token.NewFileSet()
		err := astToGo(buffer, f)
		if err != nil {
			return
		}
		//gast.Print(fset, f)
		//gast.Inspect(f, func(n gast.Node) bool {
		//	// Called recursively.
		//	gast.Print(fset, n)
		//	return true
		//})

		fmt.Print(buffer)

	} else {
		panic("bug")
	}
	return
}

//
//
// param: c
// return:
func transClass(c ast.Class, f *gast.File) (cl *gast.GenDecl) {
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
			gfi := transField(fi)
			Type.Fields.List = append(Type.Fields.List, gfi)
		}
		for _, m := range cc.Methods {
			gmeth := transFunc(m)

			//处理类接收
			recv := &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			}
			gfi := &gast.Field{
				Doc:   nil,
				Names: []*gast.Ident{gast.NewIdent("this")},
				Type: &gast.StarExpr{X: &gast.Ident{
					NamePos: 0,
					Name:    cc.Name,
					Obj:     gast.NewObj(gast.Typ, cc.Name),
				}},
				Tag:     nil,
				Comment: nil,
			}

			recv.List = append(recv.List, gfi)

			gmeth.Recv = recv
			f.Decls = append(f.Decls, gmeth)
		}

		cl.Specs = append(cl.Specs, sp)

	}
	return
}
func transFunc(fi ast.Method) (fn *gast.FuncDecl) {
	if method, ok := fi.(*ast.MethodSingle); ok {
		//处理函数参数
		params := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}
		for _, p := range method.Formals {
			params.List = append(params.List, transField(p))
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
		ret := transField(rel)
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
			ss := transStm(stm)
			if ss != nil {
				body.List = append(body.List, ss)
			}

		}

		fn = &gast.FuncDecl{
			Doc:  nil,
			Recv: nil,
			Name: gast.NewIdent(method.Name),
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

func transField(fi ast.Field) (gfi *gast.Field) {
	if field, ok := fi.(*ast.FieldSingle); ok {
		gfi = &gast.Field{
			Doc:     nil,
			Names:   []*gast.Ident{gast.NewIdent(field.Name)},
			Type:    transType(field.Tp),
			Tag:     nil,
			Comment: nil,
		}

	}
	return
}

func transStm(s ast.Stm) (stmt gast.Stmt) {
	switch v := s.(type) {
	//变量声明
	case *ast.Decl:
		log.Info("变量声明")
		d := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.VAR,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		sp := &gast.ValueSpec{
			Doc:     nil,
			Names:   []*gast.Ident{gast.NewIdent(v.Name)},
			Type:    transType(v.Tp),
			Values:  nil,
			Comment: nil,
		}
		val := transExp(v.Value)
		if val != nil {
			sp.Values = append(sp.Values, val)
		}

		d.Specs = append(d.Specs, sp)
		stmt = &gast.DeclStmt{Decl: d}
		//赋值语句
	case *ast.Assign:
		stmt = &gast.AssignStmt{
			Lhs:    []gast.Expr{transExp(v.Left)},
			TokPos: 0,
			Tok:    token.ASSIGN,
			Rhs:    []gast.Expr{transExp(v.E)},
		}
	}

	return
}

func transExp(e ast.Exp) (expr gast.Expr) {
	switch v := e.(type) {
	case *ast.Or:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.OR,
			Y:     transExp(v.Right),
		}
	case *ast.And:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.AND,
			Y:     transExp(v.Right),
		}
	case *ast.Lt:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.LSS,
			Y:     transExp(v.Right),
		}
	case *ast.Le:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.LEQ,
			Y:     transExp(v.Right),
		}
	case *ast.Gt:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.GTR,
			Y:     transExp(v.Right),
		}
	case *ast.Ge:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.GEQ,
			Y:     transExp(v.Right),
		}
	case *ast.Eq:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.EQL,
			Y:     transExp(v.Right),
		}
	case *ast.Neq:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.NEQ,
			Y:     transExp(v.Right),
		}
	case *ast.Add:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.ADD,
			Y:     transExp(v.Right),
		}
	case *ast.Sub:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.SUB,
			Y:     transExp(v.Right),
		}
	case *ast.Times:
		expr = &gast.BinaryExpr{
			X:     transExp(v.Left),
			OpPos: 0,
			Op:    token.MUL,
			Y:     transExp(v.Right),
		}
		//case *ast.Dot:
		//	expr = &gast.BinaryExpr{
		//		X:     transExp(v.Left),
		//		OpPos: 0,
		//		Op:    token.EQL,
		//		Y:     transExp(v.Right),
		//	}
		//call := &gast.CallExpr{
		//	Fun:      nil,
		//	Lparen:   0,
		//	Args:     nil,
		//	Ellipsis: 0,
		//	Rparen:   0,
		//}

	}

	return
}

func transType(t ast.Type) (Type gast.Expr) {
	switch v := t.(type) {
	case *ast.Void:
		return nil
	case *ast.String:
		return gast.NewIdent("string")
	case *ast.Integer:
		return gast.NewIdent("int64")
	case *ast.Int:
		return gast.NewIdent("int")
	case *ast.IntArray:
		return gast.NewIdent("[]int")
	case *ast.HashType:
		return &gast.MapType{
			Map:   0,
			Key:   transType(v.Key),
			Value: transType(v.Value),
		}
	case *ast.ListType:
		return &gast.ArrayType{
			Lbrack: 0,
			Len:    nil,
			Elt:    transType(v.Ele),
		}
	case *ast.ClassType:
		return &gast.Ident{
			NamePos: 0,
			Name:    v.Name,
			Obj:     gast.NewObj(gast.Typ, v.Name),
		}
	case *ast.Boolean:
		return gast.NewIdent("bool")
	default:
		log.Info(v.String())
		panic("impossible")
	}
}

func astToGo(dst *bytes.Buffer, node interface{}) error {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			log.Info(err)
		}
	}

	addNewline()

	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		return err
	}

	addNewline()

	return nil
}
