# クイックスタートガイド 🚀

## 5分でゲームを起動！

### 前提条件

- Go 1.18以上がインストール済み
- Git がインストール済み

### Step 1: リポジトリをクローン

```bash
git clone https://github.com/remmakoshino/kirby-inspired-go.git
cd kirby-inspired-go
```

### Step 2: 依存関係のインストール

#### Linux (WSL含む)

```bash
# 必要なライブラリをインストール
sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev xorg-dev

# Go依存関係をダウンロード
go mod download
```

#### macOS

```bash
# Xcodeコマンドラインツールのインストール
xcode-select --install

# Go依存関係をダウンロード
go mod download
```

#### Windows

```bash
# TDM-GCCまたはMinGW-w64が必要です
# インストール後、以下を実行:
go mod download
```

### Step 3: ゲームを実行

#### 方法1: Makefileを使用（推奨）

```bash
make run
```

#### 方法2: 直接実行

```bash
go run ./cmd/game
```

#### 方法3: ビルドしてから実行

```bash
# ビルド
go build -o bin/kirby-game ./cmd/game

# 実行
./bin/kirby-game
```

## 操作方法 🎮

### 基本操作

| キー | 動作 |
|------|------|
| ← → (または A/D) | 左右移動 |
| スペース (または W) | ジャンプ |
| X (または J) | 攻撃/能力使用 |
| R | リスタート（ゲームオーバー時） |

### ゲームのルール

1. **敵を倒す**: 敵の上から踏みつけて倒す
2. **能力をコピー**: 敵を倒すとその能力を獲得
3. **スコアを稼ぐ**: 敵を倒してハイスコアを目指す
4. **体力管理**: HPが0になるとゲームオーバー

### 敵の種類

- **緑の敵**: 地上を歩く（スピード能力）
- **紫の敵**: 空を飛ぶ（飛行能力）
- **黄色の敵**: ジャンプする（ハイジャンプ能力）

## トラブルシューティング 🔧

### ビルドエラー

**エラー: "Package gl was not found"**

```bash
# Linux/WSL
sudo apt-get install pkg-config libgl1-mesa-dev

# macOS
xcode-select --install
```

**エラー: "X11/Xlib.h: No such file"**

```bash
sudo apt-get install libx11-dev xorg-dev
```

### 実行エラー

**WSLでウィンドウが表示されない**

```bash
# X11サーバーの確認
echo $DISPLAY

# 設定されていない場合
export DISPLAY=:0
```

## 次のステップ 📚

1. **コードを読む**: `internal/` ディレクトリからスタート
2. **カスタマイズ**: キャラクターの色や能力を変更
3. **新機能を追加**: 新しい敵や能力を実装
4. **コントリビュート**: [CONTRIBUTING.md](CONTRIBUTING.md) を参照

## 詳細なドキュメント

- [開発環境セットアップ](docs/DEVELOPMENT.md)
- [GitHub公開手順](docs/GITHUB_PUBLISH.md)
- [コントリビューションガイド](CONTRIBUTING.md)

## ヘルプ

問題が解決しない場合:

1. [GitHub Issues](https://github.com/remmakoshino/kirby-inspired-go/issues) で検索
2. 新しいIssueを作成
3. [開発ドキュメント](docs/DEVELOPMENT.md) を参照

---

**楽しんでください！ 🌟**
