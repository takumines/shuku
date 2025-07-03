package compress_test

import (
	"os"
	"testing"

	"shuku/cmd/shuku/compress"

	"github.com/urfave/cli/v2"
)

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

	// Test with unsupported format by creating a dummy .webp file
	dummyWebpFile := "../../../testdata/dummy.webp"
	// Create a dummy file with .webp extension
	err := os.WriteFile(dummyWebpFile, []byte("dummy content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}
	defer os.Remove(dummyWebpFile)
	
	args := []string{"app", "compress", "--input", dummyWebpFile, "--output", "output.webp"}
	err = app.Run(args)
	
	// WebPは未対応なので、特定のエラーメッセージが発生することが期待される動作
	if err == nil {
		t.Fatal("Expected unsupported format to fail, but it succeeded")
	}
	
	// 期待するエラーメッセージ
	expectedErrorMsg := "サポートされていない画像形式です: .webp。現在はJPEG、PNG形式に対応しています。"
	
	// エラーメッセージが期待通りかどうかを確認
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
	}
	
	t.Logf("Unsupported format failed with correct error message: %v", err)
}