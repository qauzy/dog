package ast

//
//import (
//	"fmt"
//	"strconv"
//)
//
//type PrettyPrintVisitor struct {
//	indentLevel int
//}
//
//func NewPP() *PrettyPrintVisitor {
//	pp := new(PrettyPrintVisitor)
//	pp.indentLevel = 2
//
//	return pp
//}
//
//func (this *PrettyPrintVisitor) indent() {
//	this.indentLevel += 2
//}
//
//func (this *PrettyPrintVisitor) unIndent() {
//	this.indentLevel -= 2
//	if this.indentLevel <= 0 {
//		panic("indent error")
//	}
//}
//
//func (this *PrettyPrintVisitor) sayln(s string) {
//	fmt.Println(s)
//}
//func (this *PrettyPrintVisitor) say(s string) {
//	fmt.Print(s)
//}
//
//func (this *PrettyPrintVisitor) printSpeaces() {
//	i := this.indentLevel
//	for ; i != 0; i-- {
//		fmt.Printf(" ")
//	}
//}
//
//func (this *PrettyPrintVisitor) visitDec(e Field) {
//	switch v := e.(type) {
//	case *FieldSingle:
//		v.Tp.accept(this)
//		this.say(" " + v.Names)
//	default:
//		panic("wrong type")
//	}
//}
//
//func (this *PrettyPrintVisitor) visitMethod(e Y) {
//	switch v := e.(type) {
//	case *MethodSingle:
//		this.indent()
//		this.say("  public ")
//		v.RetType.accept(this)
//		this.say(" " + v.Names + "(")
//		for i, dec := range v.Formals {
//			switch d := dec.(type) {
//			case *FieldSingle:
//				if i != 0 {
//					this.say(",")
//				}
//				d.accept(this)
//			default:
//				panic("wrong type")
//			}
//		}
//		this.sayln(")")
//		this.sayln("  {")
//		for _, dec := range v.Locals {
//			switch d := dec.(type) {
//			case *FieldSingle:
//				this.printSpeaces()
//				d.Tp.accept(this)
//				this.sayln(" " + d.Names + ";")
//			default:
//				panic("wrong type")
//			}
//		}
//		this.say("\n")
//		for _, stm := range v.Stms {
//			stm.accept(this)
//		}
//		this.say("    return ")
//		v.RetExp.accept(this)
//		this.sayln(";")
//		this.sayln("  }")
//		this.unIndent()
//	default:
//		panic("wrong type")
//	}
//}
//
//func (this *PrettyPrintVisitor) visitClass(e Class) {
//	switch v := e.(type) {
//	case *ClassSingle:
//		this.say("class " + v.Names)
//		if v.Extends != "" {
//			this.sayln(" extends " + v.Extends)
//		} else {
//			this.sayln("")
//		}
//		this.sayln("{")
//		for _, dec := range v.Fields {
//			switch d := dec.(type) {
//			case *FieldSingle:
//				this.say("  ")
//				d.Tp.accept(this)
//				this.say(" ")
//				this.sayln(d.Names + ";")
//			default:
//				panic("wrong type")
//			}
//		}
//		for _, mtd := range v.Methods {
//			switch m := mtd.(type) {
//			case *MethodSingle:
//				m.accept(this)
//			default:
//				panic("wrong type")
//			}
//		}
//		this.sayln("}")
//	default:
//		panic("wrong type")
//	}
//}
//
//func (this *PrettyPrintVisitor) visitMain(e MainClass) {
//	switch v := e.(type) {
//	case *MainClassSingle:
//		this.sayln("class " + v.Names)
//		this.sayln("{")
//		this.sayln("  public static void main (String [] " + v.Args + ")")
//		this.sayln("  {")
//		this.indent()
//		v.Stms.accept(this)
//		this.unIndent()
//		this.sayln("  }")
//		this.sayln("}")
//	default:
//		panic("wrong type")
//	}
//}
//
//func (this *PrettyPrintVisitor) visitProg(e File) {
//	switch v := e.(type) {
//	case *FileSingle:
//		v.Mainclass.accept(this)
//		for _, vv := range v.Classes {
//			vv.accept(this)
//		}
//		fmt.Println("\n\n")
//	default:
//		panic("need FileSingle")
//	}
//}
//func (this *PrettyPrintVisitor) visitExp(e Exp) {
//	this.say("(")
//	switch v := e.(type) {
//	case *Add:
//		v.Left.accept(this)
//		this.say(" + ")
//		v.Right.accept(this)
//	case *LAnd:
//		v.Left.accept(this)
//		this.say(" && ")
//		v.Right.accept(this)
//	case *Times:
//		v.Left.accept(this)
//		this.say(" * ")
//		v.Right.accept(this)
//	case *IndexExpr:
//		v.X.accept(this)
//		this.say("[")
//		v.Index.accept(this)
//		this.say("]")
//	case *Call:
//		v.Callee.accept(this)
//		this.say("." + v.MethodName + "(")
//		for s, arg := range v.ArgsList {
//			arg.accept(this)
//			if s != len(v.ArgsList)-1 {
//				this.say(",")
//			}
//		}
//		this.say(")")
//	case *False:
//		this.say("false")
//	case *True:
//		this.say("true")
//	case *DefExpr:
//		this.say(v.Names)
//	case *Length:
//		v.X.accept(this)
//		this.say(".length")
//	case *Lt:
//		v.Left.accept(this)
//		this.say(" < ")
//		v.Right.accept(this)
//	case *NewIntArray:
//		this.say("new int[ ")
//		v.Size.accept(this)
//		this.say("]")
//	case *NewObject:
//		this.say("new " + v.T.String() + "()")
//	case *Not:
//		this.say("!")
//		v.Values.accept(this)
//	case *Num:
//		this.say(strconv.Itoa(v.Values))
//	case *Sub:
//		v.Left.accept(this)
//		this.say(" - ")
//		v.Right.accept(this)
//	case *This:
//		this.say("this")
//	default:
//		fmt.Printf("%T\n", v)
//		panic("wrong type")
//
//	}
//	this.say(")")
//}
//func (this *PrettyPrintVisitor) visitStm(e Stm) {
//	switch v := e.(type) {
//	case *Assign:
//		this.printSpeaces()
//		this.say(v.Names)
//		this.say(" = ")
//		v.Values.accept(this)
//		this.sayln(";")
//	case *AssignArray:
//		this.printSpeaces()
//		this.say(v.Names + "[")
//		v.Index.accept(this)
//		this.say("] = ")
//		v.Values.accept(this)
//		this.sayln(";")
//	case *Block:
//		this.printSpeaces()
//		this.sayln("{")
//		this.indent()
//		for _, stm := range v.Stms {
//			stm.accept(this)
//		}
//		this.unIndent()
//		this.printSpeaces()
//		this.sayln("}")
//	case *If:
//		this.printSpeaces()
//		this.say("if (")
//		v.Condition.accept(this)
//		this.sayln(") {")
//		this.indent()
//		v.Body.accept(this)
//		this.unIndent()
//		this.printSpeaces()
//		this.sayln("}else {")
//		this.indent()
//		v.Elsee.accept(this)
//		this.unIndent()
//		this.printSpeaces()
//		this.sayln("}")
//	case *Print:
//		this.printSpeaces()
//		this.say("System.out.println(")
//		v.Values.accept(this)
//		this.say(")")
//		this.sayln(";")
//	case *While:
//		this.printSpeaces()
//		this.say("while (")
//		v.Values.accept(this)
//		this.sayln(")")
//		this.indent()
//		v.Body.accept(this)
//		this.unIndent()
//	}
//}
//
//func (this *PrettyPrintVisitor) visitType(e Type) {
//	switch v := e.(type) {
//	case *Int:
//		this.say("int")
//	case *Boolean:
//		this.say("boolean")
//	case *IntArray:
//		this.say("int[]")
//	case *ClassType:
//		this.say(v.Names)
//	}
//}
//
//func (this *PrettyPrintVisitor) visit(e Acceptable) {
//	switch v := e.(type) {
//	case Field:
//		this.visitDec(v)
//	case Class:
//		this.visitClass(v)
//	case MainClass:
//		this.visitMain(v)
//	case File:
//		this.visitProg(v)
//	case Exp:
//		this.visitExp(v)
//	case Stm:
//		this.visitStm(v)
//	case Type:
//		this.visitType(v)
//	case Y:
//		this.visitMethod(v)
//	default:
//		panic("wrong type")
//	}
//}
//
//func (this *PrettyPrintVisitor) DumpProg(e Acceptable) {
//	this.visit(e)
//}
