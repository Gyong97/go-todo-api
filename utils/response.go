package utils

import (
	"net/http"

	"go_study/model"

	"github.com/gin-gonic/gin"

	"os"
)

// 1. 성공 응답 (200 OK)
func SendSuccess(c *gin.Context, data interface{}) {
	hostname, _ := os.Hostname()
	c.JSON(http.StatusOK, model.WebResponse{
		Code:    http.StatusOK,
		Message: "Success from" + hostname,
		Data:    data,
	})
}
func SendSuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, model.WebResponse{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	})
}

// 2. 생성 성공 응답 (201 Created) - 주로 POST 요청 시
func SendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, model.WebResponse{
		Code:    http.StatusCreated,
		Message: "Created",
		Data:    data,
	})
}

// 3. 에러 응답 (400, 404, 500 등)
func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, model.WebResponse{
		Code:    status,
		Message: message,
		Data:    nil, // 에러 날 땐 데이터 없음
	})
}
