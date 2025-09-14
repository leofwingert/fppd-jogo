package main

import (
	"time"
)

type PowerUp struct {
	X, Y    int
	Tipo    string
	Ativo   bool
	Duracao time.Duration
}
type SistemaVelocidade struct {
	VelocidadeAtiva bool
	TempoRestante   time.Duration
}

var PowerUpVelocidade = PowerUpSpeed

var posicoesPowerUp = [][2]int{
	{15, 8},
	{35, 12},
	{20, 20},
	{50, 15},
}

func gerenciadorPowerUpsSimples(acoes chan<- Acao, jogo *Jogo) {
	posicaoAtual := 0

	for {
		time.Sleep(10 * time.Second)

		x, y := posicoesPowerUp[posicaoAtual][0], posicoesPowerUp[posicaoAtual][1]

		if y < len(jogo.Mapa) && x < len(jogo.Mapa[y]) && !jogo.Mapa[y][x].tangivel {
			acoes <- Acao{Tipo: "spawn", X: x, Y: y, Elem: PowerUpVelocidade}
			acoes <- Acao{Tipo: "status", StatusMsg: "Power-up de velocidade apareceu!"}

			go func(px, py int) {
				time.Sleep(12 * time.Second)
				if jogo.Mapa[py][px].simbolo == PowerUpSpeed.simbolo {
					acoes <- Acao{Tipo: "spawn", X: px, Y: py, Elem: Vazio}
					acoes <- Acao{Tipo: "status", StatusMsg: "Power-up desapareceu..."}
				}
			}(x, y)
		}
		posicaoAtual = (posicaoAtual + 1) % len(posicoesPowerUp)
	}
}
