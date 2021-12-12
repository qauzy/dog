package codegen_go

import (
	"dog/ast"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"go/token"
	"reflect"
)

//
//
// param: s
// return:
func (this *Translation) transStm(s ast.Stm) (stmt gast.Stmt) {
	switch v := s.(type) {
	//变量声明
	case *ast.Decl:

		//log.Info("变量声明:", v.Name, "行:", v.LineNum)
		d := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.VAR,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		sp := &gast.ValueSpec{
			Doc:     nil,
			Names:   []*gast.Ident{gast.NewIdent(v.Name)},
			Type:    this.transType(v.Tp),
			Values:  nil,
			Comment: nil,
		}
		//临时变量初值
		if v.Value != nil {
			val := this.transExp(v.Value)
			if val != nil {
				sp.Values = append(sp.Values, val)
			} else {
				log.Debugf("初值为空")
			}
		}

		d.Specs = append(d.Specs, sp)
		stmt = &gast.DeclStmt{Decl: d}
		//赋值语句
	case *ast.Assign:
		stmt = &gast.AssignStmt{
			Lhs:    []gast.Expr{this.transExp(v.Left)},
			TokPos: 0,
			Tok:    token.ASSIGN,
			Rhs:    []gast.Expr{this.transExp(v.E)},
		}
	case *ast.If:
		var el gast.Stmt
		if v.Elsee != nil {
			el = this.transBlock(v.Elsee)
		} else {
			el = nil
		}

		stmt = &gast.IfStmt{
			If:   0,
			Init: nil,
			Cond: this.transExp(v.Condition),
			Body: this.transBlock(v.Body),
			Else: el,
		}

	//
	case *ast.For:
		//log.Debugf("For语句:%v", *v)
		stmt = &gast.ForStmt{
			For:  0,
			Init: this.transStm(v.Init),
			Cond: this.transExp(v.Cond),
			Post: this.transStm(v.Post),
			Body: this.transBlock(v.Body),
		}

	case *ast.Range:
		//log.Debugf("Range语句:%v", *v)
		stmt = &gast.RangeStmt{
			For:    0,
			Key:    gast.NewIdent("_"),
			Value:  this.transExp(v.Value),
			TokPos: 0,
			Tok:    token.DEFINE,
			X:      this.transExp(v.E),
			Body:   this.transBlock(v.Body),
		}
	case *ast.ExprStm:
		//log.Debugf("表达式语句:%v", v)
		stmt = &gast.ExprStmt{X: this.transExp(v.E)}
	case *ast.Throw:

		log.Debugf("Throw语句:%v", v)
		stmt = &gast.ReturnStmt{
			Return:  0,
			Results: nil,
		}
	case *ast.Return:
		//log.Debugf("Return语句:%v", v)
		result := &gast.ReturnStmt{
			Return:  0,
			Results: nil,
		}
		if v.E != nil {
			ret := this.transExp(v.E)
			if ret != nil {
				result.Results = append(result.Results, ret)
			}
		} else {
			log.Debugf("空Return语句")
		}
		return result
	case *ast.Try:
		return this.transStm(v.Test)
	case *ast.While:
		stmt = &gast.ForStmt{
			For:  0,
			Init: nil,
			Cond: this.transExp(v.E),
			Post: nil,
			Body: this.transBlock(v.Body),
		}
	case *ast.Block:
		return this.transBlock(v)
	case *ast.Binary:
		var opt token.Token
		switch v.Opt {
		case "*=":
			opt = token.MUL_ASSIGN
		case "/=":
			opt = token.QUO_ASSIGN
		case "+=":
			opt = token.ADD_ASSIGN
		case "-=":
			opt = token.SUB_ASSIGN
		case "%=":
			opt = token.REM_ASSIGN
		default:
			panic("*ast.Binary")
		}
		stmt = &gast.AssignStmt{
			Lhs:    []gast.Expr{this.transExp(v.Left)},
			TokPos: 0,
			Tok:    opt,
			Rhs:    []gast.Expr{this.transExp(v.Right)},
		}
	case *ast.Switch:
		return &gast.SwitchStmt{
			Switch: 0,
			Init:   nil,
			Tag:    this.transExp(v.E),
			Body:   this.transBlock(v.Cases),
		}
	case *ast.Case:
		cs := &gast.CaseClause{
			Case:  0,
			List:  []gast.Expr{this.transExp(v.E)},
			Colon: 0,
		}

		body := this.transStm(v.Body)
		if bb, ok := body.(*gast.BlockStmt); ok {
			cs.Body = bb.List
		} else {
			cs.Body = append(cs.Body, body)
		}
		return cs
	case *ast.Comment:
		stmt = &gast.ExprStmt{X: gast.NewIdent(v.C)}
	default:
		log.Debugf("transBlock-->%v -->%v", reflect.TypeOf(s).String(), v)
		panic("bug")

	}

	return
}

func (this *Translation) transBlock(s ast.Stm) (block *gast.BlockStmt) {
	//log.Debugf("解析Block语句")
	block = new(gast.BlockStmt)
	if bl, ok := s.(*ast.Block); ok {
		for _, st := range bl.Stms {
			if st != nil {
				if st.IsTriple() {
					sss := this.transTriple(st)
					if sss != nil && len(sss) > 0 {
						block.List = append(block.List, sss...)
					}
				} else {
					block.List = append(block.List, this.transStm(st))
				}
			}
		}
	} else {
		log.Debugf("transBlock-->%v--->%v", reflect.TypeOf(s).String(), s)
		panic("bug")
	}

	return
}
