package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (t JsonTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
} //json自动调取MarshalJSON方法，将时间格式化为字符串，这样不仅能使用time类型本身的方法，还可以json化
type UserResponse struct {
	Id       int32    `json:"id"`
	NickName string   `json:"name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
