package tool

import (
	"context"
	"fitness-mcp-server/internal/application/command/dto"
	"fitness-mcp-server/internal/application/command/handler"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// TrainingToolHandler はトレーニング記録ツールを管理します
type TrainingToolHandler struct {
	commandHandler *handler.StrengthCommandHandler
}

// NewTrainingToolHandler は新しいTrainingToolHandlerを作成します
func NewTrainingToolHandler(commandHandler *handler.StrengthCommandHandler) *TrainingToolHandler {
	return &TrainingToolHandler{
		commandHandler: commandHandler,
	}
}

// Register はトレーニング記録ツールを登録します
func (h *TrainingToolHandler) Register(s *server.MCPServer) error {
	tool := mcp.NewTool(
		"record_training",
		mcp.WithDescription(`筋トレセッションの記録を管理するツール。実施したエクササイズ、セット数、重量、回数を記録できます。

【使用例】
- ベンチプレス 80kg×10回を3セット実施した場合
- スクワット 100kg×8回で実施した場合  
- 複数のエクササイズを一つのセッションとして記録する場合`),
		mcp.WithString("date",
			mcp.Required(),
			mcp.Description("トレーニング実施日付。YYYY-MM-DD形式で指定してください。例: 2024-06-14"),
		),
		mcp.WithArray("exercises",
			mcp.Required(),
			mcp.Description(`実施したエクササイズのリスト。各エクササイズには以下を含める必要があります:

【エクササイズオブジェクト】
{
  "name": "エクササイズ名（例: ベンチプレス、スクワット、デッドリフト、ダンベルカール等）",
  "sets": [セット配列]
}

【setオブジェクト】
{
  "weight_kg": 使用重量（kg、数値）,
  "reps": 実施回数（回、整数）,
  "rpe": RPE値（1-10、省略可）
}

【RPEについて】
RPE（Rate of Perceived Exertion）は主観的運動強度です。
- 1-3: 非常に楽
- 4-6: 楽〜やや楽  
- 7-8: きつい
- 9-10: 非常にきつい〜限界`),
		),
		mcp.WithString("notes",
			mcp.Description("セッション全体のメモや備考（省略可）。例: 調子良い、フォーム意識、疲労感あり等"),
		),
	)

	toolHandler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return h.handleRecordTraining(ctx, req)
	}

	s.AddTool(tool, toolHandler)
	return nil
}

// handleRecordTraining はトレーニング記録処理を行います
func (h *TrainingToolHandler) handleRecordTraining(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// パラメータマップの取得
	paramsMap, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("パラメータが不正です"), nil
	}

	// 日付の取得
	dateStr, err := req.RequireString("date")
	if err != nil {
		return mcp.NewToolResultError("dateパラメータが必要です: " + err.Error()), nil
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return mcp.NewToolResultError("日付の形式が不正です（YYYY-MM-DD形式で入力してください）: " + err.Error()), nil
	}

	// エクササイズの解析
	exercises, err := h.parseExercises(paramsMap)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// ノート（オプション）
	notes := ""
	if notesData, exists := paramsMap["notes"]; exists {
		if notesStr, ok := notesData.(string); ok {
			notes = notesStr
		}
	}

	// RecordTrainingCommandの作成
	cmd := dto.RecordTrainingCommand{
		Date:      date,
		Exercises: exercises,
		Notes:     notes,
	}

	// バリデーション
	if err := cmd.Validate(); err != nil {
		return mcp.NewToolResultError("データが不正です: " + err.Error()), nil
	}

	result, err := h.commandHandler.RecordTraining(cmd)
	if err != nil {
		return mcp.NewToolResultError("記録に失敗しました: " + err.Error()), nil
	}

	// 結果をテキストで返す
	return mcp.NewToolResultText(
		fmt.Sprintf("記録完了: TrainingID=%v, メッセージ=%v", result.TrainingID, result.Message),
	), nil
}

// parseExercises はリクエストからエクササイズ情報を解析します
func (h *TrainingToolHandler) parseExercises(paramsMap map[string]interface{}) ([]dto.ExerciseDTO, error) {
	exercisesData, ok := paramsMap["exercises"]
	if !ok {
		return nil, fmt.Errorf("exercisesパラメータが必要です")
	}

	exercisesSlice, ok := exercisesData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("exercisesは配列である必要があります")
	}

	var exercises []dto.ExerciseDTO
	for _, exerciseData := range exercisesSlice {
		exerciseMap, ok := exerciseData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("exercise要素が不正です")
		}

		// エクササイズ名の取得
		name, ok := exerciseMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("exercise nameが必要です")
		}

		// セットの解析
		sets, err := h.parseSets(exerciseMap)
		if err != nil {
			return nil, err
		}

		exercises = append(exercises, dto.ExerciseDTO{
			Name: name,
			Sets: sets,
		})
	}

	return exercises, nil
}

// parseSets はエクササイズマップからセット情報を解析します
func (h *TrainingToolHandler) parseSets(exerciseMap map[string]interface{}) ([]dto.SetDTO, error) {
	setsData, ok := exerciseMap["sets"]
	if !ok {
		return nil, fmt.Errorf("setsが必要です")
	}

	setsSlice, ok := setsData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("setsは配列である必要があります")
	}

	var sets []dto.SetDTO
	for _, setData := range setsSlice {
		setMap, ok := setData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("set要素が不正です")
		}

		// 重量、回数の取得
		weightKg, ok := setMap["weight_kg"].(float64)
		if !ok {
			return nil, fmt.Errorf("weight_kgが必要です")
		}

		repsFloat, ok := setMap["reps"].(float64)
		if !ok {
			return nil, fmt.Errorf("repsが必要です")
		}
		reps := int(repsFloat)

		// RPE（オプション）
		var rpe *int
		if rpeData, exists := setMap["rpe"]; exists {
			if rpeFloat, ok := rpeData.(float64); ok {
				rpeInt := int(rpeFloat)
				if rpeInt < 1 || rpeInt > 10 {
					return nil, fmt.Errorf("RPEは1-10の範囲で指定してください（1:非常に楽 〜 10:限界）")
				}
				rpe = &rpeInt
			}
		}

		sets = append(sets, dto.SetDTO{
			WeightKg: weightKg,
			Reps:     reps,
			RPE:      rpe,
		})
	}

	return sets, nil
}
