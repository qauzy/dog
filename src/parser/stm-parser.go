package parser

import (
	"dog/ast"
	log "github.com/corgi-kx/logcustom"
)

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
		this.advance()
		x := ast.NewIdent(id, this.Linenum)
		exp := this.parseCallExp(x)
		switch this.current.Kind {
		//处理声明临时变量和赋值语句
		case TOKEN_ID:
			return this.parserDecl(exp)
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
			//三元表达式
			if q, ok := exp.(*ast.Question); ok {
				assign1 := ast.Assign_new(ast.Id_new(id, nil, false, this.Linenum), q.One, false, this.Linenum)
				assign2 := ast.Assign_new(ast.Id_new(id, nil, false, this.Linenum), q.Two, false, this.Linenum)
				return ast.If_new(q.E, ast.Block_new([]ast.Stm{assign1}, this.Linenum), ast.Block_new([]ast.Stm{assign2}, this.Linenum), this.Linenum)
			}
			assign.Left = ast.Id_new(id, nil, false, this.Linenum)
			assign.Value = exp

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
			//数组类型
			if this.current.Kind == TOKEN_RBRACK {
				this.eatToken(TOKEN_RBRACK) //]
				return this.parserDecl(&ast.ArrayType{Ele: exp})
			}
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
		this.eatToken(TOKEN_ID)
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
		exp := this.parseExp()
		this.eatToken(TOKEN_SEMI)
		//三元表达式
		if q, ok := exp.(*ast.Question); ok {
			assign1 := ast.Return_new(q.One, this.Linenum)
			assign2 := ast.Return_new(q.Two, this.Linenum)
			return ast.If_new(q.E, assign1, assign2, this.Linenum)
		}
		return ast.Return_new(exp, this.Linenum)
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
