package parser

import "dog/ast"

// 解析类
//
// return:
func (this *Parser) parseClassDecl() (cl ast.Class) {
	var id, extends string

	//类访问权限修饰符
	var access int
	if this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PROTECTED {
		access = this.current.Kind
		this.advance()
	}
	//枚举类型
	if this.current.Kind == TOKEN_ENUM {

		return this.parseEnumDecl(access)
	}

	if this.current.Kind == TOKEN_INTERFACE {
		return this.parseInterfaceDecl(access)
	}

	//处理abstract
	if this.current.Kind == TOKEN_ABSTRACT {
		this.advance()
	}
	this.eatToken(TOKEN_CLASS)
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
		this.parseType()
	}

	//处理implements
	if this.current.Kind == TOKEN_IMPLEMENTS {
		this.eatToken(TOKEN_IMPLEMENTS)
		extends = this.current.Lexeme
		this.eatToken(TOKEN_ID)
	}

	this.eatToken(TOKEN_LBRACE)
	classSingle := ast.NewClassSingle(access, id, extends, ast.CLASS_TYPE)
	this.currentClass = classSingle
	defer func() {
		this.currentClass = nil
	}()
	this.parseClassContext(classSingle)
	this.eatToken(TOKEN_RBRACE)
	return classSingle
}

// 解析类组
//
// return:
func (this *Parser) parseClassDecls() {

	for this.current.Kind == TOKEN_CLASS || this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PROTECTED || this.current.Kind == TOKEN_COMMENT {
		var comment string
		if this.current.Kind == TOKEN_COMMENT {
			comment = ""
			//处理注释
			for this.current.Kind == TOKEN_COMMENT {
				comment += "\n" + this.current.Lexeme
				this.advance()
			}
			if this.current.Kind == TOKEN_EOF || (this.current.Kind != TOKEN_PRIVATE && this.current.Kind != TOKEN_PUBLIC && this.current.Kind != TOKEN_PROTECTED) {
				return
			}

		}
		if this.currentFile != nil {
			this.currentFile.AddClass(this.parseClassDecl())
		} else {
			panic("currentFile is nil")
		}

	}
	return
}
