package parser

import (
	"dog/ast"
	"dog/util"
	log "github.com/corgi-kx/logcustom"
)

//
//
// param: exp 声明变量的类型
func (this *Parser) parserDecl(exp ast.Exp) ast.Stm {
	log.Debugf("*******解析临时变量声明语句*******")
	id := this.current.Lexeme
	this.eatToken(TOKEN_ID)
	decl := ast.DeclStmt_new(nil, exp, nil, this.Linenum)
	this.currentStm = decl
	defer func() {
		this.currentStm = nil
	}()

	//记录本地变量
	if this.currentMethod != nil {
		f := &ast.FieldSingle{
			Access:  0,
			Tp:      exp,
			Name:    util.GetNewId(id),
			Static:  false,
			IsField: false,
			Value:   nil,
		}
		this.currentMethod.AddLocals(f)
	}

	decl.Names = append(decl.Names, ast.NewIdent(util.GetNewId(id), this.Linenum))
	//有赋值语句
	if this.current.Kind == TOKEN_ASSIGN {
		this.assignType = exp
		//临时变量类型
		log.Debugf("*******解析临时变量声明语句(有赋值语句)*******")
		this.eatToken(TOKEN_ASSIGN)
		exp := this.parseExp()

		var call string
		var mp = new(ast.StreamStm)
		if this.CheckStreamExprs(exp, &call, mp) {
			mp.Left = decl.Names[0]
			mp.LineNum = this.Linenum
			decl.SetExtra(mp)
		}

		decl.Values = append(decl.Values, exp)
	}

	//定义多个变量
	for this.current.Kind == TOKEN_COMMER {
		this.advance()
		id = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		//记录本地变量
		if this.currentMethod != nil {
			f := &ast.FieldSingle{
				Access:  0,
				Tp:      exp,
				Name:    util.GetNewId(id),
				Static:  false,
				IsField: false,
				Value:   nil,
			}
			this.currentMethod.AddLocals(f)
		}
		decl.Names = append(decl.Names, ast.NewIdent(util.GetNewId(id), this.Linenum))

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
}

//检查List的map操作
func (this *Parser) CheckStreamExprs(exp ast.Exp, call *string, mp *ast.StreamStm) (b bool) {
	switch e := exp.(type) {
	case *ast.SelectorExpr:
		b = this.CheckStreamExprs(e.X, call, mp)
		*call = e.Sel
		return
	case *ast.CallExpr:
		b = this.CheckStreamExprs(e.Callee, call, mp)
		switch *call {
		case "Map":
			fallthrough
		case "oMap":
			fallthrough
		case "map":
			if len(e.ArgsList) == 1 {
				mp.Func = "map"
				mp.Ele = e.ArgsList[0]
				return
			} else {
				//log.Debugf("CheckStreamExprs %v", e)
				//time.Sleep(3 * time.Second)
				return false
			}
		case "stream":
			return
		case "Stream":
			return
		case "filter":
			mp.Func = "filter"
			return
		case "sorted":
			mp.Func = "sorted"
			return
		case "collect":
			if len(e.ArgsList) == 1 {
				if cc0, ok := e.ArgsList[0].(*ast.CallExpr); ok {
					if ssl0, ok := cc0.Callee.(*ast.SelectorExpr); ok {
						if ssl0.Sel == "toSet" || ssl0.Sel == "toList" || ssl0.Sel == "joining" {
							mp.ToAny = ssl0.Sel
							return
						}
					}
				}
			} else {
				return false
			}
			// Stream.of(xxxx)
		case "of":
			if len(e.ArgsList) == 1 {
				mp.List = e.ArgsList[0]
				return
			} else {
				return false
			}
		default:
			log.Debugf("false--%v-%v", *call, b)
			return
		}
	case *ast.Ident:
		if this.currentMethod != nil && this.currentMethod.GetLocals(e.Name) != nil {
			lo := this.currentMethod.GetLocals(e.Name)
			if _, ok := lo.GetDecType().(*ast.ListType); ok || *call == "stream" {
				mp.List = e
				return true
			} else {
				return false

			}
		} else if e.Name == "Stream" {
			return true
		}

	default:
		return false
	}

	return false
}
