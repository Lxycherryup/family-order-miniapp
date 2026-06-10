package model

import "family-order/server/internal/consts"

// Response 统一接口响应。
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewSuccessResponse 创建成功响应。
func NewSuccessResponse(data interface{}) Response {
	return Response{Code: consts.CodeSuccess, Msg: "success", Data: data}
}

// NewErrorResponse 创建失败响应。
func NewErrorResponse(code int, msg string) Response {
	return Response{Code: code, Msg: msg, Data: nil}
}
