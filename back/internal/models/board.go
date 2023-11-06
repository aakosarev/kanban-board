package models

type Board struct {
	Columns []*Col `json:"columns"`
}

type Col struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Tasks []*T   `json:"tasks"`
}

type T struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}
