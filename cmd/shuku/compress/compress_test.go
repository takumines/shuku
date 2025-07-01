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

// TestCompressAction_PNG tests PNG compression (should fail initially - this is our RED test)
func TestCompressAction_PNG(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
	}

	// Clean up any existing output file
	outputFile := "../../../test_output.png"
	os.Remove(outputFile)

	// Test PNG compression - this should FAIL initially (RED phase)
	args := []string{"app", "compress", "--input", "../../../testdata/test_image.png", "--output", outputFile}
	err := app.Run(args)
	
	// In RED phase, we expect this to fail
	if err == nil {
		t.Fatal("Expected PNG compression to fail, but it succeeded")
	}
	
	// Verify it fails with an error message
	if err != nil && err.Error() != "" {
		// For now, just verify it fails - we'll check the exact error message later
		t.Logf("PNG compression failed as expected: %v", err)
	}

	// Clean up (in case file was somehow created)
	os.Remove(outputFile)
}

// TestCompressAction_UnsupportedFormat tests unsupported formats
func TestCompressAction_UnsupportedFormat(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			compress.Cmd(),
		},
	}

	// Test with a non-existent file extension
	args := []string{"app", "compress", "--input", "testdata/fake.webp", "--output", "output.webp"}
	err := app.Run(args)
	
	if err == nil {
		t.Fatal("Expected unsupported format to fail, but it succeeded")
	}
	
	t.Logf("Unsupported format failed as expected: %v", err)
}