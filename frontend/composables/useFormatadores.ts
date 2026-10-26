const numero = new Intl.NumberFormat("pt-BR");
const percentual = new Intl.NumberFormat("pt-BR", { minimumFractionDigits: 2, maximumFractionDigits: 2 });

export function useFormatadores() {
  return {
    fmtNumero: (n: number) => numero.format(n),
    fmtPct: (n: number) => `${percentual.format(n)}%`,
    fmtPontos: (n: number) => `${percentual.format(n)} p.p.`,
    fmtHora: (iso: string) => new Date(iso).toLocaleTimeString("pt-BR"),
  };
}
