package menu

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// GameState はゲームの状態を表します
type GameState int

const (
	StateTitleScreen GameState = iota
	StateCharacterSelect
	StateStageSelect
	StatePlaying
	StateGameOver
	StateStageComplete
)

// PlayerCharacter はプレイ可能なキャラクター
type PlayerCharacter int

const (
	CharacterKirby PlayerCharacter = iota
	CharacterMetaKnight
)

// MenuManager はメニュー全体を管理します
type MenuManager struct {
	State              GameState
	SelectedCharacter  PlayerCharacter
	SelectedStage      int
	Window             *pixelgl.Window
	Atlas              *text.Atlas
	IMDraw             *imdraw.IMDraw
	
	// メニュー選択
	titleSelection     int
	characterSelection int
	stageSelection     int
}

// NewMenuManager は新しいメニューマネージャーを作成します
func NewMenuManager(win *pixelgl.Window) *MenuManager {
	return &MenuManager{
		State:              StateTitleScreen,
		SelectedCharacter:  CharacterKirby,
		SelectedStage:      1,
		Window:             win,
		Atlas:              text.NewAtlas(basicfont.Face7x13, text.ASCII),
		IMDraw:             imdraw.New(nil),
		titleSelection:     0,
		characterSelection: 0,
		stageSelection:     0,
	}
}

// Update はメニューの状態を更新します
func (m *MenuManager) Update(dt float64) {
	switch m.State {
	case StateTitleScreen:
		m.updateTitleScreen()
	case StateCharacterSelect:
		m.updateCharacterSelect()
	case StateStageSelect:
		m.updateStageSelect()
	}
}

// updateTitleScreen はタイトル画面の更新処理
func (m *MenuManager) updateTitleScreen() {
	// 上下キーで選択
	if m.Window.JustPressed(pixelgl.KeyUp) || m.Window.JustPressed(pixelgl.KeyW) {
		m.titleSelection--
		if m.titleSelection < 0 {
			m.titleSelection = 1 // 0: Start, 1: Exit
		}
	}
	if m.Window.JustPressed(pixelgl.KeyDown) || m.Window.JustPressed(pixelgl.KeyS) {
		m.titleSelection++
		if m.titleSelection > 1 {
			m.titleSelection = 0
		}
	}
	
	// Enterで決定
	if m.Window.JustPressed(pixelgl.KeyEnter) || m.Window.JustPressed(pixelgl.KeySpace) {
		if m.titleSelection == 0 {
			m.State = StateCharacterSelect
		} else {
			m.Window.SetClosed(true)
		}
	}
}

// updateCharacterSelect はキャラクター選択画面の更新処理
func (m *MenuManager) updateCharacterSelect() {
	// 左右キーで選択
	if m.Window.JustPressed(pixelgl.KeyLeft) || m.Window.JustPressed(pixelgl.KeyA) {
		m.characterSelection--
		if m.characterSelection < 0 {
			m.characterSelection = 1
		}
	}
	if m.Window.JustPressed(pixelgl.KeyRight) || m.Window.JustPressed(pixelgl.KeyD) {
		m.characterSelection++
		if m.characterSelection > 1 {
			m.characterSelection = 0
		}
	}
	
	// Enterで決定
	if m.Window.JustPressed(pixelgl.KeyEnter) || m.Window.JustPressed(pixelgl.KeySpace) {
		m.SelectedCharacter = PlayerCharacter(m.characterSelection)
		m.State = StateStageSelect
	}
	
	// ESCで戻る
	if m.Window.JustPressed(pixelgl.KeyEscape) {
		m.State = StateTitleScreen
	}
}

// updateStageSelect はステージ選択画面の更新処理
func (m *MenuManager) updateStageSelect() {
	// 左右キーで選択
	if m.Window.JustPressed(pixelgl.KeyLeft) || m.Window.JustPressed(pixelgl.KeyA) {
		m.stageSelection--
		if m.stageSelection < 0 {
			m.stageSelection = 1 // ステージ1, 2
		}
	}
	if m.Window.JustPressed(pixelgl.KeyRight) || m.Window.JustPressed(pixelgl.KeyD) {
		m.stageSelection++
		if m.stageSelection > 1 {
			m.stageSelection = 0
		}
	}
	
	// Enterで決定
	if m.Window.JustPressed(pixelgl.KeyEnter) || m.Window.JustPressed(pixelgl.KeySpace) {
		m.SelectedStage = m.stageSelection + 1
		m.State = StatePlaying
	}
	
	// ESCで戻る
	if m.Window.JustPressed(pixelgl.KeyEscape) {
		m.State = StateCharacterSelect
	}
}

// Draw はメニューを描画します
func (m *MenuManager) Draw() {
	m.Window.Clear(colornames.Black)
	m.IMDraw.Clear()
	
	switch m.State {
	case StateTitleScreen:
		m.drawTitleScreen()
	case StateCharacterSelect:
		m.drawCharacterSelect()
	case StateStageSelect:
		m.drawStageSelect()
	}
	
	m.IMDraw.Draw(m.Window)
}

// drawTitleScreen はタイトル画面を描画
func (m *MenuManager) drawTitleScreen() {
	width := m.Window.Bounds().W()
	height := m.Window.Bounds().H()
	
	// 背景のグラデーション
	m.IMDraw.Color = color.RGBA{R: 135, G: 206, B: 250, A: 255}
	m.IMDraw.Push(pixel.V(0, 0))
	m.IMDraw.Push(pixel.V(width, height))
	m.IMDraw.Rectangle(0)
	
	// タイトルロゴ
	titleText := text.New(pixel.V(width/2-150, height-150), m.Atlas)
	titleText.Color = colornames.White
	fmt.Fprintf(titleText, "KIRBY ADVENTURE")
	titleText.Draw(m.Window, pixel.IM.Scaled(titleText.Orig, 4))
	
	// サブタイトル
	subText := text.New(pixel.V(width/2-100, height-200), m.Atlas)
	subText.Color = colornames.Yellow
	fmt.Fprintf(subText, "Inspired RPG")
	subText.Draw(m.Window, pixel.IM.Scaled(subText.Orig, 2))
	
	// メニュー項目
	menuY := height / 2
	
	// Start Game
	startColor := colornames.White
	if m.titleSelection == 0 {
		startColor = colornames.Yellow
		// 選択カーソル
		m.IMDraw.Color = colornames.Yellow
		m.IMDraw.Push(pixel.V(width/2-100, menuY+5))
		m.IMDraw.Push(pixel.V(width/2-90, menuY+5))
		m.IMDraw.Line(3)
	}
	startText := text.New(pixel.V(width/2-70, menuY), m.Atlas)
	startText.Color = startColor
	fmt.Fprintf(startText, "START GAME")
	startText.Draw(m.Window, pixel.IM.Scaled(startText.Orig, 3))
	
	// Exit
	exitY := menuY - 60
	exitColor := colornames.White
	if m.titleSelection == 1 {
		exitColor = colornames.Yellow
		m.IMDraw.Color = colornames.Yellow
		m.IMDraw.Push(pixel.V(width/2-100, exitY+5))
		m.IMDraw.Push(pixel.V(width/2-90, exitY+5))
		m.IMDraw.Line(3)
	}
	exitText := text.New(pixel.V(width/2-70, exitY), m.Atlas)
	exitText.Color = exitColor
	fmt.Fprintf(exitText, "EXIT")
	exitText.Draw(m.Window, pixel.IM.Scaled(exitText.Orig, 3))
	
	// 操作説明
	instructionText := text.New(pixel.V(width/2-150, 50), m.Atlas)
	instructionText.Color = colornames.White
	fmt.Fprintf(instructionText, "UP/DOWN: Select  ENTER: Confirm")
	instructionText.Draw(m.Window, pixel.IM.Scaled(instructionText.Orig, 1.5))
}

// drawCharacterSelect はキャラクター選択画面を描画
func (m *MenuManager) drawCharacterSelect() {
	width := m.Window.Bounds().W()
	height := m.Window.Bounds().H()
	
	// 背景
	m.IMDraw.Color = color.RGBA{R: 50, G: 50, B: 80, A: 255}
	m.IMDraw.Push(pixel.V(0, 0))
	m.IMDraw.Push(pixel.V(width, height))
	m.IMDraw.Rectangle(0)
	
	// タイトル
	titleText := text.New(pixel.V(width/2-120, height-100), m.Atlas)
	titleText.Color = colornames.Yellow
	fmt.Fprintf(titleText, "SELECT CHARACTER")
	titleText.Draw(m.Window, pixel.IM.Scaled(titleText.Orig, 3))
	
	// カービィ
	kirbyX := width/2 - 200
	kirbyY := height / 2
	
	if m.characterSelection == 0 {
		// 選択枠
		m.IMDraw.Color = colornames.Yellow
		m.IMDraw.Push(pixel.V(kirbyX-60, kirbyY-80))
		m.IMDraw.Push(pixel.V(kirbyX+60, kirbyY+80))
		m.IMDraw.Rectangle(5)
	}
	
	// カービィのプレビュー（ピンクの丸）
	m.IMDraw.Color = color.RGBA{R: 255, G: 182, B: 193, A: 255}
	m.IMDraw.Push(pixel.V(kirbyX, kirbyY))
	m.IMDraw.Circle(40, 0)
	
	kirbyNameText := text.New(pixel.V(kirbyX-30, kirbyY-120), m.Atlas)
	kirbyNameText.Color = colornames.White
	fmt.Fprintf(kirbyNameText, "KIRBY")
	kirbyNameText.Draw(m.Window, pixel.IM.Scaled(kirbyNameText.Orig, 2))
	
	// メタナイト
	metaX := width/2 + 200
	metaY := height / 2
	
	if m.characterSelection == 1 {
		// 選択枠
		m.IMDraw.Color = colornames.Yellow
		m.IMDraw.Push(pixel.V(metaX-60, metaY-80))
		m.IMDraw.Push(pixel.V(metaX+60, metaY+80))
		m.IMDraw.Rectangle(5)
	}
	
	// メタナイトのプレビュー（紫の球体とマスク）
	m.IMDraw.Color = color.RGBA{R: 100, G: 50, B: 150, A: 255}
	m.IMDraw.Push(pixel.V(metaX, metaY))
	m.IMDraw.Circle(40, 0)
	
	// マスク
	m.IMDraw.Color = color.RGBA{R: 200, G: 200, B: 220, A: 255}
	m.IMDraw.Push(pixel.V(metaX-30, metaY))
	m.IMDraw.Push(pixel.V(metaX+30, metaY+10))
	m.IMDraw.Rectangle(0)
	
	// 目
	m.IMDraw.Color = colornames.Yellow
	m.IMDraw.Push(pixel.V(metaX, metaY+5))
	m.IMDraw.Circle(8, 0)
	
	metaNameText := text.New(pixel.V(metaX-50, metaY-120), m.Atlas)
	metaNameText.Color = colornames.White
	fmt.Fprintf(metaNameText, "META KNIGHT")
	metaNameText.Draw(m.Window, pixel.IM.Scaled(metaNameText.Orig, 2))
	
	// 操作説明
	instructionText := text.New(pixel.V(width/2-200, 50), m.Atlas)
	instructionText.Color = colornames.White
	fmt.Fprintf(instructionText, "LEFT/RIGHT: Select  ENTER: Confirm  ESC: Back")
	instructionText.Draw(m.Window, pixel.IM.Scaled(instructionText.Orig, 1.5))
}

// drawStageSelect はステージ選択画面を描画
func (m *MenuManager) drawStageSelect() {
	width := m.Window.Bounds().W()
	height := m.Window.Bounds().H()
	
	// 背景
	m.IMDraw.Color = color.RGBA{R: 30, G: 30, B: 50, A: 255}
	m.IMDraw.Push(pixel.V(0, 0))
	m.IMDraw.Push(pixel.V(width, height))
	m.IMDraw.Rectangle(0)
	
	// タイトル
	titleText := text.New(pixel.V(width/2-100, height-100), m.Atlas)
	titleText.Color = colornames.Cyan
	fmt.Fprintf(titleText, "SELECT STAGE")
	titleText.Draw(m.Window, pixel.IM.Scaled(titleText.Orig, 3))
	
	// ステージ1
	stage1X := width/2 - 200
	stage1Y := height / 2
	
	if m.stageSelection == 0 {
		m.IMDraw.Color = colornames.Yellow
		m.IMDraw.Push(pixel.V(stage1X-80, stage1Y-60))
		m.IMDraw.Push(pixel.V(stage1X+80, stage1Y+60))
		m.IMDraw.Rectangle(5)
	}
	
	// ステージ1アイコン
	m.IMDraw.Color = color.RGBA{R: 100, G: 200, B: 100, A: 255}
	m.IMDraw.Push(pixel.V(stage1X-50, stage1Y-30))
	m.IMDraw.Push(pixel.V(stage1X+50, stage1Y+30))
	m.IMDraw.Rectangle(0)
	
	stage1Text := text.New(pixel.V(stage1X-40, stage1Y-80), m.Atlas)
	stage1Text.Color = colornames.White
	fmt.Fprintf(stage1Text, "STAGE 1")
	stage1Text.Draw(m.Window, pixel.IM.Scaled(stage1Text.Orig, 2))
	
	bossText1 := text.New(pixel.V(stage1X-60, stage1Y-110), m.Atlas)
	bossText1.Color = colornames.Red
	fmt.Fprintf(bossText1, "Boss: Dedede")
	bossText1.Draw(m.Window, pixel.IM.Scaled(bossText1.Orig, 1.5))
	
	// ステージ2
	stage2X := width/2 + 200
	stage2Y := height / 2
	
	if m.stageSelection == 1 {
		m.IMDraw.Color = colornames.Yellow
		m.IMDraw.Push(pixel.V(stage2X-80, stage2Y-60))
		m.IMDraw.Push(pixel.V(stage2X+80, stage2Y+60))
		m.IMDraw.Rectangle(5)
	}
	
	// ステージ2アイコン
	m.IMDraw.Color = color.RGBA{R: 100, G: 100, B: 200, A: 255}
	m.IMDraw.Push(pixel.V(stage2X-50, stage2Y-30))
	m.IMDraw.Push(pixel.V(stage2X+50, stage2Y+30))
	m.IMDraw.Rectangle(0)
	
	stage2Text := text.New(pixel.V(stage2X-40, stage2Y-80), m.Atlas)
	stage2Text.Color = colornames.White
	fmt.Fprintf(stage2Text, "STAGE 2")
	stage2Text.Draw(m.Window, pixel.IM.Scaled(stage2Text.Orig, 2))
	
	bossText2 := text.New(pixel.V(stage2X-70, stage2Y-110), m.Atlas)
	bossText2.Color = colornames.Purple
	fmt.Fprintf(bossText2, "Boss: Meta Knight")
	bossText2.Draw(m.Window, pixel.IM.Scaled(bossText2.Orig, 1.5))
	
	// 操作説明
	instructionText := text.New(pixel.V(width/2-200, 50), m.Atlas)
	instructionText.Color = colornames.White
	fmt.Fprintf(instructionText, "LEFT/RIGHT: Select  ENTER: Confirm  ESC: Back")
	instructionText.Draw(m.Window, pixel.IM.Scaled(instructionText.Orig, 1.5))
}

// GetCharacterName はキャラクター名を返します
func (m *MenuManager) GetCharacterName() string {
	switch m.SelectedCharacter {
	case CharacterKirby:
		return "Kirby"
	case CharacterMetaKnight:
		return "Meta Knight"
	default:
		return "Unknown"
	}
}
