// internal/app/app.go
package app

import (
	"fmt"
	"os"
	"os/exec"

	"scantrix/internal/cli"
	"scantrix/internal/scanner"
	"scantrix/internal/ui"
)

type App struct {
	options cli.Options
}

func NewApp(options cli.Options) *App {
	return &App{options: options}
}

func (a *App) Run() error {
	path, files, findings, err := scanner.Prepare(a.options)
	if err != nil {
		return err
	}

	if a.options.Watch {
		return ui.RunRealtime(path, files, findings, a.options.Severity, a.options.Exclude)
	}

	if len(findings) == 0 {
		fmt.Println("✅ No vulnerabilities found.")
		return nil
	}

	return ui.Run(findings)
}

func SelfUpgrade() {
	fmt.Println("🔄 Upgrading Scantrix...")
	cmd := exec.Command("bash", "-c", "curl -sSL https://raw.githubusercontent.com/vzsigmond/scantrix/main/scripts/install.sh | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Upgrade failed: %v\n", err)
		os.Exit(1)
	}
}
