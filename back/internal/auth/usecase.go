package auth

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/models"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID int) (*models.User, error)
}
