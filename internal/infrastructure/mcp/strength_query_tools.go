package mcp

import (
	"fmt"
	"time"

	"fitness-mcp-server/internal/application/query/dto"
	"fitness-mcp-server/internal/application/query/handler"
)

// =============================================================================
// 筋トレ用MCPクエリツール定義
// =============================================================================

// RegisterStrengthQueryTools は筋トレ関連のクエリMCPツールを登録します
func RegisterStrengthQueryTools(server *MCPServer, queryHandler *handler.StrengthQueryHandler) {
	// get_trainings_by_date_range ツール
	server.RegisterTool(Tool{
		Name:        "get_trainings_by_date_range",
		Description: "指定した期間のトレーニングセッションを取得する",
		InputSchema: getTrainingsByDateRangeSchema(),
		Handler:     createGetTrainingsByDateRangeHandler(queryHandler),
	})

	// get_personal_records ツール
	server.RegisterTool(Tool{
		Name:        "get_personal_records",
		Description: "個人記録（最大重量、最大レップ数、最大ボリューム等）を取得する",
		InputSchema: getPersonalRecordsSchema(),
		Handler:     createGetPersonalRecordsHandler(queryHandler),
	})
}

// getTrainingsByDateRangeSchema はget_trainings_by_date_rangeツールのスキーマを返します
func getTrainingsByDateRangeSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"start_date": map[string]interface{}{
				"type":        "string",
				"format":      "date",
				"description": "検索開始日（YYYY-MM-DD形式）",
			},
			"end_date": map[string]interface{}{
				"type":        "string",
				"format":      "date",
				"description": "検索終了日（YYYY-MM-DD形式）",
			},
		},
		"required": []string{"start_date", "end_date"},
	}
}

// createGetTrainingsByDateRangeHandler はget_trainings_by_date_rangeツールのハンドラーを作成します
func createGetTrainingsByDateRangeHandler(queryHandler *handler.StrengthQueryHandler) ToolHandler {
	return func(arguments map[string]interface{}) (interface{}, error) {
		// 引数の取得と検証
		startDateStr, ok := arguments["start_date"].(string)
		if !ok {
			return nil, fmt.Errorf("start_date is required and must be a string")
		}

		endDateStr, ok := arguments["end_date"].(string)
		if !ok {
			return nil, fmt.Errorf("end_date is required and must be a string")
		}

		// 日付のパース
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format: %w", err)
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}

		// クエリの作成
		query := dto.GetTrainingsByDateRangeQuery{
			StartDate: startDate,
			EndDate:   endDate,
		}

		// ハンドラーの実行
		response, err := queryHandler.GetTrainingsByDateRange(query)
		if err != nil {
			return nil, fmt.Errorf("failed to get trainings: %w", err)
		}

		// レスポンスのフォーマット
		return formatTrainingsByDateRangeResponse(response), nil
	}
}

// formatTrainingsByDateRangeResponse はレスポンスを見やすい形式にフォーマットします
func formatTrainingsByDateRangeResponse(response *dto.GetTrainingsByDateRangeResponse) string {
	if response.Count == 0 {
		return fmt.Sprintf("📊 **期間: %s**\n\n❌ この期間にトレーニング記録は見つかりませんでした。", response.Period)
	}

	result := fmt.Sprintf("📊 **期間: %s**\n\n🏋️ **トレーニング記録: %d件**\n\n", response.Period, response.Count)

	for i, training := range response.Trainings {
		result += fmt.Sprintf("**%d. %s (%s)**\n",
			i+1,
			training.Date.Format("2006-01-02"),
			training.Date.Weekday())

		if training.Notes != "" {
			result += fmt.Sprintf("📝 メモ: %s\n", training.Notes)
		}

		result += fmt.Sprintf("📈 概要: %d種目, %dセット, %.1fkg総ボリューム\n\n",
			training.Summary.TotalExercises,
			training.Summary.TotalSets,
			training.Summary.TotalVolume)

		// エクササイズの詳細
		for j, exercise := range training.Exercises {
			result += fmt.Sprintf("  **%d. %s (%s)**\n", j+1, exercise.Name, exercise.Category)

			// セットの詳細
			for k, set := range exercise.Sets {
				rpeText := ""
				if set.RPE != nil {
					rpeText = fmt.Sprintf(", RPE: %d", *set.RPE)
				}
				result += fmt.Sprintf("    Set %d: %.1fkg × %d回 (休憩: %ds%s)\n",
					k+1, set.WeightKg, set.Reps, set.RestTimeSeconds, rpeText)
			}
			result += "\n"
		}

		result += "---\n\n"
	}

	return result
}

// getPersonalRecordsSchema はget_personal_recordsツールのスキーマを返します
func getPersonalRecordsSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"exercise_name": map[string]interface{}{
				"type":        "string",
				"description": "特定のエクササイズ名（省略可）。指定すると該当エクササイズの記録のみを取得します。",
			},
		},
		"required": []string{}, // exercise_nameはオプション
	}
}

// createGetPersonalRecordsHandler はget_personal_recordsツールのハンドラーを作成します
func createGetPersonalRecordsHandler(queryHandler *handler.StrengthQueryHandler) ToolHandler {
	return func(arguments map[string]interface{}) (interface{}, error) {
		// 引数の取得（オプション）
		var exerciseName *string
		if name, ok := arguments["exercise_name"].(string); ok && name != "" {
			exerciseName = &name
		}

		// クエリの作成
		query := dto.GetPersonalRecordsQuery{
			ExerciseName: exerciseName,
		}

		// ハンドラーの実行
		response, err := queryHandler.GetPersonalRecords(query)
		if err != nil {
			return nil, fmt.Errorf("failed to get personal records: %w", err)
		}

		// レスポンスのフォーマット
		return formatPersonalRecordsResponse(response), nil
	}
}

// formatPersonalRecordsResponse は個人記録レスポンスを見やすい形式にフォーマットします
func formatPersonalRecordsResponse(response *dto.GetPersonalRecordsResponse) string {
	if response.Count == 0 {
		return "🏆 **個人記録**\n\n❌ 記録が見つかりませんでした。"
	}

	result := fmt.Sprintf("🏆 **個人記録 (%d種目)**\n\n", response.Count)

	for i, record := range response.Records {
		result += fmt.Sprintf("**%d. %s (%s)**\n", i+1, record.ExerciseName, record.Category)
		result += fmt.Sprintf("📊 総セッション数: %d回 | 最終実施: %s\n\n",
			record.TotalSessions,
			record.LastPerformed.Format("2006-01-02"))

		// 最大重量
		result += fmt.Sprintf("⚖️ **最大重量**: %.1fkg\n", record.MaxWeight.Value)
		result += fmt.Sprintf("   📅 達成日: %s (ID: %s)\n",
			record.MaxWeight.Date.Format("2006-01-02"),
			record.MaxWeight.TrainingID)
		if record.MaxWeight.SetDetails != nil {
			details := record.MaxWeight.SetDetails
			rpeText := ""
			if details.RPE != nil {
				rpeText = fmt.Sprintf(", RPE: %d", *details.RPE)
			}
			result += fmt.Sprintf("   🔍 セット詳細: %.1fkg × %d回 (休憩: %ds%s)\n",
				details.WeightKg, details.Reps, details.RestTimeSeconds, rpeText)
		}

		// 最大レップ数
		result += fmt.Sprintf("\n🔥 **最大レップ数**: %.0f回\n", record.MaxReps.Value)
		result += fmt.Sprintf("   📅 達成日: %s (ID: %s)\n",
			record.MaxReps.Date.Format("2006-01-02"),
			record.MaxReps.TrainingID)
		if record.MaxReps.SetDetails != nil {
			details := record.MaxReps.SetDetails
			rpeText := ""
			if details.RPE != nil {
				rpeText = fmt.Sprintf(", RPE: %d", *details.RPE)
			}
			result += fmt.Sprintf("   🔍 セット詳細: %.1fkg × %d回 (休憩: %ds%s)\n",
				details.WeightKg, details.Reps, details.RestTimeSeconds, rpeText)
		}

		// 最大ボリューム
		result += fmt.Sprintf("\n📊 **最大ボリューム**: %.1fkg\n", record.MaxVolume.Value)
		result += fmt.Sprintf("   📅 達成日: %s (ID: %s)\n",
			record.MaxVolume.Date.Format("2006-01-02"),
			record.MaxVolume.TrainingID)
		if record.MaxVolume.SetDetails != nil {
			details := record.MaxVolume.SetDetails
			rpeText := ""
			if details.RPE != nil {
				rpeText = fmt.Sprintf(", RPE: %d", *details.RPE)
			}
			result += fmt.Sprintf("   🔍 セット詳細: %.1fkg × %d回 (休憩: %ds%s)\n",
				details.WeightKg, details.Reps, details.RestTimeSeconds, rpeText)
		}

		if i < len(response.Records)-1 {
			result += "\n---\n\n"
		} else {
			result += "\n"
		}
	}

	return result
}
