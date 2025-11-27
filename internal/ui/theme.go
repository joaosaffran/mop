package ui

import "github.com/charmbracelet/lipgloss"

// Color constants for the UI theme
const (
	// Primary colors
	ColorPrimary   = lipgloss.Color("205") // Pink/Magenta
	ColorSecondary = lipgloss.Color("81")  // Cyan

	// Text colors
	ColorText       = lipgloss.Color("250") // Light gray
	ColorTextMuted  = lipgloss.Color("241") // Dimmed gray
	ColorTextBright = lipgloss.Color("255") // White

	// Background colors
	ColorBgDark      = lipgloss.Color("236") // Dark gray
	ColorBgHighlight = lipgloss.Color("238") // Slightly lighter gray

	// Diff colors
	ColorDiffAdded   = lipgloss.Color("34")  // Darker green
	ColorDiffRemoved = lipgloss.Color("196") // Red
	ColorDiffHeader  = lipgloss.Color("81")  // Cyan
	ColorDiffHunk    = lipgloss.Color("135") // Purple
	ColorDiffMeta    = lipgloss.Color("208") // Orange
	ColorDiffContext = lipgloss.Color("250") // Light gray

	// Status colors
	ColorSuccess = lipgloss.Color("42")  // Green
	ColorError   = lipgloss.Color("196") // Red
	ColorWarning = lipgloss.Color("208") // Orange
	ColorInfo    = lipgloss.Color("81")  // Cyan
)

// Layout constants
const (
	// Padding
	PaddingHorizontal = 2
	PaddingVertical   = 1

	// Margins
	MarginBottom = 1

	// Header/Footer heights for viewport calculations
	HeaderHeight = 4
	FooterHeight = 3

	// Side panel width ratio (checklist takes this fraction of the screen)
	SidePanelWidthRatio = 0.3
	MinSidePanelWidth   = 30
	MaxSidePanelWidth   = 50
)

// Symbols used throughout the UI
const (
	SymbolCheckboxChecked   = "[x]"
	SymbolCheckboxUnchecked = "[ ]"
	SymbolCursor            = "> "
	SymbolNoCursor          = "  "
	SymbolSuccess           = "✓"
	SymbolError             = "✗"
	SymbolBullet            = "•"
)

// Pre-defined styles for common UI elements
var (
	// Title style - bold primary color
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(MarginBottom)

	// Panel styles
	StylePanelActive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)

	StylePanelInactive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorTextMuted).
				Padding(0, 1)

	StylePanelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	StylePanelTitleInactive = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorTextMuted).
				Padding(0, 1)

	// Tab styles (kept for compatibility)
	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(ColorBgDark).
			Padding(0, PaddingHorizontal)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(ColorText).
				Padding(0, PaddingHorizontal)

	// Status/footer text style
	StyleStatus = lipgloss.NewStyle().
			Foreground(ColorTextMuted)

	// Diff styles
	StyleDiffAdded = lipgloss.NewStyle().
			Foreground(ColorDiffAdded)

	StyleDiffRemoved = lipgloss.NewStyle().
				Foreground(ColorDiffRemoved)

	StyleDiffHeader = lipgloss.NewStyle().
			Foreground(ColorDiffHeader).
			Bold(true)

	StyleDiffHunk = lipgloss.NewStyle().
			Foreground(ColorDiffHunk)

	StyleDiffMeta = lipgloss.NewStyle().
			Foreground(ColorDiffMeta)

	StyleDiffContext = lipgloss.NewStyle().
				Foreground(ColorDiffContext)

	// Success/Error message styles
	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	StyleInfo = lipgloss.NewStyle().
			Foreground(ColorInfo)
)
