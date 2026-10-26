# Apuração 2026 · Presidente 2º turno

App desktop (Wails v2 + Go) com frontend Nuxt 4 (SPA) + Tailwind e mapa maplibre-gl. Sem login.

## Fonte de dados

Somente a divulgação oficial do TSE:

- `https://resultados.tse.jus.br/oficial/ele2026/6258/dados/{uf}/{uf}-c0001-e006258-u.json`
  (`6258` = eleição federal 2026, 2º turno; `c0001` = Presidente; `uf` = `br`, 27 UFs e `zz` exterior)
- Fotos: `https://resultados.tse.jus.br/oficial/ele2026/6258/fotos/br/{sqcand}.jpeg`

O Go consulta as 29 abrangências a cada 30 s (`app.go`), com GET condicional
(ETag/Last-Modified), e emite o evento `apuracao:painel` para o front. Se uma
abrangência falhar, mantém o último dado bom e avisa na tela. Se um nó do CDN
devolver uma cópia mais antiga do que a já recebida, ela é ignorada.

A malha dos estados (`frontend/public/geo/br-ufs.json`) vem da API de malhas do
IBGE e vai embutida no executável: o mapa não depende de tiles externos.

## Rodar

```bash
wails dev      # desenvolvimento (Nuxt em localhost:3000 via proxy do Wails)
wails build    # gera build/bin/apuracao.exe
```
