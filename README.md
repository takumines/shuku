# Shuku 📸

**画像ファイルを簡単に圧縮するCLIツール**

ファイルサイズを削減しながら画質を保持。JPEG、PNG、WebPの主要3形式をサポート。

## 🚀 クイックスタート

### 1. インストール

#### プリビルドバイナリ（推奨）
[GitHub Releases](https://github.com/takumines/shuku/releases)から最新版をダウンロード。

**macOS / Linux:**
```bash
curl -s https://api.github.com/repos/takumines/shuku/releases/latest | grep "browser_download_url" | grep "$(uname | tr '[:upper:]' '[:lower:]')" | cut -d '"' -f 4 | xargs curl -L -o shuku.tar.gz
tar -xzf shuku.tar.gz
sudo mv shuku /usr/local/bin/
```

**Windows:**
[リリースページ](https://github.com/takumines/shuku/releases)から`shuku_*_windows_amd64.zip`をダウンロードして展開。

#### Go install
```bash
go install github.com/takumines/shuku/cmd/shuku@latest
```

### 2. 基本的な使い方

```bash
# 最もシンプルな使用法
shuku compress -i 元の画像.jpg -o 圧縮後.jpg

# 品質を指定（1-100、数値が小さいほど高圧縮）
shuku compress -i 元の画像.jpg -o 圧縮後.jpg -q 60

# 詳細情報を表示
shuku compress -i 元の画像.jpg -o 圧縮後.jpg -v
```

## 📖 使用方法

### 基本コマンド

```bash
shuku compress -i <入力ファイル> -o <出力ファイル> [オプション]
```

| オプション | 短縮形 | 説明 | デフォルト |
|-----------|--------|------|----------|
| `--input` | `-i` | 入力ファイルパス（必須） | - |
| `--output` | `-o` | 出力ファイルパス | 元ファイル名_compressed |
| `--quality` | `-q` | 圧縮品質（1-100） | 80 |
| `--verbose` | `-v` | 詳細情報を表示 | false |

### 実用的な例

#### 1. 基本的な圧縮
```bash
# JPEG画像を圧縮
shuku compress -i photo.jpg -o compressed.jpg

# PNG画像を圧縮
shuku compress -i image.png -o compressed.png

# WebP画像を圧縮
shuku compress -i image.webp -o compressed.webp
```

#### 2. 圧縮レベルの調整
```bash
# 高品質（ファイルサイズ大）
shuku compress -i photo.jpg -o high_quality.jpg -q 90

# 標準品質（バランス）
shuku compress -i photo.jpg -o standard.jpg -q 70

# 高圧縮（ファイルサイズ小）
shuku compress -i photo.jpg -o small_size.jpg -q 50
```

#### 3. 出力先を指定しない場合
```bash
# 自動的に "photo_compressed.jpg" が作成される
shuku compress -i photo.jpg

# 圧縮率を確認
shuku compress -i photo.jpg -v
```

#### 4. 複数の画像を処理
```bash
# 複数ファイルの処理例
for file in *.jpg; do
  shuku compress -i "$file" -o "compressed_$file" -q 70
done
```

## 📊 対応形式

| 形式 | 拡張子 | 圧縮設定 | 用途 |
|-----|--------|----------|------|
| **JPEG** | `.jpg`, `.jpeg` | 品質 1-100 | 写真に最適 |
| **PNG** | `.png` | パレットサイズ | 透明度が必要な画像 |
| **WebP** | `.webp` | 品質 1-100 | 最新のWeb標準 |

## 💡 Tips

### 品質設定の目安
- **90-100**: 最高品質（ファイルサイズ大）
- **70-89**: 高品質（推奨）
- **50-69**: 標準品質（Web用）
- **30-49**: 低品質（サムネイル用）

### パフォーマンス最適化
```bash
# 詳細情報で圧縮効果を確認
shuku compress -i large_image.jpg -o compressed.jpg -q 70 -v
```

出力例：
```
入力ファイル: large_image.jpg
出力ファイル: compressed.jpg
圧縮品質: 70
元のサイズ: 2048000 バイト
圧縮後のサイズ: 1024000 バイト
圧縮率: 50.00%
```

## 🔧 トラブルシューティング

### よくある問題

**Q: 「サポートされていない画像形式です」エラー**
```bash
A: 対応形式（JPEG、PNG、WebP）を確認してください
```

**Q: 出力ファイルが作成されない**
```bash
A: 出力ディレクトリの書き込み権限を確認してください
```

**Q: 圧縮後のファイルが大きくなった**
```bash
A: 既に最適化された画像や、品質設定が高すぎる可能性があります
```

### ヘルプコマンド
```bash
# 基本的なヘルプ
shuku --help

# compressコマンドの詳細
shuku compress --help

# バージョン確認
shuku version
```

## 📚 ライブラリとして使用

Goプロジェクトでライブラリとして使用することも可能です：

```go
package main

import (
    "fmt"
    "github.com/takumines/shuku/pkg/shuku"
)

func main() {
    options := shuku.Options{
        Quality: 70,
    }
    
    err := shuku.CompressFile("input.jpg", "output.jpg", options)
    if err != nil {
        fmt.Printf("圧縮エラー: %v\n", err)
        return
    }
    
    fmt.Println("圧縮が完了しました！")
}
```

## 🛠️ 開発

### 必要条件
- Go 1.22以上

### ビルド方法
```bash
git clone https://github.com/takumines/shuku.git
cd shuku
go build -o shuku cmd/shuku/main.go cmd/shuku/root.go
```

### テスト実行
```bash
go test ./...
```

## 📄 ライセンス

MIT License
