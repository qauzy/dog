package codegen_go

import (
	"dog/ast"
	"dog/cfg"
	"dog/util"
	gast "go/ast"
	"go/token"
	"strconv"
	"strings"
)

func (this *Translation) transNameExp(e ast.Exp) (expr *gast.Ident) {
	switch v := e.(type) {
	case *ast.Ident:
		//是类型标识符, 可能需要转换
		if cfg.Capitalize && nil != this.currentFile && nil != this.currentClass && (nil != this.currentFile.GetField(v.Name) || nil != this.currentClass.GetField(v.Name) || nil != this.currentClass.GetMethod(v.Name)) {
			return gast.NewIdent(util.Capitalize(v.Name))
		}
		expr = gast.NewIdent(util.GetNewId(v.Name))
	case *ast.DefExpr:
		//是类型标识符,可能需要转换
		expr = gast.NewIdent(v.Name.Name)
	default:
		this.TranslationBug("transNameExp")
	}
	return
}

func (this *Translation) transExp(e ast.Exp) (expr gast.Expr) {
	switch v := e.(type) {
	case *ast.Ident:
		//是类型标识符, 可能需要转换
		if cfg.Capitalize && !cfg.AllStatic && nil != this.currentClass && ((nil != this.currentClass.GetField(v.Name) && !this.currentClass.GetField(v.Name).IsStatic()) || (nil != this.currentClass.GetMethod(v.Name) && !this.currentClass.GetMethod(v.Name).IsStatic())) {
			return gast.NewIdent("this." + util.Capitalize(v.Name))
		} else if cfg.Capitalize && cfg.AllStatic && nil != this.currentClass && ((nil != this.currentClass.GetField(v.Name) && !this.currentClass.GetField(v.Name).IsStatic()) || (nil != this.currentClass.GetMethod(v.Name) && !this.currentClass.GetMethod(v.Name).IsStatic())) {
			return gast.NewIdent(util.Capitalize(v.Name))
		} else if cfg.Capitalize && nil != this.currentClass && ((nil != this.currentClass.GetField(v.Name) && this.currentClass.GetField(v.Name).IsStatic()) || (nil != this.currentClass.GetMethod(v.Name) && this.currentClass.GetMethod(v.Name).IsStatic())) {
			return gast.NewIdent(util.Capitalize(v.Name))
		} else if !cfg.Capitalize && nil != this.currentClass && ((nil != this.currentClass.GetField(v.Name) && !this.currentClass.GetField(v.Name).IsStatic()) || (nil != this.currentClass.GetMethod(v.Name) && !this.currentClass.GetMethod(v.Name).IsStatic())) {
			return gast.NewIdent("this." + v.Name)
		}
		//是类型标识符,可能需要转换
		expr = gast.NewIdent(util.GetNewId(v.Name))
	case *ast.Not:
		expr = &gast.UnaryExpr{
			OpPos: 0,
			Op:    token.NOT,
			X:     this.transExp(v.E),
		}
	case *ast.LOr:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.LOR,
			Y:     this.transExp(v.Right),
		}
	case *ast.LAnd:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.LAND,
			Y:     this.transExp(v.Right),
		}
	case *ast.And:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.AND,
			Y:     this.transExp(v.Right),
		}
	case *ast.Lt:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.LSS,
			Y:     this.transExp(v.Right),
		}
	case *ast.Le:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.LEQ,
			Y:     this.transExp(v.Right),
		}
	case *ast.Gt:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.GTR,
			Y:     this.transExp(v.Right),
		}
	case *ast.Ge:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.GEQ,
			Y:     this.transExp(v.Right),
		}
	case *ast.Eq:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.EQL,
			Y:     this.transExp(v.Right),
		}
	case *ast.Neq:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.NEQ,
			Y:     this.transExp(v.Right),
		}
	case *ast.Add:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.ADD,
			Y:     this.transExp(v.Right),
		}
	case *ast.Sub:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.SUB,
			Y:     this.transExp(v.Right),
		}
	case *ast.Times:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.MUL,
			Y:     this.transExp(v.Right),
		}
	case *ast.Division:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.QUO,
			Y:     this.transExp(v.Right),
		}
		return expr
	case *ast.Remainder:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.REM,
			Y:     this.transExp(v.Right),
		}
		return expr
	case *ast.This:
		return gast.NewIdent("this")
	case *ast.NewList:
		if cfg.NoGeneric {
			call := &gast.CallExpr{
				Fun:      gast.NewIdent("make"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}

			t := &gast.ArrayType{
				Lbrack: 0,
				Len:    nil,
				Elt:    this.transType(v.Ele),
			}
			call.Args = append(call.Args, t)

			len := &gast.BasicLit{
				ValuePos: 0,
				Kind:     token.INT,
				Value:    "0",
			}
			call.Args = append(call.Args, len)
			return call
		} else {
			call := &gast.CallExpr{
				Fun: &gast.IndexListExpr{
					X:       gast.NewIdent("arraylist.New"),
					Lbrack:  0,
					Indices: []gast.Expr{this.transType(v.Ele)},
					Rbrack:  0,
				},
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}
			for _, vv := range v.ArgsList {
				call.Args = append(call.Args, this.transExp(vv))
			}
			return call
		}

	case *ast.NewObject:
		call := &gast.CallExpr{
			Fun:      this.transType(v.T),
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}
		for _, a := range v.ArgsList {
			call.Args = append(call.Args, this.transExp(a))
		}
		return call
	case *ast.NewSet:
		//TODO 实现hashset等数据结构
		if cfg.NoGeneric {
			expr = &gast.CallExpr{
				Fun:      gast.NewIdent("NewSet"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}
		} else {
			call := &gast.CallExpr{
				Fun: &gast.IndexListExpr{
					X:       gast.NewIdent("hashset.New"),
					Lbrack:  0,
					Indices: []gast.Expr{this.transType(v.Ele)},
					Rbrack:  0,
				},
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}

			for _, vv := range v.ArgsList {
				call.Args = append(call.Args, this.transExp(vv))
			}
			return call

		}

	case *ast.SelectorExpr:
		//	log.Debugf("选择表达式,%v, %s", v.X, v.Sel)
		expr = &gast.SelectorExpr{
			X:   this.transExp(v.X),
			Sel: gast.NewIdent(util.Capitalize(v.Sel)),
		}
	case *ast.CallExpr:
		fn := this.transExp(v.Callee)
		//调用无参构造函数转化为new
		if cfg.Construct2New && len(v.ArgsList) == 0 {
			if id, ok := fn.(*gast.Ident); ok {
				if im := this.currentFile.GetImport(id.Name); im != nil {
					call := &gast.CallExpr{
						Fun:      gast.NewIdent("new"),
						Lparen:   0,
						Args:     nil,
						Ellipsis: 0,
						Rparen:   0,
					}
					call.Args = append(call.Args, gast.NewIdent(im.GetPack()+"."+im.GetName()))
					return call
				}
			}
		}
		//替换日志打印语句中的 {}
		if vv, ok := fn.(*gast.SelectorExpr); ok && (vv.Sel.Name == "Info" || vv.Sel.Name == "Error") {
			if vvv, ok := vv.X.(*gast.Ident); ok && vvv.Name == "log" {
				if len(v.ArgsList) >= 1 {
					if vvvv, ok := v.ArgsList[0].(*ast.Ident); ok {
						if strings.Contains(vvvv.Name, "{}") {
							vvvv.Name = strings.ReplaceAll(vvvv.Name, "{}", "%v")
							vv.Sel.Name += "f"
						}

					}
				}
			}

		} else if vv, ok := fn.(*gast.SelectorExpr); ok && (vv.Sel.Name == "Get" || vv.Sel.Name == "get") && cfg.NoGeneric {
			if vvv, ok := vv.X.(*gast.Ident); ok && (this.currentFile.GetField(vvv.Name) != nil || this.currentClass.GetField(vvv.Name) != nil || this.currentMethod.GetField(vvv.Name) != nil) {
				if len(v.ArgsList) == 1 {
					f := this.currentClass.GetField(vvv.Name)
					if f == nil {
						f = this.currentMethod.GetField(vvv.Name)
						if f == nil {
							this.currentFile.GetField(vvv.Name)
						}
					}
					_, ok1 := f.GetDecType().(*ast.ListType)
					_, ok2 := f.GetDecType().(*ast.MapType)
					if ok1 || ok2 {
						return &gast.IndexExpr{
							X:      vv.X,
							Lbrack: 0,
							Index:  this.transExp(v.ArgsList[0]),
							Rbrack: 0,
						}
					}
				}
			}

		}

		//替换name
		if f, ok := fn.(*gast.Ident); ok {
			if IdMapper[f.Name] != "" {
				f.Name = IdMapper[f.Name]
			}
		}

		call := &gast.CallExpr{
			Fun:      fn,
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}

		for _, a := range v.ArgsList {
			call.Args = append(call.Args, this.transExp(a))
		}
		if len(call.Args) == 4 {
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
			} else if set, ok := call.Fun.(*gast.SelectorExpr); ok && (set.Sel.Name == "Set") {
				if OpsForValueC, ok := set.X.(*gast.CallExpr); ok {
					if OpsForValue, ok := OpsForValueC.Fun.(*gast.SelectorExpr); ok && (OpsForValue.Sel.Name == "OpsForValue") {
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
						call.Fun = gast.NewIdent("core.SetExpireKV")
					}
				}
			}
		} else if len(call.Args) == 1 {
			if get, ok := call.Fun.(*gast.SelectorExpr); ok && (get.Sel.Name == "Get") {
				if OpsForValueC, ok := get.X.(*gast.CallExpr); ok {
					if OpsForValue, ok := OpsForValueC.Fun.(*gast.SelectorExpr); ok && (OpsForValue.Sel.Name == "OpsForValue") {
						var args []gast.Expr
						args = append(args, gast.NewIdent("true"))
						args = append(args, call.Args...)
						call.Args = args
						call.Fun = gast.NewIdent("core.GetKey")
					}
				}
			} else if get, ok := call.Fun.(*gast.SelectorExpr); ok && (get.Sel.Name == "HasKey") {
				if OpsForValueC, ok := get.X.(*gast.CallExpr); ok {
					if OpsForValue, ok := OpsForValueC.Fun.(*gast.SelectorExpr); ok && (OpsForValue.Sel.Name == "OpsForValue") {
						var args []gast.Expr
						args = append(args, gast.NewIdent("true"))
						args = append(args, call.Args...)
						call.Args = args
						call.Fun = gast.NewIdent("core.KeyExist")
					}
				} else if redisTemplate, ok := get.X.(*gast.Ident); ok && (redisTemplate.Name == "RedisTemplate") {
					var args []gast.Expr
					args = append(args, gast.NewIdent("true"))
					args = append(args, call.Args...)
					call.Args = args
					call.Fun = gast.NewIdent("core.KeyExist")
				}
			}
		}
		return call
	case *ast.DefExpr:
		//if _, ok := v.Tp.(*ast.Function); ok {
		//	if (nil != this.currentClass.GetField(v.Name.Name) && nil == this.currentMethod.GetFormal(v.Name.Name)) || nil != this.currentClass.GetMethod(v.Name.Name) {
		//		return gast.NewIdent("this." + util.Capitalize(v.Name.Name))
		//	}
		//}
		return gast.NewIdent(v.Name.Name)
	case *ast.Num:
		return &gast.BasicLit{
			ValuePos: 0,
			Kind:     token.INT,
			Value:    strconv.Itoa(v.Value),
		}
	case *ast.False:
		return gast.NewIdent("false")
	case *ast.True:
		return gast.NewIdent("true")
	case *ast.Null:
		return gast.NewIdent("nil")
	//列表,数组长度表达式
	case *ast.Length:
		call := &gast.CallExpr{
			Fun:      gast.NewIdent("len"),
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}
		call.Args = append(call.Args, this.transExp(v.Arrayref))
		return call

	//数组索引表达式
	case *ast.IndexExpr:
		return &gast.IndexExpr{
			X:      this.transExp(v.X),
			Lbrack: 0,
			Index:  this.transExp(v.Index),
			Rbrack: 0,
		}
	case *ast.NewObjectArray:

		if v.Eles != nil {

			//FIXME 需要修改
			t := &gast.ArrayType{
				Lbrack: 0,
				Len:    nil,
				Elt:    this.transType(v.T),
			}

			ct := &gast.CompositeLit{
				Type:       t,
				Lbrace:     0,
				Elts:       nil,
				Rbrace:     0,
				Incomplete: false,
			}

			for _, arg := range v.Eles {
				ct.Elts = append(ct.Elts, this.transExp(arg))
			}
			return ct

		} else {

			call := &gast.CallExpr{
				Fun:      gast.NewIdent("make"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}

			t := &gast.ArrayType{
				Lbrack: 0,
				Len:    nil,
				Elt:    this.transType(v.T),
			}
			call.Args = append(call.Args, t)
			if v.Size != nil {
				call.Args = append(call.Args, this.transExp(v.Size))
			} else {
				len := &gast.BasicLit{
					ValuePos: 0,
					Kind:     token.INT,
					Value:    "0",
				}
				call.Args = append(call.Args, len)

			}

			return call

		}
	case *ast.ClassExp:
		return this.transType(v.Name)
	case *ast.Question:
		return this.transFuncLit(v)
	case *ast.Lambda:
		return this.transLambda(v)
	case *ast.NewHash:
		if cfg.NoGeneric {
			call := &gast.CallExpr{
				Fun:      gast.NewIdent("make"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}

			t := &gast.MapType{
				Map:   0,
				Key:   this.transType(v.Key),
				Value: this.transType(v.Ele),
			}
			call.Args = append(call.Args, t)
			return call
		} else {
			call := &gast.CallExpr{
				Fun: &gast.IndexListExpr{
					X:       gast.NewIdent("hashset.New"),
					Lbrack:  0,
					Indices: []gast.Expr{this.transType(v.Ele)},
					Rbrack:  0,
				},
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}
			return call
		}

	case *ast.NewStringArray:
		if cfg.NoGeneric {
			call := &gast.CallExpr{
				Fun:      gast.NewIdent("make"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}

			t := &gast.ArrayType{
				Lbrack: 0,
				Len:    nil,
				Elt:    gast.NewIdent("string"),
			}
			call.Args = append(call.Args, t)

			len := &gast.BasicLit{
				ValuePos: 0,
				Kind:     token.INT,
				Value:    "0",
			}
			call.Args = append(call.Args, len)
			return call
		} else {
			call := &gast.CallExpr{
				Fun: &gast.IndexListExpr{
					X:       gast.NewIdent("arraylist.New"),
					Lbrack:  0,
					Indices: []gast.Expr{gast.NewIdent("string")},
					Rbrack:  0,
				},
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}
			return call
		}

	case *ast.NewIntArray:
		call := &gast.CallExpr{
			Fun: &gast.IndexListExpr{
				X:       gast.NewIdent("arraylist.New"),
				Lbrack:  0,
				Indices: []gast.Expr{gast.NewIdent("int")},
				Rbrack:  0,
			},
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}
		return call
	case *ast.Cast: //强制类型转换
		expr = &gast.TypeAssertExpr{
			X:      this.transExp(v.Right),
			Lparen: 0,
			Type:   this.transType(v.Tp),
			Rparen: 0,
		}
	case *ast.Integer:
		return gast.NewIdent("int")
	case *ast.Float:
		return gast.NewIdent("float64")
	case *ast.NewDate:
		if len(v.Params) == 1 {
			exp := this.getExpr("time.UnixMilli()")
			callExp := exp.(*gast.CallExpr)
			callExp.Args = []gast.Expr{this.transExp(v.Params[0])}
			return exp
		}
		return gast.NewIdent("time.Now()")
	case *ast.String:
		return gast.NewIdent("string")
	case *ast.Instanceof:
		x := &gast.CallExpr{
			Fun: &gast.SelectorExpr{
				X:   gast.NewIdent("reflect"),
				Sel: gast.NewIdent("TypeOf"),
			},
			Lparen:   0,
			Args:     []gast.Expr{this.transExp(v.Left)},
			Ellipsis: 0,
			Rparen:   0,
		}
		return &gast.BinaryExpr{
			X:     x,
			OpPos: 0,
			Op:    token.EQL,
			Y:     this.transExp(v.Right),
		}
	case *ast.BuilderExpr:
		call := &gast.CallExpr{
			Fun:      gast.NewIdent("new"),
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}

		call.Args = append(call.Args, this.transExp(v.X))
		return call
	case *ast.NewArrayWithArgs:
		clit := &gast.CompositeLit{
			Type:       this.transType(v.Tp),
			Lbrace:     0,
			Elts:       nil,
			Rbrace:     0,
			Incomplete: false,
		}
		for _, vv := range v.Args {
			clit.Elts = append(clit.Elts, this.transExp(vv))
		}
		return clit

	case *ast.NewArray:
		call := &gast.CallExpr{
			Fun:      gast.NewIdent("make"),
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}

		t := &gast.ArrayType{
			Lbrack: 0,
			Len:    nil,
			Elt:    this.transExp(v.Tp),
		}
		call.Args = append(call.Args, t)

		length := &gast.BasicLit{
			ValuePos: 0,
			Kind:     token.INT,
			Value:    "0",
		}
		call.Args = append(call.Args, length)
		return call

	case *ast.ArrayAssign:
		clit := &gast.CompositeLit{
			Type:       this.transType(v.Tp),
			Lbrace:     0,
			Elts:       nil,
			Rbrace:     0,
			Incomplete: false,
		}
		for _, v := range v.E {
			clit.Elts = append(clit.Elts, this.transExp(v))
		}
		return clit
	case *ast.MethodReference:
		mp := &gast.BinaryExpr{
			Op: token.COLON,
			X:  this.transExp(v.X),
			Y:  this.transExp(v.Y),
		}
		return mp
	default:
		this.TranslationBug(v)
	}

	return
}
