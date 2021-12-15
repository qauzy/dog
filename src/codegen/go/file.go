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
	"path"
	"strings"
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

	//处理全局变量
	if this.CurrentFile.ListFields() != nil {
		v := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.VAR,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		this.GolangFile.Decls = append(this.GolangFile.Decls, v)
		for _, f := range this.CurrentFile.ListFields() {
			v.Specs = append(v.Specs, this.transGlobalField(f))
		}
	}

	for _, c := range this.CurrentFile.ListClasses() {
		if c.IsEnum() {
			this.transEnum(c)
		} else {
			cl := this.transClass(c)
			this.GolangFile.Decls = append(this.GolangFile.Decls, cl)
		}

	}
}

// 带类型的变量声明
//
// param: fi
// return:
func (this *Translation) transGlobalField(fi ast.Field) (value *gast.ValueSpec) {
	this.CurrentField = fi
	if field, ok := fi.(*ast.FieldSingle); ok {
		//只处理成员变量
		var name = field.Name
		if field.IsField {
			name = Capitalize(field.Name)
		}
		value = &gast.ValueSpec{
			Doc:     nil,
			Names:   []*gast.Ident{gast.NewIdent(name)},
			Type:    this.transType(field.Tp),
			Values:  nil,
			Comment: nil,
		}
		if field.Value != nil {
			value.Values = append(value.Values, this.transExp(field.Value))
		}
	}
	return
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

func (this *Translation) WriteFile(base string, file string) (err error) {
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

	fileNameWithSuffix := path.Base(file)
	//获取文件的后缀(文件类型)
	fileType := path.Ext(fileNameWithSuffix)
	//获取文件名称(不带后缀)
	fileNameOnly := strings.TrimSuffix(fileNameWithSuffix, fileType)

	var suffix = strings.Replace(path.Dir(file), path.Dir(base), "", -1)
	//var suffix = path.Base(base)
	log.Warnf("------%v", suffix)
	var dir = "./out/" + path.Base(base) + "/" + suffix
	if !checkFileIsExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}

	var filename = dir + "/" + fileNameOnly + ".go"

	log.Warnf("写入:%v", filename)
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
