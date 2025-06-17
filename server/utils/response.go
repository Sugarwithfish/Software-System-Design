package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewResponse 创建一个新的响应
func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// 成功响应码
const (
	SUCCESS = 0
)

// 错误响应码
const (
	ERROR_PARAM        = -1  // 参数错误
	ERROR_UNAUTHORIZED = -2  // 未授权
	ERROR_FORBIDDEN    = -3  // 禁止访问a
	ERROR_NOT_FOUND    = -4  // 资源不存在
	ERROR_INTERNAL     = -5  // 内部服务器错误
	ERROR_BUSINESS     = -10 // 业务错误
)

// Response 统一响应处理
func Result(c *gin.Context, code int, msg string, data interface{}) {
	// 始终返回 HTTP 200 状态码
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, "操作成功", data)
}

// SuccessWithMsg 带自定义消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	Result(c, SUCCESS, msg, data)
}

// FailWithMsg 带消息的失败响应
func FailWithMsg(c *gin.Context, code int, msg string) {
	Result(c, code, msg, nil)
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, msg string) {
	FailWithMsg(c, ERROR_PARAM, msg)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, msg string) {
	FailWithMsg(c, ERROR_UNAUTHORIZED, msg)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, msg string) {
	FailWithMsg(c, ERROR_FORBIDDEN, msg)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, msg string) {
	FailWithMsg(c, ERROR_NOT_FOUND, msg)
}

// InternalError 内部服务器错误响应
func InternalError(c *gin.Context, msg string) {
	FailWithMsg(c, ERROR_INTERNAL, msg)
}

// BusinessError 业务错误响应
func BusinessError(c *gin.Context, msg string) {
	FailWithMsg(c, ERROR_BUSINESS, msg)
}
