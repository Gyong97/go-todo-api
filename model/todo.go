package model

import (
	"gorm.io/gorm"
)

type Todo struct {
	// gorm.Model을 상속받으면 ID, CreatedAt, UpdatedAt, DeletedAt이 자동 포함
	gorm.Model

	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}
