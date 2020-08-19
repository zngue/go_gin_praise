package praise

import (
	"github.com/zngue/go_gin_praise/app/praise"
	"github.com/zngue/go_gin_praise/app/request"
)
func (p *Praise) GetNumByMouthRank(request request.PraiseRequest) *praise.RankList {
	p1:=p.Praise(request)
	return  p1.GetRange(request.Sort,p1.ZachMouthKey)
}
