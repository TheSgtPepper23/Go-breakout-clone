package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDHT   float32 = 800
	HEIGHT  float32 = 600
	DEG2RAD float64 = math.Pi / 180
)

type Ball struct {
	position rl.Vector2
	speedX   float32
	speedY   float32
}

func (b *Ball) calcBounceAngle() float64 {
	return 75 * DEG2RAD * -1
}

func (b *Ball) calcNewSpeed(angle float64) {
	b.speedX = b.speedX * float32(math.Sin(angle))
	b.speedY = -b.speedY * float32(math.Cos(angle))
}

type Player struct {
	position rl.Vector2
	lives    int
	speed    float32
	width    int32
}

type Game struct {
	STATUS    int
	player    Player
	score     int32
	ball      Ball
	ballSpeed float32
}

func initializeGame() Game {
	return Game{
		player: Player{
			position: rl.Vector2{
				X: WIDHT / 2,
				Y: 560,
			},
			speed: 400,
			width: 80,
		},
		ball: Ball{
			position: rl.Vector2{
				X: WIDHT / 2,
				Y: (HEIGHT / 2) - 20,
			},
			speedX: 200,
			speedY: 200,
		},
	}
}

func (g *Game) Update() {
	delta := rl.GetFrameTime()

	// player movement
	if (rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && g.player.position.X > 0 {
		g.player.position.X -= g.player.speed * delta
	}

	if (rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && g.player.position.X+float32(g.player.width) < WIDHT {
		g.player.position.X += g.player.speed * delta
	}

	// ball movement
	g.ball.position.X += g.ball.speedX * delta
	g.ball.position.Y += g.ball.speedY * delta

	if rl.CheckCollisionRecs(rl.NewRectangle(g.ball.position.X, g.ball.position.Y, 20, 20), rl.NewRectangle(g.player.position.X, g.player.position.Y, float32(g.player.width), 20)) {
		angle := g.ball.calcBounceAngle()
		g.ball.calcNewSpeed(angle)
	}

	if g.ball.position.X <= 0 {
		g.ball.position.X = 1
		g.ball.speedX *= -1
	}

	if g.ball.position.X+20 >= WIDHT {
		g.ball.position.X = WIDHT - 20
		g.ball.speedX *= -1
	}

	if g.ball.position.Y <= 0 {
		g.ball.position.Y = 1
		g.ball.speedY *= -1
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	// draw the player
	rl.DrawRectangle(int32(g.player.position.X), int32(g.player.position.Y), g.player.width, 20, rl.Red)

	// draw the ball
	rl.DrawRectangle(int32(g.ball.position.X), int32(g.ball.position.Y), 20, 20, rl.Blue)
	rl.EndDrawing()
}

func main() {
	rl.InitWindow(int32(WIDHT), int32(HEIGHT), "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	g := initializeGame()
	for !rl.WindowShouldClose() {
		g.Update()
		g.Draw()
	}
}
