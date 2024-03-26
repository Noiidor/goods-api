package dto

import "goods-api/internal/models"

type CreateForm struct {
	Project_id uint64 `form:"project_id" validate:"required,number,min=1"`
}

type CreatePayload struct {
	Name string `json:"name" validate:"required"`
}

type CreateResponse struct {
	Good models.Good `json:"good"`
}
