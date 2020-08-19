package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zngue/go_gin_praise/app/api"
)

func PraiseRouter(r *gin.RouterGroup)  {
	p:=r.Group("praise")
	{
		p.GET("praise",api.Praise)
	}
}
