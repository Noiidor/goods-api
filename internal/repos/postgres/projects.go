package repos

import (
	"goods-api/internal/models"
	"goods-api/internal/repos"

	"github.com/jmoiron/sqlx"
)

type projectsRepo struct {
	db *sqlx.DB
}

func NewProjectsRepo(db *sqlx.DB) repos.ProjectsRepo {
	return &projectsRepo{db: db}
}

func (r projectsRepo) Insert(name string) (models.Project, error) {
	var project models.Project

	query := `
		INSERT INTO projects (name) 
		VALUES ($1) 
		RETURNING *
	`

	err := r.db.QueryRowx(query, name).StructScan(&project)
	if err != nil {
		return models.Project{}, err
	}

	return project, nil
}

func (r projectsRepo) Exists(projectId uint64) (bool, error) {
	query := `
	SELECT EXISTS
		(SELECT 1 FROM projects 
		 WHERE id = $1);
	`

	var exists bool
	err := r.db.QueryRowx(query, projectId).Scan(&exists)

	return exists, err
}
