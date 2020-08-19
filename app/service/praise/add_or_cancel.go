package praise

import "github.com/zngue/go_gin_praise/app/request"
//点赞或者取消点赞
func (p *Praise)  AddOrCancel(request request.PraiseRequest) (int64,error) {
	return p.Praise(request).Praise()
}



