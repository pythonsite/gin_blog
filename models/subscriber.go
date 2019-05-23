package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

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
