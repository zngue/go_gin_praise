# go_gin_praise 点赞
##操作文档


## 使用教程
```go

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zngue/go_gin_praise/router"
	"github.com/zngue/go_tool/src/gin_run"
)

func main()  {

	gin_run.GinRun(func(group *gin.RouterGroup) {
		router.PraiseRouter(group)
	}, func(db *gorm.DB) {


	})
}

```

###部分参数json参数如下所示

```json
{
    "id":7, int64//点赞id
    "type_id":1, int64//点赞类型
    "from":"user",  string//点赞来源
    "single_num":1, //int64 单次点赞增加数量 默认1或者减少
    "action_type":1,
        //1 点赞或者取消点赞
    	//2 获取点赞数量
    	//3 判断是否点赞
    	//4 批量获取点赞数量
    	//5 批量获取点赞状态
    	//6 获取点赞排行榜
    	//7 获取点赞年排行榜
    	//8 获取点赞月排行榜
    "device":{
        "uuid":"100004", //string  设备信息
        "user_id" : 10, //int64 用户id
        "union_id": "asdkaldsjl" //unionid string
    },
    "page":1,//分页数据
    "page_size":2 //分页数据
}

```
```go
type PraiseRequest struct {
	ID   int64 `json:"id" form:"id"`//投票id
	Sort bool 	`json:"sort" form:"sort"` //排序  true  倒叙 高->低 false 顺序 低->高 获取排行榜的时候使用
	IDArr []int64  `json:"id_arr" form:"id_arr"`
	TypeID int64 `json:"type_id" form:"type_id"`//类型
	From string `json:"from" form:"from"` //来源  比如 app  比如 news
	ALLLimit int `json:"all_limit" form:"user_limit"`//总限制数量 1
	DayLimit int  `json:"day_limit" form:"day_limit"`//每日限制数量 0 表示不处理
	ExpireAtKey int  `json:"expire_at" form:"expire_at"` // key值有效期 默认永久 -1
	SingleNum int64  `json:"single_num" form:"single_num"`//单次投票增加或者减少数量 默认1
	IsYear bool `json:"is_year" form:"is_year"`
	IsMouth bool `json:"is_mouth" form:"is_mouth"`
	Device Device `json:"device" form:"device"`//设备信息
	ActionType int `json:"action_type" form:"action_type"`//操作类型
	//1 点赞或者取消点赞
	//2 获取点赞数量
	//3 判断是否点赞
	//4 批量获取点赞数量
	//5 批量获取点赞状态
	//6 获取点赞排行榜
	//7 获取点赞年排行榜
	//8 获取点赞月排行榜
	request.Page //分页参数  page   page_size
}

//设备信息
type Device struct {
	UuID string  `json:"uuid" form:"uuid"`
	UserID int `json:"user_id" form:"user_id"`
	UnionID string `json:"union_id" form:"union_id"`
}
```
