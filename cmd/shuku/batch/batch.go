package batch

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/takumines/shuku/internal/batch"
	"github.com/takumines/shuku/pkg/shuku"
	"github.com/urfave/cli/v2"
)

// Cmd returns the batch command.
func Cmd() *cli.Command {
	return &cli.Command{
		Name:    "batch",
		Usage:   "Compress multiple images in a directory.",
		Aliases: []string{"b"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Input directory path",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output directory path (optional, defaults to same as input with '_compressed' suffix)",
			},
			&cli.IntFlag{
				Name:    "quality",
				Aliases: []string{"q"},
				Value:   80,
				Usage:   "JPEG/WebP quality (0-100)",
			},
			&cli.IntFlag{
				Name:  "palette-size",
				Value: 256,
				Usage: "PNG palette size (8, 16, 32, 64, 128, 256)",
			},
			&cli.IntFlag{
				Name:    "workers",
				Aliases: []string{"w"},
				Value:   runtime.NumCPU(),
				Usage:   "Number of parallel workers",
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "Process directories recursively",
			},
			&cli.StringFlag{
				Name:  "include",
				Usage: "File patterns to include (comma-separated, e.g., '*.jpg,*.png')",
				Value: "*.jpg,*.jpeg,*.png,*.webp",
			},
			&cli.StringFlag{
				Name:  "exclude",
				Usage: "File patterns to exclude (comma-separated, e.g., '*_thumb*,*_backup*')",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Show detailed information",
			},
			&cli.BoolFlag{
				Name:  "stats",
				Usage: "Show compression statistics",
			},
		},
		Action: batchAction,
	}
}

// batchAction is the action for the batch command.
func batchAction(c *cli.Context) error {
	// 入力ディレクトリを取得
	inputDir := c.String("input")
	if inputDir == "" {
		return cli.Exit("入力ディレクトリが指定されていません。--input または -i オプションで指定してください。", 1)
	}

	// 入力ディレクトリの存在確認
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("入力ディレクトリが存在しません: %s", inputDir), 1)
	}

	// オプションの設定
	options := shuku.Options{
		Quality:     c.Int("quality"),
		PaletteSize: c.Int("palette-size"),
	}

	// バッチプロセッサーの設定
	processor := batch.NewProcessor(c.Int("workers"), c.String("output"))
	processor.SetRecursive(c.Bool("recursive"))

	// 包含パターンの設定
	if includePatterns := c.String("include"); includePatterns != "" {
		patterns := strings.Split(includePatterns, ",")
		for i, pattern := range patterns {
			patterns[i] = strings.TrimSpace(pattern)
		}
		processor.SetIncludePatterns(patterns)
	}

	// 除外パターンの設定
	if excludePatterns := c.String("exclude"); excludePatterns != "" {
		patterns := strings.Split(excludePatterns, ",")
		for i, pattern := range patterns {
			patterns[i] = strings.TrimSpace(pattern)
		}
		processor.SetExcludePatterns(patterns)
	}

	verbose := c.Bool("verbose")
	showStats := c.Bool("stats")

	// 詳細情報の表示
	if verbose {
		fmt.Printf("入力ディレクトリ: %s\n", inputDir)
		if outputDir := c.String("output"); outputDir != "" {
			fmt.Printf("出力ディレクトリ: %s\n", outputDir)
		} else {
			fmt.Println("出力ディレクトリ: 各ファイルと同じディレクトリ")
		}
		fmt.Printf("圧縮品質: %d\n", options.Quality)
		fmt.Printf("PNGパレットサイズ: %d\n", options.PaletteSize)
		fmt.Printf("並行ワーカー数: %d\n", c.Int("workers"))
		fmt.Printf("再帰処理: %s\n", boolToString(c.Bool("recursive")))
		fmt.Printf("包含パターン: %s\n", c.String("include"))
		if excludePatterns := c.String("exclude"); excludePatterns != "" {
			fmt.Printf("除外パターン: %s\n", excludePatterns)
		}
		fmt.Println()
	}

	fmt.Println("バッチ圧縮を開始しています...")

	// バッチ処理の実行
	results, err := processor.ProcessDirectory(inputDir, options)
	if err != nil {
		return cli.Exit(fmt.Sprintf("バッチ処理エラー: %v", err), 1)
	}

	// 結果の表示
	if len(results) == 0 {
		fmt.Println("処理対象のファイルが見つかりませんでした。")
		return nil
	}

	// 統計情報の計算
	stats := batch.CalculateStatistics(results)

	// 詳細表示の場合は各ファイルの結果を表示
	if verbose {
		fmt.Println("\n=== 処理結果詳細 ===")
		for _, result := range results {
			if result.Error != nil {
				fmt.Printf("❌ %s: %v\n", result.Job.InputPath, result.Error)
			} else {
				compressionRatio := 0.0
				if result.OriginalSize > 0 {
					compressionRatio = 100.0 - (float64(result.CompressedSize) / float64(result.OriginalSize) * 100.0)
				}
				fmt.Printf("✅ %s → %s (%.2f%% 圧縮)\n",
					result.Job.InputPath,
					result.Job.OutputPath,
					compressionRatio)
			}
		}
		fmt.Println()
	}

	// 統計情報の表示
	if showStats || verbose {
		fmt.Println("=== 圧縮統計 ===")
		fmt.Printf("処理ファイル数: %d\n", stats.TotalFiles)
		fmt.Printf("成功: %d, 失敗: %d\n", stats.SuccessFiles, stats.FailedFiles)

		if stats.SuccessFiles > 0 {
			fmt.Printf("元のサイズ合計: %s\n", formatFileSize(stats.TotalOriginalSize))
			fmt.Printf("圧縮後サイズ合計: %s\n", formatFileSize(stats.TotalCompressedSize))
			fmt.Printf("全体圧縮率: %.2f%%\n", stats.CompressionRatio)
		}
	} else {
		// 簡潔な結果表示
		fmt.Printf("バッチ圧縮が完了しました！\n")
		fmt.Printf("処理ファイル数: %d (成功: %d, 失敗: %d)\n", stats.TotalFiles, stats.SuccessFiles, stats.FailedFiles)
		if stats.SuccessFiles > 0 {
			fmt.Printf("全体圧縮率: %.2f%%\n", stats.CompressionRatio)
		}
	}

	// エラーがあった場合の終了コード
	if stats.FailedFiles > 0 {
		fmt.Printf("\n⚠️  %d個のファイルでエラーが発生しました。詳細は --verbose オプションで確認してください。\n", stats.FailedFiles)
		return cli.Exit("", 1)
	}

	return nil
}

// boolToString converts bool to Japanese string
func boolToString(b bool) string {
	if b {
		return "有効"
	}
	return "無効"
}

// formatFileSize formats file size in human readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
