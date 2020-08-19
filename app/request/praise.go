package request

import (
	"errors"
	"github.com/zngue/go_tool/src/common/request"
	"reflect"
)

type PraiseRequest struct {
	ID   int64 `json:"id" form:"id"`//投票id
	Sort bool 	`json:"sort" form:"sort"` //排序  true  倒叙 高->低 false 顺序 低->高
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
	request.Page
}

//设备信息
type Device struct {
	UuID string  `json:"uuid" form:"uuid"`
	UserID int `json:"user_id" form:"user_id"`
	UnionID string `json:"union_id" form:"union_id"`
}

func (p *PraiseRequest) Check()  (bool,error) {
	var s []interface{}
	if p.From=="" || p.TypeID==0   {
		return  false,errors.New("from   and type_id should")
	}
	if !p.IDArrCheck() {
		return false,errors.New(" id or id_arr should ")
	}
	if p.ActionType==2 || p.ActionType==4 || p.ActionType==6 || p.ActionType==7 || p.ActionType==8 {
		return true,nil
	}
	elem := reflect.ValueOf(&p.Device).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		if elem.Field(i).Interface()!=nil {
			s= append(s, elem.Field(i).Interface())
		}
	}
	if len(s)>0 {
		return true,nil
	}else{
		return false,errors.New("uuid  or user_id or union_id should")
	}
}

func (p *PraiseRequest) IDArrCheck() bool  {
	if p.ActionType==4 || p.ActionType==5 {
		if len(p.IDArr)==0 {
			return false
		}
	}else if p.ActionType==1 || p.ActionType==2 || p.ActionType==3 {
		if p.ID==0 {
			return  false
		}
	}
	return true
}