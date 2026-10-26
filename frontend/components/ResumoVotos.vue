<script setup lang="ts">
const props = defineProps<{ resultado: ResultadoApuracao }>();

const { fmtNumero, fmtPct } = useFormatadores();
const itens = computed(() => {
  const r = props.resultado;
  return [
    { rotulo: "Eleitorado", valor: r.eleitorado },
    { rotulo: "Comparecimento", valor: r.comparecimento, pct: r.percComparecimento },
    { rotulo: "Abstenção", valor: r.abstencao, pct: r.percAbstencao },
    { rotulo: "Válidos", valor: r.validos, pct: r.percValidos },
    { rotulo: "Brancos", valor: r.brancos, pct: r.percBrancos },
    { rotulo: "Nulos", valor: r.nulos, pct: r.percNulos },
  ];
});
</script>

<template>
  <dl class="grid grid-cols-3 gap-px overflow-hidden rounded-xl border border-borda bg-borda">
    <div v-for="item in itens" :key="item.rotulo" class="bg-cartao px-3 py-2.5">
      <dt class="text-xs text-tinta-dim">{{ item.rotulo }}</dt>
      <dd class="mt-0.5 text-sm font-semibold tabular-nums">
        {{ fmtNumero(item.valor) }}
        <span v-if="item.pct !== undefined" class="ml-1 text-xs font-normal text-tinta-muted">{{ fmtPct(item.pct) }}</span>
      </dd>
    </div>
  </dl>
</template>
