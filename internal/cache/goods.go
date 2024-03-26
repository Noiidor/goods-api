package cache

import "goods-api/internal/models"

type GoodsCache interface {
	CacheSet(goods []models.Good)
	TryGetSetMembers(offset, limit int) []models.Good
	DeleteMember(id uint64)
	UpdateMember(good models.Good)
	AddMember(good models.Good)
	DeleteSet()
}
