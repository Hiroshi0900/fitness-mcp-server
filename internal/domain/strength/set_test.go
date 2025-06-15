package strength

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// セット実行コンテキストのテスト
// =============================================================================

func TestWeight_NewWeight(t *testing.T) {
	tests := []struct {
		name        string
		kg          float64
		exceptError bool
	}{
		{
			name:        "正常系:有効な重量",
			kg:          95.5,
			exceptError: false,
		},
		{
			name:        "正常系:0kg",
			kg:          0,
			exceptError: false,
		},
		{
			name:        "異常系:負の重量",
			kg:          -10,
			exceptError: true,
		},
		{
			name:        "異常系:過大な重量",
			kg:          1001,
			exceptError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			weight, err := NewWeight(tt.kg)

			// Assert
			if tt.exceptError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.kg, weight.Kg())
				assert.Contains(t, weight.String(), "kg")
			}
		})
	}
}

func TestWeight_Equals(t *testing.T) {
	// Arrange
	weight1, _ := NewWeight(95.5)
	weight2, _ := NewWeight(95.5)
	weight3, _ := NewWeight(100.0)

	// Act & Assert
	assert.True(t, weight1.Equals(weight2))
	assert.False(t, weight1.Equals(weight3))
}

func TestReps_NewReps(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		wantError bool
	}{
		{
			name:      "有効な反復回数",
			count:     8,
			wantError: false,
		},
		{
			name:      "0回",
			count:     0,
			wantError: true,
		},
		{
			name:      "負の回数",
			count:     -1,
			wantError: true,
		},
		{
			name:      "過大な回数",
			count:     501,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			reps, err := NewReps(tt.count)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.count, reps.Count())
				assert.Contains(t, reps.String(), "回")
			}
		})
	}
}

func TestRPE_NewRPE(t *testing.T) {
	tests := []struct {
		name      string
		rating    int
		wantError bool
	}{
		{
			name:      "有効なRPE",
			rating:    8,
			wantError: false,
		},
		{
			name:      "最小値",
			rating:    1,
			wantError: false,
		},
		{
			name:      "最大値",
			rating:    10,
			wantError: false,
		},
		{
			name:      "範囲外(低)",
			rating:    0,
			wantError: true,
		},
		{
			name:      "範囲外(高)",
			rating:    11,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			rpe, err := NewRPE(tt.rating)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.rating, rpe.Rating())
				assert.Contains(t, rpe.String(), "RPE")
			}
		})
	}
}

func TestSet_NewSet(t *testing.T) {
	tests := []struct {
		name      string
		weight    float64
		reps      int
		rpeRating *int
		expectRPE bool
	}{
		{
			name:      "RPEありのセット作成",
			weight:    95.0,
			reps:      8,
			rpeRating: &[]int{8}[0],
			expectRPE: true,
		},
		{
			name:      "RPEなしのセット作成",
			weight:    100.0,
			reps:      5,
			rpeRating: nil,
			expectRPE: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			weight, _ := NewWeight(tt.weight)
			reps, _ := NewReps(tt.reps)
			var rpe *RPE
			if tt.rpeRating != nil {
				rpeVal, _ := NewRPE(*tt.rpeRating)
				rpe = &rpeVal
			}

			// When
			set := NewSet(weight, reps, rpe)

			// Then
			assert.True(t, set.Weight().Equals(weight))
			assert.True(t, set.Reps().Equals(reps))
			if tt.expectRPE {
				assert.NotNil(t, set.RPE())
				assert.Equal(t, *tt.rpeRating, set.RPE().Rating())
			} else {
				assert.Nil(t, set.RPE())
			}
		})
	}
}

func TestSet_String(t *testing.T) {
	tests := []struct {
		name      string
		weight    float64
		reps      int
		rpeRating *int
		expected  string
	}{
		{
			name:      "RPEありのセット文字列表示",
			weight:    95.0,
			reps:      8,
			rpeRating: &[]int{8}[0],
			expected:  "95.0kg × 8回 RPE 8",
		},
		{
			name:      "RPEなしのセット文字列表示",
			weight:    100.0,
			reps:      5,
			rpeRating: nil,
			expected:  "100.0kg × 5回",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			weight, _ := NewWeight(tt.weight)
			reps, _ := NewReps(tt.reps)
			var rpe *RPE
			if tt.rpeRating != nil {
				rpeVal, _ := NewRPE(*tt.rpeRating)
				rpe = &rpeVal
			}
			set := NewSet(weight, reps, rpe)

			// When
			result := set.String()

			// Then
			assert.Equal(t, tt.expected, result)
		})
	}
}
