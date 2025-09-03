// main.go - Loop principal do jogo com gerenciador concorrente
package main

import (
	"os"
)

func gerenciadorJogo(jogo *Jogo, acoes <-chan Acao) {
	for acao := range acoes {
		switch acao.Tipo {
		case "spawn":
			jogo.Mapa[acao.Y][acao.X] = acao.Elem

		case "status":
			jogo.StatusMsg = "Você entrou no portal!"

		case "moverPersonagem":
			nx, ny := jogo.PosX+acao.DX, jogo.PosY+acao.DY
			if jogoPodeMoverPara(jogo, nx, ny) {
				jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, acao.DX, acao.DY)
				jogo.PosX, jogo.PosY = nx, ny
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

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Inicia sistema de portais concorrente - aparecem e somem a cada 5 segundos
	go portal(11, 3, acoes) // Portal 1 em (11,3)
	go portal(43, 3, acoes) // Portal 2 em (43,3)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, acoes); !continuar {
			break
		}
	}

}
