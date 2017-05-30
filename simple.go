package kuji_redis

import (
	"github.com/go-redis/redis"
	"github.com/katsew/kuji"
)

type SimpleStrategy struct {
	client *redis.Client
	util   *kuji.Util
}

func NewSimpleStrategy(o *redis.Options) SimpleStrategy {
	return SimpleStrategy{
		client: redis.NewClient(o),
		util:   &kuji.Util{},
	}
}

func (s SimpleStrategy) PickOneByKey(key string) (string, error) {
	length := s.client.LLen(key)
	pickingIndex := s.util.RandomNumberFromRange(0, length.Val())
	res := s.client.LIndex(key, pickingIndex)
	return res.Result()
}

func (s SimpleStrategy) PickOneByKeyAndIndex(key string, index int64) (string, error) {
	res := s.client.LIndex(key, index)
	return res.Result()
}

func (s SimpleStrategy) PickAndDeleteOneByKey(key string) (string, error) {
	length := s.client.LLen(key)
	pickingIndex := s.util.RandomNumberFromRange(0, length.Val())
	res := s.client.LIndex(key, pickingIndex)
	s.client.LRem(key, 1, pickingIndex)
	return res.Result()
}

func (s SimpleStrategy) RegisterCandidatesWithKey(key string, candidates []kuji.KujiCandidate) (int64, error) {
	var res *redis.IntCmd
	for _, v := range s.util.SpreadCandidates(candidates) {
		res = s.client.LPush(key, v)
	}
	return res.Result()
}

func (s SimpleStrategy) Len(key string) (int64, error) {
	res := s.client.LLen(key)
	return res.Result()
}

func (s SimpleStrategy) List(key string) ([]string, error) {
	res := s.client.LRange(key, 0, -1)
	return res.Result()
}
