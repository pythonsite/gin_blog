package models

import (
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"

)

type Tag struct {
	BaseModel
	Name string
	Total int    `gorm:"-"`
}

func ListTagByPostId(id string)([]*Tag, error) {
	var tags []*Tag
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	rows,err := DB.Raw("select t.* from tags t inner join post_tags pt on t.id=pt.tag_id where pt.post_id=?", uint(pid)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		_ = DB.ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}

func ListTag()([]*Tag, error) {
	var tags []*Tag
	rows, err := DB.Raw("select t.*,count(*) total from tags t inner join post_tags pt on t.id = pt.tag_id inner join posts p on pt.post_id = p.id where p.is_published = ? group by pt.tag_id", true).Rows()
	logs.Info("rows:%#v", rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var tag Tag
		_ = DB.ScanRows(rows, &tag)
		logs.Info("%#v", tag)
		tags = append(tags, &tag)
	}
	logs.Info("all tags:%#v",tags)
	return tags, nil
}

func MustListTag()[]*Tag {
	tags, _ := ListTag()
	return tags
}

func GetTagsStr() (tagstr string) {
	tags, err := ListTag()
	if err != nil {
		return
	}
	tagNames := make([]string, 0)
	for _, tag := range tags {
		tagNames = append(tagNames,tag.Name)
	}
	tagstr = strings.Join(tagNames, ",")
	return
}

func CountTag() int {
	var count int
	DB.Model(&Tag{}).Count(&count)
	return count
}

func (tag *Tag) Insert() error {
	return DB.FirstOrCreate(tag, "name=?", tag.Name).Error
}

