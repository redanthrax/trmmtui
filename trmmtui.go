package main

import (
	"strings"
  "fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
  "github.com/charmbracelet/bubbles/list"
	term "github.com/nsf/termbox-go"
)

var (
  focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle = focusedStyle.Copy()
  noStyle = lipgloss.NewStyle()
  focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type client struct {
  name string
}

func (c client) FilterValue() string { return c.name }

type model struct {
	mode string
  url string
  key string
  focusIndex int
  inputs []textinput.Model
  cursorMode cursor.Mode
	clients list.Model
	rightColumn string
  apiErr error
}

func main() {
	// Initialize the application with a starting model
	program := tea.NewProgram(initialModel(), tea.WithAltScreen())

	// Start the event loop
	if err := program.Start(); err != nil {
		panic(err)
	}
}

func (m model) Init() tea.Cmd {
  return textinput.Blink
}

// Initialize the initial model
func initialModel() model {
  clients := []list.Item {
    client { name: "Client1"},
    client { name: "Client2"},
  }

  m := model{
		mode: "app",
    inputs: make([]textinput.Model, 2),
    clients: list.New(clients, list.NewDefaultDelegate(), 0, 0),
		rightColumn: "Table here",
	}

  var t textinput.Model
  for i := range m.inputs {
    t = textinput.New()
    t.Cursor.Style = cursorStyle
    t.CharLimit = 255

    switch i {
    case 0:
      t.Placeholder = "https://api.company.com"
      t.Focus()
      t.PromptStyle = focusedStyle
      t.TextStyle = focusedStyle
    case 1:
      t.Placeholder = "XXXXXCFDA2WBCXH0XTELBR5KAI69XXXX"
      t.EchoMode = textinput.EchoPassword
      t.EchoCharacter = '*'
    }

    m.inputs[i] = t
  }

  return m
}

// Update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case "login":
			switch msg.String() {
      case "tab", "shift+tab", "enter", "up", "down":
        s := msg.String()
        if s == "enter" && m.focusIndex == len(m.inputs) {
          //submit button focused, do the login
          if m.url != "" && m.key != "" {
            //do login and loadup here
            //validate good or bad
            m.mode = "app"
          }
        }

        if s == "up" || s == "shift+tab" {
          m.focusIndex--
        } else {
          m.focusIndex++
        }

        if m.focusIndex > len(m.inputs) {
          m.focusIndex = 0
        } else if m.focusIndex < 0 {
          m.focusIndex = len(m.inputs)
        }

        cmds := make([]tea.Cmd, len(m.inputs))
        for i := 0; i <= len(m.inputs)-1; i++ {
          if i == m.focusIndex {
            cmds[i] = m.inputs[i].Focus()
            m.inputs[i].PromptStyle = focusedStyle
            m.inputs[i].TextStyle = focusedStyle
            continue
          }

          m.inputs[i].Blur()
          m.inputs[i].PromptStyle = noStyle
          m.inputs[i].TextStyle = noStyle
        }

        return m, tea.Batch(cmds...)
      default:
        if m.focusIndex == 0 {
          m.url += string(msg.Runes)
        } else if m.focusIndex == 1 {
          m.key += string(msg.Runes)
        }
			}
		}

    switch msg.String() {
    case "ctrl+c", "esc", "q":
      return m, tea.Quit
    }
	}

  cmd := m.updateInputs(msg)
  return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
  cmds := make([]tea.Cmd, len(m.inputs))
  for i := range m.inputs {
    m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
  }

  return tea.Batch(cmds...)
}

// View function
func (m model) View() string {
	if m.mode == "login" {
		// login view
    var b strings.Builder
    for i := range m.inputs {
      b.WriteString(m.inputs[i].View())
      if i < len(m.inputs)-1 {
        b.WriteRune('\n')
      }
    }

    button := &blurredButton
    if m.focusIndex == len(m.inputs) {
      button = &focusedButton
    }

    fmt.Fprintf(&b, "\n\n%s\n\n", *button)
    return b.String()
	} else {
		// get terminal dimensions
		w, h := term.Size()

		// define the proportions for each column
		leftWidth := w / 4
		rightWidth := w * 3 / 4

		// create UI styles using lipgloss
		leftStyle := lipgloss.NewStyle().Width(leftWidth).Height(h).Border(lipgloss.NormalBorder())
		rightStyle := lipgloss.NewStyle().Width(rightWidth).Height(h).Border(lipgloss.NormalBorder())

		// layout the UI
		return leftStyle.Render(m.clients.View()) + rightStyle.Render(m.rightColumn)
	}
}

// A helper function to mask user input
func maskString(input string) string {
	var masked string
	for range input {
		masked += "*"
	}
	return masked
}
