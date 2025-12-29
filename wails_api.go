package main

import (
	"enodia/internal/apps"
	"enodia/internal/firewall"
	"fmt"
)

// GetInstalledApps returns all discovered applications
func (a *App) GetInstalledApps() []apps.InstalledApp {
	return a.installedApps
}

// RefreshApps re-discovers applications
func (a *App) RefreshApps() []apps.InstalledApp {
	a.installedApps = apps.DiscoverApps()
	return a.installedApps
}

// BlockFile blocks a single executable
func (a *App) BlockFile(path string) string {
	if a.fw == nil {
		return "Error: Firewall not available"
	}
	if err := a.fw.BlockApp(path); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return "Blocked"
}

// UnblockFile removes block for a single executable
func (a *App) UnblockFile(path string) string {
	if a.fw == nil {
		return "Error: Firewall not available"
	}
	if err := a.fw.UnblockApp(path); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return "Unblocked"
}

// BlockFiles blocks multiple executables
func (a *App) BlockFiles(paths []string) map[string]string {
	result := make(map[string]string)
	if a.fw == nil {
		for _, p := range paths {
			result[p] = "Error: Firewall not available"
		}
		return result
	}
	for path, err := range a.fw.BlockApps(paths) {
		if err != nil {
			result[path] = fmt.Sprintf("Error: %v", err)
		} else {
			result[path] = "Blocked"
		}
	}
	return result
}

// UnblockFiles removes blocks for multiple executables
func (a *App) UnblockFiles(paths []string) map[string]string {
	result := make(map[string]string)
	if a.fw == nil {
		for _, p := range paths {
			result[p] = "Error: Firewall not available"
		}
		return result
	}
	for path, err := range a.fw.UnblockApps(paths) {
		if err != nil {
			result[path] = fmt.Sprintf("Error: %v", err)
		} else {
			result[path] = "Unblocked"
		}
	}
	return result
}

// BlockInstalledApp blocks an app based on its type (Win32 or Store)
func (a *App) BlockInstalledApp(app apps.InstalledApp) string {
	if a.fw == nil {
		return "Error: Firewall not available"
	}
	if app.AppType == "store" && app.PackageSID != "" {
		if err := a.fw.BlockStoreApp(app.PackageSID, app.Name); err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return "Blocked"
	}
	if len(app.Executables) == 0 {
		return "No executables to block"
	}
	for _, err := range a.fw.BlockApps(app.Executables) {
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
	}
	return "Blocked"
}

// UnblockInstalledApp unblocks an app based on its type
func (a *App) UnblockInstalledApp(app apps.InstalledApp) string {
	if a.fw == nil {
		return "Error: Firewall not available"
	}
	if app.AppType == "store" {
		if err := a.fw.UnblockStoreApp(app.Name); err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return "Unblocked"
	}
	if len(app.Executables) == 0 {
		return "No executables to unblock"
	}
	for _, err := range a.fw.UnblockApps(app.Executables) {
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
	}
	return "Unblocked"
}

// GetBlockedApps returns all apps with Enodia firewall rules
func (a *App) GetBlockedApps() []firewall.BlockedApp {
	if a.fw == nil {
		return []firewall.BlockedApp{}
	}
	blocked, _ := a.fw.GetBlockedApps()
	return blocked
}
