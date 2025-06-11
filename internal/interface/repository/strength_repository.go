package repository

import (
	"time"

	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/domain/strength"
)

// StrengthTrainingRepository は筋トレデータの永続化を担当するインターフェース
type StrengthTrainingRepository interface {
	// Initialize はデータベースの初期化を行います
	Initialize() error

	// Save は筋トレセッションを保存します
	Save(training *strength.StrengthTraining) error

	// FindByID はIDで筋トレセッションを検索します
	FindByID(id shared.TrainingID) (*strength.StrengthTraining, error)

	// FindByDateRange は指定した期間の筋トレセッションを検索します
	FindByDateRange(start, end time.Time) ([]*strength.StrengthTraining, error)

	// FindByDate は指定した日の筋トレセッションを検索します
	FindByDate(date time.Time) ([]*strength.StrengthTraining, error)

	// FindAll は全ての筋トレセッションを検索します（テスト用）
	FindAll() ([]*strength.StrengthTraining, error)

	// Update は既存の筋トレセッションを更新します
	Update(training *strength.StrengthTraining) error

	// Delete は筋トレセッションを削除します
	Delete(id shared.TrainingID) error

	// ExistsById はIDの筋トレセッションが存在するかチェックします
	ExistsById(id shared.TrainingID) (bool, error)

	// GetPersonalRecords は個人記録を取得します
	GetPersonalRecords(exerciseName *string) ([]PersonalRecordResult, error)
}


// TODO: 以下の構造体は本来はドメインに持つのが望ましいのではないかと思われる（もしくはリポジトリを呼び出すユースケース）
//  なので、移動したい


// TODO: リポジトリ専用を置く必要があるのかは少し微妙な気がするので後で検討する
type PersonalRecordResult struct {
	ExerciseName  string
	Category      string
	MaxWeight     PersonalRecordDetail
	MaxReps       PersonalRecordDetail
	MaxVolume     PersonalRecordDetail
	TotalSessions int
	LastPerformed time.Time
}

// PersonalRecordDetail は記録詳細（Repository層用）
type PersonalRecordDetail struct {
	Value      float64
	Date       time.Time
	TrainingID string
	SetDetails *SetDetails
}

// SetDetails はセット詳細（Repository層用）
type SetDetails struct {
	WeightKg        float64
	Reps            int
	RestTimeSeconds int
	RPE             *int
}
