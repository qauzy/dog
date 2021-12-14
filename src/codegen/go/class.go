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
			//FIXME 是否排除static
			if !fi.IsStatic() && !cc.IsEnum() {
				gfi := this.transField(fi)
				Type.Fields.List = append(Type.Fields.List, gfi)
				this.buildFieldFunc(fi)
			}
		}
		for _, m := range cc.Methods {
			gmeth := this.transFunc(m)

			this.GolangFile.Decls = append(this.GolangFile.Decls, gmeth)
		}

		cl.Specs = append(cl.Specs, sp)

	}
	return
}

// 枚举转换
//
// param: c
func (this *Translation) transEnum(c ast.Class) {
	this.CurrentClass = c

	if cc, ok := c.(*ast.ClassSingle); ok {
		//1 定义枚举类型为int
		t := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		sp := &gast.TypeSpec{
			Doc:     nil,
			Name:    gast.NewIdent(cc.GetName()),
			Assign:  0,
			Type:    gast.NewIdent("int"),
			Comment: nil,
		}
		t.Specs = append(t.Specs, sp)
		this.GolangFile.Decls = append(this.GolangFile.Decls, t)

		//2 解析枚举元素
		v := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.CONST,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		this.GolangFile.Decls = append(this.GolangFile.Decls, v)
		for idx, fi := range cc.Fields {
			value := &gast.ValueSpec{
				Doc:     nil,
				Names:   []*gast.Ident{gast.NewIdent(fi.GetName())},
				Type:    nil,
				Values:  nil,
				Comment: nil,
			}
			if idx == 0 {
				value.Type = gast.NewIdent(cc.GetName())
				value.Values = append(value.Values, gast.NewIdent("iota"))
			}

			v.Specs = append(v.Specs, value)
		}

	}
}
