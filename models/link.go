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

func ListLinks()([]*Link,error) {
	var links []*Link
	err := DB.Order("sort asc").Find(&Link{}).Error
	return links, err
}

func MustListLinks() []*Link {
	links, _ := ListLinks()
	return  links
}