package render

import (
	"fmt"
	"race-game/src/config"
	"race-game/src/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// UI maneja la interfaz de usuario
type UI struct {
	assets *Assets
}

// NewUI crea una nueva interfaz
func NewUI(assets *Assets) *UI {
	return &UI{
		assets: assets,
	}
}

// DrawMenu dibuja el menú principal
func (u *UI) DrawMenu(screen *ebiten.Image) {
	// Título
	titleText := "RACE GAME"
	ebitenutil.DebugPrintAt(screen, titleText, config.ScreenWidth/2-60, 100)
	
	// Instrucciones
	instructionText := "Presiona ESPACIO para comenzar"
	ebitenutil.DebugPrintAt(screen, instructionText, config.ScreenWidth/2-140, config.ScreenHeight/2)
	
	// Créditos
	creditsText := "Concurrency Patterns Demo"
	ebitenutil.DebugPrintAt(screen, creditsText, config.ScreenWidth/2-100, config.ScreenHeight-50)
}

// DrawCountdown dibuja el countdown
func (u *UI) DrawCountdown(screen *ebiten.Image, value int) {
	text := fmt.Sprintf("%d", value)
	if value == 0 {
		text = "GO!"
	}
	
	// Dibujar texto grande en el centro
	x := config.ScreenWidth/2 - 20
	y := config.ScreenHeight/2 - 50
	
	// Crear un efecto de texto grande (dibujarlo múltiples veces con offset)
	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			ebitenutil.DebugPrintAt(screen, text, x+dx, y+dy)
		}
	}
}

// DrawHUD dibuja el HUD durante la carrera
func (u *UI) DrawHUD(screen *ebiten.Image, race *core.Race) {
	manager := race.GetManager()
	
	// Tiempo transcurrido
	duration := race.GetRaceDuration()
	timeText := fmt.Sprintf("Tiempo: %.2fs", duration.Seconds())
	ebitenutil.DebugPrintAt(screen, timeText, 10, 10)
	
	// Autos que han terminado
	finishedCount := manager.GetFinishedCount()
	totalCars := len(manager.GetCars())
	progressText := fmt.Sprintf("Terminados: %d/%d", finishedCount, totalCars)
	ebitenutil.DebugPrintAt(screen, progressText, 10, 30)
	
	// Posiciones en tiempo real
	positions := manager.GetSharedState().GetAllPositions()
	y := 60
	ebitenutil.DebugPrintAt(screen, "Posiciones:", 10, y)
	y += 20
	
	for id := 0; id < totalCars; id++ {
		if state, exists := positions[id]; exists {
			status := ""
			if state.Finished {
				status = " [TERMINÓ]"
			}
			posText := fmt.Sprintf("Auto %d: %.0f px%s", id, state.Position, status)
			
			// Color según el estado
			if state.Finished {
				// No hay una forma fácil de cambiar el color con DebugPrintAt
				// pero podemos indicarlo con el texto
			}
			
			ebitenutil.DebugPrintAt(screen, posText, 10, y)
			y += 20
		}
	}
}

// DrawResults dibuja la pantalla de resultados
func (u *UI) DrawResults(screen *ebiten.Image, race *core.Race) {
	// Título
	titleText := "¡CARRERA TERMINADA!"
	ebitenutil.DebugPrintAt(screen, titleText, config.ScreenWidth/2-100, 80)
	
	// Resultados
	results := race.GetManager().GetResults()
	y := 150
	
	ebitenutil.DebugPrintAt(screen, "RESULTADOS:", config.ScreenWidth/2-60, y)
	y += 40
	
	for _, result := range results {
		positionText := ""
		switch result.Position {
		case 1:
			positionText = "🥇"
		case 2:
			positionText = "🥈"
		case 3:
			positionText = "🥉"
		default:
			positionText = fmt.Sprintf("%d°", result.Position)
		}
		
		text := fmt.Sprintf("%s  Auto %d - Tiempo: %.2fs", 
			positionText, result.ID, result.FinalTime.Seconds())
		
		ebitenutil.DebugPrintAt(screen, text, config.ScreenWidth/2-150, y)
		y += 30
	}
	
	// Instrucción para reiniciar
	y += 40
	restartText := "Presiona R para reiniciar o ESC para salir"
	ebitenutil.DebugPrintAt(screen, restartText, config.ScreenWidth/2-180, y)
}

// CheckMenuInput verifica el input en el menú
func (u *UI) CheckMenuInput() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// CheckRestartInput verifica el input para reiniciar
func (u *UI) CheckRestartInput() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyR)
}

// CheckExitInput verifica el input para salir
func (u *UI) CheckExitInput() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}