use std::fs;
use std::path::Path;

/// Files more than this many folders below the install folder are ignored.
const MAX_DEPTH: u32 = 2;

/// Finds .exe files inside an install folder, up to two folders deep.
pub(super) fn find_executables(install_path: &str) -> Vec<String> {
    let mut executables = Vec::new();
    collect_executables(Path::new(install_path), 0, &mut executables);
    executables
}

fn collect_executables(dir: &Path, depth: u32, executables: &mut Vec<String>) {
    let Ok(entries) = fs::read_dir(dir) else {
        return;
    };

    for entry in entries.flatten() {
        let path = entry.path();

        if path.is_dir() {
            if depth < MAX_DEPTH {
                collect_executables(&path, depth + 1, executables);
            }
        } else if path
            .extension()
            .is_some_and(|ext| ext.eq_ignore_ascii_case("exe"))
        {
            executables.push(path.to_string_lossy().into_owned());
        }
    }
}