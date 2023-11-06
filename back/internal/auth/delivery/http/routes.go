package http

func (h *AuthHandlers) MapRoutes() {
	h.group.POST("/signup", h.Signup())
	h.group.POST("/login", h.Login())
	h.group.POST("/logout", h.Logout())
}
