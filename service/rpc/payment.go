package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"goshop/service-shop/model/payment"
	"goshop/service-shop/pkg/db"

	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/shoppb"
)

type Payment struct {
}

func NewPayment() *Payment {
	return &Payment{}
}

func (p *Payment) DeletePayment(ctx context.Context, req *shoppb.DeletePaymentReq) (*basepb.AnyRes, error) {
	res := &basepb.AnyRes{
		Id:    req.Id,
		State: 0,
	}
	if req.Id <= 0 {
		return res, fmt.Errorf("id不合法")
	}

	if err := db.Conn.Table(payment.GetTableName()).Where("id = ?", req.Id).Delete(nil).Error; err != nil {
		return res, fmt.Errorf("删除失败， err: %v", err)
	}

	res.State = 1

	return res, nil
}

func (p *Payment) AddPayment(ctx context.Context, req *shoppb.Payment) (*basepb.AnyRes, error) {
	buf := &payment.Payment{
		Code:   req.Code,
		Name:   req.Name,
		Params: req.Params,
		Status: int8(req.Status),
		Sort:   req.Sort,
	}

	res := &basepb.AnyRes{
		Id:    0,
		State: 0,
	}

	if err := db.Conn.Create(buf).Error; err != nil {
		return res, fmt.Errorf("添加失败, err: %v", err)
	}

	res.Id = buf.Id
	res.State = 1

	return res, nil
}

func (p *Payment) EditPayment(ctx context.Context, req *shoppb.Payment) (*basepb.AnyRes, error) {
	buf := map[string]interface{}{
		"id":     req.Id,
		"name":   req.Name,
		"code":   req.Code,
		"status": req.Status,
		"params": req.Params,
		"sort":   req.Sort,
	}

	res := &basepb.AnyRes{
		Id:    req.Id,
		State: 0,
	}
	if err := db.Conn.Table(payment.GetTableName()).Where("id = ?", req.Id).Update(buf).Error; err != nil {
		return res, err
	}

	res.State = 1

	return res, nil
}

func (p *Payment) GetPaymentDetail(ctx context.Context, req *shoppb.PaymentCodeReq) (*shoppb.Payment, error) {
	if len(req.Code) == 0 {
		return nil, fmt.Errorf("codeId不合法")
	}

	info := &payment.Payment{}
	if err := db.Conn.Table(payment.GetTableName()).Where("code = ?", req.Code).Find(info).Error; err != nil {
		return nil, fmt.Errorf("获取明细失败, err: %v", err)
	}

	resList := &shoppb.Payment{}
	_ = json.Unmarshal(func() []byte { buf, _ := json.Marshal(info); return buf }(), resList)

	return resList, nil
}

func (p *Payment) GetPaymentList(ctx context.Context, req *shoppb.ListPaymentReq) (*shoppb.ListPaymentRes, error) {
	list := make([]*payment.Payment, 0, 4)

	query := db.Conn.Table(payment.GetTableName())
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	resList := make([]*shoppb.Payment, 0, len(list))
	_ = json.Unmarshal(func() []byte { buf, _ := json.Marshal(&list); return buf }(), &resList)

	return &shoppb.ListPaymentRes{Payments: resList}, nil
}
