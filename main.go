package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 160
	screenHeight = 120
	tileSize     = 3
)

// Point представляет координаты на экране.
type Point struct {
	X int
	Y int
}

// Snake представляет змею в игре.
type Snake struct {
	Body        []Point // Тело змеи (координаты каждого сегмента).
	Direction   Point   // Направление движения змеи.
	GrowCounter int     // Счетчик роста - определяет, сколько сегментов нужно добавить к змее.
}

// NewSnake создает новый экземпляр змеи с начальными параметрами.
func NewSnake() *Snake {
	return &Snake{
		Body: []Point{
			{
				X: screenWidth / tileSize / 2,
				Y: screenHeight / tileSize / 2,
			},
		},
		Direction: Point{
			X: 1,
			Y: 0,
		},
	}
}

// Move выполняет шаг движения змеи.
func (s *Snake) Move() {
	// Вычисляем новую голову змеи.
	newHead := Point{
		X: s.Body[0].X + s.Direction.X,
		Y: s.Body[0].Y + s.Direction.Y,
	}

	// Добавляем новую голову в начало тела змеи.
	s.Body = append([]Point{newHead}, s.Body...)

	// Если счетчик роста больше 0, уменьшаем его, иначе убираем хвост (делаем змею длиннее).
	if s.GrowCounter > 0 {
		s.GrowCounter--
	} else {
		s.Body = s.Body[:len(s.Body)-1]
	}
}

// Food представляет еду в игре.
type Food struct {
	Position Point // Позиция еды на экране.
}

// NewFood создает новый экземпляр еды со случайным положением.
func NewFood() *Food {
	return &Food{
		Position: Point{
			X: rand.Intn(screenWidth / tileSize),
			Y: rand.Intn(screenHeight / tileSize),
		},
	}
}


func main() {
	// Инициализация генератора случайных чисел на основе текущего времени.
	rand.Seed(time.Now().UnixNano())

	// Создание новой игры.
	game := &Game{
		snake:    NewSnake(),
		food:     NewFood(),
		gameOver: false,
		ticks:    0,
		speed:    10,
	}

	// Установка размеров окна и заголовка.
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
