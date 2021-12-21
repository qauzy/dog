package parser

import (
	"dog/ast"
)

//
func (this *Parser) parseInterfaceDecl(access int) (cl ast.Class) {
	var id, extends string

	//类访问权限修饰符
	this.eatToken(TOKEN_INTERFACE)
	id = this.current.Lexeme
	this.eatToken(TOKEN_ID)

	//FIXME 泛型忽略
	if this.current.Kind == TOKEN_LT {
		this.eatToken(TOKEN_LT)
		this.eatToken(TOKEN_ID)
		this.eatToken(TOKEN_GT)
	}

	//处理extends
	if this.current.Kind == TOKEN_EXTENDS {
		this.eatToken(TOKEN_EXTENDS)
		extends = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			this.eatToken(TOKEN_ID)
			for this.current.Kind == TOKEN_COMMER {
				this.eatToken(TOKEN_COMMER)
				this.advance()
			}
			this.eatToken(TOKEN_GT)
		}
		for this.current.Kind == TOKEN_COMMER {
			this.eatToken(TOKEN_COMMER)
			this.eatToken(TOKEN_ID)
			if this.current.Kind == TOKEN_LT {
				this.eatToken(TOKEN_LT)
				this.eatToken(TOKEN_ID)
				for this.current.Kind == TOKEN_COMMER {
					this.eatToken(TOKEN_COMMER)
					this.advance()
				}
				this.eatToken(TOKEN_GT)
			}

		}

	}

	this.eatToken(TOKEN_LBRACE)
	classSingle := ast.NewClassSingle(access, id, extends, ast.INTERFACE_TYPE)
	this.currentClass = classSingle
	defer func() {
		this.currentClass = nil
	}()

	for this.TypeToken() ||
		this.current.Kind == TOKEN_COMMENT ||
		this.current.Kind == TOKEN_QUERY || //解析jpa的  @Query 注解,
		this.current.Kind == TOKEN_PUBLIC ||
		this.current.Kind == TOKEN_ID {

		for this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}
		var stms []ast.Stm
		if this.current.Kind == TOKEN_QUERY {
			q := this.parseQuery()
			if q != nil {
				stms = append(stms, q)
				stms = append(stms, ast.Return_new(nil, this.Linenum))

			}
		}

		if stms == nil {
			stms = append(stms, ast.Return_new(nil, this.Linenum))
		}

		if this.current.Kind == TOKEN_PUBLIC {
			this.advance()
		}

		//空接口
		if this.current.Kind == TOKEN_RBRACE {
			this.eatToken(TOKEN_RBRACE)
			return classSingle
		}
		tp := this.parseType()
		id = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		this.eatToken(TOKEN_LPAREN)
		formals := this.parseFormalList(false)
		this.eatToken(TOKEN_RPAREN)
		this.eatToken(TOKEN_SEMI)
		this.currentMethod = ast.NewMethodSingle(tp, id, formals, stms, false, false, false, "")
		classSingle.AddMethod(this.currentMethod)
	}
	this.eatToken(TOKEN_RBRACE)
	return classSingle
}
