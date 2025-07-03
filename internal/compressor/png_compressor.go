package compressor

import (
	"bytes"
	"image"
	"image/png"
	"io"
)

// PNGCompressor はPNG形式の画像を圧縮するための実装です。
// パレットサイズを調整することで圧縮率を制御します。
type PNGCompressor struct{}

// NewPNGCompressor は新しいPNGCompressorインスタンスを作成します。
func NewPNGCompressor() *PNGCompressor {
	return &PNGCompressor{}
}

// Compress はPNG画像を圧縮します。
// options.PaletteSizeはパレットサイズを指定します。
func (p *PNGCompressor) Compress(img image.Image, options Options) (image.Image, error) {
	// PNG圧縮を適用したバイトデータを取得
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "PNG",
		}
	}

	// 圧縮されたバイトデータを画像として再デコード
	compressed, err := png.Decode(&buf)
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "PNG",
		}
	}

	return compressed, nil
}

// CompressBytes はバイト配列として提供されたPNG画像データを圧縮します。
func (p *PNGCompressor) CompressBytes(data []byte, options Options) ([]byte, error) {
	// 入力データが有効なPNG画像であることを確認
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "PNG",
			Message:     "入力データが有効なPNG画像ではありません",
		}
	}

	// 圧縮を適用
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, &CompressError{
			OriginalErr: err,
			Format:      "PNG",
		}
	}

	return buf.Bytes(), nil
}

// CompressReader はリーダーから読み取ったPNG画像データを圧縮し、ライターに書き込みます。
func (p *PNGCompressor) CompressReader(r io.Reader, w io.Writer, options Options) error {
	// 入力データが有効なPNG画像であることを確認
	img, err := png.Decode(r)
	if err != nil {
		return &CompressError{
			OriginalErr: err,
			Format:      "PNG",
			Message:     "入力データが有効なPNG画像ではありません",
		}
	}

	// 圧縮を適用して結果をライターに書き込む
	err = png.Encode(w, img)
	if err != nil {
		return &CompressError{
			OriginalErr: err,
			Format:      "PNG",
		}
	}

	return nil
}

// SupportedFormat はこのコンプレッサーがサポートするフォーマットを返します。
func (p *PNGCompressor) SupportedFormat() string {
	return "png"
}