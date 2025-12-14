package handler

import (
	"go_study/model"
	"go_study/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	// repository에서 데이터를 가져옴 (직접 접근 X)
	todos := repository.GetAll()
	c.JSON(http.StatusOK, todos)
}

func AddTodo(c *gin.Context) {
	var newTodo model.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// repository 통해 저장
	createdTodo, err := repository.AddTodo(newTodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
		return
	}

	c.JSON(http.StatusCreated, createdTodo)
}
