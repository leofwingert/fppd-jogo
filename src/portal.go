package main

import (
	"time"
)

func portal(x, y int, acoes chan<- Acao) {
	for {
		acoes <- Acao{Tipo: "spawn", X: x, Y: y, Elem: Portal}
		time.Sleep(5 * time.Second)

		acoes <- Acao{Tipo: "spawn", X: x, Y: y, Elem: Vazio}
		time.Sleep(5 * time.Second)
	}
}

var destinosPortais = map[[2]int][2]int{
	{11, 3}: {43, 3}, // portal em (11,3) leva para (43,3)
	{43, 3}: {11, 3}, // portal em (43,3) leva para (11,3)
}
