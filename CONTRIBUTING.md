# Contributing to Cura

First off, thank you for considering contributing to Cura! 

## Development Workflow
1. **Fork** the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. Ensure the Go code is formatted: `go fmt ./...`.
4. Ensure the React frontend follows the existing styling in `App.css`.

## Pull Request Process
1. Update the `README.md` with details of changes to the interface or configuration if applicable.
2. The PR should target the `main` branch.
3. Once the PR is merged, your changes will be included in the next pre-release.

## Roadmap
* [ ] Implement Process Whitelist UI.
* [ ] Add System Tray (Minimize to Tray) support.
* [ ] CPU Ceiling enforcement logic.
* [ ] Dark/Light theme toggle.

```

---

### Final Step for GitHub

1. Save these files as `README.md`, `LICENSE`, and `CONTRIBUTING.md` in your project root.
2. Run these commands to push them:
```powershell
git add README.md LICENSE CONTRIBUTING.md
git commit -m "docs: add readme, license, and contributing guidelines"
git push origin main

```
