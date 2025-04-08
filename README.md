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


## 🚀 Quick Start

```bash
git clone https://github.com/vzsigmond/scantrix.git
cd scantrix
./bin/scantrix --watch ./tests/fixtures
```

## 🛠 Usage

```bash
scantrix [--watch] [--exclude="regex"] [--severity=critical|warning|info] [--git=url] /path/to/project
```

Examples:

```bash
scantrix ./myapp --exclude="node_modules|tests"
scantrix --watch ./myapp "
scantrix --git https://github.com/laravel/laravel --severity=critical
```

## 📂 Log Output

Scantrix writes logs to `logs/debug.log`.


## 🧩 Planned Features

- Better support for frameworks & CMSes like: Drupal, Laravel, Worpdress.
- CVE vulnerability
- GitHub Action integration
- CI/CD mode (non-TUI)
- Custom rules config


## 🧑‍💻 Contributing

Pull requests are welcome! Fork the repo and create a new branch for features or fixes.



