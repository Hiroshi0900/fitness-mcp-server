package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"fitness-mcp-server/internal/domain/running"
	"fitness-mcp-server/internal/interface/repository"
)

// RunningRepository はSQLiteを使ったランニングRepository実装（書き込み専用）
type RunningRepository struct {
	db *sql.DB
}

// NewRunningRepository は新しいSQLite RunningRepositoryを作成します
func NewRunningRepository(db *sql.DB) repository.RunningRepository {
	return &RunningRepository{db: db}
}

// Save はランニングセッションを保存します
func (r *RunningRepository) Save(session *running.RunningSession) error {
	log.Printf("Saving running session: %s", session.ID().String()[:8])

	// 心拍数の処理（オプション）
	var heartRateBPM *int
	if session.HeartRate() != nil {
		bpm := session.HeartRate().BPM()
		heartRateBPM = &bpm
	}

	_, err := r.db.Exec(`
		INSERT INTO running_sessions (
			id, date, distance_km, duration_seconds, pace_seconds_per_km, 
			heart_rate_bpm, run_type, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		session.ID().String(),
		session.Date(),
		session.Distance().Km(),
		int(session.Duration().Value().Seconds()),
		session.Pace().SecondsPerKm(),
		heartRateBPM,
		session.RunType().String(),
		session.Notes(),
	)
	if err != nil {
		log.Printf("Failed to save running session: %v", err)
		return fmt.Errorf("failed to save running session: %w", err)
	}

	log.Printf("Successfully saved running session: %s", session.ID().String()[:8])
	return nil
}

// コンパイル時のインターフェース実装チェック
var _ repository.RunningRepository = (*RunningRepository)(nil)
