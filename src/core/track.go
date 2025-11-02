package core

import (
	"race-game/src/config"
	"sync"
)

// SharedRaceState contiene el estado compartido de la carrera
// Este es el recurso compartido que necesita sincronización
type SharedRaceState struct {
	mu        sync.RWMutex
	positions map[int]CarState // map[carID]CarState
}

// NewSharedRaceState crea un nuevo estado compartido
func NewSharedRaceState() *SharedRaceState {
	return &SharedRaceState{
		positions: make(map[int]CarState),
	}
}

// GetAllPositions devuelve una copia de todas las posiciones de forma segura
// Usado por el renderer para dibujar los autos
func (s *SharedRaceState) GetAllPositions() map[int]CarState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Crear una copia para evitar race conditions
	copy := make(map[int]CarState, len(s.positions))
	for k, v := range s.positions {
		copy[k] = v
	}
	
	return copy
}

// GetPosition devuelve la posición de un auto específico
func (s *SharedRaceState) GetPosition(carID int) (CarState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	state, exists := s.positions[carID]
	return state, exists
}

// Track representa la pista de carreras
type Track struct {
	length      float64
	numLanes    int
	startLineX  float64
	finishLineX float64
	laneHeight  int
}

// NewTrack crea una nueva pista
func NewTrack() *Track {
	return &Track{
		length:      config.TrackLength,
		numLanes:    config.NumLanes,
		startLineX:  config.StartLineX,
		finishLineX: config.FinishLineX,
		laneHeight:  config.LaneHeight,
	}
}

// GetLaneY devuelve la posición Y del centro de un carril
func (t *Track) GetLaneY(lane int) int {
	// Calcular el offset para centrar la pista verticalmente
	totalTrackHeight := t.numLanes * t.laneHeight
	offsetY := (config.ScreenHeight - totalTrackHeight) / 2
	
	return offsetY + lane*t.laneHeight + t.laneHeight/2
}

// GetLength devuelve la longitud de la pista
func (t *Track) GetLength() float64 {
	return t.length
}

// GetNumLanes devuelve el número de carriles
func (t *Track) GetNumLanes() int {
	return t.numLanes
}

// GetStartLineX devuelve la posición X de la línea de salida
func (t *Track) GetStartLineX() float64 {
	return t.startLineX
}

// GetFinishLineX devuelve la posición X de la línea de meta
func (t *Track) GetFinishLineX() float64 {
	return t.finishLineX
}

// IsFinished verifica si un auto cruzó la meta
func (t *Track) IsFinished(position float64) bool {
	return position >= t.finishLineX
}