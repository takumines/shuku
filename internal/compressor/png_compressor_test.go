package compressor

import (
	"bytes"
	"image/png"
	"testing"
)

func TestPNGCompressor_Compress(t *testing.T) {
	// テスト用の画像を作成
	img := createTestImage(300, 200)

	compressor := NewPNGCompressor()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"標準圧縮", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{}

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

func TestPNGCompressor_CompressBytes(t *testing.T) {
	// テスト用の画像をPNGバイトデータに変換
	img := createTestImage(300, 200)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to create test PNG data: %v", err)
	}
	pngData := buf.Bytes()

	compressor := NewPNGCompressor()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"標準圧縮", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{}

			// 圧縮を実行
			compressed, err := compressor.CompressBytes(pngData, opts)

			// エラー検証
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(compressed) == 0 {
				t.Error("CompressBytes() returned empty data without error")
			}

			// 結果が有効なPNGデータかどうか確認
			if !tt.wantErr {
				_, err := png.Decode(bytes.NewReader(compressed))
				if err != nil {
					t.Errorf("CompressBytes() produced invalid PNG data: %v", err)
				}
			}
		})
	}

	// 無効なデータでのテスト
	t.Run("無効なデータ", func(t *testing.T) {
		invalidData := []byte("invalid png data")
		_, err := compressor.CompressBytes(invalidData, Options{})
		if err == nil {
			t.Error("CompressBytes() with invalid data should return error")
		}
	})
}

func TestPNGCompressor_CompressReader(t *testing.T) {
	// テスト用の画像をPNGバイトデータに変換
	img := createTestImage(300, 200)
	var inputBuf bytes.Buffer
	if err := png.Encode(&inputBuf, img); err != nil {
		t.Fatalf("Failed to create test PNG data: %v", err)
	}

	compressor := NewPNGCompressor()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"標準圧縮", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 入力と出力のバッファを準備
			inputReader := bytes.NewReader(inputBuf.Bytes())
			var outputBuf bytes.Buffer

			// 圧縮を実行
			err := compressor.CompressReader(inputReader, &outputBuf, Options{})

			// エラー検証
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && outputBuf.Len() == 0 {
				t.Error("CompressReader() wrote no data without error")
			}

			// 結果が有効なPNGデータかどうか確認
			if !tt.wantErr {
				_, err := png.Decode(bytes.NewReader(outputBuf.Bytes()))
				if err != nil {
					t.Errorf("CompressReader() produced invalid PNG data: %v", err)
				}
			}
		})
	}

	// 無効なデータでのテスト
	t.Run("無効なデータ", func(t *testing.T) {
		invalidReader := bytes.NewReader([]byte("invalid png data"))
		var outputBuf bytes.Buffer
		err := compressor.CompressReader(invalidReader, &outputBuf, Options{})
		if err == nil {
			t.Error("CompressReader() with invalid data should return error")
		}
	})
}

func TestPNGCompressor_SupportedFormat(t *testing.T) {
	compressor := NewPNGCompressor()
	format := compressor.SupportedFormat()
	if format != "png" {
		t.Errorf("SupportedFormat() = %v, want %v", format, "png")
	}
}

// Note: createTestImage関数はjpeg_compressor_test.goで定義されているため、ここでは定義しない
