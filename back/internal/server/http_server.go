package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	stackSize      = 1 << 10 // 1 KB
	bodyLimit      = "2M"
	readTimeout    = 15 * time.Second
	writeTimeout   = 15 * time.Second
)

func (s *Server) runHttpServer() error {
	s.mapRoutes()
	s.echo.Server.ReadTimeout = readTimeout
	s.echo.Server.WriteTimeout = writeTimeout
	s.echo.Server.MaxHeaderBytes = maxHeaderBytes

	return s.echo.Start(s.cfg.Http.Port)
}

func (s *Server) mapRoutes() {
	s.echo.Use(s.m.RequestLoggerMiddleware)
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.BodyLimit(bodyLimit))
}
