// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"fmt"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w': dy = -1 
	case 'a': dx = -1 
	case 's': dy = 1 
	case 'd': dx = 1 
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}

	// envia posição atual para o gerenciador de portais (não bloqueante)
	select {
	case jogo.PosChan <- [2]int{jogo.PosX, jogo.PosY}:
	default:
	}

	// Teleporte automático se entrou em portal
	if jogo.PortalAtivo {
		if jogo.PosX == jogo.PortalA.x && jogo.PosY == jogo.PortalA.y {
			jogo.PosX, jogo.PosY = jogo.PortalB.x, jogo.PortalB.y
			jogo.StatusMsg = "Você entrou no portal A -> B!"
			return
		}
		if jogo.PosX == jogo.PortalB.x && jogo.PosY == jogo.PortalB.y {
			jogo.PosX, jogo.PosY = jogo.PortalA.x, jogo.PortalA.y
			jogo.StatusMsg = "Você entrou no portal B -> A!"
			return
		}
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	// Atualmente apenas exibe uma mensagem de status
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
	if jogo.PortalAtivo {
		if jogo.PosX == jogo.PortalA.x && jogo.PosY == jogo.PortalA.y {
			jogo.PosX, jogo.PosY = jogo.PortalB.x, jogo.PortalB.y
			jogo.StatusMsg = "Você entrou no portal!"
			return
		}
		if jogo.PosX == jogo.PortalB.x && jogo.PosY == jogo.PortalB.y {
			jogo.PosX, jogo.PosY = jogo.PortalA.x, jogo.PortalA.y
			jogo.StatusMsg = "Você entrou no portal!"
			return
		}
	}
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
	}

	return true // Continua o jogo
}
