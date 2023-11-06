package middleware

import (
	"github.com/aakosarev/kanban-board/back/internal/auth"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/session"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"github.com/labstack/echo/v4"
	"time"
)

type Manager struct {
	sessionUseCase session.UseCase
	authUseCase    auth.UseCase
	cfg            *config.Config
	origins        []string
	logger         logger.Logger
}

func NewManager(sessionUseCase session.UseCase, authUseCase auth.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *Manager {
	return &Manager{sessionUseCase: sessionUseCase, authUseCase: authUseCase, cfg: cfg, origins: origins, logger: logger}
}

func (m *Manager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)

		req := c.Request()
		res := c.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start)
		m.logger.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, s)
		return err
	}
}
