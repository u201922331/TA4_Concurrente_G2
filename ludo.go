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
	Nro_casillas = 21
)

var Fichas_metidas int
var posicion int
var direccionRemota string

func Meta() {
	if posicion >= Nro_casillas {
		Fichas_metidas = Fichas_metidas + 1
		fmt.Println("Metio una ficha")
		posicion = 0
	}
}
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
	if r <= 0 {
		r = 0
	}
	if r >= Nro_casillas-1 {
		r = Nro_casillas - 1
	}
	return r
}

func GenTablero(casillas int) []rune {
	tablero := make([]rune, casillas)
	for i := range tablero {
		tablero[i] = '_' // . -> Casillas en blanco
	}
	tablero[0] = '#'              // # -> Inicio
	tablero[len(tablero)-1] = '#' // $ -> Fin

	umbral := int(float64(casillas) * 0.3) // Solo se llenará el 30% de las casillas en blanco con casillas especiales

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

func enviar(num int) {
	//conectarnos con el nodo remoto y enviar el mensaje con el número
	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()

	fmt.Fprintln(con, num) //envía el número al nodo remoto
}

func manejador(con net.Conn) {
	defer con.Close()
	//leer el mensaje que es enviado de un nodo anterior
	br := bufio.NewReader(con)
	msg, _ := br.ReadString('\n')
	//el msg es un dato que representa un valor numérico
	msg = strings.TrimSpace(msg)
	num, _ := strconv.Atoi(msg) //recuperar el numero que es enviado
	num = num + 1
	//aplicar la lógica del juego
	if Fichas_metidas == 4 {
		//finaliza el Juego y el nodo pierde
		fmt.Println("Ganaste uwu")
	} else {
		a := Dado()
		posicion = posicion + a
		enviar(num)
		fmt.Printf("dado: %d \tacumulado: %d \t%d\n", a, posicion, Fichas_metidas)
		Meta()
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

	//manejar las conexiones entrantes
	ln, _ := net.Listen("tcp", direccionLocal)
	defer ln.Close()

	Fichas_metidas = 0
	posicion = 0
	for {
		//manejo concurrente de las conexiones entrantes
		con, _ := ln.Accept()
		go manejador(con)
	}

}
