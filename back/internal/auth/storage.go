package auth

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/models"
)

type Storage interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	FindByID(ctx context.Context, userID int) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
