package firewall

import (
	"fmt"
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// createBlockRule creates a single firewall rule for a Win32 app
func createBlockRule(rules *ole.IDispatch, exePath string, direction int) error {
	prefix := RULE_PREFIX_OUT
	dirName := "OUT"
	if direction == NET_FW_RULE_DIR_IN {
		prefix = RULE_PREFIX_IN
		dirName = "IN"
	}

	ruleName := prefix + exePath
	_, _ = oleutil.CallMethod(rules, "Remove", ruleName)

	unknownRule, err := oleutil.CreateObject(PROGID_RULE)
	if err != nil {
		return fmt.Errorf("failed to create rule: %w", err)
	}
	defer unknownRule.Release()

	rule, err := unknownRule.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to query interface: %w", err)
	}
	defer rule.Release()

	oleutil.PutProperty(rule, "Name", ruleName)
	oleutil.PutProperty(rule, "Description", "Blocked by Enodia")
	oleutil.PutProperty(rule, "ApplicationName", exePath)
	oleutil.PutProperty(rule, "Protocol", int32(NET_FW_IP_PROTOCOL_ANY))
	oleutil.PutProperty(rule, "Direction", int32(direction))
	oleutil.PutProperty(rule, "Action", int32(NET_FW_ACTION_BLOCK))
	oleutil.PutProperty(rule, "Profiles", int32(NET_FW_PROFILE2_ALL))
	oleutil.PutProperty(rule, "Enabled", true)

	if _, err := oleutil.CallMethod(rules, "Add", rule); err != nil {
		return fmt.Errorf("failed to add rule: %w", err)
	}

	log.Printf("[Enodia] Created %s rule: %s", dirName, exePath)
	return nil
}

// createBlockRuleForPackage creates a firewall rule for UWP/Store apps
func createBlockRuleForPackage(rules *ole.IDispatch, packageSID, displayName string, direction int) error {
	prefix := RULE_PREFIX_OUT + "PKG-"
	dirName := "OUT"
	if direction == NET_FW_RULE_DIR_IN {
		prefix = RULE_PREFIX_IN + "PKG-"
		dirName = "IN"
	}

	ruleName := prefix + displayName
	_, _ = oleutil.CallMethod(rules, "Remove", ruleName)

	unknownRule, err := oleutil.CreateObject(PROGID_RULE)
	if err != nil {
		return fmt.Errorf("failed to create rule: %w", err)
	}
	defer unknownRule.Release()

	rule, err := unknownRule.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to query interface: %w", err)
	}
	defer rule.Release()

	oleutil.PutProperty(rule, "Name", ruleName)
	oleutil.PutProperty(rule, "Description", "Blocked by Enodia")
	oleutil.PutProperty(rule, "LocalAppPackageId", packageSID)
	oleutil.PutProperty(rule, "Direction", int32(direction))
	oleutil.PutProperty(rule, "Action", int32(NET_FW_ACTION_BLOCK))
	oleutil.PutProperty(rule, "Profiles", int32(NET_FW_PROFILE2_ALL))
	oleutil.PutProperty(rule, "Enabled", true)

	if _, err := oleutil.CallMethod(rules, "Add", rule); err != nil {
		return fmt.Errorf("failed to add UWP rule: %w", err)
	}

	log.Printf("[Enodia] Created %s UWP rule: %s", dirName, displayName)
	return nil
}
