use winreg::enums::*;
use winreg::RegKey;
use std::collections::HashMap;

/// Opens the registry locations where Windows lists installed applications.
/// There are three separate spots because of 32-bit/64-bit and per-user vs
/// machine-wide installs.
fn find_uninstall_root_keys() -> Vec<RegKey> {
    let mut root_keys = Vec::new();
    let hklm = RegKey::predef(HKEY_LOCAL_MACHINE);
    let hkcu = RegKey::predef(HKEY_CURRENT_USER);

    // 64-bit apps, machine-wide
    if let Ok(key) = hklm.open_subkey_with_flags(
        r"SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall",
        KEY_READ | KEY_WOW64_64KEY,
    ) {
        root_keys.push(key);
    }

    // 32-bit apps, machine-wide
    if let Ok(key) = hklm.open_subkey_with_flags(
        r"SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall",
        KEY_READ | KEY_WOW64_32KEY,
    ) {
        root_keys.push(key);
    }

    // per-user installs
    if let Ok(key) = hkcu.open_subkey(r"SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall") {
        root_keys.push(key);
    }

    root_keys
}

/// Walks every uninstall root key, opens each app's own subkey inside it,
/// and pulls out the display name and icon path where available.
fn collect_installed_apps() -> Vec<HashMap<String, String>> {
    let mut installed_apps = Vec::new();
    let uninstall_root_keys = find_uninstall_root_keys();

    for uninstall_root in uninstall_root_keys.iter() {
        for app_key_name in uninstall_root.enum_keys().flatten() {
            if let Ok(app_registry_key) = uninstall_root.open_subkey(&app_key_name) {
                let mut app_info = HashMap::new();

                if let Ok(display_name) = app_registry_key.get_value::<String, _>("DisplayName") {
                    app_info.insert(String::from("app_name"), display_name);
                }
                if let Ok(display_icon) = app_registry_key.get_value::<String, _>("DisplayIcon") {
                    app_info.insert(String::from("app_icon"), display_icon);
                }
                if let Ok(install_location) = app_registry_key.get_value::<String, _>("InstallLocation") {
                    app_info.insert(String::from("app_location"), install_location);
                }
                if let Ok(publisher) = app_registry_key.get_value::<String, _>("Publisher") {
                    app_info.insert(String::from("app_publisher"), publisher);
                }

                // only keep entries that actually have a display name —
                // many uninstall subkeys are patches/updates with no name at all
                if app_info.contains_key("app_name") {
                    installed_apps.push(app_info);
                }
            }
        }
    }

    installed_apps
}