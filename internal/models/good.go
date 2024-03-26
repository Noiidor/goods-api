package models

import "time"

type Good struct {
	ID          uint64    `db:"id" json:"id"`
	Project_id  uint64    `db:"project_id" json:"project_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Priority    int       `db:"priority" json:"priority"`
	Removed     bool      `db:"removed" json:"removed"`
	Created_at  time.Time `db:"created_at" json:"created_at"`
}
