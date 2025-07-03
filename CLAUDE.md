# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。
ファイルの末尾には必ず改行を入れてください。

## Claude Codeの動作設定

- **言語設定**: ユーザーとのすべてのやり取りは日本語で行う
- **コミュニケーション**: 技術的な説明、エラーメッセージ、進捗報告はすべて日本語で提供する

## プロジェクト概要

Shukuは様々な画像形式（JPEG、PNG、WebP）を圧縮するCLIツールとGoライブラリです。デュアルユース設計パターンに従い、コマンドラインツールとして使用することも、ライブラリとしてインポートすることも可能です。

## アーキテクチャ

### 核となる設計原則
- **デュアルユース設計**: CLIツールとインポート可能ライブラリの両方として機能
- **インターフェースベースアーキテクチャ**: 異なる画像形式に対して`Compressor`インターフェースを使用
- **依存性注入**: 拡張性のためにコンプレッサーをマップに登録
- **明確な分離**: CLI関連は`cmd/`、パブリックAPIは`pkg/`、内部実装は`internal/`

### パッケージ構造
```
cmd/shuku/          - CLIアプリケーションのエントリーポイントとコマンド
pkg/shuku/          - パブリックライブラリAPI (CompressFile, Compress, CompressImage)
internal/compressor/ - 画像圧縮実装 (JPEG, PNG, WebP)
internal/optimizer/ - 最適化アルゴリズム（計画中）
internal/worker/    - 並行処理ユーティリティ（計画中）
testdata/          - テスト画像と生成ツール
```

### 主要コンポーネント
- `pkg/shuku/shuku.go`: 形式自動検出機能付きのメインパブリックAPI
- `internal/compressor/interface.go`: 共通`Compressor`インターフェース
- `internal/compressor/*_compressor.go`: 形式固有の実装
- `cmd/shuku/compress/compress.go`: CLI圧縮コマンド実装

## 開発コマンド

### ビルドと実行
```bash
# CLIツールをビルド
go build -o shuku cmd/shuku/main.go

# ビルドせずに実行
go run cmd/shuku/main.go compress -i input.jpg -o output.jpg

# CLI機能をテスト
./cmd/shuku/shuku compress --input testdata/test_image.jpg --output compressed.jpg
```

### テスト
```bash
# 全テストを実行
go test ./...

# 特定パッケージのテストを実行
go test ./internal/compressor
go test ./pkg/shuku

# 詳細出力でテストを実行
go test -v ./...

# ベンチマークを実行
go test -bench=. ./internal/compressor
```

### コード品質
```bash
# コードフォーマット（コミット前に必ず実行）
gofmt -w .

# 一般的な問題をチェック
go vet ./...

# テストとvetを一緒に実行
go test -vet=all ./...
```

## 現在の実装状況

### 完了済み機能
- 品質制御付きJPEG圧縮
- PNG圧縮（実装済みだがCLIサポートは限定的）
- urfave/cli/v2を使用したCLIインターフェース
- ライブラリ使用のためのパブリックAPI
- 圧縮設定用のOptions構造体
- テストデータ生成ツール

### 既知の制限事項
- CLIは現在JPEGファイルのみサポート（`cmd/shuku/compress/compress.go`の87-88行目）
- PNGコンプレッサーは実装済みだがCLIに統合されていない
- WebPサポートはまだ実装されていない
- エラーメッセージが日本語と英語が混在

### 新しい画像形式サポートの追加方法
1. `Compressor`インターフェースを実装する`internal/compressor/<format>_compressor.go`を作成
2. `pkg/shuku/shuku.go`のinit()関数でコンプレッサーを登録
3. `detectImageFormat()`関数に形式検出ロジックを追加
4. `cmd/shuku/compress/compress.go`のCLI検証を更新

## コードスタイルガイドライン

- Go標準規約に従い`gofmt`を使用
- 全パッケージにパッケージレベルドキュメントが必要
- エクスポートされた関数と型にはドキュメントコメントが必須
- 拡張性のためインターフェースベース設計を使用
- 外部依存関係を最小化（現在はurfave/cli/v2のみ）
- 一貫性のためエラーメッセージは英語で統一
- 包括的カバレッジのためテーブル駆動テストを使用

## テスト戦略

- 実装と同じパッケージ内にユニットテスト
- `tests/`ディレクトリに統合テスト
- パフォーマンス重要な圧縮コードにベンチマークテスト
- `testdata/`ディレクトリにテストデータ
- テスト画像作成には`testdata/generate_test_image.go`を使用

## 開発ロードマップ

開発進捗とマイルストーン管理については、プロジェクトルートの`TODO.md`を参照してください。