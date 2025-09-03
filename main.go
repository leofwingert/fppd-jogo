// main.go - Loop principal do jogo com gerenciador concorrente
package main

import (
	"os"
	"time"
)

func gerenciadorJogo(jogo *Jogo, acoes <-chan Acao) {
	for acao := range acoes {
		switch acao.Tipo {
		case "spawn":
			jogo.Mapa[acao.Y][acao.X] = acao.Elem

		case "spawnInimigo":
			// Spawna inimigo apenas se a posição estiver vazia
			if jogo.Mapa[acao.Y][acao.X].simbolo == ' ' {
				jogo.Mapa[acao.Y][acao.X] = acao.Elem
			}

		case "status":
			if acao.StatusMsg != "" {
				jogo.StatusMsg = acao.StatusMsg
			} else {
				jogo.StatusMsg = "Você entrou no portal!"
			}

		case "moverPersonagem":
			nx, ny := jogo.PosX+acao.DX, jogo.PosY+acao.DY
			if jogoPodeMoverPara(jogo, nx, ny) {
				// Verifica se está pisando em um power-up antes de mover
				if jogo.Mapa[ny][nx].simbolo == PowerUpSpeed.simbolo {
					// Coleta automaticamente o power-up
					jogo.Mapa[ny][nx] = Vazio
					jogo.VelocidadeAtiva = true
					jogo.StatusMsg = "VELOCIDADE ATIVADA! Move 2x mais rápido por 8 segundos!"

					// Agenda desativação da velocidade em uma goroutine separada
					go func(g *Jogo) {
						time.Sleep(8 * time.Second)
						g.VelocidadeAtiva = false
						g.StatusMsg = "Efeito de velocidade terminou."
					}(jogo)
				}

				jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, acao.DX, acao.DY)
				jogo.PosX, jogo.PosY = nx, ny
			}

		case "moverInimigo":
			// Move inimigo apenas se a posição de destino for válida
			nx, ny := acao.X+acao.DX, acao.Y+acao.DY

			// Verifica limites e colisão
			if nx >= 0 && ny >= 0 && ny < len(jogo.Mapa) && nx < len(jogo.Mapa[ny]) {
				if !jogo.Mapa[ny][nx].tangivel && jogo.Mapa[ny][nx].simbolo != Inimigo.simbolo {
					// Move o inimigo apenas se destino estiver livre
					if jogo.Mapa[acao.Y][acao.X].simbolo == Inimigo.simbolo {
						jogo.Mapa[acao.Y][acao.X] = Vazio
						jogo.Mapa[ny][nx] = Inimigo
					}
				}
				// Se não pode mover, o inimigo fica na posição atual (sem criar cópias)
			}

		case "teleportarPersonagem":
			// Remove o personagem da posição atual
			jogo.Mapa[jogo.PosY][jogo.PosX] = jogo.UltimoVisitado

			// Atualiza a posição do personagem
			jogo.PosX, jogo.PosY = acao.X, acao.Y

			// Guarda o que estava na nova posição e coloca o personagem lá
			jogo.UltimoVisitado = jogo.Mapa[jogo.PosY][jogo.PosX]
			jogo.Mapa[jogo.PosY][jogo.PosX] = Personagem

			jogo.StatusMsg = "Você usou o portal!"

		case "ativarVelocidade":
			// Ativa efeito de velocidade no jogador
			jogo.VelocidadeAtiva = true
			if acao.StatusMsg != "" {
				jogo.StatusMsg = acao.StatusMsg
			}

		case "desativarVelocidade":
			// Desativa efeito de velocidade
			jogo.VelocidadeAtiva = false
			if acao.StatusMsg != "" {
				jogo.StatusMsg = acao.StatusMsg
			}

		}

		interfaceDesenharJogo(jogo)
	}
}

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Inicializa canal de ações e gerenciador concorrente
	acoes := make(chan Acao)
	go gerenciadorJogo(&jogo, acoes)

	// Canais para comunicação com inimigos
	alertaInimigos := make(chan [2]int, 10)
	comandosInimigo1 := make(chan string, 5)

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Inicia sistema de portais concorrente
	go portal(11, 3, acoes) // Portal 1 em (11,3)
	go portal(43, 3, acoes) // Portal 2 em (43,3)

	// Inicia sistema de inimigos concorrente
	inimigo1 := InimigoPatrulheiro{
		X: 43, Y: 14,
		DirecaoX: 1, DirecaoY: 0,
		Alcance:     5,
		Velocidade:  1 * time.Second,
		PatrulhaMin: [2]int{10, 8},
		PatrulhaMax: [2]int{25, 15},
	}
	go inimigoPatrulheiro(inimigo1, acoes, alertaInimigos, comandosInimigo1)
	go sistemaVigilancia(&jogo, alertaInimigos)

	// Inicia sistema de power-ups
	go gerenciadorPowerUpsSimples(acoes, &jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, acoes); !continuar {
			break
		}
	}

}
