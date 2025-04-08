# 🛡️ Scantrix

Scantrix is a powerful, blazing-fast code scanner built in Go, designed to help developers detect security vulnerabilities in real-time.

![Repo Size](https://img.shields.io/github/repo-size/vzsigmond/scantrix)
![Latest Tag](https://img.shields.io/github/v/tag/vzsigmond/scantrix)

It currently supports:

- Remote Code Execution (RCE)
- SQL Injection
- Cross-Site Scripting (XSS)
- CSRF
- Insecure Cryptography
- Open Redirects
- ...and more!

> ⚠️ **This project is currently in Alpha and under active development. Use at your own risk.**

## 🔧 Features

- Scans your local codebase **recursively**
- Detects common vulnerabilities using regex rules
- Supports `--exclude` filters and severity-based filtering
- Beautiful interactive **TUI** powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Generates real-time results from file changes (`--watch`)
- Scan directly from git repos via `--git <repo-url>`
- Written in Go for maximum performance
- Checks for latest version and suggests `--self-upgrade`

## 🚀 Installation

### 🔁 Auto Installer (Linux/macOS/WSL)

```bash
curl -sSL https://raw.githubusercontent.com/vzsigmond/scantrix/main/scripts/install.sh | bash
```

### 🪟 Windows PowerShell

```powershell
irm https://raw.githubusercontent.com/vzsigmond/scantrix/main/scripts/install.ps1 | iex
```

### 🍫 Chocolatey (soon)

```powershell
choco install scantrix
```

## 🛠 Usage

```bash
scantrix --path ./myapp [--watch] [--exclude="regex"] [--severity=critical|warning|info"]
scantrix --git https://github.com/laravel/laravel --severity=critical
```

Examples:

```bash
scantrix --path ./myapp --exclude="node_modules|tests"
scantrix --watch --severity=warning --path ./myapp
scantrix --git https://github.com/drupal/drupal
```

## 🆙 Self-Upgrade

```bash
scantrix --self-upgrade
```

Scantrix will also notify you when a newer version is available.

---

## 📂 Log Output

Scantrix writes logs to `logs/debug.log` **only if** `--debug` is provided.

---

## 📦 Building From Source

```bash
git clone https://github.com/vzsigmond/scantrix.git
cd scantrix
go build -o scantrix ./cmd/scantrix
```

---

## 🧩 Planned Features

- Better support for frameworks & CMSes like: Drupal, Laravel, WordPress
- CVE vulnerability references
- GitHub Action integration
- CI/CD mode (non-TUI)
- Custom rule config

---

## 🧑‍💻 Contributing

Pull requests are welcome! Fork the repo and create a new branch for features or fixes.
