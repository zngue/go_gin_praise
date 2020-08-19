package praise

import "github.com/zngue/go_gin_praise/app/request"

func (p *Praise) GetNumByIDs(request request.PraiseRequest) map[int64]string  {
	return p.Praise(request).GetIDArrNum(request.IDArr)
}
