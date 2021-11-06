package parser

import "fmt"
import "os"
import "strconv"

type Token struct {
	Kind    int
	Lexeme  string
	LineNum int
}

func newToken(kind int, lexeme string, lineNum int) *Token {
	return &Token{kind, lexeme, lineNum}
}

var tokenMap map[string]int
var tMap map[int]string

func initTokenMap() {
	tokenMap = make(map[string]int)
	tokenMap["+"] = TOKEN_ADD
	tokenMap["&&"] = TOKEN_AND
	tokenMap["="] = TOKEN_ASSIGN
	tokenMap["boolean"] = TOKEN_BOOLEAN
	tokenMap["class"] = TOKEN_CLASS
	tokenMap[","] = TOKEN_COMMER
	tokenMap["."] = TOKEN_DOT
	tokenMap["else"] = TOKEN_ELSE
	tokenMap["EOF"] = TOKEN_EOF
	tokenMap["extends"] = TOKEN_EXTENDS
	tokenMap["false"] = TOKEN_FALSE
	//id
	tokenMap["if"] = TOKEN_IF
	tokenMap["int"] = TOKEN_INT
	tokenMap["{"] = TOKEN_LBRACE
	tokenMap["["] = TOKEN_LBRACK
	tokenMap["length"] = TOKEN_LENGTH
	tokenMap["("] = TOKEN_LPAREN
	tokenMap["<"] = TOKEN_LT
	tokenMap["main"] = TOKEN_MAIN
	tokenMap["new"] = TOKEN_NEW
	tokenMap["!"] = TOKEN_NOT
	//num
	tokenMap["out"] = TOKEN_OUT
	tokenMap["println"] = TOKEN_PRINTLN
	tokenMap["public"] = TOKEN_PUBLIC
	tokenMap["abstract"] = TOKEN_ABSTRACT
	tokenMap["transient"] = TOKEN_TRANSIENT

	tokenMap["protected"] = TOKEN_PROTECTED
	tokenMap["}"] = TOKEN_RBRACE
	tokenMap["]"] = TOKEN_RBRACK
	tokenMap["return"] = TOKEN_RETURN
	tokenMap[")"] = TOKEN_RPAREN
	tokenMap[";"] = TOKEN_SEMI
	tokenMap["static"] = TOKEN_STATIC
	tokenMap["String"] = TOKEN_STRING
	tokenMap["-"] = TOKEN_SUB
	tokenMap["System"] = TOKEN_SYSTEM
	tokenMap["this"] = TOKEN_THIS
	tokenMap["*"] = TOKEN_TIMES
	tokenMap["true"] = TOKEN_TRUE
	tokenMap["void"] = TOKEN_VOID
	tokenMap["while"] = TOKEN_WHILE
	tokenMap["package"] = TOKEN_PACKAGE
	tokenMap["import"] = TOKEN_IMPORT
	tokenMap["@"] = TOKEN_AT
	tokenMap["implements"] = TOKEN_IMPLEMENTS
	tokenMap["private"] = TOKEN_PRIVATE

	tokenMap["final"] = TOKEN_FINAL

	tMap = make(map[int]string)
	tMap[TOKEN_ADD] = "TOKEN_ADD"
	tMap[TOKEN_AND] = "TOKEN_AND"
	tMap[TOKEN_ASSIGN] = "TOKEN_ASSIGN"
	tMap[TOKEN_BOOLEAN] = "TOKEN_BOOLEAN"
	tMap[TOKEN_CLASS] = "TOKEN_CLASS"
	tMap[TOKEN_COMMER] = "TOKEN_COMMER"
	tMap[TOKEN_DOT] = "TOKEN_DOT"
	tMap[TOKEN_ELSE] = "TOKEN_ELSE"
	tMap[TOKEN_EOF] = "TOKEN_EOF"
	tMap[TOKEN_EXTENDS] = "TOKEN_EXTENDS"
	tMap[TOKEN_FALSE] = "TOKEN_FALSE"
	tMap[TOKEN_IF] = "TOKEN_IF"
	tMap[TOKEN_INT] = "TOKEN_INT"
	tMap[TOKEN_ID] = "TOKEN_ID"
	tMap[TOKEN_LBRACE] = "TOKEN_LBRACE"
	tMap[TOKEN_LBRACK] = "TOKEN_LBRACK"
	tMap[TOKEN_LENGTH] = "TOKEN_LENGTH"
	tMap[TOKEN_LPAREN] = "TOKEN_LPAREN"
	tMap[TOKEN_LT] = "TOKEN_LT"
	tMap[TOKEN_MAIN] = "TOKEN_MAIN"
	tMap[TOKEN_NEW] = "TOKEN_NEW"
	tMap[TOKEN_NUM] = "TOKEN_NUM"
	tMap[TOKEN_NOT] = "TOKEN_NOT"
	tMap[TOKEN_OUT] = "TOKEN_OUT"
	tMap[TOKEN_PRINTLN] = "TOKEN_PRINTLN"
	tMap[TOKEN_PUBLIC] = "TOKEN_PUBLIC"
	tMap[TOKEN_ABSTRACT] = "TOKEN_ABSTRACT"
	tMap[TOKEN_TRANSIENT] = "TOKEN_TRANSIENT"

	tMap[TOKEN_RBRACE] = "TOKEN_RBRACE"
	tMap[TOKEN_RBRACK] = "TOKEN_RBRACK"
	tMap[TOKEN_RETURN] = "TOKEN_RETURN"
	tMap[TOKEN_RPAREN] = "TOKEN_RPAREN"
	tMap[TOKEN_SEMI] = "TOKEN_SEMI"
	tMap[TOKEN_STATIC] = "TOKEN_STATIC"
	tMap[TOKEN_STRING] = "TOKEN_STRING"
	tMap[TOKEN_SUB] = "TOKEN_SUB"
	tMap[TOKEN_SYSTEM] = "TOKEN_SYSTEM"
	tMap[TOKEN_TRUE] = "TOKEN_TRUE"
	tMap[TOKEN_THIS] = "TOKEN_THIS"
	tMap[TOKEN_TIMES] = "TOKEN_TIMES"
	tMap[TOKEN_VOID] = "TOKEN_VOID"
	tMap[TOKEN_WHILE] = "TOKEN_WHILE"
	tMap[TOKEN_PACKAGE] = "TOKEN_PACKAGE"
	tMap[TOKEN_IMPORT] = "TOKEN_IMPORT"
	tMap[TOKEN_AT] = "TOKEN_AT"
	tMap[TOKEN_IMPLEMENTS] = "TOKEN_IMPLEMENTS"
	tMap[TOKEN_PRIVATE] = "TOKEN_PRIVATE"
	tMap[TOKEN_PROTECTED] = "TOKEN_PROTECTED"
	tMap[TOKEN_FINAL] = "TOKEN_FINAL"

}

type Kind int

const (
	TOKEN_ADD = iota
	TOKEN_AND
	TOKEN_ASSIGN
	TOKEN_BOOLEAN
	TOKEN_CLASS
	TOKEN_COMMER
	TOKEN_DOT
	TOKEN_ELSE
	TOKEN_EOF
	TOKEN_EXTENDS
	TOKEN_FALSE
	TOKEN_ID
	TOKEN_IF
	TOKEN_INT
	TOKEN_LENGTH
	TOKEN_LBRACE
	TOKEN_LBRACK
	TOKEN_LPAREN
	TOKEN_LT
	TOKEN_MAIN
	TOKEN_NEW
	TOKEN_NOT
	TOKEN_NUM
	TOKEN_OUT
	TOKEN_PRINTLN
	TOKEN_PUBLIC
	TOKEN_PROTECTED
	TOKEN_DEFAULT
	TOKEN_PRIVATE
	TOKEN_ABSTRACT
	TOKEN_TRANSIENT

	TOKEN_RBRACE
	TOKEN_RBRACK
	TOKEN_RETURN
	TOKEN_RPAREN
	TOKEN_SEMI
	TOKEN_STATIC
	TOKEN_STRING
	TOKEN_SUB
	TOKEN_SYSTEM
	TOKEN_THIS
	TOKEN_TIMES
	TOKEN_TRUE
	TOKEN_VOID
	TOKEN_WHILE
	TOKEN_PACKAGE
	TOKEN_IMPORT
	TOKEN_AT
	TOKEN_IMPLEMENTS
	TOKEN_FINAL
)

func (this *Token) ToString() string {
	var s string
	if this.LineNum == 0 {
		fmt.Println("error")
		os.Exit(0)
	}

	s = ": " + this.Lexeme + " at LINE:" + strconv.Itoa(this.LineNum)
	return tMap[this.Kind] + s
}
