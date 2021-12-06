package codegen_go

import (
	"dog/ast"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"go/token"
	"reflect"
)

func TransGo(p ast.File) (f *gast.File) {
	trans := NewTranslation(p)
	trans.ParseClasses()
	trans.WriteFile()
	return trans.GolangFile
}

// 带类型的变量声明
//
// param: fi
// return:
func (this *Translation) transField(fi ast.Field) (gfi *gast.Field) {
	this.CurrentField = fi
	if field, ok := fi.(*ast.FieldSingle); ok {
		//只处理成员变量
		var name = field.Name
		if field.IsField {
			name = Capitalize(field.Name)
		}
		gfi = &gast.Field{
			Doc:     nil,
			Names:   []*gast.Ident{gast.NewIdent(name)},
			Type:    transType(field.Tp),
			Tag:     nil,
			Comment: nil,
		}

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

// 三元表达式
//
// param: v
// return:
func (this *Translation) transTriple(s ast.Stm) (stmts []gast.Stmt) {
	if !s.IsTriple() {
		panic("should triple expr")
	}
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
		d.Specs = append(d.Specs, sp)
		stmts = append(stmts, &gast.DeclStmt{Decl: d})

		if vv, ok := v.Value.(*ast.Question); ok {
			q := &gast.IfStmt{
				If:   0,
				Init: nil,
				Cond: this.transExp(vv.E),
				Body: &gast.BlockStmt{
					Lbrace: 0,
					List: []gast.Stmt{&gast.AssignStmt{
						Lhs:    []gast.Expr{gast.NewIdent(v.Name)},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []gast.Expr{this.transExp(vv.One)}}},
				},

				Else: &gast.BlockStmt{
					Lbrace: 0,
					List: []gast.Stmt{&gast.AssignStmt{
						Lhs:    []gast.Expr{gast.NewIdent(v.Name)},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []gast.Expr{this.transExp(vv.Two)}}},
					Rbrace: 0,
				},
			}

			stmts = append(stmts, q)
		} else {
			panic("should triple expr")
		}

		//赋值语句
	case *ast.Assign:
		if vv, ok := v.E.(*ast.Question); ok {
			q := &gast.IfStmt{
				If:   0,
				Init: nil,
				Cond: this.transExp(vv.E),
				Body: &gast.BlockStmt{
					Lbrace: 0,
					List: []gast.Stmt{&gast.AssignStmt{
						Lhs:    []gast.Expr{this.transExp(v.Left)},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []gast.Expr{this.transExp(vv.One)}}},
				},

				Else: &gast.BlockStmt{
					Lbrace: 0,
					List: []gast.Stmt{&gast.AssignStmt{
						Lhs:    []gast.Expr{this.transExp(v.Left)},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []gast.Expr{this.transExp(vv.Two)}}},
					Rbrace: 0,
				},
			}

			stmts = append(stmts, q)
		} else {
			panic("should triple expr")
		}
	//
	case *ast.ExprStm:
		//log.Debugf("三元表达式语句:%v", v)
		stmt := &gast.ExprStmt{X: this.transExp(v.E)}
		stmts = append(stmts, stmt)
	default:
		panic("bug")
		log.Debugf("transBlock-->%v -->%v", reflect.TypeOf(s).String(), v)

	}

	return

}

func transType(t ast.Type) (Type gast.Expr) {
	switch v := t.(type) {
	case *ast.Void:
		return nil
	case *ast.String:
		return gast.NewIdent("string")
	case *ast.StringArray:
		return gast.NewIdent("[]string")
	case *ast.Integer:
		return gast.NewIdent("int64")
	case *ast.Int:
		return gast.NewIdent("int")
	case *ast.IntArray:
		return gast.NewIdent("[]int")
	case *ast.HashType:
		return &gast.MapType{
			Map:   0,
			Key:   transType(v.Key),
			Value: transType(v.Value),
		}
	case *ast.ListType:
		return &gast.ArrayType{
			Lbrack: 0,
			Len:    nil,
			Elt:    transType(v.Ele),
		}
	case *ast.ClassType:
		return &gast.Ident{
			NamePos: 0,
			Name:    v.Name,
			Obj:     gast.NewObj(gast.Typ, v.Name),
		}
	case *ast.ObjectType:
		//return gast.NewIdent("interface{}")
		return &gast.InterfaceType{
			Interface: 0,
			Methods: &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			},
			Incomplete: false,
		}
	case *ast.Boolean:
		return gast.NewIdent("bool")
		//泛型
		//TODO 先用接口替代
	case *ast.GenericType:
		return &gast.InterfaceType{
			Interface: 0,
			Methods: &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			},
			Incomplete: false,
		}
	default:

		panic("impossible")
		//log.Info(v.String())
	}
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
