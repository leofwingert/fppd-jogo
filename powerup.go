// powerup.go - Sistema de power-ups concorrentes
package main

import (
	"math/rand"
	"time"
)

// PowerUp representa um power-up no jogo
type PowerUp struct {
	X, Y    int
	Tipo    string
	Ativo   bool
	Duracao time.Duration
}

// Estados do sistema de velocidade
type SistemaVelocidade struct {
	VelocidadeAtiva bool
	TempoRestante   time.Duration
}

// Elemento visual do power-up
var PowerUpVelocidade = PowerUpSpeed

// Gerenciador simplificado que apenas spawna power-ups automaticamente
func gerenciadorPowerUpsSimples(acoes chan<- Acao, jogo *Jogo) {
	for {
		time.Sleep(15 * time.Second) // Spawna a cada 15 segundos

		// Encontra posição vazia aleatória para spawnar power-up
		x, y := encontrarPosicaoVazia(jogo)
		if x != -1 && y != -1 {
			// Spawna o power-up
			acoes <- Acao{Tipo: "spawn", X: x, Y: y, Elem: PowerUpVelocidade}
			acoes <- Acao{Tipo: "status", StatusMsg: "Power-up de velocidade apareceu!"}

			// Timeout: power-up desaparece após 12 segundos se não coletado
			go func(px, py int) {
				time.Sleep(12 * time.Second)
				// Verifica se ainda é um power-up antes de remover
				if jogo.Mapa[py][px].simbolo == PowerUpSpeed.simbolo {
					acoes <- Acao{Tipo: "spawn", X: px, Y: py, Elem: Vazio}
					acoes <- Acao{Tipo: "status", StatusMsg: "Power-up desapareceu..."}
				}
			}(x, y)
		}
	}
}
// Encontra uma posição vazia aleatória no mapa para spawnar power-up
func encontrarPosicaoVazia(jogo *Jogo) (int, int) {
	tentativas := 50 // Máximo de tentativas para evitar loop infinito

	for i := 0; i < tentativas; i++ {
		x := rand.Intn(len(jogo.Mapa[0])-2) + 1 // Evita bordas
		y := rand.Intn(len(jogo.Mapa)-2) + 1

		// Verifica se posição está vazia e não é tangível
		if !jogo.Mapa[y][x].tangivel && jogo.Mapa[y][x].simbolo == ' ' {
			// Verifica se não está muito perto do jogador (pelo menos 5 células)
			distX := abs(x - jogo.PosX)
			distY := abs(y - jogo.PosY)
			if distX+distY >= 5 {
				return x, y
			}
		}
	}

	return -1, -1 // Não encontrou posição válida
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
