package repos

import "goods-api/internal/models"

type ProjectsRepo interface {
	Insert(name string) (models.Project, error)
	Exists(projectId uint64) (bool, error)
}
