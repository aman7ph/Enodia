export interface InstalledApp {
  id: string;
  name: string;
  publisher: string;
  installPath: string;
  executables: string[];
  iconBase64: string;
  appType: string;
  packageFamilyName: string;
  packageSid: string;
}

export interface BlockedApp {
  appPath: string;
  displayName: string;
  inboundBlocked: boolean;
  outboundBlocked: boolean;
}
