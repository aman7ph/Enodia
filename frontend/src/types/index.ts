
export interface InstalledApp {
  id: string;          
  name: string;         
  publisher: string;   
  installPath: string;  
  executables: string[];
  iconBase64: string;   
  appType: string;
  packageFamilyName: string;  // For Store/UWP apps
  packageSID: string;         // App Container SID for firewall blocking
}


export interface BlockedApp {
  appPath: string;       
  displayName: string;    
  inboundBlocked: boolean; 
  outboundBlocked: boolean;
}