package praise

import "github.com/zngue/go_gin_praise/app/request"

func (p *Praise) IsPraiseByIDs(request request.PraiseRequest) map[int64]bool  {
	return p.Praise(request).GetIDArrStatus(request.IDArr)
}
