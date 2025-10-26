#!/bin/bash

# WSL2環境セットアップスクリプト
# このスクリプトは、WSL2でゲームを実行するために必要なライブラリをインストールします

set -e

echo "================================"
echo "WSL2環境セットアップ"
echo "Kirby-Inspired RPG"
echo "================================"
echo ""

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 関数: 成功メッセージ
success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# 関数: 情報メッセージ
info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# 関数: 警告メッセージ
warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Step 1: パッケージリストの更新
info "Step 1: パッケージリストを更新中..."
sudo apt-get update
success "パッケージリスト更新完了"

echo ""

# Step 2: 必要なライブラリのインストール
info "Step 2: 必要なライブラリをインストール中..."
echo "以下のパッケージがインストールされます:"
echo "  - libgl1-mesa-dev (OpenGL)"
echo "  - xorg-dev (X11開発ライブラリ)"
echo "  - libx11-dev, libxrandr-dev, libxcursor-dev (X11関連)"
echo "  - libxinerama-dev, libxi-dev (X11拡張)"
echo "  - pkg-config (パッケージ設定ツール)"
echo ""

sudo apt-get install -y \
    libgl1-mesa-dev \
    xorg-dev \
    libx11-dev \
    libxrandr-dev \
    libxcursor-dev \
    libxinerama-dev \
    libxi-dev \
    pkg-config

success "ライブラリのインストール完了"

echo ""

# Step 3: Go依存関係のダウンロード
info "Step 3: Go依存関係をダウンロード中..."
go mod download
success "Go依存関係のダウンロード完了"

echo ""

# Step 4: X11サーバーの確認
info "Step 4: X11サーバーの設定を確認中..."

if [ -z "$DISPLAY" ]; then
    warning "DISPLAY環境変数が設定されていません"
    echo ""
    echo "Windowsのバージョンを確認してください："
    echo ""
    echo "【Windows 11の場合】"
    echo "  - WSLgが標準搭載されています"
    echo "  - 通常は追加設定不要です"
    echo ""
    echo "【Windows 10の場合】"
    echo "  1. VcXsrvをインストール"
    echo "     https://sourceforge.net/projects/vcxsrv/"
    echo ""
    echo "  2. XLaunchを起動（設定）"
    echo "     - Display number: 0"
    echo "     - Disable access control: チェック"
    echo ""
    echo "  3. 以下のコマンドを実行（またはWSLを再起動）"
    echo "     export DISPLAY=\$(cat /etc/resolv.conf | grep nameserver | awk '{print \$2}'):0"
    echo ""
    
    read -p "DISPLAY環境変数を今すぐ設定しますか？ (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        NAMESERVER=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}')
        export DISPLAY="${NAMESERVER}:0"
        echo "export DISPLAY=\"${NAMESERVER}:0\"" >> ~/.bashrc
        success "DISPLAY環境変数を設定しました: $DISPLAY"
        info "次回から自動的に設定されます"
    fi
else
    success "DISPLAY環境変数が設定されています: $DISPLAY"
fi

echo ""

# Step 5: ビルドテスト
info "Step 5: ビルドテストを実行中..."
if make build; then
    success "ビルドテスト成功！"
else
    warning "ビルドに失敗しました。エラーを確認してください。"
    exit 1
fi

echo ""
echo "================================"
echo "セットアップ完了！🎉"
echo "================================"
echo ""
echo "次のステップ:"
echo ""
echo "1. ゲームを起動:"
echo "   make run"
echo ""
echo "2. または直接実行:"
echo "   ./bin/kirby-game"
echo ""
echo "3. X11サーバーの確認:"
echo "   echo \$DISPLAY"
echo ""
echo "トラブルシューティング:"
echo "  - ウィンドウが表示されない場合:"
echo "    1. X11サーバー（VcXsrv/WSLg）が起動しているか確認"
echo "    2. DISPLAY環境変数を確認: echo \$DISPLAY"
echo "    3. ファイアウォール設定を確認"
echo ""
echo "詳細な情報:"
echo "  - README.md"
echo "  - docs/DEVELOPMENT.md"
echo ""
success "すべての準備が整いました！"
