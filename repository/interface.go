package repository

import "go_study/model"

// TodoRepository 인터페이스 (계약서)
// "이 기능을 구현한 녀석이라면 누구든 내 저장소가 될 수 있어!"
type TodoRepository interface {
	Save(t model.Todo) (model.Todo, error)
	GetAll() []model.Todo
	Update(id string) (model.Todo, error)
	Delete(id string) error
}
