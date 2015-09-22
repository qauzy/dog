package cfg

import (
	"../control"
	"fmt"
	"log"
	"os"
	"strconv"
)

func CodegenCfg(p Program) {
	var f_indentLevel int = 2
	var f_outputName string
	var f_fd *os.File
	/**
	 * Record new declaration generated by genVar
	 */
	var f_redec map[string]bool
	/**
	 * Record all classes with their field
	 */
	var f_class_dec map[string][]*Tuple

	var trans func(Acceptable)

	say := func(s string) {
		f_fd.WriteString(s)
	}

	sayln := func(s string) {
		say(s)
		f_fd.WriteString("\n")
	}

	indent := func() {
		f_indentLevel += 2
	}

	unIndent := func() {
		f_indentLevel -= 2
	}

	printSpeaces := func() {
		i := f_indentLevel
		for i != 0 {
			say(" ")
			i--
		}
	}

	getVar := func(v string) string {
		if f_redec[v] == true {
			return "frame." + v
		} else {
			return v
		}
	}

	outputGcMap := func(method Method) {
		var m *MethodSingle
		if v, ok := method.(*MethodSingle); ok {
			m = v
		} else {
			panic("impossible")
		}
		say("char * " + m.classId + "_" + m.name + "_arguments_gc_map = \"")
		for _, dec := range m.formals {
			if d, ok := dec.(*DecSingle); ok {
				if t := d.tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
					say("1")
				} else {
					say("0")
				}
			} else {
				panic("impossible")
			}
		}
		sayln("\";")
		//locals map
		i := 0
		for _, dec := range m.locals {
			if d, ok := dec.(*DecSingle); ok {
				if t := d.tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
					i++
				}
			} else {
				panic("impossible")
			}
		}
		sayln("int " + m.classId + "_" + m.name + "_locals_gc_map= " + strconv.Itoa(i) + ";")
		sayln("")
	}

	outputMainGcStack := func(mm MainMethod) {
		var m *MainMethodSingle
		if v, ok := mm.(*MainMethodSingle); ok {
			m = v
		} else {
			panic("impossible")
		}
		sayln("struct Tiger_main_gc_frame")
		sayln("{")
		indent()
		printSpeaces()
		sayln("void *prev_;")
		printSpeaces()
		sayln("char *arguments_gc_map;")
		printSpeaces()
		sayln("int *arguments_base_address;")
		printSpeaces()
		sayln("int locals_gc_map;")

		for _, dec := range m.locals {
			if d, ok := dec.(*DecSingle); ok {
				if t := d.tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
					printSpeaces()
					trans(d)
					sayln(";")
				}
			} else {
				panic("impossible")
			}
		}
		unIndent()
		sayln("};\n")
	}

	outputGcStack := func(mm Method) {
		var m *MethodSingle
		if v, ok := mm.(*MethodSingle); ok {
			m = v
		} else {
			panic("impossible")
		}
		sayln("struct " + m.classId + "_" + m.name + "_gc_frame")
		sayln("{")
		indent()
		printSpeaces()
		sayln("void *prev_;")
		printSpeaces()
		sayln("char *arguments_gc_map;")
		printSpeaces()
		sayln("int *arguments_base_address;")
		printSpeaces()
		sayln("int locals_gc_map;")
		for _, dec := range m.locals {
			if d, ok := dec.(*DecSingle); ok {
				if t := d.tp.GetType(); t == TYPE_INTARRAY || t == TYPE_CLASSTYPE {
					printSpeaces()
					trans(d)
					sayln(";")
				}
			} else {
				panic("impossible")
			}
		}
		unIndent()
		sayln("};\n")
	}

	outputVtable := func(v Vtable) {
		var vt *VtableSingle
		if vv, ok := v.(*VtableSingle); ok {
			vt = vv
		} else {
			panic("impossible")
		}
		sayln("struct " + vt.id + "_vtable " + vt.id + "_vtable_ =")
		sayln("{")
		//According to the f_class_dec, generate class
		//Gc map
		locals := f_class_dec[vt.id]
		printSpeaces()
		say("\"")
		for _, t := range locals {
			_, ok := t.tp.(*ClassType)
			_, ok2 := t.tp.(*ClassType)
			if ok || ok2 {
				say("1")
			} else {
				say("0")
			}
		}
		sayln("\",")
		for _, f := range vt.methods {
			say("  ")
			sayln(f.class_name + "_" + f.id + ",")
		}
		sayln("};\n")
	}

	trans_Type := func(tt Type) {
		switch t := tt.(type) {
		case *IntType:
			say("int")
		case *IntArrayType:
			say("int*")
		case *ClassType:
			say("struct " + t.id + "*")
		default:
			panic("impossible")
		}
	}

	trans_Dec := func(dd Dec) {
		switch d := dd.(type) {
		case *DecSingle:
			trans(d.tp)
			say(" ")
			say(d.id)
		default:
			panic("impossilbe")
		}
	}

	trans_Transfer := func(tt Transfer) {
		switch t := tt.(type) {
		case *Goto:
			printSpeaces()
			say("goto " + t.label.String() + ";\n")
		case *If:
			printSpeaces()
			say("if (")
			trans(t.cond)
			say(")\n")
			printSpeaces()
			say("  goto " + t.truee.String() + ";\n")
			printSpeaces()
			say("else\n")
			printSpeaces()
			say("  goto " + t.falsee.String() + ";\n")
		case *Return:
			printSpeaces()
			//XXX this is for the gc
			sayln("previous = frame.prev_;")
			printSpeaces()
			say("return ")
			trans(t.op)
			say(";\n")
		default:
			panic("impossilbe")
		}
	}

	trans_Block := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			sayln("//block start")
			say(b.label.String() + ":\n")
			for _, s := range b.stms {
				trans(s)
				say("\n")
			}
			trans(b.transfer)
			sayln("//block en")
		default:
			panic("impossible")
		}
	}

	trans_Operand := func(oo Operand) {
		switch o := oo.(type) {
		case *Int:
			say(strconv.Itoa(o.value))
		case *Var:
			if o.isField == false {
				if f_redec[o.id] == false {
					say(o.id)
				} else {
					say("frame." + o.id)
				}
			} else {
				say("this->" + o.id)
			}
		default:
			panic("impossible")
		}
	}

	trans_Stm := func(ss Stm) {
		switch s := ss.(type) {
		case *Add:
			printSpeaces()
			say(getVar(s.dst) + " = ")
			trans(s.left)
			say(" + ")
			trans(s.right)
			sayln(";")
		case *And:
			printSpeaces()
			say(getVar(s.dst) + " = ")
			trans(s.left)
			say(" && ")
			trans(s.right)
			sayln(";")
		case *AssignArray:
			printSpeaces()
			if s.isField == false {
				say(getVar(s.dst) + "[")
			} else {
				say("this->" + s.dst + "[")
			}
			trans(s.index)
			say("+4]=")
			trans(s.exp)
			sayln(";")
		case *ArraySelect:
			printSpeaces()
			say(s.id + " = ")
			trans(s.array)
			say("[")
			trans(s.index)
			say("+4]")
			sayln(";")
		case *InvokeVirtual:
			printSpeaces()
			say(getVar(s.dst) + " = " + getVar(s.obj))
			say("->vptr->" + s.f + "(" + getVar(s.obj))
			for _, x := range s.args {
				say(", ")
				trans(x)
			}
			say(");")
		case *Length:
			printSpeaces()
			say(getVar(s.dst) + " = ")
			trans(s.array)
			say("[2];\n")
		case *Lt:
			printSpeaces()
			say(getVar(s.dst) + " = ")
			trans(s.left)
			say(" < ")
			trans(s.right)
			say(";")
		case *Move:
			printSpeaces()
			if s.IsField == false {
				say(getVar(s.dst) + " = ")
			} else {
				say("this->" + s.dst + " = ")
			}
			trans(s.src)
			say(";")
		case *NewIntArray:
			printSpeaces()
			say(getVar(s.dst) + " = (int*)Tiger_new_array(")
			trans(s.exp)
			sayln(");")
		case *NewObject:
			printSpeaces()
			say(getVar(s.dst) +
				" = ((struct " + s.c +
				"*)(Tiger_new(&" + s.c +
				"_vtable_, sizeof(struct " + s.c +
				"))));")
		case *Not:
			printSpeaces()
			say(getVar(s.dst) + " = !(")
			trans(s.exp)
			sayln(");")
		case *Print:
			printSpeaces()
			say("System_out_println(")
			trans(s.arg)
			sayln(");")
		case *Sub:
			printSpeaces()
			say(getVar(s.dst) + " = ")
			trans(s.left)
			say(" - ")
			trans(s.right)
			say(";")
		case *Times:
			printSpeaces()
			say(getVar(s.dst) + " = ")
			trans(s.left)
			say(" * ")
			trans(s.right)
			say(";")
		default:
			fmt.Printf("%T\n", s)
			panic("impossilbe")
		}
	}

	trans_Vtable := func(vv Vtable) {
		switch v := vv.(type) {
		case *VtableSingle:
			sayln("struct " + v.id + "_vtable")
			sayln("{")
			printSpeaces()
			sayln("char* " + v.id + "_gc_map;")
			for _, f := range v.methods {
				say("  ")
				trans(f.ret_type)
				say(" (*" + f.id + ")(")
				for idx, dec := range f.args {
					if idx != 0 {
						say(", ")
					}
					trans(dec)
				}
				sayln(");")
			}

			sayln("};\n")

		default:
			panic("impossible")
		}
	}

	trans_Method := func(mm Method) {
		switch m := mm.(type) {
		case *MethodSingle:
			f_redec = make(map[string]bool)
			trans(m.ret_type)
			say(" " + m.classId + "_" + m.name + "(")
			for idx, dec := range m.formals {
				if idx != 0 {
					say(", ")
				}
				trans(dec)
			}
			sayln(")")

			sayln("{")
			printSpeaces()
			sayln("struct " + m.classId + "_" + m.name + "_gc_frame frame;")
			printSpeaces()
			sayln("frame.prev_ = previous;")
			printSpeaces()
			sayln("previous = &frame;")
			printSpeaces()
			sayln("frame.arguments_gc_map = " + m.classId + "_" + m.name + "_arguments_gc_map;")
			printSpeaces()
			sayln("frame.arguments_base_address = (int*)&this;")
			printSpeaces()
			sayln("frame.locals_gc_map = " + m.classId + "_" + m.name + "_locals_gc_map;")

			for _, dec := range m.locals {
				if d, ok := dec.(*DecSingle); ok {
					t := d.tp.GetType()
					printSpeaces()
					if t != TYPE_INTARRAY && t != TYPE_CLASSTYPE {
						trans(dec)
						sayln(";")
					} else {
						f_redec[d.id] = true
						sayln("frame." + d.id + "=0;")
					}
				} else {
					panic("impossible")
				}
			}

			sayln("")
			for _, b := range m.blocks {
				trans(b)
			}
			sayln("}")

		default:
			panic("impossible")
		}
	}

	trans_MainMethod := func(mm MainMethod) {
		switch m := mm.(type) {
		case *MainMethodSingle:
			f_redec = make(map[string]bool)
			sayln("int Tiger_main()")
			sayln("{")

			printSpeaces()
			sayln("struct Tiger_main_gc_frame frame;")
			printSpeaces()
			sayln("frame.prev_ = previous;")
			printSpeaces()
			sayln("previous = &frame;")
			printSpeaces()
			sayln("frame.arguments_gc_map = 0;")
			printSpeaces()
			sayln("frame.arguments_base_address = 0;")
			printSpeaces()
			sayln("frame.locals_gc_map = Tiger_main_locals_gc_map;")

			for _, dec := range m.locals {
				if d, ok := dec.(*DecSingle); ok {
					printSpeaces()
					t := d.tp.GetType()
					if t != TYPE_INTARRAY && t != TYPE_CLASSTYPE {
						trans(d)
						sayln(";")
					} else {
						f_redec[d.id] = true
						sayln("frame." + d.id + "=0;")
					}
				} else {
					panic("impossible")
				}
			}
			for _, b := range m.blocks {
				trans(b)
			}
			sayln("\n}\n")
		default:
			panic("impossible")
		}
	}

	trans_Class := func(cc Class) {
		switch c := cc.(type) {
		case *ClassSingle:
			locals := make([]*Tuple, 0)
			sayln("struct " + c.id)
			sayln("{")
			sayln("  struct " + c.id + "_vtable *vptr;")

			printSpeaces()
			sayln("int isObjOrArray;")
			printSpeaces()
			sayln("int length;")
			printSpeaces()
			sayln("void *forwarding;")

			for _, t := range c.decs {
				say("  ")
				trans(t.tp)
				say("  ")
				sayln(t.id + ";")
				locals = append(locals, t)
			}
			//store all field info for generate v-table
			f_class_dec[c.id] = locals
			sayln("};")
		default:
			panic("impossible")
		}
	}

	trans_Program := func(pp Program) {
		switch p := pp.(type) {
		case *ProgramSingle:
			sayln("// This is automatically generated by the Dog compiler.")
			sayln("// Do NOT modify!\n")
			sayln("extern void *previous;")
			sayln("extern void *Tiger_new_array(int);")
			sayln("extern void *Tiger_new(void*, int);")
			sayln("extern int System_out_println(int);")

			sayln("//structure")
			for _, c := range p.classes {
				trans(c)
			}
			sayln("//vtable structures")
			for _, v := range p.vtables {
				trans(v)
			}
			sayln("\n//method decls")
			for _, m := range p.methods {
				if mm, ok := m.(*MethodSingle); ok {
					trans(mm.ret_type)
					say(" " + mm.classId + "_" + mm.name + "(")
					for idx, d := range mm.formals {
						if idx != 0 {
							say(", ")
						}
						trans(d)
					}
					sayln(");")
				} else {
					panic("impossible")
				}
			}
			sayln("//vtable")
			for _, v := range p.vtables {
				outputVtable(v)
			}
			sayln("")
			sayln("//GC stack frames")
			outputMainGcStack(p.main_method)
			for _, method := range p.methods {
				outputGcStack(method)
			}
			sayln("// memory GC maps")
			sayln("int Tiger_main_locals_gc_map = 1;\n")
			for _, m := range p.methods {
				outputGcMap(m)
			}
			sayln("// methods")
			for _, m := range p.methods {
				trans(m)
			}
			sayln("")
			sayln("// main")
			trans(p.main_method)
			sayln("")
			say("\n\n")
		default:
			panic("impossible")
		}
	}

	trans = func(e Acceptable) {
		switch v := e.(type) {
		case Block:
			trans_Block(v)
		case Class:
			trans_Class(v)
		case Dec:
			trans_Dec(v)
		case MainMethod:
			trans_MainMethod(v)
		case Method:
			trans_Method(v)
		case Operand:
			trans_Operand(v)
		case Program:
			trans_Program(v)
		case Stm:
			trans_Stm(v)
		case Transfer:
			trans_Transfer(v)
		case Type:
			trans_Type(v)
		case Vtable:
			trans_Vtable(v)
		default:
			fmt.Printf("%T\n", v)
			panic("impossible")
		}
	}

	if control.CodeGen_outputName != "" {
		f_outputName = control.CodeGen_outputName
	} else if control.CodeGen_fileName != "" {
		f_outputName = control.CodeGen_fileName + ".c"
	} else {
		f_outputName = "a.c"
	}
	f_fd, err := os.Create(f_outputName)
	if err != nil {
		log.Fatal(err)
	}
	defer f_fd.Close()
	f_class_dec = make(map[string][]*Tuple)
	trans(p)
}