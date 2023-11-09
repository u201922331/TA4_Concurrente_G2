package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var remoteHost string

func Dados() int {
	d1 := rand.Intn(6) + 1
	d2 := rand.Intn(6) + 1
	s := int(math.Pow(-1, float64(rand.Intn(2))))

	resultado := d1 + s*d2

	fmt.Printf("Resultado: %d (", d1)
	if s < 0 {
		fmt.Printf("-")
	} else {
		fmt.Printf("+")
	}
	fmt.Printf(") %d -> %d\n\n", d2, resultado)

	return resultado
}

func Enviar(dadosResultado int) {
	conn, _ := net.Dial("tcp", remoteHost)
	defer conn.Close()
	fmt.Fprintf(conn, "%d", dadosResultado)
}

func main() {
	br := bufio.NewReader(os.Stdin)
	fmt.Print("Remote host: ")
	remoteHost, _ := br.ReadString('\n')
	remoteHost = strings.TrimSpace(remoteHost)

	for {
		fmt.Println("SE LANZARON LOS DADOS...")
		Enviar(Dados())
		time.Sleep(time.Second * 1) // Evitamos spamear la consola
	}
}
