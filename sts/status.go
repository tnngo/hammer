package sts

import "github.com/gin-gonic/gin"

type Status struct {
	Code    string      `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Details interface{} `json:"details"`
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

func (s *Status) Resp(ctx *gin.Context) {
	if s.Code != "" {
		ctx.JSON(200, nil)
	} else {
		ctx.JSON(204, s)
	}
}
