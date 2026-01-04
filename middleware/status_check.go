package middleware

import (
	"go_study/global" // global 패키지 경로 확인
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckActive : 현재 서버가 Active 상태인지 확인하는 미들웨어
func CheckActive(c *gin.Context) {
	// global 패키지의 상태 확인
	if !global.IsActive() {
		// Standby 상태라면 503 에러 리턴 및 중단
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"error": "⛔ Service Unavailable (Server is in STANDBY mode)",
		})
		return
	}

	// Active 상태라면 통과
	c.Next()
}
