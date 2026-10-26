package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Divulgacao oficial do TSE (resultados.tse.jus.br), eleicao federal 2026,
// 2o turno (cd 6258 em comum/config/ele-c.json), cargo 0001 = Presidente.
const (
	baseTSE    = "https://resultados.tse.jus.br/oficial/ele2026/6258"
	eleicaoTSE = "006258"
	cargoTSE   = "0001"
)

// Abrangencias publicadas pelo TSE: Brasil, 27 UFs e exterior (zz).
var abrangencias = []string{
	"br", "ac", "al", "am", "ap", "ba", "ce", "df", "es", "go", "ma", "mg", "ms", "mt",
	"pa", "pb", "pe", "pi", "pr", "rj", "rn", "ro", "rr", "rs", "sc", "se", "sp", "to", "zz",
}

type Candidato struct {
	Numero       string  `json:"numero"`
	Nome         string  `json:"nome"`
	NomeCompleto string  `json:"nomeCompleto"`
	Partido      string  `json:"partido"`
	Vice         string  `json:"vice"`
	VicePartido  string  `json:"vicePartido"`
	Votos        int64   `json:"votos"`
	Percentual   float64 `json:"percentual"`
	Eleito       bool    `json:"eleito"`
	Situacao     string  `json:"situacao"`
	Foto         string  `json:"foto"`
}

type Resultado struct {
	UF string `json:"uf"`
	// n = nao iniciada, p = parcial, f = finalizada.
	Andamento  string `json:"andamento"`
	Definido   bool   `json:"definido"`
	Totalizado string `json:"totalizado"`
	Gerado     string `json:"gerado"`

	UrnasTotal    int64   `json:"urnasTotal"`
	UrnasApuradas int64   `json:"urnasApuradas"`
	PercApuradas  float64 `json:"percApuradas"`

	Eleitorado         int64   `json:"eleitorado"`
	Comparecimento     int64   `json:"comparecimento"`
	PercComparecimento float64 `json:"percComparecimento"`
	Abstencao          int64   `json:"abstencao"`
	PercAbstencao      float64 `json:"percAbstencao"`

	Validos     int64   `json:"validos"`
	PercValidos float64 `json:"percValidos"`
	Brancos     int64   `json:"brancos"`
	PercBrancos float64 `json:"percBrancos"`
	Nulos       int64   `json:"nulos"`
	PercNulos   float64 `json:"percNulos"`

	Candidatos []Candidato `json:"candidatos"`

	geradoEm time.Time
}

// arquivoTSE: so os campos usados do {uf}-c0001-e006258-u.json. O TSE manda
// numeros como string, com virgula decimal nos percentuais.
type arquivoTSE struct {
	Cdabr string `json:"cdabr"`
	Dg    string `json:"dg"`
	Hg    string `json:"hg"`
	Dt    string `json:"dt"`
	Ht    string `json:"ht"`
	And   string `json:"and"`
	Md    string `json:"md"`
	Carg  []struct {
		Agr []struct {
			Par []struct {
				Sg   string `json:"sg"`
				Cand []struct {
					N      string `json:"n"`
					Sqcand string `json:"sqcand"`
					Nm     string `json:"nm"`
					Nmu    string `json:"nmu"`
					E      string `json:"e"`
					St     string `json:"st"`
					Vap    string `json:"vap"`
					Pvapn  string `json:"pvapn"`
					Vs     []struct {
						Nmu string `json:"nmu"`
						Sgp string `json:"sgp"`
					} `json:"vs"`
				} `json:"cand"`
			} `json:"par"`
		} `json:"agr"`
	} `json:"carg"`
	S struct {
		Ts   string `json:"ts"`
		St   string `json:"st"`
		Pstn string `json:"pstn"`
	} `json:"s"`
	E struct {
		Te  string `json:"te"`
		C   string `json:"c"`
		Pcn string `json:"pcn"`
		A   string `json:"a"`
		Pan string `json:"pan"`
	} `json:"e"`
	V struct {
		Vvc    string `json:"vvc"`
		Pvvcn  string `json:"pvvcn"`
		Vb     string `json:"vb"`
		Pvbn   string `json:"pvbn"`
		Tvn    string `json:"tvn"`
		Ptvnn  string `json:"ptvnn"`
	} `json:"v"`
}

func inteiro(s string) int64 {
	n, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return n
}

func decimal(s string) float64 {
	f, _ := strconv.ParseFloat(strings.Replace(strings.TrimSpace(s), ",", ".", 1), 64)
	return f
}

func juntarDataHora(d, h string) string {
	return strings.TrimSpace(d + " " + h)
}

func (a *arquivoTSE) resultado(uf string) *Resultado {
	r := &Resultado{
		UF:                 uf,
		Andamento:          a.And,
		Definido:           a.Md == "s",
		Totalizado:         juntarDataHora(a.Dt, a.Ht),
		Gerado:             juntarDataHora(a.Dg, a.Hg),
		UrnasTotal:         inteiro(a.S.Ts),
		UrnasApuradas:      inteiro(a.S.St),
		PercApuradas:       decimal(a.S.Pstn),
		Eleitorado:         inteiro(a.E.Te),
		Comparecimento:     inteiro(a.E.C),
		PercComparecimento: decimal(a.E.Pcn),
		Abstencao:          inteiro(a.E.A),
		PercAbstencao:      decimal(a.E.Pan),
		Validos:            inteiro(a.V.Vvc),
		PercValidos:        decimal(a.V.Pvvcn),
		Brancos:            inteiro(a.V.Vb),
		PercBrancos:        decimal(a.V.Pvbn),
		Nulos:              inteiro(a.V.Tvn),
		PercNulos:          decimal(a.V.Ptvnn),
		Candidatos:         []Candidato{},
	}
	r.geradoEm, _ = time.Parse("02/01/2006 15:04:05", r.Gerado)
	for _, cargo := range a.Carg {
		for _, agr := range cargo.Agr {
			for _, par := range agr.Par {
				for _, c := range par.Cand {
					cand := Candidato{
						Numero:       c.N,
						Nome:         c.Nmu,
						NomeCompleto: c.Nm,
						Partido:      par.Sg,
						Votos:        inteiro(c.Vap),
						Percentual:   decimal(c.Pvapn),
						Eleito:       c.E == "s",
						Situacao:     c.St,
						Foto:         fmt.Sprintf("%s/fotos/br/%s.jpeg", baseTSE, c.Sqcand),
					}
					if len(c.Vs) > 0 {
						cand.Vice, cand.VicePartido = c.Vs[0].Nmu, c.Vs[0].Sgp
					}
					r.Candidatos = append(r.Candidatos, cand)
				}
			}
		}
	}
	// Mesmo criterio do TSE: mais votado primeiro; empate (inclusive 0 x 0
	// antes da apuracao) resolve pelo numero, para a ordem nao "piscar".
	sort.SliceStable(r.Candidatos, func(i, j int) bool {
		if r.Candidatos[i].Votos != r.Candidatos[j].Votos {
			return r.Candidatos[i].Votos > r.Candidatos[j].Votos
		}
		return r.Candidatos[i].Numero < r.Candidatos[j].Numero
	})
	return r
}

type entradaCache struct {
	etag       string
	modificado string
	resultado  *Resultado
}

type clienteTSE struct {
	http  *http.Client
	mu    sync.Mutex
	cache map[string]entradaCache
}

func novoClienteTSE() *clienteTSE {
	return &clienteTSE{
		http:  &http.Client{Timeout: 15 * time.Second},
		cache: map[string]entradaCache{},
	}
}

// buscar baixa uma abrangencia com GET condicional (ETag/Last-Modified): sem
// mudanca o CDN do TSE responde 304 e reaproveitamos o cache.
func (c *clienteTSE) buscar(ctx context.Context, uf string) (*Resultado, error) {
	url := fmt.Sprintf("%s/dados/%s/%s-c%s-e%s-u.json", baseTSE, uf, uf, cargoTSE, eleicaoTSE)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	anterior, temAnterior := c.cache[uf]
	c.mu.Unlock()
	if temAnterior {
		if anterior.etag != "" {
			req.Header.Set("If-None-Match", anterior.etag)
		}
		if anterior.modificado != "" {
			req.Header.Set("If-Modified-Since", anterior.modificado)
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified && temAnterior {
		return anterior.resultado, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	var arquivo arquivoTSE
	if err := json.NewDecoder(resp.Body).Decode(&arquivo); err != nil {
		return nil, fmt.Errorf("JSON invalido: %w", err)
	}
	novo := arquivo.resultado(uf)
	// Bordas diferentes do CDN podem entregar uma copia mais antiga do que a
	// ja recebida (visto no AC: indice parcial, arquivo ainda zerado); nesse
	// caso mantem o mais recente.
	if temAnterior && novo.geradoEm.Before(anterior.resultado.geradoEm) {
		return anterior.resultado, nil
	}
	c.mu.Lock()
	c.cache[uf] = entradaCache{
		etag:       resp.Header.Get("ETag"),
		modificado: resp.Header.Get("Last-Modified"),
		resultado:  novo,
	}
	c.mu.Unlock()
	return novo, nil
}

type Painel struct {
	Brasil       *Resultado            `json:"brasil"`
	UFs          map[string]*Resultado `json:"ufs"`
	ConsultadoEm string                `json:"consultadoEm"`
	ProximaEm    string                `json:"proximaEm"`
	Intervalo    int                   `json:"intervalo"`
	Falhas       []string              `json:"falhas"`
	Fonte        string                `json:"fonte"`
}

// consultar busca todas as abrangencias em paralelo (8 por vez). Quem falhar
// fica com o ultimo dado bom do cache e entra em Falhas.
func (c *clienteTSE) consultar(ctx context.Context) Painel {
	painel := Painel{UFs: map[string]*Resultado{}, Falhas: []string{}, Fonte: baseTSE}
	var mu sync.Mutex
	var wg sync.WaitGroup
	vagas := make(chan struct{}, 8)
	for _, uf := range abrangencias {
		wg.Add(1)
		go func(uf string) {
			defer wg.Done()
			vagas <- struct{}{}
			defer func() { <-vagas }()
			res, err := c.buscar(ctx, uf)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				painel.Falhas = append(painel.Falhas, uf)
				c.mu.Lock()
				res = c.cache[uf].resultado
				c.mu.Unlock()
			}
			if res == nil {
				return
			}
			if uf == "br" {
				painel.Brasil = res
			} else {
				painel.UFs[uf] = res
			}
		}(uf)
	}
	wg.Wait()
	sort.Strings(painel.Falhas)
	painel.ConsultadoEm = time.Now().Format(time.RFC3339)
	return painel
}
