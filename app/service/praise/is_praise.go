package praise

import "github.com/zngue/go_gin_praise/app/request"

//判断是否点赞
func (p *Praise) IsPraise(request request.PraiseRequest)  (bool,error) {
	return p.Praise(request).IsPraise()
}
