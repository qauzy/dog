package codegen_go

import (
	"dog/ast"
	"dog/cfg"
	"dog/util"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"go/token"
	"path"
	"reflect"
)

func TransGo(p ast.File, base string, file string) (f *gast.File) {
	trans := NewTranslation(file, p)
	trans.ParseClasses()

	trans.WriteFile(base, file)

	return trans.GolangFile
}

// 带类型的变量声明
//
// param: fi
// return:
func (this *Translation) transField(fi ast.Field) (gfi *gast.Field) {
	this.currentField = fi
	if field, ok := fi.(*ast.FieldSingle); ok {
		//只处理成员变量
		gfi = &gast.Field{
			Doc:     nil,
			Names:   []*gast.Ident{this.transNameExp(field.Name)},
			Type:    this.transType(field.Tp),
			Tag:     nil,
			Comment: nil,
		}

		tag := &gast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("`gorm:\"column:%v\" json:\"%v\"`", util.SnakeString(field.Name.Name), field.Name.Name),
		}
		gfi.Tag = tag
	}
	return
}

// 添加错误返回
//
// param: fi
// return:
func (this *Translation) getErrRet() (gfi *gast.Field) {

	gfi = &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{gast.NewIdent("err")},
		Type:    gast.NewIdent("error"),
		Tag:     nil,
		Comment: nil,
	}

	return
}
func (this *Translation) TranslationBug(v interface{}) {
	var msg = fmt.Sprintf("未处理 [%v] %s,%v\n", reflect.TypeOf(v).String(), path.Base(this.file), v)
	util.Bug(msg)
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
	case *ast.DeclStmt:
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
			Names:   nil,
			Type:    this.transType(v.Tp),
			Values:  nil,
			Comment: nil,
		}

		for _, name := range v.Names {
			sp.Names = append(sp.Names, this.transNameExp(name))
		}

		//临时变量初值
		d.Specs = append(d.Specs, sp)
		stmts = append(stmts, &gast.DeclStmt{Decl: d})

		for idx, value := range v.Values {
			if vv, ok := value.(*ast.Question); ok {
				q := &gast.IfStmt{
					If:   0,
					Init: nil,
					Cond: this.transExp(vv.E),
					Body: &gast.BlockStmt{
						Lbrace: 0,
						List: []gast.Stmt{&gast.AssignStmt{
							Lhs:    []gast.Expr{this.transExp(v.Names[idx])},
							TokPos: 0,
							Tok:    token.ASSIGN,
							Rhs:    []gast.Expr{this.transExp(vv.One)}}},
					},

					Else: &gast.BlockStmt{
						Lbrace: 0,
						List: []gast.Stmt{&gast.AssignStmt{
							Lhs:    []gast.Expr{this.transExp(v.Names[idx])},
							TokPos: 0,
							Tok:    token.ASSIGN,
							Rhs:    []gast.Expr{this.transExp(vv.Two)}}},
						Rbrace: 0,
					},
				}

				stmts = append(stmts, q)
			} else {
				this.TranslationBug("必须是三目运算表达式")
			}

		}

		//赋值语句
	case *ast.Assign:
		if vv, ok := v.Value.(*ast.Question); ok {
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

func (this *Translation) transType(t ast.Exp) (Type gast.Expr) {
	switch v := t.(type) {
	case *ast.SelectorExpr:
		return &gast.SelectorExpr{
			X:   this.transExp(v.X),
			Sel: gast.NewIdent(v.Sel),
		}
	case *ast.Ident:
		if this.currentFile != nil && (this.currentFile.GetImport(v.Name) != nil) {
			pack := this.currentFile.GetImport(v.Name).GetPack()
			if pack == "time" && (v.Name == "LocalDate" || v.Name == "LocalDateTime") {
				pack = v.Name
			}
			expr := &gast.SelectorExpr{
				X:   gast.NewIdent(pack),
				Sel: gast.NewIdent(v.Name),
			}

			return expr
		}
		return gast.NewIdent(v.Name)
	case *ast.Void:
		return nil
	case *ast.String:
		return gast.NewIdent("string")
	case *ast.StringArray:
		return gast.NewIdent("[]string")
	case *ast.Integer:
		return gast.NewIdent("int")
	case *ast.Int:
		return gast.NewIdent("int")
	case *ast.MapType:

		if cfg.MapListIdxAccess {
			return &gast.MapType{
				Map:   0,
				Key:   this.transType(v.Key),
				Value: this.transType(v.Value),
			}
		} else {
			return &gast.IndexListExpr{
				X:       gast.NewIdent("*hashmap.Map"),
				Lbrack:  0,
				Indices: []gast.Expr{this.transType(v.Key), this.transType(v.Value)},
				Rbrack:  0,
			}
		}
	case *ast.ListType:

		if cfg.MapListIdxAccess {
			return &gast.ArrayType{
				Lbrack: 0,
				Len:    nil,
				Elt:    this.transType(v.Ele),
			}
		} else {
			return &gast.IndexExpr{
				X:      gast.NewIdent("*arraylist.List"),
				Lbrack: 0,
				Index:  this.transType(v.Ele),
				Rbrack: 0,
			}
		}
	case *ast.SetType:

		if cfg.MapListIdxAccess {
			return &gast.MapType{
				Map:   0,
				Key:   this.transType(v.Ele),
				Value: gast.NewIdent("string"),
			}
		} else {
			return &gast.IndexExpr{
				X:      gast.NewIdent("*hashset.Set"),
				Lbrack: 0,
				Index:  this.transType(v.Ele),
				Rbrack: 0,
			}
		}

	case *ast.ClassType:
		if this.currentFile != nil && (this.currentFile.GetImport(v.Name) != nil) {
			pack := this.currentFile.GetImport(v.Name).GetPack()
			if pack == "time" && (v.Name == "LocalDate" || v.Name == "LocalDateTime") {
				pack = v.Name
			}
			expr := &gast.SelectorExpr{
				X:   gast.NewIdent(pack),
				Sel: gast.NewIdent(util.GetNewId(v.Name)),
			}

			return expr
		}
		return &gast.Ident{
			NamePos: 0,
			Name:    v.Name,
			Obj:     gast.NewObj(gast.Typ, v.Name),
		}
	case *ast.ObjectType:
		return &gast.InterfaceType{
			Interface: 0,
			Methods: &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 1,
			},
			Incomplete: false,
		}

	case *ast.ObjectArray:
		return &gast.ArrayType{
			Lbrack: 0,
			Len:    nil,
			Elt:    gast.NewIdent("interface{}"),
		}
	case *ast.Boolean:
		return gast.NewIdent("bool")
	case *ast.Byte:
		return gast.NewIdent("byte")
	case *ast.ByteArray:
		return gast.NewIdent("[]byte")
		//泛型
		//TODO 先用接口替代
	case *ast.GenericType:
		return &gast.InterfaceType{
			Interface: 0,
			Methods: &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 1,
			},
			Incomplete: false,
		}
	case *ast.Float:
		return gast.NewIdent("float64")
	case *ast.Date:
		return gast.NewIdent("time.Time")
	case *ast.ArrayType:
		return &gast.ArrayType{
			Lbrack: 0,
			Len:    nil,
			Elt:    this.transExp(v.Ele),
		}

	default:
		//return nil
		log.Info(v)
		panic("impossible")

	}
}
