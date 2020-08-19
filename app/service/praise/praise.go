package praise

import (
	"github.com/zngue/go_gin_praise/app/praise"
	"github.com/zngue/go_gin_praise/app/request"
	"strconv"
)

type Praise struct {

}
func (Praise)  Praise(request request.PraiseRequest) *praise.Praise {
	p1:=praise.Praise{
		ID: request.ID,
		From: request.From,
		SignNum: request.SingleNum,
		TypeID: request.TypeID,
		IsYear:request.IsYear,
		IsMouth: request.IsMouth,
		Device:praise.Device{
			UnionID: request.Device.UnionID,
			UserID: strconv.Itoa(request.Device.UserID),
			Uuid: request.Device.UuID,
		},
		RedisPage:praise.RedisPage{
			Page: int64(request.Page.Page),
			PageSize: int64(request.PageSize),
		},
	}
	p1.Init()
	return &p1
}
