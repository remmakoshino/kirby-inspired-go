package entity

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// WaddleDee はワドルディ（オレンジの敵）
type WaddleDee struct {
	*Enemy
}

// NewWaddleDee は新しいワドルディを作成します
func NewWaddleDee(pos pixel.Vec) *WaddleDee {
	enemy := &Enemy{
		Position:       pos,
		Velocity:       pixel.ZV,
		Radius:         15.0,
		Health:         20,
		MaxHealth:      20,
		Type:           EnemyTypeWalker,
		Color:          color.RGBA{R: 255, G: 140, B: 60, A: 255},
		IsAlive:        true,
		MoveDirection:  1.0,
		AITimer:        0,
		PatrolDistance: 100.0,
		StartPosition:  pos,
		AnimationTime:  0,
	}
	
	return &WaddleDee{Enemy: enemy}
}

// WaddleDoo はワドルドゥ（目玉の敵）
type WaddleDoo struct {
	*Enemy
	ShootTimer    float64
	ShootCooldown float64
}

// NewWaddleDoo は新しいワドルドゥを作成します
func NewWaddleDoo(pos pixel.Vec) *WaddleDoo {
	enemy := &Enemy{
		Position:       pos,
		Velocity:       pixel.ZV,
		Radius:         16.0,
		Health:         25,
		MaxHealth:      25,
		Type:           EnemyTypeWalker,
		Color:          color.RGBA{R: 255, G: 100, B: 50, A: 255},
		IsAlive:        true,
		MoveDirection:  1.0,
		AITimer:        0,
		PatrolDistance: 80.0,
		StartPosition:  pos,
		AnimationTime:  0,
	}
	
	return &WaddleDoo{
		Enemy:         enemy,
		ShootTimer:    0,
		ShootCooldown: 3.0,
	}
}

// Draw はワドルディを描画します
func (wd *WaddleDee) Draw(imd *imdraw.IMDraw) {
	if !wd.IsAlive {
		return
	}
	
	// 本体（オレンジの丸）
	imd.Color = wd.Color
	imd.Push(wd.Position)
	imd.Circle(wd.Radius, 0)
	
	// 足（2本）
	footColor := color.RGBA{R: 200, G: 80, B: 20, A: 255}
	imd.Color = footColor
	
	leftFootPos := wd.Position.Add(pixel.V(-wd.Radius*0.4, -wd.Radius*0.9))
	rightFootPos := wd.Position.Add(pixel.V(wd.Radius*0.4, -wd.Radius*0.9))
	
	imd.Push(leftFootPos)
	imd.Circle(wd.Radius*0.25, 0)
	imd.Push(rightFootPos)
	imd.Circle(wd.Radius*0.25, 0)
	
	// 目（大きい黒目）
	imd.Color = color.RGBA{R: 20, G: 20, B: 20, A: 255}
	eyeY := wd.Position.Y + wd.Radius*0.2
	imd.Push(pixel.V(wd.Position.X-wd.Radius*0.3, eyeY))
	imd.Circle(wd.Radius*0.2, 0)
	imd.Push(pixel.V(wd.Position.X+wd.Radius*0.3, eyeY))
	imd.Circle(wd.Radius*0.2, 0)
	
	// 口（笑顔）
	imd.Color = color.RGBA{R: 150, G: 50, B: 20, A: 255}
	for i := 0; i <= 5; i++ {
		angle := math.Pi * float64(i) / 5.0
		x := wd.Position.X + wd.Radius*0.3*math.Cos(angle)*0.5
		y := wd.Position.Y - wd.Radius*0.3 - wd.Radius*0.2*math.Sin(angle)*0.3
		imd.Push(pixel.V(x, y))
	}
	imd.Polygon(2)
}

// Draw はワドルドゥを描画します
func (wd *WaddleDoo) Draw(imd *imdraw.IMDraw) {
	if !wd.IsAlive {
		return
	}
	
	// 本体（赤オレンジの丸）
	imd.Color = wd.Color
	imd.Push(wd.Position)
	imd.Circle(wd.Radius, 0)
	
	// 大きな目（特徴的）
	// 白目
	imd.Color = color.White
	imd.Push(wd.Position)
	imd.Circle(wd.Radius*0.6, 0)
	
	// 黒目（大きな瞳）
	imd.Color = color.RGBA{R: 20, G: 20, B: 80, A: 255}
	imd.Push(wd.Position)
	imd.Circle(wd.Radius*0.4, 0)
	
	// 光の反射
	imd.Color = color.White
	highlightPos := wd.Position.Add(pixel.V(-wd.Radius*0.15, wd.Radius*0.15))
	imd.Push(highlightPos)
	imd.Circle(wd.Radius*0.15, 0)
	
	// 足
	footColor := color.RGBA{R: 200, G: 60, B: 30, A: 255}
	imd.Color = footColor
	
	leftFootPos := wd.Position.Add(pixel.V(-wd.Radius*0.4, -wd.Radius*0.9))
	rightFootPos := wd.Position.Add(pixel.V(wd.Radius*0.4, -wd.Radius*0.9))
	
	imd.Push(leftFootPos)
	imd.Circle(wd.Radius*0.25, 0)
	imd.Push(rightFootPos)
	imd.Circle(wd.Radius*0.25, 0)
}

// GetAbilityTypeForWaddleDee はワドルディの能力タイプを返します（削除）
func GetAbilityTypeForWaddleDee() string {
	return "none" // 能力なし
}

// GetAbilityTypeForWaddleDoo はワドルドゥの能力タイプを返します
func GetAbilityTypeForWaddleDoo() string {
	return "beam" // ビーム能力（後で実装）
}
