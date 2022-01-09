package main

import (
	"dog/ast"
	codegen_go "dog/codegen/golang"
	"dog/control"
	"dog/parser"
	"dog/storage"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func dog_Parser(filename string, buf []byte) ast.File {
	return parser.NewParse(filename, buf).Parser()
}

func main() {
	//tool.PrintAst()
	//return
	//log.SetLogDiscardLevel(log.Levelwarn)
	args := os.Args[1:len(os.Args)]
	control.CodeGen_base = control.Do_arg(args)
	if control.CodeGen_base == "" {
		control.Usage()
		os.Exit(0)
	}
	filepath.Walk(control.CodeGen_base, PareseJava)

}

func PareseJava(file string, info os.FileInfo, err error) error {
	//目录不处理
	if info.IsDir() {
		//log.Warnf("忽略目录:%v", file)
		return nil
	}

	//非java文件不处理
	suffix := path.Ext(file)
	if !strings.EqualFold(suffix, ".java") {
		//log.Warnf("忽略文件非java:%v --%v", file, suffix)
		return nil
	}

	log.Warnf("-------->>>>处理文件:%v", file)
	control.CodeGen_fileName = file
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Info(err)
		os.Exit(0)
	}
	if control.Lexer_test == true {
		lex := parser.NewLexer(file, buf)
		tk := lex.NextToken()
		var path string
		for tk.Kind != parser.TOKEN_EOF {
			if tk.Kind == parser.TOKEN_PACKAGE {
				tk = lex.NextToken()
				for tk.Kind != parser.TOKEN_SEMI {
					path += tk.Lexeme
					tk = lex.NextToken()
				}
			}

			var isClass bool
			tk = lex.NextToken()
			if tk.Kind == parser.TOKEN_PUBLIC || // public
				tk.Kind == parser.TOKEN_PROTECTED || // protected
				tk.Kind == parser.TOKEN_PRIVATE || // private
				tk.Kind == parser.TOKEN_ABSTRACT { // abstract
				isClass = true
				tk = lex.NextToken()
			}

			if (tk.Kind == parser.TOKEN_CLASS || tk.Kind == parser.TOKEN_ENUM || tk.Kind == parser.TOKEN_INTERFACE) && isClass {

				var kind int
				switch tk.Kind {
				case parser.TOKEN_CLASS:
					kind = 0
				case parser.TOKEN_ENUM:
					kind = 1
				case parser.TOKEN_INTERFACE:
					kind = 2
				}
				tk = lex.NextToken()

				pk := &storage.PackInfo{
					Project: "bitrade",
					Name:    tk.Lexeme,
					Path:    path,
					Kind:    kind,
				}
				storage.AddPack(pk)
				log.Info(tk.ToString(), path)
			}
		}
		return nil
	}
	var Ast ast.File
	//setp1: lexer&&parser
	control.Verbose("parser", func() {
		Ast = dog_Parser(file, buf)
	}, control.VERBOSE_PASS)

	//set3: trans -- 翻译
	var Ast_go *gast.File
	control.Verbose("Transaction", func() {
		switch control.CodeGen_codegen {
		case control.Go:
			Ast_go = codegen_go.TransGo(Ast, control.CodeGen_base, file)
		case control.C:
			//Ast_c = codegen_c.TransC(Ast)
		default:
			panic("impossible")
		}
	}, control.VERBOSE_PASS)
	return nil
}
