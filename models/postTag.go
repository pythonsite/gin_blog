package models

type PostTag struct {
	BaseModel
	PostId uint
	TagId uint
}

func (pt *PostTag) Insert() error {
	return DB.FirstOrCreate(pt, "post_id=? and tag_id=?", pt.PostId, pt.TagId).Error
}

func DeletePostTagByPostId(postId uint) error {
	return DB.Delete(&PostTag{}, "post_id=?", postId).Error
}

