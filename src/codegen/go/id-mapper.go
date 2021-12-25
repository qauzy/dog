package codegen_go

var (
	IdMapper = make(map[string]string)
)

func init() {
	IdMapper["StringBuffer"] = "strutils.NewStringBuilder"
}
