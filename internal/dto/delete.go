package dto

type DeleteForm struct {
	ID         uint64 `form:"id" validate:"required,number,min=1"`
	Project_id uint64 `form:"project_id" validate:"required,number,min=1"`
}

type DeleteResponse struct {
	ID         uint64 `json:"id"`
	Project_id uint64 `json:"project_id"`
	Removed    bool   `json:"removed"`
}
