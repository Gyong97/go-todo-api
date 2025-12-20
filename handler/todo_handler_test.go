package handler

import (
	"bytes"
	"encoding/json"
	"go_study/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. Mock 객체 정의
// "mock.Mock"을 상속받으면(Embedding) 가짜 기능을 쓸 수 있습니다.
type MockTodoRepository struct {
	mock.Mock
}

// 2. 인터페이스 구현 (껍데기 만들기)
// TodoRepository 인터페이스의 4개 함수를 다 구현해야 합니다.
// 하지만 내용은 실제 DB가 아니라 "m.Called"를 통해 기록하고 가짜 값을 리턴합니다.

func (m *MockTodoRepository) Save(t model.Todo) (model.Todo, error) {
	// "Save 함수가 호출됐음!" 하고 기록하고, 인자를 받아옵니다.
	args := m.Called(t)
	// 미리 설정해둔 첫 번째 리턴값(0)을 model.Todo 타입으로 형변환해서 반환
	return args.Get(0).(model.Todo), args.Error(1)
}

func (m *MockTodoRepository) GetAll() []model.Todo {
	args := m.Called()
	return args.Get(0).([]model.Todo)
}

func (m *MockTodoRepository) Update(id string) (model.Todo, error) {
	args := m.Called(id)
	return args.Get(0).(model.Todo), args.Error(1)
}

func (m *MockTodoRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// [추가] 통계 조회 Mock
func (m *MockTodoRepository) GetStats() (int64, int64, error) {
	args := m.Called()
	// 리턴값이 int64 2개와 error 1개임
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

// [추가] 미완료 목록 조회 Mock
func (m *MockTodoRepository) GetPendingTodos() ([]model.Todo, error) {
	args := m.Called()
	// 리턴값이 []model.Todo 슬라이스와 error 1개임
	return args.Get(0).([]model.Todo), args.Error(1)
}

// ----------------------------------------------------------------
// 실제 테스트 함수
// ----------------------------------------------------------------

func TestAddTodo_Success(t *testing.T) {
	// 1. 준비 (Arrange)
	// 가짜 저장소(Mock)를 만듭니다.
	mockRepo := new(MockTodoRepository)

	// 테스트용 데이터
	inputTodo := model.Todo{Task: "Mock Test", Done: false}
	expectedTodo := model.Todo{ID: 1, Task: "Mock Test", Done: false}

	// ⭐️ 핵심: 시나리오 설정 (Expectation)
	// "Save 함수에 inputTodo가 들어오면 -> expectedTodo와 nil(에러없음)을 리턴해라!"
	mockRepo.On("Save", inputTodo).Return(expectedTodo, nil)

	// 핸들러에 가짜 저장소를 주입 (Dependency Injection)
	h := NewTodoHandler(mockRepo)

	// 2. 실행 (Act)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/todos", h.AddTodo)

	// 요청 생성 (JSON 바디)
	jsonData, _ := json.Marshal(inputTodo)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// 3. 검증 (Assert)
	// 상태 코드가 201 Created 인가?
	assert.Equal(t, http.StatusCreated, w.Code)

	// 응답 본문에 ID가 1로 박혀서 나왔는가?
	var responseTodo model.Todo
	json.Unmarshal(w.Body.Bytes(), &responseTodo)
	assert.Equal(t, uint(1), responseTodo.ID)
	assert.Equal(t, "Mock Test", responseTodo.Task)

	// ⭐️ Mock 검증: "정말로 Save 함수가 호출되었는가?"
	mockRepo.AssertExpectations(t)
}

func TestGenerateDailyReport_Accepted(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTodoRepository)
	// 비동기 고루틴 안에서 GetPendingTodos가 호출될 수도, 안 될 수도 있음 (Timing 이슈)
	// 따라서 "호출된다면 빈 배열 리턴해라" 정도로 유연하게 설정하거나,
	// 단순히 핸들러의 응답 코드만 체크할 거면 Mock 설정 없이도 202는 반환됨.

	// 안전하게: 혹시 고루틴이 엄청 빨리 돌아서 호출할까봐 넣어둠
	mockRepo.On("GetPendingTodos").Return([]model.Todo{}, nil).Maybe()

	h := NewTodoHandler(mockRepo)

	// 2. Act
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/reports", h.GenerateDailyReport)

	req, _ := http.NewRequest("POST", "/reports", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 3. Assert
	// 즉시 응답이 202 Accepted 인가?
	assert.Equal(t, http.StatusAccepted, w.Code)
}
