# MCPツール API設計

## 1. record_running ツール仕様

### 基本情報
```json
{
  "name": "record_running",
  "description": "ランニングセッションを記録するツール",
  "input_schema": {
    "type": "object",
    "properties": {
      "date": {
        "type": "string",
        "description": "ランニング実施日付（YYYY-MM-DD形式）",
        "pattern": "^\\d{4}-\\d{2}-\\d{2}$"
      },
      "distance_km": {
        "type": "number",
        "description": "走行距離（km単位）",
        "minimum": 0.01
      },
      "duration": {
        "type": ["string", "number"],
        "description": "走行時間（MM:SS形式の文字列または分数の数値）"
      },
      "run_type": {
        "type": "string",
        "description": "ランニングタイプ",
        "enum": ["Easy", "Tempo", "Interval", "Long", "Race"]
      },
      "heart_rate_bpm": {
        "type": "integer",
        "description": "平均心拍数（オプション）",
        "minimum": 1,
        "maximum": 300
      },
      "notes": {
        "type": "string",
        "description": "メモや備考（オプション）"
      }
    },
    "required": ["date", "distance_km", "duration", "run_type"]
  }
}
```

## 2. 入力パラメータ詳細

### date (必須)
- **形式**: `YYYY-MM-DD`
- **例**: `"2025-06-16"`
- **バリデーション**: 正しい日付形式、未来日は警告（記録は可能）

### distance_km (必須)
- **型**: number
- **単位**: キロメートル
- **例**: `5.0`, `10.5`, `0.8`
- **最小値**: 0.01km（10m）
- **バリデーション**: 正の数値

### duration (必須)
- **型**: string または number
- **形式1**: `"MM:SS"`（文字列）
  - 例: `"25:30"` (25分30秒)
  - 例: `"1:02:15"` (1時間2分15秒)
- **形式2**: 分数（数値）
  - 例: `25.5` (25分30秒)
  - 例: `62.25` (1時間2分15秒)

### run_type (必須)
- **型**: string
- **値**: 
  - `"Easy"`: イージーラン（ジョグ）
  - `"Tempo"`: テンポラン（閾値走）
  - `"Interval"`: インターバル走
  - `"Long"`: ロングラン（LSD）
  - `"Race"`: レース・タイムトライアル

### heart_rate_bpm (オプション)
- **型**: integer
- **単位**: bpm（beats per minute）
- **例**: `165`, `180`
- **範囲**: 1-300bpm
- **説明**: 平均心拍数

### notes (オプション)
- **型**: string
- **例**: `"気持ちよく走れた"`, `"膝の調子が良い"`
- **制限**: 特になし（文字数制限なし）

## 3. 使用例

### 基本的な使用例

```json
{
  "name": "record_running",
  "arguments": {
    "date": "2025-06-16",
    "distance_km": 5.0,
    "duration": "25:30",
    "run_type": "Easy",
    "notes": "朝ランで気持ちよく走れた"
  }
}
```

### 心拍数付きの例

```json
{
  "name": "record_running",
  "arguments": {
    "date": "2025-06-17",
    "distance_km": 8.0,
    "duration": 35.2,
    "run_type": "Tempo",
    "heart_rate_bpm": 175,
    "notes": "息が上がったが最後まで維持できた"
  }
}
```

### 長時間ランの例

```json
{
  "name": "record_running",
  "arguments": {
    "date": "2025-06-18",
    "distance_km": 21.1,
    "duration": "1:45:30",
    "run_type": "Long",
    "heart_rate_bpm": 155,
    "notes": "ハーフマラソンの距離。後半ペースが落ちた"
  }
}
```

## 4. レスポンス仕様

### 成功時レスポンス

```json
{
  "success": true,
  "session_id": "550e8400-e29b-41d4-a716-446655440001",
  "message": "ランニングセッションを記録しました",
  "details": {
    "date": "2025-06-16",
    "distance_km": 5.0,
    "duration_seconds": 1530,
    "pace": "5:06/km",
    "run_type": "Easy",
    "heart_rate_bpm": 165,
    "notes": "朝ランで気持ちよく走れた"
  }
}
```

### エラー時レスポンス

```json
{
  "success": false,
  "error_code": "INVALID_DURATION",
  "message": "時間の形式が不正です: '25:ab'（MM:SS形式または分数で入力してください）",
  "details": {
    "parameter": "duration",
    "value": "25:ab",
    "expected_formats": ["MM:SS", "number (minutes)"]
  }
}
```

## 5. エラーハンドリング

### エラーコード一覧

| エラーコード | 説明 | 例 |
|-------------|------|---|
| `INVALID_DATE` | 日付形式が不正 | `"2025-13-01"` |
| `INVALID_DISTANCE` | 距離が不正 | `-5.0`, `0` |
| `INVALID_DURATION` | 時間形式が不正 | `"25:ab"`, `-10` |
| `INVALID_RUN_TYPE` | ランニングタイプが不正 | `"Jog"`, `"Sprint"` |
| `INVALID_HEART_RATE` | 心拍数が不正 | `-10`, `400` |
| `DUPLICATE_ENTRY` | 同日同時刻の重複 | 同じ日時の記録が既に存在 |
| `DATABASE_ERROR` | データベースエラー | 接続失敗、制約違反等 |

### バリデーション詳細

#### 時間パースロジック

```go
func ParseDuration(input interface{}) (time.Duration, error) {
    switch v := input.(type) {
    case string:
        // MM:SS または H:MM:SS 形式
        if strings.Contains(v, ":") {
            return parseTimeString(v)
        }
        // 文字列の数値を分数として扱う
        minutes, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return 0, fmt.Errorf("時間の形式が不正です: '%s'", v)
        }
        return time.Duration(minutes * float64(time.Minute)), nil
    case float64:
        // 分数として扱う
        return time.Duration(v * float64(time.Minute)), nil
    default:
        return 0, fmt.Errorf("時間は文字列または数値で指定してください")
    }
}

func parseTimeString(timeStr string) (time.Duration, error) {
    parts := strings.Split(timeStr, ":")
    switch len(parts) {
    case 2: // MM:SS
        minutes, err := strconv.Atoi(parts[0])
        if err != nil {
            return 0, fmt.Errorf("分が不正です: %s", parts[0])
        }
        seconds, err := strconv.Atoi(parts[1])
        if err != nil {
            return 0, fmt.Errorf("秒が不正です: %s", parts[1])
        }
        return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
    case 3: // H:MM:SS
        hours, err := strconv.Atoi(parts[0])
        if err != nil {
            return 0, fmt.Errorf("時間が不正です: %s", parts[0])
        }
        minutes, err := strconv.Atoi(parts[1])
        if err != nil {
            return 0, fmt.Errorf("分が不正です: %s", parts[1])
        }
        seconds, err := strconv.Atoi(parts[2])
        if err != nil {
            return 0, fmt.Errorf("秒が不正です: %s", parts[2])
        }
        return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
    default:
        return 0, fmt.Errorf("時間の形式が不正です: '%s'（MM:SS または H:MM:SS 形式で入力してください）", timeStr)
    }
}
```

## 6. ペース計算ロジック

### 自動計算仕様

```go
type PaceCalculator struct{}

func (pc *PaceCalculator) CalculatePace(distance float64, duration time.Duration) (float64, error) {
    if distance <= 0 {
        return 0, fmt.Errorf("距離は正の値である必要があります")
    }
    
    // 秒/km で計算
    secondsPerKm := duration.Seconds() / distance
    return secondsPerKm, nil
}

func (pc *PaceCalculator) FormatPace(secondsPerKm float64) string {
    totalSeconds := int(secondsPerKm)
    minutes := totalSeconds / 60
    seconds := totalSeconds % 60
    return fmt.Sprintf("%d:%02d/km", minutes, seconds)
}
```

### ペース表示例

| 距離 | 時間 | ペース |
|------|------|-------|
| 5.0km | 25:30 | 5:06/km |
| 10.0km | 50:00 | 5:00/km |
| 3.0km | 15:45 | 5:15/km |

## 7. 将来拡張予定

### 予定される追加パラメータ

```json
{
  "weather": "晴れ|曇り|雨|雪",
  "temperature_celsius": 25,
  "humidity_percent": 60,
  "route_name": "皇居ラン",
  "elevation_gain_m": 50,
  "cadence_spm": 180,
  "max_heart_rate_bpm": 185
}
```

### 予定されるクエリツール

- `get_running_records`: 記録検索・一覧取得
- `get_running_stats`: 統計情報取得
- `get_running_pace_analysis`: ペース分析
- `get_running_progress`: 進捗分析

---

この設計により、直感的で使いやすいランニング記録APIを提供できます。
