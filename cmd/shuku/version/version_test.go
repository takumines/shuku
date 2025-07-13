package version

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestVersionCmd(t *testing.T) {
	// テスト用の値を設定
	originalVersion := Version
	originalCommit := Commit
	originalDate := Date

	defer func() {
		Version = originalVersion
		Commit = originalCommit
		Date = originalDate
	}()

	Version = "1.0.0"
	Commit = "abc123"
	Date = "2023-01-01"

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
		wantErr        bool
	}{
		{
			name: "version command",
			args: []string{"version"},
			expectedOutput: []string{
				"shuku version 1.0.0",
				"Commit: abc123",
				"Built: 2023-01-01",
			},
			wantErr: false,
		},
		{
			name: "version alias v",
			args: []string{"v"},
			expectedOutput: []string{
				"shuku version 1.0.0",
				"Commit: abc123",
				"Built: 2023-01-01",
			},
			wantErr: false,
		},
		{
			name: "dev version",
			args: []string{"version"},
			expectedOutput: []string{
				"shuku version dev",
				"Commit: unknown",
				"Built: unknown",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// dev版のテストケース用に値を再設定
			if tt.name == "dev version" {
				Version = "dev"
				Commit = "unknown"
				Date = "unknown"
			} else {
				Version = "1.0.0"
				Commit = "abc123"
				Date = "2023-01-01"
			}

			// 標準出力をキャプチャ
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// CLIアプリケーションを作成
			app := &cli.App{
				Name:     "shuku",
				Commands: []*cli.Command{Cmd()},
			}

			// コマンドを実行
			err := app.Run(append([]string{"shuku"}, tt.args...))

			// 標準出力を元に戻し、出力を読み取り
			w.Close()
			os.Stdout = old
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// エラーチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("version command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 出力チェック
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("version command output = %v, expected to contain %v", output, expected)
				}
			}
		})
	}
}

func TestVersionVariables(t *testing.T) {
	// バージョン変数の初期値をテスト
	if Version == "" {
		t.Error("Version should not be empty")
	}

	if Commit == "" {
		t.Error("Commit should not be empty")
	}

	if Date == "" {
		t.Error("Date should not be empty")
	}
}

func TestVersionCmdStructure(t *testing.T) {
	cmd := Cmd()

	// コマンド基本情報のテスト
	if cmd.Name != "version" {
		t.Errorf("Command name = %v, want %v", cmd.Name, "version")
	}

	if cmd.Usage != "Print the version" {
		t.Errorf("Command usage = %v, want %v", cmd.Usage, "Print the version")
	}

	// エイリアスのテスト
	expectedAliases := []string{"v"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Command aliases length = %v, want %v", len(cmd.Aliases), len(expectedAliases))
	}

	for i, alias := range expectedAliases {
		if i >= len(cmd.Aliases) || cmd.Aliases[i] != alias {
			t.Errorf("Command alias[%d] = %v, want %v", i, cmd.Aliases[i], alias)
		}
	}

	// Actionが設定されているかテスト
	if cmd.Action == nil {
		t.Error("Command action should not be nil")
	}
}

func TestVersionAction(t *testing.T) {
	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// テスト用の値を設定
	Version = "test-version"
	Commit = "test-commit"
	Date = "test-date"

	// Actionを直接実行
	cmd := Cmd()
	err := cmd.Action(&cli.Context{})

	// 標準出力を元に戻し、出力を読み取り
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// エラーチェック
	if err != nil {
		t.Errorf("Action returned error: %v", err)
	}

	// 出力チェック
	expectedLines := []string{
		fmt.Sprintf("shuku version %s", Version),
		fmt.Sprintf("Commit: %s", Commit),
		fmt.Sprintf("Built: %s", Date),
	}

	for _, expected := range expectedLines {
		if !strings.Contains(output, expected) {
			t.Errorf("Action output = %v, expected to contain %v", output, expected)
		}
	}
}
