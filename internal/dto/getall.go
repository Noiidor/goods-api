package dto

import "goods-api/internal/models"

type GetAllForm struct {
	Limit  int `form:"limit,default=10" validate:"number,min=0,max=100000"`
	Offset int `form:"offset" validate:"number,min=0"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type GetAllResponse struct {
	Meta  Meta          `json:"meta"`
	Goods []models.Good `json:"goods"`
}
