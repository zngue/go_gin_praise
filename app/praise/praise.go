package praise

import (
	"github.com/go-redis/redis"
	"github.com/zngue/go_tool/src/db"
	"github.com/zngue/go_tool/src/fun/time"
	"strconv"
	"sync"
	"reflect"
	"fmt"
)

type Praise struct {
	From           string
	ID             int64
	SignNum        int64  //单次点赞增加减少数量 默认1
	IsMouth        bool //月排行
	IsYear         bool //年排行
	TypeID         int64
	RedisPraiseKey //redis存储的key值
	Device         //设备信息表
	RedisPage		//
	DeviceStatus
}

func (p *Praise)  PageInit()  {
	if p.Page==0 {
		p.Page=1
	}
	if p.PageSize==0 {
		p.PageSize=15
	}
	p.Start = (p.Page-1)*p.PageSize
	p.Stop = (p.Page*p.PageSize)-1
}
func (p *Praise) Init() {
	year:=time.TimeToFormat(time.Time(),"2006")
	mouth:=time.TimeToFormat(time.Time(),"200601")
	p.NumKey = fmt.Sprintf("vote:NumKey:%s:%d:%d",p.From,p.TypeID,p.ID)//所有点赞
	p.ZachKey = fmt.Sprintf("vote:ZachKey:%s:%d",p.From,p.TypeID)//点赞集合

	p.NumYearKey = fmt.Sprintf("vote:NumYearKey:%s:%s:%d:%d",year,p.From,p.TypeID,p.ID)//月
	p.ZachYearKey = fmt.Sprintf("vote:ZachYearKey:%s:%s:%d",year,p.From,p.TypeID)//集合数据

	p.NumMouthKey = fmt.Sprintf("vote:NumMouthKey:%s:%s:%d:%d",mouth,p.From,p.TypeID,p.ID)//月
	p.ZachMouthKey = fmt.Sprintf("vote:ZachMouthKey:%s:%s:%d",mouth,p.From,p.TypeID)//月点赞集合

	//用户信息 或者设备信息
	if p.Uuid!="" {
		p.UuidStatusKey = fmt.Sprintf("vote:UuidStatusKey:%s:%d:%s",p.From,p.TypeID,p.Uuid)//用户id
	}
	if p.UserID!=""  && p.UserID!="0" {
		p.UserIDStatusKey = fmt.Sprintf("vote:UserIDStatusKey:%s:%d:%s",p.From,p.TypeID,p.UserID)//设备信息
	}
	if p.UnionID!="" {
		p.UnionIDStatusKey = fmt.Sprintf("vote:UnionIDStatusKey:%s:%d:%s",p.From,p.TypeID,p.UnionID)//设备信息
	}
	p.PageInit()
}
//批量获取点赞状态
func ( p *Praise ) GetIDArrStatus(ids []int64 ) map[int64]bool {
	m:=make(map[int64]bool)
	stars :=p.Rolex()
	for _ ,id :=range ids{
		p.ID=id
		b, _ := p.GetStatusDevice(stars...)
		m[id]=b
	}
	return m
}
func ( p *Praise) IsPraise() (bool,error) {
	stars:=p.Rolex()
	return p.GetStatusDevice(stars...)
}
//批量获取点赞数量
func (p *Praise) GetIDArrNum(ids []int64) map[int64]string  {
	m := make(map[int64]string)
	p.GetNum()
	for _ ,id :=range ids{
		p.ID = id
		p.Init()
		if num ,err:=p.GetNum(); err == nil {
			m[id]=num
		}else {
			m[id]="0"
		}
	}
	return m
}
//排序
func (p *Praise) GetRange(sort bool,key string)  *RankList  {
	var (
		Count int64
		Err   error
		List  []redis.Z
	)
	Count,Err=db.RedisConn.ZCard(key).Result()
	if Err!=nil {
		return &RankList{
			Err: Err,
		}
	}
	if sort {
		List,Err=db.RedisConn.ZRevRangeWithScores(key,p.Start,p.Stop).Result()
	}else{
		List,Err=db.RedisConn.ZRangeWithScores(key,p.Start,p.Stop).Result()
	}
	return &RankList{
		List: List,
		Count: Count,
		Err: Err,
	}
}
func (p *Praise) Rolex() []string {
	var s []string
	elem := reflect.ValueOf(&p.DeviceStatus).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		if elem.Field(i).String()!="" {
			s= append(s, elem.Field(i).String())
		}
	}
	return s
}
func (p *Praise) GetMouthRange()  {
	nus,_:=db.RedisConn.ZCard(p.ZachMouthKey).Result()
	db.RedisConn.ZIncrBy(p.ZachMouthKey,float64(1),strconv.Itoa(int(p.ID)))
	s,_:=db.RedisConn.ZRevRange(p.ZachMouthKey,p.Start,p.Stop).Result()
	a,_:=db.RedisConn.ZRevRangeWithScores(p.ZachMouthKey,p.Start,p.Stop).Result()
	//db.RedisConn.ZRan
	fmt.Println(p.ZachMouthKey)
	fmt.Println(nus)
	fmt.Println(s)
	for _ , val:=range a{
		fmt.Println(val.Member)
	}
	fmt.Println(len(a))
}
func (p *Praise) Praise() (int64,error) {
	var (
		num int64
		Err error
	)
	stirs :=p.Rolex()
	b,err := p.GetStatusDevice(stirs...)
	//b = false
	if err != nil   {
		return 0, err
	}
	var wg sync.WaitGroup
	if b {//已经点赞  现在取消点赞
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.StatusDel(stirs...)
		}()
		if p.IsYear {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.DelPraise(p.NumYearKey,p.ZachYearKey)
			}()
		}
		if  p.IsMouth {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.DelPraise(p.NumMouthKey,p.ZachMouthKey)
			}()
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			num,Err=p.DelPraise(p.NumKey,p.ZachKey)
		}()
		wg.Wait()
		return num,Err
	}else{//点赞
		if p.IsYear {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.AddPraise(p.NumYearKey,p.ZachYearKey)
			}()
		}
		if p.IsMouth {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.AddPraise(p.NumMouthKey,p.ZachMouthKey)
			}()
		}
		wg.Add(2)
		go func() {
			defer wg.Done()
			p.StatusAdd(stirs...)
		}()
		go func() {
			defer wg.Done()
			num,Err = p.AddPraise(p.NumKey,p.ZachKey)
		}()
		wg.Wait()
		return num,Err
	}
}
func (p *Praise) GetNum() (string,error) {
	return db.RedisConn.Get(p.NumKey).Result()
}
//获取设备点赞状态
func (p *Praise) GetStatusDevice(s ...string) (bool,error) {
	if len(s)> 0{
		for _ ,val :=range s{
			if b,err:=db.RedisConn.SIsMember(val,p.ID).Result(); err != nil {
				return false, err
			}else{
				if b {
					return true,nil
				}
			}
		}
	}
	return false,nil
}
//添加点赞key
func (p *Praise) StatusAdd( fn ...string)  bool {
	redis := db.RedisConn
	var Err []error
	if len(fn)>0 {
		for _ ,val :=range fn{
			_,Errs:=redis.SAdd(val,p.ID).Result()
			if Errs!=nil {
				Err= append(Err, Errs)
			}
		}
	}
	if len(Err)>0 {
		return false
	}
	return true
}
//删除点赞key
func (p *Praise) StatusDel(fn ...string)  bool {
	redis := db.RedisConn
	var Err []error
	if len(fn)>0 {
		for _ ,val :=range fn{
			_,Errs:=redis.SRem(val,p.ID).Result()
			if Errs!=nil {
				Err= append(Err, Errs)
			}
		}
	}
	if len(Err)>0 {
		return false
	}
	return true
}
//取消点赞
func (p *Praise) DelPraise(NumKey,ZachKey string) (int64,error)  {
	num,err:= db.RedisConn.DecrBy(NumKey,p.SignNum).Result()
	if err!=nil {
		return 0, err
	}
	p.AddZany(ZachKey)
	return num, nil
}
//点赞
func (p *Praise) AddPraise(NumKey ,ZachKey string) (int64,error){
	var (
		num int64
		err error
	)
	num,err= db.RedisConn.IncrBy(NumKey,p.SignNum).Result()
	p.AddZany(ZachKey)
	if err!=nil {
		return 0, err
	}
	return num, nil
}
//将点赞数据写入集合
func (p *Praise) AddZany(ZachKey string)  {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	if numbs,err:=p.GetNum(); err == nil {
		if nurses,errs:= strconv.Atoi(numbs); errs==nil{
			db.RedisConn.ZAdd(ZachKey,redis.Z{Score: float64(nurses),Member: p.ID}).Result()
		}
	}
}
