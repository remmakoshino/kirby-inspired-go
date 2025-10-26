package ability

import (
	"image/color"
)

// SpeedAbility はスピードアップ能力です
type SpeedAbility struct {
	BaseAbility
}

// NewSpeedAbility は新しいスピード能力を作成します
func NewSpeedAbility() *SpeedAbility {
	return &SpeedAbility{
		BaseAbility: BaseAbility{
			Name:     "Speed",
			Color:    color.RGBA{R: 255, G: 200, B: 100, A: 255},
			Cooldown: 0.5,
		},
	}
}

// Use はスピード能力を使用します
func (a *SpeedAbility) Use(player AbilityUser) {
	if !a.IsReady() {
		return
	}
	
	// 現在の速度を1.5倍に
	vel := player.GetVelocity()
	player.SetVelocity(vel.Scaled(1.5))
	
	a.StartCooldown()
}

// FlyAbility は飛行能力です
type FlyAbility struct {
	BaseAbility
	IsFlying bool
}

// NewFlyAbility は新しい飛行能力を作成します
func NewFlyAbility() *FlyAbility {
	return &FlyAbility{
		BaseAbility: BaseAbility{
			Name:     "Fly",
			Color:    color.RGBA{R: 200, G: 150, B: 255, A: 255},
			Cooldown: 0.1,
		},
		IsFlying: false,
	}
}

// Use は飛行能力を使用します（上昇）
func (a *FlyAbility) Use(player AbilityUser) {
	if !a.IsReady() {
		return
	}
	
	// 上向きの力を加える
	vel := player.GetVelocity()
	vel.Y = 200 // 上昇力
	player.SetVelocity(vel)
	
	a.IsFlying = true
	a.StartCooldown()
}

// JumpAbility は強化ジャンプ能力です
type JumpAbility struct {
	BaseAbility
}

// NewJumpAbility は新しいジャンプ能力を作成します
func NewJumpAbility() *JumpAbility {
	return &JumpAbility{
		BaseAbility: BaseAbility{
			Name:     "High Jump",
			Color:    color.RGBA{R: 255, G: 255, B: 100, A: 255},
			Cooldown: 1.0,
		},
	}
}

// Use はジャンプ能力を使用します
func (a *JumpAbility) Use(player AbilityUser) {
	if !a.IsReady() {
		return
	}
	
	// 強力なジャンプ
	vel := player.GetVelocity()
	vel.Y = 500 // 通常より高いジャンプ力
	player.SetVelocity(vel)
	
	a.StartCooldown()
}

// CreateAbility は能力タイプから能力を作成します
func CreateAbility(abilityType string) Ability {
	switch abilityType {
	case "speed":
		return NewSpeedAbility()
	case "fly":
		return NewFlyAbility()
	case "jump":
		return NewJumpAbility()
	default:
		return nil
	}
}
