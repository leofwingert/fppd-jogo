// inimigo.go - Sistema de inimigos concorrentes que patrulham o mapa
package main

import (
	"math"
	"time"
)

// InimigoPatrulheiro representa um inimigo que se move automaticamente
type InimigoPatrulheiro struct {
	X, Y        int
	DirecaoX    int // -1, 0, 1
	DirecaoY    int // -1, 0, 1
	Alcance     int // Distância de detecção do jogador
	Velocidade  time.Duration
	PatrulhaMin [2]int // Área mínima de patrulha
	PatrulhaMax [2]int // Área máxima de patrulha
}

// Inimigo patrulheiro que se move automaticamente e persegue o jogador
func inimigoPatrulheiro(inimigo InimigoPatrulheiro, acoes chan<- Acao, alertas <-chan [2]int, comandos <-chan string) {
	// Spawna o inimigo inicialmente APENAS UMA VEZ
	acoes <- Acao{Tipo: "spawnInimigo", X: inimigo.X, Y: inimigo.Y, Elem: Inimigo}

	ticker := time.NewTicker(inimigo.Velocidade)
	defer ticker.Stop()

	perseguindo := false
	targetX, targetY := 0, 0
	// Posição atual real do inimigo
	posX, posY := inimigo.X, inimigo.Y

	for {
		select {
		case posJogador := <-alertas:
			// Recebeu posição do jogador - verifica se está no alcance
			distancia := calcularDistancia(posX, posY, posJogador[0], posJogador[1])
			if distancia <= float64(inimigo.Alcance) {
				if !perseguindo {
					perseguindo = true
					acoes <- Acao{Tipo: "status", StatusMsg: "Inimigo detectou você!"}
				}
				// Atualiza target mesmo se já estava perseguindo (jogador se move)
				targetX, targetY = posJogador[0], posJogador[1]
			}

		case comando := <-comandos:
			// Recebe comandos externos (parar, continuar, etc.)
			switch comando {
			case "parar":
				perseguindo = false
			case "patrulhar":
				perseguindo = false
			}

		case <-ticker.C:
			// Movement tick
			oldX, oldY := posX, posY
			var newX, newY int

			if perseguindo {
				// Verifica se jogador ainda está no alcance
				distancia := calcularDistancia(posX, posY, targetX, targetY)
				if distancia > float64(inimigo.Alcance) {
					// Jogador saiu do alcance, volta a patrulhar
					perseguindo = false
					acoes <- Acao{Tipo: "status", StatusMsg: "Inimigo perdeu você de vista"}
				}
			}

			if perseguindo {
				// Calcula nova posição perseguindo o jogador
				newX, newY = moverEmDirecaoAo(posX, posY, targetX, targetY)
			} else {
				// Calcula nova posição patrulhando
				tempInimigo := inimigo
				tempInimigo.X, tempInimigo.Y = posX, posY
				newX, newY, inimigo.DirecaoX, inimigo.DirecaoY = patrulhar(tempInimigo)
			}

			// Envia comando de movimento via canal
			acoes <- Acao{
				Tipo: "moverInimigo",
				X:    oldX,
				Y:    oldY,
				DX:   newX - oldX,
				DY:   newY - oldY,
			}

			// Atualiza posição interna (assumindo que movimento será bem-sucedido)
			// Se falhar, não há problema pois o mapa não será alterado
			posX, posY = newX, newY
		}
	}
}

// Sistema de vigilância que monitora a posição do jogador e alerta inimigos
func sistemaVigilancia(jogo *Jogo, alertaInimigos chan<- [2]int) {
	ultimaPosX, ultimaPosY := jogo.PosX, jogo.PosY

	for {
		time.Sleep(500 * time.Millisecond) // Verifica a cada 0.5s

		// Se jogador se moveu, alerta todos os inimigos
		if jogo.PosX != ultimaPosX || jogo.PosY != ultimaPosY {
			select {
			case alertaInimigos <- [2]int{jogo.PosX, jogo.PosY}:
			default: // Non-blocking send
			}
			ultimaPosX, ultimaPosY = jogo.PosX, jogo.PosY
		}
	}
}

// Calcula distância euclidiana entre dois pontos
func calcularDistancia(x1, y1, x2, y2 int) float64 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return math.Sqrt(dx*dx + dy*dy)
}

// Move em direção a um alvo
func moverEmDirecaoAo(fromX, fromY, toX, toY int) (int, int) {
	newX, newY := fromX, fromY

	if fromX < toX {
		newX++
	} else if fromX > toX {
		newX--
	}

	if fromY < toY {
		newY++
	} else if fromY > toY {
		newY--
	}

	return newX, newY
}

// Lógica de patrulha dentro de uma área definida
func patrulhar(inimigo InimigoPatrulheiro) (int, int, int, int) {
	newX := inimigo.X + inimigo.DirecaoX
	newY := inimigo.Y + inimigo.DirecaoY
	newDirX := inimigo.DirecaoX
	newDirY := inimigo.DirecaoY

	// Verifica limites da área de patrulha e inverte direção se necessário
	if newX <= inimigo.PatrulhaMin[0] || newX >= inimigo.PatrulhaMax[0] {
		newDirX = -inimigo.DirecaoX
		newX = inimigo.X + newDirX
	}
	if newY <= inimigo.PatrulhaMin[1] || newY >= inimigo.PatrulhaMax[1] {
		newDirY = -inimigo.DirecaoY
		newY = inimigo.Y + newDirY
	}

	return newX, newY, newDirX, newDirY
}
