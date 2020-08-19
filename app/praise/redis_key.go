package praise

type RedisPraiseKey struct {
	NumKey  string//所有点赞数量
	ZachKey string //所有点赞数量有序集合
	NumMouthKey string//年
	ZachMouthKey string
	NumYearKey string//月
	ZachYearKey string
}
