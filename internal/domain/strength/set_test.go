package strength

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// セット実行コンテキストのテスト
// =============================================================================

func TestWeight_NewWeight(t *testing.T) {
	tests := []struct {
		name      string
		kg        float64
		wantError bool
	}{
		{
			name:      "有効な重量",
			kg:        95.5,
			wantError: false,
		},
		{
			name:      "0kg",
			kg:        0,
			wantError: false,
		},
		{
			name:      "負の重量",
			kg:        -10,
			wantError: true,
		},
		{
			name:      "過大な重量",
			kg:        1001,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			weight, err := NewWeight(tt.kg)

			// Assert
			if tt.wantError {
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
	// Arrange
	weight, _ := NewWeight(95.0)
	reps, _ := NewReps(8)
	restTime, _ := NewRestTime(3 * time.Minute)
	rpe, _ := NewRPE(8)

	// Act
	set := NewSet(weight, reps, restTime, &rpe)

	// Assert
	assert.True(t, set.Weight().Equals(weight))
	assert.True(t, set.Reps().Equals(reps))
	assert.True(t, set.RestTime().Equals(restTime))
	assert.NotNil(t, set.RPE())
	assert.Equal(t, 8, set.RPE().Rating())
}

func TestSet_String(t *testing.T) {
	// Arrange
	weight, _ := NewWeight(95.0)
	reps, _ := NewReps(8)
	restTime, _ := NewRestTime(3 * time.Minute)
	rpe, _ := NewRPE(8)

	tests := []struct {
		name string
		set  Set
		want string
	}{
		{
			name: "RPEありのセット",
			set:  NewSet(weight, reps, restTime, &rpe),
			want: "95.0kg × 8回 (休憩: 3m0s) RPE 8",
		},
		{
			name: "RPEなしのセット",
			set:  NewSet(weight, reps, restTime, nil),
			want: "95.0kg × 8回 (休憩: 3m0s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.set.String()

			// Assert
			assert.Equal(t, tt.want, result)
		})
	}
}
