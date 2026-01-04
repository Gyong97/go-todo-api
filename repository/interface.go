package repository

import (
	"go_study/model"

	"gorm.io/gorm"
)

// TodoRepository 인터페이스 (계약서)
// "이 기능을 구현한 녀석이라면 누구든 내 저장소가 될 수 있어!"
type TodoRepository interface {
	Save(t model.Todo) (model.Todo, error)
	GetAll() []model.Todo
	Update(id string) (model.Todo, error)
	Delete(id string) error

	// 👇 [추가] 통계 정보를 가져오는 함수 (전체 개수, 완료 개수, 에러)
	GetStats() (int64, int64, error)
	// 👇 [추가] 완료되지 않은 할 일만 가져오는 함수
	GetPendingTodos() ([]model.Todo, error)

	// 🚀 [추가] DB 연결 상태 확인용 접근자
	GetDB() *gorm.DB
}
