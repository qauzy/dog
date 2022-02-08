package storage

import "github.com/jinzhu/gorm"

var db *gorm.DB

func init() {
	db = NewDBWrite("dog.db")
}

func AddPack(p *PackInfo) (err error) {
	err = db.Save(p).Error
	return
}

func FindByName(name string) (result *PackInfo, err error) {
	result = new(PackInfo)
	err = db.Where("name = ?", name).First(result).Error
	return
}

func FindByPath(path string) (result []*PackInfo, err error) {
	err = db.Where("path = ?", path).Find(&result).Error
	return
}
func FindByPack(p string) (result []*PackInfo, err error) {
	err = db.Where("path = ?", p).Find(&result).Error
	return
}
func ListPack() (result []*PackInfo, err error) {

	return
}
