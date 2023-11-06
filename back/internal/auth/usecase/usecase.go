package usecase

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/auth"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/models"
	httpErrors "github.com/aakosarev/kanban-board/back/pkg/http_errors"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"github.com/pkg/errors"
)

type authUseCase struct {
	cfg         *config.Config
	authStorage auth.Storage
	log         logger.Logger
}

func NewAuthUseCase(cfg *config.Config, authStorage auth.Storage, log logger.Logger) auth.UseCase {
	return &authUseCase{cfg: cfg, authStorage: authStorage, log: log}
}

func (u *authUseCase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	existsUser, err := u.authStorage.FindByEmail(ctx, user.Email)
	if existsUser != nil || err == nil {
		return nil, errors.New(httpErrors.ErrEmailAlreadyExists)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, errors.Wrap(err, "authUseCase.Register.PrepareCreate")
	}

	createdUser, err := u.authStorage.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.SanitizePassword()

	return createdUser, nil
}

func (u *authUseCase) GetByID(ctx context.Context, userID int) (*models.User, error) {
	user, err := u.authStorage.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.SanitizePassword()

	return user, nil
}

func (u *authUseCase) Login(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser, err := u.authStorage.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, errors.Wrap(err, "authUseCase.Login.ComparePasswords")
	}

	foundUser.SanitizePassword()

	return foundUser, nil
}
