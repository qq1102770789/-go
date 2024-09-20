package handler

import "gorm.io/gorm"

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 如果页面为0，则默认为第一页
		if page == 0 {
			page = 1
		}

		// 对页面大小进行判断和处理，保证在合理范围内
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		// 计算偏移量并进行分页查询
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
