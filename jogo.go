// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
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
	Mapa           [][]Elemento // grade 2D representando o mapa
	PosX, PosY     int          // posição atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg      string       // mensagem para a barra de status
	PortalA        struct{ x, y int }
	PortalB        struct{ x, y int }
	PortalTimer    int64
	PortalAtivo    bool
	// canais para gerenciar portais e posições do jogador
	PosChan             chan [2]int
	TeleportChan        chan [2]int
	PortalVisibleSecond int // duração em segundos do portal visível
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
	Portal     = Elemento{'🔮', CorAzulClaro, CorPadrao, false}
)

func jogoNovo() Jogo {
	j := Jogo{UltimoVisitado: Vazio}
	j.PosChan = make(chan [2]int, 8)
	j.TeleportChan = make(chan [2]int, 4)
	j.PortalVisibleSecond = 6

	return j
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
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
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	// cria os portal
	CriarPortalRandon(jogo)

	// recria os portal
	go func() {
		for {
			time.Sleep(time.Duration(jogo.PortalVisibleSecond) * time.Second)
			// Remove portais atuais
			if jogo.PortalAtivo {
				if jogo.PortalA.y >= 0 && jogo.PortalA.y < len(jogo.Mapa) && jogo.PortalA.x >= 0 && jogo.PortalA.x < len(jogo.Mapa[jogo.PortalA.y]) {
					jogo.Mapa[jogo.PortalA.y][jogo.PortalA.x] = Vazio
				}
				if jogo.PortalB.y >= 0 && jogo.PortalB.y < len(jogo.Mapa) && jogo.PortalB.x >= 0 && jogo.PortalB.x < len(jogo.Mapa[jogo.PortalB.y]) {
					jogo.Mapa[jogo.PortalB.y][jogo.PortalB.x] = Vazio
				}
				jogo.PortalAtivo = false
				jogo.StatusMsg = "Portais desapareceram!"
			}
			time.Sleep(2 * time.Second)
			CriarPortalRandon(jogo)
		}
	}()

	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	// restaura o conteúdo anterior da célula de origem
	jogo.Mapa[y][x] = jogo.UltimoVisitado

	// se o destino for um portal, preservamos o portal em UltimoVisitado
	if jogo.Mapa[ny][nx].simbolo == Portal.simbolo {
		jogo.UltimoVisitado = Portal // quando sair, o portal volta
	} else {
		// guarda o conteúdo atual da célula de destino (normalmente Vazio)
		jogo.UltimoVisitado = jogo.Mapa[ny][nx]
		// move o elemento (normalmente personagem) para o destino
		jogo.Mapa[ny][nx] = elemento
	}
}

func CriarPortalRandon(jogo *Jogo) {
	// Protege contra mapa vazio
	if len(jogo.Mapa) == 0 || len(jogo.Mapa[0]) == 0 {
		return
	}

	w := len(jogo.Mapa[0])
	h := len(jogo.Mapa)
	var ax, ay, bx, by int


	// Gerador local para evitar uso de rand global
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Escolhe posição A livre e diferente do jogador
	tentativas := 0
	for {
		ax = r.Intn(w)
		ay = r.Intn(h)
		if jogo.Mapa[ay][ax].simbolo == Vazio.simbolo && (ax != jogo.PosX || ay != jogo.PosY) {
			break
		}
		tentativas++
		if tentativas > 100 { // evita loop infinito
			jogo.StatusMsg = "DEBUG: Não encontrou posição para portal A"
			return
		}
	}

	// Escolhe posição B livre e diferente de A e do jogador
	tentativas = 0
	for {
		bx = r.Intn(w)
		by = r.Intn(h)
		if jogo.Mapa[by][bx].simbolo == Vazio.simbolo && (bx != ax || by != ay) && (bx != jogo.PosX || by != jogo.PosY) {
			break
		}
		tentativas++
		if tentativas > 100 { // evita loop infinito
			jogo.StatusMsg = "DEBUG: Não encontrou posição para portal B"
			return
		}
	}

	// Marca portais no mapa
	jogo.Mapa[ay][ax] = Portal
	jogo.Mapa[by][bx] = Portal

	// Atualiza posições dos portais
	jogo.PortalA.x = ax
	jogo.PortalA.y = ay
	jogo.PortalB.x = bx
	jogo.PortalB.y = by

	jogo.PortalAtivo = true
	jogo.PortalTimer = time.Now().Unix()
	jogo.StatusMsg = fmt.Sprintf("Portais criados em A:(%d,%d) B:(%d,%d)!", ax, ay, bx, by)
}
