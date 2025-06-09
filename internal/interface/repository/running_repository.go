package repository

import (
	"time"

	"fitness-mcp-server/internal/domain/running"
	"fitness-mcp-server/internal/domain/shared"
)

// RunningSessionRepository はランニングセッションデータの永続化を担当するインターフェース
type RunningSessionRepository interface {
	// Save はランニングセッションを保存します
	Save(session *running.RunningSession) error

	// FindByID はIDでランニングセッションを検索します
	FindByID(id shared.SessionID) (*running.RunningSession, error)

	// FindByDateRange は指定した期間のランニングセッションを検索します
	FindByDateRange(start, end time.Time) ([]*running.RunningSession, error)

	// FindByDate は指定した日のランニングセッションを検索します
	FindByDate(date time.Time) ([]*running.RunningSession, error)

	// FindAll は全てのランニングセッションを検索します（テスト用）
	FindAll() ([]*running.RunningSession, error)

	// Update は既存のランニングセッションを更新します
	Update(session *running.RunningSession) error

	// Delete はランニングセッションを削除します
	Delete(id shared.SessionID) error

	// ExistsById はIDのランニングセッションが存在するかチェックします
	ExistsById(id shared.SessionID) (bool, error)
}

// RunningGoalRepository はランニング目標データの永続化を担当するインターフェース
type RunningGoalRepository interface {
	// Save はランニング目標を保存します
	Save(goal *running.RunningGoal) error

	// FindByID はIDでランニング目標を検索します
	FindByID(id shared.GoalID) (*running.RunningGoal, error)

	// FindByStatus は状態でランニング目標を検索します
	FindByStatus(status running.GoalStatus) ([]*running.RunningGoal, error)

	// FindActiveGoals はアクティブな目標を検索します
	FindActiveGoals() ([]*running.RunningGoal, error)

	// FindAll は全てのランニング目標を検索します
	FindAll() ([]*running.RunningGoal, error)

	// Update は既存のランニング目標を更新します
	Update(goal *running.RunningGoal) error

	// Delete はランニング目標を削除します
	Delete(id shared.GoalID) error
}

// RunningAnalyticsRepository はランニングデータの分析専用操作を担当するインターフェース
type RunningAnalyticsRepository interface {
	// GetPersonalBests は距離別の自己ベストを取得します
	GetPersonalBests() (map[string]*PersonalBest, error)

	// GetPersonalBestByDistance は指定距離の自己ベストを取得します
	GetPersonalBestByDistance(distance running.Distance) (*PersonalBest, error)

	// GetPaceProgress は指定期間のペース改善傾向を取得します
	GetPaceProgress(distance running.Distance, period time.Duration) (*PaceProgress, error)

	// GetRunningFrequency は指定期間のランニング頻度を取得します
	GetRunningFrequency(start, end time.Time) (*RunningFrequency, error)

	// GetMonthlyStats は月次統計を取得します
	GetMonthlyStats(year int, month int) (*MonthlyRunningStats, error)

	// GetGoalProgress は目標進捗を取得します
	GetGoalProgress(goalID shared.GoalID) (*GoalProgress, error)

	// GetRecentSessions は最近のセッションを取得します
	GetRecentSessions(limit int) ([]*running.RunningSession, error)

	// GetDistanceDistribution は距離別の走行回数分布を取得します
	GetDistanceDistribution(start, end time.Time) (*DistanceDistribution, error)
}

// PersonalBest は自己ベスト記録を表す構造体
type PersonalBest struct {
	Distance     running.Distance  `json:"distance"`
	BestTime     running.Duration  `json:"best_time"`
	BestPace     running.Pace      `json:"best_pace"`
	SessionID    shared.SessionID  `json:"session_id"`
	AchievedDate time.Time         `json:"achieved_date"`
}

// PaceProgress はペース改善傾向を表す構造体
type PaceProgress struct {
	Distance           running.Distance  `json:"distance"`
	Period             time.Duration     `json:"period"`
	StartPace          running.Pace      `json:"start_pace"`
	EndPace            running.Pace      `json:"end_pace"`
	PaceImprovement    float64           `json:"pace_improvement"` // 秒/kmの改善量
	TotalSessions      int               `json:"total_sessions"`
	ImprovementTrend   string            `json:"improvement_trend"` // "改善", "停滞", "悪化"
	WeeklyPaces        []WeeklyPace      `json:"weekly_paces"`
}

// WeeklyPace は週次ペースを表す構造体
type WeeklyPace struct {
	Week         int              `json:"week"`
	AveragePace  running.Pace     `json:"average_pace"`
	BestPace     running.Pace     `json:"best_pace"`
	SessionCount int              `json:"session_count"`
	Date         time.Time        `json:"date"` // その週の開始日
}

// RunningFrequency はランニング頻度の分析結果を表す構造体
type RunningFrequency struct {
	Period               time.Duration         `json:"period"`
	TotalSessions        int                   `json:"total_sessions"`
	SessionsPerWeek      float64               `json:"sessions_per_week"`
	SessionsPerMonth     float64               `json:"sessions_per_month"`
	TotalDistance        running.Distance      `json:"total_distance"`
	AverageDistance      running.Distance      `json:"average_distance"`
	TotalTime            running.Duration      `json:"total_time"`
	AverageTime          running.Duration      `json:"average_time"`
	MostActiveWeekday    string                `json:"most_active_weekday"`
	RunTypeDistribution  map[string]int        `json:"run_type_distribution"` // RunType -> 回数
}

// MonthlyRunningStats は月次統計を表す構造体
type MonthlyRunningStats struct {
	Year             int                   `json:"year"`
	Month            int                   `json:"month"`
	TotalSessions    int                   `json:"total_sessions"`
	TotalDistance    running.Distance      `json:"total_distance"`
	TotalTime        running.Duration      `json:"total_time"`
	AverageDistance  running.Distance      `json:"average_distance"`
	AveragePace      running.Pace          `json:"average_pace"`
	BestPace         running.Pace          `json:"best_pace"`
	LongestRun       running.Distance      `json:"longest_run"`
	FastestKm        running.Pace          `json:"fastest_km"` // その月の最速1kmペース
	WeeklyBreakdown  []WeeklyRunningStats  `json:"weekly_breakdown"`
}

// WeeklyRunningStats は週次統計を表す構造体
type WeeklyRunningStats struct {
	WeekNumber      int                `json:"week_number"`
	StartDate       time.Time          `json:"start_date"`
	Sessions        int                `json:"sessions"`
	Distance        running.Distance   `json:"distance"`
	Time            running.Duration   `json:"time"`
	AveragePace     running.Pace       `json:"average_pace"`
}

// GoalProgress は目標進捗を表す構造体
type GoalProgress struct {
	Goal              *running.RunningGoal  `json:"goal"`
	CurrentBestTime   *running.Duration     `json:"current_best_time"`   // その距離での現在の自己ベスト
	CurrentBestPace   *running.Pace         `json:"current_best_pace"`   // その距離での現在の自己ベストペース
	TimeToTarget      *running.Duration     `json:"time_to_target"`      // 目標タイムまでの差
	PaceToTarget      *running.Pace         `json:"pace_to_target"`      // 目標ペースまでの差
	RecentSessions    []*running.RunningSession `json:"recent_sessions"` // その距離での最近のセッション
	ProgressRate      float64               `json:"progress_rate"`       // 進捗率（0-1）
	EstimatedAchievementDate *time.Time     `json:"estimated_achievement_date"` // 推定達成日
	Recommendation    string                `json:"recommendation"`      // 推奨事項
}

// DistanceDistribution は距離別の走行回数分布を表す構造体
type DistanceDistribution struct {
	Period          time.Duration            `json:"period"`
	DistanceRanges  map[string]int           `json:"distance_ranges"` // "0-5km": 10回 など
	MostPopularDistance running.Distance     `json:"most_popular_distance"`
	AverageDistance running.Distance         `json:"average_distance"`
	DistanceVariation float64                `json:"distance_variation"` // 距離のばらつき（標準偏差）
}
