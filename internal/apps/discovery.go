package apps

import "log"

// DiscoverApps finds all installed applications (Win32 + Store)
func DiscoverApps() []InstalledApp {
	log.Println("[Enodia] Starting app discovery...")

	apps := discoverWin32Apps()
	apps = append(apps, discoverStoreApps()...)

	log.Printf("[Enodia] Discovered %d applications total", len(apps))
	return apps
}
