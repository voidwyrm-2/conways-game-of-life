# Conway's Game of Life, written in Go
This is made with [raylib-go](https://github.com/gen2brain/raylib-go)(because Raylib is the best graphics library)<br>
This is my first Go project<br>
Wow Go is fast to do stuff in

### Usage Instructions
So, you want to use it?<br>
Well too bad because the starting cells are written on line 35, like so
```go
var presets = []PixelSetting{
	{49, 49, true},
	{51, 49, true},
	{51, 50, true},
	{50, 50, true},
	{50, 51, true},
}
```
I really should add some outside config, but I messed up the surounding cell detection, so I need to fix that first