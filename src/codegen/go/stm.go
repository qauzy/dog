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
			Type:    transType(v.Tp),
			Values:  nil,
			Comment: nil,
		}
		//临时变量初值
		val := this.transExp(v.Value)
		if val != nil {
			sp.Values = append(sp.Values, val)
		} else {
			log.Debugf("初值为空")
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
			Post: nil,
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
	default:
		panic("bug")
		log.Debugf("transBlock-->%v -->%v", reflect.TypeOf(s).String(), v)

	}

	return
}
