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

				if jogo.Mapa[ny][nx].simbolo == PowerUpSpeed.simbolo {
					// coleta o power-up
					jogo.Mapa[ny][nx] = Vazio
					jogo.VelocidadeAtiva = true
					jogo.StatusMsg = "VELOCIDADE ATIVADA! Move 2x mais rápido por 8 segundos!"

					// agenda desativação da velocidade em uma goroutine separada
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
			nx, ny := acao.X+acao.DX, acao.Y+acao.DY

			// verifica limites e colisão
			if nx >= 0 && ny >= 0 && ny < len(jogo.Mapa) && nx < len(jogo.Mapa[ny]) {
				if !jogo.Mapa[ny][nx].tangivel && jogo.Mapa[ny][nx].simbolo != Inimigo.simbolo {
					// Move o inimigo apenas se destino estiver livre
					if jogo.Mapa[acao.Y][acao.X].simbolo == Inimigo.simbolo {
						jogo.Mapa[acao.Y][acao.X] = Vazio
						jogo.Mapa[ny][nx] = Inimigo
					}
				}
			}

		case "teleportarPersonagem":
			jogo.Mapa[jogo.PosY][jogo.PosX] = jogo.UltimoVisitado

			jogo.PosX, jogo.PosY = acao.X, acao.Y

			jogo.UltimoVisitado = jogo.Mapa[jogo.PosY][jogo.PosX]
			jogo.Mapa[jogo.PosY][jogo.PosX] = Personagem

			jogo.StatusMsg = "Você usou o portal!"

		case "ativarVelocidade":
			jogo.VelocidadeAtiva = true
			if acao.StatusMsg != "" {
				jogo.StatusMsg = acao.StatusMsg
			}

		case "desativarVelocidade":
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

	acoes := make(chan Acao)
	go gerenciadorJogo(&jogo, acoes)

	// canais para comunicação com inimigos
	alertaInimigos := make(chan [2]int, 10)
	comandosInimigo1 := make(chan string, 5)

	// desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	go portal(11, 3, acoes) // Portal 1 em (11,3)
	go portal(43, 3, acoes) // Portal 2 em (43,3)

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

	go gerenciadorPowerUpsSimples(acoes, &jogo)

	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, acoes); !continuar {
			break
		}
	}

}
