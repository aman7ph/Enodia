package firewall

import (
	"fmt"
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// BlockApp creates block rules for a single Win32 executable
func (m *Manager) BlockApp(exePath string) error {
	errChan := make(chan error, 1)
	m.jobs <- func(rules *ole.IDispatch) {
		if err := createBlockRule(rules, exePath, NET_FW_RULE_DIR_OUT); err != nil {
			errChan <- err
			return
		}
		if err := createBlockRule(rules, exePath, NET_FW_RULE_DIR_IN); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}
	return <-errChan
}

// UnblockApp removes block rules for a Win32 executable
func (m *Manager) UnblockApp(exePath string) error {
	errChan := make(chan error, 1)
	m.jobs <- func(rules *ole.IDispatch) {
		oleutil.CallMethod(rules, "Remove", RULE_PREFIX_OUT+exePath)
		oleutil.CallMethod(rules, "Remove", RULE_PREFIX_IN+exePath)
		errChan <- nil
	}
	return <-errChan
}

// BlockStoreApp creates block rules for a UWP/Store app
func (m *Manager) BlockStoreApp(packageSID, displayName string) error {
	errChan := make(chan error, 1)
	m.jobs <- func(rules *ole.IDispatch) {
		if err := createBlockRuleForPackage(rules, packageSID, displayName, NET_FW_RULE_DIR_OUT); err != nil {
			errChan <- err
			return
		}
		if err := createBlockRuleForPackage(rules, packageSID, displayName, NET_FW_RULE_DIR_IN); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}
	return <-errChan
}

// UnblockStoreApp removes block rules for a UWP/Store app
func (m *Manager) UnblockStoreApp(displayName string) error {
	errChan := make(chan error, 1)
	m.jobs <- func(rules *ole.IDispatch) {
		oleutil.CallMethod(rules, "Remove", RULE_PREFIX_OUT+"PKG-"+displayName)
		oleutil.CallMethod(rules, "Remove", RULE_PREFIX_IN+"PKG-"+displayName)
		errChan <- nil
	}
	return <-errChan
}

// BlockApps blocks multiple applications in batch
func (m *Manager) BlockApps(exePaths []string) map[string]error {
	resultChan := make(chan map[string]error, 1)
	m.jobs <- func(rules *ole.IDispatch) {
		results := make(map[string]error)
		for _, path := range exePaths {
			var err error
			if e := createBlockRule(rules, path, NET_FW_RULE_DIR_OUT); e != nil {
				err = fmt.Errorf("outbound: %w", e)
			}
			if e := createBlockRule(rules, path, NET_FW_RULE_DIR_IN); e != nil {
				if err != nil {
					err = fmt.Errorf("%v; inbound: %w", err, e)
				} else {
					err = fmt.Errorf("inbound: %w", e)
				}
			}
			results[path] = err
			if err == nil {
				log.Printf("[Enodia] Blocked: %s", path)
			}
		}
		resultChan <- results
	}
	return <-resultChan
}

// UnblockApps removes firewall rules for multiple applications
func (m *Manager) UnblockApps(exePaths []string) map[string]error {
	resultChan := make(chan map[string]error, 1)
	m.jobs <- func(rules *ole.IDispatch) {
		results := make(map[string]error)
		for _, path := range exePaths {
			oleutil.CallMethod(rules, "Remove", RULE_PREFIX_OUT+path)
			oleutil.CallMethod(rules, "Remove", RULE_PREFIX_IN+path)
			results[path] = nil
			log.Printf("[Enodia] Unblocked: %s", path)
		}
		resultChan <- results
	}
	return <-resultChan
}
