package shuku

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/gen2brain/webp"
)

// テスト用画像データを生成するヘルパー関数
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// グラデーション効果を作成
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(((x + y) * 255) / (width + height))
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

// JPEG画像データを生成
func createJPEGData(t *testing.T, width, height int) []byte {
	t.Helper()
	img := createTestImage(width, height)
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		t.Fatalf("JPEG画像データの生成に失敗しました: %v", err)
	}
	return buf.Bytes()
}

// PNG画像データを生成
func createPNGData(t *testing.T, width, height int) []byte {
	t.Helper()
	img := createTestImage(width, height)
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("PNG画像データの生成に失敗しました: %v", err)
	}
	return buf.Bytes()
}

// WebP画像データを生成
func createWebPData(t *testing.T, width, height int) []byte {
	t.Helper()
	img := createTestImage(width, height)
	var buf bytes.Buffer
	err := webp.Encode(&buf, img, webp.Options{Quality: 80})
	if err != nil {
		t.Fatalf("WebP画像データの生成に失敗しました: %v", err)
	}
	return buf.Bytes()
}

// テスト用一時ファイルを作成
func createTempFile(t *testing.T, data []byte, suffix string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "shuku_test_*"+suffix)
	if err != nil {
		t.Fatalf("一時ファイルの作成に失敗しました: %v", err)
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(data)
	if err != nil {
		t.Fatalf("一時ファイルへの書き込みに失敗しました: %v", err)
	}

	return tmpFile.Name()
}

func TestCompress(t *testing.T) {
	tests := []struct {
		name     string
		dataFunc func(*testing.T, int, int) []byte
		options  Options
		wantErr  bool
	}{
		{
			name:     "JPEG圧縮_標準品質",
			dataFunc: createJPEGData,
			options:  Options{Quality: 80},
			wantErr:  false,
		},
		{
			name:     "JPEG圧縮_低品質",
			dataFunc: createJPEGData,
			options:  Options{Quality: 30},
			wantErr:  false,
		},
		{
			name:     "JPEG圧縮_高品質",
			dataFunc: createJPEGData,
			options:  Options{Quality: 95},
			wantErr:  false,
		},
		{
			name:     "PNG圧縮",
			dataFunc: createPNGData,
			options:  Options{PaletteSize: 256},
			wantErr:  false,
		},
		{
			name:     "WebP圧縮_標準品質",
			dataFunc: createWebPData,
			options:  Options{Quality: 80},
			wantErr:  false,
		},
		{
			name:     "WebP圧縮_低品質",
			dataFunc: createWebPData,
			options:  Options{Quality: 50},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テストデータを生成
			data := tt.dataFunc(t, 100, 100)

			// 圧縮を実行
			result, err := Compress(data, tt.options)

			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 結果が空でないことを確認
				if len(result) == 0 {
					t.Error("Compress() 圧縮結果が空です")
				}

				// 圧縮されたデータが有効であることを確認（サイズが元データと異なること）
				if len(result) == len(data) {
					t.Log("注意: 圧縮前後でデータサイズが同じです（圧縮効果なし）")
				}
			}
		})
	}
}

func TestCompressFile(t *testing.T) {
	// テストディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "shuku_test_")
	if err != nil {
		t.Fatalf("テストディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name       string
		dataFunc   func(*testing.T, int, int) []byte
		suffix     string
		outputPath string
		options    Options
		wantErr    bool
	}{
		{
			name:       "JPEG圧縮_出力パス指定",
			dataFunc:   createJPEGData,
			suffix:     ".jpg",
			outputPath: filepath.Join(tmpDir, "output.jpg"),
			options:    Options{Quality: 70},
			wantErr:    false,
		},
		{
			name:       "PNG圧縮_出力パス指定",
			dataFunc:   createPNGData,
			suffix:     ".png",
			outputPath: filepath.Join(tmpDir, "output.png"),
			options:    Options{PaletteSize: 128},
			wantErr:    false,
		},
		{
			name:       "WebP圧縮_出力パス指定",
			dataFunc:   createWebPData,
			suffix:     ".webp",
			outputPath: filepath.Join(tmpDir, "output.webp"),
			options:    Options{Quality: 60},
			wantErr:    false,
		},
		{
			name:       "JPEG圧縮_自動出力パス",
			dataFunc:   createJPEGData,
			suffix:     ".jpg",
			outputPath: "", // 自動生成
			options:    Options{Quality: 80},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テストファイルを作成
			data := tt.dataFunc(t, 150, 150)
			inputPath := createTempFile(t, data, tt.suffix)
			defer os.Remove(inputPath)

			// 圧縮を実行
			err := CompressFile(inputPath, tt.outputPath, tt.options)

			if (err != nil) != tt.wantErr {
				t.Errorf("CompressFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 出力ファイルパスを決定
				expectedOutputPath := tt.outputPath
				if expectedOutputPath == "" {
					// 自動生成されたパスを確認
					ext := filepath.Ext(inputPath)
					expectedOutputPath = inputPath[:len(inputPath)-len(ext)] + "_compressed" + ext
				}

				// 出力ファイルが存在することを確認
				if _, err := os.Stat(expectedOutputPath); os.IsNotExist(err) {
					t.Errorf("出力ファイルが作成されませんでした: %s", expectedOutputPath)
				} else {
					// テスト後にクリーンアップ
					defer os.Remove(expectedOutputPath)

					// ファイルサイズを確認
					inputInfo, _ := os.Stat(inputPath)
					outputInfo, _ := os.Stat(expectedOutputPath)

					if outputInfo.Size() == 0 {
						t.Error("出力ファイルが空です")
					}

					t.Logf("元のサイズ: %d bytes, 圧縮後: %d bytes", inputInfo.Size(), outputInfo.Size())
				}
			}
		})
	}
}

func TestCompressImage(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		options Options
		wantErr bool
	}{
		{
			name:    "JPEG形式で画像圧縮",
			format:  "jpeg",
			options: Options{Quality: 75},
			wantErr: false,
		},
		{
			name:    "PNG形式で画像圧縮",
			format:  "png",
			options: Options{PaletteSize: 256},
			wantErr: false,
		},
		{
			name:    "WebP形式で画像圧縮",
			format:  "webp",
			options: Options{Quality: 85},
			wantErr: false,
		},
		{
			name:    "サポートされていない形式",
			format:  "bmp",
			options: Options{Quality: 80},
			wantErr: true,
		},
		{
			name:    "空の形式文字列",
			format:  "",
			options: Options{Quality: 80},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト画像を作成
			img := createTestImage(100, 100)

			// 圧縮を実行
			result, err := CompressImage(img, tt.format, tt.options)

			if (err != nil) != tt.wantErr {
				t.Errorf("CompressImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 結果が nil でないことを確認
				if result == nil {
					t.Error("CompressImage() 圧縮結果が nil です")
				}

				// 画像の基本属性を確認
				bounds := result.Bounds()
				if bounds.Dx() == 0 || bounds.Dy() == 0 {
					t.Error("CompressImage() 圧縮後の画像サイズが無効です")
				}
			}
		})
	}
}

// エラーケースのテスト
func TestCompressErrorCases(t *testing.T) {
	t.Run("無効な画像データ", func(t *testing.T) {
		invalidData := []byte("これは画像データではありません")
		_, err := Compress(invalidData, Options{Quality: 80})
		if err == nil {
			t.Error("無効な画像データに対してエラーが発生しませんでした")
		}
	})

	t.Run("空の画像データ", func(t *testing.T) {
		emptyData := []byte{}
		_, err := Compress(emptyData, Options{Quality: 80})
		if err == nil {
			t.Error("空の画像データに対してエラーが発生しませんでした")
		}
	})
}

func TestCompressFileErrorCases(t *testing.T) {
	t.Run("存在しないファイル", func(t *testing.T) {
		err := CompressFile("/path/to/nonexistent/file.jpg", "/tmp/output.jpg", Options{Quality: 80})
		if err == nil {
			t.Error("存在しないファイルに対してエラーが発生しませんでした")
		}
	})

	t.Run("拡張子なしファイル", func(t *testing.T) {
		// 拡張子なしの一時ファイルを作成
		tmpFile, err := os.CreateTemp("", "shuku_test_no_ext")
		if err != nil {
			t.Fatalf("一時ファイルの作成に失敗しました: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		err = CompressFile(tmpFile.Name(), "/tmp/output.jpg", Options{Quality: 80})
		if err == nil {
			t.Error("拡張子なしファイルに対してエラーが発生しませんでした")
		}
	})

	t.Run("無効な出力ディレクトリ", func(t *testing.T) {
		// 有効な入力ファイルを作成
		data := createJPEGData(t, 50, 50)
		inputPath := createTempFile(t, data, ".jpg")
		defer os.Remove(inputPath)

		// 無効な出力パス（存在しないディレクトリ）
		invalidOutputPath := "/nonexistent/directory/output.jpg"
		err := CompressFile(inputPath, invalidOutputPath, Options{Quality: 80})
		if err == nil {
			t.Error("無効な出力ディレクトリに対してエラーが発生しませんでした")
		}
	})
}

// detectImageFormat 関数の単体テスト
func TestDetectImageFormat(t *testing.T) {
	tests := []struct {
		name     string
		dataFunc func(*testing.T, int, int) []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "JPEG形式検出",
			dataFunc: createJPEGData,
			expected: "jpeg",
			wantErr:  false,
		},
		{
			name:     "PNG形式検出",
			dataFunc: createPNGData,
			expected: "png",
			wantErr:  false,
		},
		{
			name:     "WebP形式検出",
			dataFunc: createWebPData,
			expected: "webp",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.dataFunc(t, 50, 50)
			format, err := detectImageFormat(data)

			if (err != nil) != tt.wantErr {
				t.Errorf("detectImageFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && format != tt.expected {
				t.Errorf("detectImageFormat() = %v, want %v", format, tt.expected)
			}
		})
	}

	t.Run("無効なデータ", func(t *testing.T) {
		invalidData := []byte("無効な画像データ")
		_, err := detectImageFormat(invalidData)
		if err == nil {
			t.Error("無効なデータに対してエラーが発生しませんでした")
		}
	})

	t.Run("空のデータ", func(t *testing.T) {
		emptyData := []byte{}
		_, err := detectImageFormat(emptyData)
		if err == nil {
			t.Error("空のデータに対してエラーが発生しませんでした")
		}
	})
}
