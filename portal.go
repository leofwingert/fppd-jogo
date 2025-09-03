package main

import (
	"time"
)

// Portal aparece e desaparece automaticamente a cada 5 segundos
func portal(x, y int, acoes chan<- Acao) {
	for {
		// Aparece o portal
		acoes <- Acao{Tipo: "spawn", X: x, Y: y, Elem: Portal}

		// Aguarda 5 segundos
		time.Sleep(5 * time.Second)

		// Remove o portal
		acoes <- Acao{Tipo: "spawn", X: x, Y: y, Elem: Vazio}

		// Aguarda 5 segundos antes de aparecer novamente
		time.Sleep(5 * time.Second)
	}
}

var destinosPortais = map[[2]int][2]int{
	{11, 3}: {43, 3}, // portal em (11,3) leva para (43,3)
	{43, 3}: {11, 3}, // portal em (43,3) leva para (11,3)
}
