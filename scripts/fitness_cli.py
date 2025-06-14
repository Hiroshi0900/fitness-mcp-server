#!/usr/bin/env python3
"""
Fitness MCP Server CLI Tool
MCPサーバーとやり取りするためのコマンドラインツール
"""

import json
import subprocess
import sys
import argparse
from datetime import datetime, date
from typing import Dict, Any, List


class FitnessCLI:
    def __init__(self, server_path: str = "./mcp"):
        self.server_path = server_path
        
    def _call_server(self, method: str, params: Dict[str, Any]) -> Dict[str, Any]:
        """MCPサーバーにJSON-RPCリクエストを送信"""
        request = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": method,
            "params": params
        }
        
        try:
            process = subprocess.Popen(
                [self.server_path],
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate(input=json.dumps(request))
            
            if process.returncode != 0:
                print(f"Error: {stderr}", file=sys.stderr)
                return {"error": "Server error"}
                
            return json.loads(stdout)
            
        except Exception as e:
            print(f"Error calling server: {e}", file=sys.stderr)
            return {"error": str(e)}
    
    def record_training(self, date_str: str, exercises: List[Dict], notes: str = ""):
        """トレーニングを記録"""
        params = {
            "name": "record_training",
            "arguments": {
                "date": date_str,
                "exercises": exercises,
                "notes": notes
            }
        }
        
        result = self._call_server("tools/call", params)
        if "error" in result:
            print(f"❌ エラー: {result['error']}")
        else:
            print("✅ トレーニング記録完了!")
            if "result" in result and "content" in result["result"]:
                for content in result["result"]["content"]:
                    if content["type"] == "text":
                        print(content["text"])
    
    def get_trainings(self, start_date: str, end_date: str):
        """期間指定でトレーニング履歴を取得"""
        params = {
            "name": "get_trainings_by_date_range",
            "arguments": {
                "start_date": start_date,
                "end_date": end_date
            }
        }
        
        result = self._call_server("tools/call", params)
        if "error" in result:
            print(f"❌ エラー: {result['error']}")
        else:
            if "result" in result and "content" in result["result"]:
                for content in result["result"]["content"]:
                    if content["type"] == "text":
                        print(content["text"])
    
    def get_personal_records(self, exercise_name: str = None):
        """個人記録を取得"""
        arguments = {}
        if exercise_name:
            arguments["exercise_name"] = exercise_name
            
        params = {
            "name": "get_personal_records",
            "arguments": arguments
        }
        
        result = self._call_server("tools/call", params)
        if "error" in result:
            print(f"❌ エラー: {result['error']}")
        else:
            if "result" in result and "content" in result["result"]:
                for content in result["result"]["content"]:
                    if content["type"] == "text":
                        print(content["text"])


def main():
    parser = argparse.ArgumentParser(description="Fitness MCP Server CLI Tool")
    parser.add_argument("--server", default="./mcp", help="MCPサーバーのパス")
    
    subparsers = parser.add_subparsers(dest="command", help="利用可能なコマンド")
    
    # record コマンド
    record_parser = subparsers.add_parser("record", help="トレーニングを記録")
    record_parser.add_argument("--date", default=date.today().isoformat(), help="日付 (YYYY-MM-DD)")
    record_parser.add_argument("--notes", default="", help="メモ")
    
    # history コマンド
    history_parser = subparsers.add_parser("history", help="トレーニング履歴を表示")
    history_parser.add_argument("--start", required=True, help="開始日 (YYYY-MM-DD)")
    history_parser.add_argument("--end", required=True, help="終了日 (YYYY-MM-DD)")
    
    # records コマンド
    records_parser = subparsers.add_parser("records", help="個人記録を表示")
    records_parser.add_argument("--exercise", help="特定のエクササイズ名")
    
    # quick コマンド（クイック記録用）
    quick_parser = subparsers.add_parser("quick", help="クイック記録（簡単な入力）")
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        return
    
    cli = FitnessCLI(args.server)
    
    if args.command == "record":
        # インタラクティブな記録入力
        print("🏋️ トレーニング記録を開始します")
        exercises = []
        
        while True:
            print("\n--- エクササイズ入力 ---")
            name = input("エクササイズ名: ").strip()
            if not name:
                break
                
            print("カテゴリを選択:")
            print("1. Compound (複合種目)")
            print("2. Isolation (単関節種目)")
            print("3. Cardio (有酸素運動)")
            
            category_choice = input("選択 (1-3): ").strip()
            category_map = {"1": "Compound", "2": "Isolation", "3": "Cardio"}
            category = category_map.get(category_choice, "Compound")
            
            sets = []
            set_num = 1
            while True:
                print(f"\n--- セット {set_num} ---")
                weight = input("重量 (kg): ").strip()
                if not weight:
                    break
                    
                reps = input("回数: ").strip()
                rest_time = input("休憩時間 (秒, デフォルト60): ").strip() or "60"
                rpe = input("RPE (1-10, 省略可): ").strip()
                
                set_data = {
                    "weight_kg": float(weight),
                    "reps": int(reps),
                    "rest_time_seconds": int(rest_time)
                }
                
                if rpe:
                    set_data["rpe"] = int(rpe)
                    
                sets.append(set_data)
                set_num += 1
                
                if input("次のセットを追加? (y/N): ").lower() != 'y':
                    break
            
            if sets:
                exercises.append({
                    "name": name,
                    "category": category,
                    "sets": sets
                })
            
            if input("次のエクササイズを追加? (y/N): ").lower() != 'y':
                break
        
        if exercises:
            cli.record_training(args.date, exercises, args.notes)
        else:
            print("エクササイズが入力されませんでした。")
    
    elif args.command == "history":
        cli.get_trainings(args.start, args.end)
    
    elif args.command == "records":
        cli.get_personal_records(args.exercise)
    
    elif args.command == "quick":
        # クイック記録のサンプル
        print("🚀 クイック記録例を実行...")
        sample_exercises = [
            {
                "name": "ベンチプレス",
                "category": "Compound",
                "sets": [
                    {"weight_kg": 80.0, "reps": 10, "rest_time_seconds": 120, "rpe": 7}
                ]
            }
        ]
        cli.record_training(date.today().isoformat(), sample_exercises, "クイック記録")


if __name__ == "__main__":
    main()
