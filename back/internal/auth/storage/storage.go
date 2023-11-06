package storage

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/auth"
	"github.com/aakosarev/kanban-board/back/internal/models"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthStorage struct {
	log    logger.Logger
	client *pgxpool.Pool
}

func NewAuthStorage(log logger.Logger, client *pgxpool.Pool) auth.Storage {
	return &AuthStorage{
		log:    log,
		client: client,
	}
}

func (s *AuthStorage) Register(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO "user"(email, password)
		VALUES ($1, $2)
		RETURNING *;
	`

	u := &models.User{}

	if err := s.client.QueryRow(ctx, query, user.Email, user.Password).Scan(&u.ID, &u.Email, &u.Password); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *AuthStorage) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT *
		FROM "user"
		WHERE email = $1;
	`

	u := &models.User{}

	if err := s.client.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *AuthStorage) FindByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT *
		FROM "user"
		WHERE id = $1;
	`

	u := &models.User{}

	if err := s.client.QueryRow(ctx, query, userID).Scan(&u.ID, &u.Email, &u.Password); err != nil {
		return nil, err
	}

	return u, nil
}
