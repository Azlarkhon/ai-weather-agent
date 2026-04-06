package http

import (
	"net/http"

	"ai-agent/internal/agent"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	agent *agent.Agent
}

func NewHandler(agent *agent.Agent) *Handler {
	return &Handler{agent: agent}
}

func (h *Handler) Chat(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&req); err != nil {
		return
	}

	res, err := h.agent.Run(req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res})
}
