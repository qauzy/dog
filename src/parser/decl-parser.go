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
	this.eatToken(TOKEN_ID)
	decl := ast.DeclStmt_new(nil, exp, nil, this.Linenum)
	this.currentStm = decl
	defer func() {
		this.currentStm = nil
	}()

	decl.Names = append(decl.Names, ast.NewIdent(id, this.Linenum))
	//有赋值语句
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

	//定义多个变量
	for this.current.Kind == TOKEN_COMMER {
		this.advance()
		id = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		decl.Names = append(decl.Names, ast.NewIdent(id, this.Linenum))

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
