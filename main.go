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
	Nro_casillas      = 21
	nFichasPorJugador = 1
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

func (Tb Tablero) Mantener(num int) {
	if Tb.Jugadores[num].Posicion <= 0 {
		Tb.Jugadores[num].Posicion = 0
	}
	if Tb.Jugadores[num].Posicion >= Nro_casillas-1 {
		Tb.Jugadores[num].Posicion = Nro_casillas - 1
	}
}
func (Tb Tablero) Ganar(num int) {
	if Tb.Jugadores[num].Posicion >= Nro_casillas-1 {
		Tb.Jugadores[num].Fichas_metidas++ //aumentar cantidad de fichas metidas
		fmt.Println("Metio una ficha")

		//actualizar y mostrar posiciones en el tablero
		Tb.Tablero[Tb.Jugadores[0].Posicion] = 'a'
		Tb.Tablero[Tb.Jugadores[1].Posicion] = 'b'
		Tb.Tablero[Tb.Jugadores[2].Posicion] = 'c'
		Tb.Tablero[Tb.Jugadores[3].Posicion] = 'd'
		fmt.Printf("%c\n", Tb.Tablero)

		Tb.Jugadores[num].Posicion = 0 //reiniciar ficha

		//validar ganador
		if Tb.Jugadores[num].Fichas_metidas == nFichasPorJugador {
			fmt.Println("Ganaste uwu")
			os.Exit(0)
		}
	}
}

var direccionRemota string
var Tb = Tablero{[]Jugador{{0, 0}, {0, 0}, {0, 0}, {0, 0}}, make([]rune, Nro_casillas), 1}
var cantidad_de_jugadores = 2

func Dado() int {
	d1 := rand.Intn(6) + 1                        // Dado 1
	d2 := rand.Intn(6) + 1                        // Dado 2
	s := int(math.Pow(-1, float64(rand.Intn(2)))) // Dado Signo
	/*
		fmt.Printf("DADOS: (%d) ", d1)
		if s > 0 {
			fmt.Printf("(+)")
		} else {
			fmt.Printf("(-)")
		}
		fmt.Printf(" (%d)", d2)
	*/
	r := d1 + s*d2
	return r
}

func Enviar(x int) {
	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()
	Tb.C_j = x

	arrBytesJson, _ := json.Marshal(Tb)
	strMsgJson := string(arrBytesJson)

	fmt.Fprintln(con, strMsgJson)

	fmt.Println("Mensaje enviado: ")
	fmt.Println(strMsgJson)
}

func Manejador(con net.Conn) {
	var num int
	ce := 0
	defer con.Close()

	br := bufio.NewReader(con)
	msgJson, _ := br.ReadString('\n')

	json.Unmarshal([]byte(msgJson), &Tb)

	fmt.Println("Mensaje recibido: ")
	fmt.Println(Tb)

	num = Tb.C_j

	//lógica del juego
	if false {
		fmt.Println("xd")
	} else {
		//tirar dado y actualizar posicion
		a := Dado()
		Tb.Jugadores[num].Posicion = Tb.Jugadores[num].Posicion + a

		//mantener en los limites de la cantidad de casillas totales
		Tb.Mantener(num)

		//casillas especiales
		//	1 -> +3 espacios	2 -> -3 espacios	3 -> regresa al principio
		if Tb.Tablero[Tb.Jugadores[num].Posicion] == '1' { //49
			Tb.Jugadores[num].Posicion = Tb.Jugadores[num].Posicion + 3
			ce = 1
		} else if Tb.Tablero[Tb.Jugadores[num].Posicion] == '2' { //50
			Tb.Jugadores[num].Posicion = Tb.Jugadores[num].Posicion - 3
			ce = 2
		} else if Tb.Tablero[Tb.Jugadores[num].Posicion] == '3' { //51
			Tb.Jugadores[num].Posicion = 0
			ce = 3
		}

		//mantener en los limites de la cantidad de casillas totales
		Tb.Mantener(num)

		//D: dado    P: posicion    FM: fichas metidas    T:turno	de jugador x    CE: casilla especial
		fmt.Printf("D: %d\tP: %d\tFM: %d\tT: %d\tCE: %d\n",
			a, Tb.Jugadores[num].Posicion, Tb.Jugadores[num].Fichas_metidas, num, ce)

		//validar si llego a la meta
		Tb.Ganar(num)

		//actualizar turno
		num = num + 1
		if num == cantidad_de_jugadores {
			num = 0
		}

		Enviar(num)
	}

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
