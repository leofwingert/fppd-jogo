package main

import (
	"time"
)

func portal(x, y int, acoes chan<- Acao, usado <-chan bool) {
	for {
		acoes <- Acao{Tipo: "spawn", X: x, Y: y, Elem: Portal}

		timeout := time.After(10 * time.Second)
		aberto := true
		for aberto {
			select {
			case <-usado:
				acoes <- Acao{Tipo: "status", StatusMsg: "Você atravessou o portal a tempo!"}
			case <-timeout:
				aberto = false
				acoes <- Acao{Tipo: "status", StatusMsg: "O portal se fechou!"}
			}
		}

		acoes <- Acao{Tipo: "spawn", X: x, Y: y, Elem: Vazio}
		time.Sleep(5 * time.Second) // Espera para reaparecer
	}
}

var destinosPortais = map[[2]int][2]int{
	{11, 3}: {43, 3}, // portal em (11,3) leva para (43,3)
	{43, 3}: {11, 3}, // portal em (43,3) leva para (11,3)
}
