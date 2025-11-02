package main

import (
	"log"

	"race-game/src/config"
	"race-game/src/ui"
	"race-game/src/pkg/util"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Inicializar logger
	logger := util.GetLogger()
	logger.SetLevel(util.INFO) // Cambiar a DEBUG para más información
	
	logger.Info("=== RACE GAME - Concurrency Patterns Demo ===")
	logger.Info("Iniciando juego...")
	
	// Cargar configuración
	cfg := config.DefaultConfig()
	logger.Info("Configuración cargada: %d autos, pista de %.0f px", cfg.NumCars, cfg.TrackLength)
	
	// Crear instancia del juego
	game := render.NewGame(cfg)
	
	// Configurar ventana
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle("Race Game - Fan-Out/Fan-In Pattern")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	
	logger.Info("Ventana creada: %dx%d", cfg.ScreenWidth, cfg.ScreenHeight)
	logger.Info("Presiona ESPACIO en el menú para comenzar")
	logger.Info("===============================================")
	
	// Iniciar el game loop de Ebiten
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	
	logger.Info("Juego finalizado")
}