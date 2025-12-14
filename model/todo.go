package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	// gorm.Model 대신 필요한 것만 직접 정의합니다.
	// `gorm:"primaryKey"` 태그로 이게 PK임을 알려줍니다.
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`     // 생성 시간
	UpdatedAt time.Time      `json:"-"`              // 수정 시간 (JSON에는 숨김)
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 삭제 시간 (Soft Delete용, JSON 숨김)

	Task string `json:"task"`
	Done bool   `json:"done"`
}
