package models

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