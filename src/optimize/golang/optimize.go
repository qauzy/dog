package golang

import (
	ast "go/ast"
	"strings"
)

type OptimizeFunc func(n ast.Node)

var StandardRules = []OptimizeFunc{
	TimeOp,
	ConstantOp, //处理枚举
	TypeOp,     //变量声明类型加*号
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
func ConstantOp(n ast.Node) {
	sl, ok := n.(*ast.SelectorExpr)
	if ok {
		if s, ok := sl.X.(*ast.Ident); ok {
			if s.Name == "constant" {
				s.Name = sl.Sel.Name
			}
		}
	}
}

//类对象类型加*号
func TypeOp(n ast.Node) {
	tp, ok := n.(*ast.StructType)
	if ok && tp.Fields != nil {
		for _, fi := range tp.Fields.List {
			switch t := fi.Type.(type) {
			case *ast.Ident:
				if !strings.HasPrefix(t.Name, "*") {
					fi.Type = &ast.StarExpr{X: fi.Type}
				}
			case *ast.SelectorExpr:
				if id, ok := t.X.(*ast.Ident); ok {
					if id.Name == "entity" {
						fi.Type = &ast.StarExpr{X: fi.Type}
					} else if id.Name == "util" {
						t.X = ast.NewIdent("*cache")
					}
				}

			}

		}
	}
}
