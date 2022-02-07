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
func FindByPack(p string) (result []*PackInfo, err error) {
	err = db.Where("path = ?", p).Find(&result).Error
	return
}
func ListPack() (result []*PackInfo, err error) {

	return
}
