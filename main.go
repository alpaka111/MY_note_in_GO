package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	screenSplash screen = iota
	screenLogin
	screenMenu
	screenAddNote
	screenViewNotes
	screenSearch
	screenStats
	screenSettings
)

type sortMode int

const (
	sortByDate sortMode = iota
	sortByTitle
	sortByTags
)

type model struct {
	screen        screen
	notebook      *Notebook
	filename      string
	password      string
	passwordBuf   string
	titleBuf      string
	contentBuf    string
	tagsBuf       string
	cursor        int
	selected      int
	searchQuery   string
	err           error
	success       string
	width         int
	height        int
	splashTicks   int
	ready         bool
	animFrame     int
	showPassword  bool
	sortMode      sortMode
	viewMode      int // 0 = list, 1 = grid, 2 = detailed
	filterTag     string
	scrollOffset  int
	maxScroll     int
}

type tickMsg struct{}
type animMsg struct{}

func initialModel() model {
	return model{
		screen:   screenSplash,
		filename: "notatki.alpaka",
		sortMode: sortByDate,
		viewMode: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), animate())
}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func animate() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(_ time.Time) tea.Msg {
		return animMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tickMsg:
		if m.screen == screenSplash {
			m.splashTicks++
			if m.splashTicks > 20 {
				m.screen = screenLogin
			}
			return m, tick()
		}
		return m, nil

	case animMsg:
		m.animFrame = (m.animFrame + 1) % 4
		return m, animate()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.screen != screenLogin && m.screen != screenSplash {
				m.screen = screenMenu
				m.err = nil
				m.success = ""
				m.cursor = 0
				m.scrollOffset = 0
			}
			return m, nil
		}

		switch m.screen {
		case screenSplash:
			m.screen = screenLogin
			return m, nil
		case screenLogin:
			return m.updateLogin(msg)
		case screenMenu:
			return m.updateMenu(msg)
		case screenAddNote:
			return m.updateAddNote(msg)
		case screenViewNotes:
			return m.updateViewNotes(msg)
		case screenSearch:
			return m.updateSearch(msg)
		case screenStats:
			return m.updateStats(msg)
		case screenSettings:
			return m.updateSettings(msg)
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Inicjalizacja..."
	}

	switch m.screen {
	case screenSplash:
		return m.viewSplash()
	case screenLogin:
		return m.viewLogin()
	case screenMenu:
		return m.viewMenu()
	case screenAddNote:
		return m.viewAddNote()
	case screenViewNotes:
		return m.viewViewNotes()
	case screenSearch:
		return m.viewSearch()
	case screenStats:
		return m.viewStats()
	case screenSettings:
		return m.viewSettings()
	}

	return ""
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Błąd: %v\n", err)
		os.Exit(1)
	}
}

// === ULTIMATE STYLES ===

var (
	// Gradient color palette
	colorPalette = []lipgloss.Color{
		lipgloss.Color("#FF6B9D"), // Pink
		lipgloss.Color("#C792EA"), // Purple
		lipgloss.Color("#82AAFF"), // Blue
		lipgloss.Color("#89DDFF"), // Cyan
		lipgloss.Color("#C3E88D"), // Green
		lipgloss.Color("#FFCB6B"), // Yellow
		lipgloss.Color("#F78C6C"), // Orange
		lipgloss.Color("#FF5370"), // Red
	}

	// Base colors
	primary   = lipgloss.Color("#FF6B9D")
	secondary = lipgloss.Color("#C792EA")
	accent    = lipgloss.Color("#82AAFF")
	success   = lipgloss.Color("#C3E88D")
	warning   = lipgloss.Color("#FFCB6B")
	danger    = lipgloss.Color("#FF5370")
	text      = lipgloss.Color("#EEFFFF")
	textDim   = lipgloss.Color("#B0BEC5")
	muted     = lipgloss.Color("#676E95")
	bg        = lipgloss.Color("#1E1E2E")
	bgLight   = lipgloss.Color("#2A2A40")
	bgDark    = lipgloss.Color("#181825")

	// App title with gradient effect
	appTitleStyle = lipgloss.NewStyle().
			Foreground(text).
			Background(primary).
			Padding(0, 2).
			Bold(true).
			MarginBottom(1)

	// Box styles with depth
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondary).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	focusedBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(primary).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1).
			BorderBackground(bgLight)

	glowBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(accent).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1).
			Background(bgDark)

	// Menu styles
	menuItemStyle = lipgloss.NewStyle().
			Foreground(text).
			Padding(0, 2).
			MarginLeft(2)

	selectedMenuStyle = lipgloss.NewStyle().
				Foreground(bg).
				Background(primary).
				Padding(0, 2).
				MarginLeft(1).
				Bold(true)

	hoveredMenuStyle = lipgloss.NewStyle().
				Foreground(primary).
				Padding(0, 2).
				MarginLeft(2).
				Bold(true)

	// Input styles with focus states
	inputStyle = lipgloss.NewStyle().
			Foreground(text).
			Background(bgDark).
			Padding(0, 1).
			Width(70)

	focusedInputStyle = lipgloss.NewStyle().
				Foreground(text).
				Background(bgLight).
				Padding(0, 1).
				Width(70).
				Border(lipgloss.NormalBorder()).
				BorderForeground(accent)

	labelStyle = lipgloss.NewStyle().
			Foreground(textDim).
			Bold(true).
			MarginBottom(0)

	focusedLabelStyle = lipgloss.NewStyle().
				Foreground(primary).
				Bold(true).
				MarginBottom(0)

	// Note card styles
	noteCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondary).
			Padding(1, 2).
			MarginBottom(1).
			Width(75).
			Background(bgDark)

	selectedNoteStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(primary).
				Padding(1, 2).
				MarginBottom(1).
				Width(75).
				Background(bgLight)

	highlightNoteStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(accent).
				Padding(1, 2).
				MarginBottom(1).
				Width(75).
				Background(bgLight)

	noteTitleStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			Underline(true)

	noteMetaStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true)

	noteContentStyle = lipgloss.NewStyle().
				Foreground(textDim).
				MarginTop(1)

	// Tag styles with colors
	tagStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(bg).Background(primary).Padding(0, 1).MarginRight(1),
		lipgloss.NewStyle().Foreground(bg).Background(secondary).Padding(0, 1).MarginRight(1),
		lipgloss.NewStyle().Foreground(bg).Background(accent).Padding(0, 1).MarginRight(1),
		lipgloss.NewStyle().Foreground(bg).Background(success).Padding(0, 1).MarginRight(1),
		lipgloss.NewStyle().Foreground(bg).Background(warning).Padding(0, 1).MarginRight(1),
	}

	// Status styles
	successStyle = lipgloss.NewStyle().
			Foreground(success).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(danger).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warning).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(accent)

	// Help styles
	helpStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	// Stats styles
	statCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(1, 2).
			Margin(0, 1).
			Width(20).
			Align(lipgloss.Center)

	statNumberStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	statLabelStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true)

	// Progress bar
	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(muted)

	progressFullStyle = lipgloss.NewStyle().
				Foreground(success)
)

func renderGradientText(text string, colors []lipgloss.Color) string {
	if len(text) == 0 {
		return ""
	}

	var result strings.Builder
	colorCount := len(colors)

	for i, char := range text {
		colorIdx := (i * colorCount) / len(text)
		if colorIdx >= colorCount {
			colorIdx = colorCount - 1
		}
		style := lipgloss.NewStyle().Foreground(colors[colorIdx])
		result.WriteString(style.Render(string(char)))
	}

	return result.String()
}

func renderHeader(title string, subtitle string) string {
	logo := `
   ▄▄▄       ██▓     ██▓███   ▄▄▄       ██ ▄█▀▄▄▄      
  ▒████▄    ▓██▒    ▓██░  ██▒▒████▄     ██▄█▒▒████▄    
  ▒██  ▀█▄  ▒██░    ▓██░ ██▓▒▒██  ▀█▄  ▓███▄░▒██  ▀█▄  
  ░██▄▄▄▄██ ▒██░    ▒██▄█▓▒ ▒░██▄▄▄▄██ ▓██ █▄░██▄▄▄▄██ 
   ▓█   ▓██▒░██████▒▒██▒ ░  ░ ▓█   ▓██▒▒██▒ █▄▓█   ▓██▒
   ▒▒   ▓▒█░░ ▒░▓  ░▒▓▒░ ░  ░ ▒▒   ▓▒█░▒ ▒▒ ▓▒▒▒   ▓▒█░
    ▒   ▒▒ ░░ ░ ▒  ░░▒ ░       ▒   ▒▒ ░░ ░▒ ▒░ ▒   ▒▒ ░
    ░   ▒     ░ ░   ░░         ░   ▒   ░ ░░ ░  ░   ▒   
        ░  ░    ░  ░               ░  ░░  ░        ░  ░`

	gradientLogo := renderGradientText(logo, colorPalette)

	var b strings.Builder
	b.WriteString(gradientLogo)
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(80).
		Render("🦙 " + title + " 🦙"))
	b.WriteString("\n")

	if subtitle != "" {
		b.WriteString(lipgloss.NewStyle().
			Foreground(muted).
			Italic(true).
			Align(lipgloss.Center).
			Width(80).
			Render(subtitle))
		b.WriteString("\n")
	}

	separator := strings.Repeat("━", 80)
	b.WriteString(lipgloss.NewStyle().Foreground(secondary).Render(separator))
	b.WriteString("\n")

	return b.String()
}

func renderFooter(help string) string {
	var b strings.Builder
	separator := strings.Repeat("━", 80)
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(secondary).Render(separator))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(help))
	return b.String()
}

func renderHelp(keys ...string) string {
	var parts []string
	for i := 0; i < len(keys); i += 2 {
		if i+1 < len(keys) {
			key := helpKeyStyle.Render(keys[i])
			desc := keys[i+1]
			parts = append(parts, key+" "+desc)
		}
	}
	return strings.Join(parts, " │ ")
}

func renderProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	filled := (current * width) / total
	empty := width - filled

	bar := progressFullStyle.Render(strings.Repeat("█", filled)) +
		progressEmptyStyle.Render(strings.Repeat("░", empty))

	percentage := (current * 100) / total
	label := fmt.Sprintf(" %d%% (%d/%d)", percentage, current, total)

	return bar + label
}

func getAnimatedCursor(frame int) string {
	cursors := []string{"▌", "▐", "▌", "▐"}
	return cursors[frame%len(cursors)]
}