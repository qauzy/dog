package sqlx

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

func parseEx(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}
