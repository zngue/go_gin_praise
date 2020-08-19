package praise

import "github.com/zngue/go_gin_praise/app/request"

func (p *Praise) GetNum(request request.PraiseRequest) (string,error) {
	return p.Praise(request).GetNum()
}
