package usecase

import (
	"fmt"
	"log"

	"fitness-mcp-server/internal/application/dto"
	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/interface/repository"
)

type StrengthTrainingUsecaseImpl struct {
	strengthRepo repository.StrengthTrainingRepository
}

func NewStrengthTrainingUsecase(strengthRepo repository.StrengthTrainingRepository) *StrengthTrainingUsecaseImpl {
	return &StrengthTrainingUsecaseImpl{strengthRepo: strengthRepo}
}

func (u *StrengthTrainingUsecaseImpl) RecordTraining(cmd dto.RecordTrainingCommand) (*dto.RecordTrainingResult, error) {
	log.Printf("Recording training session for date: %s", cmd.Date.Format("2006-01-02"))

	training, err := cmd.ToStrengthTraining()
	if err != nil {
		return nil, fmt.Errorf("failed to create training entity: %w", err)
	}

	if err := u.strengthRepo.Save(training); err != nil {
		return nil, fmt.Errorf("failed to save training: %w", err)
	}

	log.Printf("Successfully recorded training with ID: %s", training.ID().String())

	return &dto.RecordTrainingResult{
		TrainingID: training.ID().String(),
		Date:       training.Date(),
		Message:    fmt.Sprintf("筋トレセッション（%d種目、%dセット）を記録しました", training.ExerciseCount(), training.TotalSets()),
	}, nil
}

func (u *StrengthTrainingUsecaseImpl) UpdateTraining(cmd dto.UpdateTrainingCommand) (*dto.UpdateTrainingResult, error) {
	log.Printf("Updating training session with ID: %s", cmd.ID)

	training, err := cmd.ToStrengthTraining()
	if err != nil {
		return nil, fmt.Errorf("failed to create training entity: %w", err)
	}

	if err := u.strengthRepo.Update(training); err != nil {
		return nil, fmt.Errorf("failed to update training: %w", err)
	}

	log.Printf("Successfully updated training with ID: %s", training.ID().String())

	return &dto.UpdateTrainingResult{
		TrainingID: training.ID().String(),
		Date:       training.Date(),
		Message:    "筋トレセッションを更新しました",
	}, nil
}

func (u *StrengthTrainingUsecaseImpl) DeleteTraining(cmd dto.DeleteTrainingCommand) (*dto.DeleteTrainingResult, error) {
	log.Printf("Deleting training session with ID: %s", cmd.ID)

	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// IDを変換してTrainingIDに変換
	trainingID, err := shared.NewTrainingIDFromString(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid training ID: %w", err)
	}

	if err := u.strengthRepo.Delete(trainingID); err != nil {
		return nil, fmt.Errorf("failed to delete training: %w", err)
	}

	log.Printf("Successfully deleted training with ID: %s", cmd.ID)

	return &dto.DeleteTrainingResult{
		TrainingID: cmd.ID,
		Message:    "筋トレセッションを削除しました",
	}, nil
}
