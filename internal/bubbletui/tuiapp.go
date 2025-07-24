package bubbletui

import (
  "context"
  "fmt"
  "strings"

  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/bubbles/textinput"
  "github.com/charmbracelet/lipgloss"
  "github.com/fatih/color"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

// chatMsg holds one line in the chat history.
type chatMsg struct {
  user    bool   // true if “You: …”, false if “Agent: …”
  content string // the text
}

// chatModel is our Bubble Tea model.
type chatModel struct {
  ctx     context.Context
  App     app.App
  width   int
  height  int
  history []chatMsg
  input   textinput.Model
}

// NewChatModel constructs the model.
func NewChatModel(ctx context.Context, a app.App) chatModel {
  ti := textinput.New()
  ti.Placeholder = "Type your message here"
  ti.CharLimit = 512
  ti.Width = 50
  ti.Prompt = "> "
  ti.Focus()

  return chatModel{
    ctx:     ctx,
    App:     a,
    history: make([]chatMsg, 0),
    input:   ti,
  }
}

// Init tells Bubble Tea to start the textinput caret blinking.
func (m chatModel) Init() tea.Cmd {
  return textinput.Blink
}

// Update handles incoming events: window resize, key presses, and textinput.
func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {

  case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    // make the input line span almost the full width
    m.input.Width = msg.Width - 4
    return m, nil

  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c", "q":
      return m, tea.Quit
    case "enter":
      // send it
      userLine := strings.TrimSpace(m.input.Value())
      if userLine != "" {
        // append user message
        m.history = append(m.history, chatMsg{user: true, content: userLine})
        // clear input
        m.input.SetValue("")

        // call your core SendMessage → must return (reply string, err error)
        reply, err := m.App.SendMessage(m.ctx, userLine)
        if err != nil {
          m.history = append(m.history,
            chatMsg{user: false, content: fmt.Sprintf("[error] %v", err)},
          )
        } else {
          m.history = append(m.history, chatMsg{user: false, content: reply})
        }
      }
      return m, textinput.Blink
    }
  }

  // delegate everything else to the textinput
  var cmd tea.Cmd
  m.input, cmd = m.input.Update(msg)
  return m, cmd
}

// View renders the screen: a header, the scrollable history, and the input line.
func (m chatModel) View() string {
  var b strings.Builder

  header := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("36")).
    Render("🐬 Dolphin Chat")
  b.WriteString(header + "\n\n")

  youStyle := color.New(color.FgCyan, color.Bold).SprintFunc()
  agStyle := color.New(color.FgGreen, color.Bold).SprintFunc()

  // Print history.  You could add real scrolling logic here if
  // len(history) > available lines.
  for _, cm := range m.history {
    if cm.user {
      b.WriteString(youStyle("You: ") + cm.content + "\n")
    } else {
      b.WriteString(agStyle("Agent: ") + cm.content + "\n")
    }
  }

  // Leave a blank line, then the input box
  b.WriteString("\n" + m.input.View())

  // hint
  b.WriteString("\n\n" + lipgloss.NewStyle().Faint(true).
    Render("Enter to send • q or Ctrl+C to quit"))

  return b.String()
}

// RunChatTUI is what you call from main().
func RunChatTUI(ctx context.Context, a app.App) error {
  p := tea.NewProgram(
    NewChatModel(ctx, a),
    tea.WithAltScreen(),       // clear screen / restore on exit
  )
  _, err := p.Run()
  return err
}
