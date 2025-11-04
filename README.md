🏎️ Race Game - Simulador de Carreras Concurrente en Go

Un simulador de carreras que demuestra patrones de concurrencia avanzados en Go, incluyendo Fan-Out/Fan-In, Barrier Pattern, y sincronización con canales.

 Tabla de Contenidos

Características
Arquitectura
Patrones de Concurrencia
Estructura del Proyecto
Instalación
Uso
Documentación Técnica
Conceptos Clave

Características

🚗 Simulación multi-auto: Cada auto corre en su propia goroutine independiente
🏁 Inicio sincronizado: Implementación del Barrier Pattern para inicio simultáneo
📊 Visualización en tiempo real: Renderizado de posiciones y estadísticas durante la carrera
🎮 Interfaz interactiva: Menús, countdown, HUD y pantalla de resultados
🔒 Thread-safe: Uso correcto de mutexes para acceso a estado compartido

 Arquitectura

El proyecto sigue una arquitectura en capas:
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

Patrones de Concurrencia

1. Fan-Out Pattern 
Múltiples goroutines se lanzan desde un punto central para ejecutar trabajo en paralelo.
Implementación en manager.go:

func (m *RaceManager) StartRace() {
    // Lanzar una goroutine por cada auto

	for _, car := range m.cars {
		m.wg.Add(1)
		go m.runCar(car)
	}
}

¿Qué hace?

Crea N goroutines (una por auto)
Cada goroutine ejecuta su propia lógica de carrera
Trabajan concurrentemente sin bloquearse entre sí

2. Fan-In Pattern 
Múltiples goroutines envían resultados a un punto central de recolección.

Implementación:
En car.go - Envío de resultados:
gofunc (c *Car) finish() {
    result := RaceResult{
        ID:        c.id,
        Lane:      c.lane,
        FinalTime: duration,
        FinalPos:  c.position,
    }
    
    c.resultChan <- result  // Envía al canal compartido
}

En manager.go - Recolección:
gofunc (m *RaceManager) collectResults() {
    for i := 0; i < expectedResults; i++ {
        result := <-m.resultChan  // Recibe de múltiples autos
        result.Position = position
        position++
        m.results = append(m.results, result)
    }
}
¿Qué hace?

Todos los autos envían sus resultados al mismo canal resultChan
Una goroutine recolectora recibe y procesa todos los resultados
Asigna posiciones (1°, 2°, 3°, etc.) según orden de llegada

3. Barrier Pattern (Sincronización de Inicio)
Sincroniza múltiples goroutines para que empiecen exactamente al mismo tiempo.

Implementación:
En car.go - Los autos esperan:
gofunc (c *Car) Run() {
    // Esperar la señal de inicio (BLOQUEANTE)
    <-c.startSignal
    
    // Todos arrancan aquí simultáneamente
    c.startTime = time.Now()
}

En manager.go - El manager libera:
gofunc (m *RaceManager) ReleaseStartSignal() {
    close(m.startSignal)  // Cierra el canal = todos se desbloquean
}

## 📁 Estructura del Proyecto
```
race-game/

│    
├── src/
├   ├───main.go                 # Punto de entrada
│   ├── config/
│   │   └── config.go           # Configuración y constantes
│   ├── core/                   # Lógica de concurrencia
│   │   ├── car.go              # Goroutine del auto, Fan-In
│   │   ├── manager.go          # Fan-Out, coordinación
│   │   ├── race.go             # Orquestación de la carrera
│   │   └── track.go            # Estado compartido, pista
│   ├── ui/                 # Capa de presentación
│   │   ├── assets.go           # Carga de recursos gráficos
│   │   ├── game.go             # Loop principal de Ebiten
│   │   └── ui.go               # Interfaz de usuario
│   └── pkg/
│       └── util/
│           └── logger.go       # Sistema de logging
└── README.md

 Instalación
Requisitos previos

Go 1.21+: Descargar Go
Dependencias de Ebiten: Guía de instalación

Pasos

Clonar el repositorio:

bashgit clone https://github.com/BryanChanona/race-game.git
cd race-game

Instalar dependencias:

go mod download

Ejecutar el juego:

go run src/main.go

 Uso
Controles

Tecla               Acción
ESPACIO     Iniciar carrera desde el menú
R           Reiniciar carrera después de terminar
ESC         Salir del juego

Flujo del Juego

Menú Principal → Presiona ESPACIO
Countdown (3, 2, 1, GO!) → Automático
Carrera en Progreso → Los autos corren solos
Resultados Finales → Muestra el podio y tiempos