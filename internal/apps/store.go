package apps

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// discoverStoreApps finds Microsoft Store / MSIX apps using PowerShell
func discoverStoreApps() []InstalledApp {
	var apps []InstalledApp

	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		`Get-AppxPackage | Where-Object { $_.IsFramework -eq $false } | Select-Object Name, Publisher, InstallLocation, PackageFamilyName | ConvertTo-Json`)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("[Enodia] Warning: Could not get Store apps: %v", err)
		return apps
	}

	var storeApps []struct {
		Name              string `json:"Name"`
		Publisher         string `json:"Publisher"`
		InstallLocation   string `json:"InstallLocation"`
		PackageFamilyName string `json:"PackageFamilyName"`
	}

	if err := json.Unmarshal(output, &storeApps); err != nil {
		var single struct {
			Name              string `json:"Name"`
			Publisher         string `json:"Publisher"`
			InstallLocation   string `json:"InstallLocation"`
			PackageFamilyName string `json:"PackageFamilyName"`
		}
		if err := json.Unmarshal(output, &single); err != nil {
			log.Printf("[Enodia] Warning: Could not parse Store apps: %v", err)
			return apps
		}
		storeApps = append(storeApps, single)
	}

	for _, sa := range storeApps {
		if isSystemStoreApp(sa.Name, sa.Publisher) {
			continue
		}

		app := InstalledApp{
			ID:                generateID(sa.Name + sa.InstallLocation),
			Name:              cleanStoreName(sa.Name),
			Publisher:         cleanPublisher(sa.Publisher),
			InstallPath:       sa.InstallLocation,
			AppType:           "store",
			PackageFamilyName: sa.PackageFamilyName,
			PackageSID:        getPackageSID(sa.PackageFamilyName),
		}

		if sa.InstallLocation != "" {
			app.Executables = findExecutables(sa.InstallLocation)
			app.IconBase64 = extractIconBase64(sa.InstallLocation)
		}

		apps = append(apps, app)
	}

	log.Printf("[Enodia] Found %d Store apps", len(apps))
	return apps
}

// getPackageSID gets the App Container SID for a UWP app from the registry
func getPackageSID(packageFamilyName string) string {
	if packageFamilyName == "" {
		return ""
	}

	mappingsPath := `Software\Classes\Local Settings\Software\Microsoft\Windows\CurrentVersion\AppContainer\Mappings`
	key, err := registry.OpenKey(registry.CURRENT_USER, mappingsPath, registry.READ)
	if err != nil {
		return ""
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return ""
	}

	for _, sid := range subkeys {
		if !strings.HasPrefix(sid, "S-1-15-2-") {
			continue
		}

		subkey, err := registry.OpenKey(key, sid, registry.READ)
		if err != nil {
			continue
		}

		moniker, _, err := subkey.GetStringValue("Moniker")
		subkey.Close()

		if err == nil && (strings.EqualFold(moniker, packageFamilyName) ||
			strings.Contains(strings.ToLower(moniker), strings.ToLower(strings.Split(packageFamilyName, "_")[0]))) {
			return sid
		}
	}
	return ""
}
