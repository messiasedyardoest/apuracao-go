<script setup lang="ts">
const props = defineProps<{ resultado: ResultadoApuracao; cores: Record<string, string> }>();

const { fmtNumero, fmtPct } = useFormatadores();
const fotoFalhou = reactive<Record<string, boolean>>({});
const temVotos = computed(() => props.resultado.candidatos.some((c) => c.votos > 0));

function iniciais(nome: string) {
  return nome.split(/\s+/).slice(0, 2).map((p) => p[0]).join("");
}
</script>

<template>
  <div>
    <div class="grid gap-3" :class="resultado.candidatos.length > 2 ? 'grid-cols-3' : 'grid-cols-2'">
      <article
        v-for="(c, i) in resultado.candidatos"
        :key="c.numero"
        class="relative overflow-hidden rounded-2xl border border-borda bg-cartao p-4"
        :style="{ boxShadow: `inset 0 3px 0 ${i === 0 && temVotos ? cores[c.numero] : 'transparent'}` }"
      >
        <div class="flex items-start gap-3">
          <img
            v-if="!fotoFalhou[c.numero]"
            :src="c.foto"
            :alt="c.nome"
            class="h-16 w-12 shrink-0 rounded-lg bg-realce object-cover"
            @error="fotoFalhou[c.numero] = true"
          >
          <div
            v-else
            class="grid h-16 w-12 shrink-0 place-items-center rounded-lg text-sm font-semibold text-white"
            :style="{ background: cores[c.numero] }"
          >
            {{ iniciais(c.nome) }}
          </div>
          <div class="min-w-0">
            <p class="truncate text-sm font-semibold" :title="c.nomeCompleto">{{ c.nome }}</p>
            <p class="text-xs text-tinta-dim">{{ c.numero }} · {{ c.partido }}</p>
            <p v-if="c.vice" class="truncate text-xs text-tinta-muted">Vice: {{ c.vice }}</p>
          </div>
          <span
            v-if="c.eleito"
            class="ml-auto shrink-0 rounded-full px-2 py-0.5 text-xs font-semibold text-white"
            style="background: var(--success)"
          >Eleito</span>
        </div>
        <p class="mt-4 text-4xl font-semibold tabular-nums tracking-tight" :style="{ color: cores[c.numero] }">
          {{ fmtPct(c.percentual) }}
        </p>
        <p class="text-sm tabular-nums text-tinta-dim">{{ fmtNumero(c.votos) }} votos</p>
      </article>
    </div>

    <div class="relative mt-4">
      <div class="flex h-3 overflow-hidden rounded-full bg-sem-votos">
        <div
          v-for="c in resultado.candidatos"
          :key="c.numero"
          class="h-full transition-[width] duration-700"
          :style="{ width: `${c.percentual}%`, background: cores[c.numero] }"
        />
      </div>
      <div v-if="resultado.candidatos.length === 2" class="absolute -top-1 left-1/2 h-5 w-px bg-tinta" />
    </div>
    <p class="mt-1.5 text-xs text-tinta-muted">Percentual sobre os votos válidos{{ resultado.candidatos.length === 2 ? " · linha = 50%" : "" }}</p>
  </div>
</template>
