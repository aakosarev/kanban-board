package kanban

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateColumn() echo.HandlerFunc
	DeleteColumn() echo.HandlerFunc
	ChangeNameColumn() echo.HandlerFunc

	CreateTask() echo.HandlerFunc
	DeleteTask() echo.HandlerFunc
	ChangeDescriptionTask() echo.HandlerFunc
	ChangeColumnIDTask() echo.HandlerFunc

	GetKanbanBoardByUserID() echo.HandlerFunc
}
