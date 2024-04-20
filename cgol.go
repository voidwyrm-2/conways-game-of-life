package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const bx = 100
const by = 100

var tickrate int = 10

const pixelscale int = 5

var updatepixels bool = false

/*
type CellSetting struct {
	x     int
	y     int
	state bool
}
*/

// taken from https://github.com/ThePrimeagen/game-of-life-vwm/blob/master/cmd/server.go
var dirs = [][]int{
	{-1, -1},
	{0, -1},
	{1, -1},
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
}

// taken from https://github.com/ThePrimeagen/game-of-life-vwm/blob/master/cmd/server.go
func getNeighborCount(state [by][bx]bool, x, y int) int {
	count := 0
	for _, d := range dirs {
		row := y + d[0]
		col := x + d[1]

		if row < 0 || row >= len(state) {
			continue
		}

		if col < 0 || col >= len(state[0]) {
			continue
		}

		if state[row][col] {
			count++
		}
	}
	return count
}

// taken from https://github.com/ThePrimeagen/game-of-life-vwm/blob/master/cmd/server.go
func stepBoard(state [by][bx]bool) [by][bx]bool {
	var out [by][bx]bool
	for y, row := range state {
		for x := range row {
			out[y][x] = state[y][x]

			count := getNeighborCount(state, x, y)
			if count < 2 || count > 3 {
				out[y][x] = false
			} else if count == 3 {
				out[y][x] = true
			}
		}
	}
	return out
}

func clamp32(number, min, max int32) int32 {
	if number < min {
		return min
	} else if number > max {
		return max
	} else {
		return number
	}
}

func clamp(number, min, max int) int {
	if number < min {
		return min
	} else if number > max {
		return max
	} else {
		return number
	}
}

var windowX int32 = 1000
var windowY int32 = 650

func main() {

	var savestate [by][bx]bool
	var cansave = false

	cannotLoadMessageTick := rl.NewColor(uint8(120), uint8(120), uint8(120), uint8(0))

	LoadedMessageTick := rl.NewColor(uint8(120), uint8(120), uint8(120), uint8(0))

	savedMessageTick := rl.NewColor(uint8(120), uint8(120), uint8(120), uint8(0))

	/*// Create new parser object
	parser := argparse.NewParser("clilexer", "lexes and prints the given string")
	// Create string flag
	boardxy := parser.IntList("xy", "boardxy", &argparse.Options{Required: true, Help: "dimentions of"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil || len(boardxy) < 2 {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}*/

	var board [by][bx]bool
	//bpoint := &board
	/*
		var presets = []CellSetting{
			{9, 8, true},
			{10, 9, true},
			{10, 10, true},
			{10, 11, true},
			{9, 11, true},
			{8, 11, true},
			{7, 11, true},
			{6, 10, true},
			{6, 8, true},
		}
	*/
	rl.InitWindow(windowX, windowY, "Conway's Game of Life")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	/*
		for _, p := range presets {
			board[p.y][p.x] = p.state
		}
	*/

	var tick int = 0

	for !rl.WindowShouldClose() {
		tickrate = clamp(tickrate, 0, 60)
		tick++

		if rl.IsKeyPressed(rl.KeySpace) && updatepixels {
			updatepixels = false
		} else if rl.IsKeyPressed(rl.KeySpace) && !updatepixels {
			updatepixels = true
		}

		if rl.IsKeyPressed(rl.KeyS) {
			cansave = true
			savestate = board
			savedMessageTick.A = 150
		}

		if rl.IsKeyPressed(rl.KeyL) && cansave && savestate != board {
			board = savestate
			//cansave = false
			LoadedMessageTick.A = 150
		} else if rl.IsKeyPressed(rl.KeyL) && (!cansave || savestate == board) {
			cannotLoadMessageTick.A = 150
		}

		if rl.IsKeyPressed(rl.KeyUp) {
			tickrate++
		} else if rl.IsKeyPressed(rl.KeyDown) {
			tickrate--
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			mx := clamp32(rl.GetMouseX()/int32(pixelscale), 0, bx-1)
			my := clamp32(rl.GetMouseY()/int32(pixelscale), 0, by-1)
			if board[my][mx] {
				board[my][mx] = false
			} else {
				board[my][mx] = true
			}
		}

		if updatepixels && tick > tickrate {
			board = stepBoard(board)
			tick = 0
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		//rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.White)

		for y, b1 := range board {
			for x, b2 := range b1 {
				if !b2 {
					continue
				}
				pixcol := rl.White

				rl.DrawRectangle(int32(x*pixelscale), int32(y*pixelscale), int32(pixelscale), int32(pixelscale), pixcol)
			}
		}

		if !updatepixels {
			rl.DrawText("currently paused, press Space to unpause", 2, 5, 20, rl.Gray)
		}

		rl.DrawText(fmt.Sprint(tickrate), 2, 25, 20, rl.Gray)

		rl.DrawText("current state saved", 2, windowY-30, 20, savedMessageTick)

		rl.DrawText("loaded saved state", 2, windowY-60, 20, LoadedMessageTick)

		rl.DrawText("unable to load save", 2, windowY-90, 20, cannotLoadMessageTick)

		if savedMessageTick.A > 0 {
			savedMessageTick.A--
		}

		if cannotLoadMessageTick.A > 0 {
			cannotLoadMessageTick.A--
		}

		if LoadedMessageTick.A > 0 {
			LoadedMessageTick.A--
		}

		rl.EndDrawing()
	}
}
