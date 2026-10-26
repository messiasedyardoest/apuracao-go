<script setup lang="ts">
import "maplibre-gl/dist/maplibre-gl.css";
import maplibregl from "maplibre-gl";
import type { ExpressionSpecification, LngLatBoundsLike, MapLayerMouseEvent, StyleSpecification } from "maplibre-gl";

// Mapa de UFs colorido pelo lider de cada estado (tom mais forte = vantagem
// maior). Sem mapa base externo: so a malha IBGE de public/geo/br-ufs.json
// (gerada com sigla, nome e centroide oficial), entao nao depende de tiles.
const props = defineProps<{
  ufs: Record<string, ResultadoApuracao>;
  cores: Record<string, string>;
  selecionada: string | null;
}>();
const emit = defineEmits<{ selecionar: [uf: string | null] }>();

type PropsUF = { codarea: string; sigla: string; nome: string; centro: [number, number] };
type Colecao = GeoJSON.FeatureCollection<GeoJSON.Polygon | GeoJSON.MultiPolygon, PropsUF>;

const LIMITES_BR: LngLatBoundsLike = [[-74.1, -33.8], [-34.7, 5.3]];
const SEM_VOTOS = "#ddd6cf";
const ESTILO: StyleSpecification = {
  version: 8,
  sources: {},
  layers: [{ id: "fundo", type: "background", paint: { "background-color": "#f7f4f0" } }],
};
// 50% dos validos = tom claro; 75%+ = cor cheia.
const OPACIDADE: ExpressionSpecification = [
  "case",
  ["boolean", ["feature-state", "temVotos"], false],
  ["interpolate", ["linear"], ["number", ["feature-state", "forca"], 0], 0, 0.4, 1, 0.95],
  1,
];

const { fmtNumero, fmtPct } = useFormatadores();
const alvo = ref<HTMLDivElement | null>(null);
const erro = ref("");
let mapa: maplibregl.Map | null = null;
let carregado = false;
let hover: string | undefined;
const limites = new Map<string, LngLatBoundsLike>();
const dica = new maplibregl.Popup({ closeButton: false, closeOnClick: false, offset: 12, maxWidth: "380px" });

function limitesDe(coordenadas: unknown, caixa = [Infinity, Infinity, -Infinity, -Infinity]): number[] {
  if (Array.isArray(coordenadas) && typeof coordenadas[0] === "number") {
    const [lng, lat] = coordenadas as [number, number];
    return [Math.min(caixa[0]!, lng), Math.min(caixa[1]!, lat), Math.max(caixa[2]!, lng), Math.max(caixa[3]!, lat)];
  }
  for (const parte of coordenadas as unknown[]) caixa = limitesDe(parte, caixa);
  return caixa;
}

function escapar(texto: string) {
  return texto.replace(/[&<>"']/g, (c) => `&#${c.charCodeAt(0)};`);
}

onMounted(async () => {
  try {
    const geo = await $fetch<Colecao>("/geo/br-ufs.json");
    for (const f of geo.features) {
      const c = limitesDe(f.geometry.coordinates);
      limites.set(f.properties.sigla, [[c[0]!, c[1]!], [c[2]!, c[3]!]]);
    }
    montarMapa(geo);
  } catch (e) {
    erro.value = e instanceof Error ? e.message : "Falha ao carregar o mapa";
  }
});

onBeforeUnmount(() => {
  dica.remove();
  mapa?.remove();
  mapa = null;
});

function montarMapa(geo: Colecao) {
  if (!alvo.value) return;
  mapa = new maplibregl.Map({
    container: alvo.value,
    style: ESTILO,
    bounds: LIMITES_BR,
    fitBoundsOptions: { padding: 24 },
    dragRotate: false,
    pitchWithRotate: false,
    touchPitch: false,
    renderWorldCopies: false,
    attributionControl: false,
    minZoom: 2.5,
    maxZoom: 8,
  });
  mapa.touchZoomRotate.disableRotation();
  mapa.addControl(new maplibregl.NavigationControl({ showCompass: false }), "top-right");

  mapa.on("load", () => {
    if (!mapa) return;
    mapa.addSource("ufs", { type: "geojson", data: geo, promoteId: "sigla" });
    mapa.addLayer({
      id: "ufs-preenchimento",
      type: "fill",
      source: "ufs",
      paint: {
        "fill-color": ["to-color", ["coalesce", ["feature-state", "cor"], SEM_VOTOS]],
        "fill-opacity": OPACIDADE,
      },
    });
    mapa.addLayer({
      id: "ufs-borda",
      type: "line",
      source: "ufs",
      paint: { "line-color": "#ffffff", "line-width": 1 },
    });
    // Contorno de hover/selecao numa camada propria, por cima das bordas
    // brancas dos vizinhos.
    mapa.addLayer({
      id: "ufs-destaque",
      type: "line",
      source: "ufs",
      paint: {
        "line-color": "#1a0e0c",
        "line-width": ["case", ["boolean", ["feature-state", "selecionada"], false], 2.5, 1.5],
        "line-opacity": [
          "case",
          ["any", ["boolean", ["feature-state", "selecionada"], false], ["boolean", ["feature-state", "hover"], false]],
          1,
          0,
        ],
      },
    });

    // Rotulos em HTML: o estilo nao tem glyphs (fontes) para camada symbol.
    for (const f of geo.features) {
      const rotulo = document.createElement("span");
      rotulo.textContent = f.properties.sigla.toUpperCase();
      rotulo.style.cssText = "pointer-events:none;font:600 11px Inter,system-ui,sans-serif;color:#fff;text-shadow:0 0 3px rgba(26,14,12,.55)";
      new maplibregl.Marker({ element: rotulo }).setLngLat(f.properties.centro).addTo(mapa);
    }

    mapa.on("mousemove", "ufs-preenchimento", (e: MapLayerMouseEvent) => {
      const sigla = e.features?.[0]?.properties?.sigla as string | undefined;
      if (!mapa || !sigla) return;
      if (hover && hover !== sigla) mapa.setFeatureState({ source: "ufs", id: hover }, { hover: false });
      hover = sigla;
      mapa.setFeatureState({ source: "ufs", id: sigla }, { hover: true });
      mapa.getCanvas().style.cursor = "pointer";
      dica.setLngLat(e.lngLat).setHTML(conteudoDica(sigla)).addTo(mapa);
    });
    mapa.on("mouseleave", "ufs-preenchimento", () => {
      if (!mapa) return;
      if (hover) mapa.setFeatureState({ source: "ufs", id: hover }, { hover: false });
      hover = undefined;
      mapa.getCanvas().style.cursor = "";
      dica.remove();
    });
    mapa.on("click", (e) => {
      const sigla = mapa?.queryRenderedFeatures(e.point, { layers: ["ufs-preenchimento"] })[0]?.properties?.sigla as string | undefined;
      emit("selecionar", sigla && sigla !== props.selecionada ? sigla : null);
    });

    carregado = true;
    pintar();
    marcarSelecao(props.selecionada, null);
  });
}

function conteudoDica(sigla: string) {
  const r = props.ufs[sigla];
  const nome = escapar(NOMES_UF[sigla] ?? sigla.toUpperCase());
  if (!r) return `<strong>${nome}</strong>`;
  const temVotos = r.candidatos.some((c) => c.votos > 0);
  // Grade de 3 colunas (nome | votos | %) para os numeros alinharem.
  const linhas = temVotos
    ? `<div style="display:grid;grid-template-columns:1fr auto auto;column-gap:14px;row-gap:4px;align-items:baseline">${r.candidatos
        .map(
          (c) =>
            `<span style="white-space:nowrap"><span style="display:inline-block;width:8px;height:8px;border-radius:2px;margin-right:6px;background:${props.cores[c.numero] ?? SEM_VOTOS}"></span>${escapar(c.nome)}</span>` +
            `<span style="text-align:right;color:#6b635c;white-space:nowrap">${fmtNumero(c.votos)} votos</span>` +
            `<strong style="text-align:right">${fmtPct(c.percentual)}</strong>`,
        )
        .join("")}</div>` +
      `<div style="color:#6b635c;margin-top:8px;padding-top:6px;border-top:1px solid rgba(26,14,12,.08);white-space:nowrap">Válidos ${fmtNumero(r.validos)} · Brancos ${fmtNumero(r.brancos)} · Nulos ${fmtNumero(r.nulos)}</div>`
    : `<div style="color:#a39c95">Sem votos apurados</div>`;
  return `<div style="font-size:12px;color:#1a0e0c;font-variant-numeric:tabular-nums;min-width:220px"><div style="font-weight:600;font-size:13px">${nome}</div><div style="color:#6b635c;margin:2px 0 8px;white-space:nowrap">${fmtPct(r.percApuradas)} das urnas apuradas · ${fmtNumero(r.urnasApuradas)} de ${fmtNumero(r.urnasTotal)}</div>${linhas}</div>`;
}

function pintar() {
  if (!mapa || !carregado) return;
  for (const sigla of limites.keys()) {
    const lider = props.ufs[sigla]?.candidatos[0];
    const temVotos = (lider?.votos ?? 0) > 0;
    mapa.setFeatureState(
      { source: "ufs", id: sigla },
      {
        temVotos,
        cor: temVotos ? (props.cores[lider!.numero] ?? SEM_VOTOS) : SEM_VOTOS,
        forca: temVotos ? Math.min(1, Math.max(0, (lider!.percentual - 50) / 25)) : 0,
      },
    );
  }
  // Chegou dado novo com o mouse parado sobre uma UF: atualiza a dica aberta.
  if (hover && dica.isOpen()) dica.setHTML(conteudoDica(hover));
}

function marcarSelecao(nova: string | null, antiga: string | null) {
  if (!mapa || !carregado) return;
  if (antiga && limites.has(antiga)) mapa.setFeatureState({ source: "ufs", id: antiga }, { selecionada: false });
  if (nova && limites.has(nova)) {
    mapa.setFeatureState({ source: "ufs", id: nova }, { selecionada: true });
    mapa.fitBounds(limites.get(nova)!, { padding: 60, duration: 700, maxZoom: 6.5 });
  } else if (antiga) {
    mapa.fitBounds(LIMITES_BR, { padding: 24, duration: 700 });
  }
}

watch(() => [props.ufs, props.cores], pintar);
watch(() => props.selecionada, (nova, antiga) => marcarSelecao(nova, antiga ?? null));
</script>

<template>
  <div class="relative h-full w-full">
    <div ref="alvo" class="h-full w-full" />
    <div v-if="erro" class="absolute inset-0 grid place-items-center p-6 text-center text-sm text-tinta-dim">
      Não foi possível carregar o mapa: {{ erro }}
    </div>
  </div>
</template>
