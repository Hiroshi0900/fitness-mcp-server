package converter

import (
	query_dto "fitness-mcp-server/internal/application/query/dto"
	"fmt"
)

// FormatQueryResponse はクエリレスポンスを見やすい形式にフォーマットします
func FormatQueryResponse(response *query_dto.GetTrainingsByDateRangeResponse) string {
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

		result += fmt.Sprintf("📈 概要: %d種目, %dセット, %.1fkg総ボリューム\n",
			training.Summary.TotalExercises,
			training.Summary.TotalSets,
			training.Summary.TotalVolume)

		// エクササイズの概要のみ（詳細は省略）
		for _, exercise := range training.Exercises {
			result += fmt.Sprintf("  • %s (%s): %d sets\n",
				exercise.Name, exercise.Category, len(exercise.Sets))
		}
		result += "\n"
	}

	return result
}

// FormatPersonalRecordsResponse は個人記録レスポンスを見やすい形式にフォーマットします
func FormatPersonalRecordsResponse(response *query_dto.GetPersonalRecordsResponse) string {
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

		// 最大ボリューム
		result += fmt.Sprintf("\n📊 **最大ボリューム**: %.1fkg\n", record.MaxVolume.Value)
		result += fmt.Sprintf("   📅 達成日: %s (ID: %s)\n",
			record.MaxVolume.Date.Format("2006-01-02"),
			record.MaxVolume.TrainingID)

		if i < len(response.Records)-1 {
			result += "\n---\n\n"
		} else {
			result += "\n"
		}
	}

	return result
}
