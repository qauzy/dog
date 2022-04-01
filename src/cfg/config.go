package cfg

var (
	ConstructFieldFunc = false //构建Get,Set函数
	AllStatic          = false //所有变量和方法按静态类型处理
	AppendContext      = false //添加*gin.Contex
	FieldAccess        = false //Get,Set函数转换为直接成员访问
	NoGeneric          = false //禁用泛型
	DropResult         = false //去掉返回值
	OneFold            = true  //独立文件夹
	ConstructNewFunc   = false //是否构建构造New函数
	Construct2New      = true  //是否转化调用无惨构造函数为调用new
	Capitalize         = true  //类成员大写开头
	StarClassTypeParam = true  //非原生类型函数参数带*
	StarClassTypeDecl  = true  //非原生类型变量带*
	ParseOnly          = false //只解析，不翻译(获取全局信息)
	ProjectName        = "bitrade"
	SourceBase         = "/opt/code/ZTuoExchange_framework"                                        //待转换源代码工程目录
	SourcePath         = "/opt/code/ZTuoExchange_framework/core/src/main/java/cn/ztuo/bitrade/dao" //待转换源代码目录
	TargetPath         = "/opt/code/bitrade/core"                                                  //目标目录
	ImportBase         = "bitrade/core"
)
