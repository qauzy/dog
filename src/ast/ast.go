package ast

/*--------------------interface------------------*/
type Class interface {
	accept(v Visitor)
	_class()
	AddField(f Field)
	AddMethod(m Method)
	ListFields() []Field
	GetField(name string) (f Field)
	GetMethod(name string) (m Method)
	ListMethods() []Method
	GetName() string
}

type Field interface {
	accept(v Visitor)
	GetDecType() int
	String() string
	GetName() string
}

type Exp interface {
	accept(v Visitor)
	_exp()
}

type MainClass interface {
	accept(v Visitor)
	_mainclass()
}

type Method interface {
	accept(v Visitor)
	_method()
	GetName() string
	//AddField(f Field)
	GetFormal(name string) (f Field)
	IsConstruct() bool
}

type File interface {
	accept(v Visitor)
	_prog()
	GetName() string
	GetClasses() []Class
}

type Stm interface {
	IsTriple() bool
	accept(v Visitor)
	_stm()
}

type Type interface {
	accept(v Visitor)
	Gettype() int
	String() string
}

/*------------------ struct -----------------------*/

/*Field*/ /*{{{*/
type FieldSingle struct {
	Access  int
	Tp      Type
	Name    string
	IsField bool
	Stms    Stm //处理声明变量时的初始化语句
}

func (this *FieldSingle) accept(v Visitor) {
	v.visit(this)
}

func (this *FieldSingle) GetDecType() int {
	return this.Tp.Gettype()
}
func (this *FieldSingle) String() string {
	s := this.Name + " " + this.Tp.String()
	return s
}

func (this *FieldSingle) GetName() string {
	return this.Name
}

/*}}}*/

/* MainClass {{{*/
type MainClassSingle struct {
	Name string
	Args string
	Stms Stm
}

func (this *MainClassSingle) accept(v Visitor) {
	v.visit(this)
}

func (this *MainClassSingle) _mainclass() {
}

/*}}}*/

/* ClassSingle {{{*/
type ClassSingle struct {
	Access     int
	Name       string
	Extends    string
	Fields     []Field
	FieldsMap  map[string]Field
	Methods    []Method
	MethodsMap map[string]Method
}

func (this *ClassSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *ClassSingle) _class() {

}
func (this *ClassSingle) GetName() string {
	return this.Name
}

func (this *ClassSingle) AddField(f Field) {
	this.FieldsMap[f.GetName()] = f
	this.Fields = append(this.Fields, f)
}

func (this *ClassSingle) GetField(name string) (f Field) {
	f = this.FieldsMap[name]
	return
}

func (this *ClassSingle) ListFields() []Field {
	return this.Fields
}

func (this *ClassSingle) AddMethod(m Method) {
	this.MethodsMap[m.GetName()] = m
	this.Methods = append(this.Methods, m)
}
func (this *ClassSingle) GetMethod(name string) (m Method) {
	m = this.MethodsMap[name]
	return
}

func (this *ClassSingle) ListMethods() []Method {
	return this.Methods
}

func NewClassSingle(Access int, Name string, Extends string) (cl *ClassSingle) {
	cl = &ClassSingle{
		Access:     Access,
		Name:       Name,
		Extends:    Extends,
		Fields:     nil,
		FieldsMap:  make(map[string]Field),
		Methods:    nil,
		MethodsMap: make(map[string]Method),
	}
	return
}

/*}}}*/

//Method  /*{{{*/

func NewMethodSingle(RetType Type, Name string, Formals []Field, Stms []Stm, Construct bool) (f *MethodSingle) {

	FormalsMap := make(map[string]Field)
	for _, f := range Formals {
		FormalsMap[f.GetName()] = f
	}
	f = &MethodSingle{
		RetType:    RetType,
		Name:       Name,
		Formals:    Formals,
		FormalsMap: FormalsMap,
		Stms:       Stms,
		Construct:  Construct,
	}
	return
}

type MethodSingle struct {
	RetType    Type
	Name       string // the name of whitch class belong to
	Formals    []Field
	FormalsMap map[string]Field
	Locals     []Field
	Stms       []Stm
	Construct  bool
}

func (this *MethodSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *MethodSingle) _method() {
}

func (this *MethodSingle) GetName() string {
	return this.Name
}

func (this *MethodSingle) IsConstruct() bool {
	return this.Construct
}
func (this *MethodSingle) GetFormal(name string) (f Field) {
	f = this.FormalsMap[name]
	return
}

/*}}}*/

/*Prog*/ /*{{{*/
type FileSingle struct {
	Name      string // identifier name
	Mainclass MainClass
	Classes   []Class
}

func (this *FileSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *FileSingle) _prog() {
}

func (this *FileSingle) GetClasses() []Class {
	return this.Classes
}
func (this *FileSingle) GetName() string {
	return this.Name
}

/*}}}*/

/*Exp*/ /*{{{*/

type Exp_T struct {
	LineNum int
}

//Exp.Add /*{{{*/
type Add struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Add_new(l Exp, r Exp, line int) *Add {
	e := new(Add)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *Add) accept(v Visitor) {
	v.visit(this)
}
func (this *Add) _exp() {
} /*}}}*/

//Exp.AutoAdd /*{{{*/
type AutoAdd struct {
	Left  Exp
	Right Exp
	Exp_T
}

func AutoAdd_new(l Exp, r Exp, line int) *AutoAdd {
	e := new(AutoAdd)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *AutoAdd) accept(v Visitor) {
	v.visit(this)
}
func (this *AutoAdd) _exp() {
} /*}}}*/

//Exp.AutoSub /*{{{*/
type AutoSub struct {
	Left  Exp
	Right Exp
	Exp_T
}

func AutoSub_new(l Exp, r Exp, line int) *AutoSub {
	e := new(AutoSub)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *AutoSub) accept(v Visitor) {
	v.visit(this)
}
func (this *AutoSub) _exp() {
} /*}}}*/

//Exp.Lambda /*{{{*/
type Lambda struct {
	Formals []Field
	Stms    []Stm
	RetExp  Exp
	Exp_T
}

func Lambda_new(Formals []Field, Stms []Stm, line int) *Lambda {
	n := new(Lambda)
	n.Formals = Formals
	n.Stms = Stms
	n.LineNum = line
	return n
}

func (this *Lambda) accept(v Visitor) {
	v.visit(this)
}
func (this *Lambda) _exp() {
} /*}}}*/

//Exp.Question /*{{{*/
type Question struct {
	E   Exp
	One Exp
	Two Exp
	Exp_T
}

func Question_new(l Exp, one Exp, two Exp, line int) *Question {
	n := new(Question)
	n.E = l
	n.One = one
	n.Two = two
	n.LineNum = line
	return n
}

func (this *Question) accept(v Visitor) {
	v.visit(this)
}
func (this *Question) _exp() {
} /*}}}*/

//Exp.Or /*{{{*/
type Or struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Or_new(l Exp, r Exp, line int) *Or {
	n := new(Or)
	n.Left = l
	n.Right = r
	n.LineNum = line
	return n
}

func (this *Or) accept(v Visitor) {
	v.visit(this)
}
func (this *Or) _exp() {
} /*}}}*/

//Exp.And /*{{{*/
type And struct {
	Left  Exp
	Right Exp
	Exp_T
}

func And_new(l Exp, r Exp, line int) *And {
	n := new(And)
	n.Left = l
	n.Right = r
	n.LineNum = line
	return n
}

func (this *And) accept(v Visitor) {
	v.visit(this)
}
func (this *And) _exp() {
} /*}}}*/

//Exp.Enum /*{{{*/
type Enum struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Enum_new(l Exp, r Exp, line int) *And {
	n := new(And)
	n.Left = l
	n.Right = r
	n.LineNum = line
	return n
}

func (this *Enum) accept(v Visitor) {
	v.visit(this)
}
func (this *Enum) _exp() {
} /*}}}*/

//Exp.Fcon /*{{{*/
type Fcon struct {
	Init      Stm
	Condition Exp
	Post      Exp
	Exp_T
}

func Fcon_new(Init Stm, Condition Exp, Post Exp, line int) *Fcon {
	n := new(Fcon)
	n.Init = Init
	n.Condition = Condition
	n.Post = Post
	n.LineNum = line
	return n
}

func (this *Fcon) accept(v Visitor) {
	v.visit(this)
}
func (this *Fcon) _exp() {
} /*}}}*/

//Exp.Time  /*{{{*/
type Times struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Times_new(l Exp, r Exp, line int) *Times {
	n := new(Times)
	n.Left = l
	n.Right = r
	n.LineNum = line
	return n
}

func (this *Times) accept(v Visitor) {
	v.visit(this)
}
func (this *Times) _exp() {
}

/*}}}*/

//Exp.ArraySelect /*{{{*/
type ArraySelect struct {
	Arrayref Exp
	Index    Exp
	Exp_T
}

func ArraySelect_new(array Exp, index Exp, line int) *ArraySelect {
	e := new(ArraySelect)
	e.Arrayref = array
	e.Index = index
	e.LineNum = line
	return e
}

func (this *ArraySelect) accept(v Visitor) {
	v.visit(this)
}
func (this *ArraySelect) _exp() {
}

/*}}}*/

//Exp.Call /*{{{*/
type Call struct {
	Callee     Exp //new Sub().MethodName(ArgsList)
	MethodName string
	ArgsList   []Exp
	Firsttype  string
	ArgsType   []Type
	Rt         Type
	Exp_T
}

func Call_new(callee Exp, m string, args []Exp,
	ftp string, argstype []Type,
	rt Type, line int) *Call {
	e := new(Call)
	e.Callee = callee
	e.MethodName = m
	e.ArgsList = args
	e.Firsttype = ftp
	e.ArgsType = argstype
	e.Rt = rt
	e.LineNum = line
	return e
}

func (this *Call) accept(v Visitor) {
	v.visit(this)
}
func (this *Call) _exp() {
}

/*}}}*/

//Exp.SelectorExpr /*{{{*/
type SelectorExpr struct {
	X   Exp
	Sel string
	Exp_T
}

func SelectorExpr_new(x Exp, sel string, line int) *SelectorExpr {
	e := new(SelectorExpr)
	e.X = x
	e.Sel = sel
	e.LineNum = line
	return e
}

func (this *SelectorExpr) accept(v Visitor) {
	v.visit(this)
}
func (this *SelectorExpr) _exp() {
}

/*}}}*/

//Exp.CallExpr /*{{{*/
type CallExpr struct {
	Callee   Exp //new Sub().MethodName(ArgsList)
	ArgsList []Exp
	Exp_T
}

func CallExpr_new(callee Exp, args []Exp, line int) *CallExpr {
	e := new(CallExpr)
	e.Callee = callee
	e.ArgsList = args
	e.LineNum = line
	return e
}

func (this *CallExpr) accept(v Visitor) {
	v.visit(this)
}
func (this *CallExpr) _exp() {
}

/*}}}*/

//Exp.False /*{{{*/
type False struct {
	Exp_T
}

func False_new(line int) *False {
	e := new(False)
	e.LineNum = line
	return e
}

func (this *False) accept(v Visitor) {
	v.visit(this)
}
func (this *False) _exp() {
}

/*}}}*/

//Exp.True   /*{{{*/
type True struct {
	Exp_T
}

func True_new(line int) *True {
	e := new(True)
	e.LineNum = line
	return e
}

func (this *True) accept(v Visitor) {
	v.visit(this)
}
func (this *True) _exp() {
}

/*}}}*/

//Exp.Null   /*{{{*/
type Null struct {
	Exp_T
}

func Null_new(line int) *Null {
	e := new(Null)
	e.LineNum = line
	return e
}

func (this *Null) accept(v Visitor) {
	v.visit(this)
}
func (this *Null) _exp() {
}

/*}}}*/

//Exp.Id /*{{{*/
type Id struct {
	Name      string
	Tp        Type
	IsField   bool
	Statement bool //指示是否同时声明
	Exp_T
}

func Id_new(name string, tp Type, isField bool, line int) *Id {
	e := new(Id)
	e.Name = name
	e.Tp = tp
	e.IsField = isField
	e.LineNum = line
	return e
}

func (this *Id) accept(v Visitor) {
	v.visit(this)
}
func (this *Id) _exp() {
}

/*}}}*/
//Exp.Class /*{{{*/
type ClassExp struct {
	Name Exp
	Exp_T
}

func ClassExp_new(Name Exp, line int) *ClassExp {
	e := new(ClassExp)
	e.Name = Name
	e.LineNum = line
	return e
}

func (this *ClassExp) accept(v Visitor) {
	v.visit(this)
}
func (this *ClassExp) _exp() {
}

//Exp.len /*{{{*/
//len(arrayref)
type Length struct {
	Arrayref Exp
	Exp_T
}

func Length_new(array Exp, line int) *Length {
	e := new(Length)
	e.Arrayref = array
	e.LineNum = line
	return e
}

func (this *Length) accept(v Visitor) {
	v.visit(this)
}
func (this *Length) _exp() {
}

/*}}}*/

//Exp.lt    /*{{{*/
// left < right
type Lt struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Lt_new(l Exp, r Exp, line int) *Lt {
	e := new(Lt)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *Lt) accept(v Visitor) {
	v.visit(this)
}
func (this *Lt) _exp() {
}

/*}}}*/

//Exp.le    /*{{{*/
// left <= right
type Le struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Le_new(l Exp, r Exp, line int) *Le {
	e := new(Le)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *Le) accept(v Visitor) {
	v.visit(this)
}
func (this *Le) _exp() {
}

/*}}}*/

//Exp.gt    /*{{{*/
// left > right
type Gt struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Gt_new(l Exp, r Exp, line int) *Gt {
	e := new(Gt)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *Gt) accept(v Visitor) {
	v.visit(this)
}
func (this *Gt) _exp() {
}

/*}}}*/

//Exp.ge    /*{{{*/
// left >= right
type Ge struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Ge_new(l Exp, r Exp, line int) *Ge {
	e := new(Ge)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *Ge) accept(v Visitor) {
	v.visit(this)
}
func (this *Ge) _exp() {
}

/*}}}*/

//Exp.eq    /*{{{*/
// left < right
type Eq struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Eq_new(l Exp, r Exp, line int) *Eq {
	e := new(Eq)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *Eq) accept(v Visitor) {
	v.visit(this)
}
func (this *Eq) _exp() {
}

/*}}}*/

//Exp.eq    /*{{{*/
// left < right
type Neq struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Neq_new(l Exp, r Exp, line int) *Neq {
	e := new(Neq)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *Neq) accept(v Visitor) {
	v.visit(this)
}
func (this *Neq) _exp() {
}

/*}}}*/

//Exp.NewObjectArray   /*{{{*/
type NewObjectArray struct {
	T    Type
	Eles Exp
	Size Exp
	Exp_T
}

func NewObjectArray_new(t Type, eles Exp, Size Exp, line int) *NewObjectArray {
	e := new(NewObjectArray)
	e.T = t
	e.Eles = eles
	e.Size = Size
	e.LineNum = line
	return e
}

func (this *NewObjectArray) accept(v Visitor) {
	v.visit(this)
}
func (this *NewObjectArray) _exp() {
}

/*}}}*/

//Exp.NewIntArray   /*{{{*/
type NewIntArray struct {
	Size Exp
	Exp_T
}

func NewIntArray_new(size Exp, line int) *NewIntArray {
	e := new(NewIntArray)
	e.Size = size
	e.LineNum = line
	return e
}

func (this *NewIntArray) accept(v Visitor) {
	v.visit(this)
}
func (this *NewIntArray) _exp() {
}

/*}}}*/

//Exp.NewStringArray   /*{{{*/
type NewStringArray struct {
	Size Exp
	Exp_T
}

func NewStringArray_new(size Exp, line int) *NewStringArray {
	e := new(NewStringArray)
	e.Size = size
	e.LineNum = line
	return e
}

func (this *NewStringArray) accept(v Visitor) {
	v.visit(this)
}
func (this *NewStringArray) _exp() {
}

/*}}}*/

//Exp.NewObject /*{{{*/
type NewObject struct {
	T        Type
	ArgsList []Exp //带初值的初始化
	Exp_T
}

func NewObjectWithArgsList_new(t Type, ArgsList []Exp, line int) *NewObject {
	e := new(NewObject)
	e.T = t
	e.LineNum = line
	e.ArgsList = ArgsList
	return e
}

func NewObject_new(t Type, line int) *NewObject {
	e := new(NewObject)
	e.T = t
	e.LineNum = line
	return e
}

func (this *NewObject) accept(v Visitor) {
	v.visit(this)
}
func (this *NewObject) _exp() {
}

/*}}}*/

//Exp.NewHash /*{{{*/
type NewHash struct {
	Key Type
	Ele Type
	Exp_T
}

func NewHash_new(key Type, ele Type, line int) *NewHash {
	e := new(NewHash)
	e.Key = key
	e.Ele = ele
	e.LineNum = line
	return e
}

func (this *NewHash) accept(v Visitor) {
	v.visit(this)
}
func (this *NewHash) _exp() {
}

/*}}}*/

//Exp.NewList /*{{{*/
type NewList struct {
	Ele      Type
	ArgsList []Exp //带初值的初始化
	Exp_T
}

func NewList_new(Ele Type, ArgsList []Exp, line int) *NewList {
	e := new(NewList)
	e.Ele = Ele
	e.ArgsList = ArgsList
	e.LineNum = line
	return e
}

func (this *NewList) accept(v Visitor) {
	v.visit(this)
}
func (this *NewList) _exp() {
}

/*}}}*/

//Exp.NewList /*{{{*/
type NewSet struct {
	Ele      Type
	ArgsList []Exp //带初值的初始化
	Exp_T
}

func NewSet_new(Ele Type, ArgsList []Exp, line int) *NewSet {
	e := new(NewSet)
	e.Ele = Ele
	e.ArgsList = ArgsList
	e.LineNum = line
	return e
}

func (this *NewSet) accept(v Visitor) {
	v.visit(this)
}
func (this *NewSet) _exp() {
}

/*}}}*/

//Exp.Not   /*{{{*/
// !expp
type Not struct {
	E Exp
	Exp_T
}

func Not_new(exp Exp, line int) *Not {
	e := new(Not)
	e.E = exp
	e.LineNum = line
	return e
}

func (this *Not) accept(v Visitor) {
	v.visit(this)
}
func (this *Not) _exp() {
}

/*}}}*/

//Exp.Num   /*{{{*/
type Num struct {
	Value int
	Exp_T
}

func Num_new(value int, line int) *Num {
	e := new(Num)
	e.Value = value
	e.LineNum = line
	return e
}

func (this *Num) accept(v Visitor) {
	v.visit(this)
}
func (this *Num) _exp() {
}

/*}}}*/

//Exp.Cast   /*{{{*/
type Cast struct {
	Tp    Type
	Right Exp
	Exp_T
}

func Cast_new(Tp Type, r Exp, line int) *Cast {
	e := new(Cast)
	e.Tp = Tp
	e.LineNum = line
	return e
}

func (this *Cast) accept(v Visitor) {
	v.visit(this)
}
func (this *Cast) _exp() {
}

/*}}}*/

//Exp.Sub   /*{{{*/
type Sub struct {
	Left  Exp
	Right Exp
	Exp_T
}

func Sub_new(l Exp, r Exp, line int) *Sub {
	e := new(Sub)
	e.Left = l
	e.Right = r
	e.LineNum = line
	return e
}

func (this *Sub) accept(v Visitor) {
	v.visit(this)
}
func (this *Sub) _exp() {
}

/*}}}*/

type This struct {
	Exp_T
}

func This_new(line int) *This {
	e := new(This)
	e.LineNum = line
	return e
}

func (this *This) accept(v Visitor) {
	v.visit(this)
}

func (this *This) _exp() {
}

//Exp end/*}}}*/

//Stm   /*{{{*/
type Stm_T struct {
	isTriple bool
	LineNum  int
}

func (this *Stm_T) IsTriple() bool {
	return this.isTriple
}
func (this *Stm_T) SetTriple() {
	this.isTriple = true
}

//Stm.Decl    /*{{{*/
type Decl struct {
	Name  string
	Tp    Type
	Value Exp
	Stm_T
}

func Decl_new(name string, tp Type, Value Exp, line int) *Decl {
	s := new(Decl)
	s.Name = name
	s.Tp = tp
	s.Value = Value
	s.LineNum = line
	return s
}

func (this *Decl) accept(v Visitor) {
	v.visit(this)
}
func (this *Decl) _stm() {
}

/*}}}*/

//Stm.Assign    /*{{{*/
type Assign struct {
	Left    Exp //左边可能是一个包含声明语句的
	E       Exp
	Tp      Type
	IsField bool
	Stm_T
}

func Assign_new(Left Exp, exp Exp, tp Type, isField bool, line int) *Assign {
	s := new(Assign)
	s.Left = Left
	s.E = exp
	s.Tp = tp
	s.IsField = isField
	s.LineNum = line
	return s
}

func (this *Assign) accept(v Visitor) {
	v.visit(this)
}
func (this *Assign) _stm() {
}

/*}}}*/

//Stm.ExprStm    /*{{{*/
type ExprStm struct {
	E Exp
	Stm_T
}

func ExprStm_new(exp Exp, line int) *ExprStm {
	s := new(ExprStm)
	s.E = exp
	s.LineNum = line
	return s
}
func (this *ExprStm) accept(v Visitor) {
	v.visit(this)
}
func (this *ExprStm) _stm() {
}

/*}}}*/

//Stm.Return    /*{{{*/
type Return struct {
	E Exp
	Stm_T
}

func Return_new(exp Exp, line int) *Return {
	s := new(Return)
	s.E = exp
	s.LineNum = line
	return s
}

func (this *Return) accept(v Visitor) {
	v.visit(this)
}
func (this *Return) _stm() {
}

/*}}}*/

//Stm.AssignArray   /*{{{*/
type AssignArray struct {
	// id[index] = e
	Name    string
	Index   Exp
	E       Exp
	Tp      Type
	IsField bool
	Stm_T
}

func AssignArray_new(name string,
	index Exp, exp Exp, tp Type,
	isField bool, line int) *AssignArray {
	s := new(AssignArray)
	s.Name = name
	s.Index = index
	s.E = exp
	s.Tp = tp
	s.IsField = isField
	s.LineNum = line
	return s
}

func (this *AssignArray) accept(v Visitor) {
	v.visit(this)
}
func (this *AssignArray) _stm() {
}

/*}}}*/

//Stm.Block /*{{{*/
type Block struct {
	Stms []Stm
	Stm_T
}

func Block_new(stms []Stm, line int) *Block {
	s := new(Block)
	s.Stms = stms
	s.LineNum = line
	return s
}

func (this *Block) accept(v Visitor) {
	v.visit(this)
}
func (this *Block) _stm() {
}

/*}}}*/

//Stm.Try    /*{{{*/
type Try struct {
	Stm_T
	Test      Stm
	Condition []Exp
	Catches   []Stm
	Finally   Stm
}

func Try_new(test Stm, cond []Exp, catches []Stm, Finally Stm, line int) *Try {
	s := new(Try)
	s.Condition = cond
	s.Test = test
	s.Catches = catches
	s.Finally = Finally
	s.LineNum = line
	return s
}

func (this *Try) accept(v Visitor) {
	v.visit(this)
}
func (this *Try) _stm() {
}

/*}}}*/

//Stm.If    /*{{{*/
type If struct {
	Stm_T
	Condition Exp
	Body      Stm
	Elsee     Stm
}

func If_new(cond Exp, Body Stm, elsee Stm, line int) *If {
	s := new(If)
	s.Condition = cond
	s.Body = Body
	s.Elsee = elsee
	s.LineNum = line
	return s
}

func (this *If) accept(v Visitor) {
	v.visit(this)
}
func (this *If) _stm() {
}

/*}}}*/

//Stm.Throw/*{{{*/
type Throw struct {
	Stm_T
	E Exp
}

func Throw_new(exp Exp, line int) *Throw {
	s := new(Throw)
	s.E = exp
	s.LineNum = line
	return s
}

func (this *Throw) accept(v Visitor) {
	v.visit(this)
}
func (this *Throw) _stm() {
}

/*}}}*/

//Stm.Print /*{{{*/
type Print struct {
	Stm_T
	E Exp
}

func Print_new(exp Exp, line int) *Print {
	s := new(Print)
	s.E = exp
	s.LineNum = line
	return s
}

func (this *Print) accept(v Visitor) {
	v.visit(this)
}
func (this *Print) _stm() {
}

/*}}}*/

//Stm.While   /*{{{*/
type While struct {
	Stm_T
	E    Exp
	Body Stm
}

func While_new(exp Exp, body Stm, line int) *While {
	s := new(While)
	s.E = exp
	s.Body = body
	s.LineNum = line
	return s
}

func (this *While) accept(v Visitor) {
	v.visit(this)
}
func (this *While) _stm() {
}

/*}}}*/

//Stm.For   /*{{{*/
type For struct {
	Init Stm
	Stm_T
	Cond Exp
	Post Exp
	Body Stm
}

func For_new(Init Stm, Condition Exp, Post Exp, body Stm, line int) *For {
	s := new(For)
	s.Init = Init
	s.Cond = Condition
	s.Post = Post
	s.Body = body
	s.LineNum = line
	return s
}

func (this *For) accept(v Visitor) {
	v.visit(this)
}
func (this *For) _stm() {
}

/*}}}*/

//Stm.Range   /*{{{*/
type Range struct {
	Value Exp
	Stm_T
	E    Exp
	Body Stm
}

func Range_new(Value Exp, exp Exp, body Stm, line int) *Range {
	s := new(Range)
	s.Value = Value
	s.E = exp
	s.Body = body
	s.LineNum = line
	return s
}

func (this *Range) accept(v Visitor) {
	v.visit(this)
}
func (this *Range) _stm() {
}

/*}}}*/

//Stm end/*}}}*/

//Type   /*{{{*/

//Type.Int  /*{{{*/
const (
	TYPE_INT = iota
	TYPE_Integer
	TYPE_BOOLEAN
	TYPE_VOID
	TYPE_INTARRAY
	TYPE_CLASS
	TYPE_STRING
	TYPE_STRINGARRAY
	TYPE_LIST
	TYPE_MAP
	TYPE_GENERIC
	TYPE_OBJECT
	TYPE_OBJECTARRAY
	TYPE_FUNCTION
)

type Function struct {
	TypeKind int
}

func (this *Function) accept(v Visitor) {
	v.visit(this)
}
func (this *Function) Gettype() int {
	return this.TypeKind
}
func (this *Function) String() string {
	return "@Function"
}

type Int struct {
	TypeKind int
}

func (this *Int) accept(v Visitor) {
	v.visit(this)
}
func (this *Int) Gettype() int {
	return this.TypeKind
}
func (this *Int) String() string {
	return "@int"
}

type Integer struct {
	TypeKind int
}

func (this *Integer) accept(v Visitor) {
	v.visit(this)
}
func (this *Integer) Gettype() int {
	return this.TypeKind
}
func (this *Integer) String() string {
	return "@Integer"
}

/*}}}*/

//Type.Bool /*{{{*/
type Boolean struct {
	TypeKind int
}

func (this *Boolean) accept(v Visitor) {
	v.visit(this)
}
func (this *Boolean) Gettype() int {
	return this.TypeKind
}

func (this *Boolean) String() string {
	return "@boolean"
}

//Type.Bool /*{{{*/

//Type.Void /*{{{*/
type Void struct {
	TypeKind int
}

func (this *Void) accept(v Visitor) {
	v.visit(this)
}
func (this *Void) Gettype() int {
	return this.TypeKind
}

func (this *Void) String() string {
	return "@void"
}

//Type.IntArray /*{{{*/
type IntArray struct {
	TypeKind int
}

func (this *IntArray) accept(v Visitor) {
	v.visit(this)
}
func (this *IntArray) Gettype() int {
	return this.TypeKind
}

func (this *IntArray) String() string {
	return "@int[]"
}

/*}}}*/

//Type.ObjectArray /*{{{*/
type ObjectArray struct {
	TypeKind int
}

func (this *ObjectArray) accept(v Visitor) {
	v.visit(this)
}
func (this *ObjectArray) Gettype() int {
	return this.TypeKind
}

func (this *ObjectArray) String() string {
	return "@Object[]"
}

/*}}}*/

//Type.String /*{{{*/

type String struct {
	TypeKind int
}

func (this *String) accept(v Visitor) {
	v.visit(this)
}
func (this *String) Gettype() int {
	return this.TypeKind
}

func (this *String) String() string {
	return "@String"
}

/*}}}*/

//Type.Void /*{{{*/

type StringArray struct {
	TypeKind int
}

func (this *StringArray) accept(v Visitor) {
	v.visit(this)
}
func (this *StringArray) Gettype() int {
	return this.TypeKind
}

func (this *StringArray) String() string {
	return "@String[]"
}

/*}}}*/

//Type.ClassType    /*{{{*/
type ClassType struct {
	Name     string
	TypeKind int
}

func (this *ClassType) accept(v Visitor) {
	v.visit(this)
}
func (this *ClassType) Gettype() int {
	return this.TypeKind
}

func (this *ClassType) String() string {
	return "@" + this.Name
}

/*}}}*/

//Type.GenericType    /*{{{*/
type GenericType struct {
	Name     string
	T        Type
	TypeKind int
}

func (this *GenericType) accept(v Visitor) {
	v.visit(this)
}
func (this *GenericType) Gettype() int {
	return this.TypeKind
}

func (this *GenericType) String() string {
	return this.Name + this.T.String()
}

/*}}}*/

//Type.ObjectType    /*{{{*/
type ObjectType struct {
	TypeKind int
}

func (this *ObjectType) accept(v Visitor) {
	v.visit(this)
}
func (this *ObjectType) Gettype() int {
	return this.TypeKind
}

func (this *ObjectType) String() string {
	return "@Object"
}

/*}}}*/

//Type end/*}}}*/

//泛型
type ListType struct {
	Name     string
	Ele      Type
	TypeKind int
}

func (this *ListType) accept(v Visitor) {
	v.visit(this)
}
func (this *ListType) Gettype() int {
	return this.TypeKind
}

func (this *ListType) String() string {
	return "@" + this.Name
}

type HashType struct {
	Name     string
	Key      Type
	Value    Type
	TypeKind int
}

func (this *HashType) accept(v Visitor) {
	v.visit(this)
}
func (this *HashType) Gettype() int {
	return this.TypeKind
}

func (this *HashType) String() string {
	return "@" + this.Name
}
