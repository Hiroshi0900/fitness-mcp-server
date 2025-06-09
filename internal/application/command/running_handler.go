package command

import (
	"fmt"
	"log"
	"time"

	"fitness-mcp-server/internal/domain/running"
	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/interface/repository"
)

// =============================================================================
// ランニングコマンドハンドラー - ビジネスロジックの実行
// =============================================================================

// RunningCommandHandler はランニングに関するコマンドを処理するハンドラー
type RunningCommandHandler struct {
	sessionRepo repository.RunningSessionRepository
	goalRepo    repository.RunningGoalRepository
}

// NewRunningCommandHandler は新しいRunningCommandHandlerを作成します
func NewRunningCommandHandler(
	sessionRepo repository.RunningSessionRepository,
	goalRepo repository.RunningGoalRepository,
) *RunningCommandHandler {
	return &RunningCommandHandler{
		sessionRepo: sessionRepo,
		goalRepo:    goalRepo,
	}
}

// RecordRunningSession はランニングセッションを記録します
func (h *RunningCommandHandler) RecordRunningSession(cmd RecordRunningSessionCommand) (*RecordRunningSessionResult, error) {
	log.Printf("Recording running session for date: %s", cmd.Date.Format("2006-01-02"))

	// コマンドをドメインエンティティに変換
	session, err := cmd.ToRunningSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session entity: %w", err)
	}

	// リポジトリに保存
	if err := h.sessionRepo.Save(session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	log.Printf("Successfully recorded running session with ID: %s", session.ID().String())

	// 目標達成チェック
	goalAchieved, err := h.checkGoalAchievement(session)
	if err != nil {
		log.Printf("Warning: failed to check goal achievement: %v", err)
	}

	result := &RecordRunningSessionResult{
		SessionID:    session.ID().String(),
		Date:         session.Date(),
		Distance:     session.Distance().String(),
		Duration:     session.Duration().String(),
		Pace:         session.Pace().String(),
		Message:      fmt.Sprintf("ランニング（%s、%s、ペース: %s）を記録しました", session.Distance().String(), session.Duration().String(), session.Pace().String()),
		GoalAchieved: goalAchieved,
	}

	return result, nil
}

// UpdateRunningSession はランニングセッションを更新します
func (h *RunningCommandHandler) UpdateRunningSession(cmd UpdateRunningSessionCommand) (*UpdateRunningSessionResult, error) {
	log.Printf("Updating running session with ID: %s", cmd.ID)

	// コマンドをドメインエンティティに変換
	session, err := cmd.ToRunningSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session entity: %w", err)
	}

	// 既存のセッションが存在するかチェック
	exists, err := h.sessionRepo.ExistsById(session.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to check session existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("session not found: %s", cmd.ID)
	}

	// リポジトリで更新
	if err := h.sessionRepo.Update(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	log.Printf("Successfully updated running session with ID: %s", session.ID().String())

	return &UpdateRunningSessionResult{
		SessionID: session.ID().String(),
		Date:      session.Date(),
		Distance:  session.Distance().String(),
		Duration:  session.Duration().String(),
		Pace:      session.Pace().String(),
		Message:   fmt.Sprintf("ランニング（%s、%s、ペース: %s）を更新しました", session.Distance().String(), session.Duration().String(), session.Pace().String()),
	}, nil
}

// DeleteRunningSession はランニングセッションを削除します
func (h *RunningCommandHandler) DeleteRunningSession(cmd DeleteRunningSessionCommand) (*DeleteRunningSessionResult, error) {
	log.Printf("Deleting running session with ID: %s", cmd.ID)

	// コマンドバリデーション
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// IDをドメインオブジェクトに変換
	sessionID, err := shared.NewSessionIDFromString(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// 既存のセッションが存在するかチェック
	exists, err := h.sessionRepo.ExistsById(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to check session existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("session not found: %s", cmd.ID)
	}

	// リポジトリから削除
	if err := h.sessionRepo.Delete(sessionID); err != nil {
		return nil, fmt.Errorf("failed to delete session: %w", err)
	}

	log.Printf("Successfully deleted running session with ID: %s", cmd.ID)

	return &DeleteRunningSessionResult{
		SessionID: cmd.ID,
		Message:   "ランニングセッションを削除しました",
	}, nil
}

// CreateRunningGoal はランニング目標を作成します
func (h *RunningCommandHandler) CreateRunningGoal(cmd CreateRunningGoalCommand) (*CreateRunningGoalResult, error) {
	log.Printf("Creating running goal for event type: %s", cmd.EventType)

	// コマンドをドメインエンティティに変換
	goal, err := cmd.ToRunningGoal()
	if err != nil {
		return nil, fmt.Errorf("failed to create goal entity: %w", err)
	}

	// リポジトリに保存
	if err := h.goalRepo.Save(goal); err != nil {
		return nil, fmt.Errorf("failed to save goal: %w", err)
	}

	log.Printf("Successfully created running goal with ID: %s", goal.ID().String())

	return &CreateRunningGoalResult{
		GoalID:     goal.ID().String(),
		EventType:  goal.EventType().String(),
		TargetTime: goal.TargetTime().String(),
		TargetPace: goal.TargetPace().String(),
		Message:    fmt.Sprintf("ランニング目標（%s: %s、ペース: %s）を作成しました", goal.EventType().String(), goal.TargetTime().String(), goal.TargetPace().String()),
	}, nil
}

// UpdateRunningGoal はランニング目標を更新します
func (h *RunningCommandHandler) UpdateRunningGoal(cmd UpdateRunningGoalCommand) (*UpdateRunningGoalResult, error) {
	log.Printf("Updating running goal with ID: %s", cmd.ID)

	// コマンドバリデーション
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// IDをドメインオブジェクトに変換
	goalID, err := shared.NewGoalIDFromString(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid goal ID: %w", err)
	}

	// 既存の目標を取得
	goal, err := h.goalRepo.FindByID(goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to find goal: %w", err)
	}

	// 目標を更新
	if cmd.EventDate != nil {
		goal.SetEventDate(*cmd.EventDate)
	}
	if cmd.Description != "" {
		goal.UpdateDescription(cmd.Description)
	}

	// リポジトリで更新
	if err := h.goalRepo.Update(goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	log.Printf("Successfully updated running goal with ID: %s", goal.ID().String())

	return &UpdateRunningGoalResult{
		GoalID:  goal.ID().String(),
		Message: "ランニング目標を更新しました",
	}, nil
}

// DeleteRunningGoal はランニング目標を削除します
func (h *RunningCommandHandler) DeleteRunningGoal(cmd DeleteRunningGoalCommand) (*DeleteRunningGoalResult, error) {
	log.Printf("Deleting running goal with ID: %s", cmd.ID)

	// コマンドバリデーション
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// IDをドメインオブジェクトに変換
	goalID, err := shared.NewGoalIDFromString(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid goal ID: %w", err)
	}

	// リポジトリから削除
	if err := h.goalRepo.Delete(goalID); err != nil {
		return nil, fmt.Errorf("failed to delete goal: %w", err)
	}

	log.Printf("Successfully deleted running goal with ID: %s", cmd.ID)

	return &DeleteRunningGoalResult{
		GoalID:  cmd.ID,
		Message: "ランニング目標を削除しました",
	}, nil
}

// checkGoalAchievement は目標達成をチェックします
func (h *RunningCommandHandler) checkGoalAchievement(session *running.RunningSession) (*GoalAchieved, error) {
	// アクティブな目標を取得
	activeGoals, err := h.goalRepo.FindActiveGoals()
	if err != nil {
		return nil, fmt.Errorf("failed to find active goals: %w", err)
	}

	for _, goal := range activeGoals {
		// 目標達成チェック
		if err := goal.CheckAchievement(session); err != nil {
			log.Printf("Warning: failed to check achievement for goal %s: %v", goal.ID().String(), err)
			continue
		}

		// 達成された場合
		if goal.Status().Equals(running.Achieved) {
			// 目標を更新
			if err := h.goalRepo.Update(goal); err != nil {
				log.Printf("Warning: failed to update achieved goal %s: %v", goal.ID().String(), err)
			}

			// 達成情報を計算
			targetDuration := goal.TargetTime().Value()
			actualDuration := session.Duration().Value()
			improvement := targetDuration - actualDuration

			improvementStr := "達成！"
			if improvement > 0 {
				improvementStr = fmt.Sprintf("目標より%s早い達成！", formatDuration(improvement))
			}

			return &GoalAchieved{
				GoalID:      goal.ID().String(),
				EventType:   goal.EventType().String(),
				TargetTime:  goal.TargetTime().String(),
				ActualTime:  session.Duration().String(),
				Improvement: improvementStr,
			}, nil
		}
	}

	return nil, nil
}

// formatDuration は時間を読みやすい形式にフォーマットします
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if seconds == 0 {
		return fmt.Sprintf("%d分", minutes)
	}
	return fmt.Sprintf("%d分%d秒", minutes, seconds)
}

// RecordQuick5K は5km記録の簡単登録メソッド
func (h *RunningCommandHandler) RecordQuick5K(durationMin float64, notes string) (*RecordRunningSessionResult, error) {
	cmd := RecordRunningSessionCommand{
		Date:        time.Now(),
		DistanceKm:  5.0,
		DurationMin: durationMin,
		RunType:     "Easy",
		Notes:       notes,
	}

	return h.RecordRunningSession(cmd)
}

// RecordQuick10K は10km記録の簡単登録メソッド
func (h *RunningCommandHandler) RecordQuick10K(durationMin float64, notes string) (*RecordRunningSessionResult, error) {
	cmd := RecordRunningSessionCommand{
		Date:        time.Now(),
		DistanceKm:  10.0,
		DurationMin: durationMin,
		RunType:     "Easy",
		Notes:       notes,
	}

	return h.RecordRunningSession(cmd)
}

// CreateHalfMarathonGoal はハーフマラソン目標の簡単作成メソッド
func (h *RunningCommandHandler) CreateHalfMarathonGoal(targetTimeMin float64, eventDate *time.Time, description string) (*CreateRunningGoalResult, error) {
	cmd := CreateRunningGoalCommand{
		EventType:     "Half",
		TargetTimeMin: targetTimeMin,
		EventDate:     eventDate,
		Description:   description,
	}

	return h.CreateRunningGoal(cmd)
}
