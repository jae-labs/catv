package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"catv/ollama"
	"catv/store"
	"catv/tui"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type spinnerModel struct {
	spinner spinner.Model
	msg     string
	done    bool
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.done {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.done {
		return m.msg
	}
	return fmt.Sprintf("%s Generating flashcards...", m.spinner.View())
}

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate flashcards from markdown files",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			tui.PrintError("Please provide a file or folder with --path", nil)
			os.Exit(1)
		}

		// Print info about model, DB, and API target
		tui.PrintInfo(fmt.Sprintf("Model: %s", Cfg.Ollama.Model))
		tui.PrintInfo(fmt.Sprintf("Database: %s", Cfg.Database.Name))
		tui.PrintInfo(fmt.Sprintf("API Target: %s", Cfg.Ollama.URL))

		files, err := getMarkdownFiles(path)
		if err != nil {
			tui.PrintError("File error:", err)
			os.Exit(1)
		}

		for _, f := range files {
			absPath, _ := filepath.Abs(f)
			processed, err := Store.IsFileProcessed(absPath)
			if err != nil {
				tui.PrintError("DB query error:", err)
				continue
			}
			if processed {
				tui.PrintInfo(fmt.Sprintf("Skipping already processed: %s", absPath))
				continue
			}

			data, err := os.ReadFile(f)
			if err != nil {
				tui.PrintError("Read error:", err)
				continue
			}

			prompt := fmt.Sprintf(`You are an expert flashcard generator. Your task is to extract spaced repetition flashcards from the following markdown content.

Strictly output ONLY pairs in this format, with no extra text, explanations, or numbering:
Q: <question>
A: <answer>

Repeat for each flashcard. Do not include any other text, headers, or formatting. Do not add explanations, summaries, or comments. Only output Q: and A: pairs, one after another.

Example:
Q: What is the capital of France?
A: Paris
Q: What is 2+2?
A: 4

Markdown:
%s`, string(data))

			doneChan := make(chan string)
			go func() {
				resp, err := ollama.GenerateQA(Cfg.Ollama.Model, Cfg.Ollama.URL, prompt)
				if err != nil {
					doneChan <- fmt.Sprintf("Ollama error: %v", err)
					return
				}
				// DEBUG: Print raw Ollama response for troubleshooting
				tui.PrintInfo("Raw Ollama response:")
				// tui.PrintInfo(resp)
				qas, err := ollama.ParseFlashcards(resp)
				if err != nil {
					doneChan <- fmt.Sprintf("Ollama parsing error: %v", err)
					return
				}
				count := 0
				for _, qa := range qas {
					fc := store.Flashcard{
						File:       absPath,
						Question:   qa["question"],
						Answer:     qa["answer"],
						NextReview: time.Now(),
					}
					err := Store.InsertFlashcard(fc)
					if err != nil {
						tui.PrintError("DB insert error:", err)
					} else {
						count++
					}
				}
				if count > 0 {
					doneChan <- fmt.Sprintf("Processed: %s (%d flashcards generated)", absPath, count)
				} else {
					doneChan <- fmt.Sprintf("No flashcards inserted for: %s", absPath)
				}
			}()

			sm := spinnerModel{spinner: spinner.New(), done: false}
			p := tea.NewProgram(&sm)

			go func() {
				for msg := range doneChan {
					sm.msg = msg
					sm.done = true
					p.Quit()
				}
			}()

			p.Run()
			if sm.msg != "" {
				if sm.done {
					tui.PrintSuccess(sm.msg)
				} else {
					tui.PrintInfo(sm.msg)
				}
			}
		}
	},
}

func init() {
	GenerateCmd.Flags().StringP("path", "p", "", "Markdown file or folder to process")
}

func getMarkdownFiles(path string) ([]string, error) {
	var files []string
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && (filepath.Ext(p) == ".md" || filepath.Ext(p) == ".markdown") {
				files = append(files, p)
			}
			return nil
		})
	} else {
		if filepath.Ext(path) == ".md" || filepath.Ext(path) == ".markdown" {
			files = append(files, path)
		}
	}
	return files, err
}
