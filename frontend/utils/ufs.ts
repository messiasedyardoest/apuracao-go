export const NOMES_UF: Record<string, string> = {
  br: "Brasil", ac: "Acre", al: "Alagoas", am: "Amazonas", ap: "Amapá", ba: "Bahia",
  ce: "Ceará", df: "Distrito Federal", es: "Espírito Santo", go: "Goiás", ma: "Maranhão",
  mg: "Minas Gerais", ms: "Mato Grosso do Sul", mt: "Mato Grosso", pa: "Pará", pb: "Paraíba",
  pe: "Pernambuco", pi: "Piauí", pr: "Paraná", rj: "Rio de Janeiro", rn: "Rio Grande do Norte",
  ro: "Rondônia", rr: "Roraima", rs: "Rio Grande do Sul", sc: "Santa Catarina", se: "Sergipe",
  sp: "São Paulo", to: "Tocantins", zz: "Exterior",
};

// Cor fixa por numero (nao por posicao): quando a lideranca vira, a cor
// acompanha o candidato.
const CORES_PARTIDO: Record<string, string> = { "13": "#d62828", "22": "#1f4fd1" };
const CORES_RESERVA = ["#7c3aed", "#0e7490", "#c2410c", "#15803d"];

export function coresCandidatos(candidatos: { numero: string }[]): Record<string, string> {
  const numeros = candidatos.map((c) => c.numero).sort();
  return Object.fromEntries(numeros.map((n, i) => [n, CORES_PARTIDO[n] ?? CORES_RESERVA[i % CORES_RESERVA.length]!]));
}
