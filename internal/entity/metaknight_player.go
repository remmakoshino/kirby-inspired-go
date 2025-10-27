package entity

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	
	"github.com/remmakoshino/kirby-inspired-go/internal/ability"
)

// MetaKnightPlayer はプレイアブルキャラクターとしてのメタナイトを表します
type MetaKnightPlayer struct {
	Position      pixel.Vec
	Velocity      pixel.Vec
	Radius        float64
	Health        int
	MaxHealth     int
	Color         color.RGBA
	IsAlive       bool
	IsJumping     bool
	CanDoubleJump bool
	HasDoubleJumped bool
	
	// メタナイト専用
	CurrentAbility ability.Ability
	Abilities      []ability.Ability
	
	// アニメーション
	AnimationTime  float64
	AnimationFrame int
	
	// 攻撃関連
	IsAttacking    bool
	AttackCooldown float64
	ComboCount     int
}

// NewMetaKnightPlayer は新しいメタナイトプレイヤーを作成します
func NewMetaKnightPlayer(startPos pixel.Vec) *MetaKnightPlayer {
	mk := &MetaKnightPlayer{
		Position:        startPos,
		Velocity:        pixel.ZV,
		Radius:          25.0,
		Health:          100,
		MaxHealth:       100,
		Color:           color.RGBA{R: 100, G: 50, B: 150, A: 255},
		IsAlive:         true,
		IsJumping:       false,
		CanDoubleJump:   true,
		HasDoubleJumped: false,
		AnimationTime:   0,
		AnimationFrame:  0,
		IsAttacking:     false,
		AttackCooldown:  0,
		ComboCount:      0,
	}
	
	// メタナイト専用アビリティ
	mk.Abilities = []ability.Ability{
		ability.NewSwordAbility(),
		ability.NewTornadoAbility(),
		ability.NewCapeBarrierAbility(),
	}
	mk.CurrentAbility = mk.Abilities[0] // デフォルトは剣
	
	return mk
}

// Update はメタナイトの状態を更新します
func (mk *MetaKnightPlayer) Update(dt float64, win *pixelgl.Window) {
	if !mk.IsAlive {
		return
	}
	
	mk.AnimationTime += dt
	
	// 攻撃クールダウン
	if mk.AttackCooldown > 0 {
		mk.AttackCooldown -= dt
		if mk.AttackCooldown <= 0 {
			mk.IsAttacking = false
			mk.ComboCount = 0
		}
	}
	
	// 移動入力
	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
		mk.Velocity.X = -PlayerSpeed
	} else if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
		mk.Velocity.X = PlayerSpeed
	} else {
		mk.Velocity.X = 0
	}
	
	// ジャンプ入力
	if win.JustPressed(pixelgl.KeySpace) || win.JustPressed(pixelgl.KeyW) {
		if !mk.IsJumping {
			mk.Velocity.Y = JumpForce
			mk.IsJumping = true
			mk.HasDoubleJumped = false
		} else if mk.CanDoubleJump && !mk.HasDoubleJumped {
			mk.Velocity.Y = JumpForce
			mk.HasDoubleJumped = true
		}
	}
	
	// アビリティ切り替え（1, 2, 3キー）
	if win.JustPressed(pixelgl.Key1) {
		mk.CurrentAbility = mk.Abilities[0] // 剣
	} else if win.JustPressed(pixelgl.Key2) {
		mk.CurrentAbility = mk.Abilities[1] // トルネード
	} else if win.JustPressed(pixelgl.Key3) {
		mk.CurrentAbility = mk.Abilities[2] // マント防御
	}
	
	// 攻撃入力（Eキー）
	if win.JustPressed(pixelgl.KeyE) && mk.AttackCooldown <= 0 {
		mk.IsAttacking = true
		mk.AttackCooldown = 0.5
		mk.ComboCount++
		if mk.ComboCount > 3 {
			mk.ComboCount = 1
		}
	}
	
	// アビリティ発動（Qキー）
	if win.JustPressed(pixelgl.KeyQ) {
		mk.ActivateAbility()
	}
	
	// 重力適用
	mk.Velocity.Y -= Gravity * dt
	if mk.Velocity.Y < -MaxFallSpeed {
		mk.Velocity.Y = -MaxFallSpeed
	}
	
	// 位置更新
	mk.Position = mk.Position.Add(mk.Velocity.Scaled(dt))
}

// ActivateAbility は現在のアビリティを発動します
func (mk *MetaKnightPlayer) ActivateAbility() {
	if mk.CurrentAbility == nil {
		return
	}
	
	switch mk.CurrentAbility.GetName() {
	case "Sword":
		// 剣の特殊攻撃（突進斬り）
		mk.Velocity.X = 300
		if mk.Velocity.X < 0 {
			mk.Velocity.X = -300
		}
		
	case "Tornado":
		// トルネードダッシュ
		mk.Velocity.X = 400
		mk.Velocity.Y = 100
		
	case "Cape Barrier":
		// マント防御（一時的に無敵、実装はゲームロジック側で処理）
		// ここでは速度を0にして防御姿勢を取る
		mk.Velocity.X = 0
	}
}

// Draw はメタナイトを描画します
func (mk *MetaKnightPlayer) Draw(imd *imdraw.IMDraw) {
	if !mk.IsAlive {
		return
	}
	
	// 本体（紫の球体）
	imd.Color = mk.Color
	imd.Push(mk.Position)
	imd.Circle(mk.Radius, 0)
	
	// マスク（銀色）
	maskColor := color.RGBA{R: 200, G: 200, B: 220, A: 255}
	imd.Color = maskColor
	maskCenter := mk.Position.Add(pixel.V(0, mk.Radius*0.1))
	imd.Push(pixel.V(maskCenter.X-mk.Radius*0.6, maskCenter.Y-5))
	imd.Push(pixel.V(maskCenter.X+mk.Radius*0.6, maskCenter.Y+10))
	imd.Rectangle(0)
	
	// 目（黄色、発光するような効果）
	eyeColor := color.RGBA{R: 255, G: 255, B: 100, A: 255}
	imd.Color = eyeColor
	eyePos := mk.Position.Add(pixel.V(0, mk.Radius*0.15))
	imd.Push(eyePos)
	imd.Circle(mk.Radius*0.25, 0)
	
	// 剣（現在のアビリティが剣の場合、または攻撃中）
	if mk.CurrentAbility.GetName() == "Sword" || mk.IsAttacking {
		swordColor := color.RGBA{R: 180, G: 180, B: 200, A: 255}
		imd.Color = swordColor
		
		swordAngle := 0.0
		if mk.IsAttacking {
			// 攻撃アニメーション
			swordAngle = math.Sin(mk.AnimationTime*20) * math.Pi / 4
		}
		
		swordStart := mk.Position.Add(pixel.V(mk.Radius*0.7, 0))
		swordOffset := pixel.V(
			math.Cos(swordAngle)*35,
			math.Sin(swordAngle)*35-10,
		)
		swordEnd := swordStart.Add(swordOffset)
		
		imd.Push(swordStart)
		imd.Push(swordEnd)
		imd.Line(5)
		
		// 剣先
		tipOffset := pixel.V(
			math.Cos(swordAngle)*12,
			math.Sin(swordAngle)*12,
		)
		imd.Push(swordEnd)
		imd.Push(swordEnd.Add(tipOffset))
		imd.Line(3)
	}
	
	// マント
	capeAlpha := uint8(150)
	if mk.CurrentAbility.GetName() == "Cape Barrier" {
		capeAlpha = 255
	}
	capeColor := color.RGBA{R: 150, G: 100, B: 200, A: capeAlpha}
	imd.Color = capeColor
	
	// マントの形状（動きに応じて揺れる）
	capeSwing := math.Sin(mk.AnimationTime*5) * 5
	capeLeft := mk.Position.Add(pixel.V(-mk.Radius*0.9+capeSwing, mk.Radius*0.3))
	capeRight := mk.Position.Add(pixel.V(mk.Radius*0.9+capeSwing, mk.Radius*0.3))
	capeBottom := mk.Position.Add(pixel.V(capeSwing, -mk.Radius*1.2))
	
	imd.Push(capeLeft)
	imd.Push(capeRight)
	imd.Push(capeBottom)
	imd.Polygon(0)
	
	// トルネードエフェクト（トルネードアビリティ使用時）
	if mk.CurrentAbility.GetName() == "Tornado" && mk.AttackCooldown > 0 {
		tornadoColor := color.RGBA{R: 150, G: 255, B: 150, A: 200}
		imd.Color = tornadoColor
		for i := 0; i < 3; i++ {
			angle := mk.AnimationTime*15 + float64(i)*2*math.Pi/3
			offset := pixel.V(
				math.Cos(angle)*mk.Radius*1.3,
				math.Sin(angle)*mk.Radius*1.3,
			)
			imd.Push(mk.Position.Add(offset))
			imd.Circle(8, 0)
		}
	}
}

// TakeDamage はダメージを受けます
func (mk *MetaKnightPlayer) TakeDamage(damage int) {
	// マント防御中はダメージ軽減
	if mk.CurrentAbility != nil && mk.CurrentAbility.GetName() == "Cape Barrier" {
		damage = damage / 3
	}
	
	mk.Health -= damage
	if mk.Health <= 0 {
		mk.Health = 0
		mk.IsAlive = false
	}
}

// Heal は体力を回復します
func (mk *MetaKnightPlayer) Heal(amount int) {
	mk.Health += amount
	if mk.Health > mk.MaxHealth {
		mk.Health = mk.MaxHealth
	}
}

// ResetJump はジャンプ状態をリセットします（地面に着地した時）
func (mk *MetaKnightPlayer) ResetJump() {
	mk.IsJumping = false
	mk.HasDoubleJumped = false
}

// GetBounds は当たり判定用の矩形を返します
func (mk *MetaKnightPlayer) GetBounds() pixel.Rect {
	return pixel.R(
		mk.Position.X-mk.Radius,
		mk.Position.Y-mk.Radius,
		mk.Position.X+mk.Radius,
		mk.Position.Y+mk.Radius,
	)
}

// GetPosition は現在位置を返します
func (mk *MetaKnightPlayer) GetPosition() pixel.Vec {
	return mk.Position
}

// SetPosition は位置を設定します
func (mk *MetaKnightPlayer) SetPosition(pos pixel.Vec) {
	mk.Position = pos
}

// GetVelocity は速度を返します
func (mk *MetaKnightPlayer) GetVelocity() pixel.Vec {
	return mk.Velocity
}

// SetVelocity は速度を設定します
func (mk *MetaKnightPlayer) SetVelocity(vel pixel.Vec) {
	mk.Velocity = vel
}

// GetAbility は現在のアビリティを返します
func (mk *MetaKnightPlayer) GetAbility() ability.Ability {
	return mk.CurrentAbility
}

// SetAbility はアビリティを設定します（インターフェース実装）
func (mk *MetaKnightPlayer) SetAbility(a ability.Ability) {
	mk.CurrentAbility = a
}

// GetHealth は現在の体力を返します
func (mk *MetaKnightPlayer) GetHealth() int {
	return mk.Health
}

// GetMaxHealth は最大体力を返します
func (mk *MetaKnightPlayer) GetMaxHealth() int {
	return mk.MaxHealth
}

// GetAttackRange は攻撃範囲を返します
func (mk *MetaKnightPlayer) GetAttackRange() float64 {
	if mk.CurrentAbility != nil {
		switch mk.CurrentAbility.GetName() {
		case "Sword":
			return 60.0
		case "Tornado":
			return 80.0
		default:
			return 40.0
		}
	}
	return 40.0
}

// GetAttackDamage は攻撃ダメージを返します
func (mk *MetaKnightPlayer) GetAttackDamage() int {
	if mk.CurrentAbility != nil {
		switch mk.CurrentAbility.GetName() {
		case "Sword":
			return 20 + (mk.ComboCount * 5) // コンボで威力アップ
		case "Tornado":
			return 15
		default:
			return 10
		}
	}
	return 10
}
