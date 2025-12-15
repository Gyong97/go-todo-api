package main

import (
	"fmt"
	"go_study/handler"
	"go_study/model"
	"go_study/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. DB 연결 (Infrastructure Layer)
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.Todo{})

	// 2. Repository 생성 (인터페이스 구현체)
	// SQLiteRepository 인스턴스를 만듭니다.
	todoRepo := repository.NewSQLiteRepository(db)

	// 3. Handler 생성 (의존성 주입) ⭐
	// Handler에게 "너는 이 리포지토리를 써"라고 주입해줍니다.
	todoHandler := handler.NewTodoHandler(todoRepo)

	// 4. Gin 라우팅 설정
	r := gin.Default()

	// 이제 핸들러가 메소드이므로 인스턴스(todoHandler)를 통해 호출합니다.
	r.GET("/todos", todoHandler.GetTodos)
	r.POST("/todos", todoHandler.AddTodo)
	r.PATCH("/todos/:id", todoHandler.ToggleTodoStatus)
	r.DELETE("/todos/:id", todoHandler.DeleteTodo)

	fmt.Println("Starting Server with Dependency Injection...")
	r.Run(":8080")
}
