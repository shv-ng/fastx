package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/shv-ng/fastx/internal/config"
)

type step int

const (
	stepProjectName step = iota
	stepPythonVersion
	stepPackageManager
	stepDatabase
	stepORM
	stepMigrations
	stepDocker
	stepCI
	stepDone
)

type WizardModel struct {
	step    step
	Answers *config.Config

	textInput textinput.Model
	cursor    int
	choices   []string
}

func InitialModel() WizardModel {
	ti := textinput.New()
	ti.Placeholder = "my-app"
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(30)

	return WizardModel{
		step:      stepProjectName,
		Answers:   &config.Config{},
		textInput: ti,
		cursor:    0,
		choices:   []string{"3.10", "3.11", "3.12"},
	}
}

func (m WizardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			m = m.advanceStep()
			if m.step == stepDone {
				return m, tea.Quit
			}
			return m, nil

		case "up", "down", "j", "k":
			if m.step != stepProjectName {
				if (msg.String() == "up" || msg.String() == "k") && m.cursor > 0 {
					m.cursor--
				} else if (msg.String() == "down" || msg.String() == "j") && m.cursor < len(m.choices)-1 {
					m.cursor++

				}
			}
		}
	}
	if m.step == stepProjectName {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m WizardModel) advanceStep() WizardModel {
	switch m.step {
	case stepProjectName:
		m.Answers.ProjectName = m.textInput.Value()
		m.step = stepPythonVersion
		m.cursor = 0
		m.choices = []string{"3.10", "3.11", "3.12"}

	case stepPythonVersion:
		m.Answers.PythonVersion = m.choices[m.cursor]
		m.step = stepPackageManager
		m.cursor = 0
		m.choices = []string{"pip", "poetry", "pdm", "uv"}

	case stepPackageManager:
		m.Answers.PackageManager = m.choices[m.cursor]
		m.step = stepDatabase
		m.cursor = 0
		m.choices = []string{"PostgreSQL", "MySQL", "SQLite", "None"}

	case stepDatabase:
		m.Answers.Database = m.choices[m.cursor]
		if m.Answers.Database == "None" {
			// Skip ORM and Migrations if no DB
			m.Answers.ORM = "None"
			m.Answers.Migrations = "None"
			m.step = stepDocker
			m.cursor = 0
			m.choices = []string{"Yes", "No"}
		} else {
			m.step = stepORM
			m.cursor = 0
			m.choices = []string{"SQLAlchemy", "SQLModel", "Tortoise", "RawAsyncpg"}
		}

	case stepORM:
		m.Answers.ORM = m.choices[m.cursor]
		m.step = stepMigrations
		m.cursor = 0
		// Dynamically fetch valid migrations based on your config package
		m.choices = config.ValidMigrations(m.Answers.ORM)

	case stepMigrations:
		m.Answers.Migrations = m.choices[m.cursor]
		m.step = stepDocker
		m.cursor = 0
		m.choices = []string{"Yes", "No"}

	case stepDocker:
		m.Answers.UseDocker = m.choices[m.cursor] == "Yes"
		m.step = stepCI
		m.cursor = 0
		m.choices = []string{"Yes", "No"}

	case stepCI:
		m.Answers.UseCI = m.choices[m.cursor] == "Yes"
		m.step = stepDone
	}

	return m
}

func (m WizardModel) View() tea.View {
	if m.step == stepDone {
		s := fmt.Sprintf("\n🎉 Configuration complete!\n\nProject: %s\nPython: %s\nDB: %s\nORM: %s\nMigrations: %s\n\nStarting generation...\n",
			m.Answers.ProjectName, m.Answers.PythonVersion, m.Answers.Database, m.Answers.ORM, m.Answers.Migrations)
		v := tea.NewView(s)
		return v
	}

	var b strings.Builder

	// Render based on current step
	switch m.step {
	case stepProjectName:
		b.WriteString("\nWhat is the name of your project?\n\n")
		b.WriteString(m.textInput.View())
		b.WriteString("\n\n(esc to quit)")

	default:
		// Generic renderer for all multiple-choice steps
		fmt.Fprintf(&b, "\n%s\n\n", m.getPromptForStep(m.step))

		for i, choice := range m.choices {
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor
			}
			fmt.Fprintf(&b, "%s %s\n", cursor, choice)
		}
		b.WriteString("\n(Use arrow keys to move, enter to select, esc to quit)")
	}

	v := tea.NewView(b.String())
	return v
}

// getPromptForStep returns the question string for the current step
func (m WizardModel) getPromptForStep(s step) string {
	switch s {
	case stepPythonVersion:
		return "Which Python version would you like to use?"
	case stepPackageManager:
		return "Which package manager do you prefer?"
	case stepDatabase:
		return "Which database would you like to configure?"
	case stepORM:
		return "Which ORM would you like to use?"
	case stepMigrations:
		return "Which migration tool would you like to use?"
	case stepDocker:
		return "Do you want to generate a Dockerfile?"
	case stepCI:
		return "Do you want to generate CI/CD pipelines (e.g., GitHub Actions)?"
	default:
		return "Make a selection:"
	}
}
