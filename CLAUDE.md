# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

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
  ├── compress/     - 単一ファイル圧縮コマンド
  ├── batch/        - バッチ処理コマンド
  └── version/      - バージョン表示コマンド
pkg/shuku/          - パブリックライブラリAPI (CompressFile, Compress, CompressImage)
internal/compressor/ - 画像圧縮実装 (JPEG, PNG, WebP)
internal/batch/     - バッチ処理・並行処理実装
testdata/          - テスト画像と生成ツール
tests/             - 統合テスト
```

### 主要コンポーネント
- `pkg/shuku/shuku.go`: 形式自動検出機能付きのメインパブリックAPI
- `internal/compressor/interface.go`: 共通`Compressor`インターフェース
- `internal/compressor/*_compressor.go`: 形式固有の実装
- `internal/batch/processor.go`: バッチ処理・並行処理エンジン
- `cmd/shuku/compress/compress.go`: CLI単一ファイル圧縮コマンド
- `cmd/shuku/batch/batch.go`: CLIバッチ処理コマンド

## 開発コマンド

### ビルドと実行
```bash
# CLIツールをビルド
go build -o shuku cmd/shuku/main.go

# ビルドせずに実行
go run cmd/shuku/main.go compress -i input.jpg -o output.jpg

# CLI機能をテスト（単一ファイル）
./cmd/shuku/shuku compress --input testdata/test_image.jpg --output compressed.jpg

# CLI機能をテスト（バッチ処理）
./shuku batch -i testdata -o compressed_batch -v --stats
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

### 重要な開発ルール
- **コミット前必須**: `gofmt -w .` を必ず実行してからコミットする
- **Claude Codeへの指示**: コミットする前には必ずgofmtを実行すること

## 現在の実装状況

### 完了済み機能
- 品質制御付きJPEG圧縮（CLI・ライブラリ完全対応）
- PNG圧縮（CLI・ライブラリ完全対応）
- WebP圧縮（CLI・ライブラリ完全対応）
- **バッチ処理機能（複数ファイル一括圧縮）**
- **並行処理によるパフォーマンス最適化**
- **再帰的ディレクトリ処理**
- **フィルタリング機能（包含・除外パターン）**
- **圧縮統計表示機能**
- urfave/cli/v2を使用したCLIインターフェース
- ライブラリ使用のためのパブリックAPI
- 圧縮設定用のOptions構造体
- 画像形式自動検出機能（バイナリシグネチャベース）
- テストデータ生成ツール
- 包括的テストスイート（95.8%カバレッジ）

### 既知の制限事項
- 一部パッケージのテストカバレッジが不足（cmd/shuku: 0%, version: 0%）

### 新しい画像形式サポートの追加方法
1. `Compressor`インターフェースを実装する`internal/compressor/<format>_compressor.go`を作成
2. `pkg/shuku/shuku.go`のinit()関数でコンプレッサーを登録
3. `detectImageFormat()`関数に形式検出ロジックを追加
4. `cmd/shuku/compress/compress.go`のsupportedFormats配列に拡張子を追加

## コードスタイルガイドライン

- Go標準規約に従い`gofmt`を使用
- 全パッケージにパッケージレベルドキュメントが必要
- エクスポートされた関数と型にはドキュメントコメントが必須
- 拡張性のためインターフェースベース設計を使用
- 外部依存関係を最小化（現在はurfave/cli/v2のみ）
- エラーメッセージは日本語で統一済み
- 包括的カバレッジのためテーブル駆動テストを使用

## テスト戦略

- 実装と同じパッケージ内にユニットテスト
- `tests/`ディレクトリに統合テスト
- パフォーマンス重要な圧縮コードにベンチマークテスト
- `testdata/`ディレクトリにテストデータ
- テスト画像作成には`testdata/generate_test_image.go`を使用

## 開発ロードマップ

開発進捗とマイルストーン管理については、プロジェクトルートの`TODO.md`を参照してください。
