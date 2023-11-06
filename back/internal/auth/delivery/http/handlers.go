package http

import (
	"errors"
	"github.com/aakosarev/kanban-board/back/internal/auth"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/models"
	"github.com/aakosarev/kanban-board/back/internal/session"
	httpErrors "github.com/aakosarev/kanban-board/back/pkg/http_errors"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"github.com/aakosarev/kanban-board/back/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthHandlers struct {
	group     *echo.Group
	log       logger.Logger
	cfg       *config.Config
	v         *validator.Validate
	authUC    auth.UseCase
	sessionUC session.UseCase
}

func NewAuthHandlers(
	group *echo.Group,
	log logger.Logger,
	cfg *config.Config,
	v *validator.Validate,
	authUC auth.UseCase,
	sessionUC session.UseCase,
) *AuthHandlers {
	return &AuthHandlers{group: group, log: log, cfg: cfg, v: v, authUC: authUC, sessionUC: sessionUC}
}

func (h *AuthHandlers) Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			h.log.Errorf("(utils.ReadRequest) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createdUser, err := h.authUC.Register(c.Request().Context(), user)
		if err != nil {
			h.log.Errorf("(authUC.Register) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		sess, err := h.sessionUC.CreateSession(c.Request().Context(), &models.Session{
			UserID: createdUser.ID,
		}, h.cfg.Session.Expire)
		if err != nil {
			h.log.Errorf("(sessionUC.CreateSession) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

func (h *AuthHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			h.log.Errorf("(utils.ReadRequest) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		foundUser, err := h.authUC.Login(c.Request().Context(), &models.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			h.log.Errorf("(authUC.Login) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		sess, err := h.sessionUC.CreateSession(c.Request().Context(), &models.Session{
			UserID: foundUser.ID,
		}, h.cfg.Session.Expire)
		if err != nil {
			h.log.Errorf("(sessionUC.CreateSession) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		return c.JSON(http.StatusOK, foundUser)
	}
}

func (h *AuthHandlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(h.cfg.Session.Name)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				h.log.Error(err)
				return httpErrors.ErrorCtxResponse(c, httpErrors.Unauthorized, h.cfg.Http.DebugErrorsResponse)
			}
			h.log.Error(err)
			return httpErrors.NewInternalServerError(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.sessionUC.DeleteByID(c.Request().Context(), cookie.Value); err != nil {
			h.log.Error(err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		utils.DeleteSessionCookie(c, h.cfg.Session.Name)

		return c.NoContent(http.StatusOK)
	}
}
