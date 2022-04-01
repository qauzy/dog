package parser

import (
	"dog/ast"
	"dog/cfg"
	log "github.com/corgi-kx/logcustom"
)

type FieldFunc func(string) ast.Field

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
	//
	if this.isAnnotationClass {
		log.Infof("注解类，不处理")
		return nil
	}
	//枚举类型
	if this.current.Kind == TOKEN_ENUM {

		return this.parseEnumDecl(access)
	}

	//接口类型
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

	var generics []*ast.GenericSingle
	//FIXME 泛型忽略
	if this.current.Kind == TOKEN_LT {
		this.eatToken(TOKEN_LT)
		ge := &ast.GenericSingle{}
		ge.Name = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		if this.current.Kind == TOKEN_EXTENDS {
			this.advance()
			ge.Extends = this.current.Lexeme
			this.parseType()
		}
		generics = append(generics, ge)
		for this.current.Kind == TOKEN_COMMER {
			this.advance()

			ge = &ast.GenericSingle{}
			ge.Name = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			if this.current.Kind == TOKEN_EXTENDS {
				this.advance()
				ge.Extends = this.current.Lexeme
				this.parseType()
			}
			generics = append(generics, ge)

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

		//FIXME 忽略接口实现
		this.parseType()
		if this.current.Kind == TOKEN_COMMER {
			this.eatToken(TOKEN_COMMER)
			this.parseType()
		}

		//FIXME 泛型忽略
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			this.parseType()
			for this.current.Kind == TOKEN_COMMER {
				this.eatToken(TOKEN_COMMER)
				this.parseType()
			}
			this.eatToken(TOKEN_GT)
		}
	}

	this.eatToken(TOKEN_LBRACE)
	classSingle := ast.NewClassSingle(this.currentFile, access, id, extends, ast.CLASS_TYPE)
	for _, vv := range generics {
		log.Infof("添加泛型：%v", vv.Name)
		classSingle.AddGeneric(vv)
	}

	this.currentClass = classSingle
	this.Push(classSingle)
	this.classStack.Push(classSingle)
	defer func() {
		this.classStack.Pop()
		this.Pop()
		if this.classStack.Peek() != nil {
			this.currentClass = this.classStack.Peek().(ast.Class)
		}
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
			cl := this.parseClassDecl()
			if cl != nil {
				this.currentFile.AddClass(cl)
			}

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
	log.Infof("----------------->解析类上下文")
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
			if this.current.Kind == TOKEN_EOF || this.current.Kind != TOKEN_COMMENT {
				continue
			}
		}
		var tmp ast.FieldSingle
		var IsConstruct = false
		var IsStatic = false
		var IsAbstract = false
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
			this.current.Kind == TOKEN_ABSTRACT ||
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

			if this.current.Kind == TOKEN_ABSTRACT {
				IsAbstract = true
				this.eatToken(TOKEN_ABSTRACT)
			}

		}

		//嵌套类
		if this.current.Kind == TOKEN_CLASS {
			cl := this.parseClassDecl()
			if cl != nil {
				this.currentFile.AddClass(cl)
				continue
			}
		}

		//类静态语句
		if (this.current.Kind == TOKEN_LBRACE) && IsBlock {

			classSingle.AddMethod(this.parseMemberStatic(comment))

		} else {
			//泛型方法
			if this.current.Kind == TOKEN_LT {
				this.eatToken(TOKEN_LT)
				this.eatToken(TOKEN_ID)
				this.eatToken(TOKEN_GT)
			}

			id := this.current.Lexeme
			//处理类构造函数
			if this.currentClass.GetName() == id && prefix == false {
				this.advance()
				if this.current.Kind == TOKEN_LPAREN {
					log.Infof("处理构造函数-->%v", this.current.Lexeme)
					IsConstruct = true
					tmp.Tp = &ast.Void{ast.TYPE_VOID}
					//变量/函数名
					tmp.Name = ast.NewIdent("New"+id, this.Linenum)
				} else {
					tmp.Tp = &ast.ClassType{
						Name:     id,
						TypeKind: ast.TYPE_CLASS,
					}

					this.assignType = tmp.Tp
					//变量/函数名
					tmp.Name = ast.NewIdent(this.current.Lexeme, this.Linenum)
					this.eatToken(TOKEN_ID)
				}

			} else {
				//类型
				tmp.Tp = this.parseType()
				this.assignType = tmp.Tp
				//变量/函数名 --> 转为开头大写
				tmp.Name = ast.NewIdent(this.current.Lexeme, this.Linenum)
				if this.current.Kind == TOKEN_MAIN {
					this.eatToken(TOKEN_MAIN)
				} else {
					this.eatToken(TOKEN_ID)
				}

			}

			//成员方法
			if this.current.Kind == TOKEN_LPAREN {
				if IsStatic || cfg.AllStatic {
					classSingle.AddMethod(this.parseMemberMethod(&tmp, IsConstruct, true, IsAbstract, comment))
				} else {
					classSingle.AddMethod(this.parseMemberMethod(&tmp, IsConstruct, IsStatic, IsAbstract, comment))
				}
				//成员变量
			} else {
				if IsStatic || cfg.AllStatic {
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
func (this *Parser) parseMemberMethod(dec *ast.FieldSingle, IsConstruct bool, IsStatic bool, IsAbstract bool, comment string) (meth ast.Method) {
	log.Infof("解析函数")
	var IsThrows bool
	//左括号
	this.eatToken(TOKEN_LPAREN)
	var methodSingle = ast.NewMethodSingle(this.currentClass, dec.Tp, dec.Name, nil, nil, IsConstruct, IsStatic, IsThrows, comment)
	this.currentMethod = methodSingle
	this.methodStack.Push(methodSingle)
	this.Push(methodSingle)
	defer func() {
		this.Pop()
		this.methodStack.Pop()
		if this.methodStack.Peek() != nil {
			this.currentMethod = this.methodStack.Peek().(ast.Method)
		}
	}()
	//解析参数必须在生成currentMethod之后，因为解析需要参数作为本地变量信息 --> 函数参数会作为本地变量加入本地变量表
	methodSingle.Formals = this.parseFormalList(false)
	//右括号
	this.eatToken(TOKEN_RPAREN)

	if this.current.Kind == TOKEN_THROWS {
		this.eatToken(TOKEN_THROWS)
		this.eatToken(TOKEN_ID)
		for this.current.Kind == TOKEN_COMMER {
			this.eatToken(TOKEN_COMMER)
			this.eatToken(TOKEN_ID)
		}
		methodSingle.Throws = true
	}

	//抽象方法,结束
	if IsAbstract {
		this.eatToken(TOKEN_SEMI)
		return methodSingle
	}

	//左大括号
	this.eatToken(TOKEN_LBRACE)

	//解析本地变量和表达式
	methodSingle.Stms = this.parseStatements()
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
	this.eatToken(TOKEN_RBRACE)
	if this.current.Kind == TOKEN_SEMI {
		this.eatToken(TOKEN_SEMI)
	}

	return methodSingle
}
