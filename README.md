# Kirby-Inspired RPG 🌟

星のカービィにインスパイアされた2D RPGゲームです。Go言語とPixelライブラリを使用して開発されています。
<img width="849" height="491" alt="image" src="https://github.com/user-attachments/assets/b25f8ea4-056d-46c0-bb40-d1ed740ee22e" />

![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)

> **🪟 WSL2ユーザーの方へ**: [WSL2クイックセットアップ](#wsl2での環境構築windowsユーザー向け)を参照してください

---

**📖 [クイックスタート](QUICKSTART.md)** | **🛠️ [開発ガイド](docs/DEVELOPMENT.md)** | **🤝 [コントリビューション](CONTRIBUTING.md)**

---

![Kirby Character](https://user-images.githubusercontent.com/placeholder/kirby-character.png)

## ✨ 特徴

- **カービィ風のキャラクター**: ピンクの丸いキャラクターを操作
- **コピー能力システム**: 敵を倒して能力をコピー
- **3種類の敵キャラクター**: 歩行型、飛行型、ジャンプ型
- **プラットフォームアクション**: ジャンプして移動、敵を踏みつけて倒す
- **体力システム**: HPバーとダメージ、回復システム
- **スコアシステム**: 敵を倒してハイスコアを目指そう

## 🎮 ゲームの特徴

### プレイヤーキャラクター
- ダブルジャンプ可能
- 敵を上から踏んで倒せる
- コピー能力を使用できる
- 無敵時間あり（ダメージ後）

### コピー能力
1. **スピード能力** (黄色)
   - 移動速度が1.5倍に
   - 歩行型の敵から取得

2. **飛行能力** (紫色)
   - 空中を飛べる
   - 飛行型の敵から取得

3. **ハイジャンプ能力** (黄色)
   - 通常より高くジャンプできる
   - ジャンプ型の敵から取得

### 敵キャラクター
- **歩行型** (緑色): 地面を左右にパトロール
- **飛行型** (紫色): 空中を飛びながらプレイヤーを追跡
- **ジャンプ型** (黄色): ジャンプして移動

## 🚀 必要な環境

- Go 1.21以上
- OpenGL対応のグラフィックカード
- Linux/macOS/Windows

### 依存ライブラリ
- [Pixel](https://github.com/faiface/pixel) - 2Dゲーム開発ライブラリ
- OpenGL

## 📦 インストール

### WSL2での環境構築（Windowsユーザー向け）

WSL2でゲームを実行する場合、X11ライブラリとグラフィックスライブラリが必要です。

#### 1. 必要なライブラリをインストール

```bash
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
```

#### 2. X11サーバーの設定

**Windows 11の場合（WSLg）:**
- Windows 11のWSL2には、WSLg（WSL GUI）が標準で含まれています
- 追加設定なしでGUIアプリケーションが動作します

**Windows 10の場合:**
1. [VcXsrv](https://sourceforge.net/projects/vcxsrv/)をダウンロードしてインストール
2. XLaunchを起動し、以下の設定で実行:
   - Display number: 0
   - Start no client: チェック
   - Disable access control: チェック
3. WSLで環境変数を設定:

```bash
export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}'):0
```

永続化する場合は`~/.bashrc`に追加:

```bash
echo 'export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '"'"'{print $2}'"'"'):0' >> ~/.bashrc
source ~/.bashrc
```

### 通常のインストール手順

### 1. リポジトリをクローン

```bash
git clone https://github.com/remmakoshino/kirby-inspired-go.git
cd kirby-inspired-go
```

### 2. 依存関係をインストール

```bash
go mod download
```

### 3. ゲームをビルド

```bash
# Makefileを使用（推奨）
make build

# または直接ビルド
go build -o bin/kirby-game ./cmd/game
```

### 4. ゲームを実行

```bash
# Makefileを使用
make run

# またはビルド済みバイナリを実行
./bin/kirby-game

# または直接実行（ビルドなし）
go run ./cmd/game
```

**⚠️ WSL2で実行する場合の注意:**
- X11サーバーが起動していることを確認してください
- `echo $DISPLAY`でDISPLAY環境変数が設定されているか確認
- ゲームウィンドウが表示されない場合は、上記のX11サーバー設定を確認してください

## 🎯 操作方法

### キーボード操作

- **移動**: 矢印キー または A/D
- **ジャンプ**: スペースキー または W
- **攻撃/能力使用**: X または J
- **ゲームオーバー後リスタート**: R

### ゲームのコツ

1. **敵は上から踏んで倒そう**: 横や下から当たるとダメージを受けます
2. **能力を使いこなそう**: 敵を倒すとその敵の能力をコピーできます
3. **プラットフォームを活用**: 高い場所から攻撃すると有利です
4. **連続で敵を倒してスコアアップ**: 全ての敵を倒すと新しい波が来ます

## 📁 プロジェクト構成

```
kirby-inspired-go/
├── cmd/
│   └── game/           # メインエントリーポイント
│       └── main.go
├── internal/
│   ├── entity/         # エンティティ（プレイヤー、敵）
│   │   ├── player.go
│   │   └── enemy.go
│   ├── ability/        # コピー能力システム
│   │   ├── ability.go
│   │   └── abilities.go
│   ├── stage/          # ステージとプラットフォーム
│   │   └── stage.go
│   └── game/           # ゲームメインロジック
│       └── game.go
├── assets/             # ゲームアセット（将来使用）
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## 🛠️ 技術スタック

- **言語**: Go 1.21
- **ゲームエンジン**: Pixel (OpenGLベース)
- **グラフィックス**: imdraw (即座図形描画)
- **物理**: カスタム実装（重力、衝突判定）

## 🎨 実装された機能

- ✅ プレイヤーキャラクター（移動、ジャンプ、ダブルジャンプ）
- ✅ カービィ風のキャラクターデザイン（ピンクの丸いキャラ）
- ✅ 3種類の敵キャラクター（AI付き）
- ✅ コピー能力システム（3種類の能力）
- ✅ ステージとプラットフォーム
- ✅ 衝突判定システム
- ✅ 体力システム（HP、ダメージ、無敵時間）
- ✅ スコアシステム
- ✅ ゲームオーバーとリスタート
- ✅ UI（HPバー、スコア、能力表示）

## 🔮 今後の拡張予定

- [ ] 音楽と効果音
- [ ] より多くのステージ
- [ ] ボスキャラクター
- [ ] パワーアップアイテム
- [ ] セーブ/ロード機能
- [ ] マルチプレイヤー対応
- [ ] スプライトアニメーション

## 📝 開発メモ

### WSL2でのビルドとトラブルシューティング

**必要なパッケージ（WSL2/Linux）:**
```bash
sudo apt-get install -y \
    libgl1-mesa-dev \
    xorg-dev \
    libx11-dev \
    libxrandr-dev \
    libxcursor-dev \
    libxinerama-dev \
    libxi-dev \
    pkg-config
```

**X11サーバーの確認:**
```bash
# DISPLAY環境変数の確認
echo $DISPLAY

# X11サーバーが起動しているか確認
ps aux | grep X
```

**ゲームウィンドウが表示されない場合:**
1. X11サーバー（VcXsrvまたはWSLg）が起動しているか確認
2. DISPLAY環境変数が設定されているか確認
3. ファイアウォールでX11通信が許可されているか確認

### ビルドのトラブルシューティング

**OpenGLエラーが出る場合（Linux）:**
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**OpenGLエラーが出る場合（macOS）:**
```bash
# Xcodeコマンドラインツールをインストール
xcode-select --install
```

**Windowsの場合:**
- MinGWまたはTDM-GCCが必要です
- CGOが有効になっていることを確認してください

## 🤝 コントリビューション

プルリクエストを歓迎します！大きな変更の場合は、まずissueを開いて変更内容を議論してください。

1. フォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを開く

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 👏 謝辞

- 星のカービィシリーズ（ゲームデザインのインスピレーション）
- [Pixel](https://github.com/faiface/pixel) - 素晴らしい2Dゲームライブラリ
- Go コミュニティ

## 📧 連絡先

プロジェクトリンク: [https://github.com/remmakoshino/kirby-inspired-go](https://github.com/remmakoshino/kirby-inspired-go)

---

**楽しんでプレイしてください！ 🌟**
