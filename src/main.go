package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDHT  float32 = 800
	HEIGHT float32 = 600
)

type Ball struct {
	position rl.Vector2
	speedX   float32
	speedY   float32
	size     float32
}

func (b *Ball) calcNewSpeed(player Player) {
	b.position.Y = b.position.Y - 1
	b.speedY *= -1

	rightSide := b.position.X+(b.size/2) > player.position.X+(player.width/2)

	if (b.speedX < 0 && rightSide) || (b.speedX > 0 && !rightSide) {
		b.speedX *= -1
	}
}

func (b Ball) getRect() rl.Rectangle {
	return rl.NewRectangle(b.position.X, b.position.Y, b.size, b.size)
}

type Player struct {
	position rl.Vector2
	lives    int
	width    float32
}

func (p Player) getRect() rl.Rectangle {
	return rl.NewRectangle(p.position.X, p.position.Y, p.width, 20)
}

type Brick struct {
	position rl.Vector2
	width    float32
	height   float32
	broken   bool
}

func (b Brick) getRect() rl.Rectangle {
	return rl.NewRectangle(b.position.X, b.position.Y, b.width, b.height)
}

type BrickZone struct {
	bricks []*Brick
	zone   rl.Rectangle
}

type GAMESTATE int

const (
	GAMEOVER GAMESTATE = iota
	PAUSE
	GAMEON
	GAMEREADY
)

type Game struct {
	brickSections []BrickZone
	player        Player
	STATUS        GAMESTATE
	ball          Ball
	score         int32
	ballSpeed     float32
}

func initializeGame() Game {
	brickZones := generateBricks(10, 5)
	return Game{
		player: Player{
			position: rl.Vector2{
				X: WIDHT / 2,
				Y: 560,
			},
			width: 80,
		},
		ball: Ball{
			position: rl.Vector2{
				X: WIDHT / 2,
				Y: (HEIGHT / 2) - 20,
			},
			speedX: 200,
			speedY: 200,
			size:   20,
		},
		STATUS:        GAMEREADY,
		brickSections: brickZones,
	}
}

func generateBricks(cols, rows int) []BrickZone {
	brickWidth := WIDHT / float32(cols)
	var brickHeight float32 = 30
	brickZones := make([]BrickZone, 0)
	var currentX float32 = 0.0
	currentY := brickHeight * 2

	for i := 0; i < rows; i++ {
		bricksInZone := make([]*Brick, 0)
		for j := 0; j < cols; j++ {
			temp := Brick{
				position: rl.Vector2{
					X: currentX,
					Y: currentY,
				},
				width:  brickWidth,
				height: brickHeight,
				broken: false,
			}
			bricksInZone = append(bricksInZone, &temp)
			currentX += brickWidth
		}
		brickZones = append(brickZones, BrickZone{
			zone:   rl.NewRectangle(0, currentY-5, WIDHT, brickHeight+10),
			bricks: bricksInZone,
		})
		currentY += brickHeight
		currentX = 0
	}
	return brickZones
}

func (g *Game) Update() {
	if rl.IsKeyPressed(rl.KeySpace) {
		g.STATUS = GAMEON
	}
	if g.STATUS == GAMEON {
		delta := rl.GetFrameTime()
		mousePos := rl.GetMousePosition()
		g.player.position.X = mousePos.X - g.player.width/2
		// ball movement
		g.ball.position.X += g.ball.speedX * delta
		g.ball.position.Y += g.ball.speedY * delta

		// left wall collision
		if g.ball.position.X <= 0 {
			g.ball.position.X = 1
			g.ball.speedX *= -1
		}

		// rigth wall collision
		if g.ball.position.X+g.ball.size >= WIDHT {
			g.ball.position.X = WIDHT - g.ball.size
			g.ball.speedX *= -1
		}

		// ceiling collision
		if g.ball.position.Y <= 0 {
			g.ball.position.Y = 1
			g.ball.speedY *= -1
		}

		// goes beyond the paddle
		if g.ball.position.Y > g.player.position.Y+20 {
			g.STATUS = GAMEOVER
		}

		// hits the paddle
		if rl.CheckCollisionRecs(g.ball.getRect(), g.player.getRect()) {
			// checks which side of the paddle hit
			g.ball.calcNewSpeed(g.player)
		}

		for _, brickZone := range g.brickSections {
			if rl.CheckCollisionPointRec(rl.NewVector2(g.ball.position.X, g.ball.position.Y), brickZone.zone) ||
				rl.CheckCollisionPointRec(rl.NewVector2(g.ball.position.X+g.ball.size, g.ball.position.Y+g.ball.size), brickZone.zone) {
				for _, brick := range brickZone.bricks {
					if !brick.broken && rl.CheckCollisionRecs(brick.getRect(), g.ball.getRect()) {
						brick.broken = true
						g.ball.speedX *= 1.02
						g.ball.speedY *= 1.02
						g.ball.calcNewSpeed(g.player)
					}
				}
			}
		}
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	for _, section := range g.brickSections {
		for _, brick := range section.bricks {
			if !brick.broken {
				rl.DrawRectangleRec(brick.getRect(), rl.Yellow)
				rl.DrawRectangleLines(int32(brick.position.X), int32(brick.position.Y), int32(brick.width), int32(brick.height), rl.Black)
			}
		}
	}

	rl.DrawFPS(0, 0)

	// draw the player
	rl.DrawRectangleRec(g.player.getRect(), rl.Red)

	// draw the ball
	rl.DrawRectangleRec(g.ball.getRect(), rl.Blue)
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
