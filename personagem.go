// personagem.go - Funções para movimentação e ações do personagem
package main

import "fmt"

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	// Atualmente apenas exibe uma mensagem de status
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
}

func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo, acoes chan<- Acao) bool {
	switch ev.Tipo {
	case "sair":
		return false

	case "interagir":
		personagemInteragir(jogo)
		if jogo.UltimoVisitado.simbolo == Portal.simbolo {
			// verifica se esse portal tem destino
			origem := [2]int{jogo.PosX, jogo.PosY}
			if destino, ok := destinosPortais[origem]; ok {
				acoes <- Acao{Tipo: "teleportarPersonagem", X: destino[0], Y: destino[1]}
			} else {
				jogo.StatusMsg = "Portal sem destino!"
			}
		} else {
			jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
		}

	case "mover":
		dx, dy := 0, 0
		switch ev.Tecla {
		case 'w':
			dy = -1
		case 'a':
			dx = -1
		case 's':
			dy = 1
		case 'd':
			dx = 1
		}

		// Move uma vez
		acoes <- Acao{Tipo: "moverPersonagem", DX: dx, DY: dy}

		// Se velocidade está ativa, move novamente
		if jogo.VelocidadeAtiva {
			acoes <- Acao{Tipo: "moverPersonagem", DX: dx, DY: dy}
		}
	}
	return true
}
