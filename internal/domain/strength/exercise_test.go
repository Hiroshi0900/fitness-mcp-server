package strength

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// エクササイズ定義コンテキストのテスト
// =============================================================================

func TestExerciseName_IsBIG3(t *testing.T) {
	tests := []struct {
		name     string
		exercise ExerciseName
		isBIG3   bool
	}{
		{
			name:     "ベンチプレス",
			exercise: BenchPress,
			isBIG3:   true,
		},
		{
			name:     "スクワット",
			exercise: Squat,
			isBIG3:   true,
		},
		{
			name:     "デッドリフト",
			exercise: Deadlift,
			isBIG3:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act & Assert
			assert.Equal(t, tt.isBIG3, tt.exercise.IsBIG3())
		})
	}

	// カスタムエクササイズのテスト
	t.Run("カスタムエクササイズ", func(t *testing.T) {
		customExercise, _ := NewExerciseName("カール")
		assert.False(t, customExercise.IsBIG3())
	})
}

func TestExerciseCategory_NewExerciseCategory(t *testing.T) {
	tests := []struct {
		name      string
		category  string
		wantError bool
	}{
		{
			name:      "Compound",
			category:  "Compound",
			wantError: false,
		},
		{
			name:      "Isolation",
			category:  "Isolation",
			wantError: false,
		},
		{
			name:      "Cardio",
			category:  "Cardio",
			wantError: false,
		},
		{
			name:      "無効なカテゴリ",
			category:  "Invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			category, err := NewExerciseCategory(tt.category)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.category, category.Category())
				assert.Equal(t, tt.category, category.String())
			}
		})
	}
}

func TestExercise_NewExercise(t *testing.T) {
	// Act
	exercise := NewExercise(BenchPress, Compound)

	// Assert
	assert.True(t, exercise.Name().Equals(BenchPress))
	assert.True(t, exercise.Category().Equals(Compound))
	assert.Equal(t, 0, exercise.SetCount())
	assert.Empty(t, exercise.Sets())
}

func TestExercise_AddSet(t *testing.T) {
	// Arrange
	exercise := NewExercise(BenchPress, Compound)
	weight, _ := NewWeight(95.0)
	reps, _ := NewReps(8)
	restTime, _ := NewRestTime(3 * time.Minute)
	set := NewSet(weight, reps, restTime, nil)

	// Act
	exercise.AddSet(set)

	// Assert
	assert.Equal(t, 1, exercise.SetCount())
	sets := exercise.Sets()
	assert.Len(t, sets, 1)
	assert.True(t, sets[0].Weight().Equals(weight))
}

func TestExercise_MaxWeight(t *testing.T) {
	// Arrange
	exercise := NewExercise(BenchPress, Compound)

	// セットなしの場合
	t.Run("セットが空の場合", func(t *testing.T) {
		// Act
		_, err := exercise.MaxWeight()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no sets recorded")
	})

	// セットありの場合
	t.Run("複数セットがある場合", func(t *testing.T) {
		// Arrange
		weight1, _ := NewWeight(90.0)
		weight2, _ := NewWeight(95.0)
		weight3, _ := NewWeight(92.5)
		reps, _ := NewReps(8)
		restTime, _ := NewRestTime(3 * time.Minute)

		exercise.AddSet(NewSet(weight1, reps, restTime, nil))
		exercise.AddSet(NewSet(weight2, reps, restTime, nil))
		exercise.AddSet(NewSet(weight3, reps, restTime, nil))

		// Act
		maxWeight, err := exercise.MaxWeight()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 95.0, maxWeight.Kg())
	})
}

func TestExercise_TotalVolume(t *testing.T) {
	// Arrange
	exercise := NewExercise(BenchPress, Compound)

	weight1, _ := NewWeight(95.0)
	reps1, _ := NewReps(8)
	restTime, _ := NewRestTime(3 * time.Minute)

	weight2, _ := NewWeight(95.0)
	reps2, _ := NewReps(6)

	exercise.AddSet(NewSet(weight1, reps1, restTime, nil))
	exercise.AddSet(NewSet(weight2, reps2, restTime, nil))

	// Act
	volume := exercise.TotalVolume()

	// Assert
	expected := 95.0*8 + 95.0*6 // 760 + 570 = 1330
	assert.Equal(t, expected, volume)
}
