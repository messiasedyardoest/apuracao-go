package main

import (
	"context"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const intervalo = 30 * time.Second

// Evento emitido para o front a cada consulta ao TSE.
const eventoPainel = "apuracao:painel"

type App struct {
	ctx      context.Context
	cancelar context.CancelFunc
	tse      *clienteTSE
	forcar   chan struct{}

	mu     sync.RWMutex
	painel Painel
}

func NewApp() *App {
	return &App{
		tse:    novoClienteTSE(),
		forcar: make(chan struct{}, 1),
		painel: Painel{UFs: map[string]*Resultado{}, Falhas: []string{}, Intervalo: int(intervalo.Seconds()), Fonte: baseTSE},
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx, a.cancelar = context.WithCancel(ctx)
	go a.monitorar()
}

func (a *App) shutdown(context.Context) {
	a.cancelar()
}

// monitorar consulta o TSE ao abrir e depois a cada 30 s; Atualizar() antecipa
// a consulta e reinicia a contagem.
func (a *App) monitorar() {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-a.forcar:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
		case <-timer.C:
		}
		painel := a.tse.consultar(a.ctx)
		if a.ctx.Err() != nil {
			return
		}
		painel.Intervalo = int(intervalo.Seconds())
		painel.ProximaEm = time.Now().Add(intervalo).Format(time.RFC3339)
		a.mu.Lock()
		a.painel = painel
		a.mu.Unlock()
		runtime.EventsEmit(a.ctx, eventoPainel, painel)
		timer.Reset(intervalo)
	}
}

// Painel devolve a ultima consulta (o front chama ao montar, antes do 1o evento).
func (a *App) Painel() Painel {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.painel
}

// Atualizar pede uma consulta imediata; o resultado chega pelo evento.
func (a *App) Atualizar() {
	select {
	case a.forcar <- struct{}{}:
	default:
	}
}
