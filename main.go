package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"golang.org/x/image/font/basicfont"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	screenWidth  = 320
	screenHeight = 240
	tileSize     = 5
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

// Game представляет игровой процесс.
type Game struct {
	snake         *Snake // Змея в игре.
	food          *Food  // Еда в игре.
	score         int    // Счет игрока.
	gameOver      bool   // Флаг, указывающий на завершение игры.
	ticks         int    // Счетчик тиков (шагов) игры.
	updateCounter int    // Счетчик для обновления состояния игры.
	speed         int    // Скорость движения змеи.
}

func (g *Game) Update(screen *ebiten.Image) error {
	// Проверка, завершена ли игра.
	if g.gameOver {
		// Если игра завершена, проверяем, была ли нажата клавиша R для перезапуска.
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.restart()
		}
		return nil // Возвращаем nil, так как игра завершена, и обновление не требуется.
	}

	// Увеличиваем счетчик обновлений игры.
	g.updateCounter++

	// Проверяем, достигло ли количество обновлений заданной скорости.
	if g.updateCounter < g.speed {
		return nil // Если не достигло, возвращаем nil, обновление не требуется.
	}

	// Сбрасываем счетчик обновлений и продолжаем обновление.
	g.updateCounter = 0

	// Вызываем метод Move() для обновления позиции змеи.
	g.snake.Move()

	// Проверяем нажатие клавиш для изменения направления змеи.
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.snake.Direction.X == 0 {
		g.snake.Direction = Point{X: -1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && g.snake.Direction.X == 0 {
		g.snake.Direction = Point{X: 1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) && g.snake.Direction.Y == 0 {
		g.snake.Direction = Point{X: 0, Y: -1}
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.snake.Direction.Y == 0 {
		g.snake.Direction = Point{X: 0, Y: 1}
	}

	// Получаем координаты головы змеи.
	head := g.snake.Body[0]

	// Проверяем, не выходит ли змея за границы экрана.
	if head.X < 0 || head.Y < 0 || head.X >= screenWidth/tileSize || head.Y >= screenHeight/tileSize {
		g.gameOver = true
		g.speed = 10
	}
	// Проверяем, не столкнулась ли голова змеи с её частями.
	for _, part := range g.snake.Body[1:] {
		if head.X == part.X && head.Y == part.Y {
			g.gameOver = true
			g.speed = 10
		}
	}

	// Проверяем, не съела ли змея еду.
	if head.X == g.food.Position.X && head.Y == g.food.Position.Y {
		// Увеличиваем счёт и устанавливаем новую еду.
		g.score++
		g.snake.GrowCounter += 1
		g.food = NewFood()

		// Уменьшаем скорость (с ограничением снизу).
		if g.speed > 2 {
			g.speed--
		}
	}

	return nil
}

// Draw отвечает за отрисовку игрового состояния на экране.
func (g *Game) Draw(screen *ebiten.Image) {
	// Отрисовываем фон.
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Отрисовываем змею.
	for _, p := range g.snake.Body {
		ebitenutil.DrawRect(screen, float64(p.X*tileSize), float64(p.Y*tileSize), tileSize, tileSize, color.RGBA{0, 255, 0, 255})
	}

	// Отрисовываем еду.
	ebitenutil.DrawRect(screen, float64(g.food.Position.X*tileSize), float64(g.food.Position.Y*tileSize), tileSize, tileSize, color.RGBA{255, 0, 0, 255})

	// Создаём шрифт.
	face := basicfont.Face7x13

	// Если игра окончена, выводим текст "Game Over" и инструкцию по рестарту.
	if g.gameOver {
		text.Draw(screen, "Game Over", face, screenWidth/2-40, screenHeight/2, color.White)
		text.Draw(screen, "Press 'R' to restart", face, screenWidth/2-60, screenHeight/2+16, color.White)
	}

	// Выводим текущий счёт.
	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreText, face, 5, screenHeight-5, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	// Инициализация генератора случайных чисел на основе текущего времени.
	rand.Seed(time.Now().UnixNano())

	// Создание новой игры.
	game := &Game{
		snake:    NewSnake(),
		food:     NewFood(),
		gameOver: false,
		speed:    10,
	}

	// Установка размеров окна и заголовка.
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}


func (g *Game) restart() {
	g.snake = NewSnake()
	g.score = 0
	g.gameOver = false
	g.food = NewFood()
	g.speed = 10
}
