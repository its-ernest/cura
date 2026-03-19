# standard startup
up:
	wails dev

# rebuild and restart
b:
	wails build -clean -ldflags "-w -s"

# clean rebuild
cb:
	wails build -clean -ldflags "-w -s"
	copy "cura.log" "build/bin/cura.log"
	copy "settings.toml" "build/bin/settings.toml"
