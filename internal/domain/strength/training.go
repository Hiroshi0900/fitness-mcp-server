package strength

import (
	"fmt"
	"time"

	"fitness-mcp-server/internal/domain/shared"
)

// =============================================================================
// トレーニングセッションコンテキスト - 1回のトレーニング全体の管理
// =============================================================================

// StrengthTraining は筋トレセッションを表すエンティティ
type (
	StrengthTraining struct {
		id        shared.TrainingID // トレーニングID
		date      time.Time         // トレーニング日
		exercises []*Exercise       // エクササイズのリスト
		notes     string            // メモ
	}
)

// NewStrengthTraining は新しいStrengthTrainingを作成します
func NewStrengthTraining(id shared.TrainingID, date time.Time, notes string) *StrengthTraining {
	return &StrengthTraining{
		id:        id,
		date:      date,
		exercises: make([]*Exercise, 0),
		notes:     notes,
	}
}

// ID はトレーニングIDを返します
func (st *StrengthTraining) ID() shared.TrainingID {
	return st.id
}

// Date はトレーニング日を返します
func (st *StrengthTraining) Date() time.Time {
	return st.date
}

// Exercises は全エクササイズを返します
func (st *StrengthTraining) Exercises() []*Exercise {
	// コピーを返して不変性を保つ
	result := make([]*Exercise, len(st.exercises))
	copy(result, st.exercises)
	return result
}

// Notes はメモを返します
func (st *StrengthTraining) Notes() string {
	return st.notes
}

// AddExercise はエクササイズを追加します
func (st *StrengthTraining) AddExercise(exercise *Exercise) {
	st.exercises = append(st.exercises, exercise)
}

// UpdateNotes はメモを更新します
func (st *StrengthTraining) UpdateNotes(notes string) {
	st.notes = notes
}

// ExerciseCount はエクササイズ数を返します
func (st *StrengthTraining) ExerciseCount() int {
	return len(st.exercises)
}

// TotalSets は総セット数を返します
func (st *StrengthTraining) TotalSets() int {
	totalSets := 0
	for _, exercise := range st.exercises {
		totalSets += exercise.SetCount()
	}
	return totalSets
}

// TotalVolume は総ボリュームを返します
func (st *StrengthTraining) TotalVolume() float64 {
	totalVolume := 0.0
	for _, exercise := range st.exercises {
		totalVolume += exercise.TotalVolume()
	}
	return totalVolume
}

// GetExerciseByName は名前でエクササイズを検索します
func (st *StrengthTraining) GetExerciseByName(name ExerciseName) (*Exercise, error) {
	for _, exercise := range st.exercises {
		if exercise.name.Equals(name) {
			return exercise, nil
		}
	}
	return nil, fmt.Errorf("exercise not found: %s", name.String())
}

// String はトレーニングの文字列表現を返します
func (st *StrengthTraining) String() string {
	return fmt.Sprintf("トレーニング %s (%s) - %d種目, %dセット",
		st.id.String()[:8], st.date.Format("2006-01-02"),
		len(st.exercises), st.TotalSets())
}
