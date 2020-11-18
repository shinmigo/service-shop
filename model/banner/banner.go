package banner

import (
	"fmt"
	"goshop/service-shop/pkg/db"
	"goshop/service-shop/pkg/utils"

	"github.com/shinmigo/pb/shoppb"
)

type BannerAd struct {
	Id        uint64 `json:"id" gorm:"PRIMARY_KEY"`
	EleType   uint32 `json:"ele_type"`
	TagName   string `json:"tag_name"`
	EleInfo   string `json:"ele_info"`
	Status    shoppb.BannerAdStatus
	CreatedBy uint64          `json:"created_by"`
	UpdatedBy uint64          `json:"updated_by"`
	CreatedAt utils.JSONTime  `json:"created_at"`
	UpdatedAt utils.JSONTime  `json:"updated_at"`
	DeletedAt *utils.JSONTime `json:"deleted_at"`
}

func GetTableName() string {
	return "banner_ad"
}

func GetField() []string {
	return []string{
		"id", "ele_type", "tag_name", "ele_info", "status", "created_by", "updated_by", "created_at", "updated_at",
	}
}

func GetBannerAds(req *shoppb.ListBannerAdReq) (lists []*BannerAd, total uint64, err error) {
	lists = make([]*BannerAd, 0, req.PageSize)

	query := db.Conn.Table(GetTableName()).Select(GetField())
	if req.Id > 0 {
		query = query.Where("id = ?", req.Id)
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.EleType > 0 {
		query = query.Where("ele_type = ?", req.EleType)
	}
	if len(req.TagName) > 0 {
		query = query.Where("tag_name like ?", req.TagName+"%")
	}

	query.Where("deleted_at is null").Count(&total)
	err = query.Order("id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&lists).Error
	return lists, total, err
}

func FindBannerAdById(id uint64) (res bool, err error) {
	if id <= 0 {
		return false, fmt.Errorf("data is null")
	}

	BannerAd := BannerAd{}
	err = db.Conn.Table(GetTableName()).Select("id").Where("id = ?", id).First(&BannerAd).Error
	if err != nil {
		return false, fmt.Errorf("err: %v", err)
	}

	return true, nil
}
