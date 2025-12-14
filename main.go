package main

import (
	"fmt"
	"go_study/handler"    // 핸들러 import
	"go_study/repository" // 저장소 import

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. DB 연결 및 초기화
	repository.InitDB()

	// 2. Gin 엔진 생성
	r := gin.Default()

	// 3. 핸들러 연결
	// handler 패키지의 함수들을 가져와서 연결합니다.
	r.GET("/todos", handler.GetTodos)
	r.POST("/todos", handler.AddTodo)

	// 4. 서버 실행
	fmt.Println("Starting Organized Gin Server on :8080...")
	r.Run(":8080")
}
