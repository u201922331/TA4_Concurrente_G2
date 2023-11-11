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
La arquitectura que se utilizará para la implementación del juego es la circular también conocida como Hot Potato. El componente llamado "iniciador" mandará un mensaje al jugador del primer turno. Una vez este haya realizado su jugada tirando los dados y ubicando su ficha, enviará un mensaje con la información del estado del juego al jugador del siguiente turno. Este ciclo continuará hasta que gane uno de los jugadores gane el juego. En este punto se enviará un mensaje al tablero y este mostrará los datos relacionados al fin de la partida como qué jugador fue el ganador, el estado en el que se encontran los demás jugadores, etc. A continuación se presenta la arquitectura circular.

(![arquitectura circular](https://github.com/u201922331/TA4_Concurrente_G2/assets/117599813/94147f1e-dfff-4acd-b257-f43bfb4c51ad))

## 3. Estructura del juego 





