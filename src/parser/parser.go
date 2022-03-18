package parser

import (
	"dog/ast"
	"dog/cfg"
	"dog/control"
	"dog/storage"
	"dog/util"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	"strconv"
	"strings"
)

type Parser struct {
	*Stack
	lexer             *Lexer
	current           *Token
	pending           []*Token
	currentType       ast.Exp
	assignType        ast.Exp
	isSpecial         bool
	isField           bool
	isAnnotationClass bool
	GetField          FieldFunc
	currentFile       ast.File //当前解析的File
	classStack        *Stack
	methodStack       *Stack
	stmStack          *Stack
	currentClass      ast.Class  //当前解析的class TODO 类嵌套
	currentMethod     ast.Method //当前解析的Method	TODO 函数嵌套
	currentStm        ast.Stm    //当前解析的Stm
	Linenum           int
	ProjectPath       string //项目路径
}

func NewParse(fname string, buf []byte) *Parser {
	lexer := NewLexer(fname, buf)
	p := new(Parser)
	p.Stack = InitStack()
	p.classStack = InitStack()
	p.methodStack = InitStack()
	p.stmStack = InitStack()
	p.lexer = lexer
	p.current = p.lexer.NextToken()
	return p
}

func (this *Parser) advance() {
	if this.isAnnotationClass {
		log.Infof("注解类，不处理")
		return
	}
	if control.Lexer_dump == true {
		util.Debug(this.current.ToString())
	}
	this.Linenum = this.current.LineNum
	this.current = this.lexer.NextToken()

	//处理所有注解
	for this.current.Kind == TOKEN_AT {
		this.parseAnnotation()
	}

}
func (this *Parser) advanceOnly() {
	if control.Lexer_dump == true {
		util.Debug(this.current.ToString())
	}
	this.Linenum = this.current.LineNum
	this.current = this.lexer.NextToken()
}

func (this *Parser) eatToken(kind int) {
	if this.isAnnotationClass {
		log.Infof("注解类，不处理")
		return
	}
	if kind == this.current.Kind {
		this.advance()
	} else if TOKEN_COMMENT == this.current.Kind {
		this.advance()
		this.eatToken(kind)
	} else {
		util.ParserError(tMap[kind], tMap[this.current.Kind], this.current.LineNum, this.lexer.fname)
	}
}

//解析泛型实例化参数列表
func (this *Parser) parseTypeList() (types []ast.Exp) {
	log.Infof("解析泛型参数列表")
	types = []ast.Exp{}
	if this.current.Kind == TOKEN_GT {
		return types
	}

	if this.current.Kind == TOKEN_QUESTION {
		this.eatToken(TOKEN_QUESTION)
		types = append(types, ast.NewIdent("?", this.Linenum))
	} else {
		tp := this.parseType()
		types = append(types, tp)
	}
	//FIXME 忽略
	if this.current.Kind == TOKEN_EXTENDS {
		this.eatToken(TOKEN_EXTENDS)
		this.parseType()
	}

	for this.current.Kind == TOKEN_COMMER {
		this.advance()
		tp := this.parseType()
		types = append(types, tp)
		//FIXME 忽略
		if this.current.Kind == TOKEN_EXTENDS {
			this.eatToken(TOKEN_EXTENDS)
			this.parseType()
		}
	}
	return types
}

//
//
// return:
func (this *Parser) parseFormalList(isSingle bool) (flist []ast.Field) {
	//空参数
	if this.current.Kind == TOKEN_RPAREN {
		return
	}
	log.Debugf("解析函数参数")
	flist = []ast.Field{}
	var tp ast.Exp
	var id string
	var access int

	tp = this.parseType()
	id = this.current.Lexeme
	//可变参数
	if this.current.Kind == TOKEN_DOT {
		this.eatToken(TOKEN_DOT)
		this.eatToken(TOKEN_DOT)
		this.eatToken(TOKEN_DOT)
	}
	this.eatToken(TOKEN_ID)
	flist = append(flist, ast.NewFieldSingle(access, tp, ast.NewIdent(id, this.Linenum), nil, false, false))
	//处理注释
	if this.current.Kind == TOKEN_COMMENT {
		this.advance()
	}

	for this.current.Kind == TOKEN_COMMER && !isSingle {
		log.Debugf("解析函数 --> 需要类型推断")
		this.eatToken(TOKEN_COMMER)
		if this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}
		tp = this.parseType()
		//可变参数
		if this.current.Kind == TOKEN_DOT {
			this.eatToken(TOKEN_DOT)
			this.eatToken(TOKEN_DOT)
			this.eatToken(TOKEN_DOT)
		}
		id = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		flist = append(flist, ast.NewFieldSingle(access, tp, ast.NewIdent(id, this.Linenum), nil, false, false))

	}
	for _, vv := range flist {
		if this.Peek() != nil {
			this.Peek().AddField(vv)
		}

	}
	return flist
}

//强制类型转换
func (this *Parser) parseCastExp() ast.Exp {
	switch this.current.Kind {
	case TOKEN_LPAREN:
		this.advance()
		log.Debugf("解析 parseCastExp:%v", this.current.Lexeme)
		tp := this.parseType()
		this.eatToken(TOKEN_RPAREN)

		exp := this.parseExp()

		log.Debugf("解析 parseCastExp:%v", exp)

		return ast.Cast_new(tp, exp, this.Linenum)
	}
	return nil
}

// 成员(函数/变量)访问语句 "." 或 "(" 开头作为判定条件
//
// param: x
// return:
func (this *Parser) parseCallExp(x ast.Exp) (ret ast.Exp) {
	var builder = false
	if this.current.Kind == TOKEN_LPAREN {
		this.eatToken(TOKEN_LPAREN)
		args := this.parseExpList() //1
		this.eatToken(TOKEN_RPAREN)
		x = ast.CallExpr_new(x, args, this.Linenum)
	}
	var streamProbability int

	for this.current.Kind == TOKEN_DOT {
		var old = x
		var isListOrMapGetSet bool
		var eleType ast.Exp
		var isListAdd bool
		var isListRemove bool
		var isListOrMapClear bool

		this.eatToken(TOKEN_DOT)
		if this.current.Lexeme == "Length" || this.current.Lexeme == "length" {
			this.advance()
			if this.current.Kind == TOKEN_LPAREN {
				this.eatToken(TOKEN_LPAREN)
				this.eatToken(TOKEN_RPAREN)
			}

			return ast.Length_new(x, this.Linenum)
		}
		//else if this.current.Lexeme == "Size" || this.current.Lexeme == "size" {
		//	this.advance()
		//	this.eatToken(TOKEN_LPAREN)
		//	this.eatToken(TOKEN_RPAREN)
		//	return ast.Length_new(x, this.Linenum)
		//}
		if this.current.Kind == TOKEN_CLASS {
			this.advance()
			return ast.ClassExp_new(x, this.Linenum)
		}
		if this.current.Lexeme == "stream" ||
			this.current.Lexeme == "filter" ||
			this.current.Lexeme == "map" ||
			this.current.Lexeme == "mapToInt" ||
			this.current.Lexeme == "mapToLong" ||
			this.current.Lexeme == "mapToDouble" ||
			this.current.Lexeme == "flatMap" ||
			this.current.Lexeme == "sorted" ||
			this.current.Lexeme == "peek" ||
			this.current.Lexeme == "limit" ||
			this.current.Lexeme == "forEachOrdered" ||
			this.current.Lexeme == "toList" ||
			this.current.Lexeme == "toArray" ||
			this.current.Lexeme == "min" ||
			this.current.Lexeme == "max" ||
			this.current.Lexeme == "collect" ||
			this.current.Lexeme == "forEach" {
			streamProbability++
		}
		//处理builder注解函数
		if this.current.Lexeme == "builder" {
			builder = true
			this.eatToken(TOKEN_ID)
			this.eatToken(TOKEN_LPAREN)
			this.eatToken(TOKEN_RPAREN)
			x = ast.BuilderExpr_new(x, this.Linenum)
			//处理最后的build
		} else if this.current.Lexeme == "build" {
			this.eatToken(TOKEN_ID)
			this.eatToken(TOKEN_LPAREN)
			this.eatToken(TOKEN_RPAREN)
		} else if builder {
			x = ast.SelectorExpr_new(x, "Set"+util.Capitalize(this.current.Lexeme), this.Linenum)
			this.eatToken(TOKEN_ID)
		} else {

			//处理Map,List元素访问
			if cfg.MapListIdxAccess && (this.current.Lexeme == "get" || this.current.Lexeme == "set") {

				if id, ok := x.(*ast.Ident); ok {
					if this.CheckField(id.Name) != nil {
						el1, ok1 := this.CheckField(id.Name).GetDecType().(*ast.MapType)
						el2, ok2 := this.CheckField(id.Name).GetDecType().(*ast.ListType)

						if ok1 {
							eleType = el1.Value
						} else if ok2 {
							eleType = el2.Ele
						}
						isListOrMapGetSet = ok1 || ok2
					}
				} else if idx, ok := x.(*ast.IndexExpr); ok {
					if idx.EleType != nil {
						el1, ok1 := idx.EleType.(*ast.MapType)
						el2, ok2 := idx.EleType.(*ast.ListType)

						if ok1 {
							eleType = el1.Value
						} else if ok2 {
							eleType = el2.Ele
						}
						isListOrMapGetSet = ok1 || ok2
					}

				}

			} else if cfg.MapListIdxAccess && this.current.Lexeme == "add" {
				if id, ok := x.(*ast.Ident); ok {
					if this.CheckField(id.Name) != nil {
						_, isListAdd = this.CheckField(id.Name).GetDecType().(*ast.ListType)

					}
				}
			} else if cfg.MapListIdxAccess && this.current.Lexeme == "remove" {
				if id, ok := x.(*ast.Ident); ok {
					if this.CheckField(id.Name) != nil {
						_, isListRemove = this.CheckField(id.Name).GetDecType().(*ast.MapType)

					}
				}
			} else if cfg.MapListIdxAccess && (this.current.Lexeme == "clear") {
				if id, ok := x.(*ast.Ident); ok {
					if this.CheckField(id.Name) != nil {

						el1, ok1 := this.CheckField(id.Name).GetDecType().(*ast.MapType)
						el2, ok2 := this.CheckField(id.Name).GetDecType().(*ast.ListType)

						if ok1 {
							eleType = ast.NewHash_new(el1.Key, el1.Value, this.Linenum)
						} else if ok2 {
							eleType = ast.NewList_new(el2.Ele, nil, this.Linenum)
						}
						isListOrMapClear = ok1 || ok2
					}
				}
			}
			x = ast.SelectorExpr_new(x, this.current.Lexeme, this.Linenum)

			this.eatToken(TOKEN_ID)
		}
		if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList() //2
			for idx, v := range args {
				var call string
				var mp = new(ast.StreamStm)
				if this.CheckStreamExprs(v, &call, mp) {
					arg := ast.NewIdent(fmt.Sprintf("arg%d", idx), this.Linenum)
					mp.Left = arg
					mp.LineNum = this.Linenum
					v = arg
				}
			}
			this.eatToken(TOKEN_RPAREN)
			if isListOrMapGetSet && len(args) == 1 {
				x = ast.IndexExpr_newEx(old, args[0], eleType, this.Linenum)
				//处理List Map Get&&Set操作
			} else if isListOrMapGetSet && len(args) == 2 {
				x = ast.FakeExpr_new(ast.Assign_new(ast.IndexExpr_new(old, args[0], this.Linenum), args[1], false, this.Linenum), this.Linenum)
				//处理List Add操作
			} else if isListAdd && len(args) == 1 {
				var args1 []ast.Exp
				args1 = append(args1, old)
				args1 = append(args1, args[0])
				x = ast.FakeExpr_new(ast.Assign_new(old, ast.CallExpr_new(ast.NewIdent("append", this.Linenum), args1, this.Linenum), false, this.Linenum), this.Linenum)

			} else if isListRemove && len(args) == 1 {
				var args1 []ast.Exp
				args1 = append(args1, old)
				args1 = append(args1, args[0])
				x = ast.CallExpr_new(ast.NewIdent("delete", this.Linenum), args1, this.Linenum)

			} else if isListOrMapClear && len(args) == 0 {
				x = ast.FakeExpr_new(ast.Assign_new(old, eleType, false, this.Linenum), this.Linenum)

			} else {
				x = ast.CallExpr_new(x, args, this.Linenum)
			}

		}
		if this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}
	}
	return x
}

func (this *Parser) CovertExp(exp ast.Exp) (dst ast.Exp) {
	switch exp.(type) {
	case *ast.SelectorExpr:

	case *ast.Ident:

	}
	return
}

//AtomExp   -> (exp)
//          -> INTEGER_LITERAL
//          -> true
//          -> false
//          -> this
//          -> id
//          -> new int[exp]
//          -> new id()
func (this *Parser) parseAtomExp() ast.Exp {
	log.Debugf("解析 parseAtomExp")
	for this.current.Kind == TOKEN_COMMENT {
		log.Debugf("--------->去除注释:%v", this.current.Lexeme)
		this.advance()
	}
	switch this.current.Kind {
	case TOKEN_SUB:
		this.advance()
		if this.current.Kind == TOKEN_NUM {
			num := this.current.Lexeme
			this.advance()
			s, _ := strconv.Atoi(num)
			s = -s
			n := new(ast.Num)
			n.Value = s
			n.LineNum = this.Linenum
			return n
		} else if this.current.Kind == TOKEN_ID {
			id := this.current.Lexeme
			this.advance()
			return ast.NewIdent("-"+id, this.Linenum)

		} else {
			this.ParseBug("加法解析bug")
		}

	case TOKEN_LPAREN:
		this.eatToken(TOKEN_LPAREN)
		//可能强制类型转换
		var isCast bool
		if this.IsTypeToken() || this.currentFile.GetImport(this.current.Lexeme) != nil || this.currentClass.GetGeneric(this.current.Lexeme) != nil || this.current.Lexeme == this.currentClass.GetName() {
			log.Debugf("强制类型转换")
			isCast = true
		}
		exps := this.parseExpList()
		this.eatToken(TOKEN_RPAREN)

		//	1 Lambda表达式参数
		if this.current.Kind == TOKEN_LAMBDA {
			log.Debugf("发现 TOKEN_LAMBDA")
			this.advance()
			return this.parseLambdaExp(exps)

		} else if this.current.Kind == TOKEN_QUESTION {
			return exps[0]
		}

		//1 强制类型转换
		if len(exps) == 1 && isCast && this.current.Kind != TOKEN_SEMI && this.current.Kind != TOKEN_RPAREN && this.current.Kind != TOKEN_LBRACK {
			//是类型表达式
			var tp ast.Exp
			switch v := exps[0].(type) {
			case ast.Type:
				tp = exps[0]
			case *ast.Ident:
				tp = &ast.ClassType{
					Name:     v.Name,
					TypeKind: ast.TYPE_CLASS,
				}
			case *ast.SelectorExpr:
				tp = &ast.ClassType{
					Name:     v.Sel,
					TypeKind: ast.TYPE_CLASS,
				}
			default:
				log.Debugf("________--> %v", v)
				this.ParseBug("强制类型转换bug")
			}
			exp := this.parseExp()
			return ast.Cast_new(tp, exp, this.Linenum)

		}
		//3 小括号优先级
		return exps[0]
	case TOKEN_NUM:
		value, _ := strconv.Atoi(this.current.Lexeme)
		this.advance()
		return ast.Num_new(value, this.Linenum)
	case TOKEN_TRUE:
		this.advance()
		return &ast.True{}
	case TOKEN_FALSE:
		this.advance()
		return &ast.False{}
	case TOKEN_NULL:
		this.advance()
		return &ast.Null{}
	case TOKEN_THIS:
		this.advance()
		//调用构造函数
		if this.current.Kind == TOKEN_LPAREN {
			return this.parseCallExp(&ast.This{})
		}
		return &ast.This{}
	case TOKEN_SYSTEM:
		var x ast.Exp
		x = ast.NewIdent(this.current.Lexeme, this.Linenum)
		this.eatToken(TOKEN_SYSTEM)
		if this.current.Lexeme == "CurrentTimeMillis" {
			this.advance()
			this.advance()
			this.advance()
			return ast.NewIdent("time.Now().UnixMilli()", this.Linenum)
		}
		x = this.parseCallExp(x)
		return x
	case TOKEN_ID:
		id := this.current.Lexeme
		this.advance()
		//可能是new 泛型表达式
		if this.current.Kind == TOKEN_LT {
			if this.currentFile.GetImport(id) != nil {
				//FIXME 这里没有处理
				this.eatToken(TOKEN_LT)
				if this.current.Kind != TOKEN_GT {
					this.parseNotExp()
				}
				this.eatToken(TOKEN_GT)
			}

		}

		if this.current.Kind == TOKEN_LPAREN {
			m := ast.NewIdent(id, this.Linenum)
			return this.parseCallExp(m)
		} else if this.current.Kind == TOKEN_LAMBDA {
			log.Debugf("发现 单参数 TOKEN_LAMBDA")
			this.advance()
			return this.parseLambdaExp([]ast.Exp{ast.NewIdent(id, this.Linenum)})
		} else if this.current.Kind == TOKEN_ID {
			//tp := &ast.ClassType{this.current.Lexeme, ast.TYPE_CLASS}
			//this.eatToken(TOKEN_ID)
			//return ast.DefExpr_new(id, tp, this.Linenum)
		}
		log.Debugf("适配变量ID->%s", id)
		return ast.NewIdent(id, this.Linenum)
	case TOKEN_NEW:
		return this.parseNewExp()
	default:
		if this.IsTypeToken() {
			return this.parseType()
		}
		this.ParseBug("parseAtomExp未处理")
	}
	return nil
}

// new id[]			   -->
// new id(xx,xx,xx...) -->构造函数
// new id(xx,xx,xx...) -->泛型构造函数
// new id<xx,xx,xx>()
//
// return:
func (this *Parser) parseNewExp() ast.Exp {
	log.Debugf("----------->解析 New语句")
	this.advance()
	switch this.current.Kind {
	case TOKEN_BYTE:
		this.advance()
		this.eatToken(TOKEN_LBRACK)
		exp := this.parseExp()
		this.eatToken(TOKEN_RBRACK)
		return ast.NewArray_new(ast.NewIdent("byte", this.Linenum), exp, this.Linenum)
	case TOKEN_INT:
		this.advance()
		this.eatToken(TOKEN_LBRACK)
		exp := this.parseExp()
		this.eatToken(TOKEN_RBRACK)
		return ast.NewIntArray_new(exp, this.Linenum)
	case TOKEN_INTEGER:
		this.advance()
		this.eatToken(TOKEN_LPAREN)
		exp := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		return ast.NewObject_new(exp, this.Linenum)
	case TOKEN_STRING:
		this.advance()
		if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			return ast.CallExpr_new(ast.NewIdent("strconv.Itoa", this.Linenum), args, this.Linenum)
		}
		this.eatToken(TOKEN_LBRACK)
		//new String[]{xxx, xxxx};
		if this.current.Kind == TOKEN_RBRACK {
			this.eatToken(TOKEN_RBRACK)
			if this.current.Kind == TOKEN_LBRACE {
				this.eatToken(TOKEN_LBRACE)
				args := this.parseExpList()
				this.eatToken(TOKEN_RBRACE)
				return ast.NewStringArray_new(nil, args, this.Linenum)
			}
			return ast.NewStringArray_new(nil, nil, this.Linenum)
		}
		exp := this.parseExp()
		this.eatToken(TOKEN_RBRACK)
		return ast.NewStringArray_new(exp, nil, this.Linenum)
	case TOKEN_HASHMAP:
		this.eatToken(TOKEN_HASHMAP)
		//非泛型
		if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			this.eatToken(TOKEN_RPAREN)
			return ast.NewHash_new(&ast.String{ast.TYPE_STRING}, &ast.ObjectType{ast.TYPE_OBJECT}, this.Linenum)
		}
		this.eatToken(TOKEN_LT)
		var key ast.Exp
		var ele ast.Exp
		if this.current.Kind != TOKEN_GT {
			key = this.parseType()

			this.eatToken(TOKEN_COMMER)
			ele = this.parseType()
			//类型推到
		} else {
			t, ok := this.assignType.(*ast.MapType)
			if ok {
				key = t.Key
				ele = t.Value
			} else {
				//this.ParseBug("Hash类型存在空")
				key = &ast.String{}
				ele = &ast.ObjectType{}
			}
		}
		this.eatToken(TOKEN_GT)
		this.eatToken(TOKEN_LPAREN)
		if this.current.Kind == TOKEN_NUM {
			//FIXME 处理map容量
			this.eatToken(TOKEN_NUM)
		} else if this.current.Kind != TOKEN_RPAREN {
			this.parseNotExp()
		}
		this.eatToken(TOKEN_RPAREN)
		return ast.NewHash_new(key, ele, this.Linenum)
	case TOKEN_ARRAYLIST:
		this.eatToken(TOKEN_ARRAYLIST)
		var ele ast.Exp
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			if this.current.Kind != TOKEN_GT {
				ele = this.parseType()
			} else {
				_, ok := this.assignType.(*ast.ListType)
				if !ok {
					ele = &ast.ObjectType{ast.TYPE_OBJECT}
				} else {
					ele = this.assignType.(*ast.ListType).Ele
					this.assignType = nil
				}

			}
			this.eatToken(TOKEN_GT)

		} else {
			ele = &ast.ObjectType{ast.TYPE_OBJECT}

		}

		this.eatToken(TOKEN_LPAREN)
		args := this.parseExpList()
		this.eatToken(TOKEN_RPAREN)
		return ast.NewList_new(ele, args, this.Linenum)
	case TOKEN_DATE:
		this.eatToken(TOKEN_DATE)
		this.eatToken(TOKEN_LPAREN)
		if this.current.Kind != TOKEN_RPAREN {
			exps := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			return ast.NewDateParam_new(this.Linenum, exps)
		}
		this.eatToken(TOKEN_RPAREN)
		return ast.NewDate_new(this.Linenum)
	case TOKEN_HASHSET:
		this.eatToken(TOKEN_HASHSET)
		var ele ast.Exp
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			if this.current.Kind != TOKEN_GT {
				ele = this.parseType()
			} else {
				ele = this.assignType.(*ast.SetType).Ele
			}

			this.eatToken(TOKEN_GT)
		} else {
			ele = &ast.ObjectType{ast.TYPE_OBJECT}
		}

		this.eatToken(TOKEN_LPAREN)
		args := this.parseExpList()
		this.eatToken(TOKEN_RPAREN)
		return ast.NewSet_new(ele, args, this.Linenum)
	//带参数对象初始化
	case TOKEN_ID:
		var typeName = this.current.Lexeme
		var args []ast.Exp
		log.Debugf("-------------> %v", this.current.Lexeme)
		id := this.current.Lexeme

		defer func() {
			//直接实现接口
			if this.current.Kind == TOKEN_LBRACE {
				classSingle := ast.NewClassSingle(this.currentFile, 0, id, "", ast.CLASS_TYPE)
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
				this.eatToken(TOKEN_LBRACE)
				this.parseClassContext(classSingle)
				this.eatToken(TOKEN_RBRACE)
			}

		}()

		this.eatToken(TOKEN_ID)
		exp := this.parseCallExp(ast.NewIdent(id, this.Linenum))

		//模板
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			if this.current.Kind != TOKEN_GT {
				this.parseType()
				//泛型多个参数
				for this.current.Kind == TOKEN_COMMER {
					this.eatToken(TOKEN_COMMER)
					this.parseType()
				}
			}
			this.eatToken(TOKEN_GT)
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			return ast.CallExpr_new(exp, args, this.Linenum)
		}
		//数组
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			// new xxx[]{xxx, xxx,xxx....};
			if this.current.Kind == TOKEN_RBRACK {
				this.eatToken(TOKEN_RBRACK)
				if this.current.Kind == TOKEN_LBRACE {
					this.eatToken(TOKEN_LBRACE)
					args = this.parseExpList()
					this.eatToken(TOKEN_RBRACE)
					return ast.NewObjectArray_new(&ast.ClassType{typeName, ast.TYPE_CLASS}, args, nil, this.Linenum)
				}

			} else {
				size := this.parseExp()
				this.eatToken(TOKEN_RBRACK)
				if this.current.Kind == TOKEN_LBRACE {
					this.eatToken(TOKEN_LBRACE)
					eles := this.parseExpList()
					this.eatToken(TOKEN_RBRACE)
					return ast.NewObjectArray_new(&ast.ClassType{typeName, ast.TYPE_CLASS}, eles, size, this.Linenum)
				}
				return ast.NewObjectArray_new(&ast.ClassType{typeName, ast.TYPE_CLASS}, nil, size, this.Linenum)
			}
		}

		return exp
	default:
		this.ParseBug("未处理New类型")
		return nil
	}

}

// 解析函数调用参数列表
//
// return:
func (this *Parser) parseExpList() (args []ast.Exp) {
	args = []ast.Exp{}
	if this.current.Kind == TOKEN_RPAREN {
		return args
	}
	//判断是不是lambda是不是lambda表达式
	//（exp）-> exp
	// (exp) -> {exp}
	//可能是lambda表达式
	exp := this.parseExp()
	//带类型的变量声明
	if this.current.Kind == TOKEN_ID {
		id := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		args = append(args, ast.DefExpr_new(ast.NewIdent(id, this.Linenum), exp, this.Linenum))

		for this.current.Kind == TOKEN_COMMER {
			this.advance()
			exp = this.parseExp()
			id = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			args = append(args, ast.DefExpr_new(ast.NewIdent(id, this.Linenum), exp, this.Linenum))
		}
	} else {
		args = append(args, exp)

		for this.current.Kind == TOKEN_COMMER {
			this.advance()
			//处理注释
			for this.current.Kind == TOKEN_COMMENT {
				this.advance()
			}
			args = append(args, this.parseExp())
		}
	}

	return args
}

func (this *Parser) parseLambdaExp(args []ast.Exp) (exp ast.Exp) {
	log.Debugf("尝试解析 --> Lambda")
	//处理参数
	var fields []ast.Field
	for _, arg := range args {
		switch v := arg.(type) {
		case *ast.DefExpr:
			fields = append(fields, ast.NewFieldSingle(0, v.Tp, v.Name, nil, false, false))
		case *ast.Ident:
			fields = append(fields, ast.NewFieldSingle(0, &ast.ObjectType{ast.TYPE_OBJECT}, v, nil, false, false))
		default:
			panic("parseLambdaExp")
		}
	}
	fake := ast.FakeStm_new(this.Peek(), this.Linenum)
	for _, vv := range fields {
		fake.AddField(vv)
	}
	this.Push(fake)
	defer func() {
		this.Pop()
	}()

	if this.current.Kind == TOKEN_LBRACE {
		this.eatToken(TOKEN_LBRACE)
		stms := this.parseStatements()
		this.eatToken(TOKEN_RBRACE)
		lam := ast.Lambda_new(fields, stms, this.Linenum)
		return lam
	} else {
		stm := &ast.ExprStm{
			E: this.parseExp(),
		}
		return ast.Lambda_new(fields, []ast.Stm{stm}, this.Linenum)
	}
	return
}

//NotExp    -> AtomExp
//          -> AtomExp.id(explist)
//          -> AtomExp[exp]
//          -> AtomExp.length
func (this *Parser) parseNotExp() ast.Exp {
	log.Debugf("解析 parseNotExp")
	exp := this.parseAtomExp()
	for this.current.Kind == TOKEN_DOT ||
		this.current.Kind == TOKEN_DOUBLE_COLON ||
		this.current.Kind == TOKEN_COMMENT ||
		//FIXME 自增,自减作为语句处理
		this.current.Kind == TOKEN_DECREMENT || //后缀加
		this.current.Kind == TOKEN_INCREMENT || //后缀减
		this.current.Kind == TOKEN_LBRACK {
		for this.current.Kind == TOKEN_COMMENT {
			log.Debugf("--------->去除注释:%v", this.current.Lexeme)
			this.advance()
		}

		switch this.current.Kind {
		case TOKEN_INCREMENT:
			this.eatToken(TOKEN_INCREMENT)
			return ast.Increment_new(exp, false, this.Linenum)
		case TOKEN_DECREMENT:
			this.eatToken(TOKEN_DECREMENT)
			return ast.Decrement_new(exp, false, this.Linenum)
		//可以不断循环下去
		case TOKEN_DOT:
			log.Debugf("解析函数调用,或成员变量")
			exp = this.parseCallExp(exp)
		case TOKEN_DOUBLE_COLON:
			log.Debugf("方法引用")
			this.eatToken(TOKEN_DOUBLE_COLON)
			if this.current.Kind == TOKEN_NEW {
				this.eatToken(TOKEN_NEW)
				exp = ast.MethodReference_new(exp, ast.NewIdent("New", this.Linenum), this.Linenum)
			} else {
				m := ast.NewIdent(util.Capitalize(this.current.Lexeme), this.Linenum)
				this.eatToken(TOKEN_ID)
				exp = ast.MethodReference_new(exp, m, this.Linenum)
			}

			//数组索引操作
		case TOKEN_LBRACK: //[exp]
			this.advance()
			index := this.parseExp()
			if index != nil {
				log.Debugf("数组索引表达式 --> %v", exp)
				this.eatToken(TOKEN_RBRACK)
				if this.current.Kind == TOKEN_DOT {
					exp = ast.IndexExpr_new(exp, index, this.Linenum)
				} else {
					return ast.IndexExpr_new(exp, index, this.Linenum)
				}
			} else {
				log.Debugf("数组索引用")
				this.eatToken(TOKEN_RBRACK)
				panic("数组索引用")
			}
		}
	}
	return exp
}
func Add(aa int) {

}

//TimesExp  -> !TimesExp
//          -> NotExp
//			-> ++NotExp
//			-> --NotExp
func (this *Parser) parseTimeExp() ast.Exp {
	log.Debugf("解析 parseTimeExp")
	var exp2 ast.Exp
	var opt = this.current.Kind
	if this.current.Kind == TOKEN_INCREMENT ||
		this.current.Kind == TOKEN_DECREMENT ||
		this.current.Kind == TOKEN_NOT {
		this.advance()
		exp2 = this.parseTimeExp()
	}
	if exp2 != nil {
		switch opt {
		case TOKEN_NOT:
			return ast.Not_new(exp2, this.Linenum)
		case TOKEN_INCREMENT:
			return ast.Increment_new(exp2, true, this.Linenum)
		case TOKEN_DECREMENT:
			return ast.Decrement_new(exp2, true, this.Linenum)
		default:
			panic("不支持")
		}

	} else {
		return this.parseNotExp()
	}
}

//AddSubExp -> TimesExp * TimesExp
//          -> TimesExp / TimesExp
//          -> TimesExp % TimesExp
//          -> TimesExp
func (this *Parser) parseAddSubExp() ast.Exp {
	log.Debugf("解析 parseAddSubExp")
	left := this.parseTimeExp()
	//去除注释
	for this.current.Kind == TOKEN_COMMENT {
		log.Debugf("--------->去除注释:%v", this.current.Lexeme)
		this.advance()
	}
	for this.current.Kind == TOKEN_REM ||
		this.current.Kind == TOKEN_STAR ||
		this.current.Kind == TOKEN_QUO {
		switch this.current.Kind {
		case TOKEN_STAR:
			this.advance()
			//去除注释
			for this.current.Kind == TOKEN_COMMENT {
				log.Debugf("--------->去除注释:%v", this.current.Lexeme)
				this.advance()
			}
			right := this.parseTimeExp()
			left = ast.Times_new(left, right, this.Linenum)
		case TOKEN_QUO:
			this.advance()
			//去除注释
			for this.current.Kind == TOKEN_COMMENT {
				this.advance()
			}
			right := this.parseTimeExp()
			left = ast.Division_new(left, right, this.Linenum)
		case TOKEN_REM:
			this.advance()
			//去除注释
			for this.current.Kind == TOKEN_COMMENT {
				this.advance()
			}
			right := this.parseTimeExp()
			left = ast.Remainder_new(left, right, this.Linenum)

		}

	}
	return left
}

//LtExp -> AddSubExp + AddSubExp
//      -> AddSubExp - AddSubExp
//      -> AddSubExp
func (this *Parser) parseLtExp() ast.Exp {
	log.Debugf("解析parseLtExp")
	left := this.parseAddSubExp()
	//去除注释
	for this.current.Kind == TOKEN_COMMENT {
		log.Debugf("--------->去除注释:%v", this.current.Lexeme)
		this.advance()
	}
	for this.current.Kind == TOKEN_ADD ||
		this.current.Kind == TOKEN_SUB {
		switch this.current.Kind {
		case TOKEN_ADD:
			this.advance()
			//去除注释
			for this.current.Kind == TOKEN_COMMENT {
				log.Debugf("--------->去除注释:%v", this.current.Lexeme)
				this.advance()
			}
			right := this.parseAddSubExp()
			left = ast.Add_new(left, right, this.Linenum)
			for this.current.Kind == TOKEN_COMMENT {
				log.Debugf("--------->去除注释:%v", this.current.Lexeme)
				this.advance()
			}
		case TOKEN_SUB:
			this.advance()
			//去除注释
			for this.current.Kind == TOKEN_COMMENT {
				log.Debugf("--------->去除注释:%v", this.current.Lexeme)
				this.advance()
			}
			right := this.parseAddSubExp()
			left = ast.Sub_new(left, right, this.Linenum)
			for this.current.Kind == TOKEN_COMMENT {
				log.Debugf("--------->去除注释:%v", this.current.Lexeme)
				this.advance()
			}
		}
	}
	return left
}

//EqExp    -> EqExp == EqExp || EqExp != EqExp
//          -> EqExp
func (this *Parser) parseEqExp() ast.Exp {
	log.Debugf("解析parseEqExp")
	left := this.parseLtExp()
	for this.current.Kind == TOKEN_LT || this.current.Kind == TOKEN_LE || this.current.Kind == TOKEN_GT || this.current.Kind == TOKEN_GE {
		opt := this.current.Kind
		this.advance()
		//泛型应该在前面被拦截处理 -- > 原理是类型不能被用来比较大小
		right := this.parseLtExp()
		switch opt {
		case TOKEN_LT:
			return ast.Lt_new(left, right, this.Linenum)
		case TOKEN_LE:
			return ast.Le_new(left, right, this.Linenum)
		case TOKEN_GT:
			return ast.Gt_new(left, right, this.Linenum)
		case TOKEN_GE:
			return ast.Ge_new(left, right, this.Linenum)

		}
	}
	return left
}

//AndExp    -> EqExp == EqExp  EqExp != EqExp
//          -> EqExp
func (this *Parser) parseAndExp() ast.Exp {
	log.Debugf("解析 parseAndExp")
	left := this.parseEqExp()

	for this.current.Kind == TOKEN_EQ || this.current.Kind == TOKEN_NE {
		opt := this.current.Kind
		this.advance()
		right := this.parseEqExp()
		switch opt {
		case TOKEN_EQ:
			return ast.Eq_new(left, right, this.Linenum)
		case TOKEN_NE:
			return ast.Neq_new(left, right, this.Linenum)
		}
	}
	return left
}

//Exp -> AndExp & AndExp
//    -> AndExp
func (this *Parser) parseOrExp() ast.Exp {
	log.Debugf("解析 parseOrExp")
	left := this.parseAndExp()
	for this.current.Kind == TOKEN_AND {
		this.advance()
		right := this.parseAndExp()
		left = ast.And_new(left, right, this.Linenum)
	}
	return left
}

//AndExp    -> EqExp == EqExp  EqExp != EqExp
//          -> EqExp
func (this *Parser) parseLAndExp() ast.Exp {
	log.Debugf("解析 parseLAndExp")
	left := this.parseOrExp()
	for this.current.Kind == TOKEN_OR {
		this.advance()
		right := this.parseOrExp()
		left = ast.Or_new(left, right, this.Linenum)
	}
	return left
}

//Exp -> AndExp && AndExp
//    -> AndExp
func (this *Parser) parseLOrExp() ast.Exp {
	log.Debugf("解析 parseLOrExp")
	left := this.parseLAndExp()
	for this.current.Kind == TOKEN_LAND {
		this.advance()
		right := this.parseLAndExp()
		left = ast.LAnd_new(left, right, this.Linenum)
	}
	return left
}

//OrExp    -> OrExp || OrExp
//          -> OrExp
func (this *Parser) parseQuestionExp() ast.Exp {
	log.Debugf("解析 parseQuestionExp")
	left := this.parseLOrExp()
	for this.current.Kind == TOKEN_LOR {
		log.Debugf("TOKEN_LOR")
		this.advance()
		right := this.parseLOrExp()
		left = ast.LOr_new(left, right, this.Linenum)
	}
	return left
}

//OrExp    -> OrExp || OrExp
//          -> OrExp
func (this *Parser) parseExp() ast.Exp {
	log.Debugf("解析 parseExp")

	//数组赋值语句
	if this.current.Kind == TOKEN_LBRACE {
		log.Infof("TOKEN_LBRACK --> 数组赋值语句")
		this.eatToken(TOKEN_LBRACE)
		if this.current.Kind == TOKEN_RBRACE {
			this.eatToken(TOKEN_RBRACE)
			return ast.ArrayAssign_new([]ast.Exp{ast.NewIdent("", this.Linenum)}, this.currentType, this.Linenum)
		}
		exps := this.parseExpList()
		this.eatToken(TOKEN_RBRACE)
		return ast.ArrayAssign_new(exps, this.currentType, this.Linenum)
	}
	log.Debugf("解析 parseQuestionExp --> %v", this.current.Lexeme)
	left := this.parseQuestionExp()
	//
	if this.current.Kind == TOKEN_QUESTION {
		log.Debugf("发现TOKEN_QUESTION")
		this.advance()
		log.Infof("TOKEN_QUESTION --> 解析第一个表达式")
		one := this.parseExp()
		this.eatToken(TOKEN_COLON)
		log.Infof("TOKEN_QUESTION --> 解析第二个表达式")
		two := this.parseExp()
		return ast.Question_new(left, one, two, this.Linenum)
	}
	if this.current.Kind == TOKEN_INSTANCEOF {
		this.eatToken(TOKEN_INSTANCEOF)
		right := this.parseQuestionExp()
		return ast.Instanceof_new(left, right, this.Linenum)
	}

	return left
}

func (this *Parser) parseMemberVarDecl(tmp *ast.FieldSingle, IsStatic bool) (dec ast.Field) {
	var value ast.Exp
	if this.current.Kind == TOKEN_ASSIGN {
		this.eatToken(TOKEN_ASSIGN)
		this.assignType = tmp.Tp
		value = this.parseExp()
	}
	this.eatToken(TOKEN_SEMI)
	dec = ast.NewFieldSingle(tmp.Access, tmp.Tp, tmp.Name, value, IsStatic, true)
	return dec
}

// 类静态语句
//
// param: comment
// return:
func (this *Parser) parseMemberStatic(comment string) (meth ast.Method) {

	var methodSingle = ast.NewMethodSingle(this.currentClass, &ast.Void{}, ast.NewIdent("init", this.Linenum), nil, nil, false, true, false, comment)
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

	this.eatToken(TOKEN_LBRACE)
	//解析本地变量和表达式
	methodSingle.Stms = this.parseStatements()

	this.eatToken(TOKEN_RBRACE)
	return methodSingle
}

func (this *Parser) parseAnnotation() {
	this.eatToken(TOKEN_AT)
	if this.current.Kind == TOKEN_INTERFACE {
		this.isAnnotationClass = true
		log.Infof("发现注解类，不处理")
		return
	}
	//TODO 不忽略的注解
	if this.current.Lexeme == "Query" {
		this.current.Kind = TOKEN_QUERY
		return
	}
	this.eatToken(TOKEN_ID)
	//带参数的注解
	if this.current.Kind == TOKEN_LPAREN {
		this.eatToken(TOKEN_LPAREN)

		if this.current.Kind == TOKEN_RPAREN {
			this.eatToken(TOKEN_RPAREN)
			return
		}
		for {
			this.parseExp() //id
			if this.current.Kind == TOKEN_ASSIGN {
				this.advance() // =
				if this.current.Kind == TOKEN_LBRACE {
					this.advanceOnly()
					for {
						if this.current.Kind == TOKEN_AT {
							this.parseAnnotation()
						} else {
							this.parseExp() //id
						}
						if this.current.Kind == TOKEN_COMMER {
							this.eatToken(TOKEN_COMMER)
						} else {
							break
						}
					}

					this.eatToken(TOKEN_RBRACE)
				} else {
					this.parseExp() //id
				}

			}
			if this.current.Kind == TOKEN_COMMER {
				this.eatToken(TOKEN_COMMER)
			} else {
				break
			}

		}
		this.eatToken(TOKEN_RPAREN)
		if this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}
	}

}

func (this *Parser) parseProgram() ast.File {
	var name string
	//处理package
	if this.current.Kind == TOKEN_PACKAGE {
		this.advance()
		for this.current.Kind != TOKEN_SEMI {
			name = this.current.Lexeme
			this.advance()
		}
		this.advance()
	}
	this.currentFile = ast.NewFileSingle(name, nil)
	this.Push(this.currentFile)
	defer func() {
		this.currentFile = nil
		this.Pop()
	}()
	//处理import
	for this.current.Kind == TOKEN_IMPORT {
		this.advance()
		var id string
		var path string
		var pack string
		for this.current.Kind != TOKEN_SEMI {
			var dot string
			var name string
			if this.current.Kind == TOKEN_ID {
				name = this.current.Lexeme
				id = this.current.Lexeme
				this.eatToken(TOKEN_ID)
			} else if this.current.Kind == TOKEN_STAR {
				this.eatToken(TOKEN_STAR)
				path = strings.TrimSuffix(path, ".")
				cls, err := storage.FindByPack(path)
				if err != nil {
					this.ParseBug(path)
				}
				for _, vv := range cls {
					im := &ast.ImportSingle{
						Pack: pack,
						Name: vv.Name,
						Path: path,
					}
					log.Debugf("导入:%v --> %v", path, vv.Name)
					this.currentFile.AddImport(im)
				}
				continue
			} else if this.current.Kind == TOKEN_ARRAYLIST ||
				this.current.Kind == TOKEN_LIST ||
				this.current.Kind == TOKEN_DATE ||
				this.current.Kind == TOKEN_ASSERT ||
				this.current.Kind == TOKEN_MAP ||
				this.current.Kind == TOKEN_HASHMAP ||
				this.current.Kind == TOKEN_SET ||
				this.current.Kind == TOKEN_HASHSET {
				this.advance()
			} else if this.current.Kind == TOKEN_DOT {
				pack = id
				dot = "."
				this.eatToken(TOKEN_DOT)
			} else if this.current.Kind == TOKEN_STATIC {
				this.advance()
			} else {
				this.ParseBug("导入bug")
			}
			path += name + dot

		}
		im := &ast.ImportSingle{
			Pack: pack,
			Name: id,
			Path: path,
		}
		this.eatToken(TOKEN_SEMI)
		log.Debugf("导入:%v --> %v", path, id)
		this.currentFile.AddImport(im)

	}
	var comment string
	if this.current.Kind == TOKEN_COMMENT {
		comment = ""
		//处理注释
		for this.current.Kind == TOKEN_COMMENT {
			comment += "\n" + this.current.Lexeme
			this.advance()
		}
		if this.current.Kind == TOKEN_EOF || (this.current.Kind != TOKEN_CLASS && this.current.Kind != TOKEN_PRIVATE && this.current.Kind != TOKEN_PUBLIC && this.current.Kind != TOKEN_PROTECTED) {
			return this.currentFile
		}

	}
	this.parseClassDecls()
	this.eatToken(TOKEN_EOF)
	return this.currentFile
}
func (this *Parser) CheckField(name string) ast.Field {
	if this.Peek() != nil {
		return this.Peek().GetField(name)
	} else {
		log.Debugf("------------------------->检查本地变量:%v", name)
		return nil
	}
}

func (this *Parser) Parser() ast.File {
	p := this.parseProgram()
	return p
}

func (this *Parser) ParseBug(info string) {
	var msg = fmt.Sprintf("%s:%d:%s\n", this.lexer.fname, this.Linenum, info)
	log.Errorf(this.current.Lexeme)
	util.Bug(msg)
}
