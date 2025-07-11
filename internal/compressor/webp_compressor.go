package compressor

import (
	"bytes"
	"image"
	"io"

	"github.com/gen2brain/webp"
)

// WebPCompressor はWebP形式の画像を圧縮するための実装です。
// WebP品質設定を調整することで圧縮率を制御します。
type WebPCompressor struct{}

// NewWebPCompressor は新しいWebPCompressorインスタンスを作成します。
func NewWebPCompressor() *WebPCompressor {
	return &WebPCompressor{}
}

// Compress はWebP画像を圧縮します。
// options.Qualityは0-100の値を使用して圧縮品質を指定します。
// 値が低いほどファイルサイズは小さくなりますが、画質は劣化します。
func (w *WebPCompressor) Compress(img image.Image, options Options) (image.Image, error) {
	// WebP圧縮を適用したバイトデータを取得
	var buf bytes.Buffer
	quality := w.validateQuality(options.Quality)

	err := webp.Encode(&buf, img, webp.Options{
		Quality: quality,
	})
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "WebP",
		}
	}

	// 圧縮されたバイトデータを画像として再デコード
	compressed, err := webp.Decode(&buf)
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "WebP",
		}
	}

	return compressed, nil
}

// CompressBytes はバイト配列として提供されたWebP画像データを圧縮します。
func (w *WebPCompressor) CompressBytes(data []byte, options Options) ([]byte, error) {
	// 入力データが有効なWebP画像であることを確認
	img, err := webp.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "WebP",
			Message:     "入力データが有効なWebP画像ではありません",
		}
	}

	// 圧縮を適用
	var buf bytes.Buffer
	quality := w.validateQuality(options.Quality)

	err = webp.Encode(&buf, img, webp.Options{
		Quality: quality,
	})
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "WebP",
		}
	}

	return buf.Bytes(), nil
}

// CompressReader はリーダーから読み取ったWebP画像データを圧縮し、ライターに書き込みます。
func (w *WebPCompressor) CompressReader(r io.Reader, wr io.Writer, options Options) error {
	// 入力データが有効なWebP画像であることを確認
	img, err := webp.Decode(r)
	if err != nil {
		return &CompressError{
			OriginalErr: err,
			Format:      "WebP",
			Message:     "入力データが有効なWebP画像ではありません",
		}
	}

	// 圧縮を適用して結果をライターに書き込む
	quality := w.validateQuality(options.Quality)

	err = webp.Encode(wr, img, webp.Options{
		Quality: quality,
	})
	if err != nil {
		return &CompressError{
			OriginalErr: err,
			Format:      "WebP",
		}
	}

	return nil
}

// SupportedFormat はこのコンプレッサーがサポートするフォーマットを返します。
func (w *WebPCompressor) SupportedFormat() string {
	return "webp"
}

// validateQuality は品質パラメータが有効な範囲（0-100）に収まるように調整します。
func (w *WebPCompressor) validateQuality(quality int) int {
	if quality < 0 {
		return 0
	}
	if quality > 100 {
		return 100
	}
	return quality
}
