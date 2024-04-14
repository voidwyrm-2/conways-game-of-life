package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const bx = 100
const by = 100

const pixelscale int = 3

var updatepixels bool = false

type CellSetting struct {
	x     int
	y     int
	state bool
}

func main() {
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

	var board [bx][by]bool
	//bpoint := &board
	var presets = []CellSetting{
		{49, 49, true},
		{51, 49, true},
		{51, 50, true},
		{50, 50, true},
		{50, 51, true},
	}
	rl.InitWindow(800, 450, "Conway's Game of Life")
	defer rl.CloseWindow()

	rl.SetTargetFPS(12)

	for _, p := range presets {
		board[p.y][p.x] = p.state
	}

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeySpace) && updatepixels {
			updatepixels = false
		} else if rl.IsKeyPressed(rl.KeySpace) && !updatepixels {
			updatepixels = true
		}

		if !updatepixels {
			rl.DrawText("Currently paused, press Space to upause", 190, 200, 20, rl.Gray)
		}

		if updatepixels {
			var bylen int = len(board)
			var bxlen int = len(board[0])
			for y := range board {
				for x := range board[y] {
					var survivalQuotent int = 0
					//var indexmargin int = 0

					// check for cells horizontally and vertically
					if x+1 < bxlen {
						if board[y][x+1] {
							survivalQuotent += 1
						}
					}
					if x-1 > -1 {
						if board[y][x-1] {
							survivalQuotent += 1
						}
					}
					if y+1 < bylen {
						if board[y+1][x] {
							survivalQuotent += 1
						}
					}
					if y-1 > -1 {
						if board[y-1][x] {
							survivalQuotent += 1
						}
					}

					// check for cells diagonally
					if y+1 < bylen && x+1 < bxlen {
						if board[y+1][x+1] {
							survivalQuotent += 1
						}
					}
					if y+1 < bylen && x-1 > -1 {
						if board[y+1][x-1] {
							survivalQuotent += 1
						}
					}
					if y-1 > -1 && x+1 < bxlen {
						if board[y-1][x+1] {
							survivalQuotent += 1
						}
					}
					if y-1 > -1 && x-1 > -1 {
						if board[y-1][x-1] {
							survivalQuotent += 1
						}
					}

					if !board[y][x] {
						if survivalQuotent == 3 {
							board[y][x] = true
						}
					} else {
						if survivalQuotent < 2 || survivalQuotent > 3 {
							board[y][x] = false
						} else if survivalQuotent == 2 || survivalQuotent == 3 {
							continue
						}
					}
				}
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		//rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.White)

		for y, b1 := range board {
			for x, b2 := range b1 {
				pixcol := rl.Black
				if b2 {
					pixcol = rl.White
				}

				rl.DrawRectangle(int32(x*pixelscale), int32(y*pixelscale), int32(pixelscale), int32(pixelscale), pixcol)
			}
		}

		rl.EndDrawing()
	}
}
