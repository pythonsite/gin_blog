package models

type Page struct {
	BaseModel
	Title string
	Body string
	View int
	IsPublished bool
}

func (page *Page) Insert() error {
	return DB.Create(page).Error
}

func (page *Page) Update() error {
	return DB.Model(page).Updates(map[string]interface{}{
		"title": page.Title,
		"body": page.Body,
		"is_published": page.IsPublished,
	}).Error
}

func (page *Page) UpdateView() error {
	return DB.Model(page).Updates(map[string]interface{}{
		"view": page.View,
	}).Error
}

func (page *Page) Delete() error {
	return DB.Delete(page).Error
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

func ListAllPage()([]*Page, error) {
	return listPage(false)
}

