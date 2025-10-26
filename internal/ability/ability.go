package ability

import (
	"image/color"

	"github.com/faiface/pixel"
)

// Ability はコピー能力のインターフェースです
type Ability interface {
	Use(player AbilityUser)
	GetName() string
	GetColor() color.RGBA
	Update(dt float64)
}

// AbilityUser は能力を使用できるエンティティのインターフェース
type AbilityUser interface {
	GetPosition() pixel.Vec
	GetVelocity() pixel.Vec
	SetVelocity(v pixel.Vec)
}

// BaseAbility は能力の基本構造体
type BaseAbility struct {
	Name     string
	Color    color.RGBA
	Cooldown float64
	CurrentCooldown float64
}

// Update は能力の状態を更新します
func (a *BaseAbility) Update(dt float64) {
	if a.CurrentCooldown > 0 {
		a.CurrentCooldown -= dt
		if a.CurrentCooldown < 0 {
			a.CurrentCooldown = 0
		}
	}
}

// IsReady は能力が使用可能かを返します
func (a *BaseAbility) IsReady() bool {
	return a.CurrentCooldown <= 0
}

// StartCooldown はクールダウンを開始します
func (a *BaseAbility) StartCooldown() {
	a.CurrentCooldown = a.Cooldown
}

// GetName は能力の名前を返します
func (a *BaseAbility) GetName() string {
	return a.Name
}

// GetColor は能力の色を返します
func (a *BaseAbility) GetColor() color.RGBA {
	return a.Color
}
