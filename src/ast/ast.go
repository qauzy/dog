package ast

/*--------------------interface------------------*/
type Class interface {
	accept(v Visitor)
	_class()
}

type Dec interface {
	accept(v Visitor)
	GetDecType() int
	String() string
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
}

type Program interface {
	accept(v Visitor)
	_prog()
}

type Stm interface {
	accept(v Visitor)
	_stm()
}

type Type interface {
	accept(v Visitor)
	Gettype() int
	String() string
}

/*------------------ struct -----------------------*/

/*Dec*/ /*{{{*/
type DecSingle struct {
	Access  int
	Tp      Type
	Name    string
	IsField bool
	Stms    Stm //处理声明变量时的初始化语句
}

func (this *DecSingle) accept(v Visitor) {
	v.visit(this)
}

func (this *DecSingle) GetDecType() int {
	return this.Tp.Gettype()
}
func (this *DecSingle) String() string {
	s := this.Name + " " + this.Tp.String()
	return s
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
	Access  int
	Name    string
	Extends string
	Decs    []Dec
	Methods []Method
}

func (this *ClassSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *ClassSingle) _class() {
}

/*}}}*/

//Method  /*{{{*/
type MethodSingle struct {
	RetType Type
	Name    string // the name of whitch class belong to
	Formals []Dec
	Locals  []Dec
	Stms    []Stm
	RetExp  Exp
}

func (this *MethodSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *MethodSingle) _method() {
}

/*}}}*/

/*Prog*/ /*{{{*/
type ProgramSingle struct {
	Mainclass MainClass
	Classes   []Class
}

func (this *ProgramSingle) accept(v Visitor) {
	v.visit(this)
}
func (this *ProgramSingle) _prog() {
} /*}}}*/

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
//点调用
//Exp.Dot /*{{{*/
type Dot struct {
	Callee     Exp //new Sub().MethodName(ArgsList)
	MethodName string
	ArgsList   []Exp
	Firsttype  string
	ArgsType   []Type
	Rt         Type
	Exp_T
}

func Dot_new(callee Exp, m string, args []Exp,
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

func (this *Dot) accept(v Visitor) {
	v.visit(this)
}
func (this *Dot) _exp() {
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

//Exp.Id /*{{{*/
type Id struct {
	Name      string
	Tp        Type
	IsField   bool
	Statement bool //指示是否同时声明
	Exp_T
}

func Id_Ex_new(name string, tp Type, isField bool, Statement bool, line int) *Id {
	e := new(Id)
	e.Name = name
	e.Tp = tp
	e.IsField = isField
	e.Statement = Statement
	e.LineNum = line
	return e
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

//Exp.NewObject /*{{{*/
type NewObject struct {
	Name     string
	ArgsList []Exp //带初值的初始化
	Exp_T
}

func NewObjectWithArgsList_new(name string, ArgsList []Exp, line int) *NewObject {
	e := new(NewObject)
	e.Name = name
	e.LineNum = line
	e.ArgsList = ArgsList
	return e
}

func NewObject_new(name string, line int) *NewObject {
	e := new(NewObject)
	e.Name = name
	e.LineNum = line
	return e
}

func (this *NewObject) accept(v Visitor) {
	v.visit(this)
}
func (this *NewObject) _exp() {
}

type NewHash struct {
	Key string
	Ele string
	Exp_T
}

func NewHash_new(key string, ele string, line int) *NewHash {
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

type NewList struct {
	Ele string
	Exp_T
}

func NewList_new(Ele string, line int) *NewList {
	e := new(NewList)
	e.Ele = Ele
	e.LineNum = line
	return e
}

func (this *NewList) accept(v Visitor) {
	v.visit(this)
}
func (this *NewList) _exp() {
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
	LineNum int
}

//Stm.Assign    /*{{{*/
type Assign struct {
	Name    string
	Left    Exp //左边可能是一个包含声明语句的
	E       Exp
	Tp      Type
	IsField bool
	Stm_T
}

func Assign_new(name string, exp Exp, tp Type, isField bool, line int) *Assign {
	s := new(Assign)
	s.Name = name
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

//Stm.If    /*{{{*/
type If struct {
	Stm_T
	Condition Exp
	Thenn     Stm
	Elsee     Stm
}

func If_new(cond Exp, then Stm, elsee Stm, line int) *If {
	s := new(If)
	s.Condition = cond
	s.Thenn = then
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
	Stm_T
	E    Exp
	Body Stm
}

func For_new(exp Exp, body Stm, line int) *For {
	s := new(For)
	s.E = exp
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

//Stm end/*}}}*/

//Type   /*{{{*/

//Type.Int  /*{{{*/
const (
	TYPE_INT = iota
	TYPE_BOOLEAN
	TYPE_INTARRAY
	TYPE_CLASS
	TOKEN_STRING
	TYPE_LIST
	TYPE_MAP
)

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

//Type end/*}}}*/

//泛型
type ListType struct {
	Name     string
	Ele      string
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
	Key      string
	Ele      string
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
