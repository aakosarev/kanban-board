package http

func (h *KanbanHandlers) MapRoutes() {
	h.columnGroup.POST("/create", h.CreateColumn())
	h.columnGroup.DELETE("/:column_id", h.DeleteColumn())
	h.columnGroup.PATCH("/:column_id/update_name", h.ChangeNameColumn())

	h.taskGroup.POST("/create", h.CreateTask())
	h.taskGroup.DELETE("/:task_id", h.DeleteTask())
	h.taskGroup.PATCH("/:task_id/update_description", h.ChangeDescriptionTask())
	h.taskGroup.PATCH("/:task_id/update_column_id", h.ChangeColumnIDTask())

	h.boardGroup.GET("/:user_id", h.GetKanbanBoardByUserID())
}
