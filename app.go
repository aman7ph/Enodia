package main

import (
	"context"
	"enodia/internal/apps"
	"enodia/internal/firewall"
)

// App struct holds application state
type App struct {
	ctx           context.Context
	fw            *firewall.Manager
	installedApps []apps.InstalledApp
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{}
}

// startup is called when the app launches
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.fw = firewall.NewManager()
	a.installedApps = apps.DiscoverApps()
}

// shutdown is called when the app closes
func (a *App) shutdown(ctx context.Context) {
	if a.fw != nil {
		a.fw.Close()
	}
}
