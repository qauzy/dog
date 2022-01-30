package parser

import (
	"dog/ast"
	log "github.com/corgi-kx/logcustom"
)

//
//
// param: exp 声明变量的类型
func (this *Parser) parserDecl(exp ast.Exp) ast.Stm {
	log.Debugf("*******解析临时变量声明语句*******")
	id := this.current.Lexeme
	id = GetNewId(id)
	this.eatToken(TOKEN_ID)
	decl := ast.DeclStmt_new(nil, exp, nil, this.Linenum)
	this.currentStm = decl
	defer func() {
		this.currentStm = nil
	}()

	decl.Names = append(decl.Names, ast.NewIdent(GetNewId(id), this.Linenum))
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
											//this.ParseBug(fmt.Sprintf("====%v", c1.ArgsList))
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
		decl.Names = append(decl.Names, ast.NewIdent(GetNewId(id), this.Linenum))

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
