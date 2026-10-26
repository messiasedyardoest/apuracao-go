<script setup lang="ts">
const props = defineProps<{
  ufs: Record<string, ResultadoApuracao>;
  cores: Record<string, string>;
  selecionada: string | null;
}>();
const emit = defineEmits<{ selecionar: [uf: string | null] }>();

const { fmtPct, fmtPontos } = useFormatadores();

// Ordem alfabetica pelo nome, exterior no fim.
const linhas = computed(() =>
  Object.values(props.ufs)
    .map((r) => {
      const [lider, segundo] = r.candidatos;
      const temVotos = (lider?.votos ?? 0) > 0;
      return {
        uf: r.uf,
        nome: NOMES_UF[r.uf] ?? r.uf.toUpperCase(),
        resultado: r,
        lider: temVotos ? lider : undefined,
        vantagem: temVotos && segundo ? lider!.percentual - segundo.percentual : undefined,
      };
    })
    .sort((a, b) => (a.uf === "zz" ? 1 : b.uf === "zz" ? -1 : a.nome.localeCompare(b.nome, "pt-BR"))),
);

function alternar(uf: string) {
  emit("selecionar", uf === props.selecionada ? null : uf);
}
</script>

<template>
  <section class="cartao overflow-hidden">
    <header class="px-5 pb-3 pt-4">
      <h2 class="text-base font-semibold">Por estado</h2>
      <p class="text-xs text-tinta-muted">Clique numa linha ou no mapa para ver o placar da UF</p>
    </header>
    <div class="max-h-[460px] overflow-y-auto">
      <table class="w-full text-sm">
        <thead class="sticky top-0 z-10 bg-cartao text-xs text-tinta-muted">
          <tr>
            <th class="px-5 py-2 text-left font-medium">UF</th>
            <th class="px-3 py-2 text-left font-medium">Urnas apuradas</th>
            <th class="px-3 py-2 text-left font-medium">Líder</th>
            <th class="px-5 py-2 text-right font-medium">Vantagem</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="l in linhas"
            :key="l.uf"
            tabindex="0"
            class="cursor-pointer border-t border-borda outline-none transition-colors hover:bg-realce focus-visible:bg-realce"
            :class="{ 'bg-realce': l.uf === selecionada }"
            @click="alternar(l.uf)"
            @keydown.enter="alternar(l.uf)"
          >
            <td class="whitespace-nowrap px-5 py-2">
              <span class="inline-block w-7 font-semibold uppercase">{{ l.uf }}</span>
              <span class="text-tinta-dim">{{ l.nome }}</span>
            </td>
            <td class="px-3 py-2">
              <div class="flex items-center gap-2">
                <div class="h-1.5 w-16 overflow-hidden rounded-full bg-sem-votos">
                  <div class="h-full rounded-full bg-tinta" :style="{ width: `${l.resultado.percApuradas}%` }" />
                </div>
                <span class="tabular-nums text-tinta-dim">{{ fmtPct(l.resultado.percApuradas) }}</span>
              </div>
            </td>
            <td class="whitespace-nowrap px-3 py-2">
              <template v-if="l.lider">
                <span class="mr-1.5 inline-block h-2 w-2 rounded-sm" :style="{ background: cores[l.lider.numero] }" />
                <span>{{ l.lider.nome }}</span>
                <span class="ml-1.5 tabular-nums text-tinta-dim">{{ fmtPct(l.lider.percentual) }}</span>
              </template>
              <span v-else class="text-tinta-muted">sem votos apurados</span>
            </td>
            <td class="whitespace-nowrap px-5 py-2 text-right tabular-nums">
              {{ l.vantagem === undefined ? "—" : fmtPontos(l.vantagem) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
