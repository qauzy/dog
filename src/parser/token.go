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
	tokenMap["+="] = TOKEN_ADD_ASSIGN
	tokenMap["++"] = TOKEN_AUTOADD
	tokenMap["-"] = TOKEN_SUB
	tokenMap["-="] = TOKEN_SUB_ASSIGN
	tokenMap["--"] = TOKEN_AUTOSUB

	tokenMap["&&"] = TOKEN_AND
	tokenMap["="] = TOKEN_ASSIGN
	tokenMap[","] = TOKEN_COMMER
	tokenMap["."] = TOKEN_DOT
	tokenMap[":"] = TOKEN_COLON
	tokenMap["?"] = TOKEN_QUESTION
	tokenMap["{"] = TOKEN_LBRACE
	tokenMap["["] = TOKEN_LBRACK
	tokenMap["("] = TOKEN_LPAREN
	tokenMap["<"] = TOKEN_LT
	tokenMap["<="] = TOKEN_LE
	tokenMap["=="] = TOKEN_EQ
	tokenMap[">"] = TOKEN_GT
	tokenMap[">="] = TOKEN_GE
	tokenMap["||"] = TOKEN_OR
	tokenMap["!="] = TOKEN_NE
	tokenMap["!"] = TOKEN_NOT
	tokenMap["}"] = TOKEN_RBRACE
	tokenMap["]"] = TOKEN_RBRACK
	tokenMap[")"] = TOKEN_RPAREN
	tokenMap[";"] = TOKEN_SEMI
	tokenMap["*"] = TOKEN_MUL
	tokenMap["*="] = TOKEN_MUL_ASSIGN
	tokenMap["/"] = TOKEN_QUO
	tokenMap["/="] = TOKEN_QUO_ASSIGN
	tokenMap["%"] = TOKEN_REM
	tokenMap["%="] = TOKEN_REM_ASSIGN
	tokenMap["@"] = TOKEN_AT
	tokenMap["->"] = TOKEN_LAMBDA

	tokenMap["if"] = TOKEN_IF
	tokenMap["else"] = TOKEN_ELSE
	tokenMap["while"] = TOKEN_WHILE
	tokenMap["for"] = TOKEN_FOR
	tokenMap["throws"] = TOKEN_THROWS
	tokenMap["throw"] = TOKEN_THROW

	tokenMap["switch"] = TOKEN_SWITCH
	tokenMap["case"] = TOKEN_CASE

	tokenMap["try"] = TOKEN_TRY
	tokenMap["catch"] = TOKEN_CATCH
	tokenMap["finally"] = TOKEN_FINALLY

	tokenMap["main"] = TOKEN_MAIN
	tokenMap["new"] = TOKEN_NEW
	tokenMap["false"] = TOKEN_FALSE
	tokenMap["true"] = TOKEN_TRUE
	tokenMap["null"] = TOKEN_NULL

	tokenMap["class"] = TOKEN_CLASS
	tokenMap["package"] = TOKEN_PACKAGE
	tokenMap["import"] = TOKEN_IMPORT
	tokenMap["implements"] = TOKEN_IMPLEMENTS

	tokenMap["extends"] = TOKEN_EXTENDS
	tokenMap["public"] = TOKEN_PUBLIC
	tokenMap["protected"] = TOKEN_PROTECTED
	tokenMap["private"] = TOKEN_PRIVATE
	tokenMap["final"] = TOKEN_FINAL
	tokenMap["abstract"] = TOKEN_ABSTRACT
	tokenMap["transient"] = TOKEN_TRANSIENT
	tokenMap["this"] = TOKEN_THIS
	tokenMap["static"] = TOKEN_STATIC
	tokenMap["length"] = TOKEN_LENGTH
	tokenMap["Size"] = TOKEN_SIZE
	tokenMap["return"] = TOKEN_RETURN
	tokenMap["System"] = TOKEN_SYSTEM

	tokenMap["out"] = TOKEN_OUT
	tokenMap["println"] = TOKEN_PRINTLN
	tokenMap["EOF"] = TOKEN_EOF

	//类型关键字
	tokenMap["int"] = TOKEN_INT
	tokenMap["byte"] = TOKEN_BYTE
	tokenMap["short"] = TOKEN_SHORT
	tokenMap["long"] = TOKEN_LONG
	tokenMap["float"] = TOKEN_FLOAT
	tokenMap["double"] = TOKEN_DOUBLE
	tokenMap["char"] = TOKEN_CHAR
	tokenMap["boolean"] = TOKEN_BOOLEAN
	tokenMap["void"] = TOKEN_VOID

	tokenMap["String"] = TOKEN_STRING
	tokenMap["Long"] = TOKEN_LONG
	tokenMap["Object"] = TOKEN_OBJECT
	tokenMap["Integer"] = TOKEN_INTEGER
	tokenMap["Set"] = TOKEN_SET
	tokenMap["HashSet"] = TOKEN_HASHSET
	tokenMap["Map"] = TOKEN_MAP
	tokenMap["HashMap"] = TOKEN_HASHMAP
	tokenMap["List"] = TOKEN_LIST
	tokenMap["ArrayList"] = TOKEN_ARRAYLIST

	tMap = make(map[int]string)

	tMap[TOKEN_ADD] = "TOKEN_ADD"               // +
	tMap[TOKEN_AUTOADD] = "TOKEN_AUTOADD"       // ++
	tMap[TOKEN_AND] = "TOKEN_AND"               // &&
	tMap[TOKEN_OR] = "TOKEN_OR"                 // ||
	tMap[TOKEN_NE] = "TOKEN_NE"                 // !=
	tMap[TOKEN_ASSIGN] = "TOKEN_ASSIGN"         // =
	tMap[TOKEN_COMMER] = "TOKEN_COMMER"         // ,
	tMap[TOKEN_DOT] = "TOKEN_DOT"               // .
	tMap[TOKEN_COLON] = "TOKEN_COLON"           // :
	tMap[TOKEN_QUESTION] = "TOKEN_QUESTION"     // ?
	tMap[TOKEN_AUTOSUB] = "TOKEN_AUTOSUB"       // --
	tMap[TOKEN_SUB] = "TOKEN_SUB"               // -
	tMap[TOKEN_LAMBDA] = "TOKEN_LAMBDA"         // ->
	tMap[TOKEN_MUL] = "TOKEN_MUL"               // *
	tMap[TOKEN_QUO] = "TOKEN_QUO"               // /
	tMap[TOKEN_QUO_ASSIGN] = "TOKEN_QUO_ASSIGN" // /=
	tMap[TOKEN_REM] = "TOKEN_REM"               // %
	tMap[TOKEN_REM_ASSIGN] = "TOKEN_REM_ASSIGN" // %=
	tMap[TOKEN_RBRACE] = "TOKEN_RBRACE"         // }
	tMap[TOKEN_RBRACK] = "TOKEN_RBRACK"         // ]
	tMap[TOKEN_RPAREN] = "TOKEN_RPAREN"         // )
	tMap[TOKEN_SEMI] = "TOKEN_SEMI"             // ;
	tMap[TOKEN_LBRACE] = "TOKEN_LBRACE"         // {
	tMap[TOKEN_LBRACK] = "TOKEN_LBRACK"         // [
	tMap[TOKEN_LPAREN] = "TOKEN_LPAREN"         // (
	tMap[TOKEN_LT] = "TOKEN_LT"                 // <
	tMap[TOKEN_LE] = "TOKEN_LE"                 // <=
	tMap[TOKEN_GE] = "TOKEN_GE"                 // >=
	tMap[TOKEN_EQ] = "TOKEN_EQ"                 // ==
	tMap[TOKEN_GT] = "TOKEN_GT"                 // >
	tMap[TOKEN_NOT] = "TOKEN_NOT"               // !
	tMap[TOKEN_AT] = "TOKEN_AT"                 // @

	tMap[TOKEN_IF] = "TOKEN_IF"
	tMap[TOKEN_ELSE] = "TOKEN_ELSE"
	tMap[TOKEN_FOR] = "TOKEN_FOR"
	tMap[TOKEN_WHILE] = "TOKEN_WHILE"

	tMap[TOKEN_TRUE] = "TOKEN_TRUE"
	tMap[TOKEN_FALSE] = "TOKEN_FALSE"
	tMap[TOKEN_TRY] = "TOKEN_TRY"
	tMap[TOKEN_CATCH] = "TOKEN_CATCH"
	tMap[TOKEN_FINALLY] = "TOKEN_FINALLY"

	tMap[TOKEN_CLASS] = "TOKEN_CLASS" //class
	tMap[TOKEN_PUBLIC] = "TOKEN_PUBLIC"
	tMap[TOKEN_PRIVATE] = "TOKEN_PRIVATE"
	tMap[TOKEN_PROTECTED] = "TOKEN_PROTECTED"
	tMap[TOKEN_PACKAGE] = "TOKEN_PACKAGE"
	tMap[TOKEN_IMPORT] = "TOKEN_IMPORT"
	tMap[TOKEN_IMPLEMENTS] = "TOKEN_IMPLEMENTS"
	tMap[TOKEN_FINAL] = "TOKEN_FINAL"
	tMap[TOKEN_ABSTRACT] = "TOKEN_ABSTRACT"
	tMap[TOKEN_TRANSIENT] = "TOKEN_TRANSIENT"
	tMap[TOKEN_EXTENDS] = "TOKEN_EXTENDS"

	tMap[TOKEN_THIS] = "TOKEN_THIS"
	tMap[TOKEN_NULL] = "TOKEN_NULL"
	tMap[TOKEN_ID] = "TOKEN_ID"

	tMap[TOKEN_LENGTH] = "TOKEN_LENGTH"
	tMap[TOKEN_SIZE] = "TOKEN_SIZE"
	tMap[TOKEN_MAIN] = "TOKEN_MAIN"
	tMap[TOKEN_NEW] = "TOKEN_NEW"
	tMap[TOKEN_THROWS] = "TOKEN_THROWS"
	tMap[TOKEN_THROW] = "TOKEN_THROW"
	tMap[TOKEN_SWITCH] = "TOKEN_SWITCH"
	tMap[TOKEN_CASE] = "TOKEN_CASE"

	tMap[TOKEN_NUM] = "TOKEN_NUM"
	tMap[TOKEN_OUT] = "TOKEN_OUT"
	tMap[TOKEN_PRINTLN] = "TOKEN_PRINTLN"

	tMap[TOKEN_RETURN] = "TOKEN_RETURN"
	tMap[TOKEN_STATIC] = "TOKEN_STATIC"
	tMap[TOKEN_SYSTEM] = "TOKEN_SYSTEM"
	tMap[TOKEN_EOF] = "TOKEN_EOF"

	//类型关键字
	tMap[TOKEN_VOID] = "TOKEN_VOID"
	tMap[TOKEN_INT] = "TOKEN_INT"
	tMap[TOKEN_BYTE] = "TOKEN_BYTE"
	tMap[TOKEN_SHORT] = "TOKEN_SHORT"
	tMap[TOKEN_LONG] = "TOKEN_LONG"
	tMap[TOKEN_FLOAT] = "TOKEN_FLOAT"
	tMap[TOKEN_DOUBLE] = "TOKEN_DOUBLE"
	tMap[TOKEN_CHAR] = "TOKEN_CHAR"
	tMap[TOKEN_BOOLEAN] = "TOKEN_BOOLEAN"

	tMap[TOKEN_STRING] = "TOKEN_STRING"
	tMap[TOKEN_OBJECT] = "TOKEN_OBJECT"
	tMap[TOKEN_INTEGER] = "TOKEN_INTEGER"
	tMap[TOKEN_SET] = "TOKEN_SET"
	tMap[TOKEN_HASHSET] = "TOKEN_HASHSET"
	tMap[TOKEN_MAP] = "TOKEN_MAP"
	tMap[TOKEN_HASHMAP] = "TOKEN_HASHMAP"
	tMap[TOKEN_LIST] = "TOKEN_LIST"
	tMap[TOKEN_ARRAYLIST] = "TOKEN_ARRAYLIST"

}

type Kind int

const (
	//运算符
	TOKEN_ADD        = iota // +
	TOKEN_ADD_ASSIGN        // +=
	TOKEN_AUTOADD           // ++
	TOKEN_SUB               // -
	TOKEN_SUB_ASSIGN        // -=
	TOKEN_AUTOSUB           // --
	TOKEN_AND               // &&
	TOKEN_ASSIGN            // =
	TOKEN_COMMER            // ,
	TOKEN_DOT               // .
	TOKEN_COLON             // :
	TOKEN_QUESTION          // ?
	TOKEN_LBRACE            // {
	TOKEN_RBRACE            // }
	TOKEN_LBRACK            // [
	TOKEN_RBRACK            // ]
	TOKEN_LPAREN            // (
	TOKEN_RPAREN            // )
	TOKEN_LT                // <
	TOKEN_LE                // <=
	TOKEN_EQ                // ==
	TOKEN_GT                // >
	TOKEN_GE                // >=
	TOKEN_OR                // ||
	TOKEN_NE                // !=
	TOKEN_NOT               // !
	TOKEN_MUL               // *
	TOKEN_MUL_ASSIGN        // *=
	TOKEN_QUO               // /
	TOKEN_QUO_ASSIGN        // /=
	TOKEN_REM               // %
	TOKEN_REM_ASSIGN        // %=
	TOKEN_SEMI              // ;
	TOKEN_AT                // @

	//关键字

	TOKEN_IF    // if
	TOKEN_ELSE  // else
	TOKEN_WHILE // while
	TOKEN_FOR   // for

	TOKEN_EOF     //
	TOKEN_CLASS   // class
	TOKEN_EXTENDS // extends
	TOKEN_FALSE   // false
	TOKEN_TRUE    // true
	TOKEN_ID

	TOKEN_TRY     // try
	TOKEN_CATCH   // catch
	TOKEN_FINALLY // finally
	TOKEN_THROWS  // throws
	TOKEN_THROW   // throw

	TOKEN_SWITCH // switch
	TOKEN_CASE   // case
	TOKEN_NULL   // null
	TOKEN_LENGTH // length
	TOKEN_SIZE   // size
	TOKEN_MAIN   // main
	TOKEN_NEW    // new

	TOKEN_PUBLIC    // public
	TOKEN_PROTECTED // protected
	TOKEN_PRIVATE   // private
	TOKEN_DEFAULT   // default
	TOKEN_ABSTRACT  // abstract
	TOKEN_TRANSIENT
	TOKEN_RETURN // return
	TOKEN_STATIC // static
	TOKEN_LAMBDA
	TOKEN_SYSTEM
	TOKEN_THIS // this

	TOKEN_PACKAGE    // package
	TOKEN_IMPORT     // import
	TOKEN_IMPLEMENTS // implements
	TOKEN_FINAL      // final
	TOKEN_CHARS

	TOKEN_NUM //
	TOKEN_OUT
	TOKEN_PRINTLN

	//类型TOKEN
	TOKEN_INT
	TOKEN_BYTE
	TOKEN_SHORT
	TOKEN_LONG
	TOKEN_FLOAT
	TOKEN_DOUBLE
	TOKEN_CHAR
	TOKEN_BOOLEAN

	TOKEN_VOID
	TOKEN_SET
	TOKEN_HASHSET
	TOKEN_MAP
	TOKEN_HASHMAP
	TOKEN_LIST
	TOKEN_ARRAYLIST

	TOKEN_STRING
	TOKEN_OBJECT
	TOKEN_INTEGER
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
