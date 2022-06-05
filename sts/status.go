package sts

type Status struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
