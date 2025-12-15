package handler

import (
	"go_study/model"
	"go_study/repository"
	"net/http"

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
