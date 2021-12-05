package codegen_go

import (
	"bytes"
	"dog/ast"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
)

type Translation struct {
	CurrentFile   ast.File
	CurrentClass  ast.Class
	CurrentMethod ast.Method
	CurrentField  ast.Field
	CurrentStm    ast.Stm
	CurrentExp    ast.Exp
	GolangFile    *gast.File
}

func NewTranslation(j ast.File) (translation *Translation) {
	translation = &Translation{
		CurrentFile: j,
		GolangFile: &gast.File{
			Doc:        nil,
			Package:    0,
			Name:       gast.NewIdent(j.GetName()),
			Decls:      nil,
			Scope:      nil,
			Imports:    nil,
			Unresolved: nil,
			Comments:   nil,
		},
	}
	return
}

func (this *Translation) ParseClasses() {
	for _, c := range this.CurrentFile.GetClasses() {
		cl := this.transClass(c)
		this.GolangFile.Decls = append(this.GolangFile.Decls, cl)
	}
}

func (this *Translation) astToGo(dst *bytes.Buffer, node interface{}) error {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			log.Info(err)
		}
	}

	addNewline()

	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		return err
	}

	addNewline()

	return nil
}

func (this *Translation) WriteFile() (err error) {
	header := ""
	buffer := bytes.NewBufferString(header)

	err = this.astToGo(buffer, this.GolangFile)
	if err != nil {
		return
	}
	//fset := token.NewFileSet()
	//gast.Print(fset, this.GolangFile)
	//gast.Inspect(f, func(n gast.Node) bool {
	//	// Called recursively.
	//	gast.Print(fset, n)
	//	return true
	//})
	//var filename = "D:\\code\\dog\\src\\codegen\\go\\example\\test.go"
	var filename = "/opt/google/code/dog-comp/src/codegen/go/example/test.go"
	var f *os.File
	/***************************** 第一种方式: 使用 io.WriteString 写入文件 ***********************************************/
	if checkFileIsExist(filename) { //如果文件存在
		f, err = os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0666) //打开文件
	} else {
		f, err = os.Create(filename) //创建文件
	}
	if err != nil {
		return
	}
	_, err = io.WriteString(f, buffer.String()) //写入文件(字符串)
	return
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
