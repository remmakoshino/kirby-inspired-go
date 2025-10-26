package entity

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// EnemyType は敵のタイプを表します
type EnemyType int

const (
	EnemyTypeWalker EnemyType = iota // 歩くタイプ
	EnemyTypeFlyer                   // 飛ぶタイプ
	EnemyTypeJumper                  // ジャンプするタイプ
)

// Enemy は敵キャラクターを表します
type Enemy struct {
	Position       pixel.Vec
	Velocity       pixel.Vec
	Radius         float64
	Health         int
	MaxHealth      int
	Type           EnemyType
	Color          color.RGBA
	IsAlive        bool
	MoveDirection  float64 // -1 (左) or 1 (右)
	
	// AI関連
	AITimer        float64
	PatrolDistance float64
	StartPosition  pixel.Vec
	
	// アニメーション
	AnimationTime  float64
}

// NewEnemy は新しい敵を作成します
func NewEnemy(pos pixel.Vec, enemyType EnemyType) *Enemy {
	e := &Enemy{
		Position:       pos,
		Velocity:       pixel.ZV,
		Radius:         15.0,
		Health:         30,
		MaxHealth:      30,
		Type:           enemyType,
		IsAlive:        true,
		MoveDirection:  1.0,
		AITimer:        0,
		PatrolDistance: 100.0,
		StartPosition:  pos,
		AnimationTime:  0,
	}
	
	// タイプ別の設定
	switch enemyType {
	case EnemyTypeWalker:
		e.Color = color.RGBA{R: 100, G: 200, B: 100, A: 255}
		e.Radius = 15.0
	case EnemyTypeFlyer:
		e.Color = color.RGBA{R: 200, G: 100, B: 200, A: 255}
		e.Radius = 12.0
	case EnemyTypeJumper:
		e.Color = color.RGBA{R: 200, G: 200, B: 100, A: 255}
		e.Radius = 18.0
	}
	
	return e
}

// Update は敵の状態を更新します
func (e *Enemy) Update(dt float64, playerPos pixel.Vec, stageWidth, stageHeight float64) {
	if !e.IsAlive {
		return
	}
	
	e.AnimationTime += dt
	e.AITimer += dt
	
	// タイプ別のAI
	switch e.Type {
	case EnemyTypeWalker:
		e.updateWalkerAI(dt, playerPos)
	case EnemyTypeFlyer:
		e.updateFlyerAI(dt, playerPos)
	case EnemyTypeJumper:
		e.updateJumperAI(dt, playerPos)
	}
	
	// 重力適用（飛行タイプ以外）
	if e.Type != EnemyTypeFlyer {
		e.Velocity.Y -= Gravity * dt
		if e.Velocity.Y < -MaxFallSpeed {
			e.Velocity.Y = -MaxFallSpeed
		}
	}
	
	// 位置更新
	e.Position = e.Position.Add(e.Velocity.Scaled(dt))
	
	// 地面との衝突（簡易版）
	if e.Position.Y-e.Radius <= 0 && e.Type != EnemyTypeFlyer {
		e.Position.Y = e.Radius
		e.Velocity.Y = 0
	}
	
	// 画面端処理
	if e.Position.X-e.Radius < 0 {
		e.Position.X = e.Radius
		e.MoveDirection = 1.0
	} else if e.Position.X+e.Radius > stageWidth {
		e.Position.X = stageWidth - e.Radius
		e.MoveDirection = -1.0
	}
	
	// 画面外に落ちた場合
	if e.Position.Y < -100 {
		e.IsAlive = false
	}
}

// updateWalkerAI は歩行タイプのAIを更新します
func (e *Enemy) updateWalkerAI(dt float64, playerPos pixel.Vec) {
	const walkSpeed = 50.0
	
	// パトロール
	distance := e.Position.X - e.StartPosition.X
	if math.Abs(distance) > e.PatrolDistance {
		e.MoveDirection *= -1
	}
	
	e.Velocity.X = walkSpeed * e.MoveDirection
}

// updateFlyerAI は飛行タイプのAIを更新します
func (e *Enemy) updateFlyerAI(dt float64, playerPos pixel.Vec) {
	const flySpeed = 60.0
	
	// プレイヤーに向かって緩やかに移動
	direction := playerPos.Sub(e.Position).Unit()
	
	// 一定距離以上離れている場合のみ追跡
	if playerPos.Sub(e.Position).Len() > 100 {
		e.Velocity = direction.Scaled(flySpeed * 0.5)
	} else {
		// サインカーブで上下に動く
		e.Velocity.X = flySpeed * e.MoveDirection
		e.Velocity.Y = math.Sin(e.AnimationTime*2) * 30
	}
	
	// 定期的に方向転換
	if e.AITimer > 3.0 {
		e.MoveDirection *= -1
		e.AITimer = 0
	}
}

// updateJumperAI はジャンプタイプのAIを更新します
func (e *Enemy) updateJumperAI(dt float64, playerPos pixel.Vec) {
	const jumpSpeed = 40.0
	
	// 地面にいる時のみジャンプ
	if e.Position.Y-e.Radius <= 1 {
		e.Velocity.X = jumpSpeed * e.MoveDirection
		
		// 定期的にジャンプ
		if e.AITimer > 2.0 {
			e.Velocity.Y = 300.0
			e.AITimer = 0
			
			// たまに方向転換
			if rand.Float64() < 0.3 {
				e.MoveDirection *= -1
			}
		}
	}
}

// Draw は敵を描画します
func (e *Enemy) Draw(imd *imdraw.IMDraw) {
	if !e.IsAlive {
		return
	}
	
	// タイプ別の描画
	switch e.Type {
	case EnemyTypeWalker:
		e.drawWalker(imd)
	case EnemyTypeFlyer:
		e.drawFlyer(imd)
	case EnemyTypeJumper:
		e.drawJumper(imd)
	}
}

// drawWalker は歩行タイプの敵を描画します
func (e *Enemy) drawWalker(imd *imdraw.IMDraw) {
	// 本体
	imd.Color = e.Color
	imd.Push(e.Position)
	imd.Circle(e.Radius, 0)
	
	// 目
	imd.Color = color.RGBA{R: 50, G: 50, B: 50, A: 255}
	eyeY := e.Position.Y + e.Radius*0.3
	imd.Push(pixel.V(e.Position.X-e.Radius*0.4, eyeY))
	imd.Circle(e.Radius*0.15, 0)
	imd.Push(pixel.V(e.Position.X+e.Radius*0.4, eyeY))
	imd.Circle(e.Radius*0.15, 0)
}

// drawFlyer は飛行タイプの敵を描画します
func (e *Enemy) drawFlyer(imd *imdraw.IMDraw) {
	// 本体
	imd.Color = e.Color
	imd.Push(e.Position)
	imd.Circle(e.Radius, 0)
	
	// 羽
	wingColor := color.RGBA{R: 220, G: 150, B: 220, A: 200}
	imd.Color = wingColor
	
	wingOffset := math.Sin(e.AnimationTime*10) * e.Radius * 0.3
	leftWing := e.Position.Add(pixel.V(-e.Radius*1.2, wingOffset))
	rightWing := e.Position.Add(pixel.V(e.Radius*1.2, wingOffset))
	
	imd.Push(leftWing)
	imd.Circle(e.Radius*0.6, 0)
	imd.Push(rightWing)
	imd.Circle(e.Radius*0.6, 0)
	
	// 目
	imd.Color = color.White
	imd.Push(pixel.V(e.Position.X-e.Radius*0.3, e.Position.Y+e.Radius*0.2))
	imd.Circle(e.Radius*0.2, 0)
	imd.Push(pixel.V(e.Position.X+e.Radius*0.3, e.Position.Y+e.Radius*0.2))
	imd.Circle(e.Radius*0.2, 0)
}

// drawJumper はジャンプタイプの敵を描画します
func (e *Enemy) drawJumper(imd *imdraw.IMDraw) {
	// 本体（少し縦長）
	imd.Color = e.Color
	
	imd.Push(e.Position)
	imd.Circle(e.Radius*0.8, 0)
	
	// 目
	imd.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	imd.Push(pixel.V(e.Position.X-e.Radius*0.3, e.Position.Y+e.Radius*0.4))
	imd.Circle(e.Radius*0.2, 0)
	imd.Push(pixel.V(e.Position.X+e.Radius*0.3, e.Position.Y+e.Radius*0.4))
	imd.Circle(e.Radius*0.2, 0)
	
	// 瞳
	imd.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	imd.Push(pixel.V(e.Position.X-e.Radius*0.3, e.Position.Y+e.Radius*0.3))
	imd.Circle(e.Radius*0.1, 0)
	imd.Push(pixel.V(e.Position.X+e.Radius*0.3, e.Position.Y+e.Radius*0.3))
	imd.Circle(e.Radius*0.1, 0)
}

// TakeDamage はダメージを受けます
func (e *Enemy) TakeDamage(damage int) {
	e.Health -= damage
	if e.Health <= 0 {
		e.Health = 0
		e.IsAlive = false
	}
}

// GetBounds は当たり判定用の矩形を返します
func (e *Enemy) GetBounds() pixel.Rect {
	return pixel.R(
		e.Position.X-e.Radius,
		e.Position.Y-e.Radius,
		e.Position.X+e.Radius,
		e.Position.Y+e.Radius,
	)
}

// GetAbilityType は敵が持つ能力のタイプを返します
func (e *Enemy) GetAbilityType() string {
	switch e.Type {
	case EnemyTypeWalker:
		return "speed"
	case EnemyTypeFlyer:
		return "fly"
	case EnemyTypeJumper:
		return "jump"
	default:
		return "none"
	}
}
