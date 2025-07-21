package compress_test

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/takumines/shuku/cmd/shuku/compress"

	"github.com/urfave/cli/v2"
)

// createTestImage creates a test JPEG image file
func createTestImage(t *testing.T, filePath string) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Create a simple pattern
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			r := uint8((x * 255) / 100)
			g := uint8((y * 255) / 100)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		t.Fatalf("Failed to encode JPEG: %v", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("Failed to write test image: %v", err)
	}
}

// TestCompressAction_JPEG tests JPEG compression (existing functionality)
func TestCompressAction_JPEG(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
	}

	// Clean up any existing output file
	outputFile := "../../../test_output.jpg"
	os.Remove(outputFile)

	// Test JPEG compression
	args := []string{"app", "compress", "--input", "../../../testdata/test_image.jpg", "--output", outputFile}
	err := app.Run(args)
	if err != nil {
		t.Fatalf("JPEG compression failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Clean up
	os.Remove(outputFile)
}

// TestCompressAction_PNG tests PNG compression (now should succeed - GREEN phase)
func TestCompressAction_PNG(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
	}

	// Clean up any existing output file
	outputFile := "../../../test_output.png"
	os.Remove(outputFile)

	// Test PNG compression - this should succeed in GREEN phase
	args := []string{"app", "compress", "--input", "../../../testdata/test_image.png", "--output", outputFile}
	err := app.Run(args)
	if err != nil {
		t.Fatalf("PNG compression failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Clean up
	os.Remove(outputFile)
}

// TestCompressAction_WebP tests WebP compression
func TestCompressAction_WebP(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
	}

	// Clean up any existing output file
	outputFile := "../../../test_output.webp"
	os.Remove(outputFile)

	// Test WebP compression - this should succeed
	args := []string{"app", "compress", "--input", "../../../testdata/test_image.webp", "--output", outputFile}
	err := app.Run(args)
	if err != nil {
		t.Fatalf("WebP compression failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Clean up
	os.Remove(outputFile)
}

// TestCompressAction_UnsupportedFormat tests unsupported formats
func TestCompressAction_UnsupportedFormat(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			// テスト中はexit処理をスキップ
		},
	}

	// Test with unsupported format by creating a dummy .bmp file
	dummyBmpFile := "../../../testdata/dummy.bmp"
	// Create a dummy file with .bmp extension
	err := os.WriteFile(dummyBmpFile, []byte("dummy content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}
	defer os.Remove(dummyBmpFile)

	args := []string{"app", "compress", "--input", dummyBmpFile, "--output", "output.bmp"}
	err = app.Run(args)

	// BMPは未対応なので、特定のエラーメッセージが発生することが期待される動作
	if err == nil {
		t.Fatal("Expected unsupported format to fail, but it succeeded")
	}

	// 期待するエラーメッセージ
	expectedErrorMsg := "サポートされていない画像形式です: .bmp。現在はJPEG、PNG、WebP形式に対応しています。"

	// エラーメッセージが期待通りかどうかを確認
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
	}

	t.Logf("Unsupported format failed with correct error message: %v", err)
}

// TestCompressAction_NoInput tests missing input parameter
func TestCompressAction_NoInput(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			// テスト中はexit処理をスキップ
		},
	}

	// Test with no input parameter
	args := []string{"app", "compress"}
	err := app.Run(args)

	// Should fail with input missing error
	if err == nil {
		t.Fatal("Expected missing input to fail, but it succeeded")
	}

	expectedErrorMsg := "Required flag \"input\" not set"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
	}
}

// TestCompressAction_FileNotFound tests non-existent input file
func TestCompressAction_FileNotFound(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			// テスト中はexit処理をスキップ
		},
	}

	// Test with non-existent file
	args := []string{"app", "compress", "--input", "non_existent_file.jpg"}
	err := app.Run(args)

	// Should fail with file not found error
	if err == nil {
		t.Fatal("Expected file not found to fail, but it succeeded")
	}

	if !strings.Contains(err.Error(), "が見つかりません") {
		t.Errorf("Expected file not found error message, but got '%s'", err.Error())
	}
}

// TestCompressAction_VerboseMode tests verbose mode functionality
func TestCompressAction_VerboseMode(t *testing.T) {
	// Create temporary directory and test image
	tempDir, err := os.MkdirTemp("", "compress_verbose_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "input.jpg")
	outputFile := filepath.Join(tempDir, "output.jpg")
	createTestImage(t, inputFile)

	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with verbose mode
	args := []string{"app", "compress", "--input", inputFile, "--output", outputFile, "--verbose"}
	err = app.Run(args)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("Verbose mode compression failed: %v", err)
	}

	// Check verbose output contains expected information
	if !strings.Contains(output, "入力ファイル:") {
		t.Errorf("Expected verbose output to contain input file info, got: %s", output)
	}
	if !strings.Contains(output, "出力ファイル:") {
		t.Errorf("Expected verbose output to contain output file info, got: %s", output)
	}
	if !strings.Contains(output, "圧縮品質:") {
		t.Errorf("Expected verbose output to contain quality info, got: %s", output)
	}
	if !strings.Contains(output, "圧縮率:") {
		t.Errorf("Expected verbose output to contain compression ratio, got: %s", output)
	}
}

// TestCompressAction_DefaultOutput tests default output file name generation
func TestCompressAction_DefaultOutput(t *testing.T) {
	// Create temporary directory and test image
	tempDir, err := os.MkdirTemp("", "compress_default_output_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "test.jpg")
	createTestImage(t, inputFile)

	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
	}

	// Test without specifying output (should use default)
	args := []string{"app", "compress", "--input", inputFile}
	err = app.Run(args)

	if err != nil {
		t.Fatalf("Default output compression failed: %v", err)
	}

	// Check that default output file was created
	expectedOutput := filepath.Join(tempDir, "test_compressed.jpg")
	if _, err := os.Stat(expectedOutput); os.IsNotExist(err) {
		t.Errorf("Expected default output file '%s' was not created", expectedOutput)
	}
}

// TestCompressAction_CustomQuality tests custom quality setting
func TestCompressAction_CustomQuality(t *testing.T) {
	// Create temporary directory and test image
	tempDir, err := os.MkdirTemp("", "compress_quality_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "input.jpg")
	outputFile := filepath.Join(tempDir, "output.jpg")
	createTestImage(t, inputFile)

	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
	}

	// Test with custom quality setting
	args := []string{"app", "compress", "--input", inputFile, "--output", outputFile, "--quality", "50"}
	err = app.Run(args)

	if err != nil {
		t.Fatalf("Custom quality compression failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file with custom quality was not created: %s", outputFile)
	}
}
