# 🛡️ Cura
**Precision System Resource Enforcement for Windows.**

Cura is a lightweight, Go-powered utility designed to protect your system's stability. By monitoring real-time RAM and CPU usage, Cura automatically identifies and terminates background "idle" processes when your custom-defined memory caps are breached.

![Cura Dashboard Preview](https://via.placeholder.com/800x450?text=Cura+Dashboard+Interface)

## Features
* **The Cap:** Set a hard memory reserve percentage to ensure your OS always has breathing room.
* **Smart Enforcement:** Targets low-CPU "idle" processes first to avoid interrupting active work.
* **Real-time Telemetry:** High-fidelity dashboard built with React and Wails.
* **System Protection:** Built-in whitelist for critical Windows services and user-defined apps.

## Getting Started

### Prerequisites
* [Go](https://golang.org/doc/install) (1.21+)
* [Node.js](https://nodejs.org/en/download/) (v18+)
* [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Installation
1. Clone the repository:
   ```bash
   git clone [https://github.com/your-username/cura.git](https://github.com/your-username/cura.git)
   cd cura

```

2. Install dependencies and run in dev mode:
```bash
wails dev

```


3. Build for production:
```bash
wails build -platform windows/amd64 -clean -ldflags "-s -w"

```



## Configuration

Cura uses a `settings.toml` file in the root directory to persist user preferences across sessions.

```toml
[enforcement]
is_enforced = true
memory_cap = 80.0
cpu_ceiling = 70.0

```

## Safety Disclaimer

Cura has the power to terminate processes. While it includes a default protection list for system-critical tasks, use it with caution. Always whitelist your unsaved work-heavy applications.

```

LICENSE: MIT