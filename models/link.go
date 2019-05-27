package models

import (
	"github.com/jinzhu/gorm"
)

type Link struct {
	gorm.Model
	Name string
	Url string
	Sort int  `gorm:"default:'0'"`
	View int
}

func (link *Link) Insert() error {
	return DB.FirstOrCreate(link, "url=?", link.Url).Error
}

func (link *Link) Update() error {
	return DB.Save(link).Error
}

func (link *Link) Delete() error {
	return DB.Delete(link).Error
}

func ListLinks()([]*Link,error) {
	var links []*Link
	err := DB.Order("sort asc").Find(&links).Error
	return links, err
}

func MustListLinks() []*Link {
	links, _ := ListLinks()
	return  links
}

func GetLinkById(id uint) (*Link, error) {
	var link Link
	err := DB.FirstOrCreate(&link, "id=?", id).Error
	return &link, err
}



