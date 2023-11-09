package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	nCasillas         = 72
	nFichasPorJugador = 4
)

func Clamp(val, min, max int) int {
	return int(math.Max(math.Min(float64(val), float64(min)), float64(max)))
}

type Jugador struct {
	Nombre            string
	FichasPos         []int
	FichasCompletadas int
}

type Juego struct {
	Tablero   string
	Jugadores []Jugador
}

func (juego *Juego) GetChar(i int) rune {
	return []rune(juego.Tablero)[i]
}

func (juego *Juego) Restringir(pId, ficha int) {
	juego.Jugadores[pId].FichasPos[ficha] = Clamp(juego.Jugadores[pId].FichasPos[ficha], 0, nCasillas-1)
}

func (juego *Juego) Print() {
	// Imprimimos el tablero + la posición de la ficha actual de cada jugador + número de fichas completadas
	fmt.Printf("\nESTADO DEL JUEGO\n====\n%s\n", juego.Tablero)
	for i := 0; i < len(juego.Jugadores); i++ {
		for j := 0; j < juego.Jugadores[i].FichasPos[juego.Jugadores[i].FichasCompletadas]-1; j++ {
			fmt.Print(" ")
		}
		fmt.Print("¡")
		for j := juego.Jugadores[i].FichasPos[juego.Jugadores[i].FichasCompletadas] + 1; j < nCasillas; j++ {
			fmt.Print(" ")
		}
		fmt.Printf("%d\n====\n", juego.Jugadores[i].FichasCompletadas)
	}
}

func (juego *Juego) CheckMeta() (bool, int) {
	for i := 0; i < len(juego.Jugadores); i++ {
		for j := 0; j < nFichasPorJugador; j++ {
			if juego.Jugadores[i].FichasPos[j] == nCasillas-1 {
				juego.Jugadores[i].FichasCompletadas++
			}
		}
		if juego.Jugadores[i].FichasCompletadas == nFichasPorJugador {
			return true, i
		}
	}
	return false, -1
}

var juego Juego
var win bool   // Determinar si algún jugador ganó
var pWinId int // Id del jugador que ganó. Mientras no haya uno, este será -1

var direccionRemota string

func Init() {
	tmpTablero := make([]string, nCasillas)
	for casilla := range tmpTablero {
		tmpTablero[casilla] = "."
	}
	tmpTablero[0] = "#"
	tmpTablero[len(tmpTablero)-1] = "@"

	umbral := nCasillas * 4 / 10
	for i := 0; i < umbral; i++ {
		idx := rand.Intn(nCasillas-2) + 1

		switch rand.Intn(3) + 1 {
		case 1:
			tmpTablero[idx] = "1"
		case 2:
			tmpTablero[idx] = "2"
		case 3:
			tmpTablero[idx] = "3"
		}
	}

	juego.Tablero = strings.Join(tmpTablero, "")
}

func Enviar(n int) {
	conn, _ := net.Dial("tcp", direccionRemota)
	defer conn.Close()
	fmt.Fprintf(conn, "%d", n)
}

func Manejador(conn net.Conn) {
	defer conn.Close()

	pId := len(juego.Jugadores)

	// Configuración inicial del jugador
	fmt.Printf("\n======\n¡Jugador %d se unió a la sesión!\n======\n\n", pId+1)
	juego.Jugadores = append(juego.Jugadores, Jugador{strconv.Itoa(pId + 1), make([]int, nFichasPorJugador), 0})
	for i := 0; i < nFichasPorJugador; i++ {
		juego.Jugadores[pId].FichasPos[i] = 0
	}

	for {
		if win, pWinId = juego.CheckMeta(); win {
			break
		}
		// Obtener resultado del dado
		br := bufio.NewReader(conn)
		str, _ := br.ReadString('\n')
		dados, _ := strconv.Atoi(strings.TrimSpace(str))
		fmt.Printf("El jugador %s se movió %d casillas\n", juego.Jugadores[pId].Nombre, dados)

		// Alias para el jugador actual
		fichaActual := juego.Jugadores[pId].FichasCompletadas

		// Limitar movimiento del jugador a los límites del tablero
		juego.Jugadores[pId].FichasPos[fichaActual] = juego.Jugadores[pId].FichasPos[fichaActual] + dados
		juego.Restringir(pId, fichaActual)

		// Aplicar bonus
		switch juego.GetChar(juego.Jugadores[pId].FichasPos[fichaActual]) {
		case '1':
			juego.Jugadores[pId].FichasPos[fichaActual] = juego.Jugadores[pId].FichasPos[fichaActual] + 3
		case '2':
			juego.Jugadores[pId].FichasPos[fichaActual] = juego.Jugadores[pId].FichasPos[fichaActual] - 3
		case '3':
			juego.Jugadores[pId].FichasPos[fichaActual] = 0
		}
		// Volvemos a corregir la posición del jugador
		juego.Restringir(pId, fichaActual)
	}

	if pWinId == pId {
		fmt.Printf("El ganador es el jugador %d!\n", pId)
	}
}

func main() {
	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del nodo actual: ")
	strPuertoLocal, _ := br.ReadString('\n')
	strPuertoLocal = strings.TrimSpace(strPuertoLocal)
	direccionLocal := fmt.Sprintf("localhost:%s", strPuertoLocal)

	fmt.Printf("Ingrese el puerto del nodo destino: ")
	strPuertoRemoto, _ := br.ReadString('\n')
	strPuertoRemoto = strings.TrimSpace(strPuertoRemoto)
	direccionRemota = fmt.Sprintf("localhost:%s", strPuertoRemoto)

	Init()

	ln, _ := net.Listen("tcp", direccionLocal)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go Manejador(conn)
	}
}
