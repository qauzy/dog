package parser

import (
	"dog/ast"
	"dog/storage"
	log "github.com/corgi-kx/logcustom"
)

func (this *Parser) IsTypeToken() (b bool) {
	return this.current.Kind > TOKEN_TYPE_START && this.current.Kind < TOKEN_TYPE_END
}

func (this *Parser) ExtraToken() (b bool) {
	return this.current.Kind > TOKEN_EXTRA_START && this.current.Kind < TOKEN_EXTRA_END
}

func (this *Parser) parseType() ast.Exp {
	defer func() {
		log.Debugf("解析类型:%v", this.currentType)
	}()

	//处理final
	if this.current.Kind == TOKEN_FINAL {
		this.advance()
	}

	switch this.current.Kind {
	case TOKEN_CHAR:
		this.advance()
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.ArrayType{Ele: &ast.Char{}}
		} else {
			this.currentType = &ast.Char{}
		}

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
			this.currentType = &ast.ArrayType{Ele: &ast.Int{}}
		} else {
			this.currentType = &ast.Int{}
		}
	case TOKEN_OBJECT:
		this.eatToken(TOKEN_OBJECT)
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.ArrayType{Ele: &ast.ObjectType{}}
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
			this.currentType = &ast.ArrayType{Ele: &ast.Integer{}}
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
		this.currentType = &ast.ArrayType{Ele: &ast.Int{}}
	case TOKEN_SET:
		name := this.current.Lexeme
		this.eatToken(TOKEN_SET)
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			ele := this.parseType()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.SetType{name, ele, ast.TYPE_LIST}
		} else {
			this.assignType = &ast.ObjectType{ast.TYPE_OBJECT}
			this.currentType = &ast.SetType{name, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_LIST}
		}

	case TOKEN_HASHSET:
		//处理泛型
		name := this.current.Lexeme
		this.eatToken(TOKEN_HASHSET)
		this.eatToken(TOKEN_LT)
		ele := this.parseType()
		this.eatToken(TOKEN_ID)
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}

	case TOKEN_LIST:
		name := this.current.Lexeme
		this.eatToken(TOKEN_LIST)
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			if this.current.Kind == TOKEN_QUESTION {
				this.eatToken(TOKEN_QUESTION)
				//TODO 有继承规范
				if this.current.Kind == TOKEN_EXTENDS {
					this.eatToken(TOKEN_EXTENDS)
					ele := this.parseType()
					this.eatToken(TOKEN_GT)
					this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
					//没有继承规范
				} else {
					this.eatToken(TOKEN_GT)
					this.currentType = &ast.ListType{name, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_LIST}
				}

			} else {
				ele := this.parseType()
				this.eatToken(TOKEN_GT)
				this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
			}

		} else {
			this.currentType = &ast.ListType{name, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_LIST}
		}

	case TOKEN_ARRAYLIST:
		//处理泛型
		name := this.current.Lexeme
		this.eatToken(TOKEN_ARRAYLIST)
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			ele := this.parseType()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
		} else {
			this.currentType = &ast.ListType{name, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_LIST}
		}

	case TOKEN_MAP:
		name := this.current.Lexeme
		this.eatToken(TOKEN_MAP)
		//Map.Entry
		if this.current.Kind == TOKEN_DOT {
			this.advance()
			ast.SelectorExpr_new(ast.NewIdent(name, this.Linenum), this.current.Lexeme, this.Linenum)
			name = this.current.Lexeme
			this.eatToken(TOKEN_ID)
		}

		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			key := this.parseType()
			this.eatToken(TOKEN_COMMER)
			value := this.parseType()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.MapType{name, key, value, ast.TYPE_MAP}
		} else {
			this.currentType = &ast.MapType{name, &ast.String{ast.TYPE_STRING}, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_MAP}
		}
	case TOKEN_HASHMAP:
		name := this.current.Lexeme
		this.eatToken(TOKEN_HASHMAP)
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			key := this.parseType()
			if key == nil {
				this.eatToken(TOKEN_GT)
				this.currentType = &ast.MapType{name, &ast.String{ast.TYPE_STRING}, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_MAP}
			} else {
				this.eatToken(TOKEN_COMMER)
				value := this.parseType()
				this.eatToken(TOKEN_GT)
				this.currentType = &ast.MapType{name, key, value, ast.TYPE_MAP}
			}

		} else {
			this.currentType = &ast.MapType{name, &ast.String{ast.TYPE_STRING}, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_MAP}
		}

		//泛型
	case TOKEN_LT:
		this.eatToken(TOKEN_LT)
		tp := this.parseTypeList()
		this.eatToken(TOKEN_GT)
		//
		this.currentType = &ast.GenericType{ast.NewIdent("", this.Linenum), tp, ast.TYPE_GENERIC}
	default:
		//FIXME 类型可能带包名前缀
		name := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		if this.current.Kind == TOKEN_DOT {
			id := ast.NewIdent(name, this.Linenum)
			this.parseCallExp(id)
			name = this.current.Lexeme
		}

		//数组类型
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.ArrayType{&ast.ClassType{name, ast.TYPE_CLASS}, ast.TYPE_ARRAY}

			//多为数组
			for this.current.Kind == TOKEN_LBRACK {
				this.eatToken(TOKEN_LBRACK)
				this.eatToken(TOKEN_RBRACK)
				this.currentType = &ast.ArrayType{this.currentType, ast.TYPE_ARRAY}
			}

			return this.currentType
		}

		if this.current.Kind != TOKEN_LT {
			ttp, err := storage.FindByName(name)
			if err == nil && ast.KEY(ttp.Kind) == ast.INTERFACE_TYPE {
				this.currentType = &ast.InterfaceType{name, ast.TYPE_INTERFACE}
			} else {
				this.currentType = &ast.ClassType{name, ast.TYPE_CLASS}
			}

		} else {
			this.eatToken(TOKEN_LT)
			var tp []ast.Exp
			tp = this.parseTypeList()

			this.eatToken(TOKEN_GT)
			this.currentType = &ast.GenericType{ast.NewIdent(name, this.Linenum), tp, ast.TYPE_GENERIC}
		}
	}
	log.Debugf("解析类型：%v", this.currentType)
	return this.currentType
}
