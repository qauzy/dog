package parser

import "dog/ast"

//Assert.xxx
func (this *Parser) parseAssertExp() ast.Stm {
	this.eatToken(TOKEN_ASSERT)
	this.eatToken(TOKEN_DOT)
	id := this.current.Lexeme
	var opt string
	switch this.current.Kind{
	case TOKEN_IS_TRUE:
		opt = "true"
	case TOKEN_NOT_NULL:
		opt = 'nil'

	}
}
