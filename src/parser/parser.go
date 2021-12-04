package parser

import (
	"dog/ast"
	"dog/control"
	"dog/util"
	log "github.com/corgi-kx/logcustom"
	"strconv"
)

type Parser struct {
	lexer       *Lexer
	current     *Token
	pending     []*Token
	currentNext *Token
	currentType ast.Type
	assignType  ast.Type
	isSpecial   bool
	isField     bool
	fpBak       int    //用于记录测试前的指针
	currentBak  *Token //用于记录测试前的token
	Linenum     int
}

func NewParse(fname string, buf []byte) *Parser {
	lexer := NewLexer(fname, buf)
	p := new(Parser)
	p.lexer = lexer
	p.current = p.lexer.NextToken()
	return p
}

func (this *Parser) TestIn() {
	this.fpBak = this.lexer.fp
	this.currentBak = this.current
}

func (this *Parser) TestOut() {
	this.lexer.fp = this.fpBak
	this.current = this.currentBak
}

func (this *Parser) advance() {
	if control.Lexer_dump == true {
		log.Info(this.current.ToString())
	}
	this.Linenum = this.current.LineNum
	this.current = this.lexer.NextToken()

	//处理注解
	if this.current.Kind == TOKEN_AT {
		this.eatToken(TOKEN_AT)
		this.eatToken(TOKEN_ID)
		if this.current.Kind == TOKEN_LPAREN {
			for this.current.Kind != TOKEN_RPAREN {
				this.advance()
			}
			this.advance()
		}

	}

}

func (this *Parser) eatToken(kind int) {
	if kind == this.current.Kind {
		this.advance()
	} else {
		util.ParserError(tMap[kind], tMap[this.current.Kind], this.current.LineNum)
	}
}
func (this *Parser) parseType() ast.Type {
	switch this.current.Kind {
	case TOKEN_INT:
		this.eatToken(TOKEN_INT)
		if this.current.Kind == TOKEN_LBRACK {
			this.eatToken(TOKEN_LBRACK)
			this.eatToken(TOKEN_RBRACK)
			this.currentType = &ast.IntArray{ast.TYPE_INTARRAY}
		} else {
			this.currentType = &ast.Int{}
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
			this.currentType = &ast.Integer{}
		}
	case TOKEN_VOID:
		this.eatToken(TOKEN_VOID)
		this.currentType = &ast.Void{ast.TYPE_VOID}
	case TOKEN_BOOLEAN:
		this.eatToken(TOKEN_BOOLEAN)
		this.currentType = &ast.Boolean{ast.TYPE_BOOLEAN}
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
		ele := this.parseType()
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
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
		this.eatToken(TOKEN_LT)
		ele := this.parseType()
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.ListType{name, ele, ast.TYPE_LIST}
	case TOKEN_ARRAYLIST:
		//处理泛型
		name := this.current.Lexeme
		this.eatToken(TOKEN_ARRAYLIST)
		this.eatToken(TOKEN_LT)
		ele := this.parseType()
		this.eatToken(TOKEN_ID)
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
			log.Infof("****************解析map***********************")
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
			this.TestIn()
			this.eatToken(TOKEN_LT)
			tp := this.parseType()
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.GenericType{name, tp, ast.TYPE_GENERIC}
		}
	}
	log.Infof("解析类型:%s", this.currentType.String())
	return this.currentType
}

//
//
// return:
func (this *Parser) parseFormalList() (flist []ast.Field) {
	log.Infof("解析函数参数")
	flist = []ast.Field{}
	var tp ast.Type
	var id string
	var access int

	if this.current.Kind == TOKEN_ID ||
		this.current.Kind == TOKEN_INT ||
		this.current.Kind == TOKEN_LONG ||
		this.current.Kind == TOKEN_STRING ||
		this.current.Kind == TOKEN_LIST ||
		this.current.Kind == TOKEN_INTEGER ||
		this.current.Kind == TOKEN_MAP ||
		this.current.Kind == TOKEN_BOOLEAN {
		var nonType = false
		pre := this.current.Lexeme

		tp = this.parseType()

		//不用类型推断
		if this.current.Kind == TOKEN_ID {
			id = this.current.Lexeme
			id = GetNewId(id)
			this.eatToken(TOKEN_ID)
			flist = append(flist, &ast.FieldSingle{access, tp, id, this.isField, nil})
			//lambda等需要类型推断
			//TODO 类型推断
		} else {
			log.Infof("解析函数 --> 需要类型推断")
			nonType = true
			flist = append(flist, &ast.FieldSingle{access, &ast.ObjectType{}, pre, this.isField, nil})
		}

		for this.current.Kind == TOKEN_COMMER {
			this.eatToken(TOKEN_COMMER)
			if nonType {
				log.Infof("解析函数 --> 需要类型推断")
				pre = this.current.Lexeme
				pre = GetNewId(pre)
				this.eatToken(TOKEN_ID)
				flist = append(flist, &ast.FieldSingle{access, &ast.ObjectType{}, pre, this.isField, nil})
			} else {
				tp = this.parseType()
				id = this.current.Lexeme
				id = GetNewId(id)
				this.eatToken(TOKEN_ID)
				flist = append(flist, &ast.FieldSingle{access, tp, id, this.isField, nil})
			}

		}
	}
	return flist
}

//强制类型转换
func (this *Parser) parseCastExp() ast.Exp {
	switch this.current.Kind {
	case TOKEN_LPAREN:
		this.advance()
		log.Infof("parseCastExp:%v", this.Linenum)
		tp := this.parseType()
		this.eatToken(TOKEN_RPAREN)

		exp := this.parseExp()

		return ast.Cast_new(tp, exp, this.Linenum)
	}
	return nil
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
	log.Infof("解析 parseAtomExp")
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
			//return &ast.Num{s, nil}
			return n
		} else {
			panic("error")
		}
	case TOKEN_LPAREN:
		this.advance()
		exp := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		return exp
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
		for this.current.Kind == TOKEN_DOT {
			this.eatToken(TOKEN_DOT)
			x = ast.SelectorExpr_new(x, this.current.Lexeme, this.Linenum)

			this.eatToken(TOKEN_ID)
			if this.current.Kind == TOKEN_LPAREN {
				this.eatToken(TOKEN_LPAREN)
				args := this.parseExpList()
				this.eatToken(TOKEN_RPAREN)
				x = ast.CallExpr_new(x, args, this.Linenum)
			}
		}
		return x
	case TOKEN_ID:
		id := this.current.Lexeme
		id = GetNewId(id)
		this.TestIn()
		this.advance()
		//声明一个临时变量的语句
		if this.current.Kind == TOKEN_ID {
			this.TestOut()
			tp := this.parseType()
			id := this.current.Lexeme
			this.advance()
			return ast.Id_new(id, tp, false, this.Linenum)
			//函数调用
		} else if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			m := ast.Id_new(id, &ast.Function{ast.TYPE_FUNCTION}, false, this.Linenum)
			return ast.CallExpr_new(m, args, this.Linenum)
		}
		log.Infof("适配变量ID->%s", id)
		return ast.Id_new(id, this.currentType, false, this.Linenum)
	case TOKEN_STRING:
		//FIXME 这里处理字符串字面值
		log.Infof("解析 TOKEN_STRING")
	case TOKEN_INT:
		id := this.current.Lexeme
		this.advance()
		//声明一个临时变量的语句
		if this.current.Kind == TOKEN_ID {
			log.Infof("parseAtomExp->TOKEN_INT")
			id = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			return ast.Id_new(id, &ast.Int{ast.TYPE_INT}, false, this.Linenum)

		}
		return ast.Id_new(id, this.currentType, false, this.Linenum)
		//表达式里出现Map --eg: Map.class
	case TOKEN_MAP:
		id := this.current.Lexeme
		this.eatToken(TOKEN_MAP)
		return ast.Id_new(id, this.currentType, false, this.Linenum)

	case TOKEN_NEW:
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
				//this.currentType = &ast.HashType{name, &ast.String{ast.TYPE_STRING}, &ast.ObjectType{ast.TYPE_OBJECT}, ast.TYPE_MAP}
			}
			this.eatToken(TOKEN_LT)
			var key ast.Type
			var ele ast.Type
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
			this.eatToken(TOKEN_RPAREN)
			return ast.NewHash_new(key, ele, this.Linenum)
		case TOKEN_ARRAYLIST:
			this.eatToken(TOKEN_ARRAYLIST)
			this.eatToken(TOKEN_LT)
			var ele ast.Type
			if this.current.Kind != TOKEN_GT {
				ele = this.parseType()
			} else {
				ele = this.assignType.(*ast.ListType).Ele
			}

			this.eatToken(TOKEN_GT)
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			return ast.NewList_new(ele, args, this.Linenum)

		case TOKEN_HASHSET:
			this.eatToken(TOKEN_HASHSET)
			this.eatToken(TOKEN_LT)
			var ele ast.Type
			if this.current.Kind != TOKEN_GT {
				ele = this.parseType()
			} else {
				ele = this.assignType.(*ast.ListType).Ele
			}

			this.eatToken(TOKEN_GT)
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			return ast.NewSet_new(ele, args, this.Linenum)
			//带参数对象初始化
		case TOKEN_ID:
			s := this.current.Lexeme
			this.advance()
			//模板
			if this.current.Kind == TOKEN_LT {
				this.eatToken(TOKEN_LT)
				this.eatToken(TOKEN_GT)
			}

			//数组
			if this.current.Kind == TOKEN_LBRACK {
				this.eatToken(TOKEN_LBRACK)
				size := this.parseExp()
				this.eatToken(TOKEN_RBRACK)
				if this.current.Kind == TOKEN_LBRACE {
					this.eatToken(TOKEN_LBRACE)
					eles := this.parseExp()
					this.eatToken(TOKEN_RBRACE)
					return ast.NewObjectArray_new(&ast.ClassType{s, ast.TYPE_CLASS}, eles, size, this.Linenum)
				}
				return ast.NewObjectArray_new(&ast.ClassType{s, ast.TYPE_CLASS}, nil, size, this.Linenum)
			}
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			return ast.NewObjectWithArgsList_new(&ast.ClassType{s, ast.TYPE_CLASS}, args, this.Linenum)
		default:
			log.Infof("********%v", this.current.Lexeme)
			panic("parser error1")
		}
	default:
		log.Infof("********%v", this.current.Lexeme)
		panic("parser error2")
	}
	return nil
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
	if this.current.Kind == TOKEN_LPAREN {
		args = append(args, this.parseLambdaExp())
	} else {
		args = append(args, this.parseExp())
	}

	for this.current.Kind == TOKEN_COMMER {
		this.advance()
		if this.current.Kind == TOKEN_LPAREN {
			args = append(args, this.parseLambdaExp())
		} else {
			args = append(args, this.parseExp())
		}
	}
	return args
}

func (this *Parser) parseLambdaExp() (exp ast.Exp) {
	log.Infof("解析 --> Lambda")
	this.eatToken(TOKEN_LPAREN)
	args := this.parseFormalList()
	this.eatToken(TOKEN_RPAREN)

	this.eatToken(TOKEN_LAMBDA)
	if this.current.Kind == TOKEN_LBRACE {
		this.eatToken(TOKEN_LBRACE)
		stms := this.parseStatements()
		this.eatToken(TOKEN_RBRACE)
		return ast.Lambda_new(args, stms, this.Linenum)
	} else {
		stm := this.parseStatement()
		return ast.Lambda_new(args, []ast.Stm{stm}, this.Linenum)
	}

	return
}

//NotExp    -> AtomExp
//          -> AtomExp.id(explist)
//          -> AtomExp[exp]
//          -> AtomExp.length
func (this *Parser) parseNotExp() ast.Exp {
	log.Infof("解析 parseNotExp")
	exp := this.parseAtomExp()
	for this.current.Kind == TOKEN_DOT ||
		this.current.Kind == TOKEN_AUTOSUB ||
		this.current.Kind == TOKEN_AUTOADD ||
		this.current.Kind == TOKEN_LBRACK {
		switch this.current.Kind {
		case TOKEN_AUTOSUB:
			this.eatToken(TOKEN_AUTOSUB)
			return ast.AutoSub_new(exp, nil, this.Linenum)
		case TOKEN_AUTOADD:
			this.eatToken(TOKEN_AUTOADD)
			return ast.AutoAdd_new(exp, nil, this.Linenum)
		//可以不断循环下去
		case TOKEN_DOT:
			log.Infof("解析函数调用,或成员变量")
			this.advance()
			if this.current.Kind == TOKEN_LENGTH {
				this.advance()
				return ast.Length_new(exp, this.Linenum)
			}
			if this.current.Kind == TOKEN_CLASS {
				this.advance()
				return ast.ClassExp_new(exp, this.Linenum)
			}
			c := this.current.Lexeme

			exp = ast.SelectorExpr_new(exp, this.current.Lexeme, this.Linenum)
			//点之后必须这个
			this.eatToken(TOKEN_ID)

			//成员函数
			if this.current.Kind == TOKEN_LPAREN {
				this.eatToken(TOKEN_LPAREN)
				log.Infof("解析函数调用参数-->%s", c)
				args := this.parseExpList()
				log.Infof("解析函数调用参数-->%s", c)
				this.eatToken(TOKEN_RPAREN)
				exp = ast.CallExpr_new(exp, args, this.Linenum)
				//成员变量
			}
			//数组索引操作
		case TOKEN_LBRACK: //[exp]
			this.advance()
			index := this.parseExp()
			if index != nil {
				log.Infof("数组索引表达式")
				this.eatToken(TOKEN_RBRACK)
				return ast.ArraySelect_new(exp, index, this.Linenum)
			} else {
				log.Infof("数组索引用")
				this.eatToken(TOKEN_RBRACK)

			}

		default:
			panic("need TOKEN_NOT or TOKEN_LBRACK")
		}
	}
	return exp
}

//TimesExp  -> !TimesExp
//          -> NotExp
func (this *Parser) parseTimeExp() ast.Exp {
	log.Infof("解析 parseTimeExp")
	var exp2 ast.Exp
	var opt = this.current.Kind
	for this.current.Kind == TOKEN_NOT ||
		this.current.Kind == TOKEN_AUTOADD ||
		this.current.Kind == TOKEN_AUTOSUB {
		this.advance()
		exp2 = this.parseTimeExp()
	}
	if exp2 != nil {
		switch opt {
		case TOKEN_NOT:
			return ast.Not_new(exp2, this.Linenum)
		case TOKEN_AUTOADD:
			return ast.AutoAdd_new(nil, exp2, this.Linenum)
		case TOKEN_AUTOSUB:
			return ast.AutoSub_new(nil, exp2, this.Linenum)
		default:
			panic("不支持")
		}

	} else {
		return this.parseNotExp()
	}
}

//AddSubExp -> TimesExp * TimesExp
//          -> TimesExp
func (this *Parser) parseAddSubExp() ast.Exp {
	log.Infof("解析 parseAddSubExp")
	left := this.parseTimeExp()
	for this.current.Kind == TOKEN_MUL {
		this.advance()
		right := this.parseTimeExp()
		return ast.Times_new(left, right, this.Linenum)
	}
	return left
}

//LtExp -> AddSubExp + AddSubExp
//      -> AddSubExp - AddSubExp
//      -> AddSubExp
func (this *Parser) parseLtExp() ast.Exp {
	log.Infof("解析parseLtExp")
	left := this.parseAddSubExp()
	for this.current.Kind == TOKEN_ADD ||
		this.current.Kind == TOKEN_SUB {
		switch this.current.Kind {
		case TOKEN_ADD:
			this.advance()
			right := this.parseAddSubExp()
			return ast.Add_new(left, right, this.Linenum)
		case TOKEN_SUB:
			this.advance()
			right := this.parseAddSubExp()
			return ast.Sub_new(left, right, this.Linenum)
		default:
			panic("need TOKEN_ADD or TOKEN_SUB")
		}
	}
	return left
}

//EqExp    -> EqExp == EqExp || EqExp != EqExp
//          -> EqExp
func (this *Parser) parseEqExp() ast.Exp {
	log.Infof("解析parseEqExp")
	left := this.parseLtExp()
	for this.current.Kind == TOKEN_LT || this.current.Kind == TOKEN_LE || this.current.Kind == TOKEN_GT || this.current.Kind == TOKEN_GE {
		opt := this.current.Kind
		this.advance()
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
	log.Infof("解析 parseAndExp")
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

//Exp -> AndExp && AndExp
//    -> AndExp
func (this *Parser) parseOrExp() ast.Exp {
	log.Infof("解析 parseOrExp")
	left := this.parseAndExp()
	for this.current.Kind == TOKEN_AND {
		this.advance()
		right := this.parseAndExp()
		left = ast.And_new(left, right, this.Linenum)
	}
	return left
}

//OrExp    -> OrExp || OrExp
//          -> OrExp
func (this *Parser) parseQuestionExp() ast.Exp {
	log.Infof("解析 parseQuestionExp")
	left := this.parseOrExp()
	for this.current.Kind == TOKEN_OR {
		log.Infof("TOKEN_OR")
		this.advance()
		right := this.parseOrExp()
		left = ast.Or_new(left, right, this.Linenum)
	}

	return left
}

//OrExp    -> OrExp || OrExp
//          -> OrExp
func (this *Parser) parseExp() ast.Exp {
	log.Infof("解析 parseExp")
	left := this.parseQuestionExp()
	for this.current.Kind == TOKEN_QUESTION {
		log.Infof("发现TOKEN_QUESTION")
		this.advance()
		one := this.parseQuestionExp()
		this.eatToken(TOKEN_COLON)
		two := this.parseQuestionExp()
		return ast.Question_new(left, one, two, this.Linenum)
	}

	return left
}

//
//
// return:
func (this *Parser) parseStatement() ast.Stm {
	log.Infof("*******解析代码段*******")
	switch this.current.Kind {
	case TOKEN_BOOLEAN:
		fallthrough
	case TOKEN_STRING:
		fallthrough
	case TOKEN_LONG:
		fallthrough
	case TOKEN_INT:
		fallthrough
	case TOKEN_SET:
		fallthrough
	case TOKEN_HASHSET:
		fallthrough
	case TOKEN_LIST:
		fallthrough
	case TOKEN_ARRAYLIST:
		fallthrough
	case TOKEN_MAP:
		fallthrough
	case TOKEN_HASHMAP:
		tp := this.parseType()
		id := this.current.Lexeme
		id = GetNewId(id)

		this.eatToken(TOKEN_ID)
		decl := ast.Decl_new(id, tp, nil, this.Linenum)
		//有赋值语句
		if this.current.Kind == TOKEN_ASSIGN {
			this.assignType = tp
			//临时变量类型
			log.Infof("*******解析临时变量声明语句(有赋值语句)*******")
			this.eatToken(TOKEN_ASSIGN)
			exp := this.parseExp()
			//三元表达式
			if _, ok := exp.(*ast.Question); ok {
				decl.SetTriple()
			}
			decl.Value = exp
		} else {
			log.Infof("*******解析临时变量声明语句(无赋值语句)*******")
		}
		this.eatToken(TOKEN_SEMI)
		return decl

	case TOKEN_LBRACE: //{
		log.Infof("*******解析代码段*******")
		this.eatToken(TOKEN_LBRACE)
		stms := this.parseStatements()
		this.eatToken(TOKEN_RBRACE)
		return ast.Block_new(stms, this.Linenum)
	case TOKEN_THIS:
		exp := this.parseExp()
		this.eatToken(TOKEN_SEMI)
		exprStm := ast.ExprStm_new(exp, this.Linenum)
		return exprStm
	case TOKEN_ID:
		id := this.current.Lexeme

		this.TestIn() //进入测试模式

		tp := this.parseType()
		switch this.current.Kind {
		//处理声明临时变量和赋值语句
		case TOKEN_ID:
			log.Infof("*******解析临时变量声明语句*******")
			id = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			decl := ast.Decl_new(id, tp, nil, this.Linenum)
			//有赋值语句
			if this.current.Kind == TOKEN_ASSIGN {
				//临时变量类型
				log.Infof("*******解析临时变量声明语句(有赋值语句)*******")
				this.eatToken(TOKEN_ASSIGN)
				exp := this.parseExp()
				//三元表达式
				if _, ok := exp.(*ast.Question); ok {
					decl.SetTriple()
				}
				decl.Value = exp
			}
			this.eatToken(TOKEN_SEMI)
			return decl
			//都统一为赋值语句
		case TOKEN_LPAREN:
			fallthrough
		case TOKEN_DOT:
			log.Infof("*******解析函数调用*******")
			this.TestOut() //恢复测试数据
			exp := this.parseExp()
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
		case TOKEN_ASSIGN:
			this.eatToken(TOKEN_ASSIGN)
			exp := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			assign := new(ast.Assign)
			assign.Left = ast.Id_new(id, nil, false, this.Linenum)
			assign.Name = id
			assign.E = exp
			//三元表达式
			if _, ok := exp.(*ast.Question); ok {
				assign.SetTriple()
			}
			return assign
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
			this.eatToken(TOKEN_LT)
		default:
			log.Infof("parseStatement-->%v", this.current.Lexeme)
			panic("bug1")

		}
	case TOKEN_IF:
		log.Infof("********TOKEN_IF***********")
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
		log.Infof("********TOKEN_TRY***********")
		this.eatToken(TOKEN_TRY)
		test := this.parseStatement()
		var conditions []ast.Exp
		var catches []ast.Stm
		var finally ast.Stm
		for this.current.Kind == TOKEN_CATCH {
			this.eatToken(TOKEN_CATCH)
			this.eatToken(TOKEN_LPAREN)
			condition := this.parseExp()
			conditions = append(conditions, condition)
			this.eatToken(TOKEN_RPAREN)
			catch := this.parseStatement()
			catches = append(catches, catch)
		}

		if this.current.Kind == TOKEN_FINALLY {
			this.eatToken(TOKEN_FINALLY)
			finally = this.parseStatement()
		}
		return ast.Try_new(test, conditions, catches, finally, this.Linenum)
	case TOKEN_WHILE:
		log.Infof("********TOKEN_WHILE***********")
		this.eatToken(TOKEN_WHILE)
		this.eatToken(TOKEN_LPAREN)
		exp := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		body := this.parseStatement()
		return ast.While_new(exp, body, this.Linenum)
	case TOKEN_FOR:
		log.Infof("********TOKEN_FOR***********")
		this.eatToken(TOKEN_FOR)
		this.eatToken(TOKEN_LPAREN)
		this.TestIn()

		exp := this.parseExp()
		//说明是声明语句
		if this.current.Kind == TOKEN_ID {
			this.TestOut()
			this.parseType()
			exp = this.parseExp()
		}

		//for循环三段式
		if this.current.Kind == TOKEN_ASSIGN {
			log.Infof("********TOKEN_FOR--> 解析初始化语句 ***********")
			Init := new(ast.Assign)
			Init.Left = exp
			//临时变量类型
			this.eatToken(TOKEN_ASSIGN)
			exp1 := this.parseExp()
			Init.E = exp1
			this.eatToken(TOKEN_SEMI)
			//
			log.Infof("********TOKEN_FOR--> 解析条件语句 ***********")
			Condition := this.parseExp()
			this.eatToken(TOKEN_SEMI)

			log.Infof("********TOKEN_FOR--> 解析更新语句 ***********")
			Post := this.parseExp()
			this.eatToken(TOKEN_RPAREN)
			body := this.parseStatement()
			return ast.For_new(Init, Condition, Post, body, this.Linenum)

			//枚举式
		} else if this.current.Kind == TOKEN_COLON {
			log.Infof("*******for循环枚举*************")
			this.eatToken(TOKEN_COLON)
			var right ast.Exp

			//处理强制类型转换
			if this.current.Kind == TOKEN_LPAREN {
				right = this.parseCastExp()
			} else {
				right = this.parseOrExp()
			}
			this.eatToken(TOKEN_RPAREN)

			body := this.parseStatement()

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
		if this.current.Kind == TOKEN_SEMI {
			this.eatToken(TOKEN_SEMI)
			return ast.Return_new(nil, this.Linenum)
		}
		log.Infof("<<<>>>解析return,%v", this.current.Lexeme)
		e := this.parseExp()
		this.eatToken(TOKEN_SEMI)
		return ast.Return_new(e, this.Linenum)
	default:
		log.Infof("token error->%s", this.current.Lexeme)
		panic("token error")
	}
	return nil
}

func (this *Parser) parseStatements() []ast.Stm {
	stms := []ast.Stm{}
	for this.current.Kind == TOKEN_LBRACE ||
		this.current.Kind == TOKEN_ID ||
		this.current.Kind == TOKEN_LIST ||
		this.current.Kind == TOKEN_ARRAYLIST ||
		this.current.Kind == TOKEN_MAP ||
		this.current.Kind == TOKEN_HASHMAP ||
		this.current.Kind == TOKEN_IF ||
		this.current.Kind == TOKEN_TRY ||
		this.current.Kind == TOKEN_WHILE ||
		this.current.Kind == TOKEN_FOR ||
		this.current.Kind == TOKEN_THROW ||
		this.current.Kind == TOKEN_RETURN ||
		this.current.Kind == TOKEN_BOOLEAN ||
		this.current.Kind == TOKEN_INT ||
		this.current.Kind == TOKEN_LONG ||
		this.current.Kind == TOKEN_STRING ||
		this.current.Kind == TOKEN_SET ||
		this.current.Kind == TOKEN_HASHSET ||
		this.current.Kind == TOKEN_THIS ||
		this.current.Kind == TOKEN_SYSTEM {
		log.Infof(">>>>>>>>>>>>>>>>>>>>>>>>>解析代码段:%s", this.current.Lexeme)
		stms = append(stms, this.parseStatement())
	}
	return stms
}
func (this *Parser) parseMemberVarDecl(tmp *ast.FieldSingle) ast.Field {
	var dec *ast.FieldSingle
	var assign *ast.Assign

	if this.current.Kind == TOKEN_ASSIGN {
		this.eatToken(TOKEN_ASSIGN)
		e := this.parseExp()
		this.isSpecial = false
		assign := new(ast.Assign)
		assign.Name = tmp.Name
		assign.E = e
	}
	dec = &ast.FieldSingle{tmp.Access, tmp.Tp, tmp.Name, true, assign}
	this.eatToken(TOKEN_SEMI)
	return dec
}

// 解析成员变量和成员方法
//
// return:
func (this *Parser) parseClassContext() (decs []ast.Field, methods []ast.Method) {

	//每次循环解析一个成员变量或一个成员函数
	for this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PROTECTED ||
		this.current.Kind == TOKEN_BOOLEAN || this.current.Kind == TOKEN_INT || this.current.Kind == TOKEN_STRING ||
		this.current.Kind == TOKEN_ID {
		//
		var tmp ast.FieldSingle

		//访问修饰符 [其他修饰符] 类型 变量名 = 值;
		//处理 访问修饰符
		if this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PROTECTED {
			log.Info("处理访问修饰符:", this.current.ToString())
			//1 扫描访问修饰符
			tmp.Access = this.current.Kind
			this.advance()
		} else {
			tmp.Access = TOKEN_DEFAULT
		}

		//处理 其他修饰符(忽略)
		if this.current.Kind == TOKEN_STATIC {
			this.eatToken(TOKEN_STATIC)
		}

		if this.current.Kind == TOKEN_FINAL {
			this.eatToken(TOKEN_FINAL)
		}

		if this.current.Kind == TOKEN_TRANSIENT {
			this.eatToken(TOKEN_TRANSIENT)
		}

		//类型
		tmp.Tp = this.parseType()

		//变量/函数名
		tmp.Name = this.current.Lexeme
		this.eatToken(TOKEN_ID)

		//成员方法
		if this.current.Kind == TOKEN_LPAREN {
			methods = append(methods, this.parseMemberMethod(&tmp))
			//成员变量
		} else {
			decs = append(decs, this.parseMemberVarDecl(&tmp))
		}

	}
	return
}

func (this *Parser) parseMemberMethod(dec *ast.FieldSingle) ast.Method {
	log.Infof("*******解析成员函数*******")
	//左括号
	this.eatToken(TOKEN_LPAREN)
	//解析参数
	formals := this.parseFormalList()
	//右括号
	this.eatToken(TOKEN_RPAREN)

	if this.current.Kind == TOKEN_THROWS {
		this.eatToken(TOKEN_THROWS)
		this.eatToken(TOKEN_ID)
	}
	//做大括号
	this.eatToken(TOKEN_LBRACE)
	var stms []ast.Stm
	var locals []ast.Field

	//解析本地变量和表达式
	stms = this.parseStatements()
	var retExp ast.Exp
	if this.current.Kind == TOKEN_RETURN {
		this.eatToken(TOKEN_RETURN)
		retExp = this.parseExp()
		this.eatToken(TOKEN_SEMI)
	}

	this.eatToken(TOKEN_RBRACE)

	return &ast.MethodSingle{dec.Tp, dec.Name, formals, locals, stms, retExp}
}

//解析类
func (this *Parser) parseClassDecl() ast.Class {
	var id, extends string

	//类访问权限修饰符
	var access int
	if this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PROTECTED {
		access = this.current.Kind
		this.advance()
	}
	//处理abstract
	if this.current.Kind == TOKEN_ABSTRACT {
		this.advance()
	}

	this.eatToken(TOKEN_CLASS)
	id = this.current.Lexeme
	this.eatToken(TOKEN_ID)

	//处理extends
	if this.current.Kind == TOKEN_EXTENDS {
		this.eatToken(TOKEN_EXTENDS)
		extends = this.current.Lexeme
		this.eatToken(TOKEN_ID)
	}

	//处理implements
	if this.current.Kind == TOKEN_IMPLEMENTS {
		this.eatToken(TOKEN_IMPLEMENTS)
		extends = this.current.Lexeme
		this.eatToken(TOKEN_ID)
	}

	this.eatToken(TOKEN_LBRACE)

	//处理成员变量
	//decs := this.parseVarDecls()
	decs, methods := this.parseClassContext()
	//处理方法
	//methods := this.parseMethodDecls()

	this.eatToken(TOKEN_RBRACE)
	return &ast.ClassSingle{access, id, extends, decs, methods}
}

// 解析类组
//
// return:
func (this *Parser) parseClassDecls() []ast.Class {
	classes := []ast.Class{}
	for this.current.Kind == TOKEN_CLASS || this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PROTECTED {
		classes = append(classes, this.parseClassDecl())
	}
	return classes
}
func (this *Parser) parseAnnotation() {
	this.eatToken(TOKEN_AT)

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

	//处理import
	for this.current.Kind == TOKEN_IMPORT {
		this.advance()
		for this.current.Kind != TOKEN_SEMI {
			this.advance()
		}
		this.advance()
	}

	classes := this.parseClassDecls()
	this.eatToken(TOKEN_EOF)
	return &ast.FileSingle{name, nil, classes}
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
