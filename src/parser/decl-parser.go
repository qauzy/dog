package parser

import (
	"dog/ast"
	"dog/util"
	log "github.com/corgi-kx/logcustom"
)

//
//
// param: exp 声明变量的类型
func (this *Parser) parserDecl(exp ast.Exp) ast.Stm {
	log.Debugf("*******解析临时变量声明语句*******")
	id := this.current.Lexeme
	this.eatToken(TOKEN_ID)
	decl := ast.DeclStmt_new(nil, exp, nil, this.Linenum)
	this.currentStm = decl
	defer func() {
		this.currentStm = nil
	}()

	//记录本地变量
	if this.currentMethod != nil {
		f := &ast.FieldSingle{
			Access:  0,
			Tp:      exp,
			Name:    util.GetNewId(id),
			Static:  false,
			IsField: false,
			Value:   nil,
		}
		this.currentMethod.AddLocals(f)
	}

	decl.Names = append(decl.Names, ast.NewIdent(util.GetNewId(id), this.Linenum))
	//有赋值语句
	if this.current.Kind == TOKEN_ASSIGN {
		this.assignType = exp
		//临时变量类型
		log.Debugf("*******解析临时变量声明语句(有赋值语句)*******")
		this.eatToken(TOKEN_ASSIGN)
		exp := this.parseExp()
		var toany string
		if c0, ok := exp.(*ast.CallExpr); ok {
			if sl0, ok := c0.Callee.(*ast.SelectorExpr); ok {
				if sl0.Sel == "collect" {
					if len(c0.ArgsList) == 1 {
						if cc0, ok := c0.ArgsList[0].(*ast.CallExpr); ok {
							if ssl0, ok := cc0.Callee.(*ast.SelectorExpr); ok {
								if ssl0.Sel == "toSet" || ssl0.Sel == "toList" {
									toany = ssl0.Sel
								}

							}
						}

					}
					if c1, ok := sl0.X.(*ast.CallExpr); ok {
						if sl1, ok := c1.Callee.(*ast.SelectorExpr); ok {
							if sl1.Sel == "oMap" || sl1.Sel == "map" || sl1.Sel == "Map" || sl1.Sel == "sorted" {
								if c2, ok := sl1.X.(*ast.CallExpr); ok {
									if sl2, ok := c2.Callee.(*ast.SelectorExpr); ok {
										if sl2.Sel == "stream" {
											if len(c1.ArgsList) == 1 {
												log.Debugf("*******解析map语句*******")
												decl.SetExtra(ast.MapStm_new(decl.Names[0], sl2.X, c1.ArgsList[0], toany, this.Linenum))
											}
										}
									}
								}

							} else {
								this.ParseBug(sl1.Sel)
							}
						}
					}
				}
			}
		}
		decl.Values = append(decl.Values, exp)
	}

	//定义多个变量
	for this.current.Kind == TOKEN_COMMER {
		this.advance()
		id = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		//记录本地变量
		if this.currentMethod != nil {
			f := &ast.FieldSingle{
				Access:  0,
				Tp:      exp,
				Name:    util.GetNewId(id),
				Static:  false,
				IsField: false,
				Value:   nil,
			}
			this.currentMethod.AddLocals(f)
		}
		decl.Names = append(decl.Names, ast.NewIdent(util.GetNewId(id), this.Linenum))

		if this.current.Kind == TOKEN_ASSIGN {
			//临时变量类型
			log.Debugf("*******解析临时变量声明语句(有赋值语句)*******")
			this.eatToken(TOKEN_ASSIGN)
			exp := this.parseExp()
			//三元表达式
			if _, ok := exp.(*ast.Question); ok {
				decl.SetTriple()
			}
			decl.Values = append(decl.Values, exp)
		}

	}
	this.eatToken(TOKEN_SEMI)
	return decl
}

//检查List的map操作
func (this *Parser) CheckSelectorExprs(exp ast.Exp, expext *string, call string, mp ast.MapStm) (b bool) {
	switch e := exp.(type) {
	case *ast.SelectorExpr:
		switch e.Sel {
		case "stream":
			if *expext == "stream" {
				return this.CheckSelectorExprs(e.X, expext, "stream", mp)
			} else {
				return false
			}
		case "map":
			return this.CheckSelectorExprs(e.X, expext, "map", mp)
		case "sorted":
			return this.CheckSelectorExprs(e.X, expext, "sorted", mp)
		case "collect":
			return this.CheckSelectorExprs(e.X, expext, "", mp)
		}
		return
	case *ast.CallExpr:
		switch call {
		case "stream":
			this.CheckSelectorExprs(e.Callee, expext, "", mp)
			return
		case "map":
			if len(e.ArgsList) == 1 {
				mp.Ele = e.ArgsList[0]
				return this.CheckSelectorExprs(e.Callee, expext, "", mp)
			}
		case "sorted":
			return this.CheckSelectorExprs(e.Callee, expext, "", mp)
		case "collect":
			if len(e.ArgsList) == 1 {
				if cc0, ok := e.ArgsList[0].(*ast.CallExpr); ok {
					if ssl0, ok := cc0.Callee.(*ast.SelectorExpr); ok {
						if ssl0.Sel == "toSet" || ssl0.Sel == "toList" {
							mp.ToAny = ssl0.Sel
							return this.CheckSelectorExprs(e.Callee, expext, "", mp)
						}
					}
				}
			}
			return false
		}
		return this.CheckSelectorExprs(e.Callee, expext, "", mp)
	case *ast.Ident:
		if this.currentMethod != nil && this.currentMethod.GetLocals(e.Name) != nil {
			lo := this.currentMethod.GetLocals(e.Name)
			if _, ok := lo.GetDecType().(*ast.ListType); ok {
				mp.List = e
				*expext = "stream"
				return true
			}
		}
	}
	return
}
