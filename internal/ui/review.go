package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaosaffran/mob/internal/llm"
)

// ChecklistItem represents an item in the review checklist
type ChecklistItem struct {
	Description string
}

// ReviewModel is the Bubble Tea model for the review UI
type ReviewModel struct {
	diff            string
	rawDiff         string // unhighlighted diff for LLM
	diffStat        string
	checklistItems  []ChecklistItem
	checked         map[int]bool
	cursor          int
	recsCursor      int // cursor for recommendations panel
	viewport        viewport.Model
	focusedPanel    string // "checklist", "diff", or "recommendations"
	ready           bool
	width           int
	height          int
	allChecked      bool
	issue           string
	sidebarWidth    int
	recommendations []llm.Recommendation
	loadingRecs     bool
	recsError       string
	showModal       bool   // whether to show the recommendation modal
	modalContent    string // content to display in the modal
	modalTitle      string // title for the modal
}

// NewReviewModel creates a new review model
func NewReviewModel(diff, diffStat, issue string, items []ChecklistItem) ReviewModel {
	return ReviewModel{
		diff:           highlightDiff(diff),
		rawDiff:        diff,
		diffStat:       diffStat,
		checklistItems: items,
		checked:        make(map[int]bool),
		cursor:         0,
		focusedPanel:   "checklist",
		issue:          issue,
		loadingRecs:    true,
	}
}

// recommendationsMsg is sent when recommendations are loaded
type recommendationsMsg struct {
	recommendations []llm.Recommendation
	err             error
}

// loadRecommendations fetches recommendations from LLM
func loadRecommendations(diff string) tea.Cmd {
	return func() tea.Msg {
		recs, err := llm.GetRecommendations(diff)
		return recommendationsMsg{recommendations: recs, err: err}
	}
}

// highlightDiff applies syntax highlighting to git diff output
func highlightDiff(diff string) string {
	var result strings.Builder
	lines := strings.Split(diff, "\n")

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			// File headers
			result.WriteString(StyleDiffHeader.Render(line))
		case strings.HasPrefix(line, "@@"):
			// Hunk headers (@@ -1,3 +1,4 @@)
			result.WriteString(StyleDiffHunk.Render(line))
		case strings.HasPrefix(line, "+"):
			// Added lines
			result.WriteString(StyleDiffAdded.Render(line))
		case strings.HasPrefix(line, "-"):
			// Removed lines
			result.WriteString(StyleDiffRemoved.Render(line))
		case strings.HasPrefix(line, "diff --git"):
			// Diff command header
			result.WriteString(StyleDiffMeta.Render(line))
		case strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "new file") || strings.HasPrefix(line, "deleted file"):
			// Index and file mode info
			result.WriteString(StyleDiffMeta.Render(line))
		default:
			// Context lines
			result.WriteString(StyleDiffContext.Render(line))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// Init implements tea.Model
func (m ReviewModel) Init() tea.Cmd {
	return loadRecommendations(m.rawDiff)
}

// Update implements tea.Model
func (m ReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case recommendationsMsg:
		m.loadingRecs = false
		if msg.err != nil {
			m.recsError = msg.err.Error()
		} else {
			m.recommendations = msg.recommendations
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate sidebar width
		m.sidebarWidth = int(float64(msg.Width) * SidePanelWidthRatio)
		if m.sidebarWidth < MinSidePanelWidth {
			m.sidebarWidth = MinSidePanelWidth
		}
		if m.sidebarWidth > MaxSidePanelWidth {
			m.sidebarWidth = MaxSidePanelWidth
		}

		// Calculate diff panel width (account for borders and gap)
		diffWidth := msg.Width - m.sidebarWidth - 5

		verticalMargin := HeaderHeight + FooterHeight

		if !m.ready {
			m.viewport = viewport.New(diffWidth, msg.Height-verticalMargin-2)
			m.viewport.SetContent(m.diff)
			m.ready = true
		} else {
			m.viewport.Width = diffWidth
			m.viewport.Height = msg.Height - verticalMargin - 2
			m.viewport.SetContent(m.diff)
		}

	case tea.KeyMsg:
		// Handle modal close first
		if m.showModal {
			switch msg.String() {
			case "enter", "esc", "q":
				m.showModal = false
			}
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			// Only quit if not in recommendations panel, otherwise close modal behavior
			if m.focusedPanel != "recommendations" {
				return m, tea.Quit
			}

		case "tab":
			// Cycle focus between panels: checklist -> diff -> recommendations -> checklist
			switch m.focusedPanel {
			case "checklist":
				m.focusedPanel = "diff"
			case "diff":
				m.focusedPanel = "recommendations"
			case "recommendations":
				m.focusedPanel = "checklist"
			}

		case "up", "k":
			switch m.focusedPanel {
			case "checklist":
				if m.cursor > 0 {
					m.cursor--
				}
			case "recommendations":
				if m.recsCursor > 0 {
					m.recsCursor--
				}
			case "diff":
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}

		case "down", "j":
			switch m.focusedPanel {
			case "checklist":
				if m.cursor < len(m.checklistItems)-1 {
					m.cursor++
				}
			case "recommendations":
				if m.recsCursor < len(m.recommendations)-1 {
					m.recsCursor++
				}
			case "diff":
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}

		case "enter":
			switch m.focusedPanel {
			case "checklist":
				m.checked[m.cursor] = !m.checked[m.cursor]
				m.allChecked = m.areAllChecked()
			case "recommendations":
				// Open modal with full recommendation
				if len(m.recommendations) > 0 && m.recsCursor < len(m.recommendations) {
					rec := m.recommendations[m.recsCursor]
					m.modalTitle = rec.Title
					m.modalContent = fmt.Sprintf("Severity: %s\n\n%s", rec.Severity, rec.Description)
					m.showModal = true
				}
			}

		case " ":
			if m.focusedPanel == "checklist" {
				m.checked[m.cursor] = !m.checked[m.cursor]
				m.allChecked = m.areAllChecked()
			}

		case "pgup", "pgdown", "home", "end":
			if m.focusedPanel == "diff" {
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}
	}

	return m, cmd
}

func (m ReviewModel) areAllChecked() bool {
	for i := range m.checklistItems {
		if !m.checked[i] {
			return false
		}
	}
	return true
}

// View implements tea.Model
func (m ReviewModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	// If modal is open, render it on top
	if m.showModal {
		return m.renderModal()
	}

	// Header
	title := StyleTitle.Render(fmt.Sprintf("Review for issue #%s", m.issue))

	// Status bar
	var statusText string
	if m.allChecked {
		statusText = StyleSuccess.Render(fmt.Sprintf("%s All items checked - Ready to update!", SymbolSuccess))
	} else {
		checked := 0
		for _, v := range m.checked {
			if v {
				checked++
			}
		}
		statusText = StyleStatus.Render(fmt.Sprintf("Checked: %d/%d", checked, len(m.checklistItems)))
	}

	// Calculate panel heights (checklist and recommendations share the right side)
	rightPanelHeight := m.height - HeaderHeight - FooterHeight - 2
	checklistHeight := rightPanelHeight / 2
	recsHeight := rightPanelHeight - checklistHeight - 3 // -3 for gap between panels

	// Build checklist panel
	checklistContent := m.renderChecklistContent()
	var checklistPanel string
	var checklistTitle string
	if m.focusedPanel == "checklist" {
		checklistTitle = StylePanelTitle.Render("Checklist")
		checklistPanel = StylePanelActive.
			Width(m.sidebarWidth - 2).
			Height(checklistHeight).
			Render(checklistContent)
	} else {
		checklistTitle = StylePanelTitleInactive.Render("Checklist")
		checklistPanel = StylePanelInactive.
			Width(m.sidebarWidth - 2).
			Height(checklistHeight).
			Render(checklistContent)
	}

	// Build recommendations panel
	recsContent := m.renderRecommendationsContent()
	var recsTitle string
	var recsPanel string
	if m.focusedPanel == "recommendations" {
		recsTitle = StylePanelTitle.Render("AI Recommendations (LLVM Standards)")
		recsPanel = StylePanelActive.
			Width(m.sidebarWidth - 2).
			Height(recsHeight).
			Render(recsContent)
	} else {
		recsTitle = StylePanelTitleInactive.Render("AI Recommendations (LLVM Standards)")
		recsPanel = StylePanelInactive.
			Width(m.sidebarWidth - 2).
			Height(recsHeight).
			Render(recsContent)
	}

	// Stack checklist and recommendations vertically
	rightSideTitles := lipgloss.JoinVertical(lipgloss.Left,
		checklistTitle,
	)
	rightSidePanels := lipgloss.JoinVertical(lipgloss.Left,
		checklistPanel,
		"",
		recsTitle,
		recsPanel,
	)

	// Build diff panel
	var diffPanel string
	var diffTitle string
	diffWidth := m.width - m.sidebarWidth - 5
	if m.focusedPanel == "diff" {
		diffTitle = StylePanelTitle.Render("Diff")
		diffPanel = StylePanelActive.
			Width(diffWidth).
			Height(rightPanelHeight).
			Render(m.viewport.View())
	} else {
		diffTitle = StylePanelTitleInactive.Render("Diff")
		diffPanel = StylePanelInactive.
			Width(diffWidth).
			Height(rightPanelHeight).
			Render(m.viewport.View())
	}

	// Combine panels side by side (diff on left, checklist+recs on right)
	titles := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(diffWidth).Render(diffTitle),
		"  ",
		rightSideTitles,
	)
	panels := lipgloss.JoinHorizontal(lipgloss.Top, diffPanel, "  ", rightSidePanels)

	// Footer
	footer := StyleStatus.Render(fmt.Sprintf("Tab: switch panel %s ↑/↓: navigate %s Space/Enter: toggle/view %s q: quit", SymbolBullet, SymbolBullet, SymbolBullet))

	return fmt.Sprintf("%s\n%s\n%s\n%s\n\n%s", title, statusText, titles, panels, footer)
}

// renderModal renders the recommendation detail modal
func (m ReviewModel) renderModal() string {
	modalWidth := m.width / 2
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > m.width-10 {
		modalWidth = m.width - 10
	}

	modalHeight := m.height / 2
	if modalHeight < 10 {
		modalHeight = 10
	}

	// Modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		MarginBottom(1)

	// Build modal content
	content := titleStyle.Render(m.modalTitle) + "\n\n" + m.modalContent + "\n\n" +
		StyleStatus.Render("Press Enter or Esc to close")

	modal := modalStyle.Render(content)

	// Center the modal
	centeredModal := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)

	return centeredModal
}

// renderRecommendationsContent renders the AI recommendations
func (m ReviewModel) renderRecommendationsContent() string {
	var sb strings.Builder

	if m.loadingRecs {
		sb.WriteString(StyleStatus.Render("Loading recommendations..."))
		return sb.String()
	}

	if m.recsError != "" {
		sb.WriteString(StyleError.Render(fmt.Sprintf("Error: %s", m.recsError)))
		return sb.String()
	}

	if len(m.recommendations) == 0 {
		sb.WriteString(StyleStatus.Render("No recommendations"))
		return sb.String()
	}

	for i, rec := range m.recommendations {
		// Cursor indicator
		cursor := SymbolNoCursor
		if m.recsCursor == i && m.focusedPanel == "recommendations" {
			cursor = SymbolCursor
		}

		// Severity indicator
		severityStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(llm.SeverityColor(rec.Severity)))
		severityIcon := "●"

		// Title with severity
		titleLine := fmt.Sprintf("%s%s %s", cursor, severityStyle.Render(severityIcon), rec.Title)
		sb.WriteString(titleLine)
		sb.WriteString("\n")

		// Description (truncate if needed)
		desc := rec.Description
		maxLen := m.sidebarWidth - 8
		if maxLen > 0 && len(desc) > maxLen {
			desc = desc[:maxLen-3] + "..."
		}
		sb.WriteString(StyleStatus.Render(fmt.Sprintf("    %s", desc)))

		if i < len(m.recommendations)-1 {
			sb.WriteString("\n\n")
		}
	}

	return sb.String()
}

func (m ReviewModel) renderChecklistContent() string {
	var sb strings.Builder

	for i, item := range m.checklistItems {
		cursor := SymbolNoCursor
		if m.cursor == i && m.focusedPanel == "checklist" {
			cursor = SymbolCursor
		}

		checked := SymbolCheckboxUnchecked
		if m.checked[i] {
			checked = SymbolCheckboxChecked
		}

		// Truncate description if too long
		desc := item.Description
		maxLen := m.sidebarWidth - 10
		if maxLen > 0 && len(desc) > maxLen {
			desc = desc[:maxLen-3] + "..."
		}

		sb.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checked, desc))
	}

	sb.WriteString("\n")
	sb.WriteString(StyleStatus.Render("─── Stats ───"))
	sb.WriteString("\n")
	sb.WriteString(m.diffStat)

	return sb.String()
}

// IsReviewComplete returns true if all items are checked
func (m ReviewModel) IsReviewComplete() bool {
	return m.allChecked
}

// RunReview starts the review UI and returns whether the review was completed
func RunReview(diff, diffStat, issue string, items []ChecklistItem) (bool, error) {
	model := NewReviewModel(diff, diffStat, issue, items)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("error running review UI: %w", err)
	}

	if m, ok := finalModel.(ReviewModel); ok {
		return m.IsReviewComplete(), nil
	}

	return false, nil
}
