package storage

type PackInfo struct {
	Project string `gorm:"unique_index:path_uni"`
	Name    string `gorm:"unique_index:path_uni"`
	Path    string `gorm:"unique_index:path_uni"`
	Kind    int    `gorm:"unique_index:path_uni"` //0:class 1: enum 2: interface
}
