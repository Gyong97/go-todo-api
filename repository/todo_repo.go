package repository

import (
	"go_study/model"
	"log"

	"github.com/glebarez/sqlite" // Pure Go SQLite 드라이버
	"gorm.io/gorm"
)

// 전역 DB 객체 (Connection Pool 관리)
var db *gorm.DB

// DB 초기화 및 연결 (서버 켤 때 한 번만 실행)
func InitDB() {
	var err error
	// 1. SQLite DB 파일 열기 (없으면 생성됨)
	// gorm.Open은 내부적으로 Connection Pool을 생성합니다.
	db, err = gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 2. Auto Migration (핵심!) ⭐
	// C에서는 "CREATE TABLE..." 쿼리를 직접 짜서 테이블을 만들었죠?
	// GORM은 구조체(Todo)를 보고 알아서 테이블을 생성/수정해줍니다.
	err = db.AutoMigrate(&model.Todo{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

// [INSERT] 데이터 추가
func AddTodo(t model.Todo) (model.Todo, error) {
	// INSERT INTO todos (task, done) VALUES (...)
	result := db.Create(&t)
	return t, result.Error
}

// [SELECT] 전체 조회
func GetAll() []model.Todo {
	var todos []model.Todo
	// SELECT * FROM todos
	db.Find(&todos)
	return todos
}
