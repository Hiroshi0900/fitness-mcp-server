package shared

import (
	"fmt"

	"github.com/google/uuid"
)

// TrainingID は筋トレセッションを一意に識別するID
type TrainingID struct {
	value string
}

// NewTrainingID は新しいTrainingIDを生成します
func NewTrainingID() TrainingID {
	return TrainingID{value: uuid.New().String()}
}

// NewTrainingIDFromString は文字列からTrainingIDを作成します
func NewTrainingIDFromString(s string) (TrainingID, error) {
	if s == "" {
		return TrainingID{}, fmt.Errorf("id cannot be empty")
	}
	// UUIDの形式チェック
	if _, err := uuid.Parse(s); err != nil {
		return TrainingID{}, fmt.Errorf("invalid uuid format: %w", err)
	}
	return TrainingID{value: s}, nil
}

// String はIDの文字列表現を返します
func (id TrainingID) String() string {
	return id.value
}

// IsEmpty はIDが空かどうかを判定します
func (id TrainingID) IsEmpty() bool {
	return id.value == ""
}

// Equals は2つのIDが等しいかを判定します
func (id TrainingID) Equals(other TrainingID) bool {
	return id.value == other.value
}

// SessionID はランニングセッションを一意に識別するID
type SessionID struct {
	value string
}

// NewSessionID は新しいSessionIDを生成します
func NewSessionID() SessionID {
	return SessionID{value: uuid.New().String()}
}

// NewSessionIDFromString は文字列からSessionIDを作成します
func NewSessionIDFromString(s string) (SessionID, error) {
	if s == "" {
		return SessionID{}, fmt.Errorf("id cannot be empty")
	}
	if _, err := uuid.Parse(s); err != nil {
		return SessionID{}, fmt.Errorf("invalid uuid format: %w", err)
	}
	return SessionID{value: s}, nil
}

// String はIDの文字列表現を返します
func (id SessionID) String() string {
	return id.value
}

// IsEmpty はIDが空かどうかを判定します
func (id SessionID) IsEmpty() bool {
	return id.value == ""
}

// Equals は2つのIDが等しいかを判定します
func (id SessionID) Equals(other SessionID) bool {
	return id.value == other.value
}

// GoalID は目標を一意に識別するID
type GoalID struct {
	value string
}

// NewGoalID は新しいGoalIDを生成します
func NewGoalID() GoalID {
	return GoalID{value: uuid.New().String()}
}

// NewGoalIDFromString は文字列からGoalIDを作成します
func NewGoalIDFromString(s string) (GoalID, error) {
	if s == "" {
		return GoalID{}, fmt.Errorf("id cannot be empty")
	}
	if _, err := uuid.Parse(s); err != nil {
		return GoalID{}, fmt.Errorf("invalid uuid format: %w", err)
	}
	return GoalID{value: s}, nil
}

// String はIDの文字列表現を返します
func (id GoalID) String() string {
	return id.value
}

// IsEmpty はIDが空かどうかを判定します
func (id GoalID) IsEmpty() bool {
	return id.value == ""
}

// Equals は2つのIDが等しいかを判定します
func (id GoalID) Equals(other GoalID) bool {
	return id.value == other.value
}
