export namespace main {
	
	export class ExtResult {
	    payload: string;
	    publicKey: string;
	    log: string;
	
	    static createFrom(source: any = {}) {
	        return new ExtResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.payload = source["payload"];
	        this.publicKey = source["publicKey"];
	        this.log = source["log"];
	    }
	}
	export class GenResult {
	    image: string;
	    privKey: string;
	    log: string;
	
	    static createFrom(source: any = {}) {
	        return new GenResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.image = source["image"];
	        this.privKey = source["privKey"];
	        this.log = source["log"];
	    }
	}

}

export namespace stego {
	
	export class AnalysisResult {
	    realData: string;
	    connMap: string;
	    bitGrid: string;
	    dataPath: string;
	
	    static createFrom(source: any = {}) {
	        return new AnalysisResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.realData = source["realData"];
	        this.connMap = source["connMap"];
	        this.bitGrid = source["bitGrid"];
	        this.dataPath = source["dataPath"];
	    }
	}

}

