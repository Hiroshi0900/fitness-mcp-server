# Go言語でMCPサーバを実装してみた - 筋トレ記録をAIに覚えてもらう

何番煎じか、という感じはしますが、MCPサーバを実際に作ってみて理解を深めようということで書いていこうと思います。

---

## そもそもMCPサーバとは？

まずMCPサーバとは何か？ということについて触れておきます。

**参考:**
- https://zenn.dev/cloud_ace/articles/model-context-protocol
- https://docs.anthropic.com/ja/docs/agents-and-tools/mcp

MCPとは、Model Context Protocolの略称です。LLMと外部のデータやツールを接続するためのプロトコルのことを指しています。

最初はなんとなくわかりそうでわからなかった自分がいましたが、AnthropicのドキュメントでAIアプリケーション用のUSB-Cポートのようなものと説明されており、実際に使用する中で理解が深まりました。

### MCPサーバが注目される理由

こちらの記事が参考になります：https://zenn.dev/zamis/articles/73fe4c6e9f289e

優れたLLMでも理解できない情報があります。MCPサーバは、各ツールや仕組みごとにAIに知識を渡すための機構を提供します。社内ツールとの連携などがまさにその一例です。

AIは非常に賢く、最近ではWeb検索も利用できるようになって最新情報を取得できるようになっていますが、秘匿されている情報や本当に最新の重要な情報にはアクセスできません。そこでMCPサーバが活躍します。

## 今回作るもの

今回は筋トレの記録機能を作ります。

シンプルなCRUDともいえますが、実際にトレーニングの相談や経過報告をChatGPTでよく行っていることもあり、過去の記録や目標が記憶から消されることがしばしばあるため、パーソナルトレーナーのような立て付けで使えるといいなと思い、この題材を選びました。

生成AIは一定期間は会話のコンテキストを把握してくれて、ある程度の情報は残してくれますが、会話が続くと以前の記憶が消えてしまいます。それを補完したい意図もあります。

## 技術スタック

- **言語**: Go
- **DB**: SQLite
- **MCP**: https://github.com/mark3labs/mcp-go
- **CQRS（簡易）**
- **DDD（軽量）**

軽量DDDやCQRSを使った構成にしていますが、これは後の拡張性を考慮して選定しました。今回の本題とは逸れるため、詳細はほとんど割愛します。

SQLiteの採用について、大きなこだわりはありませんが、まずローカルで動かすことを目指している背景と、再起動ごとにデータが初期化されてしまうのを避けたかったため採用しました。クリーンアーキテクチャな構成にもしているので、リモートで動かす際は別のDBクライアントへの切り替えも検討しています。

**mark3labs/mcp-go**を今回MCPサーバのライブラリとして採用しています。これは様々なライブラリと比較検討して選んだわけではなく、一番有名そうで、特に選ばない理由もなさそう、という程度の背景です。

また、今回は微妙な調整を私が行っていますが、テストや実装、ドキュメント整備などはClaudeCodeで行っています。雑に組んでしまっている部分もあるため、細かい点はご容赦ください。

## 実装について

今回のアプリケーション全体は以下のような構成になります。

```shell
.
├── cmd/mcp/main.go              # MCPサーバのエントリーポイント
├── data/fitness.db              # SQLiteデータベースファイル
├── internal/
│   ├── application/
│   │   ├── command/             # CQRSのCommand側（書き込み処理）
│   │   │   ├── dto/            
│   │   │   ├── handler/        
│   │   │   └── usecase/        
│   │   └── query/               # CQRSのQuery側（読み込み処理）
│   │       ├── dto/            
│   │       ├── handler/        
│   │       └── usecase/        
│   ├── domain/
│   │   ├── strength/            # 筋トレドメイン
│   │   ├── running/             # 今回は未実装
│   │   └── shared/             
│   ├── infrastructure/
│   │   ├── query/sqlite/        # クエリ用SQLite実装
│   │   └── repository/sqlite/   # コマンド用SQLite実装
│   └── interface/               # インターフェース層
│       ├── mcp/                 # MCP関連の処理
│       │   ├── converter/       # レスポンス整形処理
│       │   └── tool/            # MCPツール実装
│       ├── query/              
│       └── repository/         
└── mcp                          # MCPサーバ実行ファイル
```

※さらに拡張する可能性があるため、執筆時点での構成です

今回お話するのは基本的に`cmd/mcp/main.go`と`interface/mcp/tool`になります。`interface/mcp/converter`はリクエスト・レスポンスを整えるヘルパー処理のようなもので、MCPサーバを立てる文脈とはそれるため、触れません。

## 実装するMCPサーバの機能

今回のMCPサーバが提供する機能を以下にまとめました：

| Tool名 | 機能概要 | 主要パラメータ | 用途・効果 |
|--------|----------|----------------|------------|
| `record_training` | 筋トレの記録 | `date`（実施日）<br>`exercises`（エクササイズ配列）<br>`notes`（メモ） | トレーニング実績をAIが記憶<br>継続的な進捗管理が可能 |
| `get_trainings_by_date_range` | 期間指定での履歴取得 | `start_date`（開始日）<br>`end_date`（終了日） | 過去の実績を基にした<br>トレーニング相談・分析 |
| `get_personal_records` | 個人記録の参照 | `exercise_name`（種目名、任意） | 現在の最高記録を考慮した<br>目標設定とアドバイス |

### 記録できるデータの詳細

各エクササイズでは以下の詳細情報を記録できます：

| 項目 | 説明 | 例 |
|------|------|-----|
| エクササイズ名 | 実施した種目名 | ベンチプレス、スクワット、デッドリフト |
| カテゴリ | 種目の分類 | Compound（複合）/ Isolation（単関節）/ Cardio（有酸素） |
| 重量 | 使用重量（kg） | 80.0 |
| 回数 | 実施回数 | 10 |
| 休憩時間 | セット間休憩（秒） | 180 |
| RPE | 主観的運動強度（1-10） | 8（きつい） |

機能の命名や項目は大枠を私が決めましたが、細かいところはAIに委ねました。私はRPEを気にして筋トレを行ったことはありませんし、「Isolationをやっているな」なども考えたことはありませんが、有酸素と筋トレを計測したいと伝えたらこのような項目にしてくれました。

## コードについて

### アプリケーション層のハンドラー

まず、ビジネスロジックを処理するアプリケーション層のハンドラーです：

```go:internal/application/command/handler/strength_handler.go
package handler

import (
	"fitness-mcp-server/internal/application/command/dto"
	"fitness-mcp-server/internal/application/command/usecase"
)

// StrengthCommandHandler は筋トレに関するコマンドを処理するハンドラー
type StrengthCommandHandler struct {
	usecase usecase.StrengthTrainingUsecase
}

// NewStrengthCommandHandler は新しいStrengthCommandHandlerを作成します
func NewStrengthCommandHandler(usecase usecase.StrengthTrainingUsecase) *StrengthCommandHandler {
	return &StrengthCommandHandler{
		usecase: usecase,
	}
}

// RecordTraining は筋トレセッションを記録します
func (h *StrengthCommandHandler) RecordTraining(cmd dto.RecordTrainingCommand) (*dto.RecordTrainingResult, error) {
	return h.usecase.RecordTraining(cmd)
}
```

ここではユースケースを利用できるAPIのハンドラーとして定義しています。ユースケースは永続化処理を行うだけなので、ここでは割愛します。

### MCPツールの実装

次に、MCPサーバの核心部分であるツール実装です：

```go:internal/interface/mcp/tool/training_tool.go
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

// Register はトレーニング記録ツールを登録します
func (h *TrainingToolHandler) Register(s *server.MCPServer) error {
	tool := mcp.NewTool(
		"record_training",
		mcp.WithDescription(`筋トレセッションの記録を管理するツール。実施したエクササイズ、セット数、重量、回数、休憩時間を記録できます。

【使用例】
- ベンチプレス 80kg×10回を3セット実施した場合
- スクワット 100kg×8回、休憩180秒で実施した場合  
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
  "category": "エクササイズカテゴリ（必須）",
  "sets": [セット配列]
}

【categoryの選択肢】
- "Compound": 複合種目（ベンチプレス、スクワット、デッドリフト等）
- "Isolation": 単関節種目（ダンベルカール、レッグエクステンション等）  
- "Cardio": 有酸素運動（ランニング、バイク等）

【setオブジェクト】
{
  "weight_kg": 使用重量（kg、数値）,
  "reps": 実施回数（回、整数）,
  "rest_time_seconds": 休憩時間（秒、整数）,
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

// handleRecordTraining は実際の処理（パラメータ解析〜ビジネスロジック呼び出し）
// 詳細実装は省略...
```

#### MCPツール実装のポイント

通常のAPIで言うところの、サーバへのエンドポイントを設定する部分です。MCPライブラリの詳細はパッケージの内容を見ていただく方が良いと思いますが、以下のようなことを行っています：

1. **`mcp.NewTool`として新しい機能を定義する**
   - 機能名、リクエストで受け取る情報、機能の説明などを登録
   - 機能の説明にあたる内容は今回AIに拡充させましたが、AIはここの内容を見てどのMCPサーバをどのように使うかを決定します

2. **ハンドラー関数を作成し、`AddTool`で設定を完了**

**MCPツール実装の特徴：**

**1. Tool定義とハンドラーの分離**
```go
tool := mcp.NewTool(...) // Tool定義（AIが認識するパラメータや説明）
toolHandler := func(...) { ... } // 実際のリクエスト処理ロジック
s.AddTool(tool, toolHandler) // サーバへの登録
```

**2. MCPサーバとビジネスロジックの分離**
```go
type TrainingToolHandler struct {
    commandHandler *handler.StrengthCommandHandler // ビジネスロジックへの依存
}
```
MCPツール層は既存のアプリケーション層のハンドラーを呼び出すだけで、ビジネスロジックは一切持っていません。これにより、**既存のWebAPIやCLIアプリケーションがあれば、MCP層を追加するだけで簡単にMCP化できる**設計になっています。

**3. 詳細なパラメータ説明**
AIが理解しやすいように、Tool定義のDescriptionに詳細な説明を記載しています。AIはここの内容を見て、どのMCPサーバをどのように使うかを決定します。

### MCPサーバのメイン処理

最後に、MCPサーバの起動処理です：

```go:cmd/mcp/main.go
func main() {
	// 設定の初期化
	cfg := config.NewConfig()

	// データベースディレクトリを作成
	if err := cfg.EnsureDatabaseDir(); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// 依存関係の初期化
	dependencies, err := initializeDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// MCPサーバの作成
	s := server.NewMCPServer(
		cfg.MCP.Name,
		cfg.MCP.Version,
		server.WithToolCapabilities(false),
	)

	// ツールの登録
	if err := registerAllTools(s, dependencies); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// サーバの起動
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Dependencies はアプリケーションの依存関係を表します
type Dependencies struct {
	CommandHandler *handler.StrengthCommandHandler
	QueryHandler   *query_handler.StrengthQueryHandler
}

// registerAllTools はすべてのツールを登録します
func registerAllTools(s *server.MCPServer, deps *Dependencies) error {
	// トレーニング記録ツール
	trainingTool := tool.NewTrainingToolHandler(deps.CommandHandler)
	if err := trainingTool.Register(s); err != nil {
		return fmt.Errorf("failed to register training tool: %w", err)
	}

	// 期間指定クエリツール
	queryTool := tool.NewQueryToolHandler(deps.QueryHandler)
	if err := queryTool.Register(s); err != nil {
		return fmt.Errorf("failed to register query tool: %w", err)
	}

	// 個人記録ツール
	recordTool := tool.NewRecordToolHandler(deps.QueryHandler)
	if err := recordTool.Register(s); err != nil {
		return fmt.Errorf("failed to register record tool: %w", err)
	}

	return nil
}
```

今回はDIライブラリをまだ入れていないため、DBの初期化処理なども現状含んでいます。サーバをセットアップして起動する処理になります。

#### main.goの特徴

**1. 責務の明確化**
main関数はMCPサーバの起動フローのみに専念し、各処理は専用関数に委譲されています。これにより、MCPサーバを立ち上げる際の全体的な流れが一目で理解できます。

**2. 依存関係の一元管理**
アプリケーション層のハンドラーを`Dependencies`構造体で管理し、MCPツール層とビジネスロジック層の依存関係を明確化しています。これにより、新しいツールを追加する際も、この構造体を通じて必要なハンドラーにアクセスできます。

**3. ツール登録の抽象化**
各ツールは`Register()`メソッドで自分自身をMCPサーバに登録する責務を持っています。main.goは「どのツールを登録するか」だけを知っていれば良く、ツール固有の登録ロジックは各ツールが管理しています。

**4. 既存アプリケーションのMCP化の容易さ**
最も重要なポイントとして、既存のCRUDアプリケーションをMCP化する際は、main.goの構造をコピーして、`Dependencies`と`registerAllTools`を自分のアプリケーション用に書き換えるだけで済むという点があります。ビジネスロジック部分は一切変更する必要がありません。

最終的には、以下の処理で起動されています：
```go
// MCPサーバの作成
s := server.NewMCPServer(cfg.MCP.Name, cfg.MCP.Version, server.WithToolCapabilities(false))

// ツールの登録
if err := registerAllTools(s, dependencies); err != nil {
    log.Fatalf("Failed to register tools: %v", err)
}

// サーバの起動
if err := server.ServeStdio(s); err != nil {
    fmt.Printf("Server error: %v\n", err)
}
```

- `NewMCPServer`でサーバを初期化します。サーバ名（サービス名）とバージョンを設定して初期化を行っています
- `server.ServeStdio`でサーバを起動しています

## ローカルでの検証

ローカルではアプリケーションをビルドしてコマンドで実行可能です。

事前にビルドが必要です：
```shell
❯ go build ./cmd/mcp
```

実行例：
```shell
❯ echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "record_training", "arguments": {"date": "2025-06-13", "exercises": [{"name": "ベンチプレス", "category": "Compound", "sets": [{"weight_kg": 2, "reps": 8, "rest_time_seconds": 180, "rpe": 8}]}], "notes": "ローカルテスト"}}}' | ./mcp 

# 実行結果
2025/06/14 15:12:56 Initialized SQLite repository at: /Users/xxxxxx/develop/mcp/fitness-mcp-server/data/fitness.db
2025/06/14 15:12:56 Initialized SQLite query service at: /Users/xxxxxx/develop/mcp/fitness-mcp-server/data/fitness.db
2025/06/14 15:12:56 Recording training session for date: 2025-06-13
2025/06/14 15:12:56 Successfully recorded training with ID: 230dc699-4570-487b-8f22-65546fc81215
{"jsonrpc":"2.0","id":1,"result":{"content":[{"type":"text","text":"記録完了: TrainingID=230dc699-4570-487b-8f22-65546fc81215, メッセージ=筋トレセッション（1種目、1セット）を記録しました"}]}}
```

## AIツールへの設定

今回はClaudeCodeを使っているので、設定ファイルに以下のように追加します：

```json
"fitness-mcp-server": {
  "command": "docker",
  "args": [
    "run",
    "--rm",
    "-i",
    "-v",
    "/Users/sakemihiroshi/develop/mcp/fitness-mcp-server/data:/app/data",
    "fitness-mcp-server-fitness-mcp"
  ],
  "env": {
    "MCP_DATA_DIR": "/app/data"
  }
}
```

実際に送ってみました。
![](https://storage.googleapis.com/zenn-user-upload/7fd9005f284c-20250614.png)

また、テストデータも含まれますが、取得も正常に行うことができているようです。
![](https://storage.googleapis.com/zenn-user-upload/730ad8757784-20250614.png)

少し、リクエスト内容が特殊だったり、私のチャットが適当すぎるのでAIも困惑していますが、無事に保存できているようです。

## 最後に

今回初めてMCPサーバを立ててみました。公開されているMCPサーバはよく使っていて便利だなと思っていましたし、実際作ってみることができてよかったです。

また、作ってみて感じたのですが、最終的なサーバの設定が変わるだけでロジックは従来のアプリケーションとほとんど同じだったので、既存のアプリケーションの一部の機能をMCPサーバ化するのもありなのかなと思いました。

もう少し機能を追加したり、リモートへのセットアップなども今後やってみようと思います。ありがとうございました。