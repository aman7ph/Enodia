use winreg::enums::*;
use winreg::RegKey;

/// Opens the registry locations where Windows lists installed applications.
/// There are three separate spots because of 32-bit/64-bit and per-user vs
/// machine-wide installs.
pub(super) fn find_uninstall_root_keys() -> Vec<RegKey> {
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