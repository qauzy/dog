package parser

import (
	"dog/util"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	"os"
)

type Lexer struct {
	fname   string
	s       string
	lineNum int
	buf     []byte
	/**
	 * the buf index, can increse or decrese to implements reset
	 */
	fp  int
	fpp int //用于测试
}

func NewLexer(fname string, buf []byte) *Lexer {
	initTokenMap()
	lex := Lexer{}
	lex.fname = fname
	lex.s = ""
	lex.lineNum = 1
	lex.buf = buf
	lex.fp = 0

	return &lex
}

func (this *Lexer) NextToken() *Token {
	var t *Token
	t = nil
	//拿不到token就继续读入字符
	for t == nil {
		t = this.nextTokenInternal()
	}
	return t
}

func (this *Lexer) expectKeyword(expect string) bool {
	reset := this.fp
	for _, e := range expect {
		if e == int32(this.buf[this.fp]) {
			this.fp++
			continue
		} else {
			this.fp = reset
			return false
		}
	}
	return true
}

func (this *Lexer) expectIdOrKey(c byte) *Token {
	kind, exist := tokenMap[this.s]
	if exist {
		tk := newToken(kind, this.s, this.lineNum)
		this.s = ""
		this.fp--
		return tk
	} else if this.s == "" {
		if c != ' ' {
			kk := tokenMap[string(c)]
			tk := newToken(kk, string(c), this.lineNum)
			return tk
		} else {
			return nil
		}
	} else {
		tk := newToken(TOKEN_ID, this.s, this.lineNum)
		this.s = ""
		this.fp--
		return tk
	}
}

// 忽略注解
//
// param: c
func (this *Lexer) lex_Annotation(c byte) {

	var ss string
	for c != '\n' && this.fp < len(this.buf) {
		c = this.buf[this.fp]
		this.fp++
		ss += string(c)
	}
	this.lineNum++
	fmt.Println("注解:", ss)
}
func (this *Lexer) lex_String(c byte) string {
	var ss string
	var st = this.fp

	c = this.buf[this.fp]
	for c != '\n' && c != '"' && this.fp < len(this.buf) {
		this.fp++
		c = this.buf[this.fp]

	}
	if c != '"' && this.fp >= len(this.buf) {
		util.ParserError("\"", "", this.lineNum)
	}
	var ed = this.fp
	//处理字符串末尾的"
	this.fp++
	ss = string(this.buf[st-1 : ed+1])
	fmt.Println("字符串:", ss)
	return ss
}

func (this *Lexer) lex_Comments(c byte) {
	ex := this.buf[this.fp]
	this.fp++
	if ex == '/' {
		for ex != '\n' && this.fp < len(this.buf) {
			ex = this.buf[this.fp]
			this.fp++
		}
		if this.fp == len(this.buf) {
			this.fp--
			return
		} else {
			this.lineNum++
		}
	} else if ex == '*' {
		ex = this.buf[this.fp]
		for (c != '*' || ex != '/') && this.fp < len(this.buf) {
			c = ex
			ex = this.buf[this.fp]
			this.fp++
			if ex == '\n' {
				this.lineNum++
			}
		}
		if this.fp == len(this.buf) {
			log.Info("error")
			os.Exit(0)
		}
	} else {
		log.Info("error")
		os.Exit(0)
	}
}

func (this *Lexer) lex_Num(c byte) string {
	var s string
	s += string(c)

	for {
		next := this.buf[this.fp]
		this.fp++
		var f = false
		if next >= '0' && next <= '9' || next == '.' && !f {
			if next == '.' {
				f = true
			}
			s += string(next)
			continue
		}

		//999abc is not number
		if (next == '_') || (next >= 'a' && next <= 'z') ||
			(next >= 'A' && next <= 'Z' && next != 'L') {
			fmt.Println("ilegal number")
			os.Exit(0)
		}

		if next == 'L' {
			this.fp++
		}
		break
	}

	this.fp--
	return s

}

func (this *Lexer) nextTokenInternal() *Token {
	if this.fp == len(this.buf) {
		return newToken(TOKEN_EOF, "EOF", this.lineNum)
	}

	c := this.buf[this.fp]
	this.fp++

	//换行处理
	for c == '\t' || '\n' == c || '\r' == c {
		if c == '\n' {
			this.lineNum++
		}
		if this.fp >= len(this.buf) {
			return newToken(TOKEN_EOF, "EOF", this.lineNum)
		}
		c = this.buf[this.fp]
		this.fp++
	}

	//文档末尾
	if this.fp >= len(this.buf) {
		return newToken(TOKEN_EOF, "EOF", this.lineNum)
	}
	//fallthrough强制执行后面的case代码
	switch c {
	case '&':
		if this.s == "" {
			if this.expectKeyword("&") {
				return newToken(TOKEN_AND, "&&", this.lineNum)
			} else {
				panic("expect &&")
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '|':
		if this.s == "" {
			if this.expectKeyword("|") {
				return newToken(TOKEN_OR, "||", this.lineNum)
			} else {
				panic("expect ||")
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '@':
		return newToken(TOKEN_AT, "@", this.lineNum)
		//this.lex_Annotation(c)
	case '+':
		if this.s == "" {
			if this.expectKeyword("+") {
				return newToken(TOKEN_AUTOADD, "++", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '=':
		if this.s == "" {
			if this.expectKeyword("=") {
				return newToken(TOKEN_EQ, "==", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '!':
		if this.s == "" {
			if this.expectKeyword("=") {
				return newToken(TOKEN_NE, "!=", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '-':
		if this.s == "" {
			if this.expectKeyword("-") {
				return newToken(TOKEN_AUTOSUB, "--", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}

	case '<':
		if this.s == "" {
			if this.expectKeyword("=") {
				return newToken(TOKEN_LE, "<=", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '>':
		if this.s == "" {
			if this.expectKeyword("=") {
				return newToken(TOKEN_GE, ">=", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case ' ':
		fallthrough
	case '?':
		fallthrough
	case ',':
		fallthrough
	case '.':
		fallthrough
	case '{':
		fallthrough
	case '[':
		fallthrough
	case '(':
		fallthrough

	case ':':
		fallthrough
	case '}':
		fallthrough
	case ']':
		fallthrough
	case ')':
		fallthrough
	case ';':
		fallthrough

	case '*':
		return this.expectIdOrKey(c)
	case '0':
		fallthrough
	case '1':
		fallthrough
	case '2':
		fallthrough
	case '3':
		fallthrough
	case '4':
		fallthrough
	case '5':
		fallthrough
	case '6':
		fallthrough
	case '7':
		fallthrough
	case '8':
		fallthrough
	case '9':
		if this.s == "" {
			return newToken(TOKEN_NUM, this.lex_Num(c), this.lineNum)
		}
		this.s += string(c)
	case '/':
		this.lex_Comments(c)
		//字符串
	case '"':
		if this.s == "" {
			return newToken(TOKEN_ID, this.lex_String(c), this.lineNum)
		}

	default:
		this.s += string(c)
	}

	return nil

}
