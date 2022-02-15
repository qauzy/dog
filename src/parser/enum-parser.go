package parser

import (
	"dog/ast"
	log "github.com/corgi-kx/logcustom"
)

func (this *Parser) parseEnumDecl(access int) (cl ast.Class) {
	var id, extends string

	this.eatToken(TOKEN_ENUM)
	id = this.current.Lexeme
	this.eatToken(TOKEN_ID)
	//处理implements
	if this.current.Kind == TOKEN_IMPLEMENTS {
		this.eatToken(TOKEN_IMPLEMENTS)
		extends = this.current.Lexeme
		this.eatToken(TOKEN_ID)
	}
	this.eatToken(TOKEN_LBRACE)
	classSingle := ast.NewClassSingle(this.currentFile, access, id, extends, ast.ENUM_TYPE)

	this.currentClass = classSingle
	defer func() {
		this.currentClass = nil
	}()
	//处理枚举变量
	for {
		var comment string
		//处理注释
		if this.current.Kind == TOKEN_COMMENT {
			comment = ""
			for this.current.Kind == TOKEN_COMMENT {
				comment += "\n" + this.current.Lexeme
				this.advance()
			}
			if this.current.Kind == TOKEN_EOF || (this.current.Kind != TOKEN_COMMENT) {
				continue
			}
			log.Infof("注释-->%v", comment)
		}

		id = this.current.Lexeme
		this.eatToken(TOKEN_ID)
		//FIXME 只支持一个值枚举
		if this.current.Kind == TOKEN_LPAREN {
			this.eatToken(TOKEN_LPAREN)
			args := this.parseExpList()
			classSingle.AddField(ast.NewFieldEnum(access, nil, id, args, false, true))
			this.eatToken(TOKEN_RPAREN)
		} else {
			classSingle.AddField(ast.NewFieldEnum(access, nil, id, nil, false, true))
		}

		for this.current.Kind == TOKEN_COMMENT {
			this.advance()
		}
		if this.current.Kind == TOKEN_SEMI {
			this.eatToken(TOKEN_SEMI)
			//兼容多一个SEMI
			if this.current.Kind == TOKEN_SEMI {
				this.eatToken(TOKEN_SEMI)
			}
			break
			//兼容多一个逗号的情况
		} else if this.current.Kind == TOKEN_COMMER {
			this.eatToken(TOKEN_COMMER)
		}

	}
	//
	this.parseClassContext(classSingle)

	this.eatToken(TOKEN_RBRACE)

	return classSingle
}
