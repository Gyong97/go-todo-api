package repository

import (
	"encoding/json"
	"fmt"
	"go_study/model"
	"os"
	"sync"
)

var (
	todos    []model.Todo
	nextID   = 1
	mu       sync.Mutex
	filename = "todos.json"
)

// (파일 저장/로드 함수는 기존과 동일합니다)
func SaveTodos() error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadTodos() {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Error reading file:", err)
		return
	}
	if err := json.Unmarshal(data, &todos); err != nil {
		fmt.Println("Error parsing json:", err)
		return
	}
	for _, t := range todos {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
	fmt.Println("Loaded", len(todos), "todos from file.")
}

func GetAll() []model.Todo {
	mu.Lock()
	defer mu.Unlock()
	return todos
}

func AddTodo(t model.Todo) (model.Todo, error) {
	mu.Lock()
	defer mu.Unlock()

	t.ID = nextID
	nextID++
	t.Done = false
	todos = append(todos, t)

	err := SaveTodos()

	return t, err
}
