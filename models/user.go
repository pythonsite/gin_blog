package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

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

func (user *User) Insert() error {
	return DB.Create(user).Error
}

func (user *User) Update() error {
	return DB.Save(user).Error
}

func (user *User) Delete() error {
	return DB.Delete(user).Error
}

func (user *User) Lock() error {
	return DB.Model(user).Update(map[string]interface{}{
		"lock_state": user.LockState,
	}).Error
}

func GetUserByUsername(username string)(*User, error) {
	var user User
	err := DB.First(&user, "email=?",username).Error
	return &user, err
}

func GetUser(id interface{})(*User,error) {
	var user User
	err := DB.First(&user, id).Error
	return &user, err
}

func ListUsers()([]*User, error) {
	var users []*User
	err := DB.Find(&users, "is_admin=?", false).Error
	return users,err
}