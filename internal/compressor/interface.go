// Package compressor は画像圧縮のためのアルゴリズムを提供します。
// 様々な画像形式に対応し、量子化や最適化を適用することで
// ファイルサイズを削減します。
package compressor

import (
	"image"
	"io"
)

// Options は圧縮オプションを表します。
type Options struct {
	// Quality はJPEG圧縮の品質設定です（0-100）
	Quality int
	// PaletteSize はPNG圧縮のパレットサイズです
	PaletteSize int
}

// DefaultOptions はデフォルトのオプションを返します。
func DefaultOptions() Options {
	return Options{
		Quality:     80,  // JPEG品質のデフォルト値
		PaletteSize: 256, // PNGパレットサイズのデフォルト値
	}
}

// Compressor は画像圧縮の共通インターフェースを定義します。
type Compressor interface {
	// Compress は画像を圧縮し、結果の画像を返します。
	Compress(img image.Image, options Options) (image.Image, error)

	// CompressBytes はバイトスライスとして提供された画像データを圧縮し、
	// 圧縮されたデータをバイトスライスとして返します。
	CompressBytes(data []byte, options Options) ([]byte, error)

	// CompressReader は Reader からの画像データを圧縮し、
	// 圧縮されたデータを Writer に書き込みます。
	CompressReader(r io.Reader, w io.Writer, options Options) error

	// SupportedFormat はこのコンプレッサーがサポートする画像形式を返します。
	SupportedFormat() string
}
