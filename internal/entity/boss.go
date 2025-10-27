package entity

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

// BossType はボスのタイプ
type BossType int

const (
	BossDedede BossType = iota
	BossMetaKnight
)

// Boss はボスキャラクターを表します
type Boss struct {
	Position      pixel.Vec
	Velocity      pixel.Vec
	Radius        float64
	Health        int
	MaxHealth     int
	Type          BossType
	Color         color.RGBA
	IsAlive       bool
	
	// AI関連
	AIState       string // "idle", "attack", "jump", "special"
	AITimer       float64
	AttackTimer   float64
	AttackCooldown float64
	
	// アニメーション
	AnimationTime  float64
	AnimationFrame int
	
	// 攻撃パターン
	AttackPattern  int
	PhaseLevel     int // 体力に応じたフェーズ
}

// NewDededeBoss はデデデ大王ボスを作成します
func NewDededeBoss(pos pixel.Vec) *Boss {
	return &Boss{
		Position:       pos,
		Velocity:       pixel.ZV,
		Radius:         50.0,
		Health:         200,
		MaxHealth:      200,
		Type:           BossDedede,
		Color:          color.RGBA{R: 255, G: 200, B: 50, A: 255},
		IsAlive:        true,
		AIState:        "idle",
		AITimer:        0,
		AttackTimer:    0,
		AttackCooldown: 2.0,
		AnimationTime:  0,
		AnimationFrame: 0,
		AttackPattern:  0,
		PhaseLevel:     1,
	}
}

// NewMetaKnightBoss はメタナイトボスを作成します
func NewMetaKnightBoss(pos pixel.Vec) *Boss {
	return &Boss{
		Position:       pos,
		Velocity:       pixel.ZV,
		Radius:         35.0,
		Health:         150,
		MaxHealth:      150,
		Type:           BossMetaKnight,
		Color:          color.RGBA{R: 100, G: 50, B: 150, A: 255},
		IsAlive:        true,
		AIState:        "idle",
		AITimer:        0,
		AttackTimer:    0,
		AttackCooldown: 1.5,
		AnimationTime:  0,
		AnimationFrame: 0,
		AttackPattern:  0,
		PhaseLevel:     1,
	}
}

// Update はボスの状態を更新します
func (b *Boss) Update(dt float64, playerPos pixel.Vec, stageWidth, stageHeight float64) {
	if !b.IsAlive {
		return
	}
	
	b.AnimationTime += dt
	b.AITimer += dt
	b.AttackTimer += dt
	
	// フェーズレベルの更新
	healthPercent := float64(b.Health) / float64(b.MaxHealth)
	if healthPercent < 0.3 {
		b.PhaseLevel = 3
	} else if healthPercent < 0.6 {
		b.PhaseLevel = 2
	} else {
		b.PhaseLevel = 1
	}
	
	// タイプ別のAI
	switch b.Type {
	case BossDedede:
		b.updateDededeAI(dt, playerPos)
	case BossMetaKnight:
		b.updateMetaKnightAI(dt, playerPos)
	}
	
	// 重力適用
	b.Velocity.Y -= Gravity * dt
	if b.Velocity.Y < -MaxFallSpeed {
		b.Velocity.Y = -MaxFallSpeed
	}
	
	// 位置更新
	b.Position = b.Position.Add(b.Velocity.Scaled(dt))
	
	// 地面との衝突
	if b.Position.Y-b.Radius <= 0 {
		b.Position.Y = b.Radius
		b.Velocity.Y = 0
	}
	
	// 画面端処理
	if b.Position.X-b.Radius < 0 {
		b.Position.X = b.Radius
		b.Velocity.X = 0
	} else if b.Position.X+b.Radius > stageWidth {
		b.Position.X = stageWidth - b.Radius
		b.Velocity.X = 0
	}
}

// updateDededeAI はデデデ大王のAIを更新
func (b *Boss) updateDededeAI(dt float64, playerPos pixel.Vec) {
	distance := playerPos.Sub(b.Position).Len()
	
	switch b.AIState {
	case "idle":
		if b.AITimer > 1.0 {
			// ランダムに攻撃パターンを選択
			pattern := rand.Intn(3)
			switch pattern {
			case 0:
				b.AIState = "hammer_attack"
			case 1:
				b.AIState = "jump_attack"
			case 2:
				b.AIState = "charge"
			}
			b.AITimer = 0
		}
		
	case "hammer_attack":
		// ハンマー攻撃
		if distance < 100 {
			// プレイヤーに向かって移動
			direction := playerPos.Sub(b.Position).Unit()
			b.Velocity.X = direction.X * 100
		}
		
		if b.AITimer > 1.5 {
			b.AIState = "idle"
			b.AITimer = 0
		}
		
	case "jump_attack":
		// ジャンプ攻撃
		if b.Position.Y-b.Radius <= 1 && b.AITimer > 0.5 {
			b.Velocity.Y = 400.0 + float64(b.PhaseLevel)*50
			b.Velocity.X = (playerPos.X - b.Position.X) * 2
		}
		
		if b.AITimer > 2.0 {
			b.AIState = "idle"
			b.AITimer = 0
		}
		
	case "charge":
		// 突進攻撃
		direction := playerPos.Sub(b.Position).Unit()
		b.Velocity.X = direction.X * 200 * float64(b.PhaseLevel)
		
		if b.AITimer > 1.0 {
			b.AIState = "idle"
			b.AITimer = 0
			b.Velocity.X = 0
		}
	}
}

// updateMetaKnightAI はメタナイトのAIを更新
func (b *Boss) updateMetaKnightAI(dt float64, playerPos pixel.Vec) {
	distance := playerPos.Sub(b.Position).Len()
	
	switch b.AIState {
	case "idle":
		if b.AITimer > 0.5 {
			pattern := rand.Intn(4)
			switch pattern {
			case 0:
				b.AIState = "sword_combo"
			case 1:
				b.AIState = "tornado_slash"
			case 2:
				b.AIState = "dash_attack"
			case 3:
				b.AIState = "cape_defense"
			}
			b.AITimer = 0
		}
		
	case "sword_combo":
		// 剣の連続攻撃
		if distance < 80 {
			direction := playerPos.Sub(b.Position).Unit()
			b.Velocity.X = direction.X * 150
		}
		
		if b.AITimer > 1.2 {
			b.AIState = "idle"
			b.AITimer = 0
		}
		
	case "tornado_slash":
		// トルネード斬り
		b.Velocity.X = math.Cos(b.AnimationTime*10) * 200
		b.Velocity.Y = 100
		
		if b.AITimer > 1.5 {
			b.AIState = "idle"
			b.AITimer = 0
		}
		
	case "dash_attack":
		// ダッシュ攻撃
		direction := playerPos.Sub(b.Position).Unit()
		b.Velocity.X = direction.X * 300 * float64(b.PhaseLevel)
		
		if b.AITimer > 0.8 {
			b.AIState = "idle"
			b.AITimer = 0
			b.Velocity.X = 0
		}
		
	case "cape_defense":
		// マント防御（一時的に無敵）
		b.Velocity.X = 0
		
		if b.AITimer > 2.0 {
			b.AIState = "idle"
			b.AITimer = 0
		}
	}
}

// Draw はボスを描画します
func (b *Boss) Draw(imd *imdraw.IMDraw) {
	if !b.IsAlive {
		return
	}
	
	switch b.Type {
	case BossDedede:
		b.drawDedede(imd)
	case BossMetaKnight:
		b.drawMetaKnight(imd)
	}
}

// drawDedede はデデデ大王を描画
func (b *Boss) drawDedede(imd *imdraw.IMDraw) {
	// 本体（黄色の大きな体）
	imd.Color = b.Color
	imd.Push(b.Position)
	imd.Circle(b.Radius, 0)
	
	// 王冠
	crownColor := color.RGBA{R: 255, G: 215, B: 0, A: 255}
	imd.Color = crownColor
	crownTop := b.Position.Add(pixel.V(0, b.Radius+10))
	imd.Push(pixel.V(crownTop.X-20, crownTop.Y-10))
	imd.Push(pixel.V(crownTop.X+20, crownTop.Y-10))
	imd.Push(pixel.V(crownTop.X, crownTop.Y+15))
	imd.Polygon(0)
	
	// 目
	imd.Color = color.RGBA{R: 50, G: 50, B: 50, A: 255}
	eyeY := b.Position.Y + b.Radius*0.2
	imd.Push(pixel.V(b.Position.X-b.Radius*0.3, eyeY))
	imd.Circle(b.Radius*0.15, 0)
	imd.Push(pixel.V(b.Position.X+b.Radius*0.3, eyeY))
	imd.Circle(b.Radius*0.15, 0)
	
	// くちばし
	beakColor := color.RGBA{R: 255, G: 140, B: 0, A: 255}
	imd.Color = beakColor
	beakCenter := b.Position.Add(pixel.V(0, -b.Radius*0.2))
	imd.Push(pixel.V(beakCenter.X-15, beakCenter.Y+10))
	imd.Push(pixel.V(beakCenter.X+15, beakCenter.Y+10))
	imd.Push(pixel.V(beakCenter.X, beakCenter.Y-15))
	imd.Polygon(0)
	
	// ハンマー（攻撃時）
	if b.AIState == "hammer_attack" {
		hammerColor := color.RGBA{R: 139, G: 69, B: 19, A: 255}
		imd.Color = hammerColor
		hammerPos := b.Position.Add(pixel.V(b.Radius*0.8, b.Radius*0.5))
		// 柄
		imd.Push(hammerPos)
		imd.Push(hammerPos.Add(pixel.V(30, -30)))
		imd.Line(8)
		// ヘッド
		imd.Color = color.RGBA{R: 150, G: 150, B: 150, A: 255}
		hammerHead := hammerPos.Add(pixel.V(30, -30))
		imd.Push(hammerHead)
		imd.Circle(15, 0)
	}
}

// drawMetaKnight はメタナイトを描画
func (b *Boss) drawMetaKnight(imd *imdraw.IMDraw) {
	// 本体（紫の球体）
	imd.Color = b.Color
	imd.Push(b.Position)
	imd.Circle(b.Radius, 0)
	
	// マスク（銀色）
	maskColor := color.RGBA{R: 200, G: 200, B: 220, A: 255}
	imd.Color = maskColor
	maskCenter := b.Position.Add(pixel.V(0, b.Radius*0.1))
	imd.Push(pixel.V(maskCenter.X-b.Radius*0.6, maskCenter.Y-5))
	imd.Push(pixel.V(maskCenter.X+b.Radius*0.6, maskCenter.Y+10))
	imd.Rectangle(0)
	
	// 目（黄色）
	imd.Color = colornames.Yellow
	eyePos := b.Position.Add(pixel.V(0, b.Radius*0.15))
	imd.Push(eyePos)
	imd.Circle(b.Radius*0.2, 0)
	
	// 剣（常に表示）
	swordColor := color.RGBA{R: 180, G: 180, B: 200, A: 255}
	imd.Color = swordColor
	swordStart := b.Position.Add(pixel.V(b.Radius*0.7, 0))
	swordEnd := swordStart.Add(pixel.V(40, -10))
	imd.Push(swordStart)
	imd.Push(swordEnd)
	imd.Line(6)
	
	// 剣先
	imd.Push(swordEnd)
	imd.Push(swordEnd.Add(pixel.V(15, -5)))
	imd.Line(4)
	
	// マント（防御時に強調）
	capeAlpha := uint8(150)
	if b.AIState == "cape_defense" {
		capeAlpha = 255
	}
	capeColor := color.RGBA{R: 150, G: 100, B: 200, A: capeAlpha}
	imd.Color = capeColor
	
	// マントの形状
	capeLeft := b.Position.Add(pixel.V(-b.Radius*0.9, b.Radius*0.3))
	capeRight := b.Position.Add(pixel.V(b.Radius*0.9, b.Radius*0.3))
	capeBottom := b.Position.Add(pixel.V(0, -b.Radius*1.2))
	
	imd.Push(capeLeft)
	imd.Push(capeRight)
	imd.Push(capeBottom)
	imd.Polygon(0)
	
	// トルネードエフェクト（トルネード攻撃時）
	if b.AIState == "tornado_slash" {
		tornadoColor := color.RGBA{R: 150, G: 255, B: 150, A: 200}
		imd.Color = tornadoColor
		for i := 0; i < 3; i++ {
			angle := b.AnimationTime*10 + float64(i)*2*math.Pi/3
			offset := pixel.V(
				math.Cos(angle)*b.Radius*1.5,
				math.Sin(angle)*b.Radius*1.5,
			)
			imd.Push(b.Position.Add(offset))
			imd.Circle(10, 0)
		}
	}
}

// TakeDamage はダメージを受けます
func (b *Boss) TakeDamage(damage int) {
	// 防御状態の場合はダメージ軽減
	if b.AIState == "cape_defense" && b.Type == BossMetaKnight {
		damage = damage / 2
	}
	
	b.Health -= damage
	if b.Health <= 0 {
		b.Health = 0
		b.IsAlive = false
	}
}

// GetBounds は当たり判定用の矩形を返します
func (b *Boss) GetBounds() pixel.Rect {
	return pixel.R(
		b.Position.X-b.Radius,
		b.Position.Y-b.Radius,
		b.Position.X+b.Radius,
		b.Position.Y+b.Radius,
	)
}

// IsAttacking は攻撃中かどうかを返します
func (b *Boss) IsAttacking() bool {
	return b.AIState != "idle" && b.AIState != "cape_defense"
}
