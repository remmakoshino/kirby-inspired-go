package entity

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/remmakoshino/kirby-inspired-go/internal/ability"
)

const (
	PlayerRadius      = 20.0
	PlayerSpeed       = 200.0
	JumpForce         = 400.0
	Gravity           = 800.0
	MaxFallSpeed      = 500.0
	GroundFriction    = 0.85
	AirFriction       = 0.95
	MaxJumps          = 2 // ダブルジャンプ可能
)

// Player はプレイヤーキャラクターを表します
type Player struct {
	Position     pixel.Vec
	Velocity     pixel.Vec
	Radius       float64
	Health       int
	MaxHealth    int
	IsGrounded   bool
	JumpCount    int
	IsFacingLeft bool
	
	// コピー能力関連
	CurrentAbility ability.Ability
	
	// アニメーション関連
	AnimationState string
	AnimationTime  float64
	
	// 無敵時間（ダメージ後）
	InvincibleTime float64
}

// NewPlayer は新しいプレイヤーを作成します
func NewPlayer(startPos pixel.Vec) *Player {
	return &Player{
		Position:       startPos,
		Velocity:       pixel.ZV,
		Radius:         PlayerRadius,
		Health:         100,
		MaxHealth:      100,
		IsGrounded:     false,
		JumpCount:      0,
		IsFacingLeft:   false,
		CurrentAbility: nil,
		AnimationState: "idle",
		AnimationTime:  0,
		InvincibleTime: 0,
	}
}

// Update はプレイヤーの状態を更新します
func (p *Player) Update(dt float64, input PlayerInput, stageWidth, stageHeight float64) {
	// 無敵時間の更新
	if p.InvincibleTime > 0 {
		p.InvincibleTime -= dt
	}
	
	// アニメーション時間の更新
	p.AnimationTime += dt
	
	// 左右移動
	if input.MoveLeft {
		p.Velocity.X -= PlayerSpeed * dt * 10
		p.IsFacingLeft = true
		p.AnimationState = "walk"
	} else if input.MoveRight {
		p.Velocity.X += PlayerSpeed * dt * 10
		p.IsFacingLeft = false
		p.AnimationState = "walk"
	} else if p.IsGrounded {
		p.AnimationState = "idle"
	}
	
	// ジャンプ
	if input.Jump && p.JumpCount < MaxJumps {
		p.Velocity.Y = JumpForce
		p.JumpCount++
		p.IsGrounded = false
		p.AnimationState = "jump"
	}
	
	// 攻撃（能力使用）
	if input.Attack && p.CurrentAbility != nil {
		p.CurrentAbility.Use(p)
	}
	
	// 重力適用
	if !p.IsGrounded {
		p.Velocity.Y -= Gravity * dt
		if p.Velocity.Y < -MaxFallSpeed {
			p.Velocity.Y = -MaxFallSpeed
		}
	}
	
	// 摩擦適用
	if p.IsGrounded {
		p.Velocity.X *= GroundFriction
	} else {
		p.Velocity.X *= AirFriction
	}
	
	// 速度制限
	if math.Abs(p.Velocity.X) > PlayerSpeed {
		if p.Velocity.X > 0 {
			p.Velocity.X = PlayerSpeed
		} else {
			p.Velocity.X = -PlayerSpeed
		}
	}
	
	// 位置更新
	p.Position = p.Position.Add(p.Velocity.Scaled(dt))
	
	// 画面端の処理
	if p.Position.X-p.Radius < 0 {
		p.Position.X = p.Radius
		p.Velocity.X = 0
	} else if p.Position.X+p.Radius > stageWidth {
		p.Position.X = stageWidth - p.Radius
		p.Velocity.X = 0
	}
	
	// 地面との衝突（簡易版）
	if p.Position.Y-p.Radius <= 0 {
		p.Position.Y = p.Radius
		p.Velocity.Y = 0
		p.IsGrounded = true
		p.JumpCount = 0
	} else {
		p.IsGrounded = false
	}
	
	// 天井との衝突
	if p.Position.Y+p.Radius > stageHeight {
		p.Position.Y = stageHeight - p.Radius
		p.Velocity.Y = 0
	}
	
	// 画面外に落ちた場合
	if p.Position.Y < -100 {
		p.TakeDamage(20)
		p.Position = pixel.V(stageWidth/2, 200)
		p.Velocity = pixel.ZV
	}
}

// Draw はプレイヤーを描画します（カービィ風のピンクキャラクター）
func (p *Player) Draw(imd *imdraw.IMDraw) {
	// 無敵時間中は点滅
	if p.InvincibleTime > 0 && int(p.InvincibleTime*10)%2 == 0 {
		return
	}
	
	// 体の色（ピンク）
	bodyColor := color.RGBA{R: 255, G: 182, B: 193, A: 255}
	if p.CurrentAbility != nil {
		bodyColor = p.CurrentAbility.GetColor()
	}
	
	// 本体（丸い体）
	imd.Color = bodyColor
	imd.Push(p.Position)
	imd.Circle(p.Radius, 0)
	
	// 足（小さな楕円）
	footColor := color.RGBA{R: 180, G: 60, B: 80, A: 255}
	imd.Color = footColor
	
	leftFootPos := p.Position.Add(pixel.V(-p.Radius*0.4, -p.Radius*0.8))
	rightFootPos := p.Position.Add(pixel.V(p.Radius*0.4, -p.Radius*0.8))
	
	imd.Push(leftFootPos)
	imd.Circle(p.Radius*0.3, 0)
	imd.Push(rightFootPos)
	imd.Circle(p.Radius*0.3, 0)
	
	// 目
	eyeColor := color.RGBA{R: 20, G: 20, B: 80, A: 255}
	imd.Color = eyeColor
	
	eyeOffsetX := p.Radius * 0.3
	if p.IsFacingLeft {
		eyeOffsetX = -eyeOffsetX
	}
	
	leftEyePos := p.Position.Add(pixel.V(-eyeOffsetX, p.Radius*0.2))
	rightEyePos := p.Position.Add(pixel.V(eyeOffsetX, p.Radius*0.2))
	
	// 目の白目
	imd.Color = color.White
	imd.Push(leftEyePos)
	imd.Circle(p.Radius*0.25, 0)
	imd.Push(rightEyePos)
	imd.Circle(p.Radius*0.25, 0)
	
	// 瞳
	imd.Color = eyeColor
	pupilOffset := pixel.V(0, -p.Radius*0.05)
	imd.Push(leftEyePos.Add(pupilOffset))
	imd.Circle(p.Radius*0.15, 0)
	imd.Push(rightEyePos.Add(pupilOffset))
	imd.Circle(p.Radius*0.15, 0)
	
	// 口（笑顔）
	drawSmile(imd, p.Position.Add(pixel.V(0, -p.Radius*0.2)), p.Radius*0.5, footColor)
	
	// 頬の赤み
	cheekColor := color.RGBA{R: 255, G: 150, B: 170, A: 200}
	imd.Color = cheekColor
	
	leftCheekPos := p.Position.Add(pixel.V(-p.Radius*0.7, 0))
	rightCheekPos := p.Position.Add(pixel.V(p.Radius*0.7, 0))
	
	imd.Push(leftCheekPos)
	imd.Circle(p.Radius*0.2, 0)
	imd.Push(rightCheekPos)
	imd.Circle(p.Radius*0.2, 0)
}

// drawSmile は笑顔の口を描画します
func drawSmile(imd *imdraw.IMDraw, center pixel.Vec, width float64, col color.Color) {
	imd.Color = col
	
	// 弧を描くための複数の点
	for i := 0; i <= 10; i++ {
		angle := math.Pi * float64(i) / 10.0
		x := center.X + width*math.Cos(angle)*0.5
		y := center.Y - width*math.Sin(angle)*0.3
		imd.Push(pixel.V(x, y))
	}
	imd.Polygon(2)
}

// TakeDamage はダメージを受けます
func (p *Player) TakeDamage(damage int) {
	if p.InvincibleTime > 0 {
		return
	}
	
	p.Health -= damage
	if p.Health < 0 {
		p.Health = 0
	}
	
	// 無敵時間を設定
	p.InvincibleTime = 1.5
}

// Heal は体力を回復します
func (p *Player) Heal(amount int) {
	p.Health += amount
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
}

// SetAbility はコピー能力を設定します
func (p *Player) SetAbility(ab ability.Ability) {
	p.CurrentAbility = ab
}

// ClearAbility はコピー能力をクリアします
func (p *Player) ClearAbility() {
	p.CurrentAbility = nil
}

// GetBounds は当たり判定用の矩形を返します
func (p *Player) GetBounds() pixel.Rect {
	return pixel.R(
		p.Position.X-p.Radius,
		p.Position.Y-p.Radius,
		p.Position.X+p.Radius,
		p.Position.Y+p.Radius,
	)
}

// GetPosition はプレイヤーの位置を返します（AbilityUserインターフェース実装）
func (p *Player) GetPosition() pixel.Vec {
	return p.Position
}

// GetVelocity はプレイヤーの速度を返します（AbilityUserインターフェース実装）
func (p *Player) GetVelocity() pixel.Vec {
	return p.Velocity
}

// SetVelocity はプレイヤーの速度を設定します（AbilityUserインターフェース実装）
func (p *Player) SetVelocity(v pixel.Vec) {
	p.Velocity = v
}

// PlayerInput はプレイヤーの入力を表します
type PlayerInput struct {
	MoveLeft  bool
	MoveRight bool
	Jump      bool
	Attack    bool
	UseAbility bool
}
