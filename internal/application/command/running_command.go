package command

import (
	"fmt"
	"time"

	"fitness-mcp-server/internal/domain/running"
	"fitness-mcp-server/internal/domain/shared"
)

// =============================================================================
// ランニングコマンド定義 - 状態変更操作
// =============================================================================

// RecordRunningSessionCommand はランニングセッション記録コマンド
type RecordRunningSessionCommand struct {
	Date         time.Time `json:"date"`
	DistanceKm   float64   `json:"distance_km"`
	DurationMin  float64   `json:"duration_minutes"`
	RunType      string    `json:"run_type"`
	HeartRateBPM *int      `json:"heart_rate_bpm,omitempty"` // オプション
	Notes        string    `json:"notes"`
}

// UpdateRunningSessionCommand はランニングセッション更新コマンド
type UpdateRunningSessionCommand struct {
	ID           string    `json:"id"`
	Date         time.Time `json:"date"`
	DistanceKm   float64   `json:"distance_km"`
	DurationMin  float64   `json:"duration_minutes"`
	RunType      string    `json:"run_type"`
	HeartRateBPM *int      `json:"heart_rate_bpm,omitempty"` // オプション
	Notes        string    `json:"notes"`
}

// DeleteRunningSessionCommand はランニングセッション削除コマンド
type DeleteRunningSessionCommand struct {
	ID string `json:"id"`
}

// CreateRunningGoalCommand はランニング目標作成コマンド
type CreateRunningGoalCommand struct {
	EventType       string     `json:"event_type"`        // "5K", "10K", "Half", "Marathon", "Custom"
	TargetTimeMin   float64    `json:"target_time_minutes"`
	CustomDistanceKm *float64  `json:"custom_distance_km,omitempty"` // Custom時のみ
	EventDate       *time.Time `json:"event_date,omitempty"`         // オプション
	Description     string     `json:"description"`
}

// UpdateRunningGoalCommand はランニング目標更新コマンド
type UpdateRunningGoalCommand struct {
	ID          string     `json:"id"`
	EventDate   *time.Time `json:"event_date,omitempty"`
	Description string     `json:"description"`
}

// DeleteRunningGoalCommand はランニング目標削除コマンド
type DeleteRunningGoalCommand struct {
	ID string `json:"id"`
}

// Result構造体群

// RecordRunningSessionResult はランニングセッション記録結果
type RecordRunningSessionResult struct {
	SessionID    string        `json:"session_id"`
	Date         time.Time     `json:"date"`
	Distance     string        `json:"distance"`     // "5.00km"
	Duration     string        `json:"duration"`     // "25:30"
	Pace         string        `json:"pace"`         // "5:06/km"
	Message      string        `json:"message"`
	GoalAchieved *GoalAchieved `json:"goal_achieved,omitempty"` // 目標達成時
}

// UpdateRunningSessionResult はランニングセッション更新結果
type UpdateRunningSessionResult struct {
	SessionID string    `json:"session_id"`
	Date      time.Time `json:"date"`
	Distance  string    `json:"distance"`
	Duration  string    `json:"duration"`
	Pace      string    `json:"pace"`
	Message   string    `json:"message"`
}

// DeleteRunningSessionResult はランニングセッション削除結果
type DeleteRunningSessionResult struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// CreateRunningGoalResult はランニング目標作成結果
type CreateRunningGoalResult struct {
	GoalID       string `json:"goal_id"`
	EventType    string `json:"event_type"`
	TargetTime   string `json:"target_time"`
	TargetPace   string `json:"target_pace"`
	Message      string `json:"message"`
}

// UpdateRunningGoalResult はランニング目標更新結果
type UpdateRunningGoalResult struct {
	GoalID  string `json:"goal_id"`
	Message string `json:"message"`
}

// DeleteRunningGoalResult はランニング目標削除結果
type DeleteRunningGoalResult struct {
	GoalID  string `json:"goal_id"`
	Message string `json:"message"`
}

// GoalAchieved は目標達成情報
type GoalAchieved struct {
	GoalID      string `json:"goal_id"`
	EventType   string `json:"event_type"`
	TargetTime  string `json:"target_time"`
	ActualTime  string `json:"actual_time"`
	Improvement string `json:"improvement"` // "目標より30秒早い"
}

// Validation methods

// Validate はRecordRunningSessionCommandの妥当性検証を行います
func (cmd *RecordRunningSessionCommand) Validate() error {
	if cmd.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	if cmd.DistanceKm <= 0 {
		return fmt.Errorf("distance must be positive")
	}
	if cmd.DurationMin <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	if cmd.RunType == "" {
		return fmt.Errorf("run type is required")
	}
	if cmd.HeartRateBPM != nil && (*cmd.HeartRateBPM <= 0 || *cmd.HeartRateBPM > 250) {
		return fmt.Errorf("heart rate must be between 1 and 250 bpm")
	}
	return nil
}

// Validate はUpdateRunningSessionCommandの妥当性検証を行います
func (cmd *UpdateRunningSessionCommand) Validate() error {
	if cmd.ID == "" {
		return fmt.Errorf("session ID is required")
	}
	if cmd.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	if cmd.DistanceKm <= 0 {
		return fmt.Errorf("distance must be positive")
	}
	if cmd.DurationMin <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	if cmd.RunType == "" {
		return fmt.Errorf("run type is required")
	}
	if cmd.HeartRateBPM != nil && (*cmd.HeartRateBPM <= 0 || *cmd.HeartRateBPM > 250) {
		return fmt.Errorf("heart rate must be between 1 and 250 bpm")
	}
	return nil
}

// Validate はDeleteRunningSessionCommandの妥当性検証を行います
func (cmd *DeleteRunningSessionCommand) Validate() error {
	if cmd.ID == "" {
		return fmt.Errorf("session ID is required")
	}
	return nil
}

// Validate はCreateRunningGoalCommandの妥当性検証を行います
func (cmd *CreateRunningGoalCommand) Validate() error {
	if cmd.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if cmd.TargetTimeMin <= 0 {
		return fmt.Errorf("target time must be positive")
	}
	if cmd.EventType == "Custom" && (cmd.CustomDistanceKm == nil || *cmd.CustomDistanceKm <= 0) {
		return fmt.Errorf("custom distance is required for custom event type")
	}
	if cmd.Description == "" {
		return fmt.Errorf("description is required")
	}
	return nil
}

// Validate はUpdateRunningGoalCommandの妥当性検証を行います
func (cmd *UpdateRunningGoalCommand) Validate() error {
	if cmd.ID == "" {
		return fmt.Errorf("goal ID is required")
	}
	return nil
}

// Validate はDeleteRunningGoalCommandの妥当性検証を行います
func (cmd *DeleteRunningGoalCommand) Validate() error {
	if cmd.ID == "" {
		return fmt.Errorf("goal ID is required")
	}
	return nil
}

// Domain conversion methods

// ToRunningSession はコマンドからRunningSessionエンティティを生成します
func (cmd *RecordRunningSessionCommand) ToRunningSession() (*running.RunningSession, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 距離を作成
	distance, err := running.NewDistance(cmd.DistanceKm)
	if err != nil {
		return nil, fmt.Errorf("invalid distance: %w", err)
	}

	// 時間を作成
	duration, err := running.NewDurationFromMinutes(cmd.DurationMin)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}

	// ランニングタイプを作成
	runType, err := running.NewRunType(cmd.RunType)
	if err != nil {
		return nil, fmt.Errorf("invalid run type: %w", err)
	}

	// 新しいセッションIDを生成
	sessionID := shared.NewSessionID()

	// セッションを作成
	session, err := running.NewRunningSession(sessionID, cmd.Date, distance, duration, runType, cmd.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// 心拍数を設定（オプション）
	if cmd.HeartRateBPM != nil {
		heartRate, err := running.NewHeartRate(*cmd.HeartRateBPM)
		if err != nil {
			return nil, fmt.Errorf("invalid heart rate: %w", err)
		}
		session.SetHeartRate(heartRate)
	}

	return session, nil
}

// ToRunningSession はUpdateRunningSessionCommandからRunningSessionエンティティを生成します
func (cmd *UpdateRunningSessionCommand) ToRunningSession() (*running.RunningSession, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// IDを復元
	sessionID, err := shared.NewSessionIDFromString(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 距離を作成
	distance, err := running.NewDistance(cmd.DistanceKm)
	if err != nil {
		return nil, fmt.Errorf("invalid distance: %w", err)
	}

	// 時間を作成
	duration, err := running.NewDurationFromMinutes(cmd.DurationMin)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}

	// ランニングタイプを作成
	runType, err := running.NewRunType(cmd.RunType)
	if err != nil {
		return nil, fmt.Errorf("invalid run type: %w", err)
	}

	// セッションを作成
	session, err := running.NewRunningSession(sessionID, cmd.Date, distance, duration, runType, cmd.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// 心拍数を設定（オプション）
	if cmd.HeartRateBPM != nil {
		heartRate, err := running.NewHeartRate(*cmd.HeartRateBPM)
		if err != nil {
			return nil, fmt.Errorf("invalid heart rate: %w", err)
		}
		session.SetHeartRate(heartRate)
	}

	return session, nil
}

// ToRunningGoal はコマンドからRunningGoalエンティティを生成します
func (cmd *CreateRunningGoalCommand) ToRunningGoal() (*running.RunningGoal, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 新しい目標IDを生成
	goalID := shared.NewGoalID()

	// 目標時間を作成
	targetTime, err := running.NewDurationFromMinutes(cmd.TargetTimeMin)
	if err != nil {
		return nil, fmt.Errorf("invalid target time: %w", err)
	}

	var goal *running.RunningGoal

	if cmd.EventType == "Custom" {
		// カスタム距離の目標
		distance, err := running.NewDistance(*cmd.CustomDistanceKm)
		if err != nil {
			return nil, fmt.Errorf("invalid custom distance: %w", err)
		}
		goal, err = running.NewCustomRunningGoal(goalID, distance, targetTime, cmd.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to create custom goal: %w", err)
		}
	} else {
		// 標準距離の目標
		eventType, err := running.NewEventType(cmd.EventType)
		if err != nil {
			return nil, fmt.Errorf("invalid event type: %w", err)
		}
		goal, err = running.NewRunningGoal(goalID, eventType, targetTime, cmd.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to create goal: %w", err)
		}
	}

	// イベント日を設定（オプション）
	if cmd.EventDate != nil {
		goal.SetEventDate(*cmd.EventDate)
	}

	return goal, nil
}
