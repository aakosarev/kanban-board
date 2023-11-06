package kanban

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/models"
)

type UseCase interface {
	CreateColumn(ctx context.Context, column *models.Column) (*models.Column, error)
	DeleteColumn(ctx context.Context, id int) error
	ChangeNameColumn(ctx context.Context, column *models.Column) (*models.Column, error)

	CreateTask(ctx context.Context, task *models.Task) (*models.Task, error)
	DeleteTask(ctx context.Context, id int) error
	ChangeDescriptionTask(ctx context.Context, task *models.Task) (*models.Task, error)
	ChangeColumnIDTask(ctx context.Context, task *models.Task) (*models.Task, error)

	GetKanbanBoardByUserID(ctx context.Context, userID int) (*models.Board, error)
}
