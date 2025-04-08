package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700"))

	CriticalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F5F"))

	WarningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFAF00"))

	InfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5FAFFF"))

	FileStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#AAAAAA"))

	AdviceStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Margin(1, 0)
)
