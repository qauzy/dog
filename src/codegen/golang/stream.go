package codegen_go

import (
	"go/ast"
	"go/token"
)

func (this *Translation) OptimitcStreamStm(src ast.Stmt) (dst ast.Stmt) {

	switch stm := src.(type) {
	case *ast.AssignStmt:
		lhs := stm.Lhs[0]
		rhs := stm.Rhs[0]
		if call, ok := rhs.(*ast.CallExpr); ok {
			// xxx.stream().min(xxx).orElse(x);
			if orElse, ok := call.Fun.(*ast.SelectorExpr); ok && (orElse.Sel.Name == "OrElse") && len(call.Args) == 1 {
				if call2, ok := orElse.X.(*ast.CallExpr); ok {
					//xxx.stream().min(xxx).orElse(x)
					if min, ok := call2.Fun.(*ast.SelectorExpr); ok && (min.Sel.Name == "Min") && len(call2.Args) == 1 {

					} else if findFirst, ok := call2.Fun.(*ast.SelectorExpr); ok && (findFirst.Sel.Name == "FindFirst") && len(call2.Args) == 0 {

					}
				}

			} else if get, ok := call.Fun.(*ast.SelectorExpr); ok && (get.Sel.Name == "Get") && len(call.Args) == 0 {
				if call2, ok := orElse.X.(*ast.CallExpr); ok {
					// xxx.stream().findFirst().get()
					if findFirst, ok := call2.Fun.(*ast.SelectorExpr); ok && (findFirst.Sel.Name == "FindFirst") && len(call2.Args) == 0 {

					}
				}

			} else if get, ok := call.Fun.(*ast.SelectorExpr); ok && (get.Sel.Name == "Count") && len(call.Args) == 0 {
				if call2, ok := orElse.X.(*ast.CallExpr); ok {
					// xxx.stream().findFirst().get()
					if filter, ok := call2.Fun.(*ast.SelectorExpr); ok && (filter.Sel.Name == "Filter") && len(call2.Args) == 1 {

					}
				}

			} else if parseObject, ok := call.Fun.(*ast.SelectorExpr); ok && (parseObject.Sel.Name == "ParseObject") && len(call.Args) == 2 {
				if vvvv, ok := parseObject.X.(*ast.Ident); ok && (vvvv.Name == "JSON") {
					//转换json解析
					call.Args = []ast.Expr{call.Args[0], lhs}
					call.Fun = ast.NewIdent("mdata.Cjson.Unmarshal")
					as := &ast.AssignStmt{
						Lhs:    []ast.Expr{ast.NewIdent("err")},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []ast.Expr{rhs},
					}
					fk := &FakeBlock{}
					fk.List = append(fk.List, as)
					fk.List = append(fk.List, this.GetErrReturn())
					return fk

				}
			}

		}
	case *ast.DeclStmt:
	}

	return src
}
