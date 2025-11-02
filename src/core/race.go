package core

import (
	"race-game/src/config"
	"race-game/src/pkg/util"
	"time"
)

// RaceState representa el estado de la carrera
type RaceState int

const (
	RaceStateReady RaceState = iota
	RaceStateCountdown
	RaceStateRunning
	RaceStateFinished
)

// Race orquesta una carrera completa
type Race struct {
	manager        *RaceManager
	state          RaceState
	countdownValue int
	startTime      time.Time
	endTime        time.Time
	logger         *util.Logger
}

// NewRace crea una nueva carrera
func NewRace(cfg *config.GameConfig) *Race {
	return &Race{
		manager:        NewRaceManager(cfg),
		state:          RaceStateReady,
		countdownValue: 3,
		logger:         util.GetLogger(),
	}
}

// Start inicia la carrera
func (r *Race) Start() {
	if r.state != RaceStateReady {
		r.logger.Warning("Race: Intento de iniciar carrera en estado incorrecto: %v", r.state)
		return
	}
	
	r.logger.Info("Race: Iniciando secuencia de carrera")
	r.state = RaceStateCountdown
	
	// Iniciar las goroutines de los autos (pero esperarán la señal)
	r.manager.StartRace()
	
	// Goroutine para manejar el countdown
	go r.handleCountdown()
}

// handleCountdown maneja la cuenta regresiva
func (r *Race) handleCountdown() {
	r.logger.Info("Race: Countdown iniciado: %d", r.countdownValue)
	
	for r.countdownValue > 0 {
		time.Sleep(1 * time.Second)
		r.countdownValue--
		r.logger.Info("Race: Countdown: %d", r.countdownValue)
	}
	
	// ¡Luz verde!
	r.logger.Info("Race: ¡GO! Liberando señal de inicio")
	r.startTime = time.Now()
	r.state = RaceStateRunning
	
	// Liberar la señal de inicio (todos los autos arrancan)
	r.manager.ReleaseStartSignal()
}

// Update actualiza el estado de la carrera
func (r *Race) Update() {
	if r.state == RaceStateRunning {
		// Verificar si la carrera ha terminado
		if r.manager.IsRaceFinished() {
			r.endTime = time.Now()
			r.state = RaceStateFinished
			r.logger.Info("Race: Carrera finalizada. Duración total: %v", r.GetRaceDuration())
		}
	}
}

// GetState devuelve el estado actual de la carrera
func (r *Race) GetState() RaceState {
	return r.state
}

// GetCountdownValue devuelve el valor actual del countdown
func (r *Race) GetCountdownValue() int {
	return r.countdownValue
}

// GetManager devuelve el manager de la carrera
func (r *Race) GetManager() *RaceManager {
	return r.manager
}

// GetRaceDuration devuelve la duración de la carrera
func (r *Race) GetRaceDuration() time.Duration {
	if r.state == RaceStateFinished {
		return r.endTime.Sub(r.startTime)
	}
	if r.state == RaceStateRunning {
		return time.Since(r.startTime)
	}
	return 0
}

// IsReady verifica si la carrera está lista para comenzar
func (r *Race) IsReady() bool {
	return r.state == RaceStateReady
}

// IsRunning verifica si la carrera está en curso
func (r *Race) IsRunning() bool {
	return r.state == RaceStateRunning
}

// IsFinished verifica si la carrera ha terminado
func (r *Race) IsFinished() bool {
	return r.state == RaceStateFinished
}

// IsCountdown verifica si está en countdown
func (r *Race) IsCountdown() bool {
	return r.state == RaceStateCountdown
}