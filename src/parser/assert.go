package parser

import "dog/ast"

//Assert.xxx
func (this *Parser) parseAssertExp() ast.Stm {
	if this.current.Kind == TOKEN_ASSERT {
		this.eatToken(TOKEN_ASSERT)
		this.eatToken(TOKEN_DOT)
	}
	opt := this.current.Lexeme
	this.advance()
	this.eatToken(TOKEN_LPAREN)
	cond := this.parseExp()
	this.eatToken(TOKEN_COMMER)
	exp := this.parseExp()
	this.eatToken(TOKEN_RPAREN)
	this.eatToken(TOKEN_SEMI)
	return ast.Assert_new(cond, exp, opt, this.Linenum)
}
