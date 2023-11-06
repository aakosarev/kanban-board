package middleware

import (
	"context"
	httpErrors "github.com/aakosarev/kanban-board/back/pkg/http_errors"
	"github.com/aakosarev/kanban-board/back/pkg/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (m *Manager) AuthSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(m.cfg.Session.Name)
		if err != nil {
			m.logger.Errorf("AuthSessionMiddleware RequestID: %s, Error: %s",
				utils.GetRequestID(c),
				err.Error(),
			)
			if err == http.ErrNoCookie {
				return httpErrors.NewUnauthorizedError(c, err, m.cfg.Http.DebugErrorsResponse)
			}
			return httpErrors.NewUnauthorizedError(c, httpErrors.Unauthorized, m.cfg.Http.DebugErrorsResponse)
		}

		sid := cookie.Value

		sess, err := m.sessionUseCase.GetSessionByID(c.Request().Context(), cookie.Value)
		if err != nil {
			m.logger.Errorf("GetSessionByID RequestID: %s, CookieValue: %s, Error: %s",
				utils.GetRequestID(c),
				cookie.Value,
				err.Error(),
			)
			return httpErrors.NewUnauthorizedError(c, httpErrors.Unauthorized, m.cfg.Http.DebugErrorsResponse)
		}

		user, err := m.authUseCase.GetByID(c.Request().Context(), sess.UserID)
		if err != nil {
			m.logger.Errorf("GetByID RequestID: %s, Error: %s",
				utils.GetRequestID(c),
				err.Error(),
			)
			return httpErrors.NewUnauthorizedError(c, httpErrors.Unauthorized, m.cfg.Http.DebugErrorsResponse)
		}

		c.Set("sid", sid)
		c.Set("uid", sess.SessionID)
		c.Set("user", user)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))

		m.logger.Info(
			"SessionMiddleware, RequestID: %s,  IP: %s, UserID: %s, CookieSessionID: %s",
			utils.GetRequestID(c),
			utils.GetIPAddress(c),
			strconv.Itoa(user.ID),
			cookie.Value,
		)

		return next(c)
	}
}
