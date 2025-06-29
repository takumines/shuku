package compressor

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"testing"
)

func TestJPEGCompressor_Compress(t *testing.T) {
	// テスト用の画像を作成
	img := createTestImage(300, 200)

	compressor := NewJPEGCompressor()

	tests := []struct {
		name    string
		quality int
		wantErr bool
	}{
		{"高品質", 90, false},
		{"中品質", 50, false},
		{"低品質", 10, false},
		{"品質超過", 101, false}, // validateQualityで100に調整される
		{"品質不足", -10, false}, // validateQualityで0に調整される
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{Quality: tt.quality}

			// 圧縮を実行
			compressed, err := compressor.Compress(img, opts)

			// エラー検証
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && compressed == nil {
				t.Error("Compress() returned nil image without error")
			}

			// 画像サイズが変わっていないか確認
			if !tt.wantErr && (compressed.Bounds().Dx() != img.Bounds().Dx() ||
				compressed.Bounds().Dy() != img.Bounds().Dy()) {
				t.Errorf("Compress() image dimensions changed: got %v, want %v",
					compressed.Bounds(), img.Bounds())
			}
		})
	}
}

func TestJPEGCompressor_CompressBytes(t *testing.T) {
	// テスト用の画像をJPEGバイトデータに変換
	img := createTestImage(300, 200)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("Failed to create test JPEG data: %v", err)
	}
	jpegData := buf.Bytes()

	compressor := NewJPEGCompressor()

	tests := []struct {
		name    string
		quality int
		wantErr bool
	}{
		{"高品質", 90, false},
		{"低品質", 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{Quality: tt.quality}

			// 圧縮を実行
			compressed, err := compressor.CompressBytes(jpegData, opts)

			// エラー検証
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(compressed) == 0 {
				t.Error("CompressBytes() returned empty data without error")
			}

			// 低品質の方がファイルサイズが小さいことを確認（品質によるサイズ変化の検証）
			if !tt.wantErr && tt.quality < 90 && len(compressed) >= len(jpegData) {
				t.Errorf("CompressBytes() with quality %d should produce smaller data than original", tt.quality)
			}
		})
	}

	// 無効なデータでのテスト
	t.Run("無効なデータ", func(t *testing.T) {
		invalidData := []byte("invalid jpeg data")
		_, err := compressor.CompressBytes(invalidData, Options{Quality: 80})
		if err == nil {
			t.Error("CompressBytes() with invalid data should return error")
		}
	})
}

func TestJPEGCompressor_CompressReader(t *testing.T) {
	// テスト用の画像をJPEGバイトデータに変換
	img := createTestImage(300, 200)
	var inputBuf bytes.Buffer
	if err := jpeg.Encode(&inputBuf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("Failed to create test JPEG data: %v", err)
	}

	compressor := NewJPEGCompressor()

	tests := []struct {
		name    string
		quality int
		wantErr bool
	}{
		{"標準品質", 80, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 入力と出力のバッファを準備
			inputReader := bytes.NewReader(inputBuf.Bytes())
			var outputBuf bytes.Buffer

			// 圧縮を実行
			err := compressor.CompressReader(inputReader, &outputBuf, Options{Quality: tt.quality})

			// エラー検証
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && outputBuf.Len() == 0 {
				t.Error("CompressReader() wrote no data without error")
			}
		})
	}

	// 無効なデータでのテスト
	t.Run("無効なデータ", func(t *testing.T) {
		invalidReader := bytes.NewReader([]byte("invalid jpeg data"))
		var outputBuf bytes.Buffer
		err := compressor.CompressReader(invalidReader, &outputBuf, Options{Quality: 80})
		if err == nil {
			t.Error("CompressReader() with invalid data should return error")
		}
	})
}

func TestJPEGCompressor_SupportedFormat(t *testing.T) {
	compressor := NewJPEGCompressor()
	format := compressor.SupportedFormat()
	if format != "jpeg" {
		t.Errorf("SupportedFormat() = %v, want %v", format, "jpeg")
	}
}

func TestJPEGCompressor_validateQuality(t *testing.T) {
	compressor := NewJPEGCompressor()
	tests := []struct {
		input int
		want  int
	}{
		{50, 50},   // 通常の値はそのまま
		{0, 0},     // 最小値
		{100, 100}, // 最大値
		{-10, 0},   // 下限を下回る場合は0に
		{110, 100}, // 上限を超える場合は100に
	}

	for _, tt := range tests {
		got := compressor.validateQuality(tt.input)
		if got != tt.want {
			t.Errorf("validateQuality(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

// テスト用の単色画像を作成する
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)
	return img
}
