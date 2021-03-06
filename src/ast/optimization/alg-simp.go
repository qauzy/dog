package ast_opt

//
//import (
//	"dog/ast"
//)
//
//func AlgSimp(prog ast.File) ast.File {
//	var main_class ast.MainClass
//	var classes []ast.Class
//	var new_class ast.Class
//	var methods []ast.Y
//	var stm ast.Stm
//	var method ast.Y
//	var stms []ast.Stm
//	var exp ast.Exp
//	var is_0 bool
//
//	var opt func(e ast.Acceptable)
//
//	opt_Exp := func(ee ast.Exp) {
//		switch e := ee.(type) {
//		case *ast.Add:
//			var left_0 bool
//			var right_0 bool
//			opt(e.Left)
//			if is_0 {
//				left_0 = true
//			} else {
//				left_0 = false
//			}
//			left := exp
//			opt(e.Right)
//			if is_0 {
//				right_0 = true
//			} else {
//				right_0 = false
//			}
//			right := exp
//			if left_0 && right_0 {
//				exp = ast.Num_new(0, e.LineNum)
//			} else if left_0 {
//				exp = right
//			} else if right_0 {
//				exp = left
//			} else {
//				exp = e
//			}
//		case *ast.And:
//			is_0 = false //XXX
//			opt(e.Left)
//			left := exp
//			opt(e.Right)
//			right := exp
//			exp = ast.And_new(left, right, e.LineNum)
//		case *ast.ArraySelect:
//			opt(e.Index)
//			index := exp
//			exp = ast.ArraySelect_new(e.X, index, e.LineNum)
//		case *ast.Call:
//			is_0 = false
//			args := make([]ast.Exp, 0)
//			opt(e.Callee)
//			callee := exp
//			for _, arg := range e.ArgsList {
//				opt(arg)
//				args = append(args, exp)
//			}
//			exp = ast.Call_new(callee,
//				e.MethodName,
//				args,
//				e.Firsttype,
//				e.ArgsType,
//				e.Rt,
//				e.LineNum)
//		case *ast.False:
//			is_0 = false
//			exp = e
//		case *ast.Id:
//			is_0 = false
//			exp = e
//		case *ast.Length:
//			is_0 = false
//			opt(e.X)
//			array := exp
//			exp = ast.Length_new(array, e.LineNum)
//		case *ast.Lt:
//			is_0 = false
//			opt(e.Left)
//			left := exp
//			opt(e.Right)
//			right := exp
//			exp = ast.Lt_new(left, right, e.LineNum)
//		case *ast.NewIntArray:
//			is_0 = false
//			opt(e.Size)
//			size := exp
//			exp = ast.NewIntArray_new(size, e.LineNum)
//		case *ast.NewObject:
//			is_0 = false
//			exp = e
//		case *ast.Not:
//			is_0 = false
//			opt(e.Values)
//			new_exp := exp
//			exp = ast.Not_new(new_exp, e.LineNum)
//		case *ast.Num:
//			//XXX important!
//			is_0 = false
//			if e.Values == 0 {
//				is_0 = true
//			} else {
//				is_0 = false
//			}
//			exp = e
//		case *ast.Sub:
//			var left_0 bool
//			var right_0 bool
//			is_0 = false
//			opt(e.Left)
//			//left := exp
//			if is_0 {
//				left_0 = true
//			} else {
//				left_0 = false
//			}
//			opt(e.Right)
//			//right := exp
//			if is_0 {
//				right_0 = true
//			} else {
//				right_0 = false
//			}
//			if left_0 && right_0 {
//				exp = ast.Num_new(0, e.LineNum)
//				/*XXX sub can't opt right
//				                //XXX issue #11
//							} else if right_0 {
//								exp = left
//							} else if left_0 {
//								exp = right
//				*/
//			} else {
//				exp = e
//			}
//		case *ast.This:
//			is_0 = false
//			exp = e
//		case *ast.Times:
//			var left_0 bool
//			var right_0 bool
//			is_0 = false
//			opt(e.Left)
//			//left := exp
//			if is_0 {
//				left_0 = true
//			} else {
//				left_0 = false
//			}
//			opt(e.Right)
//			// right := exp
//			if is_0 {
//				right_0 = true
//			} else {
//				right_0 = false
//			}
//			if left_0 || right_0 {
//				exp = ast.Num_new(0, e.LineNum)
//			} else {
//				exp = e
//			}
//		case *ast.True:
//			is_0 = false
//			exp = e
//		default:
//			panic("impossible")
//		}
//	}
//
//	opt_Stm := func(ss ast.Stm) {
//		switch s := ss.(type) {
//		case *ast.Assign:
//			opt_Exp(s.Values)
//			stm = ast.Assign_new(s.Names, exp, s.Tp, s.IsField, s.LineNum)
//		case *ast.AssignArray:
//			opt(s.Values)
//			ee := exp
//			opt(s.Index)
//			index := exp
//			stm = ast.AssignArray_new(s.Names, index, ee, s.Tp, s.IsField, s.LineNum)
//		case *ast.Block:
//			ss := make([]ast.Stm, 0)
//			for _, s0 := range s.Stms {
//				opt(s0)
//				ss = append(ss, stm)
//			}
//			stm = ast.Block_new(ss, s.LineNum)
//		case *ast.If:
//			opt(s.Condition)
//			cond := exp
//			opt(s.Body)
//			thenn := stm
//			opt(s.Elsee)
//			elsee := stm
//			stm = ast.If_new(cond, thenn, elsee, s.LineNum)
//		case *ast.Print:
//			opt(s.Values)
//			stm = ast.Print_new(exp, s.LineNum)
//		case *ast.While:
//			opt(s.Values)
//			cond := exp
//			opt(s.Body)
//			body := stm
//			stm = ast.While_new(cond, body, s.LineNum)
//		default:
//			panic("impossible")
//		}
//	}
//
//	opt_MainClass := func(m ast.MainClass) {
//		if mc, ok := m.(*ast.MainClassSingle); ok {
//			opt_Stm(mc.Stms)
//			main_class = &ast.MainClassSingle{mc.Names, mc.Args, stm}
//		} else {
//			panic("impossible")
//		}
//	}
//
//	opt_Method := func(mm ast.Y) {
//		if m, ok := mm.(*ast.MethodSingle); ok {
//			stms = make([]ast.Stm, 0)
//			for _, s := range m.Stms {
//				opt(s)
//				stms = append(stms, stm)
//			}
//			opt(m.RetExp)
//			method = &ast.MethodSingle{m.RetType,
//				m.Names,
//				m.Formals,
//				m.Locals,
//				stms,
//				exp}
//		} else {
//			panic("impossible")
//		}
//	}
//
//	opt_Class := func(cc ast.Class) {
//		if c, ok := cc.(*ast.ClassSingle); ok {
//			methods = make([]ast.Y, 0)
//			for _, m := range c.Methods {
//				opt(m)
//				methods = append(methods, method)
//			}
//			new_class = &ast.ClassSingle{c.Access, c.Names, c.Extends, c.Fields, methods}
//		} else {
//			panic("impossible")
//		}
//	}
//
//	opt = func(e ast.Acceptable) {
//		switch v := e.(type) {
//		case ast.Exp:
//			opt_Exp(v)
//		case ast.Stm:
//			opt_Stm(v)
//		case ast.Y:
//			opt_Method(v)
//		case ast.MainClass:
//			opt_MainClass(v)
//		case ast.Class:
//			opt_Class(v)
//		case ast.Field:
//			//no need
//		case ast.Type:
//			//no need
//		default:
//			panic("impossible")
//		}
//	}
//
//	var Ast ast.File
//	if p, ok := prog.(*ast.FileSingle); ok {
//		opt(p.Mainclass)
//		classes = make([]ast.Class, 0)
//		for _, c := range p.Classes {
//			opt(c)
//			classes = append(classes, new_class)
//		}
//		Ast = &ast.FileSingle{p.Names, main_class, classes}
//	} else {
//		panic("impossible")
//	}
//
//	return Ast
//}
