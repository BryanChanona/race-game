# 🏎️ Race Game - Simulador de Carreras Concurrente en Go

Un simulador de carreras que demuestra patrones de concurrencia avanzados en Go, incluyendo **Fan-Out/Fan-In**, **Barrier Pattern**, y sincronización con canales.

---

##  Tabla de Contenidos

- [Características](#-características)
- [Arquitectura](#-arquitectura)
- [Patrones de Concurrencia](#-patrones-de-concurrencia)
  - [Fan-Out Pattern](#1-fan-out-pattern-dispersión)
  - [Fan-In Pattern](#2-fan-in-pattern-convergencia)
  - [Barrier Pattern](#3-barrier-pattern-sincronización-de-inicio)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Instalación](#-instalación)
- [Uso](#-uso)
- [Tecnologías](#-tecnologías)

---

##  Características

-  **Simulación multi-auto**: Cada auto corre en su propia goroutine independiente
-  **Inicio sincronizado**: Implementación del Barrier Pattern para inicio simultáneo
-  **Visualización en tiempo real**: Renderizado de posiciones y estadísticas durante la carrera
-  **Interfaz interactiva**: Menús, countdown, HUD y pantalla de resultados
-  **Thread-safe**: Uso correcto de mutexes para acceso a estado compartido

---

##  Arquitectura

El proyecto sigue una arquitectura en capas:

```
┌─────────────────────────────────────────┐
│            UI Layer (ui)                │
│           screen, ui, assets            │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│          Core Layer (core)              │
│  Race, Manager, Car, Track              │
│  Lógica de concurrencia y sincronización│
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│       Configuration (config)            │
│  Constantes y configuración del juego   │
└─────────────────────────────────────────┘
```

### Capas del Sistema

| Capa | Responsabilidad | Componentes |
|------|----------------|-------------|
| **UI Layer** | Presentación y renderizado | `screen.go`, `ui.go`, `assets.go` |
| **Core Layer** | Lógica de negocio y concurrencia | `car.go`, `manager.go`, `race.go`, `track.go` |
| **Config Layer** | Configuración del sistema | `config.go` |

---

##  Patrones de Concurrencia

### 1. Fan-Out Pattern 

**Concepto**: Múltiples goroutines se lanzan desde un punto central para ejecutar trabajo en paralelo.

**Implementación en `manager.go`**:

```go
func (m *RaceManager) StartRace() {
    // Lanzar una goroutine por cada auto
    for _, car := range m.cars {
        m.wg.Add(1)
        go m.runCar(car)
    }
}
```

**¿Qué hace?**
- Crea N goroutines (una por auto)
- Cada goroutine ejecuta su propia lógica de carrera
- Trabajan concurrentemente sin bloquearse entre sí

---

### 2. Fan-In Pattern 

**Concepto**: Múltiples goroutines envían resultados a un punto central de recolección.

**Implementación**:

**En `car.go` - Envío de resultados**:
```go
func (c *Car) finish() {
    result := RaceResult{
        ID:        c.id,
        Lane:      c.lane,
        FinalTime: duration,
        FinalPos:  c.position,
    }
    
    c.resultChan <- result  // Envía al canal compartido
}
```

**En `manager.go` - Recolección**:
```go
func (m *RaceManager) collectResults() {
    for i := 0; i < expectedResults; i++ {
        result := <-m.resultChan  // Recibe de múltiples autos
        result.Position = position
        position++
        m.results = append(m.results, result)
    }
}
```

**¿Qué hace?**
- Todos los autos envían sus resultados al mismo canal `resultChan`
- Una goroutine recolectora recibe y procesa todos los resultados
- Asigna posiciones (1°, 2°, 3°, etc.) según orden de llegada

---

### 3. Barrier Pattern (Sincronización de Inicio)

**Concepto**: Sincroniza múltiples goroutines para que empiecen exactamente al mismo tiempo.

**Implementación**:

**En `car.go` - Los autos esperan**:
```go
func (c *Car) Run() {
    // Esperar la señal de inicio (BLOQUEANTE)
    <-c.startSignal
    
    // Todos arrancan aquí simultáneamente
    c.startTime = time.Now()
}
```

**En `manager.go` - El manager libera**:
```go
func (m *RaceManager) ReleaseStartSignal() {
    close(m.startSignal)  // Cierra el canal = todos se desbloquean
}
```

**¿Qué hace?**
- Todas las goroutines esperan bloqueadas en un canal cerrado
- Al cerrar el canal, todas se liberan simultáneamente
- Garantiza un inicio justo de la carrera

---

##  Estructura del Proyecto

```
race-game/
│
├── src/
│   ├── main.go                 # Punto de entrada de la aplicación
│   │
│   ├── config/
│   │   └── config.go           # Configuración y constantes globales
│   │
│   ├── core/                   # Lógica de concurrencia
│   │   ├── car.go              # Goroutine del auto, Fan-In
│   │   ├── manager.go          # Fan-Out, coordinación de carreras
│   │   ├── race.go             # Orquestación de la carrera
│   │   └── track.go            # Estado compartido, gestión de pista
│   │
│   ├── ui/                     # Capa de presentación
│   │   ├── assets.go           # Carga de recursos gráficos
│   │   ├── game.go             # Loop principal de Ebiten
│   │   └── ui.go               # Renderizado de interfaz de usuario
│   │
│   └── pkg/
│       └── util/
│           └── logger.go       # Sistema de logging
│
├── go.mod                      # Dependencias del proyecto
├── go.sum                      # Checksums de dependencias
└── README.md                   # Este archivo
```

---

##  Instalación

### Requisitos Previos

- **Go 1.21+**: [Descargar Go](https://golang.org/dl/)
- **Dependencias de Ebiten**: [Guía de instalación](https://ebitengine.org/en/documents/install.html)

### Pasos de Instalación

1. **Clonar el repositorio**:
```bash
git clone https://github.com/BryanChanona/race-game.git
cd race-game
```

2. **Instalar dependencias**:
```bash
go mod download
```

3. **Ejecutar el juego**:
```bash
go run src/main.go
```

---

##  Uso

### Controles

| Tecla | Acción |
|-------|--------|
| `ESPACIO` | Iniciar carrera desde el menú |
| `R` | Reiniciar carrera después de terminar |
| `ESC` | Salir del juego |

### Flujo del Juego

1. **Menú Principal** → Presiona `ESPACIO`
2. **Countdown** (3, 2, 1, GO!) → Automático
3. **Carrera en Progreso** → Los autos corren solos
4. **Resultados Finales** → Muestra el podio y tiempos

---

##  Tecnologías

- **Lenguaje**: Go 1.21+
- **Motor Gráfico**: [Ebiten](https://ebitengine.org/)
- **Patrones**: Concurrencia con goroutines, canales y mutexes



##  Autor

**Briyan Chanona**
- GitHub: [@BryanChanona](https://github.com/BryanChanona)

---

