package main

import (
	"bufio"
	"os"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa            [][]Elemento // grade 2D representando o mapa
	PosX, PosY      int          // posição atual do personagem
	UltimoVisitado  Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg       string       // mensagem para a barra de status
	VelocidadeAtiva bool         // indica se o efeito de velocidade está ativo
}

type Acao struct {
	Tipo      string
	X, Y      int
	DX, DY    int
	Elem      Elemento
	StatusMsg string
}

var (
	Personagem   = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo      = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede       = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao    = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio        = Elemento{' ', CorPadrao, CorPadrao, false}
	Portal       = Elemento{'⚛', CorMagenta, CorPadrao, false}
	PowerUpSpeed = Elemento{'★', CorAmarelo, CorPadrao, false}
)

func jogoNovo() Jogo {
	return Jogo{UltimoVisitado: Vazio}
}

func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	jogo.UltimoVisitado = jogo.Mapa[jogo.PosY][jogo.PosX]

	return nil
}

func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	if jogo.Mapa[y][x].tangivel {
		return false
	}

	return true
}

func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	elemento := jogo.Mapa[y][x] 

	jogo.Mapa[y][x] = jogo.UltimoVisitado
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]
	jogo.Mapa[ny][nx] = elemento
}
