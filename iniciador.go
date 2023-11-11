package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
)

const (
	Nro_casillas = 21
)

type Jugador struct {
	Posicion       int
	Fichas_metidas int
}
type Juego struct {
	Jugadores        []Jugador
	Tablero          []rune
	CurrentJugadorId int
}

func GenTablero(casillas int) []rune {
	tablero := make([]rune, casillas)
	for i := range tablero {
		tablero[i] = '_' // . -> Casillas en blanco
	}
	tablero[0] = '#'              // # -> Inicio
	tablero[len(tablero)-1] = '#' // $ -> Fin

	umbral := int(float64(casillas) * 0.5) // Solo se llenará el 50% de las casillas en blanco con casillas especiales

	for i := 0; i < umbral; i++ {
		idx := rand.Intn(int(casillas))
		// No afectar las casillas inicial y final
		if idx == 0 || idx == len(tablero)-1 {
			continue
		}

		switch rand.Intn(3) + 1 {
		case 1:
			tablero[idx] = '1' // +3 espacios
		case 2:
			tablero[idx] = '2' // -3 espacios
		case 3:
			tablero[idx] = '3' // regresa al principio
		}
	}
	return tablero
}

var juego = Juego{[]Jugador{{0, 0}, {0, 0}, {0, 0}, {0, 0}}, GenTablero(Nro_casillas), 1}

func Enviar(direccionRemota string, currentJugadorId int) {
	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()
	juego.CurrentJugadorId = currentJugadorId

	arrBytesJson, _ := json.MarshalIndent(juego, "", "\t")
	strMsgJson := string(arrBytesJson)
	//fmt.Println("Mensaje enviado: ")
	fmt.Println(strMsgJson)
	fmt.Fprintln(con, strMsgJson)
}

func main() {
	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del nodo remoto: ")
	puertoRemoto, _ := br.ReadString('\n')
	puertoRemoto = strings.TrimSpace(puertoRemoto)
	direccionRemota := fmt.Sprintf("localhost:%s", "8000")

	Enviar(direccionRemota, 0)
}
