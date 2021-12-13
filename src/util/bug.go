package util

import (
	log "github.com/corgi-kx/logcustom"
	"os"
	"path"
)

func Bug(info string, filename string, linenum int) {
	log.Errorf("ERROR> %s:%d:%s\n", path.Base(filename), linenum, info)
	os.Exit(0)
}

func ParserError(expect string, current string, linenum int) {
	log.Errorf("Expect: <%s>, but got <%s> at line:%d\n", expect, current, linenum)
	os.Exit(0)
}
