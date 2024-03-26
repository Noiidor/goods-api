package models

import "time"

type Project struct {
	ID         uint64    `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Created_at time.Time `db:"created_at" json:"created_at"`
}
