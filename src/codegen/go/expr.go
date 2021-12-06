package codegen_go

import (
	"dog/ast"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func (this *Translation) transExp(e ast.Exp) (expr gast.Expr) {
	switch v := e.(type) {
	case *ast.Not:
		expr = &gast.UnaryExpr{
			OpPos: 0,
			Op:    token.NOT,
			X:     this.transExp(v.E),
		}
	case *ast.Or:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.LOR,
			Y:     this.transExp(v.Right),
		}
	case *ast.And:
		expr = &gast.BinaryExpr{
			X:     this.transExp(v.Left),
			OpPos: 0,
			Op:    token.LAND,
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

	case *ast.This:
		//log.Debugf("This表达式")
		return gast.NewIdent("this")
	case *ast.NewList:
		//log.Debugf("初始化List表达式")
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
			Elt:    transType(v.Ele),
		}
		call.Args = append(call.Args, t)

		len := &gast.BasicLit{
			ValuePos: 0,
			Kind:     token.INT,
			Value:    "0",
		}
		call.Args = append(call.Args, len)
		return call
		//TODO 这里需要自己构造一个初始化函数
	case *ast.NewObject:
		call := &gast.CallExpr{
			Fun:      transType(v.T),
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
		expr = &gast.CallExpr{
			Fun:      gast.NewIdent("NewSet"),
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}
	case *ast.SelectorExpr:
		//	log.Debugf("选择表达式,%v, %s", v.X, v.Sel)
		expr = &gast.SelectorExpr{
			X:   this.transExp(v.X),
			Sel: gast.NewIdent(Capitalize(v.Sel)),
		}
	case *ast.CallExpr:
		fn := this.transExp(v.Callee)
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
		return call
	case *ast.Id:
		if _, ok := v.Tp.(*ast.Function); ok {
			if this.CurrentClass.GetMethod(v.Name) != nil {
				v.Name = "this." + Capitalize(v.Name)
			} else {
				v.Name = Capitalize(v.Name)
			}

		}
		return gast.NewIdent(v.Name)
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
	case *ast.ArraySelect:
		return &gast.IndexExpr{
			X:      this.transExp(v.Arrayref),
			Lbrack: 0,
			Index:  this.transExp(v.Index),
			Rbrack: 0,
		}
	case *ast.NewObjectArray:

		if v.Eles == nil {
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
				Elt:    transType(v.T),
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
			panic("NewObjectArray bug")
		}
	case *ast.ClassExp:
		return gast.NewIdent("class")
	case *ast.Question:
		return this.transFuncLit(v)
	case *ast.Lambda:
		return this.transLambda(v)
	case *ast.NewHash:
		call := &gast.CallExpr{
			Fun:      gast.NewIdent("make"),
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}

		t := &gast.MapType{
			Map:   0,
			Key:   transType(v.Key),
			Value: transType(v.Ele),
		}
		call.Args = append(call.Args, t)
		return call
	case *ast.NewStringArray:

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
	case *ast.NewIntArray:

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
			Elt:    gast.NewIdent("int"),
		}
		call.Args = append(call.Args, t)

		len := &gast.BasicLit{
			ValuePos: 0,
			Kind:     token.INT,
			Value:    "0",
		}
		call.Args = append(call.Args, len)
		return call
	default:
		log.Debugf("%v --> %v", reflect.TypeOf(v).String(), v)
		panic("bug")
	}

	return
}