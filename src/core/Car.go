package core

import (
	"math/rand"
	"race-game/src/config"
	"race-game/src/pkg/util"

	"sync"
	"time"
)

// CarState representa el estado de un auto
type CarState struct {
	ID       int
	Position float64
	Lane     int
	Speed    float64
	Finished bool
}

// RaceResult representa el resultado de un auto al terminar
type RaceResult struct {
	ID          int
	Lane        int
	Position    int // 1 = primero, 2 = segundo, etc.
	FinalTime   time.Duration
	FinalPos    float64
}

// Car representa un automóvil en la carrera
type Car struct {
	id           int
	lane         int
	position     float64
	speed        float64
	finished     bool
	startTime    time.Time
	finishTime   time.Time
	
	// Canales para comunicación
	startSignal  <-chan struct{}
	resultChan   chan<- RaceResult
	
	// Sincronización
	mu           sync.RWMutex
	
	// Acceso al estado compartido
	sharedState  *SharedRaceState
	
	logger       *util.Logger
}

// NewCar crea un nuevo auto
func NewCar(id, lane int, startSignal <-chan struct{}, resultChan chan<- RaceResult, sharedState *SharedRaceState) *Car {
	return &Car{
		id:          id,
		lane:        lane,
		position:    config.StartLineX,
		speed:       0,
		finished:    false,
		startSignal: startSignal,
		resultChan:  resultChan,
		sharedState: sharedState,
		logger:      util.GetLogger(),
	}
}

// Run es la goroutine principal del auto (FAN-OUT)
// Cada auto corre en su propia goroutine independiente
func (c *Car) Run() {
	defer func() {
		c.logger.Debug("Auto %d: Goroutine terminada", c.id)
	}()
	
	c.logger.Info("Auto %d: Esperando señal de inicio...", c.id)
	
	// Esperar la señal de inicio (sincronización con Barrier Pattern)
	<-c.startSignal
	
	c.startTime = time.Now()
	c.logger.Info("Auto %d: ¡Arrancando desde carril %d!", c.id, c.lane)
	
	// Asignar velocidad aleatoria (simula diferentes potencias de motor)
	c.mu.Lock()
	c.speed = config.MinSpeed + rand.Float64()*(config.MaxSpeed-config.MinSpeed)
	c.mu.Unlock()
	
	c.logger.Debug("Auto %d: Velocidad asignada = %.2f px/tick", c.id, c.speed)
	
	// Ticker para actualizar la posición periódicamente
	ticker := time.NewTicker(time.Duration(config.UpdateInterval) * time.Millisecond)
	defer ticker.Stop()
	
	// Loop principal del auto
	for {
		select {
		case <-ticker.C:
			// Actualizar posición
			if c.updatePosition() {
				// El auto ha terminado la carrera
				c.finish()
				return
			}
		}
	}
}

// updatePosition actualiza la posición del auto
// Retorna true si el auto cruzó la meta
func (c *Car) updatePosition() bool {
	c.mu.Lock()
	
	// Pequeña variación aleatoria en la velocidad (simula imperfecciones)
	variation := (rand.Float64() - 0.5) * 0.2
	c.position += c.speed + variation
	
	hasFinished := c.position >= config.FinishLineX
	
	// Actualizar el estado compartido (SECCIÓN CRÍTICA)
	c.sharedState.mu.Lock()
	c.sharedState.positions[c.id] = CarState{
		ID:       c.id,
		Position: c.position,
		Lane:     c.lane,
		Speed:    c.speed,
		Finished: hasFinished,
	}
	c.sharedState.mu.Unlock()
	
	c.mu.Unlock()
	
	return hasFinished
}

// finish marca el auto como terminado y envía el resultado (FAN-IN)
func (c *Car) finish() {
	c.mu.Lock()
	c.finished = true
	c.finishTime = time.Now()
	duration := c.finishTime.Sub(c.startTime)
	c.mu.Unlock()
	
	c.logger.Info("Auto %d: ¡Cruzó la meta! Tiempo: %v", c.id, duration)
	
	// Enviar resultado al canal (FAN-IN: todos convergen en un solo canal)
	result := RaceResult{
		ID:        c.id,
		Lane:      c.lane,
		FinalTime: duration,
		FinalPos:  c.position,
	}
	
	c.resultChan <- result
}

// GetState devuelve el estado actual del auto de forma segura
func (c *Car) GetState() CarState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return CarState{
		ID:       c.id,
		Position: c.position,
		Lane:     c.lane,
		Speed:    c.speed,
		Finished: c.finished,
	}
}

