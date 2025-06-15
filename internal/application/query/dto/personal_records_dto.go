package dto

import "time"

// PersonalRecordQueryResult はクエリ層専用の個人記録結果
type PersonalRecordQueryResult struct {
	ExerciseName  string
	MaxWeight     PersonalRecordQueryDetail
	MaxReps       PersonalRecordQueryDetail
	MaxVolume     PersonalRecordQueryDetail
	TotalSessions int
	LastPerformed time.Time
}

// PersonalRecordQueryDetail は記録詳細（Query層専用）
type PersonalRecordQueryDetail struct {
	Value      float64
	Date       time.Time
	TrainingID string
	SetDetails *SetQueryDetails
}

// SetQueryDetails はセット詳細（Query層専用）
type SetQueryDetails struct {
	WeightKg float64
	Reps     int
	RPE      *int
}
