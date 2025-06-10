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

// StrengthTrainingQueryRepository は筋トレデータの読み取り専用操作を担当するインターフェース
type StrengthTrainingQueryRepository interface {
	// GetPersonalRecordsByExercise はエクササイズ別の自己ベストを取得します
	GetPersonalRecordsByExercise(exerciseName strength.ExerciseName) (strength.Weight, error)

	// GetProgressAnalysis は指定したエクササイズの進捗分析を取得します
	GetProgressAnalysis(exerciseName strength.ExerciseName, period time.Duration) (*ProgressAnalysis, error)

	// GetTrainingFrequency は指定期間のトレーニング頻度を取得します
	GetTrainingFrequency(start, end time.Time) (*TrainingFrequency, error)

	// GetVolumeAnalysis は指定期間のボリューム分析を取得します
	GetVolumeAnalysis(start, end time.Time) (*VolumeAnalysis, error)

	// GetRecentTrainings は最近のトレーニングを取得します
	GetRecentTrainings(limit int) ([]*strength.StrengthTraining, error)

	// GetPersonalRecords は個人記録を取得します
	GetPersonalRecords(exerciseName *string) ([]PersonalRecordResult, error)
}

// TODO: 以下の構造体は本来はドメインに持つのが望ましいのではないかと思われる（もしくはリポジトリを呼び出すユースケース）
//  なので、移動したい

// ProgressAnalysis は進捗分析の結果を表す構造体
type ProgressAnalysis struct {
	ExerciseName     strength.ExerciseName
	Period           time.Duration
	StartWeight      strength.Weight
	EndWeight        strength.Weight
	WeightIncrease   float64 // kg
	StartVolume      float64
	EndVolume        float64
	VolumeIncrease   float64 // %
	TotalSessions    int
	ImprovementTrend string // "上昇", "停滞", "下降"
}

// TrainingFrequency はトレーニング頻度の分析結果を表す構造体
type TrainingFrequency struct {
	Period               time.Duration
	TotalSessions        int
	SessionsPerWeek      float64
	SessionsPerMonth     float64
	MostActiveWeekday    string
	AverageSessionLength time.Duration
}

// VolumeAnalysis はボリューム分析の結果を表す構造体
type VolumeAnalysis struct {
	Period                  time.Duration
	TotalVolume             float64
	AverageVolumePerSession float64
	VolumeByExercise        map[string]float64 // エクササイズ名 -> ボリューム
	VolumeByWeek            []WeeklyVolume
	VolumeGrowthRate        float64 // 週次成長率 %
}

// WeeklyVolume は週次ボリュームを表す構造体
type WeeklyVolume struct {
	Week   int // 週番号
	Volume float64
	Date   time.Time // その週の開始日
}

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
