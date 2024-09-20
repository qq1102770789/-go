package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32          `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"column:add_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_time"`
	IsDeleted bool           `gorm:"column:is_deleted"`
}

/*
1.密文 2.密文不可以被反向解析
    1.对称加密
	2.非对称加密
	3.md5加密 信息摘要算法
	密码如果不可以反解，用户找回密码就比较麻烦，所以密码需要加密存储。

*/

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"` //便于修改
	Gender   string     `gorm:"column:gender;default:'male';type:varchar(6) comment'female表示女性，male表示男性'"`
	Role     int        `gorm:"column:role;default:1;type:int comment'1表示普通用户，2表示管理员'"`
}
