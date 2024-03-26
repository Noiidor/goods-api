package dto

type PriorityForm struct {
	ID         uint64 `form:"id" validate:"required,number,min=1"`
	Project_id uint64 `form:"project_id" validate:"required,number,min=1"`
}

type PriorityPayload struct {
	New_priority int `json:"new_priority" validate:"required,number"`
}

type PriorityResponse struct {
	Priorities []Priority `json:"priorities"`
}

type Priority struct {
	ID       uint64 `json:"id"`
	Priority int    `json:"priority"`
}
