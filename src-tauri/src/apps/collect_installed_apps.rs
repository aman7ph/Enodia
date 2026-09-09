use std::collections::HashMap;

use super::find_uninstall_root_keys::find_uninstall_root_keys;

/// Walks every uninstall root key, opens each app's own subkey inside it,
/// and pulls out the display name and icon path where available.
pub fn collect_installed_apps() -> Vec<HashMap<String, String>> {
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