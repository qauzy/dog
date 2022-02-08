package cfg

var (
	ConstructFieldFunc = true  //构建Get,Set函数
	AppendContext      = false //添加*gin.Contex
	FieldAccess        = true  //Get,Set函数转换为直接成员访问
	DropResult         = false //去掉返回值
	OneFold            = false //独立文件夹
	ConstructNewFunc   = false //是否构建构造New函数
	Construct2New      = true  //是否转化调用无惨构造函数为调用new
	ProjectName        = "bitrade"
	SourceBase         = "/opt/code/ZTuoExchange_framework"                                        //待转换源代码工程目录
	SourcePath         = "/opt/code/ZTuoExchange_framework/core/src/main/java/cn/ztuo/bitrade/dao" //待转换源代码目录
	TargetPath         = "/opt/3code/actJob/memberxxl/mdata"                                       //目标目录
	ImportBase         = "bitrade/core"
)
