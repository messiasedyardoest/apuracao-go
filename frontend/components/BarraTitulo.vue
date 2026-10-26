<script setup lang="ts">
import { IconArrowsMaximize, IconArrowsMinimize, IconChartBar, IconCopy, IconMinus, IconSquare, IconX } from "@tabler/icons-vue";
import {
  Quit,
  WindowFullscreen,
  WindowIsFullscreen,
  WindowIsMaximised,
  WindowMinimise,
  WindowToggleMaximise,
  WindowUnfullscreen,
} from "~~/wailsjs/runtime/runtime";

// A janela e Frameless (main.go): esta barra substitui a do Windows. Arrasta
// pela area com --wails-draggable: drag; F11 alterna tela cheia, Esc sai.
const temJanela = ref(false);
const maximizada = ref(false);
const telaCheia = ref(false);

// Maximizar pelo Win+Seta ou entrar em tela cheia tambem redimensiona a
// pagina: o evento resize basta para manter os icones em dia.
async function sincronizar() {
  [maximizada.value, telaCheia.value] = await Promise.all([WindowIsMaximised(), WindowIsFullscreen()]);
}

async function alternarTelaCheia() {
  if (await WindowIsFullscreen()) WindowUnfullscreen();
  else WindowFullscreen();
}

function teclas(e: KeyboardEvent) {
  if (e.key === "F11") {
    e.preventDefault();
    alternarTelaCheia();
  } else if (e.key === "Escape" && telaCheia.value) {
    WindowUnfullscreen();
  }
}

function duploClique() {
  if (!telaCheia.value) WindowToggleMaximise();
}

onMounted(() => {
  // Fora do Wails (navegador puro) nao ha janela nativa para controlar.
  temJanela.value = "runtime" in window;
  if (!temJanela.value) return;
  sincronizar();
  window.addEventListener("resize", sincronizar);
  window.addEventListener("keydown", teclas);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", sincronizar);
  window.removeEventListener("keydown", teclas);
});
</script>

<template>
  <header
    class="flex h-10 shrink-0 select-none items-center border-b border-borda bg-cartao"
    style="--wails-draggable: drag"
    @dblclick="duploClique"
  >
    <div class="flex items-center gap-2 px-4 text-sm">
      <IconChartBar :size="16" style="color: var(--red)" />
      <span class="font-semibold">Apuração 2026</span>
      <span class="text-tinta-muted">Presidente · 2º turno</span>
    </div>

    <div v-if="temJanela" class="ml-auto flex h-full" style="--wails-draggable: no-drag" @dblclick.stop>
      <button
        class="controle-janela"
        :title="telaCheia ? 'Sair da tela cheia (F11)' : 'Tela cheia (F11)'"
        @click="alternarTelaCheia"
      >
        <IconArrowsMinimize v-if="telaCheia" :size="16" />
        <IconArrowsMaximize v-else :size="16" />
      </button>
      <template v-if="!telaCheia">
        <button class="controle-janela" title="Minimizar" @click="WindowMinimise()">
          <IconMinus :size="16" />
        </button>
        <button class="controle-janela" :title="maximizada ? 'Restaurar' : 'Maximizar'" @click="WindowToggleMaximise()">
          <IconCopy v-if="maximizada" :size="14" />
          <IconSquare v-else :size="14" />
        </button>
      </template>
      <button class="controle-janela controle-fechar" title="Fechar" @click="Quit()">
        <IconX :size="16" />
      </button>
    </div>
  </header>
</template>
