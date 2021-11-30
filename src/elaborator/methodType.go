package elaborator

import (
	"dog/ast"
	"fmt"
)

type MethodType struct {
	retType  ast.Type
	argsType []ast.Field
}

func methodType_dump(m *MethodType) {
	fmt.Printf("        retType: ")
	fmt.Println(m.retType)
	for _, dec := range m.argsType {
		fmt.Printf("        ")
		fmt.Println(dec)
	}
}
