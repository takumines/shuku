.PHONY: fmt check vet test build clean lint

# フォーマット（すべての環境で統一）
fmt:
	gofmt -w .

# フォーマットチェック（CI用）
check-fmt:
	@# Windows環境では先にフォーマットを実行して行末文字を正規化
	@if command -v cmd >/dev/null 2>&1; then \
		echo "Windows環境を検出：フォーマットを正規化しています..."; \
		gofmt -w .; \
	fi
	@if [ $$(gofmt -l . | wc -l) -gt 0 ]; then \
		echo "The following files are not formatted correctly:"; \
		gofmt -l .; \
		echo "Please run 'make fmt' to fix formatting issues."; \
		exit 1; \
	fi

# Go vet
vet:
	go vet ./...

# go mod tidy
tidy:
	go mod tidy

# go mod tidyチェック（CI用）
check-tidy:
	go mod tidy
	@if [ -n "$$(git status --porcelain go.mod go.sum)" ]; then \
		echo "go mod tidy made changes. Please run 'make tidy' and commit the changes."; \
		git diff go.mod go.sum; \
		exit 1; \
	fi

# テスト実行
test:
	go test -v ./...

# CLI ビルド
build:
	@if [ "$$(uname)" = "Windows_NT" ] || [ "$${OS}" = "Windows_NT" ]; then \
		go build -o shuku.exe ./cmd/shuku; \
	else \
		go build -o shuku ./cmd/shuku; \
	fi

# CLI機能テスト
test-cli: build
	@if [ -f "./shuku.exe" ]; then \
		./shuku.exe help || true; \
		./shuku.exe compress --help; \
	else \
		./shuku help || true; \
		./shuku compress --help; \
	fi

# CIで実行する全チェック
ci: check-fmt vet check-tidy test build test-cli

# クリーンアップ
clean:
	@if [ -f "./shuku" ]; then rm ./shuku; fi
	@if [ -f "./shuku.exe" ]; then rm ./shuku.exe; fi
