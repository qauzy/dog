package codegen_go

import (
	"dog/ast"
	"dog/cfg"
	log "github.com/corgi-kx/logcustom"
	"github.com/xwb1989/sqlparser"
	gast "go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

func (this *Translation) GetErrReturn() gast.Stmt {
	l := &gast.ExprStmt{X: gast.NewIdent("log.Errorf(\"json.Unmarshal,err=%v\", err)")}
	blk := &gast.BlockStmt{
		Lbrace: 0,
		List:   []gast.Stmt{l, &gast.ReturnStmt{}},
		Rbrace: 0,
	}

	iferr := &gast.IfStmt{
		If:   0,
		Init: nil,
		Cond: gast.NewIdent("err != nil"),
		Body: blk,
		Else: nil,
	}
	return iferr
}
func (this *Translation) transDefine(s ast.Stm) (stmt gast.Stmt) {
	switch v := s.(type) {
	//变量声明
	case *ast.DeclStmt:
		var names, values []gast.Expr
		for _, name := range v.Names {
			names = append(names, this.transExp(name))
		}

		for _, value := range v.Values {
			if v.Values != nil {
				values = append(values, this.transExp(value))
			}
		}
		stmt = &gast.AssignStmt{
			Lhs:    names,
			TokPos: 0,
			Tok:    token.DEFINE,
			Rhs:    values,
		}
	default:
		this.TranslationBug(v)
	}
	return
}

//
//
// param: s
// return:
func (this *Translation) transStm(s ast.Stm) (stmt gast.Stmt) {
	log.Debugf("transStm = %v", s)
	switch v := s.(type) {
	//变量声明
	case *ast.DeclStmt:

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
			Names:   []*gast.Ident{},
			Type:    this.transType(v.Tp),
			Values:  nil,
			Comment: nil,
		}

		_, ok := sp.Type.(*gast.SelectorExpr)
		if ok && cfg.StarClassTypeDecl {
			sp.Type = &gast.StarExpr{
				Star: 0,
				X:    sp.Type,
			}
		}
		if len(v.Names) == len(v.Values) && v.GetExtra() == nil {
			if len(v.Values) == 1 {
				//?语句,用if语句实现
				if vv, ok := v.Values[0].(*ast.Question); ok {
					blk := &FakeBlock{}
					sp.Names = append(sp.Names, this.transNameExp(v.Names[0]))
					d.Specs = append(d.Specs, sp)
					stmt = &gast.DeclStmt{Decl: d}
					blk.List = append(blk.List, stmt)
					q := &gast.IfStmt{
						If:   0,
						Init: nil,
						Cond: this.transExp(vv.E),
						Body: &gast.BlockStmt{
							Lbrace: 0,
							List: []gast.Stmt{&gast.AssignStmt{
								Lhs:    []gast.Expr{this.transExp(v.Names[0])},
								TokPos: 0,
								Tok:    token.ASSIGN,
								Rhs:    []gast.Expr{this.transExp(vv.One)}}},
						},

						Else: &gast.BlockStmt{
							Lbrace: 0,
							List: []gast.Stmt{&gast.AssignStmt{
								Lhs:    []gast.Expr{this.transExp(v.Names[0])},
								TokPos: 0,
								Tok:    token.ASSIGN,
								Rhs:    []gast.Expr{this.transExp(vv.Two)}}},
							Rbrace: 0,
						},
					}
					blk.List = append(blk.List, q)
					return blk
					//json解析语句
				} else if call, ok := v.Values[0].(*ast.CallExpr); ok {
					expp := this.transExp(call)
					if vv, ok := expp.(*gast.CallExpr); ok && len(vv.Args) == 2 {
						//处理json
						if vvv, ok := vv.Fun.(*gast.SelectorExpr); ok && (vvv.Sel.Name == "ParseObject") {
							if vvvv, ok := vvv.X.(*gast.Ident); ok && (vvvv.Name == "JSON") {
								//转换json解析
								vv.Args = []gast.Expr{vv.Args[0], this.transExp(v.Names[0])}
								vv.Fun = gast.NewIdent("json.Unmarshal")
								as := &gast.AssignStmt{
									Lhs:    []gast.Expr{gast.NewIdent("err")},
									TokPos: 0,
									Tok:    token.ASSIGN,
									Rhs:    []gast.Expr{expp},
								}

								//组装语句
								blk := &FakeBlock{}
								sp.Names = append(sp.Names, this.transNameExp(v.Names[0]))
								d.Specs = append(d.Specs, sp)
								stmt = &gast.DeclStmt{Decl: d}
								blk.List = append(blk.List, stmt)
								blk.List = append(blk.List, as)
								blk.List = append(blk.List, this.GetErrReturn())
								return blk

							}
						}

					}

				} else if _, ok := v.Values[0].(*ast.Null); ok {
					v.Values = nil
				} else {
					sp.Type = nil
				}
			}

		}
		for _, name := range v.Names {
			sp.Names = append(sp.Names, this.transNameExp(name))
		}
		if v.GetExtra() == nil {

			for _, value := range v.Values {
				if v.Values != nil {
					sp.Values = append(sp.Values, this.transExp(value))
				}
			}
		} else {
			log.Debugf("v.GetExtra() = %v", v.GetExtra())
		}

		d.Specs = append(d.Specs, sp)
		stmt = &gast.DeclStmt{Decl: d}
		return this.OptimitcStreamStm(stmt)
		//赋值语句
	case *ast.Assign:
		//TODO stream 语句
		expp := this.transExp(v.Value)

		//?语句,用if语句实现
		if vv, ok := v.Value.(*ast.Question); ok {
			q := &gast.IfStmt{
				If:   0,
				Init: nil,
				Cond: this.transExp(vv.E),
				Body: &gast.BlockStmt{
					Lbrace: 0,
					List: []gast.Stmt{&gast.AssignStmt{
						Lhs:    []gast.Expr{this.transExp(v.Left)},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []gast.Expr{this.transExp(vv.One)}}},
				},

				Else: &gast.BlockStmt{
					Lbrace: 0,
					List: []gast.Stmt{&gast.AssignStmt{
						Lhs:    []gast.Expr{this.transExp(v.Left)},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []gast.Expr{this.transExp(vv.Two)}}},
					Rbrace: 0,
				},
			}

			return q
		} else {

			stmt = &gast.AssignStmt{
				Lhs:    []gast.Expr{this.transExp(v.Left)},
				TokPos: 0,
				Tok:    token.ASSIGN,
				Rhs:    []gast.Expr{expp},
			}

			return this.OptimitcStreamStm(stmt)

		}

	case *ast.If:
		var el gast.Stmt
		var Init gast.Stmt
		if v.Elsee != nil {
			el = this.transStm(v.Elsee)

		} else {
			el = nil
		}
		if v.Init != nil {
			Init = &gast.AssignStmt{
				Lhs: []gast.Expr{
					gast.NewIdent("_"),
					gast.NewIdent("ok"),
				},
				TokPos: 0,
				Tok:    token.DEFINE,
				Rhs:    []gast.Expr{this.transExp(v.Init)},
			}
		}

		stmt = &gast.IfStmt{
			If:   0,
			Init: Init,
			Cond: this.transExp(v.Condition),
			Body: this.transBlock(v.Body),
			Else: el,
		}

	//
	case *ast.For:
		//log.Debugf("For语句:%v", *v)
		stmt = &gast.ForStmt{
			For:  0,
			Init: this.transDefine(v.Init),
			Cond: this.transExp(v.Cond),
			Post: this.transStm(v.Post),
			Body: this.transBlock(v.Body),
		}

	case *ast.Range:
		//log.Debugf("Range语句:%v", *v)
		stmt = &gast.RangeStmt{
			For:    0,
			Key:    gast.NewIdent("_"),
			Value:  this.transExp(v.Value),
			TokPos: 0,
			Tok:    token.DEFINE,
			X:      this.transExp(v.E),
			Body:   this.transBlock(v.Body),
		}
	case *ast.ExprStm:
		expp := this.transExp(v.E)
		if vv, ok := expp.(*gast.CallExpr); ok && len(vv.Args) == 2 {
			if vvv, ok := vv.Fun.(*gast.SelectorExpr); ok && (vvv.Sel.Name == "Put" || vvv.Sel.Name == "put") && cfg.NoGeneric {
				idx := &gast.IndexExpr{
					X:      vvv.X,
					Lbrack: 0,
					Index:  vv.Args[0],
					Rbrack: 0,
				}

				as := &gast.AssignStmt{
					Lhs:    []gast.Expr{idx},
					TokPos: 0,
					Tok:    token.ASSIGN,
					Rhs:    []gast.Expr{vv.Args[1]},
				}
				return as
			}

			//处理 xxx.Stream().ForEach
		} else if call, ok := expp.(*gast.CallExpr); ok && len(call.Args) == 1 {
			if Delete, ok := call.Fun.(*gast.SelectorExpr); ok && (Delete.Sel.Name == "Delete") {
				call.Fun = gast.NewIdent("core.DelKey")
			}

		} else if call, ok := expp.(*gast.CallExpr); ok && len(call.Args) == 4 {

			if SetIfAbsent, ok := call.Fun.(*gast.SelectorExpr); ok && (SetIfAbsent.Sel.Name == "SetIfAbsent") {
				dulFun := &gast.CallExpr{
					Fun:      gast.NewIdent("time.Duration"),
					Lparen:   0,
					Args:     []gast.Expr{call.Args[2]},
					Ellipsis: 0,
					Rparen:   0,
				}
				call.Args[2] = &gast.BinaryExpr{
					X:     dulFun,
					OpPos: 0,
					Op:    token.MUL,
					Y:     call.Args[3],
				}
				call.Args = call.Args[:3]
				call.Fun = gast.NewIdent("core.SetNX")
			}
		}
		//log.Debugf("表达式语句:%v", v)

		stmt = &gast.ExprStmt{X: expp}
		return this.OptimitcStreamStm(stmt)
	case *ast.Throw:
		stmt = &gast.ReturnStmt{
			Return:  0,
			Results: []gast.Expr{this.transExp(v.E)},
		}
	case *ast.Return:
		//log.Debugf("Return语句:%v", v)
		result := &gast.ReturnStmt{
			Return:  0,
			Results: nil,
		}
		if v.E != nil {
			ret := this.transExp(v.E)
			if ret != nil {
				result.Results = append(result.Results, ret)
			}
		}
		return result
	case *ast.Try:
		block := &FakeBlock{}

		fn := &gast.FuncLit{
			Type: &gast.FuncType{
				Func:    0,
				Params:  nil,
				Results: nil,
			},
			Body: nil,
		}
		//处理返回值
		results := &gast.FieldList{
			Opening: 0,
			List:    nil,
			Closing: 0,
		}
		results.List = append(results.List, this.getField(gast.NewIdent("err"), gast.NewIdent("error")))

		//try 语句转为闭包函数
		fn.Type.Results = results
		try := this.transBlock(v.Body)
		try.List = append(try.List, &gast.ReturnStmt{})

		fn.Body = try
		call := &gast.CallExpr{
			Fun:      fn,
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}
		as := &gast.AssignStmt{
			Lhs:    []gast.Expr{gast.NewIdent("err")},
			TokPos: 0,
			Tok:    token.DEFINE,
			Rhs:    []gast.Expr{call},
		}
		block.List = append(block.List, as)

		for _, vv := range v.Catches {
			ifStmt := &gast.IfStmt{
				If:   0,
				Init: nil,
				Cond: this.transExp(vv.Test[0]), //FIXME 兼容多个Exception
				Body: nil,
				Else: nil,
			}
			stms := this.transStm(vv.Body)
			//为了兼容处理要翻译为多个stm的语句
			if blk, ok := stms.(*FakeBlock); ok {
				bl := &gast.BlockStmt{}
				for _, st := range blk.List {
					bl.List = append(bl.List, st)
				}
				block.List = append(block.List, ifStmt)
			} else if blk, ok := stms.(*gast.BlockStmt); ok {
				ifStmt.Body = blk
				block.List = append(block.List, ifStmt)
			} else {
				ifStmt.Body = &gast.BlockStmt{
					Lbrace: 0,
					List:   []gast.Stmt{stms},
					Rbrace: 0,
				}
				block.List = append(block.List, ifStmt)
			}
		}
		return block
	case *ast.While:
		stmt = &gast.ForStmt{
			For:  0,
			Init: nil,
			Cond: this.transExp(v.E),
			Post: nil,
			Body: this.transBlock(v.Body),
		}
	case *ast.Block:
		return this.transBlock(v)
	case *ast.Binary:
		var opt token.Token
		switch v.Opt {
		case "*=":
			opt = token.MUL_ASSIGN
		case "/=":
			opt = token.QUO_ASSIGN
		case "+=":
			opt = token.ADD_ASSIGN
		case "-=":
			opt = token.SUB_ASSIGN
		case "%=":
			opt = token.REM_ASSIGN
		default:
			panic("*ast.Binary")
		}
		stmt = &gast.AssignStmt{
			Lhs:    []gast.Expr{this.transExp(v.Left)},
			TokPos: 0,
			Tok:    opt,
			Rhs:    []gast.Expr{this.transExp(v.Right)},
		}
	case *ast.Switch:
		return &gast.SwitchStmt{
			Switch: 0,
			Init:   nil,
			Tag:    this.transExp(v.E),
			Body:   this.transBlock(v.Cases),
		}
	case *ast.Case:
		var ca []gast.Expr
		if v.E != nil {
			ca = []gast.Expr{this.transExp(v.E)}
		}

		cs := &gast.CaseClause{
			Case:  0,
			List:  ca,
			Colon: 0,
		}

		stms := this.transStm(v.Body)
		//为了兼容处理要翻译为多个stm的语句
		if blk, ok := stms.(*FakeBlock); ok {
			for _, st := range blk.List {
				cs.Body = append(cs.Body, st)
			}
		} else if bb, ok := stms.(*gast.BlockStmt); ok {
			cs.Body = bb.List
		} else {
			cs.Body = append(cs.Body, stms)
		}

		return cs
	case *ast.Comment:
		stmt = &gast.ExprStmt{X: gast.NewIdent(v.C)}
	case *ast.Assert:
		cond := this.transExp(v.Cond)
		block := new(gast.BlockStmt)
		block.List = append(block.List, &gast.ExprStmt{X: this.transExp(v.E)})
		switch v.Opt {
		case "isTrue":
			if e, ok := cond.(*gast.UnaryExpr); ok {
				if e.Op == token.NOT {
					cond = e.X
				}
			} else {
				cond = &gast.BinaryExpr{
					X:     cond,
					OpPos: 0,
					Op:    token.EQL,
					Y:     gast.NewIdent("false"),
				}
			}

		case "isNull":
			cond = &gast.BinaryExpr{
				X:     cond,
				OpPos: 0,
				Op:    token.NEQ,
				Y:     gast.NewIdent("nil"),
			}
		case "notNull":
			cond = &gast.BinaryExpr{
				X:     cond,
				OpPos: 0,
				Op:    token.EQL,
				Y:     gast.NewIdent("nil"),
			}
		case "hasText":
			cond = &gast.BinaryExpr{
				X:     cond,
				OpPos: 0,
				Op:    token.EQL,
				Y:     gast.NewIdent("\"\""),
			}
		default:
			this.TranslationBug("Assert语句转换bug")
		}
		stmt = &gast.IfStmt{
			If:   0,
			Init: nil,
			Cond: cond,
			Body: block,
			Else: nil,
		}

	case *ast.Print:
	//stmt = &gast.ExprStmt{X: gast.NewIdent("fmt.Print")}
	case *ast.StreamStm:
		block := new(gast.BlockStmt)

		stmt = &gast.RangeStmt{
			For:    0,
			Key:    gast.NewIdent("_"),
			Value:  gast.NewIdent("vo"),
			TokPos: 0,
			Tok:    token.DEFINE,
			X:      this.transExp(v.List),
			Body:   block,
		}
		//
		//mp :=
		//
		lf := this.transExp(v.Left)
		if v.ToAny == "toList" {

			call := &gast.CallExpr{
				Fun:      gast.NewIdent("append"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}
			call.Args = append(call.Args, lf)
			call.Args = append(call.Args, gast.NewIdent("vo"))

			as := &gast.AssignStmt{
				Lhs:    []gast.Expr{lf},
				TokPos: 0,
				Tok:    token.ASSIGN,
				Rhs:    []gast.Expr{call},
			}
			block.List = append(block.List, as)
		} else if v.ToAny == "toSet" {
			var method *gast.Ident
			if mr, ok := v.Ele.(*ast.MethodReference); ok {
				method = this.transNameExp(mr.Y)
			}
			el := &gast.CallExpr{
				Fun: &gast.SelectorExpr{
					X:   gast.NewIdent("vo"),
					Sel: method,
				},
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}

			idx := &gast.IndexExpr{
				X:      lf,
				Lbrack: 0,
				Index:  el,
				Rbrack: 0,
			}

			as := &gast.AssignStmt{
				Lhs:    []gast.Expr{idx},
				TokPos: 0,
				Tok:    token.ASSIGN,
				Rhs: []gast.Expr{&gast.BasicLit{
					ValuePos: 0,
					Kind:     token.STRING,
					Value:    "\"1\"",
				}},
			}
			block.List = append(block.List, as)
		} else if v.ToAny == "joining" {
			call := &gast.CallExpr{
				Fun:      gast.NewIdent("strconv.Itoa"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}
			call.Args = append(call.Args, gast.NewIdent("vo"))

			as := &gast.AssignStmt{
				Lhs:    []gast.Expr{lf},
				TokPos: 0,
				Tok:    token.ADD_ASSIGN,
				Rhs:    []gast.Expr{call},
			}
			block.List = append(block.List, as)
		}
	case *ast.AssignArray:
		idx := &gast.IndexExpr{
			X:      gast.NewIdent(v.Name),
			Lbrack: 0,
			Index:  this.transExp(v.Index),
			Rbrack: 0,
		}

		as := &gast.AssignStmt{
			Lhs:    []gast.Expr{idx},
			TokPos: 0,
			Tok:    token.ASSIGN,
			Rhs:    []gast.Expr{this.transExp(v.E)},
		}
		return as
	case *ast.Query:
		stm, err := sqlparser.Parse(v.SQL)
		if err != nil {
			log.Errorf("Query=%v", err)
		}

		// Otherwise do something with stmt
		switch stm := stm.(type) {
		case *sqlparser.Select:
			_ = stm
		case *sqlparser.Insert:
		}

		var args string
		for idx, arg := range this.currentMethod.ListFormal() {
			args += ","
			args += arg.GetName()
			//不是原生的要处理下
			if !v.NativeQuery {
				v.SQL = strings.Replace(v.SQL, ":"+arg.GetName(), "?", -1)
				v.SQL = strings.Replace(v.SQL, "?"+strconv.Itoa(idx+1), "?", 1)
			}
		}
		exe := `
	//FIXME 非原生sql，需要处理
	eng := this.DBWrite().Exec(` + v.SQL + args + `)
	err = eng.Error`

		stmt = &gast.ExprStmt{X: gast.NewIdent(exe)}
	default:
		this.TranslationBug(v)

	}

	return
}

func (this *Translation) transBlock(s ast.Stm) (block *gast.BlockStmt) {
	//log.Debugf("解析Block语句")
	block = new(gast.BlockStmt)
	if bl, ok := s.(*ast.Block); ok {
		for _, st := range bl.Stms {
			if st != nil {
				if st.IsTriple() {
					sss := this.transTriple(st)
					if sss != nil && len(sss) > 0 {
						block.List = append(block.List, sss...)
					}
				} else {
					stms := this.transStm(st)
					//为了兼容处理要翻译为多个stm的语句
					if blk, ok := stms.(*FakeBlock); ok {
						for _, st := range blk.List {
							block.List = append(block.List, st)
						}
					} else {
						block.List = append(block.List, stms)
					}

				}
			}
		}
	} else {
		log.Debugf("transBlock-->%v--->%v", reflect.TypeOf(s).String(), s)
		panic("bug")
	}

	return
}
