package main

import (
	"math"
	"time"
)
type InimigoPatrulheiro struct {
	X, Y        int
	DirecaoX    int // -1, 0, 1
	DirecaoY    int // -1, 0, 1
	Alcance     int // Distância de detecção do jogador
	Velocidade  time.Duration
	PatrulhaMin [2]int // Área mínima de patrulha
	PatrulhaMax [2]int // Área máxima de patrulha
}

func inimigoPatrulheiro(inimigo InimigoPatrulheiro, acoes chan<- Acao, alertas <-chan [2]int, comandos <-chan string) {
	//spawna o inimigo
	acoes <- Acao{Tipo: "spawnInimigo", X: inimigo.X, Y: inimigo.Y, Elem: Inimigo}

	ticker := time.NewTicker(inimigo.Velocidade)
	defer ticker.Stop()

	perseguindo := false
	targetX, targetY := 0, 0
	posX, posY := inimigo.X, inimigo.Y

	for {
		select {
		case posJogador := <-alertas:
			distancia := calcularDistancia(posX, posY, posJogador[0], posJogador[1])
			if distancia <= float64(inimigo.Alcance) {
				if !perseguindo {
					perseguindo = true
					acoes <- Acao{Tipo: "status", StatusMsg: "Inimigo detectou você!"}
				}
				targetX, targetY = posJogador[0], posJogador[1]
			}

		case comando := <-comandos:
			switch comando {
			case "parar":
				perseguindo = false
			case "patrulhar":
				perseguindo = false
			}

		case <-ticker.C:
			oldX, oldY := posX, posY
			var newX, newY int

			if perseguindo {
				distancia := calcularDistancia(posX, posY, targetX, targetY)
				if distancia > float64(inimigo.Alcance) {
					perseguindo = false
					acoes <- Acao{Tipo: "status", StatusMsg: "Inimigo perdeu você de vista"}
				}
			}

			if perseguindo {
				newX, newY = moverEmDirecaoAo(posX, posY, targetX, targetY)
			} else {
				tempInimigo := inimigo
				tempInimigo.X, tempInimigo.Y = posX, posY
				newX, newY, inimigo.DirecaoX, inimigo.DirecaoY = patrulhar(tempInimigo)
			}

			acoes <- Acao{
				Tipo: "moverInimigo",
				X:    oldX,
				Y:    oldY,
				DX:   newX - oldX,
				DY:   newY - oldY,
			}

			posX, posY = newX, newY
		}
	}
}

func sistemaVigilancia(jogo *Jogo, alertaInimigos chan<- [2]int) {
	ultimaPosX, ultimaPosY := jogo.PosX, jogo.PosY

	for {
		time.Sleep(500 * time.Millisecond) // Verifica a cada 0.5s

		if jogo.PosX != ultimaPosX || jogo.PosY != ultimaPosY {
			select {
			case alertaInimigos <- [2]int{jogo.PosX, jogo.PosY}:
			default:
			}
			ultimaPosX, ultimaPosY = jogo.PosX, jogo.PosY
		}
	}
}

func calcularDistancia(x1, y1, x2, y2 int) float64 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return math.Sqrt(dx*dx + dy*dy)
}

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

func patrulhar(inimigo InimigoPatrulheiro) (int, int, int, int) {
	newX := inimigo.X + inimigo.DirecaoX
	newY := inimigo.Y + inimigo.DirecaoY
	newDirX := inimigo.DirecaoX
	newDirY := inimigo.DirecaoY

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
