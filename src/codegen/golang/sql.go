package codegen_go

import (
	"dog/log"
	"fmt"
	"github.com/qauzy/sqlparser"
	gast "go/ast"
	"strings"
)

func (this *Translation) transSQl(sql string, args string) (stmt gast.Stmt) {
	sql = strings.Trim(sql, "\"")
	stm, err := sqlparser.Parse(sql)
	if err != nil {
		log.Errorf("Query=%v,err=%v", sql, err)
	}

	switch st := stm.(type) {
	case *sqlparser.Select:
		exe := `eng := this.DBWrite()`
		exe = exe + fmt.Sprintf(".Table(\"%v\")", sqlparser.String(st.From))

		if st.SelectExprs != nil {
			exe = exe + fmt.Sprintf(".Select(\"%v\")", sqlparser.String(st.SelectExprs))
		}

		if st.Where != nil {
			var sql = sqlparser.String(st.Where.Expr)
			exe = exe + fmt.Sprintf(".Where(\"%v\"%v)", sql, args)
		}

		if st.GroupBy != nil {
			exe = exe + fmt.Sprintf(".Group(\"%v\")", sqlparser.String(st.GroupBy))
		}

		if st.Limit != nil && st.Limit.Offset != nil {
			exe = exe + fmt.Sprintf(".Offset(%v)", sqlparser.String(st.Limit.Offset))
		}
		if st.Limit != nil && st.Limit.Rowcount != nil {
			exe = exe + fmt.Sprintf(".Limit(%v)", sqlparser.String(st.Limit.Rowcount))
		}
		if st.OrderBy != nil {
			var ors string
			for _, vv := range st.OrderBy {
				if ors != "" {
					ors += ","
				}
				ors += sqlparser.String(vv)
			}
			exe = exe + fmt.Sprintf(".Order(\"%v\")", ors)
		}
		exe = exe + ".Find(&result)"
		exe += `
	err = eng.Error`

		stmt = &gast.ExprStmt{X: gast.NewIdent(exe)}
		return stmt

	case *sqlparser.Insert:

	case *sqlparser.Update:
		//exe := `eng := this.DBWrite()`
		//exe = exe + fmt.Sprintf(".Table(\"%v\")", sqlparser.String(st.TableExprs))
		//if st.Where != nil {
		//	exe = exe + fmt.Sprintf(".Where(\"%v\")", sqlparser.String(st.Where))
		//}
		//if st.Exprs != nil {
		//	exe = exe + fmt.Sprintf(".Where(\"%v\")", sqlparser.String(st.Exprs))
		//}
	case *sqlparser.Delete:
		//exe := `eng := this.DBWrite()`
		//exe = exe + fmt.Sprintf(".Table(\"%v\")", sqlparser.String(st.TableExprs))
		//if st.Where != nil {
		//	exe = exe + fmt.Sprintf(".Where(\"%v\")", sqlparser.String(st.Where))
		//}
		//if st.Exprs != nil {
		//	exe = exe + fmt.Sprintf(".Where(\"%v\")", sqlparser.String(st.Exprs))
		//}

	}

	return stmt
}
