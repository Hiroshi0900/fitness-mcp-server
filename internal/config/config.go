package config

import (
	"os"
	"path/filepath"
)

// Config はアプリケーションの設定を管理します
type Config struct {
	Database DatabaseConfig `json:"database"`
	MCP      MCPConfig      `json:"mcp"`
}

// DatabaseConfig はデータベース関連の設定です
type DatabaseConfig struct {
	SQLitePath string `json:"sqlite_path"`
}

// MCPConfig はMCPサーバー関連の設定です
type MCPConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// NewConfig は新しい設定を作成します
func NewConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			SQLitePath: getDefaultDatabasePath(),
		},
		MCP: MCPConfig{
			Name:        "fitness-mcp-server",
			Version:     "1.0.0",
			Description: "筋トレ・ランニング記録管理MCPサーバー",
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

// EnsureDatabaseDir はデータベースディレクトリが存在することを確認します
func (c *Config) EnsureDatabaseDir() error {
	dbDir := filepath.Dir(c.Database.SQLitePath)
	return os.MkdirAll(dbDir, 0755)
}
