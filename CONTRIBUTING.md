# コントリビューションガイド

Kirby-Inspired RPGへのコントリビューションに興味を持っていただきありがとうございます！

## 開発環境のセットアップ

1. **リポジトリのフォーク**
   - GitHubでこのリポジトリをフォークしてください

2. **ローカルにクローン**
   ```bash
   git clone https://github.com/YOUR_USERNAME/kirby-inspired-go.git
   cd kirby-inspired-go
   ```

3. **依存関係のインストール**
   ```bash
   make deps
   # または
   go mod download
   ```

4. **ゲームの実行確認**
   ```bash
   make run
   # または
   go run ./cmd/game
   ```

## 開発フロー

1. **ブランチの作成**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **コードの変更**
   - 適切なディレクトリ構造に従ってコードを記述
   - コメントは日本語または英語で記述

3. **フォーマット**
   ```bash
   make fmt
   ```

4. **テストの実行**
   ```bash
   make test
   ```

5. **コミット**
   ```bash
   git add .
   git commit -m "feat: 新機能の説明"
   ```

## コミットメッセージのガイドライン

コミットメッセージは以下の形式で記述してください：

```
<type>: <subject>

<body>
```

### Type
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響を与えない変更（空白、フォーマットなど）
- `refactor`: バグ修正や機能追加を含まないコード変更
- `test`: テストの追加や修正
- `chore`: ビルドプロセスやツールの変更

### 例
```
feat: プレイヤーに新しい攻撃アニメーションを追加

ダッシュ攻撃のアニメーションを実装しました。
新しいスプライトシートを使用しています。
```

## プルリクエストのガイドライン

1. **わかりやすいタイトル**
   - 変更内容を簡潔に説明

2. **詳細な説明**
   - 何を変更したか
   - なぜ変更したか
   - どのようにテストしたか

3. **スクリーンショット**
   - UI変更の場合は必須

4. **関連するIssue**
   - 関連するIssueがあれば参照

## コーディング規約

### Go言語の規約
- [Effective Go](https://golang.org/doc/effective_go.html)に従う
- `go fmt`でフォーマット
- 適切なエラーハンドリング
- パブリックな関数・型にはコメントを記述

### プロジェクト固有の規約

1. **ディレクトリ構造**
   ```
   internal/
   ├── entity/    # ゲームエンティティ
   ├── ability/   # コピー能力
   ├── stage/     # ステージとレベル
   └── game/      # ゲームメインロジック
   ```

2. **命名規則**
   - 構造体: PascalCase (例: `Player`, `EnemyType`)
   - 関数: camelCase (例: `updatePosition`, `checkCollision`)
   - 定数: UPPER_SNAKE_CASE (例: `MAX_HEALTH`, `JUMP_FORCE`)

3. **コメント**
   - 日本語または英語
   - 複雑なロジックには必ず説明を追加

## 機能追加のアイデア

以下のような機能追加を歓迎します：

- [ ] 新しいコピー能力
- [ ] 新しい敵タイプ
- [ ] 新しいステージ
- [ ] ボスキャラクター
- [ ] パワーアップアイテム
- [ ] サウンドエフェクト
- [ ] BGM
- [ ] スプライトアニメーション

## 質問やサポート

- **Issue**: バグ報告や機能リクエストはGitHub Issueで
- **Discussion**: 一般的な質問や議論はGitHub Discussionsで

## ライセンス

コントリビューションは[MIT License](LICENSE)の下で公開されます。
