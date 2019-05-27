package models

import "strconv"

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

func CountComment() int {
	var count int
	DB.Model(&Comment{}).Count(&count)
	return count
}

func ListUnreadComment() ([]*Comment,error) {
	var comments []*Comment
	err := DB.Where("read_state=?",false).Order("created_at desc").Find(&comments).Error
	return comments,err
}

func MustListUnreadComment() []*Comment {
	comments, _ := ListUnreadComment()
	return comments
}

func ListCommentByPostID(postId string)([]*Comment, error) {
	pid, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		return  nil,err
	}
	var comments []*Comment
	rows, err := DB.Raw("select c.*,u.github_login_id nick_name,u.avatar_url,u.github_url from comments c inner join users u on c.user_id = u.id where c.post_id = ? order by created_at desc", uint(pid)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		_ = DB.ScanRows(rows, &comment)
		comments = append(comments, &comment)
	}
	return comments, err
}