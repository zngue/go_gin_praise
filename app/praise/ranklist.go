package praise

import "github.com/go-redis/redis"

type RankList struct {
	List []redis.Z
	Err error
	Count int64
}
