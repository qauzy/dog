package parser

import (
	"dog/ast"
	"dog/control"
	"dog/util"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	"os"
	"path"
	"strconv"
)

type Parser struct {
	lexer         *Lexer
	current       *Token
	pending       []*Token
	currentType   ast.Exp
	assignType    ast.Exp
	isSpecial     bool
	isField       bool
	currentFile   ast.File   //当前解析的File
	currentClass  ast.Class  //当前解析的class TODO 类嵌套
	currentMethod ast.Method //当前解析的Method	TODO 函数嵌套
	Linenum       int
}

func NewParse(fname string, buf []byte) *Parser {
	lexer := NewLexer(fname, buf)
	p := new(Parser)
	p.lexer = lexer
	p.current = p.lexer.NextToken()
	return p
}

func (this *Parser) advance() {
	if control.Lexer_dump == true {
		util.Debug(this.current.ToString())
	}
	this.Linenum = this.current.LineNum
	this.current = this.lexer.NextToken()

	//处理所有注解
	for this.current.Kind == TOKEN_AT {
		this.parseAnnotation()
		//this.eatToken(TOKEN_AT)
		//this.eatToken(TOKEN_ID)
		//if this.current.Kind == TOKEN_LPAREN {
		//	for this.current.Kind != TOKEN_RPAREN {
		//		this.advance()
		//	}
		//	this.advance()
		//}
	}

}

func (this *Parser) eatToken(kind int) {
	if kind == this.current.Kind {
		this.advance()
	} else if TOKEN_COMMENT == this.current.Kind {
		this.advance()
		this.eatToken(kind)
	} else {
		util.ParserError(tMap[kind], tMap[this.current.Kind], this.current.LineNum, this.lexer.fname)
	}
}
func (this *Parser) parseType() ast.Exp {
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
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			ele := this.parseNotExp()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
		} else {
			this.assignType = &ast.ObjectType{ast.TYPE_OBJECT}
			this.currentType = &ast.ListType{name, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_LIST}
		}

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
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			ele := this.parseNotExp()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
		} else {
			this.currentType = &ast.ListType{name, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_LIST}
		}

	case TOKEN_ARRAYLIST:
		//处理泛型
		name := this.current.Lexeme
		this.eatToken(TOKEN_ARRAYLIST)
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			ele := this.parseNotExp()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
		} else {
			this.currentType = &ast.ListType{name, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_LIST}
		}

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
			this.currentType = &ast.GenericType{ast.NewIdent(name, this.Linenum), tp, ast.TYPE_GENERIC}
		}
	}
	log.Debugf("解析类型:%s", this.currentType)
	return this.currentType
}

//解析泛型实例化参数列表
func (this *Parser) parseTypeList() (types []ast.Exp) {
	log.Infof("解析泛型参数列表")
	types = []ast.Exp{}
	if this.current.Kind == TOKEN_GT {
		return types
	}
	tp := this.parseType()
	types = append(types, tp)

	for this.current.Kind == TOKEN_COMMER {
		this.advance()
		tp := this.parseType()
		types = append(types, tp)
	}
	return types
}

//
//
// return:
func (this *Parser) parseFormalList(isSingle bool) (flist []ast.Field) {
	log.Debugf("解析函数参数")
	flist = []ast.Field{}
	var tp ast.Exp
	var id string
	var access int

	if this.TypeToken() ||
		this.current.Kind == TOKEN_ID {
		var nonType = false
		pre := this.current.Lexeme

		tp = this.parseNotExp()
		log.Debugf("解析函数 --> 需要类型推断 -> %v", this.current.Lexeme)
		//不用类型推断
		if this.current.Kind == TOKEN_ID {
			id = this.current.Lexeme
			id = GetNewId(id)
			this.eatToken(TOKEN_ID)
			flist = append(flist, ast.NewFieldSingle(access, tp, id, nil, false, false))

			//lambda等需要类型推断
			//TODO 类型推断
		} else {
			log.Debugf("解析函数 --> 需要类型推断")
			nonType = true
			flist = append(flist, ast.NewFieldSingle(access, &ast.ObjectType{}, id, nil, false, false))
		}
		if this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}

		for this.current.Kind == TOKEN_COMMER && !isSingle {
			this.eatToken(TOKEN_COMMER)
			if this.current.Kind == TOKEN_COMMENT {
				this.advance()
			}
			if nonType {
				log.Debugf("解析函数 --> 需要类型推断")
				pre = this.current.Lexeme
				pre = GetNewId(pre)
				this.eatToken(TOKEN_ID)
				flist = append(flist, ast.NewFieldSingle(access, &ast.ObjectType{}, pre, nil, false, false))
			} else {
				tp = this.parseNotExp()
				id = this.current.Lexeme
				id = GetNewId(id)
				this.eatToken(TOKEN_ID)
				flist = append(flist, ast.NewFieldSingle(access, tp, id, nil, false, false))
			}

		}
	}
	return flist
}

//强制类型转换
func (this *Parser) parseCastExp() ast.Exp {
	log.Infof("****************** - > %v", this.current.Lexeme)
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

// 成员(函数/变量)访问语句
//
// param: x
// return:
func (this *Parser) parseCallExp(x ast.Exp) (ret ast.Exp) {
	for this.current.Kind == TOKEN_DOT {
		this.eatToken(TOKEN_DOT)
		if this.current.Kind == TOKEN_LENGTH {
			this.advance()
			if this.current.Kind == TOKEN_LPAREN {
				this.eatToken(TOKEN_LPAREN)
				this.eatToken(TOKEN_RPAREN)
			}

			return ast.Length_new(x, this.Linenum)
		} else if this.current.Kind == TOKEN_SIZE {
			this.advance()
			this.eatToken(TOKEN_LPAREN)
			this.eatToken(TOKEN_RPAREN)
			return ast.Length_new(x, this.Linenum)
		}
		if this.current.Kind == TOKEN_CLASS {
			this.advance()
			return ast.ClassExp_new(x, this.Linenum)
		}

		x = ast.SelectorExpr_new(x, this.current.Lexeme, this.Linenum)

		this.eatToken(TOKEN_ID)
		if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			x = ast.CallExpr_new(x, args, this.Linenum)
		}
		if this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}
	}
	return x
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
		} else {
			this.ParseBug("加法解析bug")
		}

	case TOKEN_LPAREN:

		this.eatToken(TOKEN_LPAREN)
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
		if len(exps) == 1 && (this.current.Kind == TOKEN_ID || this.current.Kind == TOKEN_LPAREN) {
			//是类型表达式
			var tp ast.Exp
			switch v := exps[0].(type) {
			case ast.Type:
				tp = exps[0]
			case *ast.Id:
				tp = &ast.ClassType{
					Name:     v.Name,
					TypeKind: ast.TYPE_CLASS,
				}
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
			if this.current.Kind == TOKEN_LPAREN {
				this.eatToken(TOKEN_LPAREN)
				exp := this.parseExp()
				this.eatToken(TOKEN_RPAREN)
				return ast.Cast_new(tp, exp, this.Linenum)
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
		return &ast.This{}
	case TOKEN_SYSTEM:
		var x ast.Exp
		x = ast.Id_new(this.current.Lexeme, nil, false, this.Linenum)
		this.eatToken(TOKEN_SYSTEM)
		x = this.parseCallExp(x)
		return x
	case TOKEN_ID:
		id := this.current.Lexeme
		newId := GetNewId(id)
		this.advance()
		log.Infof("------------->TOKEN_ID-->%v", this.current.Lexeme)
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
			log.Infof("------------->函数调用-->%v", id)
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			m := ast.NewIdent(newId, this.Linenum)
			return ast.CallExpr_new(m, args, this.Linenum)

		} else if this.current.Kind == TOKEN_LAMBDA {
			log.Debugf("发现 单参数 TOKEN_LAMBDA")
			this.advance()
			return this.parseLambdaExp([]ast.Exp{ast.NewIdent(newId, this.Linenum)})
		}
		log.Debugf("适配变量ID->%s", id)
		return ast.NewIdent(id, this.Linenum)
	case TOKEN_NEW:
		return this.parseNewExp()
	default:
		if this.TypeToken() {
			return this.parseType()
		}
		this.ParseBug("未处理")
		os.Exit(0)
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
	this.advance()
	switch this.current.Kind {
	case TOKEN_INT:
		this.advance()
		this.eatToken(TOKEN_LBRACK)
		exp := this.parseExp()
		this.eatToken(TOKEN_RBRACK)
		return ast.NewIntArray_new(exp, this.Linenum)
	case TOKEN_STRING:
		this.advance()
		this.eatToken(TOKEN_LBRACK)
		exp := this.parseExp()
		this.eatToken(TOKEN_RBRACK)
		return ast.NewStringArray_new(exp, this.Linenum)
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
			t, ok := this.currentType.(*ast.HashType)
			if ok {
				key = t.Key
				ele = t.Value
			}
		}
		this.eatToken(TOKEN_GT)
		this.eatToken(TOKEN_LPAREN)
		if this.current.Kind == TOKEN_NUM {
			//FIXME 处理map容量
			this.eatToken(TOKEN_NUM)
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
				ele = this.assignType.(*ast.ListType).Ele
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
				ele = this.assignType.(*ast.ListType).Ele
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
		var typeName string
		var args []ast.Exp
		log.Debugf("-------------> %v", this.current.Lexeme)
		exp := this.parseNotExp()
		log.Debugf("-------------> %v", this.current.Lexeme)
		switch v := exp.(type) {
		case *ast.Id:
			typeName = v.Name
		case *ast.Ident:
			typeName = v.Name
		case *ast.SelectorExpr:
			typeName = v.Sel
		case *ast.CallExpr:
			//调用构造函数
			return v
		}
		//模板
		if this.current.Kind == TOKEN_LT {
			this.eatToken(TOKEN_LT)
			if this.current.Kind != TOKEN_GT {
				this.parseNotExp()
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
			size := this.parseNotExp()
			this.eatToken(TOKEN_RBRACK)
			if this.current.Kind == TOKEN_LBRACE {
				this.eatToken(TOKEN_LBRACE)
				eles := this.parseNotExp()
				this.eatToken(TOKEN_RBRACE)
				return ast.NewObjectArray_new(&ast.ClassType{typeName, ast.TYPE_CLASS}, eles, size, this.Linenum)
			}
			return ast.NewObjectArray_new(&ast.ClassType{typeName, ast.TYPE_CLASS}, nil, size, this.Linenum)
		}
		return ast.NewObjectWithArgsList_new(&ast.ClassType{typeName, ast.TYPE_CLASS}, args, this.Linenum)
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
		args = append(args, ast.Id_new(id, exp, false, this.Linenum))

		for this.current.Kind == TOKEN_COMMER {
			this.advance()
			exp = this.parseExp()
			id = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			args = append(args, ast.Id_new(id, exp, false, this.Linenum))
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
		case *ast.Id:
			fields = append(fields, ast.NewFieldSingle(0, v.Tp, v.Name, nil, false, false))
		case *ast.Ident:
			fields = append(fields, ast.NewFieldSingle(0, &ast.ObjectType{ast.TYPE_OBJECT}, v.Name, nil, false, false))
		default:
			panic("parseLambdaExp")
		}
	}

	if this.current.Kind == TOKEN_LBRACE {
		this.eatToken(TOKEN_LBRACE)
		stms := this.parseStatements()
		this.eatToken(TOKEN_RBRACE)
		return ast.Lambda_new(fields, stms, this.Linenum)
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
		this.current.Kind == TOKEN_COMMENT ||
		//FIXME 自增,自减作为语句处理
		//this.current.Kind == TOKEN_AUTOSUB || //后缀加
		//this.current.Kind == TOKEN_AUTOADD || //后缀减
		this.current.Kind == TOKEN_LBRACK {
		switch this.current.Kind {

		//可以不断循环下去
		case TOKEN_DOT:
			log.Debugf("解析函数调用,或成员变量")
			exp = this.parseCallExp(exp)
			//数组索引操作
		case TOKEN_LBRACK: //[exp]
			this.advance()
			index := this.parseExp()
			if index != nil {
				log.Debugf("数组索引表达式 --> %v", exp)
				this.eatToken(TOKEN_RBRACK)
				if this.current.Kind == TOKEN_DOT {
					exp = ast.ArraySelect_new(exp, index, this.Linenum)
				} else {
					return ast.ArraySelect_new(exp, index, this.Linenum)
				}
			} else {
				log.Debugf("数组索引用")
				this.eatToken(TOKEN_RBRACK)
				panic("数组索引用")
			}
		case TOKEN_COMMENT:
			exp = this.parseAtomExp()
		}
	}
	return exp
}

//TimesExp  -> !TimesExp
//          -> NotExp
func (this *Parser) parseTimeExp() ast.Exp {
	log.Debugf("解析 parseTimeExp")
	var exp2 ast.Exp
	var opt = this.current.Kind
	for this.current.Kind == TOKEN_NOT {
		this.advance()
		exp2 = this.parseTimeExp()
	}
	if exp2 != nil {
		switch opt {
		case TOKEN_NOT:
			return ast.Not_new(exp2, this.Linenum)
		//case TOKEN_AUTOADD:
		//	return ast.AutoAdd_new(nil, exp2, this.Linenum)
		//case TOKEN_AUTOSUB:
		//	return ast.AutoSub_new(nil, exp2, this.Linenum)
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
		this.advance()
	}
	for this.current.Kind == TOKEN_REM ||
		this.current.Kind == TOKEN_MUL ||
		this.current.Kind == TOKEN_QUO {
		switch this.current.Kind {
		case TOKEN_MUL:
			this.advance()
			//去除注释
			for this.current.Kind == TOKEN_COMMENT {
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
		this.advance()
	}
	for this.current.Kind == TOKEN_ADD ||
		this.current.Kind == TOKEN_SUB {
		switch this.current.Kind {
		case TOKEN_ADD:
			this.advance()
			//去除注释
			for this.current.Kind == TOKEN_COMMENT {
				this.advance()
			}

			right := this.parseAddSubExp()
			left = ast.Add_new(left, right, this.Linenum)
		case TOKEN_SUB:
			this.advance()
			//去除注释
			for this.current.Kind == TOKEN_COMMENT {
				this.advance()
			}
			right := this.parseAddSubExp()
			left = ast.Sub_new(left, right, this.Linenum)
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
	log.Debugf("解析 parseLAndExp --> %v", this.current.Lexeme)
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
	left := this.parseQuestionExp()
	//
	if id, ok := left.(*ast.Id); ok {
		//说明是成员变量或成员函数
		if (nil != this.currentClass.GetField(id.Name) && nil == this.currentMethod.GetFormal(id.Name)) || nil != this.currentClass.GetMethod(id.Name) {
			id.Name = "this." + Capitalize(id.Name)
		}
	}
	for this.current.Kind == TOKEN_QUESTION {
		log.Debugf("发现TOKEN_QUESTION")
		this.advance()
		log.Infof("TOKEN_QUESTION --> 解析第一个表达式")
		one := this.parseQuestionExp()
		this.eatToken(TOKEN_COLON)
		log.Infof("TOKEN_QUESTION --> 解析第二个表达式")
		two := this.parseQuestionExp()
		return ast.Question_new(left, one, two, this.Linenum)
	}
	if this.current.Kind == TOKEN_INSTANCEOF {
		this.eatToken(TOKEN_INSTANCEOF)
		right := this.parseQuestionExp()
		return ast.Instanceof_new(left, right, this.Linenum)
	}
	return left
}

//
//
// return:
func (this *Parser) parseStatement() ast.Stm {
	log.Debugf("*******解析代码段******* --> %v", this.current.Lexeme)
	switch this.current.Kind {
	case TOKEN_IS_NULL:
		fallthrough
	case TOKEN_NOT_EMPTY:
		fallthrough
	case TOKEN_HAS_TEXT:
		fallthrough
	case TOKEN_NOT_NULL:
		fallthrough
	case TOKEN_IS_TRUE:
		fallthrough
	case TOKEN_ASSERT:
		return this.parseAssertExp()
	case TOKEN_COMMENT:
		stm := ast.Comment_new(this.current.Lexeme, this.Linenum)
		this.advance()
		return stm

	case TOKEN_LBRACE: //{
		log.Debugf("*******解析代码段*******")
		this.eatToken(TOKEN_LBRACE)
		stms := this.parseStatements()
		this.eatToken(TOKEN_RBRACE)
		return ast.Block_new(stms, this.Linenum)
	case TOKEN_THIS:
		exp := this.parseExp()
		if this.current.Kind == TOKEN_ASSIGN {
			this.eatToken(TOKEN_ASSIGN)
			right := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			assign := new(ast.Assign)
			assign.Left = exp
			assign.Value = right
			//三元表达式
			if _, ok := right.(*ast.Question); ok {
				assign.SetTriple()
			}
			return assign

		}
		this.eatToken(TOKEN_SEMI)
		exprStm := ast.ExprStm_new(exp, this.Linenum)

		return exprStm
	case TOKEN_ID:
		//1 调用函数表达式
		//2 临时变量声明
		//3 赋值语句左边表达式
		//4 泛型类型声明
		id := this.current.Lexeme
		exp := this.parseNotExp()
		switch this.current.Kind {
		//处理声明临时变量和赋值语句
		case TOKEN_ID:
			log.Debugf("*******解析临时变量声明语句*******")
			id = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			decl := ast.DeclStmt_new(nil, exp, nil, this.Linenum)
			decl.Names = append(decl.Names, ast.NewIdent(id, this.Linenum))
			//有赋值语句
			if this.current.Kind == TOKEN_ASSIGN {
				//临时变量类型
				log.Debugf("*******解析临时变量声明语句(有赋值语句)*******")
				this.eatToken(TOKEN_ASSIGN)
				exp := this.parseExp()
				//三元表达式
				if _, ok := exp.(*ast.Question); ok {
					decl.SetTriple()
				}
				decl.Values = append(decl.Values, exp)
			}

			//定义多个变量
			for this.current.Kind == TOKEN_COMMER {
				this.advance()
				id = this.current.Lexeme
				this.eatToken(TOKEN_ID)
				decl.Names = append(decl.Names, ast.NewIdent(id, this.Linenum))

				if this.current.Kind == TOKEN_ASSIGN {
					//临时变量类型
					log.Debugf("*******解析临时变量声明语句(有赋值语句)*******")
					this.eatToken(TOKEN_ASSIGN)
					exp := this.parseExp()
					//三元表达式
					if _, ok := exp.(*ast.Question); ok {
						decl.SetTriple()
					}
					decl.Values = append(decl.Values, exp)
				}

			}
			this.eatToken(TOKEN_SEMI)
			return decl

		case TOKEN_LPAREN:
			//直接函数调用语句
			fallthrough
		case TOKEN_DOT:
			log.Debugf("*******解析函数调用*******")
			exp := this.parseExp()
			//有赋值语句
			if this.current.Kind == TOKEN_ASSIGN {

				//临时变量类型
				log.Debugf("*******解析临时变量声明语句(有赋值语句)*******")
				this.eatToken(TOKEN_ASSIGN)
				right := this.parseExp()
				this.eatToken(TOKEN_SEMI)
				assign := ast.Assign_new(exp, right, false, this.Linenum)
				//三元表达式
				if _, ok := right.(*ast.Question); ok {
					assign.SetTriple()
				}
				return assign
			} else {
				this.eatToken(TOKEN_SEMI)
				exprStm := ast.ExprStm_new(exp, this.Linenum)
				//检查表达式是不是三元表达式
				if fn, ok := exp.(*ast.CallExpr); ok {
					for _, v := range fn.ArgsList {
						//输入参数有三元表达式
						if _, ok := v.(*ast.Question); ok {
							exprStm.SetTriple()
						}
					}
				}
				return exprStm

			}

		case TOKEN_ASSIGN:
			this.eatToken(TOKEN_ASSIGN)
			exp := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			assign := new(ast.Assign)

			//说明是成员变量
			if this.currentClass.GetField(id) != nil {
				id = "this." + id
			}
			assign.Left = ast.Id_new(id, nil, false, this.Linenum)
			assign.Value = exp
			//三元表达式
			if _, ok := exp.(*ast.Question); ok {
				assign.SetTriple()
			}
			return assign

		case TOKEN_QUO_ASSIGN:
			left := ast.Id_new(id, nil, false, this.Linenum)
			this.eatToken(TOKEN_QUO_ASSIGN)
			right := this.parseExp()
			this.eatToken(TOKEN_SEMI)

			return ast.Binary_new(left, right, "/=", this.Linenum)
		case TOKEN_MUL_ASSIGN:

			left := ast.Id_new(id, nil, false, this.Linenum)
			this.eatToken(TOKEN_MUL_ASSIGN)
			right := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			return ast.Binary_new(left, right, "*=", this.Linenum)
		case TOKEN_SUB_ASSIGN:

			left := ast.Id_new(id, nil, false, this.Linenum)
			this.eatToken(TOKEN_SUB_ASSIGN)
			right := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			return ast.Binary_new(left, right, "-=", this.Linenum)
		case TOKEN_ADD_ASSIGN:
			left := ast.Id_new(id, nil, false, this.Linenum)
			this.eatToken(TOKEN_ADD_ASSIGN)
			right := this.parseExp()
			this.eatToken(TOKEN_SEMI)

			return ast.Binary_new(left, right, "+=", this.Linenum)
		case TOKEN_REM_ASSIGN:
			left := ast.Id_new(id, nil, false, this.Linenum)
			this.eatToken(TOKEN_REM_ASSIGN)
			right := this.parseExp()
			this.eatToken(TOKEN_SEMI)

			return ast.Binary_new(left, right, "%=", this.Linenum)

			//处理的是后缀加
		case TOKEN_AUTOADD:
			log.Debugf("处理累加")
			this.eatToken(TOKEN_AUTOADD)
			//特殊的for语句才不需要分号
			if !this.isSpecial {
				this.isSpecial = false
				this.eatToken(TOKEN_SEMI)
			}
			left := ast.Id_new(id, nil, false, this.Linenum)

			return ast.Binary_new(left, &ast.Num{Value: 1}, "+=", this.Linenum)
			//处理的是后缀减
		case TOKEN_AUTOSUB:
			this.eatToken(TOKEN_AUTOSUB)
			if !this.isSpecial {
				this.isSpecial = false
				this.eatToken(TOKEN_SEMI)
			}
			left := ast.Id_new(id, nil, false, this.Linenum)
			return ast.Binary_new(left, &ast.Num{Value: 1}, "-=", this.Linenum)
		case TOKEN_LBRACK:
			this.eatToken(TOKEN_LBRACK) //[
			index := this.parseExp()
			this.eatToken(TOKEN_RBRACK) //]
			this.eatToken(TOKEN_ASSIGN)
			exp := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			return ast.AssignArray_new(id, index, exp, nil, false, this.Linenum)
		case TOKEN_LT:
			this.eatToken(TOKEN_LT)
			tp := this.parseType()
			this.eatToken(TOKEN_GT)
			id = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			decl := ast.DeclStmt_new(nil, tp, nil, this.Linenum)
			decl.Names = append(decl.Names, ast.NewIdent(id, this.Linenum))
			//有赋值语句
			if this.current.Kind == TOKEN_ASSIGN {
				//临时变量类型
				log.Debugf("*******解析临时变量声明语句(有赋值语句)*******")
				this.eatToken(TOKEN_ASSIGN)
				exp := this.parseExp()
				//三元表达式
				if _, ok := exp.(*ast.Question); ok {
					decl.SetTriple()
				}
				decl.Values = append(decl.Values, exp)
			}
			this.eatToken(TOKEN_SEMI)
			return decl

		case TOKEN_SEMI:
			this.eatToken(TOKEN_SEMI)
			return ast.ExprStm_new(exp, this.Linenum)
		case TOKEN_COMMENT:
			this.advance()
		default:
			this.ParseBug("parseStatement代码段解析bug")
		}
	case TOKEN_IF:
		log.Debugf("********TOKEN_IF***********")
		this.eatToken(TOKEN_IF)
		this.eatToken(TOKEN_LPAREN)
		condition := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		body := this.parseStatement()

		//不是block,说明没有大括号
		if _, ok := body.(*ast.Block); !ok {
			body = ast.Block_new([]ast.Stm{body}, this.Linenum)
		}

		if this.current.Kind == TOKEN_ELSE {
			this.eatToken(TOKEN_ELSE)
			elsee := this.parseStatement()
			if _, ok := elsee.(*ast.Block); !ok {
				elsee = ast.Block_new([]ast.Stm{elsee}, this.Linenum)
			}
			return ast.If_new(condition, body, elsee, this.Linenum)
		} else {
			return ast.If_new(condition, body, nil, this.Linenum)
		}
	case TOKEN_TRY:
		log.Debugf("********TOKEN_TRY***********")
		this.eatToken(TOKEN_TRY)
		body := this.parseStatement()
		var catches []*ast.Catch
		var finally ast.Stm
		for this.current.Kind == TOKEN_CATCH {
			this.eatToken(TOKEN_CATCH)
			this.eatToken(TOKEN_LPAREN)
			var excepts []ast.Exp
			e := this.parseNotExp()
			excepts = append(excepts, e)

			//处理多个异常共用一个处理方法的表达式
			for this.current.Kind == TOKEN_OR {
				this.advance()
				e := this.parseNotExp()
				excepts = append(excepts, e)
			}
			id := this.current.Lexeme
			this.eatToken(TOKEN_ID)

			this.eatToken(TOKEN_RPAREN)
			catchBody := this.parseStatement()
			catch := ast.Catch_new(excepts, id, catchBody, this.Linenum)
			catches = append(catches, catch)
		}

		if this.current.Kind == TOKEN_FINALLY {
			this.eatToken(TOKEN_FINALLY)
			finally = this.parseStatement()
		}
		return ast.Try_new(body, catches, finally, this.Linenum)
	case TOKEN_WHILE:
		log.Debugf("********TOKEN_WHILE***********")
		this.eatToken(TOKEN_WHILE)
		this.eatToken(TOKEN_LPAREN)
		exp := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		body := this.parseStatement()
		return ast.While_new(exp, body, this.Linenum)

	case TOKEN_SWITCH:
		log.Debugf("********TOKEN_SWITCH***********")
		this.eatToken(TOKEN_SWITCH)
		this.eatToken(TOKEN_LPAREN)
		exp := this.parseExp()
		this.eatToken(TOKEN_RPAREN)

		this.eatToken(TOKEN_LBRACE)
		body := this.parseStatements()
		this.eatToken(TOKEN_RBRACE)
		return ast.Switch_new(exp, ast.Block_new(body, this.Linenum), this.Linenum)
	case TOKEN_CASE:
		this.eatToken(TOKEN_CASE)
		exp := this.parseExp()
		this.eatToken(TOKEN_COLON)
		body := this.parseStatements()
		return ast.Case_new(exp, ast.Block_new(body, this.Linenum), this.Linenum)
	case TOKEN_FOR:
		log.Debugf("********TOKEN_FOR***********")
		this.eatToken(TOKEN_FOR)
		this.eatToken(TOKEN_LPAREN)
		var exp = this.parseExp()
		var id ast.Exp
		var Init ast.Stm
		//说明是声明语句
		if this.current.Kind == TOKEN_ID {
			id = this.parseExp()
		}

		//for循环三段式
		if this.current.Kind == TOKEN_ASSIGN {
			log.Debugf("********TOKEN_FOR--> 解析初始化语句 ***********")
			//临时变量类型
			if id != nil {
				this.eatToken(TOKEN_ASSIGN)
				value := this.parseExp()
				this.eatToken(TOKEN_SEMI)
				decl := ast.DeclStmt_new(nil, exp, nil, this.Linenum)

				decl.Names = append(decl.Names, id)
				decl.Values = append(decl.Values, value)
				Init = decl
			} else {
				this.eatToken(TOKEN_ASSIGN)
				value := this.parseExp()
				this.eatToken(TOKEN_SEMI)
				Init = ast.Assign_new(exp, value, false, this.Linenum)

			}

			//
			log.Debugf("********TOKEN_FOR--> 解析条件语句 ***********")
			Condition := this.parseExp()
			this.eatToken(TOKEN_SEMI)

			log.Debugf("********TOKEN_FOR--> 解析更新语句 ***********")
			this.isSpecial = true
			Post := this.parseStatement()
			this.eatToken(TOKEN_RPAREN)
			body := this.parseStatement()
			return ast.For_new(Init, Condition, Post, body, this.Linenum)

			//枚举式
		} else if this.current.Kind == TOKEN_COLON {
			log.Debugf("*******for循环枚举*************")
			this.eatToken(TOKEN_COLON)
			var right ast.Exp

			//处理强制类型转换
			if this.current.Kind == TOKEN_LPAREN {
				right = this.parseCastExp()
			} else {
				right = this.parseLOrExp()
			}
			this.eatToken(TOKEN_RPAREN)

			body := this.parseStatement()

			if id != nil {
				exp = id
			}

			return ast.Range_new(exp, right, body, this.Linenum)
		}

	case TOKEN_SYSTEM:
		this.eatToken(TOKEN_SYSTEM)
		this.eatToken(TOKEN_DOT)
		this.eatToken(TOKEN_OUT)
		this.eatToken(TOKEN_DOT)
		this.eatToken(TOKEN_PRINTLN)
		this.eatToken(TOKEN_LPAREN)
		e := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		this.eatToken(TOKEN_SEMI)
		return ast.Print_new(e, this.Linenum)
	case TOKEN_THROW:
		this.eatToken(TOKEN_THROW)
		e := this.parseExp()
		this.eatToken(TOKEN_SEMI)
		return ast.Throw_new(e, this.Linenum)
	case TOKEN_RETURN:
		this.eatToken(TOKEN_RETURN)
		//空return
		if this.current.Kind == TOKEN_SEMI {
			this.eatToken(TOKEN_SEMI)
			return ast.Return_new(nil, this.Linenum)
		}
		log.Debugf("------>解析return,%v", this.current.Lexeme)
		e := this.parseExp()
		this.eatToken(TOKEN_SEMI)
		return ast.Return_new(e, this.Linenum)
	default:
		if this.TypeToken() {
			tp := this.parseType()
			id := this.current.Lexeme
			id = GetNewId(id)

			this.eatToken(TOKEN_ID)
			decl := ast.DeclStmt_new(nil, tp, nil, this.Linenum)

			decl.Names = append(decl.Names, ast.NewIdent(id, this.Linenum))
			//有赋值语句
			if this.current.Kind == TOKEN_ASSIGN {
				this.assignType = tp
				//临时变量类型
				log.Debugf("*******解析临时变量声明语句(有赋值语句)*******")
				this.eatToken(TOKEN_ASSIGN)
				exp := this.parseExp()
				//三元表达式
				if _, ok := exp.(*ast.Question); ok {
					decl.SetTriple()
				}
				decl.Values = append(decl.Values, exp)
			} else {
				log.Debugf("*******解析临时变量声明语句(无赋值语句)*******")
			}
			this.eatToken(TOKEN_SEMI)
			return decl

		}
		this.ParseBug("代码段解析bug")
	}
	return nil
}

func (this *Parser) parseStatements() []ast.Stm {
	stms := []ast.Stm{}
	for this.TypeToken() ||
		this.ExtraToken() ||
		this.current.Kind == TOKEN_ID ||
		this.current.Kind == TOKEN_LBRACE ||
		this.current.Kind == TOKEN_COMMENT ||
		this.current.Kind == TOKEN_IF ||
		this.current.Kind == TOKEN_TRY ||
		this.current.Kind == TOKEN_WHILE ||
		this.current.Kind == TOKEN_FOR ||
		this.current.Kind == TOKEN_THROW ||
		this.current.Kind == TOKEN_RETURN ||
		this.current.Kind == TOKEN_THIS ||
		this.current.Kind == TOKEN_SWITCH ||
		this.current.Kind == TOKEN_CASE ||
		this.current.Kind == TOKEN_SYSTEM {
		log.Infof("****************** parseStatements **********************-->%v", this.current.Lexeme)
		stms = append(stms, this.parseStatement())
	}
	return stms
}
func (this *Parser) parseMemberVarDecl(tmp *ast.FieldSingle, IsStatic bool) (dec ast.Field) {
	var value ast.Exp
	if this.current.Kind == TOKEN_ASSIGN {
		this.eatToken(TOKEN_ASSIGN)
		value = this.parseExp()
	}
	this.eatToken(TOKEN_SEMI)
	dec = ast.NewFieldSingle(tmp.Access, tmp.Tp, tmp.Name, value, IsStatic, true)
	return dec
}

// 解析类上下文
//
// return:
func (this *Parser) parseClassContext(classSingle *ast.ClassSingle) {

	//每次循环解析一个成员变量或一个成员函数
	for this.TypeToken() ||
		this.current.Kind == TOKEN_ID ||
		this.current.Kind == TOKEN_PRIVATE ||
		this.current.Kind == TOKEN_PUBLIC ||
		this.current.Kind == TOKEN_PROTECTED ||
		this.current.Kind == TOKEN_FINAL ||
		this.current.Kind == TOKEN_COMMENT ||
		this.current.Kind == TOKEN_STATIC {

		log.Debugf("解析类上下文...... -- >%v", this.current.Lexeme)
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
			//log.Infof("注释-->%v", comment)
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

			//处理类构造函数
			if this.currentClass.GetName() == this.current.Lexeme && prefix == false {
				log.Infof("处理构造函数-->%v", this.current.Lexeme)
				IsConstruct = true
				tmp.Tp = &ast.Void{ast.TYPE_VOID}
				//变量/函数名
				tmp.Name = "New" + this.current.Lexeme
			} else {
				//类型
				tmp.Tp = this.parseType()
				this.assignType = tmp.Tp
				//变量/函数名
				tmp.Name = this.current.Lexeme
			}

			this.eatToken(TOKEN_ID)

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
	log.Debugf("*******解析成员函数*******")
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

// 类静态语句
//
// param: comment
// return:
func (this *Parser) parseMemberStatic(comment string) (meth ast.Method) {

	this.currentMethod = ast.NewMethodSingle(&ast.Void{}, "init", nil, nil, false, true, false, comment)

	var stms []ast.Stm
	this.eatToken(TOKEN_LBRACE)
	//解析本地变量和表达式
	stms = this.parseStatements()

	this.eatToken(TOKEN_RBRACE)

	return ast.NewMethodSingle(&ast.Void{}, "init", nil, stms, false, true, false, comment)
}
func (this *Parser) parseEnumDecl(access int) (cl ast.Class) {
	var id, extends string

	this.eatToken(TOKEN_ENUM)
	id = this.current.Lexeme
	this.eatToken(TOKEN_ID)
	//处理implements
	if this.current.Kind == TOKEN_IMPLEMENTS {
		this.eatToken(TOKEN_IMPLEMENTS)
		extends = this.current.Lexeme
		this.eatToken(TOKEN_ID)
	}
	this.eatToken(TOKEN_LBRACE)
	classSingle := ast.NewClassSingle(access, id, extends, true)
	this.currentClass = classSingle
	//处理枚举变量
	for {
		var comment string
		//处理注释
		if this.current.Kind == TOKEN_COMMENT {
			comment = ""
			for this.current.Kind == TOKEN_COMMENT {
				comment += "\n" + this.current.Lexeme
				this.advance()
			}
			if this.current.Kind == TOKEN_EOF || (this.current.Kind != TOKEN_PRIVATE && this.current.Kind != TOKEN_PUBLIC && this.current.Kind != TOKEN_PROTECTED && this.current.Kind != TOKEN_ID) {
				return
			}
			log.Infof("注释-->%v", comment)
		}

		id = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		//FIXME 只支持一个值枚举
		if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			l := this.parseExpList()
			value := l[0]
			classSingle.AddField(ast.NewFieldSingle(access, nil, id, value, false, true))
			this.eatToken(TOKEN_RPAREN)
		} else {
			classSingle.AddField(ast.NewFieldSingle(access, nil, id, nil, false, true))
		}

		for this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}
		if this.current.Kind == TOKEN_SEMI {
			this.eatToken(TOKEN_SEMI)
			break
		} else {

			this.eatToken(TOKEN_COMMER)
		}

	}
	//
	this.parseClassContext(ast.NewClassSingle(access, id, extends, false))

	this.eatToken(TOKEN_RBRACE)

	return classSingle
}

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
	classSingle := ast.NewClassSingle(access, id, extends, false)
	this.currentClass = classSingle

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
func (this *Parser) parseAnnotation() {
	this.eatToken(TOKEN_AT)
	this.eatToken(TOKEN_ID)
	//带参数的注解
	if this.current.Kind == TOKEN_LPAREN {
		this.eatToken(TOKEN_LPAREN)

		if this.current.Kind == TOKEN_RPAREN {
			this.eatToken(TOKEN_RPAREN)
			return
		}
		for {
			this.parseNotExp() //id
			if this.current.Kind == TOKEN_ASSIGN {
				this.advance()     // =
				this.parseNotExp() //id
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

func (this *Parser) parseMainClass() ast.MainClass {
	//
	this.eatToken(TOKEN_CLASS)
	id := this.current.Lexeme
	this.eatToken(TOKEN_ID)
	this.eatToken(TOKEN_LBRACE)
	this.eatToken(TOKEN_PUBLIC)
	this.eatToken(TOKEN_STATIC)
	this.eatToken(TOKEN_VOID)
	this.eatToken(TOKEN_MAIN)
	this.eatToken(TOKEN_LPAREN)
	this.eatToken(TOKEN_STRING)
	this.eatToken(TOKEN_LBRACK)
	this.eatToken(TOKEN_RBRACK)
	arg := this.current.Lexeme
	this.eatToken(TOKEN_ID)
	this.eatToken(TOKEN_RPAREN)
	this.eatToken(TOKEN_LBRACE)
	stm := this.parseStatement()
	this.eatToken(TOKEN_RBRACE)
	this.eatToken(TOKEN_RBRACE)
	return &ast.MainClassSingle{id, arg, stm}
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
			} else if this.current.Kind == TOKEN_MUL {
				this.eatToken(TOKEN_MUL)
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
		log.Warnf("导入:%v --> %v", path, id)
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

func (this *Parser) Parser() ast.File {
	p := this.parseProgram()
	return p
}

func GetNewId(id string) (nId string) {
	if id == "map" {
		nId = "oMap"
	} else if id == "type" {
		nId = "oType"
	} else {
		nId = id
	}

	return
}
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				//log.Info("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func (this *Parser) ParseBug(info string) {
	var msg = fmt.Sprintf("[%v] %s:%d:%s\n", this.current.Lexeme, path.Base(this.lexer.fname), this.Linenum, info)
	util.Bug(msg)
}
