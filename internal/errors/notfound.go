package errors

import "fmt"

type GoodNotFoundError struct {
	ID, Project_id uint64
}

func (e *GoodNotFoundError) Error() string {
	return fmt.Sprintf("good with id: %v, project_id: %v not found", e.ID, e.Project_id)
}

type ProjectNotFoundError struct {
	Project_id uint64
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project with id: %v not found", e.Project_id)
}
