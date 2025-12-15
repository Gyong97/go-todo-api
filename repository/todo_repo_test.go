package repository

import (
	"fmt"
	"go_study/model"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 테스트용 인메모리 DB 생성 헬퍼 함수
func newTestSQLiteRepository() *SQLiteRepository {
	// ":memory:" 옵션을 쓰면 파일이 안 생기고 RAM에서만 동작합니다. (속도 엄청 빠름)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Todo{})

	// 우리가 만든 생성자 함수를 이용해 Repository 인스턴스 반환
	return NewSQLiteRepository(db)
}

func TestSQLiteRepository_SaveAndGet(t *testing.T) {
	// 1. 준비 (Arrange) - 가짜(Mock)가 아니라 "메모리 DB"를 쓴 구현체
	repo := newTestSQLiteRepository()
	newTodo := model.Todo{Task: "Integration Test", Done: false}

	// 2. 실행 (Act)
	// 이제 전역 함수가 아니라 repo 객체의 메소드를 호출합니다.
	created, err := repo.Save(newTodo)

	// 3. 검증 (Assert)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID) // ID가 1 이상이어야 함

	// 조회 테스트
	todos := repo.GetAll()
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, "Integration Test", todos[0].Task)
}

func TestSQLiteRepository_UpdateAndDelete(t *testing.T) {
	// 1. 준비
	repo := newTestSQLiteRepository()
	todo := model.Todo{Task: "To be deleted", Done: false}
	saved, _ := repo.Save(todo)

	// 2. 실행 - 업데이트
	// saved.ID는 uint 타입이므로 string으로 변환하거나,
	// 기존 Update 함수가 string을 받도록 되어 있다면 맞춰줘야 합니다.
	// (이전 단계에서 id를 string으로 받기로 했었죠? GORM은 알아서 처리해줍니다)

	// ID를 문자열로 변환 (fmt.Sprint 사용 등)해야 하지만,
	// 편의상 Update/Delete 메소드 시그니처가 string을 받는다고 가정하고 진행합니다.
	// 만약 에러가 난다면 strconv.Itoa(int(saved.ID))를 쓰면 됩니다.

	// 여기서는 편의상 ID를 이용해 직접 조회 후 업데이트 테스트 로직을 수행한다고 가정
	// 실제로는 Update 함수의 인자 타입에 맞춰주세요.
	updated, err := repo.Update(fmt.Sprint(saved.ID))

	assert.NoError(t, err)
	assert.True(t, updated.Done, "상태가 true로 바뀌어야 함")

	// 3. 실행 - 삭제
	err = repo.Delete(fmt.Sprint(saved.ID))
	assert.NoError(t, err)

	// 삭제 확인
	todos := repo.GetAll()
	assert.Equal(t, 0, len(todos), "데이터가 비어 있어야 함")
}
