package parser

import (
	"dog/ast"
	"dog/util"
	log "github.com/corgi-kx/logcustom"
)

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
		if this.current.Kind == TOKEN_EXTENDS {
			this.advance()
			this.parseType()
		}
		for this.current.Kind == TOKEN_COMMER {
			this.advance()
			this.eatToken(TOKEN_ID)
			if this.current.Kind == TOKEN_EXTENDS {
				this.advance()
				this.parseType()
			}

		}

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
		//FIXME 泛型忽略
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			this.eatToken(TOKEN_ID)
			for this.current.Kind == TOKEN_COMMER {
				this.eatToken(TOKEN_COMMER)
				this.eatToken(TOKEN_ID)
			}
			this.eatToken(TOKEN_GT)
		}
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

// 解析类上下文
//
// return:
func (this *Parser) parseClassContext(classSingle *ast.ClassSingle) {

	//每次循环解析一个成员变量或一个成员函数
	for this.IsTypeToken() ||
		this.current.Kind == TOKEN_ID ||
		this.current.Kind == TOKEN_PRIVATE ||
		this.current.Kind == TOKEN_PUBLIC ||
		this.current.Kind == TOKEN_PROTECTED ||
		this.current.Kind == TOKEN_FINAL ||
		this.current.Kind == TOKEN_COMMENT ||
		this.current.Kind == TOKEN_STATIC {
		var comment string
		//处理注释
		if this.current.Kind == TOKEN_COMMENT {
			comment = ""
			for this.current.Kind == TOKEN_COMMENT {
				comment += "\n" + this.current.Lexeme
				this.advance()
			}
			if this.current.Kind == TOKEN_EOF || (this.current.Kind != TOKEN_PRIVATE && this.current.Kind != TOKEN_PUBLIC && this.current.Kind != TOKEN_PROTECTED) {
				return
			}
		}
		var tmp ast.FieldSingle
		var IsConstruct = false
		var IsStatic = false
		var IsBlock = false
		var prefix = false

		//访问修饰符 [其他修饰符] 类型 变量名 = 值;
		//处理 访问修饰符
		if this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PROTECTED {
			//1 扫描访问修饰符
			tmp.Access = this.current.Kind
			this.advance()
		} else {
			tmp.Access = TOKEN_DEFAULT
		}
		//处理成员修饰符
		for this.current.Kind == TOKEN_STATIC ||
			this.current.Kind == TOKEN_TRANSIENT ||
			this.current.Kind == TOKEN_SYNCHRONIZED ||
			this.current.Kind == TOKEN_FINAL {

			//处理 其他修饰符(忽略)
			if this.current.Kind == TOKEN_STATIC {
				IsStatic = true
				prefix = true
				this.eatToken(TOKEN_STATIC)
				if this.current.Kind == TOKEN_LBRACE {
					IsBlock = true
				}
			}

			if this.current.Kind == TOKEN_FINAL {
				prefix = true
				this.eatToken(TOKEN_FINAL)
			}

			if this.current.Kind == TOKEN_TRANSIENT {
				prefix = true
				this.eatToken(TOKEN_TRANSIENT)
			}

			if this.current.Kind == TOKEN_SYNCHRONIZED {
				prefix = true
				this.eatToken(TOKEN_SYNCHRONIZED)
			}

		}

		//类静态语句
		if (this.current.Kind == TOKEN_LBRACE) && IsBlock {

			classSingle.AddMethod(this.parseMemberStatic(comment))

		} else {

			id := this.current.Lexeme
			//处理类构造函数
			if this.currentClass.GetName() == id && prefix == false {
				this.advance()
				if this.current.Kind == TOKEN_LPAREN {
					log.Infof("处理构造函数-->%v", this.current.Lexeme)
					IsConstruct = true
					tmp.Tp = &ast.Void{ast.TYPE_VOID}
					//变量/函数名
					tmp.Name = "New" + id
				} else {
					tmp.Tp = &ast.ClassType{
						Name:     id,
						TypeKind: ast.TYPE_CLASS,
					}

					this.assignType = tmp.Tp
					//变量/函数名
					tmp.Name = this.current.Lexeme
					this.eatToken(TOKEN_ID)
				}

			} else {
				//类型
				tmp.Tp = this.parseType()
				this.assignType = tmp.Tp
				//变量/函数名 --> 转为开头大写
				tmp.Name = util.Capitalize(this.current.Lexeme)
				this.eatToken(TOKEN_ID)
			}

			//成员方法
			if this.current.Kind == TOKEN_LPAREN {
				classSingle.AddMethod(this.parseMemberMethod(&tmp, IsConstruct, IsStatic, comment))
				//成员变量

			} else {
				if IsStatic {
					this.currentFile.AddField(this.parseMemberVarDecl(&tmp, IsStatic))
				} else {
					classSingle.AddField(this.parseMemberVarDecl(&tmp, IsStatic))
				}

			}
		}

	}
	return
}

// 解析成员函数
//
// param: dec
// param: IsConstruct
// return:
func (this *Parser) parseMemberMethod(dec *ast.FieldSingle, IsConstruct bool, IsStatic bool, comment string) (meth ast.Method) {
	var IsThrows bool
	//左括号
	this.eatToken(TOKEN_LPAREN)
	//解析参数
	formals := this.parseFormalList(false)

	this.currentMethod = ast.NewMethodSingle(dec.Tp, dec.Name, formals, nil, IsConstruct, IsStatic, IsThrows, comment)
	//右括号
	this.eatToken(TOKEN_RPAREN)

	if this.current.Kind == TOKEN_THROWS {
		this.eatToken(TOKEN_THROWS)
		this.eatToken(TOKEN_ID)
		for this.current.Kind == TOKEN_COMMER {
			this.eatToken(TOKEN_COMMER)
			this.eatToken(TOKEN_ID)
		}
		IsThrows = true
	}
	//左大括号
	this.eatToken(TOKEN_LBRACE)
	var stms []ast.Stm

	//解析本地变量和表达式
	stms = this.parseStatements()

	this.eatToken(TOKEN_RBRACE)

	return ast.NewMethodSingle(dec.Tp, dec.Name, formals, stms, IsConstruct, IsStatic, IsThrows, comment)
}
