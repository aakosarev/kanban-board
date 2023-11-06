package usecase

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/models"
	"github.com/aakosarev/kanban-board/back/internal/session"
)

type sessionUC struct {
	sessionStorage session.Storage
	cfg            *config.Config
}

func NewSessionUseCase(sessionStorage session.Storage, cfg *config.Config) session.UseCase {
	return &sessionUC{sessionStorage: sessionStorage, cfg: cfg}
}

func (u *sessionUC) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	return u.sessionStorage.CreateSession(ctx, session, expire)
}

func (u *sessionUC) DeleteByID(ctx context.Context, sessionID string) error {
	return u.sessionStorage.DeleteByID(ctx, sessionID)
}

func (u *sessionUC) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	return u.sessionStorage.GetSessionByID(ctx, sessionID)
}
