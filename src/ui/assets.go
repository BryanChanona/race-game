package render

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"race-game/src/config"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

// Assets contiene todos los recursos gráficos del juego
type Assets struct {
	carSprites    []*ebiten.Image
	trackBg       *ebiten.Image
	startLine     *ebiten.Image
	finishLine    *ebiten.Image
}

// NewAssets crea y carga todos los assets
func NewAssets() *Assets {
	a := &Assets{
		carSprites: make([]*ebiten.Image, config.NumCars),
	}
	
	
	
	for i := 0; i < config.NumCars; i++ {
		a.carSprites[i] = MustLoadImage("assets/formula.png")
	}
	
	// Crear fondo de pista
	a.trackBg = createTrackBackground()
	
	// Crear líneas de salida y meta
	a.startLine = createLine(color.RGBA{255, 255, 255, 255}, 5)
	a.finishLine = createLine(color.RGBA{255, 255, 255, 255}, 5)
	
	return a
}


// MustLoadImage carga una imagen del sistema embed y la convierte a *ebiten.Image
func MustLoadImage(name string) *ebiten.Image {
	fmt.Println("Cargando imagen:", name)

	f, err := assets.Open(name)
	if err != nil {
		panic("No se pudo abrir " + name + ": " + err.Error())
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic("Error al decodificar " + name + ": " + err.Error())
	}

	return ebiten.NewImageFromImage(img)
}

// createTrackBackground crea el fondo de la pista
func createTrackBackground() *ebiten.Image {
	img := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	
	// Fondo gris oscuro
	img.Fill(color.RGBA{50, 50, 50, 255})
	
	// Dibujar carriles
	totalTrackHeight := config.NumLanes * config.LaneHeight
	offsetY := (config.ScreenHeight - totalTrackHeight) / 2
	
	// Dibujar líneas de carril (blanco)
	lineColor := color.RGBA{200, 200, 200, 255}
	
	for i := 0; i <= config.NumLanes; i++ {
		y := offsetY + i*config.LaneHeight
		
		// Línea horizontal
		for x := 0; x < config.ScreenWidth; x++ {
			img.Set(x, y, lineColor)
		}
	}
	
	// Dibujar línea central discontinua en cada carril
	dashColor := color.RGBA{255, 255, 255, 100}
	for lane := 0; lane < config.NumLanes; lane++ {
		y := offsetY + lane*config.LaneHeight + config.LaneHeight/2
		
		for x := 0; x < config.ScreenWidth; x += 40 {
			for dx := 0; dx < 20 && x+dx < config.ScreenWidth; dx++ {
				img.Set(x+dx, y, dashColor)
			}
		}
	}
	
	return img
}

// createLine crea una línea vertical
func createLine(col color.RGBA, width int) *ebiten.Image {
	height := config.NumLanes * config.LaneHeight
	img := ebiten.NewImage(width, height)
	img.Fill(col)
	return img
}

// GetCarSprite devuelve el sprite de un auto
func (a *Assets) GetCarSprite(carID int) *ebiten.Image {
	if carID >= 0 && carID < len(a.carSprites) {
		return a.carSprites[carID]
	}
	return a.carSprites[0]
}

// GetTrackBackground devuelve el fondo de la pista
func (a *Assets) GetTrackBackground() *ebiten.Image {
	return a.trackBg
}

// GetStartLine devuelve la línea de salida
func (a *Assets) GetStartLine() *ebiten.Image {
	return a.startLine
}

// GetFinishLine devuelve la línea de meta
func (a *Assets) GetFinishLine() *ebiten.Image {
	return a.finishLine
}