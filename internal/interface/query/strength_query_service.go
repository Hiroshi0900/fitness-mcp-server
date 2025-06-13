package query

import (
	"time"

	"fitness-mcp-server/internal/application/query/dto"
	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/domain/strength"
)

// StrengthQueryService は筋トレデータの読み取り専用サービスインターフェース
type StrengthQueryService interface {
	// FindByID はIDで筋トレセッションを検索します
	FindByID(id shared.TrainingID) (*strength.StrengthTraining, error)
	
	// FindByDateRange は指定した期間の筋トレセッションを検索します
	FindByDateRange(start, end time.Time) ([]*strength.StrengthTraining, error)
	
	// FindByDate は指定した日の筋トレセッションを検索します
	FindByDate(date time.Time) ([]*strength.StrengthTraining, error)
	
	// FindAll は全ての筋トレセッションを検索します
	FindAll() ([]*strength.StrengthTraining, error)
	
	// GetPersonalRecords は個人記録を取得します
	GetPersonalRecords(exerciseName *string) ([]dto.PersonalRecordQueryResult, error)
	
	// ExistsById はIDの筋トレセッションが存在するかチェックします
	ExistsById(id shared.TrainingID) (bool, error)
}
