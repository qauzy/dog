package codegen_go

import (
	"dog/ast"
	"dog/util"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	gast "go/ast"
	"go/token"
	"path"
	"reflect"
	"strings"
)

func TransGo(p ast.File, base string, file string) (f *gast.File) {
	trans := NewTranslation(file, p)
	trans.ParseClasses()

	trans.WriteFile(base, file)

	return trans.GolangFile
}

// 带类型的变量声明
//
// param: fi
// return:
func (this *Translation) transField(fi ast.Field) (gfi *gast.Field) {
	this.CurrentField = fi
	if field, ok := fi.(*ast.FieldSingle); ok {
		//只处理成员变量
		var name = field.Name
		if field.IsField {
			name = Capitalize(field.Name)
		}
		gfi = &gast.Field{
			Doc:     nil,
			Names:   []*gast.Ident{gast.NewIdent(GetNewId(name))},
			Type:    this.transType(field.Tp),
			Tag:     nil,
			Comment: nil,
		}

		tag := &gast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("`gorm:\"column:%s\" json:\"%s\"`", SnakeString(name), field.Name),
		}
		gfi.Tag = tag
	}
	return
}

// 添加错误返回
//
// param: fi
// return:
func (this *Translation) getErrRet() (gfi *gast.Field) {

	gfi = &gast.Field{
		Doc:     nil,
		Names:   []*gast.Ident{gast.NewIdent("err")},
		Type:    gast.NewIdent("error"),
		Tag:     nil,
		Comment: nil,
	}

	return
}
func (this *Translation) TranslationBug(v interface{}) {
	var msg = fmt.Sprintf("未处理 [%v] %s\n", reflect.TypeOf(v).String(), path.Base(this.file))
	util.Bug(msg)
}

// 三元表达式
//
// param: v
// return:
func (this *Translation) transTriple(s ast.Stm) (stmts []gast.Stmt) {
	if !s.IsTriple() {
		panic("should triple expr")
	}
	switch v := s.(type) {
	//变量声明
	case *ast.DeclStmt:
		d := &gast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.VAR,
			Lparen: 0,
			Specs:  nil,
			Rparen: 0,
		}
		sp := &gast.ValueSpec{
			Doc:     nil,
			Names:   nil,
			Type:    this.transType(v.Tp),
			Values:  nil,
			Comment: nil,
		}

		for _, name := range v.Names {
			sp.Names = append(sp.Names, this.transNameExp(name))
		}

		//临时变量初值
		d.Specs = append(d.Specs, sp)
		stmts = append(stmts, &gast.DeclStmt{Decl: d})

		for idx, value := range v.Values {
			if vv, ok := value.(*ast.Question); ok {
				q := &gast.IfStmt{
					If:   0,
					Init: nil,
					Cond: this.transExp(vv.E),
					Body: &gast.BlockStmt{
						Lbrace: 0,
						List: []gast.Stmt{&gast.AssignStmt{
							Lhs:    []gast.Expr{this.transExp(v.Names[idx])},
							TokPos: 0,
							Tok:    token.ASSIGN,
							Rhs:    []gast.Expr{this.transExp(vv.One)}}},
					},

					Else: &gast.BlockStmt{
						Lbrace: 0,
						List: []gast.Stmt{&gast.AssignStmt{
							Lhs:    []gast.Expr{this.transExp(v.Names[idx])},
							TokPos: 0,
							Tok:    token.ASSIGN,
							Rhs:    []gast.Expr{this.transExp(vv.Two)}}},
						Rbrace: 0,
					},
				}

				stmts = append(stmts, q)
			} else {
				this.TranslationBug("必须是三目运算表达式")
			}

		}

		//赋值语句
	case *ast.Assign:
		if vv, ok := v.Value.(*ast.Question); ok {
			q := &gast.IfStmt{
				If:   0,
				Init: nil,
				Cond: this.transExp(vv.E),
				Body: &gast.BlockStmt{
					Lbrace: 0,
					List: []gast.Stmt{&gast.AssignStmt{
						Lhs:    []gast.Expr{this.transExp(v.Left)},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []gast.Expr{this.transExp(vv.One)}}},
				},

				Else: &gast.BlockStmt{
					Lbrace: 0,
					List: []gast.Stmt{&gast.AssignStmt{
						Lhs:    []gast.Expr{this.transExp(v.Left)},
						TokPos: 0,
						Tok:    token.ASSIGN,
						Rhs:    []gast.Expr{this.transExp(vv.Two)}}},
					Rbrace: 0,
				},
			}

			stmts = append(stmts, q)
		} else {
			panic("should triple expr")
		}
	//
	case *ast.ExprStm:
		//log.Debugf("三元表达式语句:%v", v)
		stmt := &gast.ExprStmt{X: this.transExp(v.E)}
		stmts = append(stmts, stmt)
	default:
		panic("bug")
		log.Debugf("transBlock-->%v -->%v", reflect.TypeOf(s).String(), v)

	}

	return

}

func (this *Translation) transType(t ast.Exp) (Type gast.Expr) {
	switch v := t.(type) {
	case *ast.SelectorExpr:
		return &gast.SelectorExpr{
			X:   this.transExp(v.X),
			Sel: gast.NewIdent(v.Sel),
		}
	case *ast.Id:
		if this.CurrentFile != nil && (this.CurrentFile.GetImport(v.Name) != nil) {
			pack := this.CurrentFile.GetImport(v.Name).GetPack()
			expr := &gast.SelectorExpr{
				X:   gast.NewIdent(pack),
				Sel: gast.NewIdent(v.Name),
			}

			return expr
		}
		return gast.NewIdent(v.Name)
	case *ast.Ident:
		if this.CurrentFile != nil && (this.CurrentFile.GetImport(v.Name) != nil) {
			pack := this.CurrentFile.GetImport(v.Name).GetPack()
			expr := &gast.SelectorExpr{
				X:   gast.NewIdent(pack),
				Sel: gast.NewIdent(v.Name),
			}

			return expr
		}
		return gast.NewIdent(v.Name)
	case *ast.Void:
		return nil
	case *ast.String:
		return gast.NewIdent("string")
	case *ast.StringArray:
		return gast.NewIdent("[]string")
	case *ast.Integer:
		return gast.NewIdent("int64")
	case *ast.Int:
		return gast.NewIdent("int")
	case *ast.IntArray:
		return gast.NewIdent("[]int")
	case *ast.HashType:
		return &gast.MapType{
			Map:   0,
			Key:   this.transType(v.Key),
			Value: this.transType(v.Value),
		}
	case *ast.ListType:
		return &gast.ArrayType{
			Lbrack: 0,
			Len:    nil,
			Elt:    this.transType(v.Ele),
		}
	case *ast.ClassType:
		if this.CurrentFile != nil && (this.CurrentFile.GetImport(v.Name) != nil) {
			pack := this.CurrentFile.GetImport(v.Name).GetPack()
			expr := &gast.SelectorExpr{
				X:   gast.NewIdent(pack),
				Sel: gast.NewIdent(v.Name),
			}

			return expr
		}
		return &gast.Ident{
			NamePos: 0,
			Name:    v.Name,
			Obj:     gast.NewObj(gast.Typ, v.Name),
		}
	case *ast.ObjectType:
		return &gast.InterfaceType{
			Interface: 0,
			Methods: &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 1,
			},
			Incomplete: false,
		}

	case *ast.ObjectArray:
		return &gast.ArrayType{
			Lbrack: 0,
			Len:    nil,
			Elt:    gast.NewIdent("interface{}"),
		}
	case *ast.Boolean:
		return gast.NewIdent("bool")
	case *ast.Byte:
		return gast.NewIdent("byte")
	case *ast.ByteArray:
		return gast.NewIdent("[]byte")
		//泛型
		//TODO 先用接口替代
	case *ast.GenericType:
		return &gast.InterfaceType{
			Interface: 0,
			Methods: &gast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 1,
			},
			Incomplete: false,
		}
	case *ast.Float:
		return gast.NewIdent("float64")
	case *ast.Date:
		return gast.NewIdent("time.Time")
	case *ast.ArrayType:
		return &gast.ArrayType{
			Lbrack: 0,
			Len:    nil,
			Elt:    this.transExp(v.Ele),
		}

	default:
		//return nil
		log.Info(v)
		panic("impossible")

	}
}

func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				//log.Info("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func DeCapitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 65 && vv[i] <= 90 { // 后文有介绍
				vv[i] += 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				//log.Info("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

/**
 * 驼峰转蛇形 snake string
 * @description XxYy to xx_yy , XxYY to xx_y_y
 * @date 2020/7/30
 * @param s 需要转换的字符串
 * @return string
 **/
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

/**
 * 蛇形转驼峰
 * @description xx_yy to XxYx  xx_y_y to XxYY
 * @date 2020/7/30
 * @param s要转换的字符串
 * @return string
 **/
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
