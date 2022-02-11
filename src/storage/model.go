package storage

type PackInfo struct {
	Id      int    `gorm:"-;primary_key;AUTO_INCREMENT"` //
	Project string `gorm:"unique_index:path_uni"`        //项目名
	Name    string `gorm:"unique_index:path_uni"`        //类名
	Path    string `gorm:"unique_index:path_uni"`        //包名
	Kind    int    `gorm:"unique_index:path_uni"`        //0:class 1: enum 2: interface
}

type FieldInfo struct {
	Id      int    `gorm:"-;primary_key;AUTO_INCREMENT"` //
	PackId  int    `gorm:"unique_index:path_uni"`        //类ID
	Project string `gorm:"unique_index:path_uni"`        //项目名
	Name    string `gorm:"unique_index:path_uni"`        //成员名
	Type    string `gorm:"unique_index:path_uni"`        //成员属性
	Kind    int    `gorm:"unique_index:path_uni"`        //成员类型 0:Member Var 1: Method 2: Param Var
}
