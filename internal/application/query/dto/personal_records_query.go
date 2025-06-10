package dto

import "time"

type (
	GetPersonalRecordsQuery struct {
		ExerciseName *string `json:"exercise_name,omitempty"` // オプション: 特定のエクササイズ名でフィルタリング
	}

	GetPersonalRecordsResponse struct {
		Records []PersonalRecord `json:"records"` // 個人記録のリスト
		Count   int              `json:"count"`   // レコードの総数
	}

	// PersonalRecord は個人記録のデータ転送オブジェクト
	PersonalRecord struct {
		ExerciseName  string               `json:"exercise_name"`  // エクササイズ名
		Category      string               `json:"category"`       // カテゴリ
		MaxWeight     PersonalRecordDetail `json:"max_weight"`     // 最大重量
		MaxReps       PersonalRecordDetail `json:"max_reps"`       // 最大レップ数
		MaxVolume     PersonalRecordDetail `json:"max_volume"`     // 最大ボリューム
		TotalSessions int                  `json:"total_sessions"` // セッション数
		LastPerformed time.Time            `json:"last_performed"` // 最終実施日時
	}

	// PersonalRecordDetail は個人記録の詳細情報
	PersonalRecordDetail struct {
		Value      float64   `json:"value"`                 // 記録値（kgやレップ数）
		Date       time.Time `json:"date"`                  // 記録日時
		TrainingID string    `json:"training_id"`           // 関連するトレーニングセッションのID
		SetDetails *SetInfo  `json:"set_details,omitempty"` // オプション: セットの詳細情報
	}

	// SetInfo はセットの詳細情報
	SetInfo struct {
		WeightKg        float64 `json:"weight_kg"`         // 重量（kg）
		Reps            int     `json:"reps"`              // レップ数
		RestTimeSeconds int     `json:"rest_time_seconds"` // 休憩時間（秒）
		RPE             *int    `json:"rpe,omitempty"`     // オプション: RPE（Rate of Perceived Exertion）
	}
)
