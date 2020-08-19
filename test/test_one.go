package test

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/zngue/go_tool/src/db"
	"github.com/zngue/go_tool/src/fun/time"
	"sync"
)


func main()  {

	db.InitDB()

	p:=Parise{
		From:    "user",
		ID:      101,
		Device:Device{
			UserID: "100",
			Uuid: "",
		},
		TypeID:  2,
		SignNum: 1,
	}
	p.KeySet()
	user:=p.PraiseAll()

	fmt.Println(user)

}

type Parise struct {
	From    string
	ID      int64
	SignNum int64  //单次点赞增加减少数量 默认1
	IsMouth bool //月排行
	IsYear bool //年排行
	TypeID  int64
	RedisPariseKey  //redis存储的key值
	Device  //设备信息表
}
type Device struct {
	UserID  string
	Uuid    string
	UnionID string
	UuidStatusKey string
	UnionIDStatusKey string
	UserIDStatusKey string
}
type RedisPariseKey struct {
	NumKey  string//所有点赞数量
	ZachKey string //所有点赞数量有序集合
	NumMouthKey string//年
	ZachMouthKey string
	NumYearKey string//月
	ZachYearKey string
}
func (p *Parise) KeySet() {
	year:=time.TimeToFormat(time.Time(),"2006")
	mouth:=time.TimeToFormat(time.Time(),"200601")
	p.NumKey = fmt.Sprintf("vote:NumKey:%s:%d:%d",p.From,p.TypeID,p.ID)//所有点赞
	p.ZachKey = fmt.Sprintf("vote:ZachKey:%s:%d:%d",p.From,p.TypeID,p.ID)//点赞集合

	p.NumYearKey = fmt.Sprintf("vote:NumYearKey:%s:%s:%d:%d",year,p.From,p.TypeID,p.ID)//月
	p.ZachYearKey = fmt.Sprintf("vote:ZachYearKey:%s:%s:%d:%d",year,p.From,p.TypeID,p.ID)//月点赞集合

	p.NumMouthKey = fmt.Sprintf("vote:NumMouthKey:%s:%s:%d:%d",mouth,p.From,p.TypeID,p.ID)//月
	p.ZachMouthKey = fmt.Sprintf("vote:ZachMouthKey:%s:%s:%d:%d",mouth,p.From,p.TypeID,p.ID)//月点赞集合

	//用户信息 或者设备信息
	p.UuidStatusKey = fmt.Sprintf("vote:UuidStatusKey:%s:%d:%s:%d",p.From,p.TypeID,p.Uuid,p.ID)//用户id
	p.UserIDStatusKey = fmt.Sprintf("vote:UserIDStatusKey:%s:%d:%s:%d",p.From,p.TypeID,p.UserID,p.ID)//设备信息
	p.UnionIDStatusKey = fmt.Sprintf("vote:UnionIDStatusKey:%s:%d:%s:%d",p.From,p.TypeID,p.UserID,p.ID)//设备信息
}

type PraiseStatusChange func(redis redis.Pipeliner,info *Parise)
//typeB true 可以点赞

func ( p *Parise) PariseStatusChangeAdd(info *PraiseALlInfo,wg *sync.WaitGroup, changeFn PraiseStatusChange) *PraiseALlInfo  {
	defer wg.Done()
	pipe:=db.RedisConn.TxPipeline()
	changeFn(pipe,p)
	if _,eerr:=pipe.Exec();eerr!=nil{
		info.StatusErr=eerr
		return info
	}
	return info
}
func (p *Parise) PariseStatusChangeDel(info *PraiseALlInfo,wg *sync.WaitGroup, changeFn PraiseStatusChange) *PraiseALlInfo {
	defer wg.Done()
	pipe:=db.RedisConn.TxPipeline()
	changeFn(pipe,p)
	if _,eerr:=pipe.Exec();eerr!=nil{
		info.StatusErr=eerr
		return info
	}
	return info

}
func (p *Parise) PariseStatus(typeB bool ,info *PraiseALlInfo,wg *sync.WaitGroup, change ...PraiseStatusChange) *PraiseALlInfo {
	defer wg.Done()
	pipe:=db.RedisConn.TxPipeline()
	if typeB {
		if p.Uuid!="" {
			pipe.SAdd(p.UuidStatusKey,p.ID)
		}
		if p.UserID!="" {
			pipe.SAdd(p.UserIDStatusKey,p.ID)
		}
		if p.UnionID!="" {
			pipe.SAdd(p.UnionIDStatusKey,p.ID)
		}
	}else{
		if p.Uuid!="" {
			pipe.SRem(p.UuidStatusKey,p.ID)
		}
		if p.UnionID!="" {
			pipe.SRem(p.UnionIDStatusKey,p.ID)
		}
		if p.UserID!="" {
			pipe.SRem(p.UserIDStatusKey,p.ID)
		}
	}
	if _,err:=pipe.Exec(); err != nil {
		pipe.Discard()
		info.StatusErr = err
		return info
	}
	return info
}

type PraiseALlInfo struct {
	Num int64
	PraiseErr error
	StatusErr error
	Err error
}

func  (p *Parise) PraiseAll() *PraiseALlInfo {
	var (
		ps PraiseALlInfo
		wg sync.WaitGroup
	)
	var typeB bool=false
	if !p.IsParise() { //
		typeB=true
	}
	wg.Add(2)
	go p.Praise(typeB,&ps,&wg)
	go p.PariseStatus(typeB,&ps,&wg)
	p.PariseStatusChangeAdd(&ps,&wg, func(redis redis.Pipeliner, info *Parise) {
		if  info.UnionID!=""{
			redis.SAdd(info.UnionIDStatusKey,p.ID)
		}
		if info.UserID!="" {
			redis.SAdd(info.UserIDStatusKey,p.ID)
		}
		if info.Uuid!="" {
			redis.SAdd(info.UuidStatusKey,p.ID)
		}
	})
	wg.Wait()

	if ps.PraiseErr==nil && ps.StatusErr==nil {
		return &ps
	}else {
		ps.Err=errors.New("操作失败")
		return &ps
	}
}


//点赞或者取消点赞
func (p *Parise) Praise(typeB bool,info *PraiseALlInfo,wg *sync.WaitGroup)  *PraiseALlInfo {
	defer wg.Done()
	nums,err:=p.LiKePraise(p.NumKey,p.ZachKey,typeB)//所有
	if p.IsMouth {
		p.LiKePraise(p.NumMouthKey,p.ZachMouthKey,typeB)//每月
	}
	if p.IsYear {
		p.LiKePraise(p.NumYearKey,p.ZachYearKey,typeB)//每年
	}
	if err!=nil {
		info.PraiseErr=err
		return info
	}
	info.Num=nums
	return info
}
func (p *Parise)IsParise() bool {
	var (
		b bool
		b1 bool
	)
	if p.Uuid!="" {
		b=p.IsSignParise(p.UuidStatusKey,p.ID)
	}
	if p.UserID!="" {
		b1=p.IsSignParise(p.UserIDStatusKey,p.ID)
	}
	fmt.Println(b)
	fmt.Println(b1)
	if b1||b {
		return true
	}else {
		return false
	}
}
func (p *Parise) IsSignParise(key string,vals int64) bool {
	val,_:=db.RedisConn.SIsMember(key,vals).Result()
	if val {
		return true
	}else {
		return false
	}
}

//rediskey
func (p *Parise) LiKePraise(NumKey,ZachKey string,typeB bool) (int64,error) {
	pipe:=db.RedisConn
	var err error
	var nums int64

	if typeB {
		nums,err=pipe.IncrBy(NumKey,p.SignNum).Result()
	}else {
		nums,err=pipe.DecrBy(NumKey,p.SignNum).Result()
	}
	if err != nil {
		return 0,err
	}else {
		_,err:=pipe.ZAdd(ZachKey,redis.Z{Score: float64(nums),Member: p.ID}).Result()
		if err!=nil {
			return 0,err
		}
		return nums,nil
	}
}














