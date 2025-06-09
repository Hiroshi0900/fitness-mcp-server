package running

import (
	"fmt"
	"time"

	"fitness-mcp-server/internal/domain/shared"
)

// =============================================================================
// ランニングセッションコンテキスト - 1回のランニング全体の管理
// =============================================================================

// RunType はランニングの種類を表す値オブジェクト
type RunType struct {
	value string
}

// 定義済みランニングタイプの定数
var (
	Easy     = RunType{value: "Easy"}     // イージーラン
	Tempo    = RunType{value: "Tempo"}    // テンポラン
	Interval = RunType{value: "Interval"} // インターバル
	Long     = RunType{value: "Long"}     // ロングラン
	Race     = RunType{value: "Race"}     // レース
)

// NewRunType はランニングタイプを作成します
func NewRunType(runType string) (RunType, error) {
	validTypes := []string{"Easy", "Tempo", "Interval", "Long", "Race"}
	for _, valid := range validTypes {
		if runType == valid {
			return RunType{value: runType}, nil
		}
	}
	return RunType{}, fmt.Errorf("invalid run type: %s", runType)
}

// Type はランニングタイプを返します
func (rt RunType) Type() string {
	return rt.value
}

// String はランニングタイプの文字列表現を返します
func (rt RunType) String() string {
	return rt.value
}

// Equals は2つのランニングタイプが等しいかを判定します
func (rt RunType) Equals(other RunType) bool {
	return rt.value == other.value
}

// RunningSession はランニングセッションを表すエンティティ
type RunningSession struct {
	id        shared.SessionID // セッションID
	date      time.Time        // ランニング日
	distance  Distance         // 距離
	duration  Duration         // 時間
	pace      Pace             // ペース
	heartRate *HeartRate       // 心拍数（オプション）
	runType   RunType          // ランニングタイプ
	notes     string           // メモ
}

// NewRunningSession は新しいRunningSessionを作成します
func NewRunningSession(
	id shared.SessionID,
	date time.Time,
	distance Distance,
	duration Duration,
	runType RunType,
	notes string,
) (*RunningSession, error) {
	// ペースを計算
	pace, err := CalculatePace(distance, duration)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate pace: %w", err)
	}

	return &RunningSession{
		id:        id,
		date:      date,
		distance:  distance,
		duration:  duration,
		pace:      pace,
		heartRate: nil,
		runType:   runType,
		notes:     notes,
	}, nil
}

// ID はセッションIDを返します
func (rs *RunningSession) ID() shared.SessionID {
	return rs.id
}

// Date はランニング日を返します
func (rs *RunningSession) Date() time.Time {
	return rs.date
}

// Distance は距離を返します
func (rs *RunningSession) Distance() Distance {
	return rs.distance
}

// Duration は時間を返します
func (rs *RunningSession) Duration() Duration {
	return rs.duration
}

// Pace はペースを返します
func (rs *RunningSession) Pace() Pace {
	return rs.pace
}

// HeartRate は心拍数を返します（オプション）
func (rs *RunningSession) HeartRate() *HeartRate {
	return rs.heartRate
}

// RunType はランニングタイプを返します
func (rs *RunningSession) RunType() RunType {
	return rs.runType
}

// Notes はメモを返します
func (rs *RunningSession) Notes() string {
	return rs.notes
}

// SetHeartRate は心拍数を設定します
func (rs *RunningSession) SetHeartRate(heartRate HeartRate) {
	rs.heartRate = &heartRate
}

// UpdateNotes はメモを更新します
func (rs *RunningSession) UpdateNotes(notes string) {
	rs.notes = notes
}

// UpdateDistance は距離と時間を更新し、ペースを再計算します
func (rs *RunningSession) UpdateDistance(distance Distance, duration Duration) error {
	pace, err := CalculatePace(distance, duration)
	if err != nil {
		return fmt.Errorf("failed to calculate pace: %w", err)
	}
	
	rs.distance = distance
	rs.duration = duration
	rs.pace = pace
	return nil
}

// String はランニングセッションの文字列表現を返します
func (rs *RunningSession) String() string {
	return fmt.Sprintf("ランニング %s (%s) - %s, %s, ペース: %s",
		rs.id.String()[:8], rs.date.Format("2006-01-02"),
		rs.distance.String(), rs.duration.String(), rs.pace.String())
}
