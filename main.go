package main

import (
	"go_study/config"
	"go_study/handler"
	"go_study/middleware"
	"go_study/model"
	"go_study/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	_ "go_study/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Go Todo API
// @version         1.0
// @description     이것은 Go로 만든 Todo 리스트 API 문서입니다.
// @contact.name    Gyong97
// @contact.email   gyong97@example.com
// @host            localhost:8080
// @BasePath        /
func main() {
	// 설정 로드
	config.LoadConfig()
	// 로거 초기화
	middleware.InitLogger()
	// 프로그램 종료 시 버퍼 비우기
	defer middleware.Log.Sync()

	// 1. DB 연결 (Infrastructure Layer)
	db, err := gorm.Open(sqlite.Open(config.AppConfig.Database.File), &gorm.Config{})
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
	// Default()는 기본 로거를 포함하므로, 우리가 만든 걸 쓰려면 New()로 빈 깡통을 만듦
	r := gin.New()

	// 미들웨어 부착
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger())

	r.StaticFile("/", "./index.html")

	// 이제 핸들러가 메소드이므로 인스턴스(todoHandler)를 통해 호출합니다.
	r.GET("/todos", todoHandler.GetTodos)
	r.POST("/todos", todoHandler.AddTodo)
	r.PATCH("/todos/:id", todoHandler.ToggleTodoStatus)
	r.DELETE("/todos/:id", todoHandler.DeleteTodo)

	r.POST("/reports", todoHandler.GenerateDailyReport)
	r.GET("/dashboard", todoHandler.GetDashboard)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	middleware.Log.Info("Starting Server with Dependency Injection...")
	r.Run(config.AppConfig.Server.Port)
}
