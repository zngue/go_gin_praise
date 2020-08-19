package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zngue/go_gin_praise/app/request"
	praise2 "github.com/zngue/go_gin_praise/app/service/praise"
	"github.com/zngue/go_tool/src/common/response"
)

func Praise(ctx *gin.Context)  {
	var praise request.PraiseRequest//接收参数处理
	if err:=ctx.BindJSON(&praise); err != nil {
		response.HttpFailWithParameter(ctx,err.Error())
		return
	}
	if praise.SingleNum==0 {
		praise.SingleNum=1
	}
	p:=new(praise2.Praise)
	var (
		res interface{}
		err error
	)
	if b ,err:=praise.Check();!b {
		response.HttpFailWithParameter(ctx,err.Error())
	}
	switch praise.ActionType {
	case 1://点赞
		if res,err=p.AddOrCancel(praise); err == nil {
			response.HttpOk(ctx,res)
			return
		}
		break
	case 2:
		if res ,err= p.GetNum(praise); err == nil {
			response.HttpOk(ctx,res)
			return
		}
		break
	case 3:
		if res ,err=p.IsPraise(praise); err == nil {
			response.HttpOk(ctx,res)
			return
		}
		break
	case 4:
		response.HttpOk(ctx,p.GetNumByIDs(praise))
		return
	case 5:
		response.HttpOk(ctx,p.IsPraiseByIDs(praise))
		return
	case 6:
		response.HttpOk(ctx,p.GetNumRank(praise))
		return
	case 7:
		response.HttpOk(ctx,p.GetNumByYearRank(praise))
		return
	case 8:
		response.HttpOk(ctx,p.GetNumByMouthRank(praise))
		return
	default:
		response.HttpFailWithParameter(ctx,"参数不完整")
		return
	}
	response.HttpOk(ctx,err.Error())
	return
}
