package models


type Page struct {
	BaseModel
	Title string
	Body string
	View int
	IsPublished bool
}

func listPage(published bool) ([]*Page, error) {
	var pages []*Page
	var err error
	if published {
		err = DB.Where("is_published=?", true).Find(&pages).Error
	} else {
		err = DB.Find(&pages).Error
	}
	return pages, err
}

func ListPublishedPage() ([]*Page, error) {
	return listPage(true)
}

func CountPage() int {
	var count int
	DB.Model(&Page{}).Count(&count)
	return count
}