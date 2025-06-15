package strength

import (
	"testing"
	"time"

	"fitness-mcp-server/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// トレーニングセッションコンテキストのテスト
// =============================================================================

func TestStrengthTraining_NewStrengthTraining(t *testing.T) {
	// Arrange
	id := shared.NewTrainingID()
	date := time.Now()
	notes := "今日はベンチプレス95kgで8回できた"

	// Act
	training := NewStrengthTraining(id, date, notes)

	// Assert
	assert.True(t, training.ID().Equals(id))
	assert.Equal(t, date, training.Date())
	assert.Equal(t, notes, training.Notes())
	assert.Equal(t, 0, training.ExerciseCount())
	assert.Empty(t, training.Exercises())
}

func TestStrengthTraining_AddExercise(t *testing.T) {
	// Arrange
	id := shared.NewTrainingID()
	training := NewStrengthTraining(id, time.Now(), "")
	exercise := NewExercise(BenchPress)

	// Act
	training.AddExercise(exercise)

	// Assert
	assert.Equal(t, 1, training.ExerciseCount())
	exercises := training.Exercises()
	assert.Len(t, exercises, 1)
	assert.True(t, exercises[0].Name().Equals(BenchPress))
}

func TestStrengthTraining_UpdateNotes(t *testing.T) {
	// Arrange
	id := shared.NewTrainingID()
	training := NewStrengthTraining(id, time.Now(), "古いメモ")
	newNotes := "新しいメモ"

	// Act
	training.UpdateNotes(newNotes)

	// Assert
	assert.Equal(t, newNotes, training.Notes())
}

func TestStrengthTraining_GetExerciseByName(t *testing.T) {
	// Arrange
	id := shared.NewTrainingID()
	training := NewStrengthTraining(id, time.Now(), "")

	benchPress := NewExercise(BenchPress)
	squat := NewExercise(Squat)

	training.AddExercise(benchPress)
	training.AddExercise(squat)

	t.Run("存在するエクササイズ", func(t *testing.T) {
		// Act
		exercise, err := training.GetExerciseByName(BenchPress)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, exercise)
		assert.True(t, exercise.Name().Equals(BenchPress))
	})

	t.Run("存在しないエクササイズ", func(t *testing.T) {
		// Act
		exercise, err := training.GetExerciseByName(Deadlift)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, exercise)
		assert.Contains(t, err.Error(), "exercise not found")
	})
}

func TestStrengthTraining_TotalSets(t *testing.T) {
	// Arrange
	id := shared.NewTrainingID()
	training := NewStrengthTraining(id, time.Now(), "")

	exercise1 := NewExercise(BenchPress)
	exercise2 := NewExercise(Squat)

	weight, _ := NewWeight(95.0)
	reps, _ := NewReps(8)
	set := NewSet(weight, reps, nil)

	exercise1.AddSet(set)
	exercise1.AddSet(set)
	exercise2.AddSet(set)

	training.AddExercise(exercise1)
	training.AddExercise(exercise2)

	// Act
	totalSets := training.TotalSets()

	// Assert
	assert.Equal(t, 3, totalSets)
}

func TestStrengthTraining_TotalVolume(t *testing.T) {
	// Arrange
	id := shared.NewTrainingID()
	training := NewStrengthTraining(id, time.Now(), "")

	exercise := NewExercise(BenchPress)
	weight, _ := NewWeight(95.0)
	reps, _ := NewReps(8)
	set := NewSet(weight, reps, nil)

	exercise.AddSet(set)
	exercise.AddSet(set)
	training.AddExercise(exercise)

	// Act
	totalVolume := training.TotalVolume()

	// Assert
	expected := 95.0 * 8 * 2 // 1520
	assert.Equal(t, expected, totalVolume)
}

// =============================================================================
// 統合テスト
// =============================================================================

func TestStrengthTraining_IntegrationTest(t *testing.T) {
	// Arrange: 実際のベンチプレストレーニングをシミュレート
	id := shared.NewTrainingID()
	date := time.Date(2025, 1, 15, 18, 30, 0, 0, time.UTC)
	training := NewStrengthTraining(id, date, "今日は95kg目指す")

	// ベンチプレスエクササイズを作成
	benchPress := NewExercise(BenchPress)

	// ウォームアップセット
	warmupWeight, _ := NewWeight(60.0)
	warmupReps, _ := NewReps(10)
	benchPress.AddSet(NewSet(warmupWeight, warmupReps, nil))

	// メインセット
	mainWeight, _ := NewWeight(95.0)
	mainReps1, _ := NewReps(8)
	mainReps2, _ := NewReps(6)
	mainReps3, _ := NewReps(5)

	rpe8, _ := NewRPE(8)
	rpe9, _ := NewRPE(9)
	rpe10, _ := NewRPE(10)

	benchPress.AddSet(NewSet(mainWeight, mainReps1, &rpe8))
	benchPress.AddSet(NewSet(mainWeight, mainReps2, &rpe9))
	benchPress.AddSet(NewSet(mainWeight, mainReps3, &rpe10))

	training.AddExercise(benchPress)

	// Act & Assert: 各種計算が正しく行われることを確認
	assert.Equal(t, 1, training.ExerciseCount())
	assert.Equal(t, 4, training.TotalSets())

	// ボリューム計算: 60*10 + 95*8 + 95*6 + 95*5 = 600 + 760 + 570 + 475 = 2405
	expectedVolume := 60.0*10 + 95.0*8 + 95.0*6 + 95.0*5
	assert.Equal(t, expectedVolume, training.TotalVolume())

	// 最大重量確認
	maxWeight, err := benchPress.MaxWeight()
	assert.NoError(t, err)
	assert.Equal(t, 95.0, maxWeight.Kg())

	// エクササイズ検索
	foundExercise, err := training.GetExerciseByName(BenchPress)
	assert.NoError(t, err)
	assert.NotNil(t, foundExercise)

	// メモ更新
	training.UpdateNotes("95kg達成！次は97.5kgに挑戦")
	assert.Equal(t, "95kg達成！次は97.5kgに挑戦", training.Notes())
}
