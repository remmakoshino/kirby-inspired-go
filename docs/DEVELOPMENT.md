# 開発環境セットアップガイド

## WSL (Windows Subsystem for Linux)での開発

WSL環境では、OpenGLとX11ライブラリが必要です。

### 必要なパッケージのインストール

```bash
sudo apt-get update
sudo apt-get install -y \
    libgl1-mesa-dev \
    xorg-dev \
    libx11-dev \
    libxrandr-dev \
    libxcursor-dev \
    libxinerama-dev \
    libxi-dev
```

### X11サーバーの設定

WSLでGUIアプリケーションを実行するには、X11サーバーが必要です。

#### オプション1: WSLg（Windows 11）

Windows 11のWSL2では、WSLgが自動的に有効になっています。

```bash
# 環境変数の確認
echo $DISPLAY
# 出力例: :0
```

#### オプション2: VcXsrv（Windows 10）

1. [VcXsrv](https://sourceforge.net/projects/vcxsrv/)をダウンロードしてインストール
2. XLaunchを起動
3. WSLで環境変数を設定:

```bash
export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}'):0
```

`.bashrc`に追加して永続化:

```bash
echo 'export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '"'"'{print $2}'"'"'):0' >> ~/.bashrc
source ~/.bashrc
```

### ビルドの実行

```bash
cd /mnt/c/Users/renma/kirby-inspired-go
go build -o bin/kirby-game ./cmd/game
```

### ゲームの実行

```bash
./bin/kirby-game
```

## Linux（ネイティブ）での開発

### Ubuntu/Debian

```bash
sudo apt-get update
sudo apt-get install -y \
    golang \
    libgl1-mesa-dev \
    xorg-dev
```

### Fedora/RHEL

```bash
sudo dnf install -y \
    golang \
    mesa-libGL-devel \
    libX11-devel \
    libXrandr-devel \
    libXcursor-devel \
    libXinerama-devel \
    libXi-devel
```

## macOSでの開発

### Homebrewでの環境構築

```bash
# Homebrewのインストール（未インストールの場合）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Goのインストール
brew install go

# Xcodeコマンドラインツールのインストール
xcode-select --install
```

### ビルドと実行

```bash
cd ~/kirby-inspired-go
go build -o bin/kirby-game ./cmd/game
./bin/kirby-game
```

## Windowsでの開発（ネイティブ）

### 必要なツール

1. **Go**: [公式サイト](https://golang.org/dl/)からインストーラーをダウンロード
2. **MinGW-w64**: C/C++コンパイラ
   - [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)をインストール
   - または[MinGW-w64](https://www.mingw-w64.org/)

### 環境変数の設定

```powershell
# PowerShellで確認
$env:CGO_ENABLED = "1"
go env CGO_ENABLED
```

### ビルド

```powershell
cd C:\Users\renma\kirby-inspired-go
go build -o bin\kirby-game.exe .\cmd\game
```

### 実行

```powershell
.\bin\kirby-game.exe
```

## トラブルシューティング

### エラー: "Package gl was not found"

```bash
# pkg-configのインストール
sudo apt-get install pkg-config

# OpenGLライブラリの再インストール
sudo apt-get install --reinstall libgl1-mesa-dev
```

### エラー: "X11/Xlib.h: No such file or directory"

```bash
# X11開発ライブラリのインストール
sudo apt-get install libx11-dev xorg-dev
```

### WSLでゲームウィンドウが表示されない

```bash
# DISPLAYの確認
echo $DISPLAY

# X11サーバーが起動しているか確認
ps aux | grep X

# 再設定
export DISPLAY=:0
```

### ビルドは成功するが実行できない

```bash
# 実行権限の確認
chmod +x bin/kirby-game

# 依存ライブラリの確認
ldd bin/kirby-game
```

## 開発ツール

### VS Code設定

`.vscode/settings.json`:

```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.formatTool": "gofmt",
    "go.testFlags": ["-v"],
    "files.exclude": {
        "**/.git": true,
        "**/bin": true
    }
}
```

### 推奨拡張機能

- Go (golang.go)
- GitLens
- Error Lens

## パフォーマンス最適化

### ビルドフラグ

```bash
# リリースビルド（最適化）
go build -ldflags="-s -w" -o bin/kirby-game ./cmd/game

# デバッグビルド
go build -gcflags="all=-N -l" -o bin/kirby-game-debug ./cmd/game
```

### プロファイリング

```bash
# CPUプロファイリング
go build -o bin/kirby-game ./cmd/game
./bin/kirby-game -cpuprofile=cpu.prof

# メモリプロファイリング
./bin/kirby-game -memprofile=mem.prof
```

## 次のステップ

環境構築が完了したら：

1. `make run` でゲームを起動
2. 操作方法を確認（README.md参照）
3. コードを変更して独自の機能を追加
4. `CONTRIBUTING.md`を読んでコントリビューション

楽しい開発を！🎮
