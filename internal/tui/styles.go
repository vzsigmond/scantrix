package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	TopBarHeight    = 3
	BottomBarHeight = 3
	LeftPaneWidth   = 30

	MinWindowWidth  = 50
	MinWindowHeight = 15
)

// TopBarStyle: a rounded border for the top bar.
var TopBarStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FFFFFF")).
	Padding(0, 1)

// Smaller style for text in the top bar.
var TopBarTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#EEEEEE"))

// LeftPaneStyle: normal border on the right side.
var LeftPaneStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderRight(true)

// RightPaneStyle: some padding, no border.
var RightPaneStyle = lipgloss.NewStyle().
	Padding(0, 1)

// BottomBarStyle: a rounded border, darker background, white text.
var BottomBarStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FFFFFF")).
	Padding(0, 1).
	Background(lipgloss.Color("#353533")).
	Foreground(lipgloss.Color("#FFFFFF"))

// A naive gradient for the top bar title. Adjust colors as needed.
func GradientText(text string) string {
	startHex := "#F25D94" // pinkish
	endHex := "#EDFF82"   // greenish

	rs := []rune(text)
	if len(rs) < 2 {
		return text
	}

	var sb strings.Builder
	for i, r := range rs {
		t := float64(i) / float64(len(rs)-1)
		c := interpolateColor(startHex, endHex, t)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(c))
		sb.WriteString(style.Render(string(r)))
	}
	return sb.String()
}

// Simple linear interpolation for two hex colors
func interpolateColor(start, end string, t float64) string {
	sr, sg, sb := parseHexColor(start)
	er, eg, eb := parseHexColor(end)

	r := sr + t*(er-sr)
	g := sg + t*(eg-sg)
	b := sb + t*(eb-sb)

	return fmt.Sprintf("#%02x%02x%02x", int(r), int(g), int(b))
}

func parseHexColor(hex string) (float64, float64, float64) {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return float64(r), float64(g), float64(b)
}
