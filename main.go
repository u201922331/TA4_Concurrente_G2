package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"strings"
)

const (
	nCasillas         = 21
	nFichasPorJugador = 2
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

func (juego *Juego) Print() {
	fmt.Println(juego.Tablero)
	for _, jugador := range juego.Jugadores {
		for i := 0; i < jugador.Posicion; i++ {
			fmt.Print(" ")
		}
		fmt.Print("¡")
		for i := jugador.Posicion + 1; i < nCasillas; i++ {
			fmt.Print(" ")
		}
		fmt.Print("\n\n")
	}
}

func (juego *Juego) Mantener(idx int) {
	juego.Jugadores[idx].Posicion = int(math.Max(math.Min(float64(juego.Jugadores[idx].Posicion), nCasillas-1), 0))
}

func (juego *Juego) Ganar(num int) {
	if juego.Jugadores[num].Posicion >= nCasillas-1 {
		juego.Jugadores[num].FichasMetidas++ //aumentar cantidad de fichas metidas
		fmt.Println("¡Metio una ficha!")

		//actualizar y mostrar posiciones en el tablero
		juego.Tablero[juego.Jugadores[0].Posicion] = 'a'
		juego.Tablero[juego.Jugadores[1].Posicion] = 'b'
		juego.Tablero[juego.Jugadores[2].Posicion] = 'c'
		juego.Tablero[juego.Jugadores[3].Posicion] = 'd'
		fmt.Printf("%c\n", juego.Tablero)

		juego.Jugadores[num].Posicion = 0 //reiniciar ficha

		//validar ganador
		if juego.Jugadores[num].FichasMetidas == nFichasPorJugador {
			fmt.Println("¡Ganaste!")
			juego.WinFlag = true
			os.Exit(0)
		}
	}
}

var direccionRemota string
var juego = Juego{[]Jugador{{0, 0}, {0, 0}, {0, 0}, {0, 0}}, make([]rune, nCasillas), 1, false}
var cantidad_de_jugadores = 2

func Dado() int {
	d1 := rand.Intn(6) + 1                        // Dado 1
	d2 := rand.Intn(6) + 1                        // Dado 2
	s := int(math.Pow(-1, float64(rand.Intn(2)))) // Dado Signo
	return d1 + s*d2
}

func Enviar(currentJugadorId int) {
	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()
	juego.CurrentJugadorId = currentJugadorId

	arrBytesJson, _ := json.MarshalIndent(juego, "", "\t")
	strMsgJson := string(arrBytesJson)

	fmt.Fprintln(con, strMsgJson)

	/*
		fmt.Println("Mensaje enviado: ")
		fmt.Println(strMsgJson)
	*/
}

func Manejador(con net.Conn) {
	var num int
	ce := 0
	defer con.Close()

	br := bufio.NewReader(con)
	msgJson, _ := br.ReadString('\n')

	json.Unmarshal([]byte(msgJson), &juego)

	/*
		fmt.Println("Mensaje recibido: ")
		fmt.Println(juego)
	*/

	num = juego.CurrentJugadorId

	//lógica del juego
	// ==================
	// tirar dado y actualizar posicion
	dado := Dado()
	juego.Jugadores[num].Posicion = juego.Jugadores[num].Posicion + dado
	juego.Mantener(num)

	//casillas especiales
	//	1 -> +3 espacios	2 -> -3 espacios	3 -> regresa al principio
	if juego.Tablero[juego.Jugadores[num].Posicion] == '1' {
		juego.Jugadores[num].Posicion = juego.Jugadores[num].Posicion + 3
		ce = 1
	} else if juego.Tablero[juego.Jugadores[num].Posicion] == '2' {
		juego.Jugadores[num].Posicion = juego.Jugadores[num].Posicion - 3
		ce = 2
	} else if juego.Tablero[juego.Jugadores[num].Posicion] == '3' {
		juego.Jugadores[num].Posicion = 0
		ce = 3
	}

	//mantener en los limites de la cantidad de casillas totales
	juego.Mantener(num)

	//D: dado    P: posicion    FM: fichas metidas    T:turno	de jugador x    CE: casilla especial
	fmt.Printf("D: %d\tP: %d\tFM: %d\tT: %d\tCE: %d\n",
		dado, juego.Jugadores[num].Posicion, juego.Jugadores[num].FichasMetidas, num, ce)

	//validar si llego a la meta
	juego.Ganar(num)

	//actualizar turno
	num = num + 1
	if num == cantidad_de_jugadores {
		num = 0
	}

	Enviar(num)
}

func main() {
	//fijar el puerto del nodo local
	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del nodo actual: ")
	strPuertoLocal, _ := br.ReadString('\n')
	strPuertoLocal = strings.TrimSpace(strPuertoLocal)
	direccionLocal := fmt.Sprintf("localhost:%s", strPuertoLocal)

	//fijar el puerto del nodo destino
	fmt.Print("Ingrese el puerto del nodo destino: ")
	strPuertoRemoto, _ := br.ReadString('\n')
	strPuertoRemoto = strings.TrimSpace(strPuertoRemoto)
	direccionRemota = fmt.Sprintf("localhost:%s", strPuertoRemoto)

	//conexiones entrantes
	ln, _ := net.Listen("tcp", direccionLocal)
	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go Manejador(con)
	}

}
