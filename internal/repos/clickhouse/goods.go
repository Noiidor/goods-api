package repos

import (
	"fmt"
	"goods-api/internal/models"
	"goods-api/internal/repos"

	"github.com/jmoiron/sqlx"
)

type clickhouseAnalytics struct {
	db *sqlx.DB
}

func NewGoodsAnalytics(db *sqlx.DB) repos.GoodsAnalytics {
	return &clickhouseAnalytics{db: db}
}

func (c clickhouseAnalytics) Insert(good models.Good) error {
	query := `
	INSERT INTO goods(id, project_id, name, description, priority, removed, event_time) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := c.db.Exec(query, good.ID, good.Project_id, good.Name, good.Description, good.Priority, good.Removed, good.Created_at)
	if err != nil {
		return err
	}

	return nil
}

func (c clickhouseAnalytics) InsertBatch(goods []models.Good) error {
	query := `
	INSERT INTO goods(id, project_id, name, description, priority, removed, event_time) 
	VALUES
	`

	for _, v := range goods { // Сюда по-хорошему прикрутить sql санитайзер
		query += fmt.Sprintf("(%v, %v, '%v', '%v', %v, %v, '%v')\n", v.ID, v.Project_id, v.Name, v.Description, v.Priority, v.Removed, v.Created_at.Format("2006-01-02 15:04:05"))
	} // да, выглядит убого

	_, err := c.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
