package models

type Task struct {
	ID          int    `json:"id" validate:"omitempty"`
	ColumnID    int    `json:"column_id" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}
