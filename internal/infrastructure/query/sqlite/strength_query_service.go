package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"fitness-mcp-server/internal/application/query/dto"
	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/domain/strength"
	"fitness-mcp-server/internal/interface/query"

	_ "modernc.org/sqlite"
)

// StrengthQueryService はSQLiteを使った筋トレクエリサービス実装
type StrengthQueryService struct {
	db *sql.DB
}

// NewStrengthQueryService は新しいSQLite クエリサービスを作成します
func NewStrengthQueryService(db *sql.DB) *StrengthQueryService {
	return &StrengthQueryService{db: db}
}

// FindByID はIDで筋トレセッションを検索します
func (s *StrengthQueryService) FindByID(id shared.TrainingID) (*strength.StrengthTraining, error) {
	// 筋トレセッションを取得
	row := s.db.QueryRow(`
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
	exercises, err := s.findExercisesByTrainingID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load exercises: %w", err)
	}

	for _, exercise := range exercises {
		training.AddExercise(exercise)
	}

	return training, nil
}

// FindByDateRange は指定した期間の筋トレセッションを検索します
func (s *StrengthQueryService) FindByDateRange(start, end time.Time) ([]*strength.StrengthTraining, error) {
	// 筋トレセッションを一括取得
	trainingRows, err := s.db.Query(`
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
	exercisesByTraining, err := s.findExercisesByTrainingIDs(trainingIDs)
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
func (s *StrengthQueryService) FindByDate(date time.Time) ([]*strength.StrengthTraining, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	return s.FindByDateRange(startOfDay, endOfDay)
}

// FindAll は全ての筋トレセッションを検索します
func (s *StrengthQueryService) FindAll() ([]*strength.StrengthTraining, error) {
	rows, err := s.db.Query(`
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

		training, err := s.FindByID(id)
		if err != nil {
			return nil, err
		}

		trainings = append(trainings, training)
	}

	return trainings, rows.Err()
}

// ExistsById はIDの筋トレセッションが存在するかチェックします
func (s *StrengthQueryService) ExistsById(id shared.TrainingID) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM strength_trainings WHERE id = ?`, id.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return count > 0, nil
}

// GetPersonalRecords は個人記録を取得します
func (s *StrengthQueryService) GetPersonalRecords(exerciseName *string) ([]dto.PersonalRecordQueryResult, error) {
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

	rows, err := s.db.Query(query, exerciseName)
	if err != nil {
		return nil, fmt.Errorf("failed to query personal records: %w", err)
	}
	defer rows.Close()

	var records []dto.PersonalRecordQueryResult
	for rows.Next() {
		var record dto.PersonalRecordQueryResult
		var maxWeightDetails, maxRepsDetails, maxVolumeDetails dto.SetQueryDetails
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

// プライベートヘルパーメソッド

// findExercisesByTrainingID はトレーニングIDでエクササイズを検索します
func (s *StrengthQueryService) findExercisesByTrainingID(trainingID shared.TrainingID) ([]*strength.Exercise, error) {
	rows, err := s.db.Query(`
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
		sets, err := s.findSetsByExerciseID(id)
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
func (s *StrengthQueryService) findExercisesByTrainingIDs(trainingIDs []string) (map[string][]*strength.Exercise, error) {
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

	rows, err := s.db.Query(query, args...)
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
	setsByExercise, err := s.findSetsByExerciseIDs(exerciseIDs)
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
func (s *StrengthQueryService) findSetsByExerciseID(exerciseID int64) ([]strength.Set, error) {
	rows, err := s.db.Query(`
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
func (s *StrengthQueryService) findSetsByExerciseIDs(exerciseIDs []int64) (map[int64][]strength.Set, error) {
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

	rows, err := s.db.Query(query, args...)
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

// コンパイル時のインターフェース実装チェック
var _ query.StrengthQueryService = (*StrengthQueryService)(nil)
