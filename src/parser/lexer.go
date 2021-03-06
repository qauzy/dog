package parser

import (
	"dog/util"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	"os"
	"path"
)

type Lexer struct {
	fname   string
	s       string
	lineNum int
	buf     []byte
	/**
	 * the buf index, can increse or decrese to implements reset
	 */
	fp int
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

func (this *Lexer) ExpectKeyword(expect string) bool {
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

	for c != '\n' && (c != '"' || (c == '"' && this.buf[this.fp-1] == '\\')) && this.fp < len(this.buf) {
		this.fp++
		c = this.buf[this.fp]
	}
	if c != '"' && this.fp >= len(this.buf) {
		util.ParserError("\"", string(c), this.lineNum, this.fname)
	}
	var ed = this.fp
	//处理字符串末尾的"
	this.fp++
	ss = string(this.buf[st-1 : ed+1])
	//log.Debugf("字符串:%v", ss)
	return ss
}

func (this *Lexer) lex_Comments(c byte) (comment string) {

	ex := this.buf[this.fp]
	var st = this.fp - 1
	var ed = 0
	//行注释
	if ex == '/' {
		for ex != '\n' && this.fp+1 < len(this.buf) {
			this.fp++
			ex = this.buf[this.fp]
		}
		if this.fp == len(this.buf) {
			this.fp--

			//换行符
		} else {
			//this.fp++
			//this.lineNum++
		}
		//注释块
	} else if ex == '*' {
		comment = `/*`
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
			log.Errorf("end of file,line:%v", this.lineNum)
			return
		}
	}
	ed = this.fp
	if ed > st+1 {
		comment = string(this.buf[st:ed])
	}
	return
}

func (this *Lexer) lex_Num(c byte) string {
	var s string
	s += string(c)
	var f = false
	var hex = false
	for {
		next := this.buf[this.fp]
		this.fp++

		//16进制
		if s == "0" && (next == 'x' || next == 'X') {
			hex = true
			s += string(next)
			continue
		}

		if next >= '0' && next <= '9' || hex && ((next >= 'a' && next <= 'f') || (next >= 'A' && next <= 'F')) || next == '.' && !f {
			if next == '.' {
				f = true
			}
			s += string(next)
			continue
		}

		//999abc is not number
		if (c < '0' && c > '9' && next == '_') || (next >= 'a' && next <= 'z' && (next != 'f') && (next != 'd') && (next != 'l')) ||
			(next >= 'A' && next <= 'Z' && next != 'L' && (next != 'D') && (next != 'F')) {
			log.Errorf("ilegal number,line:%v, %v, %v", this.lineNum, string(next), s)
			os.Exit(0)
		}

		if next == 'L' || next == 'l' || next == 'd' || next == 'D' || next == '_' || (next == 'f' || next == 'F') && !hex {
			this.fp++
		}
		break
	}

	this.fp--
	//log.Warnf("数字:%v", s)
	return s

}

func (this *Lexer) nextTokenInternal() *Token {
	if this.fp >= len(this.buf) {
		return newToken(TOKEN_EOF, "EOF", this.lineNum)
	}

	c := this.buf[this.fp]
	this.fp++

	//换行处理
	for c == '\t' || '\n' == c || '\r' == c {
		if c == '\n' {
			if this.s != "" {
				if (this.s[0] >= 'a' && this.s[0] <= 'z') || (this.s[0] >= 'A' && this.s[0] <= 'Z') || this.s[0] == '_' {
					return this.expectIdOrKey(c)
				}
			}
			this.lineNum++

		} else if c == '\t' || c == ' ' {
			if this.s != "" {
				if (this.s[0] >= 'a' && this.s[0] <= 'z') || (this.s[0] >= 'A' && this.s[0] <= 'Z') || this.s[0] == '_' {
					return this.expectIdOrKey(c)
				}
			}
		}
		if this.fp >= len(this.buf) {
			return newToken(TOKEN_EOF, "EOF", this.lineNum)
		}
		c = this.buf[this.fp]
		this.fp++
	}
	////文档末尾
	//if this.fp > len(this.buf) {
	//	return newToken(TOKEN_EOF, "EOF", this.lineNum)
	//}
	//fallthrough强制执行后面的case代码
	switch c {
	case '\'':
		var st = this.fp - 1
		if this.s == "" {
			this.fp++
			c = this.buf[this.fp]

			if c != '\'' && this.buf[this.fp-1] != '\\' {
				util.ParserError("'", string(c), this.lineNum, this.fname)
			} else {
				this.fp++
			}

			var ss = string(this.buf[st-1 : this.fp+1])
			this.fp++
			return newToken(TOKEN_ID, ss, this.lineNum)
		}
	case '&':
		if this.s == "" {
			if this.expectKeyword("&") {
				return newToken(TOKEN_LAND, "&&", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '|':
		if this.s == "" {
			if this.expectKeyword("|") {
				return newToken(TOKEN_LOR, "||", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '+':
		if this.s == "" {
			if this.expectKeyword("+") {
				return newToken(TOKEN_INCREMENT, "++", this.lineNum)
			} else if this.expectKeyword("=") {
				return newToken(TOKEN_ADD_ASSIGN, "+=", this.lineNum)
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
	case ':':
		if this.s == "" {
			if this.expectKeyword(":") {
				return newToken(TOKEN_DOUBLE_COLON, "::", this.lineNum)
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
				return newToken(TOKEN_DECREMENT, "--", this.lineNum)
			} else if this.expectKeyword("=") {
				return newToken(TOKEN_SUB_ASSIGN, "-=", this.lineNum)
			} else if this.expectKeyword(">") {
				return newToken(TOKEN_LAMBDA, "->", this.lineNum)
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
			} else if this.expectKeyword("<") {
				return newToken(TOKEN_SHL, "<<", this.lineNum)
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
	case '@':
		//return newToken(TOKEN_AT, "@", this.lineNum)
		fallthrough
	case '\t':
		fallthrough
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
		if this.s == "" {
			if this.expectKeyword("=") {
				return newToken(TOKEN_QUO_ASSIGN, "/=", this.lineNum)
			} else {
				comment := this.lex_Comments(c)
				if comment != "" {
					//log.Infof("注释:%v", comment)
					return newToken(TOKEN_COMMENT, comment, this.lineNum)
				} else {
					return this.expectIdOrKey(c)
				}
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '%':
		if this.s == "" {
			if this.expectKeyword("=") {
				return newToken(TOKEN_REM_ASSIGN, "%=", this.lineNum)
			} else {
				return this.expectIdOrKey(c)
			}
		} else {
			return this.expectIdOrKey(c)
		}
	case '"':
		if this.s == "" {
			return newToken(TOKEN_ID, this.lex_String(c), this.lineNum)
		}
	//字符串
	default:
		this.s += string(c)
	}

	return nil

}
func (this *Lexer) LexerBug(info string) {
	var msg = fmt.Sprintf("[%v] %d:%s\n", path.Base(this.fname), this.lineNum, info)
	util.Bug(msg)
}
