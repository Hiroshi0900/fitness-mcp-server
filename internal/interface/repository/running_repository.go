package repository

import (
	"fitness-mcp-server/internal/domain/running"
)

// RunningRepository はランニングデータの永続化を担当するインターフェース（書き込み専用）
type RunningRepository interface {
	// Save はランニングセッションを保存します
	Save(session *running.RunningSession) error
}
