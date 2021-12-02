package main

import (
	"dog/ast"
	codegen_go "dog/codegen/go"
	"dog/control"
	"dog/parser"
	"dog/util"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"io/ioutil"
	"os"
)

func dog_Parser(filename string, buf []byte) ast.File {
	return parser.NewParse(filename, buf).Parser()
}

func main() {
	args := os.Args[1:len(os.Args)]
	filename := control.Do_arg(args)
	if filename == "" {
		control.Usage()
		os.Exit(0)
	}
	control.CodeGen_fileName = filename
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Info(err)
		os.Exit(0)
	}
	if control.Lexer_test == true {
		lex := parser.NewLexer(filename, buf)
		tk := lex.NextToken()
		for tk.Kind != parser.TOKEN_EOF {
			log.Info(tk.ToString())
			tk = lex.NextToken()
		}
		log.Info(tk.ToString())
		os.Exit(0)
	}
	var Ast ast.File
	//setp1: lexer&&parser
	control.Verbose("parser", func() {
		Ast = dog_Parser(filename, buf)
	}, control.VERBOSE_PASS)
	if control.Ast_dumpAst == true {
		ast.NewPP().DumpProg(Ast)
	}
	//step2: elaborate -- 细化
	control.Verbose("Elaborate", func() {
		//elaborator.Elaborate(Ast)
	}, control.VERBOSE_PASS)

	control.Verbose("ast-Opt", func() {
		//Ast = ast_opt.Opt(Ast)
	}, control.VERBOSE_PASS)

	//set3: trans -- 翻译
	var Ast_go *gast.File
	control.Verbose("Transaction", func() {
		switch control.CodeGen_codegen {
		case control.Go:
			Ast_go = codegen_go.TransGo(Ast)
		case control.C:
			//Ast_c = codegen_c.TransC(Ast)
		case control.Bytecode:
			util.Todo()
		case control.Dalvik:
			util.Todo()
		case control.X86:
			util.Todo()
		default:
			panic("impossible")
		}
	}, control.VERBOSE_PASS)
	//step4: codegen -- 代码生成
	if control.Optimization_Level <= 1 {
		control.Verbose("CodeGen", func() {
			//codegen_go.TransGo(Ast_go)
		}, control.VERBOSE_PASS)
	} else {

		////setp5:optimization -- 优化
		////Ast_c -> Ast_cfg
		//var Ast_cfg cfg.File
		//control.Verbose("TransCfg", func() {
		//	Ast_cfg = cfg.TransCfg(Ast_c)
		//}, control.VERBOSE_PASS)
		//if control.Visualize_format != control.None {
		//	cfg.Visualize(Ast_cfg)
		//}
		//Ast_cfg = cfg_opt.Opt(Ast_cfg)
		//util.Assert(Ast_cfg != nil, func() { panic("impossible") })
		//
		////CodegenCfg
		//control.Verbose("GenCfg", func() {
		//	cfg.CodegenCfg(Ast_cfg)
		//}, control.VERBOSE_PASS)

	}

}
