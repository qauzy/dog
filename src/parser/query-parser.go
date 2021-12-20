package parser

import "dog/ast"

//解析 @Query 注解
func (this *Parser) parseQuery() (q ast.Exp) {
	this.eatToken(TOKEN_QUERY)
	this.eatToken(TOKEN_LPAREN)

	//获得sql 有多重形式
	//1 value = "xxx",nativeQuery = xxx
	id := this.current.Lexeme
	var sql = ""
	this.advance()
	if this.current.Kind == TOKEN_ASSIGN {
		this.eatToken(TOKEN_ASSIGN)
		if id == "value" {
			sql += this.current.Lexeme
		}
		this.advance()
		for this.current.Kind == TOKEN_ADD {
			this.advance()
			if id == "value" {
				sql += this.current.Lexeme
			}
			this.advance()

		}
		for this.current.Kind == TOKEN_COMMER {
			this.eatToken(TOKEN_COMMER)
			id = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			this.eatToken(TOKEN_ASSIGN)
			this.advance()
			for this.current.Kind == TOKEN_ADD {
				this.advance()
				if id == "value" {
					sql += this.current.Lexeme
				}
				this.advance()
			}
		}

		// 2 xxx + xxxx
	} else if this.current.Kind == TOKEN_ADD {
		sql += id
		for this.current.Kind == TOKEN_ADD {
			this.advance()
			sql += this.current.Lexeme
			this.advance()
		}
		//3 "xxxx"
	} else {
		sql += id
	}

	this.eatToken(TOKEN_RPAREN)

	return
}
