package repos

import (
	"goods-api/internal/models"
)

type GoodsRepo interface {
	Insert(projectId uint64, name string) (models.Good, error)
	Update(id, projectId uint64, good models.Good) (models.Good, error)
	Delete(id, projectId uint64) (models.Good, error)
	GetListOffset(limit, offset int) ([]models.Good, error)
	GetAll() ([]models.Good, error)
	Exists(id, projectId uint64) (bool, error)
	Reprioritize(id, projectId uint64, newPriority int) ([]models.Good, error)
}

type GoodsAnalytics interface {
	Insert(good models.Good) error
	InsertBatch(goods []models.Good) error
}
