package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 전역 로거 변수 (편의상)
var Log *zap.Logger

// 로거 초기화 (파일 저장 + 터미널 출력)
func InitLogger() {
	// 1. 로그 파일 설정 (Lumberjack)
	writeSyncer := getLogWriter()

	// 2. 로그 인코더 설정 (JSON 형식)
	encoder := getEncoder()

	// 3. Core 생성 (터미널 + 파일 동시에 출력하려면 MultiWriteSyncer 사용)
	// zapcore.AddSync(os.Stdout): 터미널에도 출력
	// writeSyncer: 파일에도 출력
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.InfoLevel)

	// 4. 로거 생성
	// AddCaller: 로그 찍은 파일명과 라인 수 표시 (logger.go:45)
	Log = zap.New(core, zap.AddCaller())

	defer Log.Sync()
}

// 인코더 설정 (JSON 포맷, 시간 포맷 등)
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// 시간 포맷을 사람이 읽기 편하게 (2025-12-20T10:00:00.000Z)
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 로그 파일 저장 설정 (Log Rotation)
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/server.log", // 로그 파일 위치
		MaxSize:    10,                  // 파일 하나당 최대 크기 (MB)
		MaxBackups: 5,                   // 백업 파일 최대 개수
		MaxAge:     30,                  // 최대 보관 일수
		Compress:   false,               // 백업 파일 압축 여부 (gzip)
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 2. Gin 미들웨어 함수
func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// [전처리] 타이머 시작
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// --- 핸들러로 요청을 넘김 (Next) ---
		c.Next()
		// -------------------------------

		// [후처리] 타이머 종료 및 로깅
		end := time.Now()
		latency := end.Sub(start) // 처리 소요 시간

		// 로그에 남길 필드들 정의 (Structured Logging)
		Log.Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),            // HTTP 상태 코드 (200, 404 등)
			zap.String("method", c.Request.Method),          // GET, POST
			zap.String("path", path),                        // /todos
			zap.String("query", query),                      // ?id=1
			zap.String("ip", c.ClientIP()),                  // 요청자 IP
			zap.String("user-agent", c.Request.UserAgent()), // 브라우저 정보
			zap.Duration("latency", latency),                // 소요 시간 (중요!)
		)
	}
}
