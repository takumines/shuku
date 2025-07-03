// Package shuku は画像圧縮のためのAPIを提供します。
package shuku

import (
	"errors"
	"image"
	"os"
	"path/filepath"
	"strings"

	"shuku/internal/compressor"
)

// 画像形式に対応するコンプレッサーのマップ
var compressors = map[string]compressor.Compressor{}

// コンプレッサーの登録
func init() {
	// JPEGコンプレッサーを登録
	jpegCompressor := compressor.NewJPEGCompressor()
	compressors[jpegCompressor.SupportedFormat()] = jpegCompressor
	compressors["jpg"] = jpegCompressor // jpgも同じコンプレッサーで対応

	// PNGコンプレッサーを登録
	pngCompressor := compressor.NewPNGCompressor()
	compressors[pngCompressor.SupportedFormat()] = pngCompressor
}

// Compress はバイトスライスとして提供された画像データを圧縮します。
// 画像形式は入力データから自動検出されます。
func Compress(data []byte, options Options) ([]byte, error) {
	// 画像形式の検出
	format, err := detectImageFormat(data)
	if err != nil {
		return nil, err
	}

	// 対応するコンプレッサーを取得
	comp, ok := compressors[format]
	if !ok {
		return nil, errors.New("サポートされていない画像形式です: " + format)
	}

	// 内部オプションに変換
	internalOpts := compressor.Options{
		Quality:     options.Quality,
		PaletteSize: options.PaletteSize,
	}

	// 圧縮を実行
	return comp.CompressBytes(data, internalOpts)
}

// CompressImage は画像インターフェースを圧縮します。
// 画像形式はSupportedFormat()から取得されます。
func CompressImage(img image.Image, format string, options Options) (image.Image, error) {
	// 形式を小文字に変換
	format = strings.ToLower(format)

	// 対応するコンプレッサーを取得
	comp, ok := compressors[format]
	if !ok {
		return nil, errors.New("サポートされていない画像形式です: " + format)
	}

	// 内部オプションに変換
	internalOpts := compressor.Options{
		Quality:     options.Quality,
		PaletteSize: options.PaletteSize,
	}

	// 圧縮を実行
	return comp.Compress(img, internalOpts)
}

// CompressFile はファイルパスを指定して画像ファイルを圧縮します。
// 出力ファイルが指定されていない場合は、入力ファイルの名前に "_compressed" を追加します。
func CompressFile(inputPath, outputPath string, options Options) error {
	// 入力ファイルを開く
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// 出力パスが指定されていない場合は、デフォルトのパスを生成
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		baseName := strings.TrimSuffix(inputPath, ext)
		outputPath = baseName + "_compressed" + ext
	}

	// 出力ファイルを作成
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 画像形式を拡張子から取得
	format := strings.ToLower(strings.TrimPrefix(filepath.Ext(inputPath), "."))
	if format == "" {
		return errors.New("ファイル拡張子から画像形式を判別できません")
	}

	// 対応するコンプレッサーを取得
	comp, ok := compressors[format]
	if !ok {
		return errors.New("サポートされていない画像形式です: " + format)
	}

	// 内部オプションに変換
	internalOpts := compressor.Options{
		Quality:     options.Quality,
		PaletteSize: options.PaletteSize,
	}

	// 圧縮を実行
	return comp.CompressReader(inputFile, outputFile, internalOpts)
}

// 画像データからフォーマットを検出する関数
func detectImageFormat(data []byte) (string, error) {
	// JPEGのシグネチャを確認
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpeg", nil
	}

	// PNGのシグネチャを確認
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return "png", nil
	}

	// WebPのシグネチャを確認（RIFFとWEBP）
	if len(data) >= 12 && data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return "webp", nil
	}

	return "", errors.New("サポートされていない、または認識できない画像形式です")
}
