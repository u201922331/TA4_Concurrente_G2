package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	nCasillas = 21
)

type Jugador struct {
	Posicion      int
	FichasMetidas int
}
type Juego struct {
	Jugadores        []Jugador
	Tablero          []rune
	CurrentJugadorId int
	WinFlag          bool
}

func (juego *Juego) GenTablero(n int) {
	tmp := strings.Split(strings.Repeat("_", n), "")
	tmp[0] = "#"
	tmp[len(tmp)-1] = "#"

	for i := 0; i < int(float64(n)*0.5); i++ { // Umbral de 50%
		idx := rand.Intn(n-2) + 1                 // Escogemos una casilla aleatoria (omitiendo inicio y final)
		tmp[idx] = strconv.Itoa(rand.Intn(3) + 1) // 1: +3 | 2: -3 | 3: Al inicio
	}

	juego.Tablero = []rune{}
	for _, str := range tmp {
		juego.Tablero = append(juego.Tablero, []rune(str)...)
	}
}

var juego = Juego{[]Jugador{{0, 0}, {0, 0}, {0, 0}, {0, 0}}, []rune{}, 1, false}

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
	juego.GenTablero(nCasillas)

	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del nodo remoto: ")
	puertoRemoto, _ := br.ReadString('\n')
	puertoRemoto = strings.TrimSpace(puertoRemoto)
	direccionRemota := fmt.Sprintf("localhost:%s", "8000")

	Enviar(direccionRemota, 0)
}
