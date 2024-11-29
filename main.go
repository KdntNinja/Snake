package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define the direction type and constants for snake movement
type direction int

const (
	up direction = iota
	down
	left
	right
)

// Define the model struct to hold the game state
type model struct {
	snake     []position
	direction direction
	food      position
	width     int
	height    int
	gameOver  bool
	tickCount int
}

// Define the position struct to represent coordinates
type position struct {
	x int
	y int
}

// Define styles for the game elements using lipgloss
var (
	borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	snakeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	headStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	foodStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

// Initialize the game model with default values
func initialModel() model {
	width, height := 20, 10
	snake := []position{{x: width / 2, y: height / 2}}
	food := position{x: rand.Intn(width), y: rand.Intn(height)}
	return model{
		snake:     snake,
		direction: right,
		food:      food,
		width:     width,
		height:    height,
		gameOver:  false,
		tickCount: 0,
	}
}

// Init the program with a tick command
func (m model) Init() tea.Cmd {
	return tick()
}

// Define the tick command to update the game state periodically
func tick() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return t
	})
}

// Update the game state based on the received message
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameOver {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key press events to change the snake's direction
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "w":
			if m.direction != down {
				m.direction = up
			}
		case "s":
			if m.direction != up {
				m.direction = down
			}
		case "a":
			if m.direction != right {
				m.direction = left
			}
		case "d":
			if m.direction != left {
				m.direction = right
			}
		}
	case time.Time:
		// Update the snake's position every second tick
		m.tickCount++
		if m.tickCount%2 == 0 {
			head := m.snake[0]
			newHead := head

			// Move the snake's head in the current direction
			switch m.direction {
			case up:
				newHead.y--
			case down:
				newHead.y++
			case left:
				newHead.x--
			case right:
				newHead.x++
			}

			// Check for collisions with the walls or itself
			if newHead.x < 0 || newHead.x >= m.width || newHead.y < 0 || newHead.y >= m.height {
				m.gameOver = true
				return m, nil
			}

			for _, p := range m.snake {
				if p == newHead {
					m.gameOver = true
					return m, nil
				}
			}

			// Update the snake's body
			m.snake = append([]position{newHead}, m.snake...)

			// Check if the snake has eaten the food
			if newHead == m.food {
				m.food = position{x: rand.Intn(m.width), y: rand.Intn(m.height)}
			} else {
				m.snake = m.snake[:len(m.snake)-1]
			}
		}
		return m, tick()
	}

	return m, nil
}

// View the game view as a string
func (m model) View() string {
	if m.gameOver {
		return "Game Over!\nPress q to quit.\n"
	}

	// Create the game board with borders
	board := make([][]rune, m.height+2)
	for i := range board {
		board[i] = make([]rune, m.width+2)
		for j := range board[i] {
			if i == 0 || i == m.height+1 {
				board[i][j] = '-'
			} else if j == 0 || j == m.width+1 {
				board[i][j] = '|'
			} else {
				board[i][j] = ' '
			}
		}
	}

	// Draw the snake on the board
	for i, p := range m.snake {
		if i == 0 {
			board[p.y+1][p.x+1] = '0'
		} else {
			board[p.y+1][p.x+1] = 'O'
		}
	}

	// Draw the food on the board
	board[m.food.y+1][m.food.x+1] = 'X'

	// Convert the board to a string with styles
	s := ""
	for _, row := range board {
		for _, cell := range row {
			switch cell {
			case '-':
				s += borderStyle.Render(string(cell))
			case '|':
				s += borderStyle.Render(string(cell))
			case 'O':
				s += snakeStyle.Render(string(cell))
			case '0':
				s += headStyle.Render(string(cell))
			case 'X':
				s += foodStyle.Render(string(cell))
			default:
				s += string(cell)
			}
		}
		s += "\n"
	}

	return s
}

// Main function to run the program
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
