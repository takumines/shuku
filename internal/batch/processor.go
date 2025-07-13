// Package batch はディレクトリ内の複数画像ファイルの一括圧縮機能を提供します。
package batch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/takumines/shuku/pkg/shuku"
)

// Job は単一の画像圧縮ジョブを表します。
type Job struct {
	InputPath  string
	OutputPath string
	Options    shuku.Options
}

// Result はジョブの実行結果を表します。
type Result struct {
	Job            Job
	OriginalSize   int64
	CompressedSize int64
	Error          error
}

// Processor はバッチ処理を管理します。
type Processor struct {
	WorkerCount  int      // 並行処理数
	OutputDir    string   // 出力ディレクトリ
	Recursive    bool     // 再帰的処理フラグ
	IncludeGlobs []string // 処理対象ファイルパターン
	ExcludeGlobs []string // 除外ファイルパターン
}

// NewProcessor は新しいProcessorインスタンスを作成します。
func NewProcessor(workerCount int, outputDir string) *Processor {
	if workerCount <= 0 {
		workerCount = 4 // デフォルトワーカー数
	}

	return &Processor{
		WorkerCount:  workerCount,
		OutputDir:    outputDir,
		Recursive:    false,
		IncludeGlobs: []string{"*.jpg", "*.jpeg", "*.png", "*.webp"},
		ExcludeGlobs: []string{},
	}
}

// SetRecursive は再帰的処理を有効または無効にします。
func (p *Processor) SetRecursive(recursive bool) {
	p.Recursive = recursive
}

// SetIncludePatterns は処理対象ファイルパターンを設定します。
func (p *Processor) SetIncludePatterns(patterns []string) {
	p.IncludeGlobs = patterns
}

// SetExcludePatterns は除外ファイルパターンを設定します。
func (p *Processor) SetExcludePatterns(patterns []string) {
	p.ExcludeGlobs = patterns
}

// ProcessDirectory はディレクトリ内の画像ファイルを一括圧縮します。
func (p *Processor) ProcessDirectory(inputDir string, options shuku.Options) ([]Result, error) {
	// 入力ディレクトリの存在確認
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("入力ディレクトリが存在しません: %s", inputDir)
	}

	// 出力ディレクトリの作成
	if p.OutputDir != "" {
		if err := os.MkdirAll(p.OutputDir, 0755); err != nil {
			return nil, fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
		}
	}

	// 処理対象ファイルを収集
	jobs, err := p.collectJobs(inputDir, options)
	if err != nil {
		return nil, fmt.Errorf("ファイル収集エラー: %v", err)
	}

	if len(jobs) == 0 {
		return []Result{}, nil
	}

	// 並行処理でジョブを実行
	return p.executeJobs(jobs), nil
}

// collectJobs は処理対象ファイルを収集してJobsを作成します。
func (p *Processor) collectJobs(inputDir string, options shuku.Options) ([]Job, error) {
	var jobs []Job

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリの場合
		if info.IsDir() {
			// ルートディレクトリではない かつ 再帰処理が無効の場合はスキップ
			if path != inputDir && !p.Recursive {
				return filepath.SkipDir
			}
			return nil
		}

		// ファイルのフィルタリング
		if p.shouldIncludeFile(path) {
			outputPath := p.generateOutputPath(path, inputDir)
			jobs = append(jobs, Job{
				InputPath:  path,
				OutputPath: outputPath,
				Options:    options,
			})
		}

		return nil
	}

	err := filepath.Walk(inputDir, walkFunc)
	return jobs, err
}

// shouldIncludeFile はファイルが処理対象かどうかを判定します。
func (p *Processor) shouldIncludeFile(filePath string) bool {
	fileName := filepath.Base(filePath)

	// 除外パターンのチェック
	for _, pattern := range p.ExcludeGlobs {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return false
		}
	}

	// 包含パターンのチェック
	for _, pattern := range p.IncludeGlobs {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}
	}

	return false
}

// generateOutputPath は出力ファイルパスを生成します。
func (p *Processor) generateOutputPath(inputPath, inputDir string) string {
	if p.OutputDir == "" {
		// 出力ディレクトリが指定されていない場合は、入力ファイルと同じディレクトリに "_compressed" を追加
		ext := filepath.Ext(inputPath)
		base := strings.TrimSuffix(inputPath, ext)
		return base + "_compressed" + ext
	}

	// 入力ディレクトリからの相対パスを計算
	relPath, err := filepath.Rel(inputDir, inputPath)
	if err != nil {
		// エラーの場合はファイル名のみを使用
		relPath = filepath.Base(inputPath)
	}

	return filepath.Join(p.OutputDir, relPath)
}

// executeJobs は並行処理でジョブを実行します。
func (p *Processor) executeJobs(jobs []Job) []Result {
	jobChan := make(chan Job, len(jobs))
	resultChan := make(chan Result, len(jobs))

	// ワーカーを起動
	var wg sync.WaitGroup
	for i := 0; i < p.WorkerCount; i++ {
		wg.Add(1)
		go p.worker(&wg, jobChan, resultChan)
	}

	// ジョブをチャネルに送信
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// ワーカーの完了を待機
	wg.Wait()
	close(resultChan)

	// 結果を収集
	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// worker は単一のワーカーゴルーチンを実装します。
func (p *Processor) worker(wg *sync.WaitGroup, jobChan <-chan Job, resultChan chan<- Result) {
	defer wg.Done()

	for job := range jobChan {
		result := p.processJob(job)
		resultChan <- result
	}
}

// processJob は単一のジョブを処理します。
func (p *Processor) processJob(job Job) Result {
	result := Result{Job: job}

	// 入力ファイルのサイズを取得
	if inputInfo, err := os.Stat(job.InputPath); err == nil {
		result.OriginalSize = inputInfo.Size()
	}

	// 出力ディレクトリを作成
	outputDir := filepath.Dir(job.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		result.Error = fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
		return result
	}

	// 圧縮処理を実行
	err := shuku.CompressFile(job.InputPath, job.OutputPath, job.Options)
	if err != nil {
		result.Error = fmt.Errorf("圧縮処理エラー: %v", err)
		return result
	}

	// 出力ファイルのサイズを取得
	if outputInfo, err := os.Stat(job.OutputPath); err == nil {
		result.CompressedSize = outputInfo.Size()
	}

	return result
}

// Statistics はバッチ処理の統計情報を表します。
type Statistics struct {
	TotalFiles          int
	SuccessFiles        int
	FailedFiles         int
	TotalOriginalSize   int64
	TotalCompressedSize int64
	CompressionRatio    float64
}

// CalculateStatistics は処理結果から統計情報を計算します。
func CalculateStatistics(results []Result) Statistics {
	stats := Statistics{
		TotalFiles: len(results),
	}

	for _, result := range results {
		if result.Error != nil {
			stats.FailedFiles++
		} else {
			stats.SuccessFiles++
			stats.TotalOriginalSize += result.OriginalSize
			stats.TotalCompressedSize += result.CompressedSize
		}
	}

	if stats.TotalOriginalSize > 0 {
		stats.CompressionRatio = 100.0 - (float64(stats.TotalCompressedSize) / float64(stats.TotalOriginalSize) * 100.0)
	}

	return stats
}
