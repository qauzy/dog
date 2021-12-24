package tool

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func PrintAst() {
	// 这就是上一章的代码.
	src := `
package main
func (this *airdropDao) Save(m *entity.Airdrop) (err error) {
	err = this.DBWrite().Save(m).Error
	return
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)

}
