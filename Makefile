# standard startup
up:
	wails dev

# rebuild and restart
b:
	wails build -clean -ldflags "-w -s"

# clean rebuild
cb:
	wails build -clean -ldflags "-w -s" -platform windows/amd64,windows/arm64
	type nul > "build/bin/cura.log"
	copy "settings.toml" "build/bin/settings.toml"
	copy "routines/" "build/bin/routines/"
