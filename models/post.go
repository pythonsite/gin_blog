package models

import (
	"database/sql"
	"strconv"
	"time"
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

func ListPostArchives()([]*QrArchive, error) {
	var archives []*QrArchive
	querysql := `select strftime('%Y-%m', created_at) as month, count(*) as total from posts where is_published=? group by month order by month desc`
	rows, err := DB.Raw(querysql, true).Rows()
	if err != nil {
		return  nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var archive QrArchive
		var month string
		_ = rows.Scan(&month, &archive.Total)
		archive.ArchiveDate, _ = time.Parse("2006-01",month)
		archive.Year = archive.ArchiveDate.Year()
		archive.Month = int(archive.ArchiveDate.Month())
		archives = append(archives, &archive)
	}
	return archives, nil
}

func MustListPostArchives() []*QrArchive {
	archives, _ := ListPostArchives()
	return archives
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

