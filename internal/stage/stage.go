package stage

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Platform はプラットフォームを表します
type Platform struct {
	Rect  pixel.Rect
	Color color.RGBA
}

// NewPlatform は新しいプラットフォームを作成します
func NewPlatform(x, y, width, height float64) *Platform {
	return &Platform{
		Rect:  pixel.R(x, y, x+width, y+height),
		Color: color.RGBA{R: 100, G: 150, B: 100, A: 255},
	}
}

// Draw はプラットフォームを描画します
func (p *Platform) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.Color
	imd.Push(p.Rect.Min, p.Rect.Max)
	imd.Rectangle(0)
}

// Stage はステージ全体を表します
type Stage struct {
	Width      float64
	Height     float64
	Platforms  []*Platform
	Background color.RGBA
}

// NewStage は新しいステージを作成します
func NewStage(width, height float64) *Stage {
	return &Stage{
		Width:      width,
		Height:     height,
		Platforms:  make([]*Platform, 0),
		Background: color.RGBA{R: 135, G: 206, B: 235, A: 255}, // 空色
	}
}

// AddPlatform はプラットフォームを追加します
func (s *Stage) AddPlatform(platform *Platform) {
	s.Platforms = append(s.Platforms, platform)
}

// Draw はステージを描画します
func (s *Stage) Draw(imd *imdraw.IMDraw) {
	// 背景は別途描画されるため、ここではプラットフォームのみ
	for _, platform := range s.Platforms {
		platform.Draw(imd)
	}
	
	// 地面
	ground := pixel.R(0, 0, s.Width, 5)
	imd.Color = color.RGBA{R: 139, G: 69, B: 19, A: 255} // 茶色
	imd.Push(ground.Min, ground.Max)
	imd.Rectangle(0)
}

// CheckCollision はエンティティとプラットフォームの衝突をチェックします
func (s *Stage) CheckCollision(entityBounds pixel.Rect, velocity pixel.Vec) (pixel.Vec, bool) {
	newPos := velocity
	collided := false
	
	for _, platform := range s.Platforms {
		if entityBounds.Intersects(platform.Rect) {
			// 上からの衝突（プラットフォームに着地）
			if velocity.Y < 0 && entityBounds.Min.Y > platform.Rect.Max.Y-10 {
				newPos.Y = platform.Rect.Max.Y - entityBounds.Min.Y
				collided = true
			}
			// 下からの衝突
			if velocity.Y > 0 && entityBounds.Max.Y < platform.Rect.Min.Y+10 {
				newPos.Y = platform.Rect.Min.Y - entityBounds.Max.Y
			}
			// 横からの衝突
			if velocity.X > 0 {
				newPos.X = platform.Rect.Min.X - entityBounds.Max.X
			} else if velocity.X < 0 {
				newPos.X = platform.Rect.Max.X - entityBounds.Min.X
			}
		}
	}
	
	return newPos, collided
}

// CreateDefaultStage はデフォルトのステージを作成します
func CreateDefaultStage(width, height float64) *Stage {
	stage := NewStage(width, height)
	
	// プラットフォームを配置
	stage.AddPlatform(NewPlatform(100, 100, 200, 20))
	stage.AddPlatform(NewPlatform(400, 150, 200, 20))
	stage.AddPlatform(NewPlatform(700, 100, 200, 20))
	stage.AddPlatform(NewPlatform(250, 250, 150, 20))
	stage.AddPlatform(NewPlatform(550, 300, 150, 20))
	stage.AddPlatform(NewPlatform(150, 400, 200, 20))
	stage.AddPlatform(NewPlatform(500, 450, 250, 20))
	
	return stage
}
