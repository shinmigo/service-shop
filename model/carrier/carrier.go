package carrier

import (
	"fmt"
	"goshop/service-shop/pkg/db"
	"goshop/service-shop/pkg/utils"

	"github.com/shinmigo/pb/shoppb"
)

type Carrier struct {
	CarrierId uint64 `json:"carrier_id" gorm:"PRIMARY_KEY"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Sort      uint32 `json:"sort"`
	Status    shoppb.CarrierStatus
	CreatedBy uint64          `json:"created_by"`
	UpdatedBy uint64          `json:"updated_by"`
	CreatedAt utils.JSONTime  `json:"created_at"`
	UpdatedAt utils.JSONTime  `json:"updated_at"`
	DeletedAt *utils.JSONTime `json:"deleted_at"`
}

func GetTableName() string {
	return "carrier"
}

func GetField() []string {
	return []string{
		"carrier_id", "name", "code", "sort", "name", "status", "created_by", "updated_by", "created_at", "updated_at",
	}
}

func GetCarriers(req *shoppb.ListCarrierReq, page, pageSize uint64) ([]*Carrier, uint64, error) {
	var total uint64

	rows := make([]*Carrier, 0, req.PageSize)

	query := db.Conn.Table(GetTableName()).Select(GetField())
	if req.Id > 0 {
		query = query.Where("carrier_id = ?", req.Id)
	}

	if req.Name != "" {
		query = query.Where("name like ?", "%"+req.Name+"%")
	}

	if req.Code != "" {
		query = query.Where("code like ?", "%"+req.Code+"%")
	}

	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)
	err := query.Order("carrier_id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error
	if err != nil {
		return nil, total, err
	}

	return rows, total, nil
}

func ExistCarrierById(id uint64) (bool, error) {
	if id == 0 {
		return false, fmt.Errorf("carrier is null")
	}
	carrier := Carrier{}
	err := db.Conn.Select("carrier_id").Where("carrier_id=?", id).First(&carrier).Error
	if err != nil {
		return false, fmt.Errorf("err: %v", err)
	}

	return true, nil
}
