package parser

import (
	"dog/ast"
	"strings"
)

//解析 @Query 注解
func (this *Parser) parseQuery() (stm ast.Stm) {
	this.eatToken(TOKEN_QUERY)
	this.eatToken(TOKEN_LPAREN)
	var q string
	var nativeQuery bool
	//获得sql 有多重形式
	//1 value = "xxx",nativeQuery = xxx
	id := this.current.Lexeme
	this.advance()
	if this.current.Kind == TOKEN_ASSIGN {
		this.eatToken(TOKEN_ASSIGN)
		if id == "value" {
			q += this.current.Lexeme
		} else if id == "nativeQuery" {
			if this.current.Kind == TOKEN_TRUE {
				nativeQuery = true
			}
		}

		this.advance()
		for this.current.Kind == TOKEN_ADD {
			this.advance()
			if id == "value" {
				q = strings.TrimSuffix(q, "\"") + strings.TrimPrefix(this.current.Lexeme, "\"")
			} else if id == "nativeQuery" {
				if this.current.Kind == TOKEN_TRUE {
					nativeQuery = true
				}
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
					q = strings.TrimSuffix(q, "\"") + strings.TrimPrefix(this.current.Lexeme, "\"")
				} else if id == "nativeQuery" {
					if this.current.Kind == TOKEN_TRUE {
						nativeQuery = true
					}
				}
				this.advance()
			}
		}

		// 2 xxx + xxxx
	} else if this.current.Kind == TOKEN_ADD {
		q = strings.TrimSuffix(q, "\"") + strings.TrimPrefix(id, "\"")
		for this.current.Kind == TOKEN_ADD {
			this.advance()
			q = strings.TrimSuffix(q, "\"") + strings.TrimPrefix(this.current.Lexeme, "\"")
			this.advance()
		}
		//3 "xxxx"
	} else {
		q += id
	}

	this.eatToken(TOKEN_RPAREN)
	stm = ast.Query_new(q, nativeQuery, this.Linenum)
	return
}
