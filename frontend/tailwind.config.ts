import type { Config } from "tailwindcss";

export default {
  content: [
    "./app.vue",
    "./components/**/*.{vue,js,ts}",
    "./composables/**/*.{js,ts}",
    "./utils/**/*.{js,ts}",
  ],
  theme: {
    extend: {
      colors: {
        fundo: "var(--bg)",
        cartao: "var(--shell)",
        tinta: "var(--text)",
        "tinta-dim": "var(--text-dim)",
        "tinta-muted": "var(--text-muted)",
        borda: "var(--border)",
        realce: "var(--tag-bg)",
        "sem-votos": "var(--sem-votos)",
      },
    },
  },
  plugins: [],
} satisfies Config;
