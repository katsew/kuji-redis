package kuji_redis

import (
	"github.com/go-redis/redis"
	"github.com/katsew/kuji"
)

type ShuffleStrategy struct {
	client *redis.Client
	util *kuji.Util
}

func NewShuffleStrategy(o *redis.Options) (ShuffleStrategy) {
	return ShuffleStrategy{
		client: redis.NewClient(o),
		util: &kuji.Util{},
	}
}

func (s ShuffleStrategy) PickOneByKey(key string) (string, error) {
	length := s.client.LLen(key)
	pickingIndex := s.util.RandomNumberFromRange(0, length.Val())
	res := s.client.LIndex(key, pickingIndex)
	return res.Result()
}

func (s ShuffleStrategy) PickOneByKeyAndIndex(key string, index int64) (string, error) {
	res := s.client.LIndex(key, index)
	return res.Result()
}

func (s ShuffleStrategy) PickAndDeleteOneByKey(key string) (string, error) {
	res := s.client.LPop(key)
	return res.Result()
}

func (s ShuffleStrategy) RegisterCandidatesWithKey(key string, c []kuji.KujiCandidate) (int64, error) {
	var res *redis.IntCmd
	for _, v := range s.util.Shuffle(s.util.SpreadCandidates(c)) {
		res = s.client.RPush(key, v)
	}
	return res.Result()
}

func (s ShuffleStrategy) Len(key string) (int64, error) {
	res := s.client.LLen(key)
	return res.Result()
}

func (s ShuffleStrategy) List(key string) ([]string, error) {
	res := s.client.LRange(key, 0, -1)
	return res.Result()
}
