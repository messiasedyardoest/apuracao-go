export default defineNuxtConfig({
  compatibilityDate: "2026-04-01",
  // App desktop (Wails): SPA estatica embarcada no binario Go, sem servidor Node.
  // `nuxt generate` -> .output/public (ver //go:embed em main.go).
  ssr: false,
  devtools: { enabled: false },
  css: ["~/assets/css/main.css"],
  modules: ["@nuxtjs/tailwindcss"],
  sourcemap: { server: false, client: false },
  experimental: { appManifest: false },
  vite: {
    build: { reportCompressedSize: false },
    optimizeDeps: { include: ["maplibre-gl", "@tabler/icons-vue"] },
  },
  app: {
    head: {
      title: "Apuração 2026",
      htmlAttrs: { lang: "pt-BR" },
      link: [
        { rel: "preconnect", href: "https://fonts.googleapis.com" },
        { rel: "preconnect", href: "https://fonts.gstatic.com", crossorigin: "" },
        { rel: "stylesheet", href: "https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" },
      ],
    },
  },
});
