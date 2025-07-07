# Shuku

画像ファイルを効率的に圧縮するCLIツールとGoライブラリです。

### 特徴

- 画像ファイルの効率的な圧縮
- シンプルなCLIインターフェース
- 品質制御可能な圧縮
- ライブラリとしても利用可能

## インストール

### プリビルドバイナリを使用（推奨）

[GitHub Releases](https://github.com/takumines/shuku/releases)から最新のバイナリをダウンロードしてください。

```bash
# Linux (amd64)
wget https://github.com/takumines/shuku/releases/latest/download/shuku_*_linux_amd64.tar.gz
tar -xzf shuku_*_linux_amd64.tar.gz
sudo mv shuku /usr/local/bin/

# macOS (amd64)
wget https://github.com/takumines/shuku/releases/latest/download/shuku_*_darwin_amd64.tar.gz
tar -xzf shuku_*_darwin_amd64.tar.gz
sudo mv shuku /usr/local/bin/

# Windows
# GitHubリリースページからshuku_*_windows_amd64.zipをダウンロードして展開
```

### Go installを使用

```bash
go install github.com/takumines/shuku/cmd/shuku@latest
```

## 使用方法

### 基本的な使用方法

```bash
# JPEG画像を圧縮
shuku compress -i input.jpg -o output.jpg

# PNG画像を圧縮
shuku compress -i input.png -o output.png

# 品質を指定して圧縮（0-100、デフォルト: 80）
shuku compress -i input.jpg -o output.jpg --quality 50

# 詳細情報を表示
shuku compress -i input.jpg -o output.jpg -v

# バージョン情報を表示
shuku version
```

## 対応形式

- ✅ JPEG (.jpg, .jpeg)
- ✅ PNG (.png)
- 🚧 WebP (.webp) - 開発中

## 開発

### 必要条件

- Go 1.22以上

詳細な開発情報は[CLAUDE.md](./CLAUDE.md)を参照してください。
