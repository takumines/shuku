package compressor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
)

// JPEGCompressor はJPEG形式の画像を圧縮するための実装です。
// JPEG品質設定を調整することで圧縮率を制御します。
type JPEGCompressor struct{}

// NewJPEGCompressor は新しいJPEGCompressorインスタンスを作成します。
func NewJPEGCompressor() *JPEGCompressor {
	return &JPEGCompressor{}
}

// Compress はJPEG画像を圧縮します。
// options.Qualityは0-100の値を使用して圧縮品質を指定します。
// 値が低いほどファイルサイズは小さくなりますが、画質は劣化します。
func (j *JPEGCompressor) Compress(img image.Image, options Options) (image.Image, error) {
	// JPEG圧縮を適用したバイトデータを取得
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: j.validateQuality(options.Quality),
	})
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "JPEG",
		}
	}

	// 圧縮されたバイトデータを画像として再デコード
	compressed, err := jpeg.Decode(&buf)
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "JPEG",
		}
	}

	return compressed, nil
}

// CompressBytes はバイト配列として提供されたJPEG画像データを圧縮します。
func (j *JPEGCompressor) CompressBytes(data []byte, options Options) ([]byte, error) {
	// 入力データが有効なJPEG画像であることを確認
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "JPEG",
			Message:     "入力データが有効なJPEG画像ではありません",
		}
	}

	// 圧縮を適用
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: j.validateQuality(options.Quality),
	})
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "JPEG",
		}
	}

	return buf.Bytes(), nil
}

// CompressReader はリーダーから読み取ったJPEG画像データを圧縮し、ライターに書き込みます。
func (j *JPEGCompressor) CompressReader(r io.Reader, w io.Writer, options Options) error {
	// 入力データが有効なJPEG画像であることを確認
	img, err := jpeg.Decode(r)
	if err != nil {
		return &CompressError{
			OriginalErr: err,
			Format:      "JPEG",
			Message:     "入力データが有効なJPEG画像ではありません",
		}
	}

	// 圧縮を適用して結果をライターに書き込む
	err = jpeg.Encode(w, img, &jpeg.Options{
		Quality: j.validateQuality(options.Quality),
	})
	if err != nil {
		return &CompressError{
			OriginalErr: err,
			Format:      "JPEG",
		}
	}

	return nil
}

// SupportedFormat はこのコンプレッサーがサポートするフォーマットを返します。
func (j *JPEGCompressor) SupportedFormat() string {
	return "jpeg"
}

// validateQuality は品質パラメータが有効な範囲（0-100）に収まるように調整します。
func (j *JPEGCompressor) validateQuality(quality int) int {
	if quality < 0 {
		return 0
	}
	if quality > 100 {
		return 100
	}
	return quality
}

// CompressError は圧縮処理で発生したエラーを表します。
type CompressError struct {
	OriginalErr error
	Format      string
	Message     string
}

// Error はエラーメッセージを返します。
func (e *CompressError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("圧縮エラー(%s): %s - %v", e.Format, e.Message, e.OriginalErr)
	}
	return fmt.Sprintf("圧縮エラー(%s): %v", e.Format, e.OriginalErr)
}
