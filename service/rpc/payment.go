package rpc

import (
	"context"
	
	"github.com/shinmigo/pb/shoppb"
	"goshop/service-shop/model/payment"
)

type Payment struct {
}

func NewPayment() *Payment {
	return &Payment{}
}

func (p *Payment) GetPaymentList(ctx context.Context, req *shoppb.ListPaymentReq) (*shoppb.ListPaymentRes, error) {
	rows, err := payment.GetPaymentList(int32(req.Status))
	
	if err != nil {
		return nil, err
	}
	
	list := make([]*shoppb.Payment, 0, len(rows))
	for k := range rows {
		list = append(list, &shoppb.Payment{
			Id:     rows[k].Id,
			Code:   shoppb.PaymentCode(rows[k].Code),
			Name:   rows[k].Name,
			Params: rows[k].Params,
			Status: shoppb.PaymentStatus(rows[k].Status),
		})
	}
	return &shoppb.ListPaymentRes{
		Payments: list,
	}, nil
}

func (p *Payment) GetPaymentDetail(ctx context.Context, req *shoppb.PaymentCodeReq) (*shoppb.Payment, error) {
	row, err := payment.GetOneByCode(int32(req.Code), int32(shoppb.PaymentStatus_Open))
	
	if err != nil {
		return nil, err
	}
	
	return &shoppb.Payment{
		Id:     row.Id,
		Code:   shoppb.PaymentCode(row.Code),
		Name:   row.Name,
		Params: row.Params,
		Status: shoppb.PaymentStatus(row.Status),
	}, nil
}
