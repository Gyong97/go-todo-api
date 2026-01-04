package global

import "sync/atomic"

// 서버 상태 상수
const (
	Standby int32 = 0
	Active  int32 = 1
)

// 현재 서버 모드 (기본값: Standby)
// 안전한 동시성 제어를 위해 atomic을 사용합니다.
var serverMode int32 = Standby

// SetActive: 서버를 Active 상태로 변경
func SetActive() {
	atomic.StoreInt32(&serverMode, Active)
}

// SetStandby: 서버를 Standby 상태로 변경
func SetStandby() {
	atomic.StoreInt32(&serverMode, Standby)
}

// IsActive: 현재 Active 상태인지 확인
func IsActive() bool {
	return atomic.LoadInt32(&serverMode) == Active
}
