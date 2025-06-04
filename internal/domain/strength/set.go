package strength

import (
	"fmt"
	"time"
)

// =============================================================================
// セット実行コンテキスト - 測定値と1回のセット
// =============================================================================

type (
	Weight struct {
		value float64 // kg単位
	}

	Reps struct {
		value int
	}

	RestTime struct {
		duration time.Duration
	}
	// RPE は主観的運動強度を表す値オブジェクト
	RPE struct {
		value int
	}
	Set struct {
		weight   Weight   // 重量
		reps     Reps     // 反復回数
		restTime RestTime // 休憩時間
		rpe      *RPE     // オプショナル
	}
)

// NewWeight は重量を作成します
func NewWeight(kg float64) (Weight, error) {
	if kg < 0 {
		return Weight{}, fmt.Errorf("weight cannot be negative: %f", kg)
	}
	if kg > 1000 { // 現実的な上限設定
		return Weight{}, fmt.Errorf("weight is too large: %f", kg)
	}
	return Weight{value: kg}, nil
}

// Kg は重量をkg単位で返します
func (w Weight) Kg() float64 {
	return w.value
}

// String は重量の文字列表現を返します
func (w Weight) String() string {
	return fmt.Sprintf("%.1fkg", w.value)
}

// Equals は2つの重量が等しいかを判定します
func (w Weight) Equals(other Weight) bool {
	return w.value == other.value
}

// NewReps は反復回数を作成します
func NewReps(count int) (Reps, error) {
	if count <= 0 {
		return Reps{}, fmt.Errorf("reps must be positive: %d", count)
	}
	if count > 500 { // 現実的な上限設定
		return Reps{}, fmt.Errorf("reps is too large: %d", count)
	}
	return Reps{value: count}, nil
}

// Count は反復回数を返します
func (r Reps) Count() int {
	return r.value
}

// String は反復回数の文字列表現を返します
func (r Reps) String() string {
	return fmt.Sprintf("%d回", r.value)
}

// Equals は2つの反復回数が等しいかを判定します
func (r Reps) Equals(other Reps) bool {
	return r.value == other.value
}

// NewRestTime は休憩時間を作成します
func NewRestTime(d time.Duration) (RestTime, error) {
	if d < 0 {
		return RestTime{}, fmt.Errorf("rest time cannot be negative: %v", d)
	}
	if d > 30*time.Minute { // 現実的な上限設定
		return RestTime{}, fmt.Errorf("rest time is too long: %v", d)
	}
	return RestTime{duration: d}, nil
}

// Duration は休憩時間を返します
func (rt RestTime) Duration() time.Duration {
	return rt.duration
}

// String は休憩時間の文字列表現を返します
func (rt RestTime) String() string {
	return rt.duration.String()
}

// Equals は2つの休憩時間が等しいかを判定します
func (rt RestTime) Equals(other RestTime) bool {
	return rt.duration == other.duration
}

// NewRPE は主観的運動強度を作成します
func NewRPE(rating int) (RPE, error) {
	if rating < 1 || rating > 10 {
		return RPE{}, fmt.Errorf("RPE must be between 1 and 10: %d", rating)
	}
	return RPE{value: rating}, nil
}

// Rating は主観的運動強度の値を返します
func (rpe RPE) Rating() int {
	return rpe.value
}

// String は主観的運動強度の文字列表現を返します
func (rpe RPE) String() string {
	return fmt.Sprintf("RPE %d", rpe.value)
}

// Equals は2つの主観的運動強度が等しいかを判定します
func (rpe RPE) Equals(other RPE) bool {
	return rpe.value == other.value
}

// NewSet は新しいSetを作成します
func NewSet(weight Weight, reps Reps, restTime RestTime, rpe *RPE) Set {
	return Set{
		weight:   weight,
		reps:     reps,
		restTime: restTime,
		rpe:      rpe,
	}
}

// Weight は重量を返します
func (s Set) Weight() Weight {
	return s.weight
}

// Reps は反復回数を返します
func (s Set) Reps() Reps {
	return s.reps
}

// RestTime は休憩時間を返します
func (s Set) RestTime() RestTime {
	return s.restTime
}

// RPE は主観的運動強度を返します（オプショナル）
func (s Set) RPE() *RPE {
	return s.rpe
}

// String はセットの文字列表現を返します
func (s Set) String() string {
	rpeStr := ""
	if s.rpe != nil {
		rpeStr = fmt.Sprintf(" %s", s.rpe.String())
	}
	return fmt.Sprintf("%s × %s (休憩: %s)%s",
		s.weight.String(), s.reps.String(), s.restTime.String(), rpeStr)
}
