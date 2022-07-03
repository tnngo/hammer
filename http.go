package hammer

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tnngo/lad"
)

type Http struct {
	s    *gin.Engine
	Port uint16
	Name string
}

func (h *Http) Run() {
	if h.s == nil {
		h.s = gin.Default()
	}

	if err := h.s.Run(fmt.Sprintf(":%d", h.Port)); err != nil {
		lad.L().Error(err.Error())
	}
}

func (h *Http) Server() *gin.Engine {
	return h.s
}

const (
	defaultHttpPort = 8080
)

func newHttpServer() *gin.Engine {
	return gin.Default()
}

func NewHttpServer() hammer {
	return &Http{
		s:    newHttpServer(),
		Port: defaultHttpPort,
	}
}
