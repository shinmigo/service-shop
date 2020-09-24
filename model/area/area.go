package area

import (
	"goshop/service-shop/pkg/db"

	"github.com/jinzhu/gorm"
)

type Area struct {
	Id       uint64 `json:"id" gorm:"PRIMARY_KEY"`
	ParentId uint64 `json:"parent_id"`
	Level    uint8  `json:"level"`
	Code     uint64 `json:"code"`
	Name     string `json:"name"`
}

func GetTableName() string {
	return "area"
}

func GetField() []string {
	return []string{
		"id", "parent_id", "level", "code", "name",
	}

}

func GetAreaList(level uint8, parentId int64) ([]*Area, error) {
	conditions := make([]func(db *gorm.DB) *gorm.DB, 0, 4)
	if level > 0 {
		if level > 4 {
			level = 1
		}
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("level = ?", level)
		})
	} else {
		// 只取三级
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("level < 4")
		})
	}

	if parentId > -1 {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id = ?", parentId)
		})
	}

	rows := make([]*Area, 0, 32)
	if err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Scopes(conditions...).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}
