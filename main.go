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

	"go_study/cron"
	"go_study/global"
	"os"
	"strings"
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

	// 🚀 [수정] 환경변수로 초기 상태 결정
	// Docker Compose에서 넣어준 값을 읽어옵니다.
	initialRole := strings.ToLower(os.Getenv("INITIAL_ROLE"))

	if initialRole == "active" {
		global.SetActive()
		middleware.Log.Info("🚀 서버가 ACTIVE 모드로 시작됩니다.")
	} else {
		global.SetStandby()
		middleware.Log.Info("💤 서버가 STANDBY 모드로 시작됩니다.")
	}

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

	cron.StartStatsJob(todoRepo)
	// 4. Gin 라우팅 설정
	// Default()는 기본 로거를 포함하므로, 우리가 만든 걸 쓰려면 New()로 빈 깡통을 만듦
	r := gin.New()

	// 🚀 [추가] 정적 파일(HTML/CSS) 서빙 설정
	// "./static" 폴더를 "/view"라는 주소로 연결하거나, 파일 하나를 특정 주소에 연결
	r.Static("/static", "./static")          // static 폴더 공개
	r.StaticFile("/", "./static/index.html") // 루트(/) 접속 시 index.html 보여주기

	// 미들웨어 부착
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger())

	// 이제 핸들러가 메소드이므로 인스턴스(todoHandler)를 통해 호출합니다.
	api := r.Group("/todos")
	api.Use(middleware.CheckActive)
	{
		api.GET("", todoHandler.GetTodos)
		api.POST("", todoHandler.AddTodo)
		api.PATCH("/:id", todoHandler.ToggleTodoStatus)
		api.DELETE("/:id", todoHandler.DeleteTodo)
	}

	r.POST("/reports", todoHandler.GenerateDailyReport)
	r.GET("/dashboard", todoHandler.GetDashboard)

	// healthcheck, active-stanby 구조
	r.GET("/health", todoHandler.HealthCheck)
	// 🚀 [추가] 관리자용 승격 API (Admin 그룹으로 묶는 게 좋음)
	admin := r.Group("/admin")
	{
		admin.POST("/promote", todoHandler.PromoteToActive)
		admin.POST("/demote", todoHandler.DemoteToStandby)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	middleware.Log.Info("Starting Server with Dependency Injection...")
	r.Run(config.AppConfig.Server.Port)
}
