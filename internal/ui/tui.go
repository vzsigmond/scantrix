// internal/ui/tui.go
package ui

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"scantrix/internal/scanner"
	"scantrix/internal/types"
	"scantrix/internal/logger"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/fsnotify/fsnotify"
)

type item struct {
	f scanner.Finding
}

func (i item) Title() string {
	icon := map[string]string{
		"critical": "❌",
		"warning":  "⚠️",
		"info":     "ℹ️",
	}[i.f.Severity]
	return logger.Sprintf("%s %s", icon, i.f.Title)
}

func (i item) Description() string {
	return logger.Sprintf("File: %s:%d\nAdvice: %s", i.f.File, i.f.Line, i.f.Advice)
}

func (i item) FilterValue() string {
	return i.f.Title
}

type model struct {
	list        list.Model
	allFindings []scanner.Finding
	width       int
	height      int
	sub         chan []scanner.Finding
}

func New(findings []scanner.Finding, sub chan []scanner.Finding) model {
	items := convertFindings(findings)
	l := list.New(items, list.NewDefaultDelegate(), 70, 20)
	l.Title = "Scantrix – Vulnerabilities"
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	return model{
		list:        l,
		allFindings: findings,
		sub:         sub,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForUpdate(m.sub),
	)
}

type findingsUpdateMsg []scanner.Finding

func waitForUpdate(sub chan []scanner.Finding) tea.Cmd {
	return func() tea.Msg {
		findings := <-sub
		logger.Log("✅ Received %d updated findings", len(findings))
		return findingsUpdateMsg(findings)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			return m.withFilter("critical"), nil
		case "w":
			return m.withFilter("warning"), nil
		case "i":
			return m.withFilter("info"), nil
		case "a":
			return m.withFilter(""), nil
		}

	case findingsUpdateMsg:
		logger.Log("✅ TUI received findings update")
		m.allFindings = msg
		m.list.SetItems(convertFindings(msg))
		return m, waitForUpdate(m.sub)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString(m.list.View())
	b.WriteString("\n\n[ q ] quit  [ c ] critical  [ w ] warning  [ i ] info  [ a ] all\n")
	return b.String()
}

func (m model) withFilter(severity string) model {
	filtered := make([]list.Item, 0)
	for _, f := range m.allFindings {
		if severity == "" || f.Severity == severity {
			filtered = append(filtered, item{f})
		}
	}
	m.list.SetItems(filtered)
	return m
}

func Run(findings []scanner.Finding) error {
	logger.Log("⚡ Entered Run()")
	sub := make(chan []scanner.Finding, 1)
	m := New(findings, sub)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()

	return err
}

func RunRealtime(path string, rules []types.Rule, exclude *regexp.Regexp, severity string) error {
	logger.Log("⚡ Entered RunRealtime()")
	initialFindings, _ := scanner.ScanDirectory(path, rules, exclude, severity)
	sub := make(chan []scanner.Finding, 1)
	m := New(initialFindings, sub)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logger.Log("❌ watcher.Events channel closed")
					return
				}
				logger.Log("📁 EVENT: %s", event.String())
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 &&
					!strings.HasSuffix(event.Name, "~") &&
					!strings.Contains(event.Name, ".sw") {
					logger.Log("📡 Rescanning due to: %s", event.Name)
					newFindings, _ := scanner.ScanDirectory(path, rules, exclude, severity)
					sub <- newFindings
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Log("Watcher error: %v", err)
			}
		}
	}()

	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			logger.Log("📂 Watching: %s", p)
			return watcher.Add(p)
		}
		return nil
	})
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func convertFindings(findings []scanner.Finding) []list.Item {
	items := make([]list.Item, len(findings))
	for i, f := range findings {
		items[i] = item{f}
	}
	return items
}
