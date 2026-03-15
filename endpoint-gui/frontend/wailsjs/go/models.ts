export namespace main {
	
	export class UIScanEvent {
	    current_file?: string;
	    progress_percent?: number;
	    virus_found?: boolean;
	    threat_name?: string;
	
	    static createFrom(source: any = {}) {
	        return new UIScanEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current_file = source["current_file"];
	        this.progress_percent = source["progress_percent"];
	        this.virus_found = source["virus_found"];
	        this.threat_name = source["threat_name"];
	    }
	}

}

