package rpc

import (
	"context"

	"goshop/service-shop/model/area"

	"github.com/shinmigo/pb/shoppb"
)

type Area struct {
}

func NewMArea() *Area {
	return &Area{}
}

func (m *Area) GetAreaList(ctx context.Context, req *shoppb.ListAreaReq) (*shoppb.ListAreaRes, error) {
	rows, err := area.GetAreaList(0, -1)
	if err != nil {
		return nil, err
	}

	provSclie := make([]*area.Area, 0, 64)
	citySclie := make(map[uint64][]*area.Area, 64)
	counSclie := make(map[uint64][]*shoppb.Coun, 64)
	for k := range rows {
		switch rows[k].Level {
		case 1: // 省
			provSclie = append(provSclie, rows[k])
			break
		case 2: // 市
			if len(citySclie[rows[k].ParentId]) == 0 {
				citySclie[rows[k].ParentId] = make([]*area.Area, 0, 64)
			}
			citySclie[rows[k].ParentId] = append(citySclie[rows[k].ParentId], rows[k])
			break
		case 3: // 区
			if len(counSclie[rows[k].ParentId]) == 0 {
				counSclie[rows[k].ParentId] = make([]*shoppb.Coun, 0, 64)
			}
			counSclie[rows[k].ParentId] = append(counSclie[rows[k].ParentId], &shoppb.Coun{
				Label: rows[k].Name,
				Value: rows[k].Code,
			})
			break
		default:
			break
		}
	}

	prov := make([]*shoppb.Prov, 0, len(provSclie))
	for k := range provSclie {
		provChildren := make([]*shoppb.City, 0, 64)
		if _, city := citySclie[provSclie[k].Id]; city {
			for _, c := range citySclie[provSclie[k].Id] {
				if _, coun := counSclie[c.Id]; !coun {
					continue
				}
				buf := &shoppb.City{
					Label:    c.Name,
					Value:    c.Code,
					Children: counSclie[c.Id],
				}
				provChildren = append(provChildren, buf)
			}
		}
		prov = append(prov, &shoppb.Prov{
			Label:    provSclie[k].Name,
			Value:    provSclie[k].Code,
			Children: provChildren,
		})
	}

	return &shoppb.ListAreaRes{
		Areas: prov,
	}, nil
}

func (m *Area) GetAreaNameByCodes(ctx context.Context, req *shoppb.AreaCodeReq) (*shoppb.AreaNameRes, error) {
	rows, err := area.GetAreaNameByCodes(req.Codes)
	if err != nil {
		return nil, err
	}

	return &shoppb.AreaNameRes{
		Codes: rows,
	}, nil
}
