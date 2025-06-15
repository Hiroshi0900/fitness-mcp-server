package dto

import (
	"time"

	"fitness-mcp-server/internal/domain/strength"
)

// =============================================================================
// Query系のDTO定義
// =============================================================================

// GetTrainingsByDateRangeQuery は期間指定でトレーニングを取得するクエリ
type GetTrainingsByDateRangeQuery struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// GetTrainingsByDateRangeResponse は期間指定トレーニング取得のレスポンス
type GetTrainingsByDateRangeResponse struct {
	Trainings []*TrainingDTO `json:"trainings"`
	Count     int            `json:"count"`
	Period    string         `json:"period"`
}

// TrainingDTO はトレーニングセッションのDTO
type TrainingDTO struct {
	ID        string         `json:"id"`
	Date      time.Time      `json:"date"`
	Exercises []*ExerciseDTO `json:"exercises"`
	Notes     string         `json:"notes"`
	Summary   *SummaryDTO    `json:"summary"`
}

// ExerciseDTO はエクササイズのDTO
type ExerciseDTO struct {
	Name string    `json:"name"`
	Sets []*SetDTO `json:"sets"`
}

// SetDTO はセットのDTO
type SetDTO struct {
	WeightKg float64 `json:"weight_kg"`
	Reps     int     `json:"reps"`
	RPE      *int    `json:"rpe,omitempty"`
}

// SummaryDTO はトレーニングセッションの概要DTO
type SummaryDTO struct {
	TotalExercises int     `json:"total_exercises"`
	TotalSets      int     `json:"total_sets"`
	TotalVolume    float64 `json:"total_volume"`
	Duration       string  `json:"duration"`
}

// =============================================================================
// ドメインエンティティからDTOへの変換関数
// =============================================================================

// TrainingToDTO はStrengthTrainingをTrainingDTOに変換します
func TrainingToDTO(training *strength.StrengthTraining) *TrainingDTO {
	exercises := make([]*ExerciseDTO, 0, len(training.Exercises()))
	for _, exercise := range training.Exercises() {
		exercises = append(exercises, ExerciseToDTO(exercise))
	}

	return &TrainingDTO{
		ID:        training.ID().String(),
		Date:      training.Date(),
		Exercises: exercises,
		Notes:     training.Notes(),
		Summary: &SummaryDTO{
			TotalExercises: training.ExerciseCount(),
			TotalSets:      training.TotalSets(),
			TotalVolume:    training.TotalVolume(),
			Duration:       "", // 実装時に追加
		},
	}
}

// ExerciseToDTO はExerciseをExerciseDTOに変換します
func ExerciseToDTO(exercise *strength.Exercise) *ExerciseDTO {
	sets := make([]*SetDTO, 0, len(exercise.Sets()))
	for _, set := range exercise.Sets() {
		sets = append(sets, SetToDTO(set))
	}

	return &ExerciseDTO{
		Name: exercise.Name().String(),
		Sets: sets,
	}
}

// SetToDTO はSetをSetDTOに変換します
func SetToDTO(set strength.Set) *SetDTO {
	dto := &SetDTO{
		WeightKg: set.Weight().Kg(),
		Reps:     set.Reps().Count(),
	}

	if rpe := set.RPE(); rpe != nil {
		value := rpe.Rating()
		dto.RPE = &value
	}

	return dto
}
