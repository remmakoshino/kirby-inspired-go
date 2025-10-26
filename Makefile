.PHONY: build run clean test deps

# ビルド設定
BINARY_NAME=kirby-game
MAIN_PATH=./cmd/game
BUILD_DIR=./bin

# デフォルトターゲット
all: deps build

# 依存関係のインストール
deps:
	@echo "依存関係をインストール中..."
	go mod download
	go mod tidy

# ビルド
build:
	@echo "ゲームをビルド中..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "ビルド完了: $(BUILD_DIR)/$(BINARY_NAME)"

# 実行
run: build
	@echo "ゲームを起動中..."
	$(BUILD_DIR)/$(BINARY_NAME)

# 直接実行（ビルドなし）
run-direct:
	@echo "ゲームを直接起動中..."
	go run $(MAIN_PATH)

# クリーン
clean:
	@echo "ビルド成果物を削除中..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "クリーン完了"

# テスト
test:
	@echo "テストを実行中..."
	go test -v ./...

# フォーマット
fmt:
	@echo "コードをフォーマット中..."
	go fmt ./...

# リント
lint:
	@echo "リントを実行中..."
	golangci-lint run ./...

# ヘルプ
help:
	@echo "使用可能なコマンド:"
	@echo "  make deps       - 依存関係をインストール"
	@echo "  make build      - ゲームをビルド"
	@echo "  make run        - ゲームをビルドして実行"
	@echo "  make run-direct - ゲームを直接実行（ビルドなし）"
	@echo "  make clean      - ビルド成果物を削除"
	@echo "  make test       - テストを実行"
	@echo "  make fmt        - コードをフォーマット"
	@echo "  make lint       - リントを実行"
