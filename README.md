### Tarea Académica 4

# Informe de implementación del juego Ludo modificado
#### Integrantes
- Daniel Ulises Barrionuevo Gutierrez (u201922128)
- Ramiro Chavez Caituiro (u201524658)
- Nander Emanuel Melendez Huamanchumo (u201922331)

## 1. Descripción y reglas de juego
Este informe presenta la implementación del juego Ludo modificado realizado con el lenguaje de programación GO. El juego se desarrolla en un tablero que contienen obstáculos y casillas especiales. El objetivo es que los jugadores muevan sus fichas desde el inicio hasta la meta. El juego fue implementado utilizando conceptos de concurrencia, comunicación entre procesos(canales) y sistemas distribuídos. Las reglas de juego son las siguientes : 
### 1.1 Tablero del laberinto :
- El tablero está dividido en casillas con caminos, obstáculos.
- Cada jugador tiene cuatro personajes que comienzan en puntos de partida específicos.
### 1.2 Turnos y movimientos: 
- Los jugadores se turnan para lanzar un dado y mover a sus personajes.
- Los jugadores lanzan tres dados, dos dados normales (del 1 al 6) y uno con la operación (sumar o restar) para determinar cuántos pasos pueden avanzar o retroceder en su turno.
- Los jugadores pueden mover un solo personaje por turno. 
- Los personajes deben avanzar exactamente la cantidad de pasos indicados por la operación de los dados (valor del primer dado y operador (+ -) con el valor del segundo dado). 
### 1.3 Obstáculos:
- El laberinto está lleno de obstáculos como paredes, trampas y criaturas que bloquean el paso de los personajes en varias casillas.
- Si al personaje le toca avanzar hacia una casilla con obstáculo entonces el jugador pierde el turno y continua el siguiente jugador.
### 1.4 Objetivo:
- El objetivo es llevar a los cuatro personajes desde los puntos de partida hasta la meta en el menor número de turnos posible.
- El primer jugador en llevar a todos sus personajes a la meta gana el juego.
### 1.5 Uso de sistemas distribuídos
Los jugadores y tablero están representados como entidades distribuídas los cuáles llevarán a cabo la comunicación entre ellas de manera concurrente y sincronizada a traves de los puertos de cada nodo instanciado. la implementación del juego basado en sistemas distribuídos se llevará a cabo bajo la arquitectura circular(Hot Potato).

## 2. Arquitectura 
La arquitectura que se utilizará para la implementación del juego es la circular también conocida como Hot Potato. El componente llamado "iniciador" mandará un mensaje al jugador del primer turno. Una vez este haya realizado su jugada tirando los dados y ubicando su ficha, enviará un mensaje con la información del estado del juego al jugador del siguiente turno. Este ciclo continuará hasta que uno de los jugadores gane el juego. En este punto se enviará un mensaje al tablero y este mostrará los datos relacionados al fin de la partida como qué jugador fue el ganador, el estado en el que se encontran los demás jugadores, etc. A continuación se presenta la arquitectura circular.

![arquitectura circular](https://github.com/u201922331/TA4_Concurrente_G2/assets/117599813/94147f1e-dfff-4acd-b257-f43bfb4c51ad)

## 3. Estructura del juego 
A continuación se presentará la estructura del código implementado del juego de ludo modificado. 

### 3.1. Tipos de datos : 
- **type Jugador:** Representa a un jugador de la partida. Este esta conformado por los variables posición y fichas metidas. Ambas son de tipo de dato entero
- **type Tablero:** Este struct contiene la información sobre los jugadores y el tablero del juego. Esta conformado por las variables Jugadores, Tablero y c_j
  
**Sección código:**
```go
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
```
  
### 3.2 Generación del tablero :
- **func GenTablero:** Esta función genera un tablero con un número específico de casillas, con algunas de ellas marcadas como especiales. Las casillas especiales pueden tener efectos como avanzar, retroceder o regresar al principio.
  
**Sección código:**
```go
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
```
  
### 3.3 Funciones de mantenimiento :
- **func mantener:** Esta función asegura que la posición de un jugador esté dentro de los límites del tablero.
- **func ganar:** Esta función controla si un jugador ha ganado al llegar a la meta y aumenta la cantidad de fichas metidas
- 
**Sección código:**
```go
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
```

### 3.4 Tirar el dado : 
**Func dado:** Esta función simula el lanzamiento de dos dados y una moneda el cuál genera un valor como resultado del lanzamiento. Se realiza una suma entre el valor generado por el primer dado sumado al producto del valor generado por el lanzamiento de la moneda (1 o -1) y el valor generado por el segundo dado. De esa manera se determina cuantos pasas debe avanzar o retroceder la ficha del jugador en cuestión

**Sección código:**
```go
func Dado() int {
	d1 := rand.Intn(6) + 1                        // Dado 1
	d2 := rand.Intn(6) + 1                        // Dado 2
	s := int(math.Pow(-1, float64(rand.Intn(2)))) // Dado Signo
	return d1 + s*d2
}
```

### 3.5 Envío de mensajes:
**func enviar_i:** Esta funcionalidad permite que el iniciador envíe mensajes al nodo principal a través de los protocolos TCP. Los mensajes que se empaquetan en formato JSON contienen la información acerca del estado del juego.

**Sección código:**
```go
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
```

### 3.6 Gestor de comunicación y lógica del juego:
**func Manejador:** Esta función maneja las conexiones entrantes y la comunicación entre nodos, lo que permite que los jugadores interactúen y jueguen en conjunto. El parámetro que admite esta función es "net.Conn" la cuál se trata de una conexión de red que representa una conexión TCP que proviene de otro nodo del sistema distribuído. Esta funcionalidad utiliza un lector "br" para poder leer un mensaje JSON desde la conexión entrante. Este mensaje contiene información sobre el estado del juego, incluyendo la posición de los jugadores y otros detalles. También es necesario el uso de la función "json.Unmarshal" el cuál permite la deserialización del mesaje JSON recibido y actualiza la estructura "Tablero" con la información recibida. "func Manejador" también contiene la lógica del juego el cuál consta de los siguientes pasos:
- Tirar un dado para determinar cuántos espacios avanzará el jugador actual.
- Actualizar la posición del jugador.
- Gestionar casillas especiales que pueden modificar la posición del jugador.
- Validar si el jugador ha ganado (llegado a la meta) y realizar las acciones correspondientes.

Con respecto a la comunicación entre nodos, una vez que se ha realizado la actualización del estado del juego en el nodo actual, la función prepara el estado actualizado para ser transmitido a otros nodos. SeSerializa el estado actualizado en formato JSON y abre una conexión saliente a través de la función Enviar y envía el estado actualizado a otros nodos. Esto permite que los nodos mantengan una sincronización constante del estado del juego.

**Sección código:**
```go
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
```

### 3.7. Consola  :
- **func main_** El programa pide insertar el puerto del nodo remoto. Es de esa manera que se ejecuta el primer paso de la comunicación entre el iniciador y los nodos principales.
  
**Sección código:**
```go
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
```
## 4. Resultados:

**main.go (J1)**
```cmd
Ingrese el puerto del nodo actual: 8000
Ingrese el puerto del nodo destino: 8001
D: 5    P: 5    FM: 0   T: 1    CE: 0
D: -2   P: 0    FM: 0   T: 0    CE: 0
D: -4   P: 1    FM: 0   T: 1    CE: 0
D: 8    P: 8    FM: 0   T: 0    CE: 0
D: 2    P: 3    FM: 0   T: 1    CE: 0
D: -1   P: 7    FM: 0   T: 0    CE: 0
D: 0    P: 3    FM: 0   T: 1    CE: 0
D: -1   P: 6    FM: 0   T: 0    CE: 0
D: 10   P: 13   FM: 0   T: 1    CE: 0
D: -1   P: 5    FM: 0   T: 0    CE: 0
D: 5    P: 18   FM: 0   T: 1    CE: 0
D: 6    P: 11   FM: 0   T: 0    CE: 0
D: 2    P: 20   FM: 0   T: 1    CE: 0
¡Metio una ficha!
[d           a         b]
D: 0    P: 11   FM: 0   T: 0    CE: 0
D: 5    P: 5    FM: 1   T: 1    CE: 0
D: 2    P: 13   FM: 0   T: 0    CE: 0
D: -2   P: 3    FM: 1   T: 1    CE: 0
D: 3    P: 16   FM: 0   T: 0    CE: 0
```

**main.go (J2)**
```cmd
Ingrese el puerto del nodo actual: 8001
Ingrese el puerto del nodo destino: 8000
D: 2    P: 2    FM: 0   T: 1    CE: 0
D: 7    P: 7    FM: 0   T: 0    CE: 0
D: 0    P: 2    FM: 0   T: 1    CE: 0
D: 5    P: 12   FM: 0   T: 0    CE: 0
D: 9    P: 11   FM: 0   T: 1    CE: 0
D: 5    P: 17   FM: 0   T: 0    CE: 0
D: -3   P: 8    FM: 0   T: 1    CE: 0
D: 1    P: 18   FM: 0   T: 0    CE: 0
D: 6    P: 14   FM: 0   T: 1    CE: 0
D: 8    P: 20   FM: 0   T: 0    CE: 0
¡Metio una ficha!
[d              b      a]
D: 5    P: 19   FM: 0   T: 1    CE: 0
D: 3    P: 3    FM: 1   T: 0    CE: 0
D: 4    P: 20   FM: 0   T: 1    CE: 0
¡Metio una ficha!
[d   a           b      b]
D: 2    P: 5    FM: 1   T: 0    CE: 0
D: -1   P: 0    FM: 1   T: 1    CE: 0
D: 7    P: 12   FM: 1   T: 0    CE: 0
D: 4    P: 4    FM: 1   T: 1    CE: 0
D: 8    P: 20   FM: 1   T: 0    CE: 0
¡Metio una ficha!
[d   a b          b      a]
¡Ganaste!
```

**iniciador.go**
```cmd
Ingrese el puerto del nodo remoto: 8000
{
        "Jugadores": [
                {
                        "Posicion": 0,    
                        "FichasMetidas": 0
                },
                {
                        "Posicion": 0,    
                        "FichasMetidas": 0
                },
                {
                        "Posicion": 0,    
                        "FichasMetidas": 0
                },
                {
                        "Posicion": 0,    
                        "FichasMetidas": 0
                }
        ],
        "Tablero": [
                35,
                50,
                95,
                95,
                51,
                49,
                95,
                95,
                95,
                95,
                49,
                95,
                95,
                95,
                95,
                51,
                50,
                50,
                50,
                49,
                35
        ],
        "CurrentJugadorId": 0,
        "WinFlag": false
}
```

## 5. Conclusiones:
- El proyecto aborda una de las cuestiones fundamentales en sistemas distribuidos: la comunicación y la sincronización entre nodos. El proyecto en cuestión ilustra cómo se pueden establecer conexiones y transmitir información entre múltiples nodos. Esto es esencial en sistemas distribuidos, donde múltiples dispositivos o sistemas deben trabajar juntos.
- El proyecto fomenta tanto la cooperación como la competencia entre los jugadores. Los jugadores compiten por ser los primeros en llegar a la meta mientras colaboran en mantener el estado del juego sincronizado entre los nodos. Esto refleja la naturaleza de muchas aplicaciones del mundo real, donde sistemas distribuidos deben cooperar para lograr objetivos comunes.
- La inclusión de casillas especiales que afectan la posición de los jugadores añade un elemento interesante al juego. Esto refleja la necesidad en sistemas distribuidos de gestionar situaciones inesperadas o excepcionales. Los nodos deben adaptarse y tomar decisiones en función de los eventos que ocurren.
- El proyecto ofrece una oportunidad de aprendizaje para comprender los conceptos de sistemas distribuidos, comunicación entre nodos y programación concurrente.
- Este proyecto puede servir como base para futuras mejoras y experimentos. Por ejemplo, se podrían explorar técnicas de escalabilidad para manejar más jugadores o casillas especiales más complejas. Esto puede ayudar a los desarrolladores a comprender los desafíos de implementar sistemas distribuidos a mayor escala.






