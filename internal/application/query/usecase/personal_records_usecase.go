package usecase

import (
	"fmt"

	query_dto "fitness-mcp-server/internal/application/query/dto"
	"fitness-mcp-server/internal/interface/repository"
)

// personalRecordsUsecaseImpl は個人記録に関するクエリユースケース
type (
	PersonalRecordsUsecase interface {
		GetPersonalRecords(query query_dto.GetPersonalRecordsQuery) (*query_dto.GetPersonalRecordsResponse, error)
	}
	personalRecordsUsecaseImpl struct {
		repo repository.StrengthTrainingRepository
	}
)

// NewPersonalRecordsUsecase は新しいPersonalRecordsUsecaseを作成します
func NewPersonalRecordsUsecase(repo repository.StrengthTrainingRepository) PersonalRecordsUsecase {
	return &personalRecordsUsecaseImpl{
		repo: repo,
	}
}

// GetPersonalRecords は個人記録を取得します
func (u *personalRecordsUsecaseImpl) GetPersonalRecords(query query_dto.GetPersonalRecordsQuery) (*query_dto.GetPersonalRecordsResponse, error) {
	// リポジトリから生データを取得
	repoResults, err := u.repo.GetPersonalRecords(query.ExerciseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get personal records: %w", err)
	}

	// Repository結果をDTOに変換
	records := make([]query_dto.PersonalRecord, len(repoResults))
	for i, repoResult := range repoResults {
		records[i] = convertToDTO(repoResult)
	}

	response := &query_dto.GetPersonalRecordsResponse{
		Records: records,
		Count:   len(records),
	}

	return response, nil
}

// convertToDTO はRepository結果をDTOに変換します
func convertToDTO(repoResult repository.PersonalRecordResult) query_dto.PersonalRecord {
	return query_dto.PersonalRecord{
		ExerciseName: repoResult.ExerciseName,
		Category:     repoResult.Category,
		MaxWeight: query_dto.PersonalRecordDetail{
			Value:      repoResult.MaxWeight.Value,
			Date:       repoResult.MaxWeight.Date,
			TrainingID: repoResult.MaxWeight.TrainingID,
			SetDetails: convertSetDetailsToDTO(repoResult.MaxWeight.SetDetails),
		},
		MaxReps: query_dto.PersonalRecordDetail{
			Value:      repoResult.MaxReps.Value,
			Date:       repoResult.MaxReps.Date,
			TrainingID: repoResult.MaxReps.TrainingID,
			SetDetails: convertSetDetailsToDTO(repoResult.MaxReps.SetDetails),
		},
		MaxVolume: query_dto.PersonalRecordDetail{
			Value:      repoResult.MaxVolume.Value,
			Date:       repoResult.MaxVolume.Date,
			TrainingID: repoResult.MaxVolume.TrainingID,
			SetDetails: convertSetDetailsToDTO(repoResult.MaxVolume.SetDetails),
		},
		TotalSessions: repoResult.TotalSessions,
		LastPerformed: repoResult.LastPerformed,
	}
}

// convertSetDetailsToDTO はRepository SetDetailsをDTO SetInfoに変換します
func convertSetDetailsToDTO(setDetails *repository.SetDetails) *query_dto.SetInfo {
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
