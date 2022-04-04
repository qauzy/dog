package codegen_go

import (
	"dog/ast"
	"dog/util"
	"fmt"
	gast "go/ast"
	"go/token"
	"strings"
)

//dao构造函数
func (this *Translation) getNewDaoFunc(c ast.Class) (fn *gast.FuncDecl) {
	var init gast.Stmt // 构造函数的初始化语句
	//处理函数参数
	params := &gast.FieldList{
		Opening: 0,
		List:    nil,
		Closing: 0,
	}
	params.List = append(params.List, this.getField(gast.NewIdent("db"), gast.NewIdent("*db.DB")))
	//处理返回值
	results := &gast.FieldList{
		Opening: 0,
		List:    nil,
		Closing: 0,
	}

	results.List = append(results.List, this.getField(gast.NewIdent("dao"), gast.NewIdent(c.GetName())))

	var body = &gast.BlockStmt{
		Lbrace: 0,
		List:   nil,
		Rbrace: 0,
	}

	//初始化语句
	val := &gast.UnaryExpr{
		OpPos: 0,
		Op:    token.AND,
		X: &gast.CompositeLit{
			Type:       gast.NewIdent(util.DeCapitalize(c.GetName())),
			Lbrace:     0,
			Elts:       []gast.Expr{gast.NewIdent("db")},
			Rbrace:     0,
			Incomplete: false,
		},
	}

	init = &gast.AssignStmt{
		Lhs:    []gast.Expr{gast.NewIdent("dao")},
		TokPos: 0,
		Tok:    token.ASSIGN,
		Rhs:    []gast.Expr{val},
	}

	body.List = append(body.List, init)

	retStm := &gast.ReturnStmt{
		Return:  0,
		Results: nil,
	}

	body.List = append(body.List, retStm)

	fn = &gast.FuncDecl{
		Doc:  nil,
		Recv: nil,
		Name: gast.NewIdent("New" + c.GetName()),
		Type: &gast.FuncType{
			Func:    0,
			Params:  params,
			Results: results,
		},
		Body: body,
	}
	return
}

func (this *Translation) getSaveDao(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *adminAccessLogDao) Save(m *entity.@) (result *entity.AdminAccessLog, err error) {
	err = this.DBWrite().Save(m).Error
	return
}
`
	src = strings.Replace(src, "@", strings.Replace(c.GetName(), "Dao", "", 1), 1)
	fn = this.getFunc(src)
	fn.Recv.List[0].Type = gast.NewIdent("*" + util.DeCapitalize(c.GetName()))
	fn.Type.Results.List[0].Type = gast.NewIdent("*entity." + strings.Replace(c.GetName(), "Dao", "", 1))
	return

	return
}

func (this *Translation) getFindByIdDao(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *adminDao) FindById(id int64) (result *entity.Admin,err error) {
	err = this.DBRead().Where("id = ?", id).First(&result).Error
	return
}
`
	fn = this.getFunc(src)
	fn.Recv.List[0].Type = gast.NewIdent("*" + util.DeCapitalize(c.GetName()))
	fn.Type.Results.List[0].Type = gast.NewIdent("*entity." + strings.Replace(c.GetName(), "Dao", "", 1))
	return
}

func (this *Translation) getDeleteByIdDao(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *dataDictionaryDao) DeleteById(id int64) (count int64, err error) {
	d := this.DBRead().Where("id = ?", id).Delete(@)
	err = d.Error
	count = d.RowsAffected
	return
}
`
	src = strings.Replace(src, "@", "entity."+strings.Replace(c.GetName(), "Dao", "", 1)+"{}", 1)
	fn = this.getFunc(src)
	fn.Recv.List[0].Type = gast.NewIdent("*" + util.DeCapitalize(c.GetName()))
	return
}

func (this *Translation) getFindAllDao(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func (this *adminAccessLogDao) FindAll(qp *types.QueryParam) (result []*entity.AdminAccessLog, err error) {
	d := this.DBRead()
	if qp != nil {
		d = qp.BuildQuery(d)
	}
	d = d.Find(&result)
	err = d.Error
	return
}
`
	fn = this.getFunc(src)
	fn.Recv.List[0].Type = gast.NewIdent("*" + util.DeCapitalize(c.GetName()))
	fn.Type.Results.List[0].Type = gast.NewIdent("arraylist.List[*entity." + strings.Replace(c.GetName(), "Dao", "", 1) + "]")
	return
}

// 构建New对象函数
//
// param: c
// return:
func (this *Translation) getNewService(c ast.Class) (fn *gast.FuncDecl) {
	src := `
func NewBusinessAuthApplyService(BusinessAuthApplyDao *dao.BusinessAuthApplyDao) (ret *BusinessAuthApplyService) {
	ret = new(@)
}
`
	src = strings.Replace(src, "@", c.GetName(), 1)

	fn = this.getFunc(src)
	fn.Name = gast.NewIdent("New" + c.GetName())
	fn.Type.Params.List = nil
	for _, fi := range c.ListFields() {
		param := this.transField(fi)
		//参数小写
		for _, v := range param.Names {
			v.Name = util.DeCapitalize(v.Name)
		}

		fn.Type.Params.List = append(fn.Type.Params.List, param)
		fn.Body.List = append(fn.Body.List, &gast.ExprStmt{gast.NewIdent(fmt.Sprintf("ret.%s = %s", util.Capitalize(fi.GetName()), util.DeCapitalize(fi.GetName())))})
	}
	fn.Body.List = append(fn.Body.List, &gast.ReturnStmt{})
	fn.Type.Results.List[0].Type = gast.NewIdent("*" + c.GetName())
	return
}
