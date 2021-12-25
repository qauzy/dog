package golang

import ast "go/ast"

type OptimizeFunc func(n ast.Node)

var StandardRules = []OptimizeFunc{
	TimeOp,
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
