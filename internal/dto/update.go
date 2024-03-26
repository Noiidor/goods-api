package dto

import "goods-api/internal/models"

type UpdateForm struct {
	ID         uint64 `form:"id" validate:"required,number,min=1"`
	Project_id uint64 `form:"project_id" validate:"required,number,min=1"`
}

type UpdatePayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateResponse struct {
	Good models.Good `json:"good"`
}
