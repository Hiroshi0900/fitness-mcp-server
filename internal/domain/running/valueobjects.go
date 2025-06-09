package running

import (
	"fmt"
	"time"
)

// =============================================================================
// ランニング測定値コンテキスト - 距離、時間、ペース、心拍数
// =============================================================================

// Distance は距離を表す値オブジェクト
type Distance struct {
	value float64 // km単位
}

// NewDistance は距離を作成します
func NewDistance(km float64) (Distance, error) {
	if km <= 0 {
		return Distance{}, fmt.Errorf("distance must be positive: %f", km)
	}
	if km > 1000 { // 現実的な上限設定
		return Distance{}, fmt.Errorf("distance is too large: %f", km)
	}
	return Distance{value: km}, nil
}

// Km は距離をkm単位で返します
func (d Distance) Km() float64 {
	return d.value
}

// Meters は距離をメートル単位で返します
func (d Distance) Meters() float64 {
	return d.value * 1000
}

// String は距離の文字列表現を返します
func (d Distance) String() string {
	return fmt.Sprintf("%.2fkm", d.value)
}

// Equals は2つの距離が等しいかを判定します
func (d Distance) Equals(other Distance) bool {
	return d.value == other.value
}

// Duration は時間を表す値オブジェクト
type Duration struct {
	duration time.Duration
}

// NewDuration は時間を作成します
func NewDuration(d time.Duration) (Duration, error) {
	if d <= 0 {
		return Duration{}, fmt.Errorf("duration must be positive: %v", d)
	}
	if d > 24*time.Hour { // 現実的な上限設定
		return Duration{}, fmt.Errorf("duration is too long: %v", d)
	}
	return Duration{duration: d}, nil
}

// NewDurationFromMinutes は分数から時間を作成します
func NewDurationFromMinutes(minutes float64) (Duration, error) {
	if minutes <= 0 {
		return Duration{}, fmt.Errorf("duration must be positive: %f minutes", minutes)
	}
	d := time.Duration(minutes * float64(time.Minute))
	return NewDuration(d)
}

// Value は時間を返します
func (d Duration) Value() time.Duration {
	return d.duration
}

// Minutes は時間を分単位で返します
func (d Duration) Minutes() float64 {
	return d.duration.Minutes()
}

// Seconds は時間を秒単位で返します
func (d Duration) Seconds() float64 {
	return d.duration.Seconds()
}

// String は時間の文字列表現を返します（MM:SS形式）
func (d Duration) String() string {
	totalSeconds := int(d.duration.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// Equals は2つの時間が等しいかを判定します
func (d Duration) Equals(other Duration) bool {
	return d.duration == other.duration
}

// Pace はペースを表す値オブジェクト（分/km）
type Pace struct {
	minutesPerKm float64
}

// NewPace はペースを作成します
func NewPace(minutesPerKm float64) (Pace, error) {
	if minutesPerKm <= 0 {
		return Pace{}, fmt.Errorf("pace must be positive: %f", minutesPerKm)
	}
	if minutesPerKm > 30 { // 現実的な上限設定（30分/km）
		return Pace{}, fmt.Errorf("pace is too slow: %f", minutesPerKm)
	}
	return Pace{minutesPerKm: minutesPerKm}, nil
}

// CalculatePace は距離と時間からペースを計算します
func CalculatePace(distance Distance, duration Duration) (Pace, error) {
	if distance.Km() <= 0 {
		return Pace{}, fmt.Errorf("distance must be positive")
	}
	minutesPerKm := duration.Minutes() / distance.Km()
	return NewPace(minutesPerKm)
}

// MinutesPerKm はペースを分/km単位で返します
func (p Pace) MinutesPerKm() float64 {
	return p.minutesPerKm
}

// SecondsPerKm はペースを秒/km単位で返します
func (p Pace) SecondsPerKm() float64 {
	return p.minutesPerKm * 60
}

// String はペースの文字列表現を返します（M:SS/km形式）
func (p Pace) String() string {
	totalSeconds := int(p.SecondsPerKm())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d/km", minutes, seconds)
}

// Equals は2つのペースが等しいかを判定します
func (p Pace) Equals(other Pace) bool {
	return p.minutesPerKm == other.minutesPerKm
}

// IsFasterThan は他のペースより速いかを判定します
func (p Pace) IsFasterThan(other Pace) bool {
	return p.minutesPerKm < other.minutesPerKm
}

// HeartRate は心拍数を表す値オブジェクト
type HeartRate struct {
	value int // bpm
}

// NewHeartRate は心拍数を作成します
func NewHeartRate(bpm int) (HeartRate, error) {
	if bpm <= 0 {
		return HeartRate{}, fmt.Errorf("heart rate must be positive: %d", bpm)
	}
	if bpm > 250 { // 現実的な上限設定
		return HeartRate{}, fmt.Errorf("heart rate is too high: %d", bpm)
	}
	return HeartRate{value: bpm}, nil
}

// BPM は心拍数をbpm単位で返します
func (hr HeartRate) BPM() int {
	return hr.value
}

// String は心拍数の文字列表現を返します
func (hr HeartRate) String() string {
	return fmt.Sprintf("%dbpm", hr.value)
}

// Equals は2つの心拍数が等しいかを判定します
func (hr HeartRate) Equals(other HeartRate) bool {
	return hr.value == other.value
}

// CalculateTimeForDistance は指定したペースで指定した距離を走る時間を計算します
func CalculateTimeForDistance(pace Pace, distance Distance) Duration {
	totalMinutes := pace.MinutesPerKm() * distance.Km()
	duration := time.Duration(totalMinutes * float64(time.Minute))
	return Duration{duration: duration}
}

// CalculateDistanceForTime は指定したペースで指定した時間走る距離を計算します
func CalculateDistanceForTime(pace Pace, duration Duration) (Distance, error) {
	km := duration.Minutes() / pace.MinutesPerKm()
	return NewDistance(km)
}
