// window.go e injetado pelo Wails em tempo de execucao. Os bindings gerados
// (wailsjs/go/**/*.js) tem // @ts-check e o acessam direto; sem esta
// declaracao o typecheck falha nesses arquivos, que nao podem ser editados.
// E `any` de proposito: os tipos de verdade estao nos App.d.ts gerados ao lado.
declare global {
  interface Window {
    go: any;
  }
}

export {};
