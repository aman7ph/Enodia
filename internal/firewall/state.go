package firewall

import (
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// GetBlockedApps returns a list of applications with their block status
func (m *Manager) GetBlockedApps() ([]BlockedApp, error) {
	resultChan := make(chan []BlockedApp, 1)
	errChan := make(chan error, 1)

	m.jobs <- func(rules *ole.IDispatch) {
		appMap := make(map[string]*BlockedApp)

		err := oleutil.ForEach(rules, func(v *ole.VARIANT) error {
			ruleDispatch := v.ToIDispatch()
			defer ruleDispatch.Release()

			nameVar, _ := oleutil.GetProperty(ruleDispatch, "Name")
			name := nameVar.ToString()

			if !strings.HasPrefix(name, RULE_PREFIX) {
				return nil
			}

			appPathVar, _ := oleutil.GetProperty(ruleDispatch, "ApplicationName")
			enabledVar, _ := oleutil.GetProperty(ruleDispatch, "Enabled")

			appPath := appPathVar.ToString()
			enabled := enabledVar.Val != 0

			// For PKG rules, extract the display name from the rule name
			if strings.Contains(name, "PKG-") {
				// Rule name format: "Enodia-OUT-PKG-AppName" or "Enodia-IN-PKG-AppName"
				parts := strings.SplitN(name, "PKG-", 2)
				if len(parts) == 2 {
					appPath = "PKG-" + parts[1]
				}
			}

			if appPath == "" {
				return nil
			}

			app, exists := appMap[appPath]
			if !exists {
				app = &BlockedApp{
					AppPath:     appPath,
					DisplayName: extractDisplayName(appPath),
				}
				appMap[appPath] = app
			}

			if strings.HasPrefix(name, RULE_PREFIX_OUT) {
				app.OutboundBlocked = enabled
			} else if strings.HasPrefix(name, RULE_PREFIX_IN) {
				app.InboundBlocked = enabled
			}

			return nil
		})

		if err != nil {
			errChan <- err
			return
		}

		result := make([]BlockedApp, 0, len(appMap))
		for _, app := range appMap {
			result = append(result, *app)
		}
		resultChan <- result
	}

	select {
	case res := <-resultChan:
		return res, nil
	case err := <-errChan:
		return nil, err
	}
}

// extractDisplayName gets a user-friendly name from the exe path
func extractDisplayName(exePath string) string {
	base := filepath.Base(exePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return strings.Title(strings.ToLower(name))
}
