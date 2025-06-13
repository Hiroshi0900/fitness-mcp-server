package dto

import (
	"fmt"
	"time"

	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/domain/strength"
)

// =============================================================================
// DTOマッパー - ドメインオブジェクトとDTOの変換処理
// =============================================================================

// ToStrengthTraining はRecordTrainingCommandからStrengthTrainingエンティティを生成します
func (cmd *RecordTrainingCommand) ToStrengthTraining() (*strength.StrengthTraining, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 新しいトレーニングIDを生成
	trainingID := shared.NewTrainingID()
	training := strength.NewStrengthTraining(trainingID, cmd.Date, cmd.Notes)

	// エクササイズを追加
	for _, exerciseDTO := range cmd.Exercises {
		exercise, err := exerciseDTO.ToExercise()
		if err != nil {
			return nil, fmt.Errorf("failed to create exercise: %w", err)
		}
		training.AddExercise(exercise)
	}

	return training, nil
}

// ToStrengthTraining はUpdateTrainingCommandからStrengthTrainingエンティティを生成します
func (cmd *UpdateTrainingCommand) ToStrengthTraining() (*strength.StrengthTraining, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// IDを復元
	trainingID, err := shared.NewTrainingIDFromString(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid training ID: %w", err)
	}

	training := strength.NewStrengthTraining(trainingID, cmd.Date, cmd.Notes)

	// エクササイズを追加
	for _, exerciseDTO := range cmd.Exercises {
		exercise, err := exerciseDTO.ToExercise()
		if err != nil {
			return nil, fmt.Errorf("failed to create exercise: %w", err)
		}
		training.AddExercise(exercise)
	}

	return training, nil
}

// ToExercise はExerciseDTOからExerciseエンティティを生成します
func (dto *ExerciseDTO) ToExercise() (*strength.Exercise, error) {
	if err := dto.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// エクササイズ名を作成
	exerciseName, err := strength.NewExerciseName(dto.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid exercise name: %w", err)
	}

	// エクササイズカテゴリを作成
	exerciseCategory, err := strength.NewExerciseCategory(dto.Category)
	if err != nil {
		return nil, fmt.Errorf("invalid exercise category: %w", err)
	}

	exercise := strength.NewExercise(exerciseName, exerciseCategory)

	// セットを追加
	for _, setDTO := range dto.Sets {
		set, err := setDTO.ToSet()
		if err != nil {
			return nil, fmt.Errorf("failed to create set: %w", err)
		}
		exercise.AddSet(set)
	}

	return exercise, nil
}

// ToSet はSetDTOからSetを生成します
func (dto *SetDTO) ToSet() (strength.Set, error) {
	if err := dto.Validate(); err != nil {
		return strength.Set{}, fmt.Errorf("validation failed: %w", err)
	}

	// 重量を作成
	weight, err := strength.NewWeight(dto.WeightKg)
	if err != nil {
		return strength.Set{}, fmt.Errorf("invalid weight: %w", err)
	}

	// レップ数を作成
	reps, err := strength.NewReps(dto.Reps)
	if err != nil {
		return strength.Set{}, fmt.Errorf("invalid reps: %w", err)
	}

	// 休憩時間を作成
	restTime, err := strength.NewRestTime(time.Duration(dto.RestTimeSeconds) * time.Second)
	if err != nil {
		return strength.Set{}, fmt.Errorf("invalid rest time: %w", err)
	}

	// RPEを作成（オプション）
	var rpe *strength.RPE
	if dto.RPE != nil {
		rpeValue, err := strength.NewRPE(*dto.RPE)
		if err != nil {
			return strength.Set{}, fmt.Errorf("invalid RPE: %w", err)
		}
		rpe = &rpeValue
	}

	return strength.NewSet(weight, reps, restTime, rpe), nil
}

// FromStrengthTraining はStrengthTrainingエンティティからTrainingSessionDTOを生成します
func FromStrengthTraining(training *strength.StrengthTraining) *TrainingSessionDTO {
	exercises := make([]ExerciseDTO, 0, len(training.Exercises()))
	
	for _, exercise := range training.Exercises() {
		exerciseDTO := FromExercise(exercise)
		exercises = append(exercises, *exerciseDTO)
	}

	return &TrainingSessionDTO{
		ID:        training.ID().String(),
		Date:      training.Date(),
		Exercises: exercises,
		Notes:     training.Notes(),
		CreatedAt: training.Date(), // 暫定的にDateを使用
		UpdatedAt: training.Date(), // 暫定的にDateを使用
	}
}

// FromExercise はExerciseエンティティからExerciseDTOを生成します
func FromExercise(exercise *strength.Exercise) *ExerciseDTO {
	sets := make([]SetDTO, 0, len(exercise.Sets()))
	
	for _, set := range exercise.Sets() {
		setDTO := FromSet(set)
		sets = append(sets, *setDTO)
	}

	return &ExerciseDTO{
		Name:     exercise.Name().String(),
		Category: exercise.Category().String(),
		Sets:     sets,
	}
}

// FromSet はSetからSetDTOを生成します
func FromSet(set strength.Set) *SetDTO {
	var rpe *int
	if set.RPE() != nil {
		rpeValue := set.RPE().Rating()
		rpe = &rpeValue
	}

	return &SetDTO{
		WeightKg:        set.Weight().Kg(),
		Reps:            set.Reps().Count(),
		RestTimeSeconds: int(set.RestTime().Duration().Seconds()),
		RPE:             rpe,
	}
}
