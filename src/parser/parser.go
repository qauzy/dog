package parser

import (
	"dog/ast"
	"dog/control"
	"dog/util"
	"fmt"
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
	Linenum     int
}

func NewParse(fname string, buf []byte) *Parser {
	lexer := NewLexer(fname, buf)
	p := new(Parser)
	p.lexer = lexer
	p.current = p.lexer.NextToken()

	return p
}

func (this *Parser) getFP() int {
	return this.lexer.fp
}

func (this *Parser) resetFP(fp int) {
	this.lexer.fp = fp
}

func (this *Parser) advance() {
	if control.Lexer_dump == true {
		fmt.Println(this.current.ToString())
	}
	this.Linenum = this.current.LineNum
	this.current = this.lexer.NextToken()
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
			this.currentType = &ast.StringArray{ast.TOKEN_STRING}
		} else {
			this.currentType = &ast.String{}
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
		this.eatToken(TOKEN_LT)
		key := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		this.eatToken(TOKEN_COMMER)
		ele := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.HashType{name, key, ele, ast.TYPE_MAP}
	case TOKEN_HASHMAP:
		name := this.current.Lexeme
		this.eatToken(TOKEN_HASHMAP)
		this.eatToken(TOKEN_LT)
		key := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		this.eatToken(TOKEN_COMMER)
		ele := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		this.eatToken(TOKEN_GT)
		this.currentType = &ast.HashType{name, key, ele, ast.TYPE_MAP}

	default:
		name := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		this.currentType = &ast.ClassType{name, ast.TYPE_CLASS}
	}
	log.Infof("解析类型:%s", this.currentType.String())
	return this.currentType
}

func (this *Parser) parseFormalList() []ast.Dec {
	flist := []ast.Dec{}
	var tp ast.Type
	var id string
	var access int

	if this.current.Kind == TOKEN_ID ||
		this.current.Kind == TOKEN_INT ||
		this.current.Kind == TOKEN_LIST ||
		this.current.Kind == TOKEN_MAP ||
		this.current.Kind == TOKEN_BOOLEAN {
		tp = this.parseType()
		id = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		flist = append(flist, &ast.DecSingle{access, tp, id, this.isField, nil})

		for this.current.Kind == TOKEN_COMMER {
			this.eatToken(TOKEN_COMMER)
			tp = this.parseType()
			id = this.current.Lexeme
			this.eatToken(TOKEN_ID)
			flist = append(flist, &ast.DecSingle{access, tp, id, this.isField, nil})
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
	case TOKEN_THIS:
		this.advance()
		return &ast.This{}
	case TOKEN_SYSTEM:
		var m ast.Exp
		m = ast.Id_new(this.current.Lexeme, nil, false, this.Linenum)
		this.eatToken(TOKEN_SYSTEM)
		for this.current.Kind == TOKEN_DOT {
			this.eatToken(TOKEN_DOT)
			exp := ast.Id_new(this.current.Lexeme, nil, false, this.Linenum)
			this.eatToken(TOKEN_ID)
			if this.current.Kind == TOKEN_LPAREN {
				this.eatToken(TOKEN_LPAREN)
				args := this.parseExpList()
				this.eatToken(TOKEN_RPAREN)
				m = ast.Dot_new(m, exp, args, "", nil, nil, this.Linenum)
			} else {
				m = ast.Dot_new(m, exp, nil, "", nil, nil, this.Linenum)
			}
		}
		return m
	case TOKEN_ID:
		id := this.current.Lexeme
		tp := this.parseType()
		//声明一个临时变量的语句
		if this.current.Kind == TOKEN_ID {
			id := this.current.Lexeme
			this.advance()
			return ast.Id_new(id, tp, false, this.Linenum)
			//函数调用
		} else if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			m := ast.Id_new(id, tp, false, this.Linenum)
			return ast.Dot_new(nil, m, args, "", nil, nil, this.Linenum)
		}
		return ast.Id_new(id, this.currentType, false, this.Linenum)
	case TOKEN_STRING:
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
	case TOKEN_NEW:
		this.advance()
		switch this.current.Kind {
		case TOKEN_INT:
			this.advance()
			this.eatToken(TOKEN_LBRACK)
			exp := this.parseExp()
			this.eatToken(TOKEN_RBRACK)
			return ast.NewIntArray_new(exp, this.Linenum)
		case TOKEN_HASHMAP:
			this.eatToken(TOKEN_HASHMAP)
			this.eatToken(TOKEN_LT)
			var key = ""
			var ele = ""
			if this.current.Kind == TOKEN_ID {
				key = this.current.Lexeme
				this.eatToken(TOKEN_ID)
				this.eatToken(TOKEN_COMMER)
				ele = this.current.Lexeme
				this.eatToken(TOKEN_ID)
			} else {
				key = this.currentType.(*ast.HashType).Key
				ele = this.currentType.(*ast.HashType).Ele
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
				this.eatToken(TOKEN_RBRACK)
				this.eatToken(TOKEN_LBRACE)
				exp := this.parseExp()
				this.eatToken(TOKEN_RBRACE)
				return ast.NewObjectArray_new(exp, this.Linenum)
			}
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			this.eatToken(TOKEN_RPAREN)
			return ast.NewObjectWithArgsList_new(s, args, this.Linenum)
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

func (this *Parser) parseExpList() []ast.Exp {
	args := []ast.Exp{}
	if this.current.Kind == TOKEN_RPAREN {
		return args
	}

	args = append(args, this.parseExp())
	for this.current.Kind == TOKEN_COMMER {
		this.advance()
		args = append(args, this.parseExp())
	}
	return args
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
			//else ast.Call
			var right ast.Exp
			right = ast.Id_new(this.current.Lexeme, nil, false, this.Linenum)
			//点之后必须这个
			this.eatToken(TOKEN_ID)
			//成员函数
			if this.current.Kind == TOKEN_LPAREN {
				this.eatToken(TOKEN_LPAREN)
				args := this.parseExpList()
				log.Infof(this.current.Lexeme)
				this.eatToken(TOKEN_RPAREN)
				exp = ast.Dot_new(exp, right, args, "", nil, nil, this.Linenum)

				//成员变量
			} else {
				exp = ast.Dot_new(exp, right, nil, "", nil, nil, this.Linenum)
			}
		case TOKEN_LBRACK: //[exp]
			this.advance()
			index := this.parseExp()
			this.eatToken(TOKEN_RBRACK)
			return ast.ArraySelect_new(exp, index, this.Linenum)
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
	for this.current.Kind == TOKEN_TIMES {
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
func (this *Parser) parseExp() ast.Exp {
	left := this.parseOrExp()
	for this.current.Kind == TOKEN_OR {
		log.Infof("发现TOKEN_OR")
		this.advance()
		right := this.parseOrExp()
		left = ast.Or_new(left, right, this.Linenum)
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
		this.eatToken(TOKEN_ID)
		assign := new(ast.Assign)
		assign.Left = ast.Id_new(id, tp, false, this.Linenum)
		assign.Name = id
		//有赋值语句
		if this.current.Kind == TOKEN_ASSIGN {
			//临时变量类型
			this.assignType = tp
			log.Infof("*******解析临时变量声明语句(有赋值语句)*******")
			this.eatToken(TOKEN_ASSIGN)
			exp := this.parseExp()
			assign.E = exp
		} else {
			log.Infof("*******解析临时变量声明语句(无赋值语句)*******")
		}
		this.eatToken(TOKEN_SEMI)
		return assign

	case TOKEN_LBRACE: //{
		log.Infof("*******解析代码段*******")
		this.eatToken(TOKEN_LBRACE)
		stms := this.parseStatements()
		this.eatToken(TOKEN_RBRACE)
		return ast.Block_new(stms, this.Linenum)
	case TOKEN_THIS:
		exp := this.parseExp()
		this.eatToken(TOKEN_SEMI)
		assign := new(ast.Assign)
		//assign.Left = ast.Id_new(id, tp, false, this.Linenum)--->直接点调用,没有赋值语句,只是作为一个承载
		//assign.Name = id
		assign.E = exp
		return assign
	case TOKEN_ID:
		id := this.current.Lexeme

		fp := this.getFP()
		cur := this.current
		tp := this.parseType()
		switch this.current.Kind {
		//处理声明临时变量和赋值语句
		case TOKEN_ID:
			log.Infof("*******解析临时变量声明语句*******")
			id := this.current.Lexeme
			this.eatToken(TOKEN_ID)
			if this.current.Kind == TOKEN_ASSIGN {
				this.eatToken(TOKEN_ASSIGN)
				exp := this.parseExp()
				this.eatToken(TOKEN_SEMI)
				assign := new(ast.Assign)
				assign.Left = ast.Id_new(id, tp, false, this.Linenum)
				assign.Name = id
				assign.E = exp
				return assign
			}
			this.eatToken(TOKEN_SEMI)
			//都统一为赋值语句
		case TOKEN_LPAREN:
			fallthrough
		case TOKEN_DOT:
			log.Infof("*******解析函数调用*******")
			this.resetFP(fp)
			this.current = cur
			exp := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			assign := new(ast.Assign)
			//assign.Left = ast.Id_new(id, tp, false, this.Linenum)--->直接点调用,没有赋值语句,只是作为一个承载
			//assign.Name = id
			assign.E = exp
			return assign
		case TOKEN_ASSIGN:
			this.eatToken(TOKEN_ASSIGN)
			exp := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			assign := new(ast.Assign)
			assign.Name = id
			assign.E = exp
			return assign
		case TOKEN_LBRACK:
			this.eatToken(TOKEN_LBRACK) //[
			index := this.parseExp()
			this.eatToken(TOKEN_RBRACK) //]
			this.eatToken(TOKEN_ASSIGN)
			exp := this.parseExp()
			this.eatToken(TOKEN_SEMI)
			return ast.AssignArray_new(id, index, exp, nil, false, this.Linenum)
		default:
			log.Infof("parseStatement:%v", this.current.Lexeme)
			panic("bug1")

		}
	case TOKEN_IF:
		log.Infof("********TOKEN_IF***********")
		this.eatToken(TOKEN_IF)
		this.eatToken(TOKEN_LPAREN)
		condition := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		thenn := this.parseStatement()
		if this.current.Kind == TOKEN_ELSE {
			this.eatToken(TOKEN_ELSE)
			elsee := this.parseStatement()
			return ast.If_new(condition, thenn, elsee, this.Linenum)
		} else {
			return ast.If_new(condition, thenn, nil, this.Linenum)
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
		//init := this.parseStatement()
		var Init ast.Stm
		exp := this.parseExp()
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
			Update := this.parseExp()
			log.Infof("*******1111**********")
			//可能是for循环枚举
			exp = ast.Fcon_new(Init, Condition, Update, this.Linenum)

		} else if this.current.Kind == TOKEN_COLON {
			this.eatToken(TOKEN_COLON)
			var right ast.Exp

			//处理强制类型转换
			if this.current.Kind == TOKEN_LPAREN {
				right = this.parseCastExp()
			} else {
				right = this.parseOrExp()
			}
			log.Infof("*******for循环枚举*************")
			exp = ast.Enum_new(exp, right, this.Linenum)
		}

		this.eatToken(TOKEN_RPAREN)
		log.Infof("解析for语句body")
		body := this.parseStatement()
		return ast.For_new(Init, exp, body, this.Linenum)
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
func (this *Parser) parseMemberVarDecl(tmp *ast.DecSingle) ast.Dec {
	var dec *ast.DecSingle
	var assign *ast.Assign

	if this.current.Kind == TOKEN_ASSIGN {
		this.eatToken(TOKEN_ASSIGN)
		e := this.parseExp()
		this.isSpecial = false
		assign := new(ast.Assign)
		assign.Name = tmp.Name
		assign.E = e
	}
	dec = &ast.DecSingle{tmp.Access, tmp.Tp, tmp.Name, this.isField, assign}
	this.eatToken(TOKEN_SEMI)
	return dec
}

// 解析成员变量和成员方法
//
// return:
func (this *Parser) parseClassContext() (decs []ast.Dec, methods []ast.Method) {

	//每次循环解析一个成员变量或一个成员函数
	for this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PROTECTED ||
		this.current.Kind == TOKEN_BOOLEAN || this.current.Kind == TOKEN_INT || this.current.Kind == TOKEN_STRING ||
		this.current.Kind == TOKEN_ID {
		//
		var tmp ast.DecSingle

		//访问修饰符 [其他修饰符] 类型 变量名 = 值;
		//处理 访问修饰符
		if this.current.Kind == TOKEN_PUBLIC || this.current.Kind == TOKEN_PRIVATE || this.current.Kind == TOKEN_PROTECTED {
			fmt.Println("处理访问修饰符:", this.current.ToString())
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

func (this *Parser) parseMemberMethod(dec *ast.DecSingle) ast.Method {
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
	var locals []ast.Dec

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

func (this *Parser) parseProgram() ast.Program {
	//处理package
	if this.current.Kind == TOKEN_PACKAGE {
		this.advance()
		for this.current.Kind != TOKEN_SEMI {
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

	////解析主入口类
	//main_class := this.parseMainClass()
	////解析类描述

	classes := this.parseClassDecls()
	this.eatToken(TOKEN_EOF)
	return &ast.ProgramSingle{nil, classes}
}

func (this *Parser) Parser() ast.Program {
	p := this.parseProgram()
	return p
}
