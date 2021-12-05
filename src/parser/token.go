package parser

import (
	log "github.com/corgi-kx/logcustom"
	"os"
)
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
	tokenMap["++"] = TOKEN_AUTOADD
	tokenMap["&&"] = TOKEN_AND
	tokenMap["="] = TOKEN_ASSIGN
	tokenMap["boolean"] = TOKEN_BOOLEAN
	tokenMap["class"] = TOKEN_CLASS
	tokenMap[","] = TOKEN_COMMER
	tokenMap["."] = TOKEN_DOT
	tokenMap[":"] = TOKEN_COLON
	tokenMap["?"] = TOKEN_QUESTION
	tokenMap["else"] = TOKEN_ELSE
	tokenMap["EOF"] = TOKEN_EOF
	tokenMap["extends"] = TOKEN_EXTENDS
	tokenMap["false"] = TOKEN_FALSE
	//id
	tokenMap["if"] = TOKEN_IF
	tokenMap["try"] = TOKEN_TRY
	tokenMap["catch"] = TOKEN_CATCH
	tokenMap["finally"] = TOKEN_FINALLY
	tokenMap["null"] = TOKEN_NULL

	tokenMap["int"] = TOKEN_INT
	tokenMap["Object"] = TOKEN_OBJECT
	tokenMap["Long"] = TOKEN_LONG
	tokenMap["Integer"] = TOKEN_INTEGER

	tokenMap["Set"] = TOKEN_SET
	tokenMap["HashSet"] = TOKEN_HASHSET

	tokenMap["Map"] = TOKEN_MAP
	tokenMap["HashMap"] = TOKEN_HASHMAP
	tokenMap["List"] = TOKEN_LIST
	tokenMap["ArrayList"] = TOKEN_ARRAYLIST

	tokenMap["{"] = TOKEN_LBRACE
	tokenMap["["] = TOKEN_LBRACK
	tokenMap["length"] = TOKEN_LENGTH
	tokenMap["Size"] = TOKEN_SIZE
	tokenMap["("] = TOKEN_LPAREN
	tokenMap["<"] = TOKEN_LT
	tokenMap["<="] = TOKEN_LE
	tokenMap["=="] = TOKEN_EQ
	tokenMap[">"] = TOKEN_GT
	tokenMap[">="] = TOKEN_GE
	tokenMap["||"] = TOKEN_OR
	tokenMap["!="] = TOKEN_NE
	tokenMap["main"] = TOKEN_MAIN
	tokenMap["new"] = TOKEN_NEW
	tokenMap["throws"] = TOKEN_THROWS
	tokenMap["throw"] = TOKEN_THROW
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
	tokenMap["--"] = TOKEN_AUTOSUB
	tokenMap["System"] = TOKEN_SYSTEM
	tokenMap["->"] = TOKEN_LAMBDA
	tokenMap["this"] = TOKEN_THIS
	tokenMap["*"] = TOKEN_MUL
	tokenMap["/"] = TOKEN_QUO

	tokenMap["true"] = TOKEN_TRUE
	tokenMap["void"] = TOKEN_VOID
	tokenMap["while"] = TOKEN_WHILE
	tokenMap["for"] = TOKEN_FOR

	tokenMap["package"] = TOKEN_PACKAGE
	tokenMap["import"] = TOKEN_IMPORT
	tokenMap["@"] = TOKEN_AT
	tokenMap["implements"] = TOKEN_IMPLEMENTS
	tokenMap["private"] = TOKEN_PRIVATE

	tokenMap["final"] = TOKEN_FINAL

	tMap = make(map[int]string)

	tMap[TOKEN_ADD] = "TOKEN_ADD"
	tMap[TOKEN_AUTOADD] = "TOKEN_AUTOADD"
	tMap[TOKEN_AND] = "TOKEN_AND"
	tMap[TOKEN_OR] = "TOKEN_OR"
	tMap[TOKEN_NE] = "TOKEN_NE"
	tMap[TOKEN_ASSIGN] = "TOKEN_ASSIGN"
	tMap[TOKEN_BOOLEAN] = "TOKEN_BOOLEAN"
	tMap[TOKEN_CLASS] = "TOKEN_CLASS"
	tMap[TOKEN_COMMER] = "TOKEN_COMMER"
	tMap[TOKEN_DOT] = "TOKEN_DOT"
	tMap[TOKEN_COLON] = "TOKEN_COLON"
	tMap[TOKEN_QUESTION] = "TOKEN_QUESTION"
	tMap[TOKEN_ELSE] = "TOKEN_ELSE"
	tMap[TOKEN_EOF] = "TOKEN_EOF"
	tMap[TOKEN_EXTENDS] = "TOKEN_EXTENDS"
	tMap[TOKEN_FALSE] = "TOKEN_FALSE"
	tMap[TOKEN_IF] = "TOKEN_IF"
	tMap[TOKEN_TRY] = "TOKEN_TRY"
	tMap[TOKEN_CATCH] = "TOKEN_CATCH"
	tMap[TOKEN_FINALLY] = "TOKEN_FINALLY"
	tMap[TOKEN_NULL] = "TOKEN_NULL"
	tMap[TOKEN_INT] = "TOKEN_INT"
	tMap[TOKEN_OBJECT] = "TOKEN_OBJECT"
	tMap[TOKEN_LONG] = "TOKEN_LONG"

	tMap[TOKEN_INTEGER] = "TOKEN_INTEGER"
	tMap[TOKEN_SET] = "TOKEN_SET"
	tMap[TOKEN_HASHSET] = "TOKEN_HASHSET"
	tMap[TOKEN_MAP] = "TOKEN_MAP"
	tMap[TOKEN_HASHMAP] = "TOKEN_HASHMAP"
	tMap[TOKEN_LIST] = "TOKEN_LIST"
	tMap[TOKEN_ARRAYLIST] = "TOKEN_ARRAYLIST"
	tMap[TOKEN_ID] = "TOKEN_ID"
	tMap[TOKEN_LBRACE] = "TOKEN_LBRACE"
	tMap[TOKEN_LBRACK] = "TOKEN_LBRACK"
	tMap[TOKEN_LENGTH] = "TOKEN_LENGTH"
	tMap[TOKEN_SIZE] = "TOKEN_SIZE"
	tMap[TOKEN_LPAREN] = "TOKEN_LPAREN"
	tMap[TOKEN_LT] = "TOKEN_LT"
	tMap[TOKEN_LE] = "TOKEN_LE"
	tMap[TOKEN_GE] = "TOKEN_GE"
	tMap[TOKEN_EQ] = "TOKEN_EQ"
	tMap[TOKEN_GT] = "TOKEN_GT"
	tMap[TOKEN_MAIN] = "TOKEN_MAIN"
	tMap[TOKEN_NEW] = "TOKEN_NEW"
	tMap[TOKEN_THROWS] = "TOKEN_THROWS"
	tMap[TOKEN_THROW] = "TOKEN_THROW"
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
	//tMap[TOKEN_STRING] = "TOKEN_STRING"
	tMap[TOKEN_AUTOSUB] = "TOKEN_AUTOSUB"
	tMap[TOKEN_SUB] = "TOKEN_SUB"
	tMap[TOKEN_SYSTEM] = "TOKEN_SYSTEM"
	tMap[TOKEN_LAMBDA] = "TOKEN_LAMBDA"
	tMap[TOKEN_TRUE] = "TOKEN_TRUE"
	tMap[TOKEN_THIS] = "TOKEN_THIS"
	tMap[TOKEN_MUL] = "TOKEN_MUL"
	tMap[TOKEN_QUO] = "TOKEN_QUO"

	tMap[TOKEN_VOID] = "TOKEN_VOID"
	tMap[TOKEN_WHILE] = "TOKEN_WHILE"
	tMap[TOKEN_FOR] = "TOKEN_FOR"
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
	TOKEN_AUTOADD
	TOKEN_AND
	TOKEN_ASSIGN
	TOKEN_BOOLEAN
	TOKEN_CLASS
	TOKEN_COMMER
	TOKEN_DOT
	TOKEN_COLON
	TOKEN_QUESTION
	TOKEN_ELSE
	TOKEN_EOF
	TOKEN_EXTENDS
	TOKEN_FALSE
	TOKEN_ID
	TOKEN_IF
	TOKEN_TRY
	TOKEN_CATCH
	TOKEN_FINALLY
	TOKEN_NULL

	TOKEN_LENGTH
	TOKEN_SIZE
	TOKEN_LBRACE
	TOKEN_LBRACK
	TOKEN_LPAREN
	TOKEN_LT
	TOKEN_LE
	TOKEN_EQ
	TOKEN_GT
	TOKEN_GE
	TOKEN_OR
	TOKEN_NE
	TOKEN_MAIN
	TOKEN_NEW
	TOKEN_THROWS
	TOKEN_THROW
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

	TOKEN_SUB
	TOKEN_AUTOSUB
	TOKEN_LAMBDA
	TOKEN_SYSTEM
	TOKEN_THIS
	TOKEN_MUL
	TOKEN_QUO
	TOKEN_TRUE
	TOKEN_VOID
	TOKEN_WHILE
	TOKEN_FOR
	TOKEN_PACKAGE
	TOKEN_IMPORT
	TOKEN_AT
	TOKEN_IMPLEMENTS
	TOKEN_FINAL
	TOKEN_CHARS

	//类型TOKEN
	TOKEN_STRING
	TOKEN_INT
	TOKEN_OBJECT
	TOKEN_LONG
	TOKEN_INTEGER

	TOKEN_SET
	TOKEN_HASHSET
	TOKEN_MAP
	TOKEN_HASHMAP
	TOKEN_LIST
	TOKEN_ARRAYLIST
)

func (this *Token) ToString() string {
	var s string
	if this.LineNum == 0 {
		log.Info("error")
		os.Exit(0)
	}

	s = ": " + this.Lexeme + " at LINE:" + strconv.Itoa(this.LineNum)
	return tMap[this.Kind] + s
}
