package session

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/models"
)

type UseCase interface {
	CreateSession(ctx context.Context, session *models.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
