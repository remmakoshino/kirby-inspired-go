package ability

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
)

// InhaleAbility は吸い込み能力です
type InhaleAbility struct {
	BaseAbility
	IsInhaling  bool
	InhaleRange float64
	InhaleForce float64
	InhaleTime  float64
}

// NewInhaleAbility は新しい吸い込み能力を作成します
func NewInhaleAbility() *InhaleAbility {
	return &InhaleAbility{
		BaseAbility: BaseAbility{
			Name:     "Inhale",
			Color:    color.RGBA{R: 255, G: 182, B: 193, A: 255},
			Cooldown: 0.5,
		},
		IsInhaling:  false,
		InhaleRange: 80.0,
		InhaleForce: 300.0,
		InhaleTime:  0,
	}
}

// Use は吸い込み能力を使用します
func (a *InhaleAbility) Use(player AbilityUser) {
	if !a.IsReady() {
		return
	}
	
	a.IsInhaling = true
	a.InhaleTime = 2.0 // 2秒間吸い込み可能
	a.StartCooldown()
}

// Update は吸い込みアニメーションの更新
func (a *InhaleAbility) Update(dt float64) {
	a.BaseAbility.Update(dt)
	
	if a.IsInhaling {
		a.InhaleTime -= dt
		if a.InhaleTime <= 0 {
			a.IsInhaling = false
			a.InhaleTime = 0
		}
	}
}

// StopInhale は吸い込みを停止します
func (a *InhaleAbility) StopInhale() {
	a.IsInhaling = false
	a.InhaleTime = 0
}

// HammerAbility はハンマー能力です
type HammerAbility struct {
	BaseAbility
	AttackRange  float64
	AttackDamage int
}

// NewHammerAbility は新しいハンマー能力を作成します
func NewHammerAbility() *HammerAbility {
	return &HammerAbility{
		BaseAbility: BaseAbility{
			Name:     "Hammer",
			Color:    color.RGBA{R: 255, G: 150, B: 50, A: 255},
			Cooldown: 0.8,
		},
		AttackRange:  50.0,
		AttackDamage: 30,
	}
}

// Use はハンマー攻撃を使用します
func (a *HammerAbility) Use(player AbilityUser) {
	if !a.IsReady() {
		return
	}
	
	a.StartCooldown()
}

// SwordAbility は剣能力です
type SwordAbility struct {
	BaseAbility
	AttackRange  float64
	AttackDamage int
	ComboCount   int
}

// NewSwordAbility は新しい剣能力を作成します
func NewSwordAbility() *SwordAbility {
	return &SwordAbility{
		BaseAbility: BaseAbility{
			Name:     "Sword",
			Color:    color.RGBA{R: 200, G: 200, B: 255, A: 255},
			Cooldown: 0.4,
		},
		AttackRange:  60.0,
		AttackDamage: 20,
		ComboCount:   0,
	}
}

// Use は剣攻撃を使用します
func (a *SwordAbility) Use(player AbilityUser) {
	if !a.IsReady() {
		return
	}
	
	a.ComboCount++
	if a.ComboCount > 3 {
		a.ComboCount = 1
	}
	
	a.StartCooldown()
}

// TornadoAbility はトルネード突進能力です
type TornadoAbility struct {
	BaseAbility
	IsActive     bool
	Duration     float64
	RemainingTime float64
	Speed        float64
}

// NewTornadoAbility は新しいトルネード能力を作成します
func NewTornadoAbility() *TornadoAbility {
	return &TornadoAbility{
		BaseAbility: BaseAbility{
			Name:     "Tornado",
			Color:    color.RGBA{R: 150, G: 255, B: 150, A: 255},
			Cooldown: 2.0,
		},
		IsActive:     false,
		Duration:     1.5,
		RemainingTime: 0,
		Speed:        400.0,
	}
}

// Use はトルネード突進を使用します
func (a *TornadoAbility) Use(player AbilityUser) {
	if !a.IsReady() || a.IsActive {
		return
	}
	
	a.IsActive = true
	a.RemainingTime = a.Duration
	a.StartCooldown()
}

// Update はトルネードの状態を更新
func (a *TornadoAbility) Update(dt float64) {
	a.BaseAbility.Update(dt)
	
	if a.IsActive {
		a.RemainingTime -= dt
		if a.RemainingTime <= 0 {
			a.IsActive = false
			a.RemainingTime = 0
		}
	}
}

// CapeBarrierAbility はマントバリア能力です
type CapeBarrierAbility struct {
	BaseAbility
	IsActive     bool
	Duration     float64
	RemainingTime float64
}

// NewCapeBarrierAbility は新しいマントバリア能力を作成します
func NewCapeBarrierAbility() *CapeBarrierAbility {
	return &CapeBarrierAbility{
		BaseAbility: BaseAbility{
			Name:     "Cape Barrier",
			Color:    color.RGBA{R: 200, G: 150, B: 255, A: 255},
			Cooldown: 3.0,
		},
		IsActive:     false,
		Duration:     2.0,
		RemainingTime: 0,
	}
}

// Use はマントバリアを使用します
func (a *CapeBarrierAbility) Use(player AbilityUser) {
	if !a.IsReady() || a.IsActive {
		return
	}
	
	a.IsActive = true
	a.RemainingTime = a.Duration
	a.StartCooldown()
}

// Update はバリアの状態を更新
func (a *CapeBarrierAbility) Update(dt float64) {
	a.BaseAbility.Update(dt)
	
	if a.IsActive {
		a.RemainingTime -= dt
		if a.RemainingTime <= 0 {
			a.IsActive = false
			a.RemainingTime = 0
		}
	}
}

// CreateAbilityFromType は能力タイプから能力を作成します（更新版）
func CreateAbilityFromType(abilityType string) Ability {
	switch abilityType {
	case "speed":
		return NewSpeedAbility()
	case "fly":
		return NewFlyAbility()
	case "jump":
		return NewJumpAbility()
	case "inhale":
		return NewInhaleAbility()
	case "hammer":
		return NewHammerAbility()
	case "sword":
		return NewSwordAbility()
	case "tornado":
		return NewTornadoAbility()
	case "cape":
		return NewCapeBarrierAbility()
	default:
		return nil
	}
}

// InhaleEffect は吸い込みエフェクトを計算します
func InhaleEffect(playerPos, enemyPos pixel.Vec, inhaleForce, inhaleRange float64) pixel.Vec {
	distance := playerPos.Sub(enemyPos).Len()
	
	if distance > inhaleRange {
		return pixel.ZV
	}
	
	// 距離に応じて吸い込む力を調整
	strength := inhaleForce * (1.0 - distance/inhaleRange)
	direction := playerPos.Sub(enemyPos).Unit()
	
	return direction.Scaled(strength)
}

// IsInInhaleRange は敵が吸い込み範囲内かチェック
func IsInInhaleRange(playerPos, enemyPos pixel.Vec, inhaleRange, playerFacing float64) bool {
	diff := enemyPos.Sub(playerPos)
	distance := diff.Len()
	
	if distance > inhaleRange {
		return false
	}
	
	// プレイヤーの向きに応じて判定
	if playerFacing > 0 && diff.X < 0 {
		return false // 右向きだが敵が左にいる
	}
	if playerFacing < 0 && diff.X > 0 {
		return false // 左向きだが敵が右にいる
	}
	
	// 角度チェック（前方120度の範囲）
	angle := math.Atan2(diff.Y, diff.X)
	facingAngle := 0.0
	if playerFacing < 0 {
		facingAngle = math.Pi
	}
	
	angleDiff := math.Abs(angle - facingAngle)
	if angleDiff > math.Pi {
		angleDiff = 2*math.Pi - angleDiff
	}
	
	return angleDiff < math.Pi/3 // 60度以内
}
