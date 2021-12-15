package main

import (
	"dog/ast"
	codegen_go "dog/codegen/go"
	"dog/control"
	"dog/parser"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
)

func dog_Parser(filename string, buf []byte) ast.File {
	return parser.NewParse(filename, buf).Parser()
}

func main() {
	//log.SetLogDiscardLevel(log.Levelwarn)
	args := os.Args[1:len(os.Args)]
	base := control.Do_arg(args)
	if base == "" {
		control.Usage()
		os.Exit(0)
	}
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		log.Warnf("-------->>>>处理文件:%v", path)
		control.CodeGen_fileName = path
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			log.Info(err)
			os.Exit(0)
		}
		if control.Lexer_test == true {
			lex := parser.NewLexer(path, buf)
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
			Ast = dog_Parser(path, buf)
		}, control.VERBOSE_PASS)

		//set3: trans -- 翻译
		var Ast_go *gast.File
		control.Verbose("Transaction", func() {
			switch control.CodeGen_codegen {
			case control.Go:
				Ast_go = codegen_go.TransGo(Ast, base, path)
			case control.C:
				//Ast_c = codegen_c.TransC(Ast)
			default:
				panic("impossible")
			}
		}, control.VERBOSE_PASS)
		return nil
	})

}
