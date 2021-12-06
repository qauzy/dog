package codegen_go

import (
	"dog/ast"
	gast "go/ast"
	"go/token"
)

//
//
// param: c
// return:
func (this *Translation) transClass(c ast.Class) (cl *gast.GenDecl) {
	this.CurrentClass = c
	if cc, ok := c.(*ast.ClassSingle); ok {
		cl = &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}

		sp := &gast.TypeSpec{
			Doc:     nil,
			Name:    gast.NewIdent(cc.Name),
			Assign:  0,
			Type:    nil,
			Comment: nil,
		}
		Type := &gast.StructType{
			Struct: 0,
			Fields: &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			},
			Incomplete: false,
		}
		sp.Type = Type

		for _, fi := range cc.Fields {
			gfi := this.transField(fi)
			Type.Fields.List = append(Type.Fields.List, gfi)
		}
		for _, m := range cc.Methods {
			gmeth := this.transFunc(m)

			//处理类接收
			recv := &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			}

			gfi := &gast.Field{
				Doc:   nil,
				Names: []*gast.Ident{gast.NewIdent("this")},
				Type: &gast.StarExpr{X: &gast.Ident{
					NamePos: 0,
					Name:    cc.Name,
					Obj:     gast.NewObj(gast.Typ, cc.Name),
				}},
				Tag:     nil,
				Comment: nil,
			}

			recv.List = append(recv.List, gfi)

			gmeth.Recv = recv
			this.GolangFile.Decls = append(this.GolangFile.Decls, gmeth)
		}

		cl.Specs = append(cl.Specs, sp)

	}
	return
}