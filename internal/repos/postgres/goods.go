package repos

import (
	"fmt"
	"goods-api/internal/models"
	"goods-api/internal/repos"
	"strings"

	"github.com/jmoiron/sqlx"
)

type goodsRepo struct {
	db *sqlx.DB
}

func NewGoodsRepo(db *sqlx.DB) repos.GoodsRepo {
	return &goodsRepo{db: db}
}

func (r goodsRepo) Exists(id, projectId uint64) (bool, error) {
	query := `
	SELECT EXISTS
		(SELECT 1 FROM goods 
			WHERE id = $1
			AND project_id = $2);
	`

	var exists bool
	err := r.db.QueryRowx(query, id, projectId).Scan(&exists)

	return exists, err
}

func (r goodsRepo) Insert(projectId uint64, name string) (models.Good, error) {
	var good models.Good

	query := `
	INSERT INTO goods (project_id, name) 
	VALUES ($1, $2) 
	RETURNING *
	`
	transaction, err := r.db.Beginx() // Транзакция тут оверкилл, всё равно только один запрос, нечего оборачивать,
	if err != nil {                   // разве что под расширение, тогда сгодится
		return models.Good{}, err
	}

	err = transaction.QueryRowx(query, projectId, name).StructScan(&good)
	if err != nil {
		transaction.Rollback()
		return models.Good{}, err
	}

	err = transaction.Commit()
	if err != nil {
		return models.Good{}, err
	}

	return good, nil
}

func (r goodsRepo) Update(id, projectId uint64, good models.Good) (models.Good, error) {
	var updatedGood models.Good

	query := ` 
		UPDATE goods
		SET
	` // UPDATE + транзакция и так ставит блокировку, тут ничего дополнительного не надо

	columns := make([]string, 0) // Около-универсальный способ апдейта любых полей в структуре, без рефлексии
	params := make([]any, 0)

	params = append(params, id)
	params = append(params, projectId)

	if good.Name != "" {
		columns = append(columns, fmt.Sprintf("name=$%v", len(params)+1))
		params = append(params, good.Name)
	}
	if good.Description != "" {
		columns = append(columns, fmt.Sprintf("description=$%v", len(params)+1))
		params = append(params, good.Description)
	}
	// Надо обновить новое поле - добавить один if, по-моему неплохо

	query += strings.Join(columns, ", ")

	query += `
	WHERE id = $1 AND project_id = $2
	RETURNING *
	`

	transaction, err := r.db.Beginx()
	if err != nil {
		return models.Good{}, err
	}

	err = transaction.QueryRowx(query, params...).StructScan(&updatedGood)
	if err != nil {
		transaction.Rollback()
		return models.Good{}, err
	}

	err = transaction.Commit()
	if err != nil {
		return models.Good{}, err
	}

	return updatedGood, nil
}

func (r goodsRepo) Delete(id, projectId uint64) (models.Good, error) {
	var good models.Good

	query := `
	UPDATE goods
	SET removed = true
	WHERE id = $1 AND project_id = $2
	RETURNING *;
	`

	transaction, err := r.db.Beginx()
	if err != nil {
		return models.Good{}, err
	}

	err = transaction.QueryRowx(query, id, projectId).StructScan(&good)
	if err != nil {
		transaction.Rollback()
		return models.Good{}, err
	}

	err = transaction.Commit()
	if err != nil {
		return models.Good{}, err
	}

	return good, nil
}

func (r goodsRepo) GetListOffset(limit, offset int) ([]models.Good, error) {
	goods := make([]models.Good, 0, limit)

	query := `
	SELECT *
	FROM goods
	ORDER BY id
	OFFSET $1
	LIMIT $2
	`
	err := r.db.Select(&goods, query, offset, limit)
	if err != nil {
		return goods, err
	}

	return goods, nil
}

func (r goodsRepo) Reprioritize(id, projectId uint64, newPriority int) ([]models.Good, error) {

	query := `
	WITH counting AS (
		SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS num
		FROM goods
		WHERE id >= $2 AND project_id = $3
	)
	UPDATE goods g
	SET priority = $1 + c.num - 1
	FROM counting c
	WHERE g.id = c.id
	RETURNING g.*;
	`

	transaction, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	rows, err := transaction.Queryx(query, newPriority, id, projectId)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	goods := make([]models.Good, 0)

	for rows.Next() {
		var good models.Good
		err = rows.StructScan(&good)
		goods = append(goods, good)
	}
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	if err = rows.Err(); err != nil {
		transaction.Rollback()
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return goods, nil

}

func (r goodsRepo) GetAll() ([]models.Good, error) {
	goods := make([]models.Good, 0, 1000)

	query := `
	SELECT *
	FROM goods
	ORDER BY id
	`
	err := r.db.Select(&goods, query)
	if err != nil {
		return goods, err
	}

	return goods, nil
}
