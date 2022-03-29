package parser

import (
	"dog/ast"
	"dog/cfg"
	log "github.com/corgi-kx/logcustom"
)

//
//
// return:
func (this *Parser) parseStatement() ast.Stm {
	log.Debugf("*******解析代码段******* --> %v", this.current.Lexeme)
	defer func() {
		if this.current.Kind == TOKEN_SEMI {
			this.eatToken(TOKEN_SEMI)
		}
	}()
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
	case TOKEN_BREAK:
		this.advance()
		if this.current.Kind != TOKEN_SEMI {
			var id = ast.NewIdent(this.current.Lexeme, this.Linenum)
			this.eatToken(TOKEN_ID)
			return ast.BranchStmt_new(id, TOKEN_BREAK, this.Linenum)
		}
		return ast.ExprStm_new(ast.NewIdent("break", this.Linenum), this.Linenum)
	case TOKEN_GOTO:
		this.advance()
		var id = ast.NewIdent(this.current.Lexeme, this.Linenum)
		this.eatToken(TOKEN_ID)
		return ast.BranchStmt_new(id, TOKEN_GOTO, this.Linenum)
	case TOKEN_COMMENT:
		stm := ast.Comment_new(this.current.Lexeme, this.Linenum)
		this.advance()
		return stm
	case TOKEN_SUPER:
		id := ast.NewIdent(this.current.Lexeme, this.Linenum)
		this.eatToken(TOKEN_SUPER)
		exp := this.parseCallExp(id)

		return ast.ExprStm_new(exp, this.Linenum)

	case TOKEN_SYNCHRONIZED:
		this.eatToken(TOKEN_SYNCHRONIZED)
		this.eatToken(TOKEN_LPAREN)
		exp := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		body := this.parseStatement()
		return ast.Sync_new(exp, body, this.Linenum)

	case TOKEN_NEW:
		exp := this.parseNewExp()
		if this.current.Kind == TOKEN_DOT {
			exp = this.parseCallExp(exp)
		}

		return ast.ExprStm_new(exp, this.Linenum)
	case TOKEN_LBRACE: //{
		log.Debugf("*******解析代码段*******")
		this.eatToken(TOKEN_LBRACE)
		stms := this.parseStatements()
		this.eatToken(TOKEN_RBRACE)
		return ast.Block_new(stms, this.Linenum)
	case TOKEN_THIS:
		//1 调用构造函数
		//2 调用成员函数
		//3 调用成员变量

		exp := this.parseExp()

		if this.current.Kind == TOKEN_ASSIGN {
			this.eatToken(TOKEN_ASSIGN)
			if vv, ok := exp.(*ast.SelectorExpr); ok {
				if this.currentClass.GetField(vv.Sel) != nil {
					f := this.currentClass.GetField(vv.Sel)
					this.assignType = f.GetDecType()
				}
			}
			right := this.parseExp()

			assign := ast.Assign_new(exp, right, "=", false, this.Linenum)

			//三元表达式
			if _, ok := right.(*ast.Question); ok {
				assign.SetTriple()
			}
			return assign

		}

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
		// goto statement
		case TOKEN_COLON:
			this.advance()
			var stmt = this.parseStatement()
			return ast.LabeledStmt_new(x, stmt, this.Linenum)
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

				assign := ast.Assign_new(exp, right, "=", false, this.Linenum)
				//三元表达式
				if _, ok := right.(*ast.Question); ok {
					assign.SetTriple()
				}
				return assign
			} else {

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
			if this.Peek().GetField(id) != nil {
				this.assignType = this.Peek().GetField(id).GetDecType()
			}
			right := this.parseExp()

			//三元表达式
			if q, ok := exp.(*ast.Question); ok {
				assign1 := ast.Assign_new(ast.NewIdent(id, this.Linenum), q.One, "=", false, this.Linenum)
				assign2 := ast.Assign_new(ast.NewIdent(id, this.Linenum), q.Two, "=", false, this.Linenum)
				return ast.If_new(q.E, ast.Block_new([]ast.Stm{assign1}, this.Linenum), ast.Block_new([]ast.Stm{assign2}, this.Linenum), this.Linenum)
			}
			var left = ast.NewIdent(id, this.Linenum)
			var call string
			var mp = new(ast.StreamStm)
			if this.CheckStreamExprs(right, &call, mp) {
				mp.Left = left
				mp.LineNum = this.Linenum
				return mp
			}
			assign := ast.Assign_new(exp, right, "=", false, this.Linenum)
			return assign

		case TOKEN_QUO_ASSIGN:
			left := ast.NewIdent(id, this.Linenum)
			this.eatToken(TOKEN_QUO_ASSIGN)
			right := this.parseExp()

			return ast.Assign_new(left, right, "/=", false, this.Linenum)
		case TOKEN_MUL_ASSIGN:

			left := ast.NewIdent(id, this.Linenum)
			this.eatToken(TOKEN_MUL_ASSIGN)
			right := this.parseExp()

			return ast.Assign_new(left, right, "*=", false, this.Linenum)
		case TOKEN_SUB_ASSIGN:

			left := ast.NewIdent(id, this.Linenum)
			this.eatToken(TOKEN_SUB_ASSIGN)
			right := this.parseExp()

			return ast.Assign_new(left, right, "-=", false, this.Linenum)
		case TOKEN_ADD_ASSIGN:
			left := ast.NewIdent(id, this.Linenum)
			this.eatToken(TOKEN_ADD_ASSIGN)
			right := this.parseExp()

			return ast.Assign_new(left, right, "+=", false, this.Linenum)
		case TOKEN_REM_ASSIGN:
			left := ast.NewIdent(id, this.Linenum)
			this.eatToken(TOKEN_REM_ASSIGN)
			right := this.parseExp()

			return ast.Assign_new(left, right, "%=", false, this.Linenum)

			//处理的是前缀加
		case TOKEN_INCREMENT:
			log.Debugf("处理累加")
			this.eatToken(TOKEN_INCREMENT)
			//特殊的for语句才不需要分号
			if this.isSpecial {
				this.isSpecial = false
			}
			left := ast.NewIdent(id, this.Linenum)

			return ast.Assign_new(left, &ast.Num{Value: 1}, "+=", false, this.Linenum)
			//处理的是后缀减
		case TOKEN_DECREMENT:
			this.eatToken(TOKEN_DECREMENT)

			left := ast.NewIdent(id, this.Linenum)
			return ast.Assign_new(left, &ast.Num{Value: 1}, "-=", false, this.Linenum)
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

			return ast.AssignArray_new(id, index, exp, nil, false, this.Linenum)
		case TOKEN_LT:
			tp := &ast.ClassType{id, ast.TYPE_CLASS}
			this.eatToken(TOKEN_LT)
			this.parseType()
			for this.current.Kind == TOKEN_COMMER {
				this.eatToken(TOKEN_COMMER)
				this.parseType()
			}
			this.eatToken(TOKEN_GT)
			this.currentType = &ast.ClassType{id, ast.TYPE_CLASS}
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

			return decl

		case TOKEN_SEMI:
			//因为List && Map Set 转换中间产物
			if fk, ok := exp.(*ast.FakeExpr); ok {
				return fk.Stm
			}
			return ast.ExprStm_new(exp, this.Linenum)
		case TOKEN_COMMENT:
			this.advance()
		default:
			this.ParseBug("parseStatement代码段解析bug")
		}
	case TOKEN_IF:
		log.Debugf("********TOKEN_IF***********")
		var fake = ast.FakeStm_new(this.Peek(), this.Linenum)
		this.currentStm = fake
		this.stmStack.Push(fake)
		this.Push(fake)
		defer func() {
			this.Pop()
			this.stmStack.Pop()
			if this.stmStack.Peek() != nil {
				this.currentStm = this.stmStack.Peek().(ast.Stm)
			} else {
				this.currentStm = nil
			}
		}()

		this.eatToken(TOKEN_IF)
		this.eatToken(TOKEN_LPAREN)
		condition := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		//if 条件后面不允许注释
		if this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}
		var Init ast.Exp
		if cfg.NoGeneric {
			if cl, ok := condition.(*ast.CallExpr); ok && len(cl.ArgsList) == 1 {
				if sl, ok := cl.Callee.(*ast.SelectorExpr); ok {
					if sl.Sel == "containsKey" {
						if ident, ok := sl.X.(*ast.Ident); ok {
							if this.CheckField(ident.Name) != nil {
								if _, ok := this.CheckField(ident.Name).GetDecType().(*ast.MapType); ok {
									Init = ast.IndexExpr_new(ident, cl.ArgsList[0], this.Linenum)
									condition = ast.NewIdent("ok", this.Linenum)
								}

							}
						}
					}

				}

			}
		}

		body := this.parseStatement()

		//不是block,说明没有大括号
		if _, ok := body.(*ast.Block); !ok {
			body = ast.Block_new([]ast.Stm{body}, this.Linenum)
		}
		if this.current.Kind == TOKEN_ELSE {
			this.eatToken(TOKEN_ELSE)
			elsee := this.parseStatement()
			if _, ok := elsee.(*ast.Block); !ok {
				if _, ok := elsee.(*ast.If); !ok {
					elsee = ast.Block_new([]ast.Stm{elsee}, this.Linenum)
				}
			}
			if Init != nil {
				return ast.If_newEx(Init, condition, body, elsee, this.Linenum)
			}

			return ast.If_new(condition, body, elsee, this.Linenum)
		} else {
			if Init != nil {
				return ast.If_newEx(Init, condition, body, nil, this.Linenum)
			}
			return ast.If_new(condition, body, nil, this.Linenum)
		}
	case TOKEN_TRY:
		log.Debugf("********TOKEN_TRY***********")
		var fake = ast.FakeStm_new(this.Peek(), this.Linenum)
		this.currentStm = fake
		this.stmStack.Push(fake)
		this.Push(fake)
		defer func() {
			this.Pop()
			this.stmStack.Pop()
			if this.stmStack.Peek() != nil {
				this.currentStm = this.stmStack.Peek().(ast.Stm)
			} else {
				this.currentStm = nil
			}
		}()
		this.eatToken(TOKEN_TRY)
		var resource ast.Stm
		if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			//资源初始化语句
			this.isSpecial = true
			resource = this.parseStatement()
			this.isSpecial = false
			this.eatToken(TOKEN_RPAREN)
		}
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
		return ast.Try_new(resource, body, catches, finally, this.Linenum)
	case TOKEN_WHILE:
		log.Debugf("********TOKEN_WHILE***********")
		var fake = ast.FakeStm_new(this.Peek(), this.Linenum)
		this.currentStm = fake
		this.stmStack.Push(fake)
		this.Push(fake)
		defer func() {
			this.Pop()
			this.stmStack.Pop()
			if this.stmStack.Peek() != nil {
				this.currentStm = this.stmStack.Peek().(ast.Stm)
			} else {
				this.currentStm = nil
			}
		}()
		this.eatToken(TOKEN_WHILE)
		this.eatToken(TOKEN_LPAREN)
		exp := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		body := this.parseStatement()
		return ast.While_new(exp, body, false, this.Linenum)
	case TOKEN_DO:
		log.Debugf("********Do***********")
		var fake = ast.FakeStm_new(this.Peek(), this.Linenum)
		this.currentStm = fake
		this.stmStack.Push(fake)
		this.Push(fake)
		defer func() {
			this.Pop()
			this.stmStack.Pop()
			if this.stmStack.Peek() != nil {
				this.currentStm = this.stmStack.Peek().(ast.Stm)
			} else {
				this.currentStm = nil
			}
		}()
		this.eatToken(TOKEN_DO)
		body := this.parseStatement()
		this.eatToken(TOKEN_WHILE)
		this.eatToken(TOKEN_LPAREN)
		exp := this.parseExp()
		this.eatToken(TOKEN_RPAREN)
		return ast.While_new(exp, body, true, this.Linenum)
	case TOKEN_SWITCH:
		log.Debugf("********TOKEN_SWITCH***********")
		var fake = ast.FakeStm_new(this.Peek(), this.Linenum)
		this.currentStm = fake
		this.stmStack.Push(fake)
		this.Push(fake)
		defer func() {
			this.Pop()
			this.stmStack.Pop()
			if this.stmStack.Peek() != nil {
				this.currentStm = this.stmStack.Peek().(ast.Stm)
			} else {
				this.currentStm = nil
			}
		}()
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
	case TOKEN_DEFAULT:
		this.eatToken(TOKEN_DEFAULT)
		this.eatToken(TOKEN_COLON)
		body := this.parseStatements()
		return ast.Case_new(nil, ast.Block_new(body, this.Linenum), this.Linenum)
	case TOKEN_FOR:
		log.Debugf("********TOKEN_FOR***********")
		var fake = ast.FakeStm_new(this.Peek(), this.Linenum)
		this.currentStm = fake
		this.stmStack.Push(fake)
		this.Push(fake)
		defer func() {
			this.Pop()
			this.stmStack.Pop()
			if this.stmStack.Peek() != nil {
				this.currentStm = this.stmStack.Peek().(ast.Stm)
			} else {
				this.currentStm = nil
			}
		}()
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
				Init = ast.Assign_new(exp, value, "=", false, this.Linenum)

			}

			//
			log.Debugf("********TOKEN_FOR--> 解析条件语句 ***********,%v", this.current.Lexeme)
			Condition := this.parseExp()
			this.eatToken(TOKEN_SEMI)

			log.Debugf("********TOKEN_FOR--> 解析更新语句 ***********")
			var Post ast.Stm
			if this.current.Kind == TOKEN_RPAREN {
				log.Debugf("********TOKEN_FOR--> 空更新语句 ***********")
			} else {
				this.isSpecial = true
				Post = this.parseStatement()
				this.isSpecial = false
			}

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

		return ast.Print_new(e, this.Linenum)
	case TOKEN_THROW:
		this.eatToken(TOKEN_THROW)
		e := this.parseExp()

		return ast.Throw_new(e, this.Linenum)
		//处理的是后缀加
	case TOKEN_INCREMENT:
		log.Debugf("处理累加")
		this.eatToken(TOKEN_INCREMENT)
		id := this.current.Lexeme
		this.eatToken(TOKEN_ID)
		//特殊的for语句才不需要分号

		left := ast.NewIdent(id, this.Linenum)

		return ast.Assign_new(left, &ast.Num{Value: 1}, "+=", false, this.Linenum)
	case TOKEN_RETURN:
		this.eatToken(TOKEN_RETURN)
		//空return
		if this.current.Kind == TOKEN_SEMI {

			return ast.Return_new(nil, this.Linenum)
		}
		log.Debugf("------>解析return,%v", this.current.Lexeme)
		exp := this.parseExp()

		//三元表达式
		if q, ok := exp.(*ast.Question); ok {
			assign1 := ast.Return_new(q.One, this.Linenum)
			assign2 := ast.Return_new(q.Two, this.Linenum)
			return ast.If_new(q.E, ast.Block_new([]ast.Stm{assign1}, this.Linenum), ast.Block_new([]ast.Stm{assign2}, this.Linenum), this.Linenum)
		}
		return ast.Return_new(exp, this.Linenum)
	default:
		if this.IsTypeToken() {
			tp := this.parseType()

			return this.parserDecl(tp)
		}
		this.ParseBug("代码段解析bug")
	}
	return nil
}

func (this *Parser) parseStatements() []ast.Stm {
	stms := []ast.Stm{}
	for this.IsTypeToken() ||
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
		this.current.Kind == TOKEN_DEFAULT ||
		this.current.Kind == TOKEN_NEW ||
		this.current.Kind == TOKEN_FINAL ||
		this.current.Kind == TOKEN_SYNCHRONIZED ||
		this.current.Kind == TOKEN_SUPER ||
		this.current.Kind == TOKEN_BREAK ||
		this.current.Kind == TOKEN_GOTO ||
		this.current.Kind == TOKEN_DO ||
		this.current.Kind == TOKEN_SYSTEM {
		if this.current.Kind == TOKEN_FINAL {
			this.eatToken(TOKEN_FINAL)
		}
		stms = append(stms, this.parseStatement())
	}
	return stms
}
