package core

import (
	"race-game/src/config"
	"race-game/src/pkg/util"
	"sort"
	"sync"
)

// RaceManager coordina todas las goroutines (autos)
// Implementa los patrones Fan-Out (crear goroutines) y Fan-In (recolectar resultados)
type RaceManager struct {
	cars         []*Car
	track        *Track
	sharedState  *SharedRaceState
	
	// Canales para sincronización
	startSignal  chan struct{}  // Canal para el Barrier Pattern (inicio simultáneo)
	resultChan   chan RaceResult // Canal para Fan-In (recolección de resultados)
	
	// Control de estado
	results      []RaceResult
	resultsMu    sync.Mutex
	finishedCount int
	
	logger       *util.Logger
	wg           sync.WaitGroup // WaitGroup para esperar a todas las goroutines
}

// NewRaceManager crea un nuevo manager de carrera
func NewRaceManager(cfg *config.GameConfig) *RaceManager {
	track := NewTrack()
	sharedState := NewSharedRaceState()
	
	manager := &RaceManager{
		cars:        make([]*Car, 0, cfg.NumCars),
		track:       track,
		sharedState: sharedState,
		startSignal: make(chan struct{}),          // Canal sin buffer para broadcast
		resultChan:  make(chan RaceResult, cfg.NumCars), // Buffer = número de autos
		results:     make([]RaceResult, 0, cfg.NumCars),
		logger:      util.GetLogger(),
	}
	
	// Crear los autos (aún no los arrancamos)
	for i := 0; i < cfg.NumCars; i++ {
		car := NewCar(i, i, manager.startSignal, manager.resultChan, sharedState)
		manager.cars = append(manager.cars, car)
		manager.logger.Debug("Manager: Auto %d creado en carril %d", i, i)
	}
	
	return manager
}

// StartRace inicia la carrera (FAN-OUT: dispara todas las goroutines)
func (m *RaceManager) StartRace() {
	m.logger.Info("Manager: Iniciando carrera con %d autos", len(m.cars))
	
	// FAN-OUT: Lanzar una goroutine por cada auto
	for _, car := range m.cars {
		m.wg.Add(1)
		go m.runCar(car)
	}
	
	m.logger.Info("Manager: Todas las goroutines de autos lanzadas")
	
	// Goroutine separada para recolectar resultados (FAN-IN)
	go m.collectResults()
	
	m.logger.Debug("Manager: Goroutine de recolección iniciada")
}

// runCar ejecuta la lógica de un auto dentro de una goroutine
func (m *RaceManager) runCar(c *Car) {
	defer m.wg.Done()
	c.Run()
}

// ReleaseStartSignal libera a todos los autos para que empiecen a correr
// Implementa el Barrier Pattern: todos esperan hasta que se cierra el canal
func (m *RaceManager) ReleaseStartSignal() {
	m.logger.Info("Manager: ¡Luz verde! Cerrando canal de inicio")
	close(m.startSignal) // Cerrar el canal permite que TODOS lean simultáneamente
}

// collectResults recolecta los resultados de todos los autos (FAN-IN)
// Esta es la goroutine que implementa el patrón Fan-In
func (m *RaceManager) collectResults() {
	m.logger.Debug("Manager: Recolector de resultados iniciado")
	
	expectedResults := len(m.cars)
	position := 1
	
	// Esperar a que todos los autos terminen
	for i := 0; i < expectedResults; i++ {
		result := <-m.resultChan // Recibir resultado del canal (BLOQUEANTE)
		
		m.logger.Info("Manager: Recibido resultado del Auto %d (posición %d)", result.ID, position)
		
		// Asignar posición de llegada
		result.Position = position
		position++
		
		// Almacenar resultado de forma thread-safe
		m.resultsMu.Lock()
		m.results = append(m.results, result)
		m.finishedCount++
		m.resultsMu.Unlock()
		
		m.logger.Debug("Manager: %d/%d autos han terminado", m.finishedCount, expectedResults)
	}
	
	m.logger.Info("Manager: Todos los autos han terminado la carrera")
	
	// Ordenar resultados por posición
	m.resultsMu.Lock()
	sort.Slice(m.results, func(i, j int) bool {
		return m.results[i].Position < m.results[j].Position
	})
	m.resultsMu.Unlock()
}

// GetSharedState devuelve el estado compartido (para el renderer)
func (m *RaceManager) GetSharedState() *SharedRaceState {
	return m.sharedState
}

// GetTrack devuelve la pista
func (m *RaceManager) GetTrack() *Track {
	return m.track
}

// GetResults devuelve los resultados de forma thread-safe
func (m *RaceManager) GetResults() []RaceResult {
	m.resultsMu.Lock()
	defer m.resultsMu.Unlock()
	
	// Retornar una copia para evitar modificaciones externas
	copy := make([]RaceResult, len(m.results))
	for i, r := range m.results {
		copy[i] = r
	}
	
	return copy
}

// GetFinishedCount devuelve cuántos autos han terminado
func (m *RaceManager) GetFinishedCount() int {
	m.resultsMu.Lock()
	defer m.resultsMu.Unlock()
	return m.finishedCount
}

// IsRaceFinished verifica si la carrera ha terminado
func (m *RaceManager) IsRaceFinished() bool {
	m.resultsMu.Lock()
	defer m.resultsMu.Unlock()
	return m.finishedCount >= len(m.cars)
}

// Wait espera a que todas las goroutines terminen
func (m *RaceManager) Wait() {
	m.wg.Wait()
	m.logger.Info("Manager: Todas las goroutines han finalizado")
}

// GetCars devuelve la lista de autos
func (m *RaceManager) GetCars() []*Car {
	return m.cars
}