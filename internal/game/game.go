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
	"github.com/remmakoshino/kirby-inspired-go/internal/menu"
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
	MetaKnight *entity.MetaKnightPlayer
	Enemies  []*entity.Enemy
	WaddleDees []*entity.WaddleDee
	WaddleDoos []*entity.WaddleDoo
	Boss     *entity.Boss
	Stage    *stage.Stage
	IMDraw   *imdraw.IMDraw
	Score    int
	GameOver bool
	Victory  bool
	
	// メニューシステム
	MenuManager *menu.MenuManager
	
	// ステージ情報
	CurrentStage int
	PlayerCharacter string // "Kirby" or "MetaKnight"
	
	// UI関連
	Atlas *text.Atlas
}

// NewGame は新しいゲームを作成します
func NewGame(win *pixelgl.Window) *Game {
	rand.Seed(time.Now().UnixNano())
	
	// テキスト描画用のアトラス
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	
	// メニューマネージャーの作成
	menuMgr := menu.NewMenuManager(win)
	
	return &Game{
		Window:      win,
		IMDraw:      imdraw.New(nil),
		Score:       0,
		GameOver:    false,
		Victory:     false,
		Atlas:       atlas,
		MenuManager: menuMgr,
		CurrentStage: 0,
		PlayerCharacter: "",
	}
}

// InitializeStage はステージを初期化します
func (g *Game) InitializeStage(stageNum int, character string) {
	g.CurrentStage = stageNum
	g.PlayerCharacter = character
	g.GameOver = false
	g.Victory = false
	
	// ステージ作成
	g.Stage = stage.CreateDefaultStage(WindowWidth, WindowHeight)
	
	// キャラクター作成
	startPos := pixel.V(WindowWidth/2, 200)
	if character == "Kirby" {
		g.Player = entity.NewPlayer(startPos)
		g.MetaKnight = nil
	} else {
		g.MetaKnight = entity.NewMetaKnightPlayer(startPos)
		g.Player = nil
	}
	
	// ステージに応じた敵とボスを配置
	g.setupStageEnemies(stageNum)
}

// setupStageEnemies はステージごとの敵を配置
func (g *Game) setupStageEnemies(stageNum int) {
	g.Enemies = []*entity.Enemy{}
	g.WaddleDees = []*entity.WaddleDee{}
	g.WaddleDoos = []*entity.WaddleDoo{}
	
	if stageNum == 1 {
		// ステージ1: ワドルディ中心
		g.WaddleDees = append(g.WaddleDees,
			entity.NewWaddleDee(pixel.V(200, 150)),
			entity.NewWaddleDee(pixel.V(400, 200)),
			entity.NewWaddleDee(pixel.V(700, 150)),
		)
		
		g.WaddleDoos = append(g.WaddleDoos,
			entity.NewWaddleDoo(pixel.V(500, 250)),
		)
		
		// ボス: デデデ大王
		g.Boss = entity.NewDededeBoss(pixel.V(WindowWidth-200, 200))
		
	} else if stageNum == 2 {
		// ステージ2: ワドルドゥ中心
		g.WaddleDees = append(g.WaddleDees,
			entity.NewWaddleDee(pixel.V(300, 150)),
		)
		
		g.WaddleDoos = append(g.WaddleDoos,
			entity.NewWaddleDoo(pixel.V(250, 200)),
			entity.NewWaddleDoo(pixel.V(600, 250)),
			entity.NewWaddleDoo(pixel.V(800, 150)),
		)
		
		// ボス: メタナイト
		g.Boss = entity.NewMetaKnightBoss(pixel.V(WindowWidth-200, 200))
	}
	
	// 基本の敵も配置
	g.Enemies = append(g.Enemies,
		entity.NewEnemy(pixel.V(350, 200), entity.EnemyTypeFlyer),
		entity.NewEnemy(pixel.V(650, 180), entity.EnemyTypeJumper),
	)
}

// Update はゲームの状態を更新します
func (g *Game) Update(dt float64) {
	// メニュー画面の処理
	if g.MenuManager.State != menu.StatePlaying {
		g.MenuManager.Update(dt)
		
		// ゲーム開始時の初期化
		if g.MenuManager.State == menu.StatePlaying && g.Stage == nil {
			character := "Kirby"
			if g.MenuManager.SelectedCharacter == menu.CharacterMetaKnight {
				character = "MetaKnight"
			}
			g.InitializeStage(g.MenuManager.SelectedStage, character)
		}
		return
	}
	
	if g.GameOver || g.Victory {
		// ゲームオーバー/クリア時はRキーでメニューに戻る
		if g.Window.JustPressed(pixelgl.KeyR) {
			g.MenuManager.State = menu.StateTitleScreen
			g.Stage = nil
			g.Score = 0
		}
		return
	}
	
	// プレイヤー更新
	if g.Player != nil {
		input := entity.PlayerInput{
			MoveLeft:   g.Window.Pressed(pixelgl.KeyLeft) || g.Window.Pressed(pixelgl.KeyA),
			MoveRight:  g.Window.Pressed(pixelgl.KeyRight) || g.Window.Pressed(pixelgl.KeyD),
			Jump:       g.Window.JustPressed(pixelgl.KeySpace) || g.Window.JustPressed(pixelgl.KeyW),
			Attack:     g.Window.JustPressed(pixelgl.KeyX) || g.Window.JustPressed(pixelgl.KeyJ),
			UseAbility: g.Window.JustPressed(pixelgl.KeyZ) || g.Window.JustPressed(pixelgl.KeyK),
		}
		g.Player.Update(dt, input, g.Stage.Width, g.Stage.Height)
		g.checkPlayerPlatformCollision()
		
		if g.Player.CurrentAbility != nil {
			g.Player.CurrentAbility.Update(dt)
		}
		
		// ゲームオーバー判定
		if g.Player.Health <= 0 {
			g.GameOver = true
		}
	} else if g.MetaKnight != nil {
		g.MetaKnight.Update(dt, g.Window)
		g.checkMetaKnightPlatformCollision()
		
		// ゲームオーバー判定
		if g.MetaKnight.Health <= 0 {
			g.GameOver = true
		}
	}
	
	// 敵の更新
	playerPos := pixel.ZV
	if g.Player != nil {
		playerPos = g.Player.Position
	} else if g.MetaKnight != nil {
		playerPos = g.MetaKnight.Position
	}
	
	for _, enemy := range g.Enemies {
		if enemy.IsAlive {
			enemy.Update(dt, playerPos, g.Stage.Width, g.Stage.Height)
		}
	}
	
	for _, waddleDee := range g.WaddleDees {
		if waddleDee.IsAlive {
			waddleDee.Update(dt, playerPos, g.Stage.Width, g.Stage.Height)
		}
	}
	
	for _, waddleDoo := range g.WaddleDoos {
		if waddleDoo.IsAlive {
			waddleDoo.Update(dt, playerPos, g.Stage.Width, g.Stage.Height)
		}
	}
	
	// ボスの更新
	if g.Boss != nil && g.Boss.IsAlive {
		g.Boss.Update(dt, playerPos, g.Stage.Width, g.Stage.Height)
	}
	
	// 衝突判定
	g.checkCollisions()
	
	// 勝利判定（ボスを倒した）
	if g.Boss != nil && !g.Boss.IsAlive && !g.Victory {
		g.Victory = true
		g.Score += 1000
	}
}

// checkPlayerPlatformCollision はプレイヤーとプラットフォームの衝突をチェックします
func (g *Game) checkPlayerPlatformCollision() {
	if g.Player == nil {
		return
	}
	
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

// checkMetaKnightPlatformCollision はメタナイトとプラットフォームの衝突をチェック
func (g *Game) checkMetaKnightPlatformCollision() {
	if g.MetaKnight == nil {
		return
	}
	
	playerBounds := g.MetaKnight.GetBounds()
	
	for _, platform := range g.Stage.Platforms {
		if playerBounds.Intersects(platform.Rect) {
			// 上から着地
			if g.MetaKnight.Velocity.Y <= 0 && 
			   g.MetaKnight.Position.Y-g.MetaKnight.Radius > platform.Rect.Max.Y-5 {
				g.MetaKnight.Position.Y = platform.Rect.Max.Y + g.MetaKnight.Radius
				g.MetaKnight.Velocity.Y = 0
				g.MetaKnight.ResetJump()
			}
			// 下から衝突
			if g.MetaKnight.Velocity.Y > 0 && 
			   g.MetaKnight.Position.Y+g.MetaKnight.Radius < platform.Rect.Min.Y+5 {
				g.MetaKnight.Position.Y = platform.Rect.Min.Y - g.MetaKnight.Radius
				g.MetaKnight.Velocity.Y = 0
			}
		}
	}
}

// checkCollisions は衝突判定を行います
func (g *Game) checkCollisions() {
	var playerBounds pixel.Rect
	var playerPos pixel.Vec
	var isKirby bool
	
	if g.Player != nil {
		playerBounds = g.Player.GetBounds()
		playerPos = g.Player.Position
		isKirby = true
	} else if g.MetaKnight != nil {
		playerBounds = g.MetaKnight.GetBounds()
		playerPos = g.MetaKnight.Position
		isKirby = false
	} else {
		return
	}
	
	// 通常の敵との衝突
	for _, enemy := range g.Enemies {
		if !enemy.IsAlive {
			continue
		}
		
		enemyBounds := enemy.GetBounds()
		
		if playerBounds.Intersects(enemyBounds) {
			// プレイヤーが上から踏んだ場合
			if playerPos.Y > enemy.Position.Y+10 {
				enemy.TakeDamage(30)
				if isKirby {
					g.Player.Velocity.Y = 200
				} else {
					g.MetaKnight.Velocity.Y = 200
				}
				g.Score += 10
				
				if !enemy.IsAlive {
					abilityType := enemy.GetAbilityType()
					if isKirby {
						g.Player.SetAbility(ability.CreateAbility(abilityType))
					}
					g.Score += 50
				}
			} else {
				// 横や下から当たった場合はダメージ
				if isKirby {
					g.Player.TakeDamage(10)
				} else {
					g.MetaKnight.TakeDamage(10)
				}
			}
		}
	}
	
	// ワドルディとの衝突
	for _, waddleDee := range g.WaddleDees {
		if !waddleDee.IsAlive {
			continue
		}
		
		enemyBounds := waddleDee.GetBounds()
		
		if playerBounds.Intersects(enemyBounds) {
			if playerPos.Y > waddleDee.Position.Y+10 {
				waddleDee.TakeDamage(25)
				if isKirby {
					g.Player.Velocity.Y = 200
				} else {
					g.MetaKnight.Velocity.Y = 200
				}
				g.Score += 15
			} else {
				if isKirby {
					g.Player.TakeDamage(8)
				} else {
					g.MetaKnight.TakeDamage(8)
				}
			}
		}
	}
	
	// ワドルドゥとの衝突
	for _, waddleDoo := range g.WaddleDoos {
		if !waddleDoo.IsAlive {
			continue
		}
		
		enemyBounds := waddleDoo.GetBounds()
		
		if playerBounds.Intersects(enemyBounds) {
			if playerPos.Y > waddleDoo.Position.Y+10 {
				waddleDoo.TakeDamage(30)
				if isKirby {
					g.Player.Velocity.Y = 200
				} else {
					g.MetaKnight.Velocity.Y = 200
				}
				g.Score += 20
			} else {
				if isKirby {
					g.Player.TakeDamage(12)
				} else {
					g.MetaKnight.TakeDamage(12)
				}
			}
		}
	}
	
	// ボスとの衝突
	if g.Boss != nil && g.Boss.IsAlive {
		bossBounds := g.Boss.GetBounds()
		
		if playerBounds.Intersects(bossBounds) {
			// ボスが攻撃中の場合
			if g.Boss.IsAttacking() {
				if isKirby {
					g.Player.TakeDamage(20)
				} else {
					g.MetaKnight.TakeDamage(20)
				}
			}
		}
		
		// プレイヤーの攻撃がボスに当たる
		attackRange := 50.0
		if !isKirby && g.MetaKnight.IsAttacking {
			attackRange = g.MetaKnight.GetAttackRange()
		}
		
		distance := playerPos.Sub(g.Boss.Position).Len()
		if distance < attackRange+g.Boss.Radius {
			// 攻撃判定
			if (isKirby && g.Window.JustPressed(pixelgl.KeyX)) ||
			   (!isKirby && g.MetaKnight.IsAttacking) {
				damage := 10
				if !isKirby {
					damage = g.MetaKnight.GetAttackDamage()
				}
				g.Boss.TakeDamage(damage)
				g.Score += 5
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
	g.Window.Clear(colornames.Skyblue)
	
	// メニュー画面の描画
	if g.MenuManager.State != menu.StatePlaying {
		g.MenuManager.Draw()
		return
	}
	
	if g.Stage == nil {
		return
	}
	
	g.Window.Clear(g.Stage.Background)
	g.IMDraw.Clear()
	
	// ステージ描画
	g.Stage.Draw(g.IMDraw)
	
	// 敵描画
	for _, enemy := range g.Enemies {
		enemy.Draw(g.IMDraw)
	}
	
	for _, waddleDee := range g.WaddleDees {
		waddleDee.Draw(g.IMDraw)
	}
	
	for _, waddleDoo := range g.WaddleDoos {
		waddleDoo.Draw(g.IMDraw)
	}
	
	// ボス描画
	if g.Boss != nil {
		g.Boss.Draw(g.IMDraw)
	}
	
	// プレイヤー描画
	if g.Player != nil {
		g.Player.Draw(g.IMDraw)
	} else if g.MetaKnight != nil {
		g.MetaKnight.Draw(g.IMDraw)
	}
	
	// IMDrawを画面に反映
	g.IMDraw.Draw(g.Window)
	
	// UI描画
	g.drawUI()
	
	// ゲームオーバー画面
	if g.GameOver {
		g.drawGameOver()
	}
	
	// 勝利画面
	if g.Victory {
		g.drawVictory()
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
	
	var currentHP, maxHP int
	if g.Player != nil {
		currentHP = g.Player.Health
		maxHP = g.Player.MaxHealth
	} else if g.MetaKnight != nil {
		currentHP = g.MetaKnight.Health
		maxHP = g.MetaKnight.MaxHealth
	}
	
	fmt.Fprintf(hpText, "HP: %d/%d", currentHP, maxHP)
	hpText.Draw(g.Window, pixel.IM.Scaled(hpText.Orig, 2))
	
	// HPバー
	g.drawHealthBar(currentHP, maxHP)
	
	// ボスHPバー
	if g.Boss != nil && g.Boss.IsAlive {
		g.drawBossHealthBar()
	}
	
	// 能力表示（カービィのみ）
	if g.Player != nil && g.Player.CurrentAbility != nil {
		abilityText := text.New(pixel.V(10, WindowHeight-90), g.Atlas)
		abilityText.Color = colornames.Yellow
		fmt.Fprintf(abilityText, "Ability: %s", g.Player.CurrentAbility.GetName())
		abilityText.Draw(g.Window, pixel.IM.Scaled(abilityText.Orig, 2))
	}
	
	// メタナイトの場合、現在のアビリティ表示
	if g.MetaKnight != nil && g.MetaKnight.CurrentAbility != nil {
		abilityText := text.New(pixel.V(10, WindowHeight-90), g.Atlas)
		abilityText.Color = colornames.Yellow
		fmt.Fprintf(abilityText, "Weapon: %s", g.MetaKnight.CurrentAbility.GetName())
		abilityText.Draw(g.Window, pixel.IM.Scaled(abilityText.Orig, 2))
	}
	
	// ステージ表示
	stageText := text.New(pixel.V(WindowWidth-150, WindowHeight-30), g.Atlas)
	stageText.Color = colornames.White
	fmt.Fprintf(stageText, "Stage %d", g.CurrentStage)
	stageText.Draw(g.Window, pixel.IM.Scaled(stageText.Orig, 2))
	
	// 操作説明
	controlText := text.New(pixel.V(10, 30), g.Atlas)
	controlText.Color = colornames.White
	if g.Player != nil {
		fmt.Fprintf(controlText, "Arrow/WASD: Move  Space: Jump  X: Attack  Z: Ability")
	} else {
		fmt.Fprintf(controlText, "Arrow/WASD: Move  Space: Jump  E: Attack  Q: Special  1/2/3: Switch")
	}
	controlText.Draw(g.Window, pixel.IM.Scaled(controlText.Orig, 1.5))
}

// drawHealthBar はHPバーを描画します
func (g *Game) drawHealthBar(currentHP, maxHP int) {
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
	hpRatio := float64(currentHP) / float64(maxHP)
	g.IMDraw.Color = color.RGBA{R: 255, G: 100, B: 100, A: 255}
	g.IMDraw.Push(pixel.V(barX, barY))
	g.IMDraw.Push(pixel.V(barX+barWidth*hpRatio, barY+barHeight))
	g.IMDraw.Rectangle(0)
	
	g.IMDraw.Draw(g.Window)
}

// drawBossHealthBar はボスのHPバーを描画します
func (g *Game) drawBossHealthBar() {
	barWidth := 400.0
	barHeight := 30.0
	barX := WindowWidth/2 - barWidth/2
	barY := WindowHeight - 50.0
	
	// 背景
	g.IMDraw.Color = color.RGBA{R: 80, G: 80, B: 80, A: 255}
	g.IMDraw.Push(pixel.V(barX, barY))
	g.IMDraw.Push(pixel.V(barX+barWidth, barY+barHeight))
	g.IMDraw.Rectangle(0)
	
	// HP
	hpRatio := float64(g.Boss.Health) / float64(g.Boss.MaxHealth)
	g.IMDraw.Color = color.RGBA{R: 200, G: 50, B: 50, A: 255}
	g.IMDraw.Push(pixel.V(barX, barY))
	g.IMDraw.Push(pixel.V(barX+barWidth*hpRatio, barY+barHeight))
	g.IMDraw.Rectangle(0)
	
	g.IMDraw.Draw(g.Window)
	
	// ボス名表示
	bossName := "Boss"
	if g.Boss.Type == entity.BossDedede {
		bossName = "King Dedede"
	} else if g.Boss.Type == entity.BossMetaKnight {
		bossName = "Meta Knight"
	}
	
	bossText := text.New(pixel.V(barX+barWidth/2-50, barY+barHeight+5), g.Atlas)
	bossText.Color = colornames.Red
	fmt.Fprintf(bossText, "%s", bossName)
	bossText.Draw(g.Window, pixel.IM.Scaled(bossText.Orig, 2))
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
	restartText := text.New(pixel.V(WindowWidth/2-140, WindowHeight/2-100), g.Atlas)
	restartText.Color = colornames.Yellow
	fmt.Fprintf(restartText, "Press R to Return to Menu")
	restartText.Draw(g.Window, pixel.IM.Scaled(restartText.Orig, 2))
}

// drawVictory は勝利画面を描画します
func (g *Game) drawVictory() {
	// 半透明の黄色背景
	g.IMDraw.Color = color.RGBA{R: 255, G: 215, B: 0, A: 180}
	g.IMDraw.Push(pixel.V(0, 0))
	g.IMDraw.Push(pixel.V(WindowWidth, WindowHeight))
	g.IMDraw.Rectangle(0)
	g.IMDraw.Draw(g.Window)
	
	// 勝利テキスト
	victoryText := text.New(pixel.V(WindowWidth/2-100, WindowHeight/2), g.Atlas)
	victoryText.Color = colornames.Gold
	fmt.Fprintf(victoryText, "STAGE CLEAR!")
	victoryText.Draw(g.Window, pixel.IM.Scaled(victoryText.Orig, 4))
	
	// スコア表示
	finalScoreText := text.New(pixel.V(WindowWidth/2-80, WindowHeight/2-50), g.Atlas)
	finalScoreText.Color = colornames.White
	fmt.Fprintf(finalScoreText, "Score: %d", g.Score)
	finalScoreText.Draw(g.Window, pixel.IM.Scaled(finalScoreText.Orig, 2))
	
	// リスタート案内
	restartText := text.New(pixel.V(WindowWidth/2-140, WindowHeight/2-100), g.Atlas)
	restartText.Color = colornames.White
	fmt.Fprintf(restartText, "Press R to Return to Menu")
	restartText.Draw(g.Window, pixel.IM.Scaled(restartText.Orig, 2))
}

// spawnNewWave は使用されなくなりました（削除）
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
