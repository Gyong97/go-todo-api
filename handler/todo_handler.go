package handler

import (
	"fmt"
	"go_study/model"
	"go_study/repository"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TodoHandler 구조체
// 핵심: 구체적인 *SQLiteRepository가 아니라, 추상적인 인터페이스를 가집니다.
type TodoHandler struct {
	repo repository.TodoRepository // 인터페이스 타입!
}

// 생성자: 외부에서 리포지토리를 주입(Injection) 받습니다.
func NewTodoHandler(r repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: r}
}

// 이제 모든 핸들러 함수는 TodoHandler의 메소드가 됩니다.
func (h *TodoHandler) GetTodos(c *gin.Context) {
	// h.repo를 통해 호출 (실제 뒤에 SQLite가 있는지 Mock이 있는지 모름)
	todos := h.repo.GetAll()
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) AddTodo(c *gin.Context) {
	var newTodo model.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTodo, err := h.repo.Save(newTodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
		return
	}
	c.JSON(http.StatusCreated, createdTodo)
}

func (h *TodoHandler) ToggleTodoStatus(c *gin.Context) {
	id := c.Param("id")
	updatedTodo, err := h.repo.Update(id)

	if err != nil {
		// 1. 데이터가 없어서 난 에러인지 확인
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		// 2. 그 외의 DB 에러 (연결 끊김, 제약조건 위반 등)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	c.JSON(http.StatusOK, updatedTodo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(id); err != nil {
		// 에러 종류 확인: "데이터가 없어서 에러난 거야?"
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		// 그 외의 진짜 에러 (DB 다운 등)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// [POST] /reports - 무거운 리포트 생성 작업 (비동기)
func (h *TodoHandler) GenerateDailyReport(c *gin.Context) {
	// 1. 즉시 응답 (Non-blocking)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "리포트 생성 요청이 접수되었습니다. (백그라운드 처리 중)",
	})

	// 2. 백그라운드 작업 (Goroutine)
	go func() {
		log.Println("📝 [Background] 리포트 데이터 수집 시작...")

		// ✨ 진짜 DB 조회: 미완료된 할 일 가져오기
		pendingTodos, err := h.repo.GetPendingTodos()
		if err != nil {
			log.Printf("❌ 리포트 생성 실패: %v\n", err)
			return
		}

		// 리포트 내용 작성 (파일로 저장하거나 이메일 보내는 척)
		reportContent := fmt.Sprintf("=== Daily Report ===\n남은 할 일: %d건\n", len(pendingTodos))
		for _, t := range pendingTodos {
			reportContent += fmt.Sprintf("- [ ] %s\n", t.Task)
		}

		// 시간 조금 걸리는 척 (리포트 파일 생성 시뮬레이션)
		time.Sleep(2 * time.Second)

		log.Printf("✅ [Background] 리포트 생성 완료!\n%s", reportContent)
		// 실제로는 여기서 smtp.SendMail() 등을 호출함
	}()
}

// [GET] /dashboard - 병렬 처리 예제
func (h *TodoHandler) GetDashboard(c *gin.Context) {
	// 결과를 모을 채널 생성 (문자열이 지나다니는 파이프)
	// 버퍼(2)를 주어서 송신자가 블로킹되지 않게 함
	results := make(chan string, 2)

	// WaitGroup 생성 (스레드 조인용 카운터)
	var wg sync.WaitGroup

	// "나 2개 기다릴 거야" 설정
	wg.Add(2)

	// --- [작업 1] 사용자 프로필 조회 ---
	go func() {
		defer wg.Done() // 함수 끝나면 무조건 카운트 -1

		time.Sleep(1 * time.Second) // 1초 걸리는 척
		log.Println("👤 프로필 조회 완료")
		results <- "User Profile: Gyong97 (Level 99)" // 채널에 데이터 쏘기
	}()

	// --- [작업 2] 통계 집계 ---
	go func() {
		defer wg.Done() // 함수 끝나면 무조건 카운트 -1
		// ✨ 진짜 DB 조회!
		total, done, err := h.repo.GetStats()
		if err != nil {
			results <- fmt.Sprintf("Stats Error: %v", err)
			return
		}
		// 통계 문자열 생성
		statsMsg := fmt.Sprintf("Stats: Total %d / Done %d (Rate: %.0f%%)",
			total, done, float64(done)/float64(total)*100)
		results <- statsMsg
	}()

	// --- [중요 패턴] 기다리기 & 채널 닫기 ---
	// 메인 고루틴이 멈추면 안 되니까, "기다리는 역할"도 별도 고루틴에게 시킴
	go func() {
		wg.Wait()      // 작업 2개가 다 끝날 때까지 대기
		close(results) // 다 끝났으면 파이프 입구를 막음 (그래야 받는 쪽 반복문이 끝남)
	}()

	// --- [결과 수집] ---
	var responseData []string

	// 채널이 닫힐 때까지 데이터를 계속 꺼냄 (Range Loop)
	for msg := range results {
		responseData = append(responseData, msg)
	}

	// 클라이언트 응답
	c.JSON(http.StatusOK, gin.H{
		"dashboard": responseData,
	})
}
