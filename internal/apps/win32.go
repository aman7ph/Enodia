package apps

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

// discoverWin32Apps reads installed apps from the Windows Registry
func discoverWin32Apps() []InstalledApp {
	var apps []InstalledApp

	regPaths := []struct {
		root registry.Key
		path string
	}{
		{registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
	}

	for _, regPath := range regPaths {
		key, err := registry.OpenKey(regPath.root, regPath.path, registry.READ)
		if err != nil {
			continue
		}

		subkeys, err := key.ReadSubKeyNames(-1)
		if err != nil {
			key.Close()
			continue
		}

		for _, subkeyName := range subkeys {
			subkey, err := registry.OpenKey(key, subkeyName, registry.READ)
			if err != nil {
				continue
			}

			name, _, _ := subkey.GetStringValue("DisplayName")
			publisher, _, _ := subkey.GetStringValue("Publisher")
			installPath, _, _ := subkey.GetStringValue("InstallLocation")
			iconPath, _, _ := subkey.GetStringValue("DisplayIcon")
			subkey.Close()

			if name == "" || isSystemApp(name, publisher, installPath) {
				continue
			}

			app := InstalledApp{
				ID:          generateID(name + installPath),
				Name:        name,
				Publisher:   publisher,
				InstallPath: installPath,
				AppType:     "win32",
			}

			if installPath != "" {
				app.Executables = findExecutables(installPath)
			}
			if iconPath != "" {
				app.IconBase64 = extractIconBase64(iconPath)
			}

			apps = append(apps, app)
		}
		key.Close()
	}

	log.Printf("[Enodia] Found %d Win32 apps", len(apps))
	return apps
}
