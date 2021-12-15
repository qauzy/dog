package parser

import (
	"dog/ast"
	log "github.com/corgi-kx/logcustom"
)

func (this *Parser) parseTypeV2() ast.Exp {
	switch this.current.Kind {
	case TOKEN_FLOAT:
		fallthrough
	case TOKEN_DOUBLE:
		this.advance()
		this.currentType = &ast.Float{}
	case TOKEN_SHORT:
		fallthrough
	case TOKEN_INT:
		this.advance()
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.IntArray{ast.TYPE_INTARRAY}
		} else {
			this.currentType = &ast.Int{}
		}
	case TOKEN_OBJECT:
		this.eatToken(TOKEN_OBJECT)
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.ObjectArray{ast.TYPE_OBJECTARRAY}
		} else {
			this.currentType = &ast.ObjectType{ast.TYPE_OBJECT}
		}
	case TOKEN_LONG:
		fallthrough
	case TOKEN_INTEGER:
		this.advance()
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.IntArray{ast.TYPE_INTARRAY}
		} else {
			this.currentType = &ast.Integer{ast.TYPE_INT}
		}
	case TOKEN_BYTE:
		this.advance()
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.ByteArray{ast.TYPE_BYTEARRAY}
		} else {
			this.currentType = &ast.Byte{ast.TYPE_BYTE}
		}
	case TOKEN_VOID:
		this.eatToken(TOKEN_VOID)
		this.currentType = &ast.Void{ast.TYPE_VOID}
	case TOKEN_BOOLEAN:
		this.eatToken(TOKEN_BOOLEAN)
		this.currentType = &ast.Boolean{ast.TYPE_BOOLEAN}
	case TOKEN_DATE:
		this.eatToken(TOKEN_DATE)
		this.currentType = &ast.Date{ast.TYPE_DATE}
	case TOKEN_STRING:
		this.eatToken(TOKEN_STRING)
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.StringArray{ast.TYPE_STRINGARRAY}
		} else {
			this.currentType = &ast.String{ast.TYPE_STRING}
		}
	case TOKEN_LBRACK:
		this.eatToken(TOKEN_LBRACK)
		this.eatToken(TOKEN_RBRACK)
		this.eatToken(TOKEN_INT)
		this.currentType = &ast.IntArray{ast.TYPE_INTARRAY}
	case TOKEN_SET:
		name := this.current.Lexeme
		this.eatToken(TOKEN_SET)
		this.eatToken(TOKEN_LT)
		ele := this.parseNotExp()
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
	case TOKEN_HASHSET:
		//处理泛型
		name := this.current.Lexeme
		this.eatToken(TOKEN_HASHSET)
		this.eatToken(TOKEN_LT)
		ele := this.parseNotExp()
		this.eatToken(TOKEN_ID)
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}

	case TOKEN_LIST:
		name := this.current.Lexeme
		this.eatToken(TOKEN_LIST)
		this.eatToken(TOKEN_LT)
		ele := this.parseNotExp()
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
	case TOKEN_ARRAYLIST:
		//处理泛型
		name := this.current.Lexeme
		this.eatToken(TOKEN_ARRAYLIST)
		this.eatToken(TOKEN_LT)
		ele := this.parseNotExp()
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
	case TOKEN_MAP:
		name := this.current.Lexeme
		this.eatToken(TOKEN_MAP)
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			key := this.parseType()
			this.eatToken(TOKEN_COMMER)
			value := this.parseType()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.HashType{name, key, value, ast.TYPE_MAP}
		} else {
			this.currentType = &ast.HashType{name, &ast.String{ast.TYPE_STRING}, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_MAP}
		}
	case TOKEN_HASHMAP:
		name := this.current.Lexeme
		this.eatToken(TOKEN_HASHMAP)
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			key := this.parseType()
			if key == nil {
				this.eatToken(TOKEN_GT)
				this.currentType = &ast.HashType{name, &ast.String{ast.TYPE_STRING}, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_MAP}
			} else {
				this.eatToken(TOKEN_COMMER)
				value := this.parseType()
				this.eatToken(TOKEN_GT)
				this.currentType = &ast.HashType{name, key, value, ast.TYPE_MAP}
			}

		} else {
			this.currentType = &ast.HashType{name, &ast.String{ast.TYPE_STRING}, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_MAP}
		}

	default:
		name := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		if this.current.Kind != TOKEN_LT {
			this.currentType = &ast.ClassType{name, ast.TYPE_CLASS}
		} else {
			this.eatToken(TOKEN_LT)
			tp := this.parseTypeList()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.GenericType{name, tp, ast.TYPE_GENERIC}
		}
	}
	log.Debugf("解析类型:%s", this.currentType)
	return this.currentType
}
