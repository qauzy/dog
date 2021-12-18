package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	l = log.New(os.Stderr, "", 0)
)

func Bug(info string) {
	var msg = joint(4, "[BUG]", info)
	l.Output(3, msg)
	os.Exit(0)
}

func Debug(info string) {
	var msg = joint(4, "[DEBUG]", info)
	l.Output(5, msg)
}

func ParserError(expect string, current string, linenum int, file string) {
	var msg = joint(4, "[ERROR]", fmt.Sprintf("Expect: <%s>, but got <%s> at line:%d file:%s\n", expect, current, linenum, file))
	l.Output(2, msg)
	os.Exit(0)
}

func joint(level int, prefix, message string) string {
	now := time.Now().Format("2006/01/02 15:04:05")
	filename, funcname, line := getpProcInfo(level)
	s := fmt.Sprint(prefix, ": ", now, " ", filename, ":", line, ":", funcname, ": ", message)
	return s
}

//获取打印日志的进程信息
func getpProcInfo(level int) (filename, funcname string, line int) {
	pc, filename, line, ok := runtime.Caller(level)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
		filename = filepath.Base(filename)           // /full/path/basename.go => basename.go
	}
	return
}
