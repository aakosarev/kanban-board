package models

type Column struct {
	ID     int    `json:"id" validate:"omitempty"`
	UserID int    `json:"user_id" validate:"required"`
	Name   string `json:"name" validate:"omitempty"`
}
