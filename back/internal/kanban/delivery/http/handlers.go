package http

import (
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/kanban"
	"github.com/aakosarev/kanban-board/back/internal/models"
	httpErrors "github.com/aakosarev/kanban-board/back/pkg/http_errors"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"github.com/aakosarev/kanban-board/back/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type KanbanHandlers struct {
	taskGroup   *echo.Group
	columnGroup *echo.Group
	boardGroup  *echo.Group
	log         logger.Logger
	cfg         *config.Config
	v           *validator.Validate
	kanbanUC    kanban.UseCase
}

func NewKanbanHandlers(
	taskGroup *echo.Group,
	columnGroup *echo.Group,
	boardGroup *echo.Group,
	log logger.Logger,
	cfg *config.Config,
	v *validator.Validate,
	kanbanUC kanban.UseCase,
) *KanbanHandlers {
	return &KanbanHandlers{taskGroup: taskGroup, columnGroup: columnGroup, boardGroup: boardGroup, log: log, cfg: cfg, v: v, kanbanUC: kanbanUC}
}

func (h *KanbanHandlers) CreateColumn() echo.HandlerFunc {
	return func(c echo.Context) error {
		column := &models.Column{}
		if err := utils.ReadRequest(c, column); err != nil {
			h.log.Errorf("(utils.ReadRequest) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createdColumn, err := h.kanbanUC.CreateColumn(c.Request().Context(), column)
		if err != nil {
			h.log.Errorf("(kanbanUC.CreateColumn) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusCreated, createdColumn)
	}
}

func (h *KanbanHandlers) DeleteColumn() echo.HandlerFunc {
	return func(c echo.Context) error {
		columnIDStr := c.Param("column_id")
		columnID, err := strconv.Atoi(columnIDStr)
		if err != nil {
			h.log.Errorf("(KanbanHandlers.DeleteColumn.Atoi) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest, h.cfg.Http.DebugErrorsResponse)
		}

		if err = h.kanbanUC.DeleteColumn(c.Request().Context(), columnID); err != nil {
			h.log.Errorf("(kanbanUC.DeleteColumn) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *KanbanHandlers) ChangeNameColumn() echo.HandlerFunc {
	return func(c echo.Context) error {
		columnIDStr := c.Param("column_id")
		columnID, err := strconv.Atoi(columnIDStr)
		if err != nil {
			h.log.Errorf("(KanbanHandlers.ChangeNameColumn.Atoi) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest, h.cfg.Http.DebugErrorsResponse)
		}

		column := &models.Column{}
		if err := utils.ReadRequest(c, column); err != nil {
			h.log.Errorf("(utils.ReadRequest) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		column.ID = columnID

		updatedColumn, err := h.kanbanUC.ChangeNameColumn(c.Request().Context(), column)
		if err != nil {
			h.log.Errorf("(kanbanUC.ChangeNameColumn) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, updatedColumn)
	}
}

func (h *KanbanHandlers) CreateTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		task := &models.Task{}
		if err := utils.ReadRequest(c, task); err != nil {
			h.log.Errorf("(utils.ReadRequest) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createdTask, err := h.kanbanUC.CreateTask(c.Request().Context(), task)
		if err != nil {
			h.log.Errorf("(kanbanUC.CreateTask) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusCreated, createdTask)
	}
}

func (h *KanbanHandlers) DeleteTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		taskIDStr := c.Param("task_id")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			h.log.Errorf("(KanbanHandlers.DeleteTask.Atoi) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest, h.cfg.Http.DebugErrorsResponse)
		}

		if err = h.kanbanUC.DeleteTask(c.Request().Context(), taskID); err != nil {
			h.log.Errorf("(kanbanUC.DeleteTask) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *KanbanHandlers) ChangeDescriptionTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		taskIDStr := c.Param("task_id")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			h.log.Errorf("(KanbanHandlers.ChangeDescriptionTask.Atoi) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest, h.cfg.Http.DebugErrorsResponse)
		}

		task := &models.Task{}
		if err := utils.ReadRequest(c, task); err != nil {
			h.log.Errorf("(utils.ReadRequest) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		task.ID = taskID

		updatedTask, err := h.kanbanUC.ChangeDescriptionTask(c.Request().Context(), task)
		if err != nil {
			h.log.Errorf("(kanbanUC.ChangeDescriptionTask) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, updatedTask)
	}
}

func (h *KanbanHandlers) ChangeColumnIDTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		taskIDStr := c.Param("task_id")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			h.log.Errorf("(KanbanHandlers.ChangeColumnIDTask.Atoi) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest, h.cfg.Http.DebugErrorsResponse)
		}

		task := &models.Task{}
		if err := utils.ReadRequest(c, task); err != nil {
			h.log.Errorf("(utils.ReadRequest) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		task.ID = taskID

		updatedTask, err := h.kanbanUC.ChangeColumnIDTask(c.Request().Context(), task)
		if err != nil {
			h.log.Errorf("(kanbanUC.ChangeColumnIDTask) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, updatedTask)
	}
}

func (h *KanbanHandlers) GetKanbanBoardByUserID() echo.HandlerFunc {
	return func(c echo.Context) error {
		userIDStr := c.Param("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			h.log.Errorf("(KanbanHandlers.GetKanbanBoardByUserID.Atoi) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest, h.cfg.Http.DebugErrorsResponse)
		}

		board, err := h.kanbanUC.GetKanbanBoardByUserID(c.Request().Context(), userID)
		if err != nil {
			h.log.Errorf("(kanbanUC.GetKanbanBoardByUserID) err: {%v}", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, board)
	}
}
