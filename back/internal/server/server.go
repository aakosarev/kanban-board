package server

import (
	"context"
	authHttp "github.com/aakosarev/kanban-board/back/internal/auth/delivery/http"
	authS "github.com/aakosarev/kanban-board/back/internal/auth/storage"
	authUC "github.com/aakosarev/kanban-board/back/internal/auth/usecase"
	"github.com/aakosarev/kanban-board/back/internal/config"
	kanbanHttp "github.com/aakosarev/kanban-board/back/internal/kanban/delivery/http"
	kanbanS "github.com/aakosarev/kanban-board/back/internal/kanban/storage"
	kanbanUC "github.com/aakosarev/kanban-board/back/internal/kanban/usecase"
	"github.com/aakosarev/kanban-board/back/internal/middleware"
	sessionS "github.com/aakosarev/kanban-board/back/internal/session/storage"
	sessionUC "github.com/aakosarev/kanban-board/back/internal/session/usecase"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	cfg            *config.Config
	log            logger.Logger
	v              *validator.Validate
	echo           *echo.Echo
	m              *middleware.Manager
	ps             *http.Server
	redisClient    *redis.Client
	postgresClient *pgxpool.Pool
	doneCh         chan struct{}
}

func NewServer(cfg *config.Config, log logger.Logger, redisClient *redis.Client, postgresClient *pgxpool.Pool) *Server {
	return &Server{cfg: cfg, log: log, v: validator.New(), redisClient: redisClient, postgresClient: postgresClient, echo: echo.New(), doneCh: make(chan struct{})}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := s.v.StructCtx(ctx, s.cfg); err != nil {
		return errors.Wrap(err, "cfg validate")
	}

	sessionStorage := sessionS.NewSessionStorage(s.redisClient, s.cfg)
	authStorage := authS.NewAuthStorage(s.log, s.postgresClient)
	kanbanStorage := kanbanS.NewKanbanStorage(s.log, s.postgresClient)

	authUseCase := authUC.NewAuthUseCase(s.cfg, authStorage, s.log)
	sessionUseCase := sessionUC.NewSessionUseCase(sessionStorage, s.cfg)
	kanbanUseCase := kanbanUC.NewKanbanUseCase(s.cfg, kanbanStorage, s.log)

	authHandlers := authHttp.NewAuthHandlers(s.echo.Group(s.cfg.Http.AuthPath), s.log, s.cfg, s.v, authUseCase, sessionUseCase)
	kanbanHandlers := kanbanHttp.NewKanbanHandlers(s.echo.Group(s.cfg.Http.TaskPath), s.echo.Group(s.cfg.Http.ColumnPath), s.echo.Group(s.cfg.Http.BoardPath), s.log, s.cfg, s.v, kanbanUseCase)

	s.m = middleware.NewManager(sessionUseCase, authUseCase, s.cfg, []string{"*"}, s.log)

	authHandlers.MapRoutes()
	kanbanHandlers.MapRoutes()

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Errorf("(s.runHttpServer) err: {%v}", err)
			cancel()
		}
	}()
	s.log.Infof("%s is listening on PORT: {%s}", GetMicroserviceName(s.cfg), s.cfg.Http.Port)

	<-ctx.Done()
	s.waitShootDown(waitShotDownDuration)

	if err := s.echo.Shutdown(ctx); err != nil {
		s.log.Warnf("(Shutdown) err: {%v}", err)
	}

	<-s.doneCh
	s.log.Infof("%s server exited properly", GetMicroserviceName(s.cfg))
	return nil
}
