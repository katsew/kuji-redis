package kuji_redis

import (
	"github.com/go-redis/redis"
	"github.com/katsew/kuji"
	"fmt"
)

type KujiRedisStrategy struct {
	client *redis.Client
	util *kuji.Util
}

func NewKujiRedisStrategy(o *redis.Options) (KujiRedisStrategy) {
	rs := KujiRedisStrategy{
		client: redis.NewClient(o),
		util: &kuji.Util{},
	}
	return rs
}

func (rs KujiRedisStrategy) PickOneByKey(s string) (string, error) {
	length := rs.client.LLen(s)
	random := rs.util.RandomNumberFromRange(0, length.Val())
	fmt.Printf("Output rand: %d\n", random)
	res := rs.client.LIndex(s, random)
	return res.Result()
}

func (rs KujiRedisStrategy) RegisterCandidatesWithKey(s string, c []kuji.KujiCandidate) (int64, error) {
	var res *redis.IntCmd
	for _, v := range rs.util.SpreadCandidates(c) {
		res = rs.client.LPush(s, v)
	}
	return res.Result()
}