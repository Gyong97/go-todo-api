package model

// WebResponse: 모든 API 응답의 표준 포맷
type WebResponse struct {
	Code    int         `json:"code" example:"200"`        // 결과 코드 (HTTP Status와 동일하게 쓰거나 커스텀 코드 사용)
	Message string      `json:"message" example:"Success"` // 사람이 읽을 수 있는 메시지
	Data    interface{} `json:"data"`                      // 실제 데이터 (Payload) - 무엇이든 들어갈 수 있음(Generic)
}
