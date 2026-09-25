use serde::Serialize;
use super::find_uninstall_root_keys::find_uninstall_root_keys;

#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
pub struct InstalledApp {
    pub id: String,
    pub name: String,
    pub publisher: String,
    pub install_path: String,
    pub executables: Vec<String>,
    pub icon_base64: String,
    pub app_type: String,
    pub package_family_name: String,
    pub package_sid: String,    
}
/// Walks every uninstall root key, opens each app's own subkey inside it,
/// and pulls out the display name and icon path where available.
pub fn collect_installed_apps() -> Vec<InstalledApp> {
    let mut installed_apps = Vec::new();
    let uninstall_root_keys = find_uninstall_root_keys();

    for uninstall_root in uninstall_root_keys.iter() {
        for app_key_name in uninstall_root.enum_keys().flatten() {
            if let Ok(app_registry_key) = uninstall_root.open_subkey(&app_key_name) {
                let name = app_registry_key
                    .get_value::<String, _>("DisplayName")
                    .unwrap_or_default();

                if name.is_empty() {
                    continue;
                }

                let publisher = app_registry_key
                    .get_value::<String, _>("Publisher")
                    .unwrap_or_default();

                let install_path = app_registry_key
                    .get_value::<String, _>("InstallLocation")
                    .unwrap_or_default();

                installed_apps.push(InstalledApp {
                    id: app_key_name.clone(),
                    name,
                    publisher,
                    install_path,
                    executables: Vec::new(),
                    icon_base64: String::new(),
                    app_type: String::from("win32"),
                    package_family_name: String::new(),
                    package_sid: String::new(),
                });
            }
        }
    }

    installed_apps
}