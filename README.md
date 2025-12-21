# Go Todo REST API (Production Ready) 🚀

Go 언어와 Gin 프레임워크로 구축된 **엔터프라이즈급 RESTful API 서버**입니다.  
초기 프로토타입에서 발전하여, 현재는 **계층형 아키텍처**, **의존성 주입(DI)**, **설정 관리**, **구조화된 로깅**, **API 문서화** 등 실무에서 요구되는 핵심 요소들을 모두 갖추고 있습니다.

## 🛠 Tech Stack

| Category | Technology | Description |
| :--- | :--- | :--- |
| **Language** | Go (1.25+) | Backend Language |
| **Framework** | Gin Gonic | High-performance Web Framework |
| **Database** | SQLite & GORM | Embedded DB & ORM |
| **Config** | Viper | Configuration Management (YAML/Env) |
| **Logging** | Zap & Lumberjack | Structured Logging & Log Rotation |
| **Docs** | Swagger (Swag) | API Documentation Generator |
| **Deploy** | Docker | Multi-stage Container Build |

## 🏗 Architecture & Key Features

이 프로젝트는 **관심사의 분리(Separation of Concerns)** 원칙을 철저히 준수합니다.

### 1. Layered Architecture
* **Handler (`/handler`)**: HTTP 요청 처리, 파라미터 검증, 응답 표준화.
* **Repository (`/repository`)**: DB 접근 추상화 (Interface 사용).
* **Model (`/model`)**: 데이터 엔티티 및 DTO 정의.
* **Middleware (`/middleware`)**: 로깅, 에러 복구(Recovery) 등의 공통 관심사 처리.

### 2. Production-Ready Features
* **Standardized Response**: 모든 API 응답을 `WebResponse` 구조체(`code`, `message`, `data`)로 통일하여 클라이언트 예측 가능성 확보.
* **Dependency Injection (DI)**: `main.go`에서 의존성을 주입하여 결합도를 낮추고 테스트 용이성 확보.
* **Configuration**: 하드코딩을 제거하고 `config.yaml`을 통해 환경 설정 관리.
* **Concurrency**:
    * `POST /reports`: 고루틴(Goroutine)을 이용한 **비동기(Async) 작업 처리**.
    * `GET /dashboard`: 채널(Channel)과 WaitGroup을 이용한 **병렬(Parallel) 데이터 조회**.

## 📂 Project Structure

```bash
.
├── config/             # Viper Configuration Loader
├── docs/               # Swagger Documentation (Auto-generated)
├── handler/            # Controller Logic & DTOs
├── middleware/         # Zap Logger & Global Middlewares
├── model/              # DB Entity & WebResponse Struct
├── repository/         # DB Access Interface & Implementation
├── utils/              # Helper Functions (Response wrappers)
├── config.yaml         # Configuration File
├── Dockerfile          # Multi-stage Build Dockerfile
└── main.go             # Application Entry Point
```

🔌 API Documentation (Swagger)
서버 실행 후 아래 주소에서 대화형 API 문서를 확인할 수 있습니다.

👉 Swagger UI: http://localhost:8080/swagger/index.html

Standard Response Format
모든 API는 아래와 같은 JSON 형식으로 응답합니다.

```JSON
{
  "code": 200,
  "message": "Success",
  "data": { ... } // Payload
}
```

🚀 How to Run
Option 1. Local Run
Go가 설치된 환경에서 실행합니다.

```Bash
# 1. 의존성 설치
go mod tidy

# 2. 문서 생성 (코드 변경 시)
swag init

# 3. 서버 실행
go run main.go
```

Option 2. Docker Run
Docker가 설치된 환경에서 실행합니다.

```Bash
# 1. 이미지 빌드
docker build -t go-todo-api .

# 2. 컨테이너 실행
docker run -p 8080:8080 go-todo-api
```

Created by Gyong97