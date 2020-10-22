package payment

import (
	"fmt"
	
	"goshop/service-shop/pkg/db"
)

type Payment struct {
	Id     uint64 `json:"id" gorm:"PRIMARY_KEY"`
	Code   int32
	Name   string
	Params string
	Status int32
}

func GetTableName() string {
	return "payment"
}

func GetField() []string {
	return []string{
		"id", "code", "name", "params", "status",
	}
	
}

func GetPaymentList(status int32) ([]*Payment, error) {
	rows := make([]*Payment, 0, 8)
	
	query := db.Conn.Table(GetTableName()).Select(GetField())
	
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	
	err := query.Find(&rows).Order("sort desc").Error
	
	if err != nil {
		return nil, err
	}
	
	return rows, nil
}

func GetOneByCode(code int32, status int32) (*Payment, error) {
	if code == 0 {
		return nil, fmt.Errorf("code is null")
	}
	row := &Payment{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("code = ? and status = ?", code, status).
		First(row).Error
	
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}
