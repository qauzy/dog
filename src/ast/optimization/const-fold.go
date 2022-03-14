package ast_opt

//
//import (
//	"dog/ast"
//)
//
//func ConstFold(prog ast.File) ast.File {
//	var class ast.Class
//	var main_class ast.MainClass
//	var method ast.Y
//	var exp ast.Exp
//	var stm ast.Stm
//	var opt func(e ast.Acceptable)
//
//	opt_Exp := func(ee ast.Exp) {
//		switch e := ee.(type) {
//		case *ast.Add:
//			opt(e.Left)
//			left := exp
//			opt(e.Right)
//			right := exp
//			l_num, ok := left.(*ast.Num)
//			r_num, ok2 := right.(*ast.Num)
//			if ok && ok2 {
//				new_num := l_num.Values + r_num.Values
//				exp = ast.Num_new(new_num, e.LineNum)
//			} else {
//				exp = ast.Add_new(left, right, e.LineNum)
//			}
//		case *ast.And:
//			opt(e.Left)
//			left := exp
//			opt(e.Right)
//			right := exp
//			var left_v, right_v bool
//			left_isBool := true
//			right_isBool := true
//			if _, ok := left.(*ast.True); ok {
//				left_v = true
//			} else if _, ok := left.(*ast.False); ok {
//				left_v = false
//			} else {
//				left_isBool = false
//			}
//
//			if _, ok := right.(*ast.True); ok {
//				right_v = true
//			} else if _, ok := right.(*ast.False); ok {
//				right_v = false
//			} else {
//				right_isBool = false
//			}
//
//			if right_isBool && left_isBool {
//				if left_v && right_v {
//					exp = ast.True_new(e.LineNum)
//				} else {
//					exp = ast.False_new(e.LineNum)
//				}
//			} else {
//				exp = ast.And_new(left, right, e.LineNum)
//			}
//		case *ast.ArraySelect:
//			opt(e.Index)
//			index := exp
//			exp = ast.ArraySelect_new(e.X, index, e.LineNum)
//		case *ast.Call:
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
//			exp = e
//		case *ast.Id:
//			exp = e
//		case *ast.Length:
//			opt(e.X)
//			array := exp
//			exp = ast.Length_new(array, e.LineNum)
//		case *ast.Lt:
//			opt(e.Left)
//			left := exp
//			opt(e.Right)
//			right := exp
//			l_num, ok := left.(*ast.Num)
//			r_num, ok2 := right.(*ast.Num)
//			if ok && ok2 {
//				if l_num.Values < r_num.Values {
//					exp = ast.True_new(e.LineNum)
//				} else {
//					exp = ast.False_new(e.LineNum)
//				}
//			} else {
//				exp = ast.Lt_new(left, right, e.LineNum)
//			}
//		case *ast.NewIntArray:
//			opt(e.Size)
//			size := exp
//			exp = ast.NewIntArray_new(size, e.LineNum)
//		case *ast.NewObject:
//			exp = e
//		case *ast.Not:
//			opt(e.Values)
//			not_e := exp
//			if _, ok := not_e.(*ast.True); ok {
//				exp = ast.False_new(e.LineNum)
//			} else if _, ok := not_e.(*ast.False); ok {
//				exp = ast.True_new(e.LineNum)
//			} else {
//				exp = ast.Not_new(not_e, e.LineNum)
//			}
//		case *ast.Num:
//			exp = e
//		case *ast.Sub:
//			opt(e.Left)
//			left := exp
//			opt(e.Right)
//			right := exp
//			l_num, ok := left.(*ast.Num)
//			r_num, ok2 := right.(*ast.Num)
//			if ok && ok2 {
//				new_num := l_num.Values - r_num.Values
//				exp = ast.Num_new(new_num, e.LineNum)
//			} else {
//				exp = ast.Sub_new(left, right, e.LineNum)
//			}
//		case *ast.This:
//			exp = e
//		case *ast.Times:
//			opt(e.Left)
//			left := exp
//			opt(e.Right)
//			right := exp
//			l_num, ok := left.(*ast.Num)
//			r_num, ok2 := right.(*ast.Num)
//			if ok && ok2 {
//				new_num := l_num.Values * r_num.Values
//				exp = ast.Num_new(new_num, e.LineNum)
//			} else {
//				exp = ast.Times_new(left, right, e.LineNum)
//			}
//		case *ast.True:
//			exp = e
//		default:
//			panic("impossible")
//		}
//	}
//
//	opt_Stm := func(ss ast.Stm) {
//		switch s := ss.(type) {
//		case *ast.Assign:
//			opt(s.Values)
//			stm = ast.Assign_new(s.Names, exp, s.Tp, s.IsField, s.LineNum)
//		case *ast.AssignArray:
//			opt(s.Index)
//			index := exp
//			opt(s.Values)
//			e := exp
//			stm = ast.AssignArray_new(s.Names, index, e, s.Tp, s.IsField, s.LineNum)
//		case *ast.Block:
//			bstms := make([]ast.Stm, 0)
//			for _, _s := range s.Stms {
//				opt(_s)
//				bstms = append(bstms, stm)
//			}
//			stm = ast.Block_new(bstms, s.LineNum)
//		case *ast.Print:
//			opt(s.Values)
//			stm = ast.Print_new(exp, s.LineNum)
//		case *ast.If:
//			opt(s.Condition)
//			cond := exp
//			opt(s.Body)
//			thenn := stm
//			opt(s.Elsee)
//			elsee := stm
//			stm = ast.If_new(cond, thenn, elsee, s.LineNum)
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
//	opt_Method := func(mm ast.Y) {
//		switch m := mm.(type) {
//		case *ast.MethodSingle:
//			stms := make([]ast.Stm, 0)
//			for _, s := range m.Stms {
//				opt(s)
//				stms = append(stms, stm)
//			}
//			opt(m.RetExp)
//			ret_exp := exp
//			method = &ast.MethodSingle{m.RetType,
//				m.Names,
//				m.Formals,
//				m.Locals,
//				stms,
//				ret_exp}
//		default:
//			panic("impossible")
//		}
//	}
//
//	opt_MainClass := func(mc ast.MainClass) {
//		switch c := mc.(type) {
//		case *ast.MainClassSingle:
//			opt(c.Stms)
//			main_class = &ast.MainClassSingle{c.Names, c.Args, stm}
//		default:
//			panic("impossible")
//		}
//	}
//
//	opt_Class := func(cc ast.Class) {
//		switch c := cc.(type) {
//		case *ast.ClassSingle:
//			methods := make([]ast.Y, 0)
//			for _, m := range c.Methods {
//				opt(m)
//				methods = append(methods, method)
//			}
//			class = &ast.ClassSingle{c.Access, c.Names, c.Extends, c.Fields, methods}
//		default:
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
//		case ast.Type:
//			//no need
//		case ast.Field:
//			//no need
//		default:
//			panic("impossible")
//		}
//	}
//
//	var Ast *ast.FileSingle
//	if p, ok := prog.(*ast.FileSingle); ok {
//		opt(p.Mainclass)
//		classes := make([]ast.Class, 0)
//		for _, c := range p.Classes {
//			opt(c)
//			classes = append(classes, class)
//		}
//		Ast = &ast.FileSingle{p.Names, main_class, classes}
//	} else {
//		panic("impossible")
//	}
//
//	return Ast
//}
