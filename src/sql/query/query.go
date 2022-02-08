package query

import (
	"dog/ast"
	"dog/util"
	"fmt"
)

// Query represents a parsed sql
type Query struct {
	Type        Type
	TableName   string
	TableAliase string
	Conditions  []Condition
	Updates     map[string]string
	Inserts     [][]string
	Fields      []string // Used for SELECT (i.e. SELECTed field names) and INSERT (INSERTEDed field names)
	Aliases     map[string]string
}

func (q *Query) ToQorm(f []ast.Field) (exe string, err error) {
	tail := `
	err = eng.Error`

	exe = `
	eng := this.DBRead().Table("` + util.SnakeString(q.TableName) + `")`

	switch q.Type {
	case Select:
		exe = `
	//FIXME 非原生sql1，需要处理
	eng := this.DBRead().Table("` + util.SnakeString(q.TableName) + `")`
		conds := ""

		for idx, v := range q.Conditions {
			if v.Operand2 != "?" {
				conds += fmt.Sprintf(".Where(\"%s %s %s\")", v.Operand1, OperatorSrc[v.Operator], v.Operand2)
			} else if len(f) > idx {
				conds += fmt.Sprintf(".Where(\"%s %s %s\",%s)", v.Operand1, OperatorSrc[v.Operator], v.Operand2, f[idx].GetName())
			}

		}
		exe += conds

	// Update represents an UPDATE sql
	case Update:
	// Insert represents an INSERT sql
	case Insert:
	// Delete represents a DELETE sql
	case Delete:
	}
	exe += tail

	return
}

// Type is the type of SQL sql, e.g. SELECT/UPDATE
type Type int

const (
	// UnknownType is the zero value for a Type
	UnknownType Type = iota
	// Select represents a SELECT sql
	Select
	// Update represents an UPDATE sql
	Update
	// Insert represents an INSERT sql
	Insert
	// Delete represents a DELETE sql
	Delete
)

// TypeString is a string slice with the names of all types in order
var TypeString = []string{
	"UnknownType",
	"Select",
	"Update",
	"Insert",
	"Delete",
}

// Operator is between operands in a condition
type Operator int

const (
	// UnknownOperator is the zero value for an Operator
	UnknownOperator Operator = iota
	// Eq -> "="
	Eq
	// Ne -> "!="
	Ne
	// Gt -> ">"
	Gt
	// Lt -> "<"
	Lt
	// Gte -> ">="
	Gte
	// Lte -> "<="
	Lte
)

// OperatorString is a string slice with the names of all operators in order
var OperatorString = []string{
	"UnknownOperator",
	"Eq",
	"Ne",
	"Gt",
	"Lt",
	"Gte",
	"Lte",
}

var OperatorSrc = []string{
	"UnknownOperator",
	"=",
	"!=",
	">",
	"<",
	">=",
	"<=",
}

// Condition is a single boolean condition in a WHERE clause
type Condition struct {
	// Operand1 is the left hand side operand
	Operand1 string
	// Operand1IsField determines if Operand1 is a literal or a field name
	Operand1IsField bool
	// Operator is e.g. "=", ">"
	Operator Operator
	// Operand1 is the right hand side operand
	Operand2 string
	// Operand2IsField determines if Operand2 is a literal or a field name
	Operand2IsField bool
}
