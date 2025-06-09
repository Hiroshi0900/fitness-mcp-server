package running

import (
	"fmt"
	"time"

	"fitness-mcp-server/internal/domain/shared"
)

// =============================================================================
// ランニング目標コンテキスト - 目標設定と進捗管理
// =============================================================================

// EventType はイベントの種類を表す値オブジェクト
type EventType struct {
	value string
}

// 定義済みイベントタイプの定数
var (
	FiveK      = EventType{value: "5K"}      // 5km
	TenK       = EventType{value: "10K"}     // 10km
	HalfMarathon = EventType{value: "Half"}  // ハーフマラソン
	Marathon   = EventType{value: "Marathon"} // フルマラソン
	Custom     = EventType{value: "Custom"}   // カスタム距離
)

// NewEventType はイベントタイプを作成します
func NewEventType(eventType string) (EventType, error) {
	validTypes := []string{"5K", "10K", "Half", "Marathon", "Custom"}
	for _, valid := range validTypes {
		if eventType == valid {
			return EventType{value: eventType}, nil
		}
	}
	return EventType{}, fmt.Errorf("invalid event type: %s", eventType)
}

// Type はイベントタイプを返します
func (et EventType) Type() string {
	return et.value
}

// String はイベントタイプの文字列表現を返します
func (et EventType) String() string {
	return et.value
}

// Equals は2つのイベントタイプが等しいかを判定します
func (et EventType) Equals(other EventType) bool {
	return et.value == other.value
}

// GetStandardDistance は標準的な距離を返します
func (et EventType) GetStandardDistance() (Distance, error) {
	switch et.value {
	case "5K":
		return NewDistance(5.0)
	case "10K":
		return NewDistance(10.0)
	case "Half":
		return NewDistance(21.0975) // ハーフマラソンの正確な距離
	case "Marathon":
		return NewDistance(42.195) // フルマラソンの正確な距離
	default:
		return Distance{}, fmt.Errorf("no standard distance for event type: %s", et.value)
	}
}

// GoalStatus は目標の状態を表す値オブジェクト
type GoalStatus struct {
	value string
}

// 定義済み目標状態の定数
var (
	Active    = GoalStatus{value: "Active"}    // アクティブ
	Achieved  = GoalStatus{value: "Achieved"}  // 達成済み
	Paused    = GoalStatus{value: "Paused"}    // 一時停止
	Cancelled = GoalStatus{value: "Cancelled"} // キャンセル
)

// NewGoalStatus は目標状態を作成します
func NewGoalStatus(status string) (GoalStatus, error) {
	validStatuses := []string{"Active", "Achieved", "Paused", "Cancelled"}
	for _, valid := range validStatuses {
		if status == valid {
			return GoalStatus{value: status}, nil
		}
	}
	return GoalStatus{}, fmt.Errorf("invalid goal status: %s", status)
}

// Status は目標状態を返します
func (gs GoalStatus) Status() string {
	return gs.value
}

// String は目標状態の文字列表現を返します
func (gs GoalStatus) String() string {
	return gs.value
}

// Equals は2つの目標状態が等しいかを判定します
func (gs GoalStatus) Equals(other GoalStatus) bool {
	return gs.value == other.value
}

// IsActive はアクティブな状態かを判定します
func (gs GoalStatus) IsActive() bool {
	return gs.Equals(Active)
}

// RunningGoal はランニング目標を表すエンティティ
type RunningGoal struct {
	id           shared.GoalID // 目標ID
	eventType    EventType     // イベントタイプ
	targetTime   Duration      // 目標タイム
	targetPace   Pace          // 目標ペース
	eventDate    *time.Time    // イベント日（オプション）
	status       GoalStatus    // 状態
	description  string        // 説明
	createdAt    time.Time     // 作成日時
	achievedAt   *time.Time    // 達成日時（オプション）
}

// NewRunningGoal は新しいRunningGoalを作成します
func NewRunningGoal(
	id shared.GoalID,
	eventType EventType,
	targetTime Duration,
	description string,
) (*RunningGoal, error) {
	// 標準距離を取得
	distance, err := eventType.GetStandardDistance()
	if err != nil {
		return nil, fmt.Errorf("failed to get standard distance: %w", err)
	}

	// 目標ペースを計算
	targetPace, err := CalculatePace(distance, targetTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate target pace: %w", err)
	}

	return &RunningGoal{
		id:          id,
		eventType:   eventType,
		targetTime:  targetTime,
		targetPace:  targetPace,
		eventDate:   nil,
		status:      Active,
		description: description,
		createdAt:   time.Now(),
		achievedAt:  nil,
	}, nil
}

// NewCustomRunningGoal はカスタム距離の目標を作成します
func NewCustomRunningGoal(
	id shared.GoalID,
	distance Distance,
	targetTime Duration,
	description string,
) (*RunningGoal, error) {
	// 目標ペースを計算
	targetPace, err := CalculatePace(distance, targetTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate target pace: %w", err)
	}

	return &RunningGoal{
		id:          id,
		eventType:   Custom,
		targetTime:  targetTime,
		targetPace:  targetPace,
		eventDate:   nil,
		status:      Active,
		description: description,
		createdAt:   time.Now(),
		achievedAt:  nil,
	}, nil
}

// ID は目標IDを返します
func (rg *RunningGoal) ID() shared.GoalID {
	return rg.id
}

// EventType はイベントタイプを返します
func (rg *RunningGoal) EventType() EventType {
	return rg.eventType
}

// TargetTime は目標タイムを返します
func (rg *RunningGoal) TargetTime() Duration {
	return rg.targetTime
}

// TargetPace は目標ペースを返します
func (rg *RunningGoal) TargetPace() Pace {
	return rg.targetPace
}

// EventDate はイベント日を返します（オプション）
func (rg *RunningGoal) EventDate() *time.Time {
	return rg.eventDate
}

// Status は状態を返します
func (rg *RunningGoal) Status() GoalStatus {
	return rg.status
}

// Description は説明を返します
func (rg *RunningGoal) Description() string {
	return rg.description
}

// CreatedAt は作成日時を返します
func (rg *RunningGoal) CreatedAt() time.Time {
	return rg.createdAt
}

// AchievedAt は達成日時を返します（オプション）
func (rg *RunningGoal) AchievedAt() *time.Time {
	return rg.achievedAt
}

// SetEventDate はイベント日を設定します
func (rg *RunningGoal) SetEventDate(eventDate time.Time) {
	rg.eventDate = &eventDate
}

// UpdateDescription は説明を更新します
func (rg *RunningGoal) UpdateDescription(description string) {
	rg.description = description
}

// MarkAsAchieved は目標を達成済みにマークします
func (rg *RunningGoal) MarkAsAchieved() {
	rg.status = Achieved
	now := time.Now()
	rg.achievedAt = &now
}

// MarkAsPaused は目標を一時停止にマークします
func (rg *RunningGoal) MarkAsPaused() {
	rg.status = Paused
}

// MarkAsCancelled は目標をキャンセルにマークします
func (rg *RunningGoal) MarkAsCancelled() {
	rg.status = Cancelled
}

// Resume は目標を再開します
func (rg *RunningGoal) Resume() {
	if rg.status.Equals(Paused) {
		rg.status = Active
	}
}

// IsAchievable は指定されたセッションで目標が達成可能かを判定します
func (rg *RunningGoal) IsAchievable(session *RunningSession) (bool, error) {
	// イベントタイプの標準距離を取得
	goalDistance, err := rg.eventType.GetStandardDistance()
	if err != nil {
		return false, fmt.Errorf("failed to get goal distance: %w", err)
	}

	// 距離が一致するかチェック
	if !session.Distance().Equals(goalDistance) {
		return false, nil
	}

	// セッションのタイムが目標タイム以下かチェック
	return session.Duration().Value() <= rg.targetTime.Value(), nil
}

// CheckAchievement はセッションで目標が達成されたかをチェックし、必要に応じて状態を更新します
func (rg *RunningGoal) CheckAchievement(session *RunningSession) error {
	if !rg.status.IsActive() {
		return nil // アクティブでない場合はチェックしない
	}

	achievable, err := rg.IsAchievable(session)
	if err != nil {
		return fmt.Errorf("failed to check if achievable: %w", err)
	}

	if achievable {
		rg.MarkAsAchieved()
	}

	return nil
}

// DaysUntilEvent はイベントまでの日数を返します
func (rg *RunningGoal) DaysUntilEvent() *int {
	if rg.eventDate == nil {
		return nil
	}

	days := int(rg.eventDate.Sub(time.Now()).Hours() / 24)
	return &days
}

// String は目標の文字列表現を返します
func (rg *RunningGoal) String() string {
	eventDateStr := "未設定"
	if rg.eventDate != nil {
		eventDateStr = rg.eventDate.Format("2006-01-02")
	}

	return fmt.Sprintf("目標 %s (%s) - %s: %s, ペース: %s, イベント日: %s",
		rg.id.String()[:8], rg.status.String(),
		rg.eventType.String(), rg.targetTime.String(), rg.targetPace.String(),
		eventDateStr)
}
