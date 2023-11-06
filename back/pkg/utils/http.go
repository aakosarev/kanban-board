package utils

import (
	"context"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/models"
	httpErrors "github.com/aakosarev/kanban-board/back/pkg/http_errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func CreateSessionCookie(cfg *config.Config, session string) *http.Cookie {
	return &http.Cookie{
		Name:       cfg.Session.Name,
		Value:      session,
		Path:       "/",
		RawExpires: "",
		MaxAge:     cfg.Session.Expire,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HTTPOnly,
		SameSite:   0,
	}
}

func DeleteSessionCookie(c echo.Context, sessionName string) {
	c.SetCookie(&http.Cookie{
		Name:   sessionName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

type ReqIDCtxKey struct{}

func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func GetCtxWithReqID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*15)
	ctx = context.WithValue(ctx, ReqIDCtxKey{}, GetRequestID(c))
	return ctx, cancel
}

func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

type UserCtxKey struct{}

func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(UserCtxKey{}).(*models.User)
	if !ok {
		return nil, httpErrors.Unauthorized
	}

	return user, nil
}

func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}

func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}
