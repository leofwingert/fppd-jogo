// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
const (
	CorPadrao      Cor = termbox.ColorDefault
	CorCinzaEscuro     = termbox.ColorDarkGray
	CorVermelho        = termbox.ColorRed
	CorVerde           = termbox.ColorGreen
	CorParede          = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede     = termbox.ColorDarkGray
	CorTexto           = termbox.ColorDarkGray
	CorMagenta         = termbox.ColorMagenta
	CorAmarelo         = termbox.ColorYellow
)

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// Renderiza todo o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// Desenha o personagem sobre o mapa
	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)

	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo)

	// Força a atualização do terminal
	interfaceAtualizarTela()
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal com os dados desenhados
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	statusMsg := jogo.StatusMsg
	if jogo.VelocidadeAtiva {
		statusMsg = "VELOCIDADE ATIVA!" + statusMsg
	}

	for i, c := range statusMsg {
		cor := CorTexto
		if jogo.VelocidadeAtiva && i < 22 { // Destaca o indicador de velocidade
			cor = CorVerde
		}
		termbox.SetCell(i, len(jogo.Mapa)+1, c, cor, CorPadrao)
	}

	// instruções fixas
	msg := "Use WASD para mover e E para interagir. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
	}
}
