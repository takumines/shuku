# Shuku

Shukuは、様々な画像フォーマット（JPEG、PNG、WebP）を簡単に圧縮できるCLIツールです。

## 概要

Shukuは、画像ファイルの圧縮を効率的に行うためのコマンドラインツールです。
複数の画像フォーマットに対応し、高品質な圧縮を実現します。

### 特徴

- 複数の画像フォーマットに対応（JPEG、PNG、WebP）
- シンプルなCLIインターフェース
- 高品質な圧縮アルゴリズム
- 並行処理による高速な圧縮
- カスタマイズ可能な圧縮オプション

## インストール

```bash
go install github.com/takumines/shuku/cmd/shuku@latest
```

## 使用方法

### 基本的な使用方法

```bash
# JPEG画像を圧縮
shuku compress -i input.jpg -o output.jpg

# 品質を指定して圧縮（0-100、デフォルト: 80）
shuku compress -i input.jpg -o output.jpg --quality 50

# 詳細情報を表示
shuku compress -i input.jpg -o output.jpg -v
```

## 開発進捗

### ✅ 完了済み機能
- JPEG圧縮機能（品質制御付き）
- PNG圧縮エンジン（内部実装完了）
- CLI基盤（urfave/cli/v2）
- 公開API（ライブラリとしての使用）
- テスト環境とテストデータ
- 開発ガイドライン（CLAUDE.md）

### 🚧 進行中・TODO
- [ ] PNG圧縮サポートをCLIに追加（高優先度）
- [ ] WebP圧縮機能を実装（中優先度）
- [ ] エラーメッセージを英語化（中優先度）
- [ ] 複数ファイルの一括圧縮機能を追加（低優先度）
- [ ] 使用方法ドキュメントの詳細化（低優先度）

### 📝 技術メモ
- 現在JPEGのみCLI対応、PNGは内部実装済みだがCLI制限中
- デュアルユース設計（CLI + ライブラリ）
- インターフェースベースアーキテクチャで拡張性確保

## 開発

### 必要条件

- Go 1.22以上

### セットアップ

```bash
git clone https://github.com/takumines/shuku.git
cd shuku
go mod download
```

### ビルド

```bash
make build
```

### テスト

```bash
make test
```
