# Variables
APP_NAME=cura
LAUNCHER_NAME=CuraLauncher
LD_FLAGS="-w -s"
BIN_DIR=build/bin

.PHONY: up b cb launcher clean

# Standard startup
up:
	wails dev

# Quick rebuild and restart (AMD64)
b: launcher
	wails build -clean -ldflags $(LD_FLAGS)

# Full Production Rebuild
# This compiles the launcher, the app for both architectures, and sets up the environment
cb: launcher
	@echo "--- Starting Full Production Build ---"
	wails build -clean -ldflags $(LD_FLAGS) -platform windows/amd64,windows/arm64
	
	@echo "--- Preparing Binary Environment ---"
	@if not exist "$(BIN_DIR)" mkdir "$(BIN_DIR)"
	
# Create empty log file
	type nul > "$(BIN_DIR)/$(APP_NAME).log"
	
# Copy Configuration (using /Y to suppress overwrite prompts)
	copy /Y "settings.toml" "$(BIN_DIR)\settings.toml"
	
# Copy Routines Directory (S=subdirectories, I=assume directory if destination doesn't exist, Y=overwrite)
	xcopy "routines" "$(BIN_DIR)\routines" /S /I /Y /E
	
	@echo "--- Build Complete ---"

# Build the Administrative Launcher
launcher:
	@echo "--- Compiling Admin Launcher ---"
	@if not exist "$(BIN_DIR)" mkdir "$(BIN_DIR)"
	@echo go build -ldflags=$(LD_FLAGS) -o $(BIN_DIR)/$(LAUNCHER_NAME).exe ./cmd/launcher/main.go

# Clean build artifacts
clean:
	@if exist "build\bin" rd /s /q "build\bin"
	@if exist "build\windows" rd /s /q "build\windows"