export namespace main {
	
	export class Candidato {
	    numero: string;
	    nome: string;
	    nomeCompleto: string;
	    partido: string;
	    vice: string;
	    vicePartido: string;
	    votos: number;
	    percentual: number;
	    eleito: boolean;
	    situacao: string;
	    foto: string;
	
	    static createFrom(source: any = {}) {
	        return new Candidato(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.numero = source["numero"];
	        this.nome = source["nome"];
	        this.nomeCompleto = source["nomeCompleto"];
	        this.partido = source["partido"];
	        this.vice = source["vice"];
	        this.vicePartido = source["vicePartido"];
	        this.votos = source["votos"];
	        this.percentual = source["percentual"];
	        this.eleito = source["eleito"];
	        this.situacao = source["situacao"];
	        this.foto = source["foto"];
	    }
	}
	export class Resultado {
	    uf: string;
	    andamento: string;
	    definido: boolean;
	    totalizado: string;
	    gerado: string;
	    urnasTotal: number;
	    urnasApuradas: number;
	    percApuradas: number;
	    eleitorado: number;
	    comparecimento: number;
	    percComparecimento: number;
	    abstencao: number;
	    percAbstencao: number;
	    validos: number;
	    percValidos: number;
	    brancos: number;
	    percBrancos: number;
	    nulos: number;
	    percNulos: number;
	    candidatos: Candidato[];
	
	    static createFrom(source: any = {}) {
	        return new Resultado(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uf = source["uf"];
	        this.andamento = source["andamento"];
	        this.definido = source["definido"];
	        this.totalizado = source["totalizado"];
	        this.gerado = source["gerado"];
	        this.urnasTotal = source["urnasTotal"];
	        this.urnasApuradas = source["urnasApuradas"];
	        this.percApuradas = source["percApuradas"];
	        this.eleitorado = source["eleitorado"];
	        this.comparecimento = source["comparecimento"];
	        this.percComparecimento = source["percComparecimento"];
	        this.abstencao = source["abstencao"];
	        this.percAbstencao = source["percAbstencao"];
	        this.validos = source["validos"];
	        this.percValidos = source["percValidos"];
	        this.brancos = source["brancos"];
	        this.percBrancos = source["percBrancos"];
	        this.nulos = source["nulos"];
	        this.percNulos = source["percNulos"];
	        this.candidatos = this.convertValues(source["candidatos"], Candidato);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Painel {
	    brasil?: Resultado;
	    ufs: Record<string, Resultado>;
	    consultadoEm: string;
	    proximaEm: string;
	    intervalo: number;
	    falhas: string[];
	    fonte: string;
	
	    static createFrom(source: any = {}) {
	        return new Painel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.brasil = this.convertValues(source["brasil"], Resultado);
	        this.ufs = this.convertValues(source["ufs"], Resultado, true);
	        this.consultadoEm = source["consultadoEm"];
	        this.proximaEm = source["proximaEm"];
	        this.intervalo = source["intervalo"];
	        this.falhas = source["falhas"];
	        this.fonte = source["fonte"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

