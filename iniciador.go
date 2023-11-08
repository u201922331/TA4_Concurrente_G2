package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func enviar(direccionRemota string, num int) {
	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()

	fmt.Fprintln(con, num)
}

func main() {
	//se conecta a uno de los nodos remotos de HP
	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del nodo remoto: ")
	puertoRemoto, _ := br.ReadString('\n')
	puertoRemoto = strings.TrimSpace(puertoRemoto)
	direccionRemota := fmt.Sprintf("localhost:%s", puertoRemoto)

	fmt.Print("Ingrese el número a enviar: ")
	strNum, _ := br.ReadString('\n')
	strNum = strings.TrimSpace(strNum)
	num, _ := strconv.Atoi(strNum)

	enviar(direccionRemota, num)
}
