package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/remmakoshino/kirby-inspired-go/internal/ability"
	"github.com/remmakoshino/kirby-inspired-go/internal/entity"
	"github.com/remmakoshino/kirby-inspired-go/internal/stage"
)

const (
	WindowWidth  = 1024
	WindowHeight = 768
)

// Game はゲーム全体を管理します
type Game struct {
	Window   *pixelgl.Window
	Player   *entity.Player
	Enemies  []*entity.Enemy
	Stage    *stage.Stage
	IMDraw   *imdraw.IMDraw
	Score    int
	GameOver bool
	
	// UI関連
	Atlas *text.Atlas
}

// NewGame は新しいゲームを作成します
func NewGame(win *pixelgl.Window) *Game {
	rand.Seed(time.Now().UnixNano())
	
	// ステージ作成
	stg := stage.CreateDefaultStage(WindowWidth, WindowHeight)
	
	// プレイヤー作成
	player := entity.NewPlayer(pixel.V(WindowWidth/2, 200))
	
	// 敵を配置
	enemies := []*entity.Enemy{
		entity.NewEnemy(pixel.V(200, 150), entity.EnemyTypeWalker),
		entity.NewEnemy(pixel.V(500, 200), entity.EnemyTypeFlyer),
		entity.NewEnemy(pixel.V(800, 150), entity.EnemyTypeJumper),
		entity.NewEnemy(pixel.V(300, 300), entity.EnemyTypeWalker),
		entity.NewEnemy(pixel.V(650, 350), entity.EnemyTypeFlyer),
	}
	
	// テキスト描画用のアトラス
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	
	return &Game{
		Window:   win,
		Player:   player,
		Enemies:  enemies,
		Stage:    stg,
		IMDraw:   imdraw.New(nil),
		Score:    0,
		GameOver: false,
		Atlas:    atlas,
	}
}

// Update はゲームの状態を更新します
func (g *Game) Update(dt float64) {
	if g.GameOver {
		// ゲームオーバー時はRキーでリスタート
		if g.Window.JustPressed(pixelgl.KeyR) {
			g.Restart()
		}
		return
	}
	
	// プレイヤー入力の取得
	input := entity.PlayerInput{
		MoveLeft:  g.Window.Pressed(pixelgl.KeyLeft) || g.Window.Pressed(pixelgl.KeyA),
		MoveRight: g.Window.Pressed(pixelgl.KeyRight) || g.Window.Pressed(pixelgl.KeyD),
		Jump:      g.Window.JustPressed(pixelgl.KeySpace) || g.Window.JustPressed(pixelgl.KeyW),
		Attack:    g.Window.JustPressed(pixelgl.KeyX) || g.Window.JustPressed(pixelgl.KeyJ),
		UseAbility: g.Window.JustPressed(pixelgl.KeyZ) || g.Window.JustPressed(pixelgl.KeyK),
	}
	
	// プレイヤー更新
	g.Player.Update(dt, input, g.Stage.Width, g.Stage.Height)
	
	// プラットフォームとの衝突判定
	g.checkPlayerPlatformCollision()
	
	// 敵の更新
	for _, enemy := range g.Enemies {
		if enemy.IsAlive {
			enemy.Update(dt, g.Player.Position, g.Stage.Width, g.Stage.Height)
		}
	}
	
	// 衝突判定
	g.checkCollisions()
	
	// 能力の更新
	if g.Player.CurrentAbility != nil {
		g.Player.CurrentAbility.Update(dt)
	}
	
	// ゲームオーバー判定
	if g.Player.Health <= 0 {
		g.GameOver = true
	}
	
	// 敵が全滅したら新しい敵を追加
	allDead := true
	for _, enemy := range g.Enemies {
		if enemy.IsAlive {
			allDead = false
			break
		}
	}
	if allDead {
		g.spawnNewWave()
	}
}

// checkPlayerPlatformCollision はプレイヤーとプラットフォームの衝突をチェックします
func (g *Game) checkPlayerPlatformCollision() {
	playerBounds := g.Player.GetBounds()
	
	for _, platform := range g.Stage.Platforms {
		if playerBounds.Intersects(platform.Rect) {
			// 上から着地
			if g.Player.Velocity.Y <= 0 && 
			   g.Player.Position.Y-g.Player.Radius > platform.Rect.Max.Y-5 {
				g.Player.Position.Y = platform.Rect.Max.Y + g.Player.Radius
				g.Player.Velocity.Y = 0
				g.Player.IsGrounded = true
				g.Player.JumpCount = 0
			}
			// 下から衝突
			if g.Player.Velocity.Y > 0 && 
			   g.Player.Position.Y+g.Player.Radius < platform.Rect.Min.Y+5 {
				g.Player.Position.Y = platform.Rect.Min.Y - g.Player.Radius
				g.Player.Velocity.Y = 0
			}
		}
	}
}

// checkCollisions は衝突判定を行います
func (g *Game) checkCollisions() {
	playerBounds := g.Player.GetBounds()
	
	for _, enemy := range g.Enemies {
		if !enemy.IsAlive {
			continue
		}
		
		enemyBounds := enemy.GetBounds()
		
		// プレイヤーと敵の衝突
		if playerBounds.Intersects(enemyBounds) {
			// プレイヤーが上から踏んだ場合
			if g.Player.Velocity.Y < 0 && g.Player.Position.Y > enemy.Position.Y {
				enemy.TakeDamage(30)
				g.Player.Velocity.Y = 200 // 跳ね返り
				g.Score += 10
				
				// 敵を倒したらコピー能力を獲得
				if !enemy.IsAlive {
					abilityType := enemy.GetAbilityType()
					g.Player.SetAbility(ability.CreateAbility(abilityType))
					g.Score += 50
				}
			} else {
				// 横や下から当たった場合はダメージ
				g.Player.TakeDamage(10)
			}
		}
	}
}

// spawnNewWave は新しい敵の波を生成します
func (g *Game) spawnNewWave() {
	g.Score += 100
	
	// ランダムに敵を配置
	numEnemies := 3 + rand.Intn(3)
	g.Enemies = make([]*entity.Enemy, 0, numEnemies)
	
	for i := 0; i < numEnemies; i++ {
		x := 100 + rand.Float64()*(g.Stage.Width-200)
		y := 150 + rand.Float64()*200
		enemyType := entity.EnemyType(rand.Intn(3))
		
		g.Enemies = append(g.Enemies, entity.NewEnemy(pixel.V(x, y), enemyType))
	}
}

// Draw はゲーム画面を描画します
func (g *Game) Draw() {
	g.Window.Clear(g.Stage.Background)
	g.IMDraw.Clear()
	
	// ステージ描画
	g.Stage.Draw(g.IMDraw)
	
	// 敵描画
	for _, enemy := range g.Enemies {
		enemy.Draw(g.IMDraw)
	}
	
	// プレイヤー描画
	g.Player.Draw(g.IMDraw)
	
	// IMDrawを画面に反映
	g.IMDraw.Draw(g.Window)
	
	// UI描画
	g.drawUI()
	
	// ゲームオーバー画面
	if g.GameOver {
		g.drawGameOver()
	}
}

// drawUI はUIを描画します
func (g *Game) drawUI() {
	// スコア表示
	scoreText := text.New(pixel.V(10, WindowHeight-30), g.Atlas)
	scoreText.Color = colornames.White
	fmt.Fprintf(scoreText, "Score: %d", g.Score)
	scoreText.Draw(g.Window, pixel.IM.Scaled(scoreText.Orig, 2))
	
	// HP表示
	hpText := text.New(pixel.V(10, WindowHeight-60), g.Atlas)
	hpText.Color = colornames.White
	fmt.Fprintf(hpText, "HP: %d/%d", g.Player.Health, g.Player.MaxHealth)
	hpText.Draw(g.Window, pixel.IM.Scaled(hpText.Orig, 2))
	
	// HPバー
	g.drawHealthBar()
	
	// 能力表示
	if g.Player.CurrentAbility != nil {
		abilityText := text.New(pixel.V(10, WindowHeight-90), g.Atlas)
		abilityText.Color = colornames.Yellow
		fmt.Fprintf(abilityText, "Ability: %s", g.Player.CurrentAbility.GetName())
		abilityText.Draw(g.Window, pixel.IM.Scaled(abilityText.Orig, 2))
	}
	
	// 操作説明
	controlText := text.New(pixel.V(10, 30), g.Atlas)
	controlText.Color = colornames.White
	fmt.Fprintf(controlText, "Arrow/WASD: Move  Space: Jump  X: Attack")
	controlText.Draw(g.Window, pixel.IM.Scaled(controlText.Orig, 1.5))
}

// drawHealthBar はHPバーを描画します
func (g *Game) drawHealthBar() {
	barWidth := 200.0
	barHeight := 20.0
	barX := 10.0
	barY := WindowHeight - 100.0
	
	// 背景
	g.IMDraw.Color = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	g.IMDraw.Push(pixel.V(barX, barY))
	g.IMDraw.Push(pixel.V(barX+barWidth, barY+barHeight))
	g.IMDraw.Rectangle(0)
	
	// HP
	hpRatio := float64(g.Player.Health) / float64(g.Player.MaxHealth)
	g.IMDraw.Color = color.RGBA{R: 255, G: 100, B: 100, A: 255}
	g.IMDraw.Push(pixel.V(barX, barY))
	g.IMDraw.Push(pixel.V(barX+barWidth*hpRatio, barY+barHeight))
	g.IMDraw.Rectangle(0)
	
	g.IMDraw.Draw(g.Window)
}

// drawGameOver はゲームオーバー画面を描画します
func (g *Game) drawGameOver() {
	// 半透明の黒背景
	g.IMDraw.Color = color.RGBA{R: 0, G: 0, B: 0, A: 180}
	g.IMDraw.Push(pixel.V(0, 0))
	g.IMDraw.Push(pixel.V(WindowWidth, WindowHeight))
	g.IMDraw.Rectangle(0)
	g.IMDraw.Draw(g.Window)
	
	// ゲームオーバーテキスト
	gameOverText := text.New(pixel.V(WindowWidth/2-100, WindowHeight/2), g.Atlas)
	gameOverText.Color = colornames.Red
	fmt.Fprintf(gameOverText, "GAME OVER")
	gameOverText.Draw(g.Window, pixel.IM.Scaled(gameOverText.Orig, 4))
	
	// スコア表示
	finalScoreText := text.New(pixel.V(WindowWidth/2-80, WindowHeight/2-50), g.Atlas)
	finalScoreText.Color = colornames.White
	fmt.Fprintf(finalScoreText, "Final Score: %d", g.Score)
	finalScoreText.Draw(g.Window, pixel.IM.Scaled(finalScoreText.Orig, 2))
	
	// リスタート案内
	restartText := text.New(pixel.V(WindowWidth/2-120, WindowHeight/2-100), g.Atlas)
	restartText.Color = colornames.Yellow
	fmt.Fprintf(restartText, "Press R to Restart")
	restartText.Draw(g.Window, pixel.IM.Scaled(restartText.Orig, 2))
}

// Restart はゲームを再起動します
func (g *Game) Restart() {
	g.Player = entity.NewPlayer(pixel.V(WindowWidth/2, 200))
	g.Enemies = []*entity.Enemy{
		entity.NewEnemy(pixel.V(200, 150), entity.EnemyTypeWalker),
		entity.NewEnemy(pixel.V(500, 200), entity.EnemyTypeFlyer),
		entity.NewEnemy(pixel.V(800, 150), entity.EnemyTypeJumper),
	}
	g.Score = 0
	g.GameOver = false
}

// Run はゲームループを実行します
func (g *Game) Run() {
	last := time.Now()
	
	for !g.Window.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		
		g.Update(dt)
		g.Draw()
		g.Window.Update()
	}
}
