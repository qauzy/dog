package codegen_go

import (
	"bytes"
	"dog/ast"
	"dog/cfg"
	"dog/optimize/golang"
	"dog/util"
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
	file          string
	currentFile   ast.File
	currentClass  ast.Class
	currentMethod ast.Method
	currentField  ast.Field
	currentStm    ast.Stm
	currentExp    ast.Exp
	GolangFile    *gast.File
	PkgName       string
}

func NewTranslation(file string, j ast.File) (translation *Translation) {
	translation = &Translation{
		file:        file,
		currentFile: j,
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
	if this.currentFile.ListFields() != nil {
		v := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.VAR,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}

		for _, f := range this.currentFile.ListFields() {
			if f.GetName() == "serialVersionUID" {
				continue
			}
			v.Specs = append(v.Specs, this.transGlobalField(f))
		}

		if v.Specs != nil {
			this.GolangFile.Decls = append(this.GolangFile.Decls, v)
		}

	}

	for _, c := range this.currentFile.ListClasses() {
		if c.GetType() == ast.ENUM_TYPE {
			this.transEnum(c)
		} else if c.GetType() == ast.INTERFACE_TYPE {
			this.transInterface(c)
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
	this.currentField = fi
	if field, ok := fi.(*ast.FieldSingle); ok {
		name := this.transNameExp(field.Name)
		if cfg.Capitalize {
			name.Name = util.Capitalize(name.Name)
		}
		//只处理成员变量
		value = &gast.ValueSpec{
			Doc:     nil,
			Names:   []*gast.Ident{name},
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

	gast.Inspect(this.GolangFile, func(n gast.Node) bool {
		for _, rule := range golang.StandardRules {
			rule(n)
		}
		return true
	})

	header := ""
	buffer := bytes.NewBufferString(header)

	err = this.astToGo(buffer, this.GolangFile)
	if err != nil {
		return
	}

	fileNameWithSuffix := path.Base(file)
	//获取文件的后缀(文件类型)
	fileType := path.Ext(fileNameWithSuffix)
	//获取文件名称(不带后缀)
	fileNameOnly := strings.TrimSuffix(fileNameWithSuffix, fileType)

	//删掉末尾的 "/"
	base = strings.TrimSuffix(base, "/")

	var suffix = strings.Replace(path.Dir(file), path.Dir(base), "", -1)
	//var suffix = path.Base(base)
	var dir = cfg.TargetPath + suffix

	if cfg.OneFold {
		dir += "/" + this.PkgName
	}

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
