package sts

import (
	"github.com/gin-gonic/gin"
)

type Status struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Status  string      `json:"status,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Details interface{} `json:"details,omitempty"`
	Total   int         `json:"total,omitempty"`
}

func (s *Status) SetMsg(msg string) *Status {
	s.Msg = msg
	return s
}

func (s *Status) SetData(data interface{}) *Status {
	s.Data = data
	return s
}

func (s *Status) SetDetails(details interface{}) *Status {
	s.Details = details
	return s
}

func (s *Status) SetTotal(total int) *Status {
	s.Total = total
	return s
}

func (s *Status) Resp(ctx *gin.Context) {
	switch s.Code {
	case 0:
		ctx.JSON(200, s)
	case 3:
		ctx.JSON(400, s)
	case 16:
		ctx.JSON(401, s)
	case 13:
		ctx.JSON(500, s)
	}
}
