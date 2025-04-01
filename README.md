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

WIP...

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
