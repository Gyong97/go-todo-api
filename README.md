# Go Todo API Server

이 프로젝트는 **Layered Architecture**를 적용하여 구축한 간단한 Todo REST API 서버입니다.
C 언어 기반의 백엔드 개발 경험을 바탕으로, **Go 언어의 동시성 모델(Goroutine)** 과 **모던 웹 프레임워크(Gin)** 를 학습하고 적용하는 데 초점을 두었습니다.

## 🛠 Tech Stack
- **Language:** Go (Golang)
- **Framework:** Gin Web Framework
- **Architecture:** 3-Tier Layered Architecture (Controller - Service/Repository - Model)
- **Data:** In-memory storage with File persistence (JSON)

## 🚀 Key Features
- **RESTful API:** `GET`, `POST` 메소드를 활용한 자원 관리.
- **Data Persistence:** 서버 재시작 시에도 데이터가 유지되도록 `os` 패키지를 활용한 파일 I/O 구현.
- **Concurrency Safety:** `sync.Mutex`를 사용하여 멀티 스레드 환경(Goroutine)에서의 **Race Condition 방지**.
- **Clean Architecture:** `handler`, `repository`, `model` 패키지 분리를 통한 유지보수성 확보.

## 📂 Project Structure
```
go-todo-api/ 
├── main.go # Entry Point 
├── model/ # Data Structures (Domain) 
├── repository/ # Data Access Layer (File I/O, Lock) 
└── handler/ # HTTP Request Handlers (Gin)
```

## 📝 Learning Point (From C to Go)
- **Mutex & Defer:** C 언어의 `pthread_mutex`와 달리 `defer` 키워드를 사용하여 리소스 릭(Leak)을 방지하고 데드락 위험을 줄였습니다.
- **Package Access Control:** Go의 대/소문자 접근 제어 규칙을 이해하고, Repository 패턴을 통해 전역 변수에 대한 직접 접근을 제한(Encapsulation)했습니다.
