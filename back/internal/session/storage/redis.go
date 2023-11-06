package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/models"
	"github.com/aakosarev/kanban-board/back/internal/session"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	basePrefix = "api-session"
)

type sessionStorage struct {
	redisClient *redis.Client
	basePrefix  string
	cfg         *config.Config
}

func NewSessionStorage(redisClient *redis.Client, cfg *config.Config) session.Storage {
	return &sessionStorage{redisClient: redisClient, basePrefix: basePrefix, cfg: cfg}
}

func (s *sessionStorage) CreateSession(ctx context.Context, sess *models.Session, expire int) (string, error) {
	sess.SessionID = uuid.New().String()
	sessionKey := s.createKey(sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return "", errors.WithMessage(err, "sessionStorage.CreateSession.json.Marshal")
	}
	if err = s.redisClient.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.Wrap(err, "sessionStorage.CreateSession.redisClient.Set")
	}
	return sessionKey, nil
}

func (s *sessionStorage) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	sessBytes, err := s.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "sessionStorage.GetSessionByID.redisClient.Get")
	}

	sess := &models.Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, errors.Wrap(err, "sessionStorage.GetSessionByID.json.Unmarshal")
	}
	return sess, nil
}

func (s *sessionStorage) DeleteByID(ctx context.Context, sessionID string) error {
	if err := s.redisClient.Del(ctx, sessionID).Err(); err != nil {
		return errors.Wrap(err, "sessionStorage.DeleteByID")
	}
	return nil
}

func (s *sessionStorage) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.basePrefix, sessionID)
}
