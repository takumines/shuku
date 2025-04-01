package compress

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"shuku/pkg/shuku"

	"github.com/urfave/cli/v2"
)

// Cmd returns the compress command.
func Cmd() *cli.Command {
	return &cli.Command{
		Name:    "compress",
		Usage:   "Compress an image.",
		Aliases: []string{"c"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Input image file path",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Specify the output file name (optional, defaults to input file name + '_compressed').",
			},
			&cli.IntFlag{
				Name:    "quality",
				Aliases: []string{"q"},
				Usage:   "JPEG quality (0-100)",
				Value:   80,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Show detailed information",
				Value:   false,
			},
		},
		Action: compressAction,
	}
}

// compressAction is the action for the compress command.
func compressAction(c *cli.Context) error {
	// 入力ファイルパスを取得
	inputPath := c.String("input")
	if inputPath == "" {
		return cli.Exit("入力ファイルが指定されていません。--input または -i オプションで指定してください。", 1)
	}

	// ファイルが存在するか確認
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("入力ファイル '%s' が見つかりません。", inputPath), 1)
	}

	// 出力ファイルパスを取得または生成
	outputPath := c.String("output")
	if outputPath == "" {
		// デフォルトの出力ファイル名を生成
		ext := filepath.Ext(inputPath)
		baseName := strings.TrimSuffix(inputPath, ext)
		outputPath = baseName + "_compressed" + ext
	}

	// 圧縮オプションを設定
	options := shuku.Options{
		Quality:     c.Int("quality"),
		PaletteSize: 256, // PNGの場合に使用
	}

	// 詳細表示モードが有効な場合
	verbose := c.Bool("verbose")
	if verbose {
		fmt.Printf("入力ファイル: %s\n", inputPath)
		fmt.Printf("出力ファイル: %s\n", outputPath)
		fmt.Printf("圧縮品質: %d\n", options.Quality)
	}

	// ファイル拡張子から形式を判断
	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext != ".jpg" && ext != ".jpeg" {
		return cli.Exit(fmt.Sprintf("サポートされていない画像形式です: %s。現在はJPEG形式のみ対応しています。", ext), 1)
	}

	fmt.Println("画像を圧縮しています...")

	// 圧縮処理を実行
	err := shuku.CompressFile(inputPath, outputPath, options)
	if err != nil {
		return cli.Exit(fmt.Sprintf("圧縮エラー: %v", err), 1)
	}

	// 圧縮前後のファイルサイズを取得して表示
	if verbose {
		inputInfo, err := os.Stat(inputPath)
		if err == nil {
			outputInfo, err := os.Stat(outputPath)
			if err == nil {
				inputSize := inputInfo.Size()
				outputSize := outputInfo.Size()
				reduction := 100.0 - (float64(outputSize) / float64(inputSize) * 100.0)

				fmt.Printf("元のサイズ: %d バイト\n", inputSize)
				fmt.Printf("圧縮後のサイズ: %d バイト\n", outputSize)
				fmt.Printf("圧縮率: %.2f%%\n", reduction)
			}
		}
	}

	fmt.Println("圧縮が完了しました！")
	fmt.Printf("圧縮ファイルが保存されました: %s\n", outputPath)

	return nil
}
