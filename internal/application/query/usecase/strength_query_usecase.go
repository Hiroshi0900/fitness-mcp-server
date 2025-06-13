package usecase

import (
	"fmt"
	"time"

	"fitness-mcp-server/internal/application/query/dto"
	"fitness-mcp-server/internal/interface/query"
)

// StrengthQueryUsecase は筋トレデータの読み取り系ユースケースインターフェース
type StrengthQueryUsecase interface {
	GetTrainingsByDateRange(query dto.GetTrainingsByDateRangeQuery) (*dto.GetTrainingsByDateRangeResponse, error)
}

// strengthQueryUsecaseImpl はStrengthQueryUsecaseの実装
type strengthQueryUsecaseImpl struct {
	queryService query.StrengthQueryService
}

// NewStrengthQueryUsecase は新しいStrengthQueryUsecaseを作成します
func NewStrengthQueryUsecase(queryService query.StrengthQueryService) StrengthQueryUsecase {
	return &strengthQueryUsecaseImpl{
		queryService: queryService,
	}
}

// GetTrainingsByDateRange は指定した期間のトレーニングセッションを取得します
func (u *strengthQueryUsecaseImpl) GetTrainingsByDateRange(query dto.GetTrainingsByDateRangeQuery) (*dto.GetTrainingsByDateRangeResponse, error) {
	// 入力値の検証
	if query.StartDate.After(query.EndDate) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}
	
	// 期間制限チェック（最大1年間）
	maxPeriod := 365 * 24 * time.Hour // 1年
	if query.EndDate.Sub(query.StartDate) > maxPeriod {
		return nil, fmt.Errorf("period too long: maximum 1 year allowed")
	}

	// クエリサービスからデータを取得
	trainings, err := u.queryService.FindByDateRange(query.StartDate, query.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainings by date range: %w", err)
	}

	// ドメインエンティティをDTOに変換
	trainingDTOs := make([]*dto.TrainingDTO, 0, len(trainings))
	for _, training := range trainings {
		trainingDTOs = append(trainingDTOs, dto.TrainingToDTO(training))
	}

	// 期間の文字列表現を作成
	period := fmt.Sprintf("%s to %s",
		query.StartDate.Format("2006-01-02"),
		query.EndDate.Format("2006-01-02"))

	return &dto.GetTrainingsByDateRangeResponse{
		Trainings: trainingDTOs,
		Count:     len(trainingDTOs),
		Period:    period,
	}, nil
}
