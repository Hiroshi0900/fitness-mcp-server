package usecase

import (
	"fitness-mcp-server/internal/application/command/dto"
)

// StrengthTrainingUsecase は筋トレ記録のユースケースインターフェース
type StrengthTrainingUsecase interface {
	RecordTraining(cmd dto.RecordTrainingCommand) (*dto.RecordTrainingResult, error)
	UpdateTraining(cmd dto.UpdateTrainingCommand) (*dto.UpdateTrainingResult, error)
	DeleteTraining(cmd dto.DeleteTrainingCommand) (*dto.DeleteTrainingResult, error)
}
