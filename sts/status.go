package sts

import "github.com/gin-gonic/gin"

type Status struct {
	Code    string      `json:"code"`
	Msg     string      `json:"msg"`
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
	case "OK":
		ctx.JSON(200, s)
	case "":
		ctx.JSON(204, nil)
	case "PARAM_ERROR":
		ctx.JSON(400, s)
	case "SIGN_ERROR":
		ctx.JSON(401, s)
	case "SYSTEM_ERROR":
		ctx.JSON(500, s)
	}
}
