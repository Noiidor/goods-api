package redis

import (
	"context"
	"goods-api/internal/cache"
	"goods-api/internal/errors"
	"goods-api/internal/models"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

const (
	GOODS_SET_NAME = "goods"
	GOODS_EXPIRE   = 60
)

type goodsCache struct {
	redis *redis.Client
}

func New(conn *redis.Client) cache.GoodsCache {
	return &goodsCache{redis: conn}
}

func (c goodsCache) CacheSet(goods []models.Good) {
	ctx := context.Background()

	pipe := c.redis.Pipeline()
	defer pipe.Exec(ctx)

	for _, v := range goods {
		json, err := jsoniter.MarshalToString(v)
		if err != nil {
			errors.AsyncErrors <- err
		}
		_, err = pipe.ZAdd(ctx, GOODS_SET_NAME, redis.Z{
			Score:  float64(v.ID),
			Member: json,
		}).Result()
		if err != nil {
			errors.AsyncErrors <- err
		}

	}

	_, err := pipe.Expire(ctx, GOODS_SET_NAME, time.Duration(time.Second*GOODS_EXPIRE)).Result()
	if err != nil {
		errors.AsyncErrors <- err
	}
}

func (c goodsCache) TryGetSetMembers(offset, limit int) []models.Good {
	ctx := context.Background()

	exist, err := c.redis.Exists(ctx, GOODS_SET_NAME).Result()
	if err != nil {
		errors.AsyncErrors <- err
		return nil
	}
	if exist != 1 {

		return nil
	}
	goods := make([]models.Good, 0, limit)

	if limit == 0 {
		return goods
	}

	goodJsons, err := c.redis.ZRange(ctx, GOODS_SET_NAME, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		errors.AsyncErrors <- err
		return nil
	}

	for _, v := range goodJsons {
		var good models.Good
		jsoniter.Unmarshal([]byte(v), &good)
		goods = append(goods, good)
	}

	return goods
}

func (c goodsCache) DeleteMember(id uint64) {
	ctx := context.Background()

	_, err := c.redis.ZRemRangeByScore(ctx, GOODS_SET_NAME, strconv.FormatUint(id, 10), strconv.FormatUint(id, 10)).Result()
	if err != nil {
		errors.AsyncErrors <- err
	}
}

func (c goodsCache) UpdateMember(good models.Good) {
	ctx := context.Background()

	exist, err := c.redis.Exists(ctx, GOODS_SET_NAME).Result()
	if err != nil {
		errors.AsyncErrors <- err
	}
	if exist != 1 {
		return
	}

	pipe := c.redis.Pipeline()
	defer pipe.Exec(ctx)

	_, err = pipe.ZRemRangeByScore(ctx, GOODS_SET_NAME, strconv.FormatUint(good.ID, 10), strconv.FormatUint(good.ID, 10)).Result()
	if err != nil {
		errors.AsyncErrors <- err
	}

	json, err := jsoniter.MarshalToString(good)
	if err != nil {
		errors.AsyncErrors <- err
	}

	_, err = pipe.ZAdd(ctx, GOODS_SET_NAME, redis.Z{
		Score:  float64(good.ID),
		Member: json,
	}).Result()
	if err != nil {
		errors.AsyncErrors <- err
	}
}

func (c goodsCache) AddMember(good models.Good) {
	ctx := context.Background()

	exist, err := c.redis.Exists(ctx, GOODS_SET_NAME).Result()
	if err != nil {
		errors.AsyncErrors <- err
	}
	if exist != 1 {
		return
	}

	json, err := jsoniter.MarshalToString(good)
	if err != nil {
		errors.AsyncErrors <- err
	}

	_, err = c.redis.ZAdd(ctx, GOODS_SET_NAME, redis.Z{
		Score:  float64(good.ID),
		Member: json,
	}).Result()
	if err != nil {
		errors.AsyncErrors <- err
	}
}

func (c goodsCache) DeleteSet() {
	ctx := context.Background()

	_, err := c.redis.Del(ctx, GOODS_SET_NAME).Result()
	if err != nil {
		errors.AsyncErrors <- err
	}
}
