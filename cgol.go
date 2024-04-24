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

func clamp32(number, min, max int32) int32 {
	if number < min {
		return min
	} else if number > max {
		return max
	} else {
		return number
	}
}

func getMousePixelPos() (int, int) {
	return int(clamp32(rl.GetMouseX()/int32(pixelscale), 0, bx-1)), int(clamp32(rl.GetMouseY()/int32(pixelscale), 0, by-1))
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

type contruct struct { // used for creating things that can be placed as preset constructs, such as gliders, etc
	ctype     string   // "Still Life", "Spaceship", etc
	name      string   // "Glider", "Block", etc
	rotatable bool     // for rotation, duh
	positions [][2]int // stores the cell positions relative to the mouse position, so {1, -1} would be up and to the right of the mouse
}

func (c contruct) build(board [bx][by]bool) ([bx][by]bool, bool) {
	out := board
	for _, pos := range c.positions {
		mx, my := getMousePixelPos()
		px := pos[0]
		py := pos[1]
		if my+py > len(board)-1 || mx-px > len(board[0])-1 {
			return board, true
		}
		if my+py < 0 || mx-px < 0 {
			return board, true
		}
		out[my+py][mx-px] = true
	}
	return out, false
}

var constructs = []contruct{
	// Still Lifes
	{
		"Still Life",
		"Block",
		true,
		[][2]int{
			{0, 0},
			{1, -1},
			{0, -1},
			{1, 0},
		},
	},
	{
		"Still Life",
		"Beehive",
		true,
		[][2]int{
			{1, -1},
			{1, 1},
			{0, -1},
			{0, 1},
			{-1, 0},
			{2, 0},
		},
	},
	/*
		{
			"Still Life",
			"Crab",
			true,
			[][2]int{
				{-1, 1},
				{-1, -1},
				{0, 1},
				{0, -1},
				{1, 0},
				{-2, 0},
			},
		},
	*/

	// Oscillators
	{
		"Oscillator",
		"Blinker",
		true,
		[][2]int{
			{0, 0},
			{-1, 0},
			{1, 0},
		},
	},

	// Spaceships
	{
		"Spaceship",
		"Glider",
		true,
		[][2]int{
			{0, 0},
			{1, -1},
			{-1, 0},
			{-1, -1},
			{0, 1},
		},
	},

	// Expanders
	{
		"Expander",
		"Blinker(4)",
		true,
		[][2]int{
			{0, 0},
			{1, 0},
			{0, -1},
			{0, 1},
		},
	},
	{
		"Expander",
		"Beehive(4)",
		true,
		[][2]int{
			{0, 0},
			{1, -1},
			{1, 1},
			{0, -1},
			{0, 1},
			{-1, 0},
			{2, 0},
		},
	},

	// Decayers
	{
		"Decayer",
		"Tuning Fork",
		true,
		[][2]int{
			{0, -2},
			{1, -1},
			{-1, -1},
			{1, 0},
			{-1, 0},
			{1, 1},
			{-1, 1},
			{1, 2},
			{-1, 2},
		},
	},
}

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

var holdingShift bool = false

var constructIndex int = 0

var constructorMode bool = false

var windowX int32 = 1000
var windowY int32 = 650

func main() {

	var savestate [by][bx]bool
	var cansave = false

	cannotLoadMessageTick := rl.NewColor(uint8(120), uint8(120), uint8(120), uint8(0))

	LoadedMessageTick := rl.NewColor(uint8(120), uint8(120), uint8(120), uint8(0))

	savedMessageTick := rl.NewColor(uint8(120), uint8(120), uint8(120), uint8(0))

	cBuildFailedMessageTick := rl.NewColor(uint8(120), uint8(120), uint8(120), uint8(0))

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
	rl.InitWindow(windowX, windowY, "Conway's Game of Life")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var tick int = 0

	rl.InitAudioDevice()

	music_main := rl.LoadMusicStream("./Conways Game of Life.mp3")

	pause := true

	//rl.PlayMusicStream(music_main)

	for !rl.WindowShouldClose() {
		constructIndex = clamp(constructIndex, 0, len(constructs)-1)
		tickrate = clamp(tickrate, 0, 60)
		tick++

		rl.UpdateMusicStream(music_main)

		if rl.IsKeyPressed(rl.KeyP) {
			constructorMode = !constructorMode
		}

		if rl.IsKeyPressed(rl.KeyM) {
			pause = !pause
			if pause {
				rl.PauseMusicStream(music_main)
			} else {
				rl.PlayMusicStream(music_main)
			}
		}

		if rl.IsKeyPressed(rl.KeyRight) {
			if holdingShift {
				oldCtype := constructs[constructIndex].ctype
				for constructs[constructIndex].ctype == oldCtype {
					if constructIndex+1 > len(constructs)-1 {
						break
					}
					constructIndex++
				}

			} else {
				constructIndex = clamp(constructIndex+1, 0, len(constructs)-1)
			}

		} else if rl.IsKeyPressed(rl.KeyLeft) {
			if holdingShift {
				oldCtype := constructs[constructIndex].ctype
				for constructs[constructIndex].ctype == oldCtype {
					if constructIndex-1 < 0 {
						break
					}
					constructIndex--
				}
			} else {
				constructIndex = clamp(constructIndex-1, 0, len(constructs)-1)
			}
		}

		if rl.IsKeyPressed(rl.KeySpace) && updatepixels {
			updatepixels = false
		} else if rl.IsKeyPressed(rl.KeySpace) && !updatepixels {
			updatepixels = true
		}

		if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
			holdingShift = true
		} else {
			holdingShift = false
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
			if constructorMode {
				tempboard, failed := constructs[constructIndex].build(board)
				if failed {
					cBuildFailedMessageTick.A = 150
				} else {
					board = tempboard
				}
			} else {
				mx, my := getMousePixelPos()
				if !board[my][mx] {
					board[my][mx] = true
				}

			}
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
			mx, my := getMousePixelPos()
			if board[my][mx] {
				board[my][mx] = false
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

		if constructorMode {
			rl.DrawText(fmt.Sprintf("current construct is '%s'('%s')", constructs[constructIndex].name, constructs[constructIndex].ctype), 2, 50, 15, rl.Gray)
		}

		rl.DrawText(fmt.Sprint(tickrate), 2, 25, 20, rl.Gray)

		rl.DrawText("current state saved", 2, windowY-30, 20, savedMessageTick)

		rl.DrawText("loaded saved state", 2, windowY-60, 20, LoadedMessageTick)

		rl.DrawText("unable to load save", 2, windowY-90, 20, cannotLoadMessageTick)

		rl.DrawText("unable to build construct", 2, windowY-120, 20, cBuildFailedMessageTick)

		if savedMessageTick.A > 0 {
			savedMessageTick.A--
		}

		if cannotLoadMessageTick.A > 0 {
			cannotLoadMessageTick.A--
		}

		if LoadedMessageTick.A > 0 {
			LoadedMessageTick.A--
		}

		if cBuildFailedMessageTick.A > 0 {
			cBuildFailedMessageTick.A--
		}

		rl.EndDrawing()
	}

	//rl.UnloadMusicStream(music_main)

	//rl.CloseAudioDevice()

	//rl.CloseWindow()
}
