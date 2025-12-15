# Go Todo REST API 🚀

Go 언어로 작성된 RESTful API 서버입니다.  
초기에는 JSON 파일 저장 방식을 사용했으나, 현재는 **SQLite DB**와 **GORM**을 도입하여 데이터 영속성을 보장하며, **계층형 아키텍처(Layered Architecture)**와 **의존성 주입(Dependency Injection)** 패턴을 적용해 유지보수성과 테스트 용이성을 높였습니다.

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin Gonic
- **Database**: SQLite (Embedded)
- **ORM**: GORM
- **Testing**: Testify (Assert, Mock)

## 🏗 Architecture

이 프로젝트는 **관심사의 분리(Separation of Concerns)**를 위해 3계층 구조를 따릅니다.

1.  **Handler Layer (`/handler`)**: HTTP 요청/응답 처리, 파라미터 파싱, 에러 핸들링.
2.  **Repository Layer (`/repository`)**: DB 접근 및 데이터 조작 (CRUD).
3.  **Model Layer (`/model`)**: 데이터 구조체 정의 (Entity).

### Key Features (Refactoring)
- **Dependency Injection (DI)**: `main.go`에서 의존성을 주입하여 결합도를 낮춤.
- **Interface**: `TodoRepository` 인터페이스를 통해 구현체를 추상화.
- **Unit Testing**:
    - **Handler**: `Mock` 객체를 사용하여 DB 없이 컨트롤러 로직 검증.
    - **Repository**: `In-Memory SQLite`를 사용하여 실제 쿼리 로직 검증.

## 📂 Project Structure
```
.
├── handler/            # HTTP Request Handler (Controller)
│   ├── todo_handler.go
│   └── todo_handler_test.go
├── repository/         # DB Access Layer
│   ├── interface.go    # Repository Interface definition
│   ├── todo_repo.go    # SQLite implementation
│   └── todo_repo_test.go
├── model/              # Data Models
│   └── todo.go
├── main.go             # Entry Point & Dependency Wiring
├── go.mod              # Go Modules
└── README.md
```
## 🔌 API Endpoints

| Method | Endpoint       | Description             | Body (JSON)             |
| :----- | :------------- | :---------------------- | :---------------------- |
| GET    | `/todos`       | 할 일 목록 전체 조회    | -                       |
| POST   | `/todos`       | 할 일 추가              | `{"task": "Go 공부"}`   |
| PATCH  | `/todos/:id`   | 할 일 완료 상태 토글    | -                       |
| DELETE | `/todos/:id`   | 할 일 삭제              | -                       |

## 🚀 How to Run

### 1. Prerequisite
Go 언어가 설치되어 있어야 합니다.

# 의존성 패키지 설치
go mod tidy

### 2. Run Server
go run main.go

실행 후 http://localhost:8080 에서 접속 가능합니다.

### 3. Run Tests
단위 테스트(Unit Test)와 통합 테스트(Integration Test)를 수행합니다.

go test ./... -v

---
*Created by Gyong97*