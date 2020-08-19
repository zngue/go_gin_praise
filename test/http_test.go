package test

import (
	"github.com/gin-gonic/gin"
	"github.com/zngue/go_tool/src/db"
	"github.com/zngue/go_tool/src/gin_run"
	"testing"
)

func TestHttp(t *testing.T) {


	gin_run.GinRun(func(group *gin.RouterGroup) {

	})
	defer db.ConnClose()


}
