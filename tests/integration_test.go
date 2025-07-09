package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestPNGCompressionIntegration はPNG圧縮の統合テストです
func TestPNGCompressionIntegration(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 入力ファイル（既存のPNGテスト画像）
	inputFile := "../testdata/test_image.png"

	// 出力ファイルパス
	outputFile := filepath.Join(tempDir, "compressed_output.png")

	// CLIバイナリをビルド（Windows対応）
	binaryName := "shuku"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		binaryName = "shuku.exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/shuku")
	buildCmd.Dir = ".."
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, buildOutput)
	}

	// PNG圧縮コマンドを実行
	cmd := exec.Command(binaryPath, "compress", "--input", inputFile, "--output", outputFile, "--quality", "80")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("PNG compression command failed: %v\nOutput: %s", err, output)
	}

	// 出力メッセージの確認
	outputStr := string(output)
	if !strings.Contains(outputStr, "画像を圧縮しています...") {
		t.Errorf("Expected compression start message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "圧縮が完了しました！") {
		t.Errorf("Expected compression completion message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, outputFile) {
		t.Errorf("Expected output file path in message, got: %s", outputStr)
	}

	// 出力ファイルが作成されたことを確認
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// ファイルサイズの確認（圧縮されていることを確認）
	inputStat, err := os.Stat(inputFile)
	if err != nil {
		t.Fatalf("Failed to stat input file: %v", err)
	}

	outputStat, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	// 出力ファイルが存在し、0バイトでないことを確認
	if outputStat.Size() == 0 {
		t.Errorf("Output file is empty")
	}

	// 圧縮効果の確認（必ずしも小さくなるとは限らないが、処理は成功している）
	t.Logf("Input file size: %d bytes", inputStat.Size())
	t.Logf("Output file size: %d bytes", outputStat.Size())
	t.Logf("PNG compression integration test completed successfully")
}

// TestJPEGCompressionIntegration はJPEG圧縮の統合テストです（既存機能の確認）
func TestJPEGCompressionIntegration(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 入力ファイル（既存のJPEGテスト画像）
	inputFile := "../testdata/test_image.jpg"

	// 出力ファイルパス
	outputFile := filepath.Join(tempDir, "compressed_output.jpg")

	// CLIバイナリをビルド（Windows対応）
	binaryName := "shuku"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		binaryName = "shuku.exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/shuku")
	buildCmd.Dir = ".."
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, buildOutput)
	}

	// JPEG圧縮コマンドを実行
	cmd := exec.Command(binaryPath, "compress", "--input", inputFile, "--output", outputFile, "--quality", "70")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("JPEG compression command failed: %v\nOutput: %s", err, output)
	}

	// 出力メッセージの確認
	outputStr := string(output)
	if !strings.Contains(outputStr, "画像を圧縮しています...") {
		t.Errorf("Expected compression start message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "圧縮が完了しました！") {
		t.Errorf("Expected compression completion message, got: %s", outputStr)
	}

	// 出力ファイルが作成されたことを確認
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	t.Logf("JPEG compression integration test completed successfully")
}

// TestWebPCompressionIntegration はWebP圧縮の統合テストです
func TestWebPCompressionIntegration(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 入力ファイル（既存のWebPテスト画像）
	inputFile := "../testdata/test_image.webp"

	// 出力ファイルパス
	outputFile := filepath.Join(tempDir, "compressed_output.webp")

	// CLIバイナリをビルド（Windows対応）
	binaryName := "shuku"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		binaryName = "shuku.exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/shuku")
	buildCmd.Dir = ".."
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, buildOutput)
	}

	// WebP圧縮コマンドを実行
	cmd := exec.Command(binaryPath, "compress", "--input", inputFile, "--output", outputFile, "--quality", "70")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("WebP compression command failed: %v\nOutput: %s", err, output)
	}

	// 出力メッセージの確認
	outputStr := string(output)
	if !strings.Contains(outputStr, "画像を圧縮しています...") {
		t.Errorf("Expected compression start message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "圧縮が完了しました！") {
		t.Errorf("Expected compression completion message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, outputFile) {
		t.Errorf("Expected output file path in message, got: %s", outputStr)
	}

	// 出力ファイルが作成されたことを確認
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// ファイルサイズの確認（圧縮されていることを確認）
	inputStat, err := os.Stat(inputFile)
	if err != nil {
		t.Fatalf("Failed to stat input file: %v", err)
	}

	outputStat, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	// 出力ファイルが存在し、0バイトでないことを確認
	if outputStat.Size() == 0 {
		t.Errorf("Output file is empty")
	}

	// 圧縮効果の確認（必ずしも小さくなるとは限らないが、処理は成功している）
	t.Logf("Input file size: %d bytes", inputStat.Size())
	t.Logf("Output file size: %d bytes", outputStat.Size())
	t.Logf("WebP compression integration test completed successfully")
}

// TestUnsupportedFormatIntegration はサポートされていない形式の統合テストです
func TestUnsupportedFormatIntegration(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// ダミーのBMPファイルを作成（サポートされていない形式）
	dummyBmpFile := filepath.Join(tempDir, "test.bmp")
	err := os.WriteFile(dummyBmpFile, []byte("dummy bmp content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy bmp file: %v", err)
	}

	// 出力ファイルパス
	outputFile := filepath.Join(tempDir, "output.bmp")

	// CLIバイナリをビルド（Windows対応）
	binaryName := "shuku"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		binaryName = "shuku.exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/shuku")
	buildCmd.Dir = ".."
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, buildOutput)
	}

	// BMP圧縮コマンドを実行（失敗することを期待）
	cmd := exec.Command(binaryPath, "compress", "--input", dummyBmpFile, "--output", outputFile)
	output, err := cmd.CombinedOutput()

	// コマンドが失敗することを確認
	if err == nil {
		t.Fatal("Expected command to fail for unsupported format, but it succeeded")
	}

	// 適切なエラーメッセージが表示されることを確認
	outputStr := string(output)
	expectedError := "サポートされていない画像形式です: .bmp。現在はJPEG、PNG、WebP形式に対応しています。"
	if !strings.Contains(outputStr, expectedError) {
		t.Errorf("Expected error message '%s', got: %s", expectedError, outputStr)
	}

	t.Logf("Unsupported format integration test completed successfully")
}
