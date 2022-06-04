package hammer

import (
	"fmt"
	"net/http"
	"time"

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
	h.s.Use(func(ctx *gin.Context) {

	})
	h.s.GET("/start/test/good/", func(ctx *gin.Context) {

	})

	go func() {
		for {
			lad.L().Debug("正在检测http服务...")
			time.Sleep(1 * time.Microsecond)
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/start/test/good/", h.Port))
			if err != nil {
				lad.L().Error(err.Error())
				continue
			}

			if resp.StatusCode == http.StatusOK {
				fields := make([]lad.Field, 0)

				if h.Name != "" {
					fields = append(fields, lad.String("name", h.Name))
				}
				fields = append(fields, lad.String("local_addr", fmt.Sprintf("%s:%d", localIP(), h.Port)))

				lad.L().Info("Http server started successfully", fields...)
				return
			}
		}
	}()

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
