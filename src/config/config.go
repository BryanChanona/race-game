package config

const (
	//Configuración de pantalla
	ScreenWidth  = 800
	ScreenHeight = 600

	//Configuracion de la pista

	TrackLength = 700.0 //Longitud de la pista en pixeles
	LaneHeight  = 80    //Altura de cada carril
	NumLanes    = 5     //Numero de carriles
	StartLineX  = 50    // Posición X de la linea de salida
	FinishLineX = 750   // Posición X de la linea de meta

	//Configuración de los autos

	NumCars   = 5
	CarWidth  = 60
	CarHeight = 40

	//Configurar física
	MinSpeed       = 1.0 // Velocidad mínima en píxeles por tick
	MaxSpeed       = 4.0 // Velocidad máxima en píxeles por tick
	UpdateInterval = 50  // Milisegundos entre actualizaciones de posición

	//Estados del juego 
	StateMenu     = 0
	StateCountdown = 1
	StateRacing   = 2
	StateFinished = 3
)

type GameConfig struct {
	ScreenWidth    int
	ScreenHeight   int
	TrackLength    float64
	NumCars        int
	NumLanes       int
	MinSpeed       float64
	MaxSpeed       float64
	UpdateInterval int
}

// DefaultConfig devuelve la configuración por defecto
func DefaultConfig() *GameConfig {
	return &GameConfig{
		ScreenWidth:    ScreenWidth,
		ScreenHeight:   ScreenHeight,
		TrackLength:    TrackLength,
		NumCars:        NumCars,
		NumLanes:       NumLanes,
		MinSpeed:       MinSpeed,
		MaxSpeed:       MaxSpeed,
		UpdateInterval: UpdateInterval,
	}
}