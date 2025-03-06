package main

import (
	"math/rand"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	windowWidth         = 1280
	windowHeight        = 720
	targetFPS           = 60
	batDistanceToWindow = 32
	batWidth            = 16
	batHeight           = 256
	ballSize            = 16
	ballSpeed           = 4
)

type Ball = struct {
	pos       rl.Vector2
	size      float32
	direction rl.Vector2
	velocity  float32
}

func createBall() Ball {
	return Ball{
		pos:  rl.NewVector2(windowWidth/2-ballSize/2, windowHeight/2-ballSize/2),
		size: ballSize,
		direction: rl.Vector2Normalize(rl.NewVector2(
			1,
			rand.Float32()*5-2.5,
		)),
		velocity: rand.Float32()*6 + 4,
	}
}

type Stats = struct {
	player1Points int
	player2Points int
}

func checkGoal(ball *Ball, stats *Stats) {
	if ball.pos.X <= ballSize {
		stats.player2Points += 1
		*ball = createBall()
	}
	if ball.pos.X+ball.size >= windowWidth {
		stats.player1Points += 1
		*ball = createBall()
	}
}

func updateBallPosition(ball *Ball, player1 rl.Rectangle, player2 rl.Rectangle) {
	accelerationVec := rl.Vector2Scale(ball.direction, ball.velocity)
	newPos := rl.Vector2Add(ball.pos, accelerationVec)

	// check collisions with border
	if newPos.X < ballSize {
		newPos.X = ballSize
		ball.direction.X *= -1
	}

	if newPos.X+ball.size > windowWidth {
		newPos.X = windowWidth - ball.size
		ball.direction.X *= -1
	}

	if newPos.Y < ballSize {
		newPos.Y = ballSize
		ball.direction.Y *= -1
	}

	if newPos.Y+ball.size > windowHeight {
		newPos.Y = windowHeight - ball.size
		ball.direction.Y *= -1
	}

	// check collision with players
	updateBallPositionPlayerCollision(ball, player1, &newPos)
	updateBallPositionPlayerCollision(ball, player2, &newPos)

	ball.pos = newPos
}

func updateBallPositionPlayerCollision(ball *Ball, player rl.Rectangle, newPos *rl.Vector2) {
	if rl.CheckCollisionCircleRec(ball.pos, ball.size, player) {
		if ball.pos.X > player.X+player.Width {
			newPos.X = player.X + player.Width + ballSize + 1
			ball.direction.X *= -1
		}
		if ball.pos.X < player.X {
			newPos.X = player.X - ballSize - 1
			ball.direction.X *= -1
		}
		if ball.pos.Y > player.Y+player.Height {
			newPos.Y = player.Y + player.Height + ballSize + 1
			ball.direction.Y *= -1
		}
		if ball.pos.Y < player.Y {
			newPos.Y = player.Y - ballSize - 1
			ball.direction.Y *= -1
		}
	}
}

func updatePlayerPosition(player *rl.Rectangle, isPlayer2 bool) {
	var keyUp int32 = rl.KeyW
	var keyDown int32 = rl.KeyS
	if isPlayer2 {
		keyUp = rl.KeyUp
		keyDown = rl.KeyDown
	}

	if rl.IsKeyDown(keyDown) {
		player.Y += 8
	}

	if rl.IsKeyDown(keyUp) {
		player.Y -= 8
	}

	if player.Y < 0 {
		player.Y = 0
	}

	if player.Y+player.Height > windowHeight {
		player.Y = windowHeight - player.Height
	}
}

func main() {
	// setup
	rl.InitWindow(windowWidth, windowHeight, "Pingo")
	defer rl.CloseWindow()
	rl.SetTargetFPS(targetFPS)

	// game objects
	player1 := rl.Rectangle{
		X:      batDistanceToWindow,
		Y:      windowHeight/2 - batHeight/2,
		Width:  batWidth,
		Height: batHeight,
	}

	player2 := rl.Rectangle{
		X:      windowWidth - batWidth/2 - batDistanceToWindow,
		Y:      windowHeight/2 - batHeight/2,
		Width:  batWidth,
		Height: batHeight,
	}

	ball := createBall()

	stats := Stats{
		player1Points: 0,
		player2Points: 0,
	}

	// game loop
	for !rl.WindowShouldClose() {

		// update
		updatePlayerPosition(&player1, false)
		updatePlayerPosition(&player2, true)
		updateBallPosition(&ball, player1, player2)
		checkGoal(&ball, &stats)

		// drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawRectangleRec(player1, rl.White)
		rl.DrawRectangleRec(player2, rl.White)
		rl.DrawCircleV(ball.pos, ball.size, rl.White)

		rl.DrawText(strconv.Itoa(stats.player1Points), 128, 16, 32, rl.White)
		rl.DrawText(strconv.Itoa(stats.player2Points), windowWidth-128, 16, 32, rl.White)

		rl.EndDrawing()
	}
}
