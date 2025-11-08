package commands

import (
	"catv/internal/tui"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var ReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review flashcards",
	Run: func(cmd *cobra.Command, args []string) {
		flashcards, err := Store.GetFlashcardsForReview()
		if err != nil {
			tui.PrintError("DB query error:", err)
			return
		}

		if len(flashcards) == 0 {
			tui.PrintInfo("No flashcards due for review.")
			return
		}

		// Run Bubble Tea TUI for review
		model := tui.NewReviewModel(flashcards)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running review TUI:", err)
		}

		// After review, update DB only for flashcards that were actually answered
		for i, fc := range flashcards {
			if model.FlashcardWasCorrect(i) && model.FlashcardRevisitIn(i) > 0 {
				fc.RevisitIn = model.FlashcardRevisitIn(i)
				if err := Store.UpdateFlashcard(fc); err != nil {
					tui.PrintError("DB update error:", err)
				} else {
					tui.PrintSuccess(fmt.Sprintf("Updated flashcard %d: revisitin=%d", fc.ID, fc.RevisitIn))
				}
			} else if !model.FlashcardWasCorrect(i) && model.FlashcardRevisitIn(i) > 0 {
				// Incorrect -> revisit sooner (e.g. tomorrow => 1)
				fc.RevisitIn = 1
				if err := Store.UpdateFlashcard(fc); err != nil {
					tui.PrintError("DB update error:", err)
				} else {
					tui.PrintSuccess(fmt.Sprintf("Marked flashcard %d incorrect: revisitin set to 1", fc.ID))
				}
			}
		}
	},
}
