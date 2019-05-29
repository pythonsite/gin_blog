package models

import (
	"database/sql"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"strconv"
)

type Post struct {
	BaseModel
	Title string
	Body string
	View int
	IsPublished bool
	Tags []*Tag			`gorm:"-"`
	Comments []*Comment	`gorm:"-"`
	CommentTotal int
}

func listPost(tag string, published bool, pageIndex, pageSize int)([]*Post, error) {
	var posts []*Post
	var err error
	if len(tag) > 0 {
		tagId, err := strconv.ParseUint(tag,10, 64)
		if err != nil {
			return nil, err
		}
		var rows *sql.Rows
		if published {
			if pageIndex > 0 {
				rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? and p.is_published = ? order by created_at desc limit ? offset ?", tagId, true, pageSize, (pageIndex-1)*pageSize).Rows()
			} else {
				rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id=? and p.is_published=? order by create_at desc", tagId, true).Rows()
			}
		} else {
			rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id=pt.post_id where pt.tag_id=? order by created_at desc", tagId).Rows()
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var post Post
			_ = DB.ScanRows(rows, &post)
			posts = append(posts, &post)
		}
	} else {
		if published {
			if pageIndex > 0 {
				err = DB.Where("is_published=?",true).Order("created_at desc").Limit(pageSize).Offset((pageIndex-1)*pageSize).Find(&posts).Error
			} else {
				err = DB.Where("is_published = ?", true).Order("created_at desc").Find(&posts).Error
			}
		} else {
			err = DB.Order("created_at desc").Find(&posts).Error
		}
	}
	return posts, err
}

func ListPublishedPost(tag string, pageIndex, pageSize int)([]*Post,error) {
	return listPost(tag, true, pageIndex, pageSize)
}

func ListAllPost(tag string)([]*Post, error) {
	return listPost(tag, false, 0, 0)
}

func CountPostByTag(tag string)(count int,err error) {
	var (
		tagId uint64
	)
	if len(tag) > 0 {
		tagId, err = strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return
		}
		err = DB.Raw("select count(1) from posts p inner join post_tags pt on p.id=pt.post_id where pt.tag_id=? and p.is_published=?", tagId, true).Row().Scan(&count)
	} else {
		err = DB.Raw("select count(1) from posts p where p.is_published=?", true).Row().Scan(&count)
	}
	return
}

func ListMaxReadPost() (posts []*Post, err error) {
	err = DB.Where("is_published=?",true).Order("view desc").Limit(5).Find(&posts).Error
	return
}

func MustListMaxReadPost() (posts []*Post) {
	posts, _ = ListMaxReadPost()
	return
}

func ListMaxCommentPost()(posts []*Post, err error) {
	var (
		rows *sql.Rows
	)
	rows, err = DB.Raw("select p.*,c.total comment_total from posts p inner join (select post_id, count(1) total from comments group by post_id) c on p.id=c.post_id order by c.total desc limit 5").Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		_ = DB.ScanRows(rows, &post)
		posts = append(posts, &post)

	}
	return
}

func MustListMaxCommentPost() (posts []*Post) {
	posts, _ = ListMaxCommentPost()
	return
}

func CountPost() int {
	var count int
	DB.Model(&Post{}).Count(&count)
	return count
}

func (post *Post) Insert() error {
	return DB.Create(post).Error
}

func (post *Post) Update() error {
	return DB.Model(post).Updates(map[string]interface{}{
		"title": post.Title,
		"body": post.Body,
		"is_published": post.IsPublished,
	}).Error
}

func (post *Post) Delete() error {
	return DB.Delete(post).Error
}

func GetPostById(id string) (*Post, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var post Post
	err = DB.First(&post, "id=?", pid).Error
	return &post, err
}

func (post *Post) UpdateView() error {
	return DB.Model(post).Updates(map[string]interface{}{
		"view": post.View,
	}).Error
}

func (post *Post) Excerpt() template.HTML {
	policy := bluemonday.StrictPolicy()
	sanitized := policy.Sanitize(string(blackfriday.Run([]byte(post.Body))))
	runnes := []rune(sanitized)
	if len(runnes) > 300 {
		sanitized = string(runnes[:300])
	}
	excerpt := template.HTML(sanitized + "...")
	return excerpt
}




