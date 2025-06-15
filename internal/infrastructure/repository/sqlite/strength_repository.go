package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	"fitness-mcp-server/internal/domain/shared"
	"fitness-mcp-server/internal/domain/strength"
	"fitness-mcp-server/internal/interface/repository"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// StrengthRepository はSQLiteを使った筋トレRepository実装（書き込み専用）
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
	log.Printf("Creating SQLite repository with path: %s", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Printf("Failed to open SQLite database: %v", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping SQLite database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// SQLiteの設定
	db.SetMaxOpenConns(10)           // 複数接続を許可
	db.SetMaxIdleConns(2)            // アイドル接続数
	db.SetConnMaxLifetime(time.Hour) // 接続の最大生存時間

	log.Printf("SQLite database opened successfully: %s", dbPath)
	return NewStrengthTrainingRepository(db)
}

// Initialize はデータベースの初期化（テーブル作成）を行います
func (r *StrengthRepository) Initialize() error {
	return r.migrate()
}

// migrate はマイグレーションを実行します
func (r *StrengthRepository) migrate() error {
	log.Printf("Starting database migration...")

	// マイグレーション状態追跡テーブルを作成
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("Failed to create migrations table: %v", err)
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// マイグレーションファイルのリスト（順序重要）
	migrations := []struct {
		version  string
		filename string
	}{
		{"001", "migrations/001_initial_schema.sql"},
		{"002", "migrations/002_remove_resttime_category.sql"},
		{"003", "migrations/003_add_running_tables.sql"},
	}

	for _, migration := range migrations {
		// 既に適用済みかチェック
		var count int
		err := r.db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, migration.version).Scan(&count)
		if err != nil {
			log.Printf("Failed to check migration status for %s: %v", migration.version, err)
			return fmt.Errorf("failed to check migration status for %s: %w", migration.version, err)
		}

		if count > 0 {
			log.Printf("Migration %s already applied, skipping", migration.version)
			continue
		}

		log.Printf("Executing migration: %s", migration.filename)

		migrationSQL, err := migrationFiles.ReadFile(migration.filename)
		if err != nil {
			log.Printf("Failed to read migration file %s: %v", migration.filename, err)
			return fmt.Errorf("failed to read migration file %s: %w", migration.filename, err)
		}

		_, err = r.db.Exec(string(migrationSQL))
		if err != nil {
			log.Printf("Failed to execute migration %s: %v", migration.filename, err)
			return fmt.Errorf("failed to execute migration %s: %w", migration.filename, err)
		}

		// マイグレーション完了を記録
		_, err = r.db.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, migration.version)
		if err != nil {
			log.Printf("Failed to record migration completion for %s: %v", migration.version, err)
			return fmt.Errorf("failed to record migration completion for %s: %w", migration.version, err)
		}

		log.Printf("Successfully executed migration: %s", migration.filename)
	}

	log.Printf("Database migration completed successfully")
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

// プライベートヘルパーメソッド

// saveExercise はエクササイズを保存し、IDを返します
func (r *StrengthRepository) saveExercise(tx *sql.Tx, trainingID shared.TrainingID, exercise *strength.Exercise, order int) (int64, error) {
	result, err := tx.Exec(`
		INSERT INTO exercises (training_id, name, exercise_order) 
		VALUES (?, ?, ?)`,
		trainingID.String(),
		exercise.Name().String(),
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
		INSERT INTO sets (exercise_id, weight_kg, reps, rpe, set_order) 
		VALUES (?, ?, ?, ?, ?)`,
		exerciseID,
		set.Weight().Kg(),
		set.Reps().Count(),
		rpe,
		order,
	)
	return err
}

// コンパイル時のインターフェース実装チェック
var _ repository.StrengthTrainingRepository = (*StrengthRepository)(nil)
