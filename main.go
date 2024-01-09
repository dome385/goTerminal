package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct{
	BorderColor lipgloss.Color
	InputField lipgloss.Style
}

type Question struct {
	question string
	answer string
	input Input
}

func DefaultStyles() *Styles {
	s:= new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type Main struct{
	index int
	questions []Question
	width int
	height int
	styles *Styles
	done bool
}

func newQuestion (q string) Question {
	return Question{question: q}
}

func newShortQuestion(q string) Question {
	question := newQuestion(q)
	model := NewShortAnswerField()
	question.input = model
	return question
}

func newLongQuestion(q string) Question {
	question := newQuestion(q)
	model := NewLongAnswerField()
	question.input = model
	return question
}

func New(questions []Question) *Main{
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "Trage deine Antwort hier ein"
	answerField.Focus()
	return &Main{
		questions: questions, 
		styles: styles}
}

func (m Main) Init() tea.Cmd {
	return m.questions[m.index].input.Blink
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.questions[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			current.answer = current.input.Value()
			m.Next()
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m Main) View() string {
	current := m.questions[m.index]
	if m.done {
		var output string
		for _, q := range m.questions {
			output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
		}
		return output
	}
	if m.width == 0 {
		return "loading..."
	}
	// stack some left-aligned strings together in the center of the window
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			current.question,
			m.styles.InputField.Render(current.input.View()),
		),
	)
}

	
	
	
func (m *Main) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func main() {
	questions := []Question{
		newShortQuestion("Wie ist dein Name?"),
		newShortQuestion("Was ist dein Geburtsdatum?"), 
		newLongQuestion("Was wünscht du dir für die Zukunft?" )}
		main := New(questions)
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
		p := tea.NewProgram(*main, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
}