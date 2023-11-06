package usecase

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/kanban"
	"github.com/aakosarev/kanban-board/back/internal/models"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
)

type kanbanUseCase struct {
	cfg           *config.Config
	kanbanStorage kanban.Storage
	log           logger.Logger
}

func NewKanbanUseCase(cfg *config.Config, kanbanStorage kanban.Storage, log logger.Logger) kanban.UseCase {
	return &kanbanUseCase{cfg: cfg, kanbanStorage: kanbanStorage, log: log}
}

func (kuc *kanbanUseCase) CreateColumn(ctx context.Context, column *models.Column) (*models.Column, error) {
	createdColumn, err := kuc.kanbanStorage.CreateColumn(ctx, column)
	if err != nil {
		return nil, err
	}

	return createdColumn, nil
}

func (kuc *kanbanUseCase) DeleteColumn(ctx context.Context, id int) error {
	return kuc.kanbanStorage.DeleteColumn(ctx, id)
}

func (kuc *kanbanUseCase) ChangeNameColumn(ctx context.Context, column *models.Column) (*models.Column, error) {
	updatedColumn, err := kuc.kanbanStorage.ChangeNameColumn(ctx, column)
	if err != nil {
		return nil, err
	}

	return updatedColumn, nil
}

func (kuc *kanbanUseCase) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	createdTask, err := kuc.kanbanStorage.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (kuc *kanbanUseCase) DeleteTask(ctx context.Context, id int) error {
	return kuc.kanbanStorage.DeleteTask(ctx, id)
}

func (kuc *kanbanUseCase) ChangeDescriptionTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	updatedTask, err := kuc.kanbanStorage.ChangeDescriptionTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (kuc *kanbanUseCase) ChangeColumnIDTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	updatedTask, err := kuc.kanbanStorage.ChangeColumnIDTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (kuc *kanbanUseCase) GetKanbanBoardByUserID(ctx context.Context, userID int) (*models.Board, error) {
	board, err := kuc.kanbanStorage.GetKanbanBoardByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return board, nil
}
