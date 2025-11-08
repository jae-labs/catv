package tui

import (
	"catv/internal/store"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ReviewModel manages the state for the review session
// Views: question, answer, correct/incorrect, revisitIn, done

type viewState int

const (
	viewQuestion viewState = iota
	viewAnswer
	viewRevisitIn
	viewDone
)

var (
	questionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	answerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	infoStyle     = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("244"))
	successStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("34"))
)

type ReviewModel struct {
	flashcards []store.Flashcard
	current    int
	view       viewState
	resultMsg  string
	quitting   bool
	correct    []bool
	revisitIn  []int
}

func NewReviewModel(flashcards []store.Flashcard) *ReviewModel {
	return &ReviewModel{
		flashcards: flashcards,
		current:    0,
		view:       viewQuestion,
		correct:    make([]bool, len(flashcards)),
		revisitIn:  make([]int, len(flashcards)),
	}
}

func (m *ReviewModel) Init() tea.Cmd {
	return nil
}

func (m *ReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
		switch m.view {
		case viewQuestion:
			if msg.String() == "enter" {
				m.view = viewAnswer
			}
		case viewAnswer:
			if msg.String() == "c" {
				m.correct[m.current] = true
				m.view = viewRevisitIn
			} else if msg.String() == "i" {
				m.correct[m.current] = false
				m.resultMsg = "Marked incorrect. Card will not be scheduled for repetition."
				m.nextCard()
			}
		case viewRevisitIn:
			switch msg.String() {
			case "1":
				m.revisitIn[m.current] = 1
				m.resultMsg = "Revisit in 1 day"
				m.nextCard()
			case "3":
				m.revisitIn[m.current] = 3
				m.resultMsg = "Revisit in 3 days"
				m.nextCard()
			case "7":
				m.revisitIn[m.current] = 7
				m.resultMsg = "Revisit in 7 days"
				m.nextCard()
			case "9":
				m.revisitIn[m.current] = 9
				m.resultMsg = "Revisit in 9 days"
				m.nextCard()
			}
		case viewDone:
			if msg.String() == "q" {
				m.quitting = true
				return m, tea.Quit
			}
		}
	case tea.QuitMsg:
		m.quitting = true
	}
	return m, nil
}

func (m *ReviewModel) nextCard() {
	m.current++
	if m.current >= len(m.flashcards) {
		m.view = viewDone
		return
	}
	m.view = viewQuestion
	m.resultMsg = ""
}

// Results API for review command
func (m *ReviewModel) FlashcardWasCorrect(idx int) bool {
	if idx < 0 || idx >= len(m.correct) {
		return false
	}
	return m.correct[idx]
}

func (m *ReviewModel) FlashcardRevisitIn(idx int) int {
	if idx < 0 || idx >= len(m.revisitIn) {
		return 0
	}
	return m.revisitIn[idx]
}

func (m *ReviewModel) View() string {
	if m.quitting {
		return "Goodbye!"
	}
	total := len(m.flashcards)
	current := m.current + 1
	// If review is done, show total instead of current+1
	if m.view == viewDone {
		current = total
	}
	width := 60
	height := 10
	frame := lipgloss.NewStyle().Width(width).Height(height).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63"))

	// Count correct and incorrect answers
	correctCount := 0
	incorrectCount := 0
	for i := range m.flashcards {
		// Only count cards that have been answered (i.e., where correct/incorrect has been set)
		if m.view == viewDone || i < m.current {
			if m.FlashcardWasCorrect(i) {
				correctCount++
			} else {
				incorrectCount++
			}
		}
	}

	// Bottom bar: ✅ left, current/total center, ❌ right
	bottomBar := lipgloss.NewStyle().Width(width - 2).Align(lipgloss.Center).Render(
		fmt.Sprintf("%-10s%s%10s", fmt.Sprintf("✅ %d", correctCount), fmt.Sprintf("%d/%d", current, total), fmt.Sprintf("❌ %d", incorrectCount)),
	)

	exitMsg := infoStyle.Render("Press q to exit at any time.")

	var content string
	switch m.view {
	case viewQuestion:
		content = fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", questionStyle.Render("Question:"), m.flashcards[m.current].Question, infoStyle.Render("Press Enter to reveal answer..."), bottomBar)
	case viewAnswer:
		content = fmt.Sprintf("%s\n\n%s\n\n%s\n%s", answerStyle.Render("Answer:"), m.flashcards[m.current].Answer, infoStyle.Render("Was your answer correct? [c]orrect / [i]ncorrect\n"), bottomBar)
	case viewRevisitIn:
		content = fmt.Sprintf("%s\n\n%s\n%s", infoStyle.Render("Revisit in (days): [1]  [3]  [7]  [9]"), m.resultMsg, bottomBar)
	case viewDone:
		content = fmt.Sprintf("%s\n%s", successStyle.Render("Review complete 🎉🎉🎉\n"), bottomBar)
	}
	return frame.Render(content) + "\n" + exitMsg
}
