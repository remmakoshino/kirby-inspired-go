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

#### WSL2 / Linux (Ubuntu/Debian)

```bash
# 必要なライブラリをインストール
sudo apt-get update
sudo apt-get install -y \
    libgl1-mesa-dev \
    xorg-dev \
    libx11-dev \
    libxrandr-dev \
    libxcursor-dev \
    libxinerama-dev \
    libxi-dev \
    pkg-config

# Go依存関係をダウンロード
go mod download
```

**WSL2でのX11設定:**

- **Windows 11**: WSLgが自動的に有効（追加設定不要）
- **Windows 10**: VcXsrvをインストールし、以下を実行:

```bash
export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}'):0
```

#### Linux (Fedora/RHEL)

```bash
sudo dnf install -y \
    golang \
    mesa-libGL-devel \
    libX11-devel \
    libXrandr-devel \
    libXcursor-devel \
    libXinerama-devel \
    libXi-devel

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

### メニュー操作

| キー | 動作 |
|------|------|
| ↑ ↓ (または W/S) | メニュー選択 |
| Enter | 決定 |
| Esc | 戻る |

### カービィの操作

| キー | 動作 |
|------|------|
| ← → (または A/D) | 左右移動 |
| スペース (または W) | ジャンプ（2回でダブルジャンプ） |
| X (または J) | 攻撃 |
| Z (または K) | アビリティ使用 |
| R | メニューに戻る（ゲームオーバー/クリア時） |

### メタナイトの操作

| キー | 動作 |
|------|------|
| ← → (または A/D) | 左右移動 |
| スペース (または W) | ジャンプ（2回でダブルジャンプ） |
| E | 攻撃（連続入力でコンボ） |
| Q | 特殊技 |
| 1 | 剣に切り替え |
| 2 | トルネードに切り替え |
| 3 | マント防御に切り替え |
| R | メニューに戻る（ゲームオーバー/クリア時） |

### ゲームの流れ

1. **タイトル画面**: START GAME を選択
2. **キャラクター選択**: カービィまたはメタナイトを選択
3. **ステージ選択**: Stage 1 (デデデ大王) または Stage 2 (メタナイト) を選択
4. **ゲームプレイ**: 敵を倒し、ボスを撃破
5. **クリア/ゲームオーバー**: Rキーでメニューに戻る

### ゲームのルール

1. **敵を倒す**: 敵の上から踏みつけて倒す
2. **能力をコピー** (カービィのみ): 敵を倒すとその能力を獲得
3. **ボスを倒す**: 各ステージのボスを撃破してクリア
4. **スコアを稼ぐ**: 敵を倒してハイスコアを目指す
5. **体力管理**: HPが0になるとゲームオーバー

### 敵の種類

#### 通常敵
- **緑の敵**: 地上を歩く（スピード能力）
- **紫の敵**: 空を飛ぶ（飛行能力）
- **黄色の敵**: ジャンプする（ハイジャンプ能力）

#### 新敵
- **ワドルディ** (オレンジの丸): パトロール型、体力20
- **ワドルドゥ** (単眼): 射撃攻撃可能、体力25

#### ボス
- **デデデ大王** (ステージ1): ハンマー/ジャンプ/突進、体力200
- **メタナイト** (ステージ2): 剣/トルネード/ダッシュ/防御、体力150

### コピー能力（カービィ専用）

- **吸い込み**: Zキーで敵を吸い込み、能力をコピー
- **ハンマー**: 強力な近接攻撃
- **剣**: 連続攻撃可能
- **トルネード**: 高速突進
- **マント防御**: 一時的に防御力アップ

## トラブルシューティング 🔧

### ビルドエラー

**エラー: "Package gl was not found"**

```bash
# Linux/WSL
sudo apt-get install pkg-config libgl1-mesa-dev xorg-dev

# macOS
xcode-select --install
```

**エラー: "X11/Xlib.h: No such file"**

```bash
# Linux/WSL
sudo apt-get install libx11-dev xorg-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev
```

### 実行エラー

**WSLでウィンドウが表示されない**

1. X11サーバーの確認:
```bash
# DISPLAY環境変数の確認
echo $DISPLAY

# 表示されない場合は設定
export DISPLAY=:0
```

2. **Windows 11**: WSLgが有効か確認
   - Windows Update で最新版に更新

3. **Windows 10**: VcXsrvが起動しているか確認
   - XLaunchを実行
   - "Disable access control"にチェック

**パーミッションエラー**

```bash
chmod +x bin/kirby-game
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
