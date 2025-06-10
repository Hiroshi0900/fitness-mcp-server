package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"
	"time"

	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/domain/strength"
	"fitness-mcp-server/internal/interface/repository"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// StrengthRepository はSQLiteを使った筋トレRepository実装
type StrengthRepository struct {
	db *sql.DB
}

// NewStrengthTrainingRepository は新しいSQLite Repositoryを作成します
func NewStrengthTrainingRepository(db *sql.DB) (*StrengthRepository, error) {
	repo := &StrengthRepository{db: db}
	if err := repo.migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	return repo, nil
}

// NewStrengthRepository はファイルパスからSQLite Repositoryを作成します
func NewStrengthRepository(dbPath string) (repository.StrengthTrainingRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQLiteの設定
	db.SetMaxOpenConns(10)           // 複数接続を許可
	db.SetMaxIdleConns(2)            // アイドル接続数
	db.SetConnMaxLifetime(time.Hour) // 接続の最大生存時間

	return NewStrengthTrainingRepository(db)
}

// Initialize はデータベースの初期化（テーブル作成）を行います
func (r *StrengthRepository) Initialize() error {
	return r.migrate()
}

// migrate はマイグレーションを実行します
func (r *StrengthRepository) migrate() error {
	migrationSQL, err := migrationFiles.ReadFile("migrations/001_initial_schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = r.db.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

// Close はデータベース接続を閉じます
func (r *StrengthRepository) Close() error {
	return r.db.Close()
}

// Save は筋トレセッションを保存します
func (r *StrengthRepository) Save(training *strength.StrengthTraining) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 筋トレセッションを保存
	_, err = tx.Exec(`
		INSERT INTO strength_trainings (id, date, notes) 
		VALUES (?, ?, ?)`,
		training.ID().String(),
		training.Date(),
		training.Notes(),
	)
	if err != nil {
		return fmt.Errorf("failed to save training: %w", err)
	}

	// エクササイズを保存
	for exerciseOrder, exercise := range training.Exercises() {
		exerciseID, err := r.saveExercise(tx, training.ID(), exercise, exerciseOrder)
		if err != nil {
			return fmt.Errorf("failed to save exercise: %w", err)
		}

		// セットを保存
		for setOrder, set := range exercise.Sets() {
			if err := r.saveSet(tx, exerciseID, set, setOrder); err != nil {
				return fmt.Errorf("failed to save set: %w", err)
			}
		}
	}

	return tx.Commit()
}

// saveExercise はエクササイズを保存し、IDを返します
func (r *StrengthRepository) saveExercise(tx *sql.Tx, trainingID shared.TrainingID, exercise *strength.Exercise, order int) (int64, error) {
	result, err := tx.Exec(`
		INSERT INTO exercises (training_id, name, category, exercise_order) 
		VALUES (?, ?, ?, ?)`,
		trainingID.String(),
		exercise.Name().String(),
		exercise.Category().String(),
		order,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// saveSet はセットを保存します
func (r *StrengthRepository) saveSet(tx *sql.Tx, exerciseID int64, set strength.Set, order int) error {
	var rpe *int
	if set.RPE() != nil {
		rpeValue := set.RPE().Rating()
		rpe = &rpeValue
	}

	_, err := tx.Exec(`
		INSERT INTO sets (exercise_id, weight_kg, reps, rest_time_seconds, rpe, set_order) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		exerciseID,
		set.Weight().Kg(),
		set.Reps().Count(),
		int(set.RestTime().Duration().Seconds()),
		rpe,
		order,
	)
	return err
}

// FindByID はIDで筋トレセッションを検索します
func (r *StrengthRepository) FindByID(id shared.TrainingID) (*strength.StrengthTraining, error) {
	// 筋トレセッションを取得
	row := r.db.QueryRow(`
		SELECT id, date, notes 
		FROM strength_trainings 
		WHERE id = ?`, id.String())

	var idStr string
	var date time.Time
	var notes string

	if err := row.Scan(&idStr, &date, &notes); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("training not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to scan training: %w", err)
	}

	trainingID, err := shared.NewTrainingIDFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid training ID: %w", err)
	}

	training := strength.NewStrengthTraining(trainingID, date, notes)

	// エクササイズを取得
	exercises, err := r.findExercisesByTrainingID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load exercises: %w", err)
	}

	for _, exercise := range exercises {
		training.AddExercise(exercise)
	}

	return training, nil
}

// findExercisesByTrainingID はトレーニングIDでエクササイズを検索します
func (r *StrengthRepository) findExercisesByTrainingID(trainingID shared.TrainingID) ([]*strength.Exercise, error) {
	rows, err := r.db.Query(`
		SELECT id, name, category 
		FROM exercises 
		WHERE training_id = ? 
		ORDER BY exercise_order`, trainingID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*strength.Exercise
	for rows.Next() {
		var id int64
		var name, category string

		if err := rows.Scan(&id, &name, &category); err != nil {
			return nil, err
		}

		exerciseName, err := strength.NewExerciseName(name)
		if err != nil {
			return nil, fmt.Errorf("invalid exercise name: %w", err)
		}

		exerciseCategory, err := strength.NewExerciseCategory(category)
		if err != nil {
			return nil, fmt.Errorf("invalid exercise category: %w", err)
		}

		exercise := strength.NewExercise(exerciseName, exerciseCategory)

		// セットを取得
		sets, err := r.findSetsByExerciseID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to load sets: %w", err)
		}

		for _, set := range sets {
			exercise.AddSet(set)
		}

		exercises = append(exercises, exercise)
	}

	return exercises, rows.Err()
}

// findExercisesByTrainingIDs は複数のトレーニングIDでエクササイズを一括取得します
func (r *StrengthRepository) findExercisesByTrainingIDs(trainingIDs []string) (map[string][]*strength.Exercise, error) {
	if len(trainingIDs) == 0 {
		return make(map[string][]*strength.Exercise), nil
	}

	// IN句用のプレースホルダーを生成
	placeholders := make([]string, len(trainingIDs))
	args := make([]interface{}, len(trainingIDs))
	for i, id := range trainingIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, training_id, name, category 
		FROM exercises 
		WHERE training_id IN (%s) 
		ORDER BY training_id, exercise_order`,
		strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// エクササイズIDのリストを取得
	var exerciseIDs []int64
	exerciseDataMap := make(map[int64]struct {
		trainingID string
		name       string
		category   string
	})
	exercisesByTraining := make(map[string][]*strength.Exercise)

	for rows.Next() {
		var exerciseID int64
		var trainingID, name, category string

		if err := rows.Scan(&exerciseID, &trainingID, &name, &category); err != nil {
			return nil, err
		}

		exerciseIDs = append(exerciseIDs, exerciseID)
		exerciseDataMap[exerciseID] = struct {
			trainingID string
			name       string
			category   string
		}{trainingID, name, category}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 一括でセットを取得
	setsByExercise, err := r.findSetsByExerciseIDs(exerciseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load sets: %w", err)
	}

	// エクササイズオブジェクトを作成
	for exerciseID, data := range exerciseDataMap {
		exerciseName, err := strength.NewExerciseName(data.name)
		if err != nil {
			return nil, fmt.Errorf("invalid exercise name: %w", err)
		}

		exerciseCategory, err := strength.NewExerciseCategory(data.category)
		if err != nil {
			return nil, fmt.Errorf("invalid exercise category: %w", err)
		}

		exercise := strength.NewExercise(exerciseName, exerciseCategory)

		// セットを追加
		if sets, exists := setsByExercise[exerciseID]; exists {
			for _, set := range sets {
				exercise.AddSet(set)
			}
		}

		exercisesByTraining[data.trainingID] = append(exercisesByTraining[data.trainingID], exercise)
	}

	return exercisesByTraining, nil
}

// findSetsByExerciseID はエクササイズIDでセットを検索します
func (r *StrengthRepository) findSetsByExerciseID(exerciseID int64) ([]strength.Set, error) {
	rows, err := r.db.Query(`
		SELECT weight_kg, reps, rest_time_seconds, rpe 
		FROM sets 
		WHERE exercise_id = ? 
		ORDER BY set_order`, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []strength.Set
	for rows.Next() {
		var weightKg float64
		var reps int
		var restTimeSeconds int
		var rpe *int

		if err := rows.Scan(&weightKg, &reps, &restTimeSeconds, &rpe); err != nil {
			return nil, err
		}

		weight, err := strength.NewWeight(weightKg)
		if err != nil {
			return nil, fmt.Errorf("invalid weight: %w", err)
		}

		repsObj, err := strength.NewReps(reps)
		if err != nil {
			return nil, fmt.Errorf("invalid reps: %w", err)
		}

		restTime, err := strength.NewRestTime(time.Duration(restTimeSeconds) * time.Second)
		if err != nil {
			return nil, fmt.Errorf("invalid rest time: %w", err)
		}

		var rpeObj *strength.RPE
		if rpe != nil {
			rpeValue, err := strength.NewRPE(*rpe)
			if err != nil {
				return nil, fmt.Errorf("invalid RPE: %w", err)
			}
			rpeObj = &rpeValue
		}

		set := strength.NewSet(weight, repsObj, restTime, rpeObj)
		sets = append(sets, set)
	}

	return sets, rows.Err()
}

// findSetsByExerciseIDs は複数のエクササイズIDでセットを一括取得します
func (r *StrengthRepository) findSetsByExerciseIDs(exerciseIDs []int64) (map[int64][]strength.Set, error) {
	if len(exerciseIDs) == 0 {
		return make(map[int64][]strength.Set), nil
	}

	// IN句用のプレースホルダーを生成
	placeholders := make([]string, len(exerciseIDs))
	args := make([]interface{}, len(exerciseIDs))
	for i, id := range exerciseIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT exercise_id, weight_kg, reps, rest_time_seconds, rpe 
		FROM sets 
		WHERE exercise_id IN (%s) 
		ORDER BY exercise_id, set_order`,
		strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	setsByExercise := make(map[int64][]strength.Set)
	for rows.Next() {
		var exerciseID int64
		var weightKg float64
		var reps int
		var restTimeSeconds int
		var rpe *int

		if err := rows.Scan(&exerciseID, &weightKg, &reps, &restTimeSeconds, &rpe); err != nil {
			return nil, err
		}

		weight, err := strength.NewWeight(weightKg)
		if err != nil {
			return nil, fmt.Errorf("invalid weight: %w", err)
		}

		repsObj, err := strength.NewReps(reps)
		if err != nil {
			return nil, fmt.Errorf("invalid reps: %w", err)
		}

		restTime, err := strength.NewRestTime(time.Duration(restTimeSeconds) * time.Second)
		if err != nil {
			return nil, fmt.Errorf("invalid rest time: %w", err)
		}

		var rpeObj *strength.RPE
		if rpe != nil {
			rpeValue, err := strength.NewRPE(*rpe)
			if err != nil {
				return nil, fmt.Errorf("invalid RPE: %w", err)
			}
			rpeObj = &rpeValue
		}

		set := strength.NewSet(weight, repsObj, restTime, rpeObj)
		setsByExercise[exerciseID] = append(setsByExercise[exerciseID], set)
	}

	return setsByExercise, rows.Err()
}

// FindByDateRange は指定した期間の筋トレセッションを検索します
func (r *StrengthRepository) FindByDateRange(start, end time.Time) ([]*strength.StrengthTraining, error) {
	// 筋トレセッションを一括取得
	trainingRows, err := r.db.Query(`
		SELECT id, date, notes 
		FROM strength_trainings 
		WHERE date BETWEEN ? AND ? 
		ORDER BY date DESC`, start, end)
	if err != nil {
		return nil, err
	}
	defer trainingRows.Close()

	// トレーニングIDのリストを作成
	var trainingIDs []string
	trainingDataMap := make(map[string]struct {
		id    shared.TrainingID
		date  time.Time
		notes string
	})

	for trainingRows.Next() {
		var idStr string
		var date time.Time
		var notes string

		if err := trainingRows.Scan(&idStr, &date, &notes); err != nil {
			return nil, err
		}

		id, err := shared.NewTrainingIDFromString(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid training ID: %w", err)
		}

		trainingIDs = append(trainingIDs, idStr)
		trainingDataMap[idStr] = struct {
			id    shared.TrainingID
			date  time.Time
			notes string
		}{id, date, notes}
	}

	if err := trainingRows.Err(); err != nil {
		return nil, err
	}

	if len(trainingIDs) == 0 {
		return []*strength.StrengthTraining{}, nil
	}

	// 一括でエクササイズを取得
	exercisesByTraining, err := r.findExercisesByTrainingIDs(trainingIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load exercises: %w", err)
	}

	// 結果を組み立て
	var trainings []*strength.StrengthTraining
	for _, idStr := range trainingIDs {
		data := trainingDataMap[idStr]
		training := strength.NewStrengthTraining(data.id, data.date, data.notes)

		if exercises, exists := exercisesByTraining[idStr]; exists {
			for _, exercise := range exercises {
				training.AddExercise(exercise)
			}
		}

		trainings = append(trainings, training)
	}

	return trainings, nil
}

// FindByDate は指定した日の筋トレセッションを検索します
func (r *StrengthRepository) FindByDate(date time.Time) ([]*strength.StrengthTraining, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	return r.FindByDateRange(startOfDay, endOfDay)
}

// FindAll は全ての筋トレセッションを検索します
func (r *StrengthRepository) FindAll() ([]*strength.StrengthTraining, error) {
	rows, err := r.db.Query(`
		SELECT id 
		FROM strength_trainings 
		ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainings []*strength.StrengthTraining
	for rows.Next() {
		var idStr string
		if err := rows.Scan(&idStr); err != nil {
			return nil, err
		}

		id, err := shared.NewTrainingIDFromString(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid training ID: %w", err)
		}

		training, err := r.FindByID(id)
		if err != nil {
			return nil, err
		}

		trainings = append(trainings, training)
	}

	return trainings, rows.Err()
}

// Update は既存の筋トレセッションを更新します
func (r *StrengthRepository) Update(training *strength.StrengthTraining) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 筋トレセッションを更新
	_, err = tx.Exec(`
		UPDATE strength_trainings 
		SET date = ?, notes = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`,
		training.Date(),
		training.Notes(),
		training.ID().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update training: %w", err)
	}

	// 既存のエクササイズとセットを削除
	_, err = tx.Exec(`DELETE FROM exercises WHERE training_id = ?`, training.ID().String())
	if err != nil {
		return fmt.Errorf("failed to delete old exercises: %w", err)
	}

	// 新しいエクササイズとセットを保存
	for exerciseOrder, exercise := range training.Exercises() {
		exerciseID, err := r.saveExercise(tx, training.ID(), exercise, exerciseOrder)
		if err != nil {
			return fmt.Errorf("failed to save exercise: %w", err)
		}

		for setOrder, set := range exercise.Sets() {
			if err := r.saveSet(tx, exerciseID, set, setOrder); err != nil {
				return fmt.Errorf("failed to save set: %w", err)
			}
		}
	}

	return tx.Commit()
}

// Delete は筋トレセッションを削除します
func (r *StrengthRepository) Delete(id shared.TrainingID) error {
	result, err := r.db.Exec(`DELETE FROM strength_trainings WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete training: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("training not found: %s", id.String())
	}

	return nil
}

// ExistsById はIDの筋トレセッションが存在するかチェックします
func (r *StrengthRepository) ExistsById(id shared.TrainingID) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM strength_trainings WHERE id = ?`, id.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return count > 0, nil
}

// GetPersonalRecordsByExercise はエクササイズ別の自己ベスト（最大重量）を取得します
func (r *StrengthRepository) GetPersonalRecordsByExercise(exerciseName strength.ExerciseName) (strength.Weight, error) {
	var maxWeightKg sql.NullFloat64
	err := r.db.QueryRow(`
		SELECT MAX(s.weight_kg)
		FROM exercises e
		JOIN sets s ON e.id = s.exercise_id
		WHERE e.name = ?`, exerciseName.String()).Scan(&maxWeightKg)

	if err != nil {
		return strength.Weight{}, fmt.Errorf("failed to get personal record: %w", err)
	}

	if !maxWeightKg.Valid {
		// データが存在しない場合は0kgを返す
		return strength.NewWeight(0)
	}

	return strength.NewWeight(maxWeightKg.Float64)
}

// GetProgressAnalysis は指定したエクササイズの進捗分析を取得します
func (r *StrengthRepository) GetProgressAnalysis(exerciseName strength.ExerciseName, period time.Duration) (*repository.ProgressAnalysis, error) {
	endDate := time.Now()
	startDate := endDate.Add(-period)

	// 期間の開始と終了の重量を取得
	startWeight, err := r.getWeightAtDate(exerciseName, startDate, true) // 開始日以降の最初の記録
	if err != nil {
		return nil, fmt.Errorf("failed to get start weight: %w", err)
	}

	endWeight, err := r.getWeightAtDate(exerciseName, endDate, false) // 終了日以前の最後の記録
	if err != nil {
		return nil, fmt.Errorf("failed to get end weight: %w", err)
	}

	// ボリュームデータを取得
	startVolume, err := r.getVolumeAtDate(exerciseName, startDate, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get start volume: %w", err)
	}

	endVolume, err := r.getVolumeAtDate(exerciseName, endDate, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get end volume: %w", err)
	}

	// セッション数を取得
	totalSessions, err := r.getSessionCount(exerciseName, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get session count: %w", err)
	}

	// 改善トレンドを計算
	improvementTrend := r.calculateTrend(startWeight.Kg(), endWeight.Kg())

	// 増加量を計算
	weightIncrease := endWeight.Kg() - startWeight.Kg()
	var volumeIncrease float64
	if startVolume > 0 {
		volumeIncrease = ((endVolume - startVolume) / startVolume) * 100
	}

	return &repository.ProgressAnalysis{
		ExerciseName:     exerciseName,
		Period:           period,
		StartWeight:      startWeight,
		EndWeight:        endWeight,
		WeightIncrease:   weightIncrease,
		StartVolume:      startVolume,
		EndVolume:        endVolume,
		VolumeIncrease:   volumeIncrease,
		TotalSessions:    totalSessions,
		ImprovementTrend: improvementTrend,
	}, nil
}

// GetTrainingFrequency は指定期間のトレーニング頻度を取得します
func (r *StrengthRepository) GetTrainingFrequency(start, end time.Time) (*repository.TrainingFrequency, error) {
	// 総セッション数
	var totalSessions int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM strength_trainings
		WHERE date BETWEEN ? AND ?`, start, end).Scan(&totalSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to get total sessions: %w", err)
	}

	period := end.Sub(start)
	weeks := period.Hours() / (24 * 7)
	months := period.Hours() / (24 * 30) // 30日を1ヶ月として計算

	var sessionsPerWeek, sessionsPerMonth float64
	if weeks > 0 {
		sessionsPerWeek = float64(totalSessions) / weeks
	}
	if months > 0 {
		sessionsPerMonth = float64(totalSessions) / months
	}

	// 最もアクティブな曜日を取得
	mostActiveWeekday, err := r.getMostActiveWeekday(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get most active weekday: %w", err)
	}

	// 平均セッション長は現在のスキーマでは計算できないため、デフォルト値を設定
	// 実際の実装では、セッション開始時刻と終了時刻を記録する必要があります
	averageSessionLength := time.Hour // デフォルト値

	return &repository.TrainingFrequency{
		Period:               period,
		TotalSessions:        totalSessions,
		SessionsPerWeek:      sessionsPerWeek,
		SessionsPerMonth:     sessionsPerMonth,
		MostActiveWeekday:    mostActiveWeekday,
		AverageSessionLength: averageSessionLength,
	}, nil
}

// GetVolumeAnalysis は指定期間のボリューム分析を取得します
func (r *StrengthRepository) GetVolumeAnalysis(start, end time.Time) (*repository.VolumeAnalysis, error) {
	// 総ボリューム
	var totalVolume sql.NullFloat64
	err := r.db.QueryRow(`
		SELECT SUM(s.weight_kg * s.reps)
		FROM exercises e
		JOIN sets s ON e.id = s.exercise_id
		JOIN strength_trainings st ON e.training_id = st.id
		WHERE st.date BETWEEN ? AND ?`, start, end).Scan(&totalVolume)
	if err != nil {
		return nil, fmt.Errorf("failed to get total volume: %w", err)
	}

	// セッション数
	var sessionCount int
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM strength_trainings
		WHERE date BETWEEN ? AND ?`, start, end).Scan(&sessionCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get session count: %w", err)
	}

	var averageVolumePerSession float64
	if sessionCount > 0 && totalVolume.Valid {
		averageVolumePerSession = totalVolume.Float64 / float64(sessionCount)
	}

	// エクササイズ別ボリューム
	volumeByExercise, err := r.getVolumeByExercise(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume by exercise: %w", err)
	}

	// 週次ボリューム
	volumeByWeek, err := r.getVolumeByWeek(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume by week: %w", err)
	}

	// 成長率を計算
	volumeGrowthRate := r.calculateVolumeGrowthRate(volumeByWeek)

	period := end.Sub(start)
	totalVol := 0.0
	if totalVolume.Valid {
		totalVol = totalVolume.Float64
	}

	return &repository.VolumeAnalysis{
		Period:                  period,
		TotalVolume:             totalVol,
		AverageVolumePerSession: averageVolumePerSession,
		VolumeByExercise:        volumeByExercise,
		VolumeByWeek:            volumeByWeek,
		VolumeGrowthRate:        volumeGrowthRate,
	}, nil
}

// GetRecentTrainings は最近のトレーニングを取得します
func (r *StrengthRepository) GetRecentTrainings(limit int) ([]*strength.StrengthTraining, error) {
	rows, err := r.db.Query(`
		SELECT id 
		FROM strength_trainings 
		ORDER BY date DESC 
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainings []*strength.StrengthTraining
	for rows.Next() {
		var idStr string
		if err := rows.Scan(&idStr); err != nil {
			return nil, err
		}

		id, err := shared.NewTrainingIDFromString(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid training ID: %w", err)
		}

		training, err := r.FindByID(id)
		if err != nil {
			return nil, err
		}

		trainings = append(trainings, training)
	}

	return trainings, rows.Err()
}

// ヘルパーメソッド

// getWeightAtDate は指定日付での重量を取得します
func (r *StrengthRepository) getWeightAtDate(exerciseName strength.ExerciseName, date time.Time, after bool) (strength.Weight, error) {
	var query string
	if after {
		query = `
			SELECT s.weight_kg
			FROM exercises e
			JOIN sets s ON e.id = s.exercise_id
			JOIN strength_trainings st ON e.training_id = st.id
			WHERE e.name = ? AND st.date >= ?
			ORDER BY st.date ASC, s.weight_kg DESC
			LIMIT 1`
	} else {
		query = `
			SELECT s.weight_kg
			FROM exercises e
			JOIN sets s ON e.id = s.exercise_id
			JOIN strength_trainings st ON e.training_id = st.id
			WHERE e.name = ? AND st.date <= ?
			ORDER BY st.date DESC, s.weight_kg DESC
			LIMIT 1`
	}

	var weightKg sql.NullFloat64
	err := r.db.QueryRow(query, exerciseName.String(), date).Scan(&weightKg)
	if err != nil {
		if err == sql.ErrNoRows {
			return strength.NewWeight(0)
		}
		return strength.Weight{}, err
	}

	if !weightKg.Valid {
		return strength.NewWeight(0)
	}

	return strength.NewWeight(weightKg.Float64)
}

// getVolumeAtDate は指定日付でのボリュームを取得します
func (r *StrengthRepository) getVolumeAtDate(exerciseName strength.ExerciseName, date time.Time, after bool) (float64, error) {
	var query string
	if after {
		query = `
			SELECT SUM(s.weight_kg * s.reps)
			FROM exercises e
			JOIN sets s ON e.id = s.exercise_id
			JOIN strength_trainings st ON e.training_id = st.id
			WHERE e.name = ? AND st.date >= ?
			ORDER BY st.date ASC
			LIMIT 1`
	} else {
		query = `
			SELECT SUM(s.weight_kg * s.reps)
			FROM exercises e
			JOIN sets s ON e.id = s.exercise_id
			JOIN strength_trainings st ON e.training_id = st.id
			WHERE e.name = ? AND st.date <= ?
			ORDER BY st.date DESC
			LIMIT 1`
	}

	var volume sql.NullFloat64
	err := r.db.QueryRow(query, exerciseName.String(), date).Scan(&volume)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	if !volume.Valid {
		return 0, nil
	}

	return volume.Float64, nil
}

// getSessionCount は指定期間のセッション数を取得します
func (r *StrengthRepository) getSessionCount(exerciseName strength.ExerciseName, start, end time.Time) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(DISTINCT st.id)
		FROM exercises e
		JOIN strength_trainings st ON e.training_id = st.id
		WHERE e.name = ? AND st.date BETWEEN ? AND ?`,
		exerciseName.String(), start, end).Scan(&count)
	return count, err
}

// calculateTrend は改善トレンドを計算します
func (r *StrengthRepository) calculateTrend(start, end float64) string {
	diff := end - start
	threshold := 0.01 // 1%の閾値

	if diff > threshold {
		return "上昇"
	} else if diff < -threshold {
		return "下降"
	}
	return "停滞"
}

// getMostActiveWeekday は最もアクティブな曜日を取得します
func (r *StrengthRepository) getMostActiveWeekday(start, end time.Time) (string, error) {
	rows, err := r.db.Query(`
		SELECT 
			strftime('%w', date) as weekday,
			COUNT(*) as count
		FROM strength_trainings
		WHERE date BETWEEN ? AND ?
		GROUP BY weekday
		ORDER BY count DESC
		LIMIT 1`, start, end)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	weekdays := []string{"日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"}

	if rows.Next() {
		var weekdayNum, count int
		if err := rows.Scan(&weekdayNum, &count); err != nil {
			return "", err
		}
		return weekdays[weekdayNum], nil
	}

	return "データなし", nil
}

// getVolumeByExercise はエクササイズ別ボリュームを取得します
func (r *StrengthRepository) getVolumeByExercise(start, end time.Time) (map[string]float64, error) {
	rows, err := r.db.Query(`
		SELECT 
			e.name,
			SUM(s.weight_kg * s.reps) as total_volume
		FROM exercises e
		JOIN sets s ON e.id = s.exercise_id
		JOIN strength_trainings st ON e.training_id = st.id
		WHERE st.date BETWEEN ? AND ?
		GROUP BY e.name`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	volumeByExercise := make(map[string]float64)
	for rows.Next() {
		var name string
		var volume sql.NullFloat64
		if err := rows.Scan(&name, &volume); err != nil {
			return nil, err
		}
		if volume.Valid {
			volumeByExercise[name] = volume.Float64
		}
	}

	return volumeByExercise, rows.Err()
}

// getVolumeByWeek は週次ボリュームを取得します
func (r *StrengthRepository) getVolumeByWeek(start, end time.Time) ([]repository.WeeklyVolume, error) {
	rows, err := r.db.Query(`
		SELECT 
			strftime('%Y-%W', st.date) as week_string,
			MIN(st.date) as week_start,
			SUM(s.weight_kg * s.reps) as total_volume
		FROM exercises e
		JOIN sets s ON e.id = s.exercise_id
		JOIN strength_trainings st ON e.training_id = st.id
		WHERE st.date BETWEEN ? AND ?
		GROUP BY week_string
		ORDER BY week_string`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumeByWeek []repository.WeeklyVolume
	weekNum := 1
	for rows.Next() {
		var weekString string
		var weekStart time.Time
		var volume sql.NullFloat64

		if err := rows.Scan(&weekString, &weekStart, &volume); err != nil {
			return nil, err
		}

		vol := 0.0
		if volume.Valid {
			vol = volume.Float64
		}

		volumeByWeek = append(volumeByWeek, repository.WeeklyVolume{
			Week:   weekNum,
			Volume: vol,
			Date:   weekStart,
		})
		weekNum++
	}

	return volumeByWeek, rows.Err()
}

// calculateVolumeGrowthRate は週次ボリューム成長率を計算します
func (r *StrengthRepository) calculateVolumeGrowthRate(volumeByWeek []repository.WeeklyVolume) float64 {
	if len(volumeByWeek) < 2 {
		return 0
	}

	// 線形回帰による成長率の計算（簡易版）
	n := len(volumeByWeek)
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for i, weekly := range volumeByWeek {
		x := float64(i + 1) // 週番号
		y := weekly.Volume

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	nf := float64(n)
	slope := (nf*sumXY - sumX*sumY) / (nf*sumX2 - sumX*sumX)

	// 週次成長率として返す（パーセンテージ）
	if sumY/nf > 0 {
		return (slope / (sumY / nf)) * 100
	}
	return 0
}

// GetPersonalRecords は個人記録を取得します
func (r *StrengthRepository) GetPersonalRecords(exerciseName *string) ([]repository.PersonalRecordResult, error) {
	query := `
	WITH exercise_stats AS (
		SELECT 
			e.name as exercise_name,
			e.category,
			COUNT(DISTINCT st.id) as total_sessions,
			MAX(st.date) as last_performed
		FROM exercises e
		JOIN sets s ON e.id = s.exercise_id
		JOIN strength_trainings st ON e.training_id = st.id
		WHERE ($1 IS NULL OR e.name = $1)
		GROUP BY e.name, e.category
	),
	max_weight_details AS (
		SELECT DISTINCT
			e.name as exercise_name,
			s.weight_kg,
			s.reps,
			s.rest_time_seconds,
			s.rpe,
			st.date,
			st.id as training_id,
			ROW_NUMBER() OVER (PARTITION BY e.name ORDER BY s.weight_kg DESC, st.date DESC) as rn
		FROM exercises e
		JOIN sets s ON e.id = s.exercise_id
		JOIN strength_trainings st ON e.training_id = st.id
		WHERE ($1 IS NULL OR e.name = $1)
	),
	max_reps_details AS (
		SELECT DISTINCT
			e.name as exercise_name,
			s.weight_kg,
			s.reps,
			s.rest_time_seconds,
			s.rpe,
			st.date,
			st.id as training_id,
			ROW_NUMBER() OVER (PARTITION BY e.name ORDER BY s.reps DESC, st.date DESC) as rn
		FROM exercises e
		JOIN sets s ON e.id = s.exercise_id
		JOIN strength_trainings st ON e.training_id = st.id
		WHERE ($1 IS NULL OR e.name = $1)
	),
	max_volume_details AS (
		SELECT DISTINCT
			e.name as exercise_name,
			s.weight_kg,
			s.reps,
			s.rest_time_seconds,
			s.rpe,
			st.date,
			st.id as training_id,
			(s.weight_kg * s.reps) as volume,
			ROW_NUMBER() OVER (PARTITION BY e.name ORDER BY (s.weight_kg * s.reps) DESC, st.date DESC) as rn
		FROM exercises e
		JOIN sets s ON e.id = s.exercise_id
		JOIN strength_trainings st ON e.training_id = st.id
		WHERE ($1 IS NULL OR e.name = $1)
	)
	SELECT 
		es.exercise_name,
		es.category,
		COALESCE(mwd.weight_kg, 0) as max_weight,
		COALESCE(mwd.date, '1970-01-01') as max_weight_date,
		COALESCE(mwd.training_id, '') as max_weight_training_id,
		COALESCE(mwd.weight_kg, 0) as max_weight_details_weight,
		COALESCE(mwd.reps, 0) as max_weight_details_reps,
		COALESCE(mwd.rest_time_seconds, 0) as max_weight_details_rest,
		mwd.rpe as max_weight_details_rpe,
		COALESCE(mrd.reps, 0) as max_reps,
		COALESCE(mrd.date, '1970-01-01') as max_reps_date,
		COALESCE(mrd.training_id, '') as max_reps_training_id,
		COALESCE(mrd.weight_kg, 0) as max_reps_details_weight,
		COALESCE(mrd.reps, 0) as max_reps_details_reps,
		COALESCE(mrd.rest_time_seconds, 0) as max_reps_details_rest,
		mrd.rpe as max_reps_details_rpe,
		COALESCE(mvd.volume, 0) as max_volume,
		COALESCE(mvd.date, '1970-01-01') as max_volume_date,
		COALESCE(mvd.training_id, '') as max_volume_training_id,
		COALESCE(mvd.weight_kg, 0) as max_volume_details_weight,
		COALESCE(mvd.reps, 0) as max_volume_details_reps,
		COALESCE(mvd.rest_time_seconds, 0) as max_volume_details_rest,
		mvd.rpe as max_volume_details_rpe,
		es.total_sessions,
		es.last_performed
	FROM exercise_stats es
	LEFT JOIN max_weight_details mwd ON es.exercise_name = mwd.exercise_name AND mwd.rn = 1
	LEFT JOIN max_reps_details mrd ON es.exercise_name = mrd.exercise_name AND mrd.rn = 1
	LEFT JOIN max_volume_details mvd ON es.exercise_name = mvd.exercise_name AND mvd.rn = 1
	ORDER BY es.exercise_name;`

	rows, err := r.db.Query(query, exerciseName)
	if err != nil {
		return nil, fmt.Errorf("failed to query personal records: %w", err)
	}
	defer rows.Close()

	var records []repository.PersonalRecordResult
	for rows.Next() {
		var record repository.PersonalRecordResult
		var maxWeightDetails, maxRepsDetails, maxVolumeDetails repository.SetDetails
		var maxWeightDetailsRPE, maxRepsDetailsRPE, maxVolumeDetailsRPE sql.NullInt64
		
		// 日付を文字列として受け取る
		var maxWeightDateStr, maxRepsDateStr, maxVolumeDateStr, lastPerformedStr string

		err := rows.Scan(
			&record.ExerciseName,
			&record.Category,
			&record.MaxWeight.Value,
			&maxWeightDateStr,
			&record.MaxWeight.TrainingID,
			&maxWeightDetails.WeightKg,
			&maxWeightDetails.Reps,
			&maxWeightDetails.RestTimeSeconds,
			&maxWeightDetailsRPE,
			&record.MaxReps.Value,
			&maxRepsDateStr,
			&record.MaxReps.TrainingID,
			&maxRepsDetails.WeightKg,
			&maxRepsDetails.Reps,
			&maxRepsDetails.RestTimeSeconds,
			&maxRepsDetailsRPE,
			&record.MaxVolume.Value,
			&maxVolumeDateStr,
			&record.MaxVolume.TrainingID,
			&maxVolumeDetails.WeightKg,
			&maxVolumeDetails.Reps,
			&maxVolumeDetails.RestTimeSeconds,
			&maxVolumeDetailsRPE,
			&record.TotalSessions,
			&lastPerformedStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan personal record: %w", err)
		}

		// 日付文字列をtime.Timeに変換
		if maxWeightDate, err := time.Parse("2006-01-02 15:04:05", maxWeightDateStr); err == nil {
			record.MaxWeight.Date = maxWeightDate
		} else if maxWeightDate, err := time.Parse("2006-01-02", maxWeightDateStr); err == nil {
			record.MaxWeight.Date = maxWeightDate
		}

		if maxRepsDate, err := time.Parse("2006-01-02 15:04:05", maxRepsDateStr); err == nil {
			record.MaxReps.Date = maxRepsDate
		} else if maxRepsDate, err := time.Parse("2006-01-02", maxRepsDateStr); err == nil {
			record.MaxReps.Date = maxRepsDate
		}

		if maxVolumeDate, err := time.Parse("2006-01-02 15:04:05", maxVolumeDateStr); err == nil {
			record.MaxVolume.Date = maxVolumeDate
		} else if maxVolumeDate, err := time.Parse("2006-01-02", maxVolumeDateStr); err == nil {
			record.MaxVolume.Date = maxVolumeDate
		}

		if lastPerformed, err := time.Parse("2006-01-02 15:04:05", lastPerformedStr); err == nil {
			record.LastPerformed = lastPerformed
		} else if lastPerformed, err := time.Parse("2006-01-02", lastPerformedStr); err == nil {
			record.LastPerformed = lastPerformed
		}

		// RPEの設定（NULL許可のため）
		if maxWeightDetailsRPE.Valid {
			rpe := int(maxWeightDetailsRPE.Int64)
			maxWeightDetails.RPE = &rpe
		}
		if maxRepsDetailsRPE.Valid {
			rpe := int(maxRepsDetailsRPE.Int64)
			maxRepsDetails.RPE = &rpe
		}
		if maxVolumeDetailsRPE.Valid {
			rpe := int(maxVolumeDetailsRPE.Int64)
			maxVolumeDetails.RPE = &rpe
		}

		// セット詳細の設定
		record.MaxWeight.SetDetails = &maxWeightDetails
		record.MaxReps.SetDetails = &maxRepsDetails
		record.MaxVolume.SetDetails = &maxVolumeDetails

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate personal records: %w", err)
	}

	return records, nil
}

// コンパイル時のインターフェース実装チェック
var _ repository.StrengthTrainingRepository = (*StrengthRepository)(nil)
