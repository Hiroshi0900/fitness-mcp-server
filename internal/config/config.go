package config

import (
	"os"
	"path/filepath"
	"strconv"
)

// Config はアプリケーションの設定を管理します
type Config struct {
	Database DatabaseConfig `json:"database"`
	MCP      MCPConfig      `json:"mcp"`
	Server   ServerConfig   `json:"server"`
}

// DatabaseConfig はデータベース関連の設定です
type DatabaseConfig struct {
	SQLitePath      string `json:"sqlite_path"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime_hours"`
}

// MCPConfig はMCPサーバー関連の設定です
type MCPConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// ServerConfig はサーバー関連の設定です
type ServerConfig struct {
	Environment    string `json:"environment"`
	LogLevel       string `json:"log_level"`
	RequestTimeout int    `json:"request_timeout_seconds"`
}

// NewConfig は新しい設定を作成します
func NewConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			SQLitePath:      getDefaultDatabasePath(),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 2),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME_HOURS", 1),
		},
		MCP: MCPConfig{
			Name:        getEnvString("MCP_SERVER_NAME", "fitness-mcp-server"),
			Version:     getEnvString("MCP_SERVER_VERSION", "1.0.0"),
			Description: "筋トレ・ランニング記録管理MCPサーバー",
		},
		Server: ServerConfig{
			Environment:    getEnvString("APP_ENV", "development"),
			LogLevel:       getEnvString("LOG_LEVEL", "info"),
			RequestTimeout: getEnvInt("REQUEST_TIMEOUT_SECONDS", 30),
		},
	}
}

// getDefaultDatabasePath はデフォルトのデータベースパスを取得します
func getDefaultDatabasePath() string {
	// 環境変数でデータディレクトリが指定されている場合はそれを使用
	if dataDir := os.Getenv("MCP_DATA_DIR"); dataDir != "" {
		return filepath.Join(dataDir, "fitness.db")
	}

	// 実行ファイルのディレクトリを基準にしたパスを取得
	execPath, err := os.Executable()
	if err != nil {
		// 実行ファイルパスの取得に失敗した場合は現在の作業ディレクトリを使用
		workDir, err := os.Getwd()
		if err != nil {
			// 最後のフォールバック: 一時ディレクトリを使用
			return filepath.Join(os.TempDir(), "fitness-mcp", "fitness.db")
		}
		return filepath.Join(workDir, "data", "fitness.db")
	}

	// 実行ファイルのディレクトリ配下の data/fitness.db を使用
	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, "data", "fitness.db")
}

// getEnvString は環境変数から文字列を取得します（デフォルト値付き）
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt は環境変数から整数を取得します（デフォルト値付き）
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// EnsureDatabaseDir はデータベースディレクトリが存在することを確認します
func (c *Config) EnsureDatabaseDir() error {
	dbDir := filepath.Dir(c.Database.SQLitePath)
	return os.MkdirAll(dbDir, 0755)
}

// IsDevelopment は開発環境かどうかを判定します
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction は本番環境かどうかを判定します
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// Validate は設定の妥当性をチェックします
func (c *Config) Validate() error {
	// TODO: 設定値の妥当性チェックを実装
	// 例: データベースパスの有効性、設定値の範囲チェックなど
	return nil
}
