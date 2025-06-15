package strength

import (
	"fmt"
)

// =============================================================================
// エクササイズ定義コンテキスト - エクササイズの実行履歴管理
// =============================================================================

type (
	ExerciseName struct {
		value string
	}

	Exercise struct {
		name ExerciseName // エクササイズ名
		sets []Set        // セットのリスト
	}
)

// 定義済みエクササイズ名の定数
var (
	BenchPress = ExerciseName{value: "ベンチプレス"}
	Squat      = ExerciseName{value: "スクワット"}
	Deadlift   = ExerciseName{value: "デッドリフト"}
)

// NewExerciseName はエクササイズ名を作成します
func NewExerciseName(name string) (ExerciseName, error) {
	if name == "" {
		return ExerciseName{}, fmt.Errorf("exercise name cannot be empty")
	}
	return ExerciseName{value: name}, nil
}

// Name はエクササイズ名を返します
func (en ExerciseName) Name() string {
	return en.value
}

// String はエクササイズ名の文字列表現を返します
func (en ExerciseName) String() string {
	return en.value
}

// Equals は2つのエクササイズ名が等しいかを判定します
func (en ExerciseName) Equals(other ExerciseName) bool {
	return en.value == other.value
}

// IsBIG3 はBIG3エクササイズかどうかを判定します
func (en ExerciseName) IsBIG3() bool {
	return en.Equals(BenchPress) || en.Equals(Squat) || en.Equals(Deadlift)
}

// NewExercise は新しいExerciseを作成します
func NewExercise(name ExerciseName) *Exercise {
	return &Exercise{
		name: name,
		sets: make([]Set, 0),
	}
}

// Name はエクササイズ名を返します
func (e *Exercise) Name() ExerciseName {
	return e.name
}

// Sets は全セットを返します
func (e *Exercise) Sets() []Set {
	// コピーを返して不変性を保つ
	result := make([]Set, len(e.sets))
	copy(result, e.sets)
	return result
}

// AddSet はセットを追加します
func (e *Exercise) AddSet(set Set) {
	e.sets = append(e.sets, set)
}

// SetCount はセット数を返します
func (e *Exercise) SetCount() int {
	return len(e.sets)
}

// MaxWeight は最大重量を返します
func (e *Exercise) MaxWeight() (Weight, error) {
	if len(e.sets) == 0 {
		return Weight{}, fmt.Errorf("no sets recorded")
	}

	maxWeight := e.sets[0].weight
	for _, set := range e.sets[1:] {
		if set.weight.Kg() > maxWeight.Kg() {
			maxWeight = set.weight
		}
	}
	return maxWeight, nil
}

// TotalVolume は総ボリューム（重量×回数の合計）を計算します
func (e *Exercise) TotalVolume() float64 {
	volume := 0.0
	for _, set := range e.sets {
		volume += set.weight.Kg() * float64(set.reps.Count())
	}
	return volume
}

// String はエクササイズの文字列表現を返します
func (e *Exercise) String() string {
	return fmt.Sprintf("%s - %dセット",
		e.name.String(), len(e.sets))
}
