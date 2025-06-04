package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrainingID_NewTrainingID(t *testing.T) {
	// Act
	id := NewTrainingID()

	// Assert
	assert.False(t, id.IsEmpty())
	assert.NotEqual(t, "", id.String())
	assert.Len(t, id.String(), 36) // UUID の長さ
}

func TestTrainingID_NewTrainingIDFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "有効なUUID",
			input:     "123e4567-e89b-12d3-a456-426614174000",
			wantError: false,
		},
		{
			name:      "空文字列",
			input:     "",
			wantError: true,
		},
		{
			name:      "無効なUUID形式",
			input:     "invalid-uuid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			id, err := NewTrainingIDFromString(tt.input)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.True(t, id.IsEmpty())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, id.String())
				assert.False(t, id.IsEmpty())
			}
		})
	}
}

func TestTrainingID_Equals(t *testing.T) {
	// Arrange
	id1 := NewTrainingID()
	id2 := NewTrainingID()
	id3, _ := NewTrainingIDFromString(id1.String())

	// Act & Assert
	assert.False(t, id1.Equals(id2)) // 異なるID
	assert.True(t, id1.Equals(id3))  // 同じID
}

func TestSessionID_NewSessionID(t *testing.T) {
	// Act
	id := NewSessionID()

	// Assert
	assert.False(t, id.IsEmpty())
	assert.NotEqual(t, "", id.String())
	assert.Len(t, id.String(), 36) // UUID の長さ
}

func TestSessionID_NewSessionIDFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "有効なUUID",
			input:     "123e4567-e89b-12d3-a456-426614174000",
			wantError: false,
		},
		{
			name:      "空文字列",
			input:     "",
			wantError: true,
		},
		{
			name:      "無効なUUID形式",
			input:     "invalid-uuid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			id, err := NewSessionIDFromString(tt.input)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.True(t, id.IsEmpty())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, id.String())
				assert.False(t, id.IsEmpty())
			}
		})
	}
}

func TestGoalID_NewGoalID(t *testing.T) {
	// Act
	id := NewGoalID()

	// Assert
	assert.False(t, id.IsEmpty())
	assert.NotEqual(t, "", id.String())
	assert.Len(t, id.String(), 36) // UUID の長さ
}

func TestGoalID_NewGoalIDFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "有効なUUID",
			input:     "123e4567-e89b-12d3-a456-426614174000",
			wantError: false,
		},
		{
			name:      "空文字列",
			input:     "",
			wantError: true,
		},
		{
			name:      "無効なUUID形式",
			input:     "invalid-uuid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			id, err := NewGoalIDFromString(tt.input)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.True(t, id.IsEmpty())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, id.String())
				assert.False(t, id.IsEmpty())
			}
		})
	}
}
