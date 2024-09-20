package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primarykey" json:"id"` //为什么使用int32，
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"-"default:"CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	//DeletedAt gorm.DeletedAt
	IsDeleted bool `gorm:"column:is_deleted" json:"-"`
}

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
