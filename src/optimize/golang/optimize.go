package golang

import (
	"dog/storage"
	ast "go/ast"
	"strings"
)

type OptimizeFunc func(n ast.Node)

var StandardRules = []OptimizeFunc{
	TimeOp,
	ConstantOp, //处理枚举
	//TypeOp,     //变量声明类型加*号
}

func TimeOp(n ast.Node) {
	sl, ok := n.(*ast.SelectorExpr)
	if ok {
		if s, ok := sl.X.(*ast.Ident); ok {
			if s.Name == "TimeUnit" {
				s.Name = "time"
				switch sl.Sel.Name {
				case "HOURS":
					sl.Sel.Name = "Hour"
				case "MINUTES":
					sl.Sel.Name = "Minute"
				case "SECONDS":
					sl.Sel.Name = "Second"
				case "MILLISECONDS":
					sl.Sel.Name = "Millisecond"
				case "MICROSECONDS":
					sl.Sel.Name = "Microsecond"
				}
			}
		}
	}
}
func RedisUtilOp(n ast.Node) {
	sl, ok := n.(*ast.SelectorExpr)
	if ok {
		if s, ok := sl.X.(*ast.Ident); ok {
			if s.Name == "util" && sl.Sel.Name == "RedisUtil" {
				s.Name = "*cache"
			}
		}
	}
}

//处理枚举变量
func ConstantOp(n ast.Node) {
	sl, ok := n.(*ast.SelectorExpr)
	if ok {
		if s, ok := sl.X.(*ast.Ident); ok {
			if s.Name == "constant" {
				s.Name = sl.Sel.Name
			}
		}
	} else {

	}
}

//类对象类型加*号
func TypeOp(n ast.Node) {

	switch tp := n.(type) {

	//处理类成员类型加星号
	case *ast.StructType:
		if tp.Fields != nil {
			for _, fi := range tp.Fields.List {
				switch t := fi.Type.(type) {
				case *ast.Ident:
					if !strings.HasPrefix(t.Name, "*") {
						result, err := storage.FindByName(t.Name)
						if err == nil && result.Kind == 0 {
							fi.Type = &ast.StarExpr{X: fi.Type}
						}

					}

				case *ast.SelectorExpr:
					if id, ok := t.X.(*ast.Ident); ok {
						if id.Name == "entity" {
							fi.Type = &ast.StarExpr{X: fi.Type}
						} else if id.Name == "util" {
							t.X = ast.NewIdent("*cache")
						} else {
							result, err := storage.FindByName(t.Sel.Name)
							if err == nil && result.Kind == 0 {
								fi.Type = &ast.StarExpr{X: fi.Type}
							}
						}
					}
				}

			}
		}

		//处理类成员函数参数类型加星号
	case *ast.FuncDecl:
		if tp.Recv != nil && len(tp.Recv.List) == 1 && len(tp.Recv.List[0].Names) == 1 {
			if strings.HasPrefix(tp.Recv.List[0].Names[0].Name, "Set") {

			} else if strings.HasPrefix(tp.Recv.List[0].Names[0].Name, "Get") {

			}
		}
	}

}
