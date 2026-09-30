const ALLOWED: &[&str] = &["office", "visual studio", "vscode", "edge", "teams", "onedrive"];
const SYSTEM: &[&str] = &["update", "redistributable", "runtime", ".net", "sdk", "tool"];


pub fn is_system_app(name: &str, publisher: &str, install_path: &str) -> bool {
    let lower_name = name.to_lowercase();
    let lower_publisher: String = publisher.to_lowercase();
    let lower_install_path: String = install_path.to_lowercase();

    if lower_publisher.contains("microsoft"){
        if ALLOWED.iter().any(|keyword: &&str| lower_name.contains(keyword)) {
            return false;
        }
        if SYSTEM.iter().any(|keyword: &&str| lower_name.contains(keyword)) {
            return true;
        }
    }

    if lower_install_path.contains("\\windows\\") {
        return true;
    }

    return false;

}