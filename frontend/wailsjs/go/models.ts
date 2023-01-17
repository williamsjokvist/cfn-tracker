export namespace main {
	
	export class MatchHistory {
	    cfn: string;
	    lp: number;
	    lpGain: number;
	    wins: number;
	    totalWins: number;
	    totalLosses: number;
	    totalMatches: number;
	    losses: number;
	    winRate: number;
	
	    static createFrom(source: any = {}) {
	        return new MatchHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cfn = source["cfn"];
	        this.lp = source["lp"];
	        this.lpGain = source["lpGain"];
	        this.wins = source["wins"];
	        this.totalWins = source["totalWins"];
	        this.totalLosses = source["totalLosses"];
	        this.totalMatches = source["totalMatches"];
	        this.losses = source["losses"];
	        this.winRate = source["winRate"];
	    }
	}

}

