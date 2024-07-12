package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Game represents the state of our game
type model struct {
	playerX int
	playerY int
	score int
    deaths int
	obstacles [][]int
}

func initialModel() model {
	m := model{
		playerX:   0,
		playerY:   0,
		score:     0,
        deaths: 0,
		obstacles: make([][]int, 40),
	}
	createObstacles(m)
	return m
}

func createObstacles(m model) {
	for i := range m.obstacles {
        // create x, y pairs for each obstacle
		m.obstacles[i] = make([]int, 2)
        
        // x
		m.obstacles[i][0] = rand.Intn(19)

        // set y, dependent on x value
        // avoid start
        if (m.obstacles[i][0] == 0) {
            m.obstacles[i][1] = max(1, rand.Intn(9))
        // avoid goal
        } else if (m.obstacles[i][0] == 19) {
            m.obstacles[i][1] = rand.Intn(8)
        } else {
            m.obstacles[i][1] = rand.Intn(9)
        }

        // todo: prevent overwriting existing obstacles? overwriting makes their total ~random~
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("A Moth to a Flame")
}

// Update updates the game state based on user input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w", tea.KeyUp.String():
			m.playerY = max(m.playerY-1, 0)
		case "s", tea.KeyDown.String():
			m.playerY = min(m.playerY+1, 9)
		case "a", tea.KeyLeft.String():
			m.playerX = max(m.playerX-1, 0)
		case "d", tea.KeyRight.String():
			m.playerX = min(m.playerX+1, 19)
		case "q", tea.KeyEsc.String():
			os.Exit(0)
		default:
			return m, nil
		}
	}

	if m.checkCollision() {
		m.playerX = 0
		m.playerY = 0
        m.score = 0
        m.deaths++
        m.obstacles = make([][]int, 40) // needs dry?
        createObstacles(m)
	}

	if m.goal() {
		m.playerX = 0
		m.playerY = 0
        m.obstacles = make([][]int, 40)
        createObstacles(m)
		m.score++
	}
	return m, nil
}

func (m model) goal() bool {
	return m.playerX == 19 && m.playerY == 9
}

// checkCollision checks if the player has collided with an obstacle
func (m model) checkCollision() bool {
	for _, obs := range m.obstacles {
		if m.playerX == obs[0] && m.playerY == obs[1] {
			return true
		}
	}
	return false
}

// contains checks if an obstacle is at a given position
func (m model) contains(obstacle [][]int, x int, y int) bool {
	for i := range m.obstacles {
		if obstacle[i][0] == x && obstacle[i][1] == y {
			return true
		}
	}
	return false
}

// View renders the game state as a string
func (m model) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Score: %d Deaths: %d", m.score, m.deaths))
	sb.WriteString("\n")
	for y := 0; y < 10; y++ {
		for x := 0; x < 20; x++ {
			if x == m.playerX && y == m.playerY {
				sb.WriteString("🧚")
			} else if m.contains(m.obstacles, x, y) {
				sb.WriteString("🔥")
			} else if x == 19 && y == 9 {
				sb.WriteString("⭐️")
			} else {
				sb.WriteString("⬛")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
