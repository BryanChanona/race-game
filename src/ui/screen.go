package render

import (
	"race-game/src/config"
	"race-game/src/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game implementa ebiten.Game
type Game struct {
	cfg     *config.GameConfig
	race    *core.Race
	assets  *Assets
	ui      *UI
	state   int // Estado del juego (Menu, Racing, etc.)
}

// NewGame crea una nueva instancia del juego
func NewGame(cfg *config.GameConfig) *Game {
	assets := NewAssets()
	ui := NewUI(assets)
	
	return &Game{
		cfg:    cfg,
		race:   nil,
		assets: assets,
		ui:     ui,
		state:  config.StateMenu,
	}
}

// Update actualiza la lógica del juego (llamado 60 veces por segundo)
func (g *Game) Update() error {
	switch g.state {
	case config.StateMenu:
		// Esperar input para comenzar
		if g.ui.CheckMenuInput() {
			g.startNewRace()
		}
		
	case config.StateCountdown:
		// Nada que hacer, el countdown se maneja en goroutines
		if g.race.IsRunning() {
			g.state = config.StateRacing
		}
		
	case config.StateRacing:
		// Actualizar la carrera
		g.race.Update()
		
		// Verificar si terminó
		if g.race.IsFinished() {
			g.state = config.StateFinished
		}
		
	case config.StateFinished:
		// Esperar input para reiniciar o salir
		if g.ui.CheckRestartInput() {
			g.startNewRace()
		}
		if g.ui.CheckExitInput() {
			return ebiten.Termination
		}
	}
	
	return nil
}

// Draw dibuja el juego en la pantalla
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case config.StateMenu:
		g.ui.DrawMenu(screen)
		
	case config.StateCountdown:
		g.drawRace(screen)
		g.ui.DrawCountdown(screen, g.race.GetCountdownValue())
		
	case config.StateRacing:
		g.drawRace(screen)
		g.ui.DrawHUD(screen, g.race)
		
	case config.StateFinished:
		g.drawRace(screen)
		g.ui.DrawResults(screen, g.race)
	}
}

// drawRace dibuja la carrera (pista y autos)
func (g *Game) drawRace(screen *ebiten.Image) {
	// Dibujar fondo de pista
	screen.DrawImage(g.assets.GetTrackBackground(), nil)
	
	// Dibujar línea de salida
	totalTrackHeight := config.NumLanes * config.LaneHeight
	offsetY := (config.ScreenHeight - totalTrackHeight) / 2
	
	opStart := &ebiten.DrawImageOptions{}
	opStart.GeoM.Translate(config.StartLineX, float64(offsetY))
	screen.DrawImage(g.assets.GetStartLine(), opStart)
	
	// Dibujar línea de meta
	opFinish := &ebiten.DrawImageOptions{}
	opFinish.GeoM.Translate(config.FinishLineX, float64(offsetY))
	screen.DrawImage(g.assets.GetFinishLine(), opFinish)
	
	// Dibujar autos (SECCIÓN CRÍTICA: lectura del estado compartido)
	g.drawCars(screen)
}

// drawCars dibuja todos los autos usando el estado compartido
func (g *Game) drawCars(screen *ebiten.Image) {
	if g.race == nil {
		return
	}
	
	track := g.race.GetManager().GetTrack()
	
	// Obtener todas las posiciones de forma thread-safe (RLock interno)
	positions := g.race.GetManager().GetSharedState().GetAllPositions()
	
	// Dibujar cada auto
	for id, state := range positions {
		sprite := g.assets.GetCarSprite(id)
		
		// Calcular posición en pantalla
		x := state.Position - float64(config.CarWidth/2)
		y := float64(track.GetLaneY(state.Lane) - config.CarHeight/2)
		
		// Opciones de dibujo
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		op.GeoM.Scale(1,1)
		
		// Si el auto terminó, podríamos aplicar algún efecto visual
		if state.Finished {
			// Opcional: hacer que parpadee o cambiar opacidad
			// op.ColorM.Scale(1, 1, 1, 0.7)
		}
		
		screen.DrawImage(sprite, op)
	}
}

// Layout define el tamaño de la pantalla
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

// startNewRace inicia una nueva carrera
func (g *Game) startNewRace() {
	g.race = core.NewRace(g.cfg)
	g.state = config.StateCountdown
	g.race.Start()
}