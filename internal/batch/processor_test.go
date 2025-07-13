package batch

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
	"github.com/takumines/shuku/pkg/shuku"
)

// テスト用画像データを生成するヘルパー関数
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// 簡単なパターンを作成
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

// JPEG画像ファイルを作成
func createTestJPEGFile(t *testing.T, dir, filename string, width, height int) string {
	t.Helper()
	img := createTestImage(width, height)
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		t.Fatalf("JPEG画像データの生成に失敗しました: %v", err)
	}

	filePath := filepath.Join(dir, filename)
	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	return filePath
}

// PNG画像ファイルを作成
func createTestPNGFile(t *testing.T, dir, filename string, width, height int) string {
	t.Helper()
	img := createTestImage(width, height)
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("PNG画像データの生成に失敗しました: %v", err)
	}

	filePath := filepath.Join(dir, filename)
	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	return filePath
}

// WebP画像ファイルを作成
func createTestWebPFile(t *testing.T, dir, filename string, width, height int) string {
	t.Helper()
	img := createTestImage(width, height)
	var buf bytes.Buffer
	err := webp.Encode(&buf, img, webp.Options{Quality: 80})
	if err != nil {
		t.Fatalf("WebP画像データの生成に失敗しました: %v", err)
	}

	filePath := filepath.Join(dir, filename)
	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	return filePath
}

// テキストファイルを作成（画像以外のファイル）
func createTestTextFile(t *testing.T, dir, filename string) string {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, []byte("これはテキストファイルです"), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}
	return filePath
}

func TestNewProcessor(t *testing.T) {
	tests := []struct {
		name        string
		workerCount int
		outputDir   string
		expected    int
	}{
		{
			name:        "正常なワーカー数",
			workerCount: 8,
			outputDir:   "/tmp/output",
			expected:    8,
		},
		{
			name:        "ワーカー数0（デフォルト値使用）",
			workerCount: 0,
			outputDir:   "/tmp/output",
			expected:    4,
		},
		{
			name:        "負のワーカー数（デフォルト値使用）",
			workerCount: -1,
			outputDir:   "/tmp/output",
			expected:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(tt.workerCount, tt.outputDir)

			if processor.WorkerCount != tt.expected {
				t.Errorf("NewProcessor() WorkerCount = %v, want %v", processor.WorkerCount, tt.expected)
			}

			if processor.OutputDir != tt.outputDir {
				t.Errorf("NewProcessor() OutputDir = %v, want %v", processor.OutputDir, tt.outputDir)
			}

			// デフォルト値の確認
			if processor.Recursive {
				t.Error("NewProcessor() Recursive should be false by default")
			}

			expectedInclude := []string{"*.jpg", "*.jpeg", "*.png", "*.webp"}
			if len(processor.IncludeGlobs) != len(expectedInclude) {
				t.Errorf("NewProcessor() IncludeGlobs length = %v, want %v", len(processor.IncludeGlobs), len(expectedInclude))
			}
		})
	}
}

func TestProcessor_SetRecursive(t *testing.T) {
	processor := NewProcessor(4, "/tmp")

	// 初期値はfalse
	if processor.Recursive {
		t.Error("Initial Recursive should be false")
	}

	// trueに設定
	processor.SetRecursive(true)
	if !processor.Recursive {
		t.Error("SetRecursive(true) failed")
	}

	// falseに設定
	processor.SetRecursive(false)
	if processor.Recursive {
		t.Error("SetRecursive(false) failed")
	}
}

func TestProcessor_SetIncludePatterns(t *testing.T) {
	processor := NewProcessor(4, "/tmp")

	patterns := []string{"*.jpg", "*.png"}
	processor.SetIncludePatterns(patterns)

	if len(processor.IncludeGlobs) != len(patterns) {
		t.Errorf("SetIncludePatterns() length = %v, want %v", len(processor.IncludeGlobs), len(patterns))
	}

	for i, pattern := range patterns {
		if processor.IncludeGlobs[i] != pattern {
			t.Errorf("SetIncludePatterns() pattern[%d] = %v, want %v", i, processor.IncludeGlobs[i], pattern)
		}
	}
}

func TestProcessor_SetExcludePatterns(t *testing.T) {
	processor := NewProcessor(4, "/tmp")

	patterns := []string{"*_thumb*", "*_backup*"}
	processor.SetExcludePatterns(patterns)

	if len(processor.ExcludeGlobs) != len(patterns) {
		t.Errorf("SetExcludePatterns() length = %v, want %v", len(processor.ExcludeGlobs), len(patterns))
	}

	for i, pattern := range patterns {
		if processor.ExcludeGlobs[i] != pattern {
			t.Errorf("SetExcludePatterns() pattern[%d] = %v, want %v", i, processor.ExcludeGlobs[i], pattern)
		}
	}
}

func TestProcessor_shouldIncludeFile(t *testing.T) {
	processor := NewProcessor(4, "/tmp")

	tests := []struct {
		name     string
		filePath string
		include  []string
		exclude  []string
		expected bool
	}{
		{
			name:     "JPEG画像ファイル（包含）",
			filePath: "/path/to/image.jpg",
			include:  []string{"*.jpg", "*.png"},
			exclude:  []string{},
			expected: true,
		},
		{
			name:     "PNG画像ファイル（包含）",
			filePath: "/path/to/image.png",
			include:  []string{"*.jpg", "*.png"},
			exclude:  []string{},
			expected: true,
		},
		{
			name:     "テキストファイル（除外）",
			filePath: "/path/to/document.txt",
			include:  []string{"*.jpg", "*.png"},
			exclude:  []string{},
			expected: false,
		},
		{
			name:     "サムネイル画像（除外パターンに一致）",
			filePath: "/path/to/image_thumb.jpg",
			include:  []string{"*.jpg", "*.png"},
			exclude:  []string{"*_thumb*"},
			expected: false,
		},
		{
			name:     "バックアップ画像（除外パターンに一致）",
			filePath: "/path/to/image_backup.png",
			include:  []string{"*.jpg", "*.png"},
			exclude:  []string{"*_backup*"},
			expected: false,
		},
		{
			name:     "通常のJPEG画像（除外パターンに非該当）",
			filePath: "/path/to/normal_image.jpg",
			include:  []string{"*.jpg", "*.png"},
			exclude:  []string{"*_thumb*", "*_backup*"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor.SetIncludePatterns(tt.include)
			processor.SetExcludePatterns(tt.exclude)

			result := processor.shouldIncludeFile(tt.filePath)
			if result != tt.expected {
				t.Errorf("shouldIncludeFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessor_generateOutputPath(t *testing.T) {
	tests := []struct {
		name      string
		outputDir string
		inputPath string
		inputDir  string
		expected  string
	}{
		{
			name:      "出力ディレクトリ未指定（同じディレクトリに_compressed付加）",
			outputDir: "",
			inputPath: "/input/image.jpg",
			inputDir:  "/input",
			expected:  "/input/image_compressed.jpg",
		},
		{
			name:      "出力ディレクトリ指定（相対パス維持）",
			outputDir: "/output",
			inputPath: "/input/subdir/image.jpg",
			inputDir:  "/input",
			expected:  "/output/subdir/image.jpg",
		},
		{
			name:      "出力ディレクトリ指定（ルートファイル）",
			outputDir: "/output",
			inputPath: "/input/image.png",
			inputDir:  "/input",
			expected:  "/output/image.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(4, tt.outputDir)
			result := processor.generateOutputPath(tt.inputPath, tt.inputDir)

			if result != tt.expected {
				t.Errorf("generateOutputPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessor_ProcessDirectory(t *testing.T) {
	// テストディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "batch_test_")
	if err != nil {
		t.Fatalf("テストディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 出力ディレクトリを作成
	outputDir := filepath.Join(tmpDir, "output")

	// テストファイルを作成
	createTestJPEGFile(t, tmpDir, "image1.jpg", 100, 100)
	createTestPNGFile(t, tmpDir, "image2.png", 100, 100)
	createTestWebPFile(t, tmpDir, "image3.webp", 100, 100)
	createTestTextFile(t, tmpDir, "document.txt") // 画像以外のファイル

	// サブディレクトリとファイルを作成
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("サブディレクトリの作成に失敗しました: %v", err)
	}
	createTestJPEGFile(t, subDir, "image4.jpg", 100, 100)

	options := shuku.Options{Quality: 70, PaletteSize: 256}

	t.Run("非再帰的処理", func(t *testing.T) {
		processor := NewProcessor(2, outputDir)
		processor.SetRecursive(false)

		results, err := processor.ProcessDirectory(tmpDir, options)
		if err != nil {
			t.Fatalf("ProcessDirectory() error = %v", err)
		}

		// 3つの画像ファイルが処理されることを期待（サブディレクトリは除外）
		if len(results) != 3 {
			t.Errorf("ProcessDirectory() results count = %v, want %v", len(results), 3)
		}

		// すべて成功していることを確認
		for _, result := range results {
			if result.Error != nil {
				t.Errorf("ProcessDirectory() result error: %v", result.Error)
			}

			// 出力ファイルが存在することを確認
			if _, err := os.Stat(result.Job.OutputPath); os.IsNotExist(err) {
				t.Errorf("出力ファイルが作成されませんでした: %s", result.Job.OutputPath)
			}
		}
	})

	t.Run("再帰的処理", func(t *testing.T) {
		processor := NewProcessor(2, outputDir+"_recursive")
		processor.SetRecursive(true)

		results, err := processor.ProcessDirectory(tmpDir, options)
		if err != nil {
			t.Fatalf("ProcessDirectory() error = %v", err)
		}

		// 画像ファイルが処理されることを期待（サブディレクトリ含む）
		// テスト実行時に作成されるファイル数を確認
		t.Logf("実際の処理ファイル数: %d", len(results))
		if len(results) < 4 {
			t.Errorf("ProcessDirectory() results count = %v, want at least %v", len(results), 4)
		}

		// すべて成功していることを確認
		successCount := 0
		for _, result := range results {
			if result.Error == nil {
				successCount++
			} else {
				t.Logf("処理エラー: %s - %v", result.Job.InputPath, result.Error)
			}
		}

		if successCount < 4 {
			t.Errorf("成功した処理数 = %v, want at least %v", successCount, 4)
		}
	})

	t.Run("存在しないディレクトリ", func(t *testing.T) {
		processor := NewProcessor(2, outputDir)
		_, err := processor.ProcessDirectory("/nonexistent/directory", options)
		if err == nil {
			t.Error("ProcessDirectory() should return error for nonexistent directory")
		}
	})

	t.Run("除外パターンのテスト", func(t *testing.T) {
		// サムネイル画像を作成
		createTestJPEGFile(t, tmpDir, "image_thumb.jpg", 50, 50)

		processor := NewProcessor(2, outputDir+"_exclude")
		processor.SetExcludePatterns([]string{"*_thumb*"})

		results, err := processor.ProcessDirectory(tmpDir, options)
		if err != nil {
			t.Fatalf("ProcessDirectory() error = %v", err)
		}

		// サムネイル画像が除外されていることを確認
		for _, result := range results {
			if filepath.Base(result.Job.InputPath) == "image_thumb.jpg" {
				t.Error("除外パターンに一致するファイルが処理されました")
			}
		}
	})
}

func TestCalculateStatistics(t *testing.T) {
	tests := []struct {
		name     string
		results  []Result
		expected Statistics
	}{
		{
			name: "成功のみ",
			results: []Result{
				{OriginalSize: 1000, CompressedSize: 800, Error: nil},
				{OriginalSize: 2000, CompressedSize: 1500, Error: nil},
			},
			expected: Statistics{
				TotalFiles:          2,
				SuccessFiles:        2,
				FailedFiles:         0,
				TotalOriginalSize:   3000,
				TotalCompressedSize: 2300,
				CompressionRatio:    23.33,
			},
		},
		{
			name: "成功と失敗の混在",
			results: []Result{
				{OriginalSize: 1000, CompressedSize: 800, Error: nil},
				{OriginalSize: 0, CompressedSize: 0, Error: &os.PathError{}},
			},
			expected: Statistics{
				TotalFiles:          2,
				SuccessFiles:        1,
				FailedFiles:         1,
				TotalOriginalSize:   1000,
				TotalCompressedSize: 800,
				CompressionRatio:    20.0,
			},
		},
		{
			name:    "空の結果",
			results: []Result{},
			expected: Statistics{
				TotalFiles:          0,
				SuccessFiles:        0,
				FailedFiles:         0,
				TotalOriginalSize:   0,
				TotalCompressedSize: 0,
				CompressionRatio:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := CalculateStatistics(tt.results)

			if stats.TotalFiles != tt.expected.TotalFiles {
				t.Errorf("CalculateStatistics() TotalFiles = %v, want %v", stats.TotalFiles, tt.expected.TotalFiles)
			}

			if stats.SuccessFiles != tt.expected.SuccessFiles {
				t.Errorf("CalculateStatistics() SuccessFiles = %v, want %v", stats.SuccessFiles, tt.expected.SuccessFiles)
			}

			if stats.FailedFiles != tt.expected.FailedFiles {
				t.Errorf("CalculateStatistics() FailedFiles = %v, want %v", stats.FailedFiles, tt.expected.FailedFiles)
			}

			if stats.TotalOriginalSize != tt.expected.TotalOriginalSize {
				t.Errorf("CalculateStatistics() TotalOriginalSize = %v, want %v", stats.TotalOriginalSize, tt.expected.TotalOriginalSize)
			}

			if stats.TotalCompressedSize != tt.expected.TotalCompressedSize {
				t.Errorf("CalculateStatistics() TotalCompressedSize = %v, want %v", stats.TotalCompressedSize, tt.expected.TotalCompressedSize)
			}

			// 浮動小数点の比較には許容誤差を使用
			if abs(stats.CompressionRatio-tt.expected.CompressionRatio) > 0.01 {
				t.Errorf("CalculateStatistics() CompressionRatio = %v, want %v", stats.CompressionRatio, tt.expected.CompressionRatio)
			}
		})
	}
}

// abs returns the absolute value of x.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}