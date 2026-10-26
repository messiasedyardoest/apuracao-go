<script setup lang="ts">
import { IconAlertTriangle, IconRefresh, IconX } from "@tabler/icons-vue";

const { painel, segundosRestantes, atualizando, atualizar, semWails } = useApuracao();
const { fmtNumero, fmtPct, fmtHora } = useFormatadores();

const ufSelecionada = ref<string | null>(null);
const brasil = computed(() => painel.value?.brasil);
const ufs = computed(() => painel.value?.ufs ?? {});
const escopo = computed(() => (ufSelecionada.value ? ufs.value[ufSelecionada.value] : undefined) ?? brasil.value);
// Cores pelos candidatos do Brasil: as UFs tem os mesmos, a legenda e o mapa batem.
const cores = computed(() => coresCandidatos(brasil.value?.candidatos ?? []));

const ANDAMENTO: Record<string, { rotulo: string; cor: string }> = {
  n: { rotulo: "Não iniciada", cor: "var(--text-muted)" },
  p: { rotulo: "Parcial", cor: "var(--warning)" },
  f: { rotulo: "Finalizada", cor: "var(--success)" },
};

function dataHora(totalizado: string) {
  const [data, hora] = totalizado.split(" ");
  return { data, hora };
}
</script>

<template>
  <div class="flex h-screen flex-col">
    <BarraTitulo />
    <main class="min-h-0 flex-1 overflow-y-auto">
      <div class="mx-auto max-w-[1600px] space-y-4 p-4 lg:p-6">
        <header class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h1 class="text-xl font-semibold tracking-tight">Apuração 2026 · Presidente · 2º turno</h1>
            <p class="text-sm text-tinta-dim">
              Fonte oficial: TSE (resultados.tse.jus.br)
              <template v-if="brasil?.totalizado">
                · totalizado em {{ dataHora(brasil.totalizado).data }} às {{ dataHora(brasil.totalizado).hora }}
              </template>
            </p>
          </div>
          <div class="flex flex-wrap items-center gap-3">
            <span
              v-if="brasil"
              class="inline-flex items-center gap-1.5 rounded-full border border-borda bg-cartao px-3 py-1 text-sm font-medium"
            >
              <span class="h-2 w-2 rounded-full" :style="{ background: ANDAMENTO[brasil.andamento]?.cor ?? 'var(--text-muted)' }" />
              {{ ANDAMENTO[brasil.andamento]?.rotulo ?? brasil.andamento }}
            </span>
            <span
              v-if="brasil?.definido"
              class="rounded-full px-3 py-1 text-sm font-semibold text-white"
              style="background: var(--success)"
            >Matematicamente definido</span>
            <span v-if="painel?.consultadoEm" class="text-sm tabular-nums text-tinta-dim">
              Consultado às {{ fmtHora(painel.consultadoEm) }}
              <template v-if="segundosRestantes !== null"> · próxima em {{ segundosRestantes }} s</template>
            </span>
            <button class="btn" :disabled="atualizando || semWails" @click="atualizar">
              <IconRefresh :size="16" :class="{ 'animate-spin': atualizando }" />
              Atualizar agora
            </button>
          </div>
        </header>

        <div
          v-if="painel?.falhas?.length"
          class="flex items-center gap-2 rounded-xl border px-4 py-2.5 text-sm"
          style="border-color: rgba(180, 83, 9, 0.25); background: #fff7ed; color: var(--warning)"
        >
          <IconAlertTriangle :size="18" class="shrink-0" />
          Falha ao consultar {{ painel.falhas.map((u) => u.toUpperCase()).join(", ") }} — exibindo o último dado recebido.
        </div>

        <div v-if="semWails" class="cartao p-8 text-center text-sm text-tinta-dim">
          Abra pelo aplicativo (<code>wails dev</code> ou o executável gerado): no navegador puro não há o backend Go que consulta o TSE.
        </div>

        <div v-else-if="!brasil" class="cartao p-8 text-center text-sm text-tinta-dim">
          <template v-if="painel?.consultadoEm">
            Sem resposta do TSE. Nova tentativa em {{ segundosRestantes ?? 30 }} s.
          </template>
          <template v-else>Consultando o TSE…</template>
        </div>

        <div v-else class="grid grid-cols-1 gap-4 lg:grid-cols-[minmax(0,5fr)_minmax(0,6fr)]">
          <div class="min-w-0 space-y-4">
            <section v-if="escopo" class="cartao space-y-4 p-5">
              <header class="flex items-start justify-between gap-3">
                <div>
                  <p class="text-xs font-medium uppercase tracking-wide text-tinta-muted">
                    {{ ufSelecionada ? "Estado" : "Resultado nacional" }}
                  </p>
                  <h2 class="text-lg font-semibold">{{ NOMES_UF[escopo.uf] ?? escopo.uf.toUpperCase() }}</h2>
                </div>
                <button v-if="ufSelecionada" class="btn py-1.5" @click="ufSelecionada = null">
                  Ver Brasil
                  <IconX :size="14" />
                </button>
              </header>

              <div>
                <div class="flex items-baseline justify-between text-sm">
                  <span class="text-tinta-dim">Urnas apuradas</span>
                  <span class="tabular-nums">
                    <strong class="text-base">{{ fmtPct(escopo.percApuradas) }}</strong>
                    <span class="text-tinta-muted"> · {{ fmtNumero(escopo.urnasApuradas) }} de {{ fmtNumero(escopo.urnasTotal) }}</span>
                  </span>
                </div>
                <div class="mt-1.5 h-2 overflow-hidden rounded-full bg-sem-votos">
                  <div class="h-full rounded-full bg-tinta transition-[width] duration-700" :style="{ width: `${escopo.percApuradas}%` }" />
                </div>
                <p v-if="ufSelecionada && escopo.totalizado" class="mt-1.5 text-xs text-tinta-muted">
                  Totalizado em {{ dataHora(escopo.totalizado).data }} às {{ dataHora(escopo.totalizado).hora }}
                </p>
              </div>

              <p
                v-if="escopo.andamento === 'n'"
                class="rounded-xl bg-realce px-4 py-3 text-sm text-tinta-dim"
              >
                A totalização ainda não começou aqui.
              </p>
              <PlacarCandidatos :resultado="escopo" :cores="cores" />
              <ResumoVotos :resultado="escopo" />
            </section>

            <TabelaUFs :ufs="ufs" :cores="cores" :selecionada="ufSelecionada" @selecionar="ufSelecionada = $event" />
          </div>

          <section class="cartao relative h-[560px] overflow-hidden lg:sticky lg:top-4 lg:h-[calc(100vh-2.5rem-2rem)]">
            <MapaBrasil :ufs="ufs" :cores="cores" :selecionada="ufSelecionada" @selecionar="ufSelecionada = $event" />
            <div class="pointer-events-none absolute left-4 top-4">
              <h2 class="text-base font-semibold">Líder por estado</h2>
              <p class="text-xs text-tinta-muted">Passe o mouse para detalhes</p>
            </div>
            <div class="pointer-events-none absolute bottom-4 left-4 space-y-1 rounded-xl border border-borda bg-white/90 px-3 py-2.5 text-xs">
              <div v-for="c in brasil.candidatos" :key="c.numero" class="flex items-center gap-2">
                <span class="h-2.5 w-2.5 rounded-sm" :style="{ background: cores[c.numero] }" />
                {{ c.nome }}
              </div>
              <div class="flex items-center gap-2">
                <span class="h-2.5 w-2.5 rounded-sm bg-sem-votos" />
                Sem votos apurados
              </div>
              <p class="pt-1 text-tinta-muted">Tom mais forte = vantagem maior</p>
            </div>
          </section>
        </div>
      </div>
    </main>
  </div>
</template>
