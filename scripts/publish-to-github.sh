#!/bin/bash

# GitHub公開スクリプト
# このスクリプトは、プロジェクトをGitHubに公開するための手順を実行します

set -e  # エラーが発生したら即座に終了

echo "================================"
echo "Kirby-Inspired RPG - GitHub公開"
echo "================================"
echo ""

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 関数: 成功メッセージ
success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# 関数: 警告メッセージ
warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# 関数: エラーメッセージ
error() {
    echo -e "${RED}✗ $1${NC}"
}

# Step 1: Gitリポジトリの確認
echo "Step 1: Gitリポジトリの確認..."
if [ -d .git ]; then
    warning "既存のGitリポジトリが見つかりました"
    read -p "リポジトリを再初期化しますか？ (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf .git
        git init
        success "Gitリポジトリを再初期化しました"
    else
        success "既存のリポジトリを使用します"
    fi
else
    git init
    success "Gitリポジトリを初期化しました"
fi

# Step 2: リモートリポジトリの設定
echo ""
echo "Step 2: リモートリポジトリの設定..."
REMOTE_URL="https://github.com/remmakoshino/kirby-inspired-go.git"

if git remote | grep -q "origin"; then
    warning "既存のリモート 'origin' が見つかりました"
    CURRENT_URL=$(git remote get-url origin)
    echo "現在のURL: $CURRENT_URL"
    
    if [ "$CURRENT_URL" != "$REMOTE_URL" ]; then
        git remote set-url origin "$REMOTE_URL"
        success "リモートURLを更新しました: $REMOTE_URL"
    else
        success "リモートURLは既に設定されています"
    fi
else
    git remote add origin "$REMOTE_URL"
    success "リモートリポジトリを追加しました: $REMOTE_URL"
fi

# Step 3: ファイルのステージング
echo ""
echo "Step 3: ファイルをステージング..."
git add .
success "全てのファイルをステージングしました"

# Step 4: コミット
echo ""
echo "Step 4: コミット..."
if git diff --cached --quiet; then
    warning "コミットする変更がありません"
else
    git commit -m "feat: 初回コミット - カービィ風RPGの基本実装

- プレイヤーキャラクター（ダブルジャンプ、移動）
- 3種類の敵キャラクター（歩行、飛行、ジャンプ）
- コピー能力システム（スピード、飛行、ハイジャンプ）
- ステージとプラットフォーム
- 体力システムとUI
- スコアシステム
- ゲームオーバーとリスタート"
    success "変更をコミットしました"
fi

# Step 5: ブランチの設定
echo ""
echo "Step 5: メインブランチの設定..."
git branch -M main
success "ブランチを 'main' に設定しました"

# Step 6: GitHubへプッシュ
echo ""
echo "Step 6: GitHubへプッシュ..."
echo "注意: GitHubの認証情報が必要です"
echo ""

read -p "GitHubへプッシュしますか？ (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if git push -u origin main; then
        success "GitHubへのプッシュが完了しました！"
    else
        error "プッシュに失敗しました"
        echo ""
        echo "トラブルシューティング:"
        echo "1. GitHubリポジトリが作成されているか確認"
        echo "2. 認証情報が正しいか確認"
        echo "3. 以下のコマンドを手動で実行:"
        echo "   git push -u origin main"
        exit 1
    fi
else
    warning "プッシュをスキップしました"
    echo "後で以下のコマンドで手動プッシュできます:"
    echo "  git push -u origin main"
fi

# Step 7: 完了メッセージ
echo ""
echo "================================"
echo "🎉 公開準備が完了しました！"
echo "================================"
echo ""
echo "次のステップ:"
echo "1. GitHubでリポジトリを確認"
echo "   $REMOTE_URL"
echo ""
echo "2. リポジトリの説明を追加"
echo "   Settings → About → Description"
echo ""
echo "3. トピックを追加"
echo "   go, golang, game, 2d-game, pixel, kirby, rpg"
echo ""
echo "4. 初回リリースを作成"
echo "   git tag -a v1.0.0 -m 'リリース v1.0.0'"
echo "   git push origin v1.0.0"
echo ""
echo "5. READMEを確認して、必要に応じて更新"
echo ""

success "全ての手順が完了しました！"
