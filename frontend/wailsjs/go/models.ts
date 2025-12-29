export namespace apps {
	
	export class InstalledApp {
	    id: string;
	    name: string;
	    publisher: string;
	    installPath: string;
	    executables: string[];
	    iconBase64: string;
	    appType: string;
	    packageFamilyName: string;
	    packageSID: string;
	
	    static createFrom(source: any = {}) {
	        return new InstalledApp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.publisher = source["publisher"];
	        this.installPath = source["installPath"];
	        this.executables = source["executables"];
	        this.iconBase64 = source["iconBase64"];
	        this.appType = source["appType"];
	        this.packageFamilyName = source["packageFamilyName"];
	        this.packageSID = source["packageSID"];
	    }
	}

}

export namespace firewall {
	
	export class BlockedApp {
	    appPath: string;
	    displayName: string;
	    inboundBlocked: boolean;
	    outboundBlocked: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BlockedApp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appPath = source["appPath"];
	        this.displayName = source["displayName"];
	        this.inboundBlocked = source["inboundBlocked"];
	        this.outboundBlocked = source["outboundBlocked"];
	    }
	}

}

