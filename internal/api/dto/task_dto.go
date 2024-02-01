package dto

import "github.com/gofrs/uuid"

type CreateTaskRequest struct {
	Name      string `json:"name" validate:"required,max=256"`
	Image     string `json:"image" validate:"required,max=256"`
	Namespace string `json:"namespace" validate:"required,max=64"`
	Runtime   string `json:"runtime" validate:"required,max=32"`
	Script    string `json:"script" validate:"required,max=1000"`
}

type CreateTaskResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Namespace string    `json:"namespace"`
	Runtime   string    `json:"runtime"`
	Status    string    `json:"status"`
	Script    string    `json:"script"`
	Result    string    `json:"result"`
}
