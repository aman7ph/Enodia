use std::collections::HashMap;
mod apps;

#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[tauri::command]
fn list_installed_apps() -> Vec<HashMap<String, String>> {
    apps::collect_installed_apps::collect_installed_apps()
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![greet,list_installed_apps])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
