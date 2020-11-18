package rpc

import (
	"context"
	"goshop/service-shop/model/banner"
	"goshop/service-shop/pkg/db"
	"goshop/service-shop/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/shinmigo/pb/basepb"

	"github.com/shinmigo/pb/shoppb"
)

type BannerAd struct {
}

func NewBannerAd() *BannerAd {
	return &BannerAd{}
}

func (m *BannerAd) GetBannerAdList(ctx context.Context, req *shoppb.ListBannerAdReq) (*shoppb.ListBannerAdRes, error) {
	lists, total, err := banner.GetBannerAds(req)
	if err != nil {
		return nil, err
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "The client canceled the request")
	}

	detailList := make([]*shoppb.BannerAdDetail, 0, req.PageSize)
	for k := range lists {
		detailList = append(detailList, &shoppb.BannerAdDetail{
			Id:          lists[k].Id,
			EleType:     lists[k].EleType,
			ImageUrl:    lists[k].ImageUrl,
			RedirectUrl: lists[k].RedirectUrl,
			Sort:        lists[k].Sort,
			Status:      lists[k].Status,
			TagName:     lists[k].TagName,
			CreatedBy:   lists[k].CreatedBy,
			UpdatedBy:   lists[k].UpdatedBy,
			CreatedAt:   lists[k].CreatedAt.Format(utils.TIME_STD_FORMART),
			UpdatedAt:   lists[k].UpdatedAt.Format(utils.TIME_STD_FORMART),
		})
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "The client canceled the request")
	}

	return &shoppb.ListBannerAdRes{
		Total:     total,
		BannerAds: detailList,
	}, nil
}

func (m *BannerAd) AddBannerAd(ctx context.Context, req *shoppb.BannerAd) (*basepb.AnyRes, error) {
	bannerAdd := banner.BannerAd{
		EleType:     req.EleType,
		ImageUrl:    req.ImageUrl,
		RedirectUrl: req.RedirectUrl,
		Sort:        req.Sort,
		Status:      req.Status,
		TagName:     req.TagName,
		CreatedBy:   req.AdminId,
		UpdatedBy:   req.AdminId,
	}
	err := db.Conn.Create(&bannerAdd).Error
	if err != nil {
		return nil, err
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "The client canceled the request")
	}

	return &basepb.AnyRes{
		Id:    bannerAdd.Id,
		State: 1,
	}, nil
}

func (m *BannerAd) EditBannerAd(ctx context.Context, req *shoppb.BannerAd) (*basepb.AnyRes, error) {
	ok, err := banner.FindBannerAdById(req.Id)
	if !ok {
		return nil, err
	}

	updateData := banner.BannerAd{
		EleType:     req.EleType,
		ImageUrl:    req.ImageUrl,
		RedirectUrl: req.RedirectUrl,
		Sort:        req.Sort,
		TagName:     req.TagName,
		Status:      req.Status,
		CreatedBy:   req.AdminId,
		UpdatedBy:   req.AdminId,
	}
	err = db.Conn.Table(banner.GetTableName()).Where("id = ?", req.Id).Update(updateData).Error
	if err != nil {
		return nil, err
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "The client canceled the request")
	}

	return &basepb.AnyRes{
		Id:    req.Id,
		State: 1,
	}, nil
}

func (m *BannerAd) EditBannerAdStatus(ctx context.Context, req *shoppb.EditBannerAdStatusReq) (*basepb.AnyRes, error) {
	db.Conn.Table(banner.GetTableName()).Where("id in (?)", req.Id).Updates(map[string]interface{}{
		"status":     req.Status,
		"updated_by": req.AdminId,
	})
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "The client canceled the request")
	}

	return &basepb.AnyRes{
		Id:    0,
		State: 1,
	}, nil
}

func (m *BannerAd) DelBannerAd(ctx context.Context, req *shoppb.DelBannerAdReq) (*basepb.AnyRes, error) {
	if err := db.Conn.Where("id in (?)", req.Id).Delete(&banner.BannerAd{}).Error; err != nil {
		return nil, err
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "The client canceled the request")
	}

	return &basepb.AnyRes{
		Id:    req.Id[0],
		State: 1,
	}, nil
}
