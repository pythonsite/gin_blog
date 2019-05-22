package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
	"gin_blog/system"
)

type BaseModel struct {
	ID uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdateAt time.Time
}

type Page struct {
	BaseModel
	Title string
	Body string
	View int
	IsPublished bool
}

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

type Tag struct {
	BaseModel
	Name string
	Total int 	`ogorm:"-"`
}

type PostTag struct {
	BaseModel
	PostId uint
	TagId uint
}

type User struct {
	gorm.Model
	Email string `gorm:"unique_index;default:null"`
	Telephone string `gorm:"unique_index;default:null"`
	Password string `gorm:"default:null"`
	VerifyState   string    `gorm:"default:'0'"`
	SecretKey     string    `gorm:"default:null"`
	OutTime       time.Time //过期时间
	GithubLoginId string    `gorm:"unique_index;default:null"`
	GithubUrl     string
	IsAdmin       bool      //是否是管理员
	AvatarUrl     string    // 头像链接
	NickName      string    // 昵称
	LockState     bool      `gorm:"default:'0'"` //锁定状态
}

type Comment struct {
	BaseModel
	UserID uint
	Content string
	PostID uint
	ReadState bool   `gorm:"default:'0'"` // 阅读状态
	NickName string  `gorm:"-"`
	AvatarUrl string `gorm:"-"`
	GithubUrl string `gorm:"-"`
}

// table subscribe
type Subscriber struct {
	gorm.Model
	Email          string    `gorm:"unique_index"` //邮箱
	VerifyState    bool      `gorm:"default:'0'"`  //验证状态
	SubscribeState bool      `gorm:"default:'1'"`  //订阅状态
	OutTime        time.Time //过期时间
	SecretKey      string    // 秘钥
	Signature      string    //签名
}

type Link struct {
	gorm.Model
	Name string
	Url string
	Sort int  `gorm:"default:'0'"`
	View int
}

// query result
type QrArchive struct {
	ArchiveDate time.Time //month
	Total       int       //total
	Year        int       // year
	Month       int       // month
}

func listPage(published bool) ([]*Page, error) {
	var pages []*Page
	var err error
	if published {
		err = system.DB.Where("is_published=?", true).Find(&pages).Error
	} else {
		err = system.DB.Find(&pages).Error
	}
	return pages, err
}

func ListPublishedPage() ([]*Page, error) {
	return listPage(true)
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
				rows, err = system.DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? and p.is_published = ? order by created_at desc limit ? offset ?", tagId, true, pageSize, (pageIndex-1)*pageSize).Rows()
			} else {
				rows, err = system.DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id=? and p.is_published=? order by create_at desc", tagId, true).Rows()
			}
		} else {
			rows, err = system.DB.Raw("select p.* from posts p inner join post_tags pt on p.id=pt.post_id where pt.tag_id=? order by created_at desc", tagId).Rows()
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var post Post
			_ = system.DB.ScanRows(rows, &post)
			posts = append(posts, &post)
		}
	} else {
		if published {
			if pageIndex > 0 {
				err = system.DB.Where("is_published=?",true).Order("create_at desc").Limit(pageSize).Offset((pageIndex-1)*pageSize).Find(&posts).Error
			} else {
				err = system.DB.Where("is_published = ?", true).Order("created_at desc").Find(&posts).Error
			}
		}
	}
	return posts, err
}

func ListPublishedPost(tag string, pageIndex, pageSize int)([]*Post,error) {
	return listPost(tag, true, pageIndex, pageSize)
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
		err = system.DB.Raw("select count(1) from posts p inner join post_tags pt on p.id=pt.post_id where pt.tag_id=? and p.is_published=?", tagId, true).Row().Scan(&count)
	} else {
		err = system.DB.Raw("select count(1) from posts p where p.is_published=?", true).Row().Scan(&count)
	}
	return
}

func ListTagByPostId(id string)([]*Tag, error) {
	var tags []*Tag
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	rows,err := system.DB.Raw("select t.* from tags t inner join post_tags pt on t.id=pt.tag_id where pt.post_id=?", uint(pid)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		_ = system.DB.ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}



