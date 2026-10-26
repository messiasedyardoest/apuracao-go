import { Atualizar, Painel } from "~~/wailsjs/go/main/App";
import type { main } from "~~/wailsjs/go/models";
import { EventsOn } from "~~/wailsjs/runtime/runtime";

export type PainelApuracao = main.Painel;
export type ResultadoApuracao = main.Resultado;
export type CandidatoApuracao = main.Candidato;

// O back (app.go) consulta o TSE a cada 30 s e emite este evento; aqui so
// recebemos e contamos o tempo ate a proxima consulta.
const EVENTO_PAINEL = "apuracao:painel";

export function useApuracao() {
  const painel = shallowRef<PainelApuracao | null>(null);
  const agora = ref(Date.now());
  const atualizando = ref(false);
  const semWails = ref(false);
  let pararEvento: (() => void) | undefined;
  let relogio: ReturnType<typeof setInterval> | undefined;

  function receber(novo: PainelApuracao) {
    painel.value = novo;
    atualizando.value = false;
  }

  onMounted(async () => {
    relogio = setInterval(() => (agora.value = Date.now()), 1000);
    // Fora do Wails (navegador puro em localhost:3000) nao ha backend Go.
    if (!("go" in window)) {
      semWails.value = true;
      return;
    }
    pararEvento = EventsOn(EVENTO_PAINEL, receber);
    const atual = await Painel();
    // Antes da 1a consulta o back devolve o painel vazio (sem consultadoEm).
    if (atual.consultadoEm && !painel.value) receber(atual);
  });

  onBeforeUnmount(() => {
    pararEvento?.();
    clearInterval(relogio);
  });

  const segundosRestantes = computed(() => {
    const proxima = painel.value?.proximaEm;
    if (!proxima) return null;
    return Math.max(0, Math.ceil((Date.parse(proxima) - agora.value) / 1000));
  });

  function atualizar() {
    if (semWails.value || atualizando.value) return;
    atualizando.value = true;
    Atualizar();
  }

  return { painel, segundosRestantes, atualizando, atualizar, semWails };
}
