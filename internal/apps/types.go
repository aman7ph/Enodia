package apps

// InstalledApp represents a discovered application
type InstalledApp struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Publisher         string   `json:"publisher"`
	InstallPath       string   `json:"installPath"`
	Executables       []string `json:"executables"`
	IconBase64        string   `json:"iconBase64"`
	AppType           string   `json:"appType"`
	PackageFamilyName string   `json:"packageFamilyName"`
	PackageSID        string   `json:"packageSID"`
}
