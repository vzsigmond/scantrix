package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TreeNode represents a directory or file in the tree.
type TreeNode struct {
	Name     string
	IsDir    bool
	Expanded bool
	Children []*TreeNode
	Parent   *TreeNode
	Severity string
	// If you want to store vulnerabilities or severity, add fields here.
}

// NodeItem implements list.Item so we can show the tree in a Bubble Tea list.
type NodeItem struct {
	Node     *TreeNode
	Indent   int
	TitleStr string // how it appears in the list
}

func (ni NodeItem) Title() string       { return ni.TitleStr }
func (ni NodeItem) Description() string { return "" }       // no secondary text
func (ni NodeItem) FilterValue() string { return ni.TitleStr }

// Model is our main TUI model. We store all top-level roots so we never lose them.
type Model struct {
	roots []*TreeNode // all top-level directories/files in the entire project

	list     list.Model
	viewport viewport.Model

	width, height int
	selectedNode  *TreeNode
}

// NewModel: pass in all top-level roots so we can flatten them anytime.
func NewModel(roots []*TreeNode) Model {
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0) // 0 spacing to remove big gaps

	// Flatten once initially
	items := flattenTree(roots)

	l := list.New(items, delegate, LeftPaneWidth, 10)
	l.Title = "Project Tree"
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	// Right-side viewport for details
	vp := viewport.New(50, 10)
	vp.SetContent("Select a directory or file on the left...")

	m := Model{
		roots:    roots,
		list:     l,
		viewport: vp,
	}
	return m
}

// flattenTree collects NodeItems from multiple top-level roots.
func flattenTree(roots []*TreeNode) []list.Item {
	var result []list.Item
	for _, r := range roots {
		result = append(result, flattenNode(r, 0)...)
	}
	return result
}

// flattenNode recursively collects NodeItems for a single node.
func flattenNode(n *TreeNode, indent int) []list.Item {
	title := renderNodeTitle(n, indent)
	ni := NodeItem{
		Node:     n,
		Indent:   indent,
		TitleStr: title,
	}

	var items []list.Item
	items = append(items, ni)

	if n.IsDir && n.Expanded {
		for _, c := range n.Children {
			items = append(items, flattenNode(c, indent+1)...)
		}
	}
	return items
}

func renderNodeTitle(n *TreeNode, indent int) string {
	prefix := strings.Repeat("  ", indent)
	var icon string
	if n.IsDir {
		if n.Expanded {
			icon = "📂"
		} else {
			icon = "📁"
		}
	} else {
		icon = "📄"
	}
	return fmt.Sprintf("%s%s %s", prefix, icon, n.Name)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// If terminal is too small, fallback
		if m.width < MinWindowWidth || m.height < MinWindowHeight {
			return m, nil
		}

		mainHeight := m.height - (TopBarHeight + BottomBarHeight)
		if mainHeight < 3 {
			mainHeight = 3
		}

		m.list.SetSize(LeftPaneWidth, mainHeight)
		rightWidth := m.width - LeftPaneWidth
		if rightWidth < 10 {
			rightWidth = 10
		}
		m.viewport.Width = rightWidth
		m.viewport.Height = mainHeight

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			m.toggleDir()
			return m, nil

		case "up", "down":
			prev := m.list.Index()
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			if m.list.Index() != prev {
				m.updateSelection()
			}
			return m, cmd

		case "pgup", "pgdown", "left", "right":
			// let the viewport handle scrolling
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	// pass other msgs to list + viewport
	}

	prevIndex := m.list.Index()
	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	if m.list.Index() != prevIndex {
		m.updateSelection()
	}

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(listCmd, vpCmd)
}

// toggleDir flips IsDir/Expanded if the selected item is a directory
func (m *Model) toggleDir() {
	sel := m.list.SelectedItem()
	if sel == nil {
		return
	}
	ni, ok := sel.(NodeItem)
	if !ok {
		return
	}
	n := ni.Node
	if !n.IsDir {
		return
	}

	// Flip expansion
	n.Expanded = !n.Expanded

	// Re-flatten from all top-level roots, so we don't lose siblings/parents.
	items := flattenTree(m.roots)
	m.list.SetItems(items)

	// Attempt to re-select the same node by name:
	m.restoreSelection(n.Name)
	m.updateSelection()
}

// re-select the same node name
func (m *Model) restoreSelection(name string) {
	items := m.list.Items()
	for i, it := range items {
		nodeItem, _ := it.(NodeItem)
		if nodeItem.Node.Name == name {
			m.list.Select(i)
			break
		}
	}
}

// updateSelection updates the viewport content
func (m *Model) updateSelection() {
	sel := m.list.SelectedItem()
	if sel == nil {
		m.viewport.SetContent("No selection.")
		m.selectedNode = nil
		return
	}
	ni, ok := sel.(NodeItem)
	if !ok {
		m.viewport.SetContent("Unknown selection item type.")
		m.selectedNode = nil
		return
	}
	m.selectedNode = ni.Node
	m.viewport.SetContent(buildDetails(ni.Node))
}

// buildDetails: show info about the node in the right pane
func buildDetails(n *TreeNode) string {
	if n.IsDir {
		return fmt.Sprintf("Directory: %s\n\nContains %d children.\nUse PgUp/PgDn/left/right to scroll.\n", n.Name, len(n.Children))
	}
	return fmt.Sprintf("File: %s\n\nAdd details or vulnerabilities here.\nScrolling is possible if text is long.\n", n.Name)
}

func (m Model) View() string {
	if m.width < MinWindowWidth || m.height < MinWindowHeight {
		return fmt.Sprintf("Terminal is too small (%dx%d). Resize, please.\n", m.width, m.height)
	}

	// top bar
	topBar := m.renderTopBar()

	// bottom bar
	botBar := m.renderBottomBar()

	// main area
	mainHeight := m.height - (TopBarHeight + BottomBarHeight)
	leftView := m.list.View()
	leftPane := LeftPaneStyle.Width(LeftPaneWidth).Height(mainHeight).Render(leftView)

	rightWidth := m.width - LeftPaneWidth
	if rightWidth < 10 {
		rightWidth = 10
	}
	rightPane := RightPaneStyle.Width(rightWidth).Height(mainHeight).Render(m.viewport.View())

	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	ui := lipgloss.JoinVertical(lipgloss.Left, topBar, mainArea, botBar)
	return ui
}

func (m Model) renderTopBar() string {
	title := GradientText("  Scantrix Directory Tree  ")
	info := TopBarTextStyle.Render("\nExpand/collapse with Enter, navigate with up/down.")
	barContent := title + info
	return TopBarStyle.
		Width(m.width).
		Height(TopBarHeight).
		Render(barContent)
}

func (m Model) renderBottomBar() string {
	legend := "[ up/down ] navigate  [ enter ] expand/collapse  [ pgup/pgdown/left/right ] scroll  [ q ] quit"
	return BottomBarStyle.
		Width(m.width).
		Height(BottomBarHeight).
		Render(legend)
}

// StartProgram is a helper to run the TUI. Typically invoked by a Cobra command.
func StartProgram(m Model) error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
