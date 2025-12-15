package repository

import (
	"go_study/model"

	"gorm.io/gorm"
)

// SQLiteRepository 구조체 (실제 구현체)
type SQLiteRepository struct {
	db *gorm.DB
}

// 생성자 함수: DB 연결 객체를 받아서 Repository 인스턴스를 반환
func NewSQLiteRepository(db *gorm.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// -------------------------------------------------------
// 아래 함수들은 이제 (r *SQLiteRepository)에 소속된 메소드입니다.
// 메소드 이름과 시그니처가 interface.go에 정의된 것과 똑같아야 합니다.
// -------------------------------------------------------

func (r *SQLiteRepository) Save(t model.Todo) (model.Todo, error) {
	// r.db 를 사용 (전역변수 db가 아님)
	result := r.db.Create(&t)
	return t, result.Error
}

func (r *SQLiteRepository) GetAll() []model.Todo {
	var todos []model.Todo
	r.db.Find(&todos)
	return todos
}

func (r *SQLiteRepository) Update(id string) (model.Todo, error) {
	var todo model.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		return todo, err
	}
	err := r.db.Model(&todo).Update("done", !todo.Done).Error
	return todo, err
}

func (r *SQLiteRepository) Delete(id string) error {
	// 1. 삭제 명령 실행
	result := r.db.Delete(&model.Todo{}, id)

	// 2. DB 에러 체크 (문법 에러나 커넥션 에러 등)
	if result.Error != nil {
		return result.Error
	}

	// 3. ✨ 영향받은 행 개수 체크 (C의 SQL%ROWCOUNT)
	if result.RowsAffected == 0 {
		// GORM에 정의된 "데이터 없음" 에러를 리턴합니다.
		return gorm.ErrRecordNotFound
	}
	return nil
}
