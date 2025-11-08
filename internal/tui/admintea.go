package tui

import (
	"fmt"
	"strconv"
	"strings"

	"catv/internal/store"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AdminModel provides CRUD management of flashcards.
type adminView int

const (
	adminList adminView = iota
	adminCreate
	adminEdit
	adminConfirmDelete
	adminConfirmBulkReset
)

// Key constants
const (
	keyEsc   = "esc"
	keyEnter = "enter"
	keyTab   = "tab"
)

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	PageUp    key.Binding
	PageDown  key.Binding
	Create    key.Binding
	Edit      key.Binding
	Delete    key.Binding
	BulkReset key.Binding
	Reload    key.Binding
	Help      key.Binding
	Quit      key.Binding
	Cancel    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Cancel}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},                 // first column
		{k.Create, k.Edit, k.Delete, k.BulkReset, k.Reload}, // second column
	}
}

var keys = keyMap{
	Up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
	Down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
	PageUp:    key.NewBinding(key.WithKeys("pgup", "pageup"), key.WithHelp("pgup", "page up")),
	PageDown:  key.NewBinding(key.WithKeys("pgdown", "pagedown"), key.WithHelp("pgdown", "page down")),
	Create:    key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "create")),
	Edit:      key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Delete:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	BulkReset: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "bulk reset")),
	Reload:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reload")),
	Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Cancel:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
}

type AdminModel struct {
	flashcards []store.Flashcard
	selected   int
	view       adminView
	width      int
	height     int
	help       help.Model
	keys       keyMap

	// table for list view
	table table.Model

	// form fields
	questionInput textinput.Model
	answerInput   textinput.Model
	revisitInput  textinput.Model // days until next review

	statusMsg string
	errMsg    string

	storeRef *store.Store
}

func NewAdminModel(storeRef *store.Store, flashcards []store.Flashcard) *AdminModel {
	q := textinput.New()
	q.Placeholder = "Question"
	q.Focus()
	a := textinput.New()
	a.Placeholder = "Answer"
	r := textinput.New()
	r.Placeholder = "Days (e.g. 7)"

	// Setup table columns
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Question", Width: 40},
		{Title: "Answer", Width: 30},
		{Title: "Revisit In", Width: 12},
	}

	// Convert flashcards to table rows
	rows := makeTableRows(flashcards, columns)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("205"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return &AdminModel{
		flashcards:    flashcards,
		selected:      0,
		view:          adminList,
		table:         t,
		questionInput: q,
		answerInput:   a,
		revisitInput:  r,
		storeRef:      storeRef,
		help:          help.New(),
		keys:          keys,
	}
}

// makeTableRows converts flashcards to table rows
func makeTableRows(flashcards []store.Flashcard, columns []table.Column) []table.Row {
	rows := make([]table.Row, len(flashcards))
	for i, fc := range flashcards {
		rows[i] = table.Row{
			fmt.Sprintf("%d", fc.ID),
			truncate(fc.Question, columns[1].Width),
			truncate(fc.Answer, columns[2].Width),
			fmt.Sprintf("%d days", fc.RevisitIn),
		}
	}
	return rows
}

func (m *AdminModel) Init() tea.Cmd { return nil }

func (m *AdminModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Recalculate column widths
		// Clamp ID and RevisitIn columns to a max width
		idWidth := 6
		revisitInWidth := 12
		if m.width < 50 {
			idWidth = 3
			revisitInWidth = 8
		}
		totalFixedWidth := idWidth + revisitInWidth + 6 // Account for padding
		remainingWidth := m.width - totalFixedWidth
		questionWidth := int(float64(remainingWidth) * 0.6)
		answerWidth := remainingWidth - questionWidth

		m.table.SetColumns([]table.Column{
			{Title: "ID", Width: idWidth},
			{Title: "Question", Width: questionWidth},
			{Title: "Answer", Width: answerWidth},
			{Title: "Revisit In", Width: revisitInWidth},
		})
		m.table.SetHeight(m.height - 10) // Adjust height based on terminal
		return m, nil

	case tea.KeyMsg:
		switch m.view {
	case adminList:
		return m.handleListView(msg)
	case adminCreate:
		return m.handleCreateView(msg)
	case adminEdit:
		return m.handleEditView(msg)
	case adminConfirmDelete:
		return m.handleDeleteConfirm(msg)
	case adminConfirmBulkReset:
		return m.handleBulkResetConfirm(msg)
		}
	}
	return m, nil
}

func (m *AdminModel) handleListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Create):
		m.view = adminCreate
		m.resetForm()
		return m, nil
	case key.Matches(msg, m.keys.Edit):
		if len(m.flashcards) == 0 {
			return m, nil
		}
		m.selected = m.table.Cursor()
		m.loadSelectedIntoForm()
		m.view = adminEdit
		return m, nil
	case key.Matches(msg, m.keys.Delete):
		if len(m.flashcards) == 0 {
			return m, nil
		}
		m.selected = m.table.Cursor()
		m.view = adminConfirmDelete
		return m, nil
	case key.Matches(msg, m.keys.BulkReset):
		if len(m.flashcards) == 0 {
			return m, nil
		}
		m.view = adminConfirmBulkReset
		return m, nil
	case key.Matches(msg, m.keys.Reload):
		m.reload()
		m.statusMsg = "Table refreshed"
		m.errMsg = ""
		return m, nil
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	case key.Matches(msg, m.keys.PageDown):
		newPos := m.table.Cursor() + 10
		if newPos >= len(m.flashcards) {
			newPos = len(m.flashcards) - 1
		}
		m.table.SetCursor(newPos)
		m.selected = newPos
		return m, nil
	case key.Matches(msg, m.keys.PageUp):
		newPos := m.table.Cursor() - 10
		if newPos < 0 {
			newPos = 0
		}
		m.table.SetCursor(newPos)
		m.selected = newPos
		return m, nil
	default:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		m.selected = m.table.Cursor()
		return m, cmd
	}
}

func (m *AdminModel) handleCreateView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.view = adminList
		return m, nil
	}
	if msg.String() == keyTab {
		m.cycleFocus()
		return m, nil
	}
	if msg.String() == keyEnter {
		m.createFlashcard()
		return m, nil
	}
	m.updateInputs(msg)
	return m, nil
}

func (m *AdminModel) handleEditView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.view = adminList
		return m, nil
	}
	if msg.String() == keyTab {
		m.cycleFocus()
		return m, nil
	}
	if msg.String() == keyEnter {
		m.updateFlashcard()
		return m, nil
	}
	m.updateInputs(msg)
	return m, nil
}

func (m *AdminModel) handleDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.view = adminList
		return m, nil
	}
	if msg.String() == "y" {
		m.deleteFlashcard()
		return m, nil
	}
	if msg.String() == "n" {
		m.view = adminList
	}
	return m, nil
}

func (m *AdminModel) handleBulkResetConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.view = adminList
		return m, nil
	}
	if msg.String() == "y" {
		m.bulkResetRevisitIn()
		return m, nil
	}
	if msg.String() == "n" {
		m.view = adminList
	}
	return m, nil
}

func (m *AdminModel) View() string {
	// Define color styles
	helpStyle := lipgloss.NewStyle().Faint(true)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)

	// Edit mode styles with colors
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	inputFocusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Padding(0, 1)
	inputBlurredStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	helpBar := m.help.View(m.keys)

	if m.errMsg != "" {
		helpBar += "\n" + errorStyle.Render("✗ "+m.errMsg)
	} else if m.statusMsg != "" {
		helpBar += "\n" + successStyle.Render("✓ "+m.statusMsg)
	}

	var mainContent string

	switch m.view {
	case adminList:
		var b strings.Builder

		if len(m.flashcards) == 0 {
			b.WriteString("No flashcards. Press 'c' to create.\n")
		} else {
			b.WriteString(m.table.View() + "\n")
		}
		b.WriteString("\n" + helpBar)
		mainContent = b.String()

	case adminCreate:
		qLabel := labelStyle.Render("Question:")
		qInput := m.questionInput.View()
		if m.questionInput.Focused() {
			qInput = inputFocusedStyle.Render(qInput)
		} else {
			qInput = inputBlurredStyle.Render(qInput)
		}

		aLabel := labelStyle.Render("Answer:")
		aInput := m.answerInput.View()
		if m.answerInput.Focused() {
			aInput = inputFocusedStyle.Render(aInput)
		} else {
			aInput = inputBlurredStyle.Render(aInput)
		}

		rLabel := labelStyle.Render("Revisit (days):")
		rInput := m.revisitInput.View()
		if m.revisitInput.Focused() {
			rInput = inputFocusedStyle.Render(rInput)
		} else {
			rInput = inputBlurredStyle.Render(rInput)
		}

		helpText := helpStyle.Render("[tab] next field  [enter] save  [esc] cancel")

		mainContent = fmt.Sprintf("%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n\n",
			qLabel, qInput, aLabel, aInput, rLabel, rInput, helpText)

	case adminEdit:
		title := titleStyle.Render(fmt.Sprintf("(ID %d)", m.flashcards[m.selected].ID))

		qLabel := labelStyle.Render("Question:")
		qInput := m.questionInput.View()
		if m.questionInput.Focused() {
			qInput = inputFocusedStyle.Render(qInput)
		} else {
			qInput = inputBlurredStyle.Render(qInput)
		}

		aLabel := labelStyle.Render("Answer:")
		aInput := m.answerInput.View()
		if m.answerInput.Focused() {
			aInput = inputFocusedStyle.Render(aInput)
		} else {
			aInput = inputBlurredStyle.Render(aInput)
		}

		rLabel := labelStyle.Render("Revisit (days):")
		rInput := m.revisitInput.View()
		if m.revisitInput.Focused() {
			rInput = inputFocusedStyle.Render(rInput)
		} else {
			rInput = inputBlurredStyle.Render(rInput)
		}

		helpText := helpStyle.Render("[tab] next field  [enter] update  [esc] cancel")

		mainContent = fmt.Sprintf("%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n\n",
			title, qLabel, qInput, aLabel, aInput, rLabel, rInput, helpText)

	case adminConfirmDelete:
		warning := errorStyle.Render(fmt.Sprintf("Delete Flashcard ID %d?", m.flashcards[m.selected].ID))
		options := helpStyle.Render("[y] yes  [n] no  [esc] cancel")
		mainContent = fmt.Sprintf("%s\n\n%s\n\n\n\n", warning, options)

	case adminConfirmBulkReset:
		title := titleStyle.Render("⚡ Bulk Reset RevisitIn")
		warning := errorStyle.Render(fmt.Sprintf("Set RevisitIn to 0 for ALL %d flashcards?", len(m.flashcards)))
		info := infoStyle.Render("This will make all flashcards due for immediate review.")
		options := helpStyle.Render("[y] yes  [n] no  [esc] cancel")
		mainContent = fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n\n%s", title, warning, info, options, helpBar)
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, mainContent)
}

// helpers
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func (m *AdminModel) resetForm() {
	m.questionInput.SetValue("")
	m.answerInput.SetValue("")
	m.revisitInput.SetValue("")
	m.errMsg = ""
	m.statusMsg = ""
	m.questionInput.Focus()
	m.answerInput.Blur()
	m.revisitInput.Blur()
}

func (m *AdminModel) loadSelectedIntoForm() {
	fc := m.flashcards[m.selected]
	m.questionInput.SetValue(fc.Question)
	m.answerInput.SetValue(fc.Answer)
	m.revisitInput.SetValue(strconv.Itoa(fc.RevisitIn))
	m.questionInput.Focus()
	m.answerInput.Blur()
	m.revisitInput.Blur()
}

func (m *AdminModel) updateInputs(msg tea.Msg) {
	// Update the focused input with the key message
	switch {
	case m.questionInput.Focused():
		m.questionInput, _ = m.questionInput.Update(msg)
	case m.answerInput.Focused():
		m.answerInput, _ = m.answerInput.Update(msg)
	case m.revisitInput.Focused():
		m.revisitInput, _ = m.revisitInput.Update(msg)
	}
}

func (m *AdminModel) parseRevisitDays() (int, error) {
	val := strings.TrimSpace(m.revisitInput.Value())
	if val == "" {
		return 0, fmt.Errorf("revisit days required")
	}
	d, err := strconv.Atoi(val)
	if err != nil || d < 0 {
		return 0, fmt.Errorf("invalid days")
	}
	return d, nil
}

func (m *AdminModel) createFlashcard() {
	days, err := m.parseRevisitDays()
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	fc := store.Flashcard{Question: m.questionInput.Value(), Answer: m.answerInput.Value(), File: "manual", RevisitIn: days}
	if err := m.storeRef.InsertFlashcard(fc); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.statusMsg = "Flashcard created"
	m.reload()
	m.view = adminList
}

func (m *AdminModel) updateFlashcard() {
	days, err := m.parseRevisitDays()
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	fc := m.flashcards[m.selected]
	fc.Question = m.questionInput.Value()
	fc.Answer = m.answerInput.Value()
	fc.RevisitIn = days
	if err := m.storeRef.UpdateFlashcardFull(fc); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.statusMsg = "Flashcard updated"
	m.reload()
	m.view = adminList
}

func (m *AdminModel) deleteFlashcard() {
	id := m.flashcards[m.selected].ID
	if err := m.storeRef.DeleteFlashcard(id); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.statusMsg = fmt.Sprintf("Deleted flashcard %d", id)
	m.reload()
	if m.selected >= len(m.flashcards) {
		m.selected = len(m.flashcards) - 1
	}
	m.view = adminList
}

func (m *AdminModel) reload() {
	list, err := m.storeRef.GetAllFlashcards()
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	m.flashcards = list

	// Update table with new data
	rows := makeTableRows(m.flashcards, m.table.Columns())
	m.table.SetRows(rows)

	// Adjust cursor if needed
	if m.selected >= len(m.flashcards) && len(m.flashcards) > 0 {
		m.selected = len(m.flashcards) - 1
		m.table.SetCursor(m.selected)
	}
	if m.selected < 0 {
		m.selected = 0
		m.table.SetCursor(0)
	}
}

// bulkResetRevisitIn resets RevisitIn to 0 for all flashcards
func (m *AdminModel) bulkResetRevisitIn() {
	count := 0
	for _, fc := range m.flashcards {
		fc.RevisitIn = 0
		if err := m.storeRef.UpdateFlashcard(fc); err != nil {
			m.errMsg = fmt.Sprintf("Error updating flashcard %d: %v", fc.ID, err)
			m.view = adminList
			return
		}
		count++
	}
	m.statusMsg = fmt.Sprintf("Reset RevisitIn to 0 for %d flashcards", count)
	m.errMsg = ""
	m.reload()
	m.view = adminList
}

// cycleFocus switches focus Question -> Answer -> Revisit -> Question
func (m *AdminModel) cycleFocus() {
	if m.questionInput.Focused() {
		m.questionInput.Blur()
		m.answerInput.Focus()
		m.revisitInput.Blur()
		return
	}
	if m.answerInput.Focused() {
		m.answerInput.Blur()
		m.revisitInput.Focus()
		m.questionInput.Blur()
		return
	}
	// default or revisit focused -> go to question
	m.revisitInput.Blur()
	m.questionInput.Focus()
	m.answerInput.Blur()
}
