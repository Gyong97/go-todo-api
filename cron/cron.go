package cron

import (
	"log"
	"time"

	"go_study/global"
	"go_study/repository"
)

// StartStatsJob: 1분마다 통계를 조회해서 알림을 보내는 함수
func StartStatsJob(repo repository.TodoRepository) {
	// 별도의 고루틴(일꾼)을 생성해서 메인 서버를 방해하지 않게 함
	go func() {
		// 1분(Minute) 간격으로 울리는 알람 시계 생성
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		log.Println("⏰ [Cron] 통계 리포트 작업이 시작되었습니다 (1분 간격)")

		// 알람이 울릴 때마다 루프 실행 (무한 루프)
		for range ticker.C {
			if !global.IsActive() {
				continue
			}
			// 1. DB 조회 (기존에 만들어둔 GetStats 이용)
			total, done, err := repo.GetStats()
			if err != nil {
				log.Printf("❌ [Cron] 통계 조회 실패: %v\n", err)
				continue
			}

			// 2. Slack 전송 (여기서는 로그로 흉내)
			// 실제로는 여기서 http.Post("https://hooks.slack.com/...", ...)를 호출
			log.Printf("🔔 [Slack Bot] 현재 리포트 도착! 📝 총 할 일: %d개 / ✅ 완료: %d개", total, done)
		}
	}()
}
