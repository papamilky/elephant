package main

type DesktopFile struct {
	Data
	Actions []Data
}

var (
	Name       = "desktopapplications"
	NamePretty = "Desktop Applications"
)

// TODO: watch folders for changes
func Load() {
	loadConfig()
	loadFiles()
}
