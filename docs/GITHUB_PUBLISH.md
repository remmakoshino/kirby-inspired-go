# GitHub公開手順

このドキュメントでは、Kirby-Inspired RPGをGitHubに公開する手順を説明します。

## 前提条件

- GitHubアカウントを持っていること
- Gitがインストールされていること
- リポジトリ`kirby-inspired-go`が作成されていること

## 手順

### 1. Gitリポジトリの初期化

```bash
cd /mnt/c/Users/renma/kirby-inspired-go
git init
```

### 2. .gitignoreの確認

既に`.gitignore`ファイルが作成されているので、不要なファイルが除外されます。

### 3. ファイルをステージング

```bash
git add .
```

### 4. 初回コミット

```bash
git commit -m "feat: 初回コミット - カービィ風RPGの基本実装"
```

### 5. GitHubリポジトリと連携

```bash
git branch -M main
git remote add origin https://github.com/remmakoshino/kirby-inspired-go.git
```

### 6. GitHubにプッシュ

```bash
git push -u origin main
```

## リポジトリ設定（GitHub上で）

### 1. リポジトリの説明

```
星のカービィにインスパイアされた2D RPGゲーム（Go + Pixel）
```

### 2. トピックの追加

以下のトピックを追加することをおすすめします：
- `go`
- `golang`
- `game`
- `2d-game`
- `pixel`
- `opengl`
- `kirby`
- `rpg`
- `platformer`

### 3. GitHub Pagesの設定（オプション）

- Settings → Pages
- Source: Deploy from a branch
- Branch: main / docs

### 4. Issueテンプレートの設定

`.github/ISSUE_TEMPLATE/`ディレクトリを作成して、以下のテンプレートを追加できます：

#### Bug Report
```yaml
name: バグ報告
about: バグを報告する
title: '[BUG] '
labels: bug
assignees: ''
```

#### Feature Request
```yaml
name: 機能リクエスト
about: 新機能を提案する
title: '[FEATURE] '
labels: enhancement
assignees: ''
```

### 5. プロジェクトボードの作成（オプション）

- Projects → New project
- Board template: Basic kanban
- 進捗管理に使用

## リリースの作成

### 1. バージョンタグの作成

```bash
git tag -a v1.0.0 -m "リリース v1.0.0 - 初回リリース"
git push origin v1.0.0
```

### 2. GitHub Releaseの作成

- Releases → Create a new release
- Tag: v1.0.0
- Title: v1.0.0 - 初回リリース
- Description: リリースノートを記述

## README.mdの更新

GitHubバッジを追加することで、プロジェクトの品質を示せます：

```markdown
![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Build Status](https://github.com/remmakoshino/kirby-inspired-go/workflows/Build%20and%20Test/badge.svg)
```

## 継続的な開発

### ブランチ戦略

- `main`: 安定版
- `develop`: 開発版
- `feature/*`: 機能追加
- `bugfix/*`: バグ修正

### プルリクエストフロー

1. フィーチャーブランチを作成
   ```bash
   git checkout -b feature/new-ability
   ```

2. 変更をコミット
   ```bash
   git add .
   git commit -m "feat: 新しいコピー能力を追加"
   ```

3. プッシュ
   ```bash
   git push origin feature/new-ability
   ```

4. GitHub上でPull Requestを作成

## トラブルシューティング

### プッシュできない場合

```bash
git pull origin main --rebase
git push origin main
```

### リモートURLの確認

```bash
git remote -v
```

### リモートURLの変更

```bash
git remote set-url origin https://github.com/remmakoshino/kirby-inspired-go.git
```

## 次のステップ

1. ✅ リポジトリの作成
2. ✅ 初回コミット
3. ✅ GitHubにプッシュ
4. ⬜ README.mdにバッジを追加
5. ⬜ Issueテンプレートを作成
6. ⬜ 初回リリース(v1.0.0)を作成
7. ⬜ プロジェクトボードを設定
8. ⬜ GitHub Actionsの動作確認

完了です！
