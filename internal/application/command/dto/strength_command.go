package dto

import (
	"fmt"
	"time"
)

// =============================================================================
// 筋トレコマンドDTO - 外部インターフェースとの入出力データ構造
// =============================================================================

// RecordTrainingCommand は筋トレセッション記録コマンドDTO
type RecordTrainingCommand struct {
	Date      time.Time     `json:"date"`
	Exercises []ExerciseDTO `json:"exercises"`
	Notes     string        `json:"notes"`
}

// UpdateTrainingCommand は筋トレセッション更新コマンドDTO
type UpdateTrainingCommand struct {
	ID        string        `json:"id"`
	Date      time.Time     `json:"date"`
	Exercises []ExerciseDTO `json:"exercises"`
	Notes     string        `json:"notes"`
}

// DeleteTrainingCommand は筋トレセッション削除コマンドDTO
type DeleteTrainingCommand struct {
	ID string `json:"id"`
}

// ExerciseDTO はエクササイズDTO
type ExerciseDTO struct {
	Name string   `json:"name"`
	Sets []SetDTO `json:"sets"`
}

// SetDTO はセットDTO
type SetDTO struct {
	WeightKg float64 `json:"weight_kg"`
	Reps     int     `json:"reps"`
	RPE      *int    `json:"rpe,omitempty"` // オプション
}

// Validate はRecordTrainingCommandの妥当性検証を行います
func (cmd *RecordTrainingCommand) Validate() error {
	if cmd.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	if len(cmd.Exercises) == 0 {
		return fmt.Errorf("at least one exercise is required")
	}
	for i, exercise := range cmd.Exercises {
		if err := exercise.Validate(); err != nil {
			return fmt.Errorf("exercise[%d]: %w", i, err)
		}
	}
	return nil
}

// Validate はUpdateTrainingCommandの妥当性検証を行います
func (cmd *UpdateTrainingCommand) Validate() error {
	if cmd.ID == "" {
		return fmt.Errorf("training ID is required")
	}
	if cmd.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	if len(cmd.Exercises) == 0 {
		return fmt.Errorf("at least one exercise is required")
	}
	for i, exercise := range cmd.Exercises {
		if err := exercise.Validate(); err != nil {
			return fmt.Errorf("exercise[%d]: %w", i, err)
		}
	}
	return nil
}

// Validate はDeleteTrainingCommandの妥当性検証を行います
func (cmd *DeleteTrainingCommand) Validate() error {
	if cmd.ID == "" {
		return fmt.Errorf("training ID is required")
	}
	return nil
}

// Validate はExerciseDTOの妥当性検証を行います
func (dto *ExerciseDTO) Validate() error {
	if dto.Name == "" {
		return fmt.Errorf("exercise name is required")
	}
	if len(dto.Sets) == 0 {
		return fmt.Errorf("at least one set is required")
	}
	for i, set := range dto.Sets {
		if err := set.Validate(); err != nil {
			return fmt.Errorf("set[%d]: %w", i, err)
		}
	}
	return nil
}

// Validate はSetDTOの妥当性検証を行います
func (dto *SetDTO) Validate() error {
	if dto.WeightKg <= 0 {
		return fmt.Errorf("weight must be positive")
	}
	if dto.Reps <= 0 {
		return fmt.Errorf("reps must be positive")
	}
	if dto.RPE != nil && (*dto.RPE < 1 || *dto.RPE > 10) {
		return fmt.Errorf("RPE must be between 1 and 10")
	}
	return nil
}
