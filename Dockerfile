# [Stage 1] 빌드 단계 (컴파일러 포함된 무거운 환경)
FROM golang:1.25-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 파일 먼저 복사 (캐싱 효율을 위해)
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 전체 복사
COPY . .

# 빌드 실행
# CGO_ENABLED=0: C 라이브러리 의존성 없이 순수 Go 바이너리로 빌드 (필수!)
# GOOS=linux: 리눅스용 실행 파일 생성
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# --------------------------------------------------------

# [Stage 2] 실행 단계 (실행 파일만 담을 가벼운 환경)
FROM alpine:latest

# 작업 디렉토리 설정
WORKDIR /root/

# Stage 1에서 빌드한 'server' 파일만 복사해옴
COPY --from=builder /app/server .

# 컨테이너가 8080 포트를 쓴다고 명시
EXPOSE 8080

# 컨테이너 켜지면 실행할 명령어
CMD ["./server"]