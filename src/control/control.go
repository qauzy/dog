package control

var Lexer_test = false
var Lexer_dump = true

var CodeGen_base = ""
var CodeGen_fileName = ""
var CodeGen_outputName = ""

var CodeGen_dump = false

type CodeGen_Kind int

const (
	C = iota
	Go
	Bytecode
	Dalvik
	X86
)

var CodeGen_codegen CodeGen_Kind = Go

var Ast_test bool = false
var Ast_dumpAst bool = false

var Elab_classTable bool = false
var Elab_methodTable bool = false

const (
	None = iota
	Pdf
	Ps
	Jpg
	Svg
)

var Visualize_format int = None

var Optimization_Level int = 1
