package handler

import (
	"github.com/gin-gonic/gin"
)

// PanicHandler is a test handler to trigger panics.
type PanicHandler struct{}

// NewPanicHandler creates a new instance of PanicHandler
func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

// Handle @Summary Test unexpected error
// @Description Intentionally panics to test the recovery middleware
// @Tags auth
// @Produce json
// @Router /auth/panic [get]
func (h *PanicHandler) Handle(c *gin.Context) {
	panic("This is an unexpected error designed to test the recovery middleware!")
}
