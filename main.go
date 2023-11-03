package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync" // Importar la biblioteca sync para la exclusión mutua.
)

const (
	Nro_casillas = 72
)

type Jugador struct {
	Nombre   string
	Fichas   int
	Posicion int
}

type Tablero struct {
	jugadores []Jugador
	tablero   []rune
}

func (t Tablero) MostrarJugadores(j []Jugador) {
	for i := 0; i < 4; i++ {
		nuevo := make([]rune, len(t.tablero))
		copy(nuevo, t.tablero)
		if j[i].Posicion <= 0 {
			nuevo[0] = 'X'
		}
		if j[i].Posicion > len(t.tablero)-1 {
			nuevo[len(t.tablero)-1] = 'X'
		} else {
			nuevo[j[i].Posicion] = 'X'
		}
		fmt.Printf("%c\n", nuevo)
	}
}

func Dado() int {
	d1 := rand.Intn(6) + 1                        // Dado 1
	d2 := rand.Intn(6) + 1                        // Dado 2
	s := int(math.Pow(-1, float64(rand.Intn(2)))) // Dado Signo

	fmt.Printf("DADOS: (%d) ", d1)
	if s > 0 {
		fmt.Printf("(+)")
	} else {
		fmt.Printf("(-)")
	}
	fmt.Printf(" (%d)", d2)

	return d1 + s*d2
}

func GenTablero(casillas int) []rune {
	tablero := make([]rune, casillas)
	for i := range tablero {
		tablero[i] = '_' // . -> Casillas en blanco
	}
	tablero[0] = '#'              // # -> Inicio
	tablero[len(tablero)-1] = '#' // $ -> Fin

	umbral := int(float64(casillas) * 0.4) // Solo se llenará el 40% de las casillas en blanco con casillas especiales

	for i := 0; i < umbral; i++ {
		idx := rand.Intn(int(casillas))
		// No afectar las casillas inicial y final
		if idx == 0 || idx == len(tablero)-1 {
			continue
		}

		/*
			1 -> +3 espacios
			2 -> -3 espacios
			3 -> REGRESA AL PRINCIPIO
		*/
		switch rand.Intn(3) + 1 {
		case 1:
			tablero[idx] = '1'
		case 2:
			tablero[idx] = '2'
		case 3:
			tablero[idx] = '3'
		}
	}
	return tablero
}

func Meta(j *Jugador) bool {
	return j.Posicion >= Nro_casillas-1
}

func obtenerEstadoJuego(jugadores []Jugador) string {
	estado := ""
	for _, j := range jugadores {
		estado += fmt.Sprintf("%s: Posición=%d, Fichas=%d\n", j.Nombre, j.Posicion, j.Fichas)
	}
	return estado
}

func main() {
	jugadores := []Jugador{
		{"J1", 0, 0},
		{"J2", 0, 0},
		{"J3", 0, 0},
		{"J4", 0, 0},
	}
	t := Tablero{jugadores, GenTablero(Nro_casillas)}
	var mu sync.Mutex // Agregar una exclusión mutua para garantizar la sincronización.
	jugadorActual := 0

	// Define un canal para comunicar el estado del juego
	estadoJuegoChan := make(chan string)

	// Goroutine para mostrar el estado del juego
	go func() {
		for estado := range estadoJuegoChan {
			fmt.Printf("\nEstado del juego:\n%s\n", estado)
		}
	}()

	for {
		if jugadorActual == 0 {
			fmt.Printf("----------------------------------------------\n")
		}

		j := &jugadores[jugadorActual]

		dado := Dado()
		j.Posicion = j.Posicion + dado

		if j.Posicion < 0 {
			j.Posicion = 0
		}
		if j.Posicion > Nro_casillas-1 {
			j.Posicion = Nro_casillas - 1
		}

		if t.tablero[j.Posicion] == '1' {
			fmt.Printf("\t\t¡AVANZASTE 3 CASILLAS!")
			j.Posicion = j.Posicion + 3
			if j.Posicion > Nro_casillas-1 {
				j.Posicion = Nro_casillas - 1
			}
		}
		if t.tablero[j.Posicion] == '2' {
			fmt.Printf("\t\t¡RETROCEDISTE 3 CASILLAS!")
			j.Posicion = j.Posicion - 3
			if j.Posicion < 0 {
				j.Posicion = 0
			}
		}
		if t.tablero[j.Posicion] == '3' {
			fmt.Printf("\t\t¡VUELVES AL INICIO!")
			j.Posicion = 0
		}

		fmt.Printf("\t\t%s OBTUVO: %d, SUM: %d", j.Nombre, dado, j.Posicion)

		if Meta(j) {
			mu.Lock() // Bloquear el acceso concurrente a las fichas.
			j.Fichas++
			mu.Unlock() // Desbloquear el acceso después de actualizar las fichas.
			if j.Fichas == 4 {
				fmt.Printf("\t\t¡%s GANO!, cantidad de fichas metidas: %d\n", j.Nombre, j.Fichas)
				close(estadoJuegoChan) // Cierra el canal cuando el juego ha terminado
				break
			}

			fmt.Printf("\t\t¡%s metio una ficha!, cantidad de fichas metidas: %d", j.Nombre, j.Fichas)
			j.Posicion = 0
		}

		// Después de calcular el estado del juego, envía el estado al canal
		estadoJuego := obtenerEstadoJuego(jugadores)
		estadoJuegoChan <- estadoJuego

		jugadorActual = (jugadorActual + 1) % len(jugadores)

		fmt.Printf("\n")
	}

	t.MostrarJugadores(jugadores)
}
