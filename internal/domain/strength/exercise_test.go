package strength

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// エクササイズ定義コンテキストのテスト
// =============================================================================

func TestExerciseName_NewExerciseName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "正常なエクササイズ名",
			input:     "ベンチプレス",
			wantError: false,
			expected:  "ベンチプレス",
		},
		{
			name:      "空文字列",
			input:     "",
			wantError: true,
		},
		{
			name:      "カスタムエクササイズ名",
			input:     "インクラインベンチプレス",
			wantError: false,
			expected:  "インクラインベンチプレス",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given - テストデータは上記で定義済み

			// When
			result, err := NewExerciseName(tt.input)

			// Then
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "exercise name cannot be empty")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result.String())
				assert.Equal(t, tt.expected, result.Name())
			}
		})
	}
}

func TestExerciseName_Equals(t *testing.T) {
	tests := []struct {
		name     string
		name1    ExerciseName
		name2    ExerciseName
		expected bool
	}{
		{
			name:     "同じ名前",
			name1:    BenchPress,
			name2:    BenchPress,
			expected: true,
		},
		{
			name:     "異なる名前",
			name1:    BenchPress,
			name2:    Squat,
			expected: false,
		},
		{
			name:     "カスタムエクササイズ名の比較",
			name1:    func() ExerciseName { name, _ := NewExerciseName("カール"); return name }(),
			name2:    func() ExerciseName { name, _ := NewExerciseName("カール"); return name }(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given & When & Then
			assert.Equal(t, tt.expected, tt.name1.Equals(tt.name2))
		})
	}
}

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
		{
			name:     "カスタムエクササイズ",
			exercise: func() ExerciseName { name, _ := NewExerciseName("カール"); return name }(),
			isBIG3:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given & When & Then
			assert.Equal(t, tt.isBIG3, tt.exercise.IsBIG3())
		})
	}
}

// TODO: ExerciseCategory関連のテストは削除予定
// func TestExerciseCategory_NewExerciseCategory(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		category  string
// 		wantError bool
// 	}{
// 		{
// 			name:      "Compound",
// 			category:  "Compound",
// 			wantError: false,
// 		},
// 		{
// 			name:      "Isolation",
// 			category:  "Isolation",
// 			wantError: false,
// 		},
// 		{
// 			name:      "Cardio",
// 			category:  "Cardio",
// 			wantError: false,
// 		},
// 		{
// 			name:      "無効なカテゴリ",
// 			category:  "Invalid",
// 			wantError: true,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Act
// 			category, err := NewExerciseCategory(tt.category)
//
// 			// Assert
// 			if tt.wantError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.category, category.Category())
// 				assert.Equal(t, tt.category, category.String())
// 			}
// 		})
// 	}
// }

func TestExercise_NewExercise(t *testing.T) {
	tests := []struct {
		name         string
		exerciseName ExerciseName
	}{
		{
			name:         "ベンチプレス作成",
			exerciseName: BenchPress,
		},
		{
			name:         "スクワット作成",
			exerciseName: Squat,
		},
		{
			name:         "カスタムエクササイズ作成",
			exerciseName: func() ExerciseName { name, _ := NewExerciseName("インクラインベンチプレス"); return name }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given - テストデータは上記で定義済み

			// When
			exercise := NewExercise(tt.exerciseName)

			// Then
			assert.True(t, exercise.Name().Equals(tt.exerciseName))
			assert.Equal(t, 0, exercise.SetCount())
			assert.Empty(t, exercise.Sets())
		})
	}
}

func TestExercise_AddSet(t *testing.T) {
	tests := []struct {
		name         string
		exerciseName ExerciseName
		weight       float64
		reps         int
	}{
		{
			name:         "ベンチプレスセット追加",
			exerciseName: BenchPress,
			weight:       95.0,
			reps:         8,
		},
		{
			name:         "スクワットセット追加",
			exerciseName: Squat,
			weight:       100.0,
			reps:         5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			exercise := NewExercise(tt.exerciseName)
			weight, _ := NewWeight(tt.weight)
			reps, _ := NewReps(tt.reps)
			set := NewSet(weight, reps, nil)

			// When
			exercise.AddSet(set)

			// Then
			assert.Equal(t, 1, exercise.SetCount())
			sets := exercise.Sets()
			assert.Len(t, sets, 1)
			assert.True(t, sets[0].Weight().Equals(weight))
			assert.True(t, sets[0].Reps().Equals(reps))
		})
	}
}

func TestExercise_MaxWeight(t *testing.T) {
	t.Run("セットが空の場合", func(t *testing.T) {
		// Given
		exercise := NewExercise(BenchPress)

		// When
		_, err := exercise.MaxWeight()

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no sets recorded")
	})

	t.Run("複数セットがある場合", func(t *testing.T) {
		// Given
		exercise := NewExercise(BenchPress)
		weight1, _ := NewWeight(90.0)
		weight2, _ := NewWeight(95.0)
		weight3, _ := NewWeight(92.5)
		reps, _ := NewReps(8)

		exercise.AddSet(NewSet(weight1, reps, nil))
		exercise.AddSet(NewSet(weight2, reps, nil))
		exercise.AddSet(NewSet(weight3, reps, nil))

		// When
		maxWeight, err := exercise.MaxWeight()

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 95.0, maxWeight.Kg())
	})
}

func TestExercise_TotalVolume(t *testing.T) {
	tests := []struct {
		name     string
		sets     []struct{ weight, reps float64 }
		expected float64
	}{
		{
			name:     "セットなし",
			sets:     []struct{ weight, reps float64 }{},
			expected: 0.0,
		},
		{
			name: "単一セット",
			sets: []struct{ weight, reps float64 }{
				{weight: 95.0, reps: 8},
			},
			expected: 760.0, // 95 * 8
		},
		{
			name: "複数セット",
			sets: []struct{ weight, reps float64 }{
				{weight: 95.0, reps: 8},
				{weight: 95.0, reps: 6},
				{weight: 90.0, reps: 10},
			},
			expected: 2230.0, // 760 + 570 + 900 = 2230
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			exercise := NewExercise(BenchPress)
			for _, setData := range tt.sets {
				weight, _ := NewWeight(setData.weight)
				reps, _ := NewReps(int(setData.reps))
				exercise.AddSet(NewSet(weight, reps, nil))
			}

			// When
			volume := exercise.TotalVolume()

			// Then
			assert.Equal(t, tt.expected, volume)
		})
	}
}

func TestExercise_String(t *testing.T) {
	tests := []struct {
		name     string
		exercise func() *Exercise
		expected string
	}{
		{
			name: "セットなしのエクササイズ",
			exercise: func() *Exercise {
				return NewExercise(BenchPress)
			},
			expected: "ベンチプレス - 0セット",
		},
		{
			name: "セットありのエクササイズ",
			exercise: func() *Exercise {
				ex := NewExercise(BenchPress)
				weight, _ := NewWeight(95.0)
				reps, _ := NewReps(8)
				ex.AddSet(NewSet(weight, reps, nil))
				ex.AddSet(NewSet(weight, reps, nil))
				return ex
			},
			expected: "ベンチプレス - 2セット",
		},
		{
			name: "カスタムエクササイズ",
			exercise: func() *Exercise {
				name, _ := NewExerciseName("インクラインベンチプレス")
				ex := NewExercise(name)
				weight, _ := NewWeight(80.0)
				reps, _ := NewReps(10)
				ex.AddSet(NewSet(weight, reps, nil))
				return ex
			},
			expected: "インクラインベンチプレス - 1セット",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			exercise := tt.exercise()

			// When
			result := exercise.String()

			// Then
			assert.Equal(t, tt.expected, result)
		})
	}
}
