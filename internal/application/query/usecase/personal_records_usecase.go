package usecase

import (
	"fmt"

	query_dto "fitness-mcp-server/internal/application/query/dto"
	"fitness-mcp-server/internal/interface/query"
)

// personalRecordsUsecaseImpl は個人記録に関するクエリユースケース
type (
	PersonalRecordsUsecase interface {
		GetPersonalRecords(query query_dto.GetPersonalRecordsQuery) (*query_dto.GetPersonalRecordsResponse, error)
	}
	personalRecordsUsecaseImpl struct {
		queryService query.StrengthQueryService
	}
)

// NewPersonalRecordsUsecase は新しいPersonalRecordsUsecaseを作成します
func NewPersonalRecordsUsecase(queryService query.StrengthQueryService) PersonalRecordsUsecase {
	return &personalRecordsUsecaseImpl{
		queryService: queryService,
	}
}

// GetPersonalRecords は個人記録を取得します
func (u *personalRecordsUsecaseImpl) GetPersonalRecords(query query_dto.GetPersonalRecordsQuery) (*query_dto.GetPersonalRecordsResponse, error) {
	// クエリサービスから生データを取得
	queryResults, err := u.queryService.GetPersonalRecords(query.ExerciseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get personal records: %w", err)
	}

	// Query結果をDTOに変換
	records := make([]query_dto.PersonalRecord, len(queryResults))
	for i, queryResult := range queryResults {
		records[i] = convertQueryResultToDTO(queryResult)
	}

	response := &query_dto.GetPersonalRecordsResponse{
		Records: records,
		Count:   len(records),
	}

	return response, nil
}

// convertQueryResultToDTO はQuery結果をDTOに変換します
func convertQueryResultToDTO(queryResult query_dto.PersonalRecordQueryResult) query_dto.PersonalRecord {
	return query_dto.PersonalRecord{
		ExerciseName: queryResult.ExerciseName,
		Category:     queryResult.Category,
		MaxWeight: query_dto.PersonalRecordDetail{
			Value:      queryResult.MaxWeight.Value,
			Date:       queryResult.MaxWeight.Date,
			TrainingID: queryResult.MaxWeight.TrainingID,
			SetDetails: convertSetQueryDetailsToDTO(queryResult.MaxWeight.SetDetails),
		},
		MaxReps: query_dto.PersonalRecordDetail{
			Value:      queryResult.MaxReps.Value,
			Date:       queryResult.MaxReps.Date,
			TrainingID: queryResult.MaxReps.TrainingID,
			SetDetails: convertSetQueryDetailsToDTO(queryResult.MaxReps.SetDetails),
		},
		MaxVolume: query_dto.PersonalRecordDetail{
			Value:      queryResult.MaxVolume.Value,
			Date:       queryResult.MaxVolume.Date,
			TrainingID: queryResult.MaxVolume.TrainingID,
			SetDetails: convertSetQueryDetailsToDTO(queryResult.MaxVolume.SetDetails),
		},
		TotalSessions: queryResult.TotalSessions,
		LastPerformed: queryResult.LastPerformed,
	}
}

// convertSetQueryDetailsToDTO はQuery SetDetailsをDTO SetInfoに変換します
func convertSetQueryDetailsToDTO(setDetails *query_dto.SetQueryDetails) *query_dto.SetInfo {
	if setDetails == nil {
		return nil
	}

	return &query_dto.SetInfo{
		WeightKg:        setDetails.WeightKg,
		Reps:            setDetails.Reps,
		RestTimeSeconds: setDetails.RestTimeSeconds,
		RPE:             setDetails.RPE,
	}
}
