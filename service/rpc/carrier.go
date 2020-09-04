package rpc

import (
	"context"

	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/shoppb"

	"goshop/service-shop/model/carrier"
	"goshop/service-shop/pkg/db"
	"goshop/service-shop/pkg/utils"
)

type Carrier struct {

}

func NewCarrier() *Carrier  {
	return &Carrier{}
}

func (c *Carrier) GetCarrierList(ctx context.Context, req *shoppb.ListCarrierReq) (*shoppb.ListCarrierRes, error)  {
	var page uint64 = 1
	if req.Page > 0 {
		page = req.Page
	}

	var pageSize uint64 = 10
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	rows, total, err := carrier.GetCarriers(req, page, pageSize)
	if err != nil {
		return nil, err
	}

	carrierDetails := make([]*shoppb.CarrierDetail, 0, req.PageSize)

	for k := range rows {
		carrierDetails = append(carrierDetails, &shoppb.CarrierDetail{
			CarrierId:            rows[k].CarrierId,
			Name:                 rows[k].Name,
			Code:                 rows[k].Code,
			Sort:                 rows[k].Sort,
			Status:               rows[k].Status,
			CreatedBy:            rows[k].CreatedBy,
			UpdatedBy:            rows[k].UpdatedBy,
			CreatedAt:            rows[k].CreatedAt.Format(utils.TIME_STD_FORMART),
			UpdatedAt:            rows[k].UpdatedAt.Format(utils.TIME_STD_FORMART),
		})
	}
	return &shoppb.ListCarrierRes{
		Total:                total,
		Carriers:             carrierDetails,
	}, nil
}

func (c *Carrier) AddCarrier(ctx context.Context, req *shoppb.Carrier) (*basepb.AnyRes, error)  {
	carrier := carrier.Carrier{
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		Status:    req.Status,
		CreatedBy: req.AdminId,
		UpdatedBy: req.AdminId,
	}
	if err := db.Conn.Create(&carrier).Error; err != nil {
		return nil, err
	}
	
	return &basepb.AnyRes{
		Id:                   carrier.CarrierId,
		State:                1,
	}, nil
}

func (c *Carrier) DelCarrier(ctx context.Context, req *shoppb.DelCarrierReq) (*basepb.AnyRes, error)  {
	if err := db.Conn.Where("carrier_id = ?", req.CarrierId).Delete(&shoppb.Carrier{}).Error; err != nil {
		return nil, err
	}

	return &basepb.AnyRes{
		Id:                   req.CarrierId,
		State:                1,
	}, nil
}

func (c *Carrier) EditCarrier(ctx context.Context, req *shoppb.Carrier) (*basepb.AnyRes, error) {
	if ok, err := carrier.ExistCarrierById(req.CarrierId); !ok {
		return nil, err
	}
	row := carrier.Carrier{
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		Status:    req.Status,
		CreatedBy: req.AdminId,
		UpdatedBy: req.AdminId,
	}
	if err := db.Conn.Model(&carrier.Carrier{CarrierId:req.CarrierId}).Update(row).Error; err != nil {
		return nil, err
	}
	return &basepb.AnyRes{
		Id:                   req.CarrierId,
		State:                1,
	}, nil
}

func (c *Carrier) EditCarrierStatus(ctx context.Context, req *shoppb.EditCarrierStatusReq) (*basepb.AnyRes, error)  {
	db.Conn.Table(carrier.GetTableName()).Where("carrier_id in (?)", req.CarrierId).Updates(map[string]interface{}{
		"status":     req.Status,
		"updated_by": req.AdminId,
	})
	return &basepb.AnyRes{
		Id:                   0,
		State:                1,
	}, nil
}
