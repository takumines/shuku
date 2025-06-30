# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

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

### 🎯 プロダクトゴールと戦略

**最終目標**: 実用的な画像圧縮ツールとして広く採用される  
**戦略**: 既存実装活用 → 品質基盤 → 新価値 → UX向上 → 完成品質

### 📋 マイルストーン別開発計画

#### 🎯 Milestone 1: PNG対応完了 (v0.2.0) - 基本機能の完成
*既存PNG実装を活用した最速価値提供*

- [ ] **m1_1**: PNG圧縮をCLIに統合（compress.goの形式制限削除）
- [ ] **m1_2**: PNGコンプレッサーをshuku.goのinit()に登録
- [ ] **m1_3**: PNG圧縮の統合テスト追加（CLI動作確認）
- [ ] **m1_4**: エラーメッセージ英語統一（国際化対応）

**期待成果**: 実用的な画像圧縮ツールとして機能（JPEG+PNG対応）

#### 🎯 Milestone 2: 品質基盤構築 (v0.3.0前) - 開発効率と信頼性
*長期開発効率と信頼性の確保*

- [ ] **m2_1**: GitHub Actions CI/CD設定（テスト自動化）
- [ ] **m2_2**: GoReleaser設定（バイナリ自動配布）
- [ ] **m2_3**: PNG圧縮のユニットテスト強化

**期待成果**: 安定した開発サイクル確立

#### 🎯 Milestone 3: WebP対応 (v0.3.0) - 新価値提供
*WebP対応で差別化*

- [ ] **m3_1**: WebP圧縮エンジン実装（golang.org/x/image/webp使用）
- [ ] **m3_2**: WebPコンプレッサー登録とCLI統合
- [ ] **m3_3**: WebP圧縮テスト実装

**期待成果**: モダンな画像形式対応ツール（主要3形式対応）

#### 🎯 Milestone 4: UX大幅改善 (v0.4.0) - 実用性向上
*大幅なUX改善*

- [ ] **m4_1**: 複数ファイル一括圧縮（ディレクトリ指定対応）
- [ ] **m4_2**: プログレスバー実装（大容量ファイル処理UX）
- [ ] **m4_3**: 設定ファイル対応（.shuku.yaml）

**期待成果**: 企業・プロジェクト採用レベルの実用性

#### 🎯 Milestone 5: Production Ready (v1.0.0) - 完成品質
*完成品質の提供*

- [ ] **m5_1**: 包括的ドキュメント整備
- [ ] **m5_2**: ベンチマークテスト・パフォーマンス最適化
- [ ] **m5_3**: 詳細ログ機能実装

**期待成果**: エンタープライズ品質

### ⚡ 効率化戦略

#### 🔄 並行作業戦略
- PNG統合中にCI/CD設定を並行実施
- 国際化は隙間時間で対応
- ドキュメント整備は機能開発と並行

#### 📦 既存資産最大活用
- PNG内部実装済み→CLI統合のみ
- JPEG実装パターン→WebP実装に応用
- テスト基盤→新機能テストに流用

#### 🎪 段階的価値提供
- 各マイルストーンでユーザー価値を提供
- 早期フィードバック獲得で品質向上

### 🚀 次のアクション

**今すぐ開始**: Milestone 1の m1_1 から実装開始  
→ PNG圧縮CLI統合（compress.goの形式制限削除）