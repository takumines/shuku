# Shuku

画像ファイルを効率的に圧縮するCLIツールとGoライブラリです。

### 特徴

- 画像ファイルの効率的な圧縮
- シンプルなCLIインターフェース
- 品質制御可能な圧縮
- ライブラリとしても利用可能

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

## 現在の状況

現在はJPEG形式のみ対応しています。PNG、WebP対応は開発中です。

## 開発

### 必要条件

- Go 1.22以上

詳細な開発情報は[CLAUDE.md](./CLAUDE.md)を参照してください。
