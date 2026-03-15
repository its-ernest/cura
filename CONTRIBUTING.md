# Contributing to Cura 🛡️

First off, thank you for considering contributing to Cura! It's people like you who make open-source tools better for everyone.

## Tech Stack
- **Backend:** Go 1.21+
- **Frontend:** React + Vite
- **Framework:** [Wails](https://wails.io/) (Native Go + Webview2)
- **Monitoring:** gopsutil

## Development Setup

1. **Prerequisites:**
   - Ensure you have the [Wails CLI](https://wails.io/docs/gettingstarted/installation) installed.
   - You must be on **Windows** to test the enforcement logic.

2. **Clone and Run:**
   ```bash
   git clone [https://github.com/its-ernest/cura.git](https://github.com/its-ernest/cura.git)
   cd cura
   wails dev

3. **Directory Structure:**
* `/pkg/memory`: Core enforcement logic, staleness intelligence, and process management.
* `/pkg/logging`: Thread-safe file logging.
* `/frontend`: React dashboard and UI components.

4. **Testing:**
* `/pkg/memory/utils_test.go`: This critical test must pass if you alter or improve protection of critical processes logic(Currently Windows-specific)

## 📝 Contribution Guidelines

### Code Standards

* **Production Thinking:** We don't just write code that works; we write code that is clean and modular.
* **Documentation:** Maintain concise, practical documentation
* **Version Control:** Before adding new logic, double-check resources for the correct version of dependencies (e.g., v3 vs v2).

### Pull Request Process

1. Create a branch for your feature or fix: `git checkout -b feat/your-feature-name`.
2. Ensure your code is formatted: `go fmt ./...`.
3. If you've modified the enforcement logic, verify it against the `cura.log` output.
4. Submit your PR against the `main` branch with a clear description of the changes.

## Roadmap

* [ ] Process Whitelist UI integration.
* [ ] System Tray (Minimize to Tray) support.
* [ ] Log viewer within the React dashboard.
* [ ] CPU Ceiling enforcement implementation.

---

*By contributing, you agree that your contributions will be licensed under the MIT License.*
