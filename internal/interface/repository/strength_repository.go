package repository

import (
	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/domain/strength"
)

// StrengthTrainingRepository は筋トレデータの永続化を担当するインターフェース（書き込み専用）
type StrengthTrainingRepository interface {
	// Initialize はデータベースの初期化を行います
	Initialize() error

	// Save は筋トレセッションを保存します
	Save(training *strength.StrengthTraining) error

	// Update は既存の筋トレセッションを更新します
	Update(training *strength.StrengthTraining) error

	// Delete は筋トレセッションを削除します
	Delete(id shared.TrainingID) error
}
