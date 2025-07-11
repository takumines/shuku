package compressor

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"os"
	"testing"

	"github.com/gen2brain/webp"
)

// createTestWebPImage はテスト用のWebP画像を作成します
func createTestWebPImage() ([]byte, error) {
	// 簡単なテスト画像を作成
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	red := color.RGBA{255, 0, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{red}, image.Point{}, draw.Src)

	// WebP形式でエンコード
	var buf bytes.Buffer
	err := webp.Encode(&buf, img, webp.Options{Quality: 80})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestWebPCompressor_SupportedFormat(t *testing.T) {
	compressor := NewWebPCompressor()
	expected := "webp"
	actual := compressor.SupportedFormat()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestWebPCompressor_validateQuality(t *testing.T) {
	compressor := NewWebPCompressor()

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Valid quality", 80, 80},
		{"Below minimum", -10, 0},
		{"Above maximum", 110, 100},
		{"Minimum", 0, 0},
		{"Maximum", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := compressor.validateQuality(tt.input)
			if actual != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

func TestWebPCompressor_CompressBytes(t *testing.T) {
	compressor := NewWebPCompressor()

	// テスト用WebP画像を作成
	testWebPData, err := createTestWebPImage()
	if err != nil {
		t.Fatalf("Failed to create test WebP image: %v", err)
	}

	tests := []struct {
		name         string
		data         []byte
		options      Options
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Valid WebP compression",
			data:        testWebPData,
			options:     Options{Quality: 80},
			expectError: false,
		},
		{
			name:        "High quality compression",
			data:        testWebPData,
			options:     Options{Quality: 100},
			expectError: false,
		},
		{
			name:        "Low quality compression",
			data:        testWebPData,
			options:     Options{Quality: 10},
			expectError: false,
		},
		{
			name:         "Invalid data",
			data:         []byte("invalid webp data"),
			options:      Options{Quality: 80},
			expectError:  true,
			errorMessage: "入力データが有効なWebP画像ではありません",
		},
		{
			name:         "Empty data",
			data:         []byte{},
			options:      Options{Quality: 80},
			expectError:  true,
			errorMessage: "入力データが有効なWebP画像ではありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compressor.CompressBytes(tt.data, tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorMessage != "" {
					compressErr, ok := err.(*CompressError)
					if ok && compressErr.Message != tt.errorMessage {
						t.Errorf("Expected error message %s, got %s", tt.errorMessage, compressErr.Message)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) == 0 {
					t.Errorf("Expected non-empty result")
				}
			}
		})
	}
}

func TestWebPCompressor_Compress(t *testing.T) {
	compressor := NewWebPCompressor()

	// テスト用画像を作成
	testImg := image.NewRGBA(image.Rect(0, 0, 50, 50))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(testImg, testImg.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)

	tests := []struct {
		name        string
		img         image.Image
		options     Options
		expectError bool
	}{
		{
			name:        "Valid image compression",
			img:         testImg,
			options:     Options{Quality: 80},
			expectError: false,
		},
		{
			name:        "High quality compression",
			img:         testImg,
			options:     Options{Quality: 100},
			expectError: false,
		},
		{
			name:        "Low quality compression",
			img:         testImg,
			options:     Options{Quality: 1},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compressor.Compress(tt.img, tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Expected non-nil result")
				}
			}
		})
	}
}

func TestWebPCompressor_CompressReader(t *testing.T) {
	compressor := NewWebPCompressor()

	// テスト用WebP画像を作成
	testWebPData, err := createTestWebPImage()
	if err != nil {
		t.Fatalf("Failed to create test WebP image: %v", err)
	}

	tests := []struct {
		name         string
		data         []byte
		options      Options
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Valid WebP reader compression",
			data:        testWebPData,
			options:     Options{Quality: 80},
			expectError: false,
		},
		{
			name:         "Invalid data",
			data:         []byte("invalid webp data"),
			options:      Options{Quality: 80},
			expectError:  true,
			errorMessage: "入力データが有効なWebP画像ではありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.data)
			var writer bytes.Buffer

			err := compressor.CompressReader(reader, &writer, tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorMessage != "" {
					compressErr, ok := err.(*CompressError)
					if ok && compressErr.Message != tt.errorMessage {
						t.Errorf("Expected error message %s, got %s", tt.errorMessage, compressErr.Message)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if writer.Len() == 0 {
					t.Errorf("Expected non-empty writer")
				}
			}
		})
	}
}

func TestWebPCompressor_Integration(t *testing.T) {
	compressor := NewWebPCompressor()

	// 実際のWebPファイルが存在する場合のテスト
	testFile := "../../testdata/test_image.webp"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test WebP file not found, skipping integration test")
	}

	// ファイルを読み込み
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// 圧縮テスト
	options := Options{Quality: 70}
	result, err := compressor.CompressBytes(data, options)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty compressed result")
	}

	// 結果が有効なWebP形式であることを確認
	_, err = webp.Decode(bytes.NewReader(result))
	if err != nil {
		t.Errorf("Compressed result is not valid WebP: %v", err)
	}
}
