package codegen_go

import (
	"dog/util"
	"fmt"
	"go/ast"
	"go/token"
)

//实现stream处理
func (this *Translation) OptimitcStreamStm(src ast.Stmt) (dst ast.Stmt) {

	switch stm := src.(type) {
	case *ast.AssignStmt:
		lhs := stm.Lhs[0]
		rhs := stm.Rhs[0]
		fk := &FakeBlock{}
		stmt := this.GetOptismicValue(fk, lhs, rhs)
		if stmt != nil {
			return stmt
		}
	case *ast.DeclStmt:
		//变量声明语句
		if gd, ok := stm.Decl.(*ast.GenDecl); ok {
			if gd.Tok == token.VAR {
				fk := &FakeBlock{}
				fk.List = append(fk.List, stm)
				for _, vv := range gd.Specs {
					if vs, ok := vv.(*ast.ValueSpec); ok {
						//只有单变量声明可以
						if len(vs.Names) == 1 && len(vs.Values) == 1 {
							lhs := vs.Names[0]
							rhs := vs.Values[0]
							rangeStmt := this.GetOptismicValue(fk, lhs, rhs)
							if rangeStmt != nil {
								vs.Values = nil
								return fk
							}

						}

					}

				}
				return fk
			}
		}
	case *ast.ExprStmt:
		if call0, ok := stm.X.(*ast.CallExpr); ok {
			var arg string
			if id, ok := call0.Fun.(*ast.Ident); ok {
				if len(id.Name) > 4 {
					arg = "v" + util.Capitalize(id.Name[len(id.Name)-5:])
				} else {
					arg = "v" + util.Capitalize(id.Name)
				}
			} else if sel, ok := call0.Fun.(*ast.SelectorExpr); ok {
				if len(sel.Sel.Name) > 4 {
					arg = "v" + util.Capitalize(sel.Sel.Name[len(sel.Sel.Name)-5:])
				} else {
					arg = "v" + util.Capitalize(sel.Sel.Name)
				}

			}
			if len(call0.Args) == 1 {
				lhs := ast.NewIdent(arg)
				fk := &FakeBlock{}
				rangeStmt := this.GetOptismicValue(fk, lhs, call0.Args[0])
				if rangeStmt != nil {
					call0.Args[0] = lhs
					fk.List = append(fk.List, stm)
					return fk
				}

			}

		}
	}

	return src
}

// xxx.orElse(x);
// xxx.get()
// xxx.reduce(xxx, xxx::xx);
func (this *Translation) GetOptismicValue(fk *FakeBlock, lhs ast.Expr, rhs ast.Expr) ast.Stmt {
	if call, ok := rhs.(*ast.CallExpr); ok {
		// xxx.stream().min(xxx).orElse(x);
		if orElse, ok := call.Fun.(*ast.SelectorExpr); ok && (orElse.Sel.Name == "OrElse") && len(call.Args) == 1 {
			//orElse
			as := &ast.AssignStmt{
				Lhs:    []ast.Expr{lhs},
				TokPos: 0,
				Tok:    token.ASSIGN,
				Rhs:    []ast.Expr{call.Args[0]},
			}
			fk.List = append(fk.List, as)
			stmt := this.GetStreamOpt(lhs, orElse.X)
			if stmt != nil {
				fk.List = append(fk.List, stmt)
				return fk
			}
			//stream().map(xxx::xxx).collect(Collectors.toList());
			//Stream.of(xxxx).map(xxx::xxx).collect(Collectors.toList());
		} else if collect, ok := call.Fun.(*ast.SelectorExpr); ok && (collect.Sel.Name == "Collect") && len(call.Args) == 1 {
			stmt := this.GetStreamOpt(lhs, collect.X)
			if stmt != nil {
				fk.List = append(fk.List, stmt)
				return fk
			}
		} else if get, ok := call.Fun.(*ast.SelectorExpr); ok && (get.Sel.Name == "Get") && len(call.Args) == 0 {
			if call2, ok := get.X.(*ast.CallExpr); ok {
				// xxx.stream().findFirst().get()
				if findFirst, ok := call2.Fun.(*ast.SelectorExpr); ok && (findFirst.Sel.Name == "FindFirst") && len(call2.Args) == 0 {

				}
			}

		} else if get, ok := call.Fun.(*ast.SelectorExpr); ok && (get.Sel.Name == "Count") && len(call.Args) == 0 {
			if call2, ok := get.X.(*ast.CallExpr); ok {
				// xxx.stream().findFirst().get()
				if filter, ok := call2.Fun.(*ast.SelectorExpr); ok && (filter.Sel.Name == "Filter") && len(call2.Args) == 1 {

				}
			}

			//xxx.stream().reduce(xxx, xxx::xx));
		} else if reduce, ok := call.Fun.(*ast.SelectorExpr); ok && (reduce.Sel.Name == "Reduce") && len(call.Args) == 2 {
			if call2, ok := reduce.X.(*ast.CallExpr); ok {
				if stream, ok := call2.Fun.(*ast.SelectorExpr); ok && (stream.Sel.Name == "Stream") && len(call2.Args) == 0 {
					var arg string
					if vv, ok := lhs.(*ast.Ident); ok {
						arg = vv.Name + "List"
					} else if vv, ok := lhs.(*ast.SelectorExpr); ok {
						arg = vv.Sel.Name + "List"
					}

					stmt := this.GetOptismicValue(fk, ast.NewIdent(arg), stream.X)
					if stmt != nil {
						if bin, ok := call.Args[1].(*ast.BinaryExpr); ok {
							cc := &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   lhs,
									Sel: bin.Y.(*ast.Ident),
								},
								Lparen:   0,
								Args:     nil,
								Ellipsis: 0,
								Rparen:   0,
							}
							cc.Args = append(cc.Args, ast.NewIdent("val"))
							as := &ast.AssignStmt{
								Lhs:    []ast.Expr{lhs},
								TokPos: 0,
								Tok:    token.ASSIGN,
								Rhs:    []ast.Expr{cc},
							}

							rstm := &ast.RangeStmt{
								For:    0,
								Key:    ast.NewIdent("_"),
								Value:  ast.NewIdent("val"),
								TokPos: 0,
								Tok:    token.DEFINE,
								X:      ast.NewIdent(arg),
								Body: &ast.BlockStmt{
									List: []ast.Stmt{as},
								},
							}
							fk.List = append(fk.List, rstm)
							return fk
						}

					}

				}
			}

		} else if parseObject, ok := call.Fun.(*ast.SelectorExpr); ok && (parseObject.Sel.Name == "ParseObject") && len(call.Args) == 2 {
			if vvvv, ok := parseObject.X.(*ast.Ident); ok && (vvvv.Name == "JSON") {
				//转换json解析
				call.Args = []ast.Expr{call.Args[0], lhs}
				call.Fun = ast.NewIdent("mdata.Cjson.Unmarshal")
				as := &ast.AssignStmt{
					Lhs:    []ast.Expr{ast.NewIdent("err")},
					TokPos: 0,
					Tok:    token.ASSIGN,
					Rhs:    []ast.Expr{rhs},
				}
				fk.List = append(fk.List, as)
				fk.List = append(fk.List, this.GetErrReturn())
				return fk

			}
		}

	}

	return nil
}

//根据不通操作，获取转换逻辑
func (this *Translation) GetStreamOpt(lhs ast.Expr, exp ast.Expr) ast.Stmt {
	if call2, ok := exp.(*ast.CallExpr); ok {
		if m, ok := call2.Fun.(*ast.SelectorExpr); ok && (m.Sel.Name == "Min" || m.Sel.Name == "Max") && len(call2.Args) == 1 {
			return this.GetStreamMinMax(lhs, call2)
		} else if m, ok := call2.Fun.(*ast.SelectorExpr); ok && (m.Sel.Name == "Map") && len(call2.Args) == 1 {
			return this.GetStreamMap(lhs, call2)
		}

	}
	return nil
}
func (this *Translation) GetStreamMinMax(lhs ast.Expr, exp ast.Expr) ast.Stmt {
	if call2, ok := exp.(*ast.CallExpr); ok {
		//xxx.stream().min(xxx).orElse(x)
		if m, ok := call2.Fun.(*ast.SelectorExpr); ok && (m.Sel.Name == "Min" || m.Sel.Name == "Max") && len(call2.Args) == 1 {
			var op token.Token
			if m.Sel.Name == "Min" {
				op = token.LSS
			} else {
				op = token.GTR
			}
			list := this.GetStreamList(m.X)

			if list != nil {
				ifStm := &ast.IfStmt{
					If:   0,
					Init: nil,
					Cond: &ast.BinaryExpr{
						X:     ast.NewIdent("val"),
						OpPos: 0,
						Op:    op,
						Y:     lhs,
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{&ast.AssignStmt{
							Lhs:    []ast.Expr{lhs},
							TokPos: 0,
							Tok:    token.ASSIGN,
							Rhs:    []ast.Expr{ast.NewIdent("val")},
						}},
					},
					Else: nil,
				}

				stmt := &ast.RangeStmt{
					For:    0,
					Key:    ast.NewIdent("_"),
					Value:  ast.NewIdent("val"),
					TokPos: 0,
					Tok:    token.DEFINE,
					X:      list,
					Body: &ast.BlockStmt{
						List: []ast.Stmt{ifStm},
					},
				}

				return stmt
			}

		} else if findFirst, ok := call2.Fun.(*ast.SelectorExpr); ok && (findFirst.Sel.Name == "FindFirst") && len(call2.Args) == 0 {

		}
	}
	return nil
}

func (this *Translation) GetStreamMap(lhs ast.Expr, exp ast.Expr) ast.Stmt {
	if call2, ok := exp.(*ast.CallExpr); ok {
		if oMap, ok := call2.Fun.(*ast.SelectorExpr); ok && (oMap.Sel.Name == "Map") && len(call2.Args) == 1 {
			list := this.GetStreamList(oMap.X)
			if list != nil {
				//将元素转换成新形式
				trS := this.TransMapMethod(call2.Args[0], ast.NewIdent("val"), ast.NewIdent("vvo"))

				appendC := &ast.CallExpr{
					Fun:      ast.NewIdent("append"),
					Lparen:   0,
					Args:     nil,
					Ellipsis: 0,
					Rparen:   0,
				}
				appendC.Args = append(appendC.Args, lhs)
				appendC.Args = append(appendC.Args, ast.NewIdent("vvo"))
				as := &ast.AssignStmt{
					Lhs:    []ast.Expr{lhs},
					TokPos: 0,
					Tok:    token.ASSIGN,
					Rhs:    []ast.Expr{appendC},
				}

				rangeStmt := &ast.RangeStmt{
					For:    0,
					Key:    ast.NewIdent("_"),
					Value:  ast.NewIdent("val"),
					TokPos: 0,
					Tok:    token.DEFINE,
					X:      list,
					Body: &ast.BlockStmt{
						List: []ast.Stmt{trS, as},
					},
				}

				return rangeStmt

			}
		}
	}
	return nil
}

//获取待处理的stream对象
func (this *Translation) GetStreamList(exp ast.Expr) ast.Expr {
	if call3, ok := exp.(*ast.CallExpr); ok {
		//stream()
		if len(call3.Args) == 0 {
			if stream, ok := call3.Fun.(*ast.SelectorExpr); ok && (stream.Sel.Name == "Stream") {
				return stream.X

			}
		} else if len(call3.Args) == 1 {
			if streamf, ok := call3.Fun.(*ast.SelectorExpr); ok && (streamf.Sel.Name == "Of") {
				return call3.Args[0]
			}
		}
	}
	return nil
}

//处理stram().map元素转换
func (this *Translation) TransMapMethod(method ast.Expr, ele ast.Expr, newele ast.Expr) ast.Stmt {
	if mp, ok := method.(*ast.BinaryExpr); ok {
		if id, ok := mp.Y.(*ast.Ident); ok && id.Name == "ParseInt" {
			call := &ast.CallExpr{
				Fun:      ast.NewIdent("strconv.Atoi"),
				Lparen:   0,
				Args:     nil,
				Ellipsis: 0,
				Rparen:   0,
			}
			call.Args = append(call.Args, ele)
			as := &ast.AssignStmt{
				Lhs:    []ast.Expr{newele, ast.NewIdent("_")},
				TokPos: 0,
				Tok:    token.DEFINE,
				Rhs:    []ast.Expr{call},
			}
			return as
		} else if id, ok := mp.Y.(*ast.Ident); ok && id.Name == "New" {
			if bid, ok := mp.X.(*ast.Ident); ok && bid.Name == "BigDecimal" {
				call := &ast.CallExpr{
					Fun:      ast.NewIdent("decimal.NewFromString"),
					Lparen:   0,
					Args:     nil,
					Ellipsis: 0,
					Rparen:   0,
				}
				call.Args = append(call.Args, ele)
				as := &ast.AssignStmt{
					Lhs:    []ast.Expr{newele, ast.NewIdent("_")},
					TokPos: 0,
					Tok:    token.DEFINE,
					Rhs:    []ast.Expr{call},
				}
				return as
			}

		}

	}
	panic(fmt.Sprintln("不支持转换函数:%v", method))
}
