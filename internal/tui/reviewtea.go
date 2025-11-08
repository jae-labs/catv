package tui

import (
	"catv/internal/store"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ReviewModel manages the state for the review session
// Views: question, answer, correct/incorrect, difficulty, done

type viewState int

const (
	viewQuestion viewState = iota
	viewAnswer
	viewDifficulty
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
	difficulty []int
}

func NewReviewModel(flashcards []store.Flashcard) *ReviewModel {
	return &ReviewModel{
		flashcards: flashcards,
		current:    0,
		view:       viewQuestion,
		correct:    make([]bool, len(flashcards)),
		difficulty: make([]int, len(flashcards)),
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
				m.view = viewDifficulty
			} else if msg.String() == "i" {
				m.correct[m.current] = false
				m.resultMsg = "Marked incorrect. Card will not be scheduled for repetition."
				m.nextCard()
			}
		case viewDifficulty:
			if msg.String() == "1" {
				m.difficulty[m.current] = 1
				m.resultMsg = fmt.Sprintf("Revisit in 1 day (%s)", store.NextReviewDate(1).Format("2006-01-02"))
				m.nextCard()
			} else if msg.String() == "3" {
				m.difficulty[m.current] = 3
				m.resultMsg = fmt.Sprintf("Revisit in 3 days (%s)", store.NextReviewDate(3).Format("2006-01-02"))
				m.nextCard()
			} else if msg.String() == "7" {
				m.difficulty[m.current] = 7
				m.resultMsg = fmt.Sprintf("Revisit in 7 days (%s)", store.NextReviewDate(7).Format("2006-01-02"))
				m.nextCard()
			} else if msg.String() == "14" {
				m.difficulty[m.current] = 14
				m.resultMsg = fmt.Sprintf("Revisit in 14 days (%s)", store.NextReviewDate(14).Format("2006-01-02"))
				m.nextCard()
			} else if msg.String() == "30" {
				m.difficulty[m.current] = 30
				m.resultMsg = fmt.Sprintf("Revisit in 30 days (%s)", store.NextReviewDate(30).Format("2006-01-02"))
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

func (m *ReviewModel) FlashcardDifficulty(idx int) int {
	if idx < 0 || idx >= len(m.difficulty) {
		return 0
	}
	return m.difficulty[idx]
}

func (m *ReviewModel) View() string {
	if m.quitting {
		return "Goodbye!"
	}
	total := len(m.flashcards)
	current := m.current + 1
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
	case viewDifficulty:
		content = fmt.Sprintf("%s\n\n%s\n%s", infoStyle.Render("Revisit in (days): [1]  [3]  [7]  [9]"), m.resultMsg, bottomBar)
	case viewDone:
		content = fmt.Sprintf("%s\n%s", successStyle.Render("Review complete! Press q to quit."), bottomBar)
	}
	return frame.Render(content) + "\n" + exitMsg
}

// ...existing code...
