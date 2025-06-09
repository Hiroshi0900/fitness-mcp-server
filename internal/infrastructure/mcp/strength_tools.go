package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"fitness-mcp-server/internal/application/command"
	"fitness-mcp-server/internal/application/dto"
)

// =============================================================================
// 筋トレ用MCPツール定義
// =============================================================================

// RegisterStrengthTools は筋トレ関連のMCPツールを登録します
func RegisterStrengthTools(server *MCPServer, strengthHandler *command.StrengthCommandHandler) {
	// record_training ツール
	server.RegisterTool(Tool{
		Name:        "record_training",
		Description: "筋トレセッションを記録する",
		InputSchema: getRecordTrainingSchema(),
		Handler:     createRecordTrainingHandler(strengthHandler),
	})

	// update_training ツール
	server.RegisterTool(Tool{
		Name:        "update_training",
		Description: "筋トレセッションを更新する",
		InputSchema: getUpdateTrainingSchema(),
		Handler:     createUpdateTrainingHandler(strengthHandler),
	})

	// delete_training ツール
	server.RegisterTool(Tool{
		Name:        "delete_training",
		Description: "筋トレセッションを削除する",
		InputSchema: getDeleteTrainingSchema(),
		Handler:     createDeleteTrainingHandler(strengthHandler),
	})
}

// getRecordTrainingSchema はrecord_trainingツールのスキーマを返します
func getRecordTrainingSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"date": map[string]interface{}{
				"type":        "string",
				"format":      "date",
				"description": "トレーニング日（YYYY-MM-DD形式）",
			},
			"exercises": map[string]interface{}{
				"type":        "array",
				"description": "エクササイズのリスト",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "エクササイズ名（例: ベンチプレス、スクワット、デッドリフト）",
						},
						"category": map[string]interface{}{
							"type":        "string",
							"description": "エクササイズカテゴリ（Compound, Isolation, Cardio）",
						},
						"sets": map[string]interface{}{
							"type":        "array",
							"description": "セットのリスト",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"weight_kg": map[string]interface{}{
										"type":        "number",
										"description": "重量（kg）",
										"minimum":     0,
									},
									"reps": map[string]interface{}{
										"type":        "integer",
										"description": "レップ数",
										"minimum":     1,
									},
									"rest_time_seconds": map[string]interface{}{
										"type":        "integer",
										"description": "休憩時間（秒）",
										"minimum":     0,
									},
									"rpe": map[string]interface{}{
										"type":        "integer",
										"minimum":     1,
										"maximum":     10,
										"description": "RPE（主観的運動強度、1-10）オプション",
									},
								},
								"required": []string{"weight_kg", "reps", "rest_time_seconds"},
							},
						},
					},
					"required": []string{"name", "category", "sets"},
				},
			},
			"notes": map[string]interface{}{
				"type":        "string",
				"description": "メモ（オプション）",
			},
		},
		"required": []string{"date", "exercises"},
	}
}

// getUpdateTrainingSchema はupdate_trainingツールのスキーマを返します
func getUpdateTrainingSchema() map[string]interface{} {
	schema := getRecordTrainingSchema()
	// IDフィールドを追加
	schema["properties"].(map[string]interface{})["id"] = map[string]interface{}{
		"type":        "string",
		"description": "更新するトレーニングセッションのID",
	}
	schema["required"] = []string{"id", "date", "exercises"}
	return schema
}

// getDeleteTrainingSchema はdelete_trainingツールのスキーマを返します
func getDeleteTrainingSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "削除するトレーニングセッションのID",
			},
		},
		"required": []string{"id"},
	}
}

// createRecordTrainingHandler はrecord_trainingツールのハンドラーを作成します
func createRecordTrainingHandler(strengthHandler *command.StrengthCommandHandler) ToolHandler {
	return func(arguments map[string]interface{}) (interface{}, error) {
		// 引数をDTOに変換
		argumentsBytes, err := json.Marshal(arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal arguments: %w", err)
		}

		var cmd dto.RecordTrainingCommand
		if err := json.Unmarshal(argumentsBytes, &cmd); err != nil {
			return nil, fmt.Errorf("invalid arguments: %w", err)
		}

		// ユースケースを実行
		result, err := strengthHandler.RecordTraining(cmd)
		if err != nil {
			return nil, fmt.Errorf("training recording failed: %w", err)
		}

		// 成功メッセージを作成
		response := fmt.Sprintf("✅ %s\n\n🏋️ トレーニングID: %s\n📅 日付: %s",
			result.Message,
			result.TrainingID,
			result.Date.Format("2006-01-02"))

		return response, nil
	}
}

// createUpdateTrainingHandler はupdate_trainingツールのハンドラーを作成します
func createUpdateTrainingHandler(strengthHandler *command.StrengthCommandHandler) ToolHandler {
	return func(arguments map[string]interface{}) (interface{}, error) {
		// 引数をDTOに変換
		argumentsBytes, err := json.Marshal(arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal arguments: %w", err)
		}

		var cmd dto.UpdateTrainingCommand
		if err := json.Unmarshal(argumentsBytes, &cmd); err != nil {
			return nil, fmt.Errorf("invalid arguments: %w", err)
		}

		// ユースケースを実行
		result, err := strengthHandler.UpdateTraining(cmd)
		if err != nil {
			return nil, fmt.Errorf("training update failed: %w", err)
		}

		// 成功メッセージを作成
		response := fmt.Sprintf("✅ %s\n\n🏋️ トレーニングID: %s\n📅 日付: %s",
			result.Message,
			result.TrainingID,
			result.Date.Format("2006-01-02"))

		return response, nil
	}
}

// createDeleteTrainingHandler はdelete_trainingツールのハンドラーを作成します
func createDeleteTrainingHandler(strengthHandler *command.StrengthCommandHandler) ToolHandler {
	return func(arguments map[string]interface{}) (interface{}, error) {
		// 引数をDTOに変換
		argumentsBytes, err := json.Marshal(arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal arguments: %w", err)
		}

		var cmd dto.DeleteTrainingCommand
		if err := json.Unmarshal(argumentsBytes, &cmd); err != nil {
			return nil, fmt.Errorf("invalid arguments: %w", err)
		}

		// ユースケースを実行
		result, err := strengthHandler.DeleteTraining(cmd)
		if err != nil {
			return nil, fmt.Errorf("training deletion failed: %w", err)
		}

		// 成功メッセージを作成
		response := fmt.Sprintf("✅ %s\n\n🏋️ トレーニングID: %s",
			result.Message,
			result.TrainingID)

		return response, nil
	}
}

// 便利関数: BIG3セッション記録用のヘルパー
func CreateBIG3Session(date time.Time, benchPress, squat, deadlift []dto.SetDTO, notes string) dto.RecordTrainingCommand {
	exercises := []dto.ExerciseDTO{}

	if len(benchPress) > 0 {
		exercises = append(exercises, dto.ExerciseDTO{
			Name:     "ベンチプレス",
			Category: "Compound",
			Sets:     benchPress,
		})
	}

	if len(squat) > 0 {
		exercises = append(exercises, dto.ExerciseDTO{
			Name:     "スクワット",
			Category: "Compound",
			Sets:     squat,
		})
	}

	if len(deadlift) > 0 {
		exercises = append(exercises, dto.ExerciseDTO{
			Name:     "デッドリフト",
			Category: "Compound",
			Sets:     deadlift,
		})
	}

	return dto.RecordTrainingCommand{
		Date:      date,
		Exercises: exercises,
		Notes:     notes,
	}
}
