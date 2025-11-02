package util

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// LogLevel representa el nivel de log
type LogLevel int

const (
	//iota inicia en 0, asigna automáticamente valores crecientes
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// Logger seguro para subprocesos
type Logger struct {
	mu    sync.Mutex
	level LogLevel
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger devuelve la instancia singleton del logger
func GetLogger() *Logger {
	//Nos aseguramos de que el bloque que se esté ejecutando dentro de once.Do() solo se ejecute una vez
	once.Do(func() {
		instance = &Logger{
			level: INFO,
		}
	})
	return instance
}

// SetLevel establece el nivel mínimo de log
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Debug registra un mensaje de depuración
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info registra un mensaje informativo
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warning registra una advertencia
func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

// Error registra un error
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// log es el método interno para registrar mensajes
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	levelStr := ""
	switch level {
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case WARNING:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	}

	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	log.Printf("[%s] [%s] %s", timestamp, levelStr, message)
}
