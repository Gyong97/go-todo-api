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

// [추가] 사용자가 입력할 데이터만 정의한 구조체 (DTO)
type CreateTodoInput struct {
	Task string `json:"task" binding:"required" example:"Swagger 문서 수정하기"`
}

// TodoHandler 구조체
// 핵심: 구체적인 *SQLiteRepository가 아니라, 추상적인 인터페이스를 가집니다.
type TodoHandler struct {
	repo repository.TodoRepository // 인터페이스 타입!
}

// 생성자: 외부에서 리포지토리를 주입(Injection) 받습니다.
func NewTodoHandler(r repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: r}
}

// GetTodos godoc
// @Summary     할 일 목록 조회
// @Description 저장된 모든 할 일 목록을 반환합니다.
// @Tags        Todos
// @Accept      json
// @Produce     json
// @Success     200 {object} model.WebResponse{data=[]model.Todo}
// @Router      /todos [get]
func (h *TodoHandler) GetTodos(c *gin.Context) {
	todos := h.repo.GetAll()

	// ✨ 포장지에 싸서 전달
	c.JSON(http.StatusOK, model.WebResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    todos, // 원래 데이터는 여기로!
	})
}

// AddTodo godoc
// @Summary     할 일 추가
// @Description 새로운 할 일을 목록에 추가합니다.
// @Tags        Todos
// @Accept      json
// @Produce     json
// @Param       todo body CreateTodoInput true "할 일 정보"
// @Success     201 {object} model.Todo
// @Failure     400 {object} model.WebResponse{data=nil}
// @Router      /todos [post]
func (h *TodoHandler) AddTodo(c *gin.Context) {
	var input CreateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := model.WebResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    "",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newTodo := model.Todo{
		Task: input.Task,
		Done: false, // 기본값
	}
	createdTodo, err := h.repo.Save(newTodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
		return
	}
	c.JSON(http.StatusCreated, createdTodo)
}

// ToggleTodoStatus godoc
// @Summary      할 일 완료 여부 토글
// @Description  특정 ID의 할 일 완료 상태(Done)를 반전시킵니다.
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "수정할 할 일 ID"
// @Success      200  {object}  model.Todo
// @Failure      400  {object}  model.WebResponse  "잘못된 ID 형식"
// @Failure      404  {object}  model.WebResponse  "ID를 찾을 수 없음"
// @Router       /todos/{id} [patch]
func (h *TodoHandler) ToggleTodoStatus(c *gin.Context) {
	id := c.Param("id")
	updatedTodo, err := h.repo.Update(id)

	if err != nil {
		// 1. 데이터가 없어서 난 에러인지 확인
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, model.WebResponse{
				Code:    http.StatusNotFound,
				Message: "Data not found",
				Data:    nil,
			})
			return
		}
		// 2. 그 외의 DB 에러 (연결 끊김, 제약조건 위반 등)
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: "Fail to Update",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, updatedTodo)
}

// DeleteTodo godoc
// @Summary      할 일 삭제
// @Description  특정 ID의 할 일을 영구적으로 삭제합니다.
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Param        id   path      int              true  "삭제할 할 일 ID"
// @Success      200  {object}  model.WebResponse "삭제 성공"
// @Failure      400  {object}  model.WebResponse "잘못된 ID 형식"
// @Failure      404  {object}  model.WebResponse "ID를 찾을 수 없음"
// @Failure      500  {object}  model.WebResponse "서버 내부 에러"
// @Router       /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(id); err != nil {
		// 에러 종류 확인: "데이터가 없어서 에러난 거야?"
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, model.WebResponse{
				Code:    http.StatusNotFound,
				Message: "Data not found",
				Data:    nil,
			})
			return
		}
		// 그 외의 진짜 에러 (DB 다운 등)
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: "Fail to Delete",
			Data:    nil,
		})
		return
	}

	// ✨ 데이터가 없을 때는 Data에 nil을 넣거나 생략
	c.JSON(http.StatusOK, model.WebResponse{
		Code:    http.StatusOK,
		Message: "성공적으로 삭제되었습니다.",
		Data:    nil,
	})
}

// [POST] /reports - 무거운 리포트 생성 작업 (비동기)
// GenerateDailyReport godoc
// @Summary      일일 리포트 생성 요청
// @Description  리포트 생성을 비동기로 요청합니다. (처리 결과는 이메일 발송 등)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Success      202  {object}  model.WebResponse  "요청 접수됨"
// @Router       /reports [post]
func (h *TodoHandler) GenerateDailyReport(c *gin.Context) {
	// 1. 즉시 응답 (Non-blocking)
	c.JSON(http.StatusAccepted, model.WebResponse{
		Code:    http.StatusOK,
		Message: "리포트 생성 요청이 접수되었습니다. (백그라운드 처리 중)",
		Data:    nil,
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
// GetDashboard godoc
// @Summary      대시보드 데이터 조회
// @Description  사용자 정보와 할 일 통계를 병렬로 조회하여 반환합니다.
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Success 	 200 {object} model.WebResponse{data=[]string}
// @Router       /dashboard [get]
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
	c.JSON(http.StatusOK, model.WebResponse{
		Code:    http.StatusOK,
		Message: "대시보드 데이터 조회 완료",
		Data:    responseData,
	})
}
