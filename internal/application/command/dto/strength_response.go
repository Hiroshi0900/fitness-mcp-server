package dto

import (
	"time"
)

// =============================================================================
// 筋トレレスポンスDTO - ユースケース実行結果のデータ構造
// =============================================================================

// RecordTrainingResult は筋トレセッション記録結果DTO
type RecordTrainingResult struct {
	TrainingID string    `json:"training_id"`
	Date       time.Time `json:"date"`
	Message    string    `json:"message"`
}

// UpdateTrainingResult は筋トレセッション更新結果DTO
type UpdateTrainingResult struct {
	TrainingID string    `json:"training_id"`
	Date       time.Time `json:"date"`
	Message    string    `json:"message"`
}

// DeleteTrainingResult は筋トレセッション削除結果DTO
type DeleteTrainingResult struct {
	TrainingID string `json:"training_id"`
	Message    string `json:"message"`
}

// TrainingSessionDTO は筋トレセッション表示用DTO
type TrainingSessionDTO struct {
	ID        string        `json:"id"`
	Date      time.Time     `json:"date"`
	Exercises []ExerciseDTO `json:"exercises"`
	Notes     string        `json:"notes"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// PersonalRecordDTO は個人記録DTO
type PersonalRecordDTO struct {
	ExerciseName string    `json:"exercise_name"`
	MaxWeight    float64   `json:"max_weight"`
	MaxReps      int       `json:"max_reps"`
	OneRepMax    float64   `json:"one_rep_max"`
	AchievedAt   time.Time `json:"achieved_at"`
}

// ProgressAnalysisDTO は進捗分析DTO
type ProgressAnalysisDTO struct {
	ExerciseName      string               `json:"exercise_name"`
	Period            string               `json:"period"`
	WeightProgress    WeightProgressDTO    `json:"weight_progress"`
	VolumeProgress    VolumeProgressDTO    `json:"volume_progress"`
	FrequencyAnalysis FrequencyAnalysisDTO `json:"frequency_analysis"`
}

// WeightProgressDTO は重量進捗DTO
type WeightProgressDTO struct {
	StartWeight    float64   `json:"start_weight"`
	CurrentWeight  float64   `json:"current_weight"`
	WeightIncrease float64   `json:"weight_increase"`
	ProgressRate   float64   `json:"progress_rate"` // %
	LastUpdate     time.Time `json:"last_update"`
}

// VolumeProgressDTO は総負荷量進捗DTO
type VolumeProgressDTO struct {
	StartVolume    float64 `json:"start_volume"`
	CurrentVolume  float64 `json:"current_volume"`
	VolumeIncrease float64 `json:"volume_increase"`
	ProgressRate   float64 `json:"progress_rate"` // %
}

// FrequencyAnalysisDTO は実施頻度分析DTO
type FrequencyAnalysisDTO struct {
	SessionsPerWeek  float64 `json:"sessions_per_week"`
	SessionsPerMonth float64 `json:"sessions_per_month"`
	TotalSessions    int     `json:"total_sessions"`
	ConsistencyScore float64 `json:"consistency_score"` // 0-100
}

// BIG3AnalysisDTO はBIG3分析DTO
type BIG3AnalysisDTO struct {
	BenchPress PersonalRecordDTO `json:"bench_press"`
	Squat      PersonalRecordDTO `json:"squat"`
	Deadlift   PersonalRecordDTO `json:"deadlift"`
	TotalMax   float64           `json:"total_max"`
	BodyWeight float64           `json:"body_weight"`
	Ratio      BIG3RatioDTO      `json:"ratio"`
}

// BIG3RatioDTO はBIG3比率DTO
type BIG3RatioDTO struct {
	BenchToBodyWeight    float64 `json:"bench_to_body_weight"`
	SquatToBodyWeight    float64 `json:"squat_to_body_weight"`
	DeadliftToBodyWeight float64 `json:"deadlift_to_body_weight"`
	SquatToBench         float64 `json:"squat_to_bench"`
	DeadliftToBench      float64 `json:"deadlift_to_bench"`
}
