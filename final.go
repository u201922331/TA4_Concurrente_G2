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
type Tablero struct {
	Jugadores []Jugador
	Tablero   []rune
	C_j       int
}

func GenTablero(casillas int) []rune {
	tablero := make([]rune, casillas)
	for i := range tablero {
		tablero[i] = '_' // . -> Casillas en blanco
	}
	tablero[0] = '#'              // # -> Inicio
	tablero[len(tablero)-1] = '#' // $ -> Fin

	umbral := int(float64(casillas) * 0.5) // Solo se llenará el 30% de las casillas en blanco con casillas especiales

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

var Tb = Tablero{[]Jugador{{0, 0}, {0, 0}, {0, 0}, {0, 0}}, GenTablero(Nro_casillas), 1}

func enviar_i(direccionRemota string, x int) {
	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()
	Tb.C_j = x

	arrBytesJson, _ := json.Marshal(Tb)
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

	enviar_i(direccionRemota, 0)
}
