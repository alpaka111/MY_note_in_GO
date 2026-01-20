package main

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// === SPLASH SCREEN ===
func (m model) viewSplash() string {
	animation := []string{
		"🦙",
		"🦙✨",
		"✨🦙✨",
		"✨🦙✨💫",
		"💫✨🦙✨💫",
		"✨💫🦙💫✨",
		"💫✨🦙✨💫",
		"✨🦙✨💫",
		"🦙✨",
		"🦙",
	}

	frame := m.splashTicks / 2
	if frame >= len(animation) {
		frame = len(animation) - 1
	}

	title := `
    ╔═══════════════════════════════════════════════╗
    ║                                               ║
    ║          🦙 ALPAKA NOTES v2.0 🦙             ║
    ║       Ultimate TUI Experience Edition        ║
    ║                                               ║
    ║         Secure • Beautiful • Fast            ║
    ║                                               ║
    ╚═══════════════════════════════════════════════╝
    `

	gradientTitle := renderGradientText(title, colorPalette)
	loadingBar := renderProgressBar(m.splashTicks, 20, 40)

	var b strings.Builder
	b.WriteString(gradientTitle)
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(80).
		Render(animation[frame]))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(80).
		Render(loadingBar))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(muted).
		Italic(true).
		Align(lipgloss.Center).
		Width(80).
		Render("Preparing environment..."))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

// === LOGIN SCREEN ===
func (m model) updateLogin(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if len(m.passwordBuf) == 0 {
			m.err = fmt.Errorf("password cannot be empty")
			return m, nil
		}

		notebook, err := LoadNotebook(m.filename, m.passwordBuf)
		if err != nil {
			notebook = NewNotebook(m.filename, m.passwordBuf)
			m.success = "New notebook created!"
		} else {
			m.success = fmt.Sprintf("Loaded %d notes!", len(notebook.Notes))
		}

		m.notebook = notebook
		m.password = m.passwordBuf
		m.passwordBuf = ""
		m.screen = screenMenu
		return m, nil

	case "ctrl+h":
		m.showPassword = !m.showPassword
	case "backspace":
		if len(m.passwordBuf) > 0 {
			m.passwordBuf = m.passwordBuf[:len(m.passwordBuf)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.passwordBuf += msg.String()
		}
	}
	return m, nil
}

func (m model) viewLogin() string {
	var b strings.Builder

	b.WriteString(renderHeader("WELCOME TO ALPAKA NOTES", "Your private, encrypted notebook"))
	b.WriteString("\n\n")

	infoCard := glowBoxStyle.
		Width(70).
		Render(infoStyle.Render(
			"ℹ First launch? Set a new password.\nReturning user? Enter your password.",
		))
	b.WriteString(infoCard)
	b.WriteString("\n\n")

	passwordLabel := focusedLabelStyle.Render("🔐 Password:")
	b.WriteString(passwordLabel)
	b.WriteString("\n")

	var passwordDisplay string
	if m.showPassword {
		passwordDisplay = m.passwordBuf
	} else {
		passwordDisplay = strings.Repeat("●", len(m.passwordBuf))
	}

	if len(m.passwordBuf) == 0 {
		passwordDisplay = lipgloss.NewStyle().
			Foreground(muted).
			Render("Enter password...")
	}
	passwordDisplay += getAnimatedCursor(m.animFrame)

	passwordBox := focusedBoxStyle.Width(70).Render(passwordDisplay)
	b.WriteString(passwordBox)
	b.WriteString("\n")

	toggleHint := lipgloss.NewStyle().
		Foreground(muted).
		Italic(true).
		Render(fmt.Sprintf(
			"Ctrl+H - %s password",
			map[bool]string{true: "hide", false: "show"}[m.showPassword],
		))
	b.WriteString(toggleHint)
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render(m.err.Error()))
		b.WriteString("\n\n")
	}

	securityInfo := boxStyle.
		Width(70).
		BorderForeground(success).
		Render(
			"🔒 Your data is protected with AES encryption\n" +
				"🔐 Password is never stored\n" +
				"✅ .alpaka format — for your eyes only",
		)
	b.WriteString(securityInfo)

	b.WriteString(renderFooter(renderHelp(
		"Enter", "Login",
		"Ctrl+H", "Show/Hide",
		"Ctrl+C", "Quit",
	)))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

// === MENU SCREEN ===
func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < 6 {
			m.cursor++
		}
	case "enter":
		m.err = nil
		m.success = ""
		switch m.cursor {
		case 0:
			m.screen = screenAddNote
			m.titleBuf = ""
			m.contentBuf = ""
			m.tagsBuf = ""
			m.cursor = 0
		case 1:
			m.screen = screenViewNotes
			m.selected = 0
			m.scrollOffset = 0
		case 2:
			m.screen = screenSearch
			m.searchQuery = ""
		case 3:
			m.screen = screenStats
		case 4:
			m.screen = screenSettings
		case 5:
			if err := m.notebook.Save(); err != nil {
				m.err = err
			} else {
				m.success = "Saved successfully!"
			}
		case 6:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) viewMenu() string {
	var b strings.Builder

	b.WriteString(renderHeader("Main Menu", "Cheack out your options below"))
	b.WriteString("\n")

	// Stats dashboard
	statsRow := lipgloss.JoinHorizontal(lipgloss.Top,
		statCardStyle.Render(
			statNumberStyle.Render(fmt.Sprintf("%d", len(m.notebook.Notes)))+"\n"+
				statLabelStyle.Render("📝 Notes")),
		statCardStyle.Render(
			statNumberStyle.Render(fmt.Sprintf("%d", m.notebook.CountTags()))+"\n"+
				statLabelStyle.Render("🏷️  Tags")),
		statCardStyle.Render(
			statNumberStyle.Render(fmt.Sprintf("%d", m.notebook.CountWords()))+"\n"+
				statLabelStyle.Render("📊 Words")),
	)
	b.WriteString(lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(80).
		Render(statsRow))
	b.WriteString("\n\n")

	// File info
	fileInfo := boxStyle.
		Width(70).
		BorderForeground(accent).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("📁 File: %s │ 🔐 Encrypted", m.filename))
	b.WriteString(fileInfo)
	b.WriteString("\n\n")

	// Menu items with enhanced icons
	menuItems := []struct {
		icon string
		text string
		desc string
	}{
		{"📝", "New Note", "Create a new note"},
		{"📖", "View Notes", "Browse all notes"},
		{"🔍", "Search", "Find specific notes"},
		{"📊", "Statistics", "Analyze and visualize data"},
		{"⚙️ ", "Settings", "Sorting and viewing options"},
		{"💾", "Save", "Save changes to disk"},
		{"🚪", "Exit", "Close the program"},
	}

	for i, item := range menuItems {
		itemText := fmt.Sprintf("%s  %s", item.icon, item.text)
		itemDesc := lipgloss.NewStyle().Foreground(muted).Render(" - " + item.desc)

		if m.cursor == i {
			b.WriteString(selectedMenuStyle.Render("▶ "+itemText) + itemDesc)
		} else {
			b.WriteString(menuItemStyle.Render("  "+itemText) + itemDesc)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Status messages
	if m.success != "" {
		b.WriteString(successStyle.Render("✓ " + m.success))
		b.WriteString("\n")
	}
	if m.err != nil {
		b.WriteString(errorStyle.Render("✗ " + m.err.Error()))
		b.WriteString("\n")
	}

	b.WriteString(renderFooter(renderHelp(
		"↑/↓", "Navigate",
		"j/k", "Vim",
		"Enter", "Select",
		"q", "Quit",
	)))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

// === ADD NOTE SCREEN ===
func (m model) updateAddNote(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+s":
		if len(m.titleBuf) == 0 {
			m.err = fmt.Errorf("Title cannot be empty")
			return m, nil
		}

		tags := []string{}
		if len(m.tagsBuf) > 0 {
			tags = strings.Fields(m.tagsBuf)
		}

		note := NewNote(m.titleBuf, m.contentBuf, tags)
		m.notebook.AddNote(note)

		m.screen = screenMenu
		m.success = "Note added successfully!"
		m.cursor = 0
		return m, nil

	case "tab":
		m.cursor = (m.cursor + 1) % 3
	case "shift+tab":
		m.cursor = (m.cursor - 1 + 3) % 3
	case "backspace":
		switch m.cursor {
		case 0:
			if len(m.titleBuf) > 0 {
				m.titleBuf = m.titleBuf[:len(m.titleBuf)-1]
			}
		case 1:
			if len(m.contentBuf) > 0 {
				m.contentBuf = m.contentBuf[:len(m.contentBuf)-1]
			}
		case 2:
			if len(m.tagsBuf) > 0 {
				m.tagsBuf = m.tagsBuf[:len(m.tagsBuf)-1]
			}
		}
	case "enter":
		if m.cursor == 1 {
			m.contentBuf += "\n"
		}
	default:
		if len(msg.String()) == 1 || msg.String() == "space" {
			char := msg.String()
			if char == "space" {
				char = " "
			}
			switch m.cursor {
			case 0:
				if len(m.titleBuf) < 100 {
					m.titleBuf += char
				}
			case 1:
				if len(m.contentBuf) < 10000 {
					m.contentBuf += char
				}
			case 2:
				if len(m.tagsBuf) < 200 {
					m.tagsBuf += char
				}
			}
		}
	}
	return m, nil
}

func (m model) viewAddNote() string {
	var b strings.Builder

	b.WriteString(renderHeader("NEW NOTE", "Share your thoughts"))
	b.WriteString("\n")

	// Character counters
	titleCounter := fmt.Sprintf("%d/100", len(m.titleBuf))
	contentCounter := fmt.Sprintf("%d/10000", len(m.contentBuf))
	tagsCounter := fmt.Sprintf("%d/200", len(m.tagsBuf))

	// Title field
	titleLabel := labelStyle.Render("📌 Title:")
	if m.cursor == 0 {
		titleLabel = focusedLabelStyle.Render("📌 Title:")
	}
	titleLabel += lipgloss.NewStyle().Foreground(muted).Render(" " + titleCounter)
	b.WriteString(titleLabel)
	b.WriteString("\n")

	titleContent := m.titleBuf
	if len(titleContent) == 0 && m.cursor != 0 {
		titleContent = lipgloss.NewStyle().Foreground(muted).Render("Enter a short, descriptive title...")
	}
	if m.cursor == 0 {
		titleContent += getAnimatedCursor(m.animFrame)
	}

	var titleBox string
	if m.cursor == 0 {
		titleBox = focusedBoxStyle.Width(70).Render(titleContent)
	} else {
		titleBox = boxStyle.Width(70).Render(titleContent)
	}
	b.WriteString(titleBox)
	b.WriteString("\n")

	// Content field
	contentLabel := labelStyle.Render("📄 Content:")
	if m.cursor == 1 {
		contentLabel = focusedLabelStyle.Render("📄 Content:")
	}
	contentLabel += lipgloss.NewStyle().Foreground(muted).Render(" " + contentCounter)
	b.WriteString(contentLabel)
	b.WriteString("\n")

	contentContent := m.contentBuf
	if len(contentContent) == 0 && m.cursor != 1 {
		contentContent = lipgloss.NewStyle().Foreground(muted).Render("Write your thoughts, ideas, memories...")
	}
	if m.cursor == 1 {
		contentContent += getAnimatedCursor(m.animFrame)
	}

	var contentBox string
	if m.cursor == 1 {
		contentBox = focusedBoxStyle.Width(70).Height(10).Render(contentContent)
	} else {
		contentBox = boxStyle.Width(70).Height(10).Render(contentContent)
	}
	b.WriteString(contentBox)
	b.WriteString("\n")

	// Tags field
	tagsLabel := labelStyle.Render("🏷️  Tags:")
	if m.cursor == 2 {
		tagsLabel = focusedLabelStyle.Render("🏷️  Tags:")
	}
	tagsLabel += lipgloss.NewStyle().Foreground(muted).Render(" " + tagsCounter)
	b.WriteString(tagsLabel)
	b.WriteString("\n")

	tagsContent := m.tagsBuf
	if len(tagsContent) == 0 && m.cursor != 2 {
		tagsContent = lipgloss.NewStyle().Foreground(muted).Render("work personal important idea...")
	}
	if m.cursor == 2 {
		tagsContent += getAnimatedCursor(m.animFrame)
	}

	var tagsBox string
	if m.cursor == 2 {
		tagsBox = focusedBoxStyle.Width(70).Render(tagsContent)
	} else {
		tagsBox = boxStyle.Width(70).Render(tagsContent)
	}
	b.WriteString(tagsBox)
	b.WriteString("\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render("✗ " + m.err.Error()))
		b.WriteString("\n")
	}

	b.WriteString(renderFooter(renderHelp(
		"Tab", "Next",
		"Enter", "New line",
		"Ctrl+S", "Save",
		"Esc", "Cancel",
	)))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Top,
		b.String())
}

// === VIEW NOTES SCREEN ===
func (m model) updateViewNotes(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
	case "down", "j":
		if m.selected < len(m.notebook.Notes)-1 {
			m.selected++
		}
	case "d":
		if len(m.notebook.Notes) > 0 {
			m.notebook.DeleteNote(m.selected)
			if m.selected >= len(m.notebook.Notes) && m.selected > 0 {
				m.selected--
			}
			m.success = "Note deleted"
		}
	case "v":
		m.viewMode = (m.viewMode + 1) % 3
	case "s":
		m.sortMode = (m.sortMode + 1) % 3
		m.notebook.SortNotes(m.sortMode)
	}
	return m, nil
}

func (m model) viewViewNotes() string {
	var b strings.Builder

	sortModeText := map[sortMode]string{
		sortByDate:  "Date",
		sortByTitle: "Title",
		sortByTags:  "Tags",
	}[m.sortMode]

	viewModeText := map[int]string{
		0: "List",
		1: "Grid",
		2: "Details",
	}[m.viewMode]

	b.WriteString(renderHeader("VIEW NOTES",
		fmt.Sprintf("Sortowanie: %s │ Widok: %s", sortModeText, viewModeText)))
	b.WriteString("\n")

	if len(m.notebook.Notes) == 0 {
		emptyCard := glowBoxStyle.
			Width(70).
			Align(lipgloss.Center).
			Render("📭 No notes\n\n✨ Add your first note to get started!\n\nPress Esc and select 'New Note'")
		b.WriteString(emptyCard)
	} else {
		// View modes
		notes := m.notebook.GetSortedNotes(m.sortMode)

		switch m.viewMode {
		case 0: // List view
			for i, note := range notes {
				b.WriteString(m.renderNoteCard(note, i == m.selected, false))
			}
		case 1: // Grid view
			for i := 0; i < len(notes); i += 2 {
				left := m.renderNoteCard(notes[i], i == m.selected, true)
				right := ""
				if i+1 < len(notes) {
					right = m.renderNoteCard(notes[i+1], i+1 == m.selected, true)
				}
				row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
				b.WriteString(row)
				b.WriteString("\n")
			}
		case 2: // Detailed view
			if m.selected < len(notes) {
				b.WriteString(m.renderNoteDetailed(notes[m.selected]))
			}
		}
	}

	if m.success != "" {
		b.WriteString("\n")
		b.WriteString(successStyle.Render("✓ " + m.success))
	}

	b.WriteString(renderFooter(renderHelp(
		"↑/↓", "Navigate",
		"d", "Delete",
		"v", "Change view",
		"s", "Sort",
		"Esc", "Back",
	)))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Top,
		b.String())
}

func (m model) renderNoteCard(note Note, selected bool, compact bool) string {
	var tagsStr string
	if len(note.Tags) > 0 {
		var tagBoxes []string
		for i, tag := range note.Tags {
			style := tagStyles[i%len(tagStyles)]
			tagBoxes = append(tagBoxes, style.Render(tag))
		}
		tagsStr = strings.Join(tagBoxes, "")
	}

	title := noteTitleStyle.Render(note.Title)
	meta := noteMetaStyle.Render(fmt.Sprintf("📅 %s", note.Timestamp.Format("2006-01-02 15:04")))

	var preview string
	if !compact {
		preview = noteContentStyle.Render(truncate(note.Content, 120))
	} else {
		preview = noteContentStyle.Render(truncate(note.Content, 50))
	}

	content := fmt.Sprintf("%s\n%s\n%s\n%s", title, meta, tagsStr, preview)

	var width int
	if compact {
		width = 35
	} else {
		width = 75
	}

	if selected {
		return highlightNoteStyle.Width(width).Render(content) + "\n"
	}
	return noteCardStyle.Width(width).Render(content) + "\n"
}

func (m model) renderNoteDetailed(note Note) string {
	var tagsStr string
	if len(note.Tags) > 0 {
		var tagBoxes []string
		for i, tag := range note.Tags {
			style := tagStyles[i%len(tagStyles)]
			tagBoxes = append(tagBoxes, style.Render(tag))
		}
		tagsStr = strings.Join(tagBoxes, "")
	}

	title := lipgloss.NewStyle().
		Foreground(primary).
		Bold(true).
		Underline(true).
		Render(note.Title)

	meta := noteMetaStyle.Render(
		fmt.Sprintf("📅 %s │ 📊 %d words │ 📏 %d characters",
			note.Timestamp.Format("2006-01-02 15:04:05"),
			len(strings.Fields(note.Content)),
			len(note.Content)))

	content := noteContentStyle.Render(note.Content)

	fullContent := fmt.Sprintf("%s\n\n%s\n%s\n\n%s", title, meta, tagsStr, content)

	return highlightNoteStyle.Width(75).Render(fullContent) + "\n"
}

// === SEARCH SCREEN ===
func (m model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
	default:
		if len(msg.String()) == 1 || msg.String() == "space" {
			char := msg.String()
			if char == "space" {
				char = " "
			}
			m.searchQuery += char
		}
	}
	return m, nil
}

func (m model) viewSearch() string {
	var b strings.Builder

	b.WriteString(renderHeader("SEARCH", "Find your notes instantly"))
	b.WriteString("\n")

	// Search box
	searchLabel := focusedLabelStyle.Render("🔍 Search:")
	b.WriteString(searchLabel)
	b.WriteString("\n")

	searchContent := m.searchQuery
	if len(searchContent) == 0 {
		searchContent = lipgloss.NewStyle().Foreground(muted).Render("Type a search term...")
	}
	searchContent += getAnimatedCursor(m.animFrame)

	searchBox := focusedBoxStyle.Width(70).Render(searchContent)
	b.WriteString(searchBox)
	b.WriteString("\n\n")

	// Search results
	if len(m.searchQuery) > 0 {
		results := m.notebook.Search(m.searchQuery)

		resultHeader := lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			Render(fmt.Sprintf("🎯 Found: %d notes", len(results)))
		b.WriteString(resultHeader)
		b.WriteString("\n\n")

		if len(results) == 0 {
			noResults := boxStyle.
				Width(70).
				Align(lipgloss.Center).
				BorderForeground(warning).
				Render("😕 No matching notes found\n\nTry a different search query")
			b.WriteString(noResults)
		} else {
			for _, note := range results {
				b.WriteString(m.renderNoteCard(note, false, false))
			}
		}
	} else {
		helpText := infoStyle.Render("💡 Type anything to start searching\n\nSearch includes titles, content and tags")
		helpBox := boxStyle.Width(70).Render(helpText)
		b.WriteString(helpBox)
	}

	b.WriteString(renderFooter(renderHelp(
		"Type", "Search",
		"Esc", "Back",
	)))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Top,
		b.String())
}

// === STATS SCREEN ===
func (m model) updateStats(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) viewStats() string {
	var b strings.Builder

	b.WriteString(renderHeader("STATISTICS", "Analyze your notebook"))
	b.WriteString("\n\n")

	// Main stats
	totalNotes := len(m.notebook.Notes)
	totalWords := m.notebook.CountWords()
	totalTags := m.notebook.CountTags()
	avgWordsPerNote := 0
	if totalNotes > 0 {
		avgWordsPerNote = totalWords / totalNotes
	}

	statsGrid := lipgloss.JoinHorizontal(lipgloss.Top,
		statCardStyle.
			BorderForeground(primary).
			Render(statNumberStyle.Render(fmt.Sprintf("%d", totalNotes))+"\n"+
				statLabelStyle.Render("📝 Notes")),
		statCardStyle.
			BorderForeground(secondary).
			Render(statNumberStyle.Render(fmt.Sprintf("%d", totalWords))+"\n"+
				statLabelStyle.Render("📊 Words")),
		statCardStyle.
			BorderForeground(accent).
			Render(statNumberStyle.Render(fmt.Sprintf("%d", totalTags))+"\n"+
				statLabelStyle.Render("🏷️  Tags")),
		statCardStyle.
			BorderForeground(success).
			Render(statNumberStyle.Render(fmt.Sprintf("%d", avgWordsPerNote))+"\n"+
				statLabelStyle.Render("📈 Avg. words")),
	)

	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(80).Render(statsGrid))
	b.WriteString("\n\n")

	// Tag cloud
	tagCloud := m.notebook.GetTagCloud()
	if len(tagCloud) > 0 {
		tagCloudTitle := lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			Render("🏷️  Most popular tags:")
		b.WriteString(tagCloudTitle)
		b.WriteString("\n\n")

		var tagList []string
		for tag, count := range tagCloud {
			tagStr := tagStyles[len(tagList)%len(tagStyles)].
				Render(fmt.Sprintf("%s (%d)", tag, count))
			tagList = append(tagList, tagStr)
		}

		// Sort by count
		sort.Slice(tagList, func(i, j int) bool {
			return i < j
		})

		tagDisplay := strings.Join(tagList, " ")
		tagBox := boxStyle.Width(75).Render(tagDisplay)
		b.WriteString(tagBox)
		b.WriteString("\n\n")
	}

	// Recent activity
	if totalNotes > 0 {
		recentTitle := lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			Render("📅 Recent activity:")
		b.WriteString(recentTitle)
		b.WriteString("\n\n")

		recent := m.notebook.GetRecentNotes(5)
		for _, note := range recent {
			recentItem := lipgloss.NewStyle().
				Foreground(textDim).
				Render(fmt.Sprintf("• %s - %s",
					note.Timestamp.Format("2006-01-02"),
					truncate(note.Title, 40)))
			b.WriteString(recentItem)
			b.WriteString("\n")
		}
	}

	b.WriteString(renderFooter(renderHelp(
		"Esc", "Back to menu",
	)))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Top,
		b.String())
}

// === SETTINGS SCREEN ===
func (m model) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < 2 {
			m.cursor++
		}
	case "enter", "space":
		switch m.cursor {
		case 0:
			m.sortMode = (m.sortMode + 1) % 3
		case 1:
			m.viewMode = (m.viewMode + 1) % 3
		}
	}
	return m, nil
}

func (m model) viewSettings() string {
	var b strings.Builder

	b.WriteString(renderHeader("SETTINGS", "Customize appearance and behavior"))
	b.WriteString("\n\n")

	// Settings options
	sortModeText := map[sortMode]string{
		sortByDate:  "Date (newest)",
		sortByTitle: "Title (A-Z)",
		sortByTags:  "Tags",
	}[m.sortMode]

	viewModeText := map[int]string{
		0: "List (detailed)",
		1: "Grid (compact)",
		2: "Preview (single note)",
	}[m.viewMode]

	settings := []struct {
		icon  string
		name  string
		value string
	}{
		{"📊", "Sorting", sortModeText},
		{"👁️ ", "Note view", viewModeText},
		{"💾", "File format", ".alpaka (encrypted)"},
	}

	for i, setting := range settings {
		var settingBox string
		content := fmt.Sprintf("%s %s\n%s",
			setting.icon,
			lipgloss.NewStyle().Foreground(primary).Bold(true).Render(setting.name),
			lipgloss.NewStyle().Foreground(textDim).Render("► "+setting.value))

		if m.cursor == i && i < 2 {
			settingBox = selectedNoteStyle.Width(70).Render(content)
		} else {
			settingBox = noteCardStyle.Width(70).Render(content)
		}

		b.WriteString(settingBox)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	hint := infoStyle.Render("💡 Select an option to change the setting")
	b.WriteString(boxStyle.Width(70).Render(hint))

	b.WriteString(renderFooter(renderHelp(
		"↑/↓", "Navigate",
		"Enter/Space", "Change",
		"Esc", "Back",
	)))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

// === HELPERS ===
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
